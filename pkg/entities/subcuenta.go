package entities

import (
	"gorm.io/gorm"
)

type Subcuenta struct {
	gorm.Model

	Tipo                string  `json:"tipo"`
	CuentasID           uint    `json:"cuentas_id"`
	Cbu                 string  `json:"cbu"`
	Nombre              string  `json:"nombre"`
	Email               string  `json:"email"`
	Porcentaje          float64 `json:"porcentaje"`
	Cuenta              Cuenta  `json:"cuenta" gorm:"foreignKey:CuentasID"`
	AplicaPorcentaje    bool    `json:"aplica_porcentaje"`
	AplicaCostoServicio bool    `json:"aplica_costo_servicio"`
}

// TableName sobreescribe el nombre de la tabla
func (Subcuenta) TableName() string {
	return "subcuentas"
}
