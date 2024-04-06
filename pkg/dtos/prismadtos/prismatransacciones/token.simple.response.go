package prismadtos

type PagoToken struct {
	Id                string     `json:"id,omitempty"`
	ValidationResulto bool       `josn:"validation_resulto,omitempty"`
	Status            string     `json:"status,omitempty"`
	DataUsed          string     `json:"date_used,omitempty"`
	CardNumberLength  uint64     `json:"card_number_length,omitempty"`
	Bin               string     `json:"bin,omitempty"`
	DateCreate        string     `json:"date_created,omitempty"`
	LastFourDigits    string     `json:"last_four_digits,omitempty"`
	SecurityCodeLeng  uint64     `json:"security_code_leng,omitempty"`
	ExpirationMonth   uint64     `json:"expiration_month,omitempty"`
	ExpirationYear    uint64     `json:"expiration_year,omitempty"`
	DateLastUpdated   string     `json:"date_last_updated,omitempty"`
	DateDue           string     `json:"date_due,omitempty"`
	CardHolder        CardHolder `json:"cardholder"`
}

type CardHolder struct {
	Identification Identification `json:"identification"`
	Name           string         `json:"Name,omitempty"`
}

type Identification struct {
	TypeDni   string `json:"type"`
	NumberDni string `json:"number"`
}
