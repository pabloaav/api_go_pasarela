package banco

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type rapipagoConciliacion struct {
	utilService util.UtilService
}

func NewRapipagoConciliacion(util util.UtilService) MetodoConciliacionPagos {
	return &rapipagoConciliacion{
		utilService: util,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *rapipagoConciliacion) FiltroRequestConsultaBanco(request bancodtos.RequestConciliacion) (response filtros.MovimientosBancoFiltro) {

	logs.Info("--------- Ejecutando proceso conciliacion rapipago --------------- ")
	var fechas []string
	var tipo filtros.EnumTipoOperacion
	if request.ListaRapipago != nil {
		for _, referencia := range request.ListaRapipago {
			fecha_acreditacion := referencia.Fechaacreditacion.Format("2006-01-02")
			fechas = append(fechas, fecha_acreditacion)
		}
		tipo = "rapipago"
	}
	response = filtros.MovimientosBancoFiltro{
		SubCuenta: config.COD_SUBCUENTA,
		Tipo:      tipo,
		Fechas:    fechas, // este filtro debo aplicar cuando tenga las fechas precisas (cuando el dinero ingresara al banco)
	}
	logs.Info(fechas)

	return
}

func (cl *rapipagoConciliacion) ConciliacionBanco(request bancodtos.RequestConciliacion, listaMovimientosBanco []bancodtos.ResponseMovimientosBanco) (lista bancodtos.ResponseConciliacion, movimientosIds []uint) {

	// se debe verificar la fecha y el monto que ingresaron al banco
	// si el monto no coinicide este pago pasara a estar en observacion , se debe verificar el el arancel (proveedor)
	for _, clrapipago := range request.ListaRapipago {

		logs.Info(clrapipago.Fechaacreditacion.Format("2006-01-02T00:00:00Z"))
		// 2022-08-02 00:00:00 +0000 UTC format fecha para comparar
		for _, mv := range listaMovimientosBanco {
			logs.Info(mv.Fecha)

			if clrapipago.Fechaacreditacion.Format("2006-01-02T00:00:00Z") == mv.Fecha {
				// los id de banco para luego actualizar (match en banco)
				movimientosIds = append(movimientosIds, mv.Id)
				// cabecera cierrelote
				clrapipago.BancoExternalId = int64(mv.Id)
				// & calcula la diferencia en monto del dinero que ingreso al banco con el monto calculado en el cierrelote 
				clrapipago.Difbancocl = float64(mv.Importe) - clrapipago.ImporteTotalCalculado
				if entities.Monto(clrapipago.ImporteTotalCalculado) != entities.Monto(mv.Importe) {
					// en el caso de que no coincidan los importes, se debera analizar este caso
					// arancel mal cobrado o calculado
					clrapipago.Enobservacion = true
				}

				for _, cldetalle := range clrapipago.RapipagoDetalle {
					cldetalle.Match = true
					if clrapipago.Enobservacion {
						cldetalle.Enobservacion = true
					}
				}
			}
			// if entities.Monto(clrapipago.ImporteTotalCalculado) == entities.Monto(mv.Importe) {

			// 	// los id de banco para luego actualizar (match en banco)
			// 	movimientosIds = append(movimientosIds, mv.Id)
			// 	// cabecera cierrelote
			// 	clrapipago.BancoExternalId = int64(mv.Id)

			// 	for _, cldetalle := range clrapipago.RapipagoDetalle {
			// 		cldetalle.Match = true
			// 	}
			// }
		}
	}
	lista.ListaRapipago = request.ListaRapipago

	return
}
