package linktransferencia

import (
	"encoding/json"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
)

type ResponseGetTransferenciaLink struct {
	Origen                   OrigenResponseTransferenciaLink  `json:"origen"`
	Destino                  DestinoResponseTransferenciaLink `json:"destino"`
	Importe                  string                           `json:"importe"`
	Moneda                   linkdtos.EnumMoneda              `json:"moneda"`
	Motivo                   linkdtos.EnumMotivoTransferencia `json:"motivo"`
	Referencia               string                           `json:"referencia"`
	NumeroReferenciaBancaria string                           `json:"numeroReferenciaBancaria"`
	FechaOperacion           time.Time                        `json:"fechaOperacion"`
}

func (r *ResponseGetTransferenciaLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
