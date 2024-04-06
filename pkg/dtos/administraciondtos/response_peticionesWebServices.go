package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePeticionesWebServices struct {
	Meta       dtos.Meta                     `json:"meta"`
	Peticiones []ResponsePeticionWebServices `json:"data"`
}

type ResponsePeticionWebServices struct {
	Operacion string
	Vendor    string
}

func (r *ResponsePeticionWebServices) SetPeticion(peticion entities.Webservicespeticione) {
	r.Operacion = peticion.Operacion
	r.Vendor = string(peticion.Vendor)
}
