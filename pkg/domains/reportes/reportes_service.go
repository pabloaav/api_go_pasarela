package reportes

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/commonsdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	filtros_reportes "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/reportes"
)

type ReportesService interface {
	/* 1 OBTENER CLIENTES: esto nos permite filtrar los pagos de cada clientes*/
	GetClientes(request reportedtos.RequestPagosClientes) (response administraciondtos.ResponseFacturacionPaginado, erro error)

	/* REPORTES GENERALES : reportes de todos los pagos (se generan desde el frontend) */
	GetPagosReportes(request reportedtos.RequestPagosPeriodo) (response []reportedtos.ResponsePagosPeriodo, erro error)
	ResultPagosReportes(request []reportedtos.ResponsePagosPeriodo, paginacion filtros.Paginacion) (response reportedtos.ResponseListaPagoPeriodo, erro error)

	/* REPORTES GENERADOS: estos son los reportes que se enviaran a cada cliente via email */
	GetPagosClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)                                                                /* Todos los pagos*/
	GetRendicionClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)                                                            /* Pagos acreditados */
	GetReversionesClientes(requestCliente administraciondtos.ResponseFacturacionPaginado, request reportedtos.RequestPagosClientes, filtroValidacion reportedtos.ValidacionesFiltro) (response []reportedtos.ResponseClientesReportes, erro error) /* Pagos revertidos */

	/* ENVIAR REPORTES : se envia por correo electronico los pagos a cada clientes*/
	SendPagosClientes(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error)
	SendLiquidacionClientes(request []reportedtos.ResultMovLiquidacion) (errorFile []reportedtos.ResponseCsvEmailError, erro error)

	/* REPORTES DE COBRANZAS PARA CLIENTES(batch): se genera un archivo txt */
	GetPagoItems(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItems, erro error)
	BuildPagosItems(request []reportedtos.ResponsePagosItems) (response []reportedtos.ResultPagosItems)
	ValidarEsctucturaPagosItems(request []reportedtos.ResultPagosItems) error
	SendPagosItems(ctx context.Context, request []reportedtos.ResultPagosItems, filtro reportedtos.RequestPagosClientes) error // tambien debe retornar la lista de pagos para insertar en la tabla pagoslotes(pagos que ya fueron enviados)

	// generar comprobante de liquidacion (DPEC)
	GetRecaudacion(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosLiquidacion, erro error)
	BuildPagosLiquidacion(request []reportedtos.ResponsePagosLiquidacion) (response []reportedtos.ResultPagosLiquidacion)

	// recaudacion diaria
	GetRecaudacionDiaria(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseMovLiquidacion, erro error)
	BuildMovLiquidacion(request []reportedtos.ResponseMovLiquidacion) (response []reportedtos.ResultMovLiquidacion)

	/* REPORTES COBRANZAS Y RENDICIONES (rentas)*/
	GetCobranzas(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error)         /* Pagos */
	GetCobranzasTemporal(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error) /* Pagos + Comisiones*/
	GetRendiciones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseRendiciones, erro error)     /* transferencias */
	GetReversiones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseReversiones, erro error)     /* reversiones */

	NotificacionErroresReportes(errorFile []reportedtos.ResponseCsvEmailError) (erro error)

	/* REPORTES MOVIMIENTOS-COMISIONES */
	MovimientosComisionesService(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error)
	MovimientosComisionesTemporales(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error)

	/* REPORTES COBRANZAS-CLIENTES */
	GetCobranzasClientesService(request reportedtos.RequestCobranzasClientes) (res reportedtos.ResponseCobranzasClientes, erro error)
	/* REPORTES RENDICIONES-CLIENTES */
	GetRendicionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseRendicionesClientes, erro error)
	/* REPORTES RENDICIONES-CLIENTES */
	GetReversionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseReversionesClientes, erro error)

	// Reportes Informacion general
	GetPeticiones(request reportedtos.RequestPeticiones) (response reportedtos.ResponsePeticiones, erro error)
	GetLogs(request reportedtos.RequestLogs) (response reportedtos.ResponseLogs, erro error)
	GetNotificaciones(request reportedtos.RequestNotificaciones) (response reportedtos.ResponseNotificaciones, erro error)
	GetReportesEnviadosService(request reportedtos.RequestReportesEnviados) (response reportedtos.ResponseReportesEnviados, erro error)

	EnumerarReportesEnviadosService(request reportedtos.RequestReportesEnviados) (erro error)
	CopiarNumeroReporteOriginal() (erro error)

	GetCuentaByApiKeyService(apikey string) (cuenta *entities.Cuenta, erro error)
}

type reportesService struct {
	repository     ReportesRepository
	administracion administracion.Service
	util           util.UtilService
	commons        commons.Commons
	factory        ReportesFactory
	factoryEmail   ReportesSendFactory
	store          util.Store
}

func NewService(rm ReportesRepository, adm administracion.Service, util util.UtilService, c commons.Commons, storage util.Store) ReportesService {
	reporte := &reportesService{
		repository: rm,
		// apilinkService:   link
		administracion: adm,
		util:           util,
		commons:        c,
		factory:        &procesarReportesFactory{},
		factoryEmail:   &enviaremailFactory{},
		store:          storage,
	}
	return reporte

}

func (s *reportesService) GetCuentaByApiKeyService(apikey string) (cuenta *entities.Cuenta, erro error) {
	return s.repository.GetCuentaByApiKeyRepository(apikey)
}
func (s *reportesService) GetPagosReportes(request reportedtos.RequestPagosPeriodo) (response []reportedtos.ResponsePagosPeriodo, erro error) {
	// 1 obtener estado pendiente para luego filtrar los pagos
	filtro := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtro)
	if err != nil {
		erro = err
		return
	}

	// 2  obtener channels para luego filtrar cada pago con su cierre lote
	// debin
	canalDebin, erro := s.util.FirstOrCreateConfiguracionService("CHANNEL_DEBIN", "Nombre del canal debin", "debin")
	if erro != nil {
		return
	}
	filtroChannelDebin := filtros.ChannelFiltro{
		Channel: canalDebin,
	}

	channelDebin, erro := s.administracion.GetChannelService(filtroChannelDebin)

	if erro != nil && channelDebin.Id < 1 {
		return
	}

	// offline
	canalOffline, erro := s.util.FirstOrCreateConfiguracionService("CHANNEL_OFFLINE", "Nombre del canal debin", "offline")
	if erro != nil {
		return
	}
	filtroChannelOffline := filtros.ChannelFiltro{
		Channel: canalOffline,
	}

	channelOffline, erro := s.administracion.GetChannelService(filtroChannelOffline)

	if erro != nil && channelOffline.Id < 1 {
		return
	}

	// 3 obtener pagos del periodo
	pagos, erro := s.repository.GetPagosReportes(request, estadoPendiente[0].ID)
	if erro != nil {
		return
	}

	var listaPagoApilink []string
	var listaPagoOffline []string
	var listaPagoPrisma []string
	for _, pago := range pagos {
		// logs.Info(pago.ID)
		var valorCupon entities.Monto
		if pago.PagoIntentos[len(pago.PagoIntentos)-1].Valorcupon == 0 {
			valorCupon = pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount
		} else {
			valorCupon = pago.PagoIntentos[len(pago.PagoIntentos)-1].Valorcupon
		}

		var fechaRendicion string
		var nroreferencia string
		var comision_porcentaje float64
		var comision_porcentaje_iva float64
		var importe_comision_sobre_tap float64
		var importe_comision_sobre_tap_iva float64
		var costo_fijo_transaccion float64
		var importe_rendido float64
		if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos) > 0 {
			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions) > 0 {
				comision_porcentaje = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Porcentaje * 100), 2)
				importe_comision_sobre_tap = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Monto.Float64()), 2)
			}

			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos) > 0 {
				comision_porcentaje_iva = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Porcentaje * 100
				importe_comision_sobre_tap_iva = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Monto.Float64()), 2)
			}

			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos) > 0 {
				for _, mov := range pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos {
					if mov.Tipo == "D" {
						fechaRendicion = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[1].Movimientotransferencia[0].FechaOperacion.Format("02-01-2006")
						nroreferencia = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[1].Movimientotransferencia[0].ReferenciaBancaria
					} else {
						importe_rendido = s.util.ToFixed((mov.Monto.Float64()), 4)
					}
				}
			}

		}
		if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelDebin.Id) {
			costo_fijo_transaccion = float64(pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Channel.Channelaranceles[0].Importe)
		}
		if !pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.IsZero() {
			fechaPago := pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.Format("02-01-2006")
			response = append(response, reportedtos.ResponsePagosPeriodo{
				Cliente:                 pago.PagosTipo.Cuenta.Cliente.Cliente,
				Cuenta:                  pago.PagosTipo.Cuenta.Cuenta,
				Pagotipo:                pago.PagosTipo.Pagotipo,
				IdPago:                  pago.ID,
				ExternalReference:       pago.ExternalReference,
				Estado:                  pago.PagoEstados.Nombre,
				ChannelId:               pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID,
				ExternalId:              pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID,
				TransactionId:           pago.PagoIntentos[len(pago.PagoIntentos)-1].TransactionID,
				Barcode:                 pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode,
				IdExterno:               pago.ExternalReference,
				MedioPago:               pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Mediopago,
				Pagador:                 strings.ToUpper(pago.PagoIntentos[len(pago.PagoIntentos)-1].HolderName),
				DniPagador:              pago.PagoIntentos[len(pago.PagoIntentos)-1].HolderNumber,
				Cuotas:                  uint(pago.PagoIntentos[len(pago.PagoIntentos)-1].Installmentdetail.Cuota),
				FechaPago:               fechaPago,
				FechaRendicion:          fechaRendicion,
				Amount:                  s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64()), 4),
				AmountPagado:            s.util.ToFixed((valorCupon.Float64()), 4),
				CftCoeficiente:          uint(pago.PagoIntentos[len(pago.PagoIntentos)-1].Installmentdetail.Coeficiente),
				ComisionPorcentaje:      comision_porcentaje,
				ComisionPorcentajeIva:   comision_porcentaje_iva,
				ImporteComisionSobreTap: importe_comision_sobre_tap,
				ImporteIvaComisionTap:   importe_comision_sobre_tap_iva,
				CostoFijoTransaccion:    costo_fijo_transaccion,
				ImporteRendido:          importe_rendido,
				ReferenciaBancaria:      nroreferencia,
			})
			if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelDebin.Id) {
				listaPagoApilink = append(listaPagoApilink, pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID)
			} else if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelOffline.Id) {
				listaPagoOffline = append(listaPagoOffline, pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode)
			} else {
				listaPagoPrisma = append(listaPagoPrisma, pago.PagoIntentos[len(pago.PagoIntentos)-1].TransactionID)
			}

		}
	}

	var listasPagos []reportedtos.TipoFactory
	listasPagos = append(listasPagos, reportedtos.TipoFactory{TipoApilink: listaPagoApilink}, reportedtos.TipoFactory{TipoOffline: listaPagoOffline}, reportedtos.TipoFactory{TipoPrisma: listaPagoPrisma})

	reportes, err := s.obtenerReportes(listasPagos)
	if err != nil {
		erro = err
		return
	}
	// actualzar pagos con valores obtenidos de cada cierrelote
	for i := range response {
		for j := range reportes {
			if response[i].ExternalId == reportes[j].Pago || response[i].TransactionId == reportes[j].Pago || response[i].Barcode == reportes[j].Pago {
				response[i].Nroestablecimiento = reportes[j].NroEstablecimiento
				response[i].NroLiquidacion = reportes[j].NroLiquidacion
				response[i].FechaPresentacion = reportes[j].FechaPresentacion
				response[i].FechaAcreditacion = reportes[j].FechaAcreditacion
				response[i].ArancelPorcentaje = reportes[j].ArancelPorcentaje
				response[i].RetencionIva = reportes[j].RetencionIva
				response[i].ImporteMinimo = reportes[j].Importeminimo
				response[i].ImporteMaximo = reportes[j].Importemaximo
				response[i].ArancelPorcentajeMinimo = reportes[j].ArancelPorcentajeMinimo
				response[i].ArancelPorcentajeMaximo = reportes[j].ArancelPorcentajeMaximo
				response[i].ImporteArancel = reportes[j].ImporteArancel
				response[i].ImporteArancelIva = reportes[j].ImporteArancelIva
				response[i].ImporteArancelIvaMov = reportes[j].ImporteArancalIvaMov
				response[i].ImporteCft = reportes[j].ImporteCft
				response[i].ImporteNetoCobrado = reportes[j].ImporteNetoCobrado
				response[i].Revertido = reportes[j].Revertido
				response[i].Enobservacion = reportes[j].Enobservacion
				response[i].Cantdias = reportes[j].Cantdias
			}
		}
	}

	return

}

