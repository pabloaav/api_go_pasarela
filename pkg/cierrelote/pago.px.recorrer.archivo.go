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

type pagoPxProcesarArchivos struct {
	utilService util.UtilService
}

func NewPXProcesarArchivo(util util.UtilService) MetodoProcesarArtchivos {
	return &pagoPxProcesarArchivos{utilService: util}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *pagoPxProcesarArchivos) ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo cierrelotedtos.PrismaLogArchivoResponse) {
	//logs.Info(estadosPagoExterno)
	var estado = true
	var estadoInsert = true
	var ErrorProducido string
	rutaArchivo := strings.Split(archivo.Name(), "/")
	pagoPxRegistro, err := RecorrerArchivoPx(archivo) //pagoPxRegistros,
	if err != nil {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_RECORRER_ARCHIVOS + err.Error()
		logs.Error(ErrorProducido)
		logs.Error("no se realizo la insercion de de los pagos px")
	} else {
		pagoPx := GenerarListaPxDosCuatro(rutaArchivo[len(rutaArchivo)-1], pagoPxRegistro) // pagoPxRegistros,
		err := clRepository.SaveTransactionPagoPx(pagoPx)
		if err != nil {
			estadoInsert = false
			ErrorProducido = ERROR_REGISTRO_EN_DB + err.Error()
			logs.Error(ErrorProducido)
		}
	}
	listaLogArchivo = cierrelotedtos.PrismaLogArchivoResponse{
		NombreArchivo:  rutaArchivo[len(rutaArchivo)-1], //archivo.Name(),
		ArchivoLeido:   estado,
		ArchivoMovido:  false,
		LoteInsert:     estadoInsert,
		ErrorProducido: ErrorProducido,
	}
	return
}

//pagoPxRegistros cierrelotedtos.PagoPxRegistros,
func RecorrerArchivoPx(archivo *os.File) (PagoPxRegistros []cierrelotedtos.PagoPxRegistros, erro error) {
	var registroDescripcion cierrelotedtos.EstructuraRegistros
	var PagoPxDetalleTemp []cierrelotedtos.PrismaPxDosRegistro
	readScanner := bufio.NewScanner(archivo)

	for readScanner.Scan() {

		if readScanner.Text()[11:13] == "02" && len(readScanner.Text()) == 700 {

			registroPxDos, err := convertirStrToRegistroPxDos(readScanner.Text(), registroDescripcion)
			if err != nil {
				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_02, err.Error())
				erro = errors.New(mensaje)
				return
			}

			//pagoPxRegistros.PagoPxDetalle = append(pagoPxRegistros.PagoPxDetalle, registroPxDos)
			PagoPxDetalleTemp = append(PagoPxDetalleTemp, registroPxDos)

		} else if readScanner.Text()[11:13] == "04" && len(readScanner.Text()) == 700 {

			registroPxCuatro, err := convertirStrToRegistroPxCuatro(readScanner.Text(), registroDescripcion)
			if err != nil {
				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_04, err.Error())
				erro = errors.New(mensaje)
				return
			}
			//pagoPxRegistros.PagoPxCabecera = append(pagoPxRegistros.PagoPxCabecera, registroPxCuatro)

			PagoPxRegistros = append(PagoPxRegistros, cierrelotedtos.PagoPxRegistros{
				PagoPxCabecera: registroPxCuatro,
				PagoPxDetalle:  PagoPxDetalleTemp,
			})
			PagoPxDetalleTemp = nil

		} else {
			erro = errors.New(ERROR_FORMATO_REGISTRO)
			return
		}

	}
	return
}

// func RecorrerArchivoPx(archivo *os.File) (pagoPxRegistros cierrelotedtos.PagoPxRegistros, erro error) {
// 	var registroDescripcion cierrelotedtos.EstructuraRegistros
// 	readScanner := bufio.NewScanner(archivo)
// 	for readScanner.Scan() {
// 		if readScanner.Text()[11:13] == "02" && len(readScanner.Text()) == 700 {
// 			registroPxDos, err := convertirStrToRegistroPxDos(readScanner.Text(), registroDescripcion)
// 			if err != nil {
// 				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_02, err.Error())
// 				erro = errors.New(mensaje)
// 				return
// 			}
// 			pagoPxRegistros.PagoPxDetalle = append(pagoPxRegistros.PagoPxDetalle, registroPxDos)
// 		} else if readScanner.Text()[11:13] == "04" && len(readScanner.Text()) == 700 {
// 			registroPxCuatro, err := convertirStrToRegistroPxCuatro(readScanner.Text(), registroDescripcion)
// 			if err != nil {
// 				mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_04, err.Error())
// 				erro = errors.New(mensaje)
// 				return
// 			}
// 			pagoPxRegistros.PagoPxCabecera = append(pagoPxRegistros.PagoPxCabecera, registroPxCuatro)

