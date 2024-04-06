package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponsePrismaOperacion struct {
	Id         int64
	ExternalId string
	Operacion  string
}

func (rpo *ResponsePrismaOperacion) EntityToDtos(entityPrismaOperacion entities.Prismaoperacion) {
	rpo.Id = 0
	if entityPrismaOperacion.ID > 0 {
		rpo.Id = int64(entityPrismaOperacion.ID)
	}
	rpo.ExternalId = entityPrismaOperacion.ExternalId
	rpo.Operacion = entityPrismaOperacion.Operacion
}
