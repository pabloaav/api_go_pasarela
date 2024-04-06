package webhook

import (
	"time"

	monto "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type WebhookResponse struct {
	Url                      string                     `json:"url"`
	ApikeyExterno            string                     `json:"apikey_externo"`
	ResultadoResponseWebHook []ResultadoResponseWebHook `json:"resultado"`
}
type ResultadoResponseWebHook struct {
	Id                int64       `json:"id"`
	EstadoPago        string      `json:"estado_pago"`
	Exito             bool        `json:"exito"`
	Channel           string      `json:"channel"`
	Uuid              string      `json:"uuid"`
	Description       string      `json:"description"`
	FirstDueDate      time.Time   `json:"first_due_date"`
	FirstTotal        monto.Monto `json:"first_total"`
	SecondDueDate     time.Time   `json:"second_due_date,omitempty"`
	SecondTotal       monto.Monto `json:"second_total,omitempty"`
	PayerName         string      `json:"payer_name"`
	PayerEmail        string      `json:"payer_email"`
	ExternalReference string      `json:"external_reference"`
	Metadata          string      `json:"metadata"`
	Barcode           string      `json:"barcode"`
	BarcodeUrl        string      `json:"barcode_url"`
	PdfUrl            string      `json:"pdf_url"`
	CreatedAt         time.Time   `json:"created_at"`
	ImportePagado     float64     `json:"importe_pagado"` /* pago intento */
}
