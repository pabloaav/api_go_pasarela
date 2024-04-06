package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseTrPagosCabecera struct {
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
	DetallePago                   []ResponseTrPagosDetalle
}

func (rtpc *ResponseTrPagosCabecera) EntityToDtos(entityPagoCabecera entities.Prismatrcuatropago) {
	rtpc.Id = 0
	if entityPagoCabecera.ID > 0 {
		rtpc.Id = entityPagoCabecera.ID
	}
	rtpc.Empresa = entityPagoCabecera.Empresa
	rtpc.FechaPresentacion = entityPagoCabecera.FechaPresentacion
	rtpc.TipoRegistro = entityPagoCabecera.TipoRegistro
	rtpc.Moneda = entityPagoCabecera.Moneda
	rtpc.ComercioNro = entityPagoCabecera.ComercioNro
	rtpc.EstablecimientoNro = entityPagoCabecera.EstablecimientoNro
	rtpc.LiquidacionNro = entityPagoCabecera.LiquidacionNro
	rtpc.FechaPago = entityPagoCabecera.FechaPago
	rtpc.LiquidacionTipo = entityPagoCabecera.LiquidacionTipo
	rtpc.CasaCuentaAcreditacion = entityPagoCabecera.CasaCuentaAcreditacion
	rtpc.TipoCuentaAcreditacion = entityPagoCabecera.TipoCuentaAcreditacion
	rtpc.NroCuentaAcreditacion = entityPagoCabecera.NroCuentaAcreditacion
	rtpc.CostoFinacieroExentoIva = entityPagoCabecera.CostoFinacieroExentoIva
	rtpc.SignoCostoFinacieroExtIva = entityPagoCabecera.SignoCostoFinacieroExtIva
	rtpc.RetencionPorLey_25063 = entityPagoCabecera.RetencionPorLey_25063
	rtpc.SignoRetencionPorLey = entityPagoCabecera.SignoRetencionPorLey
	rtpc.AlicuotaRetencionIb = entityPagoCabecera.AlicuotaRetencionIb
	rtpc.CargoTransferenciaBancaria = entityPagoCabecera.CargoTransferenciaBancaria
	rtpc.SignoCargoTransBancaria = entityPagoCabecera.SignoCargoTransBancaria
	rtpc.IvaSobreCargoTransBancaria = entityPagoCabecera.IvaSobreCargoTransBancaria
	rtpc.SignoIvaCargoTransBancaria = entityPagoCabecera.SignoIvaCargoTransBancaria
	rtpc.ImpuestoDbCr = entityPagoCabecera.ImpuestoDbCr
	rtpc.SignoImpuestoDbCr = entityPagoCabecera.SignoImpuestoDbCr
	rtpc.CostoFinancieroNoReducIva = entityPagoCabecera.CostoFinancieroNoReducIva
	rtpc.SignoCostoFinancieroReducIva = entityPagoCabecera.SignoCostoFinancieroReducIva
	rtpc.RetencionIvaRg_3130 = entityPagoCabecera.RetencionIvaRg_3130
	rtpc.SignoRetencionIvaRg_3130 = entityPagoCabecera.SignoRetencionIvaRg_3130
	rtpc.JurisdiccionRetencionIb = entityPagoCabecera.JurisdiccionRetencionIb
	rtpc.CargoAdicionalPlanCuotas = entityPagoCabecera.CargoAdicionalPlanCuotas
	rtpc.SignoCargoAdicional = entityPagoCabecera.SignoCargoAdicional
	rtpc.IvaCargoAdicionalPlanCuota = entityPagoCabecera.IvaCargoAdicionalPlanCuota
	rtpc.SignoIvaCargoAdicional = entityPagoCabecera.SignoIvaCargoAdicional
	rtpc.CargoAdicionalMovipag = entityPagoCabecera.CargoAdicionalMovipag
	rtpc.SignoCargiMovipag = entityPagoCabecera.SignoCargiMovipag
	rtpc.IvaCargoAdicionalMovipag = entityPagoCabecera.IvaCargoAdicionalMovipag
	rtpc.SignoIvaCargoAdicionalMovipag = entityPagoCabecera.SignoIvaCargoAdicionalMovipag
	rtpc.RetencionSello = entityPagoCabecera.RetencionSello
	rtpc.SignoRetencionSello = entityPagoCabecera.SignoRetencionSello
	rtpc.ProvinciaRetencionSello = entityPagoCabecera.ProvinciaRetencionSello
}