// 		} else {
// 			erro = errors.New(ERROR_FORMATO_REGISTRO)
// 			return
// 		}
// 	}
// 	return
// }

func convertirStrToRegistroPxDos(strLine string, registroDescripcion cierrelotedtos.EstructuraRegistros) (registroPxDos cierrelotedtos.PrismaPxDosRegistro, erro error) {
	descripcionRegistrPxDos := registroDescripcion.PxDosDescripcionRegistro()
	registroPxDos = cierrelotedtos.PrismaPxDosRegistro{
		Eclq02llEmpresa:     strLine[descripcionRegistrPxDos[0].Desde:descripcionRegistrPxDos[0].Hasta],
		Eclq02llFpres:       strLine[descripcionRegistrPxDos[1].Desde:descripcionRegistrPxDos[1].Hasta],
		Eclq02llTiporeg:     strLine[descripcionRegistrPxDos[2].Desde:descripcionRegistrPxDos[2].Hasta],
		Eclq02llMoneda:      strLine[descripcionRegistrPxDos[3].Desde:descripcionRegistrPxDos[3].Hasta],
		Eclq02llNumcom:      strLine[descripcionRegistrPxDos[4].Desde:descripcionRegistrPxDos[4].Hasta],
		Eclq02llNumest:      strLine[descripcionRegistrPxDos[5].Desde:descripcionRegistrPxDos[5].Hasta],
		Eclq02llNroliq:      strLine[descripcionRegistrPxDos[6].Desde:descripcionRegistrPxDos[6].Hasta],
		Eclq02llFpag:        strLine[descripcionRegistrPxDos[7].Desde:descripcionRegistrPxDos[7].Hasta],
		Eclq02llTipoliq:     strLine[descripcionRegistrPxDos[8].Desde:descripcionRegistrPxDos[8].Hasta],
		Eclq02llImpbruto:    strLine[descripcionRegistrPxDos[9].Desde:descripcionRegistrPxDos[9].Hasta],
		Eclq02llSigno_1:     strLine[descripcionRegistrPxDos[10].Desde:descripcionRegistrPxDos[10].Hasta],
		Eclq02llImpret:      strLine[descripcionRegistrPxDos[11].Desde:descripcionRegistrPxDos[11].Hasta],
		Eclq02llSigno_2:     strLine[descripcionRegistrPxDos[12].Desde:descripcionRegistrPxDos[12].Hasta],
		Eclq02llImpneto:     strLine[descripcionRegistrPxDos[13].Desde:descripcionRegistrPxDos[13].Hasta],
		Eclq02llSigno_3:     strLine[descripcionRegistrPxDos[14].Desde:descripcionRegistrPxDos[14].Hasta],
		Eclq02llRetesp:      strLine[descripcionRegistrPxDos[15].Desde:descripcionRegistrPxDos[15].Hasta],
		Eclq02llSigno_4:     strLine[descripcionRegistrPxDos[16].Desde:descripcionRegistrPxDos[16].Hasta],
		Eclq02llRetivaEsp:   strLine[descripcionRegistrPxDos[17].Desde:descripcionRegistrPxDos[17].Hasta],
		Eclq02llSigno_5:     strLine[descripcionRegistrPxDos[18].Desde:descripcionRegistrPxDos[18].Hasta],
		Eclq02llPercepBa:    strLine[descripcionRegistrPxDos[19].Desde:descripcionRegistrPxDos[19].Hasta],
		Eclq02llSigno_6:     strLine[descripcionRegistrPxDos[20].Desde:descripcionRegistrPxDos[20].Hasta],
		Eclq02llRetivaD1:    strLine[descripcionRegistrPxDos[21].Desde:descripcionRegistrPxDos[21].Hasta],
		Eclq02llSigno_7:     strLine[descripcionRegistrPxDos[22].Desde:descripcionRegistrPxDos[22].Hasta],
		Filler1:             strLine[descripcionRegistrPxDos[23].Desde:descripcionRegistrPxDos[23].Hasta],
		Filler2:             strLine[descripcionRegistrPxDos[24].Desde:descripcionRegistrPxDos[24].Hasta],
		Eclq02llCargoPex:    strLine[descripcionRegistrPxDos[25].Desde:descripcionRegistrPxDos[25].Hasta],
		Eclq02llSigno_9:     strLine[descripcionRegistrPxDos[26].Desde:descripcionRegistrPxDos[26].Hasta],
		Eclq02llRetivaPex1:  strLine[descripcionRegistrPxDos[27].Desde:descripcionRegistrPxDos[27].Hasta],
		Eclq02llSigno_10:    strLine[descripcionRegistrPxDos[28].Desde:descripcionRegistrPxDos[28].Hasta],
		Filler3:             strLine[descripcionRegistrPxDos[29].Desde:descripcionRegistrPxDos[29].Hasta],
		Filler4:             strLine[descripcionRegistrPxDos[30].Desde:descripcionRegistrPxDos[30].Hasta],
		Eclq02llCostoCuoemi: strLine[descripcionRegistrPxDos[31].Desde:descripcionRegistrPxDos[31].Hasta],
		Eclq02llSigno_12:    strLine[descripcionRegistrPxDos[32].Desde:descripcionRegistrPxDos[32].Hasta],
		Eclq02llRetivaCuo1:  strLine[descripcionRegistrPxDos[33].Desde:descripcionRegistrPxDos[33].Hasta],
		Eclq02llSigno_13:    strLine[descripcionRegistrPxDos[34].Desde:descripcionRegistrPxDos[34].Hasta],
		Filler5:             strLine[descripcionRegistrPxDos[35].Desde:descripcionRegistrPxDos[35].Hasta],
		Filler6:             strLine[descripcionRegistrPxDos[36].Desde:descripcionRegistrPxDos[36].Hasta],
		Eclq02llImpServ:     strLine[descripcionRegistrPxDos[37].Desde:descripcionRegistrPxDos[37].Hasta],
		Eclq02llSigno_15:    strLine[descripcionRegistrPxDos[38].Desde:descripcionRegistrPxDos[38].Hasta],
		Eclq02llIva1Xlj:     strLine[descripcionRegistrPxDos[39].Desde:descripcionRegistrPxDos[39].Hasta],
		Eclq02llSigno_16:    strLine[descripcionRegistrPxDos[40].Desde:descripcionRegistrPxDos[40].Hasta],
		Filler7:             strLine[descripcionRegistrPxDos[41].Desde:descripcionRegistrPxDos[41].Hasta],
		Filler8:             strLine[descripcionRegistrPxDos[42].Desde:descripcionRegistrPxDos[42].Hasta],
		Eclq02llCargoEdcE:   strLine[descripcionRegistrPxDos[43].Desde:descripcionRegistrPxDos[43].Hasta],
		Eclq02llSigno_18:    strLine[descripcionRegistrPxDos[44].Desde:descripcionRegistrPxDos[44].Hasta],
		Eclq02llIva1EdcE:    strLine[descripcionRegistrPxDos[45].Desde:descripcionRegistrPxDos[45].Hasta],
		Eclq02llSigno_19:    strLine[descripcionRegistrPxDos[46].Desde:descripcionRegistrPxDos[46].Hasta],
		Filler9:             strLine[descripcionRegistrPxDos[47].Desde:descripcionRegistrPxDos[47].Hasta],
		Filler10:            strLine[descripcionRegistrPxDos[48].Desde:descripcionRegistrPxDos[48].Hasta],
		Eclq02llCargoEdcB:   strLine[descripcionRegistrPxDos[49].Desde:descripcionRegistrPxDos[49].Hasta],
		Eclq02llSigno_21:    strLine[descripcionRegistrPxDos[50].Desde:descripcionRegistrPxDos[50].Hasta],
		Eclq02llIva1EdcB:    strLine[descripcionRegistrPxDos[51].Desde:descripcionRegistrPxDos[51].Hasta],
		Eclq02llSigno_22:    strLine[descripcionRegistrPxDos[52].Desde:descripcionRegistrPxDos[52].Hasta],
		Filler11:            strLine[descripcionRegistrPxDos[53].Desde:descripcionRegistrPxDos[53].Hasta],
		Filler12:            strLine[descripcionRegistrPxDos[54].Desde:descripcionRegistrPxDos[54].Hasta],
		Eclq02llCargoCitE:   strLine[descripcionRegistrPxDos[55].Desde:descripcionRegistrPxDos[55].Hasta],
		Eclq02llSigno_24:    strLine[descripcionRegistrPxDos[56].Desde:descripcionRegistrPxDos[56].Hasta],
		Eclq02llIva1CitE:    strLine[descripcionRegistrPxDos[57].Desde:descripcionRegistrPxDos[57].Hasta],
		Eclq02llllSigno_25:  strLine[descripcionRegistrPxDos[58].Desde:descripcionRegistrPxDos[58].Hasta],
		Filler13:            strLine[descripcionRegistrPxDos[59].Desde:descripcionRegistrPxDos[59].Hasta],
		Filler14:            strLine[descripcionRegistrPxDos[60].Desde:descripcionRegistrPxDos[60].Hasta],
		Eclq02llCargoCitB:   strLine[descripcionRegistrPxDos[61].Desde:descripcionRegistrPxDos[61].Hasta],
		Eclq02llSigno_27:    strLine[descripcionRegistrPxDos[62].Desde:descripcionRegistrPxDos[62].Hasta],
		Eclq02llIva1CitB:    strLine[descripcionRegistrPxDos[63].Desde:descripcionRegistrPxDos[63].Hasta],
		Eclq02llSigno_28:    strLine[descripcionRegistrPxDos[64].Desde:descripcionRegistrPxDos[64].Hasta],
		Filler15:            strLine[descripcionRegistrPxDos[65].Desde:descripcionRegistrPxDos[65].Hasta],
		Filler16:            strLine[descripcionRegistrPxDos[66].Desde:descripcionRegistrPxDos[66].Hasta],
		Eclq02llRetIva:      strLine[descripcionRegistrPxDos[67].Desde:descripcionRegistrPxDos[67].Hasta],
		Eclq02llSigno_30:    strLine[descripcionRegistrPxDos[68].Desde:descripcionRegistrPxDos[68].Hasta],
		Eclq02llRetGcias:    strLine[descripcionRegistrPxDos[69].Desde:descripcionRegistrPxDos[69].Hasta],
		Eclq02llSigno_31:    strLine[descripcionRegistrPxDos[70].Desde:descripcionRegistrPxDos[70].Hasta],
		Eclq02llRetIngbru:   strLine[descripcionRegistrPxDos[71].Desde:descripcionRegistrPxDos[71].Hasta],
		Eclq02llSigno_32:    strLine[descripcionRegistrPxDos[72].Desde:descripcionRegistrPxDos[72].Hasta],
		Filler17:            strLine[descripcionRegistrPxDos[73].Desde:descripcionRegistrPxDos[73].Hasta],
		Filler18:            strLine[descripcionRegistrPxDos[74].Desde:descripcionRegistrPxDos[74].Hasta],
		Filler19:            strLine[descripcionRegistrPxDos[75].Desde:descripcionRegistrPxDos[75].Hasta],
		Eclq02llAster:       strLine[descripcionRegistrPxDos[76].Desde:descripcionRegistrPxDos[76].Hasta],
	}
	err := registroPxDos.ValidarPxDos(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		erro = err
		return
	}
	return
}

