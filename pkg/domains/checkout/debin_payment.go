package checkout

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type debinPayment struct {
	service     apilink.AplinkService
	utilService util.UtilService
}

func NewDebinPayment(s apilink.AplinkService, util util.UtilService) PaymentMethod {
	return &debinPayment{
		service:     s,
		utilService: util,
	}
}

func (c *debinPayment) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transactionID string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error) {
	//serv := apilink.Resolve()
	// variable que vamos a guardar en la tabla en el campo holder_cbu
	var cbu string
	// variable que vamos a pasar como cuenta a apilink
	var cuentaLink linkdebin.CuentaLink
	// la cuenta va a pasar una cbu o un alias según corresponda
	if len(request.Alias) > 0 {
		cuentaLink.AliasCbu = request.Alias
		cbu = request.Alias
	} else {
		cuentaLink.Cbu = request.Cbu
		cbu = request.Cbu
	}

	comprador := linkdebin.CompradorCreateDebinLink{
		Cuit:   request.HolderCuit,
		Cuenta: cuentaLink,
	}

	// filtro := filtros.ConfiguracionFiltro{
	// 	Nombre: "CBU_CUENTA_TELCO",
	// }

	// cuentaTelco, erro := c.utilService.GetConfiguracionService(filtro)
	// if erro != nil {
	// 	logs.Error(erro.Error())
	// 	return nil, erro
	// }

	vendedor := linkdebin.VendedorCreateLink{
		Cbu: config.CBU_CUENTA_TELCO, //cuentaTelco.Valor,
	}

	debin := linkdebin.DebinCreateLink{
		ComprobanteId:         pago.ExternalReference,
		EsCuentaPropia:        request.EsCuentaPropia,
		Concepto:              linkdtos.EnumConceptoDebin(request.ConceptoAbreviado),
		TiempoExpiracion:      180, //request.TiempoExpiracion, // tiempo de expiracion en minutos
		Descripcion:           transactionID,
		Importe:               request.Importe,
		Moneda:                linkdtos.EnumMoneda(request.Moneda),
		Recurrente:            request.Recurrente,
		DescripcionPrestacion: pago.Description,
	}

	debinRequest := linkdebin.RequestDebinCreateLink{
		Comprador: comprador,
		Vendedor:  vendedor,
		Debin:     debin,
	}

	debinResult, err := c.service.CreateDebinApiLinkService(request.Uuid, debinRequest)

	pagoResultado := entities.Pagointento{
		PagosID:              int64(pago.ID),
		MediopagosID:         38,
		InstallmentdetailsID: 1,
		ExternalID:           debinResult.Id,
		ReportAt:             time.Now().Local(),
		IsAvailable:          false,
		Amount:               entities.Monto(request.Importe),
		//Valorcupon:           entities.Monto(request.Importe),
		StateComment:  string(debinResult.Estado),
		TransactionID: transactionID,
		HolderName:    request.HolderName,
		HolderType:    "CUIL/CUIT",
		HolderNumber:  request.HolderCuit,
		HolderCbu:     cbu,
	}

	if err != nil {
		logs.Error(err)
		pagoResultado.StateComment = err.Error()
	}
	// TODO: por ahora viene ACEPTADO en produccion viene INICIADO
	var iniciados = []interface{}{"ACEPTADO", "INICIADO"}

	if commons.Include(iniciados, string(debinResult.Estado)) {
		pagoResultado.PaidAt = debinResult.FechaOperacion
	}

	return &pagoResultado, nil
}
