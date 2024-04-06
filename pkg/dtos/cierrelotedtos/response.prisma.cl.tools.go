package cierrelotedtos

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePrismaCLTools struct {
	CierresLotes []PrismaCLTools `json:"cierreslotes"`
	Meta         dtos.Meta       `json:"meta"`
}

type PrismaCLTools struct {
	ClId                       uint                        `json:"cl_id"`
	Fechaoperacion             time.Time                   `json:"fechaoperacion"`
	FechaCierre                time.Time                   `json:"fecha_cierre"`
	FechaPago                  time.Time                   `json:"fecha_pago"`
	FechaCreacion              time.Time                   `json:"fecha_creacion"`
	PagoestadoexternosId       int64                       `json:"pagoestadoexternos_id"`
	ChannelarancelesId         int64                       `json:"channelaranceles_id"`
	ImpuestosId                int64                       `json:"impuestos_id"`
	PrismamovimientodetallesId int64                       `json:"prismamovimientodetalles_id"`
	PrismatrdospagosId         int64                       `json:"prismatrdospagos_id"`
	BancoExternalId            int64                       `json:"banco_external_id"`
	Tiporegistro               string                      `json:"tiporegistro"`
	ExternalmediopagoId        int64                       `json:"externalmediopago"`
	Tipooperacion              entities.EnumTipoOperacion  `json:"tipooperacion"`
	Monto                      entities.Monto              `json:"monto"`
	Montofinal                 entities.Monto              `json:"montofinal"`
	Valorpresentado            entities.Monto              `json:"valorpresentado"`
	Diferenciaimporte          entities.Monto              `json:"difernciaimporte"`
	Coeficientecalculado       float64                     `json:"coeficientecalculado"`
	Costototalporcentaje       float64                     `json:"costototalporcentaje"`
	Importeivaarancel          float64                     `json:"importeivaarancel"`
	Codigoautorizacion         string                      `json:"codigoautorizacion"`
	Nroticket                  int64                       `json:"nroticket"`
	SiteID                     int64                       `json:"site_id"`
	ExternalloteId             int64                       `json:"externallote_id"`
	Nrocuota                   int64                       `json:"nrocuota"`
	Nroestablecimiento         int64                       `json:"nroestablecimiento"`
	PagosUuid                  string                      `json:"pagos_uuid"`
	ExternalclienteID          string                      `json:"externalcliente_id"`
	Cantdias                   int                         `json:"cantdias"`
	Match                      int                         `json:"match"`
	Disputa                    bool                        `json:"disputa"`
	Enobservacion              bool                        `json:"enobservacion"`
	Channelarancel             ResponseChannelArancel      `json:"channelarancel"`
	Descripcion                string                      `json:"descripcion"`
	Descripcionpresentacion    string                      `json:"descripcionpresentacion"`
	Nombrearchivolote          string                      `json:"nombrearchivolote"`
	MovimientoPrisma           ResponseMoviminetoDetPrisma `json:"movimiento_prisma"`
	PagoPrisma                 ResponsePagoPrisma          `json:"pago_prisma"`
	Reversion                  bool                        `json:"reversion"`
	DetallemovimientoId        int64                       `json:"detallemovimineto_id"`
	DetallepagoId              int64                       `json:"detallepago_id"`
	Descripcioncontracargo     string                      `json:"descripcioncontracargo"`
	ExtbancoreversionId        int64                       `json:"extbancoreversion_id"`
	Conciliado                 bool                        `json:"conciliacion"`
	Estadomovimiento           bool                        `json:"estadomovimineto"`
	Descripcionbanco           string                      `json:"descripcionbanco"`
}

