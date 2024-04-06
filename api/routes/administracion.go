package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
	apiresponder "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/responderdto"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func AdministracionRoutes(app fiber.Router, middlewares middlewares.MiddlewareManager, service administracion.Service, util util.UtilService, banco banco.BancoService) {
	app.Post("/", middlewares.ValidarPermiso("psp.consultar.pagos"), getPagos(service))
	app.Post("/pagos", middlewares.ValidarPermiso("psp.consultar.pagos"), getPagosNew(service))
	app.Get("/items-pagos", middlewares.ValidarPermiso("psp.consultar.pagos"), getItemsPagos(service))
	//app.Use(middlewares.ValidarPermiso("pasarela.cliente")).
	app.Get("/pago/:pago", middlewares.ValidarPermiso("psp.consultar.pago"), getPago(service))
	app.Post("/consulta/pago", middlewares.ValidarPermiso("psp.consultar.pago"), postPagosConsulta(service))
	app.Post("/consulta/pagoId", postPagosConsulta(service))
	app.Get("/saldo-cuenta", middlewares.ValidarPermiso("psp.consultar.saldo_cuenta"), saldoCuenta(service))
	app.Get("/saldo-cliente", middlewares.ValidarPermiso("psp.consultar.saldo_cliente"), saldoCliente(service))
	app.Get("/plan-cuotas-info", middlewares.ValidarPermiso("psp.actualizar.plan_cuotas"), getAllPlanCuotas(service))
	// obtiene la tabla de intereses por cuotas, utilizado para ver costos desde el checkout
	app.Get("/plan-cuotas", getPlanCuotas(service))

	// obtiene impuesto que se aplica al recargo  cuando se paga con tarjeta de credito
	app.Get("/obtener-impuesto", getObtenerImpuesto(service))

	//Transferencias
	app.Get("/transferencias", middlewares.ValidarPermiso("psp.consultar.transferencias"), getTransferencias(service))
	app.Post("/transferencia-cliente", middlewares.ValidarPermiso("psp.cliente.transferencias"), transferenciaCliente(service, util))
	app.Post("/transferencia-por-movimiento", middlewares.ValidarPermiso("psp.herramienta"), transferenciaPorMovimiento(service, util))
	app.Post("/transferencia-comisiones-impuestos", middlewares.ValidarPermiso("psp.cliente.transferencias"), transferenciaComisionesImpuestos(service))

	//Movimientos
	app.Post("/movimiento-cuenta", middlewares.ValidarPermiso("psp.consultar.movimiento_cuenta"), movimientoCuenta(service))
	app.Get("/movimiento-transferencia", middlewares.ValidarPermiso("psp.consultar.transferencias"), movimientoTransferencia(service))

	//Abm Clientes
	app.Get("/clientes", middlewares.ValidarPermiso("psp.admin.clientes"), getClientes(service))
	app.Get("/clientes-herramienta", middlewares.ValidarPermiso("psp.herramienta"), getClientes(service))
	app.Post("/obtener-cliente", middlewares.ValidarPermiso("psp.admin.clientes"), getCliente(service))
	app.Get("/cliente-login", middlewares.ValidarPermiso("psp.consultar.clientelogin"), getClienteLogin(service))
	app.Post("/cliente", middlewares.ValidarPermiso("psp.admin.clientes"), postCliente(service, util))
	app.Put("/cliente", middlewares.ValidarPermiso("psp.admin.clientes"), putCliente(service, util))

	//Abm Rubros
	app.Get("/rubros", middlewares.ValidarPermiso("psp.consultar.rubros"), getRubros(service))
	app.Get("/rubro", middlewares.ValidarPermiso("psp.admin.rubros"), getRubro(service))
	app.Post("/rubro", middlewares.ValidarPermiso("psp.admin.rubros"), postRubro(service, util))
	app.Put("/rubro", middlewares.ValidarPermiso("psp.admin.rubros"), putRubro(service, util))

	//Abm PagosTipo
	app.Get("/pagos-tipo", middlewares.ValidarPermiso("psp.admin.pagotipo"), getPagosTipo(service))
	app.Get("/pago-tipo", middlewares.ValidarPermiso("psp.consultar.pagotipo"), getPagoTipo(service))
	app.Post("/pago-tipo", middlewares.ValidarPermiso("psp.admin.pagotipo"), postPagoTipo(service))
	app.Put("/pago-tipo", middlewares.ValidarPermiso("psp.admin.pagotipo"), putPagoTipo(service))
	app.Delete("/pago-tipo", middlewares.ValidarPermiso("psp.admin.pagotipo"), deletePagoTipo(service))

	//Abm PagosTipo-Channel
	app.Get("/pagostipo-channel", middlewares.ValidarPermiso("psp.consultar.pagostipochannel"), getPagosTipoChannel(service))
	app.Post("/pagostipo-channel", middlewares.ValidarPermiso("psp.crear.pagostipochannel"), postPagoTipoChannel(service))
	app.Delete("/pagostipo-channel", middlewares.ValidarPermiso("psp.bajar.pagostipochannel"), deletePagoTipoChannel(service))

	//Abm Channels
	app.Get("/channels", middlewares.ValidarPermiso("psp.consultar.canales"), getChannels(service))
	app.Get("/channel", middlewares.ValidarPermiso("psp.admin.canales"), getChannel(service))
	app.Post("/channel", middlewares.ValidarPermiso("psp.crear.canal"), postChannel(service, util))
	app.Put("/channel", middlewares.ValidarPermiso("psp.admin.canales"), putChannel(service, util))
	app.Delete("/channel", middlewares.ValidarPermiso("psp.bajar.canal"), deleteChannel(service, util))

	//Abm Configuraciones
	app.Get("/configuraciones", middlewares.ValidarPermiso("psp.admin.configuraciones"), getConfiguraciones(service))
	app.Get("/configuracion", middlewares.ValidarPermiso("psp.admin.configuraciones"), getConfiguracion(service, util))
	app.Post("/configuracion", middlewares.ValidarPermiso("psp.admin.configuraciones"), postConfiguracion(service, util))
	app.Put("/configuracion", middlewares.ValidarPermiso("psp.admin.configuraciones"), putConfiguracion(service, util))
	app.Put("/configuracion-send-email", middlewares.ValidarPermiso("psp.modificar.teminos_condiciones"), putConfiguracionSendEmail(service))

	//Abm Cuentas Comision
	app.Get("/cuentas-comision", middlewares.ValidarPermiso("psp.admin.cuentacomisiones"), getCuentasComision(service))
	app.Get("/cuenta-comision", middlewares.ValidarPermiso("psp.consultar.cuenta_comision"), getCuentaComision(service))
	app.Post("/cuenta-comision", middlewares.ValidarPermiso("psp.admin.cuentacomisiones"), postCuentaComision(service))
	app.Put("/cuenta-comision", middlewares.ValidarPermiso("psp.admin.cuentacomisiones"), putCuentaComision(service))
	app.Delete("/cuenta-comision", middlewares.ValidarPermiso("psp.bajar.cuenta_comision"), deleteCuentaComision(service))

	//Abm Cuentas
	app.Get("/cuentas", middlewares.ValidarPermiso("psp.consultar.cuentas"), getCuentas(service))
	app.Get("/cuenta", middlewares.ValidarPermiso("psp.admin.cuentas"), getCuenta(service))
	app.Post("/cuentas", middlewares.ValidarPermiso("psp.admin.cuentas"), postCuenta(service))
	app.Post("/cuentas/tipo", middlewares.ValidarPermiso("psp.crear.pagotipo"), postPagoTipoCuenta(service))
	app.Put("/cuenta", middlewares.ValidarPermiso("psp.admin.cuentas"), putCuenta(service))
	app.Put("/cuenta-setkey", middlewares.ValidarPermiso("psp.cuentas.apikey"), setApiKey(service))

	app.Get("/subcuenta", middlewares.ValidarPermiso("psp.admin.cuentas"), getSubcuenta(service))
	app.Post("/subcuenta", middlewares.ValidarPermiso("psp.admin.cuentas"), postSubcuenta(service))
	app.Get("/subcuentas", middlewares.ValidarPermiso("psp.admin.cuentas"), getSubcuentas(service))
	app.Post("/delete-subcuenta", middlewares.ValidarPermiso("psp.admin.cuentas"), deleteSubcuenta(service))

	/* Adm Impuestos */
	app.Post("/impuestos", middlewares.ValidarPermiso("psp.crear.impuestos"), postImpuesto(service))    //psp.crear.impuestos
	app.Get("/impuestos", middlewares.ValidarPermiso("psp.admin.clientes"), getImpuestos(service))      //psp.consultar.impuestos
	app.Put("/impuestos", middlewares.ValidarPermiso("psp.modificar.impuestos"), putImpuestos(service)) //psp.modificar.impuestos

	//Abm Channels-Arancel
	app.Get("/channels-arancel", middlewares.ValidarPermiso("psp.admin.aranceles"), getChannelsArancel(service))                  // psp.consultar.channels_arancel
	app.Get("/channel-arancel", middlewares.ValidarPermiso("psp.admin.aranceles"), getChannelArancel(service))                    // psp.consultar.channel_arancel
	app.Post("/channels-arancel", middlewares.ValidarPermiso("psp.admin.aranceles"), postChannelsArancel(service, util))          // psp.crear.channel_arancel
	app.Put("/channels-arancel", middlewares.ValidarPermiso("psp.admin.aranceles"), putChannelsArancel(service, util))            // psp.modificar.channel_arancel
	app.Delete("/channels-arancel", middlewares.ValidarPermiso("psp.baja.channel_arancel"), deleteChannelsArancel(service, util)) // psp.baja.channel_arancel

	//Ri BCRA
	app.Get("/informacion-supervision", middlewares.ValidarPermiso("psp.consultar.supervision.bcra"), GetInformacionSupervision(service))
	app.Post("/informacion-supervision", middlewares.ValidarPermiso("psp.crear.supervision.bcra"), PostInformacionSupervision(service))
	app.Get("/informacion-estadistica", middlewares.ValidarPermiso("psp.consultar.estadistica.bcra"), GetInformacionEstadistica(service))
	app.Post("/informacion-estadistica", middlewares.ValidarPermiso("psp.crear.estadistica.bcra"), PostInformacionEstadistica(service))
	//Solicitud de Cuenta
	app.Post("/solicitud-cuenta", PostSolicitudCuenta(service)) // middlewares.ValidarPermiso("psp.solicitar.cuenta"),
	//Consulta de Destinatario
	app.Post("/consulta-destinatario", middlewares.ValidarPermiso("psp.consultar.destinatario"), GetConsultaDestinatario(service))
	// obtener los pagos estados
	app.Get("/estados-de-pagos", middlewares.ValidarPermiso("psp.consultar.pagos"), GetEstadosDePagos(service))

	//Cuentas ApiLink  - Esto se ejecuta una sola vez al cambiar de ambiente(sanbox|homologacion|produccion)
	// app.Post("/cuenta-link", middlewares.ValidarPermiso("scan.show"), PostCrearCuentaLink(service))
	// app.Delete("/cuenta-link", middlewares.ValidarPermiso("scan.show"), DeleteCuentaLink(service))
	// app.Get("/cuenta-link", middlewares.ValidarPermiso("scan.show"), GetCuentaLink(service))

	/*
		Autor: Jose Alarcon
		Fecha: 1/7/2022
		Descripción:  Ejecutar procesos background conciliacion de apilink , transferencias y rapipago , notificacion (webhook) y transferencias automaticas
	*/
	/* Inicio*/
	// ^ Proceso que permitira calcular las comisiones de pagos del dia
	app.Get("/calcular-movimientostemporales-pagos", middlewares.ValidarPermiso("psp.notificacion.pagos"), GenerarMovimientosTemporalesPagos(service)) //

	// EJECUTAR PROCESO NOTIFICACION DE PAGOS A CLIENTES
	app.Get("/notificacion-pagos", middlewares.ValidarPermiso("psp.herramienta"), NotificacionPagos(service)) // psp.notificacion.pagos

	// ? CIERRELOTE APILINK
	// ? PASO 1 Buscar en Apilink lote de pagos, actualzar solo con los llegaron a un estado final se crean en la tabla apilinkcierrelote y se actualiza el estado de los pagos
	// ? PASO 2 Informar a clientes el cambio de estado de los debines
	// ? PASO 3 Conciliacion con banco
	// ? PASO 4 Generacion de movimientos
	app.Get("/apilink-cierre-lote", middlewares.ValidarPermiso("psp.apilink.cierre.lote"), CierreLoteApilink(service))
	app.Get("/notificar-pagos-clapilink-actualizados", middlewares.ValidarPermiso("psp.apilink.cierre.lote"), NotificarPagosClApilink(service))
	app.Get("/apilink-cierre-lote-conciliar-banco", middlewares.ValidarPermiso("psp.apilink.cierre.lote"), CierreLoteApilinkBanco(service, banco))
	app.Get("/apilink-cierre-lote-generarmov", middlewares.ValidarPermiso("psp.apilink.cierre.lote"), CierreLoteApilinkMov(service))

	app.Get("/conciliacion-banco", middlewares.ValidarPermiso("psp.herramienta"), ConciliacionTransferencias(service, banco)) // psp.conciliacion.banco

	// & CIERRELOTE RAPIPAGO
	//& 2 actualizar estado del pago con lo encontrado en cierrelote(archivo recibido)
	app.Get("/actualizar-pagos-cl", middlewares.ValidarPermiso("psp.herramienta"), ActualizarEstadosPagosClRapipago(service, banco))
	//& 3 Notificar al cliente el pago actualizado
	app.Get("/notificar-pagos-cl-actualizados", middlewares.ValidarPermiso("psp.herramienta"), NotificarPagosClRapipago(service))
	// & 4 conciliacion con banco
	app.Get("/rapipago-cierre-lote", middlewares.ValidarPermiso("psp.herramienta"), CierreLoteRapipago(service, banco))
	// & 5 generar movimientos manuales
	app.Get("/generar-movimiento-rapipago", middlewares.ValidarPermiso("psp.herramienta"), GenerarMovimientosRapipago(service, banco)) // psp.rapipago.cierre.lote                                            // psp.rapipago.cierre.lote

	// EJECUTAR PROCESO TRANSFERENCIAS AUTOMATICAS
	app.Get("/transferencias-automaticas", middlewares.ValidarPermiso("psp.transferencia.automatica"), transferenciasAutomaticas(service)) // psp.transferencia.automatica
	/*fin*/

	/* subir archivos de planes de cuotas actualizado */
	app.Post("/download-plan-cuotas", middlewares.ValidarPermiso("psp.admin.archivos"), PostDownloadPlanCuotas(service, util)) //psp.actualizar.plan_cuotas

	/* subir archivos prisma a minion */
	app.Post("/download-prisma-minio", middlewares.ValidarPermiso("psp.admin.archivos"), PostDownloadArchivos(service, util)) // psp.archivos.minion

	/* Listar Peticiones web services*/
	app.Get("/peticiones-webservices", middlewares.ValidarPermiso("psp.admin.webservices"), GetPeticiones(service)) // psp.peticiones.webservices

	/* mantenimiento */
	app.Get("/estado-mantenimiento", GetEstadoMantenimiento(service, util))

	/* Listar todos los archivos subidos por el usuario */
	app.Get("/archivos-subidos", middlewares.ValidarPermiso("psp.admin.archivos"), GetArchivosSubidos(service)) // psp.index.archivossubidos
	app.Post("/enviar-mail", middlewares.ValidarPermiso("psp.send.email"), PostEnviarMail(util))                //psp.send.email
	app.Get("/contracargos", middlewares.ValidarPermiso("psp.consultar.cliente"), GetContraCargoEnDisputa(service))
	/* Preferencias */
	app.Post("/preferencias", middlewares.ValidarPermiso("psp.consultar.cuentas"), PostPreferencias(service))     // psp.preferencias
	app.Get("/preferencias", middlewares.ValidarPermiso("psp.consultar.cuentas"), GetPreferencias(service))       // psp.preferencias
	app.Delete("/preferencias", middlewares.ValidarPermiso("psp.consultar.cuentas"), DeletePreferencias(service)) // psp.preferencias

	// & Permite generar mov de pagos en desarrollo (pruebas de clientes)
	app.Get("/generar-mov-dev", middlewares.ValidarPermiso("psp.herramienta"), GenerarMovDev(service)) //

	// ? Prmite consultar cierre de lote para herramienta wee
	app.Get("/consultar-cierrelote-rapipago", middlewares.ValidarPermiso("psp.herramienta"), getConsultarCierreloteRapipago(service))

	// Permite consultar los cierres de lote de Multipago para la herramienta Wee!
	app.Get("/consultar-cierrelote-multipago", middlewares.ValidarPermiso("psp.herramienta"), getConsultarCierreLoteMultipago(service))

	// ? Expira pagos con intentos procesando vencidos
	app.Get("/caducar-pagosintentos-offline", middlewares.ValidarPermiso("psp.herramienta"), getCaducarOfflineIntentos(service))
	app.Post("/caducar-pagos", postCaducarPagos(service))
	app.Get("/recaudacion-pdf", getRecaudacionPdf(service))

	// conciliacion automatica entre pagos exitosos y reporte de pagos
	app.Get("/conciliacion-pagos-reportes", middlewares.ValidarPermiso("psp.herramienta"), ConciliacionPagosReportes(service))

	// conciliacion automatica entre pagos exitosos y reporte de pagos
	app.Post("/asignar-bancoid-clrapipago", middlewares.ValidarPermiso("psp.herramienta"), asignar_bancoid_rapipago(service))

	//CRUD- contactosreportes
	app.Post("/contactos", middlewares.ValidarPermiso("psp.admin.clientes"), createContatosReportes(service))
	app.Delete("/contactos", middlewares.ValidarPermiso("psp.admin.clientes"), deleteContactosReportes(service))
	app.Get("/contactos", middlewares.ValidarPermiso("psp.admin.clientes"), getContactosReportes(service))
	app.Put("/contactos", middlewares.ValidarPermiso("psp.admin.clientes"), putContactosReportes(service))

	//CRUD- soporte
	//app.Post("/soportes",  middlewares.ValidarPermiso("psp.crear.cuenta"), createSoporte(service))
	//app.Put("/soportes",  middlewares.ValidarPermiso("psp.crear.cuenta"), putSoporte(service))

	//CRUD- USUARIOS
	app.Post("/user_bloqueado", postUsuarioBloqueado(service))
	app.Get("/user_bloqueado", getUsuariosBloqueados(service))
	app.Put("/user_bloqueado", putUsuarioBloqueado(service))
	app.Delete("/user_bloqueado", deleteUsuarioBloqueado(service))

	app.Post("/user_cuil", busquedaPersona(service))

	//Endpoint para consultar el estado de la api
	app.Get("/estado-api", middlewares.ValidarApikey(), estadoApi(service))

	// Historial de operaciones
	app.Get("/historial", middlewares.ValidarPermiso("psp.accesos.administrar"), getHistorialOperaciones(service))

	// Configuracion de envios de reportes y archivos
	app.Put("/update-envios", updateEnvios(service))
}

