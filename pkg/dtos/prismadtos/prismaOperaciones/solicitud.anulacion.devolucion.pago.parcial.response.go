package prismadtos

type SolicitudAnulacionDevolucionPagoParcialResponse struct {
	ID     int64  `json:"id,omitempty"`
	Amount int64  `json:"amount,omitempty"`
	Status string `json:"status,omitempty"`
}
