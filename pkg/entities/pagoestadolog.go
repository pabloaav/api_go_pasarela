package entities

import (
	"context"

	"gorm.io/gorm"
)

type Pagoestadologs struct {
	gorm.Model
	PagosID       int64 `json:"pagos_id"`
	PagoestadosID int64 `json:"pagoestados_id"`
}

func (p *Pagoestadologs) AfterSave(tx *gorm.DB) (err error) {
	var audit Auditoria
	ctxValue := tx.Statement.Context.Value(AuditUserKey{})
	if ctxValue == nil {
		return nil
		// return errors.New("no hay datos de cuenta para la transacción indicada")
	}
	audit = ctxValue.(Auditoria)
	stmt := tx.Statement
	str := tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	audit.Fila = p.ID
	audit.Query = str
	audit.Tabla = "pagoestadologs"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
