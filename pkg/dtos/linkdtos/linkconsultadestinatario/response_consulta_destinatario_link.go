package linkconsultadestinatario

import "encoding/json"

type ResponseConsultaDestinatarioLink struct {
	Titulares       []TitularLink
	EntidadBancaria EntidadBancariaLink
	Cbu             string
}

type EntidadBancariaLink struct {
	Nombre          string
	NombreAbreviado string
}

type TitularLink struct {
	IdTributario string
	Denominacion string
}

func (r *ResponseConsultaDestinatarioLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
