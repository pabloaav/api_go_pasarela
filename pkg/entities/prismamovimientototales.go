package entities

import (
	"time"

	"gorm.io/gorm"
)

type Prismamovimientototale struct {
	gorm.Model
	Empresa            string
	FechaPresentacion  time.Time
	TipoRegistro       string
	ComercioNro        string
	EstablecimientoNro string
	Codop              string
	TipoAplicacion     string
	FechaPago          time.Time
	ImporteTotal       Monto
	SignoImporteTotal  string
	Match              int
	DetalleMovimientos []Prismamovimientodetalle `gorm:"foreignkey:PrismamovimientototalesId"`
}
