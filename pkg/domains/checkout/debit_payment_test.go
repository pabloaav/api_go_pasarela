package checkout_test

import (
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

type tableDebitPaymentTest struct {
	Nombre              string
	Error               error
	Request             *dtos.ResultadoRequest
	Pago                *entities.Pago
	Cuenta              *entities.Cuenta
	MethodMocked        []callsMocked
	InstallmentsDetails *dtos.InstallmentDetailsResponse
}

func iniciarDebitPayment() []tableDebitPaymentTest {
	reqDebitoTokenNoValido := requestDebitoValido
	reqDebitoTokenNoValido.HolderDocNum = ""
	tokenReqNoValido := tokenReq
	tokenReqNoValido.Card.CardHolderIdentification.NumberDni = ""

	return []tableDebitPaymentTest{
		{
			"Debe devolver error de validación",
			fmt.Errorf("debe indicar el método por el cual va a pagar"),
			&requestNoValido,
			&pagoCreated,
			&cuentaValida,
			[]callsMocked{},
			nil,
		},
		{
			"Debe devolver error al solicitar token con error",
			fmt.Errorf("Error al solicitar token"),
			&reqDebitoTokenNoValido,
			&pagoCreated,
			&cuentaValida,
			[]callsMocked{
				{"SERV", "SolicitarToken", []interface{}{tokenReqNoValido}, nil, fmt.Errorf("Error al solicitar token")},
			},
			nil,
		},
	}
}

func TestDebitoPaymentErrors(t *testing.T) {
	assertions := assert.New(t)

	serv := mockservice.MockPrismaService{}

	payment := checkout.NewDebitPayment(&serv)

	table := iniciarDebitPayment()

	for _, v := range table {
		t.Run(v.Nombre, func(t *testing.T) {
			for _, c := range v.MethodMocked {
				switch c.mock {
				case "SERV":
					serv.On(c.nombre, c.filtro...).Return(c.retorno, c.mensaje).Once()
				}
			}

			res, err := payment.CreateResultado(v.Request, v.Pago, v.Cuenta, "aaaa-bbb-cccc-ddddd", v.InstallmentsDetails)

			if err != nil {
				assertions.EqualError(err, v.Error.Error())
			} else {
				assertions.NotZero(res.PaidAt)
			}

		})
	}
}

func TestDebitPaymentValid(t *testing.T) {
	assertions := assert.New(t)
	serv := mockservice.MockPrismaService{}
	payment := checkout.NewDebitPayment(&serv)

	t.Run("Cuando el llamado a la api devuelve error debo devolver objeto con el error en el status y paidAt en fecha 0", func(t *testing.T) {
		cuentaNoValida := cuentaValida
		cuentaNoValida.Cuenta = ""
		payReqNoValido := payReq
		payReqNoValido.PagoSimple.EstablishmentName = ""
		payReqNoValido.PagoSimple.PaymentMethodID = 31
		var installmentsDetails *dtos.InstallmentDetailsResponse
		serv.On("SolicitarToken", tokenReq).Return(pagoToken, nil).Once()
		serv.On("Payments", payReqNoValido).Return(nil, fmt.Errorf("Error al procesar pago")).Once()

		pay, _ := payment.CreateResultado(&requestDebitoValido, &pagoCreated, &cuentaNoValida, "aaaa-bbb-cccc-ddddd", installmentsDetails)

		assertions.Equal(pay.StateComment, "Error al procesar pago")
		assertions.Zero(pay.PaidAt)
	})

	t.Run("Cuando el llamado a la api devuelve status ok devuelvo el status, el tiket, el authorization code y paidAt distinto de 0", func(t *testing.T) {
		serv.On("SolicitarToken", tokenReq).Return(pagoToken, nil).Once()
		payReqDebitValid := payReq
		payReqDebitValid.PagoSimple.PaymentMethodID = 31
		var installmentsDetails *dtos.InstallmentDetailsResponse
		serv.On("Payments", payReqDebitValid).Return(payResponse, nil).Once()

		pay, _ := payment.CreateResultado(&requestDebitoValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", installmentsDetails)

		assertions.Equal(pay.StateComment, "approved")
		assertions.NotZero(pay.TicketNumber)
		assertions.NotZero(pay.AuthorizationCode)
		assertions.NotZero(pay.PaidAt)
	})
}
