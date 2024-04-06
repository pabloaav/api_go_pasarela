package administracionfake

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"gorm.io/gorm"
)

type TableDriverPrismaCierreLote struct {
	TituloPrueba                     string
	WantTable                        interface{}
	DataPruebaListaCierreLote        []entities.Prismacierrelote
	DataPruebaListaPagos             []entities.Pago
	DataFiltroPago                   filtros.PagoFiltro
	DataFiltroGetEstadoExterno       []string
	DataPruebaListaPagoEstadoExterno []entities.Pagoestadoexterno
}

type EstructuraValida struct {
	ListPrismaCierreLote []entities.Prismacierrelote
	ListPagos            []entities.Pago
	ListPagoEstadoLogs   []entities.Pagoestadologs
	Listmovimientos      []entities.Movimiento
	ListPagoIntentos     []entities.Pagointento
}

var (
	listaUuId  = []string{"3b1c9d8cd6966w4", "ce86bbf2d468aw4", "d820cfb2db8d8w4", "77c86bb9d8cf5w4"}
	filtroPago = filtros.PagoFiltro{
		Uuids:             listaUuId,
		CargaPagoIntentos: true,
		CargarPagoTipos:   true,
		CargarPagoEstado:  true}

	listaPagoEstadoExterno = []entities.Pagoestadoexterno{
		{Model: gorm.Model{ID: 12}, PagoestadosId: 7, Estado: "C", Vendor: "PRISMA"},
		{Model: gorm.Model{ID: 13}, PagoestadosId: 3, Estado: "A", Vendor: "PRISMA"},
		{Model: gorm.Model{ID: 14}, PagoestadosId: 5, Estado: "D", Vendor: "PRISMA"},
	}

	estructuraValida = EstructuraValida{
		ListPrismaCierreLote: listaCierrLotePrisma,
		ListPagos:            listaPagos,
		ListPagoEstadoLogs:   estructuraValidasEstadoLogs,
		Listmovimientos:      estructuraValidasMovimientos,
		ListPagoIntentos:     estructuraValidasPagoIntento,
	}
	/////////////////////////////////////////////Estructura de datos de Prueba/////////////////////////////////////////////////
	DataCierreloteFake = TableDriverPrismaCierreLote{
		TituloPrueba:              "verificar error al querer obter la lista de cierre de lote",
		WantTable:                 administracion.ERROR_OBTENER_CIERRE_LOTE,
		DataPruebaListaCierreLote: make([]entities.Prismacierrelote, 0),
	}
	DataPagosFake = TableDriverPrismaCierreLote{
		TituloPrueba:              "verificar error al querer obter la lista de Pagos",
		WantTable:                 administracion.ERROR_PAGO,
		DataPruebaListaCierreLote: estructuraValida.ListPrismaCierreLote,
		DataPruebaListaPagos:      make([]entities.Pago, 0),
		DataFiltroPago:            filtroPago,
	}

	GetpagoEstadoExternoFake = TableDriverPrismaCierreLote{
		TituloPrueba:                     "verificar error al querer obter la lista de Pagos estados externos",
		WantTable:                        administracion.ERROR_PAGO_ESTADO_EXTERNO,
		DataPruebaListaCierreLote:        estructuraValida.ListPrismaCierreLote,
		DataPruebaListaPagos:             estructuraValida.ListPagos,
		DataFiltroPago:                   filtroPago,
		DataFiltroGetEstadoExterno:       []string{"123", "Q", " ", "a"},
		DataPruebaListaPagoEstadoExterno: make([]entities.Pagoestadoexterno, 0),
	}

	GetpagoEstadoExternovalido = TableDriverPrismaCierreLote{
		TituloPrueba:                     "verificar que regrese todos los objetos",
		WantTable:                        estructuraValida,
		DataPruebaListaCierreLote:        estructuraValida.ListPrismaCierreLote,
		DataPruebaListaPagos:             estructuraValida.ListPagos,
		DataFiltroPago:                   filtroPago,
		DataFiltroGetEstadoExterno:       []string{"C", "C", "C", "A"},
		DataPruebaListaPagoEstadoExterno: listaPagoEstadoExterno,
	}
)
