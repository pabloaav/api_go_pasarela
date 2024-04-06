package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseCodigoRechazos struct {
	Id         uint   `json:"id"`
	ExternalId string `json:"external_id"`
	Rechazo    string `json:"rechado"`
}

func (rcr *ResponseCodigoRechazos) ToDtos(entityCodigoRechazos entities.Prismacodigorechazo) ResponseCodigoRechazos {
	rcr.Id = entityCodigoRechazos.ID
	rcr.ExternalId = entityCodigoRechazos.ExternalId
	rcr.Rechazo = entityCodigoRechazos.Rechazo
	return *rcr
}
