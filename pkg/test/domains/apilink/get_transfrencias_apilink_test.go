package apilink_tests

import (
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"

	"github.com/stretchr/testify/assert"
)

func TestGetTransferenciasApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}

	requestValidoFake := linktransferencia.RequestGetTransferenciasLink{
		Cbu:        "0340218608218026437001",
		Tamanio:    linkdtos.CienTransf,
		Pagina:     1,
		FechaDesde: time.Now(),
		FechaHasta: time.Now(),
	}
	titularFake := linktransferencia.TitularTransferenciaLink{
		IdTributario: "555",
		RazonSocial:  "Alex",
	}
	origenFake := linktransferencia.OrigenResponseTransferenciaLink{
		Cbu:         "0340218608218026437001",
		RazonSocial: "Prueba",
	}
	destinoFake := linktransferencia.DestinoResponseTransferenciaLink{
		Cbu:     "0340218608218026437001",
		Alias:   "",
		Titular: titularFake,
	}
	paginadoFake := linktransferencia.PaginadoTransferenciaLink{
		Total:       5,
		CantPaginas: 10,
	}
	transferenciaFake := linktransferencia.Transferencia{
		Origen:         origenFake,
		Destino:        destinoFake,
		Importe:        "1",
		Moneda:         linkdtos.Pesos,
		Motivo:         linkdtos.AlquilerTransf,
		Referencia:     "1234",
		FechaOperacion: time.Now(),
	}
	lista := []linktransferencia.Transferencia{transferenciaFake}
	responseFake := linktransferencia.ResponseGetTransferenciasLink{
		Transferencias: lista,
		Paginado:       paginadoFake,
	}

	t.Run("Debe retornar un error si el cbu es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetTransferenciasCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetTransferenciasApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe retornar un error si tamanio de la pagina es invalido inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetTransferenciasTamanioInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetTransferenciasApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_PAGINADO_TAMANIO)
		}

	})

	t.Run("Caso positivo deber retornar una lista de transferencias bancaria", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetTransferenciasApiLink", requerimentoId, requestValidoFake, tokenFake.AccessToken).
			Return(&responseFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.GetTransferenciasApiLinkService(requerimentoId, requestValidoFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseFake)

	})

}
