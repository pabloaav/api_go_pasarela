package banco

import (
	"fmt"

	adm "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type ConciliacionFactory interface {
	GetProcesarConciliacion(m int64) (MetodoConciliacionPagos, error)
}

type procesarConciliacionFactory struct{}

func NewRecorrerArchivos() ConciliacionFactory {
	return &procesarConciliacionFactory{}
}

func (r *procesarConciliacionFactory) GetProcesarConciliacion(m int64) (MetodoConciliacionPagos, error) {
	switch m {
	case 1:
		return NewRapipagoConciliacion(util.Resolve()), nil
	case 2:
		return NewApilinkConciliacion(util.Resolve(), adm.Resolve()), nil
	case 3:
		return NewTransferenciaConciliacion(util.Resolve()), nil
	default:
		return nil, fmt.Errorf("el tipo de conciliacion  %v, no es valido", m)

	}
}

type MetodoConciliacionPagos interface { // , rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
	FiltroRequestConsultaBanco(request bancodtos.RequestConciliacion) (response filtros.MovimientosBancoFiltro)
	ConciliacionBanco(request bancodtos.RequestConciliacion, listaMovimientosBanco []bancodtos.ResponseMovimientosBanco) (response bancodtos.ResponseConciliacion, movimientosIds []uint)
}
