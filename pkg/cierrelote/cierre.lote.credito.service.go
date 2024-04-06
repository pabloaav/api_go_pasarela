package cierrelote

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	dtosCL "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"

	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	filtrobanco "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
)

func (s *service) LeerArchivoLoteExterno(ctx context.Context, nombreDirectorio string) (archivos []fs.FileInfo, rutaArchivos string, totalArchivos int, err error) {
	var msjError prismaCierreLote.PrismaLogProcesocierreLote
	/*
		obtento lista de archivos de minio y los pasa a una carpeta temporal en el proyecto
		el servicio necesita del contexto y el nombre del directorio
	*/
	rutaArchivos, err = s.store.GetObject(ctx, nombreDirectorio) //)
	if err != nil {
		errMinioLeerObjetos := errors.New(commons.ERROR_READ_ARCHIVO + "desde MINIO" + err.Error())
		msjError = prismaCierreLote.PrismaLogProcesocierreLote{
			TipoError:     errMinioLeerObjetos.Error(),
			Descripcion:   "error al intentar obtener los archivos de cierre de lote desde minio",
			NombreArchivo: "",
		}
		notificacion := ArmarNotificacionCierreLote(msjError)
		errNotificacion := s.utilService.CreateNotificacionService(notificacion)
		if errNotificacion != nil {
			logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
		}
		return nil, "", 0, errMinioLeerObjetos
	}
	/*
		se obtiene los archivos de la carpeta temporal
	*/
	archivos, err = s.commonsService.LeerDirectorio(rutaArchivos)
	if err != nil {
		errMinioLeerDirectorio := errors.New(commons.ERROR_READ_ARCHIVO + "desde MINIO" + err.Error())
		msjError = prismaCierreLote.PrismaLogProcesocierreLote{
			TipoError:     errMinioLeerDirectorio.Error(),
			Descripcion:   "se produjo un error al intentar obtener los archivos de cierre de lote desde minio",
			NombreArchivo: "",
		}
		notificacion := ArmarNotificacionCierreLote(msjError)
		errNotificacion := s.utilService.CreateNotificacionService(notificacion)
		if errNotificacion != nil {
			logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
			logs.Error(err.Error())
		}

		return nil, "", 0, errMinioLeerDirectorio
	}
	if len(archivos) == 0 {
		err = s.commonsService.BorrarDirectorio(rutaArchivos)
		if err != nil {
			logs.Error("ERROR_AL_BORRAR_DIRECTORIO" + err.Error())
		}
		errNoExisteArchivo := errors.New(ERROR_DIRECTORIO_VACIO)
		msjError = prismaCierreLote.PrismaLogProcesocierreLote{
			TipoError:     errNoExisteArchivo.Error(),
			Descripcion:   "no existe archivo recibido por FTP desde prisma no se pudo continuar con el proceso de cierre de lotes",
			NombreArchivo: "",
		}
		notificacion := ArmarNotificacionCierreLote(msjError)
		errNotificacion := s.utilService.CreateNotificacionService(notificacion)
		if errNotificacion != nil {
			logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
		}
		return nil, "", 0, errNoExisteArchivo
	}
	totalArchivos = len(archivos)
	return
}

/*
	leer cierre de lote: al armar el proceso de background se debe tener encuenta que este servicio
	retorna una un log de los archivos que se produjo error al intentar realizar el cierre de lote
	y de los archivos que no ocacionaron error, ademas tambien retorna un error.
*/
func (s *service) LeerCierreLoteTxt(archivos []fs.FileInfo, rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno) (listaArchivo []prismaCierreLote.PrismaLogArchivoResponse, err error) {

	if len(estadosPagoExterno) == 0 {
		errObtenerEstados := errors.New(ERROR_PAGO_ESTADO_VACIO)
		err = errObtenerEstados
		return nil, err
	}
	// // se define una constante que contiene tamaño de bufer
	// const TamanioBufer = 191 // En bytes
	//var estado, estadoInsert bool
	// se recorre la lista de archivos
	var nombreDirectorio string
	for _, archivo := range archivos {
		if strings.Contains(archivo.Name(), "PX") {
			nombreDirectorio = config.DIR_PX
		} else if strings.Contains(archivo.Name(), "MX") {
			nombreDirectorio = config.DIR_MX
		} else if strings.Contains(archivo.Name(), "RP") {
			nombreDirectorio = config.DIR_RP
		} else {
			nombreDirectorio = config.DIR_CL
		}

		metodoProcesarArtchivos, err := s.factory.GetRecorrerArchivos(nombreDirectorio)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		logs.Info("Leyendo el siguiente archvio: " + archivo.Name())

		// se abre el archivo para su lectura en la ruta especificada

		archivoLote, err := AbrirArchivo(rutaArchivos, archivo.Name())
		defer func() { archivoLote.Close() }()
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		filtro := filtros.ConfiguracionFiltro{
			Nombre: "IMPUESTO_SOBRE_ARANCEL",
		}
		impuestoSobreArancel, erro := s.utilService.GetConfiguracionService(filtro)
		if erro != nil {
			err = errors.New(ERROR_OBTENER_CONFIGURACION_IMPUESTO_ARANCEL + erro.Error())
			logs.Error(err)
		}
		impuestoId, erro := strconv.Atoi(impuestoSobreArancel.Valor)
		if erro != nil {
			err = errors.New(ERROR_CONVERTIR_ENTERO + erro.Error())
			logs.Error(err)
		}
		filtroImpuesto := filtros.ImpuestoFiltro{
			Id: uint(impuestoId),
		}
		impuesto, erro := s.administracionService.GetImpuestosService(filtroImpuesto)
		if erro != nil {
			err = errors.New(ERROR_CONSULTAR_IMPUESTOS + erro.Error())
			logs.Error(err)
		}
		if len(impuesto.Impuestos) <= 0 {
			err = errors.New(ERROR_LISTA_IMPUESTOS_VACION + erro.Error())
			logs.Error(err)
		}
		logs.Info(impuesto)

		listaLogArchivo := metodoProcesarArtchivos.ProcesarArchivos(archivoLote, estadosPagoExterno, impuesto.Impuestos[0], s.repository)
		/*
			por cada uno de los archivos recoridos
			crea una lista de archivos informando:
			    - el nombre del archivo,
				- el valor true o false que indica si el archivo fue leido,
				- el valor true o false que indica si el archivo fue movido,
				- el valor true o false que indica si el contenido del archivo pudo ser procesado con exsito y se inserto en DB,
				- y existe algun error
		*/
		listaArchivo = append(listaArchivo, listaLogArchivo)
	}

	return listaArchivo, nil
}

func (s *service) MoverArchivos(ctx context.Context, rutaArchivos string, listaArchivo []prismaCierreLote.PrismaLogArchivoResponse) (countArchivo int, erro error) {
	/*
		por ultimo se mueven los archvios del directorio temporal
		a un directorio en minio dondo se almacenan los archivos
		de cierre de lote registrado en la DB
	*/
	var rutaDestino string
	for key, archivo := range listaArchivo {
		if archivo.LoteInsert {
			countArchivo++
		}
		if archivo.ArchivoLeido && archivo.LoteInsert {
			rutaDestino = config.DIR_HISTORIAL
		}
		if !archivo.LoteInsert {
			rutaDestino = config.DIR_CLERROR
		}
		/*
			se lee el contenido del archivo y se obtiene su contenido se le pasa:
			- ruta destino
			- ruta origen del archivo
			- nombre del archivo
		*/
		data, filename, filetypo, err := LeerDatosArchivo(rutaDestino, rutaArchivos, archivo.NombreArchivo)
		if err != nil {
			logs.Error(err)
			listaArchivo[key].ArchivoMovido = false
		}
		/*	necesito la data, nombre del archivo y el tipo */
		erro := s.store.PutObject(ctx, data, filename, filetypo)
		if erro != nil {
			logs.Error(ERROR_MOVER_ARCHIVO + erro.Error())
			listaArchivo[key].ArchivoMovido = false
		} else {
			listaArchivo[key].ArchivoMovido = true
		}

		if !archivo.ArchivoLeido || !listaArchivo[key].ArchivoMovido || !archivo.LoteInsert {
			notificacion := ArmarNotificacion(archivo)
			err := s.utilService.CreateNotificacionService(notificacion)
			if err != nil {
				logs.Error(ERROR_AL_CREAR_NOTIFICACION + err.Error())
			}
		}
	}
	return
}

