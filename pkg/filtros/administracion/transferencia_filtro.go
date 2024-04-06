package filtros

import "time"

type TransferenciaFiltro struct {
	Paginacion
	MovimientosIds     []uint64
	CuentaId           uint
	FechaInicio        time.Time
	FechaFin           time.Time
	ReferenciaBancaria string
	CargarPaginado     bool
	Number             uint32
	Size               uint32
	Match              uint // campo match tabla transferencias
}