/* ********************************************************************************************** */
func getRecaudacionPdf(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ruta := config.DIR_BASE + config.DIR_REPORTE

		temporal := ruta
		if _, err := os.Stat(temporal); os.IsNotExist(err) {
			err = os.MkdirAll(temporal, 0755)
			if err != nil {
				return err
			}
		}

		nombre := "prueba"
		request := struct {
			dato string
		}{
			dato: "123",
		}
		err := commons.GetRecaudacionPdf(request, ruta, nombre)
		if err != nil {
			return c.JSON("error")
		}

		return c.JSON("ok")
	}
}

func CierreLoteApilink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// obtener debines de apilink
		listas, err := service.BuildCierreLoteApiLinkService()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
			return r.Responder()
		}
		if len(listas.ListaPagos) > 0 && len(listas.ListaCLApiLink) > 0 {
			ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
			err = service.CreateCLApilinkPagosService(ctx, listas)
			if err != nil {
				logs.Error(err)
				r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
				return r.Responder()

			} else {
				// caso de exito del proceso
				r := apiresponder.NewResponse(200, "tipoconciliacion: listaCierreApiLinkBanco", "el proceso fue ejecutado con exito", c)
				return r.Responder()
			}
		}

		// caso no existen debines para procesar
		r := apiresponder.NewResponse(404, "tipoconciliacion: CierreLoteApilink", "el proceso fue ejecutado con exito. no existen debines para procesar", c)
		return r.Responder()
	}
}

