package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseMovimientosubcuenta struct {
	Subcuenta  string `json:"subcuenta"`
	Monto      uint   `json:"monto"`
	Porcentaje uint   `json:"porcentaje"`
}

func (r *ResponseMovimientosubcuenta) ToEntity(mov entities.Movimientosubcuentas) {
	r.Monto = mov.Monto
	r.Subcuenta = mov.Subcuenta.Nombre
	r.Porcentaje = mov.PorcentajeAplicado
}
