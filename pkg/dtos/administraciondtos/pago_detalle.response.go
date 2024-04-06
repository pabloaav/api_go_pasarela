package administraciondtos

import "time"

type PagoDetalleResponse struct {
	ID                int64                  `json:"id"`
	Tipo              string                 `json:"tipo"`
	Estado            string                 `json:"estado"`
	Description       string                 `json:"description"`
	FirstDueDate      time.Time              `json:"first_due_date,omitempty"`
	FirstTotal        float64                `json:"first_total"`
	SecondDueDate     time.Time              `json:"second_due_date,omitempty"`
	SecondTotal       float64                `json:"second_total,omitempty"`
	PayerName         string                 `json:"payer_name"`
	PayerEmail        string                 `json:"payer_email"`
	ExternalReference string                 `json:"external_reference"`
	Metadata          string                 `json:"metadata"`
	Uuid              string                 `json:"uuid"`
	Intentos          []PagoIntentosResponse `json:"intentos"`
}

type PagoIntentosResponse struct {
	ID           int64     `json:"id"`
	Channel      string    `json:"channel"`
	ExternalID   string    `json:"external_id"`
	PaidAt       time.Time `json:"paid_at,omitempty"`
	ReportAt     time.Time `json:"report_at,omitempty"`
	IsAvailable  bool      `json:"is_available"`
	Amount       float64   `json:"amount"`
	GrossFee     float64   `json:"gross_fee"`
	NetFee       float64   `json:"net_fee"`
	FeeIva       float64   `json:"fee_iva"`
	NetAmount    float64   `json:"net_amount"`
	StateComment string    `json:"state_comment"`
	AvailableAt  time.Time `json:"available_at,omitempty"`
	RevertedAt   time.Time `json:"reverted_at,omitempty"`
}
