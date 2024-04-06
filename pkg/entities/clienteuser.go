package entities

import "gorm.io/gorm"

type Clienteuser struct {
	gorm.Model
	ClientesId uint64
	UserId     uint64
	Cliente    Cliente `json:"cliente" gorm:"foreignKey:ClientesId"`
}
