package background

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/cierrelote"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/banco"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/reportes"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"

// 	webhooks "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
// 	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
// 	filtroCl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
// 	"github.com/robfig/cron"
// )

// // retorna true o false segun la fecha actual sea feriado, comparando con un string de fechas
// func _esFeriado(stringFechas string) (result bool) {
// 	// se toma la fecha actual en formato yyyy-mm-dd
// 	now := time.Now().UTC().Format("2006-01-02")
// 	// separador del split
// 	var sep string = ","
// 	// fechas en formato yyyy-mm-dd en tipo []string
// 	fechas := strings.Split(stringFechas, sep)

// 	result = commons.ContainStrings(fechas, now)
// 	return
// }

// func _buildNotificacion(service util.UtilService, erro error, tipo entities.EnumTipoNotificacion) {
// 	notificacion := entities.Notificacione{
// 		Tipo:        tipo,
// 		Descripcion: fmt.Sprintf("Configuración inválida. %s", erro.Error()),
// 	}
// 	service.CreateNotificacionService(notificacion)
// }

// func _buildPeriodicidad(service util.UtilService, nombreConfig string, valorConfig string, descripcionConfig string) (configuracion entities.Configuracione, erro error) {

// 	filtro := filtros.ConfiguracionFiltro{
// 		Nombre: nombreConfig,
// 	}

// 	config, erro := service.GetConfiguracionService(filtro)

// 	configuracion.Nombre = config.Nombre
// 	configuracion.Valor = config.Valor
// 	configuracion.ID = config.Id

// 	if erro != nil {
// 		_buildNotificacion(service, erro, entities.NotificacionConfiguraciones)
// 		return
// 	}

// 	if configuracion.ID == 0 {

// 		config := administraciondtos.RequestConfiguracion{
// 			Nombre:      nombreConfig,
// 			Descripcion: descripcionConfig,
// 			Valor:       valorConfig,
// 		}

// 		configuracion = config.ToEntity(false)

// 		_, erro = service.CreateConfiguracionService(config)

// 		if erro != nil {

// 			_buildNotificacion(service, erro, entities.NotificacionConfiguraciones)

// 		}

// 	}

// 	return

// }

// func BackgroudServices(service administracion.Service, cierrelote cierrelote.Service, util util.UtilService, movimientosBanco banco.BancoService, reportes reportes.ReportesService) {

// 	c := cron.New()
// 	/* TODO -> INICIO PROCESO PARA OBTENER LOS ARCHIVOS TXT DE S3 Y GUARDAR EN LA DB LA INFORMACION OBTENIDA DEL TXT*/
// 	confPrismaCierreLote, err := _buildPeriodicidad(util, "PERIODICIDAD_PRISMA_CIERRE_LOTE", "0 00 13 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de prisma")
// 	//confPrismaCierreLote.Valor = "0 30 * * * *"
// 	if err != nil {
// 		panic(err)
// 	}
// 	// obtener lista de pagos estado
// 	c.AddFunc(confPrismaCierreLote.Valor, func() {
// 		ctxaws := context.Background()
// 		/*
// 			PROCESO: Diariamente el proceso de cierre lote en backproud se divide en 4 pasos.
// 				paso 1:
// 					leer el directorio ftp de cierre de lote minio y obtener informacion de los archivos
// 					y se guarda en un directorio temporal los archivos txt existentes
// 					se obtiene todos los estados externos de prisma
// 				paso 2:
// 					se recorren uno a uno los archivos de cierre de lotes y se almacena a la bd
// 				paso 3:
// 					se mueven todos los archivos de la carpeta temporal al minio.
// 				paso 4:
// 					por ultimo se borran todos los archivos creados temporalmente y el directorio temporal
// 		*/
// 		///////////////////////////////////////PROCESO///////////////////////////////////////
// 		/* paso 1: */
// 		archivos, rutaArchivos, totalArchivos, err := cierrelote.LeerArchivoLoteExterno(ctxaws, config.DIR_KEY)
// 		if err == nil {
// 			if totalArchivos != 0 {
// 				/* se obtiene todos los estados externos de prisma */
// 				filtro := filtros.PagoEstadoExternoFiltro{
// 					Vendor:           strings.ToUpper("prisma"),
// 					CargarEstadosInt: true,
// 				}
// 				estadosPagoExterno, err := service.GetPagosEstadosExternoService(filtro)
// 				if err != nil {
// 					errObtenerEstados := errors.New("error al solicitar lista de estados de pago")
// 					err = errObtenerEstados
// 					logError := entities.Log{
// 						Tipo:          entities.EnumLog("error"),
// 						Funcionalidad: "GetPagosEstadosExternoService",
// 						Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 					}
// 					errCrearLog := service.CreateLogService(logError)
// 					if errCrearLog != nil {
// 						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 					}
// 				} else {
// 					/* paso 2: */
// 					listaArchivo, err := cierrelote.LeerCierreLoteTxt(archivos, rutaArchivos, estadosPagoExterno)
// 					if err != nil {
// 						logs.Error(err)
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionCierreLote,
// 							Descripcion: fmt.Sprintf("No se pudo realizar el cierre de lote de prisma. %s", err),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 					if len(listaArchivo) != 0 {
// 						/* paso 3: */
// 						countArchivos, err := cierrelote.MoverArchivos(ctxaws, rutaArchivos, listaArchivo)
// 						if err != nil {
// 							logs.Error(err)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al borrar los archivos temporales: %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						}
// 						/* paso 4: */

