package entities

import "gorm.io/gorm"

type Mediopagoinstallment struct {
	gorm.Model
	Nombre       string        `json:"nombre"`
	Descripcion  string        `json:"descripcion"`
	Mediopago    []Mediopago   `json:"mediopago" gorm:"foreignKey:MediopagoinstallmentsID"`
	Installments []Installment `json:"installments" gorm:"foreignKey:MediopagoinstallmentsID"`
}
