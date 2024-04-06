package reportes

import (
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

type prismaReportes struct {
	utilService util.UtilService
}

func NewReportesPrisma(util util.UtilService) ReportesPagos {
	return &prismaReportes{
		utilService: util,
	}
}

func (cl *prismaReportes) ResponseReportes(s *reportesService, listaPagos reportedtos.TipoFactory) (response []reportedtos.ResponseFactory) {

	listaCierreLote, err := s.repository.GetCierreLotePrisma(listaPagos.TipoPrisma)
	if err != nil {
		logs.Error(err)
	}
	var arrayNroLiquidacion []string
	for _, pago := range listaCierreLote {

		var nro_liquidacion string
		var nro_establecimiento string
		var fechaPres string
		var fechaAcred string
		var arancelPorcentaje float64
		var retencionIva string
		var importeArancel float64    // solo prisma  tabla prismamovimientosdetalles ->arancel
		var importeArancelIva float64 // solo prisma  tabla prismatrdospagos   ->retencion_iva
		var ImporteCft float64        // solo prisma  tablaprismamoviemientosdetalles ->importe_costo_financiero
		var importeNetoCobrado float64
		var reversion bool
		if pago.PrismamovimientodetallesId > 0 {
			nro_establecimiento = pago.Prismamovimientodetalle.MovimientoCabecera.EstablecimientoNro
			nro_liquidacion = pago.Prismamovimientodetalle.NroLiquidacion
			fechaPres = pago.Prismamovimientodetalle.MovimientoCabecera.FechaPresentacion.Format("02-01-2006")
			fechaAcred = pago.Prismamovimientodetalle.FechaPago.Format("02-01-2006")
			arancelPorcentaje = pago.Prismamovimientodetalle.PorcentDescArancel / 100
			retencionIva = pago.Prismamovimientodetalle.IdCf
			importeArancel = float64(pago.Prismamovimientodetalle.Arancel) / 100
			ImporteCft = float64(pago.Prismamovimientodetalle.ImporteCostoFinanciero) / 100
		}
		if pago.PrismatrdospagosId > 0 {
			arrayTemporal, estado := existeNroLiquidacion(pago.Prismatrdospagos.LiquidacionNro, arrayNroLiquidacion)
			if !estado {
				arrayNroLiquidacion = arrayTemporal
				importeArancelIva = pago.Prismatrdospagos.RetencionIvaD1.Float64()
				importeNetoCobrado = pago.Prismatrdospagos.ImporteNeto.Float64()
			}

		}
		if pago.Disputa {
			reversion = true
		}

		response = append(response, reportedtos.ResponseFactory{
			Pago:                 pago.ExternalclienteID,
			NroEstablecimiento:   nro_establecimiento,
			NroLiquidacion:       nro_liquidacion,
			FechaPresentacion:    fechaPres,
			FechaAcreditacion:    fechaAcred,
			ArancelPorcentaje:    arancelPorcentaje,
			RetencionIva:         retencionIva,
			ImporteArancel:       importeArancel,
			ImporteArancelIva:    importeArancelIva,
			ImporteArancalIvaMov: pago.Importeivaarancel,
			ImporteCft:           ImporteCft,
			ImporteNetoCobrado:   importeNetoCobrado,
			Revertido:            reversion,
			Enobservacion:        pago.Enobservacion,
			Cantdias:             pago.Cantdias,
		})
	}

	return
}

func existeNroLiquidacion(nroLiquidacionBuscar string, arrayNroLiquidacion []string) (array []string, estado bool) {
	if len(arrayNroLiquidacion) <= 0 {
		array = append(array, nroLiquidacionBuscar)
		estado = false
		return
	}
	for _, value := range arrayNroLiquidacion {
		if strings.Contains(value, nroLiquidacionBuscar) {
			array = arrayNroLiquidacion
			estado = true
			break
		}
	}
	if !estado {
		arrayNroLiquidacion = append(arrayNroLiquidacion, nroLiquidacionBuscar)
		array = arrayNroLiquidacion
		estado = false
	}
	return
}
