package bancodtos

type ResponseMovimientosBanco struct {
	Id                        uint   `json:"id"`
	Subcuenta                 string `json:"subcuenta"`
	Referencia                string `json:"referencia"`
	DbCr                      string `json:"debito_credito"`
	Importe                   uint64 `json:"importe"`
	Fecha                     string `json:"fecha"`
	Hora                      string `json:"hora"`
	TipoMovimiento            int64  `json:"tipo_movimiento"`
	TipoMovimientoDescripcion string `json:"tipo_movimiento_descripcion"`
	DebinId                   string `json:"debin_id"`
	Seclink                   string `json:"seclink"`
	ReferenciaTransferencia   string `json:"referencia_transferencia"`
}
