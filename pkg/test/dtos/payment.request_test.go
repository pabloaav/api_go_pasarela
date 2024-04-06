package dtos_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/stretchr/testify/assert"
)

var (
	ValidPaymentRequest = dtos.PagoRequest{
		PayerName:         "Fernando Castro",
		Description:       "Pago de obligaciones con DGR",
		FirstTotal:        100050,
		FirstDueDate:      "01-07-2021",
		ExternalReference: "111",
		SecondDueDate:     "01-08-2021",
		SecondTotal:       105050,
		PayerEmail:        "fernando.castro@telco.com.ar",
		PaymentType:       "sellos",
		Items: []entities.Pagoitems{
			{
				Quantity:    1,
				Description: "Item 1 impuestos",
				Amount:      50000,
			},
			{
				Quantity:    1,
				Description: "Item 2 intereses",
				Amount:      50050,
			},
		},
	}
)

func TestPaymentRequestValidation(t *testing.T) {

	t.Run("Debe devolver un error cuando falte el nombre del pagador", func(t *testing.T) {
		request := ValidPaymentRequest
		request.PayerName = ""

		err := request.Validar()

		assert.EqualError(t, err, "se debe indicar el nombre del pagador")
	})

	t.Run("Debe devolver un error cuando falte el concepto a pagar", func(t *testing.T) {
		request := ValidPaymentRequest
		request.Description = ""

		err := request.Validar()

		assert.EqualError(t, err, "se debe indicar el concepto a pagar")
	})

	t.Run("Debe devolver un error cuando falte el monto a pagar o sea igual o menor a 0", func(t *testing.T) {
		request := ValidPaymentRequest
		valores := []int64{-500, 0}

		for _, v := range valores {
			request.FirstTotal = v
			err := request.Validar()
			assert.EqualError(t, err, "se debe indicar el monto a pagar")
		}

	})

	t.Run("Debe devolver un error cuando falte primera fecha de vencimiento o tenga menos de 10 caracteres", func(t *testing.T) {
		request := ValidPaymentRequest
		valores := []string{"010101", "21212121", "123456789"}

		for _, v := range valores {
			request.FirstDueDate = v
			err := request.Validar()
			assert.EqualError(t, err, "se debe indicar fecha de vencimiento del pago")
		}

	})

	t.Run("Debe devolver un error cuando la primera fecha de vencimiento no tenga el formato indicado", func(t *testing.T) {
		request := ValidPaymentRequest
		valores := []string{"2006-01-02", "2121/21/21", "01/02/2006"}

		for _, v := range valores {
			request.FirstDueDate = v
			err := request.Validar()
			assert.EqualError(t, err, "el formato de fecha de vencimientos no es válido (dia-mes-año)")
		}

	})

	t.Run("Debe devolver un error cuando ingrese la segunda fecha de vencimiento y no tenga el formato indicado", func(t *testing.T) {
		request := ValidPaymentRequest
		valores := []string{"2006-01-02", "2121/21/21", "01/02/2006"}

		for _, v := range valores {
			request.SecondDueDate = v
			err := request.Validar()
			assert.EqualError(t, err, "el formato de fecha de segundo vencimientos no es válido (dia-mes-año)")
		}

	})

	t.Run("Debe devolver un error cuando falte el tipo de pago", func(t *testing.T) {
		request := ValidPaymentRequest
		request.PaymentType = ""

		err := request.Validar()

		assert.EqualError(t, err, "debe indicar el tipo de pago que desea realizar")
	})
}