func (s *reportesService) ResultPagosReportes(request []reportedtos.ResponsePagosPeriodo, paginacion filtros.Paginacion) (response reportedtos.ResponseListaPagoPeriodo, erro error) {
	var responseTemporal []reportedtos.ResultadoPagosPeriodo
	var contador int64
	var recorrerHasta int32
	for _, listaPago := range request {
		contador++
		resp := reportedtos.ResultadoPagosPeriodo{
			Cliente:                 listaPago.Cliente,
			Cuenta:                  listaPago.Cuenta,
			Pagotipo:                listaPago.Pagotipo,
			ExternalReference:       listaPago.ExternalReference,
			IdPago:                  listaPago.IdPago,
			Estado:                  listaPago.Estado,
			MedioPago:               listaPago.MedioPago,
			Pagador:                 listaPago.Pagador,
			Dni:                     listaPago.DniPagador,
			Cuotas:                  listaPago.Cuotas,
			Nroestablecimiento:      listaPago.Nroestablecimiento,
			NroLiquidacion:          listaPago.NroLiquidacion,
			FechaPago:               listaPago.FechaPago,
			FechaPresentacion:       listaPago.FechaPresentacion,
			FechaAcreditacion:       listaPago.FechaAcreditacion,
			FechaRendicion:          listaPago.FechaRendicion,
			Amount:                  listaPago.Amount,
			AmountPagado:            listaPago.AmountPagado,
			ArancelPorcentaje:       listaPago.ArancelPorcentaje,
			CftCoeficiente:          listaPago.CftCoeficiente,
			RetencionIva:            listaPago.RetencionIva,
			ImporteMinimo:           listaPago.ImporteMinimo,
			ImporteMaximo:           listaPago.ImporteMaximo,
			ArancelPorcentajeMinimo: listaPago.ArancelPorcentajeMinimo,
			ArancelPorcentajeMaximo: listaPago.ArancelPorcentajeMaximo,
			CostoFijoTransaccion:    listaPago.CostoFijoTransaccion,
			ImporteArancel:          listaPago.ImporteArancel,
			ImporteArancelIva:       listaPago.ImporteArancelIva,
			ImporteArancelIvaMov:    s.util.ToFixed(listaPago.ImporteArancelIvaMov, 2),
			ImporteCft:              listaPago.ImporteCft,
			ComisionPorcentaje:      listaPago.ComisionPorcentaje,
			ComisionPorcentajeIva:   listaPago.ComisionPorcentajeIva,
			ImporteComisionSobreTap: listaPago.ImporteComisionSobreTap,
			ImporteIvaComisionTap:   listaPago.ImporteIvaComisionTap,
			ImporteRendido:          listaPago.ImporteRendido,
			ImporteNetoCobrado:      listaPago.ImporteNetoCobrado,
			ReferenciaBancaria:      listaPago.ReferenciaBancaria,
			Revertido:               listaPago.Revertido,
			Enobservacion:           listaPago.Enobservacion,
			Cantdias:                listaPago.Cantdias,
		}
		responseTemporal = append(responseTemporal, resp)
	}
	if paginacion.Number > 0 && paginacion.Size > 0 {
		response.Meta = _setPaginacion(paginacion.Number, paginacion.Size, contador)
	}
	recorrerHasta = response.Meta.Page.To
	if response.Meta.Page.CurrentPage == response.Meta.Page.LastPage {
		recorrerHasta = response.Meta.Page.Total
	}
	if recorrerHasta == 0 {
		recorrerHasta = int32(contador)
	}

	if len(responseTemporal) > 0 {
		for i := response.Meta.Page.From; i < recorrerHasta; i++ {
			response.PagosByPeriodo = append(response.PagosByPeriodo, responseTemporal[i])
			response.TotalImporteRendidio += responseTemporal[i].ImporteRendido
		}
		response.TotalImporteRendidio = s.util.ToFixed((response.TotalImporteRendidio), 2)
	}
	return
}

func (s *reportesService) obtenerReportes(listaPagos []reportedtos.TipoFactory) (response []reportedtos.ResponseFactory, erro error) {

	for _, listaPago := range listaPagos {
		var tipoReporte string
		if len(listaPago.TipoApilink) > 0 {
			tipoReporte = "debin"
		} else if len(listaPago.TipoOffline) > 0 {
			tipoReporte = "offline"
		} else if len(listaPago.TipoPrisma) > 0 {
			tipoReporte = "prisma"
		}

		if tipoReporte != "" {
			metodoProcesarReporte, err := s.factory.GetProcesarReportes(tipoReporte)
			if err != nil {
				erro = err
				return
			}

			logs.Info("Procesando reportes tipo: " + tipoReporte)

			listaReporteProcesada := metodoProcesarReporte.ResponseReportes(s, listaPago)

			response = append(response, listaReporteProcesada...)
		}
	}

	return
}

func (s *reportesService) GetClientes(request reportedtos.RequestPagosClientes) (response administraciondtos.ResponseFacturacionPaginado, erro error) {
	filtro := filtros.ClienteFiltro{
		Id:              request.Cliente,
		CargarContactos: true,
		CargarCuentas:   true,
		ClientesIds:     request.ClientesIds,
	}
	response, erro = s.administracion.GetClientesService(filtro)
	return
}

func (s *reportesService) GetPagosClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {

	// OBTENER ESTADOS DE LOS PAGOS:

	filtroEstadoPendiente := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtroEstadoPendiente)
	if err != nil {
		erro = err
		return
	}

	//aprobado (credito , debito y offline)
	paid, erro := s.util.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}
	estado_paid, err := s.administracion.GetPagoEstado(filtroPagosEstado)
	if err != nil {
		erro = err
		return
	}
	//si no se obtiene el estado del pago no se puede seguir
	if estado_paid[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	//autorizado (debin)
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, err := s.administracion.GetPagoEstado(filtroPagoEstado)

	if err != nil {
		erro = err
		return
	}

	//si no se obtiene el estado del pago no se puede seguir
	if pagoEstadoAcreditado[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID_AUTORIZADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID_AUTORIZADO,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	var pagoestados []uint
	var fechaI time.Time
	var fechaF time.Time
	pagoestados = append(pagoestados, estado_paid[0].ID, pagoEstadoAcreditado[0].ID)
	if filtro.FechaInicio.IsZero() {
		// Entro por proceso background
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	for _, cliente := range request.Clientes {
		var cantoperaciones int64
		var totalcobrado float64
		filtroPagos := reportedtos.RequestPagosPeriodo{
			ClienteId:   uint64(cliente.Id),
			FechaInicio: fechaI,
			FechaFin:    fechaF,
			PagoEstados: pagoestados,
		}

		listaPagos, err := s.repository.GetPagosReportes(filtroPagos, estadoPendiente[0].ID)
		if err != nil {
			erro = err
			return
		}
		var pagos []reportedtos.PagosReportes
		if len(listaPagos) > 0 {
			for _, pago := range listaPagos {
				// monto := s.util.ToFixed((pago.Amount.Float64()), 4)
				cantoperaciones = cantoperaciones + 1
				totalcobrado += pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64()
				medio_pago, _ := s.commons.RemoveAccents(pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Mediopago)
				pagos = append(pagos, reportedtos.PagosReportes{
					Cuenta:    pago.PagosTipo.Cuenta.Cuenta,
					Id:        pago.ExternalReference,
					FechaPago: pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.Format("02-01-2006"),
					MedioPago: medio_pago,
					Tipo:      pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Channel.Nombre,
					Estado:    string(pago.PagoEstados.Nombre),
					Monto:     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64(), 2))),
				})
			}
		}
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Email:           cliente.Emails, //[]string{cliente.Email},
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				Fecha:           fechaI.Format("02-01-2006"),
				Pagos:           pagos,
				CantOperaciones: fmt.Sprintf("%v", cantoperaciones),
				TotalCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalcobrado, 2))),
			})
		}
	}
	return
}

func (s *reportesService) GetRendicionClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if filtro.FechaInicio.IsZero() {
		// si los filtros recibidos son ceros toman la fecha actual
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		// a las fechas se le restan un dia ya sea por backgraund o endpoint
		// fechaI = fechaI.AddDate(0, 0, int(-1))
		// fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		// fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		// fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
		fechaI = filtro.FechaInicio
		fechaF = filtro.FechaFin
	}
	logs.Info(fechaI)
	logs.Info(fechaF)
	for _, cliente := range request.Clientes {
		var totalCobrado entities.Monto
		var total entities.Monto
		var totalIVA entities.Monto
		var totalComision entities.Monto
		var totalReversion entities.Monto
		var cantOperaciones int

		// totales generales
		var rendido entities.Monto

		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId: uint64(cliente.Id),
			// ClienteId:   5,                             //Prueba con cliente 6
			FechaInicio: fechaI, // descomentar esta linea cuando se pasa a dev y produccion
			// se envian pagos del dia anterior
			FechaFin: fechaF,
		}

		// TODO se obtienen transferencias del cliente indicado en el filtro
		listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
		if err != nil {
			erro = err
			return
		}
		var pagos []*reportedtos.ResponseReportesRendiciones
		var movrevertidos []entities.Movimiento
		var pagosintentos []uint64
		var pagosintentosrevertidos []uint64
		var filtroMov reportedtos.RequestPagosPeriodo
		var totalCliente reportedtos.ResponseTotales
		if len(listaTransferencia) > 0 {
			for _, transferencia := range listaTransferencia {
				if !transferencia.Reversion {
					pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
				}
				if transferencia.Reversion {
					pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
				}
			}
			filtroMov = reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentos,
				TipoMovimiento:                  "C",
				CargarComisionImpuesto:          true,
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
			}
		}

		// en el caso de que existieran reversiones
		if len(pagosintentosrevertidos) > 0 {
			filtroRevertidos := reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentosrevertidos,
				TipoMovimiento:                  "C",
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarReversionReporte:          true,
				CargarComisionImpuesto:          true,
			}
			movrevertidos, err = s.repository.GetMovimiento(filtroRevertidos)
			if err != nil {
				erro = err
				return
			}
		}

		if len(pagosintentos) > 0 {
			mov, err := s.repository.GetMovimiento(filtroMov)
			if err != nil {
				erro = err
				return
			}
			var resulRendiciones []*reportedtos.ResponseReportesRendiciones
			for _, m := range mov {
				cantOperaciones = cantOperaciones + 1
				cantidadBoletas := len(m.Pagointentos.Pago.Pagoitems)
				total += m.Monto
				totalCobrado += m.Pagointentos.Amount

				var comision entities.Monto
				var iva entities.Monto
				if len(m.Movimientocomisions) > 0 {
					comision = m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto + m.Movimientocomisions[len(m.Movimientocomisions)-1].Montoproveedor
					iva = m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto + m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Montoproveedor
				} else {
					comision = 0
					iva = 0
				}
				totalComision += comision
				totalIVA += iva

				resulRendiciones = append(resulRendiciones, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:           m.PagointentosId,
					Cuenta:                  m.Cuenta.Cuenta,
					Id:                      m.Pagointentos.Pago.ExternalReference,
					FechaCobro:              m.Pagointentos.PaidAt.Format("02-01-2006"),
					ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))),
					ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Monto.Float64(), 2))),
					CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
					Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
					Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
					Concepto:                "Transferencia",
				})
			}

			totalCliente = reportedtos.ResponseTotales{
				// Totales: reportedtos.Totales{
				// 	CantidadOperaciones: fmt.Sprintf("%v", cantOperaciones),
				// 	TotalCobrado:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				// 	TotalRendido:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(total.Float64(), 4))),
				// },
				Detalles: resulRendiciones,
			}
		}
		// vuelvo a comparar con las transferenicas para asignar fecha de deposito
		if len(totalCliente.Detalles) > 0 {
			for _, transferencia := range listaTransferencia {
				for _, t := range totalCliente.Detalles {
					if transferencia.Movimiento.PagointentosId == t.PagoIntentoId {
						t.FechaDeposito = transferencia.FechaOperacion.Format("02-01-2006")
					}
				}
			}
		}

		if len(movrevertidos) > 0 {
			for _, mr := range movrevertidos {
				totalReversion += mr.Monto
				cantOperaciones = cantOperaciones + 1
				var comision entities.Monto
				var iva entities.Monto
				if len(mr.Movimientocomisions) > 0 {
					comision = mr.Movimientocomisions[len(mr.Movimientocomisions)-1].Monto + mr.Movimientocomisions[len(mr.Movimientocomisions)-1].Montoproveedor
					iva = mr.Movimientoimpuestos[len(mr.Movimientoimpuestos)-1].Monto + mr.Movimientoimpuestos[len(mr.Movimientoimpuestos)-1].Montoproveedor
				} else {
					comision = 0
					iva = 0
				}
				// totalComision += comision
				// totalIVA += iva
				totalCliente.Detalles = append(totalCliente.Detalles, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:     mr.PagointentosId,
					Cuenta:            mr.Cuenta.Cuenta,
					Id:                mr.Pagointentos.Pago.ExternalReference,
					ImporteDepositado: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr.Monto.Float64(), 2))),
					Concepto:          "Reversion",
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
				})
			}

		}
		pagos = totalCliente.Detalles

		rendido = total + totalReversion
		totalCliente = reportedtos.ResponseTotales{
			Totales: reportedtos.Totales{
				CantidadOperaciones: fmt.Sprintf("%v", cantOperaciones),
				TotalCobrado:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				TotalRendido:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(rendido.Float64(), 4))),
				TotalComision:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComision.Float64(), 4))),
				TotalIva:            fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalIVA.Float64(), 4))),
				TotalRevertido:      fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalReversion.Float64(), 4))),
			},
		}
		//&ultimo paso
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				Email:           cliente.Emails,
				Fecha:           fechaI.Format("02-01-2006"),
				Rendiciones:     pagos,
				CantOperaciones: totalCliente.Totales.CantidadOperaciones,
				TotalCobrado:    totalCliente.Totales.TotalCobrado,
				TotalRevertido:  totalCliente.Totales.TotalRevertido,
				RendicionTotal:  totalCliente.Totales.TotalRendido,
				TotalIva:        totalCliente.Totales.TotalIva,
				TotalComision:   totalCliente.Totales.TotalComision,
			})
		}
	}
	return
}

func (s *reportesService) GetReversionesClientes(requestCliente administraciondtos.ResponseFacturacionPaginado, request reportedtos.RequestPagosClientes, filtroValidacion reportedtos.ValidacionesFiltro) (response []reportedtos.ResponseClientesReportes, erro error) {

	for _, cliente := range requestCliente.Clientes {
		// fechaI, fechaF, err := s.commons.FormatFecha()
		// if err != nil {
		// 	return
		// }
		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId:   uint64(cliente.Id),
			FechaInicio: request.FechaInicio,
			FechaFin:    request.FechaFin,
		}
		listaPagos, err := s.repository.GetReversionesReportes(filtro, filtroValidacion)
		if err != nil {
			erro = err
			return
		}
		var pagos []reportedtos.Reversiones
		var cantOperacion int64
		if len(listaPagos) > 0 {
			for _, value := range listaPagos {
				cantOperacion = cantOperacion + 1
				var revertido reportedtos.Reversiones
				var pagoRevertido reportedtos.PagoRevertido
				var itemsRevertido []reportedtos.ItemsRevertidos
				//var itemRevertido reportedtos.ItemsRevertidos
				var intentoPagoRevertido reportedtos.IntentoPagoRevertido
				revertido.EntityToReversiones(value)
				pagoRevertido.EntityToPagoRevertido(value.PagoIntento.Pago)
				if len(value.PagoIntento.Pago.Pagoitems) > 0 {
					for _, valueItem := range value.PagoIntento.Pago.Pagoitems {
						var itemRevertido reportedtos.ItemsRevertidos
						itemRevertido.EntityToItemsRevertidos(valueItem)
						itemsRevertido = append(itemsRevertido, itemRevertido)
					}
				}
				intentoPagoRevertido.EntityToIntentoPagoRevertido(value.PagoIntento)
				pagoRevertido.Items = itemsRevertido
				pagoRevertido.IntentoPago = intentoPagoRevertido
				revertido.PagoRevertido = pagoRevertido
				pagos = append(pagos, revertido)
			}
		}

		// if len(listaPagos) > 0 {
		// 	for _, pago := range listaPagos {
		// 		pagos = append(pagos, reportedtos.Reversiones{
		// 			Cuenta:    pago.PagoIntento.Pago.PagosTipo.Cuenta.Cuenta,
		// 			Id:        pago.PagoIntento.Pago.ExternalReference,
		// 			MedioPago: pago.PagoIntento.Mediopagos.Mediopago,
		// 			Monto:     fmt.Sprintf("%v", s.util.ToFixed(entities.Monto(pago.Amount).Float64(), 4)),
		// 		})
		// 	}
		// }
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Email:           []string{cliente.Email},
				Fecha:           request.FechaInicio.Format("02-01-2006"), //fechaI.AddDate(0, 0, int(-1)).Format("02-01-2006"),
				Reversiones:     pagos,
				CantOperaciones: fmt.Sprintf("%v", cantOperacion),
				TipoArchivoPdf:  true,
			})
		}
	}

	// transformar la data adaptando al pdf

	for i, cliente := range response {
		reversionesData := transformarDatos(response[i])
		err := commons.GetReversionesPdf(reversionesData, cliente.Clientes, request.FechaInicio.Format("02-01-2006"))
		if err != nil {
			erro = err
			logs.Error(err.Error())
		}
	}

	return
}

