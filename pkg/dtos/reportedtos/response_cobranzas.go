package reportedtos

import "time"

type ResponseCobranzas struct {
	AccountId      string                    `json:"account_id"`
	ReportDate     time.Time                 `json:"report_date"`
	TotalCollected float64                   `json:"total_collected"`
	TotalGrossFee  float64                   `json:"total_gross_fee"`
	TotalNetAmount float64                   `json:"total_net_amount"`
	Data           []ResponseDetalleCobranza `json:"data"`
}

type ResponseDetalleCobranza struct {
	InformedDate      time.Time `json:"informed_date"`
	RequestId         int64     `json:"request_id"`
	ExternalReference string    `json:"external_reference"`
	PayerName         string    `json:"payer_name"`
	Description       string    `json:"description"`
	PaymentDate       time.Time `json:"payment_date"`
	Channel           string    `json:"channel"`
	AmountPaid        float64   `json:"amount_paid"`
	NetFee            float64   `json:"net_fee"`
	IvaFee            float64   `json:"iva_fee"`
	NetAmount         float64   `json:"net_amount"`
	AvailableAt       string    `json:"available_at"`
}
