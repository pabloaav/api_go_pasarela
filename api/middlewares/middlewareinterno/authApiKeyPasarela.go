package middlewareinterno

import (
	"errors"
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareManagerPasarela struct {
	Service administracion.Service
}

func (m *MiddlewareManagerPasarela) ValidarApiKeyCliente() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		apikey := c.Get("apikey")
		if len(apikey) <= 0 {
			return errors.New("Debe enviar una api key válida")
		}
		isApiKeyValid, err := m.Service.GetCuentaByApiKeyService(apikey)
		if err != nil {
			return fmt.Errorf("error " + err.Error())
		}
		if !isApiKeyValid {
			return fmt.Errorf("api-key no es valida")
		}
		return c.Next()
	}
}
