package entities

import "gorm.io/gorm"

type Auditoria struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	CuentaID  uint   `json:"cuenta_id"`
	IP        string `json:"ip"`
	Tabla     string `json:"tabla"`
	Fila      uint   `json:"fila"`
	Operacion string `json:"operacion"`
	Query     string `json:"query"`
	Resultado string `json:"resultado"`
	Origen    string `json:"origen"`
}

// TableName sobreescribe el nombre de la tabla
func (Auditoria) TableName() string {
	return "auditorias"
}

type AuditUserKey struct{}
