package reportes

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

type debinReportes struct {
	utilService util.UtilService
}

func NewReportesDebin(util util.UtilService) ReportesPagos {
	return &debinReportes{
		utilService: util,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *debinReportes) ResponseReportes(s *reportesService, listaPagos reportedtos.TipoFactory) (response []reportedtos.ResponseFactory) {

	listaCierreLote, err := s.repository.GetCierreLoteApilink(listaPagos.TipoApilink)
	if err != nil {
		logs.Error(err)
	}

	for _, pago := range listaCierreLote {
		var fechaAcred string
		if !pago.Fechaacreditacion.IsZero() {
			fechaAcred = pago.Fechaacreditacion.Format("02-01-2006")
		}
		response = append(response, reportedtos.ResponseFactory{
			Pago:               pago.DebinId,
			FechaAcreditacion:  fechaAcred,
			ImporteNetoCobrado: pago.Importe.Float64(),
		})
	}

	return
}
