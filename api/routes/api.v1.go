package routes

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares/middlewareinterno"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app fiber.Router, middlewares middlewares.MiddlewareManager, middlewaresInterno middlewareinterno.MiddlewareManagerPasarela, service administracion.Service, util util.UtilService) {
	/*
		Con este endpoint se puede consultar el estado de todos los pagos realizados
		se requiere que envie por el header de la peticion una api-key
		y enviar algunos de los posibles datos por el body:
			- si quiere consultar por un pago especifico deberia enviar
				un Uuid o ExternalReference estos atributos son de tipo String
			- puede consultar por varios pagos por lo que puede enviar:
				+ un rango de fecha, FechaDesde y FechaHasta, el formato de la misma es "DD-MM-AAAA" y el rango de fecha no debe superar los 7 dias
			- o puede enviar un array de Uuids.
	*/
	app.Post("/estado-pagos", middlewares.ValidarPermiso("psp.consultar.pago"), postConsultarEstadoPagos(service))
	app.Post("/estados-pagos", middlewaresInterno.ValidarApiKeyCliente(), postConsultarEstadoPagosForApiky(service))
}

func postConsultarEstadoPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Error: Debe enviar una api key válida")
		}

		var request administraciondtos.RequestPagosConsulta

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: parámetros incorrectos")
		}

		requestValid, err := request.IsParamsValid()
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		data, registrosEfectados, err := service.ConsultarEstadoPagosService(requestValid, api, request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		message := "datos enviados"
		if !registrosEfectados {
			message = "no se encontró  coincidencia para la consulta realizada"
		}
		total_registros := len(data)
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"total":   total_registros,
			"message": message,
		})
	}
}

func postConsultarEstadoPagosForApiky(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Error: Debe enviar una api key válida")
		}

		var request administraciondtos.RequestPagosConsulta

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: parámetros incorrectos")
		}

		requestValid, err := request.IsParamsValid()
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		data, registrosEfectados, err := service.ConsultarEstadoPagosService(requestValid, api, request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		message := "datos enviados"
		if !registrosEfectados {
			message = "no se encontró  coincidencia para la consulta realizada"
		}
		total_registros := len(data)
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"total":   total_registros,
			"message": message,
		})
	}
}
