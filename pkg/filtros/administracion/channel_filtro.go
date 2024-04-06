package filtros

type ChannelFiltro struct {
	Paginacion
	Id              uint
	Channel         string
	Channels        []string
	CargarMedioPago bool
}
