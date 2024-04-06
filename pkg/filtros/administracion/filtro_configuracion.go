package filtros

type ConfiguracionFiltro struct {
	Paginacion
	Id         uint
	Nombre     string
	Buscar     bool
	Nombrelike string
}
