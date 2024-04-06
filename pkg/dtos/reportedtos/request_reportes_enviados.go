package reportedtos

import (
	"errors"

	"github.com/jinzhu/now"
)

type RequestReportesEnviados struct {
	FechaInicio, FechaFin string
	Number                uint            `json:"number"`
	Size                  uint            `json:"size"`
	TipoReporte           EnumTipoReporte `json:"tipo_reporte"`
	Cliente               string          `json:"cliente"`
	Enum                  bool            `json:"enum"`
	SinNumero             bool            `json:"sin_numero"`
	Fecharendicion        []string         `json:"fecharendicion"`
}

type EnumTipoReporte string

const (
	pagos       EnumTipoReporte = "pagos"
	rendiciones EnumTipoReporte = "rendiciones"
	revertidos  EnumTipoReporte = "revertidos"
	todos       EnumTipoReporte = "todos"
)

func (e EnumTipoReporte) IsValid() bool {
	switch e {
	case pagos, rendiciones, revertidos, todos:
		return true
	}
	return false
}

// recuperar el string del EnumTipoReporte
func (e EnumTipoReporte) ToString() string {
	switch e {
	case pagos:
		return "pagos"
	case rendiciones:
		return "rendiciones"
	case revertidos:
		return "revertidos"
	default:
		return "todos"
	}
}

func (rre *RequestReportesEnviados) Validar() (erro error) {

	// la fecha de inicio no puede ser cero
	if len(rre.FechaInicio) == 0 {
		erro = errors.New("debe enviar una fecha de inicio")
		return
	} else {
		rre.FechaInicio += " 00:00:00"
	}

	// la fecha de fin puede ser cero. en ese caso se toma la fecha actual
	if len(rre.FechaFin) == 0 {
		// now := time.Now().UTC()
		// fin:=commons.GetDateLastMomentTime(now)
		rre.FechaFin = now.EndOfDay().String()[0:19]
	} else {
		rre.FechaFin += " 23:59:59"
	}

	if len(rre.TipoReporte) == 0 {
		rre.TipoReporte = "todos"
	}

	if !rre.TipoReporte.IsValid() {
		erro = errors.New("el tipo de reporte debe ser alguno de los definidos")
		return
	}

	return
}
