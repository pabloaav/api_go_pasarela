package cierrelote

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
)

const (
	RECHAZO_CONSUMO_PESOS                          = "0003"
	RECHAZO_DEVOLUCION_CONSUMO_PESOS               = "0004"
	CONSUMO_PESOS                                  = "0005"
	REVERSO_CONTRACARGO_PESOS                      = "1317"
	CONTRACARGO_PESOS1                             = "1507"
	CONTRACARGO_PESOS2                             = "1517"
	DEVOLUCIÓN_CONSUMO_PESOS                       = "6000"
	AJUSTE_DEBITO_ARANCEL_PESOS                    = "9100"
	AJUSTE_CREDITO_ARANCEL_PESOS                   = "9200"
	AJUSTE_DEBITO_DESCUENTO_VENTA_CAMPANIAS_PESOS  = "9300"
	AJUSTE_CREDITO_DESCUENTO_VENTA_CAMPANIAS_PESOS = "9400"
	AJUSTE_DEBITO_SERVICIO_CF_PESOS                = "9500"
	AJUSTE_CREDITO_SERVICIO_CF_PESOS               = "9600"
)

type ConciliarClMovimientosFactory interface {
	GetTipoConciliacion(m string) (MetodoConciliarClMP, error)
}

type conciliarClMovimientosFactory struct{}

func NewConciliarClMoviminetos() ConciliarClMovimientosFactory {
	return &conciliarClMovimientosFactory{}
}

func (c *conciliarClMovimientosFactory) GetTipoConciliacion(m string) (MetodoConciliarClMP, error) {

	switch m {
	case CONSUMO_PESOS:
		return NewConciliarCompras(util.Resolve()), nil
	case REVERSO_CONTRACARGO_PESOS:
		return NewConciliarReverso(util.Resolve()), nil
	case CONTRACARGO_PESOS1, CONTRACARGO_PESOS2:
		return NewConciliarContracargo(util.Resolve()), nil
	case DEVOLUCIÓN_CONSUMO_PESOS:
		return NewConciliarDevolucion(util.Resolve()), nil
	case AJUSTE_DEBITO_ARANCEL_PESOS, AJUSTE_CREDITO_ARANCEL_PESOS, AJUSTE_DEBITO_DESCUENTO_VENTA_CAMPANIAS_PESOS, AJUSTE_CREDITO_DESCUENTO_VENTA_CAMPANIAS_PESOS, AJUSTE_DEBITO_SERVICIO_CF_PESOS, AJUSTE_CREDITO_SERVICIO_CF_PESOS:
		return NewConciliarAjuste(util.Resolve()), nil
	default:
		return nil, fmt.Errorf("tipo de operacion a conciliar %v, no es valido", m)

	}

}

type MetodoConciliarClMP interface {
	ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error)
}
