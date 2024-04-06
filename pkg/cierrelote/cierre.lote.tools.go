package cierrelote

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

/*
autor: sergio
fecha: 2022-04-26
descripcion: se lee en contenido del archivo txt para luego recorrer  y convertir a la estructura de datos
*/
func RecorrerArchivo(archivoLote *os.File) (registroDetalle []prismaCierreLote.PrismaRegistroDetalle, erro error) {
	// defer archivoLote.Close()
	// se define un tipo estructura de campos del registro trailer
	var listaCamposTrailer prismaCierreLote.CampoTrailer
	// se define un tipo estructura de campos del registro detalle
	var listaCamposDetalle prismaCierreLote.CampoDetalle
	// defino un tipo estructura registro trailer
	var resultRegistroTrailer prismaCierreLote.PrismaRegistroTrailer
	// defino un tipo estructura registro trailer para luego validar
	var validarRegistroTrailer prismaCierreLote.PrismaRegistroTrailer
	fileScanner := bufio.NewScanner(archivoLote)
	logs.Info("Inicio lectura contenido de archivo")
	for fileScanner.Scan() {
		/*
			dependiendo del tipo de registro se lee la cadena de string y se convierte a la estructura de datos
		*/
		if fileScanner.Text()[0:1] == "T" && len(fileScanner.Text()) == 100 {
			resultReadTrailer, err := RecorrerTrailer(fileScanner.Text(), listaCamposTrailer)
			if err != nil {
				erro = errors.New(err.Error())
				return nil, erro
			}
			resultRegistroTrailer = resultReadTrailer
			validarRegistroTrailer.TipoRegistro = resultRegistroTrailer.TipoRegistro
			validarRegistroTrailer.IdMedioPago = resultRegistroTrailer.IdMedioPago
			validarRegistroTrailer.IdLote = resultRegistroTrailer.IdLote
			validarRegistroTrailer.Filler = resultRegistroTrailer.Filler
		} else if fileScanner.Text()[0:1] == "D" && len(fileScanner.Text()) == 190 {
			resultRegistroDetalle, err := RecorrerDetalle(fileScanner.Text(), listaCamposDetalle)
			if err != nil {
				erro = errors.New(err.Error())
				return nil, erro
			}
			registroDetalle = append(registroDetalle, resultRegistroDetalle)
			if resultRegistroDetalle.TipoOperacion == "C" {
				validarRegistroTrailer.CantidadCompras += 1
				validarRegistroTrailer.MontoCompras += resultRegistroDetalle.Monto

			} else if resultRegistroDetalle.TipoOperacion == "D" {
				validarRegistroTrailer.CantidadDevueltas += 1
				validarRegistroTrailer.MontoDevueltas += resultRegistroDetalle.Monto
			} else {
				validarRegistroTrailer.CantidadAnuladas += 1
				validarRegistroTrailer.MontoAnulacion += resultRegistroDetalle.Monto
			}
		} else {
			erro = errors.New("error en el formato del archivo, registro detalle o trailer estan mal construidos")
			return nil, erro
		}
	}

	validarRegistroTrailer.CantidadRegistros = validarRegistroTrailer.CantidadCompras + validarRegistroTrailer.CantidadDevueltas + validarRegistroTrailer.CantidadAnuladas
	validarRegistroTrailer.MontoCompras = ToFixedCL(validarRegistroTrailer.MontoCompras, 2)
	validarRegistroTrailer.MontoAnulacion = ToFixedCL(validarRegistroTrailer.MontoAnulacion, 2)
	validarRegistroTrailer.MontoDevueltas = ToFixedCL(validarRegistroTrailer.MontoDevueltas, 2)
	if resultRegistroTrailer != validarRegistroTrailer {
		// FIXME crear notificacion
		logs.Error("el registro trailer es diferente al trailer validacion")
		return nil, errors.New("error los trailer no coinciden")
	}
	//fmt.Printf("detalle %v\n", registroDetalle)
	return //registroDetalle, nil
}

