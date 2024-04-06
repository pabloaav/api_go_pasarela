package cierrelote

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type movimientoMxProcesarArchivo struct {
	utilService util.UtilService
}

func 
NewMXProcesarArchivo(util util.UtilService) MetodoProcesarArtchivos {
	return &movimientoMxProcesarArchivo{utilService: util}
}

func (m *movimientoMxProcesarArchivo) ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo cierrelotedtos.PrismaLogArchivoResponse) {
	//logs.Info(estadosPagoExterno)
	var estado = true
	var estadoInsert = true
	var ErrorProducido string
	rutaArchivo := strings.Split(archivo.Name(), "/")
	movimientoMxRegistros, err := RecorrerArchivoMx(archivo)
	if err != nil {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_FORMATO_REGISTRO_MOVIMIENTO + err.Error()
		logs.Error(ErrorProducido)
		logs.Error("no se realizo  insercion de de los pagos px")
	} else {
		movimientosMx := GenerarListasMxDetalleTotales(rutaArchivo[len(rutaArchivo)-1], movimientoMxRegistros)
		err := clRepository.SaveTransactionMovimientoMx(movimientosMx)
		if err != nil {
			estadoInsert = false
			ErrorProducido = ERROR_REGISTRO_EN_DB + err.Error()
			logs.Error(ErrorProducido)
		}
	}
	listaLogArchivo = cierrelotedtos.PrismaLogArchivoResponse{
		NombreArchivo:  rutaArchivo[len(rutaArchivo)-1],
		ArchivoLeido:   estado,
		ArchivoMovido:  false,
		LoteInsert:     estadoInsert,
		ErrorProducido: ErrorProducido,
	}
	return
}

func RecorrerArchivoMx(archivo *os.File) (movimientoMxRegistros []cierrelotedtos.MovimientoMxRegistros, erro error) {
	var registroDescripcion cierrelotedtos.EstructuraRegistros
	var movimientoMxDetalleTemp []cierrelotedtos.MovimientoMxDetalleRegistro
	readScanner := bufio.NewScanner(archivo)
	for readScanner.Scan() {

		if (readScanner.Text()[11:13] == "01" || readScanner.Text()[11:13] == "02") && len(readScanner.Text()) == 750 {
			registroMxDtalle, err := convertirStrToRegistroMxDetalle(readScanner.Text(), registroDescripcion)
			if err != nil {
				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_MIVIMIENTO_DETALLE, err.Error())
				erro = errors.New(mensaje)
				return
			}
			movimientoMxDetalleTemp = append(movimientoMxDetalleTemp, registroMxDtalle)
			//movimientoMxRegistros.MovimientoMxDetalle = append(movimientoMxRegistros.MovimientoMxDetalle, registroMxDtalle)
		} else if readScanner.Text()[11:13] == "05" && len(readScanner.Text()) == 750 {
			registroMxTotales, err := convertirStrToRegistroMxTotales(readScanner.Text(), registroDescripcion)
			if err != nil {
				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_MIVIMIENTO_TOTALES, err.Error())
				erro = errors.New(mensaje)
				return
			}
			movimientoMxRegistros = append(movimientoMxRegistros, cierrelotedtos.MovimientoMxRegistros{
				MovimientoMxDetalle: movimientoMxDetalleTemp,
				MovimientoMxTotales: registroMxTotales,
			})
			movimientoMxDetalleTemp = nil
			//movimientoMxRegistros.MovimientoMxTotales = append(movimientoMxRegistros.MovimientoMxTotales, registroMxTotales)
		} else {
			erro = errors.New(ERROR_FORMATO_REGISTRO_MOVIMIENTO)
			return
		}

	}
	return
}

