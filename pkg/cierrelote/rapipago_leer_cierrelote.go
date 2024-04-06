package cierrelote

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierreloterapipagodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type rapipagoProcesarArchivos struct {
	utilService    util.UtilService
	administracion administracion.Service
}

func NewRPProcesarArchivo(util util.UtilService, administracion administracion.Service) MetodoProcesarArtchivos {
	return &rapipagoProcesarArchivos{
		utilService:    util,
		administracion: administracion,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *rapipagoProcesarArchivos) ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo cierrelotedtos.PrismaLogArchivoResponse) {
	// logs.Info(estadosPagoExterno)

	var estado = true
	var estadoInsert = true
	var ErrorProducido string
	rutaArchivo := strings.Split(archivo.Name(), "/")
	/* recorrer y validar archivo headers, trailer , detalles*/

	// 1 verifico que el archivo no se encuentre en la base de datos
	nombreArchivo := strings.Split(rutaArchivo[len(rutaArchivo)-1], "-")
	archivoExiste, err := cl.administracion.ObtenerArchivoCierreLoteRapipago(nombreArchivo[len(nombreArchivo)-1])

	if err != nil {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_OBTENER_NOMBRE_ARCHIVOS + err.Error()
		logs.Error(ErrorProducido)
		logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
		logs.Info("ObtenerArchivoCierreLoteRapipago")
	} else if archivoExiste {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_ARCHIVO_REPETIDO
		logs.Error(ErrorProducido)
		logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
		logs.Info("archivoExiste")
	} else {
		pagosRP, err := RecorrerArchivoRP(archivo)
		if err != nil {
			estado = false
			estadoInsert = false
			ErrorProducido = ERROR_RECORRER_ARCHIVOS + err.Error()
			logs.Error(ErrorProducido)
			logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
			logs.Info("RecorrerArchivoRP verificar")
		} else {

			pagoRp, err := GenerarListaRP(cl, impuesto, nombreArchivo, pagosRP)
			if err != nil {
				estado = false
				estadoInsert = false
				ErrorProducido = ERROR_RECORRER_ARCHIVOS + err.Error()
				logs.Error(ErrorProducido)
				logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
				logs.Info("GenerarListaRP verificar")
			} else {
				err := clRepository.SaveTransactionPagoRP(pagoRp)
				if err != nil {
					estadoInsert = false
					ErrorProducido = ERROR_REGISTRO_EN_DB + err.Error()
					logs.Error(ErrorProducido)
					logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
					logs.Info("clRepository.SaveTransactionPagoRP verificar")
				}
			}
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

func RecorrerArchivoRP(archivo *os.File) (listaRP []cierreloterapipagodtos.Rapipago, erro error) {
	readScanner := bufio.NewScanner(archivo)

	var header cierreloterapipagodtos.Header
	var trailer cierreloterapipagodtos.Trailler
	var detalles []cierreloterapipagodtos.Detalles

	for readScanner.Scan() {

		/* identificar headers, detalles y trailler */
		c1 := len(readScanner.Text())
		/* 1 armar headers */
		if readScanner.Text()[0:8] == "00000000" && len(readScanner.Text()) == 73 {
			header = cierreloterapipagodtos.Header{
				IdHeader:      readScanner.Text()[0:8],
				NombreEmpresa: readScanner.Text()[8:28],
				FechaProceso:  readScanner.Text()[28:36],
				IdArchivo:     readScanner.Text()[36:56],
				FillerHeader:  readScanner.Text()[56:73],
			}
			erro = header.ValidarHeader()
			if erro != nil {
				logs.Error("error al validar el registro: " + erro.Error())
				return
			}
			/* 2 armar trailer */
		} else if readScanner.Text()[0:8] == "99999999" && len(readScanner.Text()) == 73 {

			trailer = cierreloterapipagodtos.Trailler{
				IdTrailler:    readScanner.Text()[0:8],
				CantDetalles:  readScanner.Text()[8:16],
				ImporteTotal:  readScanner.Text()[16:34],
				FillerTrailer: readScanner.Text()[34:73],
			}
			erro = trailer.ValidarTrailer()
			if erro != nil {
				logs.Error("error al validar el registro: " + erro.Error())
				return
			}
			/* detalles */
		} else if c1 == 81 {
			detallesPago := cierreloterapipagodtos.Detalles{
				FechaCobro:     readScanner.Text()[0:8],
				ImporteCobrado: readScanner.Text()[8:23],
				CodigoBarras:   readScanner.Text()[23:71],
				Clearing:       readScanner.Text()[73:81],
			}
			erro = detallesPago.ValidarDetalle()
			if erro != nil {
				logs.Error("error al validar el registro: " + erro.Error())
				return
			} else {
				detalles = append(detalles, detallesPago)

			}
		} else {
			logs.Error("error al validar el registro: " + ERROR_FORMATO_REGISTRO_RP)
			return
		}

	}

	headerTrailer := cierreloterapipagodtos.HeaderTrailer{
		Header:  header,
		Trailer: trailer,
	}

	listaRP = append(listaRP, cierreloterapipagodtos.Rapipago{
		RapipagoHeader:   headerTrailer,
		RapipagoDetalles: detalles,
	})

	return
}

func GenerarListaRP(cl *rapipagoProcesarArchivos, impuesto administraciondtos.ResponseImpuesto, nombreArchivo []string, listaRP []cierreloterapipagodtos.Rapipago) (rapipagoRegistros []entities.Rapipagocierrelote, erro error) {

	/*
		Autor: Jose Alarcon
		Descripción: Calcular comisión de pago offline : se calcula 2.75 por cada tranasacción con un minimo de 20
		Estos valores se calculan con la tabla channelsarancel y IVA(de la tabla configuiracion)
		Se debe analizar si la fecha de vigencia es menor a la fecha de proceso del archivo
	*/
	// 1 obtener medio de pago offline
	filtroChannel := filtros.ChannelFiltro{
		Channel: "OFFLINE",
	}

	channel, erro := cl.administracion.GetChannelService(filtroChannel)
	if erro != nil && int64(channel.Id) < 0 {
		logs.Info(erro)
		return
	}

	// 2 obtener channels aracel
	filtroChannelArancel := filtros.ChannelArancelFiltro{
		CargarAllMedioPago: true,
		CargarChannel:      true,
		ChannelId:          channel.Id,
	}
	channelArancel, erro := cl.administracion.GetChannelsArancelService(filtroChannelArancel)
	logs.Info(channelArancel)
	if erro != nil {
		return
	}

	// esto permitira obtener el channels arancel del proveedor vigente
	var arancel administraciondtos.ResponseChannelsAranceles
	fecha_proceso := ConvertirFormatoFechaRapipago(listaRP[0].RapipagoHeader.Header.FechaProceso)
	fecha_proceso_archivo, _ := time.Parse("2006-01-02", fecha_proceso)
	for _, ch := range channelArancel.ChannelArancel {
		fecha_arancel, _ := time.Parse("2006-01-02T00:00:00Z", ch.Fechadesde)
		if fecha_proceso_archivo.Equal(fecha_arancel) || fecha_proceso_archivo.After(fecha_arancel) {
			arancel = administraciondtos.ResponseChannelsAranceles{
				Importe:       ch.Importe,
				Importeminimo: ch.Importeminimo,
				Importemaximo: ch.Importemaximo,
			}
		}

	}

	var importeTotalHeader float64
	var cabecera entities.Rapipagocierrelote
	// var detalle []entities.Rapipagocierrelotedetalles
	for _, lista := range listaRP {
		result := ValidarTrailerDetalles(lista)
		if !result {
			erro = errors.New(ERROR_FORMATO_REGISTRO_RP_TRAILER_DETALLES)
			return
		} else {
			validar := commons.NewAlgoritmoVerificacion()
			// fecha_proceso := ConvertirFormatoFechaRapipago(lista.RapipagoHeader.Header.FechaProceso)
			fecha_clearing := ConvertirFormatoFechaRapipagoAcreditacion(lista.RapipagoDetalles[0].Clearing)
			// fecha_proc, _ := time.Parse("2006-01-02T00:00:00Z", fecha_proceso)
			fecha_acreditacion, err := time.Parse("2006-01-02", fecha_clearing)
			if err != nil {
				erro = err
				return nil, errors.New(erro.Error())
			}
			cant_dias, err := validar.CalcularDiasEntreFechas(fecha_proceso, fecha_clearing)
			if err != nil {
				erro = err
				return nil, errors.New(erro.Error())
			}

			cant_registros, err := strconv.ParseInt(lista.RapipagoHeader.Trailer.CantDetalles, 10, 64)
			if err != nil {
				erro = err
				return nil, errors.New(erro.Error())
			}
			importe_total, err := strconv.ParseInt(lista.RapipagoHeader.Trailer.ImporteTotal, 10, 64)
			if err != nil {
				erro = err
				return nil, errors.New(erro.Error())
			}
			cabecera = entities.Rapipagocierrelote{
				NombreArchivo:     nombreArchivo[len(nombreArchivo)-1],
				IdHeader:          strings.TrimSpace(lista.RapipagoHeader.Header.IdHeader),
				FechaProceso:      fecha_proceso,
				NombreEmpresa:     strings.TrimSpace(lista.RapipagoHeader.Header.NombreEmpresa),
				IdArchivo:         strings.TrimSpace(lista.RapipagoHeader.Header.IdArchivo),
				FillerHeader:      strings.TrimSpace(lista.RapipagoHeader.Header.FillerHeader),
				IdTrailer:         strings.TrimSpace(lista.RapipagoHeader.Trailer.IdTrailler),
				CantDetalles:      cant_registros,
				Fechaacreditacion: fecha_acreditacion,
				Cantdias:          cant_dias,
				ImporteTotal:      importe_total,
				ImporteMinimo:     arancel.Importeminimo,
				ImporteMaximo:     arancel.Importemaximo,
				Coeficiente:       arancel.Importe,
			}

			for _, detalleRP := range lista.RapipagoDetalles {
				fecha_cobro := ConvertirFormatoFechaRapipago(detalleRP.FechaCobro)
				importe_cobrado, err := strconv.ParseInt(detalleRP.ImporteCobrado, 10, 64)
				if err != nil {
					erro = err
					return nil, errors.New(erro.Error())
				}

				// format fecha de clearing obtenida de los detalles del archivo
				fecha_clearing_detalle := ConvertirFormatoFechaRapipagoAcreditacion(detalleRP.Clearing)

				/*
					Descripción:  Se calcula comision cobrada por rapipago
					se suman los importes de los detalles y se acumula en el header
				*/
				importeTotalCalculado := calcularComision(cl, arancel, impuesto, importe_cobrado)

				rp := entities.Rapipagocierrelotedetalles{
					FechaCobro:       fecha_cobro,
					ImporteCobrado:   importe_cobrado,
					CodigoBarras:     strings.TrimSpace(detalleRP.CodigoBarras),
					ImporteCalculado: importeTotalCalculado,
					Clearing:         fecha_clearing_detalle,
				}
				importeTotalHeader += rp.ImporteCalculado
				cabecera.ImporteTotalCalculado = importeTotalHeader
				cabecera.RapipagoDetalle = append(cabecera.RapipagoDetalle, &rp)
			}
			/* armar lista de registros */
		}
	}
	// cabecera.ImporteTotalCalculado = importeTotalCalculado

	rapipagoRegistros = append(rapipagoRegistros, cabecera)

	return
}

func ValidarTrailerDetalles(listaRP cierreloterapipagodtos.Rapipago) bool {

	importeTotal, err := strconv.ParseInt(listaRP.RapipagoHeader.Trailer.ImporteTotal, 10, 64)
	if err != nil {
		return false
	}
	registrosTotal, err := strconv.ParseInt(listaRP.RapipagoHeader.Trailer.CantDetalles, 10, 64)
	if err != nil {
		return false
	}
	// recorrer detalles sumar total y cantidad de registros
	total := 0
	totalRegistros := 0
	for _, detalle := range listaRP.RapipagoDetalles {
		totalRegistros++
		importe, err := strconv.ParseInt(detalle.ImporteCobrado, 10, 64)
		if err != nil {
			return false
		}
		total += int(importe)
	}

	if total == int(importeTotal) && totalRegistros == int(registrosTotal) {
		return true
	} else {
		return false
	}
}

func ConvertirFormatoFechaRapipago(fecha string) string {
	total := len(fecha)
	resultado := fecha[0:4] + "-" + fecha[4:6] + "-" + fecha[6:total]
	return resultado
}

func ConvertirFormatoFechaRapipagoAcreditacion(fecha string) string {
	// dia-mes-año
	total := len(fecha)
	resultado := fecha[4:total] + "-" + fecha[2:4] + "-" + fecha[0:2]
	return resultado
}

func calcularComision(cl *rapipagoProcesarArchivos, arancel administraciondtos.ResponseChannelsAranceles, impuesto administraciondtos.ResponseImpuesto, valor int64) (calculado float64) {

	/*
		Autor: Jose Alarcon
		Descripción: Calcular comisión de pago offline : se calcula 2.75 por cada tranasacción con un minimo de 20
		Estos valores se calculan con la tabla channelsarancel , IVA
	*/

	valorCalculo := cl.utilService.ToFixed((float64(valor) / 100), 2)

	// calcular comision y verificar el minimo
	comision := valorCalculo * arancel.Importe

	// verificar si el valor es menor al minimo(20)
	if comision < arancel.Importeminimo {
		iva := arancel.Importeminimo * impuesto.Porcentaje
		descuento := arancel.Importeminimo + iva
		valorfinalcalculado := valorCalculo - descuento
		calculado = cl.utilService.ToFixed(valorfinalcalculado*100, 2)
	} else if comision > arancel.Importeminimo {
		iva := comision * impuesto.Porcentaje
		descuento := comision + iva
		valorfinalcalculado := valorCalculo - descuento
		calculado = cl.utilService.ToFixed(valorfinalcalculado*100, 2)
	}
	return
}
