package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"

	//"github.com/aws/aws-sdk-go/service"
	"github.com/gofiber/fiber/v2"
)

func PrismaRoutes(app fiber.Router, service prisma.Service, middlewares middlewares.MiddlewareManager) {
	app.Get("/checkservice", middlewares.ValidarPermiso("psp.checkservicioprisma"), checkServicioPrisma(service)) //psp.checkServicioPrisma
	app.Get("/listar-pagos", middlewares.ValidarPermiso("psp.listarpagos"), consultarPagos(service))              //psp.listarPagos
	app.Get("/obtener-pagos", middlewares.ValidarPermiso("psp.obtenerpagos"), obtenerPagos(service))              //psp.obtenerPagos
	app.Post("/anulacion-devolucion-total-pago", middlewares.ValidarPermiso("psp.obtenerpagos"), anularPagoTotal(service))
	app.Get("/informacion-pago", middlewares.ValidarPermiso("psp.obtenerpagos"), getObtenerPago(service))
}

func checkServicioPrisma(service prisma.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fecha := c.Query("fecha")
		estadoServicio, err := service.CheckService()
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"fecha-consulta":  fecha,
			"estado-servicio": estadoServicio,
		})
	}
}
func consultarPagos(service prisma.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request prismadtos.ListaPagosRequest
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}
		result, err := service.ListarPagosPorFecha(request)
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"data": result,
		})
	}
}

func obtenerPagos(service prisma.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		estadoPago := c.Query("estado_pago")
		estadoPagoId, err := strconv.Atoi(estadoPago)
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		channel := c.Query("channel")
		pagos, err := service.ListarPagosService(estadoPagoId, strings.ToUpper(channel))
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"listaPagos": pagos,
		})
	}
}

func anularPagoTotal(service prisma.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request NroComprobante
		err := c.BodyParser(&request)
		if err != nil {
			logs.Error("peticion error: " + err.Error())
			return fiber.NewError(404, "Error: debe enviar un comprobante de pago valido")
		}
		response, err := service.PostAnulacionDevolucionTotalPago(request.PaymentId)
		if err != nil {
			logs.Error("peticion error: " + err.Error())
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status": "ok",
			"data":   response,
		})
	}
}

func getObtenerPago(service prisma.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestPaymentId := c.Query("ComprobantePago")
		response, err := service.GetInformePago(requestPaymentId)
		if err != nil {
			logs.Error("peticion error: " + err.Error())
			return fiber.NewError(404, "Error: no se pudo obtener informacion del pago")
		}
		if response == nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  "ok",
				"data":    response,
				"message": fmt.Sprintf("no existe informacion de pago asociado al comprobante de pago %v ", requestPaymentId),
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "ok",
			"data":    response,
			"message": "informacion de pago obtenido con exito",
		})
	}
}

type NroComprobante struct {
	PaymentId string `json:"external_id"`
}

// import (
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/pagooffline"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
// 	offlinedtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/pagoofflinedtos"

// 	//"github.com/aws/aws-sdk-go/service"
// 	"github.com/gofiber/fiber/v2"
// )

// func PrismaRoutes(app fiber.Router, service pagooffline.Service) {
// 	app.Post("/probarruta", consultarPagos(service))
// }

// func consultarPagos(service pagooffline.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		var response dtos.HttpResponse
// 		var request offlinedtos.OffLineRequestResponse
// 		err := c.BodyParser(&request)
// 		if err != nil {
// 			return fiber.NewError(404, "Error: "+err.Error())
// 		}
// 		pagoOffLine_ov := offlinedtos.New(request)

// 		if len(pagoOffLine_ov.GetErrors()) > 0 {
// 			response = dtos.HttpResponse{
// 				Type:        "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1",
// 				Title:       "error en los parametros de entradas",
// 				Status:      400,
// 				Detail:      "pagoOffLine_ov.GetErrors()",
// 				Parametters: pagoOffLine_ov.GetErrors(),
// 			}
// 			return c.JSON(&fiber.Map{
// 				"data": response,
// 			})
// 		}
// 		/// desde este punto se va para el servicio de pagooffline
// 		codigoBarra, err := service.GenerarCodigoBarra(pagoOffLine_ov)
// 		if err != nil {
// 			return fiber.NewError(404, "Error: "+err.Error())
// 		}
// 		response = dtos.HttpResponse{
// 			Type: "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1",
// 			//Title:       "error en los parametros de entradas",
// 			Status: 200,
// 			//Detail:      "pagoOffLine_ov.GetErrors()",
// 			Parametters: codigoBarra,
// 		}
// 		return c.JSON(&fiber.Map{
// 			"data": response,
// 		})

// 	}

// }
