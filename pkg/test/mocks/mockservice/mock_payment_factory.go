package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/stretchr/testify/mock"
)

type MockPaymentFactory struct {
	mock.Mock
	MockPrisma  *MockPrismaService
	MockApiLink *MockApiLinkService
}

const (
	Credit  = 1
	Debit   = 2
	Offline = 3
	Debin   = 4
)

func (mk *MockPaymentFactory) GetPaymentMethod(m int) (checkout.PaymentMethod, error) {
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
	return result.(checkout.PaymentMethod), args.Error(1)
}
