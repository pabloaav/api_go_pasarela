package entities

import "gorm.io/gorm"

type Prismacodigorechazo struct {
	gorm.Model
	ExternalId string
	Rechazo    string
}