func RecorrerDetalle(registroString string, listaCamposDetalle prismaCierreLote.CampoDetalle) (prismaCierreLote.PrismaRegistroDetalle, error) {
	listaDetalle := listaCamposDetalle.DescripcionCampos()
	trcNro := ConvertirEntero(registroString[listaDetalle[3].Desde:listaDetalle[3].Hasta], 10, 64)
	ofuscacionTCStr := CompletarStr(strconv.Itoa(int(trcNro)))
	hashTCStr := CodificarTR(ofuscacionTCStr)
	registroDetalle := prismaCierreLote.PrismaRegistroDetalle{
		TipoRegistro:    registroString[listaDetalle[0].Desde:listaDetalle[0].Hasta],
		IdTransacciones: registroString[listaDetalle[1].Desde:listaDetalle[1].Hasta], //registroString[listaDetalle[14].Desde:listaDetalle[14].Hasta],
		IdMedioPago:     ConvertirEntero(registroString[listaDetalle[2].Desde:listaDetalle[2].Hasta], 10, 64),
		//NroTarjetaCompleto: ConvertirEntero(registroString[listaDetalle[3].Desde:listaDetalle[3].Hasta], 10, 64),
		NroTarjetaCompleto: hashTCStr, //ConvertirEntero(ofuscacionTCStr, 10, 64),
		TipoOperacion:      registroString[listaDetalle[4].Desde:listaDetalle[4].Hasta],
		Fecha:              ConvertirFormatoFecha(registroString[listaDetalle[5].Desde:listaDetalle[5].Hasta]),
		//Fecha:              registroString[listaDetalle[5].Desde:listaDetalle[5].Hasta],
		Monto:       CalcularMonto(ConvertirEntero(registroString[listaDetalle[6].Desde:listaDetalle[6].Hasta], 10, 64)),
		CodAut:      registroString[listaDetalle[7].Desde:listaDetalle[7].Hasta], // ConvertirEntero(registroString[listaDetalle[7].Desde:listaDetalle[7].Hasta], 10, 64),
		NroTicket:   ConvertirEntero(registroString[listaDetalle[8].Desde:listaDetalle[8].Hasta], 10, 64),
		IdSite:      ConvertirEntero(registroString[listaDetalle[9].Desde:listaDetalle[9].Hasta], 10, 64),
		IdLote:      ConvertirEntero(registroString[listaDetalle[10].Desde:listaDetalle[10].Hasta], 10, 64),
		Cuotas:      ConvertirEntero(registroString[listaDetalle[11].Desde:listaDetalle[11].Hasta], 10, 64),
		FechaCierre: ConvertirFormatoFecha(registroString[listaDetalle[12].Desde:listaDetalle[12].Hasta]),
		//FechaCierre:        registroString[listaDetalle[12].Desde:listaDetalle[12].Hasta],
		NroEstablecimiento: ConvertirEntero(registroString[listaDetalle[13].Desde:listaDetalle[13].Hasta], 10, 64),
		IdCliente:          registroString[listaDetalle[14].Desde:listaDetalle[14].Hasta],
		Filler:             registroString[listaDetalle[15].Desde:listaDetalle[15].Hasta],
	}
	// fmt.Println(registroString)
	// fmt.Println("===============================")
	// fmt.Println(registroDetalle)
	// fmt.Println(registroString[93:101])
	err := registroDetalle.Validar()
	// logs.Info(registroDetalle.NroTarjetaCompleto)
	if err != nil {
		return registroDetalle, errors.New(err.Error())
	}
	return registroDetalle, nil
}

