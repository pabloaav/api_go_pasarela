package enumsdtos

import "errors"

type EnumVendor string

const (
	Apilink  EnumVendor = "APILINK"
	Prisma   EnumVendor = "PRISMA"
	Rapipago EnumVendor = "RAPIPAGO"
)

func (e EnumVendor) IsValid() error {
	switch e {
	case Apilink, Prisma, Rapipago:
		return nil
	}
	return errors.New("el tipo de operación es inválido")
}
