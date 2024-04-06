package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type MovimientoAcumuladoResponsePaginado struct {
	Acumulados           []MovimientoPorCuentaAcumulado `json:"data"`
	MovimientosNegativos []MovimientosNegativos         `json:"movnegativos"`
	AcumuladosNegativos  entities.Monto                 `json:"acumuladoNegativo"`
	Meta                 dtos.Meta                      `json:"meta"`
}

type MovimientoPorCuentaResponsePaginado struct {
	Acumulados           []MovimientoPorCuentaResponse `json:"data"`
	MovimientosNegativos []MovimientosNegativos        `json:"movnegativos"`
	Meta                 dtos.Meta                     `json:"meta"`
}

type MovimientoPorCuentaAcumulado struct {
	Fecha                 string                        `json:"fecha"`
	FechaDisponibleRetiro string                        `json:"fecha_disponible"`
	Acumulado             entities.Monto                `json:"acumulado"`
	Movimientos           []MovimientoPorCuentaResponse `json:"movimientos"`
}
type MovimientoPorCuentaResponse struct {
	Id                   uint                          `json:"id"`                 // Id movimiento
	PagoIntentosId       uint                          `json:"pagointentos_id"`    // id pagointento
	Identificador        string                        `json:"identificador"`      // uuid es el identificador interno de telco
	Estado               entities.EnumPagoEstado       `json:"estado"`             // estado del pago (acreditado)
	Tipo                 string                        `json:"tipo"`               // tipo Debito o Credito
	PagoId               uint                          `json:"pago_id"`            // Id del pago para poder navegar al pago
	Pagotipo             string                        `json:"pago_tipo"`          // tipo del pago ej. sello
	ExternalReference    string                        `json:"external_reference"` // Es el numero de referencia que nos pas el cliente
	MedioPago            string                        `json:"medio_pago"`         // Medio de pago utilizado en la operación
	Channels             string                        `json:"channels"`
	FechaRendicion       string                        `json:"fecha_rendicion"`
	Monto                entities.Monto                `json:"monto"`       // monto es el importe neto que es enviado al cliente
	Montopagado          entities.Monto                `json:"montopagado"` // monto pagado el importe que pago el pagador y puede variar segun medio de pago
	Montosp              entities.Monto                `json:"montosp"`     // montosp es el importe enviado por medio de la solicitud de pago
	Revertido            bool                          `json:"revertido"`
	Enobservacion        bool                          `json:"enobservacion"`
	Comisiones           MovimientoComisionResponse    `json:"comision"`
	Impuestos            MovimientoImpuestoResponse    `json:"impuestos"`
	PagoCreated_at       time.Time                     `json:"pago_created_at"`
	MovimientoCreated_at time.Time                     `json:"movimiento_created_at"`
	Retenciones          RetencionesResponse           `json:"retenciones"`
	MovimientoSubcuentas []ResponseMovimientosubcuenta `json:"movimiento_subcuentas"`
}

// en el caso de que existan movimientos negativos
type MovimientosNegativos struct {
	Id                   uint           `json:"id"`                 // Id movimiento
	PagoIntentosId       uint           `json:"pagointentos_id"`    // id pagointento
	Tipo                 string         `json:"tipo"`               // tipo Debito o Credito
	Pagotipo             string         `json:"pago_tipo"`          // tipo del pago ej. sello
	ExternalReference    string         `json:"external_reference"` // Es el numero de referencia que nos pas el cliente
	MedioPago            string         `json:"medio_pago"`         // Medio de pago utilizado en la operación
	Monto                entities.Monto `json:"monto"`
	PagoCreated_at       time.Time      `json:"pago_created_at"`
	MovimientoCreated_at time.Time      `json:"movimiento_created_at"`
	Reversion            bool           `json:"reversion"`
}

type MovimientoComisionResponse struct {
	Total   entities.Monto            `json:"total"`
	Detalle []DetalleComisionResponse `json:"detalle"`
}

type MovimientoImpuestoResponse struct {
	Total   entities.Monto            `json:"total"`
	Detalle []DetalleImpuestoResponse `json:"detalle"`
}

type RetencionesResponse struct {
	MontoTotal, RetGanancias, RetIVA, RetIIBB entities.Monto
}

type DetalleComisionResponse struct {
	Nombre     string         `json:"nombre"`
	Monto      entities.Monto `json:"monto"`
	Porcentaje float64        `json:"porcentaje"`
}

type DetalleImpuestoResponse struct {
	Nombre     string         `json:"nombre"`
	Monto      entities.Monto `json:"monto"`
	Porcentaje float64        `json:"porcentaje"`
}
