package linkdtos

import "errors"

type EnumTipoDebin string

const (
	DebinDefault EnumTipoDebin = "DEBIN"
	DebinPLF     EnumTipoDebin = "DEBINPLF"
	//La consulta DebinPLF filtrará unicamente los debines de tipo Plazo Fijo
)

func (e EnumTipoDebin) IsValid() error {
	switch e {
	case DebinDefault, DebinPLF:
		return nil
	}
	return errors.New("tipo EnumTipoDebin con formato inválido")
}
