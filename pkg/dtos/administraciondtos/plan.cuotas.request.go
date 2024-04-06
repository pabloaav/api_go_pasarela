package administraciondtos

import (
	"time"
)

type PlanCuotasRequest struct {
	Descripcion       string
	Issuer            string
	VigenciaDesde     time.Time
	VigenciaHasta     time.Time
	Installmentdetail []PlanCuotasDetalle
}

type PlanCuotasDetalle struct {
	InstallmentsID int64
	VigenciaDesde  time.Time
	Activo         bool
	Cuota          int64
	Tna            float64
	Tem            float64
	Coeficiente    float64
}

// func (planCuotas *PlanCuotasRequest) Validar() error {
// 	if commons.StringIsEmpity(planCuotas.Descripcion) {
// 		return errors.New(ERROR_STRING + " descripción")
// 	}
// 	return nil
// }
