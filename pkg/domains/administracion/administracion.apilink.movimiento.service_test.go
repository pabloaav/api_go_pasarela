package administracion_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/administracion/administracionfake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

type GetPagoIntento struct {
	filtro    filtros.PagoIntentoFiltro
	respuesta []entities.Pagointento
	erro      error
}

type GetPagoEstado struct {
	filtro    filtros.PagoEstadoFiltro
	respuesta entities.Pagoestado
	erro      error
}

type GetPagoChannel struct {
	filtro    filtros.ChannelFiltro
	respuesta entities.Channel
	erro      error
}

type TableBuildMovimientoApiLink struct {
	nombre         string
	erro           error
	listaCierre    []*entities.Apilinkcierrelote
	getPagoIntento GetPagoIntento
	getPagoEstado  GetPagoEstado
	GetPagoChannel GetPagoChannel
}

func TestBuildCierreLoteApilinkService(t *testing.T) {

	filtroPagosEstados := filtros.PagoEstadoFiltro{
		BuscarPorFinal: true,
		Final:          false,
		Nombre:         "Processing",
	}

	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{"debin"},
	}

	channelValido := entities.Channel{
		Channel: "debin",
	}

	channelValido.ID = 4

	logPagoEstado := entities.Log{
		Tipo:          entities.Warning,
		Funcionalidad: "BuildCierreLoteApiLinkService",
		Mensaje:       administracion.ERROR_PAGO_ESTADO_ID,
	}

	logChannel := entities.Log{
		Tipo:          entities.Warning,
		Funcionalidad: "BuildCierreLoteApiLinkService",
		Mensaje:       administracion.ERROR_CHANNEL_ID,
	}

	filtroPagos := filtros.PagoFiltro{
		PagoEstadosId:     uint64(administracionfake.PagoEstadoValido().ID),
		CargaPagoIntentos: true,
		CargaMedioPagos:   true,
		CargarCuenta:      true,
	}

	uuid := uuid.NewV4().String()

	pagosProcessing := administracionfake.ListaPagosProcessing()

	requestInvalido := linkdebin.RequestGetDebinesLink{
		Pagina:      1,
		Tamanio:     linkdtos.Vacio,
		Cbu:         "0110599520000003855199",
		EsComprador: false,
		FechaDesde:  pagosProcessing[0].CreatedAt,
		FechaHasta:  pagosProcessing[len(pagosProcessing)-1].CreatedAt,
		Tipo:        linkdtos.DebinDefault,
		Estado:      "",
	}

	filtroPagoEstadoExternos := filtros.PagoEstadoExternoFiltro{
		Vendor: "APILINK",
	}

	response := administracionfake.ResponseGetDebinesValido()
	pagosEstadosExternos := administracionfake.ListaPagosExternos()
	var listaCierre []*entities.Apilinkcierrelote

	for j := range response.Debines {
		cierre := entities.Apilinkcierrelote{
			Uuid:            uuid,
			DebinId:         response.Debines[j].Id,
			Concepto:        response.Debines[j].Concepto,
			Moneda:          response.Debines[j].Moneda,
			Importe:         entities.Monto(response.Debines[j].Importe),
			Estado:          response.Debines[j].Estado,
			Tipo:            response.Debines[j].Tipo,
			FechaExpiracion: response.Debines[j].FechaExpiracion,
			Devuelto:        response.Debines[j].Devuelto,
			ContracargoId:   response.Debines[j].ContraCargoId,
			CompradorCuit:   response.Debines[j].Comprador.Cuit,
			VendedorCuit:    response.Debines[j].Vendedor.Cuit,
		}

		for i := range pagosEstadosExternos {

			if linkdtos.EnumEstadoDebin(pagosEstadosExternos[i].Estado) == response.Debines[j].Estado {
				cierre.PagoestadoexternosId = uint64(pagosEstadosExternos[i].ID)
				cierre.Pagoestadoexterno = pagosEstadosExternos[i]
				break
			}
		}

		listaCierre = append(listaCierre, &cierre)

	}

	t.Run("Debe Retornar un error si no se pudo cargar el pago estado Processing", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)
		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(entities.Pagoestado{}, fmt.Errorf(administracion.ERROR_PAGO_ESTADO))

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockRepository.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_PAGO_ESTADO, err.Error())

	})

	t.Run("Debe Retornar un error si el pago estado id es menor que 1", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(entities.Pagoestado{}, nil)
		mockRepository.On("CreateLog", logPagoEstado).Return(nil)

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockRepository.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_PAGO_ESTADO_ID, err.Error())

	})

	t.Run("Debe Retornar un error si no encuentra el canal debin", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(entities.Channel{}, fmt.Errorf(administracion.ERROR_CHANNEL))

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockRepository.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_CHANNEL, err.Error())

	})

	t.Run("Debe Retornar un error si el canal debin es vacio", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(entities.Channel{}, nil)
		mockRepository.On("CreateLog", logChannel).Return(nil)

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockRepository.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_CHANNEL_ID, err.Error())

	})

	t.Run("Debe Retornar un error si no puede consultar los pagos", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(channelValido, nil)
		mockRepository.On("GetPagos", filtroPagos).Return([]entities.Pago{}, fmt.Errorf(administracion.ERROR_PAGO))

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockRepository.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_PAGO, err.Error())

	})

	t.Run("Debe Retornar una lista vacia en caso de que no encuentre debines en apilink", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(channelValido, nil)
		mockRepository.On("GetPagos", filtroPagos).Return(pagosProcessing, nil)
		//En este caso creo un error cualquiera porque me interesa saber que pasa si ocurre un error y no probar la funcion
		mockCommonsService.On("NewUUID").Return(uuid)
		mockApiLinkService.On("GetDebinesApiLinkService", uuid, requestInvalido).Return(&linkdebin.ResponseGetDebinesLink{}, nil)

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		cierre, err := service.BuildCierreLoteApiLinkService()

		mockApiLinkService.AssertExpectations(t)
		mockRepository.AssertExpectations(t)
		mockCommonsService.AssertExpectations(t)

		var cierreLote []*entities.Apilinkcierrelote
		assert.Equal(t, cierreLote, cierre)
		assert.Equal(t, nil, err)

	})

	t.Run("Debe Retornar un error si no encuentra la lista de pagos externos", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(channelValido, nil)
		mockRepository.On("GetPagos", filtroPagos).Return(pagosProcessing, nil)
		//En este caso creo un error cualquiera porque me interesa saber que pasa si ocurre un error y no probar la funcion
		mockCommonsService.On("NewUUID").Return(uuid)
		mockApiLinkService.On("GetDebinesApiLinkService", uuid, requestInvalido).Return(administracionfake.ResponseGetDebinesValido(), nil)
		mockRepository.On("GetPagosEstadosExternos", filtroPagoEstadoExternos).Return([]entities.Pagoestadoexterno{}, fmt.Errorf(administracion.ERROR_PAGO_ESTADO_EXTERNO))

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockApiLinkService.AssertExpectations(t)
		mockRepository.AssertExpectations(t)
		mockCommonsService.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_PAGO_ESTADO_EXTERNO, err.Error())

	})

	t.Run("Debe Retornar un error si no se puede guardar el cierre de lote", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(channelValido, nil)
		mockRepository.On("GetPagos", filtroPagos).Return(pagosProcessing, nil)
		//En este caso creo un error cualquiera porque me interesa saber que pasa si ocurre un error y no probar la funcion
		mockCommonsService.On("NewUUID").Return(uuid)
		mockApiLinkService.On("GetDebinesApiLinkService", uuid, requestInvalido).Return(response, nil)
		mockRepository.On("GetPagosEstadosExternos", filtroPagoEstadoExternos).Return(pagosEstadosExternos, nil)
		mockRepository.On("CreateCierreLoteApiLink", listaCierre).Return(fmt.Errorf(administracion.ERROR_CREAR_CIERRE_LOTE))

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockApiLinkService.AssertExpectations(t)
		mockRepository.AssertExpectations(t)
		mockCommonsService.AssertExpectations(t)

		assert.Equal(t, administracion.ERROR_CREAR_CIERRE_LOTE, err.Error())

	})

	t.Run("Debe retornar una lista de cierres de lote", func(t *testing.T) {

		mockRepository := new(mockrepository.MockRepositoryAdministracion)
		mockApiLinkService := new(mockservice.MockApiLinkService)
		mockCommonsService := new(mockservice.MockCommonsService)
		mockUtilService := new(mockservice.MockUtilService)
		mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
		mockStore := new(mockservice.MockStoreService)

		mockRepository.On("GetPagoEstado", filtroPagosEstados).Return(administracionfake.PagoEstadoValido(), nil)
		mockRepository.On("GetChannel", filtroChannel).Return(channelValido, nil)
		mockRepository.On("GetPagos", filtroPagos).Return(pagosProcessing, nil)
		//En este caso creo un error cualquiera porque me interesa saber que pasa si ocurre un error y no probar la funcion
		mockCommonsService.On("NewUUID").Return(uuid)
		mockApiLinkService.On("GetDebinesApiLinkService", uuid, requestInvalido).Return(response, nil)
		mockRepository.On("GetPagosEstadosExternos", filtroPagoEstadoExternos).Return(pagosEstadosExternos, nil)
		mockRepository.On("CreateCierreLoteApiLink", listaCierre).Return(nil)

		service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

		_, err := service.BuildCierreLoteApiLinkService()

		mockApiLinkService.AssertExpectations(t)
		mockRepository.AssertExpectations(t)
		mockCommonsService.AssertExpectations(t)

		assert.Equal(t, nil, err)
		assert.Equal(t, listaCierre[0].DebinId, "1idDebin")

	})

}

