package routes

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	apiresponder "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/responderdto"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtroAdm "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	filtroCl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
	"github.com/gofiber/fiber/v2"
)

func CierreLoteRoutes(app fiber.Router, serviceCL cierrelote.Service, serviceAD administracion.Service, util util.UtilService, middlewares middlewares.MiddlewareManager) {

	//////////////////////////////////////////// end-point para cieere de lotes y conciliaciones ./////////////////////////////////////////////////

	app.Get("/leer-archivo-minio", middlewares.ValidarPermiso("psp.herramienta"), getArchivosMinio(serviceCL, serviceAD, util)) //"psp.realizar.cl"

	app.Post("/conciliar-banco-cl", middlewares.ValidarPermiso("psp.herramienta"), getConciliarBancoCL(serviceCL))                               //("psp.conciliacion.bancocl"                                   // psp.conciliacion.bancocl
	app.Get("/generar-movimiento-manual", middlewares.ValidarPermiso("psp.herramienta"), getGenerarMovimientoManual(serviceCL, serviceAD))       // psp.generar.movimiento
	app.Get("/actualizar-estado-movimiento-banco", middlewares.ValidarPermiso("psp.herramienta"), getActualizarEstadoMovimientoBanco(serviceCL)) // psp.actualizar.estadomovimiento

	app.Get("/procesar-tabla-movimientos-mx", middlewares.ValidarPermiso("psp.herramienta"), getProcesarTablaMx(serviceCL)) // "psp.procesar.tablamx"
	app.Get("/procesar-tabla-pagos-px", middlewares.ValidarPermiso("psp.herramienta"), getProcesarTablaPx(serviceCL))       // "psp.procesar.tablapx"

	app.Post("/conciliacion-cl-mx", middlewares.ValidarPermiso("psp.herramienta"), getConciliacionClMx(serviceCL, serviceAD)) // "psp.conciliar.tablas"
	app.Post("/conciliacion-cl-px", middlewares.ValidarPermiso("psp.herramienta"), getConciliacionClPx(serviceCL, serviceAD)) // "psp.conciliar.tablas"
	// reportes
	app.Get("/tablas-conciliadas", middlewares.ValidarPermiso("psp.herramienta"), getPrismaTrPagos(serviceCL, serviceAD)) // "psp.conciliar.tablas"

	///////////////////////////////////////////////////////////////

	////////////////////////////////////////    end-point otros //////////////////////////////////////////////////

	app.Get("/presentaciones-prisma", middlewares.ValidarPermiso("psp.herramienta"), getMovimientosPrisma(serviceCL))
	app.Get("/one-cierre-lote-prisma", middlewares.ValidarPermiso("psp.herramienta"), getOneCierreLotePrisma(serviceCL))
	app.Get("/all-cierre-lote-prisma", middlewares.ValidarPermiso("psp.herramienta"), getAllCierreLotePrisma(serviceCL)) //
	app.Put("/edit-cierre-lote-prisma", middlewares.ValidarPermiso("psp.herramienta"), editCierreLotePrisma(serviceCL))  //middlewares.ValidarPermiso("psp.herramienta"),
	app.Delete("/delete-cierre-lote-prisma", middlewares.ValidarPermiso("psp.herramienta"), deleteCierreLotePrisma(serviceCL))

	//////////////////////////////////////// GET CIERRE LOTE API LINK ///////////////////////////////////////
	app.Get("/all-cierre-lote-apilink", middlewares.ValidarPermiso("psp.herramienta"), getAllCierreLoteApilink(serviceCL))
}

