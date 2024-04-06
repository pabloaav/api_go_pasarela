package linkcuentas

import (
	"encoding/json"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type LinkCuentasRequest struct {
	Token           linkdtos.TokenLink
	RequerimientoId string
	Request         interface{}
}

type LinkGetCuentasRequest struct {
	Token           linkdtos.TokenLink
	RequerimientoId string
}

type LinkDeleteCuenta struct {
	Cbu string `json:"cbu"`
}

func (r *LinkDeleteCuenta) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}

func (c *LinkDeleteCuenta) IsValid() (erro error) {

	erro = tools.EsCbuValido(c.Cbu, tools.ERROR_CBU)

	if erro != nil {
		return
	}

	return
}

type LinkPostCuenta struct {
	Cuit string `json:"cuit"`
	Cbu  string `json:"cbu"`
}

func (r *LinkPostCuenta) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}

func (c *LinkPostCuenta) IsValid() (erro error) {

	erro = tools.EsCbuValido(c.Cbu, tools.ERROR_CBU)

	if erro != nil {
		return
	}

	erro = tools.EsCuitValido(c.Cuit)

	if erro != nil {
		return
	}

	return

}
