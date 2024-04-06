package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponsePrismaCodigoRechazo struct {
	Id         int64
	ExternalId string
	Rechazo    string
}

func (rpcr *ResponsePrismaCodigoRechazo) EntityToDtos(entityPrismaCodigoRechazo entities.Prismacodigorechazo) {
	rpcr.Id = 0
	if entityPrismaCodigoRechazo.ID > 0 {
		rpcr.Id = int64(entityPrismaCodigoRechazo.ID)
	}
	rpcr.ExternalId = entityPrismaCodigoRechazo.ExternalId
	rpcr.Rechazo = entityPrismaCodigoRechazo.Rechazo
}
