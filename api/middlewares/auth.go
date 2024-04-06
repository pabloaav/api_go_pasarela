package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareManager struct {
	HTTPClient *http.Client
}

func (m *MiddlewareManager) ValidarPermiso(scope string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if len(bearer) <= 0 {
			return errors.New("acceso no autorizado, debe enviar un token de autenticación")
		}

		var result struct {
			Acceso string `json:"acceso"`
			ID     int64  `json:"user_id"`
		}

		base, err := url.Parse(config.AUTH)
		if err != nil {
			return fmt.Errorf("error al crear base url" + err.Error())
		}

		base.Path += "/users/permiso"

		var values struct {
			SistemaID int64  `json:"sistema_id"`
			Scope     string `json:"scope"`
		}
		//FIXME definir como vamos a setear el sistemaId
		idSitema, err := strconv.ParseInt(config.ID_SISTEMAS, 10, 64)
		if err != nil {
			return fmt.Errorf("error al convertir id sistema")
		}
		values.SistemaID = idSitema
		values.Scope = scope
		json_data, _ := json.Marshal(values)

		req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

		req.Header.Add("Authorization", bearer)
		req.Header.Add("Content-Type", "application/json")

		resp, err := m.HTTPClient.Do(req)

		if err != nil {
			return fmt.Errorf("error al enviar solicitud a api externa")
		}

		if resp.StatusCode != 200 {
			info, _ := ioutil.ReadAll(req.Body)
			erro := fmt.Errorf("acceso denegado o permisos insuficientes: %s", info)
			return fiber.NewError(403, erro.Error())
		}

		json.NewDecoder(resp.Body).Decode(&result)
		c.Set("user_id", fmt.Sprint(result.ID))

		return c.Next()
	}
}

func (m *MiddlewareManager) ValidarApikey() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		apikey := c.Get("ApiKey")
		if len(apikey) <= 0 {
			return fmt.Errorf("acceso no autorizado, debe enviar una apikey para autenticación")
		}

		var result struct {
			validez bool
		}

		base, err := url.Parse(config.AUTH)
		if err != nil {
			return fmt.Errorf("error al crear base url" + err.Error())
		}

		base.Path += "/users/validar-apikey"

		var values struct {
			// SistemaID uint `json:"sistema_id"`
		}
		sistema_id := config.ID_SISTEMAS
		// sistema_iduint64, _ := strconv.ParseUint(sistema_id, 10, 64)
		// values.SistemaID = uint(sistema_iduint64)
		json_data, _ := json.Marshal(values)

		req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("ApiKey", apikey)
		req.Header.Set("SistemaId", sistema_id)

		resp, err := m.HTTPClient.Do(req)

		if err != nil {
			return fiber.NewError(400, "Error en los datos enviados")
		}

		if resp.StatusCode != 200 {
			erro := fmt.Errorf("Api key invalida para el sistema.")
			return fiber.NewError(403, erro.Error())
		}

		json.NewDecoder(resp.Body).Decode(&result)
		c.Set("Valido", fmt.Sprint(result.validez))

		return c.Next()
	}
}
