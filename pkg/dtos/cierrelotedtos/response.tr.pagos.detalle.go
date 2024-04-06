package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseTrPagosDetalle struct {
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
	PrismaCierreLote           []ResponsePrismaCL
}

func (rtpd *ResponseTrPagosDetalle) EntityToDtos(entityPagoDetalle entities.Prismatrdospago) {
	rtpd.Id = 0
	if entityPagoDetalle.ID > 0 {
		rtpd.Id = entityPagoDetalle.ID
	}
	rtpd.PrismatrcuatropagosId = 0
	if entityPagoDetalle.PrismatrcuatropagosId > 0 {
		rtpd.PrismatrcuatropagosId = entityPagoDetalle.PrismatrcuatropagosId
	}
	rtpd.FechaPresentacion = entityPagoDetalle.FechaPresentacion
	rtpd.TipoRegistro = entityPagoDetalle.TipoRegistro
	rtpd.Moneda = entityPagoDetalle.Moneda
	rtpd.LiquidacionNro = entityPagoDetalle.LiquidacionNro
	rtpd.FechaPago = entityPagoDetalle.FechaPago
	rtpd.LiquidacionTipo = entityPagoDetalle.LiquidacionTipo
	rtpd.ImporteBruto = entityPagoDetalle.ImporteBruto
	rtpd.SignoImporteBruto = entityPagoDetalle.SignoImporteBruto
	rtpd.ImporteArancel = entityPagoDetalle.ImporteArancel
	rtpd.SignoImporteArancel = entityPagoDetalle.SignoImporteArancel
	rtpd.ImporteNeto = entityPagoDetalle.ImporteNeto
	rtpd.SignoImporteNeto = entityPagoDetalle.SignoImporteNeto
	rtpd.RetencionEspecialSobreIibb = entityPagoDetalle.RetencionEspecialSobreIibb
	rtpd.SignoRetencionEspecial = entityPagoDetalle.SignoRetencionEspecial
	rtpd.RetencionIvaEspecial = entityPagoDetalle.RetencionIvaEspecial
	rtpd.SignoRetencionIvaEspecial = entityPagoDetalle.SignoRetencionIvaEspecial
	rtpd.PercepcionIngresoBruto = entityPagoDetalle.PercepcionIngresoBruto
	rtpd.SignoPercepcionIb = entityPagoDetalle.SignoPercepcionIb
	rtpd.RetencionIvaD1 = entityPagoDetalle.RetencionIvaD1
	rtpd.SignoRetencionIva_d1 = entityPagoDetalle.SignoRetencionIva_d1
	rtpd.CostoCuotaEmitida = entityPagoDetalle.CostoCuotaEmitida
	rtpd.SignoCostoCuotaEmitida = entityPagoDetalle.SignoCostoCuotaEmitida
	rtpd.RetencionIvaCuota = entityPagoDetalle.RetencionIvaCuota
	rtpd.SignoRetencionIvaCuota = entityPagoDetalle.SignoRetencionIvaCuota
	rtpd.RetencionIva = entityPagoDetalle.RetencionIva
	rtpd.SignoRetencionIva = entityPagoDetalle.SignoRetencionIva
	rtpd.RetencionGanacias = entityPagoDetalle.RetencionGanacias
	rtpd.SignoRetencionGanacias = entityPagoDetalle.SignoRetencionGanacias
	rtpd.RetencionIngresoBruto = entityPagoDetalle.RetencionIngresoBruto
	rtpd.SignoRetencionIngresoBruto = entityPagoDetalle.SignoRetencionIngresoBruto
}
