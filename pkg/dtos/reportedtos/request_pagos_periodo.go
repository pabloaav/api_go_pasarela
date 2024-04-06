package reportedtos

import (
	"time"

	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

type RequestPagosPeriodo struct {
	Paginacion   filtros.Paginacion
	ClienteId    uint64
	CuentaId     uint64
	ApiKey       string
	FechaInicio  time.Time
	FechaFin     time.Time
	PagoIntento  uint64
	PagoIntentos []uint64
	Pagos        []uint64
	PagoEstados  []uint
	// FechaLote                    string
	TipoMovimiento                  string
	CargarComisionImpuesto          bool
	CargarMovimientosTransferencias bool
	CargarCliente                   bool
	CargarPagoIntentos              bool
	CargarCuenta                    bool
	CargarMedioPago                 bool
	CargarReversion                 bool
	CargarReversionReporte          bool
	OrdenadoFecha                   bool
	Number                          uint32
	Size                            uint32
	FiltroBarcode                   bool
	FiltroApiLink                   bool
	Barcodes                        []string
	ExternalIDs                     []string
	CargarRetenciones               bool
}
