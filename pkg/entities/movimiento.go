package entities

import (
	"context"
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"gorm.io/gorm"
)

type Movimiento struct {
	gorm.Model
	CuentasId               uint64                 `json:"cuentas_id"`
	PagointentosId          uint64                 `json:"pagointentos_id"`
	Tipo                    EnumTipoMovimiento     `json:"tipoMovimiento"`
	Monto                   Monto                  `json:"monto"`
	MotivoBaja              string                 `json:"motivo_baja"`
	Pagointentos            *Pagointento           `json:"pagointentos" gorm:"foreignKey:pagointentos_id"`
	Cuenta                  *Cuenta                `json:"cuenta" gorm:"foreignKey:cuentas_id"`
	Reversion               bool                   `json:"Reversion"`
	Enobservacion           bool                   `json:"enobservacion"`
	Movimientotransferencia []Transferencia        `json:"movimiento_transferencia" gorm:"foreignKey:MovimientosID"`
	Movimientocomisions     []Movimientocomisiones `json:"movimiento_comisiones" gorm:"foreignKey:MovimientosID"`
	Movimientoimpuestos     []Movimientoimpuestos  `json:"movimiento_impuestos" gorm:"foreignKey:MovimientosID"`
	Movimientolotes         []Movimientolotes      `json:"movimiento_lotes" gorm:"foreignKey:MovimientosID"`
	Movimientoretencions    []MovimientoRetencion  `json:"movimiento_retencions" gorm:"foreignKey:MovimientoId"`
	Retenciones             []Retencion            `gorm:"many2many:movimiento_retencions;"`
	Movimientosubcuentas    []Movimientosubcuentas `json:"movimiento_subcuenta" gorm:"foreignKey:MovimientosID"`
	// Movimientoliquidaciones []Movimientoliquidaciones `json:"movimientoliquidaciones" gorm:"foreignKey:MovimientosID"`
}

func (m *Movimiento) IsValid() error {
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

func (m *Movimiento) AddDebito(cuentaId uint64, pagoIntentosId uint64, monto Monto) error {
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

func (m *Movimiento) AddCredito(cuentaId uint64, pagoIntentosId uint64, monto Monto) error {
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

type EnumTipoMovimiento string

const (
	Debito  EnumTipoMovimiento = "D"
	Credito EnumTipoMovimiento = "C"
)

func (e EnumTipoMovimiento) IsValid() error {
	switch e {
	case Debito, Credito:
		return nil
	}
	return errors.New(ERROR_ENUM_TIPO_MOVIMIENTO)
}

const ERROR_ENUM_TIPO_MOVIMIENTO = "tipo EnumTipoMovimiento con formato inválido"
const ERROR_CUENTA = "el id de la cuenta es inválido"
const ERROR_PAGO = "el id del pago es inválido"
const ERROR_MONTO = "el monto informado es inválido"

func (ct *Movimiento) AfterSave(tx *gorm.DB) (err error) {
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
