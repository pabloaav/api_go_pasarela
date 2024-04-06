package filtros

type ChannelAranceFiltro struct {
	Paginacion
	Id              uint
	CargarRubro     bool
	CargarChannel   bool
	ChannelId       uint
	RubrosId        uint
	OrdernarChannel bool
}
