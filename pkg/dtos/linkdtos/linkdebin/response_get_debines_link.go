package linkdebin

import (
	"encoding/json"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
)

type ResponseGetDebinesLink struct {
	Paginado linkdtos.PaginadoResponseLink `json:"paginado"`
	Debines  []DebinesListLink             `json:"debines"`
}

type DebinesListLink struct {
	Id              string                     `json:"id"`
	Concepto        linkdtos.EnumConceptoDebin `json:"concepto"`
	Moneda          linkdtos.EnumMoneda        `json:"moneda"`
	Importe         int64                      `json:"importe"`
	Estado          linkdtos.EnumEstadoDebin   `json:"estado"`
	Tipo            linkdtos.EnumTipoDebin     `json:"tipo"`
	FechaExpiracion time.Time                  `json:"fechaExpiracion"`
	Devuelto        bool                       `json:"devuelto"`
	ContraCargoId   string                     `json:"contracargoId"`
	Comprador       CompradorDebinesListLink   `json:"comprador"`
	Vendedor        VendedorDebinesListLink    `json:"vendedor"`
}

type VendedorDebinesListLink struct {
	Cuit    string `json:"cuit"`
	Titular string `json:"titular"` //Titular de la Cuenta
}

type CompradorDebinesListLink struct {
	Cuit    string `json:"cuit"`
	Titular string `json:"titular"` //Titular de la Cuenta
}

func (r *ResponseGetDebinesLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
