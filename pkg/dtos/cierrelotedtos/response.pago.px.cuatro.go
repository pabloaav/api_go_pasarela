package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagoPxCuatro struct {
	Id                            uint
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
	CostoFinacieroExentoIva       entities.Monto
	SignoCostoFinacieroExtIva     string
	RetencionPorLey_25063         entities.Monto
	SignoRetencionPorLey          string
	AlicuotaRetencionIb           entities.Monto
	CargoTransferenciaBancaria    entities.Monto
	SignoCargoTransBancaria       string
	IvaSobreCargoTransBancaria    entities.Monto
	SignoIvaCargoTransBancaria    string
	ImpuestoDbCr                  entities.Monto
	SignoImpuestoDbCr             string
	CostoFinancieroNoReducIva     entities.Monto
	SignoCostoFinancieroReducIva  string
	RetencionIvaRg_3130           entities.Monto
	SignoRetencionIvaRg_3130      string
	JurisdiccionRetencionIb       string
	CargoAdicionalPlanCuotas      entities.Monto
	SignoCargoAdicional           string
	IvaCargoAdicionalPlanCuota    entities.Monto
	SignoIvaCargoAdicional        string
	CargoAdicionalMovipag         entities.Monto
	SignoCargiMovipag             string
	IvaCargoAdicionalMovipag      entities.Monto
	SignoIvaCargoAdicionalMovipag string
	RetencionSello                entities.Monto
	SignoRetencionSello           string
	ProvinciaRetencionSello       string
}
