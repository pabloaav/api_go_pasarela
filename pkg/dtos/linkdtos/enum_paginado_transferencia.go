package linkdtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type EnumPagiandoTransferencia string

const (
	DiezTransf      EnumPagiandoTransferencia = "10"
	VeinteTransf    EnumPagiandoTransferencia = "20"
	CinquentaTransf EnumPagiandoTransferencia = "50"
	CienTransf      EnumPagiandoTransferencia = "100"
)

func (e EnumPagiandoTransferencia) IsValid() error {
	switch e {
	case DiezTransf, VeinteTransf, CinquentaTransf, CienTransf, "":
		return nil
	}
	return errors.New(tools.ERROR_ENUM_PAGINADO_TAMANIO)
}
