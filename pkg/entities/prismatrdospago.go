package entities

import (
	"time"

	"gorm.io/gorm"
)

type Prismatrdospago struct {
	gorm.Model
	PrismatrcuatropagosId      uint `gorm:"column:prismatrcuatropagos_id"`
	FechaPresentacion          time.Time
	TipoRegistro               string
	Moneda                     string
	LiquidacionNro             string
	FechaPago                  time.Time
	LiquidacionTipo            string
	ImporteBruto               Monto
	SignoImporteBruto          string
	ImporteArancel             Monto
	SignoImporteArancel        string
	ImporteNeto                Monto
	SignoImporteNeto           string
	RetencionEspecialSobreIibb Monto
	SignoRetencionEspecial     string
	RetencionIvaEspecial       Monto
	SignoRetencionIvaEspecial  string
	PercepcionIngresoBruto     Monto
	SignoPercepcionIb          string
	RetencionIvaD1             Monto
	SignoRetencionIva_d1       string
	CostoCuotaEmitida          Monto
	SignoCostoCuotaEmitida     string
	RetencionIvaCuota          Monto
	SignoRetencionIvaCuota     string
	RetencionIva               Monto
	SignoRetencionIva          string
	RetencionGanacias          Monto
	SignoRetencionGanacias     string
	RetencionIngresoBruto      Monto
	SignoRetencionIngresoBruto string
	Match                      int
	CierreLotes                []Prismacierrelote `gorm:"foreignkey:prismatrdospagos_id"`
}
