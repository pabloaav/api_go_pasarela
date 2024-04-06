package dtos

// PagoResponse respuesta utilizada en la solicitud de un pago (el primer paso)
type PagoResponse struct {
	ID                int64               `json:"id"`
	Estado            string              `json:"estado"`
	Description       string              `json:"description"`
	FirstDueDate      string              `json:"first_due_date"`
	FirstTotal        float64             `json:"first_total"`
	SecondDueDate     string              `json:"second_due_date,omitempty"`
	SecondTotal       float64             `json:"second_total,omitempty"`
	PayerName         string              `json:"payer_name"`
	PayerEmail        string              `json:"payer_email"`
	ExternalReference string              `json:"external_reference"`
	Metadata          string              `json:"metadata"`
	Uuid              string              `json:"uuid"`
	CheckoutUrl       string              `json:"checkout_url"`
	PdfUrl            string              `json:"pdf_url"`
	CreatedAt         string              `json:"created_at"`
	Expiration        int64               `json:"expiration"`
	Items             []PagoResponseItems `json:"items"`
}

type PagoResponseItems struct {
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Identifier  string  `json:"identifier"`
}
