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

func TestGetTransferenciaApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}

	requestValidoFake := linktransferencia.RequestGetTransferenciaLink{
		NumeroReferenciaBancaria: "111111111222222233333333333333",
		Cbu:                      "0340218608218026437001",
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

	responseFake := linktransferencia.ResponseGetTransferenciaLink{
		Origen:         origenFake,
		Destino:        destinoFake,
		Importe:        "1",
		Moneda:         linkdtos.Pesos,
		Motivo:         linkdtos.AlquilerTransf,
		Referencia:     "1234",
		FechaOperacion: time.Now(),
	}
	t.Run("Debe Retornar un error si el cbu es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetTransferenciaCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe Retornar un error si el numero de referencia Bancaria es inválida", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetTransferenciaNumeroReferenciaInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_REFERENCIA_BANCARIA)
		}

	})

	t.Run("Caso positivo deber retornar una transferencia bancaria", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetTransferenciaApiLink", requerimentoId, requestValidoFake, tokenFake.AccessToken).
			Return(&responseFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.GetTransferenciaApiLinkService(requerimentoId, requestValidoFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseFake)

	})
}
