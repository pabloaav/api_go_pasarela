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
	customer = prismadtos.DataCustomer{
		Identification: prismadtos.IdentificationCustomer{
			Type:   "DNI",
			Number: "32880325",
		},
		Name: "Fernando Castro",
	}

	tokenOfflineReq = prismadtos.StructToken{
		DataOffline: prismadtos.OfflineTokenRequest{
			Customer: customer,
		},
		TypePay: "offline",
	}

	tokenOfflineResponse = prismadtos.OfflineTokenResponse{
		ID:          "1",
		Status:      "ok",
		DateCreated: "080721",
		DateDue:     "080721",
		Customer:    customer,
	}

	payOfflineReq = prismadtos.StructPayments{
		PagoOffline: prismadtos.PaymentsOfflineRequest{
			Customer:          customer,
			SiteTransactionID: requestOfflineValido.Uuid,
			Token:             "1",
			PaymentMethodID:   requestOfflineValido.PaymentMethodID,
			//PaymentMethodID: 26,
			Amount:      requestOfflineValido.Importe,
			Currency:    "ARS",
			PaymentType: "single",
			Email:       requestOfflineValido.HolderEmail,
			//InvoiceExpiration: invoiceExpiration,
			InvoiceExpiration: "191123",
			//CodP3:             fmt.Sprint(daysBetweenDueDates),
			CodP3: "10",
			//CodP4: fmt.Sprintf("%0*d", 3, daysBetweenDueDates+100),
			CodP4: "134",
			//Client: fmt.Sprint(cuenta.ID),
			Client: "12345678",
			//Surcharge:   int64(pago.SecondTotal) * 100,
			Surcharge:   1234567,
			PaymentMode: "offline",
		},
		TypePay: "offline",
	}

	payOfflineResponse = prismadtos.PaymentsOfflineResponse{
		ID:                      1,
		SiteTransactionID:       "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		Token:                   "1",
		PaymentMethodID:         5,
		Amount:                  105050,
		Currency:                "ARS",
		Email:                   "",
		Status:                  "invoice_generated",
		Date:                    "2021-07-07T15:04Z",
		InvoiceExpiration:       "070821",
		SecondInvoiceExpiration: "070921",
		Barcode:                 "85900121234567832880325191123000960511012345671349",
	}
)

type tableOfflinePaymentTest struct {
	Nombre             string
	Error              error
	Request            *dtos.ResultadoRequest
	Pago               *entities.Pago
	Cuenta             *entities.Cuenta
	MethodMocked       []callsMocked
	InstallmentDetails *dtos.InstallmentDetailsResponse
}

func iniciarOfflinePayment() []tableOfflinePaymentTest {
	reqOfflineTokenNoValido := requestOfflineValido
	reqOfflineTokenNoValido.HolderDocType = ""
	tokenReqNoValido := tokenOfflineReq
	tokenReqNoValido.DataOffline.Customer.Identification.Type = ""

	return []tableOfflinePaymentTest{
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
			&reqOfflineTokenNoValido,
			&pagoCreated,
			&cuentaValida,
			[]callsMocked{
				{"SERV", "SolicitarToken", []interface{}{tokenReqNoValido}, nil, fmt.Errorf("Error al solicitar token")},
			},
			nil,
		},
	}
}

func TestOfflinePaymentErrors(t *testing.T) {
	assertions := assert.New(t)

	serv := mockservice.MockPrismaService{}

	payment := checkout.NewOfflinePayment(&serv)

	table := iniciarOfflinePayment()

	for _, v := range table {
		t.Run(v.Nombre, func(t *testing.T) {
			for _, c := range v.MethodMocked {
				switch c.mock {
				case "SERV":
					serv.On(c.nombre, c.filtro...).Return(c.retorno, c.mensaje).Once()
				}
			}

			res, err := payment.CreateResultado(v.Request, v.Pago, v.Cuenta, "aaaa-bbb-cccc-ddddd", v.InstallmentDetails)

			if err != nil {
				assertions.EqualError(err, v.Error.Error())
			} else {
				assertions.NotZero(res.PaidAt)
			}

		})
	}
}

func TestOfflinePaymentValid(t *testing.T) {
	assertions := assert.New(t)
	serv := mockservice.MockPrismaService{}
	payment := checkout.NewOfflinePayment(&serv)

	t.Run("Cuando el llamado a la api devuelve error debo devolver objeto con el error en el status y paidAt en fecha 0", func(t *testing.T) {
		reqNoValido := requestOfflineValido
		reqNoValido.Importe = 0
		payReqNoValido := payOfflineReq
		payReqNoValido.PagoOffline.Amount = 0
		var installmentsDetails *dtos.InstallmentDetailsResponse
		serv.On("SolicitarToken", tokenOfflineReq).Return(tokenOfflineResponse, nil)
		serv.On("Payments", payReqNoValido).Return(nil, fmt.Errorf("Error al procesar pago"))

		pay, _ := payment.CreateResultado(&reqNoValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", installmentsDetails)

		assertions.Equal(pay.StateComment, "Error al procesar pago")
		assertions.Zero(pay.PaidAt)
	})

	t.Run("Cuando el llamado a la api devuelve status ok devuelvo el status, el codigo de barras y paidAt distinto de 0", func(t *testing.T) {
		serv.On("SolicitarToken", tokenOfflineReq).Return(tokenOfflineResponse, nil)
		serv.On("Payments", payOfflineReq).Return(payOfflineResponse, nil)
		var installmentsDetails *dtos.InstallmentDetailsResponse
		pay, _ := payment.CreateResultado(&requestOfflineValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", installmentsDetails)

		assertions.Equal(pay.StateComment, "invoice_generated")
		assertions.NotZero(pay.Barcode)
		assertions.NotZero(pay.PaidAt)
	})
}
