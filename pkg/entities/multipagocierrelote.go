package entities

import (
	"time"

	"gorm.io/gorm"
)

type Multipagoscierrelote struct {
	gorm.Model
	NombreArchivo         string
	IdHeader              string
	NombreEmpresa         string
	FechaProceso          string
	IdArchivo             string
	FillerHeader          string
	IdTrailer             string
	CantDetalles          int64
	ImporteTotal          int64
	FillerTrailer         string
	ImporteTotalCalculado float64
	BancoExternalId       int64
	PagoActualizado       bool
	Difbancocl            float64
	Fechaacreditacion     time.Time
	Cantdias              int
	ImporteMinimo         float64
	ImporteMaximo         float64
	Coeficiente           float64
	Enobservacion         bool
	MultipagoDetalle      []*Multipagoscierrelotedetalles `gorm:"foreignkey:MultipagoscierrelotesId"`
}
