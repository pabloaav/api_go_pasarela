package entities

import (
	"gorm.io/gorm"
)

type Channelarancele struct {
	gorm.Model
	ChannelsId    uint    `json:"channels_id"`
	RubrosId      int64   `json:"rubros_id"`
	Importe       float64 `json:"importe"`
	Fechadesde    string  `json:"fecha_desde"`
	Tipocalculo   string  `json:"tipocalculo"`
	Importeminimo float64 `json:"importeminimo"`
	Importemaximo float64 `json:"importemaximo"`
	Mediopagoid   int64   `json:"mediopagoid"`
	Pagocuota     bool    `json:"pagocuota"`
	Channel       Channel `json:"channel" gorm:"foreignKey:channels_id"`
	Rubro         Rubro   `json:"rubro" gorm:"foreignKey:rubros_id"`
}
