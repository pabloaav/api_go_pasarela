package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/stretchr/testify/mock"
)

type MockDebinPaymentMethod struct {
	mock.Mock
}

func (mk *MockDebinPaymentMethod) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta) (*entities.Pagointento, error) {
	args := mk.Called(request, pago, cuenta)
	result := args.Get(0)
	return result.(*entities.Pagointento), args.Error(1)
}