func (s *service) BorrarArchivos(ctx context.Context, nombreDirectorio string, rutaArchivos string, listaArchivo []prismaCierreLote.PrismaLogArchivoResponse) (erro error) {
	for _, archivoValue := range listaArchivo {
		/* antes de borrar el archivo se verifica si:
		el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
		erro = s.commonsService.BorrarArchivo(rutaArchivos, archivoValue.NombreArchivo)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.utilService.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return erro
			}
		}

		key := nombreDirectorio + "/" + archivoValue.NombreArchivo //config.DIR_KEY + "/" + archivoValue.NombreArchivo
		erro = s.store.DeleteObject(ctx, key)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "DeleteObject",
				Mensaje:       erro.Error(),
			}
			erro = s.utilService.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return erro
			}
		}
	}
	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return erro
		}
	}

	return nil
}

func (s *service) ConciliacionBancoPrisma(fechaPagoProcesar string, reversion bool, responseListprismaTrPagos []prismaCierreLote.ResponseTrPagosCabecera) (listaMovimientosBanco []bancodtos.ResponseMovimientosBanco, erro error) {

	var referenciaMovimiento []string
	var movimientosIds []uint
	var ListCierreLote []entities.Prismacierrelote

	for _, valueCabeceraPagos := range responseListprismaTrPagos {
		norEstablecimiento, err := strconv.ParseInt(valueCabeceraPagos.EstablecimientoNro, 10, 64)
		if err != nil {
			erro = errors.New(ERROR_CONVERTIR_VALOR)
			logs.Error(err.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "funcion go ParseInt",
				Mensaje:       erro.Error() + " " + err.Error(),
			}
			err = s.utilService.CreateLogService(log)
			if err != nil {
				logs.Error("error: al crear logs: " + err.Error())
				return
			}
			return
		}
		referenciaMovimiento = append(referenciaMovimiento, strconv.Itoa(int(norEstablecimiento)))
	}
	if len(referenciaMovimiento) < 1 {
		erro = errors.New(ERROR_REFERENCIAS_MOVIMIENTOS)
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "ConciliacionBancoPrisma",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
		return
	}

	/* obtengo lista de movimientos de banco relacionados con prisma */
	filtroBanco := filtrobanco.MovimientosBancoFiltro{
		SubCuenta:      config.COD_SUBCUENTA,
		Tipo:           "prisma",
		TipoMovimiento: referenciaMovimiento,
		Fecha:          fechaPagoProcesar,
		TipoOperacion:  reversion,
	}

	listaMovimientosBanco, err := s.bancoService.BuildCierreLoteApiLinkBancoService(filtroBanco)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_MOVIMIENTOS_BANCO + err.Error())
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BuildCierreLoteApiLinkBancoService",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	if len(listaMovimientosBanco) == 0 {
		requestNotificacion := entities.Notificacione{
			UserId:      0,
			Tipo:        "ConciliacionBancoCl",
			Descripcion: "No existe moviminetos informados en Banco",
		}
		s.utilService.CreateNotificacionService(requestNotificacion)
		if erro != nil {
			logs.Error("error: al crear Notificación: " + erro.Error())
			return
		}
		return
	}
	/*
		se recorre lista prisma pagos y lista de movimineto banco
		valido que el importe de movimiento banco coincida con el importe neto de pago
	*/
	logs.Info(fmt.Sprintf("cant. registros responseListprismaTrPagos - %v", len(responseListprismaTrPagos)))
	logs.Info("====================================================")
	logs.Info(fmt.Sprintf("cant. registros listaMovimientosBanco - %v", len(listaMovimientosBanco)))
	for _, valuePago := range responseListprismaTrPagos {
		for _, valueMovimientoBanco := range listaMovimientosBanco {
			bancoFecha, err := time.Parse("2006-01-02T00:00:00Z", valueMovimientoBanco.Fecha)
			if err != nil {
				erro = err
				return
			}
			if valuePago.FechaPago.Equal(bancoFecha) && valuePago.DetallePago[0].ImporteNeto == entities.Monto(valueMovimientoBanco.Importe) {
				var totalValorPresentado entities.Monto
				var totalImporteIvaArancel entities.Monto
				var msjEstadoCL string
				for _, valueCl := range valuePago.DetallePago[0].PrismaCierreLote {
					totalValorPresentado += valueCl.Valorpresentado
					totalImporteIvaArancel += entities.Monto(valueCl.Importeivaarancel * 100)
				}
				totalValorPresentado = totalValorPresentado - totalImporteIvaArancel
				// revisar crear la parte de banco
				for _, valueCl := range valuePago.DetallePago[0].PrismaCierreLote {
					if !valueCl.Reversion {
						msjEstadoCL = fmt.Sprintf("no existe conflicto")
						if totalValorPresentado != valuePago.DetallePago[0].ImporteNeto {
							valueCl.Diferenciaimporte = totalValorPresentado - valuePago.DetallePago[0].ImporteNeto
							msjEstadoCL = fmt.Sprintf("existe conflicto entre calculo de valores presentado (%v) y el informado por prisma (%v) id pagos es %v", totalValorPresentado, valuePago.DetallePago[0].ImporteNeto, valuePago.DetallePago[0].PrismatrcuatropagosId)
						}
						valueCl.Descripcion = msjEstadoCL
						valueCl.Match = 1
						valueCl.BancoExternalId = int64(valueMovimientoBanco.Id)
						ListCierreLote = append(ListCierreLote, valueCl.DtosToEntity())
					}
					if valueCl.Reversion {
						msjEstadoCL = fmt.Sprintf("no existe conflicto")
						if totalValorPresentado != valuePago.DetallePago[0].ImporteNeto {
							diferencia := totalValorPresentado - valuePago.DetallePago[0].ImporteNeto
							msjEstadoCL = fmt.Sprintf("existe diferencia %v entre valores presentado (%v) y el informado pen reversion (%v) pago id %v", diferencia, totalValorPresentado, valuePago.DetallePago[0].ImporteNeto, valuePago.DetallePago[0].PrismatrcuatropagosId)
						}
						valueCl.Descripcionbanco = msjEstadoCL
						valueCl.Conciliado = true
						valueCl.ExtbancoreversionId = int64(valueMovimientoBanco.Id)
						ListCierreLote = append(ListCierreLote, valueCl.DtosToEntity())
					}

				}

				movimientosIds = append(movimientosIds, valueMovimientoBanco.Id)
			}
		}

	}
	/* ya tengo lista de cierre lote actualizada y lista de ids movimineto banco
	 */
	logs.Info(ListCierreLote)
	logs.Info("=============================================")
	logs.Info(movimientosIds)
	if len(ListCierreLote) <= 0 && len(movimientosIds) <= 0 {
		erro = errors.New(ERROR_CONCILIACION)
		logs.Error(erro.Error())
		return
	}
	err = s.repository.ActualizarCierreLoteMatch(reversion, ListCierreLote)
	if err != nil {
		erro = errors.New(ERROR_ACTUALIZAR_CL_PRISMA + err.Error())
		logs.Error(erro.Error())
		return
	}
	/*
		envio al servicio de banco lista de ids de movimientos a conciliar
	*/
	_, err = s.bancoService.ActualizarRegistrosMatchBancoService(movimientosIds, true)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "ActualizarRegistrosMatchBancoService",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}
	return
}

func (s *service) ActualizarMovimientosBanco() (estadoResponse bool, erro error) {
	listCierreLote, err := s.repository.GetCierreLotePrismaByExternalIdAndMacht()
	if err != nil {
		erro = errors.New(ERROR_OBTENER_CL_PRISMA + err.Error())
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "GetCierreLotePrismaByExternalIdAndMacht",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}
	var movimientosIds []uint
	for _, valueCL := range listCierreLote {
		movimientosIds = append(movimientosIds, uint(valueCL.BancoExternalId))
	}
	estadoResponse, err = s.bancoService.ActualizarRegistrosMatchBancoService(movimientosIds, true)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "ActualizarRegistrosMatchBancoService",
			Mensaje:       erro.Error(),
		}
		err = s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error("error: al crear logs: " + err.Error())
			return
		}
	}
	return
}

