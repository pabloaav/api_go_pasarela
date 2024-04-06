package administracion

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

func (s *service) BuildCierreLoteApiLinkService() (res administraciondtos.RegistroClPagosApilink, erro error) {

	var listaCierretemporal []linkdebin.TemporalDebines
	// var listaPagosTemporal []cierrelotedtos.ResponsePagosApilink
	// mediante estas variables permite que debines aun no fueron registrados
	// solo se registran los debines que pasan a un estado final
	var debinId []string
	var debinTabla []string

	// 1 - Se dicidió que el estado para el debin iniciado y en curso por apilink es Processing
	processing, erro := s.utilService.FirstOrCreateConfiguracionService("DEBIN_PROCESSING", "Nombre del estado normal (iniciado/en curso) para el debin", "Processing")
	if erro != nil {
		return
	}

	filtroPagosEstado := filtros.PagoEstadoFiltro{
		BuscarPorFinal: true,
		Final:          false,
		Nombre:         processing,
	}
	pagoEstado, erro := s.repository.GetPagoEstado(filtroPagosEstado)

	if erro != nil {
		return
	}

	//Si es menor que uno el estado processing no esta presente y no se puede proseguir
	if pagoEstado.ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildCierreLoteApiLinkService",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Info(ERROR_LOG + "BuildCierreLoteApiLinkService." + erro.Error())
		}
		return
	}

	// 2 - Busco el channel debin para garantizar que obtengo solo pagos intentos de debin
	canalDebin, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_DEBIN", "Nombre del canal debin", "debin")
	if erro != nil {
		return
	}
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalDebin},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)

	if erro != nil {
		return
	}
	if channel.ID < 1 {
		erro = fmt.Errorf(ERROR_CHANNEL_ID)
		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildCierreLoteApiLinkService",
			Mensaje:       ERROR_CHANNEL_ID,
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Info(ERROR_LOG + "BuildCierreLoteApiLinkService." + erro.Error())
		}
		return
	}

	// 3 - Busco todos los pagos que tengan el estado Processing
	filtroPagos := filtros.PagoFiltro{
		PagoEstadosId:     uint64(pagoEstado.ID),
		CargaPagoIntentos: true,
		CargaMedioPagos:   true,
		CargarCuenta:      true,
	}
	pagosPendientes, _, erro := s.repository.GetPagos(filtroPagos)
	if erro != nil {
		return
	}
	//Si no hay ningún pago en estado processing es porque no hay debines entonces no hacemos nada.
	if len(pagosPendientes) < 1 {
		return
	}

	// 4 - verifico si los pagos seleccionados tiene el canal de debin eso lo hago para certificarme que solo se procesen debines
	var pagosPendientesDebin []entities.Pago
	for i := range pagosPendientes {
		if len(pagosPendientes[i].PagoIntentos) > 0 {
			// NOTE Refactorizacion solo consultar por el ultimo pago intento del pago y que la fecha de pago no sea cero
			if pagosPendientes[i].PagoIntentos[len(pagosPendientes[i].PagoIntentos)-1].Mediopagos.ChannelsID == int64(channel.ID) && !pagosPendientes[i].PagoIntentos[len(pagosPendientes[i].PagoIntentos)-1].PaidAt.IsZero() {
				pagosPendientesDebin = append(pagosPendientesDebin, pagosPendientes[i])
			}

		}
	}

	//Si la lista es vacia significa que no existen debines pendientes
	if len(pagosPendientesDebin) < 1 {
		return
	}

	// 5 - Ordena los pagos pendientes por fecha de creación // REVIEW esto no es necesario porque ya están ordenados por fecha paid_at : verificar de nuevo
	// hago eso para saber la fecha inicial y final para hacer la consulta a ApiLink
	// sort.Slice(pagosPendientesDebin, func(i, j int) bool {
	// 	return pagosPendientesDebin[i].CreatedAt.Before(pagosPendientesDebin[j].CreatedAt)
	// })

	// if erro != nil {
	// 	return
	// }

	// verificar en el primer pago el intento exitoso
	var fechaDes time.Time
	for _, paintento := range pagosPendientesDebin[0].PagoIntentos {
		if !paintento.PaidAt.IsZero() {
			fechaDes = paintento.PaidAt
			break
		}
	}
	/* NOTE  formato para consultar debines por periodo de fechas*/
	uuid := s.commonsService.NewUUID()
	request := linkdebin.RequestGetDebinesLink{
		Pagina:      1,
		Tamanio:     linkdtos.Cien,
		Cbu:         config.CBU_CUENTA_TELCO,
		EsComprador: false,
		// FechaDesde:  pagosPendientesDebin[0].PagoIntentos[0].PaidAt,                                                                                                                    //La fecha y hora del primer pago
		// FechaHasta:  pagosPendientesDebin[len(pagosPendientesDebin)-1].PagoIntentos[len(pagosPendientesDebin[len(pagosPendientesDebin)-1].PagoIntentos)-1].PaidAt.Add(time.Minute * 1), //La fecha y hora del ultimo pago mas un minuto
		FechaDesde: fechaDes,                                                                                                                                                          //La fecha y hora del primer pago
		FechaHasta: pagosPendientesDebin[len(pagosPendientesDebin)-1].PagoIntentos[len(pagosPendientesDebin[len(pagosPendientesDebin)-1].PagoIntentos)-1].PaidAt.Add(time.Minute * 1), //La fecha y hora del ultimo pago mas un minuto
		Tipo:       linkdtos.DebinDefault,
		Estado:     "",
	}
	// NOTE se consulta a apilink los debines de este periodo de fecha
	logs.Info(fmt.Sprint("se consulta a apilink desde la fecha: ", request.FechaDesde))
	logs.Info(fmt.Sprint("se consulta a apilink hasta la fecha: ", request.FechaHasta))

	response, erro := s.apilinkService.GetDebinesApiLinkService(uuid, request)
	if erro != nil || len(response.Debines) < 1 {
		return
	}
	// logs.Info(response)
	// en el caso de que la respuesta tenga mas de una pagina , se debe hacer una consulta por cada pagina
	if len(response.Debines) > 100 {
		pagina := response.Paginado.Pagina
		cantPaginas := response.Paginado.CantidadPaginas
		for i := 0; pagina != cantPaginas; i++ {
			pagina++
			request.Pagina = pagina
			debines, err := s.apilinkService.GetDebinesApiLinkService(uuid, request)
			erro = err
			if erro != nil || len(response.Debines) < 1 {
				return
			}
			response.Debines = append(response.Debines, debines.Debines...)
		}
	}

	pagoDebin, erro := s.utilService.FirstOrCreateConfiguracionService("DEBIN_PAGO_EXTERNO", "Se usa para indicar cual es el vendor de pagos externos para apilink", "APILINK")
	if erro != nil {
		return
	}

	// 7 - Consulta los estados externos para poder vincular con los estados de nuestra api
	filtroPagoEstadoExternos := filtros.PagoEstadoExternoFiltro{
		Vendor:           pagoDebin,
		CargarEstadosInt: true,
	}

	/*OBTIENE PAGOS EXTERNOS DEBINES*/
	pagosEstadosExternos, erro := s.repository.GetPagosEstadosExternos(filtroPagoEstadoExternos)

	if erro != nil {
		return
	}

	// obtener de la lista de pagosexternos solo los que son estado final
	// me permite filtrar solo los debines que estan en un estado final
	var idestadosfinales []uint64
	for _, estado := range pagosEstadosExternos {
		if estado.PagoEstados.Final {
			idestadosfinales = append(idestadosfinales, uint64(estado.ID))
		}
	}

	if len(pagosEstadosExternos) < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_EXTERNO_LISTA)
		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildCierreLoteApiLinkService",
			Mensaje:       ERROR_PAGO_ESTADO_EXTERNO_LISTA,
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Info(ERROR_LOG + "BuildCierreLoteApiLinkService." + erro.Error())
		}
		return
	}
	// NOTE control respuesta apilink con los pagos intentos:
	// situacion : en las respuestas llegaban debines que no estaban registrados en la base de datos , el la generacion de movimientos esto de controla
	/* se agrega campo referencia_banco: este campo nos permite comparar y conciliar con los movimientos del banco  */
	for j := range response.Debines {
		if len(response.Debines[j].Tipo) < 1 {
			response.Debines[j].Tipo = "DEBIN"
		}
		//NOTE el control se realiza sobre los pagos intentos de tipo debin
		for l := range pagosPendientesDebin {
			if response.Debines[j].Id == pagosPendientesDebin[l].PagoIntentos[len(pagosPendientesDebin[l].PagoIntentos)-1].ExternalID {
				// listaPagos = append(listaPagos, pagosPendientesDebin[l].PagoIntentos[len(pagosPendientesDebin[l].PagoIntentos)-1].Pago)
				cierre := entities.Apilinkcierrelote{
					Uuid:     uuid,
					DebinId:  response.Debines[j].Id,
					Concepto: response.Debines[j].Concepto,
					Moneda:   response.Debines[j].Moneda,
					Importe:  entities.Monto(response.Debines[j].Importe),
					Estado:   response.Debines[j].Estado,

					Tipo:            response.Debines[j].Tipo,
					FechaExpiracion: response.Debines[j].FechaExpiracion,
					Devuelto:        response.Debines[j].Devuelto,
					ContracargoId:   response.Debines[j].ContraCargoId,
					CompradorCuit:   response.Debines[j].Comprador.Cuit,
					VendedorCuit:    response.Debines[j].Vendedor.Cuit,
					ReferenciaBanco: commons.Concat(response.Debines[j].Id, response.Debines[j].Comprador.Cuit),
				}

				/* controla los debines consultados por apilink  y la referencia externa de nuestra base de datos*/
				for i := range pagosEstadosExternos {
					if linkdtos.EnumEstadoDebin(strings.ToUpper(pagosEstadosExternos[i].Estado)) == response.Debines[j].Estado {
						cierre.PagoestadoexternosId = uint64(pagosEstadosExternos[i].ID)
						cierre.Pagoestadoexterno = pagosEstadosExternos[i]
						break
					}
				}
				if cierre.PagoestadoexternosId == 0 {
					logs.Info("Error en estado obtenido el apilink el valor es" + response.Debines[j].Estado)
					erro = fmt.Errorf(ERROR_ESTADO_PAGO_NO_ENCONTRADO)
					return
				}

				for _, estado := range idestadosfinales {
					if cierre.PagoestadoexternosId == estado {
						pagos := cierrelotedtos.ResponsePagosApilink{
							Id:         pagosPendientesDebin[l].ID,
							Estadopago: cierre.Pagoestadoexterno.PagoestadosId,
						}
						listaCierretemporal = append(listaCierretemporal, linkdebin.TemporalDebines{
							Debines: &cierre,
							Pagos:   pagos,
						})
						debinId = append(debinId, cierre.DebinId)
						break
					}
				}

				// if cierre.PagoestadoexternosId > 0 && cierre.Pagoestadoexterno.PagoEstados.Final {
				// 	listaCierretemporal = append(listaCierretemporal, &cierre)
				// 	debinId = append(debinId, cierre.DebinId)
				// } else {
				// 	logs.Info("Error en estado obtenido el apilink el valor es" + response.Debines[j].Estado)
				// 	erro = fmt.Errorf(ERROR_ESTADO_PAGO_NO_ENCONTRADO)
				// 	return
				// }
			}

		}

	}

	// NOTE obtener debines registrados (se obtienen tambien los eliminados)
	if len(debinId) > 0 {
		filtro := linkdebin.RequestDebines{
			Debines:         debinId,
			BancoExternalId: false,
		}
		getdebines, err := s.GetConsultarDebines(filtro)
		if err != nil {
			err = erro
			return
		}

		/* comparar con la lista de debines eliminados y lista de debines obtenidas de apilink*/
		// FIXME se deben comparar con los existen y descartar los repetidos
		if len(getdebines) > 0 {
			for _, debin := range getdebines {
				debinTabla = append(debinTabla, debin.DebinId)
			}
		}
		// NOTE se guardaran solo los debines que no estan en db
		difdebines := commons.Difference(debinId, debinTabla)

		for _, lista := range listaCierretemporal {
			for _, dif := range difdebines {
				if lista.Debines.DebinId == dif {
					res.ListaCLApiLink = append(res.ListaCLApiLink, lista.Debines)
					res.ListaPagos = append(res.ListaPagos, lista.Pagos)
				}
			}
		}
	}

	return
}

// crear registros en tabla apilink y actualizar estado de los pagos(solo tiene en cuenta los pagos en estado final)
func (s *service) CreateCLApilinkPagosService(ctx context.Context, mcl administraciondtos.RegistroClPagosApilink) (erro error) {
	erro = s.repository.CreateCLApilinkPagosRepository(ctx, mcl)
	return
}

func (s *service) CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error) {
	erro = s.repository.CreateCierreLoteApiLink(cierreLotes)
	return
}

func (s *service) GetDebines(request linkdebin.RequestDebines) (response []*entities.Apilinkcierrelote, erro error) {
	response, err := s.repository.GetConsultarDebines(request)
	if err != nil {
		erro = err
		return nil, erro
	}
	return
}

