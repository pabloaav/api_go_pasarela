package cierrelotedtos

import (
	"errors"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagoPx struct {
	PagoPxCuatro ResponsePagoPxCuatro
	PagoPxDos    []ResponsePagoPxDos
}

func (rpx *ResponsePagoPx) EntityToDtosPx(pagoPxEntity []entities.Prismapxcuatroregistro) (listaPagosPx []ResponsePagoPx, erro error) {
	msgErrorFecha := errors.New(CONVERTIR_FECHA_ERROR)
	msgErrorEntero := errors.New(CONVERTIR_A_ENTERO_ERROR)
	for _, valuePagoPxCabecera := range pagoPxEntity {
		var pxCabecera ResponsePagoPxCuatro
		pxCabecera.Id = valuePagoPxCabecera.ID
		pxCabecera.Empresa = valuePagoPxCabecera.Eclq02llEmpresa_04
		fechaPresentacion, err := time.Parse("020106", valuePagoPxCabecera.Eclq02llFpres_04)
		if err != nil {
			erro = msgErrorFecha
			return
		}
		fechaPago, err := time.Parse("020106", valuePagoPxCabecera.Eclq02llFpag_04)
		if err != nil {
			erro = msgErrorFecha
			return
		}
		costoFinancieroIva, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llCfExentoIva, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		retencionLey25063, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llLey25063, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		aliluotaIngBruto, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llAliIngbru, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		cargoTransferenciaBanc, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llCargoXLiq, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		ivaCargoTransferenciaBanc, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llIva1CargoXLiq, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		importeDbCr, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llImpDbCr, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		ivaCostoFinanciero, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llCfNoReduceIva, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		retencionivaRg3130, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llSubtotalRetivaRg3130, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		cargoAdicionalPlanCuotas, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llAdicPlancuo, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}

		ivaCargoAdicionalPlanCuotas, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llIva1AdPlancuo, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}

		cargoAdicionalMovyPago, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llAdicMovpag, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}

		ivaCargoAdicionalMovyPago, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llIva1AdicMovpag, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		retencionSello, err := strconv.ParseInt(valuePagoPxCabecera.Eclq02llRetSellos, 10, 64)
		if err != nil {
			erro = msgErrorEntero
			return
		}
		pxCabecera.FechaPresentacion = fechaPresentacion
		pxCabecera.TipoRegistro = valuePagoPxCabecera.Eclq02llTiporeg_04
		pxCabecera.Moneda = valuePagoPxCabecera.Eclq02llMoneda_04
		pxCabecera.ComercioNro = valuePagoPxCabecera.Eclq02llNumcom_04
		pxCabecera.EstablecimientoNro = valuePagoPxCabecera.Eclq02llNumest_04
		pxCabecera.LiquidacionNro = valuePagoPxCabecera.Eclq02llNroliq_04
		pxCabecera.FechaPago = fechaPago
		pxCabecera.LiquidacionTipo = valuePagoPxCabecera.Eclq02llTipoliq_04
		pxCabecera.CasaCuentaAcreditacion = valuePagoPxCabecera.Eclq02llCasacta
		pxCabecera.TipoCuentaAcreditacion = valuePagoPxCabecera.Eclq02llTipcta
		pxCabecera.NroCuentaAcreditacion = valuePagoPxCabecera.Eclq02llCtabco
		pxCabecera.CostoFinacieroExentoIva = entities.Monto(costoFinancieroIva)
		pxCabecera.SignoCostoFinacieroExtIva = valuePagoPxCabecera.Eclq02llSigno_04_1
		pxCabecera.RetencionPorLey_25063 = entities.Monto(retencionLey25063)
		pxCabecera.SignoRetencionPorLey = valuePagoPxCabecera.Eclq02llSigno_04_2
		pxCabecera.AlicuotaRetencionIb = entities.Monto(aliluotaIngBruto)
		pxCabecera.CargoTransferenciaBancaria = entities.Monto(cargoTransferenciaBanc)
		pxCabecera.SignoCargoTransBancaria = valuePagoPxCabecera.Eclq02llSigno_04_8
		pxCabecera.IvaSobreCargoTransBancaria = entities.Monto(ivaCargoTransferenciaBanc)
		pxCabecera.SignoIvaCargoTransBancaria = valuePagoPxCabecera.Eclq02llSigno_04_9
		pxCabecera.ImpuestoDbCr = entities.Monto(importeDbCr)
		pxCabecera.SignoImpuestoDbCr = valuePagoPxCabecera.Eclq02llSigno_04_10
		pxCabecera.CostoFinancieroNoReducIva = entities.Monto(ivaCostoFinanciero)
		pxCabecera.SignoCostoFinancieroReducIva = valuePagoPxCabecera.Eclq02llSigno_04_11
		pxCabecera.RetencionIvaRg_3130 = entities.Monto(retencionivaRg3130)
		pxCabecera.SignoRetencionIvaRg_3130 = valuePagoPxCabecera.Eclq02llSigno_04_14
		pxCabecera.JurisdiccionRetencionIb = valuePagoPxCabecera.Eclq02llProvIngbru
		pxCabecera.CargoAdicionalPlanCuotas = entities.Monto(cargoAdicionalPlanCuotas)
		pxCabecera.SignoCargoAdicional = valuePagoPxCabecera.Eclq02llSigno_04_15
		pxCabecera.IvaCargoAdicionalPlanCuota = entities.Monto(ivaCargoAdicionalPlanCuotas)
		pxCabecera.SignoIvaCargoAdicional = valuePagoPxCabecera.Eclq02llSigno_04_16
		pxCabecera.CargoAdicionalMovipag = entities.Monto(cargoAdicionalMovyPago)
		pxCabecera.SignoCargiMovipag = valuePagoPxCabecera.Eclq02llSigno_04_27
		pxCabecera.IvaCargoAdicionalMovipag = entities.Monto(ivaCargoAdicionalMovyPago)
		pxCabecera.SignoIvaCargoAdicionalMovipag = valuePagoPxCabecera.Eclq02llSigno_04_28
		pxCabecera.RetencionSello = entities.Monto(retencionSello)
		pxCabecera.SignoRetencionSello = valuePagoPxCabecera.Eclq02llSigno_29
		pxCabecera.ProvinciaRetencionSello = valuePagoPxCabecera.Eclq02llProvSellos
		var PxDetalle []ResponsePagoPxDos
		for _, valuePagoPxDetalle := range valuePagoPxCabecera.PxDosRegistros {

			fechaPresentacion, err := time.Parse("020106", valuePagoPxDetalle.Eclq02llFpres)
			if err != nil {
				erro = msgErrorFecha
				return
			}
			fechaPago, err := time.Parse("020106", valuePagoPxDetalle.Eclq02llFpag)
			if err != nil {
				erro = msgErrorFecha
				return
			}
			eclq02llImpbruto, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llImpbruto, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}

			eclq02llImppret, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llImppret, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llImpneto, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llImpneto, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetesp, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetesp, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetivaEsp, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetivaEsp, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llPercepBa, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llPercepBa, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetivaD1, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetivaD1, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llCostoCuoemi, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llCostoCuoemi, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetivaCuo1, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetivaCuo1, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetIva, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetIva, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetGcias, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetGcias, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}
			eclq02llRetIngbru, err := strconv.ParseInt(valuePagoPxDetalle.Eclq02llRetIngbru, 10, 64)
			if err != nil {
				erro = msgErrorEntero
				return
			}

			temporalPagoPxDetalle := ResponsePagoPxDos{
				Id:                         valuePagoPxDetalle.ID,
				PrismatrcuatropagosId:      valuePagoPxDetalle.PrismapxcuatroregistrosId,
				FechaPresentacion:          fechaPresentacion,
				TipoRegistro:               valuePagoPxDetalle.Eclq02llTiporeg,
				Moneda:                     valuePagoPxDetalle.Eclq02llMoneda,
				LiquidacionNro:             valuePagoPxDetalle.Eclq02llNroliq,
				FechaPago:                  fechaPago,
				LiquidacionTipo:            valuePagoPxDetalle.Eclq02llTipoliq,
				ImporteBruto:               entities.Monto(eclq02llImpbruto),
				SignoImporteBruto:          valuePagoPxDetalle.Eclq02llSigno_1,
				ImporteArancel:             entities.Monto(eclq02llImppret),
				SignoImporteArancel:        valuePagoPxDetalle.Eclq02llSigno_2,
				ImporteNeto:                entities.Monto(eclq02llImpneto),
				SignoImporteNeto:           valuePagoPxDetalle.Eclq02llSigno_3,
				RetencionEspecialSobreIibb: entities.Monto(eclq02llRetesp),
				SignoRetencionEspecial:     valuePagoPxDetalle.Eclq02llSigno_4,
				RetencionIvaEspecial:       entities.Monto(eclq02llRetivaEsp),
				SignoRetencionIvaEspecial:  valuePagoPxDetalle.Eclq02llSigno_5,
				PercepcionIngresoBruto:     entities.Monto(eclq02llPercepBa),
				SignoPercepcionIb:          valuePagoPxDetalle.Eclq02llSigno_6,
				RetencionIvaD1:             entities.Monto(eclq02llRetivaD1),
				SignoRetencionIva_d1:       valuePagoPxDetalle.Eclq02llSigno_7,
				CostoCuotaEmitida:          entities.Monto(eclq02llCostoCuoemi),
				SignoCostoCuotaEmitida:     valuePagoPxDetalle.Eclq02llSigno_12,
				RetencionIvaCuota:          entities.Monto(eclq02llRetivaCuo1),
				SignoRetencionIvaCuota:     valuePagoPxDetalle.Eclq02llSigno_13,
				RetencionIva:               entities.Monto(eclq02llRetIva),
				SignoRetencionIva:          valuePagoPxDetalle.Eclq02llSigno_30,
				RetencionGanacias:          entities.Monto(eclq02llRetGcias),
				SignoRetencionGanacias:     valuePagoPxDetalle.Eclq02llSigno_31,
				RetencionIngresoBruto:      entities.Monto(eclq02llRetIngbru),
				SignoRetencionIngresoBruto: valuePagoPxDetalle.Eclq02llSigno_32,
			}
			PxDetalle = append(PxDetalle, temporalPagoPxDetalle)
		}
		listaPagosPx = append(listaPagosPx, ResponsePagoPx{pxCabecera, PxDetalle})
	}
	return
}

