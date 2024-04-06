package filtros

import (
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
)

const (
	layoutISO = "2006-01-02"
)

type ApilinkCierreloteFiltro struct {
	FechaInicio string // query:FechaInicio db:created_at
	FechaFin    string // query:FechaFin db:created_at
	Number      uint32 // query:Number
	Size        uint32 // query:Size
}

func (aclf *ApilinkCierreloteFiltro) Validar() error {
	if len(strings.TrimSpace(aclf.FechaInicio)) <= 0 {
		return fmt.Errorf("parametros enviados la fecha de inicio no puede ser vacío")
	}
	if len(strings.TrimSpace(aclf.FechaFin)) <= 0 {
		return fmt.Errorf("parametros enviados la fecha de fin no puede ser vacío")
	}
	return nil
}

// recibe fecha formato
func (aclf *ApilinkCierreloteFiltro) ToFiltroRequest() (apilinkRequest cierrelotedtos.ApilinkRequest) {
	apilinkRequest.FechaInicio, _ = time.Parse(layoutISO, aclf.FechaInicio)
	apilinkRequest.FechaFin, _ = time.Parse(layoutISO, aclf.FechaFin)
	apilinkRequest.FechaFin = commons.GetDateLastMomentTime(apilinkRequest.FechaFin)
	apilinkRequest.Number = aclf.Number
	apilinkRequest.Size = aclf.Size
	return
}
