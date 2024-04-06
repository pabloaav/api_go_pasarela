package routes

import (
	"fmt"
	"math"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares/middlewareinterno"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	dtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	apiresponder "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/responderdto"
	"github.com/gofiber/fiber/v2"
)

func ReporteRoutes(app fiber.Router, middlewares middlewares.MiddlewareManager, middlewaresInterno middlewareinterno.MiddlewareManagerPasarela, service reportes.ReportesService, runEndpoint util.UtilService) {
	// reporte telco
	app.Get("/cliente", middlewares.ValidarPermiso("psp.reportes.pago"), getReporte(service)) //("psp.reportes.pago")

	// reporte clente se envia por correo
	app.Get("/send-pagos", middlewares.ValidarPermiso("psp.reportes.pago"), sendPagos(service))             //("psp.reportes.pago")
	app.Get("/send-rendiciones", middlewares.ValidarPermiso("psp.reportes.pago"), sendRendiciones(service)) //("psp.reportes.rendiciones")
	app.Get("/send-reversiones", middlewares.ValidarPermiso("psp.reportes.pago"), sendReversiones(service)) //("psp.reportes.reversiones") , middlewares.ValidarPermiso("psp.reportes.reversiones")
	// Enviar por FTP(DPEC)
	app.Get("/batch-pago-items", middlewares.ValidarPermiso("psp.reportes.pago"), batchPagoItems(service)) //("psp.reportes.cobranza") DPEC

	// app.Get("/recaudaciones", middlewares.ValidarPermiso("psp.reportes.pago"), recaudacion(service))                //("psp.reportes.cobranza")
	// app.Get("/orden-pago-movimientos", middlewares.ValidarPermiso("psp.reportes.pago"), recaudacionDiaria(service)) //("psp.reportes.cobranza")
	/// endpointtemporal para generar ordenes de pago desde el inicio de los pagos de wee
	// app.Get("/orden-de-pago", ordenDePago(service)) //("psp.reportes.cobranza") DEPEC

	// Reportes para clientes por api key
	// app.Get("/cobranzas", middlewaresInterno.ValidarApiKeyCliente(), reportesCobranzas(service)) //("psp.reportes.pago")
	app.Get("/cobranzas", middlewaresInterno.ValidarApiKeyCliente(), reportesCobranzasTemporal(service))
	app.Get("/rendiciones", middlewaresInterno.ValidarApiKeyCliente(), reportesRendiciones(service)) //("psp.reportes.pago")
	app.Get("/reversiones", middlewaresInterno.ValidarApiKeyCliente(), reportesReversiones(service)) //("psp.reportes.pago")

	// Reporte Moivmiento-Comisiones
	app.Get("/movimientos-comisiones", middlewares.ValidarPermiso("psp.accesos.comisiones"), reporteMovimientosComisiones(service))
	app.Get("/movimientos-comisiones-temporales", middlewares.ValidarPermiso("psp.accesos.comisiones"), reporteMovimientosComisionesTemporales(service))

	// Reporte Cobranzas-Cliente
	app.Get("/cobrazas-cliente", middlewares.ValidarPermiso("psp.consultar.pagos"), reportesCobranzasClientes(service))

	// Verifica que los montos que devuelven distintos endpoints de cobranzas
	app.Get("/verificar-cobranzas", verificarCobranzasClienteCobranzas(service, runEndpoint))

	// Reporte Rendiciones-Cliente
	app.Get("/rendiciones-cliente", middlewares.ValidarPermiso("psp.consultar.pagos"), reportesRendicionesClientes(service))

	// Reporte Rendiciones-Cliente
	app.Get("/reversiones-cliente", middlewares.ValidarPermiso("psp.consultar.pagos"), reportesReversionesClientes(service))

	// Reportes Informacion Generales
	app.Get("/peticiones", middlewares.ValidarPermiso("psp.reportes.pago"), reportesPeticiones(service)) //("psp.reportes.pago")

	app.Get("/logs", middlewares.ValidarPermiso("psp.herramienta"), reportesLogs(service))                     // psp.herramienta
	app.Get("/notificaciones", middlewares.ValidarPermiso("psp.herramienta"), reportesNotificaciones(service)) // psp.herramienta
	// REPORTES HECHOS y ENVIADOS
	app.Get("/reportes-enviados", middlewares.ValidarPermiso("psp.herramienta"), reportesEnviados(service))
	// app.Get("/reportes-enumerar", EnumerarReportesEnviados(service)) // Funcion para asignar nros reportes (first_time)
}

