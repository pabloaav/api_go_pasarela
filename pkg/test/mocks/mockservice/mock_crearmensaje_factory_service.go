package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/stretchr/testify/mock"
)

type MockCrearMensajeServiceFactory struct {
	mock.Mock
}

const (
	EMAIL_TEMPLATE = 1
	EMAIL_ADJUNTO  = 2
)

func (mk *MockCrearMensajeServiceFactory) GetCrearMensajeMethod(m int) (util.CrearMensajeMethod, error) {
	// switch m {
	// case Credit:
	// 	return checkout.NewCreditPayment(mk.MockPrisma), nil
	// case Debit:
	// 	return checkout.NewDebitPayment(mk.MockPrisma), nil
	// case Offline:
	// 	return checkout.NewOfflinePayment(mk.MockPrisma), nil
	// case Debin:
	// 	return checkout.NewDebinPayment(mk.MockApiLink), nil
	// default:
	// 	return nil, fmt.Errorf("no se reconoce el metodo de pago número %d", m)
	// }
	args := mk.Called(m)
	result := args.Get(0)
	return result.(util.CrearMensajeMethod), args.Error(1)
}


