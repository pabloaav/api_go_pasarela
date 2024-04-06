package linkdtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type EnumMoneda string

const (
	Pesos EnumMoneda = "ARS"
	Dolar EnumMoneda = "USD"
)

func (e EnumMoneda) IsValid() error {
	switch e {
	case Pesos, Dolar:
		return nil
	}
	return errors.New(tools.ERROR_ENUM_MONEDA)
}