func (s *service) GetConsultarDebines(request linkdebin.RequestDebines) (response []linkdebin.ResponseDebinesEliminados, erro error) {
	debines, err := s.repository.GetConsultarDebines(request)
	if err != nil {
		erro = err
		return nil, erro
	}
	for _, deb := range debines {
		response = append(response, linkdebin.ResponseDebinesEliminados{
			Id:              uint64(deb.ID),
			DebinId:         deb.DebinId,
			Estado:          linkdtos.EnumEstadoDebin(deb.Estado),
			Match:           deb.Match,
			BancoExternalId: deb.BancoExternalId,
		})
	}

	responsedebines := unique(response)
	return responsedebines, nil
}

func unique(arr []linkdebin.ResponseDebinesEliminados) []linkdebin.ResponseDebinesEliminados {
	keys := make(map[string]bool)
	list := []linkdebin.ResponseDebinesEliminados{}
	for _, entry := range arr {
		if _, value := keys[entry.DebinId]; !value {
			keys[entry.DebinId] = true
			list = append(list, entry)
		}
	}
	return list
}

func (s *service) BuildNotificacionPagosService(request webhook.RequestWebhook) (listaPagos []entities.Pagotipo, erro error) {

	// 1 obtener estado de pagos finales son tanto para debin, offline y credito/debito
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		BuscarPorFinal: request.EstadoFinalPagos,
		Final:          request.EstadoFinalPagos,
	}

	pagoEstados, erro := s.repository.GetPagosEstados(filtroPagosEstado)
	if erro != nil {
		return
	}
	if len(pagoEstados) == 0 {
		erro = fmt.Errorf(ERROR_ESTADO_PAGO_NO_ENCONTRADO)

		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildNotificacionPagosService",
			Mensaje:       ERROR_ESTADO_PAGO_NO_ENCONTRADO,
		}
		err := s.utilService.CreateLogService(log)

		if err != nil {
			logs.Info(ERROR_LOG + "BuildNotificacionPagosService." + erro.Error())
		}

		return
	}

	// obtener estado pendiente
	filtro := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.repository.GetPagosEstados(filtro)
	if err != nil {
		erro = err
		return
	}
	if len(estadoPendiente) == 0 {
		erro = fmt.Errorf(ERROR_ESTADO_PAGO_NO_ENCONTRADO)

		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildNotificacionPagosService",
			Mensaje:       ERROR_ESTADO_PAGO_NO_ENCONTRADO,
		}
		err := s.utilService.CreateLogService(log)

		if err != nil {
			logs.Info(ERROR_LOG + "BuildNotificacionPagosService." + erro.Error())
		}

		return
	}

	// 3 filtrar estados menos el pendiente
	var pagosEstados []uint64
	for _, pagoEstado := range pagoEstados {
		if pagoEstado.ID != estadoPendiente[0].ID {
			pagosEstados = append(pagosEstados, uint64(pagoEstado.ID))
		}
	}

	//filtro indicar la cantidad de dias de pagos por notificar
	filtroPagos := filtros.PagoTipoFiltro{
		CargarPagos:           true,
		CargarPagosNotificado: request.PagosNotificado,
		IdCuenta:              int64(request.CuentaId),
		PagoEstadosIds:        pagosEstados,
		FechaPagoInicio:       time.Now().AddDate(0, 0, int(request.DiasPago*-1)),
		FechaPagoFin:          time.Now(),
	}

	// 4 obtener los pagos de los ultimos dias indicado en el filtro con estado procesando y pagado
	pagosNotificacion, _, erro := s.repository.GetPagosTipo(filtroPagos)

	if erro != nil {
		return
	}

	return pagosNotificacion, nil
}

func (s *service) BuildNotificacionPagosCLRapipago(filtro filtros.PagoEstadoFiltro) (response []webhook.WebhookResponse, erro error) {

	pagoEstadoAprobado, erro := s.repository.GetPagoEstado(filtro)

	if erro != nil {
		return
	}

	// & obtener canal offline
	canalOffline, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_OFFLINE", "Nombre del canal de debin", "offline")

	if erro != nil {
		return
	}
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalOffline},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)
	if erro != nil {
		return
	}

	//filtro indicar la cantidad de dias de pagos por notificar
	filtroPagos := filtros.PagoTipoFiltro{
		CargarPagos:     true,
		PagoEstadosIds:  []uint64{uint64(pagoEstadoAprobado.ID)},
		FechaPagoInicio: time.Now().AddDate(0, 0, -7),
		FechaPagoFin:    time.Now(),
	}

	// 4 obtener los pagos de los ultimos dias indicado en el filtro con estado procesando y pagado
	pagosNotificacion, _, erro := s.repository.GetPagosTipo(filtroPagos)

	for _, pagoTipo := range pagosNotificacion {
		if pagoTipo.BackUrlNotificacionPagos != "" && len(pagoTipo.Pagos) > 0 {
			url := pagoTipo.BackUrlNotificacionPagos
			var pagos []webhook.ResultadoResponseWebHook
			for _, pago := range pagoTipo.Pagos {
				if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channel.ID) {
					var importePagado entities.Monto
					last := pago.PagoIntentos[len(pago.PagoIntentos)-1]
					importePagado = entities.Monto(last.Amount)
					pagos = append(pagos, webhook.ResultadoResponseWebHook{
						Id:                int64(pago.ID),
						EstadoPago:        pago.PagoEstados.Nombre,
						Exito:             true,
						Uuid:              pago.Uuid,
						Channel:           last.Mediopagos.Channel.Nombre,
						Description:       pago.Description,
						FirstDueDate:      pago.FirstDueDate,
						FirstTotal:        pago.FirstTotal,
						SecondDueDate:     pago.SecondDueDate,
						SecondTotal:       pago.SecondTotal,
						PayerName:         pago.PayerName,
						PayerEmail:        pago.PayerEmail,
						ExternalReference: pago.ExternalReference,
						Metadata:          pago.Metadata,
						PdfUrl:            pago.PdfUrl,
						CreatedAt:         pago.CreatedAt,
						ImportePagado:     importePagado.Float64(),
					})
				}

			}
			if len(pagos) > 0 {
				response = append(response, webhook.WebhookResponse{
					Url:                      url,
					ResultadoResponseWebHook: pagos,
				})
			}
		}
	}

	if erro != nil {
		return
	}

	return
}

func (s *service) BuildNotificacionPagosCLApilink(request []linkdebin.ResponseDebinesEliminados) (response []webhook.WebhookResponse, debin []uint64, erro error) {

	estados := filtros.PagoEstadoFiltro{
		BuscarPorFinal: true,
		Final:          true,
	}
	pagoEstadofinales, erro := s.repository.GetPagosEstados(estados)
	if erro != nil {
		return
	}
	if len(pagoEstadofinales) == 0 {
		erro = fmt.Errorf(ERROR_ESTADO_PAGOS_NO_ENCONTRADO)

		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildNotificacionPagosService",
			Mensaje:       ERROR_ESTADO_PAGO_NO_ENCONTRADO,
		}
		err := s.utilService.CreateLogService(log)

		if err != nil {
			logs.Info(ERROR_LOG + "BuildNotificacionPagosService." + erro.Error())
		}

		return
	}

	var estfinal []uint64
	for _, estadosFinales := range pagoEstadofinales {
		estfinal = append(estfinal, uint64(estadosFinales.ID))
	}

	// 2 - Busco el channel debin para garantizar que obtengo solo pagos intentos de debin
	canalDebin, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_DEBIN", "Nombre del canal debin", "debin")
	if erro != nil {
		return
	}
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalDebin},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)

	if erro != nil {
		return
	}
	if channel.ID < 1 {
		erro = fmt.Errorf(ERROR_CHANNEL_ID)
		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildCierreLoteApiLinkService",
			Mensaje:       ERROR_CHANNEL_ID,
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Info(ERROR_LOG + "BuildCierreLoteApiLinkService." + erro.Error())
		}
		return
	}
	//filtro indicar la cantidad de dias de pagos por notificar
	filtroPagos := filtros.PagoTipoFiltro{
		CargarPagos:     true,
		PagoEstadosIds:  estfinal,
		FechaPagoInicio: time.Now().AddDate(0, 0, -5),
		FechaPagoFin:    time.Now(),
	}
	// 4 obtener los pagos de los ultimos dias indicado en el filtro con estado procesando y pagado
	pagosNotificacion, _, erro := s.repository.GetPagosTipo(filtroPagos)
	for _, pagoTipo := range pagosNotificacion {
		if pagoTipo.BackUrlNotificacionPagos != "" && len(pagoTipo.Pagos) > 0 {
			url := pagoTipo.BackUrlNotificacionPagos
			var pagos []webhook.ResultadoResponseWebHook
			for _, pago := range pagoTipo.Pagos {
				if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channel.ID) {
					// NOTE recorrer debines y comparar con pagos intentos
					for _, deb := range request {
						last := pago.PagoIntentos[len(pago.PagoIntentos)-1]
						var importePagado entities.Monto
						if deb.DebinId == last.ExternalID {
							debin = append(debin, deb.Id)
							importePagado = entities.Monto(last.Amount)
							pagos = append(pagos, webhook.ResultadoResponseWebHook{
								Id:                int64(pago.ID),
								EstadoPago:        pago.PagoEstados.Nombre,
								Exito:             true,
								Uuid:              pago.Uuid,
								Channel:           last.Mediopagos.Channel.Nombre,
								Description:       pago.Description,
								FirstDueDate:      pago.FirstDueDate,
								FirstTotal:        pago.FirstTotal,
								SecondDueDate:     pago.SecondDueDate,
								SecondTotal:       pago.SecondTotal,
								PayerName:         pago.PayerName,
								PayerEmail:        pago.PayerEmail,
								ExternalReference: pago.ExternalReference,
								Metadata:          pago.Metadata,
								PdfUrl:            pago.PdfUrl,
								CreatedAt:         pago.CreatedAt,
								ImportePagado:     importePagado.Float64(),
							})
						}
					}
				}
			}
			if len(pagos) > 0 {
				response = append(response, webhook.WebhookResponse{
					Url:                      url,
					ResultadoResponseWebHook: pagos,
				})
			}
		}
	}

	if erro != nil {
		return
	}

	return
}

func (s *service) CreateNotificacionPagosService(listaPagos []entities.Pagotipo) (response []webhook.WebhookResponse, erro error) {

	if len(listaPagos) > 0 {
		for _, pago := range listaPagos {
			if pago.BackUrlNotificacionPagos != "" && len(pago.Pagos) > 0 {
				url := pago.BackUrlNotificacionPagos
				var pagos []webhook.ResultadoResponseWebHook
				for _, pago := range pago.Pagos {
					// pagosupdate = append(pagosupdate, pago.ID)
					/* obtener el ultimo pago intento */
					var importePagado entities.Monto
					last := pago.PagoIntentos[len(pago.PagoIntentos)-1]
					importePagado = entities.Monto(last.Amount)
					pagos = append(pagos, webhook.ResultadoResponseWebHook{
						Id:                int64(pago.ID),
						EstadoPago:        pago.PagoEstados.Nombre,
						Exito:             true,
						Uuid:              pago.Uuid,
						Channel:           last.Mediopagos.Channel.Nombre,
						Description:       pago.Description,
						FirstDueDate:      pago.FirstDueDate,
						FirstTotal:        pago.FirstTotal,
						SecondDueDate:     pago.SecondDueDate,
						SecondTotal:       pago.SecondTotal,
						PayerName:         pago.PayerName,
						PayerEmail:        pago.PayerEmail,
						ExternalReference: pago.ExternalReference,
						Metadata:          pago.Metadata,
						PdfUrl:            pago.PdfUrl,
						CreatedAt:         pago.CreatedAt,
						ImportePagado:     importePagado.Float64(),
					})
				}
				response = append(response, webhook.WebhookResponse{
					Url:                      url,
					ResultadoResponseWebHook: pagos,
				})
			}
		}

	}
	return
}

func (s *service) NotificarPagos(listaPagos []webhook.WebhookResponse) (pagoupdate []uint) {
	// var pagosupdate []uint
	for _, webhook := range listaPagos {
		erro := s.webhook.NotificarPagos(webhook)
		if erro != nil {
			logs.Info(erro) //solo informar el error continuar enviando los pagos a los demas clientes
			log := entities.Log{
				Tipo:          entities.Error,
				Funcionalidad: "NotificarPagos",
				Mensaje:       fmt.Sprintf("webhook: no se pudo notificar pagos al cliente .: %s%s", erro, webhook.Url),
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				logs.Info(ERROR_LOG + "NotificarPagos." + erro.Error())
			}
		} else {
			logs.Info(fmt.Sprintf("webhook: se notifico con exito al cliente:%s", webhook.Url))
			for _, pago := range webhook.ResultadoResponseWebHook {
				pagoupdate = append(pagoupdate, uint(pago.Id))
			}
			logs.Info(fmt.Sprintf("webhook: se notifico con exito al cliente los siguientes pagos:%v", pagoupdate))

			// crear logs
			log := entities.Log{
				Tipo:          entities.Info,
				Funcionalidad: "NotificarPagos",
				Mensaje:       fmt.Sprintf("webhook: se notifico con exito al cliente:%s", webhook.Url),
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				logs.Info(ERROR_LOG + "NotificarPagos." + erro.Error())
			}
		}
	}
	return
}

