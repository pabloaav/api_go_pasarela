package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseImpuestos struct {
	Impuestos []ResponseImpuesto `json:"data"`
	Meta      dtos.Meta          `json:"meta"`
}

type ResponseImpuesto struct {
	Id         uint      `json:"id"`
	Impuesto   string    `json:"impuesto"`
	Porcentaje float64   `json:"porcentaje"`
	Tipo       string    `json:"tipo"`
	FechaDesde time.Time `json:"fecha_desde"`
}

func (r *ResponseImpuesto) FromImpuesto(entity entities.Impuesto) {
	r.Id = entity.ID
	r.Impuesto = entity.Impuesto
	r.Porcentaje = entity.Porcentaje
	r.Tipo = entity.Tipo
	r.FechaDesde = entity.Fechadesde
}
