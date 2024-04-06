package filtros

type FiltroPrismaMovimientoDetalle struct {
	ListIdsDetalle               []int64
	ListIdsCabecera              []int64
	FechaPago                    string
	CargarCabecera               bool
	Contracargovisa              bool
	Contracargomaster            bool
	Tipooperacion                bool
	Rechazotransaccionprincipal  bool
	Rechazotransaccionsecundario bool
	Motivoajuste                 bool
	Match                        bool
	MatchCl                      bool
}