func (s *reportesService) SendPagosClientes(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Clientes)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Clientes),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {
			var tipo_archivo string
			var contentType string
			var asunto string
			var nombreArchivo string
			var tipo string
			var titulo string
			if len(cliente.Pagos) > 0 {
				cliente.TipoReporte = "pagos"
				asunto = "Pagos realizados " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "cobranzas"
				tipo = "cobrados"
			} else if len(cliente.Rendiciones) > 0 {
				cliente.TipoReporte = "rendiciones"
				asunto = "Recaudación WEE! " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "rendiciones"
				tipo = "rendidos"
			} else if len(cliente.Reversiones) > 0 {
				cliente.TipoReporte = "revertidos"
				asunto = "Pagos revertidos " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "reversiones"
				tipo = "revertidos"
			}
			if cliente.TipoArchivoPdf {
				tipo_archivo = ".pdf"
				contentType = "application/pdf"
			} else {
				tipo_archivo = ".csv"
				contentType = "application/pdf"
				metodoConvertirCvs, err := s.factoryEmail.SendEnviarEmail(cliente.TipoReporte)
				if err != nil {
					erro = err
					return
				}
				logs.Info("Procesando reportes tipo: " + cliente.TipoReporte)
				convertircvs := metodoConvertirCvs.SendReportes(ruta, nombreArchivo, cliente)
				if convertircvs != nil {
					erro = convertircvs
					return
				}
			}

			var campo_adicional = []string{"pagos"}
			var email = cliente.Email //[]string{cliente.Email}
			filtro := utildtos.RequestDatosMail{
				Email:            email,
				Asunto:           asunto,
				From:             "Wee.ar!",
				Nombre:           cliente.Clientes,
				Mensaje:          "reportes de pagos: #0",
				CamposReemplazar: campo_adicional,
				Descripcion: utildtos.DescripcionTemplate{
					Fecha:   cliente.Fecha,
					Cliente: cliente.RazonSocial,
					Cuit:    cliente.Cuit,
				},
				Totales: utildtos.TotalesTemplate{
					Titulo:       titulo,
					TipoReporte:  tipo,
					Elemento:     "pagos",
					Cantidad:     cliente.CantOperaciones,
					TotalCobrado: cliente.TotalCobrado,
					TotalRendido: cliente.RendicionTotal,
				},
				AdjuntarEstado: true,
				Attachment: utildtos.Attachment{
					Name:        fmt.Sprintf("%s%s", nombreArchivo, tipo_archivo),
					ContentType: contentType,
					WithFile:    true,
				},
				TipoEmail: "reporte",
			}
			/*enviar archivo csv por correo*/
			/* en el caso de no registrar error al enviar correo se guardan los datos del reporte*/
			erro = s.util.EnviarMailService(filtro)
			logs.Info(erro)
			if erro != nil {
				erro = fmt.Errorf("no se no pudo enviar rendicion al %v", cliente.Clientes)
				errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
					Archivo: filtro.Attachment.Name,
					Error:   fmt.Sprintf("servicio email: %v", erro),
				})
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "EnviarMailService",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return erro
				}
				/* informar el error al enviar el emial pero se debe continuar enviando los siguientes archivos a otros clientes */
			} else {
				// guardar datos del reporte
				//si el archivo se sube correctamente se registra en tabla movimientos lotes
				pagos := reportedtos.ToEntityRegistroReporte(cliente)
				siguienteNroReporte, erro := s.repository.GetLastReporteEnviadosRepository(pagos, false)
				if erro != nil {
					mensaje := fmt.Errorf("no se pudo obtener nro reporte para el reporte de pago enviado al cliente %+v", cliente.Clientes).Error()
					logs.Info(mensaje)
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "GetLastReporteEnviadosRepository",
						Mensaje:       mensaje,
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
					}

				}

				pagos.Nro_reporte = siguienteNroReporte
				erro = s.repository.SaveGuardarDatosReporte(pagos)
				if erro != nil {
					mensaje := fmt.Errorf("no se pudieron registrar datos del reporte de pago enviado al cliente %+v", cliente.Clientes).Error()
					logs.Info(mensaje)
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "SaveGuardarDatosReporte",
						Mensaje:       mensaje,
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
					}

				}
			}

			// una vez enviado el correo se elimina el archivo csv
			erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf("%s.csv", nombreArchivo))
			if erro != nil {
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "BorrarArchivos",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return nil, erro
				}
			}
		}

	}
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) GetPagoItems(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItems, erro error) {
	// OBTENER ESTADOS DE LOS PAGOS:
	//aprobado (credito , debito y offline)
	paid, erro := s.util.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}
	estado_paid, err := s.administracion.GetPagoEstado(filtroPagosEstado)
	if err != nil {
		erro = err
		return
	}
	//si no se obtiene el estado del pago no se puede seguir
	if estado_paid[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	//autorizado (debin)
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, err := s.administracion.GetPagoEstado(filtroPagoEstado)

	if err != nil {
		erro = err
		return
	}

	//si no se obtiene el estado del pago no se puede seguir
	if pagoEstadoAcreditado[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID_AUTORIZADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID_AUTORIZADO,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	var pagoestados []uint
	var fechaI time.Time
	var fechaF time.Time
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	fecha := commons.ConvertFechaString(fechaI) // fecha de creacion del archivo
	pagoestados = append(pagoestados, estado_paid[0].ID, pagoEstadoAcreditado[0].ID)

	for _, cliente := range request.Clientes {
		if cliente.ReporteBatch {
			filtro := reportedtos.RequestPagosPeriodo{
				ClienteId:   uint64(cliente.Id),
				FechaInicio: fechaI,
				FechaFin:    fechaF,
				PagoEstados: pagoestados,
			}
			logs.Info(filtro)
			// se debe obtener los pagos aptobados y autorizados(debin)
			pagosItems, err := s.repository.GetPagosBatch(filtro)
			if err != nil {
				erro = err
				return
			}
			// obtener lote del cliente
			lote, err := s.repository.GetLastLote(filtro)
			if err != nil {
				erro = err
				return
			}

			var pagos []entities.Pago
			var idpg []uint

			// solo los pagos de tipo C y que no se informaron en algun lote
			if len(pagosItems) > 0 {
				for _, pg := range pagosItems {
					if pg.PagoIntentos[len(pg.PagoIntentos)-1].Amount > 0 && len(pg.Pagolotes) == 0 {
						pagos = append(pagos, pg)
						idpg = append(idpg, pg.ID)
					}
				}
			}
			// respuesta solo si existen  pagos para ese cliente
			if len(pagos) > 0 {
				response = append(response, reportedtos.ResponsePagosItems{
					Clientes: reportedtos.ClientesResponse{
						Id:          cliente.Id,
						Cliente:     cliente.Cliente,
						RazonSocial: cliente.NombreFantasia,
						Email:       cliente.Email,
					},
					Fecha: fecha,
					Pagos: pagos,
					PagLotes: reportedtos.PagLotes{
						Idpg:          idpg,
						Idcliente:     cliente.Id,
						Lote:          int(lote.Lote) + 1,
						Fechalote:     fecha,
						Cliente:       cliente.Cliente,
						NombreReporte: cliente.NombreReporte,
					},
				})
			}

		}
	}
	return
}

func (s *reportesService) BuildPagosItems(request []reportedtos.ResponsePagosItems) (response []reportedtos.ResultPagosItems) {
	var cabeceraArchivo reportedtos.CabeceraArchivo
	var cabeceraLote reportedtos.CabeceraLote
	var colaArchivo reportedtos.ColaArchivo
	for _, pago := range request {
		// CABECERAS

		cabeceraArchivo = reportedtos.CabeceraArchivo{
			RecordCode:   "1",
			CreateDate:   pago.Fecha,
			OrigenName:   commons.EspaciosBlanco("WEE", 25, "RIGHT"),
			ClientNumber: commons.EspaciosBlanco("", 9, "RIGHT"),
			ClientName:   commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:       commons.EspaciosBlanco("", 54, "RIGHT"),
		}
		cabeceraLote = reportedtos.CabeceraLote{
			RecordCodeLote: "3",
			CreateDateLote: pago.Fecha,
			BatchNumber:    commons.AgregarCeros(6, pago.PagLotes.Lote), // esta longitud puede variar (su longitud maxima es 6)
			Description:    commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:         commons.EspaciosBlanco("", 82, "RIGHT"),
		}
		// DETALLES
		var resultItems []reportedtos.ResultItems
		var detalle1 reportedtos.DetalleTransaccion
		var detalle2 reportedtos.DetalleDescripcion
		var payment_date string
		var payment_time string
		var fileCount int64
		var totalFileAmount entities.Monto
		for _, items := range pago.Pagos {
			payment_date = commons.ConvertFechaString(items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt)
			payment_time = fmt.Sprintf("%v%v", items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt.Hour(), items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt.Minute())
			for _, pi := range items.Pagoitems {
				fileCount = fileCount + 1
				totalFileAmount += pi.Amount
				detalle1 = reportedtos.DetalleTransaccion{
					RecordCodeTransaccion: "5",
					RecordSequence:        commons.AgregarCeros(5, 0),
					TransactionCode:       commons.AgregarCeros(2, 0),
					WorkDate:              commons.AgregarCeros(8, 0),
					TransferDate:          commons.AgregarCeros(8, 0),
					AccountNumber:         commons.EspaciosBlanco(pi.Description, 21, "RIGHT")[0:21], // pago items
					CurrencyCode:          commons.EspaciosBlanco("", 3, "RIGHT"),
					Amount:                commons.AgregarCeros(14, int(pi.Amount)), // pago items
					TerminalId:            commons.EspaciosBlanco("", 6, "RIGHT"),
					PaymentDate:           payment_date,                                        //Pago intento
					PaymentTime:           commons.AgregarCerosString(payment_time, 4, "LEFT"), // pago intento
					SeqNumber:             commons.AgregarCeros(4, 0),
					Filler:                commons.EspaciosBlanco("", 48, "RIGHT"),
				}
				detalle2 = reportedtos.DetalleDescripcion{
					RecordCodeLote: "6",
					BarCode:        commons.AgregarCerosString(pi.Identifier, 80, "LEFT")[0:80], // pago items
					TypeCode:       commons.EspaciosBlanco("", 1, "RIGHT"),
					Filler:         commons.EspaciosBlanco("", 50, "RIGHT"),
				}
				resultItems = append(resultItems, reportedtos.ResultItems{
					DetalleTransaccion: detalle1,
					DetalleDescripcion: detalle2,
				})
			}
		}
		// COLA DE ARCHIVO
		colaArchivo = reportedtos.ColaArchivo{
			RecordCodeCola:    "9",
			CreateDateCola:    pago.Fecha,
			TotalBatches:      commons.AgregarCeros(6, 0),
			FilePaymentCount:  commons.AgregarCeros(7, int(fileCount)),
			FilePaymentAmount: commons.AgregarCeros(12, int(totalFileAmount)), // total acumulado(detalles)
			Filler:            commons.AgregarCerosString("0", 38, "LEFT"),
			FileCount:         commons.AgregarCeros(7, 0),
			Filler2:           commons.EspaciosBlanco("", 53, "RIGHT"),
		}
		// RESPUESTA
		response = append(response, reportedtos.ResultPagosItems{
			PagLotes:        pago.PagLotes,
			CabeceraArchivo: cabeceraArchivo,
			CabeceraLote:    cabeceraLote,
			ResultItems:     resultItems,
			// DetalleTransaccion: detalle1,
			// DetalleDescripcion: detalle2,
			ColaArchivo: colaArchivo,
		})
	}
	return
}

func (s *reportesService) ValidarEsctucturaPagosItems(request []reportedtos.ResultPagosItems) (err error) {
	var registroDescripcion reportedtos.EstructuraRegistrosBatch
	for _, items := range request {
		//validar datos de la cabecera
		err = validarRegistroCabeceraArchivo(items.CabeceraArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroCabeceraLote(items.CabeceraLote, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroDetalleTransaccion(items.ResultItems, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_TRANSACCION, err.Error())
			return errors.New(mensaje)
		}

		// err = validarRegistroDetalleDescripcion(items.DetalleDescripcion, registroDescripcion)
		// if err != nil {
		// 	mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_DESCRIPCION, err.Error())
		// 	return errors.New(mensaje)
		// }

		err = validarColaArchivo(items.ColaArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_COLA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}
	}
	return nil
}

func validarRegistroCabeceraArchivo(cabeceraArchivo reportedtos.CabeceraArchivo, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraArchivo.ValidarCabeceraArchivo(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func validarRegistroCabeceraLote(cabeceraLote reportedtos.CabeceraLote, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraLote.ValidarCabeceraLote(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func validarRegistroDetalleTransaccion(descripcionTransaccion []reportedtos.ResultItems, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {

	for _, detalle := range descripcionTransaccion {
		if detalle.DetalleTransaccion.RecordCodeTransaccion == "5" {
			err := detalle.DetalleTransaccion.ValidarDetalleTransaccion(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		if detalle.DetalleDescripcion.RecordCodeLote == "6" {
			err := detalle.DetalleDescripcion.ValidarDetalleDescripcion(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

// func validarRegistroDetalleDescripcion(descripcionTransaccion []reportedtos.DetalleDescripcion, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
// 	for _, detalle := range descripcionTransaccion {
// 		err := detalle.ValidarDetalleDescripcion(&registroDescripcion)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

func validarColaArchivo(cabeceraLote reportedtos.ColaArchivo, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraLote.ValidarColaArchivo(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// ENVIAR ARCHIVO POR FTP ETC
func (s *reportesService) SendPagosItems(ctx context.Context, request []reportedtos.ResultPagosItems, filtro reportedtos.RequestPagosClientes) (erro error) {

	// obtener fecha de envio del archivo
	// por defecto toma el ultimo dia
	var fechaArchivo time.Time

	if filtro.FechaInicio.IsZero() {
		fechaArchivo = time.Now()
	} else {
		fechaArchivo = filtro.FechaInicio
	}

	/* en esta ruta se crearan los archivos para enviar */
	//ruta := fmt.Sprintf("..%s", config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}
	for _, pago := range request {
		// 1 CREAR EL NOMBRE DEL ARCHIVO
		// se crea con el nombre(WEE) + la fecha de hoy
		// fechaArchivo := time.Now()
		nombreArchivo := "WEE" + fechaArchivo.Format("020106")
		nombreArchivoPagosItems := commonsdtos.FileName{
			RutaBase:  ruta + "/",
			Nombre:    nombreArchivo,
			Extension: "txt",
			UsaFecha:  false,
		}
		rutaDetalle := s.commons.CreateFileName(nombreArchivoPagosItems)

		// 2 CREAR EL ARCHIVO
		_, erro = s.commons.CreateFile(rutaDetalle)
		if erro != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_CREAR_MOVIMIENTOS_LOTES, erro.Error())
			return errors.New(mensaje)
		}

		// 3 ABRIR EL ARCHIVO
		file, err := s.commons.LeerArchivo(rutaDetalle)
		if err != nil {
			err = erro
			return
		}

		// 4 ESCRIBIR EN EL ARCHIVO
		err = s.buildArchivo(file, pago)
		if err != nil {
			err = erro
			return
		}

		// 5 GUARDAR CAMBIOS y CERRAR ARCHIVO
		err = s.commons.GuardarCambios(file)
		if err != nil {
			err = erro
			return
		}

		//si el archivo se sube correctamente se registra en tabla movimientos lotes
		lotes := reportedtos.ToEntity(pago.PagLotes)
		erro = s.repository.SavePagosLotes(ctx, lotes)
		if erro != nil {
			mensaje := fmt.Errorf("no se pudieron registrar los siguiente movimientos en la tabla lotes %+v", pago.PagLotes.Idpg).Error()
			logs.Info(mensaje)
			return
		}

		// 6 ENVIAR ARCHIVO POR SFTP
		erro = s.SubirArchivo(ctx, nombreArchivoPagosItems, pago.PagLotes.NombreReporte, file)

		// 7 en el caso de que no se pueda enviar el archivo se deben dar de bajas los movimientos lotes creados
		if erro != nil {
			logs.Info("ocurrio error al enviar el archivo " + fmt.Sprintf("%v", erro))
			// 8.1 - En caso de que me tire un error se dan de bajas los movimientos lotes creados anteriormente
			err = s.repository.BajaPagosLotes(ctx, lotes, erro.Error())

			if err != nil {
				// 8.1.1 - En caso de que no se puede cancelar los movimientos aviso al usuario para que intervenga manualmente.
				mensaje_baja := fmt.Errorf("no se pudieron dar de bajas los siguientes movientos lotes %+v", pago.PagLotes.Idpg).Error()

				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionTransferencia,
					Descripcion: fmt.Sprintf("atención los siguientes movimientos de lotes no pudieron ser cancelados, movimientosId: %s", mensaje_baja),
				}
				erro = s.util.CreateNotificacionService(notificacion)
				if erro != nil {
					logs.Error(erro.Error())
				}
				erro = err
				return erro
			}
			return
		}

		// 5 UNA VEZ ENVIADO EL ARCHIVO , ELIMINAR EL ARCHIVO CREADO TEMPORALEMTE
		erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf("%s.txt", nombreArchivo))
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return erro
			}
		}

	}

	// 6 BORRAR DIRECTORIO CREADO PARA EL REPORTE
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) buildArchivo(archivo *os.File, request reportedtos.ResultPagosItems) (erro error) {

	/* 	CABECERA DE ARCHIVO */
	cabArchivo := []string{request.CabeceraArchivo.RecordCode, request.CabeceraArchivo.CreateDate,
		request.CabeceraArchivo.OrigenName, request.CabeceraArchivo.ClientNumber,
		request.CabeceraArchivo.ClientName, request.CabeceraArchivo.Filler, "\n"}
	resultcabeceraArhivo := commons.JoinString(cabArchivo)

	// escribir cabecera archivo
	erro = s.commons.EscribirArchivo(resultcabeceraArhivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE ARCHIVO */

	/* CABECERA DE LOTE */
	cabLote := []string{request.CabeceraLote.RecordCodeLote, request.CabeceraLote.CreateDateLote,
		request.CabeceraLote.BatchNumber, request.CabeceraLote.Description, request.CabeceraLote.Filler, "\n"}
	resultcabeceraLote := commons.JoinString(cabLote)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcabeceraLote, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	/* DETALLES TRANSACCION */
	for _, detalle := range request.ResultItems {
		var detalleTransaccion = []string{}
		var detalleDescripcion = []string{}
		detalleTransaccion = []string{detalle.DetalleTransaccion.RecordCodeTransaccion, detalle.DetalleTransaccion.RecordSequence, detalle.DetalleTransaccion.TransactionCode,
			detalle.DetalleTransaccion.WorkDate, detalle.DetalleTransaccion.TransferDate, detalle.DetalleTransaccion.AccountNumber, detalle.DetalleTransaccion.CurrencyCode, detalle.DetalleTransaccion.Amount,
			detalle.DetalleTransaccion.TerminalId, detalle.DetalleTransaccion.PaymentDate, detalle.DetalleTransaccion.PaymentTime, detalle.DetalleTransaccion.SeqNumber, detalle.DetalleTransaccion.Filler, "\n"}
		resultdetalleTransaccion := commons.JoinString(detalleTransaccion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleTransaccion, archivo)
		if erro != nil {
			return erro
		}

		detalleDescripcion = []string{detalle.DetalleDescripcion.RecordCodeLote, detalle.DetalleDescripcion.BarCode, detalle.DetalleDescripcion.TypeCode,
			detalle.DetalleDescripcion.Filler, "\n"}
		resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
		if erro != nil {
			return erro
		}
		/* 	END DETALLES TRANSACCION */
	}

	// /* DETALLES DESCRIPCION */
	// for _, detalle2 := range request.DetalleDescripcion {
	// 	var detalleDescripcion = []string{}
	// 	detalleDescripcion = []string{detalle2.RecordCodeLote, detalle2.BarCode, detalle2.TypeCode,
	// 		detalle2.Filler, "\n"}
	// 	resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
	// 	// escribir cabecera lote
	// 	erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
	// 	if erro != nil {
	// 		return erro
	// 	}
	// 	/* 	END DETALLES DESCRIPCION */
	// }

	/* COLA DE ARCHIVO */
	colaArchivo := []string{request.ColaArchivo.RecordCodeCola, request.ColaArchivo.CreateDateCola,
		request.ColaArchivo.TotalBatches, request.ColaArchivo.FilePaymentCount, request.ColaArchivo.FilePaymentAmount,
		request.ColaArchivo.Filler, request.ColaArchivo.FileCount, request.ColaArchivo.Filler2}
	resultcolaArchivo := commons.JoinString(colaArchivo)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcolaArchivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	return
}

func (s *reportesService) SubirArchivo(ctx context.Context, rutaArchivos commonsdtos.FileName, cliente string, archivo *os.File) (erro error) {
	// rutaDestino := config.DIR_KEY_REPORTES
	rutaDestinoReporte := strings.Replace(config.DIR_KEY_REPORTES, "*", cliente, 3)
	data, filename, filetypo, err := util.LeerArchivo(rutaDestinoReporte, rutaArchivos.RutaBase, rutaArchivos.Nombre+"."+rutaArchivos.Extension)
	if err != nil {
		erro = err
		logs.Error(err)
		return
	}

	erro = s.store.PutObject(ctx, data, filename, filetypo)
	if erro != nil {
		logs.Error("No se pudo guardar el archivo")
		return
	}

	defer archivo.Close()
	return
}

func _setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))
	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}
	return
}

func (s *reportesService) GetCobranzasTemporal(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}
	//Obtengo la fecha actual y la paso a un formato para poder compararla con la fecha que viene por parametro
	currentDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	parsedcurrentDate, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	//No puede consultar cobranzas de la fecha actual por que todavia no termina el dia, y el registro estaria incompleto
	if date.Equal(parsedcurrentDate) {
		erro = errors.New("no existen registros para la fecha ingresada")
		return
	}

	// consultar datos de la cuenta
	filtroCuenta := filtros.CuentaFiltro{
		ApiKey: apikey,
	}

	cuenta, erro := s.administracion.GetCuenta(filtroCuenta)
	if erro != nil {
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	// var pagoestados []uint
	// pagoestados = append(pagoestados, estado_paid[0].ID /* pagoEstadoAcreditado[0].ID */)
	//buscar transferencias los pagos correspondiente al cliente(apiKey)

	fechaInicioString := s.commons.GetDateFirstMoment(date)
	fechaFinString := s.commons.GetDateLastMoment(date)

	filtro := filtros_reportes.CobranzasClienteFiltro{
		FechaInicio: fechaInicioString,
		FechaFin:    fechaFinString,
		ClienteId:   int(cuenta.ClientesID),
		CuentaId:    int(cuenta.Id),
	}

	var resultado []reportedtos.DetallesPagosCobranza

	apilink, err := s.repository.CobranzasApilink(filtro)
	resultado = append(resultado, apilink...)

	rapipago, err := s.repository.CobranzasRapipago(filtro)
	resultado = append(resultado, rapipago...)

	prisma, err := s.repository.CobranzasPrisma(filtro)
	resultado = append(resultado, prisma...)

	multipago, err := s.repository.CobranzasMultipago(filtro)
	resultado = append(resultado, multipago...)

	if err != nil {
		erro = errors.New(ERROR_COBRANZAS_CLIENTES)
		return
	}

	var totalCobrado entities.Monto
	var descuentoComisionIva entities.Monto
	var totalNeto entities.Monto

	response.AccountId = commons.AgregarCerosString(fmt.Sprintf("%v", cuenta.Id), 6, "LEFT")
	response.ReportDate = date

	if len(resultado) > 0 {
		var resulRendiciones []reportedtos.ResponseDetalleCobranza
		for _, pago := range resultado {
			// netAmount := m.Monto - (comision + iva)

			amountPaid := entities.Monto(pago.TotalPago)
			netFee := entities.Monto(pago.Comision)
			ivaFee := entities.Monto(pago.Iva)
			netAmount := amountPaid - (netFee + ivaFee)

			totalCobrado += amountPaid
			totalNeto += netAmount
			descuentoComisionIva += (netFee + ivaFee)

			// Se intenta analizar la cadena como una fecha y hora en el formato RFC3339
			paymentDate, err := time.Parse(time.RFC3339, pago.FechaPago)
			if err != nil {
				fmt.Println("Error al formatear fecha de pago:", err.Error())
			}

			resulRendiciones = append(resulRendiciones, reportedtos.ResponseDetalleCobranza{
				InformedDate: date,
				// NOTE informacion del pago
				RequestId:         int64(pago.Id),
				ExternalReference: pago.Referencia,
				PayerName:         pago.PayerName,
				Description:       pago.Descripcion,
				PaymentDate:       paymentDate,
				Channel:           pago.CanalPago,
				AmountPaid:        amountPaid.Float64(),
				// NOTE disponible solo cuando hay movimientos
				NetFee:      netFee.Float64(),
				IvaFee:      ivaFee.Float64(),
				NetAmount:   netAmount.Float64(),
				AvailableAt: pago.FechaCobro,
			})

			response = reportedtos.ResponseCobranzas{
				AccountId:      commons.AgregarCerosString(fmt.Sprintf("%v", cuenta.Id), 6, "LEFT"),
				ReportDate:     date,
				TotalCollected: s.util.ToFixed(totalCobrado.Float64(), 2),
				TotalGrossFee:  s.util.ToFixed(descuentoComisionIva.Float64(), 2),
				TotalNetAmount: s.util.ToFixed(totalNeto.Float64(), 2),
				Data:           resulRendiciones,
			}
		}
	}

	return
}

func (s *reportesService) GetCobranzas(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	// obtener estado pendiente para filtrar pagos
	filtroEstadoPago := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtroEstadoPago)
	if err != nil {
		erro = err
		return
	}

	//NOTE aprobado (credito , debito y offline)
	paid, erro := s.util.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}
	estado_paid, err := s.administracion.GetPagoEstado(filtroPagosEstado)
	if err != nil {
		erro = err
		return
	}
	//si no se obtiene el estado del pago no se puede seguir
	if estado_paid[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	//NOTE pagos estados autorizado (debin)
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, err := s.administracion.GetPagoEstado(filtroPagoEstado)

	if err != nil {
		erro = err
		return
	}

	//si no se obtiene el estado del pago no se puede seguir
	if pagoEstadoAcreditado[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID_AUTORIZADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID_AUTORIZADO,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	// consultar datos de la cuenta
	filtroCuenta := filtros.CuentaFiltro{
		ApiKey: apikey,
	}

	cuenta, erro := s.administracion.GetCuenta(filtroCuenta)
	if erro != nil {
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	var pagoestados []uint
	pagoestados = append(pagoestados, estado_paid[0].ID, pagoEstadoAcreditado[0].ID)
	//buscar transferencias los pagos correspondiente al cliente(apiKey)
	filtro := reportedtos.RequestPagosPeriodo{
		ApiKey:      apikey,
		FechaInicio: date,
		FechaFin:    date,
		PagoEstados: pagoestados,
	}

	// 3 obtener pagos del periodo
	listapagos, erro := s.repository.GetPagosReportes(filtro, estadoPendiente[0].ID)
	if erro != nil {
		return
	}

	var totalCobrado entities.Monto
	var descuentoComisionIva entities.Monto
	var totalNeto entities.Monto

	response.AccountId = commons.AgregarCerosString(fmt.Sprintf("%v", cuenta.Id), 6, "LEFT")
	response.ReportDate = date

	if len(listapagos) > 0 {
		var resulRendiciones []reportedtos.ResponseDetalleCobranza
		for _, m := range listapagos {
			// controlar que este pago sea un movimiento
			var comision entities.Monto
			var iva entities.Monto
			var availableAt string
			var netAmount entities.Monto
			totalCobrado += m.PagoIntentos[len(m.PagoIntentos)-1].Amount
			if len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos) > 0 {
				if m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Tipo == "C" {
					netAmount = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Monto
					availableAt = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].CreatedAt.Format("2006-01-02 15:04:05")
					totalNeto += m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Monto
					if len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions) > 0 && len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos) > 0 {
						comision = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Monto + m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Montoproveedor
						iva = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Monto + m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Montoproveedor
					}
					descuentoComisionIva += comision + iva
				}
			}
			// netAmount := m.Monto - (comision + iva)
			resulRendiciones = append(resulRendiciones, reportedtos.ResponseDetalleCobranza{
				InformedDate: date,
				// NOTE informacion del pago
				RequestId:         int64(m.ID),
				ExternalReference: m.ExternalReference,
				PayerName:         m.PayerName,
				Description:       m.Description,
				PaymentDate:       m.PagoIntentos[len(m.PagoIntentos)-1].PaidAt,
				Channel:           m.PagoIntentos[len(m.PagoIntentos)-1].Mediopagos.Mediopago,
				AmountPaid:        s.util.ToFixed(m.PagoIntentos[len(m.PagoIntentos)-1].Amount.Float64(), 2),
				// NOTE disponible solo cuando hay movimientos
				NetFee:      s.util.ToFixed(comision.Float64(), 2),
				IvaFee:      s.util.ToFixed(iva.Float64(), 2),
				NetAmount:   s.util.ToFixed(netAmount.Float64(), 2),
				AvailableAt: availableAt,
			})

			response = reportedtos.ResponseCobranzas{
				AccountId:      commons.AgregarCerosString(fmt.Sprintf("%v", listapagos[0].PagosTipo.CuentasID), 6, "LEFT"),
				ReportDate:     date,
				TotalCollected: s.util.ToFixed(totalCobrado.Float64(), 2),
				TotalGrossFee:  s.util.ToFixed(descuentoComisionIva.Float64(), 2),
				TotalNetAmount: s.util.ToFixed(totalNeto.Float64(), 2),
				Data:           resulRendiciones,
			}
		}
	}

	return
}

func (s *reportesService) GetRendiciones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseRendiciones, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	// consultar datos de la cuenta
	filtroCuenta := filtros.CuentaFiltro{
		ApiKey: apikey,
	}

	cuenta, erro := s.administracion.GetCuenta(filtroCuenta)
	if erro != nil {
		return
	}

	//buscar transferencias los pagos correspondiente al cliente(apiKey)
	filtro := reportedtos.RequestPagosPeriodo{
		ApiKey:      apikey,
		FechaInicio: date,
		FechaFin:    date,
	}
	response.AccountId = commons.AgregarCerosString(fmt.Sprintf("%v", cuenta.Id), 6, "LEFT")
	response.ReportDate = date
	// 3 obtener pagos del periodo
	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}
	var totalCredit uint64
	var creditAmount entities.Monto
	var debitAmount entities.Monto

	// NOTE este es el total credit - debit
	var settlementAmount entities.Monto

	var pagosintentos []uint64
	var filtroMov reportedtos.RequestPagosPeriodo
	if len(listaTransferencia) > 0 {
		for _, transferencia := range listaTransferencia {
			if !transferencia.Reversion { // se descartan las operaciones que son reversiones
				pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
			}
		}
		filtroMov = reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentos,
			TipoMovimiento:                  "C",
			CargarComisionImpuesto:          true,
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
		}
	}

	if len(pagosintentos) > 0 {
		mov, err := s.repository.GetMovimiento(filtroMov)
		if err != nil {
			erro = err
			return
		}
		var resulRendiciones []reportedtos.ResponseDetalleRendiciones
		for _, m := range mov {
			if !m.Reversion {
				totalCredit = totalCredit + 1
				creditAmount += m.Pagointentos.Amount
				settlementAmount += m.Monto // total rendido
				var resultDebit entities.Monto
				// var resultDebit float64
				comision := m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto + m.Movimientocomisions[len(m.Movimientocomisions)-1].Montoproveedor
				iva := m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto + m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Montoproveedor
				resultDebit = comision + iva // total de comisiones + iva de cada operacion
				debitAmount += resultDebit
				// NOTE detalles de las operacions de rendicion
				resulRendiciones = append(resulRendiciones, reportedtos.ResponseDetalleRendiciones{
					RequestId:         m.Pagointentos.PagosID,
					ExternalReference: m.Pagointentos.Pago.ExternalReference,
					Credit:            s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2),
					Debit:             s.util.ToFixed(resultDebit.Float64(), 2),
				})
			}
		}

		response = reportedtos.ResponseRendiciones{
			AccountId:        commons.AgregarCerosString(fmt.Sprintf("%v", mov[0].CuentasId), 6, "LEFT"),
			ReportDate:       date,
			TotalCredits:     totalCredit,                                   // total de rendiciones
			CreditAmount:     s.util.ToFixed(creditAmount.Float64(), 2),     // Importe total cobrado
			TotalDebits:      totalCredit,                                   // total de comisiones + iva
			DebitAmount:      s.util.ToFixed(debitAmount.Float64(), 2),      // monto total de comisiones + iva
			SettlementAmount: s.util.ToFixed(settlementAmount.Float64(), 2), // monto neto rendido
			Data:             resulRendiciones,                              // detalle de las operaciones
		}

	}

	return
}