func NotificarPagosClApilink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		// SE CONSULTA EN LA TABLA APILINKCIERRE LOTE LOS DEBINES QUE AUN NO FUERON INFORMADOS
		filtro := linkdebin.RequestDebines{
			BancoExternalId: false,
			Pagoinformado:   true,
		}
		debines, erro := service.GetConsultarDebines(filtro)
		if erro != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
			return r.Responder()
		}
		//    buscar pagos encontrados en tabla apilinkcierrelote
		if len(debines) > 0 {
			// NOTE construir lote de pagos debin que se notificara al cliente
			pagos, debin, erro := service.BuildNotificacionPagosCLApilink(debines)
			if erro != nil {
				r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
				return r.Responder()
			}
			if len(pagos) > 0 && len(debin) > 0 {
				// NOTE enviar lote de pagos a clientes
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) > 0 {
					//NOTE Si se envian los pagos con exito se debe actualziar el campo pagoinformado en la tabla aplilinkcierrelote
					filtro := linkdebin.RequestListaUpdateDebines{
						DebinId: debin,
					}
					erro := service.UpdateCierreLoteApilink(filtro)
					if erro != nil {
						logs.Error(erro)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink pagoinformado: %s", erro),
						}
						service.CreateNotificacionService(notificacion)
						r := apiresponder.NewResponse(404, "notificacion pagos apilink: NotificarPagosClApilink", "no se pudieron actualizar pagosinformados de clapilink", c)
						return r.Responder()
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
					}
					service.CreateNotificacionService(notificacion)
					r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
					return r.Responder()
				}
			}

			data := map[string]interface{}{
				"pagos informados":           pagos,
				"debines actualizados":       debin,
				"notificacion pagos apilink": "NotificarPagosClApilink",
			}
			r := apiresponder.NewResponse(200, data, "el proceso fue ejecutado con exito", c)
			return r.Responder()
		} else {
			r := apiresponder.NewResponse(404, "notificacion pagos apilink: NotificarPagosClApilink", "no existen debines para informar", c)
			return r.Responder()
		}
	}
}

func CierreLoteApilinkBanco(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		// SE CONSULTA EN LA TABLA APILINKCIERRE LOTE LOS DEBINES QUE AUN NO FUERON INFORMADOS
		filtro := linkdebin.RequestDebines{
			BancoExternalId:  false,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(debines) > 0 {
			if err == nil {
				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 2,
					ListaApilink:     debines,
				}
				listaCierreApiLinkBanco, listaBancoId, err := banco.ConciliacionPasarelaBanco(request)
				// 1.1 conciliar lista de debines de apilink con los movimientos de banco
				// listaBancoId se utilizara para actualizar los movimientos de banco
				// listaCierreApiLinkBanco, listaBancoId, err := banco.ConciliacionBancoApliLInk(listaCierre)

				/*si no hay error guardar listaCierreloteapilink en la base de datos */
				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("error al conciliar movimiento banco y cierre loteapilink: %s", err),
					}
					service.CreateNotificacionService(notificacion)
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				} else {

					// NOTE Actualizar lista de cierreloteapilink campo banco external_id, match y fecha de acreditacion
					if len(listaCierreApiLinkBanco.ListaApilink) > 0 || len(listaCierreApiLinkBanco.ListaApilinkNoAcreditados) > 0 {
						listas := linkdebin.RequestListaUpdateDebines{
							Debines:              listaCierreApiLinkBanco.ListaApilink,
							DebinesNoAcreditados: listaCierreApiLinkBanco.ListaApilinkNoAcreditados,
						}
						erro := service.UpdateCierreLoteApilink(listas)
						if erro != nil {
							logs.Error(erro)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink: %s", erro),
							}
							service.CreateNotificacionService(notificacion)
							r := apiresponder.NewResponse(400, "conciliacion banco-apilink: CierreLoteApilinkBanco", "error al actualizar registros de cierrelote apilink", c)
							return r.Responder()
						}
					}

					if len(listaBancoId) > 0 {
						_, err := banco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros del banco: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
							return r.Responder()
							// en le caso de este error y si el pago no se actualizo a estados finales no afecta el cierre de apilink
							// el estado del pago se actualiza a estado final y no tendra en cuenta al consultar a apilink
							// ACCION : se debe actualizar manualmente el campo check en la tabla de movimientos de banco(no es obligatorio)
						}
					}
				}

			}

		} else {
			r := apiresponder.NewResponse(404, "tipoconciliacion: listaCierreApiLinkBanco", "no existen pagos con debin para conciliar", c)
			return r.Responder()
		}

		// caso de exito
		r := apiresponder.NewResponse(200, "tipoconciliacion: listaCierreApiLinkBanco", "conciliacion banco-apilink se ejecuto con exito", c)
		return r.Responder()
	}
}

// generar mov pagos debin conciliados con banco
func CierreLoteApilinkMov(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//obtener lista pagos apilink(solo los conciliados con banco)
		filtro := linkdebin.RequestDebines{
			BancoExternalId:  true,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(debines) > 0 {
			responseCierreLote, err := service.BuildMovimientoApiLink(debines)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				} else {
					// caso de exito
					r := apiresponder.NewResponse(200, "tipoconciliacion: listaCierreApiLinkBanco", "proceso ejecutado con exito", c)
					return r.Responder()
				}
			}
		}
		r := apiresponder.NewResponse(404, "tipoconciliacion: listaCierreApiLinkBanco", "no existen debines para generar movimientos", c)
		return r.Responder()
	}
}

