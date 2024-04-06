package entities

import "gorm.io/gorm"

type Gravamen struct {
	gorm.Model
	Gravamen       string
	CodigoGravamen string
	Descripcion    string
	Condiciones    []Condicion `gorm:"foreignKey:GravamensId"`
}
