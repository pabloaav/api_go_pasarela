package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseMotivosAjustes struct {
	Id            uint   `json:"id"`
	ExternalId    string `json:"external_id"`
	Motivoajustes string `json:"motivoajustes"`
}

func (rma *ResponseMotivosAjustes) ToDtos(entityCodigoRechazos entities.Prismamotivosajuste) ResponseMotivosAjustes {
	rma.Id = entityCodigoRechazos.ID
	rma.ExternalId = entityCodigoRechazos.ExternalId
	rma.Motivoajustes = entityCodigoRechazos.Motivoajustes
	return *rma
}