func verificarCobranzasClienteCobranzas(service reportes.ReportesService, runEndpoint util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		type requestControlCobranzas struct {
			Apykeys        []string `json:"apykeys"`
			FechaConsultar string   `json:"fecha_consultar"`
		}

		type ResponseControlReporte struct {
			Estado               string                         `json:"estado"`
			CuentaId             string                         `json:"cuenta_id"`
			Cuenta               string                         `json:"cuenta"`
			MontoCobranza        string                         `json:"monto_cobranza"`
			MontoCobranzaCliente string                         `json:"monto_cobranza_cliente"`
			Diferencia           string                         `json:"diferencia"`
			PagosRepetidos       []dtos.DetallesPagosCobranza   `json:"pagos_repetidos"`
			PagosNoCalculados    []dtos.ResponseDetalleCobranza `json:"pagos_no_calculados"`
		}

		var requestEndpoint requestControlCobranzas
		var responseEndpoint []ResponseControlReporte

		err := c.QueryParser(&requestEndpoint)

		if err != nil {
			return c.Status(404).JSON(&fiber.Map{
				"message": "Error parseando el request",
			})
		}

		fechaConsultar := requestEndpoint.FechaConsultar

		for _, apyKey := range requestEndpoint.Apykeys {

			cuentaCliente, _ := service.GetCuentaByApiKeyService(apyKey)

			requestCobranzaCliente := dtos.RequestCobranzasClientes{
				FechaInicio: fechaConsultar,
				FechaFin:    fechaConsultar,
				ClienteId:   int(cuentaCliente.ClientesID),
				CuentaId:    int(cuentaCliente.ID),
			}

			// Deserializar la respuesta JSON en la estructura ResponseData
			respCobranzasCliente, _ := service.GetCobranzasClientesService(requestCobranzaCliente)

			requestCobranza := dtos.RequestCobranzas{
				Date: fechaConsultar,
			}

			// Deserializar la respuesta JSON en la estructura ResponseData
			respCobranzas, _ := service.GetCobranzasTemporal(requestCobranza, apyKey)

			if respCobranzasCliente.CantidadCobranzas <= 0 {
				controlReporteCuenta := ResponseControlReporte{
					Estado:               "No hay movimientos",
					CuentaId:             respCobranzas.AccountId,
					Cuenta:               cuentaCliente.Cuenta,
					MontoCobranzaCliente: fmt.Sprint(0),
					MontoCobranza:        fmt.Sprint(0),
				}

				responseEndpoint = append(responseEndpoint, controlReporteCuenta)

				continue
			}

			fmt.Println("Cobranza varios dias: $", respCobranzasCliente.Cobranzas[0].Subtotal)
			fmt.Println("Cobranza individual: $", uint(math.Round(respCobranzas.TotalCollected*100.00)))

			if uint(math.Round(respCobranzas.TotalCollected*100.00)) != respCobranzasCliente.Cobranzas[0].Subtotal {

				// buscando aquellos pagos que no están calculados
				var pagosNoCalculados []dtos.ResponseDetalleCobranza

				for _, cobranza := range respCobranzas.Data {
					if cobranza.NetFee == 0 {
						pagosNoCalculados = append(pagosNoCalculados, cobranza)
					}
				}

				if len(pagosNoCalculados) > 0 {

					controlReporteCuenta := ResponseControlReporte{
						Estado:               "Hay Diferencia",
						CuentaId:             respCobranzas.AccountId,
						Cuenta:               cuentaCliente.Cuenta,
						MontoCobranzaCliente: fmt.Sprint(respCobranzasCliente.Cobranzas[0].Subtotal),
						MontoCobranza:        fmt.Sprint(uint(math.Round(respCobranzas.TotalCollected * 100.00))),
						Diferencia:           "Pagos No Calculados",
						PagosNoCalculados:    pagosNoCalculados,
					}

					responseEndpoint = append(responseEndpoint, controlReporteCuenta)

					continue

				}

				// buscando pagos repetidos
				var pagosRepetidos []dtos.DetallesPagosCobranza
				visto := make(map[int]bool)
				for _, pago := range respCobranzasCliente.Cobranzas[0].Pagos {
					// Buscando cobranzas repetidas
					if _, ok := visto[pago.Id]; ok {
						// Ya hemos visto una cobranza con este RequestId
						pagosRepetidos = append(pagosRepetidos, pago)
					} else {
						// Marcar como visto
						visto[pago.Id] = true
					}
				}

				if len(pagosRepetidos) > 0 {
					controlReporteCuenta := ResponseControlReporte{
						Estado:               "Hay Diferencia",
						CuentaId:             respCobranzas.AccountId,
						Cuenta:               cuentaCliente.Cuenta,
						MontoCobranzaCliente: fmt.Sprint(respCobranzasCliente.Cobranzas[0].Subtotal),
						MontoCobranza:        fmt.Sprint(uint(math.Round(respCobranzas.TotalCollected * 100.00))),
						Diferencia:           "Pagos Repetidos",
						PagosRepetidos:       pagosRepetidos,
					}

					responseEndpoint = append(responseEndpoint, controlReporteCuenta)

					continue
				}

				controlReporteCuenta := ResponseControlReporte{
					Estado:               "Hay Diferencia",
					CuentaId:             respCobranzas.AccountId,
					Cuenta:               cuentaCliente.Cuenta,
					MontoCobranzaCliente: fmt.Sprint(respCobranzasCliente.Cobranzas[0].Subtotal),
					MontoCobranza:        fmt.Sprint(uint(math.Round(respCobranzas.TotalCollected * 100.00))),
					Diferencia:           "Otra razón",
				}

				responseEndpoint = append(responseEndpoint, controlReporteCuenta)

			}

			controlReporteCuenta := ResponseControlReporte{
				Estado:               "OK",
				CuentaId:             respCobranzas.AccountId,
				Cuenta:               cuentaCliente.Cuenta,
				MontoCobranzaCliente: fmt.Sprint(respCobranzasCliente.Cobranzas[0].Subtotal),
				MontoCobranza:        fmt.Sprint(uint(math.Round(respCobranzas.TotalCollected * 100.00))),
			}

			responseEndpoint = append(responseEndpoint, controlReporteCuenta)
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    responseEndpoint,
			"message": "Se verificaron las cobranzas correctamente.",
		})

	}
}

