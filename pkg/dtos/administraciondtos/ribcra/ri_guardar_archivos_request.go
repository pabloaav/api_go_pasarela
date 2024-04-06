package ribcradtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)

type RIGuardarArchivosRequest struct {
	Ruta string
	RI   interface{}
}

func (r *RIGuardarArchivosRequest) IsValid() error {

	if len(strings.TrimSpace(r.Ruta)) < 1 {
		return fmt.Errorf(administraciondtos.ERROR_RUTA_INVALIDA)
	}

	if r.RI == nil {
		return fmt.Errorf(administraciondtos.ERROR_RI_DATOS)
	}

	return nil
}
