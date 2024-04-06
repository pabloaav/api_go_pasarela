package administraciondtos

type PlanCuotasResponse struct {
	Id                      uint                        `json:"id"`
	Descripcion             string                      `json:"descripcion"`
	MediopagoinstallmentsID int64                       `json:"mediopagoinstallments_id"`
	Installmentdetail       []PlanCuotasResponseDetalle `json:"installmentdetail"`
}

type PlanCuotasResponseDetalle struct {
	InstallmentsID uint    `json:"installments_id,omitempty"`
	Cuota          uint    `json:"cuota"`
	Tna            float64 `json:"tna"`
	Tem            float64 `json:"tem"`
	Coeficiente    float64 `json:"coeficiente"`
}
