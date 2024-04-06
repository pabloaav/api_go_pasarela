package prismadtos

import prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"

type UnPagoResponse struct {
	ID                             int64                    `json:"id"`
	SiteTransactionID              string                   `json:"site_transaction_id"`
	PaymentMethodID                int64                    `json:"payment_method_id"`
	CardBrand                      string                   `json:"card_brand"`
	Amount                         int64                    `json:"amount"`
	Currency                       string                   `json:"currency"`
	Status                         string                   `json:"status"`
	StatusDetails                  prismadtos.StatusDetails `json:"status_details"`
	Date                           string                   `json:"date"`
	Customer                       interface{}              `json:"customer"`
	Bin                            string                   `json:"bin"`
	Installments                   int64                    `json:"installments"`
	FirstInstallmentExpirationDate interface{}              `json:"first_installment_expiration_date"`
	PaymentType                    string                   `json:"payment_type"`
	SubPayments                    []interface{}            `json:"sub_payments"`
	SiteID                         string                   `json:"site_id"`
	FraudDetection                 interface{}              `json:"fraud_detection"`
	AggregateData                  interface{}              `json:"aggregate_data"`
	EstablishmentName              interface{}              `json:"establishment_name"`
	Spv                            interface{}              `json:"spv"`
	Confirmed                      interface{}              `json:"confirmed"`
	Pan                            string                   `json:"pan"`
	CustomerToken                  interface{}              `json:"customer_token"`
	CardData                       string                   `json:"card_data"`
	EmvIssuerData                  interface{}              `json:"emv_issuer_data"`
	Token                          string                   `json:"token"`
}
