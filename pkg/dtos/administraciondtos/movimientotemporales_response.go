package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type MovimientoTemporalesResponse struct {
	ListaPagosCalculado []uint                         `json:"pagointentos,omitempty"`
	ListaMovimientos    []entities.Movimientotemporale `json:"moviminetotemporales,omitempty"`
}
