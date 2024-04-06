package banco

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type transferenciaConciliacion struct {
	utilService util.UtilService
}

func NewTransferenciaConciliacion(util util.UtilService) MetodoConciliacionPagos {
	return &transferenciaConciliacion{
		utilService: util,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *transferenciaConciliacion) FiltroRequestConsultaBanco(request bancodtos.RequestConciliacion) (response filtros.MovimientosBancoFiltro) {

	logs.Info("--------- Ejecutando proceso conciliacion transferencias --------------- ")
	var movimientos []string
	var tipo filtros.EnumTipoOperacion
	var dbcr string
	// para cada transferencia, se filtra solo las ReferenciaBanco de cada transferencia
	if request.Transferencias.Transferencias != nil {

		for _, referencia := range request.Transferencias.Transferencias {
			// se adapta cada referencia banco al formato que se guarda en la tabla banco.movimientos:
			//  fecha + ultimos_cuatro_digitos de la referencia banco
			referencia_banco, erro := cl.utilService.CrearReferenciaBancoStandard(referencia.ReferenciaBanco)
			if erro != nil {
				logs.Error("error en FiltroRequestConsultaBanco: " + erro.Error())
			}
			movimientos = append(movimientos, referencia_banco)
		}

		tipo = "transferencia"
		dbcr = "1" // buscar transferencias salientes
	}
	response = filtros.MovimientosBancoFiltro{
		SubCuenta:      config.COD_SUBCUENTA,
		Tipo:           tipo,
		TipoMovimiento: movimientos,
		Dbcr:           dbcr,
	}

	return

}

func (cl *transferenciaConciliacion) ConciliacionBanco(request bancodtos.RequestConciliacion, listaMovimientosBanco []bancodtos.ResponseMovimientosBanco) (lista bancodtos.ResponseConciliacion, movimientosIds []uint) {

	if request.Transferencias.Transferencias != nil {
		for _, transferencia := range request.Transferencias.Transferencias {
			// se formatea correctamente la referencia banco antes de comparar
			referencia_banco, erro := cl.utilService.CrearReferenciaBancoStandard(transferencia.ReferenciaBanco)
			if erro != nil {
				logs.Error("error en ConciliacionBanco: " + erro.Error())
			}
			for _, itemMovimientoBanco := range listaMovimientosBanco {

				// se comparan las dos listas mediante el atributo referencia banco
				if referencia_banco == itemMovimientoBanco.ReferenciaTransferencia {
					// lista para actualizar campo estado_check en movimientos del banco
					movimientosIds = append(movimientosIds, itemMovimientoBanco.Id)
					// lista para actualizar campos match y banco_external_id en tabla pasarela.transferencias
					transferencia.Match = 1
					transferencia.BancoExternalId = int(itemMovimientoBanco.Id)
					// en la variable lista, atributo TransferenciasConciliadas, de tipo bancodtos.ResponseConciliacion se cargan transferencias-movimientos de la tabla pasarela.transferencias que se deben actualizar
					lista.TransferenciasConciliadas = append(lista.TransferenciasConciliadas, bancodtos.TransferenciasConciliadasConBanco{
						ListaIdsTransferenciasConciliadas: transferencia.ListaIdsTransferenciasAgrupadas,
						BancoExternalId:                   transferencia.BancoExternalId,
						Match:                             transferencia.Match,
					})

				} // Fin del if
			} // Fin for _, itemMovimientoBanco
		} // Fin _, transferencia
	}

	return lista, movimientosIds

}
