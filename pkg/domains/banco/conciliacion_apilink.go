package banco

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtro "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type apilinkConciliacion struct {
	utilService    util.UtilService
	administracion administracion.Service
}

func NewApilinkConciliacion(util util.UtilService, administracion administracion.Service) MetodoConciliacionPagos {
	return &apilinkConciliacion{
		utilService:    util,
		administracion: administracion,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *apilinkConciliacion) FiltroRequestConsultaBanco(request bancodtos.RequestConciliacion) (response filtros.MovimientosBancoFiltro) {

	logs.Info("--------- Ejecutando proceso conciliacion apilink --------------- ")

	var movimientos []string
	var tipo filtros.EnumTipoOperacion
	if request.ListaApilink != nil {
		for _, referencia := range request.ListaApilink {
			movimientos = append(movimientos, referencia.ReferenciaBanco)
		}
		tipo = "debin"
	}
	response = filtros.MovimientosBancoFiltro{
		SubCuenta:      config.COD_SUBCUENTA,
		Tipo:           tipo,
		TipoMovimiento: movimientos,
	}

	return

}

func (cl *apilinkConciliacion) ConciliacionBanco(request bancodtos.RequestConciliacion, listaMovimientosBanco []bancodtos.ResponseMovimientosBanco) (lista bancodtos.ResponseConciliacion, movimientosIds []uint) {

	// /* 3 obtener estado pagoexterno aprobado y procesando para comprarar en con los registros de banco */
	filtroPagoEstado := filtro.PagoEstadoExternoFiltro{
		Nombre: config.PAGOEXTERNO_APROBADO,
	}

	pagoEstadoAcreditado, erro := cl.administracion.GetPagosEstadosExternoService(filtroPagoEstado)
	if erro != nil {
		return
	}
	var listaTemporal []*entities.Apilinkcierrelote
	// NOTE La conciliacion con banco comparamos la lista obtenida en banco con la cierreloteapilink
	if request.ListaApilink != nil {
		for _, clapilink := range request.ListaApilink {
			if clapilink.Pagoestadoexterno.PagoestadosId != pagoEstadoAcreditado[len(pagoEstadoAcreditado)-1].PagoestadosId {
				lista.ListaApilinkNoAcreditados = append(lista.ListaApilinkNoAcreditados, clapilink)
			} else {
				for _, mv := range listaMovimientosBanco {
					if clapilink.ReferenciaBanco == mv.DebinId && clapilink.Importe == entities.Monto(mv.Importe) {
						// lista para actualizar campo estado_check en movimientos del banco
						fechaAcreditacion, err := time.Parse("2006-01-02T00:00:00Z", mv.Fecha)
						if err != nil {
							logs.Error(err)
						}
						movimientosIds = append(movimientosIds, mv.Id)
						clapilink.Match = 1
						clapilink.BancoExternalId = int(mv.Id)
						clapilink.Fechaacreditacion = fechaAcreditacion
						listaTemporal = append(listaTemporal, clapilink)
					}
				}
			}
		}
		lista.ListaApilink = listaTemporal
	}

	return lista, movimientosIds
}