func TestBuildMovimientoApilink(t *testing.T) {

	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

	log, table := _inicializarBuildMovimientoApilink()

	for _, v := range table {

		t.Run(v.nombre, func(t *testing.T) {

			mockRepository.On("GetPagosIntentos", v.getPagoIntento.filtro).Return(v.getPagoIntento.respuesta, v.getPagoIntento.erro).Once()
			mockRepository.On("GetPagoEstado", v.getPagoEstado.filtro).Return(v.getPagoEstado.respuesta, v.getPagoEstado.erro).Once()
			mockRepository.On("CreateLog", log).Return(nil).Once()
			mockRepository.On("GetChannel", v.GetPagoChannel.filtro).Return(v.GetPagoChannel.respuesta, v.GetPagoChannel.erro).Once()

			movimientoCierreLote, err := service.BuildMovimientoApiLink(v.listaCierre)

			if err != nil {
				assert.Equal(t, v.erro.Error(), err.Error())
			}

			if err == nil {

				assert.NotNil(t, movimientoCierreLote.ListaPagoIntentos)
				assert.Equal(t, int64(2), movimientoCierreLote.ListaPagoIntentos[0].Pago.PagoestadosID)
				assert.NotNil(t, movimientoCierreLote.ListaPagosEstadoLogs)
				assert.Len(t, movimientoCierreLote.ListaPagosEstadoLogs, 10)
				assert.NotNil(t, movimientoCierreLote.ListaPagos)
				assert.Len(t, movimientoCierreLote.ListaPagos, 10)

				assert.Len(t, movimientoCierreLote.ListaMovimientos, 1)
			}

		})
	}

}

