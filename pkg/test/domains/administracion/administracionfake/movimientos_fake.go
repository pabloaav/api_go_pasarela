package administracionfake

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

var estructuraValidasMovimientos []entities.Movimiento = []entities.Movimiento{
	{
		//  gorm.Model{ID:0,}	,
		CuentasId:      0,
		PagointentosId: 7,
		Tipo:           "D",
		Monto:          66051,
		MotivoBaja:     "",
	},
	{
		CuentasId:      0,
		PagointentosId: 8,
		Tipo:           "D",
		Monto:          76051,
		MotivoBaja:     "",
	},
	{
		CuentasId:      0,
		PagointentosId: 11,
		Tipo:           "D",
		Monto:          76051,
		MotivoBaja:     "",
	},
	{
		CuentasId:      0,
		PagointentosId: 12,
		Tipo:           "C",
		Monto:          -76051,
		MotivoBaja:     "",
	},
}
