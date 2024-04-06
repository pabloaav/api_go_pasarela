package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

const (
	PRISMA  = "prisma"
	DEBIN   = "debin"
	OFFLINE = "offline"
)

type ReportesFactory interface {
	GetProcesarReportes(m string) (ReportesPagos, error)
}

type procesarReportesFactory struct{}

func NewRecorrerArchivos() ReportesFactory {
	return &procesarReportesFactory{}
}

func (r *procesarReportesFactory) GetProcesarReportes(m string) (ReportesPagos, error) {
	switch m {
	case PRISMA:
		return NewReportesPrisma(util.Resolve()), nil
	case DEBIN:
		return NewReportesDebin(util.Resolve()), nil
	case OFFLINE:
		return NewReportesOffline(util.Resolve()), nil
	default:
		return nil, fmt.Errorf("el tipo de reportes a procesar  %v, no es valido", m)

	}
}

type ReportesPagos interface {
	ResponseReportes(s *reportesService, listaPagos reportedtos.TipoFactory) (response []reportedtos.ResponseFactory)
}
