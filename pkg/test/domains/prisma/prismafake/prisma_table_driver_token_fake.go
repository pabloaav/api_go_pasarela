package prismafake

import (
	"encoding/json"

	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

func EstructuraTokenPagoTarjetaFakeNroTarjeta() (tableDriverTestToken TableDriverTest) {

	tableDriverTestToken = TableDriverTest{

		TituloPrueba: "verificar ingreso de número de tarjeta invalido",
		WantTable:    prismatransacciones.ERROR_NUMBER_CARD, //"el número de tarjeta no es valido",
		DataPrueba:   []string{"4507990000004", "", " "},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "",
				CardExpirationMonth: "08",
				CardExpirationYear:  "21",
				SecurityCode:        "123",
				CardHolderName:      "John Doe",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "simple",
		},
	}
	return
}

func EstructuraTokenPagoTarjetaFakeExpirationMonth() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso de mes de expiración invalido",
		WantTable:    prismatransacciones.ERROR_DATE_CARD, //"la fecha de vencimiento de la tarjeta es invalido",
		DataPrueba:   []string{"", " ", "19"},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "4507990000004905",
				CardExpirationMonth: "08",
				CardExpirationYear:  "12",
				SecurityCode:        "123",
				CardHolderName:      "John Doe",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "simple",
		},
	}
	return
}

func EstructuraTokenPagoTarjetaFakeExpirationYear() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso del año de expiración invalido",
		WantTable:    prismatransacciones.ERROR_DATE_CARD, //"la fecha de vencimiento de la tarjeta es invalido",
		DataPrueba:   []string{"", " ", "12", "19"},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "4507990000004905",
				CardExpirationMonth: "05",
				CardExpirationYear:  "21",
				SecurityCode:        "123",
				CardHolderName:      "John Doe",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "simple",
		},
	}
	return
}
func EstructuraTokenPagoTarjetaFakeHolderName() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso de nombre invalido",
		WantTable:    prismatransacciones.ERROR_HOLDER_NAME, //"nombre no valido",
		DataPrueba:   []string{"", " "},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "4507990000004905",
				CardExpirationMonth: "07",
				CardExpirationYear:  "21",
				SecurityCode:        "123",
				CardHolderName:      "John Doe",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "simple",
		},
	}
	return
}

func EstructuraTokenPagoTarjetaFakeTypePay() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso de un tipo de pago invalido",
		WantTable:    "tipo de pago no valido",
		DataPrueba:   []string{"", " ", "otro tipo"},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "4507990000004905",
				CardExpirationMonth: "07",
				CardExpirationYear:  "21",
				SecurityCode:        "123",
				CardHolderName:      "",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "",
		},
	}
	return
}

func EstructuraTokenPagoTarjetaFakeCardEmpty() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar envio de estructura Card vacía cuando el tipo de pago es simple",
		WantTable:    "los datos recibidos son incorrectos",
		DataPrueba:   []string{""},
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "",
				CardExpirationMonth: "",
				CardExpirationYear:  "",
				SecurityCode:        "",
				CardHolderName:      "",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "",
					NumberDni: "",
				},
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	tabla de datos pago offline
*/
func EstructuraTokenPagoOffLineFakeDni() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso de número de documento invalido",
		WantTable:    prismatransacciones.ERROR_NRO_DOC, //"el número de documento ingresado no es valido",
		DataPrueba:   []string{"", " ", "26458", "sdsdffsf", "wswqw"},
		TokenStructura: prismatransacciones.StructToken{
			DataOffline: prismatransacciones.OfflineTokenRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "",
					},
					Name: "torres pablo",
				},
			},
			TypePay: "offline",
		},
	}
	return
}

func EstructuraTokenPagoOffLineFakeNonbre() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso del nombre invalido",
		WantTable:    prismatransacciones.ERROR_HOLDER_NAME, //"nombre no valido",
		DataPrueba:   []string{"", " "},
		TokenStructura: prismatransacciones.StructToken{
			DataOffline: prismatransacciones.OfflineTokenRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "26458796",
					},
					Name: "",
				},
			},
			TypePay: "offline",
		},
	}
	return
}

