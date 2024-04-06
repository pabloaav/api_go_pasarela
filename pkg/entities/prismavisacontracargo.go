package entities

import "gorm.io/gorm"

type Prismavisacontracargo struct {
	gorm.Model
	ExternalId  string
	Contracargo string
}