func (s *service) ObtenerMxMoviminetosServices() (MovimientosMx []prismaCierreLote.ResponseMovimientoMx, entityMovimientoMxStr []entities.Prismamxtotalesmovimiento, erro error) {
	var listaMovimientoMx prismaCierreLote.ResponseMovimientoMx
	entityMovimientoMxStr, err := s.repository.GetMovimientosMxRepository()
	if err != nil {
		erro = errors.New(err.Error())
	}

	MovimientosMx, err = listaMovimientoMx.EntityToDtosMx(entityMovimientoMxStr)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	return
}

func (s *service) ObtenerTablasRelacionadasServices() (tablasRelacionadas prismaCierreLote.ResponseTablasRelacionadas, erro error) {
	var listaCodigoRechazos prismaCierreLote.ResponseCodigoRechazos
	var listaVisaContracargo prismaCierreLote.ResponseVisaContracargo
	var listaMotivosAjustes prismaCierreLote.ResponseMotivosAjustes
	var listaOperaciones prismaCierreLote.ResponseOperaciones
	var listaMasterContracargo prismaCierreLote.ResponseMasterContracargo
	resultCodigoRechazos, err := s.repository.GetCodigosRechazoRepository()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range resultCodigoRechazos {
		reslutado := listaCodigoRechazos.ToDtos(value)
		tablasRelacionadas.ListaCodigosRechazados = append(tablasRelacionadas.ListaCodigosRechazados, reslutado)
	}

	resultVisaContracargo, err := s.repository.GetVisaContracargoRepository()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range resultVisaContracargo {
		resultado := listaVisaContracargo.ToDtos(value)
		tablasRelacionadas.ListaVisaContracargo = append(tablasRelacionadas.ListaVisaContracargo, resultado)

	}
	resultMotivosAjustes, err := s.repository.GetMotivosAjustesRepository()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range resultMotivosAjustes {
		resultado := listaMotivosAjustes.ToDtos(value)
		tablasRelacionadas.ListaMotivosAjustes = append(tablasRelacionadas.ListaMotivosAjustes, resultado)

	}
	resultOperaciones, err := s.repository.GetOperacionesRepository()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range resultOperaciones {
		resultado := listaOperaciones.ToDtos(value)
		tablasRelacionadas.ListaOperaciones = append(tablasRelacionadas.ListaOperaciones, resultado)

	}
	resultMasterContracargo, err := s.repository.GetMasterContracargoRepository()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range resultMasterContracargo {
		resultado := listaMasterContracargo.ToDtos(value)
		tablasRelacionadas.ListaMasterContracargo = append(tablasRelacionadas.ListaMasterContracargo, resultado)
	}
	return
}

func (s *service) ProcesarMovimientoMxServices(movimientosMx []prismaCierreLote.ResponseMovimientoMx, tablasRelacionadas prismaCierreLote.ResponseTablasRelacionadas) (resultMovimientosMx []prismaCierreLote.ResponseMovimientoMx) {
	for _, value := range movimientosMx {
		var temporalMovimientosMxDetalles []prismaCierreLote.ResponseMovimientoMxDetalle
		for _, valueDetalle := range value.MovimientosMxDetalles {
			if strings.TrimSpace(valueDetalle.FechaPagoOrigenAjuste) == "" {
				valueDetalle.FechaPagoOrigenAjuste = "000000"
			}
			for _, valuePrimerRechazo := range tablasRelacionadas.ListaCodigosRechazados {
				if strings.TrimSpace(valueDetalle.RechazoPrincipalId) == "" {
					valueDetalle.RechazoPrincipalId = "0"
				}
				if valueDetalle.RechazoPrincipalId == valuePrimerRechazo.ExternalId {
					valueDetalle.RechazoPrincipalId = strconv.Itoa(int(valuePrimerRechazo.Id))
					break
				}
			}
			for _, valueSegundoRechazo := range tablasRelacionadas.ListaCodigosRechazados {
				if strings.TrimSpace(valueDetalle.RechazoSecundarioId) == "" {
					valueDetalle.RechazoSecundarioId = "0"
				}
				if valueDetalle.RechazoSecundarioId == valueSegundoRechazo.ExternalId {
					valueDetalle.RechazoSecundarioId = strconv.Itoa(int(valueSegundoRechazo.Id))
					break
				}
			}
			for _, valueVisaContracargo := range tablasRelacionadas.ListaVisaContracargo {
				if strings.TrimSpace(valueDetalle.PrismavisacontracargosId) == "" {
					valueDetalle.PrismavisacontracargosId = "00" //Por defecto
				}
				if valueDetalle.PrismavisacontracargosId == valueVisaContracargo.ExternalId {
					valueDetalle.PrismavisacontracargosId = strconv.Itoa(int(valueVisaContracargo.Id))
					break
				}
			}
			for _, valueMotivoAjuste := range tablasRelacionadas.ListaMotivosAjustes {
				if strings.TrimSpace(valueDetalle.PrismamotivosajustesId) == "" {
					valueDetalle.PrismamotivosajustesId = "0" //Por defecto
				}
				if valueDetalle.PrismamotivosajustesId == valueMotivoAjuste.ExternalId {
					valueDetalle.PrismamotivosajustesId = strconv.Itoa(int(valueMotivoAjuste.Id))
					break
				}
			}
			for _, valueOperacion := range tablasRelacionadas.ListaOperaciones {
				if valueDetalle.PrismaoperacionsId == valueOperacion.ExternalId {
					valueDetalle.PrismaoperacionsId = strconv.Itoa(int(valueOperacion.Id))
					break
				}
			}
			for _, valueMasterContracargo := range tablasRelacionadas.ListaMasterContracargo {
				if strings.TrimSpace(valueDetalle.PrismamastercontracargosId) == "" || valueDetalle.PrismamastercontracargosId == "0000" {
					valueDetalle.PrismamastercontracargosId = "0" //Por defecto
				}
				if valueDetalle.PrismamastercontracargosId == valueMasterContracargo.ExternalId {
					valueDetalle.PrismamastercontracargosId = strconv.Itoa(int(valueMasterContracargo.Id))
					break
				}
			}
			temporalMovimientosMxDetalles = append(temporalMovimientosMxDetalles, valueDetalle)
		}
		temporalMovimientosMXTotales := value.MovimientosMXTotales
		resultMovimientosMx = append(resultMovimientosMx, prismaCierreLote.ResponseMovimientoMx{MovimientosMXTotales: temporalMovimientosMXTotales, MovimientosMxDetalles: temporalMovimientosMxDetalles})
	}

	return
}
func (s *service) SaveMovimientoMxServices(movimientosMx []prismaCierreLote.ResponseMovimientoMx, movimientosMxEntity []entities.Prismamxtotalesmovimiento) (erro error) {
	var resultadoentity []entities.Prismamovimientototale
	for _, value := range movimientosMx {
		entityMx, err := value.ToEntity()
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		resultadoentity = append(resultadoentity, entityMx)
	}
	err := s.repository.SaveMovimientoMxRepository(resultadoentity, movimientosMxEntity)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	return
}

func (s *service) ObtenerPxPagosServices() (pagosPx []prismaCierreLote.ResponsePagoPx, entityPagoPxStr []entities.Prismapxcuatroregistro, erro error) {
	var listaPagoPx prismaCierreLote.ResponsePagoPx
	entityPagoPxStr, err := s.repository.GetPagosPxRepository()
	if err != nil {
		erro = errors.New(err.Error())
	}
	pagosPx, err = listaPagoPx.EntityToDtosPx(entityPagoPxStr)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	return
}

func (s *service) SavePagoPxServices(pagoPx []prismaCierreLote.ResponsePagoPx, entityPagoPxStr []entities.Prismapxcuatroregistro) (erro error) {
	var resultadoentity []entities.Prismatrcuatropago
	for _, value := range pagoPx {
		entityMx, err := value.ToEntity()
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		resultadoentity = append(resultadoentity, entityMx)
	}
	err := s.repository.SavePagosPxRepository(resultadoentity, entityPagoPxStr)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	return
}