func (pclt *PrismaCLTools) EntityCLToDto(resultado *entities.Prismacierrelote) {
	pclt.ClId = resultado.Model.ID
	pclt.Fechaoperacion = resultado.Fechaoperacion
	pclt.FechaCierre = resultado.FechaCierre
	pclt.FechaPago = resultado.FechaPago
	pclt.FechaCreacion = resultado.Model.CreatedAt
	pclt.PagoestadoexternosId = resultado.PagoestadoexternosId
	pclt.ChannelarancelesId = resultado.ChannelarancelesId
	pclt.ImpuestosId = resultado.ImpuestosId
	pclt.PrismamovimientodetallesId = resultado.PrismamovimientodetallesId
	pclt.PrismatrdospagosId = resultado.PrismatrdospagosId
	pclt.BancoExternalId = resultado.BancoExternalId
	pclt.Tiporegistro = resultado.Tiporegistro
	pclt.ExternalmediopagoId = resultado.ExternalmediopagoId
	pclt.Tipooperacion = resultado.Tipooperacion
	pclt.Monto = resultado.Monto
	pclt.Montofinal = resultado.Montofinal
	pclt.Valorpresentado = resultado.Valorpresentado
	pclt.Diferenciaimporte = resultado.Diferenciaimporte
	pclt.Coeficientecalculado = resultado.Coeficientecalculado
	pclt.Costototalporcentaje = resultado.Costototalporcentaje
	pclt.Importeivaarancel = resultado.Importeivaarancel
	pclt.Codigoautorizacion = resultado.Codigoautorizacion
	pclt.Nroticket = resultado.Nroticket
	pclt.SiteID = resultado.SiteID
	pclt.ExternalloteId = resultado.ExternalloteId
	pclt.Nrocuota = resultado.Nrocuota
	pclt.Nroestablecimiento = resultado.Nroestablecimiento
	pclt.PagosUuid = resultado.PagosUuid
	pclt.ExternalclienteID = resultado.ExternalclienteID
	pclt.Cantdias = int(resultado.Cantdias)
	pclt.Match = resultado.Match
	pclt.Disputa = resultado.Disputa
	pclt.Enobservacion = resultado.Enobservacion
	//pclt.Channelarancel             = resultado.cha
	pclt.Descripcion = resultado.Descripcion
	pclt.Descripcionpresentacion = resultado.Descripcionpresentacion
	pclt.Nombrearchivolote = resultado.Nombrearchivolote
	pclt.Reversion = resultado.Reversion
	pclt.DetallemovimientoId = resultado.DetallemovimientoId
	pclt.DetallepagoId = resultado.DetallepagoId
	pclt.Descripcioncontracargo = resultado.Descripcioncontracargo
	pclt.ExtbancoreversionId = resultado.ExtbancoreversionId
	pclt.Conciliado = resultado.Conciliado
	pclt.Estadomovimiento = resultado.Estadomovimiento
	pclt.Descripcionbanco = resultado.Descripcionbanco
}

type ResponseMoviminetoDetPrisma struct {
	MovCabeceraId                uint
	MovDetalleId                 uint
	FechaOrigenCompra            time.Time
	FechaPresentacion            time.Time
	FechaPago                    time.Time
	Empresa                      string
	ComercioNro                  string
	EstablecimientoNro           string
	Codop                        string
	ImporteTotal                 entities.Monto
	SignoImporteTotal            string
	TipoRegistro                 string
	TipoAplicacion               string
	Lote                         int64
	NroCupon                     int64
	Importe                      entities.Monto
	SignoImporte                 string
	NroCuota                     string
	PlanCuota                    int64
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
	NroAutorizacionXl            string
	Contracargovisa              string
	Contracargomaster            string
	Tipooperacion                string
	Rechazotransaccionprincipal  string
	Rechazotransaccionsecundario string
	Motivoajuste                 string
	EstadoCab                    bool
	EstadoDet                    bool
}

