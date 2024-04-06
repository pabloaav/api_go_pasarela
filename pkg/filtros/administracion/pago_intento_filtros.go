package filtros

import "time"

type PagoIntentoFiltros struct {
	PagosId                     []uint
	PagoEstadosIds              []uint64
	Channel                     bool
	CargarPago                  bool
	CargarPagoTipo              bool
	CargarPagoEstado            bool
	CargarCuenta                bool
	CargarCliente               bool
	CargarPagoCalculado         bool
	ExternalId                  bool
	FechaPagoInicio             time.Time
	FechaPagoFin                time.Time
	ClienteId                   uint64
	CuentaId                    uint64
	PagoIntentoAprobado         bool
	CargarMovimientosTemporales bool
	CargarPagoItems             bool
}
