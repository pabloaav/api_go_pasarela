package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseVisaContracargo struct {
	Id          uint   `json:"id"`
	ExternalId  string `json:"external_id"`
	Contracargo string `json:"contracargo"`
}

func (rvc *ResponseVisaContracargo) ToDtos(entityCodigoRechazos entities.Prismavisacontracargo) ResponseVisaContracargo {
	rvc.Id = entityCodigoRechazos.ID
	rvc.ExternalId = entityCodigoRechazos.ExternalId
	rvc.Contracargo = entityCodigoRechazos.Contracargo
	return *rvc
}
