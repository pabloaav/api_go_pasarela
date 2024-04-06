package dtos_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/stretchr/testify/assert"
)

func TestValidarResultadoRequest(t *testing.T) {
	type tableTest struct {
		Nombre        string
		Params        dtos.ResultadoRequest
		ExpectedValue string
	}

	tabla := []tableTest{
		{
			Nombre: "Debe indicar un error si falta el channel mediante el cual va a pagar",
			Params: dtos.ResultadoRequest{
				Channel:      "",
				Uuid:         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
				Cbu:          "013322215665465654654",
				CardNumber:   "423321354654444",
				HolderDocNum: "32880325",
				HolderCuit:   "23328803259",
			},
			ExpectedValue: "debe indicar el método por el cual va a pagar",
		},
		{
			Nombre: "Debe indicar un error si falta el uuid del pago",
			Params: dtos.ResultadoRequest{
				Channel:      "debin",
				Uuid:         "",
				Cbu:          "013322215665465654654",
				CardNumber:   "423321354654444",
				HolderDocNum: "32880325",
				HolderCuit:   "23328803259",
			},
			ExpectedValue: "debe indicar el código identificador del pago",
		},
		{
			Nombre: "Debe indicar un error si falta el cbu en un pago por debin",
			Params: dtos.ResultadoRequest{
				Channel:      "debin",
				Uuid:         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
				Cbu:          "",
				CardNumber:   "423321354654444",
				HolderDocNum: "32880325",
				HolderCuit:   "23328803259",
			},
			ExpectedValue: "debe indicar el número de CBU",
		},
		{
			Nombre: "Debe indicar un error si falta el número de tarjeta cuando se paga por credito",
			Params: dtos.ResultadoRequest{
				Channel:      "credit",
				Uuid:         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
				Cbu:          "013322215665465654654",
				CardNumber:   "",
				HolderDocNum: "32880325",
				HolderCuit:   "23328803259",
			},
			ExpectedValue: "debe indicar el número de la tarjeta",
		},
		{
			Nombre: "Debe indicar un error si falta el número de tarjeta cuando se paga por debito",
			Params: dtos.ResultadoRequest{
				Channel:      "debit",
				Uuid:         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
				Cbu:          "013322215665465654654",
				CardNumber:   "",
				HolderDocNum: "32880325",
				HolderCuit:   "23328803259",
			},
			ExpectedValue: "debe indicar el número de la tarjeta",
		},
		{
			Nombre: "Debe indicar un error si falta el número de identificación cuando se paga offline",
			Params: dtos.ResultadoRequest{
				Channel:      "offline",
				Uuid:         "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
				Cbu:          "013322215665465654654",
				CardNumber:   "321321321321321",
				HolderDocNum: "",
				HolderCuit:   "",
			},
			ExpectedValue: "debe indicar un número de identificación",
		},
	}

	for _, test := range tabla {
		t.Run(test.Nombre, func(t *testing.T) {
			err := test.Params.Validar()
			assert.Equal(t, err.Error(), test.ExpectedValue)
		})
	}
}