func (s *reportesService) GetReversiones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseReversiones, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// Parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	filtro := reportedtos.RequestPagosPeriodo{
		ApiKey:          apikey,
		FechaInicio:     date,
		FechaFin:        date,
		CargarReversion: true,
	}

	listaTransferencias, err := s.repository.GetReversionesReportesRenta(filtro)
	if err != nil {
		erro = err
		return
	}

	var (
		TotalChargeback float64 // Monto total de reverciones
		pagosintentos   []uint64
		filtroMov       reportedtos.RequestPagosPeriodo
	)

	if len(listaTransferencias) > 0 {
		for _, transferencia := range listaTransferencias {
			pagosintentos = append(pagosintentos, uint64(transferencia.Movimiento.PagointentosId))
		}

		filtroMov = reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentos,
			TipoMovimiento:                  "C",
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarCliente:                   true,
			CargarMedioPago:                 true,
			OrdenadoFecha:                   true,
		}
	}

	if len(pagosintentos) > 0 {

		mov, err := s.repository.GetMovimiento(filtroMov)
		if err != nil {
			erro = err
			return
		}

		var resulRereverciones []reportedtos.ResponseDetalleReversiones

		for _, m := range mov {
			if m.Reversion {

				TotalChargeback += -(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))
				formatoDeseado := "02-01-2006 15:04:05"
				fechaFormateada := m.Pagointentos.PaidAt.Format(formatoDeseado)

				resulRereverciones = append(resulRereverciones, reportedtos.ResponseDetalleReversiones{
					InformedDate:      fechaFormateada,
					RequestID:         int(m.ID),
					ExternalReference: m.Pagointentos.Pago.ExternalReference,
					PayerName:         m.Pagointentos.Pago.PayerName,
					Description:       m.Pagointentos.Pago.Description,
					Channel:           m.Pagointentos.Mediopagos.Channel.Nombre,
					RevertedAmount:    -(m.Pagointentos.Amount.Float64()),
				})
			}

		}

		response = reportedtos.ResponseReversiones{
			AccountID:       commons.AgregarCerosString(fmt.Sprintf("%v", listaTransferencias[0].Movimiento.CuentasId), 6, "LEFT"),
			ReportDate:      request.Date,
			TotalChargeback: s.util.ToFixed(TotalChargeback, 2),
			Data:            resulRereverciones,
		}

	}

	return
}
func (s *reportesService) GetRecaudacion(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosLiquidacion, erro error) {
	var fechaI time.Time       // este seria la fecha de cobro
	var fechaProceso time.Time // fecha que se procesa el archivo
	var fechaF time.Time
	var lote int64
	var cantpagos int
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaProceso = fechaI
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaProceso = filtro.FechaInicio
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	for _, cliente := range request.Clientes {
		// este reporte es enviado solo a dpec
		if cliente.ReporteBatch {
			request := filtros.PagoIntentoFiltros{
				PagoEstadosIds:              []uint64{4, 7},
				CargarPago:                  true,
				CargarPagoTipo:              true,
				CargarPagoEstado:            true,
				CargarCuenta:                true,
				PagoIntentoAprobado:         true,
				CargarPagoCalculado:         true,
				FechaPagoInicio:             fechaI,
				FechaPagoFin:                fechaF,
				ClienteId:                   uint64(cliente.Id),
				CargarPagoItems:             true,
				CargarMovimientosTemporales: true,
				Channel:                     true,
			}

			//NOTE consultar los pagos intentos del dia anterior:
			// 1 Deben estar aprobados y el estado debe ser calculado
			// se deben comparar con los lotes informados en el archivo batch
			// se deben informar la misma cantidad en los 2 reportes cobranzas y liquidacion
			pagos, err := s.administracion.GetPagosIntentosCalculoComisionRepository(request)
			if err != nil {
				erro = err
				return
			}

			if len(pagos) > 0 {
				cantpagos = len(pagos)
				var pagos_id []uint64
				for _, pg := range pagos {
					pagos_id = append(pagos_id, uint64(pg.PagosID))
				}

				// obtener cantidad de lotes del cliente del dia
				//NOTE se debe verificar que las liquidaciones sean iguales a las cobranzas informadas
				filtro := reportedtos.RequestPagosPeriodo{
					ClienteId: uint64(cliente.Id),
					Pagos:     pagos_id,
				}
				lote, err = s.repository.GetCantidadLotes(filtro)
				if err != nil {
					erro = err
					return
				}
			}

			// respuesta solo si existen  pagos para ese cliente
			if int64(cantpagos) == lote {
				response = append(response, reportedtos.ResponsePagosLiquidacion{
					Clientes: reportedtos.Clientes{
						Id:          cliente.Id,
						Cliente:     cliente.Cliente,
						RazonSocial: cliente.RazonSocial,
						Email:       cliente.Emails,
						Cuit:        cliente.Cuit,
					},
					FechaCobro:   fechaI,
					FechaProceso: fechaProceso,
					Pagos:        pagos,
				})
			} else {
				erro = errors.New("error cantidad de pagos a informar es distinto a los lotes enviados en cobranzas")
				return
			}
		}
	}
	return
}