func EstructuraTokenPagoOffLineFakePagoOffline() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar envio de estructura pago offline vacía cuando el tipo de pago es offline",
		WantTable:    "los datos recibidos son incorrectos",
		DataPrueba:   []string{},
		TokenStructura: prismatransacciones.StructToken{
			DataOffline: prismatransacciones.OfflineTokenRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "",
						Number: "",
					},
					Name: "",
				},
			},
			TypePay: "offline",
		},
	}
	return
}

func EstructuraTokenPagoOffLineFakeTipoDni() (tableDriverTestToken TableDriverTest) {
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar que el tipo de documento ingresado sea invalid",
		WantTable:    prismatransacciones.ERROR_TIPO_DOCUMENTO, //"tipo de documento seleccionado inválido",
		DataPrueba:   []string{"D.N.I", " ", "26458", "sdsdffsf"},
		TokenStructura: prismatransacciones.StructToken{
			DataOffline: prismatransacciones.OfflineTokenRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "",
						Number: "26589452",
					},
					Name: "torres pablo",
				},
			},
			TypePay: "offline",
		},
	}
	return
}

/*
	estructura valida pago simple
*/
func EstructuraTokenPagoTarjetaFaValido() (tableDriverTestToken TableDriverTest) {
	response := prismatransacciones.PagoToken{
		Id:                "11002b27-ad1f-4750-aa88-8925c46fca07",
		ValidationResulto: false,
		Status:            "active",
		CardNumberLength:  16,
		Bin:               "450799",
		DateCreate:        "2021-06-24T11:57Z",
		LastFourDigits:    "4905",
		ExpirationMonth:   8,
		ExpirationYear:    21,
		DateDue:           "2021-06-24T12:12Z",
		CardHolder: prismatransacciones.CardHolder{
			Identification: prismatransacciones.Identification{
				TypeDni:   "dni",
				NumberDni: "25123456",
			},
			Name: "John Doe",
		},
	}
	wantResponse, _ := json.Marshal(response)
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "test de respuesta del repositorio remoto PostSolicitudTokenPago",
		WantTable:    string(wantResponse),
		TokenStructura: prismatransacciones.StructToken{
			Card: prismatransacciones.Card{
				CardNumber:          "4507990000004905",
				CardExpirationMonth: "08",
				CardExpirationYear:  "21",
				SecurityCode:        "123",
				CardHolderName:      "John Doe",
				CardHolderIdentification: prismatransacciones.CardHolderIdentification{
					TypeDni:   "DNI",
					NumberDni: "25123456",
				},
			},
			TypePay: "simple",
		},
	}
	return
}

/*
	estructura valida pago offline
*/
func EstructuraTokenPagoOffLineValido() (tableDriverTestToken TableDriverTest) {

	response := prismatransacciones.OfflineTokenResponse{
		ID:          "e1268846-f809-49e8-9b0a-a6f689b90ea4",
		Status:      "active",
		DateCreated: "2021-06-24T12:19Z",
		DateDue:     "2021-06-24T12:34Z",
		Customer: prismatransacciones.DataCustomer{
			Identification: prismatransacciones.IdentificationCustomer{
				Type:   "DNI",
				Number: "32880325",
			},
			Name: "Castro Fernando",
		},
	}
	wantResponse, _ := json.Marshal(response)
	tableDriverTestToken = TableDriverTest{
		TituloPrueba: "verificar ingreso de número de documento invalido",
		WantTable:    string(wantResponse),
		TokenStructura: prismatransacciones.StructToken{
			DataOffline: prismatransacciones.OfflineTokenRequest{
				Customer: prismatransacciones.DataCustomer{
					Identification: prismatransacciones.IdentificationCustomer{
						Type:   "DNI",
						Number: "32880325",
					},
					Name: "torres pablo",
				},
			},
			TypePay: "offline",
		},
	}
	return
}
