package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseMovimientoMxDetalle struct {
	Id                          uint
	PrismamovimientototalesId   uint
	Fechapresentacion           time.Time
	Tiporeg                     string
	PrismaoperacionsId          string
	Tipoaplic                   string
	Lote                        int64
	Numtar                      string
	FechaOrigenCompra           time.Time
	FechaPago                   time.Time
	NroCupon                    int64
	Importe                     entities.Monto
	SignoImporte                string
	NroAutorizacion             string
	NroCuota                    string
	PlanCuota                   string
	RecAcep                     string
	RechazoPrincipalId          string
	RechazoSecundarioId         string
	NroLiquidacion              string
	ContracargoOrigen           string
	PrismavisacontracargosId    string
	Moneda                      string
	IdCf                        string
	CfExentoIva                 string
	FechaPagoOrigenAjuste       string
	PrismamotivosajustesId      string
	PorcentDescArancel          float64
	Arancel                     entities.Monto
	SignoArancel                string
	TnaCf                       entities.Monto
	ImporteCostoFinanciero      entities.Monto
	SignoImporteCostoFinanciero string
	BanderaEstablecimiento      string
	PrismamastercontracargosId  string
	NroTarjetaXl                string
	NroAutorizacionXl           string
}