func (s *reportesService) BuildPagosLiquidacion(request []reportedtos.ResponsePagosLiquidacion) (response []reportedtos.ResultPagosLiquidacion) {
	var cabecera reportedtos.Clientes
	var fechaCobro string
	var fechaProceso string
	for _, pago := range request {

		// CABECERAS
		cabecera = pago.Clientes
		fechaCobro = pago.FechaCobro.Format("02-01-2006")
		fechaProceso = pago.FechaProceso.Format("02-01-2006")

		// DETALLES
		// var resultItemsMediopago []reportedtos.MedioPagoItems
		//  ? 1cobrado total y por medio de pago
		// credito
		var importecobrado entities.Monto
		var importecobradoCredito entities.Monto
		var importecobradoDebito entities.Monto
		var importecobradoDebin entities.Monto

		// debito
		var importeADepositar entities.Monto
		var importeADepositarCredito entities.Monto
		var importeADepositarDebito entities.Monto
		var importeADepositarDebin entities.Monto

		// cantidad boletas
		var cantidadTotalBoletas int
		var cantidadTotalBoletasCredit int
		var cantidadTotalBoletasDebito int
		var cantidadTotalBoletasDebin int

		// comision
		var comisionTotal entities.Monto
		var comisionCredito entities.Monto
		var comisionDebito entities.Monto
		var comisionDebin entities.Monto
		// iva
		var ivaTotal entities.Monto
		var ivaCredito entities.Monto
		var ivaDebito entities.Monto
		var ivaDebin entities.Monto

		// detales torales parciales
		var mcredit reportedtos.MedioPagoCredit
		var mdebito reportedtos.MedioPagoDebit
		var mdebin reportedtos.MedioPagoDebin
		for _, items := range pago.Pagos {
			importecobrado += items.Amount
			importeADepositar += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
			cantidadTotalBoletas += len(items.Pago.Pagoitems)
			comisionTotal += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
			ivaTotal += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

			if items.Mediopagos.Channel.ID == 1 {
				importecobradoCredito += items.Amount
				importeADepositarCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasCredit += len(items.Pago.Pagoitems)
				comisionCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mcredit = reportedtos.MedioPagoCredit{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoCredito.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarCredito.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasCredit),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionCredito.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaCredito.Float64(), 2))),
				}
			}

			if items.Mediopagos.Channel.ID == 2 {
				importecobradoDebito += items.Amount
				importeADepositarDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasDebito += len(items.Pago.Pagoitems)
				comisionDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mdebito = reportedtos.MedioPagoDebit{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebito.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebito.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasDebito),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionDebito.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaDebito.Float64(), 2))),
				}
			}

			if items.Mediopagos.Channel.ID == 4 {
				importecobradoDebin += items.Amount
				importeADepositarDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasDebin += len(items.Pago.Pagoitems)
				comisionDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mdebin = reportedtos.MedioPagoDebin{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebin.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebin.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasDebin),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionDebin.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaDebin.Float64(), 2))),
				}
			}
		}

		// RESPUESTA
		response = append(response, reportedtos.ResultPagosLiquidacion{
			Cabeceras:    cabecera,
			FechaCobro:   fechaCobro,
			FechaProceso: fechaProceso,
			MedioPagoItems: reportedtos.MedioPagoItems{
				MedioPagoCredit: mcredit,
				MedioPagoDebit:  mdebito,
				MedioPagoDebin:  mdebin,
			},
			Totales: reportedtos.TotalesALiquidar{
				ImporteCobrado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobrado.Float64(), 2))),
				ImporteADepositar:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositar.Float64(), 2))),
				CantidadTotalBoletas: fmt.Sprintf("%v", cantidadTotalBoletas),
				ComisionTotal:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionTotal.Float64(), 2))),
				IvaTotal:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaTotal.Float64(), 2))),
			},
		})
	}
	return
}

