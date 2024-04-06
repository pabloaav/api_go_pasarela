package cierrelotedtos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseMovimientoMx struct {
	MovimientosMXTotales  ResponseMovimientoMxTotal
	MovimientosMxDetalles []ResponseMovimientoMxDetalle
}

func (responseMx *ResponseMovimientoMx) EntityToDtosMx(movimientoMxEntity []entities.Prismamxtotalesmovimiento) (listaMovimientoMx []ResponseMovimientoMx, erro error) {

	for _, valueMxcabecera := range movimientoMxEntity {
		var MxCabecera ResponseMovimientoMxTotal
		MxCabecera.Id = valueMxcabecera.ID
		MxCabecera.Empresa = valueMxcabecera.Empresa
		fechaPresentacion, err := time.Parse("020106", valueMxcabecera.Fechapres)
		if err != nil {
			erro = errors.New(CONVERTIR_FECHA_ERROR)
			return
		}
		MxCabecera.FechaPresentacion = fechaPresentacion
		MxCabecera.TipoRegistro = valueMxcabecera.Tiporeg
		MxCabecera.ComercioNro = valueMxcabecera.Numcom
		MxCabecera.EstablecimientoNro = valueMxcabecera.Numest
		MxCabecera.Codop = valueMxcabecera.Codop
		MxCabecera.TipoAplicacion = valueMxcabecera.Tipoaplic
		fechaPago, err := time.Parse("020106", valueMxcabecera.Fechapago)
		if err != nil {
			erro = errors.New(CONVERTIR_FECHA_ERROR)
			return
		}
		MxCabecera.FechaPago = fechaPago
		importe, err := strconv.ParseInt(valueMxcabecera.ImporteTotal, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		MxCabecera.ImporteTotal = entities.Monto(importe)
		MxCabecera.SignoImporteTotal = valueMxcabecera.SignoImporteTotal

		var MxDetalle []ResponseMovimientoMxDetalle
		for _, valuedetalle := range valueMxcabecera.MovimientosDetalle {
			fechapresentacion, err := time.Parse("020106", valuedetalle.Fechapresentacion)
			if err != nil {
				erro = errors.New(CONVERTIR_FECHA_ERROR)
				return
			}
			nroLote, err := strconv.ParseInt(valuedetalle.Lote, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}
			fechaOrigen, err := time.Parse("020106", valuedetalle.ForigCompra)
			if err != nil {
				erro = errors.New(CONVERTIR_FECHA_ERROR)
				return
			}
			fechaPag, err := time.Parse("020106", valuedetalle.Fechapag)
			if err != nil {
				erro = errors.New(CONVERTIR_FECHA_ERROR)
				return
			}
			nroCupon, err := strconv.ParseInt(valuedetalle.Numcomp, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}
			importe, err := strconv.ParseInt(valuedetalle.Importe, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}
			//ferificar si viene vacio
			var fechaOrigenAjuste string
			if len(strings.TrimSpace(valuedetalle.FechapagAjuLqe)) > 0 {
				fecha, err := time.Parse("020106", valuedetalle.FechapagAjuLqe)
				if err != nil {
					erro = errors.New(CONVERTIR_FECHA_ERROR)
					return
				}
				fechaOrigenAjuste = fecha.Format("020106")

			} else {
				fechaOrigenAjuste = ""
			}
			porcentajeArancel, err := strconv.ParseFloat(valuedetalle.PorcdtoArancel, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_FLOAT_ERROR)
				return
			}
			arancel, err := strconv.ParseInt(valuedetalle.Arancel, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}
			tnaCf, err := strconv.ParseInt(valuedetalle.TnaCf, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}
			importeCostoFinanciero, err := strconv.ParseInt(valuedetalle.ImporteCostoFin, 10, 64)
			if err != nil {
				erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
				return
			}

			///enums
			strEstadoTransaccion := obtenerEnumTransaccionAceptada(valuedetalle.RecAcep)
			if strEstadoTransaccion == "" {
				erro = errors.New(ERROR_CONVERTIR_ENNUM + " - enum transacciones")
				return
			}
			strMoneda := obtenerEnumMoneda(valuedetalle.Moneda)
			strIdCf := obtenerEnumIdCf(strings.TrimSpace(valuedetalle.IdCf))
			strIvaExcento := obtenerEnumIvaExcento(valuedetalle.CfExentoIva)
			strEstablecimiento := obtenerEnumBanderaEstablecimiento(valuedetalle.BanderaEst)
			temporalResponseMx := ResponseMovimientoMxDetalle{
				Id: valuedetalle.ID,

				PrismamovimientototalesId:   valuedetalle.PrismamxtotalesmoviminetosId,
				Fechapresentacion:           fechapresentacion,
				Tiporeg:                     valuedetalle.Tiporeg,
				PrismaoperacionsId:          valuedetalle.Codop,
				Tipoaplic:                   valuedetalle.Tipoaplic,
				Lote:                        nroLote,
				Numtar:                      valuedetalle.Numtar,
				FechaOrigenCompra:           fechaOrigen,
				FechaPago:                   fechaPag,
				NroCupon:                    nroCupon,
				Importe:                     entities.Monto(importe),
				SignoImporte:                valuedetalle.Signo,
				NroAutorizacion:             valuedetalle.Numaut,
				NroCuota:                    valuedetalle.Numcuot,
				PlanCuota:                   valuedetalle.Plancuot,
				RecAcep:                     strEstadoTransaccion,
				RechazoPrincipalId:          valuedetalle.RechPrint,
				RechazoSecundarioId:         valuedetalle.RechSecun,
				NroLiquidacion:              valuedetalle.Nroliq,
				ContracargoOrigen:           valuedetalle.CcoOrigen,
				PrismavisacontracargosId:    valuedetalle.CcoMotivo,
				Moneda:                      strMoneda,
				IdCf:                        strIdCf,
				CfExentoIva:                 strIvaExcento,
				FechaPagoOrigenAjuste:       fechaOrigenAjuste,
				PrismamotivosajustesId:      valuedetalle.CodMotivoAjuLqe,
				PorcentDescArancel:          porcentajeArancel,
				Arancel:                     entities.Monto(arancel),
				SignoArancel:                valuedetalle.SignoArancel,
				TnaCf:                       entities.Monto(tnaCf),
				ImporteCostoFinanciero:      entities.Monto(importeCostoFinanciero),
				SignoImporteCostoFinanciero: valuedetalle.SigImporteCostoFinanciero,
				BanderaEstablecimiento:      strEstablecimiento,
				PrismamastercontracargosId:  valuedetalle.CcoMotivoMc,
				NroTarjetaXl:                valuedetalle.NumtarXl,
				NroAutorizacionXl:           valuedetalle.NumautXl,
			}
			MxDetalle = append(MxDetalle, temporalResponseMx)
		}
		listaMovimientoMx = append(listaMovimientoMx, ResponseMovimientoMx{MxCabecera, MxDetalle})
	}
	return
}

