package administraciondtos

type RequestCLRapipagoExternalId struct {
	BancoId    int64 `json:"banco_id"`
	RapipagoId int64 `json:"rapipago_id"`
}