func convertirStrToRegistroPxCuatro(strLine string, registroDescripcion cierrelotedtos.EstructuraRegistros) (registroPxCuatro cierrelotedtos.PrismaPxCuatroRegistro, erro error) {
	descripcionRegistrPxCuatro := registroDescripcion.PxCuatroDescripcionRegistro()
	registroPxCuatro = cierrelotedtos.PrismaPxCuatroRegistro{
		Eclq02llEmpresa_04:         strLine[descripcionRegistrPxCuatro[0].Desde:descripcionRegistrPxCuatro[0].Hasta],
		Eclq02llFpres_04:           strLine[descripcionRegistrPxCuatro[1].Desde:descripcionRegistrPxCuatro[1].Hasta],
		Eclq02llTiporeg_04:         strLine[descripcionRegistrPxCuatro[2].Desde:descripcionRegistrPxCuatro[2].Hasta],
		Eclq02llMoneda_04:          strLine[descripcionRegistrPxCuatro[3].Desde:descripcionRegistrPxCuatro[3].Hasta],
		Eclq02llNumcom_04:          strLine[descripcionRegistrPxCuatro[4].Desde:descripcionRegistrPxCuatro[4].Hasta],
		Eclq02llNumest_04:          strLine[descripcionRegistrPxCuatro[5].Desde:descripcionRegistrPxCuatro[5].Hasta],
		Eclq02llNroliq_04:          strLine[descripcionRegistrPxCuatro[6].Desde:descripcionRegistrPxCuatro[6].Hasta],
		Eclq02llFpag_04:            strLine[descripcionRegistrPxCuatro[7].Desde:descripcionRegistrPxCuatro[7].Hasta],
		Eclq02llTipoliq_04:         strLine[descripcionRegistrPxCuatro[8].Desde:descripcionRegistrPxCuatro[8].Hasta],
		Eclq02llCasacta:            strLine[descripcionRegistrPxCuatro[9].Desde:descripcionRegistrPxCuatro[9].Hasta],
		Eclq02llTipcta:             strLine[descripcionRegistrPxCuatro[10].Desde:descripcionRegistrPxCuatro[10].Hasta],
		Eclq02llCtabco:             strLine[descripcionRegistrPxCuatro[11].Desde:descripcionRegistrPxCuatro[11].Hasta],
		Eclq02llCfExentoIva:        strLine[descripcionRegistrPxCuatro[12].Desde:descripcionRegistrPxCuatro[12].Hasta],
		Eclq02llSigno_04_1:         strLine[descripcionRegistrPxCuatro[13].Desde:descripcionRegistrPxCuatro[13].Hasta],
		Eclq02llLey25063:           strLine[descripcionRegistrPxCuatro[14].Desde:descripcionRegistrPxCuatro[14].Hasta],
		Eclq02llSigno_04_2:         strLine[descripcionRegistrPxCuatro[15].Desde:descripcionRegistrPxCuatro[15].Hasta],
		Eclq02llAliIngbru:          strLine[descripcionRegistrPxCuatro[16].Desde:descripcionRegistrPxCuatro[16].Hasta],
		Eclq02llDtoCampania:        strLine[descripcionRegistrPxCuatro[17].Desde:descripcionRegistrPxCuatro[17].Hasta],
		Eclq02llSigno_04_3:         strLine[descripcionRegistrPxCuatro[18].Desde:descripcionRegistrPxCuatro[18].Hasta],
		Eclq02llIva1DtoCampania:    strLine[descripcionRegistrPxCuatro[19].Desde:descripcionRegistrPxCuatro[19].Hasta],
		Eclq02llSigno_04_4:         strLine[descripcionRegistrPxCuatro[20].Desde:descripcionRegistrPxCuatro[20].Hasta],
		Eclq02llRetIngbru2:         strLine[descripcionRegistrPxCuatro[21].Desde:descripcionRegistrPxCuatro[21].Hasta],
		Eclq02llSigno_04_5:         strLine[descripcionRegistrPxCuatro[22].Desde:descripcionRegistrPxCuatro[22].Hasta],
		Eclq02llAliIngbru2:         strLine[descripcionRegistrPxCuatro[23].Desde:descripcionRegistrPxCuatro[23].Hasta],
		Filler1:                    strLine[descripcionRegistrPxCuatro[24].Desde:descripcionRegistrPxCuatro[24].Hasta],
		Filler2:                    strLine[descripcionRegistrPxCuatro[25].Desde:descripcionRegistrPxCuatro[25].Hasta],
		Filler3:                    strLine[descripcionRegistrPxCuatro[26].Desde:descripcionRegistrPxCuatro[26].Hasta],
		Filler4:                    strLine[descripcionRegistrPxCuatro[27].Desde:descripcionRegistrPxCuatro[27].Hasta],
		Filler5:                    strLine[descripcionRegistrPxCuatro[28].Desde:descripcionRegistrPxCuatro[28].Hasta],
		Eclq02llTasaPex:            strLine[descripcionRegistrPxCuatro[29].Desde:descripcionRegistrPxCuatro[29].Hasta],
		Eclq02llCargoXliq:          strLine[descripcionRegistrPxCuatro[30].Desde:descripcionRegistrPxCuatro[30].Hasta],
		Eclq02llSigno_04_8:         strLine[descripcionRegistrPxCuatro[31].Desde:descripcionRegistrPxCuatro[31].Hasta],
		Eclq02llIva1CargoXliq:      strLine[descripcionRegistrPxCuatro[32].Desde:descripcionRegistrPxCuatro[32].Hasta],
		Eclq02llSigno_04_9:         strLine[descripcionRegistrPxCuatro[33].Desde:descripcionRegistrPxCuatro[33].Hasta],
		Eclq02llDealer:             strLine[descripcionRegistrPxCuatro[34].Desde:descripcionRegistrPxCuatro[34].Hasta],
		Eclq02llImpDbCr:            strLine[descripcionRegistrPxCuatro[35].Desde:descripcionRegistrPxCuatro[35].Hasta],
		Eclq02llSigno_04_10:        strLine[descripcionRegistrPxCuatro[36].Desde:descripcionRegistrPxCuatro[36].Hasta],
		Eclq02llCfNoReduceIva:      strLine[descripcionRegistrPxCuatro[37].Desde:descripcionRegistrPxCuatro[37].Hasta],
		Eclq02llSigno_04_11:        strLine[descripcionRegistrPxCuatro[38].Desde:descripcionRegistrPxCuatro[38].Hasta],
		Eclq02llPercepIbAgip:       strLine[descripcionRegistrPxCuatro[39].Desde:descripcionRegistrPxCuatro[39].Hasta],
		Eclq02llSigno_04_12:        strLine[descripcionRegistrPxCuatro[40].Desde:descripcionRegistrPxCuatro[40].Hasta],
		Eclq02llAlicPercepIbAgip:   strLine[descripcionRegistrPxCuatro[41].Desde:descripcionRegistrPxCuatro[41].Hasta],
		Eclq02llRetenIbAgip:        strLine[descripcionRegistrPxCuatro[42].Desde:descripcionRegistrPxCuatro[42].Hasta],
		Eclq02llSigno_04_13:        strLine[descripcionRegistrPxCuatro[43].Desde:descripcionRegistrPxCuatro[43].Hasta],
		Eclq02llAlicRetenIbAgip:    strLine[descripcionRegistrPxCuatro[44].Desde:descripcionRegistrPxCuatro[44].Hasta],
		Eclq02llSubtotRetivaRg3130: strLine[descripcionRegistrPxCuatro[45].Desde:descripcionRegistrPxCuatro[45].Hasta],
		Eclq02llSigno_04_14:        strLine[descripcionRegistrPxCuatro[46].Desde:descripcionRegistrPxCuatro[46].Hasta],
		Eclq02llProvIngbru:         strLine[descripcionRegistrPxCuatro[47].Desde:descripcionRegistrPxCuatro[47].Hasta],
		Eclq02llAdicPlancuo:        strLine[descripcionRegistrPxCuatro[48].Desde:descripcionRegistrPxCuatro[48].Hasta],
		Eclq02llSigno_04_15:        strLine[descripcionRegistrPxCuatro[49].Desde:descripcionRegistrPxCuatro[49].Hasta],
		Eclq02llIva1AdPlancuo:      strLine[descripcionRegistrPxCuatro[50].Desde:descripcionRegistrPxCuatro[50].Hasta],
		Eclq02llSigno_04_16:        strLine[descripcionRegistrPxCuatro[51].Desde:descripcionRegistrPxCuatro[51].Hasta],
		Eclq02llAdic_opinter:       strLine[descripcionRegistrPxCuatro[52].Desde:descripcionRegistrPxCuatro[52].Hasta],
		Eclq02llSigno_04_17:        strLine[descripcionRegistrPxCuatro[53].Desde:descripcionRegistrPxCuatro[53].Hasta],
		Eclq02llIva1Ad_opinter:     strLine[descripcionRegistrPxCuatro[54].Desde:descripcionRegistrPxCuatro[54].Hasta],
		Eclq02llSigno_04_18:        strLine[descripcionRegistrPxCuatro[55].Desde:descripcionRegistrPxCuatro[55].Hasta],
		Eclq02llAdicAltacom:        strLine[descripcionRegistrPxCuatro[56].Desde:descripcionRegistrPxCuatro[56].Hasta],
		Eclq02llSigno_04_19:        strLine[descripcionRegistrPxCuatro[57].Desde:descripcionRegistrPxCuatro[57].Hasta],
		Eclq02llIva1AdAltacom:      strLine[descripcionRegistrPxCuatro[58].Desde:descripcionRegistrPxCuatro[58].Hasta],
		Eclq02llSigno_04_20:        strLine[descripcionRegistrPxCuatro[59].Desde:descripcionRegistrPxCuatro[59].Hasta],
		Eclq02llAdicCupmanu:        strLine[descripcionRegistrPxCuatro[60].Desde:descripcionRegistrPxCuatro[60].Hasta],
		Eclq02llSigno_04_21:        strLine[descripcionRegistrPxCuatro[61].Desde:descripcionRegistrPxCuatro[61].Hasta],
		Eclq02llIva1AdCupmanu:      strLine[descripcionRegistrPxCuatro[62].Desde:descripcionRegistrPxCuatro[62].Hasta],
		Eclq02llSigno_04_22:        strLine[descripcionRegistrPxCuatro[63].Desde:descripcionRegistrPxCuatro[63].Hasta],
		Eclq02llAdicAltacomBco:     strLine[descripcionRegistrPxCuatro[64].Desde:descripcionRegistrPxCuatro[64].Hasta],
		Eclq02llSigno_04_23:        strLine[descripcionRegistrPxCuatro[65].Desde:descripcionRegistrPxCuatro[65].Hasta],
		Eclq02llIva1AdAltacomBco:   strLine[descripcionRegistrPxCuatro[66].Desde:descripcionRegistrPxCuatro[66].Hasta],
		Eclq02llSigno_04_24:        strLine[descripcionRegistrPxCuatro[67].Desde:descripcionRegistrPxCuatro[67].Hasta],
		Filler6:                    strLine[descripcionRegistrPxCuatro[68].Desde:descripcionRegistrPxCuatro[68].Hasta],
		Filler7:                    strLine[descripcionRegistrPxCuatro[69].Desde:descripcionRegistrPxCuatro[69].Hasta],
		Filler8:                    strLine[descripcionRegistrPxCuatro[70].Desde:descripcionRegistrPxCuatro[70].Hasta],
		Filler9:                    strLine[descripcionRegistrPxCuatro[71].Desde:descripcionRegistrPxCuatro[71].Hasta],
		Eclq02llAdicMovypag:        strLine[descripcionRegistrPxCuatro[72].Desde:descripcionRegistrPxCuatro[72].Hasta],
		Eclq02llSigno_04_27:        strLine[descripcionRegistrPxCuatro[73].Desde:descripcionRegistrPxCuatro[73].Hasta],
		Eclq02llIva1AdicMovypag:    strLine[descripcionRegistrPxCuatro[74].Desde:descripcionRegistrPxCuatro[74].Hasta],
		Eclq02llSigno_04_28:        strLine[descripcionRegistrPxCuatro[75].Desde:descripcionRegistrPxCuatro[75].Hasta],
		Eclq02llRetSellos:          strLine[descripcionRegistrPxCuatro[76].Desde:descripcionRegistrPxCuatro[76].Hasta],
		Eclq02llSigno_04_29:        strLine[descripcionRegistrPxCuatro[77].Desde:descripcionRegistrPxCuatro[77].Hasta],
		Eclq02llProvSellos:         strLine[descripcionRegistrPxCuatro[78].Desde:descripcionRegistrPxCuatro[78].Hasta],
		Eclq02llRetIngbru3:         strLine[descripcionRegistrPxCuatro[79].Desde:descripcionRegistrPxCuatro[79].Hasta],
		Eclq02llSigno_04_30:        strLine[descripcionRegistrPxCuatro[80].Desde:descripcionRegistrPxCuatro[80].Hasta],
		Eclq02llAliIngbru3:         strLine[descripcionRegistrPxCuatro[81].Desde:descripcionRegistrPxCuatro[81].Hasta],
		Eclq02llRetIngbru4:         strLine[descripcionRegistrPxCuatro[82].Desde:descripcionRegistrPxCuatro[82].Hasta],
		Eclq02llSigno_04_31:        strLine[descripcionRegistrPxCuatro[83].Desde:descripcionRegistrPxCuatro[83].Hasta],
		Eclq02llAliIngbru4:         strLine[descripcionRegistrPxCuatro[84].Desde:descripcionRegistrPxCuatro[84].Hasta],
		Eclq02llRetIngbru5:         strLine[descripcionRegistrPxCuatro[85].Desde:descripcionRegistrPxCuatro[85].Hasta],
		Eclq02llSigno_04_32:        strLine[descripcionRegistrPxCuatro[86].Desde:descripcionRegistrPxCuatro[86].Hasta],
		Eclq02llAliIngbru5:         strLine[descripcionRegistrPxCuatro[87].Desde:descripcionRegistrPxCuatro[87].Hasta],
		Eclq02llRetIngbru6:         strLine[descripcionRegistrPxCuatro[88].Desde:descripcionRegistrPxCuatro[88].Hasta],
		Eclq02llSigno_04_33:        strLine[descripcionRegistrPxCuatro[89].Desde:descripcionRegistrPxCuatro[89].Hasta],
		Eclq02llAliIngbru6:         strLine[descripcionRegistrPxCuatro[90].Desde:descripcionRegistrPxCuatro[90].Hasta],
		Eclq02llFiller_04_10:       strLine[descripcionRegistrPxCuatro[91].Desde:descripcionRegistrPxCuatro[91].Hasta],
		Eclq02llAster_04_11:        strLine[descripcionRegistrPxCuatro[92].Desde:descripcionRegistrPxCuatro[92].Hasta],
	}
	err := registroPxCuatro.ValidarPxCuatro(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		erro = err
		return
	}
	return
}

