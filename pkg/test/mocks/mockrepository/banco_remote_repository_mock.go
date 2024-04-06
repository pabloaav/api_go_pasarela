package mockrepository

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
	"github.com/stretchr/testify/mock"
)

type MockBancoRemoteRepository struct {
	mock.Mock
}

func (mock *MockBancoRemoteRepository) GetTokenBancoRemotRepository(token bancodtos.RequestTokenBanco) (responseToken bancodtos.ResponseTokenBanco, err error) {
	args := mock.Called(token) // Called le dice al objeto simulado que se ha llamado a un método y obtiene una serie de argumentos para devolver.
	resultado := args.Get(0)
	return resultado.(bancodtos.ResponseTokenBanco), args.Error(1)
}

func (mock *MockBancoRemoteRepository) GetMovimientosBancoRemotRepository(filtrosMovimientos filtros.MovimientosBancoFiltro, token string) (response []bancodtos.ResponseMovimientosBanco, erro error) {
	args := mock.Called(filtrosMovimientos, token)
	resultado := args.Get(0)
	return resultado.([]bancodtos.ResponseMovimientosBanco), args.Error(1)
}

func (mock *MockBancoRemoteRepository) ActualizarRegistrosMatchRepository(listaMovimientoMatch bancodtos.RequestUpdateMovimiento, token string) (response bool, erro error) {
	args := mock.Called(listaMovimientoMatch, token)
	return args.Bool(1), args.Error(1)
}

func (mock *MockBancoRemoteRepository) GetConsultarMovimientosRemoteRepository(filtro filtros.RequestMovimientos, token string) (response bancodtos.ResponseMovimientos, erro error) {
	args := mock.Called(filtro, token)
	resultado := args.Get(1)
	return resultado.(bancodtos.ResponseMovimientos), args.Error(1)

}
