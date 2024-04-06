package entities

import "gorm.io/gorm"

type Movimientoimpuestotemporale struct {
	gorm.Model
	MovimientotemporalesID uint64               `json:"movimientotemporales_id"`
	ImpuestosID            uint64               `json:"impuestos_id"`
	Monto                  Monto                `json:"monto"`
	Montoproveedor         Monto                `json:"montoproveedor"`
	Porcentaje             float64              `json:"pocentaje"`
	Movimientotemporale    *Movimientotemporale `json:"movimientotemporales" gorm:"foreignKey:movimientotemporales_id"`
	Impuesto               *Impuesto            `json:"impuesto" gorm:"foreignKey:impuestos_id"`
}
