package entities

import "gorm.io/gorm"

type Condicion struct {
	gorm.Model
	Condicion   string   `json:"condicion"`
	GravamensId uint     `json:"gravamens_id"`
	Gravamen    Gravamen `gorm:"foreignKey:GravamensId"`
	Descripcion string   `json:"descripcion"`
	Exento      bool     `json:"exento"`
}
