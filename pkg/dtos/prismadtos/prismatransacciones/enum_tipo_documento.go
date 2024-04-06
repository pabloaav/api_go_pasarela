package prismadtos

import "errors"

type EnumTipoDocumento string

const (
	Dni EnumTipoDocumento = "DNI"
	Ci  EnumTipoDocumento = "CI"
	Le  EnumTipoDocumento = "LE"
	Lc  EnumTipoDocumento = "LC"
)

func (e EnumTipoDocumento) IsValid() error {
	switch e {
	case Dni, Ci, Le, Lc:
		return nil
	}
	return errors.New(ERROR_TIPO_DOCUMENTO)
}