func ConciliacionTransferencias(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// filtro de transferencias, para traer solo aquellas que tienen el campo match en 0
		filtro := filtros.TransferenciaFiltro{}

		err := c.QueryParser(&filtro)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}
		// interesan solo las transferencias aun no matcheadas o conciliadas
		filtro.Match = 0

		// traer todas las transferencias. Response de tipo administraciondtos.TransferenciaRespons
		listaTransferencia, err := service.GetTransferencias(filtro)

		// casteo manual de []TransferenciaResponseAgrupada a []TransferenciaResponse
		// definimos una variable de tipo []TransferenciaResponse
		var listaTransferenciaResponse []administraciondtos.TransferenciaResponse
		for _, transferenciaAgrupada := range listaTransferencia.Transferencias {
			listaTransferenciaResponse = append(listaTransferenciaResponse, administraciondtos.TransferenciaResponse{
				ReferenciaBancaria:              transferenciaAgrupada.ReferenciaBancaria,
				CbuDestino:                      transferenciaAgrupada.CbuDestino,
				CbuOrigen:                       transferenciaAgrupada.CbuOrigen,
				Fecha:                           transferenciaAgrupada.Fecha,
				ReferenciaBanco:                 transferenciaAgrupada.ReferenciaBanco,
				ListaIdsTransferenciasAgrupadas: transferenciaAgrupada.IdsTransferenciasAgrupadas,
			})
		}

		if len(listaTransferencia.Transferencias) > 0 {
			request := bancodtos.RequestConciliacion{
				TipoConciliacion: 3,
				Transferencias: administraciondtos.TransferenciaResponsePaginado{
					// espera el tipo []TransferenciaResponse
					Transferencias: listaTransferenciaResponse,
				},
			}

			// listaTransferenciasMatch contiene los datos y los ids para actualizar la tabla pasarela.transferencias
			// listaIdBanco contiene los ids para actualizar la tabla banco.movimientos
			listaTransferenciasMatch, listaIdBanco, err := banco.ConciliacionPasarelaBanco(request)
			if err != nil {
				logs.Error(err)
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: fmt.Sprintf("error al conciliar transferencias con movimientos del banco: %s", err),
				}

				service.CreateNotificacionService(notificacion)
			}

			/* SI LA CANTIDAD DE TRANSFERENCIAS MATCH ES MENOR A 0/HAY ERRORES EL PROCESO DE CONCILIACION TERMINA */
			logs.Info("EL total de transferencias que se actualizaran son " + strconv.Itoa(len(listaTransferenciasMatch.TransferenciasConciliadas)))
			if len(listaTransferenciasMatch.TransferenciasConciliadas) > 0 {
				/* ACTUALIZAR CAMPO MATCH Y BANCO_EXTERNAL_ID TABLA TRANSFERENCIAS */
				err := service.UpdateTransferencias(listaTransferenciasMatch)

				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("error al actualizar transferencias: %s", err),
					}
					service.CreateNotificacionService(notificacion)
				} else {
					// Si no hubo error en la actualizacion de la tabla pasarela.transferencias
					/*ACTUALIZAR CAMPO ESTADO_CHECK EN SERVICIO BANCO*/
					response, err := banco.ActualizarRegistrosMatchBancoService(listaIdBanco, true)
					logs.Info(response)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar campo check en servicio banco: %s", err),
						}
						service.CreateNotificacionService(notificacion)
					}

				}

			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de conciliacion de transferencia. %s", err),
			}
			service.CreateNotificacionService(notificacion)

			return c.JSON(&fiber.Map{
				"error":            err,
				"tipoconciliacion": "conciliaciontransferecias",
			})
		}

		return c.JSON(&fiber.Map{
			"statusMessage":    "la conciliacion de transferecias se realizó con exito.",
			"tipoconciliacion": "conciliaciontransferecias",
		})

	}
}

func NotificacionPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request webhook.RequestWebhook

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, err := service.BuildNotificacionPagosService(request)

		if err == nil {
			pagosNotificar, err := service.CreateNotificacionPagosService(pagos)
			if err == nil {
				if len(pagosNotificar) > 0 {
					// si inicia el proceso de notifiacar al cliente
					pagosupdate := service.NotificarPagos(pagosNotificar)
					// NOTE se debe actualizar solo si el pago llegi a un estado final
					if len(pagosupdate) > 0 && request.EstadoFinalPagos { /* actualzar estado de pagos a notificado */
						err = service.UpdatePagosNoticados(pagosupdate)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionWebhook,
								Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							return fiber.NewError(400, "Error: "+err.Error())
						}
					}
					if len(pagosupdate) > 0 {
						return c.JSON(&fiber.Map{
							"data":    pagosupdate,
							"message": "se notifico con exito los siguientes pagos",
						})
					}
					if len(pagosupdate) == 0 {
						return c.JSON(&fiber.Map{
							"message": "error al notificar pagos a clientes",
						})
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintf("webhook: No existen pagos por notificar. %s", err),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            "no existen pagos por notificar",
						"tipoconciliacion": "notificacionPagos",
					})
				}
			}
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "notificacionPagos",
		})

	}
}

func transferenciasAutomaticas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

		response, erro := service.RetiroAutomaticoClientes(ctx)

		if erro != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica para los clientes. %s", erro),
			}
			err := service.CreateNotificacionService(notificacion)
			if err != nil {
				return c.JSON(&fiber.Map{
					"error":            err,
					"tipoconciliacion": "transferenciasAutomaticas",
				})
			}
		}

		// if len(response.MovimientosId) > 0 {
		// 	// enviar comisiones
		// 	uuid := uuid.NewV4()
		// 	idmovimientos := administraciondtos.RequestMovimientosId{
		// 		MovimientosId:             response.MovimientosId,
		// 		MovimimientosIdRevertidos: response.MovimimientosIdRevertidos,
		// 	}
		// 	result := service.SendTransferenciasComisiones(ctx, uuid.String(), idmovimientos)
		// 	if !result {
		// 		logs.Error(result)
		// 		notificacion := entities.Notificacione{
		// 			Tipo:        entities.NotificacionTransferencia,
		// 			Descripcion: fmt.Sprintf("error al intentar transferir comisiones impuestos telco: %v", result),
		// 		}
		// 		service.CreateNotificacionService(notificacion)
		// 	}
		// }

		return c.JSON(&fiber.Map{
			"response":         response,
			"tipoconciliacion": "transferenciasAutomaticas",
		})

	}
}

func CierreLoteRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      true,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		if len(listaCierreRapipago) > 0 {
			if err == nil {

				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 1,
					ListaRapipago:    listaCierreRapipago,
				}
				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
				listaCierreRapipago, listaBancoId, err := banco.ConciliacionPasarelaBanco(request)
				if err != nil {
					return fiber.NewError(400, "error al conciliar clrapipago con banco: "+err.Error())
				}

				if len(listaBancoId) == 0 {
					logs.Info("no existen movimientos en banco para conciliar con pagos rapipago")
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintln("no existen movimientos en banco para conciliar con pagos rapipago"),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            notificacion.Descripcion,
						"tipoconciliacion": "CierreLoteRapipago",
					})
				} else {
					/*en el caso de error a actualizar la tabla rapipagocierrelote el proceso termina */
					err := service.UpdateCierreLoteRapipago(listaCierreRapipago.ListaRapipago)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote rapipago (volver a ejecutar proceso): %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else {
						// son los registros que coincidieron en el cierrerapipago y banco
						// si no se actualiza los registros del banco se debera actualizar manualmente
						_, err := banco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar movimientos del banco - conciliacion rapipago(actualizar manualmente los siguientes movimientos): %s", err),
							}
							service.CreateNotificacionService(notificacion)
						}
					}

				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de rapipago por conciliar"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteRapipago se ejecuto con exito",
		})

	}
}

func ActualizarEstadosPagosClRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      false,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote  */
		listaPagoaRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		if len(listaPagoaRapipago) > 0 {
			if err == nil {
				listaPagosClRapipago, err := service.BuildPagosClRapipago(listaPagoaRapipago)
				if err == nil {
					// Actualizar estados del pago y cierrelote
					err = service.ActualizarPagosClRapipagoService(listaPagosClRapipago)
					if err != nil {
						logs.Error(err)
						return fiber.NewError(400, "Error: "+err.Error())
					}

				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de rapipago para actualizar"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteRapipago se ejecuto con exito",
		})

	}
}

func NotificarPagosClRapipago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.PagoEstadoFiltro

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, err := service.BuildNotificacionPagosCLRapipago(request)

		if len(pagos) > 0 {
			pagosNotificar := service.NotificarPagos(pagos)
			if len(pagosNotificar) > 0 { /* actualzar estado de pagos a notificado */
				mensaje := fmt.Sprintf("los siguientes pagos se actualizaron correctamente. %v", pagosNotificar)
				return c.JSON(&fiber.Map{
					"pagos actualizados": mensaje,
					"notificacion":       "exitosa",
					"tipoconciliacion":   "notificacionPagosClRapipago",
				})
			} else {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionWebhook,
					Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
				}
				service.CreateNotificacionService(notificacion)
				return fiber.NewError(400, "Error: "+err.Error())
			}

		} else {
			return c.JSON(&fiber.Map{
				"error":            "no existen pagos por notificar",
				"tipoconciliacion": "notificacionPagosClRapipago",
			})
		}

	}
}

func GenerarMovimientosRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: true,
			PagosNotificado:      true,
		}

		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		// Si no se guarda ningún cierre no hace falta seguir el proce
		if len(listaCierreRapipago) > 0 {
			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
			responseCierreLote, err := service.BuildRapipagoMovimiento(listaCierreRapipago)

			if err == nil {

				// 3 - Guarda los movimientos en la base de datos en caso de error se
				// repetira en el día siguiente
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				} else {
					return c.JSON(&fiber.Map{
						"mensaje":          "conciliacion exitosa",
						"tipoconciliacion": "CierreLoteRapipago",
					})
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return nil

	}
}

func GetConsultaDestinatario(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request linkconsultadestinatario.RequestConsultaDestinatarioLink

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		uuid := uuid.NewV4()

		response, err := service.GetConsultaDestinatarioService(uuid.String(), request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)

	}
}

func PostSolicitudCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.SolicitudCuentaRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		err = service.SendSolicitudCuenta(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.Status(200).JSON(&fiber.Map{

			"status":  true,
			"message": "la solicitud se recibió correctamente",
		})
	}
}

func putCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.CuentaRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateCuentaService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func setApiKey(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.CuentaRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.SetApiKeyService(ctx, &request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"api_key":    request.Apikey,
		})

	}
}

func getImpuestos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtros.ImpuestoFiltro
		// var requestPaginacion filtros.Paginacion
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: los datos de paginacion no son validos ")
		}
		// err = c.BodyParser(&request)
		// if err != nil {
		// 	return fiber.NewError(400, "Error: los datos enviado no son validos")
		// }
		// request.Paginacion = requestPaginacion
		response, err := service.GetImpuestosService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(&response)
	}
}

func postImpuesto(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.ImpuestoRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.PostImpuestoService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"id":      id,
			"status":  true,
			"message": "El impuesto se guardó correctamente.",
		})
	}
}

func putImpuestos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.ImpuestoRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateImpuestoService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func getTransferencias(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request filtros.TransferenciaFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		request.Paginacion.Number = request.Number
		request.Paginacion.Size = request.Size
		response, err := service.GetTransferencias(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

// NOTE modificar para obener solo los movimientos para transferir //
func movimientoTransferencia(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request filtros.MovimientoFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		request.CargarPago = true
		request.CargarPagoEstados = true
		request.CargarPagoIntentos = true
		request.CargarMedioPago = true
		request.AcumularPorPagoIntentos = true

		// este campo viene desde el front end
		request.CargarMovimientosNegativos = true

		movimiento, err := service.GetMovimientosAcumulados(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&movimiento)
	}
}

// Busca el valor de los pagos
func getPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoFiltro

		err := c.BodyParser(&request)

		if err != nil || request.CuentaId < 1 {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		// transformar la fecha de fin de filtro para que haga correctamente la busqueda
		if !request.FechaPagoFin.IsZero() {
			fechaFinLastMoment := commons.GetDateLastMomentTime(request.FechaPagoFin)
			request.FechaPagoFin = fechaFinLastMoment
		}

		request.CargaPagoIntentos = true
		request.CargarPagoTipos = true
		request.CargarPagoEstado = true
		request.CargaMedioPagos = true
		request.CargarChannel = true
		request.CargarCuenta = true

		data, err := service.GetPagosService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		// Evitar la paginacion para acumular saldos en el mismo rango de fecha
		request.Paginacion.Number = 0

		data2, err := service.GetPagosService(request)

		data.SaldoPendiente = data2.SaldoPendiente
		data.SaldoDisponible = data2.SaldoDisponible

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(data)
	}
}

// Busca el valor de los pagos
func getPagosNew(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoFiltro

		err := c.BodyParser(&request)

		if err != nil || request.CuentaId < 1 {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		// transformar la fecha de fin de filtro para que haga correctamente la busqueda
		if !request.FechaPagoFin.IsZero() {
			fechaFinLastMoment := commons.GetDateLastMomentTime(request.FechaPagoFin)
			request.FechaPagoFin = fechaFinLastMoment
		}

		data, err := service.GetPagos(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		var filtroSaldo filtros.PagoFiltro
		filtroSaldo.FechaPagoInicio = request.FechaPagoInicio
		filtroSaldo.FechaPagoFin = request.FechaPagoFin
		filtroSaldo.CuentaId = request.CuentaId
		filtroSaldo.PagoEstadosIds = []uint64{4, 7}

		resp, err := service.GetSaldoPagoCuenta(filtroSaldo)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		data.SaldoDisponible = resp.SaldoDisponible
		data.SaldoPendiente = resp.SaldoPendiente

		return c.JSON(data)
	}
}

// Busca los items de un pago específico
func getItemsPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoItemFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		data, err := service.GetItemsPagos(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"message": "items obtenidos correctamente",
		})
	}
}

func postPagosConsulta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Error: Debe enviar una api key válida")
		}

		var request administraciondtos.RequestPagosConsulta

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: parámetros incorrectos")
		}

		data, err := service.GetPagosConsulta(api, request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		message := "datos enviados"
		if len(*data) <= 0 {
			message = "no se encontró  coincidencia para la consulta realizada"
		}
		total_registros := len(*data)
		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"total":   total_registros,
			"message": message,
		})
	}
}

func getPago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagoID, _ := strconv.Atoi(c.Params("pago"))
		if pagoID <= 0 {
			return fiber.NewError(401, "Error: no se indicó el id del cliente a consultar")
		}

		resultado, err := service.GetPagoByID(int64(pagoID))
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    resultado,
			"message": "Datos del pago enviados.",
		})
	}
}

func getCuentas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		number, _ := strconv.Atoi(c.Query("number"))
		size, _ := strconv.Atoi(c.Query("size"))
		cliente, _ := strconv.ParseInt(c.Query("cliente"), 10, 64)

		if cliente <= 0 {
			return fiber.NewError(401, "Error: no se indicó el id del cliente a consultar")
		}

		meta, links, data, err := service.GetCuentasByCliente(cliente, number, size)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"meta":  meta,
			"links": links,
			"data":  data,
		})
	}
}

func getCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.CuentaFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetCuenta(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}

}

/*
	 func deleteSubcuenta(service administracion.Service) fiber.Handler {
		return func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Query("id"))
			if err != nil {
				return fiber.NewError(404, "Error: "+err.Error())
			}

			ctx := getContextAuditable(c)

			err = service.DeleteSubcuentaService(ctx, uint64(id))

			if err != nil {
				return fiber.NewError(400, "Error: "+err.Error())
			}

			return c.JSON(&fiber.Map{
				"resultado": true,
			})
		}
	}
*/
func deleteSubcuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request []administraciondtos.SubcuentaRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		ok, err := service.DeleteSubcuenta(ctx, &request)
		if !ok {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

func getSubcuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.CuentaFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetSubcuenta(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}

}

func getSubcuentas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		number, _ := strconv.Atoi(c.Query("number"))
		size, _ := strconv.Atoi(c.Query("size"))
		cuenta, _ := strconv.ParseInt(c.Query("cuenta"), 10, 64)

		if cuenta <= 0 {
			return fiber.NewError(401, "Error: no se indicó el id de la cuenta a consultar")
		}

		meta, links, data, err := service.GetSubcuentasByCuenta(cuenta, number, size)
		if err != nil {
			if strings.Contains(err.Error(), "no se encontraron subcuentas para") {
				return c.Status(200).JSON(&fiber.Map{
					"meta":  meta,
					"links": links,
					"data":  data,
				})
			} else {
				return fiber.NewError(400, "Error: "+err.Error())
			}

		}

		return c.JSON(&fiber.Map{
			"meta":  meta,
			"links": links,
			"data":  data,
		})
	}
}

func postSubcuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request []administraciondtos.SubcuentaRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		ok, err := service.PostSubcuenta(ctx, request)
		if !ok {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"status":  true,
			"message": "Guardado exitoso.",
		})

	}

}

func postCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.CuentaRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		if ok, err := service.PostCuenta(ctx, request); !ok {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"message": "la Cuenta se guardó correctamente.",
		})
	}
}

func getContextAuditable(c *fiber.Ctx) context.Context {
	userid := string(c.Response().Header.Peek("user_id"))
	intUserID, _ := strconv.Atoi(userid)
	userctx := entities.Auditoria{
		UserID: uint(intUserID),
		IP:     c.IP(),
	}
	ctx := context.WithValue(c.Context(), entities.AuditUserKey{}, userctx)
	return ctx
}

func postPagoTipoCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request entities.Pagotipo

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		if ok, err := service.PostPagotipo(ctx, &request); !ok {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"message": "el tipo de pago se guardó correctamente.",
		})
	}
}

func saldoCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cuenta_id, err := strconv.Atoi(c.Query("cuenta_id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		saldo, err := service.GetSaldoCuentaService(uint64(cuenta_id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"cuentas_id": saldo.CuentasId,
			"total":      saldo.Total.Float64(),
		})
	}
}

func saldoCliente(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cliente_id, err := strconv.Atoi(c.Query("cliente_id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		saldo, err := service.GetSaldoClienteService(uint64(cliente_id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"cliente_id": saldo.ClienteId,
			"total":      saldo.Total.Float64(),
		})
	}
}

// Retorna el movimiento de una cuenta este metodo se usa en la pagina de administación para mostrar los movimientos al usuario
func movimientoCuenta(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		// var request filtroCl.FiltroCierreLote
		// err := c.BodyParser(&request)
		// if err != nil {
		// 	return fiber.NewError(400, "error en los parametros recibidos por el filtro "+err.Error())
		// }
		var request filtros.MovimientoFiltro
		err := c.BodyParser(&request)

		// var request filtros.MovimientoFiltro
		// err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		request.CargarPagoIntentos = true
		request.CargarPagoEstados = true
		request.CargarPago = true
		request.CargarMedioPago = true
		request.CargarComision = true
		request.CargarImpuesto = true
		request.CargarTransferencias = true
		request.CargarRetenciones = true

		movimiento, err := service.GetMovimientos(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&movimiento)
	}
}

