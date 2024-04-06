package checkout

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/pagooffline"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// var (
// 	PrismaServiceVar prisma.Service = prisma.Resolve()
// )

type Service interface {
	// NewPago genera un pago en la base de datos y devuelve al url para acceder al checkout y pagarlo.
	NewPago(ctx context.Context, request *dtos.PagoRequest, apiKey string) (*dtos.PagoResponse, error)
	// GetPagoStatus devuelve un bool si el pago estado del pago es
	GetPagoStatus(barcode string) (status bool, erro error)
	// GetPago devuelve los datos de un pago que necesita el checkout para mostrar al usuario pagador.
	GetPago(barcode string) (*dtos.CheckoutResponse, error)
	// GetPagoResultado obtiene los parámetros del checkout y ejecuta el pago devolviendo detalles del resultado.
	GetPagoResultado(ctx context.Context, request *dtos.ResultadoRequest) (*dtos.ResultadoResponse, error)
	// CheckPrisma funcionalidad que ayuda al checkout a saber si puede contar con los servicios de Prisma.
	CheckPrisma() error
	// GetBilling devuelve el recibo del pago en un archivo pdf
	GetBilling(uuid string) (*bytes.Buffer, error)
	// GetTarjetas devuelve datos de las tarjetas
	GetTarjetas() (*[]entities.Mediopago, error)
	// GetMatenimietoSistema permite verificar si el checkout esta en mantenimiento o no
	GetMatenimietoSistema() (estado bool, fecha time.Time, erro error)

	HashOperacionTarjeta(number string, pagointento_id int64) (status bool, erro error)
	ControlTarjetaHash(number string) (status bool, erro error)

	CheckUsuario(request dtos.RequestCheckusuario) (usuarioDb entities.Usuariobloqueados, erro error)
	UpdateUsuarioBloqueo(usuarioEntity entities.Usuariobloqueados) (erro error)
}

type service struct {
	repository         Repository
	commons            commons.Commons
	payment            PaymentFactory
	prismaService      prisma.Service
	pagoOffLineService pagooffline.Service
	utilService        util.UtilService
	webhook            webhook.RemoteRepository
	store              util.Store
}

func NewService(r Repository, c commons.Commons, ps prisma.Service, polS pagooffline.Service, utilService util.UtilService, webhook webhook.RemoteRepository, s util.Store) Service {
	return &service{
		repository:         r,
		commons:            c,
		payment:            &paymentFactory{},
		prismaService:      ps,
		pagoOffLineService: polS,
		utilService:        utilService,
		webhook:            webhook,
		store:              s,
	}
}

func NewServiceWithPayment(r Repository, c commons.Commons, p PaymentFactory) Service { //, util util.UtilService
	return &service{
		repository: r,
		commons:    c,
		payment:    p,
		//utilService: util,
	}
}

func (s *service) NewPago(ctx context.Context, request *dtos.PagoRequest, apiKey string) (*dtos.PagoResponse, error) {
	// Valido campos obligatorios
	err := request.Validar()
	if err != nil {
		return nil, err
	}

	// objeto para verificar entre otras cosas validez de fechas
	validar := commons.NewAlgoritmoVerificacion()

	// se convierten las fechas a un formato correcto para la funciones de dias entre vencimientos

	primer_vencimiento := commons.ConvertirFechaYYYYMMDD(request.FirstDueDate)
	segundo_vencimiento := commons.ConvertirFechaYYYYMMDD(request.SecondDueDate)

	cant_dias, err := validar.CalcularDiasEntreFechas(primer_vencimiento, segundo_vencimiento)

	if err != nil {
		return nil, errors.New(err.Error())
	}
	if cant_dias > 99 {
		return nil, errors.New("la cantidad de dias entre fechas de vencimiento no puede ser mayor a dos dígitos")
	}

	request.ToFormatStr()
	// Valido los montos de los items con el primer total a pagar
	err = s.validarMontos(request.Items, entities.Monto(request.FirstTotal))
	if err != nil {
		return nil, err
	}

	s.repository.BeginTx()
	defer func() {
		if err != nil {
			s.repository.RollbackTx()
		} else {
			s.repository.CommitTx()
		}
	}()

	// Busco Cuenta con apiKey
	cuenta, err := s.repository.GetCuentaByApikey(apiKey)
	if err != nil {
		return nil, err
	}
	// almacento cliente id en context
	ctx = ctxWithClienteID(ctx, uint(cuenta.ClientesID))

	// Busco id de Tipo de pago
	var tipoPagoID int64

	for _, t := range *cuenta.Pagotipos {
		if strings.EqualFold(request.PaymentType, t.Pagotipo) {
			tipoPagoID = int64(t.ID)
		}
	}

	// si tipoPagoID sigue siendo 0 significa q no hay configuracion de cuentas ni pagotipos
	if tipoPagoID <= 0 {
		return nil, fmt.Errorf("en la configuración de cuentas, no hay tipo de pago correcto para %s", request.PaymentType)
	}

	// Parseo string a fechas
	fechaVencimiento, err := time.Parse("02-01-2006", request.FirstDueDate)
	if err != nil {
		return nil, fmt.Errorf("error en fecha de vencimiento %s", err.Error())
	}
	var fechaSegundoVencimiento time.Time
	if len(request.SecondDueDate) > 0 {
		fechaSegundoVencimiento, err = time.Parse("02-01-2006", request.SecondDueDate)
		if err != nil {
			return nil, fmt.Errorf("error en fecha de segundo vencimiento: %s", err.Error())
		}
	}

	// Genero UUID unico para el pago
	// label para saltar a este punto si genera un uuid repetido
reintento:
	pagoid := s.commons.NewUUID()
	//logs.Info("codigo unico " + pagoid.String())

	// Genero registro
	pago := entities.Pago{
		PagostipoID:       tipoPagoID,
		PagoestadosID:     1,
		Description:       request.Description,
		FirstDueDate:      fechaVencimiento,
		FirstTotal:        entities.Monto(request.FirstTotal),
		SecondDueDate:     fechaSegundoVencimiento,
		SecondTotal:       entities.Monto(request.SecondTotal),
		PayerName:         request.PayerName,
		PayerEmail:        request.PayerEmail,
		ExternalReference: request.ExternalReference,
		Metadata:          request.Metadata,
		Uuid:              pagoid,
		PdfUrl:            config.APP_HOST + "/checkout/bill/" + pagoid,
		Pagoitems:         request.Items,
		Expiration:        request.Expiration,
	}

	pagodb, err := s.repository.CreatePago(ctx, &pago)
	if err != nil {
		if strings.Contains(err.Error(), "uuid_UNIQUE") {
			// cuando se intenta guardar un uuid repetido, salta a la linea 119 con el label reintento
			goto reintento
		}
		return nil, fmt.Errorf("NewPago: %s", err)
	}

	if err = s.repository.CreatePagoEstadoLog(
		ctx,
		&entities.Pagoestadologs{PagosID: int64(pagodb.ID), PagoestadosID: 1},
	); err != nil {
		logs.Error(err)
	}

	// Devuelvo response adecuada
	items := make([]dtos.PagoResponseItems, 0)
	if len(request.Items) > 0 {
		for _, t := range request.Items {
			items = append(items, dtos.PagoResponseItems{
				Quantity:    int64(t.Quantity),
				Description: t.Description,
				Amount:      t.Amount.Float64(),
				Identifier:  t.Identifier,
			})
		}
	}

	response := dtos.PagoResponse{
		ID:                int64(pagodb.ID),
		Estado:            "pending",
		Description:       pagodb.Description,
		FirstDueDate:      pagodb.FirstDueDate.Format("02-01-2006"),
		FirstTotal:        pagodb.FirstTotal.Float64(),
		SecondDueDate:     pagodb.SecondDueDate.Format("02-01-2006"),
		SecondTotal:       pagodb.SecondTotal.Float64(),
		PayerName:         pagodb.PayerName,
		PayerEmail:        pagodb.PayerEmail,
		ExternalReference: pagodb.ExternalReference,
		Metadata:          pagodb.Metadata,
		Uuid:              pagodb.Uuid,
		CheckoutUrl:       config.APP_HOST + "/checkout/" + pagodb.Uuid,
		PdfUrl:            config.APP_HOST + "/checkout/" + pagodb.Uuid,
		CreatedAt:         pagodb.CreatedAt.Format("02-01-2006"),
		Items:             items,
		Expiration:        pagodb.Expiration,
	}

	return &response, nil
}

