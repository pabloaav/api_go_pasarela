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

func TestDeleteDebinApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requestDeleteDebinValido := linkdebin.RequestDeleteDebinLink{
		Cbu: "0340218608218026437001",
		Id:  "25336885223",
	}

	requestDebinFake := linkdebin.RequestGetDebinLink{}
	requestDebinFake.Cbu = requestDeleteDebinValido.Cbu
	requestDebinFake.Id = requestDeleteDebinValido.Id
	debinEstadoAcreditado := linkdebin.DebinDetalleLink{
		Estado: linkdtos.Acreditado,
	}
	debinEstadoIniciado := linkdebin.DebinDetalleLink{
		Estado: linkdtos.Iniciado,
	}
	responseDebinAcreditadoFake := linkdebin.ResponseGetDebinLink{
		Debin: debinEstadoAcreditado,
	}
	responseDebinIniciadoFake := linkdebin.ResponseGetDebinLink{
		Debin: debinEstadoIniciado,
	}

	t.Run("Debe Retornar un error si el cbu informado es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		listaInvalidos := requestDeleteDebinCbuInvalido()

		for _, r := range listaInvalidos {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.DeleteDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe Retornar un error si el id informado es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		listaInvalidos := requestDeleteDebinIdInvalido()

		for _, r := range listaInvalidos {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.DeleteDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ID)
		}

	})

	t.Run("Debe Retornar un error si se trata de eliminar un debin acreditado", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetDebinApiLink", requerimentoId, requestDebinFake, tokenFake.AccessToken).Return(&responseDebinAcreditadoFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		_, err := service.DeleteDebinApiLinkService(requerimentoId, requestDeleteDebinValido)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err.Error(), apilink.ERROR_DELETE_DEBINES_ACREDITADOS)

	})

	t.Run("Caso positivo debin eliminado correctamente", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetDebinApiLink", requerimentoId, requestDebinFake, tokenFake.AccessToken).Return(&responseDebinIniciadoFake, nil)
		mockRemoteRepository.On("DeleteDebinApiLink", requerimentoId, requestDeleteDebinValido, tokenFake.AccessToken).Return(false, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.DeleteDebinApiLinkService(requerimentoId, requestDeleteDebinValido)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, false)

	})

}