func transferenciaCliente(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cuenta_id, err := strconv.Atoi(c.Query("cuenta_id"))

		if err != nil {
			return fiber.NewError(404, "Error en el la consulta: "+err.Error())
		}

		var request administraciondtos.RequestTransferenicaCliente

		err = c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		uuid := uuid.NewV4()
		response, err := service.BuildTransferenciaCliente(ctx, uuid.String(), request, uint64(cuenta_id))
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		} else {
			// si no existen errores se deben transferir comisiones impuestos
			request := administraciondtos.RequestComisiones{
				MovimientosId: response.MovimientosIdTransferidos,
			}
			_, erro := service.SendTransferenciasComisiones(ctx, uuid.String(), request)
			mensaje := fmt.Sprintf("Error: transferencia decomisiones %+v", erro)
			logs.Error(erro.Error() + mensaje)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       mensaje,
				Funcionalidad: "SendTransferenciasComisiones",
			}
			erro = util.CreateLogService(log)
			if erro != nil {
				logs.Error(erro.Error() + mensaje)
			}
		}

		return c.JSON(&response)
	}
}

func transferenciaPorMovimiento(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request administraciondtos.RequestTransferenciaMov

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		uuid := uuid.NewV4()
		response, err := service.BuildTransferenciaClientePorMovimiento(ctx, uuid.String(), request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func transferenciaComisionesImpuestos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request administraciondtos.RequestComisiones

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		uuid := uuid.NewV4()
		result, err := service.SendTransferenciasComisiones(ctx, uuid.String(), request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"resultado": result,
		})
	}
}

func getAllPlanCuotas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		InstallmentId, err := strconv.ParseInt(c.Query("MedioInstallmentId"), 10, 64)
		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}
		planCuotas, err := service.GetAllInstallmentsById(InstallmentId)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"estatus": "ok",
			"data":    planCuotas,
		})
	}
}

func getPlanCuotas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fecha := time.Now().Format("2006-01-02")
		planCuotas, err := service.GetInteresesPlanes(fecha)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(planCuotas)
	}
}

func getObtenerImpuesto(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filtro := filtros.ConfiguracionFiltro{
			Nombre: "IMPUESTO_SOBRE_COEFICIENTE",
		}
		configuracionImpuesto, err := service.GetConfiguracionesService(filtro)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener configuración "+err.Error())
		}
		impuestoId, err := strconv.Atoi(configuracionImpuesto.Data[0].Valor)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener configuración ")
		}
		filtroImpuesto := filtros.ImpuestoFiltro{
			Id: uint(impuestoId),
		}
		impuesto, err := service.GetImpuestosService(filtroImpuesto)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(impuesto.Impuestos[0])
	}
}

func GetInformacionSupervision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request ribcradtos.GetInformacionSupervisionRequest

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ri, err := service.GetInformacionSupervision(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(ri)

	}
}

func PostInformacionSupervision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		//Cargo el archivo pdf
		file, err := c.FormFile("infespecial")

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}
		//Cargo el json de datos y luego lo convierto en un objeto
		value := c.FormValue("data")

		datos := []byte(value)

		var ri ribcradtos.BuildInformacionSupervisionRequest
		err = json.Unmarshal(datos, &ri)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ri.InfEspecial = file
		ri.Fiber = c

		response, err := service.BuildInformacionSupervision(ri)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(response)

	}
}

func GetInformacionEstadistica(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request ribcradtos.GetInformacionEstadisticaRequest

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ri, err := service.GetInformacionEstadistica(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(ri)

	}
}

func PostInformacionEstadistica(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request ribcradtos.BuildInformacionEstadisticaRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.BuildInformacionEstadistica(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(response)

	}
}

func getCliente(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ClienteFiltro

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		cliente, err := service.GetClienteService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&cliente)
	}
}

func getClienteLogin(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		logs.Info("entra en endpoint - getClienteLogin")
		var request filtros.ClienteFiltro

		err := c.QueryParser(&request)

		if request.UserId < 1 {
			return fiber.NewError(400, "Error: "+"debes informar el id del usuario")
		}

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		cliente, err := service.GetClienteService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&cliente)
	}
}

func getClientes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ClienteFiltro

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		clientes, err := service.GetClientesService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&clientes)
	}
}

func postCliente(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.ClienteRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreateClienteService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se dió de alta un nuevo cliente con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putCliente(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.ClienteRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateClienteService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Actualizar
		observacion := fmt.Sprintf("Se actualizaron los datos del cliente con ID %v", request.Id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func getRubro(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.RubroFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetRubroService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getRubros(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.RubroFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetRubrosService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postRubro(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RubroRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreateRubroService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se dió de alta un nuevo rubro con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putRubro(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RubroRequest

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateRubroService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Actualizar
		observacion := fmt.Sprintf("Se actualizaron los datos del rubro con ID %v", request.Id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func getPagoTipo(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoTipoFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		request.CargarCuenta = true

		response, err := service.GetPagoTipoService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getPagosTipo(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoTipoFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		request.CargarCuenta = true

		response, err := service.GetPagosTipoService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postPagoTipo(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestPagoTipo

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreatePagoTipoService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putPagoTipo(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestPagoTipo

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdatePagoTipoService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func deletePagoTipo(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.DeletePagoTipoService(ctx, uint64(id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

// Abm Channel

func getChannel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ChannelFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetChannelService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getChannels(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ChannelFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetChannelsService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postChannel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestChannel

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreateChannelService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se creó un nuevo channel con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putChannel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestChannel

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateChannelService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Actualizar
		observacion := fmt.Sprintf("Se actualizaron los datos del channel con ID %v", request.Id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func deleteChannel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.DeleteChannelService(ctx, uint64(id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Eliminar
		observacion := fmt.Sprintf("Se actualizaron los datos del rubro con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

// Abm Cuentas Comision

func getCuentaComision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.CuentaComisionFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		request.CargarCuenta = true
		request.CargarChannel = true

		response, err := service.GetCuentaComisionService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getCuentasComision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.CuentaComisionFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}
		request.CargarCuenta = true
		request.CargarChannel = true
		request.Channelarancel = true

		response, err := service.GetCuentasComisionService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postCuentaComision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestCuentaComision

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreateCuentaComisionService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putCuentaComision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestCuentaComision

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateCuentaComisionService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func deleteCuentaComision(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.DeleteCuentaComisionService(ctx, uint64(id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

// Abm Configuraciones

func getConfiguraciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ConfiguracionFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetConfiguracionesService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getConfiguracion(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ConfiguracionFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := util.GetConfiguracionService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postConfiguracion(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestConfiguracion

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		id, err := util.CreateConfiguracionService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se dió de alta una nueva configuración con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putConfiguracion(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestConfiguracion

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateConfiguracionService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Actualizar
		observacion := fmt.Sprintf("Se actualizaron los datos de la configuración con ID %v", request.Id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func putConfiguracionSendEmail(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestConfiguracion

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateConfiguracionSendEmailService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

// /// es para probar el cierre de lote end point temporal
// func verPagosUuid(service administracion.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Accepts("application/x-www-form-urlencoded")
// 		responseMovimientoCierreLote, err := service.BuildPrismaMovimiento()
// 		if err != nil {
// 			return fiber.NewError(400, "Error al intentar chequear servicio")
// 		}
// 		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
// 		err1 := service.CreateMovimientosService(ctx, responseMovimientoCierreLote)
// 		if err1 != nil {
// 			return fiber.NewError(401, "Error:"+err.Error())
// 		}

//			prisma := responseMovimientoCierreLote.ListaCLPrisma
//			movimiento := responseMovimientoCierreLote.ListaMovimientos
//			pagointento := responseMovimientoCierreLote.ListaPagoIntentos
//			pagos := responseMovimientoCierreLote.ListaPagos
//			pagoslogs := responseMovimientoCierreLote.ListaPagosEstadoLogs
//			reversiones := responseMovimientoCierreLote.ListaReversiones
//			return c.JSON(&fiber.Map{
//				"cierreLote":   prisma,
//				"status":       movimiento,
//				"pago-intento": pagointento,
//				"pagos":        pagos,
//				"pago-logs":    pagoslogs,
//				"reversion":    reversiones,
//			})
//		}
//	}
func PostCrearCuentaLink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request linkcuentas.LinkPostCuenta

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		err = service.CreateCuentaApilinkService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"creada": true,
		})

	}
}

func DeleteCuentaLink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request linkcuentas.LinkDeleteCuenta

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		err = service.DeleteCuentaApilinkService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"deshabilitada": true,
		})

	}
}

func GetEstadosDePagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		response, err := service.GetPagosEstadosService(false, false)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"estados": response,
		})
	}
}

func GetCuentaLink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		resp, err := service.GetCuentasApiLinkService()

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(resp)

	}
}

func PostDownloadPlanCuotas(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request administraciondtos.RequestPlanCuotas

		form, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtenr e larchivo "+err.Error())
		}
		ruta := fmt.Sprintf("..%s/plancuota", config.DOC_CL)
		if _, err := os.Stat(ruta); os.IsNotExist(err) {
			err = os.MkdirAll(ruta, 0755)
			if err != nil {
				return fiber.NewError(400, "Error: no se pudo crear el directorio "+err.Error())
			}
		}
		request.RutaFile = fmt.Sprintf("%s/%s", ruta, form.Filename)
		err = c.SaveFile(form, fmt.Sprintf(request.RutaFile))
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo guardar el archivo "+err.Error())
		}
		request.InstalmentsId = c.FormValue("InstallmentsId", "")
		request.VigenciaDesde = c.FormValue("VigenciaDesde", "")
		erro := service.CreatePlanCuotasService(request)
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo procesar el archivo"+erro.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se subieron archivos de plan de cuotas")
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.Status(200).JSON(&fiber.Map{
			"status": true,
		})

	}
}

func GetPeticiones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtros.PeticionWebServiceFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetPeticionesService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.Status(200).JSON(&fiber.Map{
			"peticiones": response,
		})
	}
}

func getPagosTipoChannel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.PagoTipoChannelFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetPagosTipoChannelService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func deletePagoTipoChannel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.DeletePagoTipoChannelService(ctx, uint64(id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

func postPagoTipoChannel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request administraciondtos.RequestPagoTipoChannel

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)
		id, err := service.CreatePagoTipoChannel(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"id": id,
		})

	}
}

func PostDownloadArchivos(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := getContextAuditable(c)

		fecha := c.FormValue("fecha")
		multipartFormData, err := c.MultipartForm()

		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener los archivos "+err.Error())
		}
		contarArchivos := 0

		for _, value := range multipartFormData.File["file"] {

			// form, err := c.FormFile("file")
			form := value

			ruta := fmt.Sprintf("..%s/%s", config.DOC_CL, config.DIR_KEY)
			if _, err := os.Stat(ruta); os.IsNotExist(err) {
				err = os.MkdirAll(ruta, 0755)
				if err != nil {
					return fiber.NewError(400, "Error: no se pudo crear el directorio "+err.Error())
				}
			}

			RutaFile := fmt.Sprintf("%s%s-%s", ruta, fecha, form.Filename)
			err = c.SaveFile(form, fmt.Sprint(RutaFile))
			if err != nil {
				return fiber.NewError(400, "Error: no se pudo guardar el archivo "+err.Error())
			}

			archivoFile := administraciondtos.ArchivoResponse{
				NombreArchivo:  fmt.Sprintf("%s-%s", fecha, form.Filename),
				ErrorProducido: ``,
			}
			var archivos []administraciondtos.ArchivoResponse
			archivos = append(archivos, archivoFile)

			// se mueven todos los archivos de la carpeta temporal al minio y luego se borran los archivos temporales
			// countArchivos, err := service.SubirArchivos(context.Background(), rutaArchivos, listaArchivo)

			countArchivos, erro := service.SubirArchivos(ctx, ruta, archivos)

			if erro != nil {
				return fiber.NewError(400, "Error: no se pudo procesar el archivo"+erro.Error())
			}
			contarArchivos = contarArchivos + countArchivos

		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se subieron archivos de cierres de lote")
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.Status(200).JSON(&fiber.Map{
			"status":   true,
			"archivos": contarArchivos,
		})

	}
}

func getChannelsArancel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ChannelArancelFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetChannelsArancelService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postChannelsArancel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestChannelsAranncel

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		id, err := service.CreateChannelsArancelService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Crear
		observacion := fmt.Sprintf("Se dió de alta un nuevo arancel con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"id": id,
		})
	}
}

func putChannelsArancel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestChannelsAranncel

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.UpdateChannelsArancelService(ctx, request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Actualizar
		observacion := fmt.Sprintf("Se actualizaron los datos del arancel con ID %v", request.Id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"modificado": true,
			"id":         request.Id,
		})
	}
}