func (s *service) validarMontos(items []entities.Pagoitems, total entities.Monto) error {
	var totalItems entities.Monto
	for _, t := range items {
		totalItems += entities.Monto(int64(t.Quantity) * int64(t.Amount))
	}
	if totalItems != total {
		return fmt.Errorf("el total de los items no coincide con el total del pago")
	}
	return nil
}

func (s *service) CheckUsuario(request dtos.RequestCheckusuario) (usuarioDb entities.Usuariobloqueados, erro error) {
	usuarioEntity := entities.Usuariobloqueados{
		Email:  request.HolderEmail,
		Nombre: request.HolderName,
		Dni:    request.HolderDocNum,
	}
	usuarioDb, err := s.repository.CheckUsuarioTRepository(usuarioEntity)
	if err != nil {
		return usuarioDb, err
	}
	//Si el usuario no existe, significa que no esta bloqueado
	if usuarioDb.ID == 0 {
		return usuarioEntity, nil
	}
	//Cada usuario puede tener un bloqueo permanente, que se realiza luego de llegar a los 3 bloqueos
	if usuarioDb.Permanente {
		logs.Info(fmt.Sprintf("Intento de pago de usuario bloqueado %s, pero con bloqueo permanente", request.HolderEmail))
		return usuarioDb, fmt.Errorf("tarjeta invalida")
	}
	//Si el usuario existe, verifico si esta bloqueado hace mas de 4 hs, para desbloqeuarlo
	now := time.Now()
	timeDiference := now.Sub(usuarioDb.FechaBloqueo)
	if timeDiference.Hours() > 4 {
		logs.Info(fmt.Sprintf("Intento de pago de usuario bloqueado %s, pero con tiempo de bloqueo mayor a 4hs, se debloquea", request.HolderEmail))
		return usuarioDb, nil
	}
	//Si el usuario existe, verifico si esta bloqueado hace menos de 4 hs, para no permitir el pago
	if timeDiference.Hours() <= 4 {
		logs.Info(fmt.Sprintf("Intento de pago de usuario bloqueado %s", request.HolderEmail))
		return usuarioDb, fmt.Errorf("tarjeta invalida")
	}
	return
}
func (s *service) UpdateUsuarioBloqueo(usuarioEntity entities.Usuariobloqueados) (erro error) {
	usuarioEntity.FechaBloqueo = time.Now()
	usuarioEntity.CantBloqueo++
	//Si el usuario no existe, significa que no esta bloqueado y lo bloqueamos
	if usuarioEntity.ID == 0 {
		erro = s.repository.AgregarUsuarioListaBloqueoRepository(usuarioEntity)
		if erro != nil {
			return erro
		}
		return nil
	}
	//Si fue bloqueado mas de 2 veces, hay que bloquearlo permanentemente
	if usuarioEntity.CantBloqueo > 2 {
		usuarioEntity.Permanente = true
	}
	erro = s.repository.UpdateUsuarioBloqueoRepository(usuarioEntity)
	if erro != nil {
		return erro
	}
	return nil
}

func ctxWithClienteID(ctx context.Context, id uint) context.Context {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)
	audit.CuentaID = id
	newCtx := context.WithValue(ctx, entities.AuditUserKey{}, audit)
	return newCtx
}

