package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Pagotipo struct {
	gorm.Model
	CuentasID                uint                 `json:"cuentas_id"`
	Pagotipo                 string               `json:"pagotipo"`
	BackUrlSuccess           string               `json:"back_url_success"`
	BackUrlPending           string               `json:"back_url_pending"`
	BackUrlRejected          string               `json:"back_url_rejected"`
	BackUrlNotificacionPagos string               `json:"back_url_notificacion_pagos"`
	Cuenta                   Cuenta               `json:"cuenta" gorm:"foreignKey:cuentas_id"`
	Pagos                    []Pago               `json:"pagos" gorm:"foreignKey:PagostipoID"`
	Pagotipochannel          []Pagotipochannel    `json:"canales" gorm:"foreignKey:PagotiposId"`
	Pagotipoinstallment      []Pagotipointallment `json:"cuotas" gorm:"foreignKey:PagotiposId"`
	ApikeyExterno            string               `json:"apikey_externo"`
}

func (ct *Pagotipo) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "pagotipos"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
