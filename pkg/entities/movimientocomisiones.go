package entities

import "gorm.io/gorm"

type Movimientocomisiones struct {
	gorm.Model
	MovimientosID       uint64          `json:"movimientos_id"`
	CuentacomisionsID   uint            `json:"cuenta_comisiones_id"`
	Monto               Monto           `json:"monto"`
	Montoproveedor      Monto           `json:"montoproveedor"`
	Porcentaje          float64         `json:"porcentaje"`
	Porcentajeproveedor float64         `json:"porcentajeproveedor"`
	Movimiento          *Movimiento     `json:"movimiento" gorm:"foreignKey:movimientos_id"`
	Cuentacomisions     *Cuentacomision `json:"cuenta_comision" gorm:"foreignKey:cuentacomisions_id"`
}
