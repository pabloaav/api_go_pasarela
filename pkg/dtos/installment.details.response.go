package dtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type InstallmentDetailsResponse struct {
	Id                  uint
	InstallmentsID      uint
	NroCuota            int64
	TasaEfectivaMensual float64
	Coeficiente         float64
	Impuesto            float64
}

func (i *InstallmentDetailsResponse) DtosToEntity() (resultEntity entities.Installmentdetail) {
	if i.Id <= 0 {
		resultEntity.ID = 0
	}
	resultEntity.ID = i.Id
	resultEntity.InstallmentsID = i.InstallmentsID
	resultEntity.Cuota = i.NroCuota
	resultEntity.Tem = i.TasaEfectivaMensual
	resultEntity.Coeficiente = i.Coeficiente
	return
}