// 						err = cierrelote.BorrarArchivos(ctxaws, config.DIR_KEY, rutaArchivos, listaArchivo)
// 						if err != nil {
// 							logs.Error(err)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al borrar los archivos temporales: %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						}
// 						var notificacion entities.Notificacione
// 						if countArchivos > 0 {
// 							notificacion = entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("fecha: %v - Se procesaron %v archivos de cierre de lote Prisma, recibido por: SFTP", time.Now().String(), countArchivos),
// 							}
// 						} else {
// 							notificacion = entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("fecha: %v - No existe movimientos de cierre de lote Prisma: %s FTP", time.Now().String(), err),
// 							}
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					} else {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionCierreLote,
// 							Descripcion: fmt.Sprintf("fecha: %v - No existe archivos de cierre de lote Prisma: %s FTP", time.Now().String(), err),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 				}
// 			}
// 		} else {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("fecha: %v - No se pudo realizar cierre de lote Prisma: %s FTP", time.Now().String(), err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 	})

// 	// /* TODO -> INICIO PROCESAR TABLA MOVIMIENTOS MX*/
// 	confProcesarTablaMx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESAR_TABLA_MX", "0 00 14 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de procesar tabla mx")
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confProcesarTablaMx.Valor, func() {

// 		movimientoMx, movimientoMxEntity, err := cierrelote.ObtenerMxMoviminetosServices()
// 		if err != nil {
// 			errObtenerEstados := errors.New("error al obtener registros de la tablas movimientos mx")
// 			err = errObtenerEstados
// 			logError := entities.Log{
// 				Tipo:          entities.EnumLog("error"),
// 				Funcionalidad: "ObtenerMxMoviminetosServices",
// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 			}
// 			errCrearLog := service.CreateLogService(logError)
// 			if errCrearLog != nil {
// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 			}
// 		} else {
// 			tablasRelaciondas, err := cierrelote.ObtenerTablasRelacionadasServices()
// 			if err != nil {
// 				errObtenerEstados := errors.New("al obtener las tablas relacionadas" + err.Error())
// 				err = errObtenerEstados
// 				logError := entities.Log{
// 					Tipo:          entities.EnumLog("error"),
// 					Funcionalidad: "ObtenerTablasRelacionadasServices",
// 					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 				}
// 				errCrearLog := service.CreateLogService(logError)
// 				if errCrearLog != nil {
// 					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 				}
// 			} else {
// 				resultadoMovimientoMx := cierrelote.ProcesarMovimientoMxServices(movimientoMx, tablasRelaciondas)
// 				if len(resultadoMovimientoMx) <= 0 {
// 					errObtenerEstados := errors.New("error: procesar movimiento mx se encuentra vacia")
// 					err = errObtenerEstados
// 					logError := entities.Log{
// 						Tipo:          entities.EnumLog("error"),
// 						Funcionalidad: "ProcesarMovimientoMxServices",
// 						Mensaje:       errObtenerEstados.Error(),
// 					}
// 					errCrearLog := service.CreateLogService(logError)
// 					if errCrearLog != nil {
// 						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 					}
// 				} else {
// 					err = cierrelote.SaveMovimientoMxServices(resultadoMovimientoMx, movimientoMxEntity)
// 					if err != nil {
// 						errObtenerEstados := errors.New("error al guardar los movimientos: " + err.Error())
// 						err = errObtenerEstados
// 						logError := entities.Log{
// 							Tipo:          entities.EnumLog("error"),
// 							Funcionalidad: "SaveMovimientoMxServices",
// 							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 						}
// 						errCrearLog := service.CreateLogService(logError)
// 						if errCrearLog != nil {
// 							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 						}
// 					} else {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionProcesoMx,
// 							Descripcion: fmt.Sprintf("fecha : %v - procesamiento de movimientos mx se realizo con exito", time.Now().String()),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 				}
// 			}
// 		}
// 	})

// 	// /* TODO -> INICIO PROCESAR TABLA PAGOS PX*/
// 	confProcesarTablaPx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESAR_TABLA_PX", "0 00 14 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de procesar tabla px")

// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confProcesarTablaPx.Valor, func() {

// 		pagosPx, entityPagoPxStr, err := cierrelote.ObtenerPxPagosServices()
// 		if err != nil {
// 			errObtenerEstados := errors.New("error al obtener registros de la tablas pagos px")
// 			err = errObtenerEstados
// 			logError := entities.Log{
// 				Tipo:          entities.EnumLog("error"),
// 				Funcionalidad: "ObtenerPxPagosServices",
// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 			}
// 			errCrearLog := service.CreateLogService(logError)
// 			if errCrearLog != nil {
// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 			}
// 		} else {
// 			err = cierrelote.SavePagoPxServices(pagosPx, entityPagoPxStr)
// 			if err != nil {
// 				errObtenerEstados := errors.New("error al guardar liquidacion de prisma")
// 				err = errObtenerEstados
// 				logError := entities.Log{
// 					Tipo:          entities.EnumLog("error"),
// 					Funcionalidad: "SavePagoPxServices",
// 					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 				}
// 				errCrearLog := service.CreateLogService(logError)
// 				if errCrearLog != nil {
// 					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 				}
// 			} else {
// 				notificacion := entities.Notificacione{
// 					Tipo:        entities.NotificacionProcesoPx,
// 					Descripcion: fmt.Sprintf("fecha : %v - procesamiento de pagos px se realizo con exito", time.Now().String()),
// 				}
// 				service.CreateNotificacionService(notificacion)
// 			}
// 		}
// 	})

// 	// Metodos Sebas 4 Filtros automaticos comentados
// 	// agregar a importaciones filtroCl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
// 	// /* TODO -> INICIO PROCESO CONCILIAR CIERRE LOTE Y MOVIMIENTOS PRISMA*/
// 	confProcesoConciliarClMx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_MX", "0 00 15 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-Mx")

// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc(confProcesoConciliarClMx.Valor, func() {

// 		filtro_clmx_compras := filtroCl.FiltroCierreLote{
// 			MatchCl:       true,
// 			MovimientosMX: true,
// 			PagosPx:       true,
// 			Banco:         true,
// 			Compras:       true,
// 			Devolucion:    false,
// 			Anulacion:     false,
// 			ContraCargo:   false,
// 			ContraCargoMx: false,
// 			ContraCargoPx: false,
// 			Reversion:     false,
// 		}

// 		filtro_clmx_devolucion := filtroCl.FiltroCierreLote{
// 			MatchCl:       true,
// 			MovimientosMX: true,
// 			PagosPx:       true,
// 			Banco:         true,
// 			Compras:       false,
// 			Devolucion:    true,
// 			Anulacion:     false,
// 			ContraCargo:   false,
// 			ContraCargoMx: false,
// 			ContraCargoPx: false,
// 			Reversion:     false,
// 		}

// 		filtro_clmx_anulacion := filtroCl.FiltroCierreLote{
// 			MatchCl:       true,
// 			MovimientosMX: true,
// 			PagosPx:       true,
// 			Banco:         true,
// 			Compras:       false,
// 			Devolucion:    false,
// 			Anulacion:     true,
// 			ContraCargo:   false,
// 			ContraCargoMx: false,
// 			ContraCargoPx: false,
// 			Reversion:     false,
// 		}

// 		filtro_clmx_cc := filtroCl.FiltroCierreLote{
// 			MatchCl:       false,
// 			MovimientosMX: false,
// 			PagosPx:       false,
// 			Banco:         false,
// 			Compras:       true,
// 			Devolucion:    false,
// 			Anulacion:     false,
// 			ContraCargo:   true,
// 			ContraCargoMx: false,
// 			ContraCargoPx: false,
// 			Reversion:     false,
// 		}

// 		var filtros_automaticos []filtroCl.FiltroCierreLote

// 		filtros_automaticos = append(filtros_automaticos, filtro_clmx_compras, filtro_clmx_devolucion, filtro_clmx_anulacion, filtro_clmx_cc)
// 		filtros_nombres := []string{"filtro compras", "filtro devolucion", "filtro anulacion", "filtro contracargo"}

// 		for i_loop, filtro_loop := range filtros_automaticos {
// 			filtro_aplicado := filtro_loop

// 			var filtro filtroCl.FiltroPrismaMovimiento
// 			filtro = filtroCl.FiltroPrismaMovimiento{
// 				Match:                        false,
// 				CargarDetalle:                true,
// 				Contracargovisa:              true,
// 				Contracargomaster:            true,
// 				Tipooperacion:                true,
// 				Rechazotransaccionprincipal:  true,
// 				Rechazotransaccionsecundario: true,
// 				Motivoajuste:                 true,
// 				ContraCargo:                  filtro_aplicado.ContraCargo,
// 				CodigosOperacion:             []string{"0005"},
// 				TipoAplicacion:               "+",
// 			}
// 			if filtro_aplicado.ContraCargo {
// 				filtro.Match = false
// 				filtro.CodigosOperacion = []string{"1507", "6000", "1517"}
// 				filtro.TipoAplicacion = "-"
// 			}
// 			listaMovimientos, codigoautorizacion, err := cierrelote.ObtenerPrismaMovimientosServices(filtro)
// 			if err != nil {
// 				errObtenerEstados := errors.New(err.Error())
// 				err = errObtenerEstados
// 				logError := entities.Log{
// 					Tipo:          entities.EnumLog("error"),
// 					Funcionalidad: "ObtenerPrismaMovimientosServices",
// 					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 				}
// 				errCrearLog := service.CreateLogService(logError)
// 				if errCrearLog != nil {
// 					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 				}
// 			} else {

// 				listaCierreLote, err := cierrelote.ObtenerCierreloteServices(filtro_aplicado, codigoautorizacion)
// 				if err != nil {
// 					errObtenerEstados := errors.New(err.Error())
// 					err = errObtenerEstados
// 					logError := entities.Log{
// 						Tipo:          entities.EnumLog("error"),
// 						Funcionalidad: "ObtenerCierreloteServices",
// 						Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 					}
// 					errCrearLog := service.CreateLogService(logError)
// 					if errCrearLog != nil {
// 						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 					}
// 				} else {

// 					listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := cierrelote.ConciliarCierreLotePrismaMovimientoServices(listaCierreLote, listaMovimientos)
// 					if err != nil {
// 						errObtenerEstados := errors.New(err.Error())
// 						err = errObtenerEstados
// 						logError := entities.Log{
// 							Tipo:          entities.EnumLog("error"),
// 							Funcionalidad: "ConciliarCierreLotePrismaMovimientoServices",
// 							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 						}
// 						errCrearLog := service.CreateLogService(logError)
// 						if errCrearLog != nil {
// 							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 						}
// 					} else {

// 						if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionConciliacionCLMx,
// 								Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con movimientos %v", time.Now().String(), filtros_nombres[i_loop]),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						} else {

// 							logs.Info("en end-point")
// 							logs.Info(listaCierreLoteProcesado)
// 							err = cierrelote.ActualizarCierreloteMoviminetosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
// 							if err != nil {
// 								errObtenerEstados := errors.New(err.Error())
// 								err = errObtenerEstados
// 								logError := entities.Log{
// 									Tipo:          entities.EnumLog("error"),
// 									Funcionalidad: "ActualizarCierreloteMoviminetosServices",
// 									Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 								}
// 								errCrearLog := service.CreateLogService(logError)
// 								if errCrearLog != nil {
// 									logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 								}
// 							} else {
// 								notificacion := entities.Notificacione{
// 									Tipo:        entities.NotificacionConciliacionCLMx,
// 									Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion movimientos con cierre lote exito  %v", time.Now().String(), filtros_nombres[i_loop]),
// 								}
// 								service.CreateNotificacionService(notificacion)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 	})

// 	// /* TODO -> INICIO PROCESO CONCILIAR CIERRE LOTE Y PAGOS PRISMA*/
// 	// confProcesoConciliarClPx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_PX", "0 35 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-px")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c.AddFunc(confProcesoConciliarClPx.Valor, func() {
// 	// 	filtro_clpx_compra := filtroCl.FiltroCierreLote{
// 	// 		MatchCl:         true,
// 	// 		MovimientosMX:   false,
// 	// 		PagosPx:         true,
// 	// 		Banco:           true,
// 	// 		EstadoFechaPago: true,
// 	// 		FechaPago:       "0000-00-00",
// 	// 		Compras:         true,
// 	// 		Devolucion:      false,
// 	// 		Anulacion:       false,
// 	// 		ContraCargo:     false,
// 	// 		ContraCargoMx:   false,
// 	// 		ContraCargoPx:   false,
// 	// 		Reversion:       false,
// 	// 	}

// 	// 	filtro_clpx_devolucion := filtroCl.FiltroCierreLote{
// 	// 		MatchCl:         true,
// 	// 		MovimientosMX:   false,
// 	// 		PagosPx:         true,
// 	// 		Banco:           true,
// 	// 		EstadoFechaPago: true,
// 	// 		FechaPago:       "0000-00-00",
// 	// 		Compras:         false,
// 	// 		Devolucion:      true,
// 	// 		Anulacion:       false,
// 	// 		ContraCargo:     false,
// 	// 		ContraCargoMx:   false,
// 	// 		ContraCargoPx:   false,
// 	// 		Reversion:       false,
// 	// 	}

// 	// 	filtro_clpx_anulacion := filtroCl.FiltroCierreLote{
// 	// 		MatchCl:         true,
// 	// 		MovimientosMX:   false,
// 	// 		PagosPx:         true,
// 	// 		Banco:           true,
// 	// 		EstadoFechaPago: true,
// 	// 		FechaPago:       "0000-00-00",
// 	// 		Compras:         false,
// 	// 		Devolucion:      false,
// 	// 		Anulacion:       true,
// 	// 		ContraCargo:     false,
// 	// 		ContraCargoMx:   false,
// 	// 		ContraCargoPx:   false,
// 	// 		Reversion:       false,
// 	// 	}

// 	// 	filtro_clpx_cc := filtroCl.FiltroCierreLote{
// 	// 		MatchCl:         true,
// 	// 		MovimientosMX:   false,
// 	// 		PagosPx:         true,
// 	// 		Banco:           true,
// 	// 		EstadoFechaPago: true,
// 	// 		FechaPago:       "0000-00-00",
// 	// 		Compras:         false,
// 	// 		Devolucion:      false,
// 	// 		Anulacion:       false,
// 	// 		ContraCargo:     true,
// 	// 		ContraCargoMx:   false,
// 	// 		ContraCargoPx:   false,
// 	// 		Reversion:       false,
// 	// 	}

// 	// 	var filtros_clpx_automaticos []filtroCl.FiltroCierreLote

// 	// 	filtros_clpx_automaticos = append(filtros_clpx_automaticos, filtro_clpx_compra, filtro_clpx_devolucion, filtro_clpx_anulacion, filtro_clpx_cc)
// 	// 	filtros_clpx_nombres := []string{"filtro compras", "filtro devolucion", "filtro anulacion", "filtro contracargo"}

// 	// 	for i_loop, filtro_loop := range filtros_clpx_automaticos {
// 	// 		filtro_clpx_aplicado := filtro_loop

// 	// 		var filtro filtroCl.FiltroPrismaMovimiento

// 	// 		filtro = filtroCl.FiltroPrismaMovimiento{
// 	// 			Match:                        false,
// 	// 			CargarDetalle:                true,
// 	// 			Contracargovisa:              true,
// 	// 			Contracargomaster:            true,
// 	// 			Tipooperacion:                true,
// 	// 			Rechazotransaccionprincipal:  true,
// 	// 			Rechazotransaccionsecundario: true,
// 	// 			Motivoajuste:                 true,
// 	// 			ContraCargo:                  filtro_clpx_aplicado.ContraCargo,
// 	// 			CodigosOperacion:             []string{"0005"},
// 	// 			TipoAplicacion:               "+",
// 	// 		}
// 	// 		if filtro_clpx_aplicado.ContraCargo {
// 	// 			filtro.Match = false
// 	// 			filtro.CodigosOperacion = []string{"1507", "6000"}
// 	// 			filtro.TipoAplicacion = "-"
// 	// 		}

// 	// 		var codigo []string
// 	// 		listaCierreLote, err := cierrelote.ObtenerCierreloteServices(filtro_clpx_aplicado, codigo)
// 	// 		if err != nil {
// 	// 			errObtenerEstados := errors.New(err.Error())
// 	// 			err = errObtenerEstados
// 	// 			logError := entities.Log{
// 	// 				Tipo:          entities.EnumLog("error"),
// 	// 				Funcionalidad: "ObtenerCierreloteServices",
// 	// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 			}
// 	// 			errCrearLog := service.CreateLogService(logError)
// 	// 			if errCrearLog != nil {
// 	// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 			}
// 	// 		} else {

// 	// 			if len(listaCierreLote) == 0 {
// 	// 				notificacion := entities.Notificacione{
// 	// 					Tipo:        entities.NotificacionConciliacionCLPx,
// 	// 					Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con pagos %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
// 	// 				}
// 	// 				service.CreateNotificacionService(notificacion)
// 	// 			} else {

// 	// 				filtroCabecera := filtroCl.FiltroPrismaMovimiento{
// 	// 					ContraCargo: filtro_clpx_aplicado.ContraCargo,
// 	// 				}
// 	// 				listaCierreLoteMovimientos, err := cierrelote.ObtenerPrismaMovimientoConciliadosServices(listaCierreLote, filtroCabecera)
// 	// 				if err != nil {
// 	// 					errObtenerEstados := errors.New(err.Error())
// 	// 					err = errObtenerEstados
// 	// 					logError := entities.Log{
// 	// 						Tipo:          entities.EnumLog("error"),
// 	// 						Funcionalidad: "ObtenerPrismaMovimientoConciliadosServices",
// 	// 						Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 					}
// 	// 					errCrearLog := service.CreateLogService(logError)
// 	// 					if errCrearLog != nil {
// 	// 						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 					}
// 	// 				} else {

// 	// 					var listaFechaPagos []string
// 	// 					for _, value := range listaCierreLoteMovimientos {
// 	// 						fechaString := value.MovimientoCabecer.FechaPago.Format("2006-01-02")
// 	// 						listaFechaPagos = append(listaFechaPagos, fechaString)
// 	// 					}
// 	// 					filtro_prisma := filtroCl.FiltroPrismaTrPagos{
// 	// 						Match:         false,
// 	// 						CargarDetalle: true,
// 	// 						Devolucion:    filtro_clpx_aplicado.Devolucion,
// 	// 						FechaPagos:    listaFechaPagos,
// 	// 					}
// 	// 					listaPrismaPago, err := cierrelote.ObtenerPrismaPagosServices(filtro_prisma)
// 	// 					if err != nil {
// 	// 						errObtenerEstados := errors.New(err.Error())
// 	// 						err = errObtenerEstados
// 	// 						logError := entities.Log{
// 	// 							Tipo:          entities.EnumLog("error"),
// 	// 							Funcionalidad: "ObtenerPrismaPagosServices",
// 	// 							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 						}
// 	// 						errCrearLog := service.CreateLogService(logError)
// 	// 						if errCrearLog != nil {
// 	// 							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 						}
// 	// 					} else {
// 	// 						listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := cierrelote.ConciliarCierreLotePrismaPagoServices(listaCierreLoteMovimientos, listaPrismaPago)
// 	// 						if err != nil {
// 	// 							errObtenerEstados := errors.New(err.Error())
// 	// 							err = errObtenerEstados
// 	// 							logError := entities.Log{
// 	// 								Tipo:          entities.EnumLog("error"),
// 	// 								Funcionalidad: "ConciliarCierreLotePrismaPagoServices",
// 	// 								Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 							}
// 	// 							errCrearLog := service.CreateLogService(logError)
// 	// 							if errCrearLog != nil {
// 	// 								logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 							}
// 	// 						} else {

// 	// 							if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
// 	// 								notificacion := entities.Notificacione{
// 	// 									Tipo:        entities.NotificacionConciliacionCLPx,
// 	// 									Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con pagos %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
// 	// 								}
// 	// 								service.CreateNotificacionService(notificacion)
// 	// 							}
// 	// 							err = cierrelote.ActualizarCierrelotePagosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
// 	// 							if err != nil {
// 	// 								errObtenerEstados := errors.New(err.Error())
// 	// 								err = errObtenerEstados
// 	// 								logError := entities.Log{
// 	// 									Tipo:          entities.EnumLog("error"),
// 	// 									Funcionalidad: "ActualizarCierrelotePagosServices",
// 	// 									Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 								}
// 	// 								errCrearLog := service.CreateLogService(logError)
// 	// 								if errCrearLog != nil {
// 	// 									logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 								}
// 	// 							} else {
// 	// 								notificacion := entities.Notificacione{
// 	// 									Tipo:        entities.NotificacionConciliacionCLPx,
// 	// 									Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion pagos con cierre lote exito %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
// 	// 								}
// 	// 								service.CreateNotificacionService(notificacion)
// 	// 							}
// 	// 						}

// 	// 					}

// 	// 				}
// 	// 			}
// 	// 		}

// 	// 	}

// 	// })

// 	// // /* TODO -> INICIO PROCESO CONCILIAR CON EL BANCO*/
// 	// confProcesoConciliarClBanco, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_BANCO", "0 40 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-banco")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// c.AddFunc(confProcesoConciliarClBanco.Valor, func() {
// 	// 	fechaActual, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
// 	// 	if err != nil {
// 	// 		errObtenerEstados := errors.New("error al parsear fecha " + err.Error())
// 	// 		err = errObtenerEstados
// 	// 		logError := entities.Log{
// 	// 			Tipo:          entities.EnumLog("error"),
// 	// 			Funcionalidad: "parseo de fecha ",
// 	// 			Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 		}
// 	// 		errCrearLog := service.CreateLogService(logError)
// 	// 		if errCrearLog != nil {
// 	// 			logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 		}
// 	// 	} else {
// 	// 		fechatemporal := fechaActual.Add(24 * -1)
// 	// 		fechaPagoProcesar := fechatemporal.Format("2006-01-02")
// 	// 		logs.Info(fechaPagoProcesar)
// 	// 		filtro := filtroCl.FiltroTablasConciliadas{
// 	// 			FechaPago: fechaPagoProcesar,
// 	// 			Match:     true,
// 	// 			Reversion: true,
// 	// 		}
// 	// 		responseListprismaTrPagos, err := cierrelote.ObtenerRepoPagosPrisma(filtro)
// 	// 		if err != nil {
// 	// 			errObtenerEstados := errors.New(err.Error())
// 	// 			err = errObtenerEstados
// 	// 			logError := entities.Log{
// 	// 				Tipo:          entities.EnumLog("error"),
// 	// 				Funcionalidad: "ActualizarCierreloteMoviminetosServices",
// 	// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 			}
// 	// 			errCrearLog := service.CreateLogService(logError)
// 	// 			if errCrearLog != nil {
// 	// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 			}
// 	// 		}
// 	// 		/* obtener configuracion periodo de acreditacion */
// 	// 		movimientoBanco, erro := cierrelote.ConciliacionBancoPrisma(fechaPagoProcesar, filtro.Reversion, responseListprismaTrPagos)
// 	// 		if erro != nil {
// 	// 			errObtenerEstados := errors.New(err.Error())
// 	// 			err = errObtenerEstados
// 	// 			logError := entities.Log{
// 	// 				Tipo:          entities.EnumLog("error"),
// 	// 				Funcionalidad: "ActualizarCierreloteMoviminetosServices",
// 	// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 			}
// 	// 			errCrearLog := service.CreateLogService(logError)
// 	// 			if errCrearLog != nil {
// 	// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 			}
// 	// 		}
// 	// 		if movimientoBanco == nil {
// 	// 			notificacion := entities.Notificacione{
// 	// 				Tipo:        entities.NotificacionConciliacionBancoCL,
// 	// 				Descripcion: fmt.Sprintf("fecha : %v - no existe movimientos en banco para conciliar con los pagos ", time.Now().String()),
// 	// 			}
// 	// 			service.CreateNotificacionService(notificacion)

// 	// 		} else {
// 	// 			notificacion := entities.Notificacione{
// 	// 				Tipo:        entities.NotificacionConciliacionBancoCL,
// 	// 				Descripcion: fmt.Sprintf("fecha : %v -  proceso de conciliacion pagos con banco exitoso ", time.Now().String()),
// 	// 			}
// 	// 			service.CreateNotificacionService(notificacion)
// 	// 		}
// 	// 	}
// 	// })

// 	// /* TODO -> INICIO PROCESO GENERAR MOVIMIENTOS */
// 	// confProcesoBuildMovimiento, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_PAGOS_PSP", "0 40 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que concilia cl con pagos pasarela")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// c.AddFunc(confProcesoBuildMovimiento.Valor, func() {
// 	// 	reversion := true
// 	// 	responseCierreLote, err := service.BuildPrismaMovimiento(reversion)
// 	// 	if err != nil {
// 	// 		errObtenerEstados := errors.New(err.Error())
// 	// 		err = errObtenerEstados
// 	// 		logError := entities.Log{
// 	// 			Tipo:          entities.EnumLog("error"),
// 	// 			Funcionalidad: "BuildPrismaMovimiento",
// 	// 			Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 		}
// 	// 		errCrearLog := service.CreateLogService(logError)
// 	// 		if errCrearLog != nil {
// 	// 			logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 		}
// 	// 	} else {

// 	// 		ctxPrueba := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
// 	// 		err = service.CreateMovimientosService(ctxPrueba, responseCierreLote)
// 	// 		if err != nil {
// 	// 			errObtenerEstados := errors.New(err.Error())
// 	// 			err = errObtenerEstados
// 	// 			logError := entities.Log{
// 	// 				Tipo:          entities.EnumLog("error"),
// 	// 				Funcionalidad: "CreateMovimientosService",
// 	// 				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
// 	// 			}
// 	// 			errCrearLog := service.CreateLogService(logError)
// 	// 			if errCrearLog != nil {
// 	// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 	// 			}
// 	// 		} else {
// 	// 			notificacion := entities.Notificacione{
// 	// 				Tipo:        entities.NotificacionConciliacionCLPx,
// 	// 				Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion pagos con cierre lote exito ", time.Now().String()),
// 	// 			}
// 	// 			service.CreateNotificacionService(notificacion)
// 	// 		}
// 	// 	}
// 	// })

// 	/*
// 		PROCESO 2:
// 			se realiza conciliacion de los archivos de cierre de lotes con los movimientos de banco
// 		PROCESO 3:
// 			se crean los diferentes objetos para registrar los movimineto recibidos en el cierre de lote
// 	*/

// 	// TODO INICIO PROCESO RAPIPAGO
// 	// & 1 Se procesan los archivos se guardan en la tabla rapipago
// 	// & 2 Se actualizan estados de los pagos con los encontrado en cierrelote(archivo recibido) -> EL estado APROBADO indica que el pagador fue a un rapipago
// 	// & 3 Se notifica el cambia de estado al cliente(se ejecuta webhook)
// 	// & 4 Conciliar con los movimientos ingresados en banco
// 	// & 5 Generar movimientos

// 	// & 2
// 	confRapipagoCierreLote, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE2", "@midnight", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
// 	logs.Info(confRapipagoCierreLote)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc("0 00 14 * * *", func() {

// 		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
// 			CargarMovConciliados: false,
// 			PagosNotificado:      false,
// 		}
// 		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote  */
// 		listaPagoaRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago no se puede continuar. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 		if len(listaPagoaRapipago) > 0 {
// 			if err == nil {
// 				listaPagosClRapipago, err := service.BuildPagosClRapipago(listaPagoaRapipago)
// 				if err == nil {
// 					// Actualizar estados del pago y cierrelote
// 					logs.Info("inicio actualizacion de pagos rapipago")
// 					err = service.ActualizarPagosClRapipagoService(listaPagosClRapipago)
// 					if err != nil {
// 						logs.Error(err)
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionCierreLote,
// 							Descripcion: fmt.Sprintln("error al actualizar estados de los pagos "),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 				}
// 			}

// 		} else {
// 			logs.Info("no existen pagos de rapipago para actualizar")
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintln("No existen pagos de rapipago para actualizar"),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})

