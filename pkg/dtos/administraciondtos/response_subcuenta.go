package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseSubcuenta struct {
	Id         uint    `json:"id"`
	CuentasID  uint    `json:"cuentas_id"`
	Nombre     string  `json:"nombre"`
	Email      string  `json:"email"`
	Tipo       string  `json:"tipo"`
	Porcentaje float64 `json:"porcentaje"`
	Cbu        string  `json:"cbu"`
}

func (r *ResponseSubcuenta) FromSubcuenta(c entities.Subcuenta) {
	r.Id = c.ID
	r.CuentasID = c.CuentasID
	r.Tipo = c.Tipo
	r.Nombre = c.Nombre
	r.Email = c.Email
	r.Porcentaje = c.Porcentaje
	r.Cbu = c.Cbu
}
