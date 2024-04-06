package reportedtos

import (
	"errors"
	"time"
)

// Nota: en el join las columnas con nombres ambiguos se pueden setear con alias
// el nombre del atributo de la struct sigue la canvencion del alias del join

type RequestReporteMovimientosComisiones struct {
	FechaInicio, FechaFin string
	Number                int  `json:"number"`
	Size                  int  `json:"size"`
	ClienteId             int  `json:"cliente_id"`
	CuentaId              int  `json:"cuenta_id"`
	FechaPago             bool `json:"fechaPago"`
	SoloReversiones       bool `json:"solo_reversiones"`
}

func (rrmc *RequestReporteMovimientosComisiones) Validar() (erro error) {

	// la fecha de inicio no puede ser cero
	if len(rrmc.FechaInicio) == 0 {
		erro = errors.New("debe enviar una fecha de inicio")
		return
	}

	if len(rrmc.FechaFin) == 0 {
		DDMMYYYY := "02-01-2006"
		now := time.Now().UTC()
		rrmc.FechaFin = now.Format(DDMMYYYY)
	}

	return
}
