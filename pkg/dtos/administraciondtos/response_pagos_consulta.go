package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagosConsulta struct {
	Uuid              string    `json:"uuid"`
	ExternalReference string    `json:"external_reference"`
	Metadata          string    `json:"metadata"`
	Fecha             time.Time `json:"fecha"`
	Estado            string    `json:"estado"`
}

func (r *ResponsePagosConsulta) SetPago(pago entities.Pago) {
	r.Uuid = pago.Uuid
	r.ExternalReference = pago.ExternalReference
	r.Metadata = pago.Metadata
	r.Fecha = pago.UpdatedAt //.Format("02-01-2006")
	r.Estado = string(pago.PagoEstados.Estado)
}
