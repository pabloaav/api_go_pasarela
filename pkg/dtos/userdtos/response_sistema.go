package userdtos

type ResponseSistema struct {
	Id           uint
	Sistema      string
	Activo       bool
	DiasPassword uint64
}