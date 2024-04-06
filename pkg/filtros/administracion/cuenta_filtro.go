package filtros

type CuentaFiltro struct {
	Id               uint
	DistintoId       uint
	Cbu              string
	Cvu              string
	CargarComisiones bool
	CargarImpuestos  bool

	TipoCuenta string
	ApiKey     string
}
