package entities

import "gorm.io/gorm"

type Movimientoimpuestos struct {
	gorm.Model
	MovimientosID  uint64      `json:"movimientos_id"`
	ImpuestosID    uint64      `json:"impuestos_id"`
	Monto          Monto       `json:"monto"`
	Montoproveedor Monto       `json:"montoproveedor"`
	Porcentaje     float64     `json:"pocentaje"`
	Movimiento     *Movimiento `json:"movimiento" gorm:"foreignKey:movimientos_id"`
	Impuesto       *Impuesto   `json:"impuesto" gorm:"foreignKey:impuestos_id"`
}
