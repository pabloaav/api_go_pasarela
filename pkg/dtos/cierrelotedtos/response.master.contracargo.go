package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseMasterContracargo struct {
	Id          uint   `json:"id"`
	ExternalId  string `json:"external_id"`
	Contracargo string `json:"contracargo"`
}

func (rmc *ResponseMasterContracargo) ToDtos(entityCodigoRechazos entities.Prismamastercontracargo) ResponseMasterContracargo {
	rmc.Id = entityCodigoRechazos.ID
	rmc.ExternalId = entityCodigoRechazos.ExternalId
	rmc.Contracargo = entityCodigoRechazos.Contracargo
	return *rmc
}
