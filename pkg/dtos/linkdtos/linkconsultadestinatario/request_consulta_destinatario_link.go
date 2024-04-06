package linkconsultadestinatario

import (
	"encoding/json"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestConsultaDestinatarioLink struct {
	Cbu   string
	Alias string
}

func (r *RequestConsultaDestinatarioLink) IsValid() (erro error) {

	if len(r.Cbu) > 0 {

		erro = tools.EsCbuValido(r.Cbu, tools.ERROR_CBU)

		if erro != nil {
			return
		}
		r.Alias = ""

	} else {
		erro = tools.EsAliasCbuValido(r.Alias)

		if erro != nil {
			return
		}
	}

	return
}

func (r *RequestConsultaDestinatarioLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
