package administracionfake

import filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"

func EstructuraConsultarPagosFake() (tableDriverTestPeyment TableDriverTestConsultarPagos) {
	tableDriverTestPeyment = TableDriverTestConsultarPagos{
		TituloPrueba: "validar datos al consultar pagos, los valores enviados no son correctos",
		WantTable:    ERROR_DATOS_REQUEST,
		Request: filtros.PagoFiltro{
			// SubCuenta:      "123",
			// Tipo:           "debi",
			// TipoMovimiento: []string{"1"},
			// Fecha:          "2020-01-01",
			PagoEstadosIds:       []uint64{2, 4},
			VisualizarPendientes: false,
			CargaPagoIntentos:    true,
			CargaMedioPagos:      false,
		},
	}
	return
}
