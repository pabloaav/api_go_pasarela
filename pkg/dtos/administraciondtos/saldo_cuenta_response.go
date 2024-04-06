package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type SaldoCuentaResponse struct {
	CuentasId uint64         `json:"cuentas_id"`
	Total     entities.Monto `json:"total"`
}
