package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseMoviminetoDetalles struct {
	Id                           int64
	PrismamovimientototalesId    uint
	PrismamastercontracargosId   uint
	PrismaoperacionsId           uint
	RechazoPrincipalId           uint
	RechazoSecundarioId          uint
	PrismavisacontracargosId     uint
	PrismamotivosajustesId       uint
	TipoRegistro                 string
	TipoAplicacion               string
	Lote                         int64
	NroTarjeta                   string
	FechaOrigenCompra            time.Time
	FechaPago                    time.Time
	NroCupon                     int64
	Importe                      entities.Monto
	SignoImporte                 string
	NroAutorizacion              string
	NroCuota                     string
	PlanCuota                    int64
	RecAcep                      string
	NroLiquidacion               string
	ContracargoOrigen            string
	Moneda                       string
	IdCf                         string
	CfExentoIva                  string
	FechaPagoOrigenAjuste        string
	PorcentDescArancel           float64
	Arancel                      entities.Monto
	SignoArancel                 string
	TnaCf                        entities.Monto
	ImporteCostoFinanciero       entities.Monto
	SignoImporteCostoFinanciero  string
	BanderaEstablecimiento       string
	NroTarjetaXl                 string
	NroAutorizacionXl            string
	Match                        int
	Contracargovisa              ResponsePrismaVisaContracargo
	Contracargomaster            ResponsePrismaMasterContracargo
	Tipooperacion                ResponsePrismaOperacion
	Rechazotransaccionprincipal  ResponsePrismaCodigoRechazo
	Rechazotransaccionsecundario ResponsePrismaCodigoRechazo
	Motivoajuste                 ResponsePrismaMotivosAjuste
	CierreLote                   ResponsePrismaCL
}

func (rmd *ResponseMoviminetoDetalles) EntityToDtos(entityMovimientoDetalle entities.Prismamovimientodetalle) {
	var visaContraCargo ResponsePrismaVisaContracargo
	var masterContracargo ResponsePrismaMasterContracargo
	var codigoOperaciones ResponsePrismaOperacion
	var codigoRechazo ResponsePrismaCodigoRechazo
	var codigoRechazoSecundario ResponsePrismaCodigoRechazo
	var motivoAjuste ResponsePrismaMotivosAjuste
	rmd.Id = 0
	if entityMovimientoDetalle.ID > 0 {
		rmd.Id = int64(entityMovimientoDetalle.ID)
	}

	rmd.PrismamovimientototalesId = 0
	if entityMovimientoDetalle.PrismamovimientototalesId > 0 {
		rmd.PrismamovimientototalesId = entityMovimientoDetalle.PrismamovimientototalesId
	}

	rmd.PrismamastercontracargosId = 0
	if entityMovimientoDetalle.PrismamastercontracargosId > 0 {
		rmd.PrismamastercontracargosId = entityMovimientoDetalle.PrismamastercontracargosId
		masterContracargo.EntityToDtos(entityMovimientoDetalle.Contracargomaster)
		rmd.Contracargomaster = masterContracargo
	}

	rmd.PrismaoperacionsId = 0
	if entityMovimientoDetalle.PrismaoperacionsId > 0 {
		rmd.PrismaoperacionsId = entityMovimientoDetalle.PrismaoperacionsId
		codigoOperaciones.EntityToDtos(entityMovimientoDetalle.Tipooperacion)
		rmd.Tipooperacion = codigoOperaciones
	}
	rmd.RechazoPrincipalId = 0
	if entityMovimientoDetalle.RechazoPrincipalId > 0 {
		rmd.RechazoPrincipalId = entityMovimientoDetalle.RechazoPrincipalId
		codigoRechazo.EntityToDtos(entityMovimientoDetalle.Rechazotransaccionprincipal)
		rmd.Rechazotransaccionprincipal = codigoRechazo
	}
	rmd.RechazoSecundarioId = 0
	if entityMovimientoDetalle.RechazoSecundarioId > 0 {
		rmd.RechazoSecundarioId = entityMovimientoDetalle.RechazoSecundarioId
		codigoRechazoSecundario.EntityToDtos(entityMovimientoDetalle.Rechazotransaccionsecundario)
		rmd.Rechazotransaccionsecundario = codigoRechazoSecundario
	}
	rmd.PrismavisacontracargosId = 0
	if entityMovimientoDetalle.PrismavisacontracargosId > 0 {
		rmd.PrismavisacontracargosId = entityMovimientoDetalle.PrismavisacontracargosId
		visaContraCargo.EntityToDtos(entityMovimientoDetalle.Contracargovisa)
		rmd.Contracargovisa = visaContraCargo
	}
	rmd.PrismamotivosajustesId = 0
	if entityMovimientoDetalle.PrismamotivosajustesId > 0 {
		rmd.PrismamotivosajustesId = entityMovimientoDetalle.PrismamotivosajustesId
		motivoAjuste.EntityToDtos(entityMovimientoDetalle.Motivoajuste)
		rmd.Motivoajuste = motivoAjuste
	}
	rmd.TipoRegistro = entityMovimientoDetalle.TipoRegistro
	rmd.TipoAplicacion = entityMovimientoDetalle.TipoAplicacion
	rmd.Lote = entityMovimientoDetalle.Lote
	rmd.NroTarjeta = entityMovimientoDetalle.NroTarjeta
	rmd.FechaOrigenCompra = entityMovimientoDetalle.FechaOrigenCompra
	rmd.FechaPago = entityMovimientoDetalle.FechaPago
	rmd.NroCupon = entityMovimientoDetalle.NroCupon
	rmd.Importe = entityMovimientoDetalle.Importe
	rmd.SignoImporte = entityMovimientoDetalle.SignoImporte
	rmd.NroAutorizacion = entityMovimientoDetalle.NroAutorizacion
	rmd.NroCuota = entityMovimientoDetalle.NroCuota
	rmd.PlanCuota = entityMovimientoDetalle.PlanCuota
	rmd.RecAcep = entityMovimientoDetalle.RecAcep
	rmd.NroLiquidacion = entityMovimientoDetalle.NroLiquidacion
	rmd.ContracargoOrigen = entityMovimientoDetalle.ContracargoOrigen
	rmd.Moneda = entityMovimientoDetalle.Moneda
	rmd.IdCf = entityMovimientoDetalle.IdCf
	rmd.CfExentoIva = entityMovimientoDetalle.CfExentoIva
	rmd.FechaPagoOrigenAjuste = entityMovimientoDetalle.FechaPagoOrigenAjuste
	rmd.PorcentDescArancel = entityMovimientoDetalle.PorcentDescArancel
	rmd.Arancel = entityMovimientoDetalle.Arancel
	rmd.SignoArancel = entityMovimientoDetalle.SignoArancel
	rmd.TnaCf = entityMovimientoDetalle.TnaCf
	rmd.ImporteCostoFinanciero = entityMovimientoDetalle.ImporteCostoFinanciero
	rmd.SignoImporteCostoFinanciero = entityMovimientoDetalle.SignoImporteCostoFinanciero
	rmd.BanderaEstablecimiento = entityMovimientoDetalle.BanderaEstablecimiento
	rmd.NroTarjetaXl = entityMovimientoDetalle.NroTarjetaXl
	rmd.NroAutorizacionXl = entityMovimientoDetalle.NroAutorizacionXl
	rmd.Match = entityMovimientoDetalle.Match
}
