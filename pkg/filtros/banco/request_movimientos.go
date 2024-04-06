package filtros

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

type RequestMovimientos struct {
	Paginacion              filtros.Paginacion
	FechaInicio             string
	FechaFin                string
	ListaIdsTipoMovimientos string
	DebitoCredito           string
}

func (rm *RequestMovimientos) ParamsToFiltro(params bancodtos.RequestParams) {
	rm.FechaInicio = params.FechaInicio
	rm.FechaFin = params.FechaFin
	rm.Paginacion.Number = params.Number
	rm.Paginacion.Size = params.Size
	if len(params.ListaIdsTipoMovimientos) != 0 {
		rm.ListaIdsTipoMovimientos = params.ListaIdsTipoMovimientos
	}
	if len(params.DebitoCredito) != 0 {
		rm.DebitoCredito = params.DebitoCredito
	}
}