func convertirStrToRegistroMxDetalle(strLine string, registroDescripcion cierrelotedtos.EstructuraRegistros) (registroMxDetalle cierrelotedtos.MovimientoMxDetalleRegistro, erro error) {
	descripcionRegistroMx := registroDescripcion.MxDetalleDescripcionRegistro()
	hashTcStr := CodificarTR(strLine[descripcionRegistroMx[12].Desde:descripcionRegistroMx[12].Hasta])
	hashTcXlStr := CodificarTR(strLine[descripcionRegistroMx[57].Desde:descripcionRegistroMx[57].Hasta])
	registroMxDetalle = cierrelotedtos.MovimientoMxDetalleRegistro{
		Empresa:                   strLine[descripcionRegistroMx[0].Desde:descripcionRegistroMx[0].Hasta],
		Fechapresentacion:         strLine[descripcionRegistroMx[1].Desde:descripcionRegistroMx[1].Hasta],
		Tiporeg:                   strLine[descripcionRegistroMx[2].Desde:descripcionRegistroMx[2].Hasta],
		Numcom:                    strLine[descripcionRegistroMx[3].Desde:descripcionRegistroMx[3].Hasta],
		Numest:                    strLine[descripcionRegistroMx[4].Desde:descripcionRegistroMx[4].Hasta],
		Codop:                     strLine[descripcionRegistroMx[5].Desde:descripcionRegistroMx[5].Hasta],
		Tipoaplic:                 strLine[descripcionRegistroMx[6].Desde:descripcionRegistroMx[6].Hasta],
		Lote:                      strLine[descripcionRegistroMx[7].Desde:descripcionRegistroMx[7].Hasta],
		Codbco:                    strLine[descripcionRegistroMx[8].Desde:descripcionRegistroMx[8].Hasta],
		Codcasa:                   strLine[descripcionRegistroMx[9].Desde:descripcionRegistroMx[9].Hasta],
		Bcoest:                    strLine[descripcionRegistroMx[10].Desde:descripcionRegistroMx[10].Hasta],
		Bcocasa:                   strLine[descripcionRegistroMx[11].Desde:descripcionRegistroMx[11].Hasta],
		Numtar:                    hashTcStr, //strLine[descripcionRegistroMx[12].Desde:descripcionRegistroMx[12].Hasta],
		ForigCompra:               strLine[descripcionRegistroMx[13].Desde:descripcionRegistroMx[13].Hasta],
		Fechapag:                  strLine[descripcionRegistroMx[14].Desde:descripcionRegistroMx[14].Hasta],
		Numcomp:                   strLine[descripcionRegistroMx[15].Desde:descripcionRegistroMx[15].Hasta],
		Importe:                   strLine[descripcionRegistroMx[16].Desde:descripcionRegistroMx[16].Hasta],
		Signo:                     strLine[descripcionRegistroMx[17].Desde:descripcionRegistroMx[17].Hasta],
		Numaut:                    strLine[descripcionRegistroMx[18].Desde:descripcionRegistroMx[18].Hasta],
		Numcuot:                   strLine[descripcionRegistroMx[19].Desde:descripcionRegistroMx[19].Hasta],
		Plancuot:                  strLine[descripcionRegistroMx[20].Desde:descripcionRegistroMx[20].Hasta],
		RecAcep:                   strLine[descripcionRegistroMx[21].Desde:descripcionRegistroMx[21].Hasta],
		RechPrint:                 strLine[descripcionRegistroMx[22].Desde:descripcionRegistroMx[22].Hasta],
		RechSecun:                 strLine[descripcionRegistroMx[23].Desde:descripcionRegistroMx[23].Hasta],
		ImpPlan:                   strLine[descripcionRegistroMx[24].Desde:descripcionRegistroMx[24].Hasta],
		Signo1:                    strLine[descripcionRegistroMx[25].Desde:descripcionRegistroMx[25].Hasta],
		McaPex:                    strLine[descripcionRegistroMx[26].Desde:descripcionRegistroMx[26].Hasta],
		Nroliq:                    strLine[descripcionRegistroMx[27].Desde:descripcionRegistroMx[27].Hasta],
		CcoOrigen:                 strLine[descripcionRegistroMx[28].Desde:descripcionRegistroMx[28].Hasta],
		CcoMotivo:                 strLine[descripcionRegistroMx[29].Desde:descripcionRegistroMx[29].Hasta],
		IdCargoliq:                strLine[descripcionRegistroMx[30].Desde:descripcionRegistroMx[30].Hasta],
		Moneda:                    strLine[descripcionRegistroMx[31].Desde:descripcionRegistroMx[31].Hasta],
		PromoBonifUsu:             strLine[descripcionRegistroMx[32].Desde:descripcionRegistroMx[32].Hasta],
		PromoBonifEst:             strLine[descripcionRegistroMx[33].Desde:descripcionRegistroMx[33].Hasta],
		IdPromo:                   strLine[descripcionRegistroMx[34].Desde:descripcionRegistroMx[34].Hasta],
		MovImporig:                strLine[descripcionRegistroMx[35].Desde:descripcionRegistroMx[35].Hasta],
		SgImporig:                 strLine[descripcionRegistroMx[36].Desde:descripcionRegistroMx[36].Hasta],
		IdCf:                      strLine[descripcionRegistroMx[37].Desde:descripcionRegistroMx[37].Hasta],
		CfExentoIva:               strLine[descripcionRegistroMx[38].Desde:descripcionRegistroMx[38].Hasta],
		Dealer:                    strLine[descripcionRegistroMx[39].Desde:descripcionRegistroMx[39].Hasta],
		CuitEst:                   strLine[descripcionRegistroMx[40].Desde:descripcionRegistroMx[40].Hasta],
		FechapagAjuLqe:            strLine[descripcionRegistroMx[41].Desde:descripcionRegistroMx[41].Hasta],
		CodMotivoAjuLqe:           strLine[descripcionRegistroMx[42].Desde:descripcionRegistroMx[42].Hasta],
		IdentifNroFactura:         strLine[descripcionRegistroMx[43].Desde:descripcionRegistroMx[43].Hasta],
		PorcdtoArancel:            strLine[descripcionRegistroMx[44].Desde:descripcionRegistroMx[44].Hasta],
		Arancel:                   strLine[descripcionRegistroMx[45].Desde:descripcionRegistroMx[45].Hasta],
		SignoArancel:              strLine[descripcionRegistroMx[46].Desde:descripcionRegistroMx[46].Hasta],
		TnaCf:                     strLine[descripcionRegistroMx[47].Desde:descripcionRegistroMx[47].Hasta],
		ImporteCostoFin:           strLine[descripcionRegistroMx[48].Desde:descripcionRegistroMx[48].Hasta],
		SigImporteCostoFinanciero: strLine[descripcionRegistroMx[49].Desde:descripcionRegistroMx[49].Hasta],
		IdTx:                      strLine[descripcionRegistroMx[50].Desde:descripcionRegistroMx[50].Hasta],
		Agencia:                   strLine[descripcionRegistroMx[51].Desde:descripcionRegistroMx[51].Hasta],
		TipoPlan:                  strLine[descripcionRegistroMx[52].Desde:descripcionRegistroMx[52].Hasta],
		BanderaEst:                strLine[descripcionRegistroMx[53].Desde:descripcionRegistroMx[53].Hasta],
		Subcodigo:                 strLine[descripcionRegistroMx[54].Desde:descripcionRegistroMx[54].Hasta],
		Filler:                    strLine[descripcionRegistroMx[55].Desde:descripcionRegistroMx[55].Hasta],
		CcoMotivoMc:               strLine[descripcionRegistroMx[56].Desde:descripcionRegistroMx[56].Hasta],
		NumtarXl:                  hashTcXlStr, //strLine[descripcionRegistroMx[57].Desde:descripcionRegistroMx[57].Hasta],
		NumautXl:                  strLine[descripcionRegistroMx[58].Desde:descripcionRegistroMx[58].Hasta],
		Filler1:                   strLine[descripcionRegistroMx[59].Desde:descripcionRegistroMx[59].Hasta],
		IdFinRegistro:             strLine[descripcionRegistroMx[60].Desde:descripcionRegistroMx[60].Hasta],
	}
	err := registroMxDetalle.ValidarMxDetalle(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		erro = err
		return
	}
	return
}