func (s *service) ObtenerCierreloteServices(filtro filtrocl.FiltroCierreLote, codigoautorizacion []string) (listaCierreLote []prismaCierreLote.ResponsePrismaCL, erro error) {
	var listaTransaccionId []string
	var dtosListaCierreLote prismaCierreLote.ResponsePrismaCL
	entityCierreLote, err := s.repository.GetCierreLoteRepository(filtro)
	if err != nil {
		logs.Info(err.Error())
		erro = errors.New(ERROR_OBTENER_LISTA_CIERRE_LOTE + " - " + err.Error())
		return
	}

	if len(codigoautorizacion) > 0 {
		for _, valueEntityCl := range entityCierreLote {
			for _, value := range codigoautorizacion {
				if strings.Contains(valueEntityCl.Codigoautorizacion, value) {
					listaTransaccionId = append(listaTransaccionId, valueEntityCl.ExternalclienteID)
					dtosListaCierreLote.EntityToDtos(valueEntityCl)
					listaCierreLote = append(listaCierreLote, dtosListaCierreLote)
				}
			}
		}
	}
	if len(codigoautorizacion) <= 0 {
		for _, valueEntityCl := range entityCierreLote {
			listaTransaccionId = append(listaTransaccionId, valueEntityCl.ExternalclienteID)
			dtosListaCierreLote.EntityToDtos(valueEntityCl)
			listaCierreLote = append(listaCierreLote, dtosListaCierreLote)
		}
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		TransaccionesId:         listaTransaccionId,
		Channel:                 true,
		CargarPago:              true,
		CargarPagoTipo:          true,
		CargarCuenta:            true,
		CargarCliente:           true,
		CargarCuentaComision:    true,
		CargarImpuestos:         true,
		ExternalId:              true,
		CargarInstallmentdetail: true,
	}
	pagosImtentos, err := s.administracionService.GetPagosIntentosByTransaccionIdService(filtroPagoIntento)
	if err != nil {
		logs.Info(err.Error())
		erro = errors.New(ERROR_OBTENER_LISTA_PAGOS_INTENTOS + " - " + err.Error())
		return
	}
	for key := range listaCierreLote {
		var installmentDetails dtosCL.ResponseInstallmentInfo
		for _, valuePI := range pagosImtentos {
			if listaCierreLote[key].ExternalclienteID == valuePI.TransactionID {
				installmentDetails.EntityToDtosForConciliacion(valuePI.Installmentdetail)
				installmentDetails.TransaccionesId = listaCierreLote[key].ExternalclienteID
				listaCierreLote[key].Istallmentsinfo = installmentDetails
				break
			}
		}
	}
	return
}

func (s *service) ObtenerPrismaMovimientosServices(filtro filtrocl.FiltroPrismaMovimiento) (listaPrismaMovimientos []prismaCierreLote.ResponseMovimientoTotales, codigoautorizacion []string, erro error) {
	var dtosListaPrismaMovimientoTemporal prismaCierreLote.ResponseMovimientoTotales
	var dtosListaMovimientoDetalleTemporal prismaCierreLote.ResponseMoviminetoDetalles
	var entityPrismaMovimiento []entities.Prismamovimientototale
	var err error
	// if !filtro.ContraCargo {
	entityPrismaMovimiento, err = s.repository.GetPrismaMovimientosRepository(filtro)
	// }
	// if filtro.ContraCargo {
	// 	entityPrismaMovimiento, err = s.repository.GetContraCargoPrismaMovimientosRepository(filtro)
	// }
	if err != nil {
		logs.Info(err.Error())
		erro = errors.New(ERROR_OBTENER_LISTA_PRISMA_MOVIMIENTOS + err.Error())
		return
	}
	var listaPrismaMovimientosTemporal []prismaCierreLote.ResponseMovimientoTotales
	for _, valueMovimientoTotales := range entityPrismaMovimiento {
		dtosListaPrismaMovimientoTemporal.EntityToDtos(valueMovimientoTotales)
		for _, valueMovimientoDetalle := range valueMovimientoTotales.DetalleMovimientos {
			dtosListaMovimientoDetalleTemporal.EntityToDtos(valueMovimientoDetalle)
			dtosListaPrismaMovimientoTemporal.DetalleMovimientos = append(dtosListaPrismaMovimientoTemporal.DetalleMovimientos, dtosListaMovimientoDetalleTemporal)
		}
		listaPrismaMovimientosTemporal = append(listaPrismaMovimientosTemporal, dtosListaPrismaMovimientoTemporal)
		dtosListaPrismaMovimientoTemporal.DetalleMovimientos = nil
	}
	listaPrismaMovimientos = append(listaPrismaMovimientos, listaPrismaMovimientosTemporal...)
	for _, value := range listaPrismaMovimientos {
		for _, value1 := range value.DetalleMovimientos {
			codigoautorizacion = append(codigoautorizacion, value1.NroAutorizacionXl[3:len(value1.NroAutorizacionXl)])
		}
	}
	return
}

func (s *service) ConciliarCierreLotePrismaMovimientoServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaPrismaMovimientos []prismaCierreLote.ResponseMovimientoTotales) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {

	var listCabeceraMoviminetosId []int64
	// recorro lista de cierre de lotes
	for _, valueCL := range listaCierreLote {
		valorCuota := valueCL.Nrocuota
		if valueCL.Nrocuota == 1 {
			valorCuota = 0
		}
		//strNroEstablecimiento := strconv.Itoa(int(valueCL.Nroestablecimiento))
		// recorro la cabecer y el detallle de los moviminetos prisma
		for _, valueCabecera := range listaPrismaMovimientos {

			for _, valueDetalle := range valueCabecera.DetalleMovimientos {
				//valueCL.ExternalloteId == valueDetalle.Lote &&

				factoryConciliacion, err := s.factoryConciliacionMovimiento.GetTipoConciliacion(valueDetalle.Tipooperacion.ExternalId)
				if err != nil {
					erro = errors.New(err.Error())
					return
				}

				CierreLoteProcesada, detalleMoviminetosIds, cabeceraMoviminetosIds, err := factoryConciliacion.ConciliarTablas(valorCuota, valueCL, valueCabecera, valueDetalle)
				if err != nil {
					return
				}
				detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, detalleMoviminetosIds...)
				listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, cabeceraMoviminetosIds...)
				listaCierreLoteProcesada = append(listaCierreLoteProcesada, CierreLoteProcesada...)
				// aqui ***1
			}

		}
	}
	logs.Info(listaCierreLoteProcesada)
	cabeceraMoviminetosIdArray = commons.RemoveDuplicateValues(listCabeceraMoviminetosId)
	return
}

func (s *service) ActualizarCierreloteMoviminetosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
	var entityCierreLote []entities.Prismacierrelote
	for _, valueCL := range listaCierreLote {
		entityCierreLote = append(entityCierreLote, valueCL.DtosToEntity())
	}
	logs.Info("en servicio")
	logs.Info(listaCierreLote)
	err := s.repository.UpdateCierreloteAndMoviminetosRepository(entityCierreLote, listaIdsCabecera, listaIdsDetalle)
	if err != nil {
		erro = errors.New(ERROR_CONCILIACION_CL_MOVIMIENTOS + " - " + err.Error())
		logs.Error(erro.Error())
		return
	}
	return
}