const (
	CONVERTIR_FECHA_ERROR    = "error al convertir fecha"
	CONVERTIR_A_ENTERO_ERROR = "error al convertir a númenro entero"
	CONVERTIR_A_FLOAT_ERROR  = "error al convertir a número flotante"
	ERROR_CONVERTIR_ENNUM    = "el codigo ennum no es valido"
)

type EnumsEstructura struct {
	Id     string
	Nombre string
}

func obtenerEnumTransaccionAceptada(id string) string {
	transaccionAceptadaArray := make([]EnumsEstructura, 0, 3)
	transaccionAceptadaArray = append(transaccionAceptadaArray,
		EnumsEstructura{Id: "0", Nombre: "CREDITO"},
		EnumsEstructura{Id: "2", Nombre: "DEBITO"},
		EnumsEstructura{Id: "1", Nombre: "TRANSACCION RECHAZADA"})
	for _, value := range transaccionAceptadaArray {
		if value.Id == id {
			return value.Nombre
		}
	}
	return ""
}

func obtenerEnumMoneda(id string) string {
	monedaArray := make([]EnumsEstructura, 0, 3)
	monedaArray = append(monedaArray,
		EnumsEstructura{Id: "A", Nombre: "peso"},
		EnumsEstructura{Id: "D", Nombre: "dolar"},
		EnumsEstructura{Id: "P", Nombre: "patacon"})
	for _, value := range monedaArray {
		if value.Id == id {
			return value.Nombre
		}
	}
	return "no especifica valor"
}

