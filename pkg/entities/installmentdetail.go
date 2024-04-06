package entities

import (
	"time"

	"gorm.io/gorm"
)

type Installmentdetail struct {
	gorm.Model
	InstallmentsID uint        `json:"installments_id"`
	Cuota          int64       `json:"cuota"`
	Tna            float64     `json:"tna"`
	Tem            float64     `json:"tem"`
	Coeficiente    float64     `json:"coeficiente"`
	Activo         bool        `json:"activo"`
	Fechadesde     time.Time   `json:"fechadesde"`
	Installment    Installment `json:"installmentdetail" gorm:"foreignKey:InstallmentsID"`
}
