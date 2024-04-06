package filtros

import "time"

type MovimientoFiltro struct {
	Paginacion
	Ids                     []uint64
	PagoIntentosIds         []uint64 //Este solo se debe usarse cuando se acumula por pago intentos
	CuentaId                uint64
	FechaInicio             time.Time
	FechaFin                time.Time
	MedioPagoId             int64
	ReferenciaPago          string
	Concepto                string
	CargarPago              bool
	CargarPagoEstados       bool
	CargarPagoIntentos      bool
	CargarMedioPago         bool
	AcumularPorPagoIntentos bool
	CargarTransferencias    bool
	CargarComision          bool
	CargarImpuesto          bool

	// cargar movimientos negativos
	CargarMovimientosNegativos bool
	CargarMovimientosDetalles  bool
	CargarRetenciones          bool
	CargarMovimientoSubcuenta  bool
	// IdsMovNeg                  []uint64
}
