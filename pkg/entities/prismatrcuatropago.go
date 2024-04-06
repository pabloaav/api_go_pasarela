package entities

import (
	"time"

	"gorm.io/gorm"
)

type Prismatrcuatropago struct {
	gorm.Model
	Empresa                       string
	FechaPresentacion             time.Time
	TipoRegistro                  string
	Moneda                        string
	ComercioNro                   string
	EstablecimientoNro            string
	LiquidacionNro                string
	FechaPago                     time.Time
	LiquidacionTipo               string
	CasaCuentaAcreditacion        string
	TipoCuentaAcreditacion        string
	NroCuentaAcreditacion         string
	CostoFinacieroExentoIva       Monto
	SignoCostoFinacieroExtIva     string
	RetencionPorLey_25063         Monto
	SignoRetencionPorLey          string
	AlicuotaRetencionIb           Monto
	CargoTransferenciaBancaria    Monto
	SignoCargoTransBancaria       string
	IvaSobreCargoTransBancaria    Monto
	SignoIvaCargoTransBancaria    string
	ImpuestoDbCr                  Monto
	SignoImpuestoDbCr             string
	CostoFinancieroNoReducIva     Monto
	SignoCostoFinancieroReducIva  string
	RetencionIvaRg_3130           Monto
	SignoRetencionIvaRg_3130      string
	JurisdiccionRetencionIb       string
	CargoAdicionalPlanCuotas      Monto
	SignoCargoAdicional           string
	IvaCargoAdicionalPlanCuota    Monto
	SignoIvaCargoAdicional        string
	CargoAdicionalMovipag         Monto
	SignoCargiMovipag             string
	IvaCargoAdicionalMovipag      Monto
	SignoIvaCargoAdicionalMovipag string
	RetencionSello                Monto
	SignoRetencionSello           string
	ProvinciaRetencionSello       string
	Match                         int
	Pagostrdos                    []Prismatrdospago `gorm:"foreignkey:PrismatrcuatropagosId"`
}
