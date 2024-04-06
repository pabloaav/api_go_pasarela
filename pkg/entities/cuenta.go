package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Cuenta struct {
	gorm.Model
	ClientesID              int64             `json:"clientes_id"`
	RubrosID                uint              `json:"rubros_id"`
	Cuenta                  string            `json:"cuenta"`
	Cbu                     string            `json:"cbu"`
	Cvu                     string            `json:"cvu"`
	Apikey                  string            `json:"apikey"`
	DiasRetiroAutomatico    int64             `json:"dias_retiro_automatico"`
	TransferenciaAutomatica bool              `json:"transferencia_automatica"`
	SplitCuentas            bool              `json:"split_cuentas"`
	Pagotipos               *[]Pagotipo       `json:"pagotipos" gorm:"foreignKey:CuentasID"`
	Cliente                 *Cliente          `json:"cliente" gorm:"foreignKey:clientes_id"`
	Rubro                   *Rubro            `json:"rubro" gorm:"foreignKey:rubros_id"`
	Cuentacomisions         *[]Cuentacomision `json:"cuentacomisiones" gorm:"foreignKey:CuentasID"`
	Subcuentas              []Subcuenta       `json:"subcuentas" gorm:"foreignKey:CuentasID"`
}

// TableName sobreescribe el nombre de la tabla
func (Cuenta) TableName() string {
	return "cuentas"
}

func (ct *Cuenta) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "cuenta"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
