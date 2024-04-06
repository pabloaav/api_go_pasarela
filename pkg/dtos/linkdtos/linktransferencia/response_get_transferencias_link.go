package linktransferencia

import (
	"encoding/json"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
)

type ResponseGetTransferenciasLink struct {
	Transferencias []Transferencia           `json:"transferencias"`
	Paginado       PaginadoTransferenciaLink `json:"paginado"`
}

type Transferencia struct {
	Origen                   OrigenResponseTransferenciaLink  `json:"origen"`
	Destino                  DestinoResponseTransferenciaLink `json:"destino"`
	Importe                  string                           `json:"importe"`
	Moneda                   linkdtos.EnumMoneda              `json:"moneda"`
	Motivo                   linkdtos.EnumMotivoTransferencia `json:"motivo"`
	Referencia               string                           `json:"referencia"`
	NumeroReferenciaBancaria string                           `json:"numeroReferenciaBancaria"`
	FechaOperacion           time.Time                        `json:"fechaOperacion"`
}

type OrigenResponseTransferenciaLink struct {
	Cbu         string `json:"cbu"`
	RazonSocial string `json:"razonSocial"`
}

type DestinoResponseTransferenciaLink struct {
	Cbu     string                   `json:"cbu"`
	Alias   string                   `json:"alias"`
	Titular TitularTransferenciaLink `json:"titular"`
}

type TitularTransferenciaLink struct {
	IdTributario string `json:"idTributario"`
	RazonSocial  string `json:"razonSocial"`
}

type PaginadoTransferenciaLink struct {
	Total       int32 `json:"total"`
	CantPaginas int32 `json:"cantPaginas"`
}

func (r *ResponseGetTransferenciasLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