func convertirStrToRegistroMxTotales(strLine string, registroDescripcion cierrelotedtos.EstructuraRegistros) (registroMxTotales cierrelotedtos.MovimientoMxTotalesRegistro, erro error) {
	descripcionRegistroMx := registroDescripcion.MxTotalesDescripcionRegistro()
	registroMxTotales = cierrelotedtos.MovimientoMxTotalesRegistro{
		Empresa:           strLine[descripcionRegistroMx[0].Desde:descripcionRegistroMx[0].Hasta],
		Fechapres:         strLine[descripcionRegistroMx[1].Desde:descripcionRegistroMx[1].Hasta],
		Tiporeg:           strLine[descripcionRegistroMx[2].Desde:descripcionRegistroMx[2].Hasta],
		Numcom:            strLine[descripcionRegistroMx[3].Desde:descripcionRegistroMx[3].Hasta],
		Numest:            strLine[descripcionRegistroMx[4].Desde:descripcionRegistroMx[4].Hasta],
		Codop:             strLine[descripcionRegistroMx[5].Desde:descripcionRegistroMx[5].Hasta],
		Tipoaplic:         strLine[descripcionRegistroMx[6].Desde:descripcionRegistroMx[6].Hasta],
		Filler:            strLine[descripcionRegistroMx[7].Desde:descripcionRegistroMx[7].Hasta],
		Fechapago:         strLine[descripcionRegistroMx[8].Desde:descripcionRegistroMx[8].Hasta],
		Libre:             strLine[descripcionRegistroMx[9].Desde:descripcionRegistroMx[9].Hasta],
		ImporteTotal:      strLine[descripcionRegistroMx[10].Desde:descripcionRegistroMx[10].Hasta],
		SignoImporteTotal: strLine[descripcionRegistroMx[11].Desde:descripcionRegistroMx[11].Hasta],
		Filler1:           strLine[descripcionRegistroMx[12].Desde:descripcionRegistroMx[12].Hasta],
		McaPex:            strLine[descripcionRegistroMx[13].Desde:descripcionRegistroMx[13].Hasta],
		Filler2:           strLine[descripcionRegistroMx[14].Desde:descripcionRegistroMx[14].Hasta],
		Aster:             strLine[descripcionRegistroMx[15].Desde:descripcionRegistroMx[15].Hasta],
	}
	err := registroMxTotales.ValidarMxTotales(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		erro = err
		return
	}
	return
}

func GenerarListasMxDetalleTotales(nombreArchivo string, movimientoMxRegistros []cierrelotedtos.MovimientoMxRegistros) (movimientosMx []entities.Prismamxtotalesmovimiento) {

	for _, valueRegistros := range movimientoMxRegistros {
		var movimientoMxDetalleEntities []entities.Prismamxdetallemovimiento
		entityMxTotales := valueRegistros.MovimientoMxTotales.MxTotalesToEntities(nombreArchivo)
		for _, valueMxDetalle := range valueRegistros.MovimientoMxDetalle {
			if valueRegistros.MovimientoMxTotales.Numest == valueMxDetalle.Numest {
				entityMxDetalle := valueMxDetalle.MxDetalleToEntities()
				movimientoMxDetalleEntities = append(movimientoMxDetalleEntities, entityMxDetalle)
			}
		}
		entityMxTotales.MovimientosDetalle = movimientoMxDetalleEntities
		movimientosMx = append(movimientosMx, entityMxTotales)
	}
	return
}