func (s *service) GetPagoStatus(uuid string) (status bool, erro error) {
	// valida uuid
	if len(uuid) <= 0 {
		return false, fmt.Errorf("debe enviar código único del pago, envió: %s", uuid)
	}

	if ok, err := s.commons.IsValidUUID(uuid); !ok {
		return false, fmt.Errorf("el identificador del pago no es válido: %w", err)
	}

	// busca pago por uuid en la base de datos
	pago, err := s.repository.GetPagoByUuid(uuid)
	if err != nil {
		return false, fmt.Errorf("error al obtener pago: %s", err.Error())
	}

	// si ya pagó (pagoestado_id != pending) devuelve un error
	if pago.PagoestadosID != 1 {
		filtroMedioPago := make(map[string]interface{})
		if len(pago.PagoIntentos) == 0 {
			return false, fmt.Errorf("el pago seleccionado no posee ningún intento de pago")
		}

		filtroMedioPago["id"] = pago.PagoIntentos[len(pago.PagoIntentos)-1].MediopagosID
		medioPago, erro := s.repository.GetMediopago(filtroMedioPago)
		if erro != nil {
			return false, fmt.Errorf("error al obtener medio de pago: %s", erro.Error())
		}
		return false, fmt.Errorf("el pago ya fue procesado a través del medio de pago %v", medioPago.Mediopago)
	}

	return true, nil
}

func (s *service) GetPago(uuid string) (*dtos.CheckoutResponse, error) {

	var secondDueDate bool
	// valida uuid
	if len(uuid) <= 0 {
		return nil, fmt.Errorf("debe enviar código único del pago, envió: %s", uuid)
	}

	if ok, err := s.commons.IsValidUUID(uuid); !ok {
		return nil, fmt.Errorf("el identificador del pago no es válido: %w", err)
	}
	// busca pago por uuid en la base de datos
	pago, err := s.repository.GetPagoByUuid(uuid)
	if err != nil {
		return nil, fmt.Errorf("error al obtener pago: %s", err.Error())
	}

	/*
		FIXME
		tal vez se tenga que controlar si el tipo de pago es offline y pagos estado es 2 no le debe
		permitir generar otro pago
	*/
	// si ya pagó devuelve un error
	if pago.PagoestadosID != 1 {
		filtroMedioPago := make(map[string]interface{})
		filtroMedioPago["id"] = pago.PagoIntentos[len(pago.PagoIntentos)-1].MediopagosID
		medioPago, erro := s.repository.GetMediopago(filtroMedioPago)
		if erro != nil {
			return nil, fmt.Errorf("error al obtener medio de pago: %s", erro.Error())
		}
		return nil, fmt.Errorf("el pago ya fue procesado a través del medio de pago %v", medioPago.Mediopago)
	}

	// TIEMPO DE EXPIRACION DEL CHECKOUT. COMPARA CON VALOR EXPIRATION DEL PAGO
	err = _getPaymentExpiration(pago)
	if err != nil {
		return nil, err
	}

	//TODO: OBTENER LOS TIPOS DE PAGOS
	channels, err := s.repository.GetPagotipoChannelByPagotipoId(pago.PagostipoID)
	if err != nil {
		return nil, err
	}
	//TODO: OBTENER LAS CUOTAS PARA LOS PAGOS CON TARJETAS
	cuotas, err := s.repository.GetPagotipoIntallmentByPagotipoId(pago.PagostipoID)
	if err != nil {
		return nil, err
	}
	// busca tipo de pago por su id
	tipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	preference, err := s.repository.GetPreferencesByIdClienteRepository(uint(tipo.Cuenta.ClientesID))
	if err != nil {
		return nil, err
	}

	var logo []byte
	if len(preference.Logo) > 0 {
		ctx := context.Background()
		logo, err = s.store.GetObjectS3(ctx, preference.Logo)
		if err != nil {
			preference.Logo = ""
			logs.Info("no se pudo recuperar logo de S3")
		}
		preference.Logo = base64.StdEncoding.EncodeToString(logo)
	}
	var channelsArrayString []string
	for _, c := range *channels {
		channelsArrayString = append(channelsArrayString, c.Channel.Channel)
	}
	var cuotasString string
	for key, value := range *cuotas {
		if len(*cuotas) == key+1 {
			cuotasString += value.Cuota
		} else {

			cuotasString += value.Cuota + ","
		}
	}
	// primer importe y fecha de vencimiento
	importe := pago.FirstTotal
	dueDate := pago.FirstDueDate
	// si se indica fecha de vencimiento, la comparamos para cobrar el segundo monto
	if !pago.FirstDueDate.IsZero() {
		hoy := time.Now().Local()
		hoyDate, err := time.Parse("2006-01-02T00:00:00Z", hoy.Format("2006-01-02T00:00:00Z"))
		if err != nil {
			return nil, fmt.Errorf("formato de fecha invalido")
		}
		if hoyDate.After(pago.FirstDueDate) {
			importe = pago.SecondTotal
			dueDate = pago.SecondDueDate
			secondDueDate = true
		}
	}
	// convierto los items del pago para mostrarlos en el frontend
	items := make([]dtos.PagoResponseItems, 0)
	if len(pago.Pagoitems) > 0 {
		for _, t := range pago.Pagoitems {
			items = append(items, dtos.PagoResponseItems{
				Quantity:    int64(t.Quantity),
				Description: t.Description,
				Amount:      t.Amount.Float64(),
			})
		}
	}
	// devuelvo los items en formato json
	byteItems, _ := json.Marshal(items)
	// armo la respuesta
	response := dtos.CheckoutResponse{
		Estado:               "pending",
		Description:          pago.Description,
		DueDate:              dueDate.Format("02-01-2006"),
		SecondDueDate:        secondDueDate,
		Total:                importe.Float64(),
		PayerName:            pago.PayerName,
		PayerEmail:           pago.PayerEmail,
		ExternalReference:    pago.ExternalReference,
		Metadata:             pago.Metadata,
		Uuid:                 pago.Uuid,
		PdfUrl:               pago.PdfUrl,
		CreatedAt:            pago.CreatedAt.String(),
		BackUrlSuccess:       tipo.BackUrlSuccess,
		BackUrlPending:       tipo.BackUrlPending,
		BackUrlRejected:      tipo.BackUrlRejected,
		IncludedChannels:     channelsArrayString, //strings.Split(tipo.IncludedChannels, ","),
		IncludedInstallments: cuotasString,        //tipo.IncludedInstallments,
		Items:                string(byteItems),
		Preference: dtos.ResponsePreference{
			Client:         preference.Cliente.Cliente,
			MainColor:      preference.Maincolor,
			SecondaryColor: preference.Secondarycolor,
			Logo:           preference.Logo,
		},
	}

	return &response, nil
}

