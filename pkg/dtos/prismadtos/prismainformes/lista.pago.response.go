package prismadtos

import prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"

type ListaPagosResponse struct {
	Limit   int64    `json:"limit"`
	Offset  int64    `json:"offset"`
	Results []Result `json:"results"`
	HasMore bool     `json:"hasMore"`
}

type Result struct {
	ID                int64                    `json:"id"`
	SiteTransactionID string                   `json:"site_transaction_id"`
	Token             string                   `json:"token"`
	UserID            interface{}              `json:"user_id"`
	CardBrand         string                   `json:"card_brand"`
	Bin               string                   `json:"bin"`
	Amount            int64                    `json:"amount"`
	Currency          string                   `json:"currency"`
	Installments      int64                    `json:"installments"`
	Description       string                   `json:"description"`
	PaymentType       string                   `json:"payment_type"`
	SubPayments       []interface{}            `json:"sub_payments"`
	Status            string                   `json:"status"`
	StatusDetails     prismadtos.StatusDetails `json:"status_details"`
	Date              string                   `json:"date"`
	MerchantID        interface{}              `json:"merchant_id"`
	FraudDetection    interface{}              `json:"fraud_detection"`
}

// type FraudDetection struct {
// }
