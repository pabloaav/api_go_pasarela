package entities

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Pago struct {
	gorm.Model
	PagostipoID       int64         `json:"pagostipo_id"`
	PagoestadosID     int64         `json:"pagoestados_id"`
	Description       string        `json:"description"`
	FirstDueDate      time.Time     `json:"first_due_date"`
	FirstTotal        Monto         `json:"first_total"`
	SecondDueDate     time.Time     `json:"second_due_date,omitempty"`
	SecondTotal       Monto         `json:"second_total,omitempty"`
	PayerName         string        `json:"payer_name"`
	PayerEmail        string        `json:"payer_email"`
	ExternalReference string        `json:"external_reference"`
	Metadata          string        `json:"metadata"`
	Uuid              string        `json:"uuid"`
	PdfUrl            string        `json:"pdf_url"`
	Notificado        bool          `json:"notificado"`
	PagosTipo         Pagotipo      `json:"pagotipo" gorm:"foreignKey:pagostipo_id"`
	PagoEstados       Pagoestado    `json:"pagoestados" gorm:"foreignKey:pagoestados_id"`
	PagoIntentos      []Pagointento `json:"pago_intentos" gorm:"foreignKey:pagos_id"`
	Pagoitems         []Pagoitems   `json:"pagoitems" gorm:"foreignKey:PagosID"`
	Pagolotes         []Pagolotes   `json:"pagos_lotes" gorm:"foreignKey:PagosID"`
	Expiration        int64         `json:"expiration"`
	IntentoFallido	  int64          `json:"intento_fallido"`
}

func (p *Pago) AfterSave(tx *gorm.DB) (err error) {
	var audit Auditoria
	ctxValue := tx.Statement.Context.Value(AuditUserKey{})
	if ctxValue == nil {
		return errors.New("no hay datos de cuenta para la transacción indicada")
	}
	audit = ctxValue.(Auditoria)
	stmt := tx.Statement
	str := tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	audit.Fila = p.ID
	audit.Query = str
	audit.Tabla = "pagos"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
