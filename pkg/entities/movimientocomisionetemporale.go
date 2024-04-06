package entities

import "gorm.io/gorm"

type Movimientocomisionetemporale struct {
	gorm.Model
	MovimientotemporalesID uint64               `json:"movimientotemporales_id"`
	CuentacomisionsID      uint                 `json:"cuenta_comisiones_id"`
	Monto                  Monto                `json:"monto"`
	Montoproveedor         Monto                `json:"montoproveedor"`
	Porcentaje             float64              `json:"porcentaje"`
	Porcentajeproveedor    float64              `json:"porcentajeproveedor"`
	Movimientotemporale    *Movimientotemporale `json:"movimientotemporales" gorm:"foreignKey:movimientotemporales_id"`
	Cuentacomisions        *Cuentacomision      `json:"cuenta_comision" gorm:"foreignKey:cuentacomisions_id"`
}
