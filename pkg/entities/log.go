package entities

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	UserId        uint    `json:"user_id"`
	Ip            string  `json:"ip"`
	Tipo          EnumLog `json:"tipo"`
	Funcionalidad string  `json:"funcionalidad"`
	Mensaje       string  `json:"mensaje"`
}

type EnumLog string

const (
	Info    EnumLog = "info"
	Warning EnumLog = "warning"
	Error   EnumLog = "error"
)
