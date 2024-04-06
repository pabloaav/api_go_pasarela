package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type TransferenciaResponsePaginado struct {
	Transferencias []TransferenciaResponse `json:"data"`
	Meta           dtos.Meta               `json:"meta"`
}

type TransferenciaResponse struct {
	Id                              uint      `json:"id"`
	MovimientosID                   uint64    `json:"movimientos_id"`
	UserId                          uint64    `json:"user_id"`
	ReferenciaBancaria              string    `json:"referencia_bancaria"` // Es la referencia que nos envia apilink luego de realizar la transferencia
	Uuid                            string    `json:"uuid"`                //Es la referencia que enviamos nosotros a apilink
	CbuDestino                      string    `json:"cbu_destino"`
	CbuOrigen                       string    `json:"cbu_origen"`
	Match                           int       `json:"match"`
	BancoExternalId                 int       `json:"banco_external_id"`
	Fecha                           time.Time `json:"fecha"`
	ReferenciaBanco                 string    `json:"referencia_banco"`
	ListaIdsTransferenciasAgrupadas []uint    `json:"lista_ids_transferencias_agrupadas"` // grupo de ids de transferencias que corresponden a un mismo codigo de referencia banco
}

func (t *TransferenciaResponse) New(transferencia entities.Transferencia) {
	t.Id = transferencia.ID
	t.MovimientosID = transferencia.MovimientosID
	t.UserId = transferencia.UserId
	t.ReferenciaBancaria = transferencia.ReferenciaBancaria
	t.Uuid = transferencia.Uuid
	t.CbuDestino = transferencia.CbuDestino
	t.CbuOrigen = transferencia.CbuOrigen
	t.Fecha = transferencia.CreatedAt
	t.Match = transferencia.Match
	t.BancoExternalId = transferencia.BancoExternalId
	t.ReferenciaBanco = transferencia.ReferenciaBanco
}

type TransferenciaRespons struct {
	Transferencias      []TransferenciaResponseAgrupada `json:"data"`
	TotalTransferencias string                          `json:"total_transferencias"`
	Meta                dtos.Meta                       `json:"meta"`
}

type TransferenciaResponseAgrupada struct {
	ReferenciaBancaria         string         `json:"referencia_bancaria"`
	CbuDestino                 string         `json:"cbu_destino"`
	CbuOrigen                  string         `json:"cbu_origen"`
	Fecha                      time.Time      `json:"fecha"`
	Monto                      entities.Monto `json:"monto"`
	MovimientoReponse          []MovimientoReponse
	ReferenciaBanco            string `json:"referencia_banco"`             // atributo para comparar con los movimientos del proyecto banco
	IdsTransferenciasAgrupadas []uint `json:"ids_transferencias_agrupadas"` // acumula los ids de los registros transferencias que corresponden a una misma transaccion o referencia banco
}

type MovimientoReponse struct {
	Concepto       string         `json:"concepto"`
	ReferenciaPago string         `json:"referenciaPago"`
	CanalPago      string         `json:"canal_pago"`
	MontoMov       entities.Monto `json:"monto_mov"`
}
