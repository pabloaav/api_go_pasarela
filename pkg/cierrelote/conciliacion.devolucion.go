package cierrelote

import (
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
)

type conciliarDevolucion struct {
	utilService util.UtilService
}

func NewConciliarDevolucion(util util.UtilService) MetodoConciliarClMP {
	return &conciliarDevolucion{
		utilService: util,
	}
}

func (c *conciliarDevolucion) ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	if cierreLote.Tipooperacion == "D" {
		strNroEstablecimiento := strconv.Itoa(int(cierreLote.Nroestablecimiento))
		if cierreLote.FechaCierre == movimientoCabecera.FechaPresentacion && strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && cierreLote.ExternalloteId == movimientoDetalle.Lote && cierreLote.Nrotarjeta == movimientoDetalle.NroTarjetaXl && cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe) && movimientoDetalle.TipoAplicacion == "-" && movimientoCabecera.Codop == movimientoDetalle.Tipooperacion.ExternalId {

			porcentajeArancelControl := c.utilService.ToFixed( cierreLote.Channelarancel.Importe * 100, 2)
			porcentajeArancelPrisma := movimientoDetalle.PorcentDescArancel / 100
			if porcentajeArancelControl != porcentajeArancelPrisma {
				cierreLote.Enobservacion = true
			}
			cierreLote.FechaPago = movimientoCabecera.FechaPago
			cierreLote.Cantdias = int(cierreLote.FechaPago.Sub(cierreLote.FechaCierre).Hours() / 24)

			cierreLote.PrismamovimientodetallesId = movimientoDetalle.Id
			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
			cabeceraMoviminetosIdArray = append(cabeceraMoviminetosIdArray, int64(movimientoDetalle.PrismamovimientototalesId))
			listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)
		}
	}
	return
}