func (rmdp *ResponseMoviminetoDetPrisma) EntityMovToDto(dtosMovimiento *entities.Prismamovimientodetalle) {
	rmdp.MovCabeceraId = dtosMovimiento.MovimientoCabecera.ID
	rmdp.MovDetalleId = dtosMovimiento.Model.ID
	rmdp.FechaOrigenCompra = dtosMovimiento.FechaOrigenCompra
	rmdp.FechaPresentacion = dtosMovimiento.MovimientoCabecera.FechaPresentacion
	rmdp.FechaPago = dtosMovimiento.FechaPago
	rmdp.Empresa = dtosMovimiento.MovimientoCabecera.Empresa
	rmdp.ComercioNro = dtosMovimiento.MovimientoCabecera.ComercioNro
	rmdp.EstablecimientoNro = dtosMovimiento.MovimientoCabecera.EstablecimientoNro
	rmdp.Codop = dtosMovimiento.MovimientoCabecera.Codop
	rmdp.ImporteTotal = dtosMovimiento.MovimientoCabecera.ImporteTotal
	rmdp.SignoImporteTotal = dtosMovimiento.MovimientoCabecera.SignoImporteTotal
	rmdp.TipoRegistro = dtosMovimiento.TipoRegistro
	rmdp.TipoAplicacion = dtosMovimiento.TipoAplicacion
	rmdp.Lote = dtosMovimiento.Lote
	rmdp.NroCupon = dtosMovimiento.NroCupon
	rmdp.Importe = dtosMovimiento.Importe
	rmdp.SignoImporte = dtosMovimiento.SignoImporte
	rmdp.NroCuota = dtosMovimiento.NroCuota
	rmdp.PlanCuota = dtosMovimiento.PlanCuota
	rmdp.NroLiquidacion = dtosMovimiento.NroLiquidacion
	rmdp.ContracargoOrigen = dtosMovimiento.ContracargoOrigen
	rmdp.Moneda = dtosMovimiento.Moneda
	rmdp.IdCf = dtosMovimiento.IdCf
	rmdp.CfExentoIva = dtosMovimiento.CfExentoIva
	rmdp.FechaPagoOrigenAjuste = dtosMovimiento.FechaPagoOrigenAjuste
	rmdp.PorcentDescArancel = dtosMovimiento.PorcentDescArancel
	rmdp.Arancel = dtosMovimiento.Arancel
	rmdp.SignoArancel = dtosMovimiento.SignoArancel
	rmdp.TnaCf = dtosMovimiento.TnaCf
	rmdp.ImporteCostoFinanciero = dtosMovimiento.ImporteCostoFinanciero
	rmdp.SignoImporteCostoFinanciero = dtosMovimiento.SignoImporteCostoFinanciero
	rmdp.BanderaEstablecimiento = dtosMovimiento.BanderaEstablecimiento
	rmdp.NroAutorizacionXl = dtosMovimiento.NroAutorizacionXl
	rmdp.Contracargovisa = fmt.Sprintf("%v-%v", dtosMovimiento.Contracargovisa.ExternalId, dtosMovimiento.Contracargovisa.Contracargo)
	rmdp.Contracargomaster = fmt.Sprintf("%v-%v", dtosMovimiento.Contracargomaster.ExternalId, dtosMovimiento.Contracargomaster.Contracargo)
	rmdp.Tipooperacion = fmt.Sprintf("%v-%v", dtosMovimiento.Tipooperacion.ExternalId, dtosMovimiento.Tipooperacion.Operacion)
	rmdp.Rechazotransaccionprincipal = fmt.Sprintf("%v-%v", dtosMovimiento.Rechazotransaccionprincipal.ExternalId, dtosMovimiento.Rechazotransaccionprincipal.Rechazo)
	rmdp.Rechazotransaccionsecundario = fmt.Sprintf("%v-%v", dtosMovimiento.Rechazotransaccionsecundario.ExternalId, dtosMovimiento.Rechazotransaccionsecundario.Rechazo)
	rmdp.Motivoajuste = fmt.Sprintf("%v-%v", dtosMovimiento.Motivoajuste.ExternalId, dtosMovimiento.Motivoajuste.Motivoajustes)
	if dtosMovimiento.MovimientoCabecera.Match == 1 {
		rmdp.EstadoCab = true
	}
	if dtosMovimiento.MovimientoCabecera.Match == 0 {
		rmdp.EstadoCab = false
	}
	if dtosMovimiento.Match == 1 {
		rmdp.EstadoDet = true
	}
	if dtosMovimiento.Match == 0 {
		rmdp.EstadoDet = false
	}

}