func (s *reportesService) SendLiquidacionClientes(request []reportedtos.ResultMovLiquidacion) (errorFile []reportedtos.ResponseCsvEmailError, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Cabeceras.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Cabeceras.Cliente)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Cabeceras.Cliente),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {

			logs.Info("Procesando reportes tipo: liquidacion diaria ")
			asunto := "Liquidación Wee! " + cliente.FechaProceso
			nombreArchivo := cliente.Cabeceras.Cliente + "-" + cliente.FechaProceso

			// crear en carpeta tempora

			erro = s.util.GetRecaudacionPdf(cliente, ruta, nombreArchivo)
			if erro != nil {
				return errorFile, erro
			}

			if cliente.Cabeceras.EnviarEmail {
				var campo_adicional = []string{"pagos"}
				var email = cliente.Cabeceras.Email //[]string{cliente.Email}
				filtro := utildtos.RequestDatosMail{
					Email:            email,
					Asunto:           asunto,
					From:             "Wee.ar!",
					Nombre:           "Wee.ar!",
					Mensaje:          "reportes de pagos: #0",
					CamposReemplazar: campo_adicional,
					AdjuntarEstado:   true,
					Attachment: utildtos.Attachment{
						Name:        fmt.Sprintf("%s.pdf", nombreArchivo),
						ContentType: "text/csv",
						WithFile:    true,
					},
					TipoEmail: "adjunto",
				}
				/*enviar archivo csv por correo*/
				erro = s.util.EnviarMailService(filtro)
				logs.Info(erro)
				if erro != nil {
					erro = fmt.Errorf("no se no pudo enviar rendicion al %v", cliente.Cabeceras.Cliente)
					errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
						Archivo: filtro.Attachment.Name,
						Error:   fmt.Sprintf("servicio email: %v", erro),
					})
					logs.Error(erro.Error())
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "EnviarMailService",
						Mensaje:       erro.Error(),
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
						// return erro
					}
					/* informar el error al enviar el emial pero se debe continuar enviando los siguientes archivos a otros clientes */
				}
			}

			// una vez enviado el correo se elimina el archivo csv
			erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf("%s.pdf", nombreArchivo))
			if erro != nil {
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "BorrarArchivos",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					return nil, erro
				}
			}
		}

	}
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) GetRecaudacionDiaria(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseMovLiquidacion, erro error) {
	var fechaI time.Time // este seria la fecha de cobr
	var fechaF time.Time
	var fechaProceso time.Time // fecha que se procesa el archivo
	var fechaRendicion string
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaProceso = fechaI
	} else {
		fechaI = filtro.FechaInicio
		fechaF = filtro.FechaFin
		fechaProceso = filtro.FechaInicio
	}

	if filtro.FechaAdicional.IsZero() {
		fechaRendicion = fechaI.AddDate(0, 0, int(+2)).Format("02-01-2006")
	} else if filtro.CargarFechaAdicional {
		fechaRendicion = filtro.FechaAdicional.Format("02-01-2006")
	}

	for _, c := range request.Clientes {
		var orden_diaria bool
		for _, cu := range c.Cuenta {
			logs.Info("procesando orden de pago cliente" + c.Cliente)
			filtroMov := reportedtos.RequestPagosPeriodo{
				FechaInicio:                     fechaI,
				FechaFin:                        fechaF,
				TipoMovimiento:                  "C",
				CargarComisionImpuesto:          true,
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarMedioPago:                 true,
				CuentaId:                        uint64(cu.Id),
			}
			mov, err := s.repository.GetRendicionReportes(filtroMov)
			if err != nil {
				erro = err
				return
			}

			if !c.OrdenDiaria && len(mov) == 0 {
				orden_diaria = false
			}
			if c.OrdenDiaria {
				orden_diaria = true
			} else if !c.OrdenDiaria && len(mov) > 0 {
				orden_diaria = true
			}

			if orden_diaria {
				var detalleliquidacion []entities.Liquidaciondetalles
				if len(mov) > 0 {
					for _, m := range mov {
						detalleliquidacion = append(detalleliquidacion, entities.Liquidaciondetalles{
							PagointentosId: int64(m.PagointentosId),
							MovimientosId:  int64(m.ID),
							CuentasId:      int64(m.CuentasId),
						})
					}
				}

				liquidacion := entities.Movimientoliquidaciones{
					ClientesID:           uint64(c.Id),
					FechaEnvio:           fechaProceso.Format("2006-01-02"),
					LiquidacioneDetalles: detalleliquidacion,
				}
				// guardar y obtener numero de liquidacion(ultimo registro ingresado)
				nroliquidacion, err := s.repository.SaveLiquidacion(liquidacion)
				if err != nil {
					erro = err
					return
				}

				// buscar fecha de rendicion en transferencia (moviminetos tipo D)
				response = append(response, reportedtos.ResponseMovLiquidacion{
					Clientes: reportedtos.Cliente{
						Id:          c.Id,
						Cliente:     c.Cliente,
						RazonSocial: c.RazonSocial,
						Email:       c.Emails,
						Cuit:        c.Cuit,
						EnviarEmail: filtro.EnviarEmail,
						EnviarPdf:   false, // descomentar esta linea cuando se desea generar pdf de dpec
					},
					Cuenta:         cu.Cuenta,
					FechaRendicion: fechaRendicion,
					FechaProceso:   fechaProceso,
					Movimientos:    mov,
					NroLiquidacion: int(nroliquidacion),
				})
			}

		}
	}
	return
}

func (s *reportesService) BuildMovLiquidacion(request []reportedtos.ResponseMovLiquidacion) (response []reportedtos.ResultMovLiquidacion) {
	var cabecera reportedtos.Cliente
	var cuenta string
	var nroliquidacion int
	// var fechaCobro string
	var fechaProceso string
	var fecharendicion string
	for _, pago := range request {

		// CABECERAS
		cabecera = pago.Clientes
		cuenta = pago.Cuenta
		nroliquidacion = pago.NroLiquidacion
		fecharendicion = pago.FechaRendicion
		fechaProceso = pago.FechaProceso.Format("02-01-2006")

		// DETALLES
		// var resultItemsMediopago []reportedtos.MedioPagoItems
		//  ? 1cobrado total y por medio de pago
		// credito
		var importecobrado entities.Monto
		var importecobradoCredito entities.Monto
		var importecobradoDebito entities.Monto
		var importecobradoDebin entities.Monto

		// debito
		var importeADepositar entities.Monto
		var importeADepositarCredito entities.Monto
		var importeADepositarDebito entities.Monto
		var importeADepositarDebin entities.Monto

		// cantidad boletas
		var cantidadTotalBoletas int
		var cantidadTotalBoletasCredit int
		var cantidadTotalBoletasDebito int
		var cantidadTotalBoletasDebin int

		// comision
		var comisionTotal entities.Monto
		var comisionCredito entities.Monto
		var comisionDebito entities.Monto
		var comisionDebin entities.Monto
		// iva
		var ivaTotal entities.Monto
		var ivaCredito entities.Monto
		var ivaDebito entities.Monto
		var ivaDebin entities.Monto

		// detales torales parciales
		var mcredit reportedtos.MedioMovCredit
		var detallecredit []reportedtos.DetalleMov

		var mdebito reportedtos.MedioMovDebit
		var detalledebit []reportedtos.DetalleMov

		var mdebin reportedtos.MedioMovDebin
		var detalledebin []reportedtos.DetalleMov

		for _, items := range pago.Movimientos {
			importecobrado += items.Pagointentos.Amount
			importeADepositar += items.Monto
			cantidadTotalBoletas += len(items.Pagointentos.Pago.Pagoitems)
			comisionTotal += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
			ivaTotal += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			if items.Pagointentos.Mediopagos.ChannelsID == 1 {
				comision := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				detallecredit = append(detallecredit, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 2))),
				})
				importecobradoCredito += items.Pagointentos.Amount
				importeADepositarCredito += items.Monto
				cantidadTotalBoletasCredit += len(items.Pagointentos.Pago.Pagoitems)
				comisionCredito += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaCredito += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
			}

			if items.Pagointentos.Mediopagos.ChannelsID == 2 {
				comision1 := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva1 := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				detalledebit = append(detalledebit, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision1.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva1.Float64(), 2))),
				})
				importecobradoDebito += items.Pagointentos.Amount
				importeADepositarDebito += items.Monto
				cantidadTotalBoletasDebito += len(items.Pagointentos.Pago.Pagoitems)
				comisionDebito += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaDebito += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			}

			if items.Pagointentos.Mediopagos.ChannelsID == 4 {
				comision2 := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva2 := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				// cantidadBoletasDebin += len(items.Pagointentos.Pago.Pagoitems)
				detalledebin = append(detalledebin, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision2.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva2.Float64(), 2))),
				})
				importecobradoDebin += items.Pagointentos.Amount
				importeADepositarDebin += items.Monto
				cantidadTotalBoletasDebin += len(items.Pagointentos.Pago.Pagoitems)
				comisionDebin += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaDebin += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			}
		}
		mcredit = reportedtos.MedioMovCredit{
			Detalle:              detallecredit,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasCredit),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoCredito.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarCredito.Float64(), 2))),
		}
		mdebito = reportedtos.MedioMovDebit{
			Detalle:              detalledebit,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasDebito),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebito.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebito.Float64(), 2))),
		}
		mdebin = reportedtos.MedioMovDebin{
			Detalle:              detalledebin,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasDebin),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebin.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebin.Float64(), 2))),
		}

		// RESPUESTA
		response = append(response, reportedtos.ResultMovLiquidacion{
			Cabeceras:      cabecera,
			NroLiquidacion: commons.AgregarCerosString(fmt.Sprintf("%v", nroliquidacion), 6, "LEFT"),
			FechaProceso:   fechaProceso,
			Cuenta:         cuenta,
			FechaRendicion: fecharendicion,
			MedioPagoItems: reportedtos.MedioMovItems{
				MedioMovCredit: mcredit,
				MedioMovDebit:  mdebito,
				MedioMovDebin:  mdebin,
			},
			Totales: reportedtos.TotalesMovLiquidar{
				ImporteCobrado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobrado.Float64(), 2))),
				ImporteADepositar:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositar.Float64(), 2))),
				CantidadTotalBoletas: fmt.Sprintf("%v", cantidadTotalBoletas),
				ComisionTotal:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionTotal.Float64(), 2))),
				IvaTotal:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaTotal.Float64(), 2))),
			},
		})
	}
	return
}

func (s *reportesService) NotificacionErroresReportes(errorFile []reportedtos.ResponseCsvEmailError) (erro error) {
	// var campos = []string{}
	// var slice_array = []string{}
	// for _, tr := range responseTransferencias {
	// 	slice_array = append(slice_array, tr.Cuenta, fmt.Sprintf("CuentaID%v ", tr.CuentaId), fmt.Sprintf("Origen%v ", tr.Origen), fmt.Sprintf("Destino%v ", tr.Destino), fmt.Sprintf("Importe%v ", tr.Importe), tr.Error, "\n")
	// }
	// mensaje := (strings.Join(slice_array, "-"))
	// var arrayEmail []string
	// // NOTE para  pruebas
	// // arrayEmail = append(arrayEmail, "jose.alarcon@telco.com.ar")
	// arrayEmail = append(arrayEmail, config.EMAIL_TELCO)
	// params := utildtos.RequestDatosMail{
	// 	Email:            arrayEmail,
	// 	Asunto:           "Error transferencias automaticas",
	// 	Nombre:           "Wee!!",
	// 	Mensaje:          mensaje,
	// 	CamposReemplazar: campos,
	// 	From:             "Wee.ar!",
	// 	TipoEmail:        "template",
	// }
	// erro = s.util.EnviarMailService(params)
	// if erro != nil {
	// 	logs.Info("Ocurrio un error al enviar correo notificación transferencias automaticas con error")
	// 	logs.Error(erro.Error())
	// }
	return
}

func transformarDatos(responseClienteReversion reportedtos.ResponseClientesReportes) (reversionesData []commons.ReversionData) {

	// Parsear manualmente las reversiones a la data struct para presentar el pdf
	for _, revers := range responseClienteReversion.Reversiones {
		var data commons.ReversionData
		// Datos del pago en header
		data.Pago.ReferenciaExterna = revers.PagoRevertido.ReferenciaExterna
		data.Pago.MedioPago = revers.MedioPago
		// formatear importe
		montorevertido_int64, _ := strconv.ParseInt(revers.Monto, 10, 64)
		montorevertido_float := entities.Monto(montorevertido_int64).Float64()
		data.Pago.Monto = util.Resolve().FormatNum(montorevertido_float)
		data.Pago.IdPago = revers.PagoRevertido.IdPago
		data.Pago.Estado = revers.PagoRevertido.PagoEstado
		// Datos del intento de pago en herader
		data.Intento.IdIntento = revers.PagoRevertido.IntentoPago.IdIntentoPago
		data.Intento.IdTransaccion = revers.PagoRevertido.IntentoPago.IdTransaccion
		fecha := strings.Split(revers.PagoRevertido.IntentoPago.FechaPago, " ")
		data.Intento.FechaPago = fecha[0]
		data.Intento.Importe = revers.PagoRevertido.IntentoPago.ImportePagado

		for _, item := range revers.PagoRevertido.Items {
			var tempItem commons.ItemsReversionData
			tempItem.Cantidad = item.Cantidad
			tempItem.Descripcion = item.Descripcion
			tempItem.Identificador = item.Identificador

			montoitem_int64, _ := strconv.ParseInt(item.Monto, 10, 64)
			montoitem_float := entities.Monto(montoitem_int64).Float64()
			tempItem.Monto = util.Resolve().FormatNum(montoitem_float)
			data.Items = append(data.Items, tempItem)
		}
		reversionesData = append(reversionesData, data)
	}
	////
	return
}

// for _, revers := range responseClienteReversion.Reversiones {
// 	var data commons.ReversionData
// 	// Datos del pago en header
// 	data.Pago.ReferenciaExterna = revers.PagoRevertido.ReferenciaExterna
// 	data.Pago.MedioPago = revers.MedioPago
// 	data.Pago.Monto = revers.Monto
// 	data.Pago.IdPago = revers.PagoRevertido.IdPago
// 	data.Pago.Estado = revers.PagoRevertido.PagoEstado
// 	// Datos del intento de pago en herader
// 	data.Intento.IdIntento = revers.PagoRevertido.IntentoPago.IdIntentoPago
// 	data.Intento.IdTransaccion = revers.PagoRevertido.IntentoPago.IdTransaccion
// 	fecha := strings.Split(revers.PagoRevertido.IntentoPago.FechaPago, " ")
// 	data.Intento.FechaPago = fecha[0]
// 	data.Intento.Importe = revers.PagoRevertido.IntentoPago.ImportePagado

// 	for _, item := range revers.PagoRevertido.Items {
// 		var tempItem commons.ItemsReversionData
// 		tempItem.Cantidad = item.Cantidad
// 		tempItem.Descripcion = item.Descripcion
// 		tempItem.Identificador = item.Identificador
// 		tempItem.Monto = item.Monto
// 		data.Items = append(data.Items, tempItem)
// 	}

