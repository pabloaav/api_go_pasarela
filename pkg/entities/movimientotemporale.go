package entities

import (
	"context"
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"gorm.io/gorm"
)

type Movimientotemporale struct {
	gorm.Model
	CuentasId           uint64                         `json:"cuentas_id"`
	PagointentosId      uint64                         `json:"pagointentos_id"`
	Tipo                EnumTipoMovimiento             `json:"tipoMovimiento"`
	Monto               Monto                          `json:"monto"`
	MotivoBaja          string                         `json:"motivo_baja"`
	Pagointentos        *Pagointento                   `json:"pagointentos" gorm:"foreignKey:pagointentos_id"`
	Cuenta              *Cuenta                        `json:"cuenta" gorm:"foreignKey:cuentas_id"`
	Reversion           bool                           `json:"Reversion"`
	Enobservacion       bool                           `json:"enobservacion"`
	Movimientocomisions []Movimientocomisionetemporale `json:"movimiento_comisiones" gorm:"foreignKey:MovimientotemporalesID"`
	Movimientoimpuestos []Movimientoimpuestotemporale  `json:"movimiento_impuestos" gorm:"foreignKey:MovimientotemporalesID"`
}

func (m *Movimientotemporale) IsValid() error {
	if m.CuentasId == 0 {
		return errors.New(ERROR_CUENTA)
	}
	if m.PagointentosId == 0 {
		return errors.New(ERROR_PAGO)
	}
	err := m.Tipo.IsValid()
	if err != nil {
		return err
	}
	// if m.Monto <= 0 {
	// 	return errors.New(ERROR_MONTO)
	// }
	return nil
}

func (m *Movimientotemporale) AddDebito(cuentaId uint64, pagoIntentosId uint64, monto Monto) error {
	m.CuentasId = cuentaId
	m.PagointentosId = pagoIntentosId
	m.Monto = monto
	m.Tipo = Debito

	err := m.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (m *Movimientotemporale) AddCredito(cuentaId uint64, pagoIntentosId uint64, monto Monto) error {
	m.CuentasId = cuentaId
	m.PagointentosId = pagoIntentosId
	m.Monto = monto
	m.Tipo = Credito
	err := m.IsValid()
	if err != nil {
		return err
	}
	return nil
}

type EnumTipoMovimientoTemporal string

const (
	Debito1  EnumTipoMovimientoTemporal = "D"
	Credito2 EnumTipoMovimientoTemporal = "C"
)

func (e EnumTipoMovimientoTemporal) IsValid() error {
	switch e {
	case Debito1, Credito2:
		return nil
	}
	return errors.New(ERROR_ENUM_TIPO_MOVIMIENTO)
}

const ERROR_ENUM_TIPO_MOVIMIENTO_TEMPORAL = "tipo EnumTipoMovimiento con formato inválido"
const ERROR_CUENTA_TEMPORAL = "el id de la cuenta es inválido"
const ERROR_PAGO_TEMPORAL = "el id del pago es inválido"
const ERROR_MONTO_TEMPORAL = "el monto informado es inválido"

func (ct *Movimientotemporale) AfterSave(tx *gorm.DB) (err error) {
	var audit Auditoria
	ctxValue := tx.Statement.Context.Value(AuditUserKey{})
	if ctxValue == nil {
		logs.Error("no hay datos de usuario para la transacción indicada")
	}
	audit = ctxValue.(Auditoria)
	stmt := tx.Statement
	str := tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	audit.Fila = ct.ID
	audit.Query = str
	audit.Tabla = "movimiento"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}
