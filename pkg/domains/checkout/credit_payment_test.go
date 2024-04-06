package checkout_test

import (
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

var (
	tokenReq = prismadtos.StructToken{
		Card: prismadtos.Card{
			CardNumber:          requestCreditoValido.CardNumber,
			CardExpirationMonth: requestCreditoValido.CardExpiration[0:2],
			CardExpirationYear:  requestCreditoValido.CardExpiration[2:4],
			SecurityCode:        requestCreditoValido.CardCode,
			CardHolderName:      requestCreditoValido.HolderName,
			CardHolderIdentification: prismadtos.CardHolderIdentification{
				TypeDni:   prismadtos.EnumTipoDocumento(requestCreditoValido.HolderDocType),
				NumberDni: requestCreditoValido.HolderDocNum,
			},
		},
		TypePay: "simple",
	}
	pagoToken = prismadtos.PagoToken{
		Id:                "1",
		ValidationResulto: true,
		Status:            "accepted",
		DataUsed:          "",
		CardNumberLength:  16,
		Bin:               "034",
		DateCreate:        "010121",
		LastFourDigits:    "1234",
		SecurityCodeLeng:  3,
		ExpirationMonth:   8,
		ExpirationYear:    2021,
		DateLastUpdated:   "010121",
		DateDue:           "010121",
		CardHolder: prismadtos.CardHolder{
			Identification: prismadtos.Identification{
				TypeDni:   "DNI",
				NumberDni: "32880325",
			},
			Name: "Fernando Castro",
		},
	}
	payReq = prismadtos.StructPayments{
		PagoSimple: prismadtos.PaymentsSimpleRequest{
			Customerid: prismadtos.Customerid{
				ID: fmt.Sprint(cuentaValida.ID),
			},
			SiteTransactionID: "aaaa-bbb-cccc-ddddd",
			Token:             pagoToken.Id,
			PaymentMethodID:   requestCreditoValido.PaymentMethodID,
			Bin:               pagoToken.Bin,
			Amount:            requestCreditoValido.Importe,
			Currency:          "ARS",
			Installments:      1,
			Description:       pagoCreated.Description,
			PaymentType:       "single",
			EstablishmentName: cuentaValida.Cuenta,
			Customeremail: prismadtos.Customeremail{
				Email: requestCreditoValido.HolderEmail,
			},
			SubPayments: make([]interface{}, 0),
		},
		TypePay: "simple",
	}
	payResponse = prismadtos.PaymentsSimpleResponse{
		ID:                1,
		SiteTransactionID: "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		Token:             "1",
		Customer: prismadtos.Customer{
			ID: "1",
		},
		PaymentMethodID: 1,
		Bin:             "034",
		Amount:          105050,
		Currency:        "ARS",
		Installments:    1,
		PaymentType:     "single",
		SubPayments:     make([]interface{}, 0),
		Status:          "approved",
		StatusDetails: prismadtos.StatusDetails{
			Ticket:                "3975",
			CardAuthorizationCode: "083146",
		},
		CardBrand:                      "Visa",
		Date:                           "2021-07-07T15:04Z",
		FirstInstallmentExpirationDate: "070821",
		SiteID:                         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		EstablishmentName:              "Principal",
	}
)

type tableCreditPaymentTest struct {
	Nombre              string
	Error               error
	Request             *dtos.ResultadoRequest
	Pago                *entities.Pago
	Cuenta              *entities.Cuenta
	MethodMocked        []callsMocked
	InstallmentsDetails *dtos.InstallmentDetailsResponse
}

func iniciarCreditPayment() []tableCreditPaymentTest {

	reqCreditoTokenNoValido := requestCreditoValido
	reqCreditoTokenNoValido.HolderDocType = ""
	tokenReqNoValido := tokenReq
	tokenReqNoValido.Card.CardHolderIdentification.TypeDni = ""

	return []tableCreditPaymentTest{
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
			&reqCreditoTokenNoValido,
			&pagoCreated,
			&cuentaValida,
			[]callsMocked{
				{"SERV", "SolicitarToken", []interface{}{tokenReqNoValido}, nil, fmt.Errorf("Error al solicitar token")},
			},
			nil,
		},
	}
}

func TestCreditPaymentErrors(t *testing.T) {
	assertions := assert.New(t)

	serv := mockservice.MockPrismaService{}
	utilServ := mockservice.MockUtilService{}

	payment := checkout.NewCreditPayment(&serv, &utilServ)

	table := iniciarCreditPayment()

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
func TestCreditPaymentValid(t *testing.T) {
	assertions := assert.New(t)
	serv := mockservice.MockPrismaService{}
	utilServ := mockservice.MockUtilService{}
	payment := checkout.NewCreditPayment(&serv, &utilServ)

	t.Run("Cuando el llamado a la api devuelve error debo devolver objeto con el error en el status y paidAt en fecha 0", func(t *testing.T) {
		cuentaNoValida := cuentaValida
		cuentaNoValida.Cuenta = ""
		payReqNoValido := payReq
		payReqNoValido.PagoSimple.EstablishmentName = ""
		serv.On("SolicitarToken", tokenReq).Return(pagoToken, nil)
		serv.On("Payments", payReqNoValido).Return(nil, fmt.Errorf("Error al procesar pago"))

		pay, _ := payment.CreateResultado(&requestCreditoValido, &pagoCreated, &cuentaNoValida, "aaaa-bbb-cccc-ddddd", &dtos.InstallmentDetailsResponse{})

		assertions.Equal(pay.StateComment, "Error al procesar pago")
		assertions.Zero(pay.PaidAt)
	})

	t.Run("Cuando el llamado a la api devuelve status ok devuelvo el status, el tiket, el authorization code y paidAt distinto de 0", func(t *testing.T) {
		serv.On("SolicitarToken", tokenReq).Return(pagoToken, nil)
		serv.On("Payments", payReq).Return(payResponse, nil)

		pay, _ := payment.CreateResultado(&requestCreditoValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", &dtos.InstallmentDetailsResponse{})

		assertions.Equal(pay.StateComment, "approved")
		assertions.NotZero(pay.TicketNumber)
		assertions.NotZero(pay.AuthorizationCode)
		assertions.NotZero(pay.PaidAt)
	})
}
