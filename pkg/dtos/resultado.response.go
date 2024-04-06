package dtos

type ResultadoResponseWebHook struct {
	Url               string              `json:"url"`
	ApikeyExterno     string              `json:"apikey_externo"`
	ResultadoResponse []ResultadoResponse `json:"resultado_response"`
}
type ResultadoResponse struct {
	ID                int64   `json:"id"`
	Estado            string  `json:"estado"`
	EstadoPago        string  `json:"estado_pago"`
	Exito             bool    `json:"exito"`
	Uuid              string  `json:"uuid"`
	Channel           string  `json:"channel"`
	Description       string  `json:"description"`
	FirstDueDate      string  `json:"first_due_date"`
	FirstTotal        float64 `json:"first_total"`
	SecondDueDate     string  `json:"second_due_date,omitempty"`
	SecondTotal       float64 `json:"second_total,omitempty"`
	PayerName         string  `json:"payer_name"`
	PayerEmail        string  `json:"payer_email"`
	ExternalReference string  `json:"external_reference"`
	Metadata          string  `json:"metadata"`
	Barcode           string  `json:"barcode"`
	BarcodeUrl        string  `json:"barcode_url"`
	PdfUrl            string  `json:"pdf_url"`
	CreatedAt         string  `json:"created_at"`
	ImportePagado     float64 `json:"importe_pagado"`
	Mensaje           string  `json:"mensaje"`
}