func (s *service) ObtenerPrismaMovimientoConciliadosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, filtroCabecera filtrocl.FiltroPrismaMovimiento) (listaCierreLoteMovimiento []prismaCierreLote.ResponseCierrreLotePrismaMovimiento, erro error) {
	var listaIdsCabecera []int64
	filtroCabecera.Match = true
	// filtroCabecera.Match := filtrocl.FiltroPrismaMovimiento{
	// 	Match: true,
	// }
	entityMovimientosCabeceraConciliados, err := s.repository.GetMovimientosConciliadosRepository(filtroCabecera)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_LISTA_PRISMA_MOVIMIENTOS + err.Error())
		logs.Error(erro.Error())
		return
	}
	for _, value := range entityMovimientosCabeceraConciliados {
		listaIdsCabecera = append(listaIdsCabecera, int64(value.ID))
	}
	filtroDetalle := filtrocl.FiltroPrismaMovimientoDetalle{
		ListIdsCabecera:              listaIdsCabecera,
		Contracargovisa:              false,
		Contracargomaster:            false,
		Tipooperacion:                false,
		Rechazotransaccionprincipal:  false,
		Rechazotransaccionsecundario: false,
		Motivoajuste:                 false,
		FechaPago:                    "0000-00-00",
		Match:                        true,
		MatchCl:                      false,
	}
	if filtroCabecera.ContraCargo {
		filtroDetalle.FechaPago = ""
	}
	entityMovimientosDetalleConciliados, err := s.repository.GetMovimientosDetalleConciliadosRepository(filtroDetalle)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_LISTA_PRISMA_MOVIMIENTOS + " - " + err.Error())
		logs.Error(erro.Error())
		return
	}
	for _, valueCabecera := range entityMovimientosCabeceraConciliados {
		var movimientosConciliadosTemporal prismaCierreLote.ResponseCierrreLotePrismaMovimiento
		movimientosConciliadosTemporal.MovimientoCabecer.EntityToDtos(valueCabecera)
		var movimientoDetalleTemporal prismaCierreLote.ResponseMoviminetoDetalles
		for _, valueDetalle := range entityMovimientosDetalleConciliados {
			if valueCabecera.ID == valueDetalle.PrismamovimientototalesId {
				for _, valueCl := range listaCierreLote {
					if filtroCabecera.ContraCargo {
						if int64(valueDetalle.ID) == valueCl.DetallemovimientoId {
							movimientoDetalleTemporal.EntityToDtos(valueDetalle)
							movimientoDetalleTemporal.CierreLote = valueCl
							movimientosConciliadosTemporal.MovimientoDetalle = append(movimientosConciliadosTemporal.MovimientoDetalle, movimientoDetalleTemporal)
							break
						}
					}
					if !filtroCabecera.ContraCargo {
						if int64(valueDetalle.ID) == valueCl.PrismamovimientodetallesId {
							movimientoDetalleTemporal.EntityToDtos(valueDetalle)
							movimientoDetalleTemporal.CierreLote = valueCl
							movimientosConciliadosTemporal.MovimientoDetalle = append(movimientosConciliadosTemporal.MovimientoDetalle, movimientoDetalleTemporal)
							break
						}
					}

				}
			}
		}
		if len(movimientosConciliadosTemporal.MovimientoDetalle) != 0 {
			listaCierreLoteMovimiento = append(listaCierreLoteMovimiento, movimientosConciliadosTemporal)
		}
	}

	return
}
func (s *service) ObtenerPrismaPagosServices(filtro filtrocl.FiltroPrismaTrPagos) (listaPrismaPago []prismaCierreLote.ResponseTrPagosCabecera, erro error) {
	entityPrismaPago, err := s.repository.GetPrismaPagosRepository(filtro)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_LISTA_PRISMA_PAGOS + " - " + err.Error())
		logs.Error(erro.Error())
		return
	}
	for _, valuePagoCabecera := range entityPrismaPago {
		var pagoTemporal prismaCierreLote.ResponseTrPagosCabecera
		pagoTemporal.EntityToDtos(valuePagoCabecera)
		for _, valuePagoDetalle := range valuePagoCabecera.Pagostrdos {
			var pagoDetalleTemporal prismaCierreLote.ResponseTrPagosDetalle
			pagoDetalleTemporal.EntityToDtos(valuePagoDetalle)
			pagoTemporal.DetallePago = append(pagoTemporal.DetallePago, pagoDetalleTemporal)
		}
		listaPrismaPago = append(listaPrismaPago, pagoTemporal)
	}
	return
}

func (s *service) ConciliarCierreLotePrismaPagoServices(listaCierreLoteMovimientos []prismaCierreLote.ResponseCierrreLotePrismaMovimiento, listaPrismaPago []prismaCierreLote.ResponseTrPagosCabecera) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detallePagosIdArray []int64, cabeceraPagosIdArray []int64, erro error) {
	var listCabeceraPagosId []int64
	var listDetallePagosId []int64
	for _, valuePagoCabecera := range listaPrismaPago {
		for _, valuePagoDetalle := range valuePagoCabecera.DetallePago {
			for _, valueClMovimiento := range listaCierreLoteMovimientos {
				if valuePagoCabecera.FechaPago == valueClMovimiento.MovimientoCabecer.FechaPago && valuePagoCabecera.FechaPresentacion == valueClMovimiento.MovimientoCabecer.FechaPresentacion && valuePagoCabecera.EstablecimientoNro == valueClMovimiento.MovimientoCabecer.EstablecimientoNro {
					for _, valueclMovimientoDetalle := range valueClMovimiento.MovimientoDetalle {
						if valueclMovimientoDetalle.CierreLote.Tipooperacion == "C" && !valueclMovimientoDetalle.CierreLote.Reversion {

							if valuePagoCabecera.FechaPago == valueClMovimiento.MovimientoCabecer.FechaPago && valuePagoCabecera.FechaPresentacion == valueClMovimiento.MovimientoCabecer.FechaPresentacion && valuePagoCabecera.EstablecimientoNro == valueClMovimiento.MovimientoCabecer.EstablecimientoNro && valuePagoDetalle.LiquidacionNro == valueclMovimientoDetalle.NroLiquidacion && valuePagoDetalle.ImporteNeto != 0 && valuePagoDetalle.ImporteBruto != 0 { //&& valuePagoDetalle.ImporteBruto == valueClMovimiento.MovimientoCabecer.ImporteTotal
								listClTemp := valueclMovimientoDetalle.CierreLote
								listClTemp.PrismatrdospagosId = int64(valuePagoDetalle.Id)
								listaCierreLoteProcesada = append(listaCierreLoteProcesada, listClTemp)
								listCabeceraPagosId = append(listCabeceraPagosId, int64(valuePagoCabecera.Id))
								listDetallePagosId = append(listDetallePagosId, int64(valuePagoDetalle.Id))
							}
						}
						if valueclMovimientoDetalle.CierreLote.Tipooperacion == "D" && !valueclMovimientoDetalle.CierreLote.Reversion {

							if valuePagoCabecera.FechaPago == valueClMovimiento.MovimientoCabecer.FechaPago && valuePagoCabecera.FechaPresentacion == valueClMovimiento.MovimientoCabecer.FechaPresentacion && valuePagoCabecera.EstablecimientoNro == valueClMovimiento.MovimientoCabecer.EstablecimientoNro && valuePagoDetalle.LiquidacionNro == valueclMovimientoDetalle.NroLiquidacion && valuePagoDetalle.ImporteNeto != 0 && valuePagoDetalle.ImporteBruto != 0 {
								listClTemp := valueclMovimientoDetalle.CierreLote
								listClTemp.PrismatrdospagosId = int64(valuePagoDetalle.Id)
								listaCierreLoteProcesada = append(listaCierreLoteProcesada, listClTemp)
								listCabeceraPagosId = append(listCabeceraPagosId, int64(valuePagoCabecera.Id))
								listDetallePagosId = append(listDetallePagosId, int64(valuePagoDetalle.Id))
							}
						}
						if valueclMovimientoDetalle.CierreLote.Tipooperacion == "C" && valueclMovimientoDetalle.CierreLote.Reversion {
							if valuePagoCabecera.FechaPago == valueClMovimiento.MovimientoCabecer.FechaPago && valuePagoCabecera.FechaPresentacion == valueClMovimiento.MovimientoCabecer.FechaPresentacion && valuePagoCabecera.EstablecimientoNro == valueClMovimiento.MovimientoCabecer.EstablecimientoNro && valuePagoDetalle.LiquidacionNro == valueclMovimientoDetalle.NroLiquidacion && valueclMovimientoDetalle.TipoAplicacion == "-" {
								listClTemp := valueclMovimientoDetalle.CierreLote
								listClTemp.DetallepagoId = int64(valuePagoDetalle.Id)
								listaCierreLoteProcesada = append(listaCierreLoteProcesada, listClTemp)
								listCabeceraPagosId = append(listCabeceraPagosId, int64(valuePagoCabecera.Id))
								listDetallePagosId = append(listDetallePagosId, int64(valuePagoDetalle.Id))

							}

						}

					}
				}
			}
		}
	}
	cabeceraPagosIdArray = commons.RemoveDuplicateValues(listCabeceraPagosId)
	detallePagosIdArray = commons.RemoveDuplicateValues(listDetallePagosId)
	return
}

func (s *service) ActualizarCierrelotePagosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
	var entityCierreLote []entities.Prismacierrelote
	for _, valueCL := range listaCierreLote {
		entityCierreLote = append(entityCierreLote, valueCL.DtosToEntity())
	}
	err := s.repository.UpdateCierreloteAndPagosRepository(entityCierreLote, listaIdsCabecera, listaIdsDetalle)
	if err != nil {
		erro = errors.New(ERROR_CONCILIACION_CL_MOVIMIENTOS + " - " + err.Error())
		logs.Error(erro.Error())
		return
	}
	return
}

func (s *service) ObtenerRepoPagosPrisma(filtro filtrocl.FiltroTablasConciliadas) (responseListprismaTrPagos []prismaCierreLote.ResponseTrPagosCabecera, erro error) {

	entityPrismaTr4Pago, err := s.repository.GetCierreLoteMatch(filtro)
	if err != nil {
		erro = errors.New(ERROR_CONCILIACION_CL_MOVIMIENTOS + " - " + err.Error())
		logs.Error(erro.Error())
		return
	}
	for _, value := range entityPrismaTr4Pago {
		var prismaTrPagosCabeceraDtos prismaCierreLote.ResponseTrPagosCabecera
		var prismaTrPagosDetalleDtos prismaCierreLote.ResponseTrPagosDetalle

		filtroCL := filtrocl.FiltroCierreLote{
			MovimientosMX: false,
			PagosPx:       false,
		}
		if !filtro.Reversion {
			filtroCL.PrismaPagoId = int64(value.Pagostrdos[0].ID)
		}
		if filtro.Reversion {
			filtroCL.DetallePagoId = int64(value.Pagostrdos[0].ID)
			filtroCL.ContraCargo = true
			filtroCL.Reversion = true
			filtroCL.ContraCargoMx = true
			filtroCL.ContraCargoPx = true
		}

		EntityCierreLote, err := s.repository.GetCierreLoteRepository(filtroCL)
		if err != nil {
			erro = errors.New(ERROR_CONSULTAR_CIERRE_LOTE + " - " + err.Error())
			logs.Error(erro.Error())
			return
		}
		var listCierreLoteTemporal []prismaCierreLote.ResponsePrismaCL
		if len(EntityCierreLote) != 0 {
			for _, valueCL := range EntityCierreLote {
				var listCierreLoteDtos prismaCierreLote.ResponsePrismaCL
				listCierreLoteDtos.EntityToDtos(valueCL)
				listCierreLoteTemporal = append(listCierreLoteTemporal, listCierreLoteDtos)
			}
			prismaTrPagosCabeceraDtos.EntityToDtos(value)
			prismaTrPagosDetalleDtos.EntityToDtos(value.Pagostrdos[0])
			prismaTrPagosDetalleDtos.PrismaCierreLote = listCierreLoteTemporal
			prismaTrPagosCabeceraDtos.DetallePago = append(prismaTrPagosCabeceraDtos.DetallePago, prismaTrPagosDetalleDtos)
			responseListprismaTrPagos = append(responseListprismaTrPagos, prismaTrPagosCabeceraDtos)
		}
	}
	return
}

//ConciliarCierreLotePrismaMovimientoServices ***1
// if valueCL.Tipooperacion == "C" {
// 	if valueCL.Fechaoperacion == valueDetalle.FechaOrigenCompra && valueCL.FechaCierre == valueCabecera.FechaPresentacion && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && strings.Contains(valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion) && valueCL.Nroticket == valueDetalle.NroCupon && valorCuota == valueDetalle.PlanCuota && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "+" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {
// 		valueCL.FechaPago = valueCabecera.FechaPago
// 		valueCL.PrismamovimientodetallesId = valueDetalle.Id
// 		detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
// 		listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
// 		listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
// 	}
// }
// if valueCL.Tipooperacion == "D" {
// 	if strings.Contains(valueDetalle.Tipooperacion.Operacion, "REVERSO") {
// 		logs.Info("no hay reverso")
// 	}

// 	if strings.Contains(valueDetalle.Tipooperacion.Operacion, "CONTRACARGO") {
// 		if valueCL.Fechaoperacion == valueDetalle.FechaOrigenCompra && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && strings.Contains(valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion) && valueCL.Nroticket == valueDetalle.NroCupon && valorCuota == valueDetalle.PlanCuota && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "-" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {
// 			valueCL.FechaPago = valueCabecera.FechaPago
// 			valueCL.PrismamovimientodetallesId = valueDetalle.Id
// 			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
// 			listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
// 			listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
// 		}
// 	}
// 	if strings.Contains(valueDetalle.Tipooperacion.Operacion, "DEVOLUCIÓN") {
// 		if valueCL.FechaCierre == valueCabecera.FechaPresentacion && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.ExternalloteId == valueDetalle.Lote && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "-" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {
// 			valueCL.FechaPago = valueCabecera.FechaPago
// 			valueCL.PrismamovimientodetallesId = valueDetalle.Id
// 			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
// 			listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
// 			listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
// 		}

// 	}

