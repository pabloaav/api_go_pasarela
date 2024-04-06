package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

func CheckoutRoutes(app fiber.Router, service checkout.Service) {
	app.Get("/:barcode", getCheckout(service)) // devuelve el checkout
	app.Post("/", newCheckout(service))        // genra la solicitud de pago
	app.Post("/pagar", postPagar(service))     // realiza el pago
	app.Post("/prisma", checkPrisma(service))
	app.Get("/bill/:barcode", getBill(service)) // pdf comprobante del pago
	app.Get("/tarjetas/all", getTarjetas(service))
	app.Get("/verificar/pago/:barcode", getVerificarPago(service))
	app.Get("/verificar/pago/estado/:barcode", getVerificarPagoEstado(service))

	// app.Get("/control/:hash", controlHash(service))
}

func getCheckout(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		code := c.Params("barcode")

		response, err := service.GetPago(code)
		if err != nil {
			return c.Render("error", &fiber.Map{"message": err.Error()})
		}

		response.BaseUrl = config.APP_HOST

		byteResponse, err := json.Marshal(response)
		if err != nil {
			logs.Error(err.Error())
		}
		var mapResponse interface{}
		err = json.Unmarshal(byteResponse, &mapResponse)
		if err != nil {
			logs.Error(err.Error())
		}
		logs.Info(mapResponse)

		return c.Render("index", mapResponse)
	}
}

func newCheckout(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}

		var request dtos.PagoRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Parámetros incorrectos: "+err.Error())
		}

		estado, fecha, err := service.GetMatenimietoSistema()

		if err != nil {
			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"message": "el sistema estara en mantenimiento hasta " + time.Now().Format(time.RFC822Z),
			})
			// return fiber.NewError(503, "Parámetros incorrectos: "+err.Error())
		}
		logs.Error(fmt.Sprintf("error fecha:%v ", fecha))
		if estado {
			fechaString := fmt.Sprintf("%v-%v-%v %v:%v:%v", fecha.Day(), fecha.Month(), fecha.Year(), fecha.Hour(), fecha.Minute(), fecha.Second())

			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"message": "el sistema estara en mantenimiento hasta " + fechaString,
			})
		}

		ctx := getCheckoutContext(c)

		res, err := service.NewPago(ctx, &request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    res,
			"message": "Solicitud de pago generada",
		})
	}
}

func postPagar(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		logs.Info(fmt.Sprint("postPagar: ", c.IP()))

		var request dtos.ResultadoRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Parámetros incorrectos: "+err.Error())
		}
		request.Ip = c.IP()
		logs.Info(request.Channel)
		logs.Info(request.Uuid)
		///// validar si el pago ya

		ctx := getCheckoutContext(c)

		res, err := service.GetPagoResultado(ctx, &request)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		return c.JSON(&fiber.Map{
			"status":         res.Exito,
			"data":           res,
			"statusMessage":  res.Estado,
			"successMesagge": res.Mensaje,
		})
	}
}

func checkPrisma(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		err := service.CheckPrisma()
		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "el servicio de prisma está funcionando correctamente.",
		})
	}
}

// Retorna un stream para mostrar un pdf con los datos del pago
func getBill(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		code := c.Params("barcode")

		file, err := service.GetBilling(code)
		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		c.Set("Content-Disposition", "filename=recibo.pdf")
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Length", fmt.Sprint(file.Len()))

		return c.SendStream(file) // c.SendStream(file)
	}
}

func getTarjetas(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		tar, err := service.GetTarjetas()
		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    tar,
			"message": "tarjetas enviadas",
		})
	}
}

func getVerificarPagoEstado(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		code := c.Params("barcode")

		_, err := service.GetPagoStatus(code) // servicio que verifica el estado del pago
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		response, err := service.GetPago(code)
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		logs.Info(response)
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "puede realizar pago",
		})
	}
}

func getVerificarPago(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		code := c.Params("barcode")
		response, err := service.GetPago(code)
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		logs.Info(response)
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "puede realizar pago",
		})
	}
}

func controlHash(service checkout.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		hash := c.Params("hash")

		response, err := service.ControlTarjetaHash(hash)
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    response,
			"message": "Control tarjeta bloqueada",
		})
	}
}

func getCheckoutContext(c *fiber.Ctx) context.Context {
	userctx := entities.Auditoria{
		IP: c.IP(),
	}
	ctx := context.WithValue(c.Context(), entities.AuditUserKey{}, userctx)
	return ctx
}
