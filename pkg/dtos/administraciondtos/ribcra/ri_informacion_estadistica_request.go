package ribcradtos

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)

type GetInformacionEstadisticaRequest struct {
	FechaInicio time.Time
	FechaFin    time.Time
}

func (r *GetInformacionEstadisticaRequest) IsValid() (erro error) {

	if r.FechaInicio.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}

	if !r.FechaFin.IsZero() && r.FechaInicio.After(r.FechaFin) {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}
	if r.FechaFin.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA)
		return
	}

	return

}

type BuildInformacionEstadisticaRequest struct {
	RiInfestadistica []RiInfestadistica
	Rectificativa    bool
	Opera            bool
	Periodo          string
}
