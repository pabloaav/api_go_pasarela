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

func TestGetDebinesApilink(t *testing.T) {

	requerimentoId := "02329451-dd5e-4eef-a44f-56bdf9008357"
	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	tokenFake := linkdtos.TokenLink{
		AccessToken: "2222333333336366",
	}
	requestValidoFake := linkdebin.RequestGetDebinesLink{
		Pagina:      5,
		Tamanio:     linkdtos.Cinco,
		Cbu:         "0340218608218026437001",
		Estado:      linkdtos.Iniciado,
		FechaDesde:  time.Now(),
		FechaHasta:  time.Now(),
		EsComprador: true,
		Tipo:        linkdtos.DebinDefault,
	}

	paginadoFake := linkdtos.PaginadoResponseLink{
		Pagina:          1,
		CantidadPaginas: 5,
	}
	vendedorFake := linkdebin.VendedorDebinesListLink{
		Cuit:    "20953043336",
		Titular: "Alex",
	}
	compradorFake := linkdebin.CompradorDebinesListLink{
		Cuit:    "20953043336",
		Titular: "Alex",
	}
	debinFake := linkdebin.DebinesListLink{
		Id:              "8215810612920320",
		Concepto:        "PLF",
		Moneda:          linkdtos.Dolar,
		Importe:         42152846,
		Estado:          linkdtos.Acreditado,
		Tipo:            "DEBINPLF",
		FechaExpiracion: time.Now(),
		Devuelto:        false,
		ContraCargoId:   "3318834243043328",
		Comprador:       compradorFake,
		Vendedor:        vendedorFake,
	}
	listaDEbines := []linkdebin.DebinesListLink{debinFake}
	responseFake := linkdebin.ResponseGetDebinesLink{
		Paginado: paginadoFake,
		Debines:  listaDEbines,
	}

	t.Run("Debe retornar un error si el cbu es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinesCbuInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinesApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_CBU)
		}

	})

	t.Run("Debe retornar un error si el tamaño del paginado es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinesTamanioInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinesApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_PAGINADO_TAMANIO)
		}

	})

	t.Run("Debe retornar un error si el estado del Debin es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinesEstadoInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinesApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_ESTADO_DEBIN)
		}

	})

	t.Run("Debe retornar un error si el tipo del debin es inválido", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		lista := requestGetDebinesTipoInvalido()

		for _, r := range lista {

			service := apilink.NewService(mockRemoteRepository, mockrepository)

			_, err := service.GetDebinesApiLinkService(requerimentoId, *r)

			mockRemoteRepository.AssertExpectations(t)

			assert.Equal(t, err.Error(), tools.ERROR_ENUM_TIPO_DEBIN)
		}

	})

	t.Run("Caso Positivo debe retornar una lista de debines", func(t *testing.T) {

		mockRemoteRepository := new(mockrepository.MockRemoteRepositoryApiLink)
		mockrepository := new(mockrepository.MockRepositoryApiLink)

		mockRemoteRepository.On("GetTokenApiLink", requerimentoId, scopes).Return(&tokenFake, nil)
		mockRemoteRepository.On("GetDebinesApiLink", requerimentoId, requestValidoFake, tokenFake.AccessToken).Return(&responseFake, nil)

		service := apilink.NewService(mockRemoteRepository, mockrepository)

		response, err := service.GetDebinesApiLinkService(requerimentoId, requestValidoFake)

		mockRemoteRepository.AssertExpectations(t)

		assert.Equal(t, err, nil)
		assert.Equal(t, response, responseFake)

	})

}
