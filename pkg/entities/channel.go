package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	Channel          string            `json:"channel"`
	Nombre           string            `json:"nombre"`
	CodigoBcra       int32             `json:"codigo_bcra"`
	Mediopagos       []Mediopago       `json:"mediopagos" gorm:"foreignKey:ChannelsID"`
	Channelaranceles []Channelarancele `gorm:"foreignKey:ChannelsId"`
}

func (ct *Channel) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "channels"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
