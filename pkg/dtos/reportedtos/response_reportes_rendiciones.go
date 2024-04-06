package reportedtos

type ResponseReportesRendiciones struct {
	PagoIntentoId           uint64
	Cuenta                  string // Nombre de la cuenta del cliente
	Id                      string // external_reference enviada por el cliente
	FechaCobro              string // fecha que el pagador realizo el pago
	FechaDeposito           string // fecha que se le envio el dinero al cliente(transferencia)
	ImporteCobrado          string // importe solicitud de pago
	ImporteDepositado       string // importe depositado al cliente
	CantidadBoletasCobradas string // pago items
	// ComisionPorcentaje      string // comision de telco cobrada al cliente
	// ComisionIva             string // iva Cobrado al cliente
	Comision    string // comision de telco cobrada al cliente
	Iva         string // iva Cobrado al cliente
	Concepto    string
	Retenciones string // retenciones
}

type Totales struct {
	CantidadOperaciones string
	TotalCobrado        string
	TotalRendido        string
	TotalIva            string
	TotalComision       string
	TotalRevertido      string
}

type ResponseTotales struct {
	Totales  Totales
	Detalles []*ResponseReportesRendiciones
}

// Seccion Para Reportes visuales y excel agrupados por fecha

type ResponseRendicionesClientes struct {
	CantidadRegistros   int
	Total               float64
	DetallesRendiciones []DetallesRendicion
}

type DetallesRendicion struct {
	Fecha                                        string
	Nombre                                       string
	CantidadOperaciones                          uint
	TotalCobrado                                 float64
	TotalRendido                                 float64
	TotalReversion                               float64
	TotalComision                                float64
	TotalIva                                     float64
	NroReporte                                   string
	Rendiciones                                  []ResponseReportesRendiciones
	TotalRetGanancias, TotalRetIVA, TotalRetIIBB string
}
