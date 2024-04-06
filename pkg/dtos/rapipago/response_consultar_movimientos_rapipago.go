package rapipago

// CheckoutResponse respuesta utilizada en el frontend del checkout
type ResponseConsultarMovimientosRapipago struct {
	FechaCobro      string
	Importe         int64
	CodigoBarras    string
	BancoExternalId int64
	Match           bool
}
