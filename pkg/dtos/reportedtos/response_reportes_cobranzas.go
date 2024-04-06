package reportedtos

type ResponseCobranzasClientes struct {
	CantidadCobranzas int
	Total             uint
	Cobranzas         []DetallesCobranza
}

type DetallesCobranza struct {
	Fecha          string
	Nombre         string
	Registros      uint
	Subtotal       uint
	TotalComision  uint
	TotalIva       uint
	TotalRetencion uint
	NroReporte     string
	Pagos          []DetallesPagosCobranza
}

type DetallesPagosCobranza struct {
	Id          int
	Cliente     string `json:"cliente"`
	Pagoestado  string `json:"pagoestado"`
	Descripcion string `json:"descripcion"`
	Referencia  string `json:"referencia"`
	PayerName   string `json:"payer_name"`
	PayerEmail  string `json:"payer_email"`
	TotalPago   uint   `json:"total_pago"`
	MedioPago   string `json:"medio_pago"`
	CanalPago   string `json:"canal_pago"`
	Cuenta      string `json:"cuenta"`
	FechaPago   string `json:"fecha_pago"`
	FechaCobro  string `json:"fecha_cobro"`
	Comision    uint   `json:"comision"`
	Iva         uint   `json:"iva"`
	Retencion   uint   `json:"retencion"`
}
