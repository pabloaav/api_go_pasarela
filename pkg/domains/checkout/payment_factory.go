package checkout

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/pagooffline"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

// Definición de constantes, cada tipo de pago tiene su id en la tabla Channels
const (
	Credit   = 1
	Debit    = 2
	Offline  = 3
	Debin    = 4
	Rapipago = 5
)

// PaymentFactory va a devolver el Método de pago según el parámetro que reciba
type PaymentFactory interface {
	GetPaymentMethod(m int) (PaymentMethod, error)
}

type paymentFactory struct{}

func NewPaymentFactory() PaymentFactory {
	return &paymentFactory{}
}

func (p *paymentFactory) GetPaymentMethod(m int) (PaymentMethod, error) {
	switch m {
	case Credit:
		return NewCreditPayment(prisma.Resolve(), util.Resolve()), nil
	case Debit:
		return NewDebitPayment(prisma.Resolve()), nil
	case Offline:
		return NewRapipagoPayment(pagooffline.Resolve()), nil
	case Debin:
		return NewDebinPayment(apilink.Resolve(), util.Resolve()), nil
	// case Rapipago:
	// 	return NewOfflinePayment(prisma.Resolve()), nil
	default:
		return nil, fmt.Errorf("no se reconoce el metodo de pago número %d", m)
	}
}

// PaymentMethod es la interfaz que cada metodo de pago va a tener que implementar
type PaymentMethod interface {
	CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transaction string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error)
}
