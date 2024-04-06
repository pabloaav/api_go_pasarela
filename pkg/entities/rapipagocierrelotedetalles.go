package entities

import (
	"gorm.io/gorm"
)

type Rapipagocierrelotedetalles struct {
	gorm.Model
	RapipagocierrelotesId int64 `json:"rapipagocierrelotes_id"`
	FechaCobro            string
	ImporteCobrado        int64
	CodigoBarras          string
	ImporteCalculado      float64
	Match                 bool
	Clearing              string
	Enobservacion         bool
	Pagoinformado         bool
	// RapipagoCabecera      Rapipagocierrelote `gorm:"foreignkey:rapipagocierrelotes_id"`
	RapipagoCabecera Rapipagocierrelote `json:"rapipagocierrelotes" gorm:"foreignKey:RapipagocierrelotesId"`
}
