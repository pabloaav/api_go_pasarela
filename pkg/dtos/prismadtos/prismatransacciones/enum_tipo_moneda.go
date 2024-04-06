package prismadtos

import (
	"errors"
)

type EnumTipoMoneda string

const (
	Pesos EnumTipoMoneda = "ARS"
	Dolar EnumTipoMoneda = "USD"
)

func (e EnumTipoMoneda) IsValid() error {
	switch e {
	case Pesos, Dolar:
		return nil
	}
	return errors.New(ERROR_CURRENCY)
}
