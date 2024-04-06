package entities

import "gorm.io/gorm"

type Contactosreporte struct {
	gorm.Model
	Email      string
	ClientesID int64 `json:"clientes_id"`
}