// 	// & 3
// 	confRapipagoCierreLoteNotificarPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE3", "@midnight", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
// 	logs.Info(confRapipagoCierreLoteNotificarPagos)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc("0 00 15 * * *", func() {

// 		request := filtros.PagoEstadoFiltro{
// 			EstadoId: 4,
// 		}
// 		pagos, err := service.BuildNotificacionPagosCLRapipago(request)

// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago no se puede continuar. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 		if len(pagos) > 0 {
// 			pagosNotificar := service.NotificarPagos(pagos)
// 			if len(pagosNotificar) > 0 { /* actualzar estado de pagos a notificado */
// 				mensaje := fmt.Sprintf("los siguientes pagos se actualizaron correctamente. %v", pagosNotificar)
// 				logs.Info(mensaje)
// 			} else {
// 				notificacion := entities.Notificacione{
// 					Tipo:        entities.NotificacionWebhook,
// 					Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
// 				}
// 				service.CreateNotificacionService(notificacion)
// 			}

// 		} else {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintln("no existen pagos por notificar"),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})

// 	// & 4
// 	confRapipagoCierreLoteConciliarBanco, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE4", "@midnight", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
// 	logs.Info(confRapipagoCierreLoteConciliarBanco)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc("0 00 05 * * *", func() {

// 		filtroMovConciliarRapipago := rapipago.RequestConsultarMovimientosRapipago{
// 			CargarMovConciliados: false,
// 			PagosNotificado:      true,
// 		}
// 		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
// 		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovConciliarRapipago)
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 		if len(listaCierreRapipago) > 0 {
// 			if err == nil {

// 				request := bancodtos.RequestConciliacion{
// 					TipoConciliacion: 1,
// 					ListaRapipago:    listaCierreRapipago,
// 				}
// 				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
// 				listaCierreRapipago, listaBancoId, err := movimientosBanco.ConciliacionPasarelaBanco(request)

// 				if len(listaBancoId) == 0 {
// 					logs.Info("no existen movimientos en banco para conciliar con pagos rapipago")
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionCierreLote,
// 						Descripcion: fmt.Sprintf("no existen movimientos en banco para conciliar con pagos rapipago: %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				} else {
// 					/*en el caso de error a actualizar la tabla rapipagocierrelote el proceso termina */
// 					err := service.UpdateCierreLoteRapipago(listaCierreRapipago.ListaRapipago)
// 					if err != nil {
// 						logs.Error(err)
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionCierreLote,
// 							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote rapipago (volver a ejecutar proceso): %s", err),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					} else {
// 						// actualiza registro movimientos del banco
// 						// si no se actualiza los registros del banco se debera actualizar manualmente
// 						_, err := movimientosBanco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
// 						if err != nil {
// 							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
// 							logs.Error(err)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al actualizar movimientos del banco - conciliacion rapipago(actualizar manualmente los siguientes movimientos): %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						}
// 					}

// 				}

// 			}

// 		} else {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintln("No existen pagos de rapipago por conciliar"),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})

// 	// & 5
// 	confRapipagoCierreLoteGenerarMovimientos, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE4", "0 00 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
// 	logs.Info(confRapipagoCierreLoteGenerarMovimientos)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc("0 00 06 * * *", func() {

// 		filtroMovMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
// 			CargarMovConciliados: true,
// 			PagosNotificado:      true,
// 		}

// 		listaCierreMovRapipago, err := service.GetCierreLoteRapipagoService(filtroMovMovRapipago)
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo obtener los pagos para generar movimientos. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 		// Si no se guarda ningún cierre no hace falta seguir el proce
// 		if len(listaCierreMovRapipago) > 0 {
// 			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
// 			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
// 			responseCierreLote, err := service.BuildRapipagoMovimiento(listaCierreMovRapipago)

// 			if err == nil {

// 				// 3 - Guarda los movimientos en la base de datos en caso de error se
// 				// repetira en el día siguiente
// 				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
// 				err = service.CreateMovimientosService(ctx, responseCierreLote)
// 				if err != nil {
// 					logs.Error(err)
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionCierreLote,
// 						Descripcion: fmt.Sprintf("No se pudo crear los movimientos clrapipago. %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}

// 			}

// 		} else {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos clrapipago"),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 	})

// 	// ^ PASO 1 y PASO 2: se establecen franjas horarias de ejecucion de proceso conciliacion pagos con apilink
// 	// 0 */5 7-16 * * * : de 7 a 16 hs se ejecuta cada 5 minutos
// 	// 0 0 */1 17-22 * * : de 17 a 22 hs se ajecuta cada 1 hora
// 	confApilinkCierreLoteMorning, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_MORNING", "0 0 */3 7-12 * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confApilinkCierreLoteMorning.Valor, func() {
// 		_myFunctConciliacionPagosApilink(service)
// 	})

// 	confApilinkCierreLoteAfternoon, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_AFTERNOON", "0 0 */5 13-22 * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confApilinkCierreLoteAfternoon.Valor, func() {
// 		_myFunctConciliacionPagosApilink(service)
// 	})

// 	confApilinkCierreNight, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_NIGHT", "0 00 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confApilinkCierreNight.Valor, func() {
// 		_myFunctConciliacionPagosApilink(service)
// 	})

// 	// ^ PASO 3: Conciliacion pagos debin con banco
// 	confConciliacionBancoDebin, err := _buildPeriodicidad(util, "PERIODICIDAD_CONCILIACION_PAGOSDEBIN_BANCO", "0 30 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
// 	if err != nil {
// 		logs.Info(err)
// 		panic(err)
// 	}
// 	c.AddFunc(confConciliacionBancoDebin.Valor, func() {
// 		filtro := linkdebin.RequestDebines{
// 			BancoExternalId:  false,
// 			CargarPagoEstado: true,
// 		}
// 		debines, err := service.GetDebines(filtro)
// 		if err != nil {
// 			errorBuildNotificacion := errors.New("error al obtener debines para conciliar con banco")
// 			err = errorBuildNotificacion
// 			logError := entities.Log{
// 				Tipo:          entities.EnumLog("error"),
// 				Funcionalidad: "GetConsultarDebines",
// 				Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
// 			}
// 			errCrearLog := service.CreateLogService(logError)
// 			if errCrearLog != nil {
// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 			}
// 		}
// 		if len(debines) > 0 {
// 			if err == nil {
// 				request := bancodtos.RequestConciliacion{
// 					TipoConciliacion: 2,
// 					ListaApilink:     debines,
// 				}
// 				listaCierreApiLinkBanco, listaBancoId, err := movimientosBanco.ConciliacionPasarelaBanco(request)
// 				if err != nil {
// 					logs.Error(err)
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionCierreLote,
// 						Descripcion: fmt.Sprintf("error al conciliar movimiento banco y cierre loteapilink: %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)

// 				} else {
// 					// NOTE Actualizar lista de cierreloteapilink campo banco external_id, match y fecha de acreditacion
// 					if len(listaCierreApiLinkBanco.ListaApilink) > 0 || len(listaCierreApiLinkBanco.ListaApilinkNoAcreditados) > 0 {
// 						listas := linkdebin.RequestListaUpdateDebines{
// 							Debines:              listaCierreApiLinkBanco.ListaApilink,
// 							DebinesNoAcreditados: listaCierreApiLinkBanco.ListaApilinkNoAcreditados,
// 						}
// 						erro := service.UpdateCierreLoteApilink(listas)
// 						if erro != nil {
// 							logs.Error(erro)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink y conciliacion con banco: %s", erro),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						}
// 					}
// 					// FIXME se debe verificar si las 2 listas son iguales ?
// 					if len(listaBancoId) > 0 {
// 						_, err := movimientosBanco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
// 						if err != nil {
// 							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
// 							logs.Error(err)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al actualizar registros del banco: %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 							// en le caso de este error y si el pago no se actualizo a estados finales no afecta el cierre de apilink
// 							// el estado del pago se actualiza a estado final y no tendra en cuenta al consultar a apilink
// 							// ACCION : se debe actualizar manualmente el campo check en la tabla de movimientos de banco(no es obligatorio)
// 						}
// 					}

// 				}

// 			}

// 		}
// 	})

// 	// ^ PASO 4: Generar movimientos debines
// 	confGenerarMovimientosDebines, err := _buildPeriodicidad(util, "PERIODICIDAD_GENERAR_MOV_DEBINES", "0 00 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
// 	if err != nil {
// 		logs.Info(err)
// 		c.Stop()
// 	}
// 	c.AddFunc(confGenerarMovimientosDebines.Valor, func() {
// 		filtro := linkdebin.RequestDebines{
// 			BancoExternalId:  true,
// 			CargarPagoEstado: true,
// 		}
// 		debines, err := service.GetDebines(filtro)
// 		if err != nil {
// 			errorBuildNotificacion := errors.New("error al obtener debines para generar movimientos")
// 			err = errorBuildNotificacion
// 			logError := entities.Log{
// 				Tipo:          entities.EnumLog("error"),
// 				Funcionalidad: "GetDebines",
// 				Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
// 			}
// 			errCrearLog := service.CreateLogService(logError)
// 			if errCrearLog != nil {
// 				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 			}
// 		}
// 		if len(debines) > 0 {
// 			responseCierreLote, err := service.BuildMovimientoApiLink(debines)
// 			if err == nil {
// 				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
// 				err = service.CreateMovimientosService(ctx, responseCierreLote)
// 				if err != nil {
// 					logs.Info("error al generar movimientos debines")
// 					logs.Error(err)
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionCierreLote,
// 						Descripcion: fmt.Sprintf("error al generar movimientos debines: %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}
// 			}
// 		}
// 	})
// 	// /*

// 	//  ^ fin CIERRELOTE APILINK

// 	// /* TODO ->CONCILIACION TRANSFERENCIAS CON SERVICIO DE BANCO*/
// 	// // confTransferencias, err := _buildPeriodicidad(util, "PERIODICIDAD_CONCILIACION_TRANSFERENCIA", "@midnight", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de conciliacion de transferencias con movimientos de banco")

// 	// // if err != nil {
// 	// // 	panic(err)
// 	// // }

// 	// // c.AddFunc(confTransferencias.Valor, func() {

// 	// // 	filtros := filtros.TransferenciaFiltro{}
// 	// // 	listaTransferencia, err := service.GetTransferencias(filtros)

// 	// // 	if len(listaTransferencia.Transferencias) > 0 {
// 	// // 		request := bancodtos.RequestConciliacion{
// 	// // 			TipoConciliacion: 3,
// 	// // 			Transferencias:   listaTransferencia,
// 	// // 		}
// 	// // 		listaTransferenciasMatch, listaIdBanco, err := movimientosBanco.ConciliacionPasarelaBanco(request)
// 	// // 		if err != nil {
// 	// // 			logs.Error(err)
// 	// // 			notificacion := entities.Notificacione{
// 	// // 				Tipo:        entities.NotificacionCierreLote,
// 	// // 				Descripcion: fmt.Sprintf("error al conciliar transferencias con movimientos del banco: %s", err),
// 	// // 			}

// 	// // 			service.CreateNotificacionService(notificacion)
// 	// // 		}

// 	// // 		/* SI LA CANTIDAD DE TRANSFERENCIAS MATCH ES MENOR A 0/HAY ERRORES EL PROCESO DE CONCILIACION TERMINA */
// 	// // 		logs.Info("EL total de transferencias que se actualizaran son " + strconv.Itoa(len(listaTransferenciasMatch.Transferencias)))
// 	// // 		if len(listaTransferenciasMatch.Transferencias) > 0 {
// 	// // 			/*ACTUALIZAR CAMPO MATCH Y EXTERNAL_BANCO_ID TABLA TRANSFERENCIA*/
// 	// // 			err := service.UpdateTransferencias(listaTransferenciasMatch)
// 	// // 			if err != nil {
// 	// // 				logs.Error(err)
// 	// // 				notificacion := entities.Notificacione{
// 	// // 					Tipo:        entities.NotificacionCierreLote,
// 	// // 					Descripcion: fmt.Sprintf("error al actualizar transferencias: %s", err),
// 	// // 				}
// 	// // 				service.CreateNotificacionService(notificacion)
// 	// // 			} else {
// 	// // 				/*ACTUALIZAR CAMPO ESTADO_CHECK EN SERVICIO BANCO*/
// 	// // 				response, err := movimientosBanco.ActualizarRegistrosMatchBancoService(listaIdBanco, true)
// 	// // 				logs.Info(response)
// 	// // 				if err != nil {
// 	// // 					logs.Error(err)
// 	// // 					notificacion := entities.Notificacione{
// 	// // 						Tipo:        entities.NotificacionCierreLote,
// 	// // 						Descripcion: fmt.Sprintf("error al actualizar campo check en servicio banco: %s", err),
// 	// // 					}
// 	// // 					service.CreateNotificacionService(notificacion)
// 	// // 				}

// 	// // 			}

// 	// // 		}

// 	// // 	}

// 	// // 	if err != nil {
// 	// // 		notificacion := entities.Notificacione{
// 	// // 			Tipo:        entities.NotificacionPagoExpirado,
// 	// // 			Descripcion: fmt.Sprintf("No se pudo realizar el proceso de conciliacion de transferencia. %s", err),
// 	// // 		}
// 	// // 		service.CreateNotificacionService(notificacion)
// 	// // 	}
// 	// // })

// 	// /* fin CONCILIACION TRANSFERENCIAS CON SERVICIO DE BANCO*/

// 	/* TODO
// 	*WEBHOOK(NOTIFICACION DE PAGOS) A CLIENTES
// 	 */
// 	confNotificacionPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_NOTIFICACION_PAGOS", "0 */12 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de notificacion de pagos a los clientes")
// 	//  @every 12h
// 	logs.Info(confNotificacionPagos)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.AddFunc(confNotificacionPagos.Valor, func() {

// 		filtroWebhook := webhooks.RequestWebhook{
// 			DiasPago:         15,
// 			PagosNotificado:  false,
// 			EstadoFinalPagos: true,
// 		}
// 		pagos, err := service.BuildNotificacionPagosService(filtroWebhook)
// 		if err == nil {
// 			pagosNotificar, err := service.CreateNotificacionPagosService(pagos)
// 			if err == nil {
// 				if len(pagosNotificar) > 0 {
// 					pagosupdate := service.NotificarPagos(pagosNotificar)
// 					if len(pagosupdate) > 0 { /* actualzar estado de pagos a notificado */
// 						err = service.UpdatePagosNoticados(pagosupdate)
// 						if err != nil {
// 							logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
// 							logs.Error(err)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionWebhook,
// 								Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)

// 						}
// 					} else {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionWebhook,
// 							Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}

// 				} else {
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionWebhook,
// 						Descripcion: fmt.Sprintln("webhook: No existen pagos por notificar"),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}
// 			}
// 		}
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de notificacion de pagos. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})
// 	/*
// 	*fin NOTIFICACION DE PAGOS A CLIENTES
// 	 */

