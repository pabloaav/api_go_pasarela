package entities

import "gorm.io/gorm"

type Movimientosubcuentas struct {
	gorm.Model
	SubcuentasID       uint      `json:"subcuentas_id"`
	MovimientosID      uint      `json:"movimientos_id"`
	Transferido        bool      `json:"transferido"`
	Monto              uint      `json:"monto"`
	PorcentajeAplicado uint      `json:"porcentaje_aplicado"`
	Subcuenta          Subcuenta `json:"subcuenta" gorm:"foreignKey:subcuentas_id"`
}