func (s *service) UpdatePagosNoticados(listaPagosNotificar []uint) (erro error) {
	return s.repository.UpdatePagosNotificados(listaPagosNotificar)
}

func (s *service) UpdateCierreLoteApilink(request linkdebin.RequestListaUpdateDebines) (erro error) {
	return s.repository.UpdateCierreloteApilink(request)
}

func (s *service) BuildMovimientoApiLink(listaCierre []*entities.Apilinkcierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

	// 1 - listadeCierre es el resultado de consultar los debines iniciados en api link
	if len(listaCierre) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	var listaCierreDebinId []string

	for i := range listaCierre {
		listaCierreDebinId = append(listaCierreDebinId, listaCierre[i].DebinId)
	}

	canalDebin, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_DEBIN", "Nombre del canal de debin", "debin")

	if erro != nil {
		return
	}

	//Busco el canal debin para filtrar
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalDebin},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)

	if erro != nil {
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		ExternalIds:          listaCierreDebinId,
		Channel:              true,
		CargarPago:           true,
		CargarPagoTipo:       true,
		CargarCuenta:         true,
		CargarCliente:        true,
		CargarCuentaComision: true,
		CargarImpuestos:      true,
	}
	// 1 - Busco los pagos intentos que corresponden a los externalids
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)

	if erro != nil {
		return
	}

	// Verifico si el canal del pago es debin porque puede
	// ocurrir de haber un otro pago con el mismo numero de externalids
	for i := range pagosIntentos {
		if pagosIntentos[i].Mediopagos.ChannelsID == int64(channel.ID) {
			movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, pagosIntentos[i])
		}
	}

	// 2 - Busco el estado acreditado
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, erro := s.repository.GetPagoEstado(filtroPagoEstado)

	if erro != nil {
		return
	}

	// Se debe encontrar un pagointento para cada debin del cierre
	/* si la lista de pago intentos es distinto a la lista de cierre de lote de apilink*/
	if len(listaCierre) != len(movimientoCierreLote.ListaPagoIntentos) {
		var listaPagosExternalId []string

		for i := range movimientoCierreLote.ListaPagoIntentos {
			listaPagosExternalId = append(listaPagosExternalId, movimientoCierreLote.ListaPagoIntentos[i].ExternalID)
		}
		mensaje := fmt.Errorf("no se encontrarion los siguientes debines %+v", commons.Difference(listaCierreDebinId, listaPagosExternalId)).Error()

		erro = fmt.Errorf(ERROR_CIERRE_PAGO_INTENTO)

		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildMovimientoApiLink",
			Mensaje:       mensaje,
		}

		err := s.utilService.CreateLogService(log)

		if err != nil {
			logs.Info(ERROR_LOG + "BuildMovimientoApiLink." + erro.Error())
		}

		return
	}
	// 3 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range movimientoCierreLote.ListaPagoIntentos {
		/* NOTE para el calculo de la comision fitrar por el id del channel y el id de la cuenta para debin funciona para prisma controlar*/
		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         channel.ID,
			CuentaId:          movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.ID,
			FechaPagoVigencia: movimientoCierreLote.ListaPagoIntentos[i].PaidAt,
			Channelarancel:    true,
		}

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			logs.Info(fmt.Sprintf("no se pudo encontrar una comision para la cuenta %s del cliente %s", movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente))
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}
		for j := range listaCierre {
			/* control de que el debin del cierre de lote sea igual al debin del pago intento*/
			if movimientoCierreLote.ListaPagoIntentos[i].ExternalID == listaCierre[j].DebinId {
				if movimientoCierreLote.ListaPagoIntentos[i].Pago.PagoestadosID == int64(listaCierre[j].Pagoestadoexterno.PagoestadosId) {
					// FIXME Hay que ver como controlar ese error si hay que abortar
					if movimientoCierreLote.ListaPagoIntentos[i].Amount != listaCierre[j].Importe {
						erro = fmt.Errorf("el monto informado no es valido")
						return
					}

					// movimientoCierreLote.ListaPagoIntentos[i].Pago.PagoestadosID = int64(listaCierre[j].Pagoestadoexterno.PagoestadosId)
					// movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, movimientoCierreLote.ListaPagoIntentos[i].Pago)

					pagoEstadoLog := entities.Pagoestadologs{
						PagosID:       movimientoCierreLote.ListaPagoIntentos[i].PagosID,
						PagoestadosID: int64(listaCierre[j].Pagoestadoexterno.PagoestadosId),
					}
					movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)

					// verificar que el pago sea acreditado , el campo match en 1 y ademas tenga el id del banco : con esto me aseguro que el movimiento se encuentra en el banco
					if (listaCierre[j].Pagoestadoexterno.PagoestadosId == uint64(pagoEstadoAcreditado.ID)) && (listaCierre[j].Match == 1 || listaCierre[j].BancoExternalId != 0) {
						movimiento := entities.Movimiento{}
						movimiento.AddCredito(uint64(movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.CuentasID), uint64(movimientoCierreLote.ListaPagoIntentos[i].ID), listaCierre[j].Importe)

						s.utilService.BuildComisiones(&movimiento, &listaCuentaComision, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, listaCierre[j].Importe)

						movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
						movimientoCierreLote.ListaPagoIntentos[i].AvailableAt = listaCierre[j].CreatedAt
						// si se concilia y generan movimientos se dan de bajas los debines de la tabla apilinkcierrelote
						movimientoCierreLote.ListaCLApiLink = append(movimientoCierreLote.ListaCLApiLink, *listaCierre[j])
					}
				}
				break
			}
		}
	}
	return
}

/*
Autor: Jose Alarcon
Descripción: se debe verificar que los movimientoe que se envian para transfererir correspodan con la fecha de retiro automatico configurada en la cuenta
*/
func (s *service) BuildTransferenciaCliente(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenicaCliente, cuentaId uint64) (response linktransferencia.ResponseTransferenciaCreateLink, err error) {

	filtroConfiguracionCuenta := filtros.ConfiguracionFiltro{
		Nombre: "CBU_CUENTA_TELCO",
	}

	cuentaTelco, err := s.utilService.GetConfiguracionService(filtroConfiguracionCuenta)

	if err != nil || len(cuentaTelco.Valor) < 1 {
		err = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
		return
	}

	request.Transferencia.Origen.Cbu = cuentaTelco.Valor

	filtroConfiguracionMotivo := filtros.ConfiguracionFiltro{
		Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
	}

	motivo, err := s.utilService.GetConfiguracionService(filtroConfiguracionMotivo)
	if err != nil {
		err = fmt.Errorf("no se pudo encontrar el motivo para la transferencia")
		return
	}
	if motivo.Id == 0 {
		request := administraciondtos.RequestConfiguracion{
			Nombre:      "MOTIVO_TRANSFERENCIA_CLIENTE",
			Valor:       "VAR",
			Descripcion: "Motivo por defecto para las transferencias a los clientes",
		}

		_, err = s.utilService.CreateConfiguracionService(request)

		if err != nil {
			return
		}
	}
	request.Transferencia.Motivo = linkdtos.EnumMotivoTransferencia(motivo.Valor)

	// 2 - Garantizo que el dinero se depositará en la cuenta del cliente
	filtroCliente := filtros.CuentaFiltro{
		Id: uint(cuentaId),
	}

	cuenta, err := s.repository.GetCuenta(filtroCliente)

	if err != nil {
		return
	}

	request.Transferencia.Destino.Cbu = cuenta.Cbu

	// 3 - Me aseguro que el importe a transferir es en pesos

	request.Transferencia.Moneda = linkdtos.Pesos

	datosClientes := administraciondtos.DatosClientes{
		NombreCliente: cuenta.Cliente.Cliente,
		EmailCliente:  cuenta.Cliente.Email,
	}

	return s.BuildTransferencia(ctx, requerimientoId, request, cuentaId, datosClientes)

}

