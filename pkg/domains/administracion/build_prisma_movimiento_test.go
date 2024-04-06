package administracion_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/administracion/administracionfake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

func TestBuildPrismaMovimientoFakeListaCierreLote(t *testing.T) {
	listaCierreLotePrismaFakeListaCierrreLote := administracionfake.DataCierreloteFake
	mockRepositoryAdministracion := new(mockrepository.MockRepositoryAdministracion)
	mockApilinkService := new(mockservice.MockApiLinkService)
	mockUtilService := new(mockservice.MockUtilService)
	mockCommonds := new(mockservice.MockCommonsService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepositoryAdministracion, mockApilinkService, mockCommonds, mockUtilService, mockRepositoryWebHook, mockStore)

	t.Run(listaCierreLotePrismaFakeListaCierrreLote.TituloPrueba, func(t *testing.T) {
		want := fmt.Sprintf("%v", listaCierreLotePrismaFakeListaCierrreLote.WantTable)
		//listaCierreLoteDB := listaCierreLotePrismaFakeListaCierrreLote.DataPruebaListaCierreLote.([]entities.Prismacierrelote) //make([]entities.Prismacierrelote{,} )
		mockRepositoryAdministracion.On("GetPrismaCierreLotes").Return(listaCierreLotePrismaFakeListaCierrreLote.DataPruebaListaCierreLote, errors.New(want))
		_, got := service.BuildPrismaMovimiento(true)
		assert.Equal(t, got.Error(), want)
	})
}

func TestBuildPrismaMovimientoFakeGetPagos(t *testing.T) {
	listaCierreLotePrismaFakeListaPago := administracionfake.DataPagosFake
	mockRepositoryAdministracion := new(mockrepository.MockRepositoryAdministracion)
	mockApilinkService := new(mockservice.MockApiLinkService)
	mockCommonds := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepositoryAdministracion, mockApilinkService, mockCommonds, mockUtilService, mockRepositoryWebHook, mockStore)
	t.Run(listaCierreLotePrismaFakeListaPago.TituloPrueba, func(t *testing.T) {
		want := fmt.Sprintf("%v", listaCierreLotePrismaFakeListaPago.WantTable)
		mockRepositoryAdministracion.On("GetPrismaCierreLotes").Return(listaCierreLotePrismaFakeListaPago.DataPruebaListaCierreLote, nil)
		mockRepositoryAdministracion.On("GetPagos", listaCierreLotePrismaFakeListaPago.DataFiltroPago).Return(listaCierreLotePrismaFakeListaPago.DataPruebaListaPagos, errors.New(want))
		_, got := service.BuildPrismaMovimiento(true)
		assert.Equal(t, got.Error(), want)
	})
}

func TestBuildPrismaMovimientoFakeGetEstadoExternos(t *testing.T) {
	listaCierreLotePrismaFakeGetPagoEstadoExterno := administracionfake.GetpagoEstadoExternoFake
	mockRepositoryAdministracion := new(mockrepository.MockRepositoryAdministracion)
	mockApilinkService := new(mockservice.MockApiLinkService)
	mockCommonds := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepositoryAdministracion, mockApilinkService, mockCommonds, mockUtilService, mockRepositoryWebHook, mockStore)
	mockRepositoryAdministracion.On("GetPrismaCierreLotes").Return(listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaCierreLote, nil)
	mockRepositoryAdministracion.On("GetPagos", listaCierreLotePrismaFakeGetPagoEstadoExterno.DataFiltroPago).Return(listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaPagos, nil)
	for k, v := range listaCierreLotePrismaFakeGetPagoEstadoExterno.DataFiltroGetEstadoExterno {
		listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaCierreLote[k].Tipooperacion = entities.EnumTipoOperacion(v)
	}
	fmt.Println(listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaCierreLote)
	for _, valueFilterEstado := range listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaCierreLote {
		t.Run(listaCierreLotePrismaFakeGetPagoEstadoExterno.TituloPrueba, func(t *testing.T) {
			want := fmt.Sprintf("%v", listaCierreLotePrismaFakeGetPagoEstadoExterno.WantTable)

			filtroEstadoExterno := filtros.PagoEstadoExternoFiltro{Nombre: string(valueFilterEstado.Tipooperacion)}
			mockRepositoryAdministracion.On("GetPagosEstadosExternos", filtroEstadoExterno).Return(listaCierreLotePrismaFakeGetPagoEstadoExterno.DataPruebaListaPagoEstadoExterno, errors.New(want))
			_, got := service.BuildPrismaMovimiento(true)
			assert.Equal(t, got.Error(), want)
		})
	}
}