// probando reporte
func getReporte(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dtos.RequestPagosPeriodo
		err := c.QueryParser(&request)
		if err != nil {
			logs.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
		request.Paginacion.Number = request.Number
		request.Paginacion.Size = request.Size
		data, err := service.GetPagosReportes(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		result, err := service.ResultPagosReportes(data, request.Paginacion)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		// Respuesta
		return c.Status(200).JSON(&fiber.Map{
			"status": "ok",
			"data":   result,
		})
	}
}

func sendPagos(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		clientes, err := service.GetClientes(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaPagosClientes, err := service.GetPagosClientes(clientes, request)
				if err != nil {
					return fiber.NewError(400, "Error: "+err.Error())
				}

				// enviar los pagos a clientes
				if len(listaPagosClientes) > 0 {
					listaErro, err := service.SendPagosClientes(listaPagosClientes)
					if err != nil {
						r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
						return r.Responder()
					} else if len(listaErro) > 0 {
						r := apiresponder.NewResponse(400, listaErro, "error: no se pudo enviar reporte de pagos al cliente", c)
						return r.Responder()
					} else {
						// caso de exito del proceso
						r := apiresponder.NewResponse(200, nil, "el proceso se ejecuto con éxito", c)
						return r.Responder()
					}
				} else {
					r := apiresponder.NewResponse(400, "notificacionPagos", "no existen pagos por enviar", c)
					return r.Responder()
				}
			}

		}

		r := apiresponder.NewResponse(404, "ReportesPagosEnviados", "no existen clientes: verifique datos enviados", c)
		return r.Responder()
	}
}

func sendRendiciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaRendicionClientes, err := service.GetRendicionClientes(clientes, request)
				if err != nil {
					return fiber.NewError(400, "Error: "+err.Error())
				}
				// enviar los pagos a clientes
				if len(listaRendicionClientes) > 0 {
					listaErro, err := service.SendPagosClientes(listaRendicionClientes)
					if err != nil {
						return fiber.NewError(400, "Error: "+err.Error())
					} else if len(listaErro) > 0 {
						return c.JSON(&fiber.Map{
							"error":    listaErro,
							"response": "nil",
						})
					}
				} else {
					return c.JSON(&fiber.Map{
						"error":            "no existen rendiciones por enviar",
						"tipoconciliacion": "notificacionPagos",
					})
				}
			}

		}
		return c.JSON(&fiber.Map{
			"error":            "",
			"tipoconciliacion": "ReportesRendicionesEnviados",
		})

	}
}

func sendReversiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request dtos.RequestPagosClientes

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: parametros recibidos no son validos")
		}
		// se valida que las fechas no sean nulas y que la fecha este antes de la fecha fin o sean iguales
		filtro, err := request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error %v", err.Error()))
		}
		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaReversionesClientes, err := service.GetReversionesClientes(clientes, request, filtro)
				if err != nil {
					return fiber.NewError(400, "Error: "+err.Error())
				}
				if len(listaReversionesClientes) <= 0 {
					return c.Status(200).JSON(&fiber.Map{
						"status":  true,
						"message": "no exiten reversiones",
					})
				}
				// enviar los pagos a clientes
				listaErro, err := service.SendPagosClientes(listaReversionesClientes)
				logs.Info(listaErro)
				if err != nil {
					return fiber.NewError(400, "Error: "+err.Error())
				}
			}

		}
		return c.JSON(&fiber.Map{
			"error":            "",
			"status":           true,
			"message":          "success",
			"tipoconciliacion": "ReportesReversionesEnviados",
		})

	}
}

func batchPagoItems(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)
		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos/pagoitems por cliente
				listaPagosItems, err := service.GetPagoItems(clientes, request)
				if err != nil {
					return fiber.NewError(400, "Error: "+err.Error())
				}
				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
				if len(listaPagosItems) > 0 {
					resultpagositems := service.BuildPagosItems(listaPagosItems)
					if len(resultpagositems) > 0 {
						err := service.ValidarEsctucturaPagosItems(resultpagositems) // validar estructura antes de crear el archivo
						if err != nil {
							return fiber.NewError(400, "Error: "+err.Error())
						} else {
							err := service.SendPagosItems(ctx, resultpagositems, request)
							if err != nil {
								return fiber.NewError(400, "Error: "+err.Error())
							}
							return c.JSON(&fiber.Map{
								"error":            "Reportes batch se ejecuto con éxito",
								"tipoconciliacion": "batchPagoItems",
							})
						}
					}
				} else {
					return c.JSON(&fiber.Map{
						"error":            "no existen pagos para informar",
						"tipoconciliacion": "batchPagoItems",
					})
				}
			}

		}
		return c.JSON(&fiber.Map{
			"error":            "no existen clientes para informar reporte",
			"tipoconciliacion": "batchPagoItems",
		})

	}
}

func reportesCobranzas(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}
		var request dtos.RequestCobranzas

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetCobranzas(request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": res,
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud cobranzas generada",
		})

	}
}

func reportesCobranzasTemporal(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}
		var request dtos.RequestCobranzas

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetCobranzasTemporal(request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": res,
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud cobranzas generada",
		})

	}
}

func reportesRendiciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}
		var request dtos.RequestCobranzas

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetRendiciones(request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": res,
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud rendiones generada",
		})

	}
}
func reportesReversiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}
		var request dtos.RequestCobranzas

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetReversiones(request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para la fecha",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud reversiones generada",
		})

	}
}

// func recaudacionDiaria(service reportes.ReportesService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		c.Accepts("application/json")

// 		var request dtos.RequestPagosClientes
// 		err := c.QueryParser(&request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
// 		}

// 		// 1 obtener lista de cliente
// 		clientes, err := service.GetClientes(request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error: "+err.Error())
// 		}
// 		if len(clientes.Clientes) > 0 {
// 			if err == nil {
// 				// obtener los pagos/pagoitems por cliente
// 				listaMovItems, err := service.GetRecaudacionDiaria(clientes, request)
// 				if err != nil {
// 					return fiber.NewError(400, "Error: "+err.Error())
// 				}
// 				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
// 				if len(listaMovItems) > 0 {
// 					// contruir datos que se enviara en la liquidacion diaria
// 					resultpagositems := service.BuildMovLiquidacion(listaMovItems)
// 					if len(resultpagositems) > 0 {
// 						logs.Info(resultpagositems)
// 						errorFile, err := service.SendLiquidacionClientes(resultpagositems)
// 						logs.Info(errorFile)
// 						if err != nil {
// 							return fiber.NewError(400, "Error: "+err.Error())
// 						}
// 						return c.JSON(&fiber.Map{
// 							"error":            "Reportes liquidacion se ejecuto con éxito",
// 							"tipoconciliacion": "batchPagoItems",
// 						})
// 					}

