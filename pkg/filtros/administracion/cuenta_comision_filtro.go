package filtros

import "time"

type CuentaComisionFiltro struct {
	Paginacion
	Id                uint
	CuentaComision    string
	CargarCuenta      bool
	CargarChannel     bool
	ChannelId         uint
	CuentaId          uint
	Mediopagoid       uint
	ExaminarPagoCuota bool
	PagoCuota         bool
	Channelarancel    bool
	FechaPagoVigencia time.Time
}
