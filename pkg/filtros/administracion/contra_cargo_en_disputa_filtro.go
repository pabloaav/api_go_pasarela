package filtros

import (
	"errors"
	"time"
)

type ContraCargoEnDisputa struct {
	Paginacion
	IdCliente           uint
	IdCuenta            uint
	CargarCuentas       bool
	CargarTiposPago     bool
	CargarPagos         bool
	CargarPagosIntentos bool
	TransactionId       []string
	FechaInicio         string `json:"fecha_inicio"`
	FechaFin            string `json:"fecha_fin"`
	FechaCreacion       bool
	FechaOperacion      bool
	FechaPago           bool
	FechaCierre         bool
}

func (r *ContraCargoEnDisputa) ValidarFechas() (estadoValidacion bool, erro error) {

	if r.FechaCreacion || r.FechaOperacion || r.FechaPago || r.FechaCierre {

		if r.FechaInicio == "" || r.FechaFin == "" {
			erro = errors.New("no se recibio las fecha para filtrar ")
			return
		}

		tInicio, error1 := time.Parse("02-01-2006", r.FechaInicio)
		if error1 != nil {
			erro = error1
			return
		}

		tFin, error2 := time.Parse("02-01-2006", r.FechaFin)
		if error2 != nil {
			erro = error2
			return
		}

		if tInicio.After(tFin) {
			erro = errors.New("se recibio fecha fin anterior a la fecha inicio")
			return
		}

		r.FechaInicio = tInicio.Format("2006-01-02")
		r.FechaFin = tFin.Format("2006-01-02")

	}

	estadoValidacion = true
	return
}
