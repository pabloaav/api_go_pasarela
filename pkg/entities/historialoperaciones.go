package entities

import (
	"gorm.io/gorm"
)

type EnumHistorial string

const (
	Crear      EnumHistorial = "crear"
	Actualizar EnumHistorial = "actualizar"
	Eliminar   EnumHistorial = "eliminar"
)

type HistorialOperaciones struct {
	gorm.Model
	UsersId       int           `json:"users_id"`
	Correo        string        `json:"correo"`
	TipoOperacion EnumHistorial `json:"tipo_operacion"`
	Observaciones string        `json:"observaciones"`
}
