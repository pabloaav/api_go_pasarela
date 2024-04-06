package administracion_test

import (
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	fake "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/administracion/administracionfake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TablePagosExpirados struct {
	Nombre        string
	Erro          error
	MetodosMocked []callsMocked
}

func _inicializarPagosExpirados() (table []TablePagosExpirados) {

	filtroConf := filtros.ConfiguracionFiltro{
		Nombre: "TIEMPO_EXPIRACION_PAGOS",
	}

	filtroPending := filtros.PagoEstadoFiltro{
		Nombre: "Pending",
	}

	configuracion := entities.Configuracione{
		Model:       gorm.Model{ID: 1},
		Nombre:      "TIEMPO_EXPIRACION_PAGOS",
		Descripcion: "Tiempo en días para que expire un pago que está en estado pending ",
		Valor:       "30",
	}

	pagoEstadoPending := entities.Pagoestado{
		Model:  gorm.Model{ID: 1},
		Estado: "Pending",
	}

	pagoEstadoExpirado := entities.Pagoestado{
		Model:  gorm.Model{ID: 6},
		Estado: "Expired",
	}

	filtroPagos := filtros.PagoFiltro{
		PagoEstadosId:    uint64(pagoEstadoPending.ID),
		TiempoExpiracion: configuracion.Valor,
	}

	filtroExpired := filtros.PagoEstadoFiltro{
		Nombre: "Expired",
	}

	pagosExpirados := []entities.Pago{
		{PagoestadosID: 1},
		{PagoestadosID: 1},
		{PagoestadosID: 1},
		{PagoestadosID: 1},
	}

	mapa := make(map[string]interface{})
	mapa["pagos"] = pagosExpirados
	mapa["pagoEstado"] = uint64(pagoEstadoExpirado.ID)

	table = []TablePagosExpirados{
		{"Debe Retornar un error si no puede recuperar el parámetro TIEMPO_EXPIRACION_PAGOS", fmt.Errorf(administracion.ERROR_CONFIGURACIONES), []callsMocked{
			{"GetConfiguracion", filtroConf, entities.Configuracione{}, fmt.Errorf(administracion.ERROR_CONFIGURACIONES)},
		}},

		{"Debe Retornar un error si no puede recuperar el estado pending", fmt.Errorf(administracion.ERROR_PAGO_ESTADO), []callsMocked{
			{"GetConfiguracion", filtroConf, configuracion, nil},
			{"GetPagoEstado", filtroPending, entities.Pagoestado{}, fmt.Errorf(administracion.ERROR_PAGO_ESTADO)},
		}},
		{"Debe Retornar un error si no puede recuperar los pagos expirados", fmt.Errorf(administracion.ERROR_PAGO), []callsMocked{
			{"GetConfiguracion", filtroConf, configuracion, nil},
			{"GetPagoEstado", filtroPending, pagoEstadoPending, nil},
			{"GetPagos", filtroPagos, []entities.Pago{}, fmt.Errorf(administracion.ERROR_PAGO)},
		}},
		{"Si no existen pagos expirados finaliza operacion", fmt.Errorf(administracion.ERROR_PAGO), []callsMocked{
			{"GetConfiguracion", filtroConf, configuracion, nil},
			{"GetPagoEstado", filtroPending, pagoEstadoPending, nil},
			{"GetPagos", filtroPagos, []entities.Pago{}, nil},
		}},
		{"Debe retornar un error si no puede recuperar el estado Expired", fmt.Errorf(administracion.ERROR_PAGO_ESTADO), []callsMocked{
			{"GetConfiguracion", filtroConf, configuracion, nil},
			{"GetPagoEstado", filtroPending, pagoEstadoPending, nil},
			{"GetPagos", filtroPagos, pagosExpirados, nil},
			{"GetPagoEstado", filtroExpired, entities.Pagoestado{}, fmt.Errorf(administracion.ERROR_PAGO_ESTADO)},
		}},
		{"Debe retornar un error si no puede modificar los estados de los pagos", fmt.Errorf(administracion.ERROR_CREAR_ESTADO_LOGS), []callsMocked{
			{"GetConfiguracion", filtroConf, configuracion, nil},
			{"GetPagoEstado", filtroPending, pagoEstadoPending, nil},
			{"GetPagos", filtroPagos, pagosExpirados, nil},
			{"GetPagoEstado", filtroExpired, pagoEstadoExpirado, nil},
			{"UpdateEstadoPagos", mapa, nil, fmt.Errorf(administracion.ERROR_CREAR_ESTADO_LOGS)},
		}},
	}

	return
}

func TestModificarEstadoPagosExpirados(t *testing.T) {
	table := _inicializarPagosExpirados()
	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for _, m := range v.MetodosMocked {
				mockRepository.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
				if m.Nombre == "UpdateEstadoPagos" {
					mockRepository.On(m.Nombre, m.Request.(map[string]interface{})["pagos"], m.Request.(map[string]interface{})["pagoEstado"]).Return(m.Erro).Once()
				}
			}

			err := service.ModificarEstadoPagosExpirados()

			if err != nil {
				assert.Equal(t, v.Erro.Error(), err.Error())
			}

		})
	}
}

// TEST CONSULTAR PAGOS Y GENERAR MOVIMIENTOS EN DEV
func TestRequestValidConsultarMovimientos(t *testing.T) {
	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)
	TableDriverTest := fake.EstructuraConsultarPagosFake()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		want := TableDriverTest.WantTable
		_, got := service.GetPagosDevService(TableDriverTest.Request)
		assert.Equal(t, got.Error(), want)
	})
}
