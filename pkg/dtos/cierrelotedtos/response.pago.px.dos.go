package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagoPxDos struct {
	Id                         uint
	PrismatrcuatropagosId      uint
	FechaPresentacion          time.Time
	TipoRegistro               string
	Moneda                     string
	LiquidacionNro             string
	FechaPago                  time.Time
	LiquidacionTipo            string
	ImporteBruto               entities.Monto
	SignoImporteBruto          string
	ImporteArancel             entities.Monto
	SignoImporteArancel        string
	ImporteNeto                entities.Monto
	SignoImporteNeto           string
	RetencionEspecialSobreIibb entities.Monto
	SignoRetencionEspecial     string
	RetencionIvaEspecial       entities.Monto
	SignoRetencionIvaEspecial  string
	PercepcionIngresoBruto     entities.Monto
	SignoPercepcionIb          string
	RetencionIvaD1             entities.Monto
	SignoRetencionIva_d1       string
	CostoCuotaEmitida          entities.Monto
	SignoCostoCuotaEmitida     string
	RetencionIvaCuota          entities.Monto
	SignoRetencionIvaCuota     string
	RetencionIva               entities.Monto
	SignoRetencionIva          string
	RetencionGanacias          entities.Monto
	SignoRetencionGanacias     string
	RetencionIngresoBruto      entities.Monto
	SignoRetencionIngresoBruto string
}
