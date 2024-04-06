package routes

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
	"github.com/gofiber/fiber/v2"
)

func BancoRoutes(app fiber.Router, service banco.BancoService, middlewares middlewares.MiddlewareManager) {
	app.Get("/obtener_movimientos", middlewares.ValidarPermiso("psp.herramienta"), getConsultarMovimientos(service))

	app.Get("/procesar_movimientos_herramienta", middlewares.ValidarPermiso("psp.herramienta"), getProcesarMovimientosHerramienta(service))
}

func getConsultarMovimientos(service banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var params bancodtos.RequestParams
		var request filtros.RequestMovimientos
		err := c.QueryParser(&params)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}
		request.ParamsToFiltro(params)
		response, err := service.GetConsultarMovimientosService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(fiber.Map{
			"data": response,
		})
	}
}

func getProcesarMovimientosHerramienta(service banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		response, err := service.GetProcesarMovimientosGeneralService()
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(fiber.Map{
			"data": response,
		})
	}
}
