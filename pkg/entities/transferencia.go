package entities

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Transferencia struct {
	gorm.Model
	MovimientosID              uint64     `json:"movimientos_id"`
	UserId                     uint64     `json:"user_id"`
	Referencia                 string     `json:"referencia"`          // es la referecnia de la transferencia al enviar la peticion
	ReferenciaBancaria         string     `json:"referencia_bancaria"` // Es la referencia que nos envia apilink luego de realizar la transferencia
	Uuid                       string     `json:"uuid"`
	CbuDestino                 string     `json:"cbu_destino"`
	CbuOrigen                  string     `json:"cbu_origen"`
	Match                      int        `json:"match"`             // utilizado parra conciliacion con banco
	BancoExternalId            int        `json:"banco_external_id"` // utilizado parra conciliacion con banco
	FechaOperacion             *time.Time `json:"fecha_operacion"`   // Es la fecha que nos envia apilink luego de realizar la transferencia
	Reversion                  bool       `json:"reversion"`
	NumeroConciliacionBancaria string     `json:"numero_conciliacion_bancaria"` // Es el numero de conciliacion que nos envia apilink luego de realizar la transferencia
	ReferenciaBanco            string     `json:"referencia_banco"`             // campo que concatena la fecha y numero de conciliacion para que sea unica
	Movimiento                 Movimiento `json:"movimiento" gorm:"foreignKey:movimientos_id"`
}

// TableName sobreescribe el nombre de la tabla
func (Transferencia) TableName() string {
	return "transferencias"
}

func (ct *Transferencia) AfterSave(tx *gorm.DB) (err error) {
	var audit Auditoria
	ctxValue := tx.Statement.Context.Value(AuditUserKey{})
	if ctxValue == nil {
		return errors.New("no hay datos de usuario para la transacción indicada")
	}
	audit = ctxValue.(Auditoria)
	stmt := tx.Statement
	str := tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	audit.Fila = ct.ID
	audit.Query = str
	audit.Tabla = "movimiento"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
