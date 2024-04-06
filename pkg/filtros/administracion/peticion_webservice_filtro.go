package filtros

type PeticionWebServiceFiltro struct {
	Paginacion
	Id                 uint
	Operacion          string
	Vendor             string
	Fecha              []string
	OrdenarPorFechaInv bool
}