// 				} else {
// 					return c.JSON(&fiber.Map{
// 						"error":            "no existen pagos para informar en liquidacion",
// 						"tipoconciliacion": "batchPagoItems",
// 					})
// 				}
// 			}

// 		}
// 		return c.JSON(&fiber.Map{
// 			"error":            "no existen clientes para informar reporte",
// 			"tipoconciliacion": "batchPagoItems",
// 		})

// 	}
// }

// func recaudacion(service reportes.ReportesService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		c.Accepts("application/json")

// 		var request dtos.RequestPagosClientes
// 		err := c.QueryParser(&request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
// 		}

// 		// 1 obtener lista de cliente
// 		clientes, err := service.GetClientes(request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error: "+err.Error())
// 		}
// 		if len(clientes.Clientes) > 0 {
// 			if err == nil {
// 				// obtener los pagos/pagoitems por cliente
// 				listaPagosItems, err := service.GetRecaudacion(clientes, request)
// 				if err != nil {
// 					return fiber.NewError(400, "Error: "+err.Error())
// 				}
// 				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
// 				if len(listaPagosItems) > 0 {
// 					// contruir datos que se enviara en la liquidacion diaria
// 					resultpagositems := service.BuildPagosLiquidacion(listaPagosItems)
// 					if len(resultpagositems) > 0 {
// 						// errorFile, err := service.SendLiquidacionClientes(resultpagositems)
// 						// logs.Info(errorFile)
// 						// if err != nil {
// 						// 	return fiber.NewError(400, "Error: "+err.Error())
// 						// }
// 						// return c.JSON(&fiber.Map{
// 						// 	"error":            "Reportes liquidacion se ejecuto con éxito",
// 						// 	"tipoconciliacion": "batchPagoItems",
// 						// })
// 					}

// 				} else {
// 					return c.JSON(&fiber.Map{
// 						"error":            "no existen pagos para informar en liquidacion",
// 						"tipoconciliacion": "batchPagoItems",
// 					})
// 				}
// 			}

// 		}
// 		return c.JSON(&fiber.Map{
// 			"error":            "no existen clientes para informar reporte",
// 			"tipoconciliacion": "batchPagoItems",
// 		})

// 	}
// }

// func ordenDePago(service reportes.ReportesService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		c.Accepts("application/json")

// 		var request dtos.RequestPagosClientes
// 		err := c.QueryParser(&request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
// 		}
// 		if len(request.ClientesString) > 0 {
// 			request.ObtenerIdsClientes()
// 		}
// 		// 1 obtener lista de cliente
// 		clientes, err := service.GetClientes(request)
// 		if err != nil {
// 			return fiber.NewError(400, "Error: "+err.Error())
// 		}
// 		// contar cantidad de dias entre fecha desde y fecha hasta
// 		// realizar bucle hasta cantidad de dias e ir sumando undia a la vez
// 		// crear un nuevo objeto reques dodne la fecha desde y hasta sea la misma y enviar a GetRecaudacionesDiarias
// 		vslorTemporal := request.FechaInicio.Sub(request.FechaFin)
// 		cantidadDias := int(vslorTemporal.Hours()/24) * -1
// 		fechaProceso := request.FechaInicio
// 		if len(clientes.Clientes) > 0 {
// 			if err == nil {
// 				// se recorre fecha a fecha segun rango de fecha enviado por query
// 				for i := 0; i <= cantidadDias; i++ {
// 					logs.Info(fechaProceso)
// 					requestTemporal := dtos.RequestPagosClientes{
// 						FechaInicio: fechaProceso,
// 						FechaFin:    fechaProceso,
// 						EnviarEmail: false,
// 					}
// 					// obtener los pagos/pagoitems por cliente
// 					listaMovItems, err := service.GetRecaudacionDiaria(clientes, requestTemporal)
// 					if err != nil {
// 						return fiber.NewError(400, "Error: "+err.Error())
// 					}
// 					// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
// 					if len(listaMovItems) > 0 && listaMovItems[len(listaMovItems)-1].Clientes.EnviarPdf {
// 						// contruir datos que se enviara en la liquidacion diaria
// 						resultpagositems := service.BuildMovLiquidacion(listaMovItems)
// 						logs.Info(resultpagositems)
// 						if len(resultpagositems) > 0 {
// 							logs.Info(resultpagositems)
// 							errorFile, err := service.SendLiquidacionClientes(resultpagositems)
// 							logs.Info(errorFile)
// 							if err != nil {
// 								return fiber.NewError(400, "Error: "+err.Error())
// 							}
// 							return c.JSON(&fiber.Map{
// 								"error":            "Reportes liquidacion se ejecuto con éxito",
// 								"tipoconciliacion": "batchPagoItems",
// 							})
// 						}

