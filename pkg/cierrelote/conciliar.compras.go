package cierrelote

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type conciliarCompras struct {
	utilService util.UtilService
}

func NewConciliarCompras(util util.UtilService) MetodoConciliarClMP {
	return &conciliarCompras{
		utilService: util,
	}
}

func (c *conciliarCompras) ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	if cierreLote.Tipooperacion == "C" {
		strNroEstablecimiento := strconv.Itoa(int(cierreLote.Nroestablecimiento))
		if cierreLote.Fechaoperacion == movimientoDetalle.FechaOrigenCompra && strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && cierreLote.Nrotarjeta == movimientoDetalle.NroTarjetaXl && strings.Contains(movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion) && cierreLote.Nroticket == movimientoDetalle.NroCupon && valorCuota == movimientoDetalle.PlanCuota && movimientoDetalle.TipoAplicacion == "+" && movimientoCabecera.Codop == movimientoDetalle.Tipooperacion.ExternalId { //&& cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe)
			if cierreLote.Monto.Int64() < int64(movimientoDetalle.Importe) || cierreLote.Monto.Int64() > int64(movimientoDetalle.Importe) {
				// mensaje := fmt.Sprintf(" monto ci $ #0 e importe  $ #1 - cl_id = #2 y movimientoDetalle_id = #3 ")
				// dataReemplazar := []string{fmt.Sprintf("%v", cierreLote.Monto.Int64()), fmt.Sprintf("%v", int64(movimientoDetalle.Importe)), fmt.Sprintf("%v", cierreLote.Id), fmt.Sprintf("%v", movimientoDetalle.Id)}
				// correos := []string{config.EMAIL_TELCO}
				// params := utildtos.RequestDatosMail{
				// 	Email:            correos,
				// 	Asunto:           "monto ",
				// 	From:             "Wee.ar!",
				// 	Nombre:           "Administracion Wee!!",
				// 	Mensaje:          mensaje,
				// 	CamposReemplazar: dataReemplazar,
				// 	TipoEmail:        "template",
				// 	AdjuntarEstado:   false,
				// }
				// err := c.utilService.EnviarMailService(params)
				// if err != nil {
				// notificacion := entities.Notificacione{
				// 	Tipo:        entities.EnumTipoNotificacion("ConciliacionClMx"),
				// 	Descripcion: fmt.Sprintf("diferencia entre monto CL $%v e importe Mov. detalles $%v. cl_id = %v y movDetalle_id = %v ", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe), cierreLote.Id, movimientoDetalle.Id),
				// }
				// err := c.utilService.CreateNotificacionService(notificacion)
				// if err != nil {
				logs.Error(fmt.Sprintf("error: existe diferencia entre monto cierre lote $%v e importe movimiento detalles $%v. cl_id = %v y movimientoDetalle_id = %v ", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe), cierreLote.Id, movimientoDetalle.Id))
				// }
				// }
				return

			}
			porcentajeArancelControl := c.utilService.ToFixed(cierreLote.Channelarancel.Importe*100, 2)
			porcentajeArancelPrisma := movimientoDetalle.PorcentDescArancel / 100
			// si los porcentajes no coinciden se marca como observacion
			if porcentajeArancelControl != porcentajeArancelPrisma {
				logs.Error(fmt.Sprintf("error: existe diferencia entre porcentaje arancel control %v y porcentaje arancel prisma %v. cl_id = %v y movimientoDetalle_id = %v", porcentajeArancelControl, porcentajeArancelPrisma, cierreLote.Id, movimientoDetalle.Id))
				cierreLote.Enobservacion = true
			}
			cierreLote.FechaPago = movimientoCabecera.FechaPago
			fecha := cierreLote.FechaCierre
			if !cierreLote.FechaCierre.Equal(movimientoCabecera.FechaPresentacion) {
				fecha = movimientoCabecera.FechaPresentacion
				cierreLote.Descripcionpresentacion = fmt.Sprintf("la fecha de cierre en CL %v se modifico por %v", cierreLote.FechaCierre, movimientoCabecera.FechaPresentacion)
				cierreLote.FechaCierre = fecha
			}
			cierreLote.Cantdias = int(cierreLote.FechaPago.Sub(fecha).Hours() / 24)
			cierreLote.PrismamovimientodetallesId = movimientoDetalle.Id
			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
			cabeceraMoviminetosIdArray = append(cabeceraMoviminetosIdArray, int64(movimientoDetalle.PrismamovimientototalesId))

			RequestValidarCF := utildtos.RequestValidarCF{
				Cupon:  cierreLote.Monto,
				Cuotas: float64(cierreLote.Nrocuota),
				Dias:   float64(cierreLote.Cantdias),
				Tna:    cierreLote.Istallmentsinfo.Tna,
				//channelArancel.importe en prisma es el porcentaje
				ArancelMonto: cierreLote.Channelarancel.Importe,
			}
			responseValidarCF := util.Resolve().ValidarCalculoCF(RequestValidarCF)
			valor_pres := entities.Monto(responseValidarCF.ValorPresente * 100)
			cierreLote.Valorpresentado = valor_pres
			importeArancel := util.Resolve().ToFixed(cierreLote.Monto.Float64()-responseValidarCF.ValorPresente, 4)
			cierreLote.Importeivaarancel = util.Resolve().ToFixed(importeArancel*0.21, 4)
			// logs.Info("========")
			// logs.Info(cierreLote.Monto.Float64())
			// logs.Info(responseValidarCF.ValorPresente)
			// logs.Info(cierreLote.Diferenciaimporte)
			// logs.Info(cierreLote.Importeivaarancel)
			// logs.Info("========")
			cierreLote.Coeficientecalculado = responseValidarCF.ValorCoeficiente
			cierreLote.Costototalporcentaje = responseValidarCF.CostoTotalPorcentaje

			listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)

		}
	}
	return
}

/*
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("============Fecha Operacion======")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Fechaoperacion, movimientoDetalle.FechaOrigenCompra)
	fmt.Println("============Fecha Presentacion===")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.FechaCierre, movimientoCabecera.FechaPresentacion)
	fmt.Println("============Establecimiento======")
	fmt.Printf("cl: %v - movimiento: %v \n", movimientoCabecera.EstablecimientoNro, strNroEstablecimiento)

	fmt.Println("============Nro.Tarjeta==========")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Nrotarjeta, movimientoDetalle.NroTarjetaXl)
	fmt.Println("============Nor.Atorizacion======")
	fmt.Printf("cl: %v - movimiento: %v \n", movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion)

	fmt.Println("============Lote=================")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.ExternalloteId, movimientoDetalle.Lote)

	fmt.Println("============Ticket===============")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Nroticket, movimientoDetalle.NroCupon)

	fmt.Println("============Cuota================")
	fmt.Printf("cl: %v - movimiento: %v \n", valorCuota, movimientoDetalle.PlanCuota)

	fmt.Println("============Importe==============")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe))
	fmt.Println("=================================")
	fmt.Println("=================================")
*/
