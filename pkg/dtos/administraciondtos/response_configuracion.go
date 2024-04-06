package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseConfiguraciones struct {
	Data []ResponseConfiguracion `json:"data"`
	Meta dtos.Meta               `json:"meta"`
}

type ResponseConfiguracion struct {
	Id          uint
	Nombre      string
	Descripcion string
	Valor       string
}

func (r *ResponseConfiguracion) FromEntity(c entities.Configuracione) {
	r.Id = c.ID
	r.Nombre = c.Nombre
	r.Descripcion = c.Descripcion
	r.Valor = c.Valor
}
