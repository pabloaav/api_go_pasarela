package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponsePrismaMasterContracargo struct {
	Id          int64
	ExternalId  string
	Contracargo string
}

func (rpmc *ResponsePrismaMasterContracargo) EntityToDtos(entityPrismaMasterContracargo entities.Prismamastercontracargo) {
	rpmc.Id = 0
	if entityPrismaMasterContracargo.ID > 0 {
		rpmc.Id = int64(entityPrismaMasterContracargo.ID)
	}
	rpmc.ExternalId = entityPrismaMasterContracargo.ExternalId
	rpmc.Contracargo = entityPrismaMasterContracargo.Contracargo
}
