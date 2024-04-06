package dtos

// CheckoutResponse respuesta utilizada en el frontend del checkout
type CheckoutResponse struct {
	Estado               string             `json:"estado"`
	Description          string             `json:"description"`
	DueDate              string             `json:"due_date"`
	SecondDueDate        bool               `json:"second_due_date"`
	Total                float64            `json:"first_total"`
	PayerName            string             `json:"payer_name"`
	PayerEmail           string             `json:"payer_email"`
	ExternalReference    string             `json:"external_reference"`
	Metadata             string             `json:"metadata"`
	Uuid                 string             `json:"uuid"`
	PdfUrl               string             `json:"pdf_url"`
	CreatedAt            string             `json:"created_at"`
	BaseUrl              string             `json:"base_url"`
	BackUrlSuccess       string             `json:"back_url_success"`
	BackUrlPending       string             `json:"back_url_pending"`
	BackUrlRejected      string             `json:"back_url_rejected"`
	IncludedChannels     []string           `json:"included_channels"`
	IncludedInstallments string             `json:"included_installments"`
	Items                string             `json:"items"`
	Preference           ResponsePreference `json:"preferences"`
}

type ResponsePreference struct {
	Client         string `json:"client"`
	MainColor      string `json:"mainColor"`
	SecondaryColor string `json:"secondaryColor"`
	Logo           string `json:"logo"`
}