func deleteChannelsArancel(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		ctx := getContextAuditable(c)

		err = service.DeleteCuentaComisionService(ctx, uint64(id))

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		bearer := c.Get("Authorization")

		accion := entities.Eliminar
		observacion := fmt.Sprintf("Se dió de baja un arancel con ID %v", id)
		util.CreateHistorialOperacionesService(bearer, observacion, accion)

		return c.JSON(&fiber.Map{
			"resultado": true,
		})
	}
}

func getChannelArancel(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ChannelAranceFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		request.CargarRubro = true
		request.CargarChannel = true

		response, err := service.GetChannelArancelService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func GetEstadoMantenimiento(service administracion.Service, util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		estado, fecha, err := util.GetMatenimietoSistemaService()
		if err != nil {
			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"message": "el sistema estara en mantenimiento hasta " + time.Now().Format(time.RFC822Z),
			})

		}

		if estado {
			fechaString := fmt.Sprintf("%v", fecha)

			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"fecha":   fecha,
				"message": "el sistema estara en mantenimiento hasta " + fechaString,
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  estado,
			"message": "el sistema esta funcionando correctamente",
		})
	}
}

func GetArchivosSubidos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request filtros.Paginacion
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: los datos de paginacion no son validos ")
		}
		listArchivos, err := service.ObtenerArchivosSubidos(request)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		if len(listArchivos.ArchivosSubidos) == 0 {
			return c.Status(204).JSON(&fiber.Map{
				"status":  true,
				"data":    listArchivos,
				"message": "la busqueda no arroja resultado",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    listArchivos,
			"message": "la busqueda fue exitosa",
		})
	}
}

func PostEnviarMail(util util.UtilService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request utildtos.RequestDatosMail
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: los parametros no son validos ")
		}
		err = util.EnviarMailService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "mensaje enviado con exito",
		})
	}
}

func GetContraCargoEnDisputa(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request filtros.ContraCargoEnDisputa
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: los paramentros no son validos.")
		}
		request.CargarCuentas = true
		request.CargarTiposPago = true
		request.CargarPagos = true
		request.CargarPagosIntentos = true

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		cierreLoteDisputa, err := service.GetCierreLoteEnDisputaServices(1, request)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener operaciones en disputa.")
		}
		if len(cierreLoteDisputa) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status":  true,
				"message": "la solicitud no se proceso por completo, no existen reversiones.",
			})
		}

		listaPagoenDisputa, err := service.GetPagosByTransactionIdsServices(request, cierreLoteDisputa)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		var message string
		message = fmt.Sprintf("%v", "datos obtenidos con exito")
		if len(listaPagoenDisputa.Cuenta.PagoTipo) == 0 {
			message = fmt.Sprintf("%v", "la solicitud se realizo con exito, pero la lista se encuentra vacia.")
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": message,
			"data":    listaPagoenDisputa,
		})
	}
}

// preferencias
func PostPreferencias(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request_preferences administraciondtos.RequestPreferences
		err := c.BodyParser(&request_preferences)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener los parámetros enviados")
		}
		// obtener el archivo del contexto
		archivo, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener el archivo "+err.Error())
		}
		// ruta donde se guarda el logo o archivo
		request_preferences.File = archivo
		request_preferences.RutaLogo = config.URL_LOGOS

		if err != nil {
			return fiber.NewError(400, "Error: no se pudo guardar el archivo "+err.Error())
		}
		err = service.PostPreferencesService(request_preferences)
		if err != nil {
			return fiber.NewError(400, "Error:"+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "OK",
			"message": "Las preferencias se guardaron correctamente.",
		})

	}
}
func GetPreferencias(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request administraciondtos.RequestPreferences
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener los parámetros enviados")
		}
		responsePreference, err := service.GetPreferencesService(request)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		if responsePreference.Client == "0" {
			return c.Status(200).JSON(&fiber.Map{
				"status":  "false",
				"message": "El usuario indicado no posee preferencias",
				"data":    nil,
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "OK",
			"message": "Las preferencias se obtuvieron correctamente",
			"data":    responsePreference,
		})
	}
}
func DeletePreferencias(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request administraciondtos.RequestPreferences
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener los parámetros enviados")
		}
		err = service.DeletePreferencesService(request)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  "OK",
			"message": "Las preferencias fueron eliminadas correctamente",
		})
	}
}

// * Permite generar mov de pagos en desarrollo (pruebas de clientes)
func GenerarMovDev(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request filtros.PagoFiltro

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		// 1 Se debe consultar los pagos en estado aprobado y procesando
		request = filtros.PagoFiltro{
			PagoEstadosIds:     []uint64{2, 4},
			CargaPagoIntentos:  true,
			CargarPagoTipos:    true,
			CargarPagoEstado:   true,
			CargaMedioPagos:    false,
			CargarChannel:      true,
			CargarCuenta:       true,
			FechaPagoInicio:    request.FechaPagoInicio,
			FechaPagoFin:       request.FechaPagoFin,
			CuentaId:           request.CuentaId,
			ExternalReferences: request.ExternalReferences,
		}
		pagos, err := service.GetPagosDevService(request)
		if len(pagos) > 0 {
			pg, err := service.UpdatePagosDevService(pagos)
			if err != nil {
				logs.Info("no se actualizo correctamente los pagos , no se puede continuar con proceso")
				return fiber.NewError(400, "Error: "+err.Error())
			} else {
				// se deben aplicar las comisiones a los pagos: aplicar al ultimo pago intento
				responseCierreLote, err := service.BuildPagosMovDev(pg)

				if err == nil {

					// 3 - Guarda los movimientos en la base de datos en caso de error se
					// repetira en el día siguiente
					ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
					err = service.CreateMovimientosService(ctx, responseCierreLote)
					if err != nil {
						logs.Error(err)
						return fiber.NewError(400, "Error: "+err.Error())
					}

				}

			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de cierre de lote de apilink. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		} else if len(pagos) < 1 {
			return c.JSON(&fiber.Map{
				"error":            "no existen debines por procesar",
				"tipoconciliacion": "listaCierreApiLinkBanco",
			})
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "listaCierreApiLinkBanco",
		})

	}
}

