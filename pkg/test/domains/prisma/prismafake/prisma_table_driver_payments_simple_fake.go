package prismafake

import (
	"encoding/json"

	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

/*
	datos de prueba PaymentSimple
*/

/*
	Fake SiteTransactionId
*/
func EstructuraPaymentsPagoTarjetaFakeSiteTransactionId() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de SiteTransactionId invalido",
		WantTable:        prismatransacciones.ERROR_SITE_TRANSACTION_ID, //"id de ttansacción es incorrecto",
		DataPruebaString: []string{" ", "", "eda59288-9f7d-425f-bee0-78a16dbe981cbngdf4g5g4d"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake token
*/
func EstructuraPaymentsPagoTarjetaFakeToken() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de de token invalido",
		WantTable:        prismatransacciones.ERROR_TOKEN_PAGO, //"token de pago no valido",
		DataPruebaString: []string{"", " ", "68af9aaa-340f-4af6-8afb-a18620fb8ab5dfgfdgdfgdg3453g34t3"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake Bin
*/
func EstructuraPaymentsPagoTarjetaFakeBin() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de número Bin invalido",
		WantTable:        prismatransacciones.ERROR_BIN, //"el número de bien es incorrecto",
		DataPruebaString: []string{"", " ", "45789625", "123", "1234gb"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake Amount
*/
func EstructuraPaymentsPagoTarjetaFakeAmount() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:  prismatransacciones.ERROR_AMOUNT, //"verificar ingreso de monto invalido",
		WantTable:     "el monto ingresado no es valido",
		DataPruebaInt: []int64{0, 000000000000000000000000000000},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            0,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake Currency
*/
func EstructuraPaymentsPagoTarjetaFakeCurrency() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de de tipo de moneda invalido",
		WantTable:        prismatransacciones.ERROR_CURRENCY, //"la moneda seleccionada no es valida",
		DataPruebaString: []string{"", " ", "A.R.S", "ars", "U.S.D", "usd", "fsfsdfsf", "SDFSDFF"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake Installments
*/
func EstructuraPaymentsPagoTarjetaFakeInstallments() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:  "verificar ingreso de número de cuotas invalido",
		WantTable:     prismatransacciones.ERROR_INSTALLMENTS, //"el valor de cuota ingresado es incorrecto",
		DataPruebaInt: []int64{000, 123, 0},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      0,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake payment type
*/
func EstructuraPaymentsPagoTarjetaFakePaymentType() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de tipo de pago invalido",
		WantTable:        prismatransacciones.ERROR_PAYMENT_TYPE, // "tipo de pago seleccionado inválido",
		DataPruebaString: []string{"saludos", "", " ", "SINGLE", "DISTRIBUTED"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "",
				EstablishmentName: "prueba desa soft",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake  establishment name
*/
func EstructuraPaymentsPagoTarjetaFakeEstablishmentName() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de número de tarjeta invalido",
		WantTable:        prismatransacciones.ERROR_NOMBRE_ESTABLECIMIENTO, //"el número de tarjeta no es valido",
		DataPruebaString: []string{"SFDSFSÑ", "", " ", "maní", "sdvsdácvSDF", "CDSercvxcvefv454534bfg4tg4a", "CDSercvxcvefv454534bfg4tgñ", "CDSercvxcvefv454534bfg4tá", "CDSercvxcvefv454534bfg-sds343434"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com.ar"},
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	Fake Email
*/
func EstructuraPaymentsPagoTarjetaFakeEmail() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de mail invalido",
		WantTable:        prismatransacciones.ERROR_EMAIL, //"email no valido",
		DataPruebaString: []string{"dsfsf@", "sdfsfsf.com", " ", "", "nombredeusuario", "@fsff.com", "ddsffsd@.com"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "123"},
				SiteTransactionID: "eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				SiteID:            "21",
				Token:             "68af9aaa-340f-4af6-8afb-a18620fb8ab5",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            12545,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probando pago",
				PaymentType:       "single",
				EstablishmentName: "prueba desa soft",

				Customeremail: prismatransacciones.Customeremail{Email: ""},
				SubPayments:   make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	estructura valida payments simple
*/
func EstructuraPaymentsPagoTarjetaFakeValido() (tableDriverTestPeyment TableDriverTestPayment) {
	response := prismatransacciones.PaymentsSimpleResponse{
		ID:                10751099,
		SiteTransactionID: "TELCO_PRUEBA2_prueba456",
		Token:             "5ef79bd4-6ab3-4be8-b0df-c9afa72e892f",
		Customer:          prismatransacciones.Customer{Email: "prisma@prismamp.com"},
		PaymentMethodID:   1,
		Bin:               "450799",
		Amount:            9598746,
		Currency:          "ars",
		Installments:      1,
		PaymentType:       "single",
		Status:            "approved",
		StatusDetails: prismatransacciones.StatusDetails{
			Ticket:                "3977",
			CardAuthorizationCode: "085027",
			AddressValidationCode: "VTE0011",
			Error: prismatransacciones.ErrorDetails{
				Type: "",
				Reason: prismatransacciones.Reason{
					ID:                    0,
					Description:           "",
					AdditionalDescription: "",
				},
			},
		},
		FraudDetection:    prismatransacciones.FraudDetection{},
		Pan:               "345425f15b2c7c4584e0044357b6394d7e",
		CardBrand:         "Visa",
		Date:              "2021-07-01T08:50Z",
		SiteID:            "99999966",
		EstablishmentName: "pago tarjeta prueba",
		CardData:          "/tokens/10751099",
	}
	wantResponse, _ := json.Marshal(response)
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba: "verificar el envio de un objeto pago con tarjeta valido",
		WantTable:    string(wantResponse),
		PaymentStructura: prismatransacciones.StructPayments{
			PagoSimple: prismatransacciones.PaymentsSimpleRequest{
				Customerid:        prismatransacciones.Customerid{ID: "1"},
				SiteTransactionID: "TELCO_PRUEBA2_prueba456",
				SiteID:            "0",
				Token:             "5ef79bd4-6ab3-4be8-b0df-c9afa72e892f",
				PaymentMethodID:   1,
				Bin:               "450799",
				Amount:            9598746,
				Currency:          "ARS",
				Installments:      1,
				Description:       "probandopago",
				PaymentType:       "single",
				Customeremail:     prismatransacciones.Customeremail{Email: "prisma@prismamp.com"},
				EstablishmentName: "pago tarjeta prueba",
				SubPayments:       make([]interface{}, 0),
			},
			TypePay: "simple",
		},
	}
	return
}