func (s *service) GetPagoResultado(ctx context.Context, request *dtos.ResultadoRequest) (*dtos.ResultadoResponse, error) {

	// validaciones basicas
	if err := request.Validar(); err != nil {
		return nil, err
	}
	request.ToFormatStr()

	if request.CardNumber != "" {
		control, err := s.ControlTarjetaHash(request.CardNumber)
		if err != nil {
			return nil, err
		}

		if control {
			return nil, fmt.Errorf("Error con Tarjeta. Comunicarse a través de nuestro correo de soporte si persiste.")
		}

	}

	// valido el metodo de pago con el de la base de datos
	channel, err := s.repository.GetChannelByName(request.Channel)
	if err != nil {
		return nil, err
	}

	// valido el medio de pago
	filtro := make(map[string]interface{})
	filtro["channels_id"] = channel.ID
	if len(request.CardBrand) > 0 {
		filtro["mediopago"] = request.CardBrand
	} else {
		return nil, fmt.Errorf("error tarjeta invalida")
	}
	/* REVIEW: revisar cuando trae el medio de pago */
	medio, err := s.repository.GetMediopago(filtro)
	if err != nil {
		logs.Error(err)
		erro1 := fmt.Errorf("¡Lo sentimos! Por el momento este medio de pago no se encuentra disponible para esta operación")
		return nil, fmt.Errorf("error en medio de pago: %s", erro1.Error())
	}

	// al request le paso el id externo de medio de pago
	request.PaymentMethodID, _ = strconv.ParseInt(medio.ExternalID, 10, 64)

	// obtengo datos del pago mediante el uuid
	pago, err := s.repository.GetPagoByUuid(request.Uuid)
	if err != nil {
		return nil, err
	}
	//Si se paga con tarjeta, tengo que realizar las siguientes verificaciones
	var usuarioDb entities.Usuariobloqueados
	if channel.ID == 1 || channel.ID == 2 {
		// if pago.PagostipoID == 10 {
		// 	logs.Info("pago usuario dpec: %s" + pago.PayerEmail)
		// 	err = errors.New("servicio no disponible,  inténtalo más tarde")
		// 	return nil, err
		// }

		//1- Verifico si el usuario esta bloqueado
		requsuario := dtos.RequestCheckusuario{
			HolderName:  pago.PayerName,
			HolderEmail: pago.PayerEmail,
		}
		usuarioDb, err = s.CheckUsuario(requsuario)
		if err != nil {
			logs.Info("pago usuario dpec intento de fraude:" + pago.PayerEmail)
			errorUsuarioBloqueado := errors.New("tarjeta invalida")
			return nil, errorUsuarioBloqueado
		}
	}

	// ver tiempo de expiracion para realizar el pago en el checkout
	err = _getPaymentExpiration(pago)
	if err != nil {
		return nil, err
	}

	// obtengo datos del tipo de pago configurado por el cliente
	tipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	// obtengo datos de la cuenta bancaria a la cual corresponde el pago
	cuenta, err := s.repository.GetCuentaById(int64(tipo.CuentasID))
	if err != nil {
		return nil, err
	}

	// calculo el importe a pagar segun las fechas de vencimiento
	importe := pago.FirstTotal

	// si se indica fecha de vencimiento, la comparamos para cobrar el segundo monto
	if !pago.FirstDueDate.IsZero() {
		fechaHoy := time.Now().Local()
		hoyDate, err := time.Parse("2006-01-02T00:00:00Z", fechaHoy.Format("2006-01-02T00:00:00Z"))
		if err != nil {
			return nil, fmt.Errorf("formato de fecha invalido")
		}
		logs.Info("fecha actual")
		logs.Info(hoyDate)
		logs.Info("fecha primer venc.")
		logs.Info(pago.FirstDueDate)
		logs.Info(hoyDate.After(pago.FirstDueDate))
		if hoyDate.After(pago.FirstDueDate) {
			importe = pago.SecondTotal
		}
	}
	// el monto lo pasamos como integer a las apis
	request.Importe = importe.Int64()
	fechaActual, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		erro := errors.New("error convertir fecha actual")
		return nil, erro
	}
	cuotas, _ := strconv.ParseInt(request.Installments, 10, 64)
	installments, err := s.repository.GetInstallmentsByMedioPagoInstallmentsId(medio.MediopagoinstallmentsID)
	if err != nil {
		return nil, err
	}
	var installmentId int64
	for _, valueInstallment := range installments {
		if valueInstallment.VigenciaHasta == nil {
			installmentId = int64(valueInstallment.ID)
			break
		}
		if (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.Equal(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Equal(*valueInstallment.VigenciaHasta)) {
			installmentId = int64(valueInstallment.ID)
			break
		}
	}

	installmentsDetails, err := s.repository.GetInstallmentDetails(installmentId, cuotas) //medio.InstallmentsID int64(installments.ID)
	if err != nil {
		return nil, err
	}
	filtroConfiguracion := filtros.ConfiguracionFiltro{
		Buscar:     true,
		Nombrelike: "IMPUESTO_SOBRE_COEFICIENTE",
	}
	configuracionImpuesto, err := s.utilService.GetConfiguracionesService(filtroConfiguracion)
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	impuestoId, err := strconv.Atoi(configuracionImpuesto[0].Valor)
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	impuesto, err := s.utilService.GetImpuestoByIdService(int64(impuestoId))
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	installmentsDetails.Impuesto = impuesto.Porcentaje
	// obtengo constructor de pago mediante patrón factory
	paymentMethod, err := s.payment.GetPaymentMethod(int(channel.ID))
	if err != nil {
		return nil, err
	}

	// genero un transaction_id
	transactionID := s.commons.NewUUID()

	// envio los datos al constructor correspondiente para procesar el pago
	resultado, err := paymentMethod.CreateResultado(request, pago, cuenta, transactionID, installmentsDetails)
	if err != nil {
		return nil, err
	}
	//Si el pago es con tarjeta y el estado es distinto de aprobado, se incrementa el contador de intentos fallidos
	if channel.ID == 1 || channel.ID == 2 {
		if resultado.StateComment != "approved" {
			pago.IntentoFallido++
			err := s.repository.UpdatePagoIntentoFallidos(ctx, pago)
			if err != nil {
				logs.Error(err)
			}
			//Si el pago tiene 3 o mas intentos fallidos, bloquear al usuario
			if pago.IntentoFallido > 2 {
				logs.Error(fmt.Sprintf("el usuario %s intento pagar mas de %d veces de forma erronea el pago id:%d", pago.PayerEmail, pago.IntentoFallido, pago.ID))
				err = s.UpdateUsuarioBloqueo(usuarioDb)
				if err != nil {
					logs.Error(err)
				}
				err = fmt.Errorf("tarjeta invalida")
				logs.Info("pago usuario dpec intento de fraude:" + pago.PayerEmail)
				return nil, err
			}
		}
	}

	// busco en la base el id de installmentdetails
	resultado.InstallmentdetailsID = int64(installmentsDetails.Id)
	//resultado.InstallmentdetailsID = s.repository.GetInstallmentDetailsID(medio.InstallmentsID, cuotas)
	resultado.MediopagosID = int64(medio.ID)

	// algunas apis devuelven el monto en entero, me aseguro q se guarde en bd el float64
	resultado.Amount = importe
	//Agrego la ip
	resultado.Ip = request.Ip
	// agrego cliente id 1 al context para la auditoria
	ctx = ctxWithClienteID(ctx, 1)

	// almaceno el resultado en la base de datos
	if ok, err := s.repository.CreateResultado(ctx, resultado); !ok {
		return nil, err
	}

	if resultado.CardLastFourDigits != "" {
		if _, err := s.HashOperacionTarjeta(request.CardNumber, int64(resultado.ID)); err != nil {
			return nil, err
		}
	}

	// actualizo el estado del pago
	// cuando el pago se procesa con exito, se le coloca una fecha a PaidAt,
	// cuando hay un error, se devuelve PaidAt con fecha 0, y no se actualiza el pago.

	if !resultado.PaidAt.IsZero() {
		if channel.Channel == "DEBIN" {
			pago.PagoestadosID = 2
		} else if channel.Channel == "OFFLINE" {
			pago.PagoestadosID = 2
		} else {
			pago.PagoestadosID = 4
		}
		if ok, err := s.repository.UpdatePago(ctx, pago); !ok {
			return nil, err
		}
		if err = s.repository.CreatePagoEstadoLog(
			ctx,
			&entities.Pagoestadologs{PagosID: int64(pago.ID), PagoestadosID: pago.PagoestadosID},
		); err != nil {
			logs.Error(err)
		}
	}
	var importePagado entities.Monto
	if resultado.Valorcupon > 0 {
		importePagado = entities.Monto(resultado.Valorcupon)
	} else {
		importePagado = entities.Monto(resultado.Amount)

	}

	estadoPago, erro := s.repository.GetPagoEstado(pago.PagoestadosID)
	if erro != nil {
		estadoPago.Nombre = "PENDIENTE"
		logs.Error(erro)
	}

	// armo respuesta para el frontend
	response := dtos.ResultadoResponse{
		ID:                int64(resultado.ID),
		Estado:            resultado.StateComment,
		EstadoPago:        estadoPago.Nombre,
		Exito:             !resultado.PaidAt.IsZero(),
		Uuid:              request.Uuid,
		Channel:           channel.Channel,
		Description:       pago.Description,
		FirstDueDate:      pago.FirstDueDate.Format("02-01-2006"),
		FirstTotal:        importe.Float64(),
		SecondDueDate:     pago.SecondDueDate.Format("02-01-2006"),
		SecondTotal:       pago.SecondTotal.Float64(),
		PayerName:         pago.PayerName,
		PayerEmail:        pago.PayerEmail,
		ExternalReference: pago.ExternalReference,
		Metadata:          pago.Metadata,
		PdfUrl:            pago.PdfUrl,
		CreatedAt:         resultado.CreatedAt.String(),
		ImportePagado:     importePagado.Float64(),
	}
	/*
		Autor: Jose Alarcon
		Fecha: 21/06/2022
		Descripción: webhook notificacion de pago al cliente
		Verificar que el cliente tenga una url configurada y que ademas para notifcar el pago sea exitoso
	*/
	if tipo.BackUrlNotificacionPagos != "" && response.Exito {
		var apikey string
		if len(tipo.ApikeyExterno) > 0 {
			apikey = tipo.ApikeyExterno
		}
		var result []dtos.ResultadoResponse
		result = append(result, response)
		notificacionPago := dtos.ResultadoResponseWebHook{
			Url:               tipo.BackUrlNotificacionPagos,
			ApikeyExterno:     apikey,
			ResultadoResponse: result,
		}
		er := s.webhook.NotificarPago(notificacionPago)
		if er != nil {
			logs.Info("webhook:no se pudo notificar el pago: %s" + er.Error())
		} else {
			peticionWebHook := dtos.RequestWebServicePeticion{
				Operacion: "NotificarPago",
				Vendor:    "WebHook",
			}
			err1 := s.utilService.CrearPeticionesService(peticionWebHook)
			if err1 != nil {
				logs.Error("no se pudo registrar la peticion" + err1.Error())
			}
		}
	}

	if response.Exito {
		// Construir el texto html del mensaje del email
		dir_url_comprobante := config.APP_HOST + "/checkout/bill/" + pago.Uuid
		url_imagen_descargaDoc := "https://img.icons8.com/?size=512&id=2mGSkp2owx0d&format=png"
		mensaje := "<ul style='list-style: none;text-align: left;display:inline-block;'><li> Fecha: <b>#4</b></li><li> Referencia: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li><li>Nro solicitud: <b>#5</b></li> <li style='padding-top:6px;' > <a href='" + dir_url_comprobante + "'><img src='" + url_imagen_descargaDoc + "' width='16' height='16'> Comprobante de Pago </a> </li></ul>"

		// TODO se modifica template de recibo de pagos
		var descripcion utildtos.DescripcionTemplate
		var detallesPago []utildtos.DetallesPago

		for _, det := range pago.Pagoitems {
			var identificador string
			if len(det.Identifier) > 0 {
				identificador = " - " + det.Identifier
			}
			detallesPago = append(detallesPago, utildtos.DetallesPago{
				Descripcion: det.Description + identificador,
				Cantidad:    fmt.Sprintf("%v", det.Quantity),
				Monto:       fmt.Sprintf("$%v", s.utilService.FormatNum(s.utilService.ToFixed(det.Amount.Float64(), 2))),
			})
		}

		descripcion = utildtos.DescripcionTemplate{
			Cliente:     cuenta.Cliente.Cliente,
			Cuit:        cuenta.Cliente.Cuit,
			Detalles:    detallesPago,
			TotalPagado: fmt.Sprintf("$%v", s.utilService.FormatNum(s.utilService.ToFixed(response.ImportePagado, 2))),
		}

		// decidir que titulo de mensaje enviar en el email
		mensajeSegunMedioPagoEmail := GetMensajeSegunMedioPago("email", medio.ChannelsID)

		/* enviar mail al usuario pagador */
		var arrayEmail []string
		var email string
		email = request.HolderEmail
		if request.HolderEmail == "" {
			email = pago.PayerEmail
		}
		arrayEmail = append(arrayEmail, email)
		params := utildtos.RequestDatosMail{
			Email:                 arrayEmail,
			Asunto:                "Información de Pago",
			Nombre:                pago.PayerName,
			Mensaje:               mensaje,
			CamposReemplazar:      []string{response.ExternalReference, pago.Uuid, medio.Mediopago, response.Description, resultado.PaidAt.Format("02-01-2006"), fmt.Sprintf("%v", pago.ID), dir_url_comprobante},
			FiltroReciboPago:      true,
			Descripcion:           descripcion,
			From:                  "Wee.ar!",
			TipoEmail:             "template",
			MensajeSegunMedioPago: mensajeSegunMedioPagoEmail,
		}
		// aqui utiliza el util service para enviar los parametros que se van a mostrar en el email y enviar el email
		erro = s.utilService.EnviarMailService(params)
		if erro != nil {
			logs.Error(erro.Error())
		}
	}

	mensajeSegunMedioPagoCheckout := GetMensajeSegunMedioPago("checkout", medio.ChannelsID)
	response.Mensaje = mensajeSegunMedioPagoCheckout.Content

	return &response, nil
}

