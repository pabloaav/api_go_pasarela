package bancodtos

type RequestParams struct {
	Number                  uint32
	Size                    uint32
	FechaInicio             string
	FechaFin                string
	ListaIdsTipoMovimientos string
	DebitoCredito           string
}