func (s *service) BuildTransferencia(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenicaCliente, cuentaId uint64, datosClientes administraciondtos.DatosClientes) (response linktransferencia.ResponseTransferenciaCreateLink, err error) {

	var fechaRetiro bool
	//var listaMovimientosComisionesImpuestos []uint64
	//var listaMovimientosRevert []uint64
	var listamovimientossincomision []uint64
	var movimientos []entities.Movimiento
	_, err = s.commonsService.IsValidUUID(requerimientoId)

	if err != nil {
		return
	}

	if cuentaId == 0 {
		err = fmt.Errorf(ERROR_CUENTA_ID)
		return
	}

	if request.Transferencia.Importe < 0 {
		err = fmt.Errorf(ERROR_IMPORTE_ENVIADO)
		return
	}

	// Recupar datos de la cuenta y verificar los dias de retiro automatico
	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(cuentaId),
	}
	cuenta, erro := s.GetCuenta(filtroCuenta)
	if erro != nil {
		return
	}
	fechaRetiroAutomatico := time.Now().AddDate(0, 0, int(-cuenta.DiasRetiroAutomatico))

	err = request.Transferencia.IsValid()

	if err != nil {
		return
	}

	// 1 - busca los movimientos por ids
	filtroMovimiento := filtros.MovimientoFiltro{
		Ids:                request.ListaMovimientosId,
		CargarPagoIntentos: true,
		CargarPago:         true,
		CargarComision:     true,
		CargarImpuesto:     true,
	}

	movimientos_response, _, err := s.repository.GetMovimientos(filtroMovimiento)

	if err != nil {
		return
	}

	// FIXME se debe controlar que los movimientos tengan comisiones y los impuestos
	totalSinComision := entities.Monto(0)
	for i := range movimientos_response {
		if len(movimientos_response[i].Movimientocomisions) > 0 && len(movimientos_response[i].Movimientoimpuestos) > 0 {
			movimientos = append(movimientos, movimientos_response[i])
		} else {
			totalSinComision += movimientos_response[i].Monto
			listamovimientossincomision = append(listamovimientossincomision, uint64(movimientos_response[i].ID))
			mensaje := fmt.Errorf("el movimiento %d no tiene comisiones/impuestos.No puede ser transferido ", movimientos_response[i].ID).Error()
			aviso := entities.Notificacione{
				Tipo:        entities.NotificacionTransferencia,
				Descripcion: mensaje,
				UserId:      0,
			}
			erro := s.utilService.CreateNotificacionService(aviso)

			if erro != nil {
				logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
			}
		}
	}

	// NOTE actualizar lista de id sacando las que no tienen comision.
	if len(listamovimientossincomision) > 0 {
		request.ListaMovimientosId = commons.DifferenceInteger(request.ListaMovimientosId, listamovimientossincomision)
		request.Transferencia.Importe = request.Transferencia.Importe - totalSinComision
	}

	/* NOTE CASO PAGOS REVERTIDOS Y NO */
	filtroMovimientoNeg := filtros.MovimientoFiltro{
		CuentaId:                   uint64(cuenta.Id),
		CargarMovimientosNegativos: true,
		CargarPagoIntentos:         true,
		CargarPago:                 true,
		AcumularPorPagoIntentos:    true,
	}
	movimientoNegativos, erro := s.repository.GetMovimientosNegativos(filtroMovimientoNeg)
	if erro != nil {
		err = erro
		return
	}
	/*CONSIDERAR 2 CASOS EN MOVIMIENTOS NEGATIVOS
	1 SI SON REVERSIONES
	2 SI LAS COMISIONES SON MAYORES AL MONTO NETO
	*/
	totalNeg := entities.Monto(0)
	if len(movimientoNegativos) > 0 {
		for i := range movimientoNegativos {
			totalNeg += movimientoNegativos[i].Monto
		}
		// comparar la lista enviada con la obtenoda en la base de datos
		if len(movimientoNegativos) != len(request.ListaMovimientosIdNeg) {
			err = errors.New(ERROR_MOVIMIENTO_LISTA_DIFERENCIA)
			stringMovimientosIdsNeg := make([]string, len(movimientoNegativos))
			stringListaMovimientosIdsNeg := make([]string, len(request.ListaMovimientosIdNeg))

			for i := range movimientoNegativos {
				stringMovimientosIdsNeg[i] = fmt.Sprint(movimientoNegativos[i].ID)
			}
			if len(request.ListaMovimientosIdNeg) == 0 {
				err = errors.New(ERROR_MOVIMIENTO_LISTA_DIFERENCIA)
			} else {
				for i := range request.ListaMovimientosIdNeg {
					stringListaMovimientosIdsNeg[i] = fmt.Sprint(request.ListaMovimientosIdNeg[i])
				}
			}

			/*
				Descripción:se comparan las lista enviada con la que se obtiene de la base de datos
			*/
			mensaje := fmt.Errorf("no se encontraron los siguientes movimientos neg. seleccionados %+v", commons.Difference(stringMovimientosIdsNeg, stringListaMovimientosIdsNeg)).Error()

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       mensaje,
				Funcionalidad: "BuildTransferencia",
			}
			erro := s.utilService.CreateLogService(log)

			if erro != nil {

				logs.Error(erro.Error() + mensaje)
			}
			return
		}
	}

	/* END CASO PAGOS REVERTIDOS */

	// 1.1 - verificar que las fechas de los movimientos correspondan con la fecha de retiro automatico
	for i := range movimientos {
		fechaRetiro = movimientos[i].CreatedAt.Before(fechaRetiroAutomatico)
		if !fechaRetiro {
			err = fmt.Errorf("el movimiento %d no se puede transferir no esta incluido en la fecha de retiro configurada", movimientos[i].ID)
			return
		}
	}

	// 2 - Si la lista de movimientos es vacia no se puede hacer nada porque no existen los movimientos seleccionados en la base de datos, si la lista encontrada no tiene la misma cantidad de elementos que la lista de entrada hay algún problema con los movimientos.
	if len(movimientos) == 0 || len(movimientos) != len(request.ListaMovimientosId) {

		err = errors.New(ERROR_MOVIMIENTO_LISTA_DIFERENCIA)

		stringMovimientosIds := make([]string, len(movimientos))
		stringListaMovimientosIds := make([]string, len(request.ListaMovimientosId))

		for i := range movimientos {
			stringMovimientosIds[i] = fmt.Sprint(movimientos[i].ID)
		}
		for i := range request.ListaMovimientosId {
			stringListaMovimientosIds[i] = fmt.Sprint(request.ListaMovimientosId[i])
		}

		mensaje := fmt.Errorf("no se encontraron los siguientes movimientos seleccionados %+v", commons.Difference(stringMovimientosIds, stringListaMovimientosIds)).Error()

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       mensaje,
			Funcionalidad: "BuildTransferencia",
		}

		erro := s.utilService.CreateLogService(log)

		if erro != nil {

			logs.Error(erro.Error() + mensaje)
		}

		return

	}
	// else {
	// 	// crear mi lista de movimientos que seran transferidos con comisiones impuestos
	// 	listaMovimientosComisionesImpuestos = append(listaMovimientosComisionesImpuestos, request.ListaMovimientosId...)
	// }

	// 3 - Busco el estado acreditado
	FiltroPagoEstado := filtros.PagoEstadoFiltro{
		BuscarPorFinal: true,
		Final:          true,
		Nombre:         config.MOVIMIENTO_ACCREDITED,
	}

	estadoAcreditado, err := s.repository.GetPagoEstado(FiltroPagoEstado)
	if err != nil {
		return
	}

	// 4 - Verifico si los pagos están en el estado acreditado y pertenecen a la cuenta del cliente
	for i := range movimientos {
		if movimientos[i].Pagointentos.Pago.PagoestadosID != int64(estadoAcreditado.ID) {
			err = fmt.Errorf("el movimiento %d que corresponde al pago %d no está acreditado", movimientos[i].ID, movimientos[i].Pagointentos.PagosID)
			return
		}
		if movimientos[i].CuentasId != cuentaId {
			err = fmt.Errorf("el movimiento %d que corresponde al pago %d no pertenece a la cuenta %d", movimientos[i].ID, movimientos[i].Pagointentos.PagosID, movimientos[i].CuentasId)
			return
		}
	}

	// Busco el saldo de los movimientos para saber si todavía pueden ser transferidos
	var listaPagoIntentos []uint64
	for i := range movimientos {
		listaPagoIntentos = append(listaPagoIntentos, movimientos[i].PagointentosId)
	}
	filtroSaldoMovimiento := filtros.MovimientoFiltro{AcumularPorPagoIntentos: true, PagoIntentosIds: listaPagoIntentos}

	saldoMovimientos, _, err := s.repository.GetMovimientos(filtroSaldoMovimiento)

	if err != nil {
		return
	}

	total := entities.Monto(0)
	for i := range saldoMovimientos {
		total += saldoMovimientos[i].Monto
	}

	// NOTE  aqui si exiten movimientos negativos en la cuenta  se debe descontar al importe que quiere transferir el cliente
	if totalNeg < 0 {
		total += totalNeg
	}

	if total != entities.Monto(request.Transferencia.Importe) {
		err = fmt.Errorf(ERROR_IMPORTE_TRANSFERENCIA)
		return
	}

	// 5 - Verifica si la cuenta tiene un saldo suficiente para realizar la transferencia
	saldoCuenta, err := s.repository.GetSaldoCuenta(cuentaId)

	if err != nil {
		return
	}

	if saldoCuenta.Total < request.Transferencia.Importe && saldoCuenta.Total > 0 {
		err = errors.New(ERROR_SALDO_CUENTA_INSUFICIENTE)
		return
	}

	// NOTE una vez que este todo correcto se debe dar de baja esos movimientos negativos

	// 6 - crea un movimiento de salida
	var listaMovimientos []*entities.Movimiento

	for i := range movimientos {
		movimiento := entities.Movimiento{
			CuentasId:      uint64(cuentaId),
			PagointentosId: uint64(movimientos[i].PagointentosId),
			Monto:          movimientos[i].Monto * -1.0,
			Tipo:           "D",
		}
		listaMovimientos = append(listaMovimientos, &movimiento)
	}

	if len(movimientoNegativos) > 0 {
		for i := range movimientoNegativos {
			var revert bool
			if !movimientoNegativos[i].Pagointentos.RevertedAt.IsZero() {
				revert = true
			}
			movimientoNeg := entities.Movimiento{
				CuentasId:      uint64(cuentaId),
				PagointentosId: uint64(movimientoNegativos[i].PagointentosId),
				Monto:          movimientoNegativos[i].Monto * -1.0,
				Tipo:           "D",
				Reversion:      revert,
			}
			listaMovimientos = append(listaMovimientos, &movimientoNeg)
		}
	}

	err = s.repository.CreateMovimientosTransferencia(ctx, listaMovimientos)

	if err != nil {
		return
	}

	// Elimino la lista de pagos para no tener conflicto en la transferencia.
	request.ListaMovimientosId = nil

	//Como referencia se usa el id de 1 movimiento que se está transfiriendo
	request.Transferencia.Referencia = strconv.FormatUint(uint64(movimientos[0].ID), 10)

	// 7 - Envío la solicitud de transferencia a apilink
	response, err = s.apilinkService.CreateTransferenciaApiLinkService(requerimientoId, request.Transferencia)

	if err != nil {
		logs.Info("ocurrio error al realizar transferencia servicio de apilink(CreateTransferenciaApiLinkService)" + fmt.Sprintf("%v", err))
		// 8.1 - En caso de que me tire un error en la transferencia doy de baja logica en los movimientos
		erro := s.repository.BajaMovimiento(ctx, listaMovimientos, err.Error())

		if erro != nil {
			// 8.1.1 - En caso de que no se puede cancelar los movimientos aviso al usuario para que intervenga manualmente.
			mensage := ""
			for i := range listaMovimientos {
				mensage += fmt.Sprintf("%d,", listaMovimientos[i].ID)
			}
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionTransferencia,
				Descripcion: fmt.Sprintf("atención los siguientes movimientos de transferencia se realizaron incorrectamente pero no pudieron ser cancelados, movimientosId: %s", mensage),
			}
			erro = s.utilService.CreateNotificacionService(notificacion)
			if erro != nil {
				logs.Error(erro.Error())
			}
		}
		return
	} else {
		// response.MovimientosIdTransferidos = listaMovimientosComisionesImpuestos
		// response.MovimientosIdReversiones = listaMovimientosRevert

		logs.Info("enviar email")
		mensaje := "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Comprobante Transferencia</h2><ul><li> Importe: <b>#0</b></li><li> Cbu Origen : <b>#1</b></li><li> Cbu Destino : <b>#2</b></li><li> Referencia bancaria: <b>#3</b></li><li> Fecha de operación: <b>#4</b></li></ul></p>"

		var arrayEmail []string
		arrayEmail = append(arrayEmail, datosClientes.EmailCliente, config.EMAIL_TELCO)
		importe := s.utilService.ToFixed(request.Transferencia.Importe.Float64(), 2)
		valorTransferido := s.utilService.FormatNum(importe)
		fechaOperacion := response.FechaOperacion.Format("02-01-2006")
		params := utildtos.RequestDatosMail{
			Email:            arrayEmail,
			Asunto:           "Comprobante Transferencia",
			Nombre:           datosClientes.NombreCliente,
			Mensaje:          mensaje,
			CamposReemplazar: []string{fmt.Sprintf("$%v", valorTransferido), request.Transferencia.Origen.Cbu, request.Transferencia.Destino.Cbu, response.NumeroReferenciaBancaria, fechaOperacion},
			From:             "Wee.ar!",
			TipoEmail:        "template",
		}
		erro = s.utilService.EnviarMailService(params)
		if erro != nil {
			logs.Error(erro.Error())
		}
	}

	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	// 8 - creo una transferencia para cada movimiento.
	// NOTE para conciliar con los registros del banco se debe agregar un campo con la fecha y numero de concilicacion bancaria que nos devuelve apilink
	var listaTransferencias []*entities.Transferencia

	for i := range listaMovimientos {
		var reversion bool
		if listaMovimientos[i].Reversion {
			reversion = true
		}
		transferencia := entities.Transferencia{
			MovimientosID:              uint64(listaMovimientos[i].ID),
			Referencia:                 request.Transferencia.Referencia,
			ReferenciaBancaria:         response.NumeroReferenciaBancaria,
			UserId:                     uint64(audit.UserID),
			Uuid:                       requerimientoId,
			CbuOrigen:                  request.Transferencia.Origen.Cbu,
			CbuDestino:                 request.Transferencia.Destino.Cbu,
			FechaOperacion:             &response.FechaOperacion,
			NumeroConciliacionBancaria: response.NumeroConciliacionBancaria,
			ReferenciaBanco:            commons.ConcatReferencia(&response.FechaOperacion, response.NumeroConciliacionBancaria),
			Reversion:                  reversion,
		}
		// lista de movimientos transferidos
		response.MovimientosIdTransferidos = append(response.MovimientosIdTransferidos, uint64(listaMovimientos[i].ID))
		listaTransferencias = append(listaTransferencias, &transferencia)
	}

	err = s.repository.CreateTransferencias(ctx, listaTransferencias)

	if err != nil {
		mensaje := fmt.Errorf("la transferencia %s por el importe de %d se realizo correctamente en apilink pero no se pudo guardar los valores en la tabla de transferencias. Ids de movimientos %v. ", response.NumeroReferenciaBancaria, request.Transferencia.Importe, request.ListaMovimientosId).Error()
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       mensaje,
			Funcionalidad: "BuildTransferencia",
			UserId:        audit.UserID,
		}

		aviso := entities.Notificacione{
			Tipo:        entities.NotificacionTransferencia,
			Descripcion: mensaje,
			UserId:      uint64(audit.UserID),
		}
		erro := s.utilService.CreateNotificacionService(aviso)

		if erro != nil {
			logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
		}

		erro = s.utilService.CreateLogService(log)

		if erro != nil {
			logs.Error(erro.Error() + mensaje)
		}
	}

	return
}

func (s *service) BuildTransferenciaClientePorMovimiento(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenciaMov) (response []linktransferencia.ResponseTransferenciaCreateLink, err error) {

	filtroConfiguracionCuenta := filtros.ConfiguracionFiltro{
		Nombre: "CBU_CUENTA_TELCO",
	}

	cuentaTelco, err := s.utilService.GetConfiguracionService(filtroConfiguracionCuenta)

	if err != nil || len(cuentaTelco.Valor) < 1 {
		err = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
		return
	}

	request.Transferencia.Origen.Cbu = cuentaTelco.Valor

	filtroConfiguracionMotivo := filtros.ConfiguracionFiltro{
		Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
	}

	motivo, err := s.utilService.GetConfiguracionService(filtroConfiguracionMotivo)
	if err != nil {
		err = fmt.Errorf("no se pudo encontrar el motivo para la transferencia")
		return
	}
	if motivo.Id == 0 {
		request := administraciondtos.RequestConfiguracion{
			Nombre:      "MOTIVO_TRANSFERENCIA_CLIENTE",
			Valor:       "VAR",
			Descripcion: "Motivo por defecto para las transferencias a los clientes",
		}

		_, err = s.utilService.CreateConfiguracionService(request)

		if err != nil {
			return
		}
	}
	request.Transferencia.Motivo = linkdtos.EnumMotivoTransferencia(motivo.Valor)

	// obteniendo movimientos id --> id cuenta

	// 1 - busca los movimientos por ids
	filtroMovimiento := filtros.MovimientoFiltro{
		Ids: request.Lista,
	}

	movimientos, _, err := s.repository.GetMovimientos(filtroMovimiento)

	movimientosPorCuenta := make(map[uint64][]entities.Movimiento)

	for _, mov := range movimientos {
		movimientosPorCuenta[mov.CuentasId] = append(movimientosPorCuenta[mov.CuentasId], mov)
	}

	for cuentaId, movimientosC := range movimientosPorCuenta {

		// 2 - Garantizo que el dinero se depositará en la cuenta del cliente
		filtroCliente := filtros.CuentaFiltro{
			Id: uint(cuentaId),
		}

		fmt.Print(movimientosC)

		cuenta, err := s.repository.GetCuenta(filtroCliente)

		if err != nil {
			fmt.Println("Error al obtener cuenta:", err)
		}

		request.Transferencia.Destino.Cbu = cuenta.Cbu

		// 3 - Me aseguro que el importe a transferir es en pesos

		request.Transferencia.Moneda = linkdtos.Pesos

		datosClientes := administraciondtos.DatosClientes{
			NombreCliente: cuenta.Cliente.Cliente,
			EmailCliente:  cuenta.Cliente.Email,
		}

		var mov_ids []uint64

		for _, v := range movimientosC {
			mov_ids = append(mov_ids, uint64(v.ID))
		}

		var request2 administraciondtos.RequestTransferenciaMov

		request2.Lista = mov_ids
		request2.Transferencia = request.Transferencia
		request2.ReferenciaBancaria = request.ReferenciaBancaria
		request2.NumeroConciliacionBancaria = request.NumeroConciliacionBancaria

		// Llama a la función y asigna los valores devueltos a variables
		transferenciaResponse, err := s.BuildTransferenciaMovimiento(ctx, requerimientoId, request2, cuentaId, datosClientes)

		if err != nil {
			fmt.Println("Error al construir transferencia:", err)
		}

		// Agrega la respuesta al slice de respuestas
		response = append(response, transferenciaResponse)
	}

	return

}

