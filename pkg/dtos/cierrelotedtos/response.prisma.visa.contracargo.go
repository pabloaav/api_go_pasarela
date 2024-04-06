package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponsePrismaVisaContracargo struct {
	Id          int64
	ExternalId  string
	Contracargo string
}

func (rpvc *ResponsePrismaVisaContracargo) EntityToDtos(entityPrismaVisaContracargo entities.Prismavisacontracargo) {
	rpvc.Id = 0
	if entityPrismaVisaContracargo.ID > 0 {
		rpvc.Id = int64(entityPrismaVisaContracargo.ID)
	}
	rpvc.ExternalId = entityPrismaVisaContracargo.ExternalId
	rpvc.Contracargo = entityPrismaVisaContracargo.Contracargo

}