// }
///logs tipo operacion C
// fmt.Println("=================================")
// fmt.Println("=================================")
// fmt.Println("======Fecha de pago=====")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.FechaPago, valueClMovimiento.MovimientoCabecer.FechaPago)
// fmt.Println("======Fecha de presentacion=====")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.FechaPresentacion, valueClMovimiento.MovimientoCabecer.FechaPresentacion)
// fmt.Println("============Establecimiento======")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.EstablecimientoNro, valueClMovimiento.MovimientoCabecer.EstablecimientoNro)
// fmt.Println("============Nro.liquidacion==========")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoDetalle.LiquidacionNro, valueclMovimientoDetalle.NroLiquidacion)
// fmt.Println("============Importe bruto==============")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoDetalle.ImporteBruto, valueClMovimiento.MovimientoCabecer.ImporteTotal)
// fmt.Println("============Importe neto==============")
// fmt.Printf("pago: %v  \n", valuePagoDetalle.ImporteNeto)
// fmt.Println("=================================")
// fmt.Println("=================================")
///logs tipo operacion D
// fmt.Println("=================================")
// fmt.Println("=================================")
// fmt.Println("======Fecha de pago=====")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.FechaPago, valueClMovimiento.MovimientoCabecer.FechaPago)
// fmt.Println("======Fecha de presentacion=====")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.FechaPresentacion, valueClMovimiento.MovimientoCabecer.FechaPresentacion)
// fmt.Println("============Establecimiento======")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoCabecera.EstablecimientoNro, valueClMovimiento.MovimientoCabecer.EstablecimientoNro)
// fmt.Println("============Nro.liquidacion==========")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoDetalle.LiquidacionNro, valueclMovimientoDetalle.NroLiquidacion)
// fmt.Println("============Importe bruto==============")
// fmt.Printf("pago: %v - movimiento: %v \n", valuePagoDetalle.ImporteBruto, valueClMovimiento.MovimientoCabecer.ImporteTotal)
// fmt.Println("============Importe neto==============")
// fmt.Printf("pago: %v  \n", valuePagoDetalle.ImporteNeto)
// fmt.Println("=================================")
// fmt.Println("=================================")
/*
func (s *service) ConciliarCierreLotePrismaMovimientoServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaPrismaMovimientos []prismaCierreLote.ResponseMovimientoTotales) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	var listCabeceraMoviminetosId []int64
	for _, valueCL := range listaCierreLote {
		valorCuota := valueCL.Nrocuota
		if valueCL.Nrocuota == 1 {
			valorCuota = 0
		}
		strNroEstablecimiento := strconv.Itoa(int(valueCL.Nroestablecimiento))
		for _, valueCabecera := range listaPrismaMovimientos {
			for _, valueDetalle := range valueCabecera.DetalleMovimientos {
				//valueCL.ExternalloteId == valueDetalle.Lote &&

				if valueCL.Tipooperacion == "C" {
					if valueCL.Fechaoperacion == valueDetalle.FechaOrigenCompra && valueCL.FechaCierre == valueCabecera.FechaPresentacion && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && strings.Contains(valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion) && valueCL.Nroticket == valueDetalle.NroCupon && valorCuota == valueDetalle.PlanCuota && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "+" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {

						fmt.Println("=================================")
						fmt.Println("=================================")
						fmt.Println("============Fecha Operacion======")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Fechaoperacion, valueDetalle.FechaOrigenCompra)
						fmt.Println("============Fecha Presentacion===")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.FechaCierre, valueCabecera.FechaPresentacion)
						fmt.Println("============Establecimiento======")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCabecera.EstablecimientoNro, strNroEstablecimiento)

						fmt.Println("============Nro.Tarjeta==========")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nrotarjeta, valueDetalle.NroTarjetaXl)
						fmt.Println("============Nor.Atorizacion======")
						fmt.Printf("cl: %v - movimiento: %v \n", valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion)

						fmt.Println("============Lote=================")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.ExternalloteId, valueDetalle.Lote)

						fmt.Println("============Ticket===============")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nroticket, valueDetalle.NroCupon)

						fmt.Println("============Cuota================")
						fmt.Printf("cl: %v - movimiento: %v \n", valorCuota, valueDetalle.PlanCuota)

						fmt.Println("============Importe==============")
						fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Monto.Int64(), int64(valueDetalle.Importe))
						fmt.Println("=================================")
						fmt.Println("=================================")

						valueCL.FechaPago = valueCabecera.FechaPago
						valueCL.PrismamovimientodetallesId = valueDetalle.Id
						detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
						listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
						listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
					}
				}
				if valueCL.Tipooperacion == "D" {
					if strings.Contains(valueDetalle.Tipooperacion.Operacion, "REVERSO") {
						fmt.Println("************************************")
						fmt.Println("**************REVEERSO**************")
					}

					if strings.Contains(valueDetalle.Tipooperacion.Operacion, "CONTRACARGO") {
						// fehca origen compra - establecimiento - nro tarjeta - nro autorizacion - nro cupon - valor cupon - valor cuota - tipo aplicacion - codigo operacion
						fmt.Println("************************************")
						fmt.Println("**************CONTRACARGO**************")
						fmt.Println("************************************")
						if valueCL.Fechaoperacion == valueDetalle.FechaOrigenCompra && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && strings.Contains(valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion) && valueCL.Nroticket == valueDetalle.NroCupon && valorCuota == valueDetalle.PlanCuota && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "-" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {
							fmt.Println("=================================")
							fmt.Println("=================================")
							fmt.Println("======Fecha Origen de compra=====")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Fechaoperacion, valueDetalle.FechaOrigenCompra)
							fmt.Println("============Establecimiento======")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCabecera.EstablecimientoNro, strNroEstablecimiento)
							fmt.Println("============Nro.cupon==========")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nroticket, valueDetalle.NroCupon)
							fmt.Println("======Codigo Autorizacion========")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Codigoautorizacion, valueDetalle.NroAutorizacionXl)
							fmt.Println("============Nro.Tarjeta==========")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nrotarjeta, valueDetalle.NroTarjetaXl)

							fmt.Println("============Importe==============")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Monto.Int64(), int64(valueDetalle.Importe))
							fmt.Println("=================================")
							fmt.Println("=================================")

							valueCL.FechaPago = valueCabecera.FechaPago
							valueCL.PrismamovimientodetallesId = valueDetalle.Id
							detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
							listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
							listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
						}
					}
					if strings.Contains(valueDetalle.Tipooperacion.Operacion, "DEVOLUCIÓN") {
						// fecha presentacion - nro establecimiento - nro lote - importe - nro tarjeta - tipo aplicacion - codigo operacion
						fmt.Println("************************************")
						fmt.Println("**************DEVOLUCIÓN**************")
						fmt.Println("************************************")
						if valueCL.FechaCierre == valueCabecera.FechaPresentacion && strings.Contains(valueCabecera.EstablecimientoNro, strNroEstablecimiento) && valueCL.ExternalloteId == valueDetalle.Lote && valueCL.Nrotarjeta == valueDetalle.NroTarjetaXl && valueCL.Monto.Int64() == int64(valueDetalle.Importe) && valueDetalle.TipoAplicacion == "-" && valueCabecera.Codop == valueDetalle.Tipooperacion.ExternalId {
							fmt.Println("=================================")
							fmt.Println("=================================")
							fmt.Println("======Fecha presentacion=====")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.FechaCierre, valueCabecera.FechaPresentacion)
							fmt.Println("============Establecimiento======")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCabecera.EstablecimientoNro, strNroEstablecimiento)

							fmt.Println("============lote======")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.ExternalloteId, valueDetalle.Lote)

							fmt.Println("============Nro.Tarjeta==========")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nrotarjeta, valueDetalle.NroTarjetaXl)

							fmt.Println("============Importe==============")
							fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Monto.Int64(), int64(valueDetalle.Importe))
							fmt.Println("=================================")
							fmt.Println("=================================")

							valueCL.FechaPago = valueCabecera.FechaPago
							valueCL.PrismamovimientodetallesId = valueDetalle.Id
							detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, valueDetalle.Id)
							listCabeceraMoviminetosId = append(listCabeceraMoviminetosId, int64(valueDetalle.PrismamovimientototalesId))
							listaCierreLoteProcesada = append(listaCierreLoteProcesada, valueCL)
						}

					}

				}
			}
		}
	}
	cabeceraMoviminetosIdArray = commons.RemoveDuplicateValues(listCabeceraMoviminetosId)
	return
}
*/

/*
para imprimir y verificar los cl y movimiento
				if valueCL.Fechaoperacion == valueDetalle.FechaOrigenCompra && valueCL.FechaCierre == valueCabecera.FechaPresentacion && valueCL.Monto.Int64() == int64(valueDetalle.Importe) {
					fmt.Println("=================================")
					fmt.Println("=================================")
					fmt.Println("============Fecha Operacion======")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Fechaoperacion, valueDetalle.FechaOrigenCompra)
					fmt.Println("============Fecha Presentacion===")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.FechaCierre, valueCabecera.FechaPresentacion)
					fmt.Println("============Establecimiento======")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCabecera.EstablecimientoNro, strNroEstablecimiento)

					fmt.Println("============Nro.Tarjeta==========")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nrotarjeta, valueDetalle.NroTarjetaXl)
					fmt.Println("============Nor.Atorizacion======")
					fmt.Printf("cl: %v - movimiento: %v \n", valueDetalle.NroAutorizacionXl, valueCL.Codigoautorizacion)

					fmt.Println("============Lote=================")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.ExternalloteId, valueDetalle.Lote)

					fmt.Println("============Ticket===============")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Nroticket, valueDetalle.NroCupon)

					fmt.Println("============Cuota================")
					fmt.Printf("cl: %v - movimiento: %v \n", valorCuota, valueDetalle.PlanCuota)

					fmt.Println("============Importe==============")
					fmt.Printf("cl: %v - movimiento: %v \n", valueCL.Monto.Int64(), int64(valueDetalle.Importe))
					fmt.Println("=================================")
					fmt.Println("=================================")

				}
*/

/////////////////////funciones///////////////////////
// func convertirFechaStringToTime(formatofecha, fechaString string) (fechaTime time.Time, erro error) {
// 	fechaTime, err := time.Parse(formatofecha, fechaString)
// 	if err != nil {
// 		erro = errors.New(ERROR_PARSER_FECHA + err.Error())
// 		logs.Error(erro.Error())
// 		return
// 	}
// 	return
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// func (s *service) ConciliacionBancoPrisma(responseListprismaTrPagos []prismaCierreLote.ResponseTrPagosCabecera) (listaMovimientosBanco []bancodtos.ResponseMovimientosBanco, erro error) {