func (s *service) CheckPrisma() error {
	check, err := s.prismaService.CheckService()
	if err != nil {
		return err
	}
	if !check {
		return fmt.Errorf("el servicio de prisma no está disponible")
	}
	return nil
}

func (s *service) GetBilling(uuid string) (*bytes.Buffer, error) {
	pago, err := s.repository.GetPagoByUuid(uuid)
	if err != nil {
		return nil, err
	}

	intento, err := s.repository.GetValidPagointentoByPagoId(int64(pago.ID))
	if err != nil {
		return nil, err
	}

	medioPago, err := s.repository.GetMediopago(map[string]interface{}{"id": intento.MediopagosID})
	if err != nil {
		return nil, err
	}

	channel, err := s.repository.GetChannelById(uint(medioPago.ChannelsID))
	if err != nil {
		return nil, err
	}

	pagotipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}
	logs.Info(pagotipo)
	cuenta, err := s.repository.GetCuentaById(int64(pagotipo.CuentasID))
	if err != nil {
		return nil, err
	}
	logs.Info(cuenta)

	cliente, err := s.repository.GetClienteByApikey(cuenta.Apikey)
	if err != nil {
		return nil, err
	}
	logs.Info(cliente)

	// generar el pdf del comprobante de pago con los datos
	file, err := _getBillingPdf(pago, cliente, channel, intento)

	if err != nil {
		fmt.Println("Could not save PDF:", err)
	}

	return &file, nil
}

