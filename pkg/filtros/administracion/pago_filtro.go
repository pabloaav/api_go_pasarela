package filtros

import (
	"time"
)

type PagoFiltro struct {
	Paginacion
	Ids                                  []uint64
	PagoEstadosId                        uint64
	PagoEstadosIds                       []uint64
	CuentaId                             uint64
	Nombre                               string
	PagosTipoId                          uint64
	MedioPagoId                          uint64
	Referencia                           string
	VisualizarPendientes                 bool
	CargaPagoIntentos                    bool
	CargaPagoIntentosByExternalReference []string
	CargaMedioPagos                      bool
	CargarChannel                        bool
	CargarPagoTipos                      bool
	CargarCuenta                         bool
	CargarPagoEstado                     bool
	Uuids                                []string
	TiempoExpiracion                     string
	ExternalReference                    string
	ExternalReferences                   []string
	Fecha                                []string
	FiltroFechaPaid                      bool
	FechaPagoInicio                      time.Time
	FechaPagoFin                         time.Time
	BuscarNotificado                     bool
	Notificado                           bool
	PagosTipoIds                         []uint64
	CargarPagosItems                     bool
	Ordenar                              bool
	Descendente                          bool
	HolderNumber                         string
}

type PagoItemFiltro struct {
	PagoId uint64
}
