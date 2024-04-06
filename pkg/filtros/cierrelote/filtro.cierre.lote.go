package filtros

type FiltroCierreLote struct {
	Nroestablecimiento string
	MatchCl            bool
	MovimientosMX      bool
	PagosPx            bool
	Banco              bool
	/*se utilizan para obtener los CL conciliados con movimiento y no fueron procesados con el banco */
	//EstadoBancoId          bool
	//EstadoPrismaMovimiento bool
	EstadoFechaPago bool
	FechaPago       string
	PrismaPagoId    int64
	Compras         bool
	Devolucion      bool
	Anulacion       bool
	ContraCargo     bool
	ContraCargoMx   bool
	ContraCargoPx   bool
	Reversion       bool
	DetallePagoId   int64
}