func (s *service) HashOperacionTarjeta(number string, pagointento_id int64) (status bool, erro error) {
	//Hashear card number
	textoPlano := number
	hash := md5.Sum([]byte(textoPlano))
	hashstring := hex.EncodeToString(hash[:])

	hasheado := entities.Uuid{
		Uuid: hashstring,
	}

	s.repository.SaveHasheado(&hasheado, uint(pagointento_id))

	return
}

func (s *service) ControlTarjetaHash(number string) (status bool, erro error) {
	//Hashear card number
	textoPlano := number
	hash := md5.Sum([]byte(textoPlano))
	hashstring := hex.EncodeToString(hash[:])

	status, err := s.repository.GetHasheado(hashstring)

	if err != nil {
		erro = err
		return
	}

	return
}

func (s *service) GetMatenimietoSistema() (estado bool, fecha time.Time, erro error) {
	estado, fecha, err := s.utilService.GetMatenimietoSistemaService()
	if err != nil {
		erro = fmt.Errorf("el servicio no está disponible")
		return
	}
	// filtro := filtros.ConfiguracionFiltro{
	// 	Nombre: "ESTADO_APLICACION",
	// }
	// estadoConfiguracion, err := s.utilService.GetConfiguracionService(filtro)
	// if err != nil {
	// 	estado = true
	// 	erro = fmt.Errorf("el servicio no está disponible")
	// 	return
	// }
	// if estadoConfiguracion.Valor != "sin valor" {
	// 	fecha, err = time.Parse(time.RFC3339, estadoConfiguracion.Valor)
	// 	if err != nil {
	// 		estado = true
	// 		logs.Error("error al convertir fecha de configuración")
	// 		erro = fmt.Errorf("el servicio no está disponible")
	// 		return
	// 	}
	// 	if !fecha.IsZero() {
	// 		estado = true
	// 		return
	// 	}
	// }
	// estado = false
	return
}

