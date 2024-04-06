package pagoofflinedtos

import "time"

type OffLineRequestResponse struct {
	CodigoEmpresa     string    `json:"codigo_empresa"`
	NumeroCliente     string    `json:"numero_cliente"`
	NumeroComprobante string    `json:"numero_comprobante"`
	Importe           int64     `json:"importe"`
	FechaPrimerVto    time.Time `json:"fecha_primer_vto"`
	ImporteRecargo    int64     `json:"importe_recargo"`
	FechaSegundoVto   time.Time `json:"fecha_segundo_vto"`
}