func (s *service) BuildTransferenciaMovimiento(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenciaMov, cuentaId uint64, datosClientes administraciondtos.DatosClientes) (response linktransferencia.ResponseTransferenciaCreateLink, err error) {
	_, err = s.commonsService.IsValidUUID(requerimientoId)

	if err != nil {
		return
	}

	if cuentaId == 0 {
		err = fmt.Errorf(ERROR_CUENTA_ID)
		return
	}

	if request.Transferencia.Importe < 0 {
		err = fmt.Errorf(ERROR_IMPORTE_ENVIADO)
		return
	}

	err = request.Transferencia.IsValid()

	if err != nil {
		return
	}

	filtroMovimiento := filtros.MovimientoFiltro{
		Ids:                request.Lista,
		CargarPagoIntentos: true,
		CargarPago:         true,
		CargarComision:     true,
		CargarImpuesto:     true,
	}

	movimientos, _, err := s.repository.GetMovimientos(filtroMovimiento)

	if err != nil {
		return
	}

	// 3 - Busco el estado acreditado
	FiltroPagoEstado := filtros.PagoEstadoFiltro{
		BuscarPorFinal: true,
		Final:          true,
		Nombre:         config.MOVIMIENTO_ACCREDITED,
	}

	estadoAcreditado, err := s.repository.GetPagoEstado(FiltroPagoEstado)
	if err != nil {
		return
	}

	// 4 - Verifico si los pagos están en el estado acreditado y pertenecen a la cuenta del cliente
	for i := range movimientos {
		if movimientos[i].Pagointentos.Pago.PagoestadosID != int64(estadoAcreditado.ID) {
			err = fmt.Errorf("el movimiento %d que corresponde al pago %d no está acreditado", movimientos[i].ID, movimientos[i].Pagointentos.PagosID)
			return
		}
		if movimientos[i].CuentasId != cuentaId {
			err = fmt.Errorf("el movimiento %d que corresponde al pago %d no pertenece a la cuenta %d", movimientos[i].ID, movimientos[i].Pagointentos.PagosID, movimientos[i].CuentasId)
			return
		}
	}

	// 6 - crea un movimiento de salida
	var listaMovimientos []*entities.Movimiento

	for i := range movimientos {
		movimiento := entities.Movimiento{
			CuentasId:      uint64(cuentaId),
			PagointentosId: uint64(movimientos[i].PagointentosId),
			Monto:          movimientos[i].Monto * -1.0,
			Tipo:           "D",
		}
		listaMovimientos = append(listaMovimientos, &movimiento)
	}

	err = s.repository.CreateMovimientosTransferencia(ctx, listaMovimientos)

	if err != nil {
		return
	}

	// Elimino la lista de pagos para no tener conflicto en la transferencia.
	request.Lista = nil

	//Como referencia se usa el id de 1 movimiento que se está transfiriendo
	request.Transferencia.Referencia = strconv.FormatUint(uint64(movimientos[0].ID), 10)

	// audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	// 8 - creo una transferencia para cada movimiento.
	// NOTE para conciliar con los registros del banco se debe agregar un campo con la fecha y numero de concilicacion bancaria que nos devuelve apilink
	var listaTransferencias []*entities.Transferencia

	for i := range listaMovimientos {
		var reversion bool
		if listaMovimientos[i].Reversion {
			reversion = true
		}

		response.FechaOperacion = time.Now().AddDate(0, 0, -1)
		response.NumeroConciliacionBancaria = request.NumeroConciliacionBancaria
		response.NumeroReferenciaBancaria = request.ReferenciaBancaria

		transferencia := entities.Transferencia{
			MovimientosID:              uint64(listaMovimientos[i].ID),
			Referencia:                 request.Transferencia.Referencia,
			ReferenciaBancaria:         response.NumeroReferenciaBancaria,
			UserId:                     1,
			Uuid:                       requerimientoId,
			CbuOrigen:                  request.Transferencia.Origen.Cbu,
			CbuDestino:                 request.Transferencia.Destino.Cbu,
			FechaOperacion:             &response.FechaOperacion,
			NumeroConciliacionBancaria: response.NumeroConciliacionBancaria,
			ReferenciaBanco:            commons.ConcatReferencia(&response.FechaOperacion, response.NumeroConciliacionBancaria),
			Reversion:                  reversion,
		}
		// lista de movimientos transferidos
		response.MovimientosIdTransferidos = append(response.MovimientosIdTransferidos, uint64(listaMovimientos[i].ID))
		listaTransferencias = append(listaTransferencias, &transferencia)
	}

	err = s.repository.CreateTransferencias(ctx, listaTransferencias)

	if err != nil {
		mensaje := fmt.Errorf("la transferencia %s por el importe de %d se realizo correctamente en apilink pero no se pudo guardar los valores en la tabla de transferencias. Ids de movimientos %v. ", response.NumeroReferenciaBancaria, request.Transferencia.Importe, request.Lista).Error()
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       mensaje,
			Funcionalidad: "BuildTransferencia",
			UserId:        1,
		}

		aviso := entities.Notificacione{
			Tipo:        entities.NotificacionTransferencia,
			Descripcion: mensaje,
			UserId:      1,
		}
		erro := s.utilService.CreateNotificacionService(aviso)

		if erro != nil {
			logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
		}

		erro = s.utilService.CreateLogService(log)

		if erro != nil {
			logs.Error(erro.Error() + mensaje)
		}
	}

	return
}

// transferencias de comisiones impuestos

func (s *service) SendTransferenciasComisiones(ctx context.Context, requerimientoId string, req administraciondtos.RequestComisiones) (res administraciondtos.ResponseTransferenciaComisiones, erro error) {

	// 1 consultar las transferencias del dia
	var controlmovimientosId []uint64
	var request administraciondtos.RequestTransferenciasComisiones
	var filtroTransferencia filtros.TransferenciaFiltro

	// NOTE depende de donde llega la solicitud
	// 1 Transferencia realizada por el cliente
	// 2 Transferencia automaticas de todos los clientes
	if len(req.MovimientosId) > 0 {
		filtroTransferencia = filtros.TransferenciaFiltro{
			MovimientosIds: req.MovimientosId,
		}
	} else {
		fecha_inicio := req.FechaInicio
		fecha_fin := req.FechaInicio
		if fecha_inicio.IsZero() {
			fecha_inicio = time.Now()
			fecha_fin = time.Now()
		}
		filtroTransferencia = filtros.TransferenciaFiltro{
			FechaInicio: fecha_inicio,
			FechaFin:    fecha_fin,
		}

	}

	transferencia, _, err := s.repository.GetTransferencias(filtroTransferencia)
	if err != nil {
		erro = err
		return
	}

	// 2 separar las que son reversiones y no
	var pagointentos []uint64
	var pagointentosrevertidos []uint64
	if len(transferencia) > 0 {
		for _, t := range transferencia {
			if !t.Reversion {
				pagointentos = append(pagointentos, t.Movimiento.PagointentosId)
			} else {
				pagointentosrevertidos = append(pagointentosrevertidos, t.Movimiento.PagointentosId)
			}
		}
	} else {
		res = administraciondtos.ResponseTransferenciaComisiones{
			Resultado: "no existe comisiones por transferir ",
		}
		return
	}

	// 3 consultar movimientos = sumar comisiones e impuestos
	// se crean variables para sumar montos
	var montoCom entities.Monto
	var montoRev entities.Monto

	filtromov := reportedtos.RequestPagosPeriodo{
		PagoIntentos:           pagointentos,
		TipoMovimiento:         "C",
		CargarComisionImpuesto: true,
		CargarReversion:        false,
	}
	var mov []entities.Movimiento
	var mov_revertidos []entities.Movimiento
	if len(pagointentos) > 0 {
		mov, erro = s.repository.GetMovimientosTransferencias(filtromov)
		if erro != nil {
			return
		}
		for _, mov := range mov {
			for _, comision := range mov.Movimientocomisions {
				montoCom += comision.Monto
			}
			for _, impuestos := range mov.Movimientoimpuestos {
				montoCom += impuestos.Monto
			}
			controlmovimientosId = append(controlmovimientosId, uint64(mov.ID))
		}
	}

	// 4 si existen reversiones restar a las comisiones
	if len(pagointentosrevertidos) > 0 {
		filtromov.PagoIntentos = pagointentosrevertidos // cambiar filtros para buscar los revertidos
		filtromov.CargarReversion = true
		mov_revertidos, erro = s.repository.GetMovimientosTransferencias(filtromov)
		if erro != nil {
			return
		}
		for _, mov_revertidos := range mov_revertidos {
			for _, comision := range mov_revertidos.Movimientocomisions {
				montoRev += comision.Monto
			}
			for _, impuestos := range mov_revertidos.Movimientoimpuestos {
				montoRev += impuestos.Monto
			}
			controlmovimientosId = append(controlmovimientosId, uint64(mov_revertidos.ID))
		}
		// restar las comisiones revertidas
		montoCom = montoCom + montoRev
		mov = append(mov, mov_revertidos...)
	}

	// 5 controlar si estas comisiones ya fueron transferidas
	filtroTransferencias := filtros.TransferenciaFiltro{
		MovimientosIds: controlmovimientosId,
	}
	transferencias, _, erro := s.repository.GetTransferenciasComisiones(filtroTransferencias)
	if erro != nil {
		mensage := ERROR_TRANSFERENCIA_COMISIONES
		logs.Info(mensage)
		return
	}

	// NOTE 6 existieron comisiones que ya fueron transferidas
	if len(transferencias) > 0 {

		// incluir en el mensaje email
		var mensaje string
		var arrayMovsInconsistentes []uint64

		// buscar los movimientos-transferencias-comisiones que estan repetidos
		for i := range transferencias {
			arrayMovsInconsistentes = append(arrayMovsInconsistentes, transferencias[i].MovimientosID)
		}
		// si hubo alguna transferencia inconsistente
		if len(arrayMovsInconsistentes) > 0 {
			erro = _enviarEmailNotificacionError(s, arrayMovsInconsistentes)
			if erro != nil {
				// si no pudo notificar mediante email se intenta crear un logs
				mensaje = fmt.Sprintf("Error: comisiones ya transferidas ids movimientos, transferir movimientos manualmente %+v", arrayMovsInconsistentes)
				logs.Error(erro.Error() + mensaje)
				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       mensaje,
					Funcionalidad: "SendTransferenciasComisiones",
				}

				erro := s.utilService.CreateLogService(log)

				if erro != nil {
					logs.Error(erro.Error() + mensaje)
				}
			}
			res = administraciondtos.ResponseTransferenciaComisiones{
				Resultado: "existen comisiones que ya fueron transferidas",
			}
			return
		}
	}

	// // monto que se va a transferir a la cuenta de telco
	// 7 datos para la transferencia
	filtroConfiguracionCuenta := filtros.ConfiguracionFiltro{
		Nombre: "CBU_CUENTA_TELCO",
	}

	cuentaTelco, err := s.utilService.GetConfiguracionService(filtroConfiguracionCuenta)

	if err != nil || len(cuentaTelco.Valor) < 1 {
		err = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
		logs.Info(err)
		return
	}

	request.Transferencia.Origen.Cbu = cuentaTelco.Valor

	filtroConfiguracionMotivo := filtros.ConfiguracionFiltro{
		Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
	}

	motivo, err := s.utilService.GetConfiguracionService(filtroConfiguracionMotivo)
	if err != nil {
		erro = fmt.Errorf("no se pudo encontrar el motivo para la transferencia")
		logs.Info(err)
		return
	}
	if motivo.Id == 0 {
		request := administraciondtos.RequestConfiguracion{
			Nombre:      "MOTIVO_TRANSFERENCIA_CLIENTE",
			Valor:       "VAR",
			Descripcion: "Motivo por defecto para las transferencias a los clientes",
		}

		_, erro = s.utilService.CreateConfiguracionService(request)

		if err != nil {
			return
		}
	}
	request.Transferencia.Motivo = linkdtos.EnumMotivoTransferencia(motivo.Valor)
	// obtener cbu destino (cuenta de telco)
	cbuDestino, erro := s.utilService.FirstOrCreateConfiguracionService("CBU_CUENTA_TELCO_DESTINO", "Cbu cuenta de telco para enviar comisiones impuestos de pagos", "0940099310007439910042")

	if erro != nil {
		return
	}
	request.Transferencia.Destino.Cbu = cbuDestino
	request.Transferencia.Moneda = linkdtos.Pesos

	_, erro = s.commonsService.IsValidUUID(requerimientoId)

	if err != nil {
		return
	}

	// 8 el monto a trasnferir debe ser mayor a 0
	if montoCom < 0 {
		erro = fmt.Errorf(ERROR_SALDO_CUENTA_INSUFICIENTE)
		logs.Info(erro)
		return
	}

	request.Transferencia.Importe = montoCom

	//Como referencia se usa el id de 1 movimiento que se está transfiriendo
	request.Transferencia.Referencia = strconv.FormatUint(uint64(mov[0].ID), 10)

	// 7 - Envío la solicitud de transferencia a apilink
	response, err := s.apilinkService.CreateTransferenciaApiLinkService(requerimientoId, request.Transferencia)

	if err != nil {
		logs.Info("ocurrio error al realizar transferencia servicio de apilink(SendTransferenciasComisiones)" + fmt.Sprintf("%v", err))

		if erro != nil {
			// 8.1.1 - En caso de que no se puede cancelar los movimientos aviso al usuario para que intervenga manualmente.
			mensage := ""
			for i := range mov {
				mensage += fmt.Sprintf("%d,", mov[i].ID)
			}
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionTransferencia,
				Descripcion: fmt.Sprintf("atención los siguientes movimientos de transferencia se realizaron incorrectamente movimientosId: %s", mensage),
			}
			erro = s.utilService.CreateNotificacionService(notificacion)
			if erro != nil {
				logs.Error(erro.Error())
			}
		}
		return
	} else {
		logs.Info("enviar email") // enviar email notificar las comisiones transferidas
		mensaje := "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Comprobante Transferencia comisiones</h2><ul><li> Importe: <b>#0</b></li><li> Cbu Origen : <b>#1</b></li><li> Cbu Destino : <b>#2</b></li><li> Referencia bancaria: <b>#3</b></li><li> Fecha de operación: <b>#4</b></li></ul></p>"
		var arrayEmail []string
		arrayEmail = append(arrayEmail, config.EMAIL_TELCO)
		importe := s.utilService.ToFixed(request.Transferencia.Importe.Float64(), 2)
		valorTransferido := s.utilService.FormatNum(importe)
		fechaOperacion := response.FechaOperacion.Format("02-01-2006")
		params := utildtos.RequestDatosMail{
			Email:            arrayEmail,
			Asunto:           "Comprobante Transferencia: Comisiones TelCo",
			Nombre:           "Wee",
			Mensaje:          mensaje,
			CamposReemplazar: []string{fmt.Sprintf("$%v", valorTransferido), request.Transferencia.Origen.Cbu, request.Transferencia.Destino.Cbu, response.NumeroReferenciaBancaria, fechaOperacion},
			From:             "Wee.ar!",
			TipoEmail:        "template",
		}
		erro := s.utilService.EnviarMailService(params)
		if erro != nil {
			logs.Error(erro.Error())
		}
	}

	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	// 8 - creo una transferencia para cada movimiento.
	// NOTE para conciliar con los registros del banco se debe agregar un campo con la fecha y numero de concilicacion bancaria que nos devuelve apilink
	var listaTransferencias []*entities.Transferenciacomisiones

	for i := range mov {
		var reversion bool
		if !mov[i].Pagointentos.RevertedAt.IsZero() {
			reversion = true
		}
		transferenciatelco := entities.Transferenciacomisiones{
			MovimientosID:              uint64(mov[i].ID),
			Referencia:                 request.Transferencia.Referencia,
			ReferenciaBancaria:         response.NumeroReferenciaBancaria,
			UserId:                     uint64(audit.UserID),
			Uuid:                       requerimientoId,
			CbuOrigen:                  request.Transferencia.Origen.Cbu,
			CbuDestino:                 request.Transferencia.Destino.Cbu,
			FechaOperacion:             &response.FechaOperacion,
			NumeroConciliacionBancaria: response.NumeroConciliacionBancaria,
			Reversion:                  reversion,
		}
		listaTransferencias = append(listaTransferencias, &transferenciatelco)
	}

	err = s.repository.CreateTransferenciasComisiones(ctx, listaTransferencias)

	if err != nil {
		mensaje := fmt.Errorf("la transferencia %s por el importe de %d se realizo correctamente en apilink pero no se pudo guardar los valores en la tabla de transferencias. Ids de movimientos %v. ", response.NumeroReferenciaBancaria, request.Transferencia.Importe, request.MovimientosIdComisiones).Error()
		logs.Info(mensaje)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       mensaje,
			Funcionalidad: "SendTransferenciasComisiones",
			UserId:        audit.UserID,
		}

		aviso := entities.Notificacione{
			Tipo:        entities.NotificacionTransferencia,
			Descripcion: mensaje,
			UserId:      uint64(audit.UserID),
		}
		erro := s.utilService.CreateNotificacionService(aviso)

		if erro != nil {
			logs.Error(erro.Error() + "no se pudo crear notificación en SendTransferenciasComisiones")
		}

		erro = s.utilService.CreateLogService(log)

		if erro != nil {
			logs.Error(erro.Error() + mensaje)
		}
	}

	res = administraciondtos.ResponseTransferenciaComisiones{
		Resultado: "transferencia exitosa",
	}
	return

}

