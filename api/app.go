package main

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/middlewares/middlewareinterno"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api/routes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/storage"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/pagooffline"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/usuario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/webhook"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func InicializarApp(clienteHttp *http.Client, clienteSql *database.MySQLClient, clienteFile *os.File) *fiber.App {
	//Servicios comunes
	fileRepository := commons.NewFileRepository(clienteFile)
	commonsService := commons.NewCommons(fileRepository)
	algoritmoVerificacionService := commons.NewAlgoritmoVerificacion()
	middlewares := middlewares.MiddlewareManager{HTTPClient: clienteHttp}
	utilRepository := util.NewUtilRepository(clienteSql)
	utilService := util.NewUtilService(utilRepository, clienteHttp)

	// runEndpoint := util.NewRunEndpoint(clienteHttp, utilService)

	//Valida si existe un correo para solicitud de nuevas cuentas si no existe lo crea.
	utilService.FirstOrCreateConfiguracionService("EMAIL_SOLICITUD_CUENTA", "Email que recibirá la solicitud de apertura de cuenta", "developmenttelco@gmail.com")

	//ApiLink
	apiLinkRemoteRepository := apilink.NewRemote(clienteHttp, utilService)
	apiLinkRepository := apilink.NewRepository(clienteSql, utilService)
	apiLinkService := apilink.NewService(apiLinkRemoteRepository, apiLinkRepository)

	// webhooks
	webhooksRepository := webhook.NewRemote(clienteHttp)

	//Store Service
	storeService := storage.NewS3Session()
	storeServiceEst := cierrelote.NewStore(storeService)

	//storeUtil := util.NewStore(storeService)

	auditoriaRespository := auditoria.NewAuditoriaRepository(clienteSql)
	auditoriaService := auditoria.AuditoriaService(auditoriaRespository)

	administracionRepository := administracion.NewRepository(clienteSql, auditoriaService, utilService)
	administracionService := administracion.NewService(administracionRepository, apiLinkService, commonsService, utilService, webhooksRepository, storeServiceEst)
	middlewaresPasarela := middlewareinterno.MiddlewareManagerPasarela{Service: administracionService} //.MiddlewareManagerPasarela{Service: administracionService}

	/* MOVIMIENTOS BANCO: servicio para consultar y validar movimientos de pagos acreditados en la cuenta de telco*/
	movimientosBancoRemoteRepository := banco.NewRemote(clienteHttp)
	movimientosBancoService := banco.NewService(movimientosBancoRemoteRepository, utilService, administracionService)

	/* REPORTES CLIENTES */
	reportesRepository := reportes.NewRepository(clienteSql, auditoriaService, utilService)
	reportesService := reportes.NewService(reportesRepository, administracionService, utilService, commonsService, storeServiceEst)

	usuarioRepository := usuario.NewRepository(clienteSql, utilService)
	usuarioRemoteRepository := usuario.NewRemote(clienteHttp)
	usuarioService := usuario.NewService(usuarioRemoteRepository, usuarioRepository)
	bancoRepository := banco.NewRemote(clienteHttp)
	bancoService := banco.NewService(bancoRepository, utilService, administracionService)
	prismaRepository := prisma.NewRepository(clienteSql)
	remoteRepository := prisma.NewRepoasitory(clienteHttp, prismaRepository)
	prismaService := prisma.NewService(remoteRepository, prismaRepository, commonsService)
	pagoOffLineService := pagooffline.NewService(algoritmoVerificacionService)
	cierreloteRepository := cierrelote.NewRepository(clienteSql, utilRepository)
	storage := storage.NewS3Session()
	reafileStore := cierrelote.NewStore(storage)
	cierreloteService := cierrelote.NewService(cierreloteRepository, commonsService, utilService, reafileStore, administracionService, movimientosBancoService)
	checkoutRepository := checkout.NewRepository(clienteSql, auditoriaService)
	checkoutService := checkout.NewService(checkoutRepository, commonsService, prismaService, pagoOffLineService, utilService, webhooksRepository, storeServiceEst)

	// //descomentar esto en servidor
	engine := html.New(filepath.Join(filepath.Base(config.DIR_BASE), "api", "views"), ".html")
	//descomentar esto en local
	// engine := html.New("views", ".html")
	engine.Delims("${", "}")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var msg string
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}

			if msg == "" {
				msg = "No se pudo procesar el llamado a la api: " + err.Error()
			}

			_ = ctx.Status(code).JSON(internalError{
				Message: msg,
			})

			return nil
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.ALLOW_ORIGIN + ", " + config.AUTH,
		AllowHeaders: "",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Corrientes Telecomunicaciones Api Servicio de Pasarela de Pagos"))
	})
	checkout := app.Group("/checkout")
	routes.CheckoutRoutes(checkout, checkoutService)
	//pagooffline := app.Group("/pagooffline")
	//routes.PrismaRoutes(pagooffline, pagoOffLineService)
	prisma := app.Group("/prisma")
	routes.PrismaRoutes(prisma, prismaService, middlewares)

	banco := app.Group("/banco")
	routes.BancoRoutes(banco, bancoService, middlewares)

	// cierre de lote
	cierrelote := app.Group("/cierrelote")
	routes.CierreLoteRoutes(cierrelote, cierreloteService, administracionService, utilService, middlewares)

	administracion := app.Group("/administracion")
	routes.AdministracionRoutes(administracion, middlewares, administracionService, utilService, movimientosBancoService)

	usuario := app.Group("/usuario")
	routes.UsuarioRoutes(usuario, middlewares, usuarioService)

	/*reportes */
	reportes := app.Group("/reporte")
	routes.ReporteRoutes(reportes, middlewares, middlewaresPasarela, reportesService, utilService)

	/* nuevo grupo de endpoint	*/
	apiv1 := app.Group("/api/v1")
	//routes.PruebaRoutes(api, middlewares, utilService)
	routes.ApiV1Routes(apiv1, middlewares, middlewaresPasarela, administracionService, utilService)

	//Procesos en segundo plano
	//background.BackgroudServices(administracionService, cierreloteService, utilService, movimientosBancoService, reportesService)
	//descomentar esto en local
	// app.Static("/", "./views")
	//descomentar esto en servidor
	app.Static("/", filepath.Join(filepath.Base(config.DIR_BASE), "api", "views"))

	return app
}

func main() {
	var HTTPTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false, // <- this is my adjustment
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	var HTTPClient = &http.Client{
		Transport: HTTPTransport,
	}

	//HTTPClient.Timeout = time.Second * 120 //Todo validar si este tiempo está bien
	clienteSQL := database.NewMySQLClient()
	osFile := os.File{}

	app := InicializarApp(HTTPClient, clienteSQL, &osFile)

	_ = app.Listen(":3300")
}

type internalError struct {
	Message string `json:"message"`
}
