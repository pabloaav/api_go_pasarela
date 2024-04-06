package dtos

import (
	"errors"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

// PagoRequest utilizado como parametros de solicitud de un pago
type PagoRequest struct {
	PayerName         string               `json:"payer_name"`
	Description       string               `json:"description"`
	FirstTotal        int64                `json:"first_total"`
	FirstDueDate      string               `json:"first_due_date"`
	ExternalReference string               `json:"external_reference"`
	SecondDueDate     string               `json:"second_due_date"`
	SecondTotal       int64                `json:"second_total"`
	PayerEmail        string               `json:"payer_email"`
	PaymentType       string               `json:"payment_type"`
	Metadata          string               `json:"metadata"`
	Items             []entities.Pagoitems `json:"items"`
	Expiration        int64                `json:"expiration"`
}

func (r *PagoRequest) Validar() error {
	if len(r.PayerName) <= 0 {
		return errors.New("se debe indicar el nombre del pagador")
	}
	if len(r.Description) <= 0 {
		return errors.New("se debe indicar el concepto a pagar")
	}
	if r.FirstTotal <= 0 {
		return errors.New("se debe indicar el monto a pagar")
	}
	if len(r.FirstDueDate) <= 9 {
		return errors.New("se debe indicar fecha de vencimiento del pago")
	}
	if r.FirstDueDate[2:3] != "-" || r.FirstDueDate[5:6] != "-" {
		return errors.New("el formato de fecha de vencimientos no es válido (dia-mes-año)")
	}
	if len(r.PaymentType) <= 0 {
		return errors.New("debe indicar el tipo de pago que desea realizar")
	}
	if len(r.ExternalReference) > 10 {
		return errors.New("el campo external reference no debe superar los 10 caracteres")
	}

	if len(r.SecondDueDate) == 0 {
		if r.SecondTotal == 0 {
			r.SecondDueDate = r.FirstDueDate
			r.SecondTotal = r.FirstTotal
		} else {
			return errors.New("se debe indicar fecha del segundo vencimiento del pago")
		}
	} else {
		if r.SecondTotal <= 0 {
			return errors.New("se debe indicar el monto del segundo vencimiento del pago")
		}
	}
	if len(r.SecondDueDate) > 0 {
		if r.SecondDueDate[2:3] != "-" || r.SecondDueDate[5:6] != "-" {
			return errors.New("el formato de fecha de segundo vencimientos no es válido (dia-mes-año)")
		}
	}
	if r.SecondTotal < r.FirstTotal {
		return errors.New("el monto del segundo pago debe ser mayor al monto del primer pago")
	}

	// la expiracion maxima es 300 minutos que son 5 horas
	var expirationMaxTime int64 = 300
	// Si la expiracion del pago es cero, se pone expirationMaxTime
	if r.Expiration == 0 {
		r.Expiration = expirationMaxTime
	}

	// Si la expiracion del pago es superior a expirationMaxTime, se pone expirationMaxTime
	if r.Expiration > expirationMaxTime {
		r.Expiration = expirationMaxTime
	}

	return nil
}

func (r *PagoRequest) ToFormatStr() {
	r.PayerName = strings.ToUpper(r.PayerName)
	r.Description = strings.ToUpper(r.Description)
}
