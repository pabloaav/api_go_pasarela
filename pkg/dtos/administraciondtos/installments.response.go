package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type InstallmentsResponse struct {
	Id                      uint       `json:"id"`
	Descripcion             string     `json:"descripcion"`
	Issuer                  string     `json:"issuer"`
	VigenciaDesde           time.Time  `json:"vigencia_desde"`
	VigenciaHasta           *time.Time `json:"vigencia_hasta"`
	MediopagoinstallmentsID int64      `json:"mediopagoinstallments_id"`
}

func (ir *InstallmentsResponse) EntityToDtos(installment entities.Installment) {
	ir.Id = installment.ID
	ir.Descripcion = installment.Descripcion
	ir.Issuer = installment.Issuer
	ir.VigenciaDesde = installment.VigenciaDesde
	ir.VigenciaHasta = installment.VigenciaHasta
	ir.MediopagoinstallmentsID = installment.MediopagoinstallmentsID
}
