package cierrelote

import (
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
)

type conciliarAjuste struct {
	utilService util.UtilService
}

func NewConciliarAjuste(util util.UtilService) MetodoConciliarClMP {
	return &conciliarAjuste{
		utilService: util,
	}
}

func (c *conciliarAjuste) ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	var listCabeceraMoviminetosId []int64
	strNroEstablecimiento := strconv.Itoa(int(cierreLote.Nroestablecimiento))
	if cierreLote.Fechaoperacion == movimientoDetalle.FechaOrigenCompra && cierreLote.FechaCierre == movimientoCabecera.FechaPresentacion && strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && cierreLote.Nrotarjeta == movimientoDetalle.NroTarjetaXl && strings.Contains(movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion) && cierreLote.Nroticket == movimientoDetalle.NroCupon && valorCuota == movimientoDetalle.PlanCuota && cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe) && movimientoDetalle.TipoAplicacion == "+" && movimientoCabecera.Codop == movimientoDetalle.Tipooperacion.ExternalId {
		porcentajeArancelControl := c.utilService.ToFixed(cierreLote.Channelarancel.Importe*100, 2)
		porcentajeArancelPrisma := movimientoDetalle.PorcentDescArancel / 100
		if porcentajeArancelControl != porcentajeArancelPrisma {
			cierreLote.Enobservacion = true
		}
		cierreLote.FechaPago = movimientoCabecera.FechaPago

		cierreLote.Cantdias = int(cierreLote.FechaPago.Sub(cierreLote.FechaCierre).Hours() / 24)
		cierreLote.PrismamovimientodetallesId = movimientoDetalle.Id
		detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
		listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(movimientoDetalle.PrismamovimientototalesId))
		listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)
	}
	// logs.Info(detalleMoviminetosIdArray)
	// logs.Info(listCabeceraMoviminetosId)
	return
}
