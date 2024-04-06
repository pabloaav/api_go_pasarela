package linkdebin

type RequestDebines struct {
	Debines          []string `json:"debines"`
	Match            bool     `json:"match"`
	BancoExternalId  bool     `json:"banco_external_id"`
	Pagoinformado    bool     `json:"pagoinformado"`
	CargarPagoEstado bool     `json:"cargar_pago_estado"`
}
