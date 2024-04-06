package administraciondtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type EnumTipoSubcuenta string

const (
	Primaria   EnumTipoSubcuenta = "principal"
	Secundaria EnumTipoSubcuenta = "secundaria"
)

func (e EnumTipoSubcuenta) IsValid() error {
	switch e {
	case Primaria, Secundaria:
		return nil
	}
	return errors.New(tools.ERROR_TIPO_CUENTA)
}
