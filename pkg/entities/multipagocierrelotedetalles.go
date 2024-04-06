package entities

import (
	"gorm.io/gorm"
)

type Multipagoscierrelotedetalles struct {
	gorm.Model
	MultipagoscierrelotesId int64 `json:"multipagoscierrelotes_id"`
	FechaCobro              string
	ImporteCobrado          int64
	CodigoBarras            string
	ImporteCalculado        float64
	Match                   bool
	Clearing                string
	Enobservacion           bool
	Pagoinformado           bool
	MultipagoCabecera       Multipagoscierrelote `json:"multipagoscierrelotes" gorm:"foreignKey:MultipagoscierrelotesId"`
}
