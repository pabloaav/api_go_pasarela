package prismafake

import (
	"encoding/json"

	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

/*
	datos de prueba PaymentsPagoOffLine
*/

/*
	fake tipo de documento ingresado
*/

func EstructuraPaymentsPagoOffLineFakeTipoDocumento() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de tipo documento",
		WantTable:        prismatransacciones.ERROR_TIPO_DOCUMENTO, //"tipo de documento seleccionado inválido",
		DataPruebaString: []string{"dni", "D.N.I", " ", "26458", "sdsdffsf"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake número del documento
*/
func EstructuraPaymentsPagoOffLineFakeDocumento() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar número de documento invalido",
		WantTable:        prismatransacciones.ERROR_NRO_DOC, //"el número de documento ingresado no es valido",
		DataPruebaString: []string{"", " ", "26458", "sdsdffsf", "wswqw", "23445vfv"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake nombre del pagoador
*/
func EstructuraPaymentsPagoOffLineFakeCustomerName() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de nombre invalido",
		WantTable:        prismatransacciones.ERROR_NOMBRE_PAGADOR, //"nombre ingresado no debe ser vacío",
		DataPruebaString: []string{"", " "},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake SiteTransactionID
*/

func EstructuraPaymentsPagoOffLineFakeSiteTransactionId() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar logitud de site transation di invalido",
		WantTable:        prismatransacciones.ERROR_SITE_TRANSACTION_ID,
		DataPruebaString: []string{"", "eda59288-9f7d-425f-bee0-78a16dbe981cbngdf4g5g4d"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake token
*/
func EstructuraPaymentsPagoOffLineFakeToken() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar envio de token invalido",
		WantTable:        prismatransacciones.ERROR_TOKEN_PAGO, //"token de pago no valido",
		DataPruebaString: []string{"", " ", "68af9aaa-340f-4af6-8afb-a18620fb8ab5dfgfdgdfgdg3453g34t3"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake validar longitud Amount
*/
func EstructuraPaymentsPagoOffLineFakeAmount() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:  "verificar ingreso de monto invalido",
		WantTable:     prismatransacciones.ERROR_AMOUNT,
		DataPruebaInt: []int64{-1, 45821152545, 345625454, 999999999, 0},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            0,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},
			TypePay: "offline",
		},
	}
	return
}

/*
	fake email
*/
func EstructuraPaymentsPagoOffLineFakeEmail() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de email invalido",
		WantTable:        prismatransacciones.ERROR_EMAIL,
		DataPruebaString: []string{"dsfsf@", "sdfsfsf.com", " ", "", "nombredeusuario", "@fsff.com", "ddsffsd@.com"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},
			TypePay: "offline",
		},
	}
	return
}

/*
	fake Currency
*/
func EstructuraPaymentsPagoOffLineFakeCurrency() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso del tipo de moneda invalido",
		WantTable:        prismatransacciones.ERROR_CURRENCY,
		DataPruebaString: []string{"", " ", "A.R.S", "ars", "U.S.D", "usd", "fsfsdfsf", "SDFSDFF"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake CodP3
*/
func EstructuraPaymentsPagoOffLineFakeCodP3() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de los días entre vencimientos, invalido",
		WantTable:        prismatransacciones.ERROR_CODP3,
		DataPruebaString: []string{"sds", "", " ", "2323", "2w", "w2", "1"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake CodP4
*/
func EstructuraPaymentsPagoOffLineFakeCodP4() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar Días después del 1º vencimiento y hasta que el cliente pueda abonar, invalido",
		WantTable:        prismatransacciones.ERROR_CODP4,
		DataPruebaString: []string{"sds", "", " ", "2323", "2w", "w2", "23"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake número de cliente
*/
func EstructuraPaymentsPagoOffLineFakeCliente() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar número de cliente invalido",
		WantTable:        prismatransacciones.ERROR_CLIENTE_NRO,
		DataPruebaString: []string{"sds", "", " ", "2323", "2324433w", "s2er34w2", "1223236743", "qwertyui"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake surchange
*/
func EstructuraPaymentsPagoOffLineFakeSurcharge() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:  "verificar recargo invalido",
		WantTable:     prismatransacciones.ERROR_SURCHANGE,
		DataPruebaInt: []int64{-14245, -1, 457896365, -01},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         0,
				PaymentMode:       "offline",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	fake paymentMode
*/

func EstructuraPaymentsPagoOffLineFakePaymentMode() (tableDriverTestPeyment TableDriverTestPayment) {
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba:     "verificar ingreso de número de tarjeta invalido",
		WantTable:        prismatransacciones.ERROR_MODO_PAGO,
		DataPruebaString: []string{"Offline", "", " ", "OFFLINE", "offLine", "Off"},
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "Castro Fernando",
				},
				SiteTransactionID: "eda5-7d4", //"eda59288-9f7d-425f-bee0-78a16dbe981cbn",
				Token:             "07f2af05-37fa-49c8-bb88-7036eb10f8ed",
				PaymentMethodID:   26,
				Amount:            1100,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "191123",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "",
			},

			TypePay: "offline",
		},
	}
	return
}

/*
	estructura valida payments offline
*/
func EstructuraPaymentsPagoOffLineFakeValido() (tableDriverTestPeyment TableDriverTestPayment) {
	response := prismatransacciones.PaymentsOfflineResponse{
		ID:                10751178,
		SiteTransactionID: "TelCo-58",
		Token:             "38b9a383-6159-4fb7-9a7e-a3effa6dfcb0",
		PaymentMethodID:   26,
		Amount:            1000,
		Currency:          "ARS",
		Email:             "asd@asd.com",
		Status:            "invoice_generated",
		StatusDetails: prismatransacciones.StatusDetails{
			Error: prismatransacciones.ErrorDetails{
				Type: "",
				Reason: prismatransacciones.Reason{
					ID:                    0,
					Description:           "",
					AdditionalDescription: "",
				},
			},
		},
		Date:                    "2021-07-01T10:11Z",
		InvoiceExpiration:       "290621",
		SecondInvoiceExpiration: "",
		Surcharge:               123456,
		Client:                  "12345678",
		Barcode:                 "85900121234567832880325290621000010001001234561343",
	}
	wantResponse, _ := json.Marshal(response)
	tableDriverTestPeyment = TableDriverTestPayment{
		TituloPrueba: "verificar el envio de un objeto pago con tarjeta valido",
		WantTable:    string(wantResponse),
		PaymentStructura: prismatransacciones.StructPayments{
			PagoOffline: prismatransacciones.PaymentsOfflineRequest{
				Customer: prismatransacciones.DataCustomer{
					Name: "Castro Fernando",
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
				},
				SiteTransactionID: "TelCo-58",
				Token:             "38b9a383-6159-4fb7-9a7e-a3effa6dfcb0",
				PaymentMethodID:   26,
				Amount:            1000,
				Currency:          "ARS",
				PaymentType:       "single",
				Email:             "asd@asd.com",
				InvoiceExpiration: "290621",
				CodP3:             "10",
				CodP4:             "134",
				Client:            "12345678",
				Surcharge:         123456,
				PaymentMode:       "offline",
			},
			TypePay: "offline",
		},
	}
	return
}
