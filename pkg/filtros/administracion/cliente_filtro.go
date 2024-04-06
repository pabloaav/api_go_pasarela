package filtros

type ClienteFiltro struct {
	Paginacion
	Id                   uint
	DistintoId           uint
	Cuit                 string
	UserId               uint
	RetiroAutomatico     bool
	CargarImpuestos      bool
	CargarCuentas        bool
	CargarRubros         bool
	CargarCuentaComision bool
	CargarTiposPago      bool
	CargarContactos      bool
	ClientesIds          []uint
	CargarEnvio          bool
}
