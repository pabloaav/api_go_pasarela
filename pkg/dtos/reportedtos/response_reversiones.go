package reportedtos

type ResponseReversiones struct {
	AccountID       string                       `json:"account_id"`
	ReportDate      string                       `json:"report_date"`
	TotalChargeback float64                      `json:"total_chargeback"`
	Data            []ResponseDetalleReversiones `json:"data"`
}

type ResponseDetalleReversiones struct {
	InformedDate      string  `json:"informed_date"`
	RequestID         int     `json:"request_id"`
	ExternalReference string  `json:"external_reference"`
	PayerName         string  `json:"payer_name"`
	Description       string  `json:"description"`
	Channel           string  `json:"channel"`
	RevertedAmount    float64 `json:"reverted_amount"`
}
