package entities

import (
	"gorm.io/gorm"
)

type Liquidaciondetalles struct {
	gorm.Model
	MovimientoliquidacionesId int64 `json:"movimientoliquidaciones_id"`
	PagointentosId            int64
	MovimientosId             int64
	CuentasId                 int64
	Movimientoliquidacion     Movimientoliquidaciones `json:"movimientoliquidaciones" gorm:"foreignKey:MovimientoliquidacionesId"`
}