// 					} else {
// 						logs.Info(len(listaMovItems))
// 						logs.Info(listaMovItems[len(listaMovItems)-1].Clientes.EnviarPdf)
// 						logs.Info("generando numero de liquidacion")
// 						// return c.JSON(&fiber.Map{
// 						// 	"error":            "no existen pagos para informar en liquidacion",
// 						// 	"tipoconciliacion": "batchPagoItems",
// 						// })
// 					}

// 					fechaProceso = fechaProceso.Add(time.Hour * 24)
// 				}

// 			}

// 		}
// 		return c.JSON(&fiber.Map{
// 			"message":          "porcesos de generacion de ordenes finalizada",
// 			"tipoconciliacion": "batchPagoItems",
// 		})

// 	}
// }

func reportesPeticiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPeticiones

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, "Error en la validación de los parámetros enviados: "+err.Error())
		}

		res, err := service.GetPeticiones(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para los filtros enviados",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud peticiones generada",
		})

	}
}

func reporteMovimientosComisiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestReporteMovimientosComisiones

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.MovimientosComisionesService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Reportes) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"data":   "Sin resultados encontrados para la busqueda",
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud reporte movimientos generada",
		})

	}
}
func reporteMovimientosComisionesTemporales(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestReporteMovimientosComisiones

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.MovimientosComisionesTemporales(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Reportes) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"data":   "Sin resultados encontrados para la busqueda",
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud reporte movimientos generada",
		})

	}
}

func reportesCobranzasClientes(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestCobranzasClientes

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetCobranzasClientesService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Cobranzas) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"data":   "Sin resultados encontrados para la busqueda",
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud reporte cobranzas clientes generada",
		})

	}
}

func reportesRendicionesClientes(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestReporteClientes

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetRendicionesClientesService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		// if len(res.Clientes) == 0 {
		// 	return c.Status(200).JSON(&fiber.Map{
		// 		"status": true,
		// 		"data":   "Sin resultados encontrados para la busqueda",
		// 	})
		// }

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud reporte rendiciones clientes generada",
		})

	}
}

func reportesReversionesClientes(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestReporteClientes

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		res, err := service.GetReversionesClientesService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.DetallesReversiones) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"data":   "Sin resultados encontrados para la busqueda",
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud reporte reversiones clientes generada",
		})

	}
}

func reportesLogs(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestLogs

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, "Error en la validación de los parámetros enviados: "+err.Error())
		}

		res, err := service.GetLogs(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para la fecha",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud logs generada",
		})

	}
}

func reportesNotificaciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestNotificaciones

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, "Error en la validación de los parámetros enviados: "+err.Error())
		}

		res, err := service.GetNotificaciones(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para la fecha",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud notificaciones generada",
		})

	}
}

/* recibe filtro de fecha y de tipo de reporte enviado */
func reportesEnviados(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// filtro de la request
		var filtroReportesEnviados dtos.RequestReportesEnviados

		// parse de los parametros de la request al filtro CierreLoteFiltro
		err := c.QueryParser(&filtroReportesEnviados)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en los parametros recibidos "+err.Error(), c)
			return r.Responder()
		}

		// Validar los parametros
		err = filtroReportesEnviados.Validar()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}

		// Enviar la consulta al servicio correspondiente
		result, err := service.GetReportesEnviadosService(filtroReportesEnviados)

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		// si no hubo resultados en la consulta, pero tampoco errores, devolver en data un string vacio
		if len(result.Reportes) == 0 {

			r := apiresponder.NewResponse(200, []string{}, "Datos de consulta enviados, sin resultados", c)
			return r.Responder()

		}

		r := apiresponder.NewResponse(200, result, "Datos de consulta enviados", c)
		return r.Responder()
	}
}

func EnumerarReportesEnviados(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		filtro := dtos.RequestReportesEnviados{
			TipoReporte: "todos",
			Enum:        true,
		}
		// Enviar la consulta al servicio correspondiente
		err := service.EnumerarReportesEnviadosService(filtro)

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		// Enviar la consulta al servicio correspondiente
		err = service.CopiarNumeroReporteOriginal()

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, "Válido", "Operacion Realizada correctamente", c)
		return r.Responder()
	}
}