//Inicio TestBuildMovimientoApilink
func _inicializarBuildMovimientoApilink() (log entities.Log, table []TableBuildMovimientoApiLink) {
	cierresLotes := _listaCierreLotes()

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		ExternalIds:    []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
		Channel:        true,
		CargarPago:     true,
		CargarPagoTipo: true,
	}

	getPagoIntentoError := GetPagoIntento{
		filtro: filtroPagoIntento,
		erro:   fmt.Errorf(administracion.ERROR_PAGO_INTENTO),
	}
	getPagoIntentoInValido := GetPagoIntento{
		filtro:    filtroPagoIntento,
		respuesta: _listaPagosIntentosInvalidos(),
	}

	getPagoIntentoValido := GetPagoIntento{
		filtro:    filtroPagoIntento,
		respuesta: _listaPagosIntentosValidos(),
		erro:      nil,
	}

	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: "Accredited",
	}

	GetPagoEstadoError := GetPagoEstado{
		filtro: filtroPagoEstado,
		erro:   fmt.Errorf(administracion.ERROR_PAGO_ESTADO),
	}

	GetPagoEstadoValido := GetPagoEstado{
		filtro: filtroPagoEstado,
		respuesta: entities.Pagoestado{
			Model:  gorm.Model{ID: 7},
			Estado: entities.Accredited,
			Final:  true,
		},
	}

	log = entities.Log{
		Tipo:          entities.Warning,
		Funcionalidad: "BuildMovimientoApiLink",
		Mensaje:       "no se encontrarion los siguientes debines [0 1 2 3 4 5 6 7 8 9 10 11]",
	}

	filtroChannel := filtros.ChannelFiltro{
		Channels: []string{"debin"},
	}

	getPagoChannelError := GetPagoChannel{filtro: filtroChannel, respuesta: entities.Channel{}, erro: fmt.Errorf(administracion.ERROR_CHANNEL)}

	channelDEbin := entities.Channel{
		Model: gorm.Model{ID: 4},
	}

	getPagoChannel := GetPagoChannel{filtro: filtroChannel, respuesta: channelDEbin, erro: nil}

	table = []TableBuildMovimientoApiLink{
		{"Debe retornar un error si la lista de cierre es vacia", fmt.Errorf(administracion.ERROR_LISTA_CIERRE_LOTE), []*entities.Apilinkcierrelote{}, GetPagoIntento{}, GetPagoEstado{}, GetPagoChannel{}},
		{"Debe retornar un error si no se puede cargar el canal debin", fmt.Errorf(administracion.ERROR_CHANNEL), cierresLotes, GetPagoIntento{}, GetPagoEstado{}, getPagoChannelError},
		{"Debe retornar un error si no se puede cargar los pagos intentos", fmt.Errorf(administracion.ERROR_PAGO_INTENTO), cierresLotes, getPagoIntentoError, GetPagoEstado{}, getPagoChannel},
		{"Debe retornar un error si no se puede cargar es pago estado Accredited", fmt.Errorf(administracion.ERROR_PAGO_ESTADO), cierresLotes, getPagoIntentoInValido, GetPagoEstadoError, getPagoChannel},
		{"Debe retornar un error si la lista de cierre no tiene la longitud de la lista de pagos intentos", fmt.Errorf(administracion.ERROR_CIERRE_PAGO_INTENTO), cierresLotes, getPagoIntentoInValido, GetPagoEstadoValido, getPagoChannel},
		{"Debe retornar una lista de pagos, pagosintentos y pagosintentoslogs cuando existe al menos una modificacion", nil, cierresLotes, getPagoIntentoValido, GetPagoEstadoValido, getPagoChannel},
	}

	return

}

