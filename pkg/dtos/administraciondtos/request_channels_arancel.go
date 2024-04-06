package administraciondtos

import (
	"errors"
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestChannelsAranncel struct {
	Id            uint
	RubrosId      uint
	ChannelsId    uint
	Importe       float64
	Fechadesde    time.Time
	Tipocalculo   EnumTipoCal
	Importeminimo float64
	Importemaximo float64
	Mediopagoid   int64
	Pagocuota     bool
}

type EnumTipoCal string

const (
	Porcentaje EnumTipoCal = "PORCENTAJE"
	Fijo       EnumTipoCal = "FIJO"
)

func (r *RequestChannelsAranncel) ToChannelsArancel(cargarId bool) (response entities.Channelarancele) {
	if cargarId {
		response.ID = r.Id
	}
	fecha := r.Fechadesde.Format("2006-01-02")
	response.RubrosId = int64(r.RubrosId)
	response.ChannelsId = r.ChannelsId
	response.Importe = r.Importe
	response.Fechadesde = fecha
	response.Tipocalculo = string(r.Tipocalculo)
	response.Importeminimo = r.Importeminimo
	response.Importemaximo = r.Importemaximo
	response.Mediopagoid = r.Mediopagoid
	response.Pagocuota = r.Pagocuota
	return
}

func (r *RequestChannelsAranncel) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}
	if r.RubrosId < 1 {
		return fmt.Errorf(tools.ERROR_RUBRO_ID)
	}

	if r.ChannelsId < 1 {
		return fmt.Errorf(tools.ERROR_CHANNEL)
	}

	if r.Fechadesde.IsZero() {
		return fmt.Errorf(tools.ERROR_VIGENCIA_DESDE)
	}

	if r.Tipocalculo != Porcentaje && r.Tipocalculo != Fijo {
		return errors.New(tools.ERROR_TIPO)
	}

	return
}