func (s *service) GetMovimientosAcumulados(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoAcumuladoResponsePaginado, erro error) {

	var fechaRetiro bool
	resp, erro := s.GetMovimientos(filtro)

	if erro != nil {
		return
	}

	if len(resp.MovimientosNegativos) > 0 {
		movimientoResponse.MovimientosNegativos = resp.MovimientosNegativos
		for _, movNeg := range movimientoResponse.MovimientosNegativos {
			movimientoResponse.AcumuladosNegativos += movNeg.Monto
		}
	}

	// Recupar datos de la cuenta y verificar los dias de retiro automatico
	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(filtro.CuentaId),
	}
	cuenta, erro := s.GetCuenta(filtroCuenta)
	if erro != nil {
		return
	}
	fechaRetiroAutomatico := time.Now().AddDate(0, 0, int(-cuenta.DiasRetiroAutomatico))

	/* acumular registros de movimientos por fecha */
	var acumulado entities.Monto
	var movimientoAcumulado []administraciondtos.MovimientoPorCuentaResponse
	var fechaAux string
	for key, value := range resp.Acumulados {
		if fechaAux == "" {
			fechaAux = value.MovimientoCreated_at.Format("2006-01-02")
		}

		if len(resp.Acumulados) == 1 {
			fechaRetiro = value.MovimientoCreated_at.Before(fechaRetiroAutomatico)
			if fechaRetiro {
				acumulado += value.Monto
				movimientoAcumulado = append(movimientoAcumulado, value)
				//AGREGADO
				movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
					Fecha:                 fechaAux,
					FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
					Acumulado:             acumulado,
					Movimientos:           movimientoAcumulado,
				}
				movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
			}
			// END
			//
			break
		}
		if fechaAux == value.MovimientoCreated_at.Format("2006-01-02") {
			fechaRetiro = value.MovimientoCreated_at.Before(fechaRetiroAutomatico)
			if fechaRetiro {
				acumulado += value.Monto
				movimientoAcumulado = append(movimientoAcumulado, value)
				if len(resp.Acumulados)-1 == key {
					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
						Fecha:                 fechaAux,
						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
						Acumulado:             acumulado,
						Movimientos:           movimientoAcumulado,
					}
					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
				}
			}
		} else {
			movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
				Fecha:                 fechaAux,
				FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
				Acumulado:             acumulado,
				Movimientos:           movimientoAcumulado,
			}
			movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
			acumulado = 0
			movimientoAcumulado = []administraciondtos.MovimientoPorCuentaResponse{}
			fechaAux = value.MovimientoCreated_at.Format("2006-01-02")
			fechaRetiro = value.MovimientoCreated_at.Before(fechaRetiroAutomatico)
			if fechaRetiro {
				acumulado += value.Monto
				movimientoAcumulado = append(movimientoAcumulado, value)
				if len(resp.Acumulados)-1 == key {
					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
						Fecha:                 fechaAux,
						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
						Acumulado:             acumulado,
						Movimientos:           movimientoAcumulado,
					}
					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
				}
			}
		}
	}
	return
}

func (s *service) GetMovimientos(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoPorCuentaResponsePaginado, erro error) {

	movimiento, total, erro := s.repository.GetMovimientos(filtro)

	if erro != nil {
		return
	}

	// filtro movimientos negativos
	if filtro.CargarMovimientosNegativos {
		movimientoNegativos, err := s.repository.GetMovimientosNegativos(filtro)
		if err != nil {
			erro = err
			return
		}
		var movNegativos []administraciondtos.MovimientosNegativos
		for _, m := range movimientoNegativos {
			mov := administraciondtos.MovimientosNegativos{
				Id:                   m.ID,
				PagoIntentosId:       uint(m.PagointentosId),
				Tipo:                 string(m.Tipo),
				Pagotipo:             m.Pagointentos.Pago.PagosTipo.Pagotipo,
				ExternalReference:    m.Pagointentos.Pago.ExternalReference,
				MedioPago:            m.Pagointentos.Mediopagos.Mediopago,
				Monto:                m.Monto,
				PagoCreated_at:       m.Pagointentos.Pago.CreatedAt,
				MovimientoCreated_at: m.CreatedAt,
				Reversion:            !(m.Pagointentos.RevertedAt.IsZero()),
			}
			movNegativos = append(movNegativos, mov)
		}
		movimientoResponse.MovimientosNegativos = movNegativos
	}

	// filtro paginacion
	if filtro.Paginacion.Number > 0 && filtro.Paginacion.Size > 0 {

		from := (filtro.Number - 1) * filtro.Size
		lastPage := math.Ceil(float64(total) / float64(filtro.Size))

		meta := dtos.Meta{
			Page: dtos.Page{
				CurrentPage: int32(filtro.Number),
				From:        int32(from),
				LastPage:    int32(lastPage),
				PerPage:     int32(filtro.Size),
				To:          int32(filtro.Number * filtro.Size),
				Total:       int32(total),
			},
		}

		movimientoResponse.Meta = meta
	}

	// FOR para cada movimiento
	for i := range movimiento {

		var listaDetalleMovimientosComision []administraciondtos.DetalleComisionResponse
		var listaDetalleMovimientoImpuestos []administraciondtos.DetalleImpuestoResponse
		var totalComision entities.Monto
		var totalImpuesto entities.Monto

		// for slice de Movimientocomisions de cada movimiento
		for j := range movimiento[i].Movimientocomisions {
			// acumular montos comisiones
			totalComision += movimiento[i].Movimientocomisions[j].Monto + movimiento[i].Movimientocomisions[j].Montoproveedor
			detalleComision := administraciondtos.DetalleComisionResponse{
				Nombre:     movimiento[i].Movimientocomisions[j].Cuentacomisions.Cuentacomision,
				Monto:      movimiento[i].Movimientocomisions[j].Monto,
				Porcentaje: movimiento[i].Movimientocomisions[j].Porcentaje,
			}
			listaDetalleMovimientosComision = append(listaDetalleMovimientosComision, detalleComision)
		}

		// for slice de Movimientoimpuestos de cada movimiento
		for j := range movimiento[i].Movimientoimpuestos {
			// acumular montos impuestos
			totalImpuesto += movimiento[i].Movimientoimpuestos[j].Monto + movimiento[i].Movimientoimpuestos[j].Montoproveedor
			detalleImpuesto := administraciondtos.DetalleImpuestoResponse{
				Nombre:     movimiento[i].Movimientoimpuestos[j].Impuesto.Impuesto,
				Monto:      movimiento[i].Movimientoimpuestos[j].Monto,
				Porcentaje: movimiento[i].Movimientoimpuestos[j].Porcentaje,
			}
			listaDetalleMovimientoImpuestos = append(listaDetalleMovimientoImpuestos, detalleImpuesto)
		}

		movimientoComision := administraciondtos.MovimientoComisionResponse{
			Total:   totalComision,
			Detalle: listaDetalleMovimientosComision,
		}
		movimientoImpuestos := administraciondtos.MovimientoImpuestoResponse{
			Total:   totalImpuesto,
			Detalle: listaDetalleMovimientoImpuestos,
		}
		var fecha_rendicion string
		if len(movimiento[i].Movimientotransferencia) > 0 {
			fecha_rendicion = movimiento[i].Movimientotransferencia[len(movimiento[i].Movimientotransferencia)-1].FechaOperacion.Format("2006-01-02")
		}
		var montopagado entities.Monto
		montopagado = movimiento[i].Pagointentos.Amount
		if movimiento[i].Pagointentos.Valorcupon != 0 {
			montopagado = movimiento[i].Pagointentos.Valorcupon
		}

		var retenciones administraciondtos.RetencionesResponse
		movs_retenciones := entities.MovimientosRetenciones(movimiento[i].Movimientoretencions)
		// Acumular los montos de retenciones para el movimiento
		retenciones.MontoTotal = movs_retenciones.GetSumatoriaRetenciones()
		// guardar los montos por tipo de retencions si los tuviere
		retenciones.RetGanancias = movs_retenciones.GetImporteRetencionByName("ganancias")
		retenciones.RetIVA = movs_retenciones.GetImporteRetencionByName("iva")
		retenciones.RetIIBB = movs_retenciones.GetImporteRetencionByName("iibb")

		var mov_subcuentas = make([]administraciondtos.ResponseMovimientosubcuenta, len(movimiento[i].Movimientosubcuentas))
		for id, mov := range movimiento[i].Movimientosubcuentas {
			mov_subcuentas[id].ToEntity(mov)
		}

		response := administraciondtos.MovimientoPorCuentaResponse{
			Id:                   movimiento[i].ID,
			PagoIntentosId:       uint(movimiento[i].PagointentosId),
			Identificador:        movimiento[i].Pagointentos.Pago.Uuid,
			Estado:               movimiento[i].Pagointentos.Pago.PagoEstados.Estado,
			Tipo:                 string(movimiento[i].Tipo),
			Monto:                movimiento[i].Monto,
			Montopagado:          montopagado,
			Montosp:              movimiento[i].Pagointentos.Amount,
			Revertido:            movimiento[i].Reversion,
			Enobservacion:        movimiento[i].Enobservacion,
			Comisiones:           movimientoComision,
			Impuestos:            movimientoImpuestos,
			PagoCreated_at:       movimiento[i].Pagointentos.Pago.CreatedAt,
			MovimientoCreated_at: movimiento[i].CreatedAt,
			FechaRendicion:       fecha_rendicion,
			Retenciones:          retenciones,
			MovimientoSubcuentas: mov_subcuentas,
		}
		// filtro cargar datos del pago
		if filtro.CargarPago {
			response.PagoId = movimiento[i].Pagointentos.Pago.ID
			response.ExternalReference = movimiento[i].Pagointentos.Pago.ExternalReference
			response.Pagotipo = movimiento[i].Pagointentos.Pago.PagosTipo.Pagotipo
		}
		// filtro cargar datos del medio de pago
		if filtro.CargarMedioPago {
			response.MedioPago = movimiento[i].Pagointentos.Mediopagos.Mediopago
			response.Channels = movimiento[i].Pagointentos.Mediopagos.Channel.Nombre
		}

		if filtro.CargarTransferencias {
			if len(movimiento[i].Movimientotransferencia) > 0 {
				response.FechaRendicion = movimiento[i].Movimientotransferencia[len(movimiento[i].Movimientotransferencia)-1].FechaOperacion.Format("02-01-2006")
			}
		}

		movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, response)
	} // Fin de for i := range movimiento

	return
}