func getArchivosMinio(serviceCL cierrelote.Service, serviceAD administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request Variable

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		// paso 1:
		// leer el directorio ftp de cierre de lote minio y obtener informacion de los archivos
		// y se guarda en un directorio temporal los archivos txt existentes
		nombreDirectorio := config.DIR_KEY
		archivos, rutaArchivos, totalArchivos, err := serviceCL.LeerArchivoLoteExterno(context.Background(), nombreDirectorio)
		if err != nil {
			return fiber.NewError(400, "Error archivos: "+err.Error())
		}
		if totalArchivos <= 0 {
			return fiber.NewError(400, "no existen archivos para realizar el cierre de lote")
		}
		// se obtiene todos los estados externos de prisma
		filtro := filtroAdm.PagoEstadoExternoFiltro{
			Vendor:           strings.ToUpper("prisma"),
			CargarEstadosInt: true,
		}
		estadosPagoExterno, err := serviceAD.GetPagosEstadosExternoService(filtro)
		if err != nil {
			return fiber.NewError(400, "error al recuperar estados pago: "+err.Error())
		}

		// ob

		// paso 2:
		// se recorren uno a uno los archivos de cierre de lotes y se almacena a la bd
		listaArchivo, err := serviceCL.LeerCierreLoteTxt(archivos, rutaArchivos, estadosPagoExterno)
		if err != nil {
			return fiber.NewError(400, "error al mover los archivos: "+err.Error())
		}
		// paso 3:
		// se mueven todos los archivos de la carpeta temporal al minio.
		countArchivos, err := serviceCL.MoverArchivos(context.Background(), rutaArchivos, listaArchivo)
		if err != nil {
			return fiber.NewError(400, "error al borrar los archivos temporales: "+err.Error())
		}
		// paso 4:
		// por ultimo se borran todos los archivos creados temporalmente y el directorio temporal
		err = serviceCL.BorrarArchivos(context.Background(), nombreDirectorio, rutaArchivos, listaArchivo)
		if err != nil {
			return fiber.NewError(400, "error al borrar los archivos temporales: "+err.Error())
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":       true,
			"data":         "res",
			"contadorfile": countArchivos,
			"message":      "lectura de archivos ok",
		})
	}
}

func getConciliarBancoCL(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request Variable
		err := c.BodyParser(&request)
		logs.Info(request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		err = request.Validar()
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		// fechaPagoProcesar, err := request.ObtenerFechaPagoProcesar("2006-01-02")
		// if err != nil {
		// 	return fiber.NewError(400, err.Error())
		// }
		/*
			obtengo los pagos con sus cierre lotes relacionados
		*/
		filtro := filtroCl.FiltroTablasConciliadas{
			FechaPago: request.FechaAcreditacion, // fechaPagoProcesar,
			Match:     true,
			Reversion: request.Reversion,
		}
		responseListprismaTrPagos, err := serviceCL.ObtenerRepoPagosPrisma(filtro)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		/* obtener configuracion periodo de acreditacion */
		movimientoBanco, erro := serviceCL.ConciliacionBancoPrisma(request.FechaAcreditacion, request.Reversion, responseListprismaTrPagos) // fechaPagoProcesar
		if erro != nil {
			return fiber.NewError(400, "Error al conciliar movimiento banco y cierre lote prisma: "+erro.Error())
		}
		if movimientoBanco == nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  true,
				"data":    movimientoBanco,
				"message": "no existe movimientos en banco para conciliar con los pagos",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    movimientoBanco,
			"message": "conciliacion de banco con cierre de lote exitoso",
		})
	}
}

func getGenerarMovimientoManual(serviceCL cierrelote.Service, serviceAD administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtroCl.FiltroCierreLote
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "error en el parametro recibido "+err.Error())
		}
		responseCierreLote, err := serviceAD.BuildPrismaMovimiento(request.Reversion)
		if err != nil {
			return fiber.NewError(400, "error BuildPrismaMovimineto: "+err.Error())
		}
		ctxPrueba := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
		err = serviceAD.CreateMovimientosService(ctxPrueba, responseCierreLote)
		if err != nil {
			return fiber.NewError(400, "error CreateMovimientosService: "+err.Error())
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    responseCierreLote,
			"message": "movimineto manual creado",
		})
	}
}

func getActualizarEstadoMovimientoBanco(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		estadoResponse, err := serviceCL.ActualizarMovimientosBanco()
		if err != nil {
			return fiber.NewError(400, "error al actualizar movimientos banco: "+err.Error())
		}
		if !estadoResponse {
			return c.Status(200).JSON(&fiber.Map{
				"status":  estadoResponse,
				"message": "no existen movimientos para actualizar",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "estadoResponse",
			"message": "actualizacion estado moviminetos manual realizado con exito",
		})
	}
}

