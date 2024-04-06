package filtros

type CobranzasClienteFiltro struct {
	FechaInicio                   string
	FechaFin                      string
	ClienteId                     int      `json:"cliente_id"`
	CuentaId                      int      `json:"cuenta_id"`
	FiltroBarcode                 []string `json:"filtro_barcode"`
	ObtenerBarcodes               bool
	ObtenerApiLinkByFechaCobro    bool
	ObtenerPrismaByFechaOperacion bool
	FiltrarPorFechaCobro          bool
}
