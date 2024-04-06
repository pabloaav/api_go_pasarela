package entities

import (
	"gorm.io/gorm"
)

type Movimientolotes struct {
	gorm.Model
	MovimientosID uint64      `json:"movimientos_id"`
	ClientesID    uint64      `json:"clientes_id"`
	Lote          int64       `json:"lote"`
	FechaEnvio    string      `json:"fecha_envio"`
	MotivoBaja    string      `json:"motivi_baja"`
	Movimiento    *Movimiento `json:"movimiento" gorm:"foreignKey:movimientos_id"`
	Cliente       *Cliente    `json:"impuesto" gorm:"foreignKey:clientes_id"`
}
