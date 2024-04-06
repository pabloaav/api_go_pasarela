package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseArchivoSubido struct {
	ArchivosSubidos []ArchivoSubido `json:"data"`
	Meta            dtos.Meta       `json:"meta"`
}

type ArchivoSubido struct {
	NombreArchivo string
	FechaSubida   time.Time
}

func (rcs *ArchivoSubido) EntityClToDtos(entityCl *entities.Prismacierrelote) {
	rcs.NombreArchivo = entityCl.Nombrearchivolote
	rcs.FechaSubida = entityCl.CreatedAt
}

func (rcs *ArchivoSubido) EntityPxToDtos(entityPx *entities.Prismapxcuatroregistro) {
	rcs.NombreArchivo = entityPx.Nombrearchivo
	rcs.FechaSubida = entityPx.CreatedAt
}

func (rcs *ArchivoSubido) EntityMxToDtos(entityMx *entities.Prismamxtotalesmovimiento) {
	rcs.NombreArchivo = entityMx.Nombrearchivo
	rcs.FechaSubida = entityMx.CreatedAt
}
