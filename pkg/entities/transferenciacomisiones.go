package entities

import (
	"time"

	"gorm.io/gorm"
)

type Transferenciacomisiones struct {
	gorm.Model
	MovimientosID              uint64     `json:"movimientos_id"`
	UserId                     uint64     `json:"user_id"`
	Referencia                 string     `json:"referencia"`          // es la referecnia de la transferencia al enviar la peticion
	ReferenciaBancaria         string     `json:"referencia_bancaria"` // Es la referencia que nos envia apilink luego de realizar la transferencia
	Uuid                       string     `json:"uuid"`
	CbuDestino                 string     `json:"cbu_destino"`
	CbuOrigen                  string     `json:"cbu_origen"`
	FechaOperacion             *time.Time `json:"fecha_operacion"`              // Es la fecha que nos envia apilink luego de realizar la transferencia
	NumeroConciliacionBancaria string     `json:"numero_conciliacion_bancaria"` // Es el numero de conciliacion que nos envia apilink luego de realizar la transferencia
	Reversion                  bool       `json:"reversion"`
	Movimiento                 Movimiento `json:"movimiento" gorm:"foreignKey:movimientos_id"`
}