func RecorrerTrailer(registroString string, listaCamposTrailer prismaCierreLote.CampoTrailer) (prismaCierreLote.PrismaRegistroTrailer, error) {

	listaTrailer := listaCamposTrailer.DescripcionCampos()
	registroTrailer := prismaCierreLote.PrismaRegistroTrailer{
		TipoRegistro:      registroString[listaTrailer[0].Desde:listaTrailer[0].Hasta],
		CantidadRegistros: ConvertirEntero(registroString[listaTrailer[1].Desde:listaTrailer[1].Hasta], 10, 64),
		IdMedioPago:       ConvertirEntero(registroString[listaTrailer[2].Desde:listaTrailer[2].Hasta], 10, 64),
		IdLote:            ConvertirEntero(registroString[listaTrailer[3].Desde:listaTrailer[3].Hasta], 10, 64),
		CantidadCompras:   ConvertirEntero(registroString[listaTrailer[4].Desde:listaTrailer[4].Hasta], 10, 64),
		MontoCompras:      CalcularMonto(ConvertirEntero(registroString[listaTrailer[5].Desde:listaTrailer[5].Hasta], 10, 64)),
		CantidadDevueltas: ConvertirEntero(registroString[listaTrailer[6].Desde:listaTrailer[6].Hasta], 10, 64),
		MontoDevueltas:    CalcularMonto(ConvertirEntero(registroString[listaTrailer[7].Desde:listaTrailer[7].Hasta], 10, 64)),
		CantidadAnuladas:  ConvertirEntero(registroString[listaTrailer[8].Desde:listaTrailer[8].Hasta], 10, 64),
		MontoAnulacion:    CalcularMonto(ConvertirEntero(registroString[listaTrailer[9].Desde:listaTrailer[9].Hasta], 10, 64)),
		Filler:            registroString[listaTrailer[10].Desde:listaTrailer[10].Hasta],
	}
	err := registroTrailer.Validar()
	if err != nil {
		return registroTrailer, errors.New(err.Error())
	}
	return registroTrailer, nil
}

func ConvertirEntero(valorConvert string, base, bitSize int) int64 {
	resultado, _ := strconv.ParseInt(valorConvert, base, bitSize)
	//resultado, _ := strconv.ParseInt(valorConvert, 10, 64)
	return resultado
}

func CalcularMonto(monto int64) float64 {
	resultado := float64(monto) / float64(100)
	return resultado
}

func ConvertirFormatoFecha(fecha string) string {
	total := len(fecha)
	resultado := fecha[0:2] + "-" + fecha[2:4] + "-" + fecha[4:total]
	return resultado
}