// 	// confPagosExpirados, err := _buildPeriodicidad(util, "PERIODICIDAD_UPDATE_PAGOS_EXPIRADOS", "@monthly", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que expira los pagos en estado pending")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c.AddFunc(confPagosExpirados.Valor, func() {
// 	// 	err := service.ModificarEstadoPagosExpirados()
// 	// 	if err != nil {
// 	// 		notificacion := entities.Notificacione{
// 	// 			Tipo:        entities.NotificacionPagoExpirado,
// 	// 			Descripcion: fmt.Sprintf("No se pudo realizar la modificación de estado de los pagos expirados. %s", err),
// 	// 		}
// 	// 		service.CreateNotificacionService(notificacion)
// 	// 	}
// 	// })

// 	// Todos: Ejemplo los martes a las 12 con 58 en punto "0 58 12 * * 2"
// 	// & Retiro automatico se ejecutara todos los dias a las 4 am

// 	/* TODO -> RETIRO AUTOMATICOS: ESTO MODIFICAR SEGUN LA CONFIGUARCION DEL CLIENTE*/
// 	confRetiroAutomatico, err := _buildPeriodicidad(util, "PERIODICIDAD_RETIRO_AUTOMATICO", "0 00 02 * * *", "Periodicidad (en formato cron) en que el sistema realiza automaticamente las transferencias a los clientes.")
// 	logs.Info(confRetiroAutomatico)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Second | Minute | Hour | Dom | Month | DowOptional | Descriptor
// 	// Dom = Day of the month
// 	// DowOptional = Day of the week Opcional
// 	c.AddFunc("0 00 08 * * 1-5", func() {
// 		var feriado bool