type ResponsePagoPrisma struct {
	PagoCabeceraId             uint
	PagoDetalleId              uint
	Empresa                    string
	FechaPresentacion          time.Time
	FechaPago                  time.Time
	TipoRegistro               string
	Moneda                     string
	ComercioNro                string
	EstablecimientoNro         string
	LiquidacionNro             string
	LiquidacionTipo            string
	RetencionSello             entities.Monto
	SignoRetencionSello        string
	ProvinciaRetencionSello    string
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

func (rpp *ResponsePagoPrisma) EntityPagoToDto(dtosPago *entities.Prismatrdospago) {
	rpp.PagoCabeceraId = dtosPago.PrismatrcuatropagosId
	rpp.PagoDetalleId = dtosPago.ID
	// rpp.Empresa = dtosPago .Empresa
	rpp.FechaPresentacion = dtosPago.FechaPresentacion
	rpp.FechaPago = dtosPago.FechaPago
	rpp.TipoRegistro = dtosPago.TipoRegistro
	rpp.Moneda = dtosPago.Moneda
	// rpp.ComercioNro = dtosPago.ComercioNro
	// rpp.EstablecimientoNro = dtosPago.EstablecimientoNro
	rpp.LiquidacionNro = dtosPago.LiquidacionNro
	rpp.LiquidacionTipo = dtosPago.LiquidacionTipo
	// rpp.RetencionSello = dtosPago.RetencionSello
	// rpp.SignoRetencionSello = dtosPago.SignoRetencionSello
	// rpp.ProvinciaRetencionSello = dtosPago.ProvinciaRetencionSello
	rpp.ImporteBruto = dtosPago.ImporteBruto
	rpp.SignoImporteBruto = dtosPago.SignoImporteBruto
	rpp.ImporteArancel = dtosPago.ImporteArancel
	rpp.SignoImporteArancel = dtosPago.SignoImporteArancel
	rpp.ImporteNeto = dtosPago.ImporteNeto
	rpp.SignoImporteNeto = dtosPago.SignoImporteNeto
	rpp.RetencionEspecialSobreIibb = dtosPago.RetencionEspecialSobreIibb
	rpp.SignoRetencionEspecial = dtosPago.SignoRetencionEspecial
	rpp.RetencionIvaEspecial = dtosPago.RetencionIvaEspecial
	rpp.SignoRetencionIvaEspecial = dtosPago.SignoRetencionIvaEspecial
	rpp.PercepcionIngresoBruto = dtosPago.PercepcionIngresoBruto
	rpp.SignoPercepcionIb = dtosPago.SignoPercepcionIb
	rpp.RetencionIvaD1 = dtosPago.RetencionIvaD1
	rpp.SignoRetencionIva_d1 = dtosPago.SignoRetencionIva_d1
	rpp.CostoCuotaEmitida = dtosPago.CostoCuotaEmitida
	rpp.SignoCostoCuotaEmitida = dtosPago.SignoCostoCuotaEmitida
	rpp.RetencionIvaCuota = dtosPago.RetencionIvaCuota
	rpp.SignoRetencionIvaCuota = dtosPago.SignoRetencionIvaCuota
	rpp.RetencionIva = dtosPago.RetencionIva
	rpp.SignoRetencionIva = dtosPago.SignoRetencionIva
	rpp.RetencionGanacias = dtosPago.RetencionGanacias
	rpp.SignoRetencionGanacias = dtosPago.SignoRetencionGanacias
	rpp.RetencionIngresoBruto = dtosPago.RetencionIngresoBruto
	rpp.SignoRetencionIngresoBruto = dtosPago.SignoRetencionIngresoBruto
}