func (s *service) GetTarjetas() (*[]entities.Mediopago, error) {
	return s.repository.GetMediosDePagos()
}

/* ---------------------------------- Funciones auxiliares ---------------------------------------------- */

func _getPaymentExpiration(pago *entities.Pago) (erro error) {
	hoy := time.Now().Local()
	tiempoExpiracion := pago.Expiration
	diferencia := hoy.Sub(pago.CreatedAt)
	minutos := diferencia.Minutes() // en float64
	if tiempoExpiracion != 0 && minutos > float64(tiempoExpiracion) {
		return fmt.Errorf("el pago expiró, vuelva a generarlo")
	}
	return
}

func formatFechaString(fecha time.Time, formatoFecha string) string {
	fechaStr := fecha.Format(formatoFecha)
	fechaArrayStr := strings.Split(fechaStr[0:10], "-")
	fechaVto := fmt.Sprintf("%v-%v-%v", fechaArrayStr[2], fechaArrayStr[1], fechaArrayStr[0])
	return fechaVto
}

func GetMensajeSegunMedioPago(proposito string, channel_id int64) (message utildtos.MensajeSegunMedioPagoStruct) {
	var Credit, Debit, Debin int64
	Credit = 1
	Debit = 2
	Debin = 4

	if channel_id == Credit || channel_id == Debit {
		message.Title = "Pago Aprobado"
		if proposito == "email" {
			message.Content = "Su pago fue procesado de manera exitosa."
		}
	} else if channel_id == Debin {
		message.Title = "Pago Procesado"
		message.Content = "Diríjase a su Home Banking para autorizar el debin generado."
	} else {
		message.Title = "Pago Procesado"
		message.Content = "Diríjase a un punto de Rapipago para completar el proceso de pago."
	}
	return
}

/* -------------------------------- Funciones PDF Comprobante ------------------------------------------- */

func getHeaderAndContent(pagoItems *[]entities.Pagoitems) (header []string, contents [][]string) {
	// La cabecera de la tabla. Los nombres de las columnas
	header = []string{"Transacción", "Producto", "Cantidad", "Precio"}

	// cantidad de items del pago
	size := len(*pagoItems)
	items := make([][]string, size)

	for i, x := range *pagoItems {
		identificador := ""
		if x.Identifier != "" {
			identificador = x.Identifier
		}
		// identificador, descripcion, cantidad, monto
		items[i] = []string{
			identificador, x.Description, fmt.Sprint(x.Quantity), strconv.FormatFloat(x.Amount.Float64(), 'f', 2, 64),
		}
	}

	contents = items

	return header, contents
}

func getTelCoGreenColor() color.Color {
	return color.Color{
		Red:   195,
		Green: 216,
		Blue:  46,
	}
}

func getHeaderTextColor() color.Color {
	return color.NewBlack()
}

func getTelCoSoftBlueColor() color.Color {
	return color.Color{
		Red:   0,
		Green: 184,
		Blue:  241,
	}
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

// size dinamico segun long del texto que se recibe
func _resolveColumnWidthSize(texto string) (colSize, colSpaceSize int) {
	long := len(texto)
	var columMaxSize int = 12
	switch {
	case long <= 30:
		colSize = 3
		colSpaceSize = columMaxSize - colSize
	case long > 30 && long < 43:
		colSize = 4
		colSpaceSize = columMaxSize - colSize
	case long >= 43:
		colSize = 5
		colSpaceSize = columMaxSize - colSize
	default:
		colSize = 4
		colSpaceSize = columMaxSize - colSize
	}
	return
}

func _buildHeading(m pdf.Maroto, cliente *entities.Cliente, intento *entities.Pagointento) {
	green := getTelCoGreenColor()
	negro := getHeaderTextColor()
	blanco := color.NewWhite()

	// RegisterHeader
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				// LOCAL
				// err := m.FileImage(filepath.FromSlash("./assets/images/cabecera_recibo.png"), props.Rect{})

				// SERVER
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "cabecera_recibo.png"), props.Rect{})

				if err != nil {
					logs.Error("_buildHeading: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
		})
	})

	texto := cliente.Cliente + " - CUIT " + cliente.Cuit

	// segun long de texto se adapta columnas
	colSize, colSpaceSize := _resolveColumnWidthSize(texto)

	m.Row(10, func() {
		m.Col(uint(colSpaceSize), func() {
			m.Text(intento.HolderName, props.Text{Size: 8, Top: 5, Left: 10})
		})
		m.SetBackgroundColor(green)
		m.Col(uint(colSize), func() {
			m.Text(texto, props.Text{
				Style: consts.Normal,
				Color: negro,
				Align: consts.Right,
				Size:  8,
				Top:   3,
				Right: 2,
			})
		})
		m.SetBackgroundColor(blanco)
	})
	m.Row(10, func() {
		m.Col(uint(colSpaceSize), func() {
			m.Text("CUIL/DNI: "+intento.HolderNumber, props.Text{Size: 8, Left: 10})
		})
		m.Col(uint(colSize), func() {
			fechaEmision := time.Now().Format("02-01-2006")

			m.Text("Fecha de emisión: "+fechaEmision, props.Text{
				Style: consts.Normal,
				Color: negro,
				Align: consts.Right,
				Size:  8,
				Top:   3,
				Right: 2,
			})
		})
		m.SetBackgroundColor(blanco)
	})
}

