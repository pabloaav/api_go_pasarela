package ribcradtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)

type RiDatosFondos struct {
	Numero             string
	Denominacion       string
	Agente             Agente
	DenominacionAgente string
	CuitAgente         string
}

func (r *RiDatosFondos) IsValid() error {

	if len(r.Numero) > 5 {
		return fmt.Errorf(administraciondtos.ERROR_RI_NUMERO_FONDO)
	}
	if len(r.Denominacion) > 50 {
		return fmt.Errorf(administraciondtos.ERROR_RI_DENOMINACION_FONDO)
	}
	err := r.Agente.IsValid()

	if err != nil {
		return err
	}

	if len(r.DenominacionAgente) > 50 {
		return fmt.Errorf(administraciondtos.ERROR_RI_DENOMINACION_AGENTE)
	}

	if len(r.CuitAgente) > 11 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CUIT_AGENTE)
	}

	return nil
}

type Agente int

const (
	AgenteAdministracion Agente = iota
	AgenteCustodia
	AgenteColocacion
)

func (e Agente) IsValid() error {
	switch e {
	case AgenteAdministracion, AgenteCustodia, AgenteColocacion:
		return nil
	}
	return fmt.Errorf(administraciondtos.ERROR_RI_AGENTE)
}

func (r *RiDatosFondos) ToString() string {

	return fmt.Sprintf("%s;%s;%v;%s;%s", r.Numero, r.Denominacion, r.Agente, r.DenominacionAgente, r.CuitAgente)

}
