package prismadtos

import prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"

type SolicitudAnulacionDevolucionResponse struct {
	ID            int64                    `json:"id,omitempty"`
	Amount        int64                    `json:"amount,omitempty"`
	SubPayments   interface{}              `json:"sub_payments,omitempty"`
	Error         PrismaErrorResponse      `json:"error,omitempty"`
	Status        string                   `json:"status,omitempty"`
	StatusDetails prismadtos.StatusDetails `json:"status_details,omitempty"`
}