//		reversionesData = append(reversionesData, data)
//	}
func (s *reportesService) MovimientosComisionesService(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.MovimientosComisionesFiltro{
		FechaInicio:     fechaInicioString,
		FechaFin:        fechaFinString,
		ClienteId:       request.ClienteId,
		Number:          request.Number,
		Size:            request.Size,
		UsarFechaPago:   request.FechaPago,
		SoloReversiones: request.SoloReversiones,
	}

	resultado, total, err := s.repository.MovimientosComisionesRepository(filtro)
	if err != nil {
		erro = errors.New(ERROR_MOVIMIENTOS_COMISIONES)
		return
	}

	res.Reportes = total
	res.SetTotales()

	res.Reportes = resultado
	for i := 0; i < len(res.Reportes); i++ {
		var reporte = res.Reportes[i]
		res.Reportes[i].PorcentajeComision = s.util.ToFixed((res.Reportes[i].PorcentajeComision), 4)
		res.Reportes[i].Subtotal = (reporte.MontoComision + reporte.MontoImpuesto)
	}

	if request.Size != 0 {
		res.LastPage = int(math.Ceil(float64(len(total)) / float64(request.Size)))
	}

	return
}

func (s *reportesService) ReversionesComisionesService(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.MovimientosComisionesFiltro{
		FechaInicio:     fechaInicioString,
		FechaFin:        fechaFinString,
		ClienteId:       request.ClienteId,
		Number:          request.Number,
		Size:            request.Size,
		SoloReversiones: true,
	}

	resultado, total, err := s.repository.MovimientosComisionesRepository(filtro)
	if err != nil {
		erro = errors.New(ERROR_MOVIMIENTOS_COMISIONES)
		return
	}

	res.Reportes = total
	res.SetTotales()

	res.Reportes = resultado
	for i := 0; i < len(res.Reportes); i++ {
		var reporte = res.Reportes[i]
		res.Reportes[i].PorcentajeComision = s.util.ToFixed((res.Reportes[i].PorcentajeComision), 4)
		res.Reportes[i].Subtotal = (reporte.MontoComision + reporte.MontoImpuesto)
	}

	if request.Size != 0 {
		res.LastPage = int(math.Ceil(float64(len(total)) / float64(request.Size)))
	}

	return
}

func (s *reportesService) MovimientosComisionesTemporales(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.MovimientosComisionesFiltro{
		FechaInicio:   fechaInicioString,
		FechaFin:      fechaFinString,
		ClienteId:     request.ClienteId,
		Number:        request.Number,
		Size:          request.Size,
		UsarFechaPago: request.FechaPago,
	}

	resultado, total, err := s.repository.MovimientosComisionesTemporales(filtro)
	if err != nil {
		erro = errors.New(ERROR_MOVIMIENTOS_COMISIONES)
		return
	}

	res.Reportes = total
	res.SetTotales()

	res.Reportes = resultado
	for i := 0; i < len(res.Reportes); i++ {
		var reporte = res.Reportes[i]
		res.Reportes[i].PorcentajeComision = s.util.ToFixed((res.Reportes[i].PorcentajeComision), 4)
		res.Reportes[i].Subtotal = (reporte.MontoComision + reporte.MontoImpuesto)
	}

	if request.Size != 0 {
		res.LastPage = int(math.Ceil(float64(len(total)) / float64(request.Size)))
	}

	return
}

func (s *reportesService) GetCobranzasClientesService(request reportedtos.RequestCobranzasClientes) (res reportedtos.ResponseCobranzasClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.CobranzasClienteFiltro{
		FechaInicio: fechaInicioString,
		FechaFin:    fechaFinString,
		ClienteId:   request.ClienteId,
		CuentaId:    request.CuentaId,
	}

	// resultado, err := s.repository.CobranzasClientesRepository(filtro)
	var resultado []reportedtos.DetallesPagosCobranza

	apilink, err := s.repository.CobranzasApilink(filtro)
	resultado = append(resultado, apilink...)

	rapipago, err := s.repository.CobranzasRapipago(filtro)
	resultado = append(resultado, rapipago...)

	prisma, err := s.repository.CobranzasPrisma(filtro)
	resultado = append(resultado, prisma...)

	multipago, err := s.repository.CobranzasMultipago(filtro)
	resultado = append(resultado, multipago...)

	if err != nil {
		erro = errors.New(ERROR_COBRANZAS_CLIENTES)
		return
	}

	var fechas []string
	for _, pago := range resultado {
		var fecha time.Time
		//Solo los pagos con rapipago tienen fechacobro, por lo tanto si tiene, debe tomar esa fecha, no la fecha del paid_at
		//Por que debemos controlar por el dia que se realizo el pago en el rapipago, no el dia que se creo el pago

		if pago.FechaCobro != "" {
			fecha, _ = time.Parse("2006-01-02T00:00:00Z", pago.FechaCobro)
		} else {
			fecha, _ = time.Parse("2006-01-02T00:00:00Z", pago.FechaPago)
		}
		fecha_pago := fecha.Format("2006-01-02")
		if !contains(fechas, fecha_pago) {
			fechas = append(fechas, fecha_pago)
			var cobranza reportedtos.DetallesCobranza

			cobranza.Fecha = fecha_pago
			cobranza.Nombre = (pago.Cliente + "-" + fecha.Format("02-01-2006"))

			entityControl := entities.Reporte{
				Tiporeporte:   "pagos",
				Fechacobranza: fecha.Format("02-01-2006"),
			}
			nroReporteUint, err := s.repository.GetLastReporteEnviadosRepository(entityControl, true)
			if err != nil {
				erro = err
				return
			}
			if nroReporteUint != 0 {
				nroReporteString := strconv.FormatUint(uint64(nroReporteUint), 10)
				cobranza.NroReporte = nroReporteString
			}

			for _, OtroPago := range resultado {
				var fechaOtro time.Time
				//Solo los pagos con rapipago tienen fechacobro, por lo tanto si tiene, debe tomar esa fecha, no la fecha del paid_at
				if OtroPago.FechaCobro != "" {
					fechaOtro, _ = time.Parse("2006-01-02T00:00:00Z", OtroPago.FechaCobro)
				} else {
					fechaOtro, _ = time.Parse("2006-01-02T00:00:00Z", OtroPago.FechaPago)
				}
				//fechaOtro, _ := time.Parse("2006-01-02T00:00:00Z", OtroPago.FechaPago)
				fecha_pagoOtro := fechaOtro.Format("2006-01-02")
				if fecha_pago == fecha_pagoOtro {
					cobranza.Pagos = append(cobranza.Pagos, OtroPago)
					cobranza.TotalIva += OtroPago.Iva
					cobranza.TotalComision += OtroPago.Comision
					cobranza.TotalRetencion += OtroPago.Retencion

					cobranza.Subtotal += OtroPago.TotalPago
					cobranza.Registros += 1
				}
			}

			res.Cobranzas = append(res.Cobranzas, cobranza)
		}
	}

	res.CantidadCobranzas = len(res.Cobranzas)

	sort.Slice(res.Cobranzas, func(i, j int) bool {
		if res.Cobranzas[i].Fecha != "" && res.Cobranzas[j].Fecha != "" {
			fecha1 := s.commons.ConvertirFechaYYYYMMDD(res.Cobranzas[i].Fecha)
			fecha2 := s.commons.ConvertirFechaYYYYMMDD(res.Cobranzas[j].Fecha)
			cambio := fecha1 > fecha2
			return cambio
		}
		return false
	})

	for _, cob := range res.Cobranzas {
		res.Total += cob.Subtotal
	}

	return
}

func (s *reportesService) GetRendicionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseRendicionesClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	filtro := reportedtos.RequestPagosPeriodo{
		ClienteId:     uint64(request.ClienteId),
		CuentaId:      uint64(request.CuentaId),
		FechaInicio:   fechaInicioTime,
		FechaFin:      fechaFinTime,
		OrdenadoFecha: true,
	}

	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}

	var pagosintentos []uint64

	type fechaDeposito struct {
		pagoIntentoId uint64
		fechaDeposito string
	}
	var depositos []fechaDeposito
	var depositosRev []fechaDeposito
	var filtroMov reportedtos.RequestPagosPeriodo
	var pagosintentosrevertidos []uint64
	var movrevertidos []entities.Movimiento
	var controlIds []string

	if len(listaTransferencia) > 0 {
		for _, transferencia := range listaTransferencia {
			transferenciaDepo := fechaDeposito{
				pagoIntentoId: transferencia.Movimiento.PagointentosId,
				fechaDeposito: transferencia.FechaOperacion.Format("02-01-2006"),
			}
			if !transferencia.Reversion {
				pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
				depositos = append(depositos, transferenciaDepo)
			}
			if transferencia.Reversion {
				pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
				depositosRev = append(depositosRev, transferenciaDepo)
			}
		}
		filtroMov = reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentos,
			TipoMovimiento:                  "C",
			CargarComisionImpuesto:          true,
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarCliente:                   true,
			OrdenadoFecha:                   true,
			CargarRetenciones:               true,
		}
	}

	// se obtienen movimientos revertidos, en el caso de que existieran reversiones
	if len(pagosintentosrevertidos) > 0 {
		filtroRevertidos := reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentosrevertidos,
			TipoMovimiento:                  "C",
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarReversionReporte:          true,
			CargarCliente:                   true,
			OrdenadoFecha:                   true,
			CargarComisionImpuesto:          true,
			CargarRetenciones:               true,
		}
		movrevertidos, err = s.repository.GetMovimiento(filtroRevertidos)
		if err != nil {
			erro = err
			return
		}
	}

	var fechas []string

	// se obtienen movimientos a partir de pagointentos positivos, y se recorre cada movimiento
	if len(pagosintentos) > 0 {
		// se obtienen MOVIMIENTOS
		mov, err := s.repository.GetMovimiento(filtroMov)
		if err != nil {
			erro = err
			return
		}

		// En este punto se tienen movimientos C positivos y movimientos revertidos en las var mov y movrevertidos

		for _, m_fecha := range mov {
			// vars para acumular retenciones por gravamen
			var (
				movRetencionGanancias, movRetencionIVA, movRetencionIIBB entities.Monto
			)

			var fecha_deposito string
			// se obtiene el valor de la var fecha_deposito que determina la fecha de cada reporte de rendicion
			for _, deposito := range depositos {
				if m_fecha.Pagointentos.ID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}
			// para no repetir reportes, se controla por la fecha y se va guardando en slice para comparar
			if !contains(fechas, fecha_deposito) {
				// se guarda la feha en slice string fechas
				fechas = append(fechas, fecha_deposito)

				// Var que representa un Reporte de rendicion
				var rendiciones reportedtos.DetallesRendicion

				// cada rendicion tiene una fecha especifica distinta
				rendiciones.Fecha = fecha_deposito
				rendiciones.Nombre = (m_fecha.Cuenta.Cliente.Cliente + "-" + fecha_deposito)

				entityControl := entities.Reporte{
					Tiporeporte:    "rendiciones",
					Fecharendicion: fecha_deposito,
				}
				// teniendo en cuenta la fecha unica se obtiene el numero de reporte correspondiente
				nroReporteUint, err := s.repository.GetLastReporteEnviadosRepository(entityControl, true)
				if err != nil {
					erro = err
					return
				}
				if nroReporteUint != 0 {
					nroReporteString := strconv.FormatUint(uint64(nroReporteUint), 10)
					rendiciones.NroReporte = nroReporteString
				}

				// Cada iteracion de este for equivale a cada operacion que conforman un reporte de cliente
				for _, m := range mov {
					// var para acumular retenciones por movimiento
					var totalRetencionPorMovimiento entities.Monto

					var fecha_deposito_otro string
					for _, deposito := range depositos {
						if m.Pagointentos.ID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {
						// Se calcula la comision e impuesto de cada movimiento
						var comision entities.Monto
						var iva entities.Monto
						if len(m.Movimientocomisions) > 0 {
							comision = m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto + m.Movimientocomisions[len(m.Movimientocomisions)-1].Montoproveedor
							iva = m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto + m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Montoproveedor
						} else {
							comision = 0
							iva = 0
						}

						// Para cada movimiento (m) se suman sus retenciones si las tuviere
						if len(m.Movimientoretencions) > 0 {
							for _, movimiento_retencion := range m.Movimientoretencions {
								totalRetencionPorMovimiento += movimiento_retencion.ImporteRetenido
							}
							movRetencionGanancias += importeRetencionByName("ganancias", m.Movimientoretencions)
							movRetencionIVA += importeRetencionByName("iva", m.Movimientoretencions)
							movRetencionIIBB += importeRetencionByName("iibb", m.Movimientoretencions)
						}

						cantidadBoletas := len(m.Pagointentos.Pago.Pagoitems)
						if !contains(controlIds, fmt.Sprint(m.PagointentosId)) {
							controlIds = append(controlIds, fmt.Sprint(m.PagointentosId))

							rendicion := reportedtos.ResponseReportesRendiciones{
								PagoIntentoId:           m.PagointentosId,
								Cuenta:                  m.Cuenta.Cuenta,
								Id:                      m.Pagointentos.Pago.ExternalReference,
								FechaCobro:              m.Pagointentos.PaidAt.Format("02-01-2006"),
								ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))),
								ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Monto.Float64(), 2))),
								CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
								Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
								Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
								Concepto:                "Transferencia",
								FechaDeposito:           fecha_deposito_otro,
								Retenciones:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRetencionPorMovimiento.Float64(), 2))),
							}

							rendiciones.Rendiciones = append(rendiciones.Rendiciones, rendicion)
							auxCobrado := m.Pagointentos.Amount.Float64()
							auxDepositado := m.Monto.Float64()

							rendiciones.TotalCobrado += auxCobrado
							rendiciones.TotalRendido += auxDepositado
							rendiciones.TotalComision += comision.Float64()
							rendiciones.TotalIva += iva.Float64()
							rendiciones.CantidadOperaciones += 1

						}

					}

				} // Fin de for _, m := range mov. Fin de cada operacion de un reporte
				// Totalizar las retenciones por cada reporte, por tipo de gravamen
				rendiciones.TotalRetGanancias = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionGanancias.Float64(), 2)))
				rendiciones.TotalRetIVA = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionIVA.Float64(), 2)))
				rendiciones.TotalRetIIBB = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionIIBB.Float64(), 2)))
				res.DetallesRendiciones = append(res.DetallesRendiciones, rendiciones)

			}
		} // Fin de for _, m_fecha := range mov. Fin de un reporte

	} // Fin de if len(pagosintentos) > 0

	if len(movrevertidos) > 0 {
		var fechasReversion []string

		for _, mr := range movrevertidos {

			// vars para acumular retenciones por gravamen
			var (
				movRetencionGanancias, movRetencionIVA, movRetencionIIBB entities.Monto
			)

			var fecha_deposito string
			for _, deposito := range depositosRev {
				if mr.Pagointentos.ID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}

			if !contains(fechasReversion, fecha_deposito) {
				fechasReversion = append(fechasReversion, fecha_deposito)

				fechas = append(fechas, fecha_deposito)
				var rendiciones reportedtos.DetallesRendicion
				for i := 0; i < len(res.DetallesRendiciones); i++ {
					if res.DetallesRendiciones[i].Fecha == fecha_deposito {
						rendiciones = res.DetallesRendiciones[i]
					}
				}

				for _, mr2 := range movrevertidos {
					var fecha_deposito_otro string
					for _, deposito := range depositosRev {
						if mr2.Pagointentos.ID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {

						// var para acumular retenciones por movimiento
						var totalRetencionPorMovimiento entities.Monto

						var comision entities.Monto
						var iva entities.Monto
						if len(mr2.Movimientocomisions) > 0 {
							comision = mr2.Movimientocomisions[len(mr2.Movimientocomisions)-1].Monto + mr2.Movimientocomisions[len(mr2.Movimientocomisions)-1].Montoproveedor
							iva = mr2.Movimientoimpuestos[len(mr2.Movimientoimpuestos)-1].Monto + mr2.Movimientoimpuestos[len(mr2.Movimientoimpuestos)-1].Montoproveedor
						} else {
							comision = 0
							iva = 0
						}

						// Para cada movimiento (m) se suman sus retenciones si las tuviere
						if len(mr2.Movimientoretencions) > 0 {
							for _, movimiento_retencion := range mr2.Movimientoretencions {
								totalRetencionPorMovimiento += movimiento_retencion.ImporteRetenido
							}
							movRetencionGanancias += importeRetencionByName("ganancias", mr2.Movimientoretencions)
							movRetencionIVA += importeRetencionByName("iva", mr2.Movimientoretencions)
							movRetencionIIBB += importeRetencionByName("iibb", mr2.Movimientoretencions)
						}

						rendicion := reportedtos.ResponseReportesRendiciones{
							PagoIntentoId:     mr2.PagointentosId,
							Cuenta:            mr2.Cuenta.Cuenta,
							Id:                mr2.Pagointentos.Pago.ExternalReference,
							ImporteDepositado: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr2.Monto.Float64(), 2))),
							Concepto:          "Reversion",
							FechaDeposito:     fecha_deposito_otro,
							Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 2))),
							Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 2))),
							Retenciones:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRetencionPorMovimiento.Float64(), 2))),
						}

						rendiciones.Rendiciones = append(rendiciones.Rendiciones, rendicion)
						auxDepositado := mr2.Monto.Float64()

						rendiciones.TotalRendido += auxDepositado
						rendiciones.CantidadOperaciones += 1
						rendiciones.TotalReversion += auxDepositado
					}

				}

				for i := 0; i < len(res.DetallesRendiciones); i++ {
					if res.DetallesRendiciones[i].Fecha == fecha_deposito {
						res.DetallesRendiciones[i] = rendiciones
					}
				}

			}

		} // Fin de for _, mr := range movrevertidos

		// fmt.Print(depositosRev)

	}

	sort.Slice(res.DetallesRendiciones, func(i, j int) bool {
		if res.DetallesRendiciones[i].Fecha != "" && res.DetallesRendiciones[j].Fecha != "" {
			return s.commons.ConvertirFecha(res.DetallesRendiciones[i].Fecha) > s.commons.ConvertirFecha(res.DetallesRendiciones[j].Fecha)
		}
		return false
	})

	for _, rendicionesFinales := range res.DetallesRendiciones {
		res.CantidadRegistros += 1
		res.Total += rendicionesFinales.TotalRendido
	}

	return
}

