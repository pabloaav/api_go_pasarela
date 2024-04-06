package prismadtos

type PaymentsOfflineResponse struct {
	ID                      int64         `json:"id"`
	SiteTransactionID       string        `json:"site_transaction_id"`
	Token                   string        `json:"token"`
	PaymentMethodID         int64         `json:"payment_method_id"`
	Amount                  int64         `json:"amount"`
	Currency                string        `json:"currency"`
	Email                   string        `json:"email"`
	Status                  string        `json:"status"`
	StatusDetails           StatusDetails `json:"status_details"`
	Date                    string        `json:"date"`
	InvoiceExpiration       string        `json:"invoice_expiration"`
	SecondInvoiceExpiration string        `json:"second_invoice_expiration"`
	Surcharge               int64         `json:"surcharge"`
	Client                  string        `json:"client"`
	Barcode                 string        `json:"barcode"`
}
