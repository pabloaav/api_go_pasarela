package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseMovimientoMxTotal struct {
	Id                 uint
	Empresa            string
	FechaPresentacion  time.Time
	TipoRegistro       string
	ComercioNro        string
	EstablecimientoNro string
	Codop              string
	TipoAplicacion     string
	FechaPago          time.Time
	ImporteTotal       entities.Monto
	SignoImporteTotal  string
}
