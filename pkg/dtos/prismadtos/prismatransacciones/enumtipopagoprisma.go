package prismadtos

import "errors"

type EnumTipoPagoPrisma string

const (
	Simple  EnumTipoPagoPrisma = "simple"
	Offline EnumTipoPagoPrisma = "offline"
)

func (e EnumTipoPagoPrisma) IsValid() error {
	switch e {
	case Simple, Offline:
		return nil
	}
	return errors.New("tipo de pago seleccionado inválido")
}
