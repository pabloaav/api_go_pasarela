package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"

type ResponseRubros struct {
	Rubros []ResponseRubro `json:"data"`
	Meta   dtos.Meta       `json:"meta"`
}

type ResponseRubro struct {
	Id    uint   `json:"id"`
	Rubro string `json:"rubro"`
}
