package entities

import "gorm.io/gorm"

type Prismamastercontracargo struct {
	gorm.Model
	ExternalId  string
	Contracargo string
}
