package checkout_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentMethod(t *testing.T) {
	assertions := assert.New(t)
	s := new(mockservice.MockPrismaService)
	u := new(mockservice.MockUtilService)
	d := new(mockservice.MockApiLinkService)

	payment := checkout.NewPaymentFactory()

	t.Run("Debe devolver un error si se pide un payment method que no existe", func(t *testing.T) {
		_, err := payment.GetPaymentMethod(0)
		assertions.EqualError(err, "no se reconoce el metodo de pago número 0")
	})

	t.Run("Debe devolver un objeto de tipo Credit Payment cuando venga 1 de parametro", func(t *testing.T) {
		x, _ := payment.GetPaymentMethod(1)
		assertions.IsType(x, checkout.NewCreditPayment(s, u))
	})

	t.Run("Debe devolver un objeto de tipo Debit Payment cuando venga 1 de parametro", func(t *testing.T) {
		x, _ := payment.GetPaymentMethod(2)
		assertions.IsType(x, checkout.NewDebitPayment(s))
	})

	t.Run("Debe devolver un objeto de tipo Offline Payment cuando venga 1 de parametro", func(t *testing.T) {
		x, _ := payment.GetPaymentMethod(3)
		assertions.IsType(x, checkout.NewOfflinePayment(s))
	})

	t.Run("Debe devolver un objeto de tipo Debin Payment cuando venga 1 de parametro", func(t *testing.T) {
		x, _ := payment.GetPaymentMethod(4)
		assertions.IsType(x, checkout.NewDebinPayment(d, u))
	})
}