func getProcesarTablaMx(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		movimientoMx, movimientoMxEntity, err := serviceCL.ObtenerMxMoviminetosServices()
		if err != nil {
			return fiber.NewError(400, "error al obtener registros de la tablas movimientos mx: "+err.Error())
		}
		tablasRelaciondas, err := serviceCL.ObtenerTablasRelacionadasServices()
		if err != nil {
			return fiber.NewError(400, "error al obtener las tablas relacionadas: "+err.Error())
		}
		resultadoMovimientoMx := serviceCL.ProcesarMovimientoMxServices(movimientoMx, tablasRelaciondas)
		if len(resultadoMovimientoMx) <= 0 {
			return fiber.NewError(400, "error: procesar movimiento mx se encuentra vacia.")
		}
		err = serviceCL.SaveMovimientoMxServices(resultadoMovimientoMx, movimientoMxEntity)
		if err != nil {
			return fiber.NewError(400, "error al guardar los movimientos: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "ok",
			"message": "proceso de tablas movimientos mx",
		})
	}
}

func getProcesarTablaPx(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		pagosPx, entityPagoPxStr, err := serviceCL.ObtenerPxPagosServices()
		if err != nil {
			return fiber.NewError(400, "error al obtener registros de la tablas pagos px: "+err.Error())
		}
		err = serviceCL.SavePagoPxServices(pagosPx, entityPagoPxStr)
		if err != nil {
			return fiber.NewError(400, "error al guardar liquidacion de prisma: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "ok",
			"message": "proceso de tablas pagos px",
		})
	}
}

func getConciliacionClMx(serviceCL cierrelote.Service, serviceAD administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtroCl.FiltroCierreLote
		var filtro filtroCl.FiltroPrismaMovimiento
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "error en los parametros recibidos por el filtro "+err.Error())
		}
		filtro = filtroCl.FiltroPrismaMovimiento{
			Match:                        false,
			CargarDetalle:                true,
			Contracargovisa:              true,
			Contracargomaster:            true,
			Tipooperacion:                true,
			Rechazotransaccionprincipal:  true,
			Rechazotransaccionsecundario: true,
			Motivoajuste:                 true,
			ContraCargo:                  request.ContraCargo,
			CodigosOperacion:             []string{"0005"},
			TipoAplicacion:               "+",
		}
		if request.ContraCargo {
			filtro.Match = false
			filtro.CodigosOperacion = []string{"1507", "6000","1517"}
			filtro.TipoAplicacion = "-"
		}
		listaMovimientos, codigoautorizacion, err := serviceCL.ObtenerPrismaMovimientosServices(filtro)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		//var codigoautorizacion []string
		// for _, value := range listaMovimientos {
		// 	for _, value1 := range value.DetalleMovimientos {
		// 		codigoautorizacion = append(codigoautorizacion, value1.NroAutorizacionXl[3:len(value1.NroAutorizacionXl)])
		// 	}
		// }

		listaCierreLote, err := serviceCL.ObtenerCierreloteServices(request, codigoautorizacion)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := serviceCL.ConciliarCierreLotePrismaMovimientoServices(listaCierreLote, listaMovimientos)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
			return fiber.NewError(400, "no existe cierre de lotes para conciliar con movimientos")
		}
		logs.Info("en end-point")
		logs.Info(listaCierreLoteProcesado)
		err = serviceCL.ActualizarCierreloteMoviminetosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "conciliacion cierre lote con movimientos exitoso",
		})
	}
}

