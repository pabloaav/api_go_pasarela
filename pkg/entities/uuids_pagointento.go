package entities

import (
	"gorm.io/gorm"
)

type UuidsPagointento struct {
	gorm.Model
	UuidsId        uint         `json:"uuids_id"`
	Uuid           *Uuid        `json:"uuid" gorm:"foreignKey:uuids_id"`
	PagointentosId uint         `json:"pagointentos_id"`
	Pagointento    *Pagointento `json:"pagointento" gorm:"foreignKey:pagointentos_id"`
}
