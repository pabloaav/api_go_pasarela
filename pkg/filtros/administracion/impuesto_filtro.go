package filtros

type ImpuestoFiltro struct {
	Paginacion
	Id              uint    `json:"id"`
	Impuesto        string  `json:"impuesto"`
	Porcentaje      float64 `json:"porcentaje"`
	Tipo            string  `json:"tipo"`
	OrdenarPorFecha bool    `json:"ordenarporfecha"`
}