// CrearListaCierreLote recibe una lista de detalle registro del cierre de lote de prisma y
// devuelve una lista de cierre lote para DB telco
func CrearListaCierreLote(listaPagoIntentos []entities.Pagointento, listaEstadoPagos []entities.Pagoestadoexterno, nombreArchivoLote string, listaDetalleCierreLote []prismaCierreLote.PrismaRegistroDetalle) (listCierreLote []entities.Prismacierrelote, err error) {

	var idEstado int64
	var channelArancelId int64
	var disputa bool
	for _, valueCierreLote := range listaDetalleCierreLote {

		if valueCierreLote.TipoOperacion == "C" {
			idEstado = ObtenerEstadoId(strings.ToUpper("Accredited"), listaEstadoPagos)
		}
		if valueCierreLote.TipoOperacion == "A" {
			idEstado = ObtenerEstadoId(strings.ToUpper("Rejected"), listaEstadoPagos)
		}
		if valueCierreLote.TipoOperacion == "D" {
			idEstado = ObtenerEstadoId(strings.ToUpper("Reverted"), listaEstadoPagos)
		}
		fechaOperacion, err := time.Parse("02-01-2006", valueCierreLote.Fecha)
		if err != nil {
			return nil, errors.New(ERROR_CONVERTIR_FECHA)
		}

		fechaCierre, err := time.Parse("02-01-2006", valueCierreLote.FechaCierre)
		if err != nil {
			return nil, errors.New(ERROR_CONVERTIR_FECHA)
		}
		for _, value := range listaPagoIntentos {
			cantidadStringCL := len(valueCierreLote.IdTransacciones)
			uuIDStringPI := value.TransactionID[0:cantidadStringCL]
			if valueCierreLote.IdTransacciones == uuIDStringPI {
				rubroId := value.Pago.PagosTipo.Cuenta.RubrosID
				for _, valueChannelArancel := range value.Mediopagos.Channel.Channelaranceles {
					//"2022-08-01T00:00:00Z"
					fechaVigencia, err := time.Parse("2006-01-02T00:00:00Z", valueChannelArancel.Fechadesde)
					if err != nil {
						return nil, errors.New(ERROR_CONVERTIR_FECHA)
					}

					if valueChannelArancel.RubrosId == int64(rubroId) && (fechaCierre.Equal(fechaVigencia) || fechaCierre.After(fechaVigencia)) {
						switch value.Mediopagos.Channel.Channel {
						case "CREDIT":
							if valueCierreLote.Cuotas == 1 {
								if !valueChannelArancel.Pagocuota && valueChannelArancel.Mediopagoid == 0 {
									channelArancelId = int64(valueChannelArancel.ID)
								}
								if valueChannelArancel.Pagocuota && valueChannelArancel.Mediopagoid != 0 && valueChannelArancel.Mediopagoid == valueCierreLote.IdMedioPago {
									channelArancelId = int64(valueChannelArancel.ID)
								}
							}
							if valueCierreLote.Cuotas > 1 {
								if valueChannelArancel.Pagocuota && valueChannelArancel.Mediopagoid == 0 {
									channelArancelId = int64(valueChannelArancel.ID)
								}
								if valueChannelArancel.Pagocuota && valueChannelArancel.Mediopagoid != 0 && valueChannelArancel.Mediopagoid == valueCierreLote.IdMedioPago {
									channelArancelId = int64(valueChannelArancel.ID)
								}
							}
						case "DEBIT":
							channelArancelId = int64(valueChannelArancel.ID)
						default:
							return nil, errors.New(ERROR_CONVERTIR_FECHA)
						}
					}
				}
				break
			}
		}
		arraysNombre := strings.Split(nombreArchivoLote, "/")
		if channelArancelId == 0 {
			return nil, errors.New(ERROR_OBTENER_ID_CHANNELARANCEL)
		}
		if valueCierreLote.TipoOperacion == "D" {
			disputa = true
		} else {
			disputa = false
		}
		listCierreLote = append(listCierreLote, entities.Prismacierrelote{
			PagoestadoexternosId:       idEstado,
			Tiporegistro:               valueCierreLote.TipoRegistro,
			PagosUuid:                  valueCierreLote.IdTransacciones,
			ExternalmediopagoId:        valueCierreLote.IdMedioPago,
			PrismamovimientodetallesId: 0,
			PrismamovimientodetalleId:  0,
			PrismatrdospagosId:         0,
			PrismapagotrdoId:           0,
			ChannelarancelesId:         channelArancelId,
			ImpuestosId:                1,
			Nrotarjeta:                 valueCierreLote.NroTarjetaCompleto, //fmt.Sprint(valueCierreLote.NroTarjetaCompleto),
			Tipooperacion:              entities.EnumTipoOperacion(valueCierreLote.TipoOperacion),
			Fechaoperacion:             fechaOperacion,
			Monto:                      entities.Monto(valueCierreLote.Monto * 100),
			Montofinal:                 entities.Monto(0), //entities.Monto(montofinal),
			Codigoautorizacion:         valueCierreLote.CodAut,
			Nroticket:                  valueCierreLote.NroTicket,
			SiteID:                     valueCierreLote.IdSite,
			ExternalloteId:             valueCierreLote.IdLote,
			Nrocuota:                   valueCierreLote.Cuotas,
			FechaCierre:                fechaCierre,
			Nroestablecimiento:         valueCierreLote.NroEstablecimiento,
			ExternalclienteID:          valueCierreLote.IdCliente[4:len(valueCierreLote.IdCliente)],
			Nombrearchivolote:          arraysNombre[len(arraysNombre)-1],
			Disputa:                    disputa,
			Cantdias:                   0,
			Reversion:                  false,
			DetallemovimientoId:        0,
			DetallepagoId:              0,
			Descripcioncontracargo:     "",
			ExtbancoreversionId:        0,
			Conciliado:                 false,
			Estadomovimiento:           false,
			Descripcionbanco:           "",
		})
	}
	return listCierreLote, nil
}

