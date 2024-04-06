package filtros

import "time"

type PagoTipoFiltro struct {
	Paginacion
	Id                     uint
	PagoTipo               string
	CargarCuenta           bool
	CargarPagos            bool
	CargarPagosIntentos    bool
	ExternalId             []string
	CargarPagosNotificado  bool
	PagoEstadosIds         []uint64
	FechaPagoInicio        time.Time
	FechaPagoFin           time.Time
	VisualizarPendientes   bool
	IdCuenta               int64
	CargarTipoPagoChannels bool
}
