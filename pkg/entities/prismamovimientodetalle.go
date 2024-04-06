package entities

import (
	"time"

	"gorm.io/gorm"
)

type Prismamovimientodetalle struct {
	gorm.Model
	PrismamovimientototalesId    uint `gorm:"column:prismamovimientototales_id"`
	PrismamastercontracargosId   uint `gorm:"column:prismamastercontracargos_id"`
	PrismaoperacionsId           uint `gorm:"column:prismaoperacions_id"`
	RechazoPrincipalId           uint `gorm:"column:rechazo_principal_id"`
	RechazoSecundarioId          uint `gorm:"column:rechazo_secundario_id"`
	PrismavisacontracargosId     uint `gorm:"column:prismavisacontracargos_id"`
	PrismamotivosajustesId       uint `gorm:"column:prismamotivosajustes_id"`
	TipoRegistro                 string
	TipoAplicacion               string
	Lote                         int64
	NroTarjeta                   string
	FechaOrigenCompra            time.Time
	FechaPago                    time.Time
	NroCupon                     int64
	Importe                      Monto
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
	Arancel                      Monto
	SignoArancel                 string
	TnaCf                        Monto
	ImporteCostoFinanciero       Monto
	SignoImporteCostoFinanciero  string
	BanderaEstablecimiento       string //enum
	NroTarjetaXl                 string
	NroAutorizacionXl            string
	Match                        int
	Contracargovisa              Prismavisacontracargo   `gorm:"foreignkey:PrismavisacontracargosId"`
	Contracargomaster            Prismamastercontracargo `gorm:"foreignkey:PrismamastercontracargosId"`
	Tipooperacion                Prismaoperacion         `gorm:"foreignkey:PrismaoperacionsId"`
	Rechazotransaccionprincipal  Prismacodigorechazo     `gorm:"foreignkey:RechazoPrincipalId"`
	Rechazotransaccionsecundario Prismacodigorechazo     `gorm:"foreignkey:RechazoSecundarioId"`
	Motivoajuste                 Prismamotivosajuste     `gorm:"foreignkey:PrismamotivosajustesId"`
	MovimientoCabecera           Prismamovimientototale  `gorm:"foreignkey:PrismamovimientototalesId"`
	CierreLote                   Prismacierrelote        `gorm:"prismamovimientodetalles_id"`
}
