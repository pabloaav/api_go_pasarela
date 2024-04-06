package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseOperaciones struct {
	Id         uint   `json:"id"`
	ExternalId string `json:"external_id"`
	Operacion  string `json:"operacion"`
}

func (ro *ResponseOperaciones) ToDtos(entityCodigoRechazos entities.Prismaoperacion) ResponseOperaciones {
	ro.Id = entityCodigoRechazos.ID
	ro.ExternalId = entityCodigoRechazos.ExternalId
	ro.Operacion = entityCodigoRechazos.Operacion
	return *ro
}