func (s *service) GetSaldoClienteService(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {

	if clienteId == 0 {
		erro = errors.New("debes informar el id del cliente")
		return
	}

	//Busco el saldo de un cliente por la lista de cuentas que tiene
	saldo, erro = s.repository.GetSaldoCliente(clienteId)
	saldo.ClienteId = clienteId

	return
}

func (s *service) GetSaldoCuentaService(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {

	if cuentaId == 0 {
		erro = errors.New("debes especificar una cuenta")
		return
	}

	saldo, erro = s.repository.GetSaldoCuenta(cuentaId)

	return
}

func (s *service) CreateMovimientosService(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error) {

	erro = s.repository.CreateMovimientosCierreLote(ctx, mcl)

	return
}

func (s *service) CreateMovimientosTemporalesService(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error) {

	erro = s.repository.CreateMovimientosTemporalesCierreLote(ctx, mcl)

	return
}

func (s *service) ActualizarPagosClRapipagoService(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error) {

	erro = s.repository.ActualizarPagosClRapipagoRepository(pagosclrapiapgo)

	return
}

func (s *service) GetTransferencias(filtro filtros.TransferenciaFiltro) (response administraciondtos.TransferenciaRespons, erro error) {

	// traer transferencias desde el repositorio con un filtro determinado
	transferencias, _, erro := s.repository.GetTransferencias(filtro)

	if erro != nil {
		return
	}

	// declaracion de variables locales
	var (
		responseTemporal             administraciondtos.TransferenciaRespons
		contador                     int
		recorrerHasta                int32
		movimientosTransferencias    []administraciondtos.MovimientoReponse
		movimientoPorCuentaAcumulado administraciondtos.TransferenciaResponseAgrupada
		num_ref                      string
		cbuOrigen                    string
		cbuDestino                   string
		fecha                        time.Time
		acumulado                    entities.Monto
		referencia_banco             string
		ids_transferencias_agrupadas []uint
	)

	// for _, t := range transferencias {
	// 	r := administraciondtos.TransferenciaResponse{}
	// 	r.New(t)
	// 	response.Transferencias = append(response.Transferencias, r)
	// }
	// var transferenicasAgrupadas   []administraciondtos.

	for key, t := range transferencias {

		// Estos atributos se repiten en cada registro de la tabla por cada movimiento distinto que existe
		if num_ref == "" {
			num_ref = t.ReferenciaBancaria
			cbuOrigen = t.CbuOrigen
			cbuDestino = t.CbuDestino
			fecha = *t.FechaOperacion
			referencia_banco = t.ReferenciaBanco
		}

		//  si en el respuesta al repositorio existe un solo registro
		if len(transferencias) == 1 {
			acumulado += t.Movimiento.Monto
			// los datos de cada movimiento
			movimientosTransferencias = append(movimientosTransferencias, administraciondtos.MovimientoReponse{
				Concepto:       t.Movimiento.Pagointentos.Pago.PagosTipo.Pagotipo,
				ReferenciaPago: t.Movimiento.Pagointentos.Pago.ExternalReference,
				CanalPago:      t.Movimiento.Pagointentos.Mediopagos.Channel.Nombre,
				MontoMov:       t.Movimiento.Monto * -1,
			})
			//AGREGADO
			movimientoPorCuentaAcumulado = administraciondtos.TransferenciaResponseAgrupada{
				ReferenciaBancaria:         num_ref,
				CbuOrigen:                  cbuOrigen,
				CbuDestino:                 cbuDestino,
				Fecha:                      fecha,
				Monto:                      acumulado * -1,
				MovimientoReponse:          movimientosTransferencias,
				ReferenciaBanco:            referencia_banco,
				IdsTransferenciasAgrupadas: ids_transferencias_agrupadas,
			}
			// Acumulado de cada transferencia
			responseTemporal.Transferencias = append(responseTemporal.Transferencias, movimientoPorCuentaAcumulado)

			break
		}

		if num_ref == t.ReferenciaBancaria {
			// acumular ids de cada transferencia
			ids_transferencias_agrupadas = append(ids_transferencias_agrupadas, t.ID)
			// acumulado cabecera
			acumulado += t.Movimiento.Monto
			// los datos de cada movimiento
			if t.Movimiento.Pagointentos != nil {
				movimientosTransferencias = append(movimientosTransferencias, administraciondtos.MovimientoReponse{
					Concepto:       t.Movimiento.Pagointentos.Pago.PagosTipo.Pagotipo,
					ReferenciaPago: t.Movimiento.Pagointentos.Pago.ExternalReference,
					CanalPago:      t.Movimiento.Pagointentos.Mediopagos.Channel.Nombre,
					MontoMov:       t.Movimiento.Monto * -1,
				})
			}

			if len(transferencias)-1 == key {
				movimientoPorCuentaAcumulado = administraciondtos.TransferenciaResponseAgrupada{
					ReferenciaBancaria:         num_ref,
					CbuOrigen:                  cbuOrigen,
					CbuDestino:                 cbuDestino,
					Fecha:                      fecha,
					Monto:                      acumulado * -1,
					MovimientoReponse:          movimientosTransferencias,
					ReferenciaBanco:            referencia_banco,
					IdsTransferenciasAgrupadas: ids_transferencias_agrupadas,
				}
				// Acumulado de cada tranaferencia
				responseTemporal.Transferencias = append(responseTemporal.Transferencias, movimientoPorCuentaAcumulado)
			}
		} else {
			//AGREGADO
			movimientoPorCuentaAcumulado = administraciondtos.TransferenciaResponseAgrupada{
				ReferenciaBancaria:         num_ref,
				CbuOrigen:                  cbuOrigen,
				CbuDestino:                 cbuDestino,
				Fecha:                      fecha,
				Monto:                      acumulado * -1,
				MovimientoReponse:          movimientosTransferencias,
				ReferenciaBanco:            referencia_banco,
				IdsTransferenciasAgrupadas: ids_transferencias_agrupadas,
			}
			// Acumulado de cada transferencia
			responseTemporal.Transferencias = append(responseTemporal.Transferencias, movimientoPorCuentaAcumulado)

			// se ponen en cero las variables de acumulacion y se asignan las del proximo conjunto de movimientos, que tienen atributos de transferencia iguales
			acumulado = 0
			movimientosTransferencias = []administraciondtos.MovimientoReponse{}
			// responseTemporal = administraciondtos.TransferenciaRespons{}
			movimientoPorCuentaAcumulado = administraciondtos.TransferenciaResponseAgrupada{}
			num_ref = t.ReferenciaBancaria
			cbuOrigen = t.CbuOrigen
			cbuDestino = t.CbuDestino
			fecha = *t.FechaOperacion
			referencia_banco = t.ReferenciaBanco
			ids_transferencias_agrupadas = nil

			acumulado += t.Movimiento.Monto

			movimientosTransferencias = append(movimientosTransferencias, administraciondtos.MovimientoReponse{
				Concepto:       t.Movimiento.Pagointentos.Pago.PagosTipo.Pagotipo,
				ReferenciaPago: t.Movimiento.Pagointentos.Pago.ExternalReference,
				CanalPago:      t.Movimiento.Pagointentos.Mediopagos.Channel.Nombre,
				MontoMov:       t.Movimiento.Monto * -1,
			})

			// acumular id de la transferencia que a este punto ya es de otro grupo
			ids_transferencias_agrupadas = append(ids_transferencias_agrupadas, t.ID)

			if len(transferencias)-1 == key {
				movimientoPorCuentaAcumulado = administraciondtos.TransferenciaResponseAgrupada{
					ReferenciaBancaria:         t.ReferenciaBancaria,
					CbuOrigen:                  cbuOrigen,
					CbuDestino:                 cbuDestino,
					Fecha:                      fecha,
					Monto:                      acumulado * -1,
					MovimientoReponse:          movimientosTransferencias,
					ReferenciaBanco:            referencia_banco,
					IdsTransferenciasAgrupadas: ids_transferencias_agrupadas,
				}
				// Acumulado de cada tranaferencia
				responseTemporal.Transferencias = append(responseTemporal.Transferencias, movimientoPorCuentaAcumulado)
			}

		}

	}

	contador = len(responseTemporal.Transferencias)
	var totalTransferencia float64
	if contador > 0 {
		for _, transf := range responseTemporal.Transferencias {
			totalTransferencia += transf.Monto.Float64()
		}

	}
	response.TotalTransferencias = s.utilService.FormatNum(s.utilService.ToFixed(totalTransferencia, 2))

	if filtro.Paginacion.Number > 0 && filtro.Paginacion.Size > 0 {
		response.Meta = _setPaginacion(filtro.Paginacion.Number, filtro.Paginacion.Size, int64(contador))
	}
	recorrerHasta = response.Meta.Page.To
	if response.Meta.Page.CurrentPage == response.Meta.Page.LastPage {
		recorrerHasta = response.Meta.Page.Total
	}
	if recorrerHasta == 0 {
		recorrerHasta = int32(len(responseTemporal.Transferencias))
	}

	if len(responseTemporal.Transferencias) > 0 {
		if len(responseTemporal.Transferencias) < int(recorrerHasta) {
			erro = errors.New("Numero de pagina incorrecto para la búsqueda")
			return
		}

		for i := response.Meta.Page.From; i < recorrerHasta; i++ {
			if int(i) < len(responseTemporal.Transferencias) {
				response.Transferencias = append(response.Transferencias, responseTemporal.Transferencias[i])
			}
		}
	}
	return

}

// actuliazar transferencias conciliacion banco
func (s *service) UpdateTransferencias(listas bancodtos.ResponseConciliacion) error {
	return s.repository.UpdateTransferencias(listas)
}

func (s *service) GetCierreLoteRapipagoService(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error) {
	response, erro = s.repository.GetConsultarMovimientosRapipago(filtro)
	if erro != nil {
		return
	}
	return
}

func (s *service) BuildRapipagoMovimiento(listaCierre []*entities.Rapipagocierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

	// * 1 - listadeCierre es el resultado de consultar la lista de rapipago cierre lote y conciliacion con banco
	if len(listaCierre) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	// * 2 generar una lista de barcode(de la listaCierre consultada en rapipago) para luego consultar los pagos intentos  */
	// solo los encontrandos en banco */
	var listaCierreRapipago []string                                       // lista de barcode de la lista de cierre lote de rapipago
	var listaCierreRapipagoDetalles []*entities.Rapipagocierrelotedetalles // *De la lista solo los pagos encontrados en banco
	for i := range listaCierre {
		if listaCierre[i].BancoExternalId != 0 {
			// verifica que los pagos de este cierre fueron transferidos
			// listaCierre[i].PagoMatch = true
			movimientoCierreLote.ListaCLRapipagoHeaders = append(movimientoCierreLote.ListaCLRapipagoHeaders, listaCierre[i])
			for j := range listaCierre[i].RapipagoDetalle {
				if listaCierre[i].RapipagoDetalle[j].Match {
					listaCierreRapipago = append(listaCierreRapipago, listaCierre[i].RapipagoDetalle[j].CodigoBarras)
					listaCierreRapipagoDetalles = append(listaCierreRapipagoDetalles, listaCierre[i].RapipagoDetalle[j])
				}
			}
		}
	}

	canalOffline, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_OFFLINE", "Nombre del canal de debin", "offline")

	if erro != nil {
		return
	}

	// * 3 Busco el canal offline para filtrar
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalOffline},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)
	if erro != nil {
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Barcode:              listaCierreRapipago,
		Channel:              true,
		CargarPago:           true,
		CargarPagoTipo:       true,
		CargarCuenta:         true,
		CargarCliente:        true,
		CargarCuentaComision: true,
		CargarImpuestos:      true,
	}
	// * 4 - Busco los pagos intentos que corresponden a los barcode de la lista de cierre
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	// * 5 Verifico que los pagos intentos correspondan al canal offline
	for i := range pagosIntentos {
		if pagosIntentos[i].Mediopagos.ChannelsID == int64(channel.ID) {
			movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, pagosIntentos[i])
		}
	}

	// * 6 - Busco el estado acreditado
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, erro := s.repository.GetPagoEstado(filtroPagoEstado)

	if erro != nil {
		return
	}

	// * 7 la lista de cierre loterapipago debe ser igual a la lista de pagos de intentos
	if len(listaCierreRapipagoDetalles) != len(movimientoCierreLote.ListaPagoIntentos) {
		var listaPagosExternalId []string
		// guardar el barcode de los pagos intentos en una lista
		for i := range movimientoCierreLote.ListaPagoIntentos {
			listaPagosExternalId = append(listaPagosExternalId, movimientoCierreLote.ListaPagoIntentos[i].ExternalID)
		}
		mensaje := fmt.Errorf("no se encontrarion los siguientes pagos %+v", commons.Difference(listaCierreRapipago, listaPagosExternalId)).Error()
		erro = fmt.Errorf(ERROR_CIERRE_PAGO_INTENTO)
		log := entities.Log{
			Tipo:          entities.Warning,
			Funcionalidad: "BuildRapipagoMovimiento",
			Mensaje:       mensaje,
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Info(ERROR_LOG + "BuildRapipagoMovimiento." + erro.Error())
		}
		return
	}

	// * 8 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range movimientoCierreLote.ListaPagoIntentos {
		/* * para el calculo de la comision fitrar por el id del channel y el id de la cuentar*/
		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         channel.ID,
			CuentaId:          movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.ID,
			FechaPagoVigencia: movimientoCierreLote.ListaPagoIntentos[i].PaidAt,
			Channelarancel:    true,
		}

		logs.Info(filtroComisionChannel)

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}
		for j := range listaCierreRapipagoDetalles {
			/* * control de que el barcode del cierre de lote sea igual al barcode del pago intento*/
			if movimientoCierreLote.ListaPagoIntentos[i].Barcode == listaCierreRapipagoDetalles[j].CodigoBarras {
				movimientoCierreLote.ListaCLRapipago = append(movimientoCierreLote.ListaCLRapipago, *listaCierreRapipagoDetalles[j])
				// FIXME Hay que ver como controlar ese error si hay que abortar
				// * controlar que los montosn coinicidan
				if movimientoCierreLote.ListaPagoIntentos[i].Amount != entities.Monto(listaCierreRapipagoDetalles[j].ImporteCobrado) {
					erro = fmt.Errorf("el monto informado no es valido verificar cierre lote rapipago %s", movimientoCierreLote.ListaPagoIntentos[i].Barcode)
					return
				}
				// * modificar el pago intento a estado acreditado
				movimientoCierreLote.ListaPagoIntentos[i].Pago.PagoestadosID = int64(pagoEstadoAcreditado.ID)
				movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, movimientoCierreLote.ListaPagoIntentos[i].Pago)

				// * crear el log de estado de pago acreditado
				pagoEstadoLog := entities.Pagoestadologs{
					PagosID:       movimientoCierreLote.ListaPagoIntentos[i].PagosID,
					PagoestadosID: int64(pagoEstadoAcreditado.ID),
				}
				movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)

				if listaCierreRapipagoDetalles[j].Match {
					movimiento := entities.Movimiento{}
					movimiento.AddCredito(uint64(movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.CuentasID), uint64(movimientoCierreLote.ListaPagoIntentos[i].ID), entities.Monto(listaCierreRapipagoDetalles[j].ImporteCobrado))

					s.utilService.BuildComisiones(&movimiento, &listaCuentaComision, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, entities.Monto(listaCierreRapipagoDetalles[j].ImporteCobrado))

					if listaCierreRapipagoDetalles[j].Enobservacion {
						movimiento.Enobservacion = true
					}
					movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
					movimientoCierreLote.ListaPagoIntentos[i].AvailableAt = listaCierreRapipagoDetalles[j].CreatedAt
				}
				break
			}
		}
	}
	return
}

