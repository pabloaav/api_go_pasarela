package prismadtos

type DeletSolicitudPagoResponse struct {
	Amount int64  `json:"amount,omitempty"`
	Status string `json:"status,omitempty"`
}
