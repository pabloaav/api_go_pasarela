package linktransferencia

import (
	"encoding/json"
	"time"
)

type ResponseTransferenciaCreateLink struct {
	NumeroReferenciaBancaria   string    `json:"numeroReferenciaBancaria"`
	FechaOperacion             time.Time `json:"fechaOperacion"`
	NumeroConciliacionBancaria string    `json:"numeroConciliacionBancaria"`
	MovimientosIdTransferidos  []uint64  `json:"movimientos_id_transferidos"`
	//MovimientosIdReversiones   []uint64  `json:"movimientos_id_revertidos"`
}

func (r *ResponseTransferenciaCreateLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
