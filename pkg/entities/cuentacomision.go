package entities

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Cuentacomision struct {
	gorm.Model
	CuentasID          uint    `json:"cuentas_id"`
	ChannelsId         uint    `json:"channels_id"`
	ChannelarancelesId uint    `json:"channelaranceles_id"`
	Cuentacomision     string  `json:"cuentacomision"`
	Comision           float64 `json:"comision"`
	Mediopagoid        uint    `json:"mediopagoid"`
	//	Iva                  float64                `json:"iva"`
	VigenciaDesde        *time.Time             `json:"vigencia_desde"`
	Cuenta               Cuenta                 `json:"cuenta" gorm:"foreignKey:cuentas_id"`
	Movimientocomisiones []Movimientocomisiones `json:"movimiento_comisiones" gorm:"foreignKey:cuentacomisions_id"`
	Channel              Channel                `json:"channel" gorm:"foreignKey:channels_id"`
	Importeminimo        float64                `json:"importeminimo"`
	Importemaximo        float64                `json:"importemaximo"`
	Pagocuota            bool                   `json:"pagocuota"`
	ChannelArancel       Channelarancele        `json:"channelarancel" gorm:"foreignKey:channelaranceles_id"`
}

func (ct *Cuentacomision) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "Cuentacomisions"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