func getConciliacionClPx(serviceCL cierrelote.Service, serviceAD administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtroCl.FiltroCierreLote
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "error en los parametros recibidos por el filtro "+err.Error())
		}
		var codigo []string
		listaCierreLote, err := serviceCL.ObtenerCierreloteServices(request, codigo)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		if len(listaCierreLote) == 0 {
			return fiber.NewError(400, "no existe reversiones o contra cargo informados")
		}
		filtroCabecera := filtroCl.FiltroPrismaMovimiento{
			ContraCargo: request.ContraCargo,
		}
		listaCierreLoteMovimientos, err := serviceCL.ObtenerPrismaMovimientoConciliadosServices(listaCierreLote, filtroCabecera)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		var listaFechaPagos []string
		for _, value := range listaCierreLoteMovimientos {
			fechaString := value.MovimientoCabecer.FechaPago.Format("2006-01-02")
			listaFechaPagos = append(listaFechaPagos, fechaString)
		}
		filtro := filtroCl.FiltroPrismaTrPagos{
			Match:         false,
			CargarDetalle: true,
			Devolucion:    request.Devolucion,
			FechaPagos:    listaFechaPagos,
		}
		listaPrismaPago, err := serviceCL.ObtenerPrismaPagosServices(filtro)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := serviceCL.ConciliarCierreLotePrismaPagoServices(listaCierreLoteMovimientos, listaPrismaPago)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
			return fiber.NewError(400, "no existe cierre de lotes para conciliar con pagos")
		}
		err = serviceCL.ActualizarCierrelotePagosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "conciliacion cierre lote con pagos exitoso",
		})
	}

}

func getPrismaTrPagos(serviceCL cierrelote.Service, serviceAD administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtroCl.FiltroTablasConciliadas
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "error en los parametros recibidos "+err.Error())
		}
		err = request.Validar()
		if err != nil {
			return fiber.NewError(400, "error: "+err.Error())
		}
		result, err := serviceCL.ObtenerRepoPagosPrisma(request)
		if err != nil {
			return fiber.NewError(400, "error al obtener registros de la tablas pagos conciliado: "+err.Error())
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  "ok",
			"data":    result,
			"message": "consulta de tablas conciliadas cl - movimiento - pago",
		})
	}
}

//////////////////// ENDPOINTS HERRAMIENTA CIERRE LOTE PRISMA //////////////////////

func getOneCierreLotePrisma(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		id_param := c.Query("id")
		// Validacion del parametro id
		Id, _ := strconv.Atoi(id_param)
		if Id <= 0 {
			r := apiresponder.NewResponse(400, nil, "Error: no se indicó el id del cierre de lote a consultar", c)
			return r.Responder()
		}
		// filtro de la request
		var oneCierreLoteFiltro filtroCl.OneCierreLoteFiltro

		// parse de los parametros de la request al filtro CierreLoteFiltro
		err := c.QueryParser(&oneCierreLoteFiltro)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en los parametros recibidos "+err.Error(), c)
			return r.Responder()
		}

		// se consulta al servicio correspondiente
		result, err := serviceCL.GetOneCierreLotePrismaService(oneCierreLoteFiltro)

		if err != nil {
			r := apiresponder.NewResponse(404, nil, err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, result, "Datos de consulta de tablas conciliadas enviados", c)
		return r.Responder()
	}
}

func getAllCierreLotePrisma(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// filtro de la request
		var cierreLoteFiltro filtroCl.CierreLoteFiltro

		// parse de los parametros de la request al filtro CierreLoteFiltro
		err := c.QueryParser(&cierreLoteFiltro)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en los parametros recibidos "+err.Error(), c)
			return r.Responder()
		}
		// Validar los parametros
		err = cierreLoteFiltro.Validar()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}
		// Enviar la consulta al servicio correspondiente
		result, err := serviceCL.GetAllCierreLotePrismaService(cierreLoteFiltro)

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		// si no hubo resultados en la consulta, pero tampoco errores, devolver en data un string vacio
		if len(result.CierresLotes) == 0 {

			r := apiresponder.NewResponse(200, []string{}, "Datos de consulta de tablas conciliadas enviados", c)
			return r.Responder()

		}

		r := apiresponder.NewResponse(200, result, "Datos de consulta de tablas conciliadas enviados", c)
		return r.Responder()
	}
}

