package entities

import (
	"gorm.io/gorm"
)

type Movimientoliquidaciones struct {
	gorm.Model
	ClientesID           uint64                `json:"clientes_id"`
	FechaEnvio           string                `json:"fecha_envio"`
	MotivoBaja           string                `json:"motivi_baja"`
	Cliente              *Cliente              `json:"impuesto" gorm:"foreignKey:clientes_id"`
	LiquidacioneDetalles []Liquidaciondetalles `json:"liquidaciondetalles" gorm:"foreignKey:movimientoliquidaciones_id"`
}