// 		/* NO TRANSFERIR FEIRADOS */
// 		filtro := filtros.ConfiguracionFiltro{
// 			Nombre: "FERIADOS",
// 		}

// 		// buscar la configuracion de dias feriados
// 		configuracion, erro := util.GetConfiguracionService(filtro)

// 		// si hay error se notifica pero se continua
// 		if erro != nil {
// 			_buildNotificacion(util, erro, entities.NotificacionConfiguraciones)
// 		}

// 		// si se obtiene el resultado de la configuracion para FERIADOS
// 		if configuracion.Id != 0 {
// 			feriado = _esFeriado(configuracion.Valor)
// 		}

// 		// si NO es feriado, hacer la transferencia
// 		if !feriado {
// 			ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

// 			_, err := service.RetiroAutomaticoClientes(ctx)

// 			if err != nil {
// 				notificacion := entities.Notificacione{
// 					Tipo:        entities.NotificacionPagoExpirado,
// 					Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica para los clientes. %s", err),
// 				}
// 				service.CreateNotificacionService(notificacion)
// 			}

// 			// if len(response.MovimientosId) > 0 {
// 			// 	uuid := uuid.NewV4()
// 			// 	idmovimientos := administraciondtos.RequestMovimientosId{
// 			// 		MovimientosId:             response.MovimientosId,
// 			// 		MovimimientosIdRevertidos: response.MovimimientosIdRevertidos,
// 			// 	}
// 			// 	result := service.SendTransferenciasComisiones(ctx, uuid.String(), idmovimientos)
// 			// 	if !result {
// 			// 		logs.Error(result)
// 			// 		notificacion := entities.Notificacione{
// 			// 			Tipo:        entities.NotificacionTransferencia,
// 			// 			Descripcion: fmt.Sprintf("error al intentar transferir comisiones impuestos telco: %v", result),
// 			// 		}
// 			// 		service.CreateNotificacionService(notificacion)
// 			// 	}
// 			// }
// 		}
// 	})