func getMovimientosPrisma(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request filtroCl.FiltroMovimientosPrisma
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "error en los parametros "+err.Error())
		}
		result, meta, err := serviceCL.GetAllMovimientosPrismaServices(request)
		if err != nil {
			return fiber.NewError(400, "error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "ok",
			"data":    result,
			"meta":    meta,
			"message": "consulta de presentaciones prima exitoso",
		})
	}
}

func editCierreLotePrisma(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request cierrelotedtos.RequestPrismaCL

		err := c.BodyParser(&request)

		// Body parsing validation
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}

		// movimientos, err := service.UpdateMovimientos(request)
		err = serviceCL.EditCierreLotePrismaService(request)
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		// if !movimientos {
		// 	return c.Status(200).JSON(&fiber.Map{
		// 		"status":  movimientos,
		// 		"message": "No se encontraron movimientos por actualizar",
		// 	})
		// }
		r := apiresponder.NewResponse(200, []string{}, "el cierre de lote se actualizó con exito", c)
		return r.Responder()
	}
}

func deleteCierreLotePrisma(serviceCL cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			r := apiresponder.NewResponse(404, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}
		if id < 1 {
			r := apiresponder.NewResponse(403, nil, "el id de cuenta comision es invalido", c)
			return r.Responder()
		}

		err = serviceCL.DeleteCierreLotePrismaService(uint64(id))

		if err != nil {
			r := apiresponder.NewResponse(400, nil, err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, true, "baja logica exitosa", c)
		return r.Responder()
	}
}

//////////////////// ENDPOINTS CIERRE LOTE APILINK //////////////////////

func getAllCierreLoteApilink(service cierrelote.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// filtro de la request
		var apilinkCierreloteFiltro filtroCl.ApilinkCierreloteFiltro

		// parse de los parametros de la request al filtro CierreLoteFiltro
		err := c.QueryParser(&apilinkCierreloteFiltro)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en los parametros recibidos "+err.Error(), c)
			return r.Responder()
		}
		// Validar los parametros
		err = apilinkCierreloteFiltro.Validar()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}
		// Enviar la consulta al servicio correspondiente
		result, err := service.GetAllCierreLoteApiLinkService(apilinkCierreloteFiltro)

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		// si no hubo resultados en la consulta, pero tampoco errores, devolver en data un string vacio
		if len(result.CierresLotes) == 0 {

			r := apiresponder.NewResponse(200, []string{}, "Datos de consulta enviados", c)
			return r.Responder()

		}

		r := apiresponder.NewResponse(200, result, "Datos de consulta enviados", c)
		return r.Responder()
	}
}

// ////////////////////////
type Variable struct {
	FechaAcreditacion string `josn:"fecha_acreditacion"`
	Reversion         bool   `josn:"reversion"`
}

func (v *Variable) Validar() error {
	erro := errors.New("parametros fecha no validos")
	format := "2006-01-02"
	regularFecha := regexp.MustCompile(`([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))$`)
	valorboolpresentacion := regularFecha.MatchString(v.FechaAcreditacion)
	if commons.StringIsEmpity(v.FechaAcreditacion) || !valorboolpresentacion {
		return erro

	}
	fechaprocesar, err := time.Parse(format, v.FechaAcreditacion)
	if err != nil {
		erro = errors.New(fmt.Sprintf("Error al analizar la cadena de fecha: %v", err))
		return erro
	}
	err = cierrelotedtos.EnumDiaSemana(fechaprocesar.Weekday().String()).IsValid()
	if err != nil {
		erro = errors.New(fmt.Sprintf("Error: %v", err))
		return erro
	}
	return nil
}

func (v *Variable) ObtenerFechaPagoProcesar(formatoFecha string) (fechaPagoProcesar string, erro error) {
	fechaActual, err := time.Parse(formatoFecha, v.FechaAcreditacion)
	if err != nil {
		logs.Error(err)
		erro = errors.New("error al parsear fecha")
		return
	}
	fechatemporal := fechaActual.Add(24 * -1)
	fechaPagoProcesar = fechatemporal.Format(formatoFecha)
	return
}
