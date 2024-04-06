package reportedtos

import (
	"errors"
	"time"
)

type RequestLogs struct {
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	Number      int       `json:"number"`
	Size        int       `json:"size"`
}

func (r *RequestLogs) ValidarFechas() (estadoValidacion bool, erro error) {

	estadoValidacion = false
	if r.FechaInicio.IsZero() || r.FechaFin.IsZero() {
		erro = errors.New(" debe enviar una fecha de inicio y de fin")
		return
	}
	if !r.FechaInicio.IsZero() && !r.FechaFin.IsZero() {
		if r.FechaInicio.After(r.FechaFin) {
			erro = errors.New("la fecha de inicio no puede ser mayor que la fecha fin ")
			return
		}
	}
	estadoValidacion = true
	return
}
