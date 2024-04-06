package prismadtos

type DataCustomer struct {
	Identification IdentificationCustomer `json:"identification"`
	Name           string                 `json:"name"`
}

type IdentificationCustomer struct {
	Type   EnumTipoDocumento `json:"type"`
	Number string            `json:"number"`
}
