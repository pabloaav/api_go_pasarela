package reportedtos

// Nota: en el join las columnas con nombres ambiguos se pueden setear con alias
// el nombre del atributo de la struct sigue la canvencion del alias del join

type ReporteMovimientosComisiones struct {
	Id                          int
	Tipo                        string
	MontoPago                   int // AS monto_pago
	MontoMovimiento             int // AS monto_movimiento
	MontoComision               int // AS monto_comision
	PorcentajeComision          float64
	MontoImpuesto               int // AS monto_impuesto
	PorcentajeImpuesto          float64
	MontoComisionproveedor      int
	PorcentajeComisionproveedor float64
	MontoImpuestoproveedor      int
	NombreCliente               string // AS nombre_cliente
	CreatedAt                   string
	FechaOperacion              string
	FechaPago                   string
	Subtotal                    int
}
type ResposeReporteMovimientosComisiones struct {
	Total    int
	Reportes []ReporteMovimientosComisiones
	LastPage int
}

func (rrmc *ResposeReporteMovimientosComisiones) SetTotales() {
	var total int
	for _, reporte := range rrmc.Reportes {
		reporte.Subtotal = (reporte.MontoComision + reporte.MontoImpuesto)
		total += reporte.Subtotal
	}
	rrmc.Total = total

}
