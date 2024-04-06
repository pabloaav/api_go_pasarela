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

func TestCreateTransferenciaApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requestTransferenciaFake := _CreateTransferenciaValida()
	responseTransferenciaFake := linktransferencia.ResponseTransferenciaCreateLink{
		NumeroReferenciaBancaria: "123366555",
		FechaOperacion:           time.Now(),
	}

	t.Run("Debe Retornar un error si el cbu de origen es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaOrigenCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe Retornar un error si el cbu de destino es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaDestinoCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe Retornar un error si el alias del cbu de destino es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaDestinoAliasInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ALIASCBULEN)
		}

	})

	t.Run("Debe retornar un error si el tipo de moneda es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaMonedaInvalida()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_MONEDA)
		}

	})

	t.Run("Debe retornar un error si el motivo es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaMotivoInvalida()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_MOTIVO)
		}

	})

	t.Run("Debe retornar un error la referencia es inválida", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTransferenciaReferenciaInvalida()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateTransferenciaApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_REFERENCIA)
		}

	})

	t.Run("Caso positivo deber retornar el numero de referencia bancaria y fecha operación", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("CreateTransferenciaApiLink", requerimentoId, requestTransferenciaFake, tokenFake.AccessToken).
			Return(&responseTransferenciaFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.CreateTransferenciaApiLinkService(requerimentoId, requestTransferenciaFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseTransferenciaFake)

	})

}