func (s *service) BuildPagosClRapipago(listaCierre []*entities.Rapipagocierrelote) (pagosclrapiapgo administraciondtos.PagosClRapipagoResponse, erro error) {

	// & 1 - listadeCierrelote
	if len(listaCierre) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}
	// & 2 generar una lista de barcode(de la listaCierre consultada en rapipago) para luego consultar los pagos para actualizar  */
	// solo los encontrandos en banco */
	var listaCierreRapipago []string                                       // lista de barcode de la lista de cierre lote de rapipago
	var listaCierreRapipagoDetalles []*entities.Rapipagocierrelotedetalles // *De la lista solo los pagos encontrados en banco
	for i := range listaCierre {
		pagosclrapiapgo.ListaCLRapipagoHeaders = append(pagosclrapiapgo.ListaCLRapipagoHeaders, listaCierre[i].ID)
		for j := range listaCierre[i].RapipagoDetalle {
			listaCierreRapipago = append(listaCierreRapipago, listaCierre[i].RapipagoDetalle[j].CodigoBarras)
			listaCierreRapipagoDetalles = append(listaCierreRapipagoDetalles, listaCierre[i].RapipagoDetalle[j])
		}
	}

	canalOffline, erro := s.utilService.FirstOrCreateConfiguracionService("CHANNEL_OFFLINE", "Nombre del canal de debin", "offline")

	if erro != nil {
		return
	}

	// & 3 Busco el canal offline para filtrar
	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{canalOffline},
	}
	channel, erro := s.repository.GetChannel(filtroChannel)
	if erro != nil {
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Barcode:              listaCierreRapipago,
		Channel:              true,
		CargarPago:           true,
		CargarPagoTipo:       true,
		CargarCuenta:         true,
		CargarCliente:        true,
		CargarCuentaComision: true,
		CargarImpuestos:      true,
	}
	// & 4 - Busco los pagos intentos que corresponden a los barcode de la lista de cierre
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	// & 5 Verifico que los pagos intentos correspondan al canal offline
	for i := range pagosIntentos {
		if pagosIntentos[i].Mediopagos.ChannelsID != int64(channel.ID) {
			erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE_RAPIPAGO)
			return
		} else {
			pagosclrapiapgo.ListaPagos = append(pagosclrapiapgo.ListaPagos, pagosIntentos[i].Pago.ID)
		}
	}

	// & 6 la lista de cierre loterapipago debe ser igual a la lista de pagos de intentos
	if len(listaCierreRapipagoDetalles) != len(pagosIntentos) {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE_RAPIPAGO)
		return
	}

	// & 1 - Busco el estado acreditado
	paid, erro := s.utilService.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}

	pagoEstadoAprobado, erro := s.repository.GetPagoEstado(filtroPagoEstado)

	if erro != nil {
		return
	}

	pagosclrapiapgo.EstadoAprobado = pagoEstadoAprobado.ID

	return
}

func (s *service) UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error) {
	erro = s.repository.UpdateCierreLoteRapipago(cierreLotes)
	return
}

func _enviarEmailNotificacionError(s *service, ids []uint64) error {

	mensaje := "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Error transferencias comisiones</h2><ul><li> Los ids de movimientos que no se pudieron procesar son: <b>#0</b></li></ul></p>"

	var arrayEmail []string

	arrayEmail = append(arrayEmail, config.EMAIL_TELCO)

	params := utildtos.RequestDatosMail{
		Email:            arrayEmail,
		Asunto:           "Error en transferencias de comisiones",
		Nombre:           "Wee",
		Mensaje:          mensaje,
		CamposReemplazar: []string{fmt.Sprintf("%v", ids)},
		From:             "Wee.ar!",
		TipoEmail:        "template",
	}
	erro := s.utilService.EnviarMailService(params)
	if erro != nil {
		logs.Error(erro.Error())
	}
	return erro
}

// func (s *service) GetMovimientosAcumulados(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoAcumuladoResponsePaginado, erro error) {

// 	var fechaRetiro bool
// 	resp, erro := s.GetMovimientos(filtro)

// 	if erro != nil {
// 		return
// 	}

// 	if len(resp.MovimientosNegativos) > 0 {
// 		movimientoResponse.MovimientosNegativos = resp.MovimientosNegativos
// 	}

// 	// Recupar datos de la cuenta y verificar los dias de retiro automatico
// 	filtroCuenta := filtros.CuentaFiltro{
// 		Id: uint(filtro.CuentaId),
// 	}
// 	cuenta, erro := s.GetCuenta(filtroCuenta)
// 	if erro != nil {
// 		return
// 	}
// 	fechaRetiroAutomatico := time.Now().AddDate(0, 0, int(-cuenta.DiasRetiroAutomatico))

// 	/* acumular registros de movimientos por fecha */
// 	var acumulado entities.Monto
// 	var movimientoAcumulado []administraciondtos.MovimientoPorCuentaResponse
// 	for i := range resp.Acumulados {
// 		if len(resp.Acumulados)-1 != i {
// 			if resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02") == resp.Acumulados[i+1].MovimientoCreated_at.Format("2006-01-02") {
// 				fechaRetiro = resp.Acumulados[i].MovimientoCreated_at.Before(fechaRetiroAutomatico)
// 				if fechaRetiro {
// 					acumulado += resp.Acumulados[i].Monto
// 					movimientoAcumulado = append(movimientoAcumulado, resp.Acumulados[i])
// 				}
// 			} else {
// 				fechaRetiro = resp.Acumulados[i].MovimientoCreated_at.Before(fechaRetiroAutomatico)
// 				if fechaRetiro {
// 					acumulado += resp.Acumulados[i].Monto
// 					movimientoAcumulado = append(movimientoAcumulado, resp.Acumulados[i])
// 					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
// 						Fecha:                 resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02"),
// 						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
// 						Acumulado:             acumulado,
// 						Movimientos:           movimientoAcumulado,
// 					}
// 					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
// 					acumulado = 0
// 					movimientoAcumulado = []administraciondtos.MovimientoPorCuentaResponse{}
// 				}
// 			}
// 		} else {

// 			if len(resp.Acumulados) == 1 {
// 				fechaRetiro = resp.Acumulados[i].MovimientoCreated_at.Before(fechaRetiroAutomatico)
// 				if fechaRetiro {
// 					acumulado += resp.Acumulados[i].Monto
// 					movimientoAcumulado = append(movimientoAcumulado, resp.Acumulados[i])
// 					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
// 						Fecha:                 resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02"),
// 						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
// 						Acumulado:             acumulado,
// 						Movimientos:           movimientoAcumulado,
// 					}
// 					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
// 					acumulado = 0
// 					movimientoAcumulado = []administraciondtos.MovimientoPorCuentaResponse{}
// 				}
// 			} else if resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02") == resp.Acumulados[i-1].MovimientoCreated_at.Format("2006-01-02") {
// 				fechaRetiro = resp.Acumulados[i].MovimientoCreated_at.Before(fechaRetiroAutomatico)
// 				if fechaRetiro {
// 					acumulado += resp.Acumulados[i].Monto
// 					movimientoAcumulado = append(movimientoAcumulado, resp.Acumulados[i])
// 					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
// 						Fecha:                 resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02"),
// 						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
// 						Acumulado:             acumulado,
// 						Movimientos:           movimientoAcumulado,
// 					}

// 					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
// 				}

// 			} else {
// 				fechaRetiro = resp.Acumulados[i].MovimientoCreated_at.Before(fechaRetiroAutomatico)
// 				if fechaRetiro {
// 					acumulado += resp.Acumulados[i].Monto
// 					movimientoAcumulado = append(movimientoAcumulado, resp.Acumulados[i])
// 					movimientoPorCuentaAcumulado := administraciondtos.MovimientoPorCuentaAcumulado{
// 						Fecha:                 resp.Acumulados[i].MovimientoCreated_at.Format("2006-01-02"),
// 						FechaDisponibleRetiro: fechaRetiroAutomatico.Format("2006-01-02"),
// 						Acumulado:             acumulado,
// 						Movimientos:           movimientoAcumulado,
// 					}
// 					movimientoResponse.Acumulados = append(movimientoResponse.Acumulados, movimientoPorCuentaAcumulado)
// 					acumulado = 0
// 					movimientoAcumulado = []administraciondtos.MovimientoPorCuentaResponse{}
// 				}
// 			}
// 		}
// 	}

// 	return
// }