// Fakes
func _listaPagosEstados() []entities.Pagoestado {
	var lista = make([]entities.Pagoestado, 7)
	lista[0] = entities.Pagoestado{
		Model:  gorm.Model{ID: 1},
		Estado: entities.Pending,
		Final:  false,
	}
	lista[1] = entities.Pagoestado{
		Model:  gorm.Model{ID: 2},
		Estado: entities.Processing,
		Final:  false,
	}
	lista[2] = entities.Pagoestado{
		Model:  gorm.Model{ID: 3},
		Estado: entities.Rejected,
		Final:  true,
	}
	lista[3] = entities.Pagoestado{
		Model:  gorm.Model{ID: 4},
		Estado: entities.Paid,
		Final:  false,
	}
	lista[4] = entities.Pagoestado{
		Model:  gorm.Model{ID: 5},
		Estado: entities.Reverted,
		Final:  true,
	}
	lista[5] = entities.Pagoestado{
		Model:  gorm.Model{ID: 6},
		Estado: entities.Expired,
		Final:  true,
	}
	lista[6] = entities.Pagoestado{
		Model:  gorm.Model{ID: 7},
		Estado: entities.Accredited,
		Final:  true,
	}
	return lista
}

func _listaCierreLotes() []*entities.Apilinkcierrelote {

	pagosExternos := administracionfake.ListaPagosExternos()
	pagosEstados := _listaPagosEstados()

	var lista = make([]*entities.Apilinkcierrelote, len(pagosExternos))

	for i := range pagosExternos {
		estado := linkdtos.Acreditado

		for j := range pagosEstados {
			if uint64(pagosEstados[j].ID) == pagosExternos[i].PagoestadosId {
				estado = linkdtos.EnumEstadoDebin(pagosEstados[j].Estado)
			}

		}

		lista[i] = &entities.Apilinkcierrelote{
			Uuid:                 "",
			DebinId:              fmt.Sprint(i),
			Concepto:             linkdtos.Alquiler,
			Moneda:               linkdtos.Pesos,
			Importe:              entities.Monto(((1000.54 * float64(i)) + 0.03) * 100),
			Estado:               estado,
			Tipo:                 linkdtos.DebinDefault,
			FechaExpiracion:      time.Now(),
			Devuelto:             true,
			CompradorCuit:        "20785695147",
			VendedorCuit:         "20546951227",
			PagoestadoexternosId: uint64(pagosExternos[i].ID),
			Pagoestadoexterno:    pagosExternos[i],
		}
	}

	return lista

}

func _listaPagosIntentosInvalidos() []entities.Pagointento {
	pagosExternos := administracionfake.ListaPagosExternos()

	var lista = make([]entities.Pagointento, len(pagosExternos)-2)

	for i := 0; i < len(pagosExternos)-2; i++ {
		lista[i] = entities.Pagointento{
			ExternalID: fmt.Sprint(i),
			Amount:     entities.Monto((1000.54 * float64(i)) + 0.03),
			PagosID:    int64(i),
		}
	}
	return lista

}

func _listaPagosIntentosValidos() []entities.Pagointento {
	pagosExternos := administracionfake.ListaPagosExternos()

	var lista = make([]entities.Pagointento, len(pagosExternos))
	var pagos = make([]entities.Pago, len(pagosExternos))
	for i := 0; i < len(pagos); i++ {
		pagos[i] = entities.Pago{
			Model:         gorm.Model{ID: uint(i)},
			PagoestadosID: 2,
			PagostipoID:   1,
			PagosTipo:     entities.Pagotipo{CuentasID: 1},
		}
	}

	for i := 0; i < len(pagosExternos); i++ {
		lista[i] = entities.Pagointento{
			ExternalID: fmt.Sprint(i),
			Amount:     entities.Monto((1000.54 * float64(i)) + 0.03),
			PagosID:    int64(i),
			Mediopagos: entities.Mediopago{
				ChannelsID: 4,
			},
			Pago: pagos[i],
		}
	}
	return lista

}

//Final TestBuildMovimientoApilink

// Inicio TestBuildTransferenciaCliente
type GetMovimientos struct {
	Filtro    filtros.MovimientoFiltro
	Respuesta []entities.Movimiento
	Erro      error
}

type GetSaldoCuenta struct {
	Filtro    uint64
	Respuesta administraciondtos.SaldoCuentaResponse
	Erro      error
}

