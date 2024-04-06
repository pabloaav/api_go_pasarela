package filtros

type MovimientosTransferenciasResponse struct {
	Id                 uint   `json:"id"`
	ReferenciaBancaria string `json:"referencia_bancaria"` // Es la referencia que nos envia apilink luego de realizar la transferencia
	Match              int    `json:"match"`
	BancoExternalId    int    `json:"banco_external_id"`
}