func obtenerEnumIdCf(id string) string {
	idCfArray := make([]EnumsEstructura, 0, 3)
	idCfArray = append(idCfArray,
		EnumsEstructura{Id: "E", Nombre: "OPERACIÓN EXCENTA DE IVA"},
		EnumsEstructura{Id: "I", Nombre: "APLICA IVA AL 21%"},
		//EnumsEstructura{Id: "", Nombre: "REDUCCIÓN DE IVA 10,5%"}
	)
	for _, value := range idCfArray {
		if value.Id == id {
			return value.Nombre
		}
	}
	return "REDUCCIÓN DE IVA 10,5%"
}

func obtenerEnumIvaExcento(id string) string {
	ivaExcentoArray := make([]EnumsEstructura, 0, 2)
	ivaExcentoArray = append(ivaExcentoArray,
		EnumsEstructura{Id: "E", Nombre: "CF Sin IVA"},
		EnumsEstructura{Id: "I", Nombre: "CF Sin Reduccion de. IVA"})
	for _, value := range ivaExcentoArray {
		if value.Id == id {
			return value.Nombre
		}
	}
	return "no especifica valor"
}

func obtenerEnumBanderaEstablecimiento(id string) string {
	establecimientosArray := make([]EnumsEstructura, 0, 5)
	establecimientosArray = append(establecimientosArray,
		EnumsEstructura{Id: "0", Nombre: "visa"},
		EnumsEstructura{Id: "900", Nombre: "visa"},
		EnumsEstructura{Id: "391", Nombre: "cabal"},
		EnumsEstructura{Id: "100", Nombre: "master"},
		EnumsEstructura{Id: "147", Nombre: "amex"})
	for _, value := range establecimientosArray {
		if value.Id == id {
			return value.Nombre
		}
	}
	return "no especifica valor"
}

