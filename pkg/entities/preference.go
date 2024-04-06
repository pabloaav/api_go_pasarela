package entities

import (
	"gorm.io/gorm"
)

type Preference struct {
	gorm.Model
	ClientesId     uint `json:"clientes_id"`
	Maincolor      string
	Secondarycolor string
	Logo           string
	Cliente        Cliente `json:"cliente" gorm:"foreignKey:ClientesId"`
}
