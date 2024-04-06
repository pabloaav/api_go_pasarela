package entities

import "gorm.io/gorm"

type PaymentRequest struct {
	gorm.Model
	Description         string         `json:"description"`
	ExternalReference   string         `json:"external_reference"`
	Ammount             float64        `json:"ammount"`
	PayerName           string         `json:"payer_name"`
	PayerLastname       string         `json:"payer_lastname"`
	PayerIdentification Identification `json:"payer_identification"`
	Status              string         `json:"status"`
	CardNumberLength    int64          `json:"card_number_length"`
	DateCreated         string         `json:"date_created"`
	Bin                 string         `json:"bin"`
	LastFourDigits      string         `json:"last_four_digits"`
	SecurityCodeLength  int64          `json:"security_code_length"`
	ExpirationMonth     int64          `json:"expiration_month"`
	ExpirationYear      int64          `json:"expiration_year"`
	DateDue             string         `json:"date_due"`
	Cardholder          Cardholder     `json:"cardholder"`
}

type Cardholder struct {
	Identification Identification `json:"identification"`
	Name           string         `json:"name"`
}

type Identification struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}
