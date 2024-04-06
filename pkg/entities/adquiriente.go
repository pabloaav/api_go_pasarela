package entities

import "gorm.io/gorm"

type Adquiriente struct {
	gorm.Model
	Adquiriente string `json:"adquiriente"`
}
