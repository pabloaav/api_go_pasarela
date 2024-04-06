package filtros

type UsuarioBloqueadoFiltro struct {
	Paginacion
	Id         uint
	Ids        []uint
	Permanente bool
}
