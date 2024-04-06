package cierrelote

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
)

type conciliarContracargo struct {
	utilService util.UtilService
}

func NewConciliarContracargo(util util.UtilService) MetodoConciliarClMP {
	return &conciliarContracargo{
		utilService: util,
	}
}

func (c *conciliarContracargo) ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	strNroEstablecimiento := strconv.Itoa(int(cierreLote.Nroestablecimiento))
	if cierreLote.Tipooperacion == "D" {
		if cierreLote.Fechaoperacion == movimientoDetalle.FechaOrigenCompra && strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && cierreLote.Nrotarjeta == movimientoDetalle.NroTarjetaXl && strings.Contains(movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion) && cierreLote.Nroticket == movimientoDetalle.NroCupon && valorCuota == movimientoDetalle.PlanCuota && cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe) && movimientoDetalle.TipoAplicacion == "-" && movimientoCabecera.Codop == movimientoDetalle.Tipooperacion.ExternalId {
			cierreLote.Enobservacion = ValidarArancel(c.utilService.ToFixed(cierreLote.Channelarancel.Importe*100, 2), movimientoDetalle.PorcentDescArancel/100)
			cierreLote.FechaPago = movimientoCabecera.FechaPago
			cierreLote.Cantdias = int(cierreLote.FechaPago.Sub(cierreLote.FechaCierre).Hours() / 24)
			cierreLote.PrismamovimientodetallesId = movimientoDetalle.Id
			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
			cabeceraMoviminetosIdArray = append(cabeceraMoviminetosIdArray, int64(movimientoDetalle.PrismamovimientototalesId))
			listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)
		}
	}
	if cierreLote.Tipooperacion == "C" {
		// cierreLote.Fechaoperacion == movimientoDetalle.FechaOrigenCompra &&
		if  strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && strings.Contains(movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion[1:len(cierreLote.Codigoautorizacion)]) && cierreLote.Nroticket == movimientoDetalle.NroCupon && cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe) && movimientoDetalle.TipoAplicacion == "-" {
			var tipoContraCargo string
			if movimientoDetalle.ContracargoOrigen == "E" {
				tipoContraCargo = fmt.Sprint("Contra cargo visa")
			}

			if movimientoDetalle.ContracargoOrigen == " " {
				tipoContraCargo = fmt.Sprint("Contra cargo mastercard")
			}

			cierreLote.Enobservacion = ValidarArancel(c.utilService.ToFixed(cierreLote.Channelarancel.Importe*100, 2), movimientoDetalle.PorcentDescArancel/100)
			cierreLote.Descripcioncontracargo = fmt.Sprintf("%v fehca pago: %v, se modificaco por fecha reversion", tipoContraCargo, cierreLote.FechaPago)
			cierreLote.FechaPago = movimientoCabecera.FechaPago
			cierreLote.DetallemovimientoId = movimientoDetalle.Id
			cierreLote.Reversion = true
			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
			cabeceraMoviminetosIdArray = append(cabeceraMoviminetosIdArray, int64(movimientoDetalle.PrismamovimientototalesId))
			listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)

		}
	}

	return
}

func ValidarArancel(porcentajeArancelControl, porcentajeArancelPrisma float64) bool {
	if porcentajeArancelControl != porcentajeArancelPrisma {
		return true
	}
	return false
}

// porcentajeArancelControl := c.utilService.ToFixed( cierreLote.Channelarancel.Importe * 100, 2)
// porcentajeArancelPrisma := movimientoDetalle.PorcentDescArancel / 100
// // controla el arancel informado por prisma con el que se encuentra registrado en nuestro sistema
// if porcentajeArancelControl != porcentajeArancelPrisma {
// 	// se pone en observacion para indicar que hay diferencia en el arancel
// 	cierreLote.Enobservacion = true
// }