func (responsePx *ResponsePagoPx) ToEntity() (entityPago entities.Prismatrcuatropago, erro error) {
	entityPago.Empresa = responsePx.PagoPxCuatro.Empresa
	entityPago.FechaPresentacion = responsePx.PagoPxCuatro.FechaPresentacion
	entityPago.TipoRegistro = responsePx.PagoPxCuatro.TipoRegistro
	entityPago.Moneda = responsePx.PagoPxCuatro.Moneda
	entityPago.ComercioNro = responsePx.PagoPxCuatro.ComercioNro
	entityPago.EstablecimientoNro = responsePx.PagoPxCuatro.EstablecimientoNro
	entityPago.LiquidacionNro = responsePx.PagoPxCuatro.LiquidacionNro
	entityPago.FechaPago = responsePx.PagoPxCuatro.FechaPago
	entityPago.LiquidacionTipo = responsePx.PagoPxCuatro.LiquidacionTipo
	entityPago.CasaCuentaAcreditacion = responsePx.PagoPxCuatro.CasaCuentaAcreditacion
	entityPago.TipoCuentaAcreditacion = responsePx.PagoPxCuatro.TipoCuentaAcreditacion
	entityPago.NroCuentaAcreditacion = responsePx.PagoPxCuatro.NroCuentaAcreditacion
	entityPago.CostoFinacieroExentoIva = responsePx.PagoPxCuatro.CostoFinacieroExentoIva
	entityPago.SignoCostoFinacieroExtIva = responsePx.PagoPxCuatro.SignoCostoFinacieroExtIva
	entityPago.RetencionPorLey_25063 = responsePx.PagoPxCuatro.RetencionPorLey_25063
	entityPago.SignoRetencionPorLey = responsePx.PagoPxCuatro.SignoRetencionPorLey
	entityPago.AlicuotaRetencionIb = responsePx.PagoPxCuatro.AlicuotaRetencionIb
	entityPago.CargoTransferenciaBancaria = responsePx.PagoPxCuatro.CargoTransferenciaBancaria
	entityPago.SignoCargoTransBancaria = responsePx.PagoPxCuatro.SignoCargoTransBancaria
	entityPago.IvaSobreCargoTransBancaria = responsePx.PagoPxCuatro.IvaSobreCargoTransBancaria
	entityPago.SignoIvaCargoTransBancaria = responsePx.PagoPxCuatro.SignoIvaCargoTransBancaria
	entityPago.ImpuestoDbCr = responsePx.PagoPxCuatro.ImpuestoDbCr
	entityPago.SignoImpuestoDbCr = responsePx.PagoPxCuatro.SignoImpuestoDbCr
	entityPago.CostoFinancieroNoReducIva = responsePx.PagoPxCuatro.CostoFinancieroNoReducIva
	entityPago.SignoCostoFinancieroReducIva = responsePx.PagoPxCuatro.SignoCostoFinancieroReducIva
	entityPago.RetencionIvaRg_3130 = responsePx.PagoPxCuatro.RetencionIvaRg_3130
	entityPago.SignoRetencionIvaRg_3130 = responsePx.PagoPxCuatro.SignoRetencionIvaRg_3130
	entityPago.JurisdiccionRetencionIb = responsePx.PagoPxCuatro.JurisdiccionRetencionIb
	entityPago.CargoAdicionalPlanCuotas = responsePx.PagoPxCuatro.CargoAdicionalPlanCuotas
	entityPago.SignoCargoAdicional = responsePx.PagoPxCuatro.SignoCargoAdicional
	entityPago.IvaCargoAdicionalPlanCuota = responsePx.PagoPxCuatro.IvaCargoAdicionalPlanCuota
	entityPago.SignoIvaCargoAdicional = responsePx.PagoPxCuatro.SignoIvaCargoAdicional
	entityPago.CargoAdicionalMovipag = responsePx.PagoPxCuatro.CargoAdicionalMovipag
	entityPago.SignoCargiMovipag = responsePx.PagoPxCuatro.SignoCargiMovipag
	entityPago.IvaCargoAdicionalMovipag = responsePx.PagoPxCuatro.IvaCargoAdicionalMovipag
	entityPago.SignoIvaCargoAdicionalMovipag = responsePx.PagoPxCuatro.SignoIvaCargoAdicionalMovipag
	entityPago.RetencionSello = responsePx.PagoPxCuatro.RetencionSello
	entityPago.SignoRetencionSello = responsePx.PagoPxCuatro.SignoRetencionSello
	entityPago.ProvinciaRetencionSello = responsePx.PagoPxCuatro.ProvinciaRetencionSello
	for _, valuePagoDetalle := range responsePx.PagoPxDos {
		entityPago.Pagostrdos = append(entityPago.Pagostrdos, entities.Prismatrdospago{
			FechaPresentacion:          valuePagoDetalle.FechaPresentacion,
			TipoRegistro:               valuePagoDetalle.TipoRegistro,
			Moneda:                     valuePagoDetalle.Moneda,
			LiquidacionNro:             valuePagoDetalle.LiquidacionNro,
			FechaPago:                  valuePagoDetalle.FechaPago,
			LiquidacionTipo:            valuePagoDetalle.LiquidacionTipo,
			ImporteBruto:               valuePagoDetalle.ImporteBruto,
			SignoImporteBruto:          valuePagoDetalle.SignoImporteBruto,
			ImporteArancel:             valuePagoDetalle.ImporteArancel,
			SignoImporteArancel:        valuePagoDetalle.SignoImporteArancel,
			ImporteNeto:                valuePagoDetalle.ImporteNeto,
			SignoImporteNeto:           valuePagoDetalle.SignoImporteNeto,
			RetencionEspecialSobreIibb: valuePagoDetalle.RetencionEspecialSobreIibb,
			SignoRetencionEspecial:     valuePagoDetalle.SignoRetencionEspecial,
			RetencionIvaEspecial:       valuePagoDetalle.RetencionIvaEspecial,
			SignoRetencionIvaEspecial:  valuePagoDetalle.SignoRetencionIvaEspecial,
			PercepcionIngresoBruto:     valuePagoDetalle.PercepcionIngresoBruto,
			SignoPercepcionIb:          valuePagoDetalle.SignoPercepcionIb,
			RetencionIvaD1:             valuePagoDetalle.RetencionIvaD1,
			SignoRetencionIva_d1:       valuePagoDetalle.SignoRetencionIva_d1,
			CostoCuotaEmitida:          valuePagoDetalle.CostoCuotaEmitida,
			SignoCostoCuotaEmitida:     valuePagoDetalle.SignoCostoCuotaEmitida,
			RetencionIvaCuota:          valuePagoDetalle.RetencionIvaCuota,
			SignoRetencionIvaCuota:     valuePagoDetalle.SignoRetencionIvaCuota,
			RetencionIva:               valuePagoDetalle.RetencionIva,
			SignoRetencionIva:          valuePagoDetalle.SignoRetencionIva,
			RetencionGanacias:          valuePagoDetalle.RetencionGanacias,
			SignoRetencionGanacias:     valuePagoDetalle.SignoRetencionGanacias,
			RetencionIngresoBruto:      valuePagoDetalle.RetencionIngresoBruto,
			SignoRetencionIngresoBruto: valuePagoDetalle.SignoRetencionIngresoBruto,
		})
	}
	return
}