func TestBuildPrismaMovimientoFakevalido(t *testing.T) {
	var unPagoEstadoExterno entities.Pagoestadoexterno
	estructurasValidas := administracionfake.GetpagoEstadoExternovalido
	mockRepositoryAdministracion := new(mockrepository.MockRepositoryAdministracion)
	mockApilinkService := new(mockservice.MockApiLinkService)
	mockCommonds := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepositoryAdministracion, mockApilinkService, mockCommonds, mockUtilService, mockRepositoryWebHook, mockStore)
	mockRepositoryAdministracion.On("GetPrismaCierreLotes").Return(estructurasValidas.DataPruebaListaCierreLote, nil)
	mockRepositoryAdministracion.On("GetPagos", estructurasValidas.DataFiltroPago).Return(estructurasValidas.DataPruebaListaPagos, nil)
	wantValido := estructurasValidas.WantTable.(administracionfake.EstructuraValida)
	var wantMovimiento []entities.Movimiento = make([]entities.Movimiento, 4)
	for _, valueFilterEstado := range estructurasValidas.DataPruebaListaCierreLote {
		filtroEstadoExterno := filtros.PagoEstadoExternoFiltro{Nombre: string(valueFilterEstado.Tipooperacion)}
		for _, v := range estructurasValidas.DataPruebaListaPagoEstadoExterno {
			if string(valueFilterEstado.Tipooperacion) == v.Estado {
				unPagoEstadoExterno = v
				break
			}
		}
		mockRepositoryAdministracion.On("GetPagosEstadosExternos", filtroEstadoExterno).Return([]entities.Pagoestadoexterno{unPagoEstadoExterno}, nil)

		for k, valorPAgo := range estructurasValidas.DataPruebaListaPagos {
			if valueFilterEstado.PagosUuid == valorPAgo.Uuid {

				if valueFilterEstado.Tipooperacion == "C" {
					wantMovimiento[k].AddDebito(uint64(valorPAgo.PagosTipo.CuentasID), uint64(valorPAgo.PagoIntentos[0].ID), valorPAgo.PagoIntentos[0].Amount)

					valorPAgo.PagoIntentos[0].AvailableAt = valueFilterEstado.FechaCierre
					valorPAgo.PagoIntentos[0].RevertedAt = time.Time{}

				} else {
					wantMovimiento[k].AddCredito(uint64(valorPAgo.PagosTipo.CuentasID), uint64(valorPAgo.PagoIntentos[0].ID), -1.00*(valorPAgo.PagoIntentos[0].Amount))
					valorPAgo.PagoIntentos[0].RevertedAt = valueFilterEstado.FechaCierre
					valorPAgo.PagoIntentos[0].AvailableAt = time.Time{}
				}

			}

		}
	}
	/*
		CL: Cierre Lote - P: Pagos - PEL: Pago Estado Logs - M: Movimientos - PI: Pago Intento
	*/
	t.Run(estructurasValidas.TituloPrueba, func(t *testing.T) {
		movimientoCierreLote, _ := service.BuildPrismaMovimiento(true)
		mockRepositoryAdministracion.AssertExpectations(t)
		assert.Equal(t, len(movimientoCierreLote.ListaCLPrisma), len(wantValido.ListPrismaCierreLote))
		assert.Equal(t, len(movimientoCierreLote.ListaPagos), len(wantValido.ListPagos))
		assert.Equal(t, len(movimientoCierreLote.ListaPagosEstadoLogs), len(wantValido.ListPagoEstadoLogs))
		assert.Equal(t, len(movimientoCierreLote.ListaMovimientos), len(wantValido.Listmovimientos))
		assert.Equal(t, len(movimientoCierreLote.ListaPagoIntentos), len(wantValido.ListPagoIntentos))

		assert.Equal(t, movimientoCierreLote.ListaCLPrisma, wantValido.ListPrismaCierreLote)
		assert.Equal(t, movimientoCierreLote.ListaPagos, wantValido.ListPagos)
		assert.Equal(t, movimientoCierreLote.ListaPagosEstadoLogs, wantValido.ListPagoEstadoLogs)
		assert.Equal(t, movimientoCierreLote.ListaMovimientos, wantMovimiento)
		assert.Equal(t, movimientoCierreLote.ListaPagoIntentos, wantValido.ListPagoIntentos)

		var gotMontoCL, gotMontoM, gotMontoPI entities.Monto
		var wantMontoCL, wantMontoM, wantMontoPI entities.Monto

		for i := 0; i < len(estructurasValidas.DataPruebaListaCierreLote); i++ {
			gotMontoCL += movimientoCierreLote.ListaCLPrisma[i].Monto
			wantMontoCL += wantValido.ListPrismaCierreLote[i].Monto
			gotMontoM += movimientoCierreLote.ListaMovimientos[i].Monto
			wantMontoM += wantValido.Listmovimientos[i].Monto
			gotMontoPI += movimientoCierreLote.ListaPagoIntentos[i].Amount
			wantMontoPI += wantValido.ListPagoIntentos[i].Amount
		}
		assert.Equal(t, gotMontoCL, wantMontoCL)
		assert.Equal(t, gotMontoPI, wantMontoPI)
		assert.Equal(t, gotMontoM, wantMontoM)

	})

}
