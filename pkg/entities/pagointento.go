package entities

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Pagointento struct {
	gorm.Model
	PagosID              int64                 `json:"pagos_id"`
	MediopagosID         int64                 `json:"mediopagos_id"`
	InstallmentdetailsID int64                 `json:"installmentdetails_id"`
	ExternalID           string                `json:"external_id"`
	PaidAt               time.Time             `json:"paid_at"`
	ReportAt             time.Time             `json:"report_at"`
	IsAvailable          bool                  `json:"is_available"`
	Amount               Monto                 `json:"amount"`
	Valorcupon           Monto                 `json:"valorcupon"`
	StateComment         string                `json:"state_comment"`
	Barcode              string                `json:"barcode"`
	BarcodeUrl           string                `json:"barcode_url"`
	AvailableAt          time.Time             `json:"available_at"`
	RevertedAt           time.Time             `json:"reverted_at"`
	HolderName           string                `json:"holder_name"`
	HolderEmail          string                `json:"holder_email"`
	HolderCbu            string                `json:"holder_cbu"`
	HolderType           string                `json:"holder_type"`
	HolderNumber         string                `json:"holder_number"`
	TicketNumber         string                `json:"ticket_number"`
	AuthorizationCode    string                `json:"authorization_code"`
	CardLastFourDigits   string                `json:"card_last_four_digits"`
	TransactionID        string                `json:"transaction_id"`
	SiteID               int64                 `json:"site_id"`
	Calculado            bool                  `json:"calculado"`
	Ip 					 string 			   `json:"ip"`  
	Mediopagos           Mediopago             `json:"mediopago" gorm:"foreignKey:mediopagos_id"`
	Pago                 Pago                  `json:"pago" gorm:"foreignKey:pagos_id"`
	Movimientos          []Movimiento          `json:"movimiento" gorm:"foreignKey:pagointentos_id"`
	Movimientotemporale  []Movimientotemporale `json:"movimientotemporales" gorm:"foreignKey:pagointentos_id"`
	Installmentdetail    Installmentdetail     `json:"Installmentdetail" gorm:"foreignKey:InstallmentdetailsID"`
}

func (p *Pagointento) AfterSave(tx *gorm.DB) (err error) {
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
	audit.Tabla = "pagointentos"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
