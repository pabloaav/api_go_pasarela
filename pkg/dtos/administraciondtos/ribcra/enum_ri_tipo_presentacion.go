package ribcradtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)

type EnumRITipoPresentacion string

const (
	Normal        EnumRITipoPresentacion = "Normal"
	Rectificativa EnumRITipoPresentacion = "Rectificativa"
)

func (e EnumRITipoPresentacion) IsValid() error {
	switch e {
	case Normal, Rectificativa:
		return nil
	}
	return fmt.Errorf(administraciondtos.ERROR_RI_TIPO_PRESENTACION)
}