// pagoPxRegistros cierrelotedtos.PagoPxRegistros,
func GenerarListaPxDosCuatro(nombreArchivo string, PagoPxRegistros []cierrelotedtos.PagoPxRegistros) (pagoPxCuatroEntities []entities.Prismapxcuatroregistro) {

	for _, valueRegistros := range PagoPxRegistros {
		var pagoPxDosEntities []entities.Prismapxdosregistro
		entityPxCabecera := valueRegistros.PagoPxCabecera.PxCuatroToEntities(nombreArchivo)
		for _, valueRegistrosDetalle := range valueRegistros.PagoPxDetalle {
			if valueRegistros.PagoPxCabecera.Eclq02llNumest_04 == valueRegistrosDetalle.Eclq02llNumest {
				entityPxDetalle := valueRegistrosDetalle.PxDosToEntities()
				pagoPxDosEntities = append(pagoPxDosEntities, entityPxDetalle)
			}
		}
		entityPxCabecera.PxDosRegistros = pagoPxDosEntities
		pagoPxCuatroEntities = append(pagoPxCuatroEntities, entityPxCabecera)

	}

	// for _, valuePagoPxCuatro := range pagoPxRegistros.PagoPxCabecera {
	// 	var pagoPxDosEntities []entities.Prismapxdosregistro
	// 	entityPxCabecera := valuePagoPxCuatro.PxCuatroToEntities(nombreArchivo)
	// 	for _, valuePagoPxDos := range pagoPxRegistros.PagoPxDetalle {
	// 		if valuePagoPxCuatro.Eclq02llNumest_04 == valuePagoPxDos.Eclq02llNumest {
	// 			entityPxDetalle := valuePagoPxDos.PxDosToEntities()
	// 			pagoPxDosEntities = append(pagoPxDosEntities, entityPxDetalle)
	// 		}
	// 	}
	// 	entityPxCabecera.PxDosRegistros = pagoPxDosEntities
	// 	pagoPxCuatroEntities = append(pagoPxCuatroEntities, entityPxCabecera)
	// }

	return
}
