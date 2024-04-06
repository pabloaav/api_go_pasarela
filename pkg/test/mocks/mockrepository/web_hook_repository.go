package mockrepository

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryWebHook struct {
	mock.Mock
}

func (mock *MockRepositoryWebHook) NotificarPago(Pago dtos.ResultadoResponseWebHook) (erro error) {
	args := mock.Called(Pago)
	return args.Error(0)
}

func (mock *MockRepositoryWebHook) NotificarPagos(listaPagos webhook.WebhookResponse) (erro error) {
	args := mock.Called(listaPagos)
	return args.Error(0)
}