// 	//  & INICIO PPROCESO ENVIO DE ARCHIVOS POR EMAIL: PAGOS, RENDICIONES , REVERSIONES Y BATCH(SOLO DPEC)
// 	// // FIXME Revisar esta funcionalidad de envios de correos
// 	// ? 1 Enviar pagos
// 	// confSendPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_PAGOS", "0 */12 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos de pagos a los clientes")
// 	// //  @every 12h
// 	// logs.Info(confSendPagos)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c.AddFunc("0 40 07 * * *", func() {
// 	// 	// 1 obtener lista de cliente
// 	// 	request := reportedtos.RequestPagosClientes{}
// 	// 	clientes, err := reportes.GetClientes(request)
// 	// 	if len(clientes.Clientes) > 0 {
// 	// 		if err == nil {
// 	// 			// obtener los pagos por cliente
// 	// 			listaPagosClientes, err := reportes.GetPagosClientes(clientes, request)
// 	// 			if err != nil {
// 	// 				notificacion := entities.Notificacione{
// 	// 					Tipo:        entities.NotificacionSendEmailCsv,
// 	// 					Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear archivos csv de clientes. %s", err),
// 	// 				}
// 	// 				service.CreateNotificacionService(notificacion)
// 	// 			}

// 	// 			// enviar los pagos a clientes
// 	// 			if len(listaPagosClientes) > 0 {
// 	// 				listaErro, err := reportes.SendPagosClientes(listaPagosClientes)
// 	// 				if err != nil {
// 	// 					notificacion := entities.Notificacione{
// 	// 						Tipo:        entities.NotificacionSendEmailCsv,
// 	// 						Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", err),
// 	// 					}
// 	// 					service.CreateNotificacionService(notificacion)
// 	// 				} else if len(listaErro) > 0 {
// 	// 					notificacion := entities.Notificacione{
// 	// 						Tipo:        entities.NotificacionSendEmailCsv,
// 	// 						Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", listaErro),
// 	// 					}
// 	// 					service.CreateNotificacionService(notificacion)
// 	// 				}
// 	// 			} else {
// 	// 				notificacion := entities.Notificacione{
// 	// 					Tipo:        entities.NotificacionSendEmailCsv,
// 	// 					Descripcion: fmt.Sprintln("No existen pagos de clientes para enviar a email"),
// 	// 				}
// 	// 				service.CreateNotificacionService(notificacion)
// 	// 			}
// 	// 		}