func (responseMx *ResponseMovimientoMx) ToEntity() (entityMovimiento entities.Prismamovimientototale, erro error) {

	entityMovimiento.Empresa = responseMx.MovimientosMXTotales.Empresa
	entityMovimiento.FechaPresentacion = responseMx.MovimientosMXTotales.FechaPresentacion
	entityMovimiento.TipoRegistro = responseMx.MovimientosMXTotales.TipoRegistro
	entityMovimiento.ComercioNro = responseMx.MovimientosMXTotales.ComercioNro
	entityMovimiento.EstablecimientoNro = responseMx.MovimientosMXTotales.EstablecimientoNro
	entityMovimiento.Codop = responseMx.MovimientosMXTotales.Codop
	entityMovimiento.TipoAplicacion = responseMx.MovimientosMXTotales.TipoAplicacion
	entityMovimiento.FechaPago = responseMx.MovimientosMXTotales.FechaPago
	entityMovimiento.ImporteTotal = responseMx.MovimientosMXTotales.ImporteTotal
	entityMovimiento.SignoImporteTotal = responseMx.MovimientosMXTotales.SignoImporteTotal
	for _, valueDetalle := range responseMx.MovimientosMxDetalles {
		idMasterContracargo, err := strconv.ParseUint(valueDetalle.PrismamastercontracargosId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		idOperacion, err := strconv.ParseUint(valueDetalle.PrismaoperacionsId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		idPrimerRechazo, err := strconv.ParseUint(valueDetalle.RechazoPrincipalId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		idSegundoRechazo, err := strconv.ParseUint(valueDetalle.RechazoSecundarioId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		idVisacontracargo, err := strconv.ParseUint(valueDetalle.PrismavisacontracargosId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		idMotivoAjuste, err := strconv.ParseUint(valueDetalle.PrismamotivosajustesId, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		// ofuscacionTarjeta := OfuscacionTarjeta(valueDetalle.Numtar, 4)
		// ofuscacionTarjetaxL := OfuscacionTarjeta(valueDetalle.NroTarjetaXl, 4)
		planCuota, err := strconv.ParseInt(valueDetalle.PlanCuota, 10, 64)
		if err != nil {
			erro = errors.New(CONVERTIR_A_ENTERO_ERROR)
			return
		}
		entityMovimiento.DetalleMovimientos = append(entityMovimiento.DetalleMovimientos, entities.Prismamovimientodetalle{
			PrismamastercontracargosId:  uint(idMasterContracargo),
			PrismaoperacionsId:          uint(idOperacion),
			RechazoPrincipalId:          uint(idPrimerRechazo),
			RechazoSecundarioId:         uint(idSegundoRechazo),
			PrismavisacontracargosId:    uint(idVisacontracargo),
			PrismamotivosajustesId:      uint(idMotivoAjuste),
			TipoRegistro:                valueDetalle.Tiporeg,
			TipoAplicacion:              valueDetalle.Tipoaplic,
			Lote:                        valueDetalle.Lote,
			NroTarjeta:                  valueDetalle.Numtar,
			FechaOrigenCompra:           valueDetalle.FechaOrigenCompra,
			FechaPago:                   valueDetalle.FechaPago,
			NroCupon:                    valueDetalle.NroCupon,
			Importe:                     valueDetalle.Importe,
			SignoImporte:                valueDetalle.SignoImporte,
			NroAutorizacion:             valueDetalle.NroAutorizacion,
			NroCuota:                    valueDetalle.NroCuota,
			PlanCuota:                   planCuota,
			RecAcep:                     valueDetalle.RecAcep,
			NroLiquidacion:              valueDetalle.NroLiquidacion,
			ContracargoOrigen:           valueDetalle.ContracargoOrigen,
			Moneda:                      valueDetalle.Moneda,
			IdCf:                        valueDetalle.IdCf,
			CfExentoIva:                 valueDetalle.CfExentoIva,
			FechaPagoOrigenAjuste:       valueDetalle.FechaPagoOrigenAjuste,
			PorcentDescArancel:          valueDetalle.PorcentDescArancel,
			Arancel:                     valueDetalle.Arancel,
			SignoArancel:                valueDetalle.SignoArancel,
			TnaCf:                       valueDetalle.TnaCf,
			ImporteCostoFinanciero:      valueDetalle.ImporteCostoFinanciero,
			SignoImporteCostoFinanciero: valueDetalle.SignoImporteCostoFinanciero,
			BanderaEstablecimiento:      valueDetalle.BanderaEstablecimiento,
			NroTarjetaXl:                valueDetalle.NroTarjetaXl,
			NroAutorizacionXl:           valueDetalle.NroAutorizacionXl,
		})
	}

	return
}

func OfuscacionTarjeta(cadenaStr string, valueStr int) (resultStr string) {
	var buildStr string
	nroTC, _ := strconv.ParseInt(cadenaStr, 10, 64)
	temporalStr := fmt.Sprintf("%d", nroTC)

	leftStr := temporalStr[0:valueStr]

	totalStr := len(temporalStr) - valueStr
	rightStr := temporalStr[totalStr:]

	totalStr = len(temporalStr) - (valueStr * 2)
	for i := 0; i < totalStr; i++ {
		buildStr += "0"
	}
	resultStr = leftStr + buildStr + rightStr
	return
}
