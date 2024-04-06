package administraciondtos

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ImpuestoRequest struct {
	Id         uint      `json:"id"`
	Impuesto   string    `json:"impuesto"`
	Porcentaje float64   `json:"porcentaje"`
	Tipo       EnumTipo  `json:"tipo"`
	Fechadesde time.Time `json:"fechadesde"`
	Activo     bool      `json:"activo"`
}

type EnumTipo string

const (
	IVA  EnumTipo = "IVA"
	IIBB EnumTipo = "IIBB"
)

func (c *ImpuestoRequest) ToImpuesto(cargarId bool) (impuesto entities.Impuesto) {
	if cargarId {
		impuesto.ID = c.Id
	}
	var porcentaje float64 = c.Porcentaje
	impuesto.Impuesto = strings.ToUpper(c.Impuesto)
	impuesto.Porcentaje = porcentaje
	impuesto.Tipo = string(c.Tipo)
	impuesto.Fechadesde = c.Fechadesde
	impuesto.Activo = true
	return

}

func (impuesto *ImpuestoRequest) Validar() error {

	const ERROR_CAMPO = "el campo es obligatorio"
	const ERROR_TIPO = "el tipo de impuesto no es valido"
	const ERROR_PORCENTAJE = "se debe informar el porcentaje"
	const ERROR_FECHA = "se debe informar la fecha"
	// digitCheckInt := regexp.MustCompile(`^[0-9]+$`)

	if commons.StringIsEmpity(impuesto.Impuesto) {
		return fmt.Errorf(ERROR_CAMPO)
	}

	if impuesto.Porcentaje < 0 {
		return fmt.Errorf(ERROR_PORCENTAJE)
	}

	if impuesto.Tipo != IVA && impuesto.Tipo != IIBB {
		return errors.New(ERROR_TIPO)
	}

	if impuesto.Fechadesde.IsZero() {
		return fmt.Errorf(ERROR_FECHA)
	}

	return nil
}
