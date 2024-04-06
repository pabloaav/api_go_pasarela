package linktransferencia

import (
	"encoding/json"
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestGetTransferenciaLink struct {
	NumeroReferenciaBancaria string `json:"numeroReferenciaBancaria"`
	Cbu                      string `json:"cbu"`
}

func (r *RequestGetTransferenciaLink) IsValid() error {
	err := tools.EsCbuValido(r.Cbu, tools.ERROR_CBU)
	if err != nil {
		return err
	}
	if len(r.NumeroReferenciaBancaria) != 30 {
		return errors.New(tools.ERROR_REFERENCIA_BANCARIA)
	}

	return nil
}

func (r *RequestGetTransferenciaLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
