package linkdtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type EnumPagiandoDebin string //Cantidad de Registros por pagina

const (
	Vacio     EnumPagiandoDebin = ""
	Cinco     EnumPagiandoDebin = "5"
	Diez      EnumPagiandoDebin = "10"
	Veinte    EnumPagiandoDebin = "20"
	Treinta   EnumPagiandoDebin = "30"
	Cuarenta  EnumPagiandoDebin = "40"
	Cinquenta EnumPagiandoDebin = "50"
	Sesenta   EnumPagiandoDebin = "60"
	Setenta   EnumPagiandoDebin = "70"
	Ochenta   EnumPagiandoDebin = "80"
	Noventa   EnumPagiandoDebin = "90"
	Cien      EnumPagiandoDebin = "100"
)

func (e EnumPagiandoDebin) IsValid() error {
	switch e {
	case Vacio, Cinco, Diez, Veinte, Treinta, Cuarenta,
		Cinquenta, Sesenta, Setenta, Ochenta, Noventa, Cien:
		return nil
	}
	return errors.New(tools.ERROR_ENUM_PAGINADO_TAMANIO)
}