// 	// 	}
// 	// 	if err != nil {
// 	// 		notificacion := entities.Notificacione{
// 	// 			Tipo:        entities.NotificacionCierreLote,
// 	// 			Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
// 	// 		}
// 	// 		service.CreateNotificacionService(notificacion)
// 	// 	}

// 	// })
// 	// 1 END ENVIAR PAGOS

// 	// & 2 Enviar archivo de rendicioN: estos se envian posterior a las transferencias del dia
// 	confSendArchivoRendicion, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_RENDICION", "0 30 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos de rendicion a los clientes")
// 	//  @every 12h
// 	logs.Info(confSendArchivoRendicion)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc("0 00 09 * * *", func() {
// 		// 1 obtener lista de cliente
// 		// 1 obtener lista de cliente
// 		request := reportedtos.RequestPagosClientes{}
// 		clientes, err := reportes.GetClientes(request)
// 		if len(clientes.Clientes) > 0 {
// 			if err == nil {

// 				// obtener los pagos por cliente los transferidos
// 				listaRendicionClientes, err := reportes.GetRendicionClientes(clientes, request)
// 				if err != nil {
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionSendEmailCsv,
// 						Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear archivos de rendicion csv de clientes. %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}

// 				if len(listaRendicionClientes) > 0 {
// 					listaErro, err := reportes.SendPagosClientes(listaRendicionClientes)
// 					if err != nil {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionSendEmailCsv,
// 							Descripcion: fmt.Sprintf("No se pudo enviar archivos rendicion csv de clientes. %s", err),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					} else if len(listaErro) > 0 {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionSendEmailCsv,
// 							Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", listaErro),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 				} else {
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionSendEmailCsv,
// 						Descripcion: fmt.Sprintln("No existen rendiciones de clientes para enviar a email"),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}
// 			}

