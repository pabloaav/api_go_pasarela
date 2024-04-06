package entities

import (
	"time"

	"gorm.io/gorm"
)

type Impuesto struct {
	gorm.Model
	Impuesto            string                `json:"impuesto"`
	Porcentaje          float64               `json:"porcentaje"`
	Tipo                string                `json:"tipo"`
	Fechadesde          time.Time             `json:"fechadesde"`
	Activo              bool                  `json:"activo"`
	Movimientoimpuestos []Movimientoimpuestos `json:"movimiento_impuestos" gorm:"foreignKey:ImpuestosID"`
}
