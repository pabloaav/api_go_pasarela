package entities

import "gorm.io/gorm"

type Prismaoperacion struct {
	gorm.Model
	ExternalId string
	Operacion  string
}
