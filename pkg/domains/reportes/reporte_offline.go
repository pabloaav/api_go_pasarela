package reportes

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type offlineReportes struct {
	utilService util.UtilService
}

func NewReportesOffline(util util.UtilService) ReportesPagos {
	return &offlineReportes{
		utilService: util,
	}
}

func (cl *offlineReportes) ResponseReportes(s *reportesService, listaPagos reportedtos.TipoFactory) (response []reportedtos.ResponseFactory) {

	listaCierreLote, err := s.repository.GetCierreLoteOffline(listaPagos.TipoOffline)
	if err != nil {
		logs.Error(err)
	}

	for _, pago := range listaCierreLote {
		var fechaAcred string
		var arancelminimo float64
		// var netoCobraso entities.Monto
		if !pago.RapipagoCabecera.Fechaacreditacion.IsZero() {
			fechaAcred = pago.RapipagoCabecera.Fechaacreditacion.Format("02-01-2006")
		}
		if pago.RapipagoCabecera.Coeficiente > 0 {
			arancelminimo = pago.RapipagoCabecera.Coeficiente * 100
		}
		netoCobrado := entities.Monto(pago.ImporteCalculado)
		response = append(response, reportedtos.ResponseFactory{
			Pago:                    pago.CodigoBarras,
			FechaAcreditacion:       fechaAcred,
			Importeminimo:           pago.RapipagoCabecera.ImporteMinimo,
			Importemaximo:           pago.RapipagoCabecera.ImporteMaximo,
			ArancelPorcentajeMinimo: arancelminimo,
			// ArancelPorcentajeMaximo: pago.RapipagoCabecera.ArancelPorcentajeMaximo,
			ImporteNetoCobrado: s.util.ToFixed(netoCobrado.Float64(), 2),
		})
	}

	return
}
