package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagos struct {
	Pagos           []ResponsePago `json:"data"`
	SaldoPendiente  entities.Monto `json:"saldo_pendiente"`
	SaldoDisponible entities.Monto `json:"saldo_disponible"`
	Meta            dtos.Meta      `json:"meta"`
}

type ResponsePago struct {
	Identificador       uint           `json:"identificador"`
	Cuenta              string         `json:"cuenta"`
	Pagotipo            string         `json:"pagotipo"`
	Fecha               time.Time      `json:"fecha"`
	ExternalReference   string         `json:"external_reference"`
	PayerName           string         `json:"payer_name"`
	Estado              string         `json:"estado"`
	NombreEstado        string         `json:"nombre_estado"`
	Amount              entities.Monto `json:"amount"`
	FechaPago           time.Time      `json:"fecha_pago"`
	Channel             string         `json:"channel"`
	NombreChannel       string         `json:"nombre_channel"`
	UltimoPagoIntentoId uint64         `json:"ultimo_pago_intento_id"`
	TransferenciaId     uint64         `json:"transferencia_id"`
	FechaTransferencia  string         `json:"fecha_transferencia"`
	ReferenciaBancaria  string         `json:"referencia_bancaria"`
	// PagoItems           []PagoItems    `json:"pago_items"`
}

type PagoItems struct {
	Descripcion   string
	Identificador string
	Cantidad      int64
	Monto         float64
}

func (pagoDTO *ResponsePago) FromPago(pago entities.Pago) {
	pagoDTO.Identificador = pago.ID
	pagoDTO.Fecha = pago.CreatedAt
	pagoDTO.ExternalReference = pago.ExternalReference
	pagoDTO.PayerName = pago.PayerName

	if pago.PagoEstados.ID > 0 {
		pagoDTO.Estado = string(pago.PagoEstados.Estado)
		pagoDTO.NombreEstado = pago.PagoEstados.Nombre

	}

	if pago.PagosTipo.ID > 0 {
		pagoDTO.Pagotipo = pago.PagosTipo.Pagotipo
		if pago.PagosTipo.Cuenta.ID > 0 {
			pagoDTO.Cuenta = pago.PagosTipo.Cuenta.Cuenta
		}
	}

	if len(pago.PagoIntentos) > 0 {
		last := len(pago.PagoIntentos) - 1
		pagoIntento := pago.PagoIntentos[last] // obteniendo ultimo pago intento
		pagoDTO.Amount = pagoIntento.Amount
		pagoDTO.FechaPago = pagoIntento.PaidAt
		if pagoIntento.Mediopagos.ID > 0 && pagoIntento.Mediopagos.Channel.ID > 0 {
			pagoDTO.Channel = pagoIntento.Mediopagos.Channel.Channel
			pagoDTO.NombreChannel = pagoIntento.Mediopagos.Channel.Nombre
		}
		pagoDTO.UltimoPagoIntentoId = uint64(pagoIntento.ID)
	}
}