func ObtenerEstadoId(estado string, listaEstadoPagos []entities.Pagoestadoexterno) (idEstado int64) {
	for _, valorEstado := range listaEstadoPagos {
		if string(valorEstado.PagoEstados.Estado) == estado {
			idEstado = int64(valorEstado.ID)
			break
		}
	}
	return idEstado
}

func ArmarNotificacion(archivo prismaCierreLote.PrismaLogArchivoResponse) (notificacion entities.Notificacione) {
	byteArchivo, _ := json.Marshal(archivo)
	notificacion.Descripcion = string(byteArchivo)
	notificacion.Tipo = "CierreLote"
	return notificacion
}

func ArmarNotificacionCierreLote(logerror prismaCierreLote.PrismaLogProcesocierreLote) (notificacion entities.Notificacione) {
	bytelogerror, _ := json.Marshal(logerror)
	notificacion.Descripcion = string(bytelogerror)
	notificacion.Tipo = "CierreLote"
	return notificacion
}

/*
AbrirArchivo recibe:
la ruta hasta el directorio que contiene los archivos y el nombre del archivo.
retorna  el erchivo abierto o un error
*/
func AbrirArchivo(ruta string, nombreArchivo string) (archivoCLLeido *os.File, erro error) {
	archivoCLLeido, err := os.Open(fmt.Sprintf("%s/%s", ruta, nombreArchivo))
	if err != nil {
		msj := "Error abriendo el archivo de nombre:" + nombreArchivo
		logs.Error(msj)
		erro = errors.New(msj)
		return
	}
	//defer archivoCLLeido.Close()
	return
}

/*
permite leer el contenido de un archivo y retorna
elcontenido del archivo en byte, el nombre del archivo y tipo del archivo, o error
*/
func LeerDatosArchivo(rutadestino string, ruta string, nombreArchivo string) (data []byte, archivonombre string, archivotipo string, erro error) {
	data, erro = ioutil.ReadFile(fmt.Sprintf("%s/%s", ruta, nombreArchivo))

	if erro != nil {
		msj := "error a leer datos del archivo:" + nombreArchivo
		logs.Error(msj)
		erro = errors.New(msj)
		return
	}
	fecha := fmt.Sprintf("%v-%v-%v_%v:%v:%v_", time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	archivo_extension := strings.Split(nombreArchivo, ".")
	archivonombre = fmt.Sprintf("%s/%v%s.%s", rutadestino, fecha, archivo_extension[0], archivo_extension[1])
	archivotipo = archivo_extension[len(archivo_extension)-1]
	return
}

func ToFixedCL(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundCL(num*output)) / output
}

func roundCL(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ExtractStr(cadenaStr string, valueStr int) (resultStr string) {
	var buildStr string
	nroTC := ConvertirEntero(cadenaStr, 10, 64)
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

func CompletarStr(cadenaStr string) (resultStr string) {
	if len(cadenaStr) >= 19 {
		resultStr = cadenaStr
		return
	}
	var strTemporal string
	for i := 0; i < 19-len(cadenaStr); i++ {
		strTemporal += "0"
	}
	resultStr = cadenaStr + strTemporal
	return
}

func CodificarTR(cadenaStr string) (resultStr string) {
	strTemporal := sha256.Sum256([]byte(cadenaStr))
	resultStr = fmt.Sprintf("%x", strTemporal)
	return
}

func SrtExtraerCeros(cadenaStr string, inicio, fin int) (resultStr string) {
	resultStr = cadenaStr[inicio:fin]
	return
}