// 	var referenciaMovimiento []string
// 	var movimientosIds []uint
// 	/* obtener configuracion periodo de acreditacion */
// 	filtro := filtros.ConfiguracionFiltro{
// 		Buscar:     true,
// 		Nombrelike: "PRISMA_PERIODO_ACREDITACION",
// 	}
// 	prismaPeriodosAcreditaciones, err := s.utilService.GetConfiguracionesService(filtro)
// 	if err != nil {
// 		erro = errors.New(ERROR_OBTENER_CONFIGURACION)
// 		logs.Error(err.Error())
// 		log := entities.Log{
// 			Tipo:          entities.EnumLog("Error"),
// 			Funcionalidad: "GetConfiguracionesService",
// 			Mensaje:       erro.Error() + " " + err.Error(),
// 		}
// 		err = s.utilService.CreateLogService(log)
// 		if err != nil {
// 			logs.Error("error: al crear logs: " + err.Error())
// 			return
// 		}
// 		return
// 	}
// 	if len(prismaPeriodosAcreditaciones) == 0 {
// 		erro = errors.New(ERROR_NO_EXISTE_CONFIGURACION)
// 		logs.Error(erro.Error())
// 		log := entities.Log{
// 			Tipo:          entities.EnumLog("Error"),
// 			Funcionalidad: "GetConfiguracionesService",
// 			Mensaje:       erro.Error() + " " + err.Error(),
// 		}
// 		err = s.utilService.CreateLogService(log)
// 		if err != nil {
// 			logs.Error("error: al crear logs: " + err.Error())
// 			return
// 		}
// 		return
// 	}
// 	var cierreLotePrisma []prismaCierreLote.PrismaClResultGroup
// 	for _, valuePeriodoAcreditacion := range prismaPeriodosAcreditaciones {
// 		cuotaArrayString := strings.Split(valuePeriodoAcreditacion.Nombre, "_")
// 		cuotaString := cuotaArrayString[len(cuotaArrayString)-1]
// 		cuota, err := strconv.Atoi(cuotaString)
// 		if err != nil {
// 			erro = errors.New(ERROR_CONVERTIR_VALOR)
// 			logs.Error(err.Error())
// 			return
// 		}

// 		/* obtengo lista de cierre de lotes por nro de establecimiento	*/
// 		listaTemporalClPrisma, err := s.repository.GetCierreLoteGroupByRepository(int64(cuota))
// 		if err != nil {
// 			erro = errors.New(ERROR_OBTENER_CL_PRISMA + err.Error())
// 			logs.Error(erro.Error())
// 			log := entities.Log{
// 				Tipo:          entities.EnumLog("Error"),
// 				Funcionalidad: "GetCierreLoteByNroEstablecimiento",
// 				Mensaje:       erro.Error(),
// 			}
// 			erro = s.utilService.CreateLogService(log)
// 			if erro != nil {
// 				logs.Error("error: al crear logs: " + erro.Error())
// 				return
// 			}
// 		}
// 		for _, valueListaTemporal := range listaTemporalClPrisma {
// 			ValorInt64, err := strconv.ParseInt(valuePeriodoAcreditacion.Valor, 10, 64)
// 			if err != nil {
// 				erro = errors.New(ERROR_CONVERTIR_VALOR)
// 				logs.Error(err.Error())
// 				return
// 			}
// 			fechaCierre, err := convertirFechaStringToTime("2006-01-02T00:00:00Z", valueListaTemporal.FechaCierre)
// 			if err != nil {
// 				erro = err
// 				return
// 			}
// 			horas := ValorInt64 * 24
// 			fechaAcreditacion := fechaCierre.Add(time.Hour * time.Duration(horas))
// 			valueListaTemporal.FechaAcreditacion = fechaAcreditacion
// 			/*
// 				FIXME: posiblemente se tendira que verificar tambien si la fecha de acreditacion es igual
// 				a la de acreditacion
// 			*/
// 			logs.Info(time.Now().After(fechaAcreditacion))
// 			if time.Now().After(fechaAcreditacion) {
// 				cierreLotePrisma = append(cierreLotePrisma, valueListaTemporal)
// 			}
// 		}
// 	}

// 	/* obtengo lista de cierre de lotes por nro de establecimiento	*/

// 	/* recorro la lista clPrisma y genero lista de string referenciaMovimiento  */
// 	for _, cl := range cierreLotePrisma {
// 		referenciaMovimiento = append(referenciaMovimiento, cl.Nroestablecimiento)
// 	}
// 	//logs.Info(referenciaMovimiento)
// 	/* obtengo lista de movimientos de banco relacionados con prisma */
// 	filtroMov := filtrobanco.MovimientosBancoFiltro{
// 		SubCuenta:      config.COD_SUBCUENTA,
// 		Tipo:           "prisma",
// 		TipoMovimiento: referenciaMovimiento,
// 	}
// 	listaMovimientosBanco, err = s.bancoService.BuildCierreLoteApiLinkBancoService(filtroMov)
// 	if err != nil {
// 		erro = errors.New(ERROR_OBTENER_MOVIMIENTOS_BANCO + err.Error())
// 		logs.Error(erro.Error())
// 		log := entities.Log{
// 			Tipo:          entities.EnumLog("Error"),
// 			Funcionalidad: "BuildCierreLoteApiLinkBancoService",
// 			Mensaje:       erro.Error(),
// 		}
// 		erro = s.utilService.CreateLogService(log)
// 		if erro != nil {
// 			logs.Error("error: al crear logs: " + erro.Error())
// 			return
// 		}
// 	}
// 	/* recorro la lista clPrisma y listaMovimientoBanco, se construye lista de movimientosIds con aquellos que coincidan
// 	   con el nro de establecimiento, fecha y monto */
// 	for keyCL, valueClPrisma := range cierreLotePrisma {
// 		for _, valueMovimientosBanco := range listaMovimientosBanco {
// 			clFecha, err := convertirFechaStringToTime("2006-01-02T00:00:00Z", valueClPrisma.FechaCierre)
// 			if err != nil {
// 				erro = err
// 				return
// 			}
// 			bancoFecha, err := time.Parse("2006-01-02T00:00:00Z", valueMovimientosBanco.Fecha)
// 			if err != nil {
// 				erro = err
// 				return
// 			}
// 			if valueClPrisma.Nroestablecimiento == valueMovimientosBanco.Referencia && clFecha.Equal(bancoFecha) && valueClPrisma.Monto == int64(valueMovimientosBanco.Importe) {
// 				movimientosIds = append(movimientosIds, valueMovimientosBanco.Id)
// 				cierreLotePrisma[keyCL].EstadoConciliacion = true
// 				cierreLotePrisma[keyCL].BancoExternalId = int64(valueMovimientosBanco.Id)
// 				break
// 			}
// 		}
// 	}

// 	/*
// 		se actualiza el campo match en la tabla cierreloteprisma y se agrega el id del movimiento banco
// 	*/
// 	var listCierreLoteConcilidada []entities.Prismacierrelote
// 	for _, valueCL := range cierreLotePrisma {
// 		if valueCL.EstadoConciliacion {
// 			clConciliar, err := s.repository.GetCierreLoteByGroup(valueCL)
// 			if err != nil {
// 				erro = errors.New(ERROR_OBTENER_CL_PRISMA + err.Error())
// 				logs.Error(erro.Error())
// 				return
// 			}
// 			for _, valueClConciliar := range clConciliar {
// 				valueClConciliar.Match = 1
// 				valueClConciliar.BancoExternalId = valueCL.BancoExternalId
// 				listCierreLoteConcilidada = append(listCierreLoteConcilidada, valueClConciliar)
// 			}
// 		}
// 	}
// 	if len(listCierreLoteConcilidada) <= 0 {
// 		erro = errors.New(ERROR_CONCILIACION)
// 		logs.Error(erro.Error())
// 		return
// 	}
// 	err = s.repository.ActualizarCierreLoteMatch(listCierreLoteConcilidada)
// 	if err != nil {
// 		erro = errors.New(ERROR_ACTUALIZAR_CL_PRISMA + err.Error())
// 		logs.Error(erro.Error())
// 		return
// 	}

// 	/*
// 		envio al servicio de banco lista de ids de movimientos a conciliar
// 	*/
// 	_, err = s.bancoService.ActualizarRegistrosMatchBancoService(movimientosIds, true)
// 	if err != nil {
// 		erro = errors.New(err.Error())
// 		logs.Error(erro.Error())
// 		log := entities.Log{
// 			Tipo:          entities.EnumLog("Error"),
// 			Funcionalidad: "ActualizarRegistrosMatchBancoService",
// 			Mensaje:       erro.Error(),
// 		}
// 		erro = s.utilService.CreateLogService(log)
// 		if erro != nil {
// 			logs.Error("error: al crear logs: " + erro.Error())
// 			return
// 		}
// 	}

// 	return
// }
