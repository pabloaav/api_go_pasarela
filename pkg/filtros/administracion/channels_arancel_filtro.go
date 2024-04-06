package filtros

type ChannelArancelFiltro struct {
	Paginacion
	RubrosId           uint
	ChannelsId         uint
	CargarRubro        bool
	ChannelId          uint
	CargarChannel      bool
	PagoCuota          bool
	MedioPagoId        uint
	CargarAllMedioPago bool
}
