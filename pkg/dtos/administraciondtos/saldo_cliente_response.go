package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type SaldoClienteResponse struct {
	ClienteId uint64         `json:"cliente_id"`
	Total     entities.Monto `json:"total"`
}
