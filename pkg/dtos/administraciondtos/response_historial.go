package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseHistorial struct {
	Historial []RegistroHistorial `json:"historial"`
	Meta      dtos.Meta           `json:"meta"`
}

type RegistroHistorial struct {
	Id            uint                   `json:"id"`
	UserId        uint                   `json:"users_id"`
	Correo        string                 `json:"correo"`
	Operacion     entities.EnumHistorial `json:"tipo_operacion"`
	Observaciones string                 `json:"observaciones"`
}

func (r *RegistroHistorial) FromEntity(h entities.HistorialOperaciones) {
	r.Id = h.ID
	r.UserId = uint(h.UsersId)
	r.Correo = h.Correo
	r.Operacion = h.TipoOperacion
	r.Observaciones = h.Observaciones
}

func (rh *ResponseHistorial) FromEntities(historial []entities.HistorialOperaciones) {
	for _, h := range historial {
		rh.Historial = append(rh.Historial, RegistroHistorial{
			Id:            h.ID,
			UserId:        uint(h.UsersId),
			Correo:        h.Correo,
			Operacion:     h.TipoOperacion,
			Observaciones: h.Observaciones,
		})
	}
	return
}
