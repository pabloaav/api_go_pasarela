package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Cliente struct {
	gorm.Model
	IvaID             int64               `json:"iva_id"`
	IibbID            int64               `json:"iibb_id"`
	Cliente           string              `json:"cliente"`
	Razonsocial       string              `json:"razonsocial"`
	Nombrefantasia    string              `json:"nombrefantasia"`
	Email             string              `json:"email"`
	Emailcontacto     string              `json:"emailcontacto"`
	Personeria        string              `json:"personeria"`
	RetiroAutomatico  bool                `json:"retiro_automatico"`
	Cuit              string              `json:"cuit"`
	ReporteBatch      bool                `json:"reporte_batch"`
	NombreReporte     string              `json:"nombre_reporte"`
	OrdenDiaria       bool                `json:"orden_diaria"`
	Iva               *Impuesto           `json:"iva" gorm:"foreignKey:iva_id"`
	Iibb              *Impuesto           `json:"iibb" gorm:"foreignKey:iibb_id"`
	Cuentas           *[]Cuenta           `json:"cuentas" gorm:"foreignKey:ClientesID"`
	Clienteusers      *[]Clienteuser      `json:"cliente_users" gorm:"foreignKey:ClientesId"`
	Contactosreportes *[]Contactosreporte `json:"contactos_reportes" gorm:"foreignKey:ClientesID"`
	SujetoRetencion   bool                `json:"sujeto_retencion"`
	Formulario_8125   bool                `json:"formulario_8125"`
	Envio  		Envio               				`gorm:"foreignKey:ClientesId"`
}

func (ct *Cliente) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "cliente"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
