package reportedtos

import (
	"errors"
	"time"
)

type RequestReporteClientes struct {
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	ClienteId   int    `json:"cliente_id"`
	CuentaId    int    `json:"cuenta_id"`
}

func (rrmc *RequestReporteClientes) Validar() (erro error) {

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
