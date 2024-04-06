package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Pagoestado struct {
	gorm.Model
	Estado EnumPagoEstado `json:"estado"`
	Final  bool           `json:"final"`
	Nombre string         `json:"nombre"`
}

type EnumPagoEstado string

const (
	Pending    EnumPagoEstado = "Pending"
	Processing EnumPagoEstado = "Processing"
	Rejected   EnumPagoEstado = "Rejected"
	Paid       EnumPagoEstado = "Paid"
	Reverted   EnumPagoEstado = "Reverted"
	Expired    EnumPagoEstado = "Expired"
	Accredited EnumPagoEstado = "Accredited"
)

func (e EnumPagoEstado) IsValid() error {
	switch e {
	case Pending, Processing, Rejected, Paid, Reverted, Expired, Accredited:
		return nil
	}
	return errors.New("tipo EnumPagoEstado con formato inválido")
}

func (p *Pagoestado) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "pagoestados"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