type BajaMovimiento struct {
	ListaMovimientos []*entities.Movimiento
	Mensaje          string
	Erro             error
}

type CreateTransferencias struct {
	listaTransferencias []*entities.Transferencia
	Erro                error
}

type CreateNotificacion struct {
	Notificacion entities.Notificacione
	Erro         error
}

type CreateLog struct {
	Request entities.Log
	Erro    error
}

type CreateMovimientosTransferencia struct {
	Filtro []*entities.Movimiento
	Erro   error
}

type CreateTransferenciaApiLinkService struct {
	Requerimiento string
	Transferencia linktransferencia.RequestTransferenciaCreateLink
	Respuesta     *linktransferencia.ResponseTransferenciaCreateLink
	Erro          error
}

type IsValidUUID struct {
	Filtro    string
	Respuesta bool
	Erro      error
}

type TableTest struct {
	Nombre                            string
	Request                           requestBuildTransferenciaCliente
	Error                             error
	IsValidUUID                       IsValidUUID
	GetMovimientos                    GetMovimientos
	CreateLog                         CreateLog
	GetPagoEstado                     GetPagoEstado
	GetSaldoCuenta                    GetSaldoCuenta
	CreateMovimientosTransferencia    CreateMovimientosTransferencia
	CreateTransferenciaApiLinkService CreateTransferenciaApiLinkService
	BajaMovimiento                    BajaMovimiento
	CreateNotificacion                CreateNotificacion
	CreateTransferencias              CreateTransferencias
}

type requestBuildTransferenciaCliente struct {
	RequerimientoId string
	Request         administraciondtos.RequestTransferenicaCliente
	CuentaId        uint64
	UserId          uint64
}

