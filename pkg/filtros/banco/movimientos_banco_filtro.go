package filtros

import "errors"

type MovimientosBancoFiltro struct {
	SubCuenta      string
	Tipo           EnumTipoOperacion
	TipoMovimiento []string
	Fecha          string
	Fechas         []string
	TipoOperacion  bool
	Dbcr           string
}

type EnumTipoOperacion string

const (
	Debin         EnumTipoOperacion = "DEBIN"
	Prisma        EnumTipoOperacion = "PRISMA"
	Transferencia EnumTipoOperacion = "TRANSFERENCIA"
)

func (e EnumTipoOperacion) IsValid() error {
	switch e {
	case Debin, Prisma, Transferencia:
		return nil
	}
	return errors.New("tipo de operacion es incorrecto")
}
