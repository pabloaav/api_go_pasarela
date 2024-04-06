package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseInstallmentInfo struct {
	TransaccionesId       string  `json:"transacciones_id"`
	InstallmentsID        uint    `json:"installments_id"`
	InstallmentsdetailsID uint    `json:"installmentsdetails_id"`
	Cuota                 int64   `json:"cuota"`
	Tna                   float64 `json:"tna"`
	Tem                   float64 `json:"tem"`
	Coeficiente           float64 `json:"coeficiente"`
}

func (rii *ResponseInstallmentInfo) EntityToDtosForConciliacion(entityInstallmentsDetails entities.Installmentdetail) {
	rii.InstallmentsdetailsID = entityInstallmentsDetails.ID
	rii.InstallmentsID = entityInstallmentsDetails.InstallmentsID
	rii.Cuota = entityInstallmentsDetails.Cuota
	rii.Tna = entityInstallmentsDetails.Tna
	rii.Tem = entityInstallmentsDetails.Tem
	rii.Coeficiente = entityInstallmentsDetails.Coeficiente
}
