package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponsePrismaMotivosAjuste struct {
	Id            int64
	ExternalId    string
	Motivoajustes string
}

func (rpma *ResponsePrismaMotivosAjuste) EntityToDtos(entityPrismaMotivoAjuste entities.Prismamotivosajuste) {
	rpma.Id = 0
	if entityPrismaMotivoAjuste.ID > 0 {
		rpma.Id = int64(entityPrismaMotivoAjuste.ID)
	}
	rpma.ExternalId = entityPrismaMotivoAjuste.ExternalId
	rpma.Motivoajustes = entityPrismaMotivoAjuste.Motivoajustes
}
