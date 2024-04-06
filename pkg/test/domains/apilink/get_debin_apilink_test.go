package apilink_tests

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"

	"github.com/stretchr/testify/assert"
)

func TestGetDebinApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requestDebinFake := linkdebin.RequestGetDebinLink{}
	requestDebinFake.Cbu = "0340218608218026437001"
	requestDebinFake.Id = "123563"
	debinEstadoIniciado := linkdebin.DebinDetalleLink{
		Estado: linkdtos.Iniciado,
	}
	responseDebinIniciadoFake := linkdebin.ResponseGetDebinLink{
		Debin: debinEstadoIniciado,
	}
	t.Run("Debe Retornar un error si el cbu es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe Retornar un error si el id es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinIdInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ID)
		}

	})

	t.Run("Caso Positivo debe retornar un debin", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetDebinApiLink", requerimentoId, requestDebinFake, tokenFake.AccessToken).Return(&responseDebinIniciadoFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.GetDebinApiLinkService(requerimentoId, requestDebinFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseDebinIniciadoFake)

	})

}
