package prismadtos

type OfflineTokenResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
	DateDue     string `json:"date_due"`
	// json que contiene nombre tipo-DNI y nro
	Customer DataCustomer `json:"customer"`
}