func _inicializarBuilTransferenciaCliente() []TableTest {

	// requestValido := requestBuildTransferenciaCliente{
	// 	RequerimientoId: uuid.NewV4().String(),
	// 	Request: administraciondtos.RequestTransferenicaCliente{
	// 		Transferencia: linktransferencia.RequestTransferenciaCreateLink{
	// 			Importe: 406,
	// 		},
	// 		ListaMovimientosId: []uint64{1, 2, 3, 4},
	// 	},
	// 	CuentaId: 1,
	// 	UserId:   1,
	// }

	requestRequerimientoInvalido := requestBuildTransferenciaCliente{
		RequerimientoId: "0",
		Request: administraciondtos.RequestTransferenicaCliente{
			Transferencia: linktransferencia.RequestTransferenciaCreateLink{
				Importe: 406,
			},
			ListaMovimientosId: []uint64{1, 2, 3, 4},
		},
		CuentaId: 1,
		UserId:   1,
	}

	// requestCuentaInvalida := requestBuildTransferenciaCliente{
	// 	RequerimientoId: uuid.NewV4().String(),
	// 	Request: administraciondtos.RequestTransferenicaCliente{
	// 		Transferencia: linktransferencia.RequestTransferenciaCreateLink{
	// 			Importe: 406,
	// 		},
	// 		ListaMovimientosId: []uint64{1, 2, 3, 4},
	// 	},
	// 	CuentaId: 0,
	// 	UserId:   1,
	// }

	// requestUserInvalido := requestBuildTransferenciaCliente{
	// 	RequerimientoId: uuid.NewV4().String(),
	// 	Request: administraciondtos.RequestTransferenicaCliente{
	// 		Transferencia: linktransferencia.RequestTransferenciaCreateLink{
	// 			Importe: 406,
	// 		},
	// 		ListaMovimientosId: []uint64{1, 2, 3, 4},
	// 	},
	// 	CuentaId: 1,
	// 	UserId:   0,
	// }

	// filtroMovimientos := filtros.MovimientoFiltro{Ids: requestValido.Request.ListaMovimientosId, CargarPago: true}

	// pagoEstadoAcreditado := entities.Pagoestado{
	// 	Model:  gorm.Model{ID: 7},
	// 	Estado: entities.Accredited,
	// }

	// listaPagosAcreditados := make([]entities.Pago, 4)
	// for i := 0; i < 4; i++ {
	// 	listaPagosAcreditados[i] = entities.Pago{
	// 		Model:         gorm.Model{ID: uint(i + 1)},
	// 		PagoestadosID: 7,
	// 		PagoEstados:   pagoEstadoAcreditado,
	// 		PagoIntentos: []entities.Pagointento{
	// 			{
	// 				Amount: (float64(i) + 1*100),
	// 			},
	// 		},
	// 	}
	// }

	// listaPagosNoAcreditados := make([]entities.Pago, 4)
	// for i := 0; i < 4; i++ {
	// 	listaPagosNoAcreditados[i] = entities.Pago{
	// 		Model:         gorm.Model{ID: uint(i + 1)},
	// 		PagoestadosID: int64(i + 5),
	// 		PagoEstados:   pagoEstadoAcreditado,
	// 		PagoIntentos: []entities.Pagointento{
	// 			{
	// 				Amount: (float64(i) + 1*100),
	// 			},
	// 		},
	// 	}
	// }

	// listaMovimientos := make([]entities.Movimiento, 4)
	// for i := 0; i < 4; i++ {
	// 	listaMovimientos[i] = entities.Movimiento{
	// 		Model:   gorm.Model{ID: uint(i + 1)},
	// 		PagointentosId: 1,
	// 		PagoIntentosId: uint64(listaPagosAcreditados[i].ID),
	// 		Pago:    listaPagosAcreditados[i],
	// 	}
	// }

	// listaMovimentosMenor := listaMovimientos[0:2]

	// listaMovimientoPagoNoAcreditado := make([]entities.Movimiento, 4)
	// for i := 0; i < 4; i++ {
	// 	listaMovimientoPagoNoAcreditado[i] = entities.Movimiento{
	// 		Model:   gorm.Model{ID: uint(i + 1)},
	// 		PagosId: uint64(listaPagosAcreditados[i].ID),
	// 		Pago:    listaPagosNoAcreditados[i],
	// 	}
	// }

	// log := entities.Log{
	// 	Tipo:          entities.Error,
	// 	Mensaje:       fmt.Sprintf("no se encontraron los siguientes movimientos seleccionados [3 4]"),
	// 	Funcionalidad: "BuildTransferenciaCliente",
	// }

	// FiltroPagoEstado := filtros.PagoEstadoFiltro{
	// 	BuscarPorFinal: true,
	// 	Final:          true,
	// 	Nombre:         config.MOVIMIENTO_ACCREDITED,
	// }

	// requestInvalidoImporte := requestBuildTransferenciaCliente{
	// 	RequerimientoId: uuid.NewV4().String(),
	// 	Request: administraciondtos.RequestTransferenicaCliente{
	// 		Transferencia: linktransferencia.RequestTransferenciaCreateLink{
	// 			Importe: 405,
	// 		},
	// 		ListaMovimientosId: []uint64{1, 2, 3, 4},
	// 	},
	// 	CuentaId: 1,
	// 	UserId:   1,
	// }

	// saldoCuentaInsuficiente := administraciondtos.SaldoCuentaResponse{
	// 	CuentasId: 1,
	// 	Total:     3,
	// }

	// saldoCuenta := administraciondtos.SaldoCuentaResponse{
	// 	CuentasId: 1,
	// 	Total:     1000,
	// }

	// requestCbuInvalido := requestBuildTransferenciaCliente{
	// 	RequerimientoId: uuid.NewV4().String(),
	// 	Request: administraciondtos.RequestTransferenicaCliente{
	// 		Transferencia: linktransferencia.RequestTransferenciaCreateLink{
	// 			Importe: 406,
	// 		},
	// 		ListaMovimientosId: []uint64{1, 2, 3, 4},
	// 	},
	// 	CuentaId: 1,
	// 	UserId:   1,
	// }

	// var movimientos []*entities.Movimiento

	// for i := range listaMovimientos {
	// 	movimiento := entities.Movimiento{
	// 		CuentasId: uint64(requestValido.CuentaId),
	// 		PagosId:   uint64(listaMovimientos[i].PagosId),
	// 		Monto:     listaMovimientos[i].Monto,
	// 		Tipo:      "C",
	// 	}
	// 	movimientos = append(movimientos, &movimiento)
	// }

	// notificacion := entities.Notificacione{
	// 	Tipo:        entities.NotificacionTransferencia,
	// 	Descripcion: fmt.Sprintf("atención los siguientes movimientos de transferencia se realizaron incorrectamente pero no pudieron ser cancelados, movimientosId: %s", "0,0,0,0,"),
	// }

	// responseCreateTransferencia := linktransferencia.ResponseTransferenciaCreateLink{
	// 	NumeroReferenciaBancaria: "152016",
	// 	FechaOperacion:           time.Now(),
	// }

	// var transferencias []*entities.Transferencia

	// for i := range movimientos {
	// 	transferencia := entities.Transferencia{
	// 		MovimientosID:      uint64(i + 1),
	// 		ReferenciaBancaria: responseCreateTransferencia.NumeroReferenciaBancaria,
	// 		UserId:             requestValido.UserId,
	// 		Uuid:               requestValido.RequerimientoId,
	// 	}
	// 	transferencias = append(transferencias, &transferencia)
	// }

	table := []TableTest{
		{
			"Debe retornar un error si el el formato del requerimientoId es envalido",
			requestRequerimientoInvalido,
			fmt.Errorf(commons.ERROR_UUID),
			IsValidUUID{
				Filtro:    requestRequerimientoInvalido.RequerimientoId,
				Respuesta: false,
				Erro:      fmt.Errorf(commons.ERROR_UUID),
			},
			GetMovimientos{},
			CreateLog{},
			GetPagoEstado{},
			GetSaldoCuenta{},
			CreateMovimientosTransferencia{},
			CreateTransferenciaApiLinkService{},
			BajaMovimiento{},
			CreateNotificacion{},
			CreateTransferencias{},
		},
		// {
		// 	"Debe retornar un error si la cuenta id es envalida",
		// 	requestCuentaInvalida,
		// 	fmt.Errorf(administracion.ERROR_CUENTA_ID),
		// 	IsValidUUID{
		// 		Filtro:    requestCuentaInvalida.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{},
		// 	CreateLog{},
		// 	GetPagoEstado{},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si el user id es invalido",
		// 	requestUserInvalido,
		// 	fmt.Errorf(administracion.ERROR_USER_ID),
		// 	IsValidUUID{
		// 		Filtro:    requestUserInvalido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{},
		// 	CreateLog{},
		// 	GetPagoEstado{},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no encuentra los movimientos",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_MOVIMIENTO),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: []entities.Movimiento{},
		// 		Erro:      fmt.Errorf(administracion.ERROR_MOVIMIENTO),
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si la cantidad de elementos de la lista de movimientos es distinta de la de ids de entrada",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_MOVIMIENTO_LISTA_DIFERENCIA),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimentosMenor,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{
		// 		Request: log,
		// 		Erro:    nil,
		// 	},
		// 	GetPagoEstado{},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no encuentra el estado acreditado",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_PAGO_ESTADO),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: entities.Pagoestado{},
		// 		erro:      fmt.Errorf(administracion.ERROR_PAGO_ESTADO),
		// 	},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si algun pago no esta en el estado acreditado",
		// 	requestValido,
		// 	fmt.Errorf("el movimiento %d que corresponde al pago %d no está acreditado", 1, 1),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientoPagoNoAcreditado,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si el total de los movimientos no corresponde con el importe a transferir",
		// 	requestInvalidoImporte,
		// 	fmt.Errorf(administracion.ERROR_IMPORTE_TRANSFERENCIA),
		// 	IsValidUUID{
		// 		Filtro:    requestInvalidoImporte.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no se puede recuperar el saldo de la cuenta",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_SALDO_CUENTA),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestValido.CuentaId,
		// 		Respuesta: administraciondtos.SaldoCuentaResponse{},
		// 		Erro:      fmt.Errorf(administracion.ERROR_SALDO_CUENTA),
		// 	},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si el saldo de la cuenta es insuficiente cuenta",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_SALDO_CUENTA_INSUFICIENTE),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestValido.CuentaId,
		// 		Respuesta: saldoCuentaInsuficiente,
		// 		Erro:      nil,
		// 	},
		// 	CreateMovimientosTransferencia{},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no se puede crear el movimiento",
		// 	requestCbuInvalido,
		// 	fmt.Errorf(administracion.ERROR_CREAR_MOVIMIENTOS),
		// 	IsValidUUID{
		// 		Filtro:    requestCbuInvalido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestCbuInvalido.CuentaId,
		// 		Respuesta: saldoCuenta,
		// 		Erro:      nil,
		// 	},
		// 	CreateMovimientosTransferencia{
		// 		Filtro: movimientos,
		// 		Erro:   fmt.Errorf(administracion.ERROR_CREAR_MOVIMIENTOS),
		// 	},
		// 	CreateTransferenciaApiLinkService{},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no se puede crear la transferencia en apilink",
		// 	requestValido,
		// 	fmt.Errorf(apilink.ERROR_CREATE_TRANSFERENCIA),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestValido.CuentaId,
		// 		Respuesta: saldoCuenta,
		// 		Erro:      nil,
		// 	},
		// 	CreateMovimientosTransferencia{
		// 		Filtro: movimientos,
		// 		Erro:   nil,
		// 	},
		// 	CreateTransferenciaApiLinkService{
		// 		Requerimiento: requestValido.RequerimientoId,
		// 		Transferencia: requestValido.Request.Transferencia,
		// 		Respuesta:     &linktransferencia.ResponseTransferenciaCreateLink{},
		// 		Erro:          fmt.Errorf(apilink.ERROR_CREATE_TRANSFERENCIA),
		// 	},
		// 	BajaMovimiento{
		// 		ListaMovimientos: movimientos,
		// 		Mensaje:          apilink.ERROR_CREATE_TRANSFERENCIA,
		// 		Erro:             fmt.Errorf(administracion.ERROR_BAJAR_MOVIMIENTOS),
		// 	},
		// 	CreateNotificacion{
		// 		Notificacion: notificacion,
		// 		Erro:         fmt.Errorf(administracion.ERROR_CREAR_NOTIFICACION),
		// 	},
		// 	CreateTransferencias{},
		// },
		// {
		// 	"Debe retornar un error si no se puede crear la transferencia",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_CREAR_TRANSFERENCIAS),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestValido.CuentaId,
		// 		Respuesta: saldoCuenta,
		// 		Erro:      nil,
		// 	},
		// 	CreateMovimientosTransferencia{
		// 		Filtro: movimientos,
		// 		Erro:   nil,
		// 	},
		// 	CreateTransferenciaApiLinkService{
		// 		Requerimiento: requestValido.RequerimientoId,
		// 		Transferencia: requestValido.Request.Transferencia,
		// 		Respuesta:     &responseCreateTransferencia,
		// 		Erro:          nil,
		// 	},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{
		// 		listaTransferencias: transferencias,
		// 		Erro:                fmt.Errorf(administracion.ERROR_CREAR_TRANSFERENCIAS),
		// 	},
		// },
		// {
		// 	"Debe retornar el numero de la transferencia en caso de succeso",
		// 	requestValido,
		// 	fmt.Errorf(administracion.ERROR_CREAR_TRANSFERENCIAS),
		// 	IsValidUUID{
		// 		Filtro:    requestValido.RequerimientoId,
		// 		Respuesta: true,
		// 		Erro:      nil,
		// 	},
		// 	GetMovimientos{
		// 		Filtro:    filtroMovimientos,
		// 		Respuesta: listaMovimientos,
		// 		Erro:      nil,
		// 	},
		// 	CreateLog{},
		// 	GetPagoEstado{
		// 		filtro:    FiltroPagoEstado,
		// 		respuesta: pagoEstadoAcreditado,
		// 		erro:      nil,
		// 	},
		// 	GetSaldoCuenta{
		// 		Filtro:    requestValido.CuentaId,
		// 		Respuesta: saldoCuenta,
		// 		Erro:      nil,
		// 	},
		// 	CreateMovimientosTransferencia{
		// 		Filtro: movimientos,
		// 		Erro:   nil,
		// 	},
		// 	CreateTransferenciaApiLinkService{
		// 		Requerimiento: requestValido.RequerimientoId,
		// 		Transferencia: requestValido.Request.Transferencia,
		// 		Respuesta:     &responseCreateTransferencia,
		// 		Erro:          nil,
		// 	},
		// 	BajaMovimiento{},
		// 	CreateNotificacion{},
		// 	CreateTransferencias{
		// 		listaTransferencias: transferencias,
		// 		Erro:                nil,
		// 	},
		// },
	}

	return table

}

