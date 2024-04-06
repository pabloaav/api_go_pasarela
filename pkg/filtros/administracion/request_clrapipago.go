package filtros

import "time"

type RequestClrapipago struct {
	Paginacion
	FechaInicio time.Time
	FechaFin    time.Time
	CodigoBarra string
}
