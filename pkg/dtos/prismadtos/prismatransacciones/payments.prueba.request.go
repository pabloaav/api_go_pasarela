package prismadtos

type PaymentsPruebaRequest struct {
	SiteTransactionID string        `json:"site_transaction_id"`
	Token             string        `json:"token"`
	PaymentMethodID   int64         `json:"payment_method_id"`
	Bin               string        `json:"bin"`
	Amount            int64         `json:"amount"`
	Currency          string        `json:"currency"`
	Installments      int64         `json:"installments"`
	Description       string        `json:"description"`
	PaymentType       string        `json:"payment_type"`
	SubPayments       []interface{} `json:"sub_payments"`
}
