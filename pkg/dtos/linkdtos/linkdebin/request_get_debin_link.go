package linkdebin

import (
	"encoding/json"
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestGetDebinLink struct {
	Cbu string `json:"cbu"` //cbu vendedor o comprador
	Id  string `json:"id"`  //Id del debin
}

func (c *RequestGetDebinLink) IsValid() error {

	err := tools.EsCbuValido(c.Cbu, tools.ERROR_CBU)
	if err != nil {
		return err
	}

	if tools.EsStringVacio(c.Id) {
		return errors.New(tools.ERROR_ID)
	}

	return nil
}

func (r *RequestGetDebinLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}