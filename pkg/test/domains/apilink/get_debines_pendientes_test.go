package apilink_tests

import (
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"

	"github.com/stretchr/testify/assert"
)

func TestGetDebinesPendientesApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requestCbuValido := "0340218608218026437001"

	vendedorFake := linkdebin.VendedorDebinesPendientesLink{
		Cuit:    "20953043637",
		Titular: "Alex",
	}
	debinPendienteFake := linkdebin.DebinesPendientesLink{
		Id:              "111111",
		Concepto:        linkdtos.DevolucionPEI,
		Moneda:          linkdtos.Dolar,
		Importe:         20,
		Estado:          linkdtos.Iniciado,
		FechaExpiracion: time.Now(),
		CbuComprador:    "0340218608218026437001",
		Vendedor:        vendedorFake,
	}
	lista := []linkdebin.DebinesPendientesLink{
		debinPendienteFake,
	}
	responseDebinPendienteFake := linkdebin.ResponseGetDebinesPendientesLink{
		Debines: lista,
	}

	t.Run("Debe Retornar un error si el cbu es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinesPendientesCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinesPendientesApiLinkService(requerimentoId, r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Caso Positivo debe retornar una lista de debines pendientes", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetDebinesPendientesApiLink", requerimentoId, requestCbuValido, tokenFake.AccessToken).Return(&responseDebinPendienteFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.GetDebinesPendientesApiLinkService(requerimentoId, requestCbuValido)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseDebinPendienteFake)

	})

}
