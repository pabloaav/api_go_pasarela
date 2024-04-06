package filtros

import (
	"errors"
	"regexp"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type FiltroTablasConciliadas struct {
	FechaPago          string
	FechaPresentacion  string
	NroEstablecimiento string
	Match              bool
	Reversion          bool
}

func (fil *FiltroTablasConciliadas) Validar() error {
	erro := errors.New("parametros de filtro no validos")
	regularFecha := regexp.MustCompile(`([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))$`)
	if !commons.StringIsEmpity(fil.FechaPago) {
		valorboolpresentacion := regularFecha.MatchString(fil.FechaPago)
		if !valorboolpresentacion {
			return erro
		}
	}
	if !commons.StringIsEmpity(fil.FechaPresentacion) {
		valorboolpago := regularFecha.MatchString(fil.FechaPresentacion)
		if !valorboolpago {
			return erro
		}
	}

	if !commons.StringIsEmpity(fil.NroEstablecimiento) {
		if len(fil.NroEstablecimiento) != 10 {
			return erro
		}
	}
	return nil
}
