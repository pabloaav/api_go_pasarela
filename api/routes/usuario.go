package routes

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/usuario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/userdtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/usuario"
	"github.com/gofiber/fiber/v2"
)

func UsuarioRoutes(app fiber.Router, middlewares middlewares.MiddlewareManager, service usuario.UsuarioService) {
	app.Get("/user", middlewares.ValidarPermiso("usuario.show"), getUsuario(service))
	app.Get("/users", middlewares.ValidarPermiso("usuario.index"), getUsuarios(service))
	app.Post("/user", middlewares.ValidarPermiso("usuario.create"), postUsuario(service))
	app.Put("/user", middlewares.ValidarPermiso("usuario.update"), putUsuario(service))
}

func getUsuario(service usuario.UsuarioService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.UserFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		bearer := c.Get("Authorization")

		requestAutorizacion := filtros.UserFiltroAutenticacion{
			Token: bearer,
			User:  request,
		}

		response, err := service.GetUsuarioService(requestAutorizacion)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getUsuarios(service usuario.UsuarioService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.UserFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		bearer := c.Get("Authorization")

		requestAutorizacion := filtros.UserFiltroAutenticacion{
			Token: bearer,
			User:  request,
		}

		response, err := service.GetUsuariosService(requestAutorizacion)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postUsuario(service usuario.UsuarioService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request userdtos.RequestUser

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		requestAutorizacion := userdtos.RequestUserAutorizacion{
			Token:   bearer,
			Request: &request,
		}

		id, err := service.CreateUsuarioService(requestAutorizacion)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putUsuario(service usuario.UsuarioService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request userdtos.RequestUserUpdate

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		requestAutorizacion := userdtos.RequestUserAutorizacion{
			Token:         bearer,
			RequestUpdate: &request,
		}

		err = service.UpdateUsuarioService(requestAutorizacion)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}
