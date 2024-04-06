package filtros

import "time"

type FiltroMovimientosPrisma struct {
	BuscarPorFechaPresentacion bool
	BuscarPorFechaPago         bool
	BuscarPorFechaCreacion     bool
	// BuscarPorFechaOrigenCompra bool
	FechaInicio time.Time
	FechaFin    time.Time
	Number      uint32
	Size        uint32
}
