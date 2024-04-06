package filtros

import "time"

type RequestClMultipago struct {
	Paginacion
	FechaInicio time.Time
	FechaFin    time.Time
	CodigoBarra string
}
