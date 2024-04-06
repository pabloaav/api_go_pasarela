package apilink_tests

import (
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

func TestCreateDebinCompradorApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"

	t.Run("Debe Retornar un error si el cbu del comprador es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateCompradorCuentaCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU_COMPRADOR)
		}

	})

	t.Run("Debe Retornar un error si el cbu esta vacio y el alias del cbu del comprador es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateCompradorCuentaAliasCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ALIASCBULEN)
		}

	})

	t.Run("Debe Retornar un error si el cuit del comprador es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateCompradorCuitInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), commons.ERROR_CUIL)
		}

	})

}

func TestCreateDebinVendedorApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"

	t.Run("Debe Retornar un error si el cbu del vendedor es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateVendedorCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU_VENDEDOR)
		}

	})

	t.Run("Debe Retornar un error si el cbu esta vacio y el alias del cbu del vendedor es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateVendedorAliasCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ALIASCBULEN)
		}

	})

}

func TestCreateDebinApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"

	t.Run("Debe Retornar un error si el comprobanteId es vacio", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateComprobanteIdInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_IDENTIFICADORDEBIN)
		}

	})

	t.Run("Debe Retornar un error si el concepto es invalido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateConceptoInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_CONCEPTO)
		}

	})

	t.Run("Debe Retornar un error si el tiempo de expiracion es invalido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateTiempoExpiracionInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_TIEMPOEXPIRACIONDEBIN)
		}

	})

	t.Run("Debe Retornar un error si el importe es negatico", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateImporteInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_IMPORTE)
		}

	})

	t.Run("Debe Retornar un error si la moneda es incorrecta", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateMonedaInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_MONEDA)
		}

	})

	t.Run("Debe Retornar un error la descripcion es invalida", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateDescripcionInvalida()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_DESCRIPCIONDEBIN)
		}

	})

	t.Run("Debe Retornar un error la descripcionPrestacion es invalida", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestCreateDescripcionPrestacionInvalida()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.CreateDebinApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_DESCRIPCIONDEBIN)
		}

	})

}

func TestCreateDebinApilinkPositivo(t *testing.T) {
	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	requestValidoFake := linkdebin.RequestDebinCreateLink{
		Comprador: _compradorValidoCreateDebin(),
		Vendedor:  _vendedorValidoCreateDebin(),
		Debin:     _debinValidoCreateDebin(),
	}

	responseFake := linkdebin.ResponseDebinCreateLink{
		Id:              "6722716440723456",
		FechaOperacion:  time.Now(),
		Estado:          "ERROR_DEBITO",
		FechaExpiracion: time.Now(),
	}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"

	t.Run("Debe Retornar un error si el comprobanteId es vacio", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)
		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("CreateDebinApiLink", requerimentoId, requestValidoFake, tokenFake.AccessToken).Return(&responseFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.CreateDebinApiLinkService(requerimentoId, requestValidoFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response.Estado, linkdtos.EnumEstadoDebin("ERROR_DEBITO"))

	})

}
