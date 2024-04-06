package reportedtos

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type RequestPagosClientes struct {
	Cliente     uint      `json:"cliente_id"`
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	// FechaReversion time.Time `json:"fecha_reversion"`
	FechaAdicional       time.Time `json:"fecha_adicional"`
	CargarFechaAdicional bool
	EnviarEmail          bool
	ClientesIds          []uint `json:"clientes_ids"`
	ClientesString       string `json:"clientes_string"`
}

type ClientesId struct {
	Id uint `json:"id"` // Id clientereporte"`
}

func (r *RequestPagosClientes) ValidarFechas() (estadoValidacion ValidacionesFiltro, erro error) {
	estadoValidacion.Cliente = true
	estadoValidacion.Inicio = true
	estadoValidacion.Fin = true
	if r.FechaInicio.IsZero() && r.FechaFin.IsZero() {
		erro = errors.New("por lomenos debe enviar una fecha de inicio")
		return
	}
	if !r.FechaInicio.IsZero() && !r.FechaFin.IsZero() {
		if r.FechaInicio.After(r.FechaFin) {
			erro = errors.New("la fecha de inicio no puede ser mayor que la fecha fin ")
			return
		}
	}
	if r.Cliente == 0 {
		estadoValidacion.Cliente = false
	}
	return
}

type ValidacionesFiltro struct {
	Inicio  bool
	Fin     bool
	Cliente bool
}

func (r *RequestPagosClientes) ObtenerIdsClientes() {
	var arrayString []string
	//if len(r.ClientesString) > 0 {
	arrayString = strings.Split(r.ClientesString, ",")
	//}
	for _, value := range arrayString {
		result, _ := strconv.ParseUint(value, 10, 32)
		r.ClientesIds = append(r.ClientesIds, uint(result))
	}
}
