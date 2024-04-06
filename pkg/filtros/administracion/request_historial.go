package filtros

import "time"

type RequestHistorial struct {
	Paginacion
	FechaInicio time.Time
	FechaFin    time.Time
	Correo      string
}
