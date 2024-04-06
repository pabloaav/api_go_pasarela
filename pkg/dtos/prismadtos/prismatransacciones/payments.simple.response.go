package prismadtos

type PaymentsSimpleResponse struct {
	ID                int64           `json:"id,omitempty"`
	SiteTransactionID string          `json:"site_transaction_id,omitempty"`
	Token             string          `json:"token,omitempty"`
	Customer          Customer        `json:"customer,omitempty"`
	PaymentMethodID   int64           `json:"payment_method_id,omitempty"`
	Bin               string          `json:"bin,omitempty"`
	Amount            int64           `json:"amount,omitempty"`
	Currency          string          `json:"currency,omitempty"`
	Installments      int64           `json:"installments,omitempty"`
	PaymentType       EnumPaymentType `json:"payment_type"`
	SubPayments       []interface{}   `json:"sub_payments,omitempty"`
	Status            string          `json:"status,omitempty"`
	Confirmed         interface{}     `json:"confirmed,omitempty"`
	StatusDetails     StatusDetails   `json:"status_details,omitempty"`
	FraudDetection    FraudDetection  `json:"fraud_detection,omitempty"`
	AggregateData     interface{}     `json:"aggregate_data,omitempty"`
	Pan               string          `json:"pan,omitempty"`
	CustomerToken     string          `json:"customer_token,omitempty"`

	CardBrand                      string      `json:"card_brand,omitempty"`
	Date                           string      `json:"date,omitempty"`
	FirstInstallmentExpirationDate interface{} `json:"first_installment_expiration_date,omitempty"`
	SiteID                         string      `json:"site_id,omitempty"`
	EstablishmentName              interface{} `json:"establishment_name,omitempty"`
	Spv                            interface{} `json:"spv,omitempty"`
	CardData                       string      `json:"card_data,omitempty"`
}

type Customer struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

type FraudDetection struct {
	Status string `json:"status,omitempty"`
}