func (s *reportesService) GetReversionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseReversionesClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	filtro := reportedtos.RequestPagosPeriodo{
		ClienteId:     uint64(request.ClienteId),
		CuentaId:      uint64(request.CuentaId),
		FechaInicio:   fechaInicioTime,
		FechaFin:      fechaFinTime,
		OrdenadoFecha: true,
	}
	filtro_validacion := reportedtos.ValidacionesFiltro{
		Inicio:  true,
		Fin:     true,
		Cliente: true,
	}

	listaPagos, err := s.repository.GetReversionesReportes(filtro, filtro_validacion)
	if err != nil {
		erro = err
		return
	}

	var pagosRevertidos []uint64
	for _, pagosRevertido := range listaPagos {
		pagosRevertidos = append(pagosRevertidos, uint64(pagosRevertido.PagointentosID))
	}

	type fechaDeposito struct {
		pagoIntentoId uint64
		fechaDeposito string
	}

	var depositosRev []fechaDeposito
	var pagosintentosrevertidos []uint64

	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}
	if len(listaPagos) > 0 {
		if len(listaTransferencia) > 0 {
			for _, transferencia := range listaTransferencia {
				if transferencia.Reversion {
					transferenciaDepo := fechaDeposito{
						pagoIntentoId: transferencia.Movimiento.PagointentosId,
						fechaDeposito: transferencia.FechaOperacion.Format("02-01-2006"),
					}
					pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
					depositosRev = append(depositosRev, transferenciaDepo)
				}
			}
		}

		var fechas []string

		for _, m_fecha := range listaPagos {
			var fecha_deposito string
			for _, deposito := range depositosRev {
				if m_fecha.PagointentosID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}
			if !contains(fechas, fecha_deposito) {
				fechas = append(fechas, fecha_deposito)
				var reversiones reportedtos.DetallesReversiones
				reversiones.Fecha = fecha_deposito
				reversiones.Nombre = (m_fecha.PagoIntento.Pago.PagosTipo.Cuenta.Cliente.Cliente + "-" + fecha_deposito)
				for _, value := range listaPagos {

					var fecha_deposito_otro string
					for _, deposito := range depositosRev {
						if value.PagointentosID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {
						var revertido reportedtos.Reversiones
						var pagoRevertido reportedtos.PagoRevertido
						var itemsRevertido []reportedtos.ItemsRevertidos
						//var itemRevertido reportedtos.ItemsRevertidos
						var intentoPagoRevertido reportedtos.IntentoPagoRevertido
						revertido.EntityToReversiones(value)
						pagoRevertido.EntityToPagoRevertido(value.PagoIntento.Pago)
						if len(value.PagoIntento.Pago.Pagoitems) > 0 {
							for _, valueItem := range value.PagoIntento.Pago.Pagoitems {
								var itemRevertido reportedtos.ItemsRevertidos
								itemRevertido.EntityToItemsRevertidos(valueItem)
								itemsRevertido = append(itemsRevertido, itemRevertido)
							}
						}
						intentoPagoRevertido.EntityToIntentoPagoRevertido(value.PagoIntento)
						pagoRevertido.Items = itemsRevertido
						pagoRevertido.IntentoPago = intentoPagoRevertido
						revertido.PagoRevertido = pagoRevertido
						reversiones.Reversiones = append(reversiones.Reversiones, revertido)

						auxMonto, _ := strconv.ParseFloat(revertido.Monto, 64)
						reversiones.TotalMonto += auxMonto
						reversiones.CantidadOperaciones += 1
					}
				}

				res.DetallesReversiones = append(res.DetallesReversiones, reversiones)
				res.CantidadRegistros += 1
				res.Total += reversiones.TotalMonto
			}
		}
	}
	sort.Slice(res.DetallesReversiones, func(i, j int) bool {
		if res.DetallesReversiones[i].Fecha != "" && res.DetallesReversiones[j].Fecha != "" {
			return s.commons.ConvertirFecha(res.DetallesReversiones[i].Fecha) > s.commons.ConvertirFecha(res.DetallesReversiones[j].Fecha)
		}
		return false
	})
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *reportesService) GetPeticiones(request reportedtos.RequestPeticiones) (response reportedtos.ResponsePeticiones, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	var listaPeticiones []entities.Webservicespeticione
	var total int64
	var err error
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}
	if request.Vendor == "" {
		request.Vendor = "ApiLink"
	}

	if request.Operacion != "" {
		listaPeticiones, total, err = s.repository.GetPeticionesReportesByOperacion(request)
	} else {
		listaPeticiones, total, err = s.repository.GetPeticionesReportes(request)
	}
	if err != nil {
		erro = err
		return
	}

	if len(listaPeticiones) > 0 {
		var resultPeticiones []reportedtos.ResponseDetallePeticion
		for _, peticion := range listaPeticiones {

			resultPeticiones = append(resultPeticiones, reportedtos.ResponseDetallePeticion{
				Operacion: peticion.Operacion,
				Fecha:     peticion.CreatedAt.Format("2006-01-02 15:04:05"),
				Vendor:    peticion.Vendor,
			})
		}

		response = reportedtos.ResponsePeticiones{
			FechaComienzo:   request.FechaInicio.Format("2006-01-02"),
			FechaFin:        request.FechaFin.Format("2006-01-02"),
			TotalPeticiones: int(total),
			Data:            resultPeticiones,
			LastPage:        int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetLogs(request reportedtos.RequestLogs) (response reportedtos.ResponseLogs, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}

	listaLogs, total, err := s.repository.GetLogs(request)
	if err != nil {
		erro = err
		return
	}

	if len(listaLogs) > 0 {
		var resultLogs []reportedtos.ResponseDetalleLog
		for _, log := range listaLogs {

			resultLogs = append(resultLogs, reportedtos.ResponseDetalleLog{
				Mensaje:       log.Mensaje,
				Fecha:         log.CreatedAt.Format("2006-01-02 15:04:05"),
				Funcionalidad: log.Funcionalidad,
				Tipo:          string(log.Tipo),
			})
		}

		response = reportedtos.ResponseLogs{
			FechaComienzo: request.FechaInicio.Format("2006-01-02"),
			FechaFin:      request.FechaFin.Format("2006-01-02"),
			TotalLogs:     int(total),
			Data:          resultLogs,
			LastPage:      int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetNotificaciones(request reportedtos.RequestNotificaciones) (response reportedtos.ResponseNotificaciones, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}

	listaNotif, total, err := s.repository.GetNotificaciones(request)
	if err != nil {
		erro = err
		return
	}

	if len(listaNotif) > 0 {
		var resultNotifs []reportedtos.ResponseDetalleNotificaciones
		for _, notif := range listaNotif {

			resultNotifs = append(resultNotifs, reportedtos.ResponseDetalleNotificaciones{
				Descripcion: notif.Descripcion,
				Fecha:       notif.CreatedAt.Format("2006-01-02 15:04:05"),
				Tipo:        string(notif.Tipo),
			})
		}

		response = reportedtos.ResponseNotificaciones{
			FechaComienzo:       request.FechaInicio.Format("2006-01-02"),
			FechaFin:            request.FechaFin.Format("2006-01-02"),
			TotalNotificaciones: int(total),
			Data:                resultNotifs,
			LastPage:            int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetReportesEnviadosService(request reportedtos.RequestReportesEnviados) (response reportedtos.ResponseReportesEnviados, erro error) {

	// se recibe respuesta del repositorio con datos de base de datos. o un error
	listaReportes, totalFilas, erro := s.repository.GetReportesEnviadosRepository(request)

	if erro != nil {
		erro = errors.New(erro.Error())
		return
	}

	// pasar las entidades a DTO response correspondiente

	var resTemporal reportedtos.ResponseReporteEnviado

	for _, reporte := range listaReportes {

		resTemporal.EntityToDto(reporte)

		response.Reportes = append(response.Reportes, resTemporal)
	}

	// paginacion
	if request.Number > 0 && request.Size > 0 {
		response.Meta = setPaginacion(uint32(request.Number), uint32(request.Size), totalFilas)
	}

	return
}

func (s *reportesService) EnumerarReportesEnviadosService(request reportedtos.RequestReportesEnviados) (erro error) {

	tiposReportes := []string{"pagos", "rendiciones"}

	for _, tipoReporte := range tiposReportes {

		request.TipoReporte = reportedtos.EnumTipoReporte(tipoReporte)

		listaReportes, _, err := s.repository.GetReportesEnviadosRepository(request)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}

		var cambios []entities.Reporte

		contador := 1
		for i, reporte := range listaReportes {
			if reporte.Tiporeporte == tipoReporte {
				if listaReportes[i].Nro_reporte != uint(contador) {
					listaReportes[i].Nro_reporte = uint(contador)
					cambios = append(cambios, listaReportes[i])
				}

				contador += 1
			}
		}

		if len(cambios) > 0 {
			erro = s.repository.GuardarReportesInfo(cambios)

			if erro != nil {
				erro = errors.New(erro.Error())
				return
			}

		}

	}

	return
}

func (s *reportesService) CopiarNumeroReporteOriginal() (erro error) {

	var request reportedtos.RequestReportesEnviados
	tiposReportes := []string{"pagos", "rendiciones"}

	for _, tipoReporte := range tiposReportes {

		request.TipoReporte = reportedtos.EnumTipoReporte(tipoReporte)
		request.SinNumero = true

		listaReportes, _, err := s.repository.GetReportesEnviadosRepository(request)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}

		for i, reporte := range listaReportes {
			if reporte.Tiporeporte == tipoReporte {
				nro_OGReporte, err := s.repository.GetLastReporteEnviadosRepository(reporte, true)
				if err != nil {
					erro = errors.New(err.Error())
					return
				}
				listaReportes[i].Nro_reporte = nro_OGReporte
			}
		}

		erro = s.repository.GuardarReportesInfo(listaReportes)

		if erro != nil {
			erro = errors.New(erro.Error())
			return
		}
	}

	return
}

func setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))

	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	return

}

// para un slice de entities.MovimientoRetencion devuelve el importe de la retencion por el nombre del gravamen correspondiente
func importeRetencionByName(gravamen_name string, mr []entities.MovimientoRetencion) (importe entities.Monto) {
	for _, item := range mr {
		if gravamen_name == item.Retencion.Condicion.Gravamen.Gravamen {
			importe = item.ImporteRetenido
			break
		}
	}

	return
}

// obtener reporte de rendicion de un slice de reportes a partir de la fecharendicion
func ObtenerReportePorFecha(fecharendicion string, reportes []entities.Reporte) (reporte entities.Reporte) {
	for _, item := range reportes {
		if fecharendicion == item.Fecharendicion {
			reporte = item
			break
		}
	}

	return
}
