package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestTransferenciaAutomatica struct {
	CuentaId      uint64
	Cuenta        string
	DatosClientes DatosClientes
	Request       RequestTransferenicaCliente
}
type RequestTransferenicaCliente struct {
	Transferencia         linktransferencia.RequestTransferenciaCreateLink `json:"transferencia,omitempty"`
	ListaMovimientosId    []uint64                                         `json:"lista_movimientos_id,omitempty"`
	ListaMovimientosIdNeg []uint64                                         `json:"lista_movimientos_id_neg,omitempty"`
}

type RequestTransferenciaMov struct {
	Transferencia              linktransferencia.RequestTransferenciaCreateLink `json:"transferencia"`
	Lista                      []uint64                                         `json:"lista"`
	ReferenciaBancaria         string                                           `json:"referencia_bancaria"`
	NumeroConciliacionBancaria string                                           `json:"numero_conciliacion_bancaria"`
	ListaMovimientosIdNeg      []uint64                                         `json:"lista_movimientos_id_neg"`
}

type ResponseTransferenciaAutomatica struct {
	CuentaId uint64         `json:"cuentaid"`
	Cuenta   string         `json:"cuenta"`
	Origen   string         `json:"origen"`
	Destino  string         `json:"destino"`
	Importe  entities.Monto `json:"importe"`
	Error    string         `json:"error"`
}

type DatosClientes struct {
	NombreCliente string
	EmailCliente  string
}