func TestBuildTransferenciaCliente(t *testing.T) {

	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	defaultContext := context.Background()
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

	table := _inicializarBuilTransferenciaCliente()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {
			mockCommonsService.On("IsValidUUID", v.IsValidUUID.Filtro).Return(v.IsValidUUID.Respuesta, v.IsValidUUID.Erro).Once()
			mockRepository.On("GetMovimientos", v.GetMovimientos.Filtro).Return(v.GetMovimientos.Respuesta, v.GetMovimientos.Erro).Once()
			mockRepository.On("CreateLog", v.CreateLog.Request).Return(v.CreateLog.Erro).Once()
			mockRepository.On("GetPagoEstado", v.GetPagoEstado.filtro).Return(v.GetPagoEstado.respuesta, v.GetPagoEstado.erro).Once()
			mockRepository.On("GetSaldoCuenta", v.GetSaldoCuenta.Filtro).Return(v.GetSaldoCuenta.Respuesta, v.GetSaldoCuenta.Erro).Once()
			mockRepository.On("CreateMovimientosTransferencia", v.CreateMovimientosTransferencia.Filtro).Return(v.CreateMovimientosTransferencia.Erro).Once()

			mockApiLinkService.On("CreateTransferenciaApiLinkService", v.CreateTransferenciaApiLinkService.Requerimiento, v.CreateTransferenciaApiLinkService.Transferencia).Return(v.CreateTransferenciaApiLinkService.Respuesta, v.CreateTransferenciaApiLinkService.Erro).Once()

			if v.CreateTransferenciaApiLinkService.Erro != nil {
				mockRepository.On("BajaMovimiento", v.BajaMovimiento.ListaMovimientos, v.CreateTransferenciaApiLinkService.Erro.Error()).Return(v.BajaMovimiento.Erro).Once()
				mockRepository.On("CreateNotificacion", v.CreateNotificacion.Notificacion).Return(v.CreateNotificacion.Erro).Once()
			}
			if len(v.CreateTransferencias.listaTransferencias) > 0 {
				mockRepository.On("CreateTransferencias", v.CreateTransferencias.listaTransferencias).Return(v.CreateTransferencias.Erro).Once()
			}

			response, err := service.BuildTransferenciaCliente(defaultContext, v.Request.RequerimientoId, v.Request.Request, v.Request.CuentaId)

			if err != nil {
				assert.Equal(t, v.Error.Error(), err.Error())
			}

			if err == nil {
				assert.NotNil(t, response)

			}

		})
	}
}