// ? permite consultar cierrelote de rapipago para herramienta wee
func getConsultarCierreloteRapipago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.RequestClrapipago

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetConsultarClRapipagoService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getCaducarOfflineIntentos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		response, err := service.GetCaducarOfflineIntentos()

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func postCaducarPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")

		var request filtros.PagoCaducadoFiltro

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		if request.PagoTipo == "" {
			return fiber.NewError(400, "Error en los parámetros enviados: "+"Pago tipo faltante")
		}

		request.CuentaApikey = api

		err = request.Validar()
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetCaducarPagosExpirados(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func GenerarMovimientosTemporalesPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request filtros.PagoIntentoFiltros
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		// 1 Se debe consultar los pagos en estado aprobado y procesando
		request = filtros.PagoIntentoFiltros{
			PagoEstadosIds:      []uint64{4, 7},
			CargarPago:          true,
			CargarPagoTipo:      true,
			CargarPagoEstado:    true,
			CargarCuenta:        true,
			PagoIntentoAprobado: true,
			FechaPagoInicio:     request.FechaPagoInicio,
			FechaPagoFin:        request.FechaPagoFin,
		}
		pagos, err := service.GetPagosCalculoMovTemporalesService(request)
		if len(pagos) > 0 {
			// se deben aplicar las comisiones a los pagos: aplicar al ultimo pago intento
			responseCierreLote, err := service.BuildPagosCalculoTemporales(pagos)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				// crear los movimientos temoorales y actualziar campo calculado en pago intento
				// esto inidica que el pago ya fue calculado y guardado en movimientostemporales
				err = service.CreateMovimientosTemporalesService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de cierre de lote de apilink. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		} else if len(pagos) < 1 {
			return c.JSON(&fiber.Map{
				"error":            "no existen pagos por procesar",
				"tipoconciliacion": "GenerarMovimientosTemporalesPagos",
			})
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "GenerarMovimientosTemporalesPagos",
		})

	}
}

func ConciliacionPagosReportes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtro := filtros.PagoFiltro{}

		erro := c.QueryParser(&filtro)

		// error del parse
		if erro != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+erro.Error())
		}

		// error en validacion
		if filtro.CuentaId == 0 || len(filtro.Fecha[0]) == 0 {
			return fiber.NewError(400, "Error en los parámetros enviados, debe enviar los parámetros requeridos.")
		}

		// traer todas las transferencias. Response de tipo administraciondtos.TransferenciaRespons
		res, erro := service.ConciliacionPagosReportesService(filtro)

		if erro != nil {

			logs.Info("ocurrio un error en el servicio ConciliacionPagosReportesService.")

			return c.Status(404).JSON(&fiber.Map{
				"data": struct {
					data []string
				}{
					data: res},
				"statusMessage":    erro.Error(),
				"tipoconciliacion": "ConciliacionPagosReportes",
			})
		}

		return c.JSON(&fiber.Map{
			"data":             res,
			"statusMessage":    "la conciliacion de pagos se realizó con exito.",
			"tipoconciliacion": "ConciliacionPagosReportes",
		})

	}
}

func asignar_bancoid_rapipago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request administraciondtos.RequestCLRapipagoExternalId

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(401, "Error: con los parametros enviados")
		}
		status := false
		bandoMovID := request.BancoId
		clRapipagoID := request.RapipagoId

		erro := service.AsignarBancoIdRapipagoService(bandoMovID, clRapipagoID)

		if erro != nil {
			return c.JSON(&fiber.Map{
				"status":        status,
				"statusMessage": erro.Error(),
			})
		}

		status = true
		return c.JSON(&fiber.Map{
			"status":        status,
			"statusMessage": "asignacion realizada correctamente",
		})

	}
}
func createContatosReportes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestContactosReportes
		err := c.BodyParser(&request)
		//Valido que existan los parametros
		if err != nil {
			return fiber.NewError(400, "Error: con los parametros enviados por body")
		}
		err = service.CreateContactosReportesService(request)

		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Contacto creado correctamente",
		})

	}
}
func getContactosReportes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		clienteid, err := strconv.Atoi(c.Query("clienteid"))
		status := false
		//Valido que existan los parametros
		if err != nil {
			return fiber.NewError(400, "Error: con los parametros enviados por Query")
		}
		var request administraciondtos.RequestContactosReportes
		request.ClienteID = uint(clienteid)
		contactos, error := service.ReadContactosReportesService(request)
		if error != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", error.Error()))
		}
		if len(contactos.EmailsContacto) == 0 {
			return c.JSON(&fiber.Map{
				"status":  status,
				"message": "No se encontraron datos para el id enviado",
				"data":    nil,
			})
		}
		status = true
		return c.JSON(&fiber.Map{
			"status":  status,
			"message": "Datos enviados correctamente",
			"data":    contactos.EmailsContacto,
		})
	}
}
func putContactosReportes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestContactosReportes
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los datos enviados por body")
		}

		err = service.UpdateContactosReportesService(request)
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Contacto actualizado correctamente",
		})

	}
}
func deleteContactosReportes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestContactosReportes
		err := c.BodyParser(&request)
		//Valido que existan los parametros
		if err != nil {
			return fiber.NewError(400, "Error: con los parametros enviados por Body")
		}
		err = service.DeleteContactosReportesService(request)
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Contacto eliminado correctamente",
		})
	}
}

// func createSoporte(service  administracion.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Accepts("application/json")
// 		var request administraciondtos.RequestSoporte
// 		err := c.BodyParser(&request)
// 		//Valido que existan los parametros
// 		if err != nil {
// 			return fiber.NewError(400, "Error: con los parametros enviados")
// 		}
// 		// obtener el archivo del contexto
// 		file, _ := c.FormFile("file")
// 		request.File=file
// 		err = service.CreateSoporteService(request)
// 		if err != nil {
// 			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
// 		}
// 		return c.JSON(&fiber.Map{
// 			"status":        true,
// 			"statusMessage": "Soporte creado correctamente",

// 		})

// 	}
// }
// func putSoporte(service  administracion.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Accepts("application/json")
// 		var request administraciondtos.RequestSoporte
// 		err := c.BodyParser(&request)
// 		//Valido que existan los parametros
// 		if err != nil {
// 			return fiber.NewError(400, "Error: con los parametros enviados")
// 		}
// 		err = service.PutSoporteService(request)
// 		if err != nil {
// 			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
// 		}
// 		return c.JSON(&fiber.Map{
// 			"status":        true,
// 			"statusMessage": "Soporte actualizado correctamente",
// 		})

// 	}
// }

func estadoApi(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		err := service.EstadoApiService()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "El sistema no se encuentra disponible",
			})
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "El sistema se encuentra disponible",
		})

	}
}

func busquedaPersona(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestFraudeControl
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		object, err := service.CallFraudePersonas(request.Cuil)
		if err != nil {
			return c.JSON(err)
		}

		return c.JSON(object)
	}
}

func postUsuarioBloqueado(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestUserBloqueado

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		err = service.CreateUsuarioBloqueadoService(request)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "No se pudo registrar el usuario como bloqueado.",
			})
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Usuario bloqueado registrado correctamente.",
		})
	}
}

func getUsuariosBloqueados(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var filtro filtros.UsuarioBloqueadoFiltro

		err := c.QueryParser(&filtro)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		usuarios, err := service.GetUsuariosBloqueadoService(filtro)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "No se pudo obtener usuarios bloqueados",
			})
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"data":          usuarios,
			"statusMessage": "Operacion exitosa de consulta de usuarios bloqueados.",
		})
	}
}

func putUsuarioBloqueado(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestUserBloqueado

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		err = service.UpdateUsuarioBloqueadoService(request)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "No se pudo actualizar el usuario bloqueado.",
			})
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Usuario bloqueado actualizado correctamente.",
		})
	}
}

func deleteUsuarioBloqueado(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request administraciondtos.RequestUserBloqueado

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		err = service.DeleteUsuarioBloqueadoService(request)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "No se pudo actualizar el usuario bloqueado.",
			})
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Usuario bloqueado actualizado correctamente.",
		})
	}
}

func getHistorialOperaciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var filtro filtros.RequestHistorial

		err := c.QueryParser(&filtro)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetHistorialOperacionesService(filtro)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func updateEnvios(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestEnvios

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}

		err = service.UpsertEnvioService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "success. id: "+fmt.Sprintf("%d", request.Id), c)
		return r.Responder()
	}
}

func getConsultarCierreLoteMultipago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request filtros.RequestClMultipago

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		response, err := service.GetConsultarClMultipagoService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}
