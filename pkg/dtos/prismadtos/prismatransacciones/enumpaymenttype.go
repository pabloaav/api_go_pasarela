package prismadtos

import "errors"

type EnumPaymentType string

const (
	Single      EnumPaymentType = "single"
	Distributed EnumPaymentType = "distributed"
)

func (e EnumPaymentType) IsValid() error {
	switch e {
	case Single, Distributed:
		return nil
	}
	return errors.New("tipo de pago seleccionado inválido")
}