// 		}
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})

// 	//NOTE archivo bacth solo DPEC (incluye los pagos solo aprobados/autorizados del dia anterior)
// 	confSendArchivoBatch, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_BATCH", "0 */12 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos batch a los clientes")
// 	//  @every 12h
// 	logs.Info(confSendArchivoBatch)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.AddFunc("0 40 07 * * *", func() {

// 		// ctx := getContextAuditable(c)
// 		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

// 		// 1 obtener lista de cliente
// 		request := reportedtos.RequestPagosClientes{}
// 		clientes, err := reportes.GetClientes(request)
// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionSendEmailCsv,
// 				Descripcion: fmt.Sprintf("No se pudo obtener clientes para construir el archivo. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}
// 		if len(clientes.Clientes) > 0 {
// 			if err == nil {
// 				// obtener los pagos/pagoitems por cliente
// 				// NOTE solo se obtiene los que son movimientos -> pagos autorizados
// 				listaPagosItems, err := reportes.GetPagoItems(clientes, request)
// 				if err != nil {
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionArchivoBatchCliente,
// 						Descripcion: fmt.Sprintf("No se pudo obtener lista de pagos de clientes para procesar archivo batch. %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}
// 				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
// 				if len(listaPagosItems) > 0 {
// 					// se crea la estructura correspondiente para el archivo
// 					resultpagositems := reportes.BuildPagosItems(listaPagosItems)
// 					if len(resultpagositems) > 0 {
// 						err := reportes.ValidarEsctucturaPagosItems(resultpagositems) // validar estructura antes de crear el archivo
// 						if err != nil {
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionArchivoBatchCliente,
// 								Descripcion: fmt.Sprintf("la estructura del archivo batch creado es incorrecta. %s", err),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						} else {
// 							// enviar archivo por sftp
// 							err := reportes.SendPagosItems(ctx, resultpagositems, request)
// 							if err != nil {
// 								notificacion := entities.Notificacione{
// 									Tipo:        entities.NotificacionArchivoBatchCliente,
// 									Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", err),
// 								}
// 								service.CreateNotificacionService(notificacion)
// 							} else {
// 								logs.Info("Archivo batch enviado con exito")
// 								notificacion := entities.Notificacione{
// 									Tipo:        entities.NotificacionArchivoBatchCliente,
// 									Descripcion: fmt.Sprintln("el archivo batch se envio con exito"),
// 								}
// 								service.CreateNotificacionService(notificacion)
// 							}

// 						}
// 					} else {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionArchivoBatchCliente,
// 							Descripcion: fmt.Sprintf("existe un error al costruir el archivo batch del cliente. %s", err),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}
// 				} else {
// 					notificacion := entities.Notificacione{
// 						Tipo:        entities.NotificacionArchivoBatchCliente,
// 						Descripcion: fmt.Sprintf("no existen pagos batch para informar al cliente. %s", err),
// 					}
// 					service.CreateNotificacionService(notificacion)
// 				}
// 			}

// 		}

// 		if err != nil {
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionArchivoBatchCliente,
// 				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo batch a clientes. %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		}

// 	})

// 	// // Caducar pagos con metodos offline que expiran
// 	// confCaducarPagosOfflineExpirados, err := _buildPeriodicidad(util, "PERIODICIDAD_CADUCAR_PAGOSOFFLINE_EXPIRADOS", "0 */12 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad caducar pagos con pagointentos offline procesando vencidos")

// 	// logs.Info(confCaducarPagosOfflineExpirados)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c.AddFunc("0 40 12 * * *", func() {

// 	// 	_, err := service.GetCaducarOfflineIntentos()

// 	// 	if err != nil {
// 	// 		notificacion := entities.Notificacione{
// 	// 			Tipo:        entities.NotificacionPagoOfflineExpirado,
// 	// 			Descripcion: fmt.Sprintf("No se pudo realizar el proceso de caducar pagos offline vencidos. %s", err),
// 	// 		}
// 	// 		service.CreateNotificacionService(notificacion)
// 	// 	}

// 	// })

// 	// END ARCHIVOS DE RENDICION

// 	//? actualizacion a pagos expirados
// 	// confPagosExpirados, err := _buildPeriodicidad(util, "PERIODICIDAD_UPDATE_PAGOS_EXPIRADOS", "@monthly", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que expira los pagos en estado pending")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c.AddFunc(confPagosExpirados.Valor, func() {
// 	// 	err := service.ModificarEstadoPagosExpirados()
// 	// 	if err != nil {
// 	// 		notificacion := entities.Notificacione{
// 	// 			Tipo:        entities.NotificacionPagoExpirado,
// 	// 			Descripcion: fmt.Sprintf("No se pudo realizar la modificación de estado de los pagos expirados. %s", err),
// 	// 		}
// 	// 		service.CreateNotificacionService(notificacion)
// 	// 	}
// 	// })
// 	c.Start()
// }

// // NOTE esta funcion realiza consultas al servicio de apilink. Ademas notifica al cliente del cambio de estado de los pagos(debines)
// func _myFunctConciliacionPagosApilink(service administracion.Service) {
// 	listas, err := service.BuildCierreLoteApiLinkService()
// 	if err != nil {
// 		logs.Error(err)
// 		notificacion := entities.Notificacione{
// 			Tipo:        entities.NotificacionCierreLote,
// 			Descripcion: fmt.Sprintf("error al consultar servicio apilink: %s", err),
// 		}
// 		service.CreateNotificacionService(notificacion)
// 	}
// 	if len(listas.ListaPagos) > 0 && len(listas.ListaCLApiLink) > 0 {
// 		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
// 		err = service.CreateCLApilinkPagosService(ctx, listas)
// 		if err != nil {
// 			logs.Error(err)
// 			notificacion := entities.Notificacione{
// 				Tipo:        entities.NotificacionCierreLote,
// 				Descripcion: fmt.Sprintf("error al crear registros en apilinkcierrelote: %s", err),
// 			}
// 			service.CreateNotificacionService(notificacion)
// 		} else {
// 			// NOTE Notificar al usuario del cambio de estado de pagos
// 			filtro := linkdebin.RequestDebines{
// 				BancoExternalId: false,
// 				Pagoinformado:   true,
// 			}
// 			debines, _ := service.GetConsultarDebines(filtro)
// 			if len(debines) > 0 {
// 				// NOTE construir lote de pagos debin que se notificara al cliente
// 				pagos, debin, erro := service.BuildNotificacionPagosCLApilink(debines)
// 				if erro != nil {
// 					errorBuildNotificacion := errors.New("error al obtener debines para notificar al cliente")
// 					err = errorBuildNotificacion
// 					logError := entities.Log{
// 						Tipo:          entities.EnumLog("error"),
// 						Funcionalidad: "GetConsultarDebines",
// 						Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
// 					}
// 					errCrearLog := service.CreateLogService(logError)
// 					if errCrearLog != nil {
// 						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
// 					}
// 				}
// 				if len(pagos) > 0 && len(debin) > 0 {
// 					// NOTE notificar lote de pagos a clientes
// 					pagosNotificar := service.NotificarPagos(pagos)
// 					if len(pagosNotificar) > 0 {
// 						//NOTE Si se envian los pagos con exito se debe actualziar el campo pagoinformado en la tabla aplilinkcierrelote
// 						filtro := linkdebin.RequestListaUpdateDebines{
// 							DebinId: debin,
// 						}
// 						// actualizar pagoinformado en tabla apilinkcierrelote
// 						erro := service.UpdateCierreLoteApilink(filtro)
// 						if erro != nil {
// 							logs.Error(erro)
// 							notificacion := entities.Notificacione{
// 								Tipo:        entities.NotificacionCierreLote,
// 								Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink pagoinformado: %s", erro),
// 							}
// 							service.CreateNotificacionService(notificacion)
// 						}
// 					} else {
// 						notificacion := entities.Notificacione{
// 							Tipo:        entities.NotificacionWebhook,
// 							Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos debines"),
// 						}
// 						service.CreateNotificacionService(notificacion)
// 					}

// 				}
// 			}
// 		}

// 	}
// }