func _buildBodyList(m pdf.Maroto, pago *entities.Pago, channel entities.Channel, intento *entities.Pagointento) {
	// set de colores
	celeste := getTelCoSoftBlueColor()
	blanco := color.NewWhite()
	darkGrayColor := getDarkGrayColor()

	// contenido y cabeceras de la tabla
	header, contents := getHeaderAndContent(&pago.Pagoitems)

	// mostrar solo cuando el channel es debito o credito
	if channel.Channel == "CREDIT" || channel.Channel == "DEBIT" {
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Verás este pago en tu resumen como TelCo Wee!!", props.Text{
					Top:   5,
					Style: consts.Italic,
					Align: consts.Center,
					Color: color.NewBlack(),
				})
			})
		})
	}

	// Referencia del pago
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Referencia de pago: "+pago.ExternalReference, props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})
	// una linea de espacio o separacion
	m.Line(5, props.Line{
		Color: blanco,
	})

	m.SetBackgroundColor(celeste)

	// Tabla de items del pago
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			GridSizes: []uint{3, 4, 2, 3},
			Size:      10,
			Style:     consts.Bold,
			Color:     blanco,
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		HeaderContentSpace:   1,
		AlternatedBackground: &blanco,
		Line:                 true,
		LineProp: props.Line{
			Color: celeste,
		},
		VerticalContentPadding: 7, // alto de fila
	})

	// Totales y vencimientos, segun canales de pago
	if channel.Channel == "OFFLINE" {
		//primer vencimiento
		m.Row(12, func() {

			m.Col(1, func() {
				m.Text("Primer Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.FirstDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0
				if intento.Valorcupon > 0 {
					total = intento.Valorcupon.Float64()
				} else {
					total = intento.Amount.Float64()
				}
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})

		//Segundo vencimiento
		m.Row(12, func() {
			// m.ColSpace(7) // Incluye un espacio en blanco de 7 columnas a la izquierda
			// m.ColSpace(1) // Incluye un espacio en blanco de 7 columnas a la izquierda
			m.Col(1, func() {
				m.Text("Segundo Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.SecondDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0

				total = pago.SecondTotal.Float64()

				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})
		// se genera el codigo de barra
		m.Row(15, func() { //15
			m.Col(6, func() {

				_ = m.Barcode(intento.Barcode, props.Barcode{
					Percent: 0,
					Proportion: props.Proportion{
						Width:  20,
						Height: 2,
					},
				})
				m.Text(intento.Barcode, props.Text{
					Top:    12,
					Family: "",
					Style:  consts.Bold,
					Size:   9,
					Align:  consts.Center,
				})
			})
			m.ColSpace(6)
		})
	}

	if channel.Channel == "CREDIT" {
		// Si existe valor cupon, y por lo tanto es un pago en cuotas, se muestra el costo financiado
		if intento.Valorcupon != 0 {
			dif := intento.Valorcupon - intento.Amount
			costo_financiero := dif.Float64()
			m.Row(7, func() {
				m.ColSpace(7) // Incluye un espacio en blanco de 7 columnas a la izquierda
				m.Col(2, func() {
					m.Text("Costo Financiero:", props.Text{
						Top:   5,
						Style: consts.Bold,
						Size:  8,
						Align: consts.Right,
					})
				})
				m.Col(3, func() {
					m.Text(strconv.FormatFloat(costo_financiero, 'f', 2, 64), props.Text{
						Top:   5,
						Style: consts.Bold,
						Size:  8,
						Align: consts.Center,
					})
				})
			})
		}
		m.Row(7, func() {
			m.ColSpace(4)
			m.Col(4, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(5, func() {
				var total float64 = 0
				if intento.Valorcupon > 0 {
					total = intento.Valorcupon.Float64()
				} else {
					total = intento.Amount.Float64()
				}
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})
	}

	if channel.Channel == "DEBIN" || channel.Channel == "DEBIT" {
		m.Row(7, func() { //15
			m.ColSpace(4)
			m.Col(4, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(5, func() {
				var total float64 = 0
				total = intento.Amount.Float64()
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})
	}

	// Medio de Pago, Numero de Solicitud, Numero de Operacion, Codigo Autorizacion
	m.Row(4, func() {
		var cardLastFour string
		cardLastFour = ""
		if channel.Channel == "CREDIT" || channel.Channel == "DEBIT" {
			cardLastFour = ". Los últimos 4 dígitos de su tarjeta son: " + intento.CardLastFourDigits
		}
		m.Col(15, func() {
			m.Text("Medio de Pago: "+channel.Nombre+cardLastFour, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})

		})
	})
	m.Row(4, func() {
		m.Col(15, func() {
			nroSolicitud := strconv.Itoa(int(pago.ID))
			m.Text("Nro. Solicitud: "+nroSolicitud, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})
	})
	m.Row(4, func() {
		m.Col(15, func() {
			nroOperacion := strconv.Itoa(int(intento.ID))
			m.Text("Nro. Op.: "+nroOperacion, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})

	})
	m.Row(4, func() {
		var codigoAutorizacion string
		if channel.Channel == "DEBIN" || channel.Channel == "OFFLINE" {
			codigoAutorizacion = intento.ExternalID
		}
		if channel.Channel == "CREDIT" || channel.Channel == "DEBIT" {
			codigoAutorizacion = intento.AuthorizationCode
		}
		m.Col(15, func() {
			m.Text("Código Autorización: "+codigoAutorizacion, props.Text{
				Top:   10,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})
	})
}

func _buildFooter(m pdf.Maroto) {
	m.RegisterFooter(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				// LOCAL
				// err := m.FileImage(filepath.FromSlash("./assets/images/footer_recibo.png"), props.Rect{
				// 	Top: 30,
				// })

				// SERVER
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "footer_recibo.png"), props.Rect{
					Top: 30,
				})

				if err != nil {
					logs.Error("_buildFooter: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
		})
	})
}

func _getBillingPdf(pago *entities.Pago, cliente *entities.Cliente, channel entities.Channel, intento *entities.Pagointento) (file bytes.Buffer, erro error) {

	// instancia de objeto PDF
	m := pdf.NewMaroto(consts.Portrait, consts.A4)

	m.SetPageMargins(0, 0, 0)

	// Header
	_buildHeading(m, cliente, intento)

	m.SetPageMargins(10, 0, 10)

	// Body
	_buildBodyList(m, pago, channel, intento)

	m.SetPageMargins(0, 0, 0)

	// Footer
	_buildFooter(m)

	// set de variables de retorno
	file, erro = m.Output()
	return
}
