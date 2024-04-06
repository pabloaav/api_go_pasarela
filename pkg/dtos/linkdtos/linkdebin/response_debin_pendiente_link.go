package linkdebin

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
)

type ResponseGetDebinesPendientesLink struct {
	Debines []DebinesPendientesLink `json:"debines"`
}

type DebinesPendientesLink struct {
	Id              string                        `json:"id"`
	Concepto        linkdtos.EnumConceptoDebin    `json:"concepto"`
	Moneda          linkdtos.EnumMoneda           `json:"moneda"`
	Importe         int64                         `json:"importe"`
	Estado          linkdtos.EnumEstadoDebin      `json:"estado"`
	FechaExpiracion time.Time                     `json:"fechaexpiracion"`
	CbuComprador    string                        `json:"cbucomprador"`
	Vendedor        VendedorDebinesPendientesLink `json:"vendedor"`
}

type VendedorDebinesPendientesLink struct {
	Cuit    string `json:"cuit"`
	Titular string `json:"titularcuenta"` //Titular de la Cuenta
}
