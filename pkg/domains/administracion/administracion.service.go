package administracion

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"net/smtp"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	webhooks "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service interface {
	//PAGOS
	GetPagoByID(pagoID int64) (*entities.Pago, error)
	PostPagotipo(ctx context.Context, pagotipo *entities.Pagotipo) (bool, error)
	GetPagosService(filtro filtros.PagoFiltro) (response administraciondtos.ResponsePagos, erro error)
	GetPagos(filtro filtros.PagoFiltro) (administraciondtos.ResponsePagos, error)
	GetItemsPagos(filtro filtros.PagoItemFiltro) ([]administraciondtos.PagoItems, error)
	GetSaldoPagoCuenta(filtro filtros.PagoFiltro) (administraciondtos.ResponsePagos, error)
	GetPagosConsulta(string, administraciondtos.RequestPagosConsulta) (*[]administraciondtos.ResponsePagosConsulta, error)
	ConsultarEstadoPagosService(requestValid administraciondtos.ParamsValidados, apiKey string, request administraciondtos.RequestPagosConsulta) (responsePagoEstado []administraciondtos.ResponseSolicitudPago, registrosAfectados bool, erro error)

	///
	//GetEstadoPagoRepository()()
	///

	// PAGOS INTENTOS
	GetPagosIntentosByTransaccionIdService(filtroPagoIntento filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error)

	//TRANSFERENCIAS
	GetTransferencias(filtro filtros.TransferenciaFiltro) (response administraciondtos.TransferenciaRespons, erro error)
	BuildTransferenciaCliente(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenicaCliente, cuentaId uint64) (response linktransferencia.ResponseTransferenciaCreateLink, erro error)
	BuildTransferenciaClientePorMovimiento(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenciaMov) (response []linktransferencia.ResponseTransferenciaCreateLink, erro error)
	/*conciliacion banco actualizar campos match con servicio banco */
	UpdateTransferencias(listas bancodtos.ResponseConciliacion) error
	SendTransferenciasComisiones(ctx context.Context, requerimientoId string, movimientosId administraciondtos.RequestComisiones) (res administraciondtos.ResponseTransferenciaComisiones, err error)

	//CLIENTES
	CreateClienteService(ctx context.Context, request administraciondtos.ClienteRequest) (id uint64, erro error)
	UpdateClienteService(ctx context.Context, cliente administraciondtos.ClienteRequest) (erro error)
	DeleteClienteService(ctx context.Context, id uint64) (erro error)
	GetClienteService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacion, erro error)
	GetClientesService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacionPaginado, erro error)

	//RUBROS
	CreateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (id uint64, erro error)
	UpdateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (erro error)
	GetRubroService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubro, erro error)
	GetRubrosService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubros, erro error)

	//ABM PAGOS TIPOS
	CreatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (id uint64, erro error)
	UpdatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (erro error)
	GetPagoTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagoTipo, erro error)
	GetPagosTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagosTipo, erro error)
	DeletePagoTipoService(ctx context.Context, id uint64) (erro error)

	//ABM CHANNELS
	CreateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (id uint64, erro error)
	UpdateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (erro error)
	GetChannelService(filtro filtros.ChannelFiltro) (channel administraciondtos.ResponseChannel, erro error)
	GetChannelsService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannels, erro error)
	DeleteChannelService(ctx context.Context, id uint64) (erro error)

	//ABM CUENTA COMISSIONES
	CreateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (id uint64, erro error)
	UpdateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (erro error)
	GetCuentaComisionService(filtro filtros.CuentaComisionFiltro) (channel administraciondtos.ResponseCuentaComision, erro error)
	GetCuentasComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentasComision, erro error)
	DeleteCuentaComisionService(ctx context.Context, id uint64) (erro error)

	//CUENTAS
	PostCuentaComision(ctx context.Context, comision *entities.Cuentacomision) error
	PostCuenta(ctx context.Context, cuenta administraciondtos.CuentaRequest) (bool, error)
	GetCuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseCuenta, erro error)
	GetCuentasByCliente(cliente int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Cuenta, error)
	UpdateCuentaService(ctx context.Context, request administraciondtos.CuentaRequest) (erro error)
	SetApiKeyService(ctx context.Context, request *administraciondtos.CuentaRequest) (erro error)
	DeleteCuentaService(ctx context.Context, id uint64) (erro error)
	GetCuentaByApiKeyService(apikey string) (reult bool, erro error)

	GetSubcuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseSubcuenta, erro error)
	PostSubcuenta(ctx context.Context, request []administraciondtos.SubcuentaRequest) (ok bool, err error)
	GetSubcuentasByCuenta(cuenta int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Subcuenta, error)
	PostEditSubcuenta(ctx context.Context, request administraciondtos.SubcuentaRequest) (ok bool, err error)
	/* 	DeleteSubcuentaService(ctx context.Context, id uint64) (erro error)
	 */DeleteSubcuenta(ctx context.Context, request *[]administraciondtos.SubcuentaRequest) (ok bool, erro error)

	/* impuestos */
	PostImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (id uint64, erro error)
	GetImpuestosService(filtro filtros.ImpuestoFiltro) (response administraciondtos.ResponseImpuestos, erro error)
	UpdateImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (erro error)

	// CHANNELS ARANCELES
	GetChannelsArancelService(filtro filtros.ChannelArancelFiltro) (response administraciondtos.ResponseChannelsArancel, erro error)
	CreateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (id uint64, erro error)
	UpdateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (erro error)
	DeleteChannelsArancelService(ctx context.Context, id uint64) (erro error)
	GetChannelArancelService(filtro filtros.ChannelAranceFiltro) (response administraciondtos.ResponseChannelsAranceles, erro error)

	/*
		Devuelve el saldo de una cuenta específica
	*/
	GetSaldoCuentaService(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error)

	/*
		Devuelve el saldo de un cliente específico
	*/
	GetSaldoClienteService(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error)

	//MOVIMIENTOS
	GetMovimientosAcumulados(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoAcumuladoResponsePaginado, erro error)
	GetMovimientos(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoPorCuentaResponsePaginado, erro error)
	BuildMovimientoApiLink(listaCierre []*entities.Apilinkcierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	CreateMovimientosService(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error)
	BuildPrismaMovimiento(reversion bool) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	CreateMovimientosTemporalesService(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error)

	/*
		Crea una notificación para los usuarios del sistema
	*/
	CreateNotificacionService(notificacion entities.Notificacione) error

	CreateLogService(log entities.Log) error

	/*
		Busca una lista de pagos estados. Si busca por final es true verifica si el estado el final.
	*/
	GetPagosEstadosService(buscarPorFinal, final bool) (estados []entities.Pagoestado, erro error)
	// servicio para consultar estadopago
	GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estadoPago []entities.Pagoestado, erro error)
	/*
		se obtiene una lista de pagos estado externos
	*/
	GetPagosEstadosExternoService(filtro filtros.PagoEstadoExternoFiltro) (estadosExternos []entities.Pagoestadoexterno, erro error)
	/*
		Construye y guarda una lista de cierre de lotes para api link
		Este proceso crea el cierre de lote a partir de las informaciones consultadas en apilink
	*/
	BuildCierreLoteApiLinkService() (response administraciondtos.RegistroClPagosApilink, erro error)

	/*
		permite obtener los planes de cuotas vigentes para un medio de pago
	*/
	GetPlanCuotas(idMedioPago uint) (response []administraciondtos.PlanCuotasResponseDetalle, erro error)

	/*
		obtiene los intereses de todos los planes existentes para informarlos
	*/
	GetInteresesPlanes(fecha string) (planes []administraciondtos.PlanCuotasResponse, erro error)
	/*
		obtiene todos los planes decuotas por id de installment
	*/
	GetAllInstallmentsById(installments_id int64) (planesCuotas []administraciondtos.InstallmentsResponse, erro error)
	//------------------------------------------------RI BCRA-----------------------------------------
	/*
		Crea archivo txt para enviar al BCRA
		https://telcodev.atlassian.net/secure/RapidBoard.jspa?rapidView=17&projectKey=PP&modal=detail&selectedIssue=PP-122
	*/
	RIInfestadistica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)
	/*
		Construye la información para supervisión prevista en la sección 69.1 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	GetInformacionSupervision(request ribcradtos.GetInformacionSupervisionRequest) (ri ribcradtos.RiInformacionSupervisionReponse, erro error)
	/*
		Guarda en un archivo zip la información de supervisión prevista en la sección 69.1 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	BuildInformacionSupervision(request ribcradtos.BuildInformacionSupervisionRequest) (ruta string, erro error)
	/*
		Construye la información para estadística prevista en la sección 69.2 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	GetInformacionEstadistica(request ribcradtos.GetInformacionEstadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)
	/*
		Guarda en un archivo zip la información de estadistica prevista en la sección 69.2 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	BuildInformacionEstadistica(request ribcradtos.BuildInformacionEstadisticaRequest) (ruta string, erro error)

	RIGuardarArchivos(request ribcradtos.RIGuardarArchivosRequest) (erro error)
	//------------------------------------------------RI BCRA-----------------------------------------

	/*
		Modifica el estado de pagos expirados
	*/
	ModificarEstadoPagosExpirados() (erro error)

	/*
		Realiza automaticamente las transferencias de acuerdo con el período informado en configuraciones
	*/
	RetiroAutomaticoClientes(ctx context.Context) (response administraciondtos.RequestMovimientosId, erro error)

	//CONFIGURACIONES
	GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (response administraciondtos.ResponseConfiguraciones, erro error)
	UpdateConfiguracionService(ctx context.Context, config administraciondtos.RequestConfiguracion) (erro error)
	UpdateConfiguracionSendEmailService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error)

	//Send Mails
	SendSolicitudCuenta(request administraciondtos.SolicitudCuentaRequest) (erro error)

	GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error)

	//Cuenta Apilink -> para los debines
	CreateCuentaApilinkService(request linkcuentas.LinkPostCuenta) (erro error)
	DeleteCuentaApilinkService(request linkcuentas.LinkDeleteCuenta) (erro error)
	GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error)

	// plan de cuotas INSTALLMENTS
	CreatePlanCuotasService(request administraciondtos.RequestPlanCuotas) (erro error)
	// notificaciones de pagos
	BuildNotificacionPagosService(request webhooks.RequestWebhook) (listaPagos []entities.Pagotipo, erro error)
	BuildNotificacionPagosCLRapipago(filtro filtros.PagoEstadoFiltro) (response []webhooks.WebhookResponse, erro error)
	CreateNotificacionPagosService(listaPagos []entities.Pagotipo) (response []webhooks.WebhookResponse, erro error)
	NotificarPagos(listaPagosNotificar []webhooks.WebhookResponse) (pagoupdate []uint)
	UpdatePagosNoticados(listaPagosNotificar []uint) (erro error)

	// & apilink cierrelote
	CreateCLApilinkPagosService(ctx context.Context, mcl administraciondtos.RegistroClPagosApilink) (erro error)
	CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error)
	GetDebines(request linkdebin.RequestDebines) (response []*entities.Apilinkcierrelote, erro error)
	GetConsultarDebines(request linkdebin.RequestDebines) (response []linkdebin.ResponseDebinesEliminados, erro error)
	BuildNotificacionPagosCLApilink(request []linkdebin.ResponseDebinesEliminados) (response []webhooks.WebhookResponse, debinID []uint64, erro error)
	UpdateCierreLoteApilink(request linkdebin.RequestListaUpdateDebines) (erro error)

	// & rapipagocierrelote
	GetCierreLoteRapipagoService(filtro rapipago.RequestConsultarMovimientosRapipago) (listaCierreRapipago []*entities.Rapipagocierrelote, erro error)
	UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error)
	BuildPagosClRapipago(listaPagosClRapipago []*entities.Rapipagocierrelote) (pagosclrapiapgo administraciondtos.PagosClRapipagoResponse, erro error)
	BuildRapipagoMovimiento(listaCierre []*entities.Rapipagocierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	ActualizarPagosClRapipagoService(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error)

	//PagoTipoChannel
	GetPagosTipoChannelService(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error)
	DeletePagoTipoChannelService(ctx context.Context, id uint64) (erro error)
	CreatePagoTipoChannel(ctx context.Context, request administraciondtos.RequestPagoTipoChannel) (id uint64, erro error)

	//Busca lista de peticiones web services
	GetPeticionesService(filtro filtros.PeticionWebServiceFiltro) (peticiones administraciondtos.ResponsePeticionesWebServices, erro error)

	//SubirArchivos
	SubirArchivos(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse) (countArchivo int, erro error)

	// Archivos Subidos
	ObtenerArchivosSubidos(filtro filtros.Paginacion) (lisArchivosSubidos administraciondtos.ResponseArchivoSubido, erro error)

	// obtener archivo de cierreloterapipago
	ObtenerArchivoCierreLoteRapipago(nombre string) (result bool, err error)

	// reversiones o contracargo //
	// obtener cierre lote con reversiones en disputa
	GetCierreLoteEnDisputaServices(estadoDisputa int, request filtros.ContraCargoEnDisputa) (cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL, erro error)
	// obtener informacion de los pagos relacionados con los cierre de lotes en disputa
	GetPagosByTransactionIdsServices(filtro filtros.ContraCargoEnDisputa, cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL) (listaRevertidos administraciondtos.ResponseOperacionesContracargo, erro error)

	/* Preferences*/
	PostPreferencesService(request administraciondtos.RequestPreferences) (erro error)
	GetPreferencesService(request administraciondtos.RequestPreferences) (responsePreference dtos.ResponsePreference, erro error)
	DeletePreferencesService(request administraciondtos.RequestPreferences) (erro error)

	/*Obtener pagos para pruebas -> utilizada para generar movimientos en dev*/
	GetPagosDevService(filtro filtros.PagoFiltro) (response []entities.Pago, erro error)
	UpdatePagosDevService(pagos []entities.Pago) (pg []uint, erro error)

	BuildPagosMovDev(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	//? Servicios que permiten consultar datos de cllote pora herramieta wee
	GetConsultarClRapipagoService(filtro filtros.RequestClrapipago) (response administraciondtos.ResponseCLRapipago, erro error)

	// Permite consultas los datos de cierres de lote para la herramienta
	GetConsultarClMultipagoService(filtro filtros.RequestClMultipago) (response administraciondtos.ResponseClMultipago, erro error)

	GetCaducarOfflineIntentos() (intentosCaducados int, erro error)
	GetCaducarPagosExpirados(filtro filtros.PagoCaducadoFiltro) (pagosCaducados int, erro error)

	// servicio para consultar contruir y generar movimientos temporales
	GetPagosCalculoMovTemporalesService(filtro filtros.PagoIntentoFiltros) (pagosid []uint, erro error)
	BuildPagosCalculoTemporales(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoTemporalesResponse, erro error)
	GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error)

	ConciliacionPagosReportesService(filtro filtros.PagoFiltro) (valoresNoEncontrados []string, erro error)

	AsignarBancoIdRapipagoService(banco_id int64, rapipago_id int64) error

	//ContactoReportes
	CreateContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error)
	ReadContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (contactosFormato administraciondtos.ResponseGetContactosReportes, err error)
	UpdateContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error)
	DeleteContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error)

	// Usuarios Bloqueados
	CreateUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error)
	GetUsuariosBloqueadoService(filtro filtros.UsuarioBloqueadoFiltro) (response administraciondtos.ResponseUsuariosBloqueados, erro error)
	UpdateUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error)
	DeleteUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error)

	//Run endpoint
	CallFraudePersonas(cuil string) (object interface{}, erro error)

	//Soporte
	// CreateSoporteService(contactoReporte administraciondtos.RequestSoporte) (err error)
	// PutSoporteService(contactoReporte administraciondtos.RequestSoporte) (err error)

	//Endpoint para verficiar si el servicio esta activo
	EstadoApiService() (err error)

	// Historial de operaciones
	GetHistorialOperacionesService(filtro filtros.RequestHistorial) (response administraciondtos.ResponseHistorial, erro error)

	UpsertEnvioService(request administraciondtos.RequestEnvios) (erro error)
}

// variable que va a manejar la instancia del servicio
var admService *service

type service struct {
	repository     Repository
	apilinkService apilink.AplinkService
	commonsService commons.Commons
	utilService    util.UtilService
	webhook        webhook.RemoteRepository
	store          util.Store
}

func NewService(r Repository, s apilink.AplinkService, c commons.Commons, u util.UtilService, webhook webhook.RemoteRepository, storage util.Store) Service {
	admService = &service{
		repository:     r,
		apilinkService: s,
		commonsService: c,
		utilService:    u,
		webhook:        webhook,
		store:          storage,
	}
	return admService
}

// Resolve devuelve la instancia antes creada
func Resolve() *service {
	return admService
}

func _setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))

	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	return

}

/* +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

func (s *service) GetCuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseCuenta, erro error) {

	resp, erro := s.repository.GetCuenta(filtro)
	if erro != nil {
		return
	}

	response.FromCuenta(resp)

	return
}

func (s *service) GetSubcuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseSubcuenta, erro error) {

	/* resp, erro := s.repository.GetCuenta(filtro)
	if erro != nil {
		return
	} */
	resp, erro := s.repository.GetSubcuenta(filtro)
	if erro != nil {
		return
	}
	resp.Porcentaje *= 100

	response.FromSubcuenta(resp)

	return
}

func (s *service) GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {
	return s.apilinkService.GetConsultaDestinatarioService(requerimientoId, request)
}

func (s *service) GetPagoByID(pagoID int64) (*entities.Pago, error) {
	return s.repository.PagoById(pagoID)
}

func (s *service) GetPagosConsulta(apikey string, req administraciondtos.RequestPagosConsulta) (*[]administraciondtos.ResponsePagosConsulta, error) {
	err := req.IsValid()
	if err != nil {
		return nil, fmt.Errorf("validación %w", err)
	}
	var pagotiposIds []uint64
	var rangoFechas []string
	//obtengo los pagostipos relacionada con
	cuenta, err := s.repository.GetCuentaByApiKey(apikey)
	if err != nil {
		return nil, errors.New("error: " + err.Error())
	}
	for _, values := range *cuenta.Pagotipos {
		pagotiposIds = append(pagotiposIds, uint64(values.ID))
	}

	if len(req.FechaDesde) > 0 {
		fechaDesde, err := time.Parse("02-01-2006", req.FechaDesde)
		if err != nil {
			return nil, fmt.Errorf("formato de fecha desde incorrecto: %w", err)
		}
		fechaHasta, err := time.Parse("02-01-2006", req.FechaHasta)
		if err != nil {
			return nil, fmt.Errorf("formato de fecha hasta incorrecto: %w", err)
		}
		if fechaHasta.Sub(fechaDesde) < 0 {
			return nil, fmt.Errorf("periodo de consulta incorrecto")
		}
		if fechaHasta.Sub(fechaDesde).Hours()/24 > 7 {
			return nil, fmt.Errorf("período de consulta mayor a 7 días")
		}
		rangoFechas = append(rangoFechas, fechaDesde.Format("2006-01-02")+" 00:00:00", fechaHasta.Format("2006-01-02")+" 23:59:59")
		// cuenta, err := s.repository.GetCuentaByApiKey(apikey)
		// if err != nil {
		// 	return nil, errors.New("error: " + err.Error())
		// }
		// for _, values := range *cuenta.Pagotipos {
		// 	pagotiposIds = append(pagotiposIds, uint64(values.ID))
		// }

	}

	var uuidList []string
	var external_references []string

	if len(req.Uuid) > 0 {
		uuidList = append(uuidList, req.Uuid)
	}
	if len(req.Uuids) > 0 {
		uuidList = append(uuidList, req.Uuids...)
	}

	if len(req.ExternalReferences) > 0 {
		external_references = append(external_references, req.ExternalReferences...)
	}

	// estados que se requiere filtrar: 1- pendiente, 2- procesando, 6- expirado
	estadosPago := []uint64{
		1,
		2,
		6,
	}

	filtro := filtros.PagoFiltro{
		Uuids:              uuidList,
		ExternalReference:  req.ExternalReference,
		ExternalReferences: external_references,
		CargarPagoEstado:   true,
		Fecha:              rangoFechas,
		Notificado:         true,
		PagosTipoIds:       pagotiposIds,
		PagoEstadosIds:     estadosPago,
		// VisualizarPendientes: true,
	}

	pago, _, err := s.repository.GetPagos(filtro)
	if err != nil {
		return nil, fmt.Errorf("consultando a la base de datos: %w", err)
	}

	res := make([]administraciondtos.ResponsePagosConsulta, len(pago))

	for i, p := range pago {
		res[i].SetPago(p)
	}

	return &res, nil
}
func (s *service) ConsultarEstadoPagosService(requestValid administraciondtos.ParamsValidados, apiKey string, request administraciondtos.RequestPagosConsulta) (responsePagoEstado []administraciondtos.ResponseSolicitudPago, registrosAfectados bool, erro error) {
	err := request.IsValid()
	if err != nil {
		erro = fmt.Errorf("validación %w", err)
		return
	}
	var pagotiposIds []uint64
	var rangoFechas []string
	var uuidList []string
	var referenciaExterna string
	//obtengo los pagostipos relacionada con
	cuenta, err := s.repository.GetCuentaByApiKey(apiKey)
	if err != nil {
		erro = errors.New("error: " + err.Error())
		return
	}
	for _, values := range *cuenta.Pagotipos {
		pagotiposIds = append(pagotiposIds, uint64(values.ID))
	}

	if requestValid.Uuuid {
		uuidList = append(uuidList, request.Uuid)
	}
	if requestValid.ExternalReference {
		referenciaExterna = request.ExternalReference
	}
	if requestValid.RangoFecha {
		if len(request.FechaDesde) > 0 && len(request.FechaHasta) > 0 {
			fechaDesde, err := time.Parse("02-01-2006", request.FechaDesde)
			if err != nil {
				erro = fmt.Errorf("formato de fecha desde incorrecto: %w", err)
				return
			}
			fechaHasta, err := time.Parse("02-01-2006", request.FechaHasta)
			if err != nil {
				erro = fmt.Errorf("formato de fecha hasta incorrecto: %w", err)
				return
			}
			if fechaHasta.Sub(fechaDesde) < 0 {
				erro = fmt.Errorf("periodo de consulta incorrecto")
				return
			}
			if fechaHasta.Sub(fechaDesde).Hours()/24 > 7 {
				erro = fmt.Errorf("período de consulta mayor a 7 días")
				return
			}
			rangoFechas = append(rangoFechas, fechaDesde.Format("2006-01-02")+" 00:00:00", fechaHasta.Format("2006-01-02")+" 23:59:59")
		}
	}
	if requestValid.Uuids {
		uuidList = append(uuidList, request.Uuids...)
	}

	filtroEstado := filtros.PagoEstadoFiltro{
		Nombre: "PENDING",
	}
	entityPagoEstado, err := s.repository.GetPagoEstado(filtroEstado)
	if err != nil {
		erro = errors.New("no se pudo obtener estado de pago" + err.Error())
		return
	}

	filtro := filtros.PagoFiltro{
		Uuids:             uuidList,
		ExternalReference: referenciaExterna,
		Fecha:             rangoFechas,
		PagosTipoIds:      pagotiposIds,
		PagoEstadosId:     uint64(entityPagoEstado.ID),
		CargarPagoTipos:   true,
		CargarPagoEstado:  true,
		CargaPagoIntentos: true,
		//Notificado:           true,
		//VisualizarPendientes: true,

	}
	entityPagos, err := s.repository.ConsultarEstadoPagosRepository(requestValid, filtro)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	if len(entityPagos) == 0 {
		registrosAfectados = false
		return
	}
	for _, value := range entityPagos {
		var temporalPagoEstado administraciondtos.ResponseSolicitudPago
		temporalPagoEstado.SolicitudEntityToDtos(value)
		// temporalPagoEstado.PagoIntento[0].GrossFee = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].GrossFee, 2)
		// temporalPagoEstado.PagoIntento[0].NetFee = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].NetFee, 2)
		// temporalPagoEstado.PagoIntento[0].FeeIva = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].FeeIva, 2)
		responsePagoEstado = append(responsePagoEstado, temporalPagoEstado)
	}
	registrosAfectados = true
	return
}
func (s *service) GetPagosIntentosByTransaccionIdService(filtroPagoIntento filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error) {
	pagosIntentos, err := s.repository.GetPagosIntentos(filtroPagoIntento)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}

	return
}

func (s *service) GetCuentasByCliente(cliente int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Cuenta, error) {
	from := (number - 1) * size
	data, total, e := s.repository.CuentaByClientePage(cliente, size, from)
	lastPage := math.Ceil(float64(total) / float64(size))

	meta := dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	links := dtos.Links{
		First: fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, 1, size),
		Last:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, meta.Page.LastPage, size),
		Next:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, (number + 1), size),
		Prev:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, (number - 1), size),
	}

	return &meta, &links, data, e
}
func (s *service) GetSubcuentasByCuenta(cuenta int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Subcuenta, error) {
	from := (number - 1) * size
	data, total, e := s.repository.SubcuentaByCuentaPage(cuenta, size, from)

	logs.Info(data)

	if len(*data) > 0 {
		modifiedData := make([]entities.Subcuenta, len(*data))
		for i, subcuenta := range *data {

			decimalPorcentaje := decimal.NewFromFloat(subcuenta.Porcentaje)
			decimalFactor := decimal.NewFromFloat(100.0)
			result := decimalPorcentaje.Mul(decimalFactor)

			subcuenta.Porcentaje, _ = result.Float64()
			modifiedData[i] = subcuenta
		}
		data = &modifiedData
	}

	lastPage := math.Ceil(float64(total) / float64(size))

	meta := dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	links := dtos.Links{
		First: fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cuenta, 1, size),
		Last:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cuenta, meta.Page.LastPage, size),
		Next:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cuenta, (number + 1), size),
		Prev:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cuenta, (number - 1), size),
	}

	return &meta, &links, data, e
}

func (s *service) PostCuenta(ctx context.Context, request administraciondtos.CuentaRequest) (ok bool, err error) {

	err = request.IsVAlid(false)

	if err != nil {
		return
	}
	/* 	Se comenta la validacion de CVU/CBU repetido. 28-12-2022 */
	//Valido si ya existe una cuenta con el cbu/cvu registrado
	// var filtro filtros.CuentaFiltro
	// if len(request.Cbu) > 0 {
	// 	filtro.Cbu = request.Cbu
	// 	response, err := s.repository.GetCuenta(filtro)

	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	if response.ID > 0 {
	// 		return false, fmt.Errorf(ERROR_CBU_REGISTRADO)
	// 	}

	// } else {
	// 	filtro.Cvu = request.Cvu
	// 	response, err := s.repository.GetCuenta(filtro)

	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	if response.ID > 0 {
	// 		return false, fmt.Errorf(ERROR_CVU_REGISTRADO)
	// 	}
	// }
	//Creo un uuid automaticamente
	request.Apikey = s.commonsService.NewUUID()

	cuenta := request.ToCuenta()

	ok, err = s.repository.SaveCuenta(ctx, &cuenta)
	if err != nil {
		return false, err
	}

	return ok, err
}

func (s *service) PostSubcuenta(ctx context.Context, request []administraciondtos.SubcuentaRequest) (ok bool, err error) {

	var subcuentas administraciondtos.ArraySubcuentaRequest
	subcuentas.ArraySubcuentas = request

	tipoOperacion, err := subcuentas.ArrayRequestTypeSaved()
	if err != nil {
		return false, err
	}

	filter := filtros.CuentaFiltro{
		Id: request[0].Id,
	}

	subcuenta, err := s.repository.GetSubcuenta(filter)
	if err != nil {
		return false, err
	}
	subcuentasByCuentasId, err := s.repository.GetSubcuentasByCuentaId(subcuenta.CuentasID)
	if err != nil {
		return false, err
	}

	switch tipoOperacion {
	case "ACTUALIZAR":
		ok, err := util.ValidateRequestUpdated(&request, &subcuentasByCuentasId)
		if err != nil {
			return ok, err
		}

	case "CREAR":
		subcuentasByCuentasId, err := s.repository.GetSubcuentasByCuentaId(uint(request[0].CuentaID))
		if err != nil {
			return false, err
		}
		ok, err := util.ValidateRequestCreated(&request, &subcuentasByCuentasId)
		if err != nil {
			return ok, err
		}

	case "CREARACTUALIZAR":
		ok, err := util.ValidateRequestCreatedUpdated(&request, &subcuentasByCuentasId)
		if err != nil {
			return ok, err
		}
	}

	//Actualizo o creo subcuentas desde el repository
	// for _, v := range request {
	// 	if v.Modificado == 1 || v.Id == 0 {

	// 		subcuenta := v.ToCuenta()
	// 		ok, err = s.repository.SaveSubcuenta(ctx, &subcuenta)
	// 		if err != nil {
	// 			return false, err
	// 		}
	// 	}
	// }

	saved, err := s.repository.GuardarSubcuentas(ctx, request)
	if err != nil {
		return
	}

	return saved, nil
}

func (s *service) PostEditSubcuenta(ctx context.Context, request administraciondtos.SubcuentaRequest) (ok bool, err error) {

	if !(request.Id <= 0 || request.CuentaID <= 0) {
		fmt.Println("entro", request.Id, request.CuentaID)
		err = request.IsVAlid(true)
		if err != nil {
			return
		}

	} else {
		err = request.IsVAlid(false)
		if err != nil {
			return
		}
	}
	var filtro filtros.CuentaFiltro
	filtro.Id = request.Id
	filtro.Cbu = request.Cbu

	subcuenta, err := s.GetSubcuenta(filtro)
	if err != nil {
		return
	}
	if len(subcuenta.Cbu) > 0 {
		err = fmt.Errorf("Ya existe una cuenta con ese CBU.")
		return
	}

	var subcuentas []administraciondtos.SubcuentaRequest
	subcuentas = append(subcuentas, request)
	ok, err = s.repository.GuardarSubcuentas(ctx, subcuentas)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) SetApiKeyService(ctx context.Context, request *administraciondtos.CuentaRequest) (erro error) {

	if request.Id < 1 {
		return fmt.Errorf(ERROR_ID)
	}

	request.Apikey = s.commonsService.NewUUID()

	cuenta := request.ToCuenta()

	erro = s.repository.SetApiKey(ctx, cuenta)

	if erro != nil {
		return erro
	}

	return
}

func (s *service) UpdateCuentaService(ctx context.Context, request administraciondtos.CuentaRequest) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}
	filtro := filtros.CuentaFiltro{
		DistintoId: request.Id,
	}

	if len(request.Cbu) > 0 {
		filtro.Cbu = request.Cbu
		response, err := s.repository.GetCuenta(filtro)

		if err != nil {
			return err
		}
		if response.ID > 0 {
			return fmt.Errorf(ERROR_CBU_REGISTRADO)
		}

	} else {
		filtro.Cvu = request.Cvu
		response, err := s.repository.GetCuenta(filtro)

		if err != nil {
			return err
		}
		if response.ID > 0 {
			return fmt.Errorf(ERROR_CVU_REGISTRADO)
		}
	}

	cuenta := request.ToCuenta()

	erro = s.repository.UpdateCuenta(ctx, cuenta)

	if erro != nil {
		return erro
	}

	return
}

func (s *service) DeleteCuentaService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	erro = s.repository.DeleteCuenta(id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	if erro != nil {
		return erro
	}

	return
}

func (s *service) GetCuentaByApiKeyService(apikey string) (result bool, erro error) {
	cuenta, err := s.repository.GetCuentaByApiKey(apikey)
	if err != nil {
		log := entities.Log{
			Tipo:          "info",
			Funcionalidad: "GetCuentaByApiKey",
			Mensaje:       err.Error(),
		}
		err = s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error("error al intentar registrar logs de erro en GetCuentaByApiKey")
		}
		erro = errors.New("api-key invalido.")
		return
	}
	result = false
	if len(cuenta.Apikey) > 0 {
		result = true
	}
	return
}

func (s *service) GetImpuestosService(filtro filtros.ImpuestoFiltro) (response administraciondtos.ResponseImpuestos, erro error) {
	impuestosEntity, totalFilas, err := s.repository.GetImpuestosRepository(filtro)
	if err != nil {
		return
	}
	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, totalFilas)
	}
	for _, valueImpuesto := range impuestosEntity {
		impuesto := administraciondtos.ResponseImpuesto{}
		impuesto.FromImpuesto(valueImpuesto)
		response.Impuestos = append(response.Impuestos, impuesto)
	}
	return
}

func (s *service) PostImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (id uint64, erro error) {

	erro = filtro.Validar()
	if erro != nil {
		return 0, erro
	}

	impuesto := filtro.ToImpuesto(false)

	id, erro = s.repository.CreateImpuestoRepository(ctx, impuesto)
	if erro != nil {
		return 0, erro
	}

	return
}

func (s *service) UpdateImpuestoService(ctx context.Context, request administraciondtos.ImpuestoRequest) (erro error) {

	erro = request.Validar()
	if erro != nil {
		return erro
	}

	impuesto := request.ToImpuesto(true)

	return s.repository.UpdateImpuestoRepository(ctx, impuesto)

}

func (s *service) PostPagotipo(ctx context.Context, pagotipo *entities.Pagotipo) (bool, error) {
	var err error

	res, err := s.repository.SavePagotipo(pagotipo)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	return res, err
}

func (s *service) PostCuentaComision(ctx context.Context, comision *entities.Cuentacomision) error {

	err := s.repository.SaveCuentacomision(comision)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return err
}

func (s *service) CreateNotificacionService(notificacion entities.Notificacione) error {

	err := s.utilService.CreateNotificacionService(notificacion)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateLogService(log entities.Log) error {
	err := s.utilService.CreateLogService(log)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetPagosEstadosService(buscarPorFinal, final bool) (estados []entities.Pagoestado, erro error) {
	filtro := filtros.PagoEstadoFiltro{BuscarPorFinal: buscarPorFinal, Final: final}
	estados, erro = s.repository.GetPagosEstados(filtro)

	return
}

func (s *service) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {
	estados, erro = s.repository.GetPagosEstados(filtro)

	return
}

func (s *service) GetPagosEstadosExternoService(filtro filtros.PagoEstadoExternoFiltro) (estadosExternos []entities.Pagoestadoexterno, erro error) {
	estadosExternos, erro = s.repository.GetPagosEstadosExternos(filtro)
	if erro != nil {
		return
	}
	return
}

func (s *service) GetPagosService(filtro filtros.PagoFiltro) (response administraciondtos.ResponsePagos, erro error) {

	pagos, total, erro := s.repository.GetPagos(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	var listaPagoIntentos []uint64

	// recorrer cada uno de los pagos obtenidos
	for _, p := range pagos {
		r := administraciondtos.ResponsePago{
			Identificador:     p.ID,
			Fecha:             p.CreatedAt,
			ExternalReference: p.ExternalReference,
			PayerName:         p.PayerName,
		}

		if p.PagoEstados.ID > 0 {
			r.Estado = string(p.PagoEstados.Estado)
			r.NombreEstado = p.PagoEstados.Nombre

		}

		if p.PagosTipo.ID > 0 {
			r.Pagotipo = p.PagosTipo.Pagotipo
			if p.PagosTipo.Cuenta.ID > 0 {
				r.Cuenta = p.PagosTipo.Cuenta.Cuenta
			}
		}

		if len(p.PagoIntentos) > 0 {
			r.Amount = p.PagoIntentos[len(p.PagoIntentos)-1].Amount
			r.FechaPago = p.PagoIntentos[len(p.PagoIntentos)-1].PaidAt
			if p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.ID > 0 && p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.ID > 0 {
				r.Channel = p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.Channel
				r.NombreChannel = p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.Nombre
			}
			r.UltimoPagoIntentoId = uint64(p.PagoIntentos[len(p.PagoIntentos)-1].ID)
			listaPagoIntentos = append(listaPagoIntentos, uint64(p.PagoIntentos[len(p.PagoIntentos)-1].ID))
		}

		if filtro.CargarPagosItems {
			var listaItems []administraciondtos.PagoItems
			for _, items := range p.Pagoitems {
				listaItems = append(listaItems, administraciondtos.PagoItems{
					Descripcion:   items.Description,
					Identificador: items.Identifier,
					Cantidad:      int64(items.Quantity),
					Monto:         float64(items.Amount),
				})
			}
			// r.PagoItems = listaItems
		}

		response.Pagos = append(response.Pagos, r)
	}

	FiltroMovimientos := filtros.MovimientoFiltro{
		PagoIntentosIds: listaPagoIntentos,
		CuentaId:        filtro.CuentaId,
	}

	movimientos, _, erro := s.repository.GetMovimientos(FiltroMovimientos)

	if erro != nil {
		return
	}

	var listaMovimientos []uint64
	for i := range movimientos {
		listaMovimientos = append(listaMovimientos, uint64(movimientos[i].ID))
	}

	filtroTransferencias := filtros.TransferenciaFiltro{
		MovimientosIds: listaMovimientos,
	}

	transferencias, _, erro := s.repository.GetTransferencias(filtroTransferencias)

	if erro != nil {
		return
	}

	for i := range response.Pagos {

		for j := range transferencias {
			if response.Pagos[i].UltimoPagoIntentoId == transferencias[j].Movimiento.PagointentosId {
				response.Pagos[i].ReferenciaBancaria = transferencias[j].ReferenciaBancaria
				response.Pagos[i].TransferenciaId = uint64(transferencias[j].ID)
				response.Pagos[i].FechaTransferencia = transferencias[j].FechaOperacion.Format("02-01-2006")
			}
		}
		// si el estado es PAID se acumula en un atributo de la struct ResponsePagos
		if response.Pagos[i].Estado == "PAID" {
			response.SaldoPendiente += response.Pagos[i].Amount
		}

		// si el estado es ACCREDITED y no tiene fecha de transferencia se acumula en un atributo de la struct ResponsePagos
		if response.Pagos[i].Estado == "ACCREDITED" && len(response.Pagos[i].FechaTransferencia) == 0 {
			response.SaldoDisponible += response.Pagos[i].Amount
		}
	}

	return

}

func (as *service) GetPagos(filter filtros.PagoFiltro) (administraciondtos.ResponsePagos, error) {
	var response administraciondtos.ResponsePagos

	// 1. OBTENER PAGOS
	pagos, total, err := as.repository.GetPagosRepository(filter)

	if err != nil {
		return response, err
	}

	if filter.Number > 0 && filter.Size > 0 {
		response.Meta = _setPaginacion(filter.Number, filter.Size, total)
	}

	response.Pagos = pagos

	return response, nil
}

func (as *service) GetItemsPagos(filter filtros.PagoItemFiltro) ([]administraciondtos.PagoItems, error) {
	// 1. OBTENER PAGOS

	if filter.PagoId == 0 {
		return nil, errors.New("debe enviar un id de pago")
	}

	pagositems, err := as.repository.GetItemsPagos(filter)

	if err != nil {
		return nil, err
	}

	return pagositems, nil
}

func (as *service) GetSaldoPagoCuenta(filter filtros.PagoFiltro) (administraciondtos.ResponsePagos, error) {
	var response administraciondtos.ResponsePagos

	// 1. OBTENER PAGOS
	pagos, _, err := as.repository.GetPagosRepository(filter)

	if err != nil {
		return response, err
	}

	// recorrer cada uno de los pagos obtenidos
	for _, pago := range pagos {

		if pago.Estado == "PAID" {
			response.SaldoPendiente += pago.Amount
		}

		// si el estado es ACCREDITED y no tiene fecha de transferencia se acumula en un atributo de la struct ResponsePagos
		if pago.Estado == "ACCREDITED" && len(pago.FechaTransferencia) == 0 {
			response.SaldoDisponible += pago.Amount
		}
	}

	return response, nil
}

func (s *service) GetPlanCuotas(idMedioPago uint) (response []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	response, erro = s.repository.GetPlanCuotasByMedioPago(idMedioPago)
	if erro != nil {
		erro = errors.New("problema al obtener plan de cuotas - " + erro.Error())
	}
	return
}

func (s *service) GetInteresesPlanes(fecha string) (planes []administraciondtos.PlanCuotasResponse, erro error) {
	// fmt.Printf("Fecha actual %v", fecha)
	fechaActual, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	response, erro := s.repository.GetInstallments(fechaActual)
	if erro != nil {
		return
	}
	for _, valueMedioPagoInstallment := range response {
		var details []administraciondtos.PlanCuotasResponseDetalle
		var installmentTemp administraciondtos.PlanCuotasResponse
		for _, valueInstallment := range valueMedioPagoInstallment.Installments {
			// logs.Info(valueInstallment.VigenciaHasta)
			if valueInstallment.VigenciaHasta == nil {
				installmentTemp = administraciondtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, administraciondtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}
			// logs.Info("============================\n")
			// fmt.Printf("after %v - before %v \n", valueInstallment.VigenciaDesde.After(fechaActual), valueInstallment.VigenciaDesde.Before(fechaActual))
			// logs.Info("============================\n")

			// fmt.Printf("%v--%v--%v \n", valueInstallment.VigenciaDesde, fechaActual, valueInstallment.VigenciaHasta)
			// logs.Info("============================\n")
			// logs.Info("============================\n")

			if (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.Equal(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Equal(*valueInstallment.VigenciaHasta)) {
				installmentTemp = administraciondtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, administraciondtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}

		}
		planes = append(planes, administraciondtos.PlanCuotasResponse{
			Id:                      installmentTemp.Id,
			Descripcion:             installmentTemp.Descripcion,
			MediopagoinstallmentsID: installmentTemp.MediopagoinstallmentsID,
			Installmentdetail:       details,
		})
	}
	return
}

func (s *service) GetAllInstallmentsById(installments_id int64) (planesCuotas []administraciondtos.InstallmentsResponse, erro error) {
	plancuotas, err := s.repository.GetAllInstallmentsById(uint(installments_id))
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range plancuotas {
		var temporalPlanCuotas administraciondtos.InstallmentsResponse
		temporalPlanCuotas.EntityToDtos(value)
		planesCuotas = append(planesCuotas, temporalPlanCuotas)
	}
	return
}

func (s *service) ModificarEstadoPagosExpirados() (erro error) {

	// Busco el tiempo de expiración de los pagos si no existe lo creo con valor de 30 dias

	filtroConf := filtros.ConfiguracionFiltro{
		Nombre: "TIEMPO_EXPIRACION_PAGOS",
	}

	configuracion, erro := s.utilService.GetConfiguracionService(filtroConf)

	if erro != nil {
		return
	}

	if configuracion.Id == 0 {

		config := administraciondtos.RequestConfiguracion{
			Nombre:      "TIEMPO_EXPIRACION_PAGOS",
			Descripcion: "Tiempo en días para que expire un pago que está en estado pending ",
			Valor:       "30",
		}

		_, erro = s.utilService.CreateConfiguracionService(config)

		if erro != nil {
			return
		}

	}

	// Busco el pagoEstado con nobre de pending

	filtroPending := filtros.PagoEstadoFiltro{
		Nombre: "Pending",
	}

	pagoEstadoPending, erro := s.repository.GetPagoEstado(filtroPending)

	if erro != nil {
		return
	}

	// Busco los pagos que están en el estado pending y que están expirados

	filtroPagos := filtros.PagoFiltro{
		PagoEstadosId:    uint64(pagoEstadoPending.ID),
		TiempoExpiracion: configuracion.Valor,
	}

	pagos, _, erro := s.repository.GetPagos(filtroPagos)

	if erro != nil {
		return
	}

	if len(pagos) == 0 {
		return
	}

	// Busco el pago estado expirado
	filtroExpired := filtros.PagoEstadoFiltro{
		Nombre: "Expired",
	}

	pagoEstadoExpired, erro := s.repository.GetPagoEstado(filtroExpired)

	if erro != nil {
		return
	}

	erro = s.repository.UpdateEstadoPagos(pagos, uint64(pagoEstadoExpired.ID))

	return

}

func (s *service) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (response administraciondtos.ResponseConfiguraciones, erro error) {

	configuraciones, total, erro := s.repository.GetConfiguraciones(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, c := range configuraciones {

		r := administraciondtos.ResponseConfiguracion{}
		r.FromEntity(c)

		response.Data = append(response.Data, r)
	}

	return
}

func (s *service) UpdateConfiguracionService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error) {

	erro = request.IsValid(true)

	if erro != nil {
		return
	}

	config := request.ToEntity(true)

	erro = s.repository.UpdateConfiguracion(ctx, config)

	if erro != nil {
		return
	}
	return

}

// ABM CLIENTES
func (s *service) GetClienteService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacion, erro error) {
	cliente, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}
	var cli administraciondtos.ResponseFacturacion
	cli.FromEntity(cliente)
	response = cli

	return
}

func (s *service) GetClientesService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacionPaginado, erro error) {

	clientes, total, erro := s.repository.GetClientes(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cliente := range clientes {

		var cli administraciondtos.ResponseFacturacion
		cli.FromEntity(cliente)

		response.Clientes = append(response.Clientes, cli)
	}

	return
}

func (s *service) CreateClienteService(ctx context.Context, request administraciondtos.ClienteRequest) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ClienteFiltro{
		Cuit: request.Cuit,
	}

	response, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}

	if response.ID > 0 {
		erro = fmt.Errorf(ERROR_CLIENTE_REGISTRADO)
		return
	}

	cliente := request.ToCliente(false)

	return s.repository.CreateCliente(ctx, cliente)

}

func (s *service) UpdateClienteService(ctx context.Context, cliente administraciondtos.ClienteRequest) (erro error) {

	erro = cliente.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ClienteFiltro{
		DistintoId: cliente.Id,
		Cuit:       cliente.Cuit,
	}

	cliente_existente_sin_modificar, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}

	if cliente_existente_sin_modificar.ID == 0 {
		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)
		return
	}

	clienteModificado := cliente.ToCliente(true)

	erro = s.repository.UpdateCliente(ctx, clienteModificado)
	if erro != nil {
		logs.Error(erro)
		return
	}

	return
}

func (s *service) DeleteClienteService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del cliente es invalido")
		return
	}

	erro = s.repository.DeleteCliente(ctx, id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	return
}

//ABMRUBROS

func (s *service) GetRubroService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubro, erro error) {
	rubro, erro := s.repository.GetRubro(filtro)

	if erro != nil {
		return
	}
	response = administraciondtos.ResponseRubro{
		Id:    rubro.ID,
		Rubro: rubro.Rubro,
	}

	return
}

func (s *service) GetRubrosService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubros, erro error) {

	rubros, total, erro := s.repository.GetRubros(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, rubro := range rubros {

		r := administraciondtos.ResponseRubro{
			Id:    rubro.ID,
			Rubro: rubro.Rubro,
		}

		response.Rubros = append(response.Rubros, r)
	}

	return
}

func (s *service) CreateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	rubro := request.ToRubro(false)

	return s.repository.CreateRubro(ctx, rubro)

}

func (s *service) UpdateRubroService(ctx context.Context, rubro administraciondtos.RubroRequest) (erro error) {

	erro = rubro.IsVAlid(true)

	if erro != nil {
		return
	}

	rubroModificado := rubro.ToRubro(true)

	return s.repository.UpdateRubro(ctx, rubroModificado)

}

//ABM PAGO TIPOS

func (s *service) GetPagoTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagoTipo, erro error) {

	pagoTipo, erro := s.repository.GetPagoTipo(filtro)

	if erro != nil {
		return
	}

	response.FromPagoTipo(pagoTipo)

	return
}

func (s *service) GetPagosTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagosTipo, erro error) {

	pagosTipo, total, erro := s.repository.GetPagosTipo(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, pagoTipo := range pagosTipo {
		var ch []administraciondtos.CanalesPago
		var cuotas []administraciondtos.CuotasPago
		for _, channel := range pagoTipo.Pagotipochannel {
			c := administraciondtos.CanalesPago{
				ChannelsId: channel.Channel.ID,
				Channel:    channel.Channel.Channel,
				Nombre:     channel.Channel.Nombre,
			}
			ch = append(ch, c)
		}

		for _, cuota := range pagoTipo.Pagotipoinstallment {
			cuo := administraciondtos.CuotasPago{
				Nro: cuota.Cuota,
			}
			cuotas = append(cuotas, cuo)
		}

		r := administraciondtos.ResponsePagoTipo{}
		r.IncludedChannels = ch
		r.IncludedInstallments = cuotas
		r.FromPagoTipo(pagoTipo)

		response.PagosTipo = append(response.PagosTipo, r)
	}

	return
}

func (s *service) CreatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	pagoTipo := request.ToPagoTipo(false)

	for _, channel := range request.IncludedChannels {
		filtro := filtros.ChannelFiltro{
			Id: uint(channel),
		}
		ch, err := s.repository.GetChannel(filtro)
		if err != nil && ch.ID == 0 {
			erro = fmt.Errorf("el id del channels es invalido")
			return 0, erro
		}
	}

	return s.repository.CreatePagoTipo(ctx, pagoTipo, request.IncludedChannels, request.IncludedInstallments)

}

func (s *service) UpdatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	pagoTipoModificado := request.ToPagoTipo(true)

	filtro := filtros.PagoTipoFiltro{
		Id:                     pagoTipoModificado.ID,
		CargarTipoPagoChannels: true,
	}
	pagotipo, err := s.repository.GetPagoTipo(filtro)
	if err != nil {
		erro = err
		return
	}

	var channels []int64
	var cuotas []string
	for _, p := range pagotipo.Pagotipochannel {
		channels = append(channels, int64(p.Channel.ID))
	}
	for _, p := range pagotipo.Pagotipoinstallment {
		cuotas = append(cuotas, p.Cuota)
	}

	channelAdd, channelDelete := commons.DifferenceInt(request.IncludedChannels, channels)
	updateChannels := administraciondtos.RequestPagoTipoChannels{
		Add:    channelAdd,
		Delete: channelDelete,
	}

	cuotasAdd, cuotasDelete := commons.DifferenceString(request.IncludedInstallments, cuotas)
	updateCuotas := administraciondtos.RequestPagoTipoCuotas{
		Add:    cuotasAdd,
		Delete: cuotasDelete,
	}

	return s.repository.UpdatePagoTipo(ctx, pagoTipoModificado, updateChannels, updateCuotas)

}

func (s *service) DeletePagoTipoService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del pago tipo es invalido")
		return
	}

	return s.repository.DeletePagoTipo(ctx, id)

}

//ABM CHANNEL

func (s *service) GetChannelService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannel, erro error) {

	channel, erro := s.repository.GetChannel(filtro)

	if erro != nil {
		return
	}

	response.FromChannel(channel)

	return
}

func (s *service) GetChannelsService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannels, erro error) {

	channels, total, erro := s.repository.GetChannels(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, channel := range channels {

		r := administraciondtos.ResponseChannel{}
		r.FromChannel(channel)

		response.Channels = append(response.Channels, r)
	}

	return
}

func (s *service) CreateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	channel := request.ToChannel(false)

	return s.repository.CreateChannel(ctx, channel)

}

func (s *service) UpdateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	channelModificado := request.ToChannel(true)

	return s.repository.UpdateChannel(ctx, channelModificado)

}

func (s *service) DeleteChannelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del channel es invalido")
		return
	}

	return s.repository.DeleteChannel(ctx, id)

}

//ABM CUENTAS COMISION

func (s *service) GetCuentaComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentaComision, erro error) {

	cuentaComision, erro := s.repository.GetCuentaComision(filtro)

	if erro != nil {
		return
	}

	response.FromCuentaComision(cuentaComision)

	return
}

func (s *service) GetCuentasComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentasComision, erro error) {

	cuentasComsion, total, erro := s.repository.GetCuentasComisiones(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cuentaComision := range cuentasComsion {

		r := administraciondtos.ResponseCuentaComision{}
		r.FromCuentaComision(cuentaComision)

		response.CuentasComision = append(response.CuentasComision, r)
	}

	return
}

func (s *service) CreateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return 0, erro
	}

	cuentaComision := request.ToCuentaComision(false)

	return s.repository.CreateCuentaComision(ctx, cuentaComision)

}

func (s *service) UpdateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (erro error) {
	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return erro
	}

	cuentaComisionModificada := request.ToCuentaComision(true)

	return s.repository.UpdateCuentaComision(ctx, cuentaComisionModificada)

}

/*
	 func (s *service) DeleteSubcuentaService(ctx context.Context, id uint64) (erro error) {

		if id < 1 {
			erro = fmt.Errorf("el id de la subcuenta es invalido")
			return
		}

		return s.repository.DeleteSubcuenta(ctx, id)

}
*/
func (s *service) DeleteSubcuenta(ctx context.Context, request *[]administraciondtos.SubcuentaRequest) (ok bool, erro error) {

	ok, erro = util.ValidateRequestDeleted(request)
	if erro != nil {
		return
	}
	// Actualizo las cuentas o elimino las cuentas
	for _, v := range *request {
		if v.Id != 0 {
			if !v.Eliminado {

				var subcuentas []administraciondtos.SubcuentaRequest
				subcuentas = append(subcuentas, v)
				ok, erro = s.repository.GuardarSubcuentas(ctx, subcuentas)
				if erro != nil {
					return false, erro
				}
			} else {
				erro = s.repository.DeleteSubcuenta(ctx, uint64(v.Id))
				if erro != nil {
					return
				}
			}
		}
	}

	return true, nil

}

func (s *service) DeleteCuentaComisionService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id de cuenta comision es invalido")
		return
	}

	return s.repository.DeleteCuentaComision(ctx, id)

}

// Solicitud de cuenta
func (s *service) SendSolicitudCuenta(solicitudRequest administraciondtos.SolicitudCuentaRequest) (erro error) {

	// Valido los datos de entrada
	erro = solicitudRequest.IsValid()

	if erro != nil {
		return
	}

	entidadSolicitud := solicitudRequest.ToSolicitudEntity()
	// guadar la data de solicitud por medio del repository
	erro = s.repository.CreateSolicitudRepository(entidadSolicitud)

	if erro != nil {
		return
	}

	// Crear mensaje
	to := []string{
		config.EMAIL_TO_SOLICITUD_CUENTA,
	}

	from := config.EMAIL_FROM_SOLICITUD_CUENTA

	t, erro := template.ParseFiles("./api/views/solicitud_cuenta.html")

	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("no se pudo crear el template. %s", erro.Error()))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}
	buf := new(bytes.Buffer)
	erro = t.Execute(buf, solicitudRequest)
	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("no se pudo ejecutar el template. %s", erro.Error()))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}

	message := s.commonsService.CreateMessage(to, from, buf.String(), "Solicitud de Cuenta")

	// password := config.EMAIL_PASS
	smtpUsername := config.SMTP_USERNAME
	smtpPassword := config.SMTP_PASSWORD

	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort

	// Authentication.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// 4- Sending email.
	erro = smtp.SendMail(address, auth, from, to, []byte(message))

	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("%s. No se pudo enviar el email de solicitud de cuenta. Solicitante: %s ,  Email: %s", erro.Error(), solicitudRequest.Razonsocial, solicitudRequest.Email))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}

	return
}
func (s *service) UpdateConfiguracionSendEmailService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error) {

	//Busco todos los clientes activos del sistema para recuperar sus correos
	filtroCliente := filtros.ClienteFiltro{}

	response, erro := s.GetClientesService(filtroCliente)

	if len(response.Clientes) < 1 {
		return
	}
	//Cargo la lista de correos para enviar
	var emails []string
	for _, c := range response.Clientes {
		if !tools.EsStringVacio(c.Email) {
			emails = append(emails, c.Email)
		}
	}

	// Crear mensaje
	to := emails

	from := config.EMAIL_FROM

	t, erro := template.ParseFiles("../api/views/terminos_condiciones.html")

	ruta := administraciondtos.TerminosCondiciones{
		Ruta: config.RUTA_BASE_HOME_PAGE + "terminos-politicas",
	}

	if erro != nil {
		s._buildLog(erro, "SendTerminosCondiciones")
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}
	buf := new(bytes.Buffer)
	erro = t.Execute(buf, ruta)
	if erro != nil {
		s._buildLog(erro, "SendTerminosCondiciones")
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}

	message := s.commonsService.CreateMessage(to, from, buf.String(), "Actualización Terminos y Condiciones")

	password := config.EMAIL_PASS
	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort

	//Modifico los terminos y condiciones para luego enviar el mensaje
	s.repository.BeginTx()

	defer func() {
		if erro != nil {
			s.repository.RollbackTx()
		}
		s.repository.CommitTx()
	}()

	erro = s.UpdateConfiguracionService(ctx, request)

	if erro != nil {
		// s.repository.RollbackTx()
		return
	}

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4- Sending email.
	erro = smtp.SendMail(address, auth, from, to, []byte(message))

	if erro != nil {
		s.repository.RollbackTx()
		s._buildNotificacion(erro, entities.NotivicacionEnvioEmail, fmt.Sprintf("%s. No se pudo enviar el email de actualización de terminos y condiciones a los clientes.", erro.Error()))
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}

	// s.repository.CommitTx()

	return
}

func (s *service) _buildNotificacion(erro error, tipo entities.EnumTipoNotificacion, descripcion string) {
	notificacion := entities.Notificacione{
		Tipo:        tipo,
		Descripcion: descripcion,
	}
	s.CreateNotificacionService(notificacion)
}

func (s *service) _buildLog(erro error, funcionalidad string) {

	log := entities.Log{
		Tipo:          entities.Error,
		Mensaje:       erro.Error(),
		Funcionalidad: funcionalidad,
	}

	err := s.CreateLogService(log)

	if err != nil {
		mensaje := fmt.Sprintf("Crear Log: %s.  %s", err.Error(), erro.Error())
		logs.Error(mensaje)
	}

}

func (s *service) RetiroAutomaticoClientes(ctx context.Context) (movimientosidcomision administraciondtos.RequestMovimientosId, erro error) {

	//Buscar todos los clientes que tengan habilitado la opción de retiro automático
	// ? se habilita transferencias para automaticas para varios clientes
	filtroCliente := filtros.ClienteFiltro{
		RetiroAutomatico: true,
		CargarCuentas:    true,
	}

	clientes, _, erro := s.repository.GetClientes(filtroCliente)

	if erro != nil {
		return
	}

	// variable creada en el caso de que existan errore al enviar transferenicas a cuentas
	// var responseTransferencias []administraciondtos.ResponseTransferenciaAutomatica
	//Si no hay ningun cliente habilitado no debe hacer nada
	if len(clientes) > 0 {
		//TODO Aquí se podría usar una go routina para buscar el motivo y la cuenta de telco.
		//Busco el motivo por defecto para las transferencias
		filtroMotivo := filtros.ConfiguracionFiltro{
			Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
		}

		motivo, erro := s.utilService.GetConfiguracionService(filtroMotivo)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		// Busco la cuenta de telco para las transferencias

		filtroCuenta := filtros.ConfiguracionFiltro{
			Nombre: "CBU_CUENTA_TELCO",
		}

		cbu, erro := s.utilService.GetConfiguracionService(filtroCuenta)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		if len(cbu.Valor) < 1 {
			erro = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
			return administraciondtos.RequestMovimientosId{}, erro
		}

		var listaTransferencias []administraciondtos.RequestTransferenciaAutomatica

		for _, c := range clientes {
			for _, cu := range *c.Cuentas {
				filtroMovimiento := filtros.MovimientoFiltro{
					AcumularPorPagoIntentos: true,
					CuentaId:                uint64(cu.ID),
					CargarPago:              true,
					CargarPagoEstados:       true,
					CargarPagoIntentos:      true,
					CargarMedioPago:         true,
					// FechaInicio:             fechaI.AddDate(0, 0, int(-cu.DiasRetiroAutomatico)),
					// FechaFin:                fechaF,
					CargarMovimientosNegativos: true,
				}
				// filtroMovimiento.CuentaId = uint64(cu.ID)
				movimientos, err := s.GetMovimientosAcumulados(filtroMovimiento)
				if err != nil {
					erro = err
					return administraciondtos.RequestMovimientosId{}, erro
				}
				if len(movimientos.Acumulados) > 0 {
					var listaIdsMovimiento []uint64
					var listaMovimientosIdNeg []uint64
					var acumulado entities.Monto
					var acumuladoNeg entities.Monto
					for _, ma := range movimientos.Acumulados {
						acumulado += ma.Acumulado
						for _, m := range ma.Movimientos {
							listaIdsMovimiento = append(listaIdsMovimiento, uint64(m.Id))
						}
					}
					for _, mn := range movimientos.MovimientosNegativos {
						acumuladoNeg += mn.Monto
						listaMovimientosIdNeg = append(listaMovimientosIdNeg, uint64(mn.Id))
					}
					if acumuladoNeg < 0 {
						acumulado += acumuladoNeg
					}
					request := administraciondtos.RequestTransferenciaAutomatica{
						CuentaId: uint64(cu.ID),
						Cuenta:   cu.Cuenta,
						DatosClientes: administraciondtos.DatosClientes{
							NombreCliente: c.Cliente,
							EmailCliente:  c.Email,
						},
						Request: administraciondtos.RequestTransferenicaCliente{
							Transferencia: linktransferencia.RequestTransferenciaCreateLink{
								Origen: linktransferencia.OrigenTransferenciaLink{
									Cbu: cbu.Valor,
								},
								Destino: linktransferencia.DestinoTransferenciaLink{
									Cbu:            cu.Cbu,
									EsMismoTitular: false,
								},
								Importe: acumulado,
								Moneda:  linkdtos.Pesos,
								Motivo:  linkdtos.EnumMotivoTransferencia(motivo.Valor),
							},
							ListaMovimientosId:    listaIdsMovimiento,
							ListaMovimientosIdNeg: listaMovimientosIdNeg,
						},
					}

					listaTransferencias = append(listaTransferencias, request)
				}

			} // fin for range cuentas
		} // fin for range clientes

		var idmovcomisiones administraciondtos.RequestMovimientosId
		for _, t := range listaTransferencias {
			if t.Request.Transferencia.Importe > 0 {
				uuid := uuid.NewV4()
				response, erro := s.BuildTransferencia(ctx, uuid.String(), t.Request, t.CuentaId, t.DatosClientes)
				if erro != nil {
					aviso := entities.Notificacione{
						Tipo:        entities.NotificacionTransferenciaAutomatica,
						Descripcion: fmt.Sprintf("atención no se pudo realizar transferencia automatica de la cuenta: %s", t.Cuenta),
						UserId:      0,
					}
					erro := s.utilService.CreateNotificacionService(aviso)
					if erro != nil {
						logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
					}
				} else {
					// acumular id de movimientos transferidos
					idmovcomisiones.MovimientosId = append(idmovcomisiones.MovimientosId, response.MovimientosIdTransferidos...)
					//idmovcomisiones.MovimimientosIdRevertidos = append(idmovcomisiones.MovimimientosIdRevertidos, response.MovimientosIdReversiones...)
				}
			}
		}
		movimientosidcomision = idmovcomisiones
	}
	return
}

func (s *service) CreateCuentaApilinkService(request linkcuentas.LinkPostCuenta) (erro error) {
	return s.apilinkService.CreateCuentaApiLinkService(request)
}

func (s *service) DeleteCuentaApilinkService(request linkcuentas.LinkDeleteCuenta) (erro error) {
	return s.apilinkService.DeleteCuentaApiLinkService(request)
}

func (s *service) GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error) {
	return s.apilinkService.GetCuentasApiLinkService()
}

func (s *service) CreatePlanCuotasService(request administraciondtos.RequestPlanCuotas) (erro error) {

	/* cnovertir datos para ser procesados */
	installmentId, fechaVigencia, err := procesarRequest(request.InstalmentsId, request.VigenciaDesde)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}

	installmentActual, err := s.repository.GetInstallmentById(uint(installmentId))
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	// fmt.Printf("%v", installmentActual.VigenciaDesde.After(fechaVigencia))
	// if installmentActual.VigenciaDesde.After(fechaVigencia) || installmentActual.VigenciaDesde.Equal(fechaVigencia) {}
	if installmentActual.VigenciaHasta != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	fechaHasta := fechaVigencia.Add(time.Hour * 24 * -1)
	installmentActual.VigenciaHasta = &fechaHasta
	fmt.Printf("%v - %v \n", fechaVigencia, fechaHasta)

	installmentNew := entities.Installment{
		MediopagoinstallmentsID: installmentActual.MediopagoinstallmentsID,
		Descripcion:             installmentActual.Descripcion,
		Issuer:                  installmentActual.Issuer,
		VigenciaDesde:           fechaVigencia,
	}

	openFile, err := os.Open(request.RutaFile)
	if err != nil {
		erro = errors.New(ERROR_LEER_ARCHIVO)
		return
	}
	readFile := csv.NewReader(openFile)
	readFile.Comma = ';'
	readFile.FieldsPerRecord = -1
	flag := false
	var listPlanescuotas []entities.Installmentdetail
	for {
		registro, err := readFile.Read()
		if err != nil && err != io.EOF {
			erro = errors.New(ERROR_LEER_ARCHIVO)
			return
		}
		if err == io.EOF {
			break
		}
		if !flag {
			flag = true
			listPlanescuotas = append(listPlanescuotas, entities.Installmentdetail{
				InstallmentsID: 0,
				Activo:         false,
				Cuota:          1,
				Tna:            0,
				Tem:            0,
				Coeficiente:    1,
				Fechadesde:     fechaVigencia,
			})
			continue
		}
		cuota, tna, tem, coeficiente, err := procesarRegistro(registro)
		if err != nil {
			erro = err
			return
		}
		listPlanescuotas = append(listPlanescuotas, entities.Installmentdetail{
			InstallmentsID: 0,
			Activo:         false,
			Cuota:          cuota,
			Tna:            tna,
			Tem:            tem,
			Coeficiente:    coeficiente,
			Fechadesde:     fechaVigencia,
		})

	}
	openFile.Close()
	err = s.repository.CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew, listPlanescuotas)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       err.Error(),
			Funcionalidad: "CreatePlanCuotasByInstallmenIdRepository - repository",
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error(err)
		}
		return
	}

	err = s.commonsService.BorrarDirectorio(request.RutaFile)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       err.Error(),
			Funcionalidad: "BorrarDirectorio - commonService",
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error(err)
		}
		return
	}
	return
}

func procesarRegistro(planCuota []string) (cuota int64, tna, tem, coeficiente float64, erro error) {

	cuota, err := strconv.ParseInt(planCuota[0], 10, 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	tna, err = strconv.ParseFloat(planCuota[1][0:len(planCuota[1])-1], 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	tem, err = strconv.ParseFloat(planCuota[2][0:len(planCuota[2])-1], 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	coeficiente, err = strconv.ParseFloat(planCuota[3], 64)

	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	return
}

func procesarRequest(installmentid, fecha string) (idInstallments int, fechaVigencia time.Time, erro error) {
	idInstallments, err := strconv.Atoi(installmentid)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	fechaVigencia, err = time.Parse("2006-01-02", fecha)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	return
}

// func formatFecha() (fechaI time.Time, fechaF time.Time, erro error) {
// startTime := time.Now()
// fechaConvert := startTime.Format("2006-01-02") //YYYY.MM.DD
// fec := strings.Split(fechaConvert, "-")
//
// dia, err := strconv.Atoi(fec[len(fec)-1])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// mes, err := strconv.Atoi(fec[1])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// anio, err := strconv.Atoi(fec[0])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// fechaI = time.Date(anio, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
// fechaF = time.Date(anio, time.Month(mes), dia, 23, 59, 59, 0, time.UTC)
//
// return
// }

func (s *service) GetPeticionesService(filtro filtros.PeticionWebServiceFiltro) (peticiones administraciondtos.ResponsePeticionesWebServices, erro error) {

	peticionesRes, total, erro := s.repository.GetPeticionesWebServices(filtro)
	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		peticiones.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, p := range peticionesRes {
		peticiones.Peticiones = append(peticiones.Peticiones, administraciondtos.ResponsePeticionWebServices{
			Operacion: p.Operacion,
			Vendor:    string(p.Vendor),
		})
	}

	return peticiones, nil
}

func (s *service) GetPagosTipoChannelService(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error) {
	return s.repository.GetPagosTipoChannelRepository(filtro)
}

func (s *service) DeletePagoTipoChannelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	erro = s.repository.DeletePagoTipoChannel(id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	if erro != nil {
		return erro
	}

	return
}

func (s *service) CreatePagoTipoChannel(ctx context.Context, request administraciondtos.RequestPagoTipoChannel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.PagoTipoFiltro{
		Id: request.PagoTipoId,
	}

	_, erro = s.repository.GetPagoTipo(filtro)

	if erro != nil {
		return
	}

	filtro2 := filtros.ChannelFiltro{
		Id: request.ChannelId,
	}

	_, erro = s.repository.GetChannel(filtro2)

	if erro != nil {
		return
	}

	filtro3 := filtros.PagoTipoChannelFiltro{
		PagoTipoId: request.PagoTipoId,
		ChannelId:  request.ChannelId,
	}

	res, erro := s.repository.GetPagosTipoChannelRepository(filtro3)

	if erro != nil {
		return
	}
	if len(res) > 0 {
		return uint64(res[0].ID), nil
	}

	pagotipochannel := request.ToPagoTipoChannel(false)

	return s.repository.CreatePagoTipoChannel(ctx, pagotipochannel)

}

func (s *service) SubirArchivos(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse) (countArchivo int, erro error) {
	/*
		por ultimo se mueven los archvios del directorio temporal
		a un directorio en minio dondo se almacenan los archivos
		de cierre de lote registrado en la DB
	*/
	var rutaDestino string
	for _, archivo := range listaArchivo {
		/*
			se lee el contenido del archivo y se obtiene su contenido se le pasa:
			- ruta destino
			- ruta origen del archivo
			- nombre del archivo
		*/
		data, filename, filetypo, err := util.LeerDatosArchivo(rutaDestino, rutaArchivos, archivo.NombreArchivo)
		filename = config.DIR_KEY + filename
		if err != nil {
			logs.Error(err)
		}
		/*	necesito la data, nombre del archivo y el tipo */
		filenameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
		erro := s.store.PutObject(ctx, data, filenameWithoutExt, filetypo)
		if erro != nil {
			logs.Error("No se pudo guardar el archivo")
		}

	}
	nombreDirectorio := config.DIR_KEY

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
				return 0, erro
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
				return 0, erro
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
			return 0, erro
		}
	}

	return 1, erro
}

func (s *service) GetChannelsArancelService(filtro filtros.ChannelArancelFiltro) (response administraciondtos.ResponseChannelsArancel, erro error) {

	channelaranc, total, erro := s.repository.GetChannelsAranceles(filtro)

	logs.Info(channelaranc)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, ca := range channelaranc {

		r := administraciondtos.ResponseChannelsAranceles{}
		r.FromChannelArancel(ca)

		response.ChannelArancel = append(response.ChannelArancel, r)
	}

	return
}

func (s *service) CreateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return 0, erro
	}

	channelsArancel := request.ToChannelsArancel(false)

	return s.repository.CreateChannelsArancel(ctx, channelsArancel)

}

func (s *service) UpdateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (erro error) {
	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return erro
	}

	channelArancelModificada := request.ToChannelsArancel(true)

	return s.repository.UpdateChannelsArancel(ctx, channelArancelModificada)

}

func (s *service) DeleteChannelsArancelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id de channel arancel es invalido")
		return
	}

	return s.repository.DeleteChannelsArancel(ctx, id)

}

func (s *service) GetChannelArancelService(filtro filtros.ChannelAranceFiltro) (response administraciondtos.ResponseChannelsAranceles, erro error) {

	channel_arance, erro := s.repository.GetChannelArancel(filtro)

	if erro != nil {
		return
	}

	response.FromChArancel(channel_arance)

	return
}

func (s *service) ObtenerArchivosSubidos(filtro filtros.Paginacion) (lisArchivosSubidos administraciondtos.ResponseArchivoSubido, erro error) {
	var contador int64
	var recorrerHasta int32
	var listaTemporalArchivo []administraciondtos.ArchivoSubido
	entityCl, err := s.repository.GetCierreLoteSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityCl) > 0 {
		for _, valueCL := range entityCl {
			contador++
			var listaClTemporal administraciondtos.ArchivoSubido
			listaClTemporal.EntityClToDtos(&valueCL)
			listaTemporalArchivo = append(listaTemporalArchivo, listaClTemporal)
		}
	}

	entityPx, err := s.repository.GetPrismaPxSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityPx) > 0 {
		for _, valuePx := range entityPx {
			contador++
			var listaPxTemporal administraciondtos.ArchivoSubido
			listaPxTemporal.EntityPxToDtos(&valuePx)
			listaTemporalArchivo = append(listaTemporalArchivo, listaPxTemporal)
		}
	}

	entityMx, err := s.repository.GetPrismaMxSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityMx) > 0 {
		for _, valueMx := range entityMx {
			contador++
			var listaMxTemporal administraciondtos.ArchivoSubido
			listaMxTemporal.EntityMxToDtos(&valueMx)
			listaTemporalArchivo = append(listaTemporalArchivo, listaMxTemporal)
		}
	}

	sort.Slice(listaTemporalArchivo, func(i, j int) bool {
		return listaTemporalArchivo[i].FechaSubida.Before(listaTemporalArchivo[j].FechaSubida)
	})

	if filtro.Number > 0 && filtro.Size > 0 {
		lisArchivosSubidos.Meta = _setPaginacion(filtro.Number, filtro.Size, contador)
	}
	recorrerHasta = lisArchivosSubidos.Meta.Page.To
	if lisArchivosSubidos.Meta.Page.CurrentPage == lisArchivosSubidos.Meta.Page.LastPage {
		recorrerHasta = lisArchivosSubidos.Meta.Page.Total
	}
	if len(listaTemporalArchivo) > 0 {
		for i := lisArchivosSubidos.Meta.Page.From; i < recorrerHasta; i++ {
			lisArchivosSubidos.ArchivosSubidos = append(lisArchivosSubidos.ArchivosSubidos, listaTemporalArchivo[i])
		}
	}
	return
}

func (s *service) ObtenerArchivoCierreLoteRapipago(nombre string) (archivo bool, err error) {

	archivo, err = s.repository.ObtenerArchivoCierreLoteRapipago(nombre)

	if err != nil {
		return
	}

	return
}

func (s *service) GetCierreLoteEnDisputaServices(estadoDisputa int, request filtros.ContraCargoEnDisputa) (cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL, erro error) {

	var clTemporal cierrelotedtos.ResponsePrismaCL
	operacionesEnDisputa, err := s.repository.ObtenerCierreLoteEnDisputaRepository(estadoDisputa, request)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	operacionesContracargo, err := s.repository.ObtenerCierreLoteContraCargoRepository(estadoDisputa, request)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}

	for _, valueCl := range operacionesEnDisputa {
		clTemporal.EntityToDtos(valueCl)
		cierreLoteDisputa = append(cierreLoteDisputa, clTemporal)
	}

	for _, valueCl := range operacionesContracargo {
		clTemporal.EntityToDtos(valueCl)
		cierreLoteDisputa = append(cierreLoteDisputa, clTemporal)
	}

	return
}

func (s *service) GetPagosByTransactionIdsServices(filtro filtros.ContraCargoEnDisputa, cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL) (listaRevertidos administraciondtos.ResponseOperacionesContracargo, erro error) {
	for _, value := range cierreLoteDisputa {
		filtro.TransactionId = append(filtro.TransactionId, value.ExternalclienteID)
	}
	// obtener cuenta y pago tipos

	// obtengo pagos por pago tipo id

	listaPagosRevertidos, err := s.repository.ObtenerPagosInDisputaRepository(filtro)
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(listaPagosRevertidos) > 0 {
		var pagostemporal []administraciondtos.ResponsePagoCC
		var pagosTipoTemporal []administraciondtos.ResponsePagotipoCC

		listaRevertidos.Cuenta.Id = listaPagosRevertidos[0].Pago.PagosTipo.Cuenta.ID
		listaRevertidos.Cuenta.Cuenta = listaPagosRevertidos[0].Pago.PagosTipo.Cuenta.Cuenta
		var pagoTipoId uint
		for _, value := range listaPagosRevertidos {
			if value.Pago.PagosTipo.ID != pagoTipoId {
				pagoTipoId = value.Pago.PagosTipo.ID
				pagosTipoTemporal = append(pagosTipoTemporal, administraciondtos.ResponsePagotipoCC{
					Id:       value.Pago.PagosTipo.ID,
					Pagotipo: value.Pago.PagosTipo.Pagotipo,
				})

			}
		}
		for _, valuePago := range listaPagosRevertidos {
			pagostemporal = append(pagostemporal, administraciondtos.ResponsePagoCC{
				Id:                  valuePago.Pago.ID,
				PagostipoID:         valuePago.Pago.PagostipoID,
				Fecha:               valuePago.Pago.FirstDueDate,
				ExternalReference:   valuePago.Pago.ExternalReference,
				PayerName:           strings.ToUpper(valuePago.Pago.PayerName),
				Estado:              "",
				NombreEstado:        "",
				Amount:              valuePago.Pago.FirstTotal,
				FechaPago:           valuePago.Pago.CreatedAt,
				Channel:             "",
				NombreChannel:       "",
				UltimoPagoIntentoId: uint64(valuePago.ID),
				TransferenciaId:     0,
				ReferenciaBancaria:  "",
				PagoIntento: administraciondtos.ResponsePagoIntentoCC{
					Id:                   valuePago.ID,
					MediopagosId:         uint(valuePago.MediopagosID),
					InstallmentdetailsId: uint(valuePago.InstallmentdetailsID),
					ExternalId:           "",
					PaidAt:               valuePago.PaidAt,
					ReortAt:              valuePago.ReportAt,
					IsAvailable:          valuePago.IsAvailable,
					Amount:               valuePago.Amount,
					Valorcupon:           valuePago.Valorcupon,
					StateComent:          valuePago.StateComment,
					Barcode:              "",
					BarcodeUrl:           "",
					AvailableAt:          valuePago.AvailableAt,
					RevertedAt:           valuePago.RevertedAt,
					HolderName:           strings.ToUpper(valuePago.HolderName),
					HolderEmail:          valuePago.HolderEmail,
					HolderType:           valuePago.HolderType,
					HolderNumber:         valuePago.HolderNumber,
					HolderCbu:            "",
					TicketNumber:         valuePago.TicketNumber,
					AuthorizationCode:    valuePago.AuthorizationCode,
					CardLastFourDigits:   valuePago.CardLastFourDigits,
					TransactionId:        valuePago.TransactionID,
					SiteId:               "",
				},
			})
		}

		for Key, valuePagoTipoTemp := range pagosTipoTemporal {
			for _, valuePagoTemp := range pagostemporal {
				if valuePagoTipoTemp.Id == uint(valuePagoTemp.PagostipoID) {
					pagosTipoTemporal[Key].Pagos = append(pagosTipoTemporal[Key].Pagos, valuePagoTemp)
				}
			}
		}
		listaRevertidos.Cuenta.PagoTipo = append(listaRevertidos.Cuenta.PagoTipo, pagosTipoTemporal...)
		for _, valueCL := range cierreLoteDisputa {
			for key1, valuePI := range listaRevertidos.Cuenta.PagoTipo {
				for key, valuePago := range valuePI.Pagos {
					if valueCL.ExternalclienteID == valuePago.PagoIntento.TransactionId {
						// listaRevertidos.Cuenta.PagoTipo[key1].Pagos[key].PagoIntento.CierreLote = administraciondtos.ResponseCLCC(valueCL)
						listaRevertidos.Cuenta.PagoTipo[key1].Pagos[key].PagoIntento.CierreLote = administraciondtos.ResponseCLCC{
							Id:                         valueCL.Id,
							PagoestadoexternosId:       valueCL.PagoestadoexternosId,
							ChannelarancelesId:         valueCL.ChannelarancelesId,
							ImpuestosId:                valueCL.ImpuestosId,
							PrismamovimientodetallesId: valueCL.PrismamovimientodetallesId,
							PrismamovimientodetalleId:  0,
							PrismatrdospagosId:         valueCL.PrismatrdospagosId,
							BancoExternalId:            valueCL.BancoExternalId,
							Tiporegistro:               valueCL.Tiporegistro,
							PagosUuid:                  valueCL.PagosUuid,
							ExternalmediopagoId:        valueCL.ExternalmediopagoId,
							Nrotarjeta:                 valueCL.Nrotarjeta,
							Tipooperacion:              valueCL.Tipooperacion,
							Fechaoperacion:             valueCL.Fechaoperacion,
							Monto:                      valueCL.Monto,
							Montofinal:                 valueCL.Montofinal,
							Codigoautorizacion:         valueCL.Codigoautorizacion,
							Nroticket:                  valueCL.Nroticket,
							SiteID:                     valueCL.SiteID,
							ExternalloteId:             valueCL.ExternalloteId,
							Nrocuota:                   valueCL.Nrocuota,
							FechaCierre:                valueCL.FechaCierre,
							Nroestablecimiento:         valueCL.Nroestablecimiento,
							ExternalclienteID:          valueCL.ExternalclienteID,
							Nombrearchivolote:          valueCL.Nombrearchivolote,
							Match:                      valueCL.Match,
							FechaPago:                  valueCL.FechaPago,
							Disputa:                    valueCL.Disputa,
							Reversion:                  valueCL.Reversion,
							DetallemovimientoId:        valueCL.DetallemovimientoId,
							DetallepagoId:              valueCL.DetallepagoId,
							Descripcioncontracargo:     valueCL.Descripcioncontracargo,
							ExtbancoreversionId:        valueCL.ExtbancoreversionId,
							Conciliado:                 valueCL.Conciliado,
							Estadomovimiento:           valueCL.Estadomovimiento,
							Descripcionbanco:           valueCL.Descripcionbanco,
						}
					}
				}
			}
		}
	}

	return
}

func (s *service) PostPreferencesService(request administraciondtos.RequestPreferences) (erro error) {

	// validamos los datos
	erro = request.Validar()
	if erro != nil {
		return
	}
	// abrir archivo , leer converir a array de bute y guardar

	buffer, err := request.File.Open()
	if err != nil {
		erro = err
		return
	}
	defer buffer.Close() // al finalizar cerrar el archivo

	// se crea contexto
	ctx := context.Background()

	// lerr archivo
	data, erro := ioutil.ReadAll(buffer) // leer el contenido del archivo
	if erro != nil {
		msj := "error a leer datos del archivo:" + request.File.Filename
		logs.Error(msj)
		erro = errors.New(msj)
		return
	}

	archivo_extension := strings.Split(request.File.Filename, ".") // divide en un array de nombre y extension
	archivonombre := fmt.Sprintf("%s/%s", request.RutaLogo, archivo_extension[0])
	archivotipo := archivo_extension[len(archivo_extension)-1]

	// Se guarda en S3
	//data = contenido []byte ,archivonombre= carpeta y nombre del archivo donde se guarda ,archivotipo= extension

	erro = s.store.PutObject(ctx, data, archivonombre, archivotipo)
	if erro != nil {
		logs.Error("No se pudo guardar el archivo")
		return
	} else {
		// en el caso de no ocurrir error guardar en base de datos
		clienteId, err := strconv.Atoi(request.ClientId)
		if err != nil {
			erro = errors.New("cliente id no válido")
			return
		}
		ruta := fmt.Sprintf("%s/%s", request.RutaLogo, request.File.Filename)
		preferenceEntity := entities.Preference{
			ClientesId:     uint(clienteId),
			Maincolor:      request.MainColor,
			Secondarycolor: request.SecondaryColor,
			Logo:           ruta, // se debe guardar la carpeta y nombre del archivo
		}

		err = s.repository.PostPreferencesRepository(preferenceEntity)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
	}

	// defer request.File.

	return
}
func (s *service) GetPreferencesService(request administraciondtos.RequestPreferences) (responsePreference dtos.ResponsePreference, erro error) {
	clienteId, err := strconv.Atoi(request.ClientId)
	if err != nil || clienteId < 1 {
		erro = errors.New("debe enviar un cliente_id válido")
		return
	}
	clienteEntity := entities.Cliente{
		Model: gorm.Model{
			ID: uint(clienteId),
		},
	}
	preference, err := s.repository.GetPreferencesRepository(clienteEntity)
	if err != nil {
		erro = fmt.Errorf("error %s", err)
		return
	}
	var logo []byte
	if len(preference.Logo) > 0 {
		ctx := context.Background()
		logo, err = s.store.GetObjectS3(ctx, preference.Logo)
		if err != nil {
			preference.Logo = ""
			logs.Info("no se pudo recuperar logo de S3")
		}
		preference.Logo = base64.StdEncoding.EncodeToString(logo)
	}
	responsePreference = dtos.ResponsePreference{
		Client:         fmt.Sprint(preference.ClientesId),
		MainColor:      preference.Maincolor,
		SecondaryColor: preference.Secondarycolor,
		Logo:           preference.Logo,
	}
	return

}

func (s *service) DeletePreferencesService(request administraciondtos.RequestPreferences) (erro error) {
	clienteId, erro := strconv.Atoi(request.ClientId)
	if erro != nil || clienteId < 1 {
		erro = errors.New("error:debe enviar un cliente_id válido")
		return
	}
	clienteEntity := entities.Cliente{
		Model: gorm.Model{
			ID: uint(clienteId),
		},
	}
	erro = s.repository.DeletePreferencesRepository(clienteEntity)
	if erro != nil {
		return fmt.Errorf("error %s", erro.Error())
	}
	return nil
}

func (s *service) GetPagosDevService(filtro filtros.PagoFiltro) (response []entities.Pago, erro error) {

	response, _, erro = s.repository.GetPagos(filtro)

	if erro != nil {
		return
	}

	return

}

func (s *service) UpdatePagosDevService(response []entities.Pago) (pg []uint, erro error) {
	// return s.repository.UpdatePagosNotificados(listaPagosNotificar)

	for _, pago := range response {
		pg = append(pg, pago.ID)
	}

	erro = s.repository.UpdatePagosDev(pg)
	if erro != nil {
		return nil, erro
	}

	return pg, nil

}

func (s *service) BuildPagosMovDev(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

	// deben existir pagos
	if len(pagos) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Channel:                 true,
		CargarPago:              true,
		CargarPagoTipo:          true,
		CargarCuenta:            true,
		CargarCliente:           true,
		CargarCuentaComision:    true,
		CargarImpuestos:         true,
		ExternalId:              true,
		CargarInstallmentdetail: true,
		PagosId:                 pagos,
	}
	// * 4 - Busco los pagos intentos que corresponden a los pagos
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	// * 5 Obtener los pagos intentos
	for i := range pagosIntentos {
		movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, pagosIntentos[i])
	}

	// * 6 - Busco el estado acreditado
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, erro := s.repository.GetPagoEstado(filtroPagoEstado)
	logs.Info(pagoEstadoAcreditado)

	if erro != nil {
		return
	}

	// var monto_pagado entities.Monto
	// * 8 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range movimientoCierreLote.ListaPagoIntentos {
		/* * para el calculo de la comision fitrar por el id del channel y el id de la cuentar*/

		var pagoCuotas bool
		var examinarPagoCuota bool
		if movimientoCierreLote.ListaPagoIntentos[i].Installmentdetail.Cuota > 1 {
			pagoCuotas = true
			examinarPagoCuota = true
		}
		var idMedioPago uint
		if movimientoCierreLote.ListaPagoIntentos[i].MediopagosID == 30 {
			idMedioPago = uint(movimientoCierreLote.ListaPagoIntentos[i].MediopagosID)
			pagoCuotas = true
			examinarPagoCuota = true
		}

		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         uint(movimientoCierreLote.ListaPagoIntentos[i].Mediopagos.ChannelsID),
			CuentaId:          movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.ID,
			Mediopagoid:       idMedioPago,
			ExaminarPagoCuota: examinarPagoCuota,
			PagoCuota:         pagoCuotas,
			Channelarancel:    true,
			FechaPagoVigencia: movimientoCierreLote.ListaPagoIntentos[i].PaidAt,
		}

		logs.Info(filtroComisionChannel)

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}

		movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, movimientoCierreLote.ListaPagoIntentos[i].Pago)

		// * crear el log de estado de pago acreditado
		pagoEstadoLog := entities.Pagoestadologs{
			PagosID:       movimientoCierreLote.ListaPagoIntentos[i].PagosID,
			PagoestadosID: int64(pagoEstadoAcreditado.ID),
		}
		movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)

		if movimientoCierreLote.ListaPagoIntentos[i].Pago.PagoestadosID == int64(pagoEstadoAcreditado.ID) {

			var importe entities.Monto
			importe = movimientoCierreLote.ListaPagoIntentos[i].Amount
			if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon != 0 {
				importe = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
			}

			movimiento := entities.Movimiento{}
			// monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Amount
			// if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon > 0 {
			// 	monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
			// }
			movimiento.AddCredito(uint64(movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.CuentasID), uint64(movimientoCierreLote.ListaPagoIntentos[i].ID), importe)

			s.utilService.BuildComisiones(&movimiento, &listaCuentaComision, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, importe)

			movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
			movimientoCierreLote.ListaPagoIntentos[i].AvailableAt = movimientoCierreLote.ListaPagoIntentos[i].CreatedAt
		}

	}
	return
}

// ? Implementacion de servicio conusltar cierrelote para herramienta wee

func (s *service) GetConsultarClRapipagoService(filtro filtros.RequestClrapipago) (response administraciondtos.ResponseCLRapipago, erro error) {

	clrapiapgo, total, erro := s.repository.GetConsultarClRapipagoRepository(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cl := range clrapiapgo {

		var detalles []administraciondtos.ClRapipagoDetalle
		for _, detalle := range cl.RapipagoDetalle {
			detalles = append(detalles, administraciondtos.ClRapipagoDetalle{
				FechaCobro:       detalle.FechaCobro,
				ImporteCobrado:   uint64(detalle.ImporteCobrado),
				ImporteCalculado: float64(detalle.ImporteCalculado),
				CodigoBarras:     detalle.CodigoBarras,
				Conciliado:       detalle.Match,
				Informado:        detalle.Pagoinformado,
			})
		}

		fechaProceso := s.commonsService.ConvertirFormatoFecha(cl.FechaProceso)
		r := administraciondtos.CLRapipago{
			IdClRapipago:             uint64(cl.ID),
			IdArchivo:                cl.NombreArchivo,
			FechaProceso:             fechaProceso,
			Detalles:                 uint64(cl.CantDetalles),
			ImporteTotal:             uint64(cl.ImporteTotal),
			ImporteTotalCalculado:    float64(cl.ImporteTotalCalculado),
			IdBanco:                  uint64(cl.BancoExternalId),
			FechaAcreditacion:        cl.Fechaacreditacion.Format("2006-01-02"),
			CantidadDiasAcreditacion: uint64(cl.Cantdias),
			ImporteMinimo:            uint64(cl.ImporteMinimo),
			Coeficiente:              cl.Coeficiente,
			EnObservacion:            cl.Enobservacion,
			DiferenciaBanco:          s.utilService.ToFixed(cl.Difbancocl, 2) / 100,
			FechaCreacion:            cl.CreatedAt.Format("2006-01-02"),
			PagoActualizado:          cl.PagoActualizado,
			ClRapipagoDetalle:        detalles,
		}

		response.ClRapipago = append(response.ClRapipago, r)
	}

	return
}

func (s *service) GetConsultarClMultipagoService(filtro filtros.RequestClMultipago) (response administraciondtos.ResponseClMultipago, erro error) {

	clmultipago, total, erro := s.repository.GetConsultarClMultipagoRepository(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cl := range clmultipago {

		var detalles []administraciondtos.ClMultipagoDetalle
		for _, detalle := range cl.MultipagoDetalle {
			detalles = append(detalles, administraciondtos.ClMultipagoDetalle{
				FechaCobro:       detalle.FechaCobro,
				ImporteCobrado:   uint64(detalle.ImporteCobrado),
				ImporteCalculado: float64(detalle.ImporteCalculado),
				CodigoBarras:     detalle.CodigoBarras,
				Conciliado:       detalle.Match,
				Informado:        detalle.Pagoinformado,
			})
		}

		fechaProceso := s.commonsService.ConvertirFormatoFecha(cl.FechaProceso)
		r := administraciondtos.ClMultipago{
			IdClMultipago:            uint64(cl.ID),
			IdArchivo:                cl.NombreArchivo,
			FechaProceso:             fechaProceso,
			Detalles:                 uint64(cl.CantDetalles),
			ImporteTotal:             uint64(cl.ImporteTotal),
			ImporteTotalCalculado:    float64(cl.ImporteTotalCalculado),
			IdBanco:                  uint64(cl.BancoExternalId),
			FechaAcreditacion:        cl.Fechaacreditacion.Format("2006-01-02"),
			CantidadDiasAcreditacion: uint64(cl.Cantdias),
			ImporteMinimo:            uint64(cl.ImporteMinimo),
			Coeficiente:              cl.Coeficiente,
			EnObservacion:            cl.Enobservacion,
			DiferenciaBanco:          s.utilService.ToFixed(cl.Difbancocl, 2) / 100,
			FechaCreacion:            cl.CreatedAt.Format("2006-01-02"),
			PagoActualizado:          cl.PagoActualizado,
			ClMultipagoDetalle:       detalles,
		}

		response.ClMultipago = append(response.ClMultipago, r)
	}

	return
}

func (s *service) GetCaducarOfflineIntentos() (intentosCaducados int, erro error) {

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		CargarPago: true,
		// CargarPagoEstado:   true,
		// Channel:            true,
		PagoEstadoIdFiltro: 2, // Estado Processing
		ChannelIdFiltro:    3, // Channel Offline
	}
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	var pagosActualizar []entities.Pago
	const idEstadoExpired = 6 // Estado Expired

	for _, pagoIntento := range pagosIntentos {
		lastFechaVencimiento := pagoIntento.Pago.FirstDueDate
		if pagoIntento.Pago.SecondDueDate.After(pagoIntento.Pago.FirstDueDate) {
			lastFechaVencimiento = pagoIntento.Pago.SecondDueDate
		}

		fechaControl := time.Now().AddDate(0, 0, -5)
		if fechaControl.After(lastFechaVencimiento) {
			pagosActualizar = append(pagosActualizar, pagoIntento.Pago)
		}

	}

	if len(pagosActualizar) > 0 {
		erro = s.repository.UpdateEstadoPagos(pagosActualizar, uint64(idEstadoExpired))
		if erro != nil {
			return
		}
	}

	return len(pagosActualizar), nil
}

func (s *service) GetCaducarPagosExpirados(filtro filtros.PagoCaducadoFiltro) (intentosCaducados int, erro error) {

	if filtro.PagoTipo != "" && filtro.PagosTipoId == 0 {
		filtroPagoTipo := filtros.PagoTipoFiltro{
			PagoTipo: filtro.PagoTipo,
		}
		pagoTipo, err := s.repository.GetPagoTipo(filtroPagoTipo)
		if err != nil {
			erro = err
			return
		}

		filtro.PagosTipoId = uint64(pagoTipo.ID)
	}

	if filtro.CuentaApikey != "" && filtro.CuentaId == 0 {

		cuenta, err := s.repository.GetCuentaByApiKey(filtro.CuentaApikey)
		if err != nil {
			erro = err
			return
		}

		filtro.CuentaId = uint64(cuenta.ID)
	}

	pagoEstadosBuscados := []uint64{1, 2}
	filtroPago := filtros.PagoFiltro{
		CargarCuenta:      true,
		CuentaId:          filtro.CuentaId,
		PagosTipoId:       filtro.PagosTipoId,
		PagoEstadosIds:    pagoEstadosBuscados,
		CargaPagoIntentos: true,
		// VisualizarPendientes: true,
	}
	pagos, _, erro := s.repository.GetPagos(filtroPago)
	if erro != nil {
		return
	}

	var pagosActualizar []entities.Pago
	const idEstadoExpired = 6 // Estado Expired
	// const minutosDia = 1440   // Minutos en un dia // Un dia para caducar para tickets

	for _, pago := range pagos {
		lastFechaVencimiento := pago.FirstDueDate
		if pago.SecondDueDate.After(pago.FirstDueDate) {
			lastFechaVencimiento = pago.SecondDueDate
		}

		hoy := time.Now().Local()

		//Codigo para expirar tickets sin tiempos de rapipago.
		// tiempoExpiracion := pago.Expiration

		// if tiempoExpiracion != 0 && tiempoExpiracion != 300 {

		// 	fechaExpiracion := pago.CreatedAt.AddDate(0, 0, 1)
		// 	diferencia := hoy.Sub(fechaExpiracion)
		// 	minutos := diferencia.Minutes() // en float64
		// 	if minutos > float64(minutosDia) {
		// 		pagosActualizar = append(pagosActualizar, pago)
		// 	}

		// } else {

		// Por defecto control con 5 días
		fechaControl := hoy.AddDate(0, 0, -5)
		if fechaControl.After(lastFechaVencimiento) {
			pagosActualizar = append(pagosActualizar, pago)
		} else {
			fechaControlSinPI := hoy.AddDate(0, 0, -1)
			if len(pago.PagoIntentos) == 0 && fechaControlSinPI.After(lastFechaVencimiento) {
				pagosActualizar = append(pagosActualizar, pago)
			}

		}
		// }

	}

	if len(pagosActualizar) > 0 {
		erro = s.repository.UpdateEstadoPagos(pagosActualizar, uint64(idEstadoExpired))
		if erro != nil {
			return
		}
	}

	return len(pagosActualizar), nil
}

func (s *service) GetPagosCalculoMovTemporalesService(filtro filtros.PagoIntentoFiltros) (pagosid []uint, erro error) {

	if filtro.FechaPagoInicio.IsZero() {
		var fechaI time.Time
		var fechaF time.Time
		// si los filtros recibidos son ceros toman la fecha actual
		fechaI, fechaF, erro = s.commonsService.FormatFecha()
		if erro != nil {
			return
		}
		// a las fechas se le restan un dia ya sea por backgraund o endpoint
		filtro.FechaPagoInicio = fechaI.AddDate(0, 0, int(-1))
		filtro.FechaPagoFin = fechaF.AddDate(0, 0, int(-1))

	} else {
		filtro.FechaPagoInicio = filtro.FechaPagoInicio.AddDate(0, 0, int(-1))
		filtro.FechaPagoFin = filtro.FechaPagoFin.AddDate(0, 0, int(-1))
	}
	logs.Info(filtro.FechaPagoInicio)
	logs.Info(filtro.FechaPagoFin)
	pagos, erro := s.repository.GetPagosIntentosCalculoComisionRepository(filtro)

	if erro != nil {
		return
	}

	for _, pg := range pagos {
		pagosid = append(pagosid, uint(pg.PagosID))
	}

	return

}

func (s *service) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {
	pagos, erro = s.repository.GetPagosIntentosCalculoComisionRepository(filtro)
	if erro != nil {
		return
	}
	return
}

func (s *service) BuildPagosCalculoTemporales(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoTemporalesResponse, erro error) {

	// deben existir pagos
	if len(pagos) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Channel:                 true,
		CargarPago:              true,
		CargarPagoTipo:          true,
		CargarCuenta:            true,
		CargarCliente:           true,
		CargarCuentaComision:    true,
		CargarImpuestos:         true,
		CargarInstallmentdetail: true,
		PagoIntentoAprobado:     true,
		PagosId:                 pagos,
	}
	// * 1 - Busco los pagos intentos que corresponden a los pagos y cargar toda informacion necesaria para el calculo de comisiones
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	if len(pagos) != len(pagosIntentos) {
		erro = fmt.Errorf(ERROR_LISTA_PAGOS_INTENTOS)
		return
	}

	// var monto_pagado entities.Monto
	// * 8 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range pagosIntentos {
		/* * para el calculo de la comision fitrar por el id del channel y el id de la cuentar*/

		var pagoCuotas bool
		var examinarPagoCuota bool
		if pagosIntentos[i].Installmentdetail.Cuota > 1 {
			pagoCuotas = true
			examinarPagoCuota = true
		}
		var idMedioPago uint
		if pagosIntentos[i].MediopagosID == 30 {
			idMedioPago = uint(pagosIntentos[i].MediopagosID)
			pagoCuotas = true
			examinarPagoCuota = true
		}

		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         uint(pagosIntentos[i].Mediopagos.ChannelsID),
			CuentaId:          pagosIntentos[i].Pago.PagosTipo.Cuenta.ID,
			Mediopagoid:       idMedioPago,
			ExaminarPagoCuota: examinarPagoCuota,
			PagoCuota:         pagoCuotas,
			Channelarancel:    true,
			FechaPagoVigencia: pagosIntentos[i].PaidAt,
		}

		logs.Info(filtroComisionChannel)

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", pagosIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, pagosIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}

		movimientoCierreLote.ListaPagosCalculado = pagos

		var importe entities.Monto
		importe = pagosIntentos[i].Amount
		if pagosIntentos[i].Valorcupon != 0 {
			importe = pagosIntentos[i].Valorcupon
		}

		movimiento := entities.Movimientotemporale{}
		// monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Amount
		// if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon > 0 {
		// 	monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
		// }
		movimiento.AddCredito(uint64(pagosIntentos[i].Pago.PagosTipo.CuentasID), uint64(pagosIntentos[i].ID), importe)

		s.utilService.BuildComisionesTemporales(&movimiento, &listaCuentaComision, pagosIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, importe)

		movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)

	}
	return
}

func (s *service) ConciliacionPagosReportesService(filtro filtros.PagoFiltro) (valoresNoEncontrados []string, erro error) {

	successPayments, erro := s.repository.GetSuccessPaymentsRepository(filtro)

	if erro != nil {
		erro = fmt.Errorf("error en consultar pagos exitosos para la conciliacion")
		return
	}

	// si no hay resultados en la busqueda de pagos exitosos, no se debe continuar
	if len(successPayments) == 0 {
		erro = fmt.Errorf("no existen pagos exitosos para ser conciliados")
		return
	}

	// necesario convertir fecha para consultar en tabla reportes
	filtro.Fecha[0] = s.commonsService.ConvertirFechaToDDMMYYYY(filtro.Fecha[0])

	reporte, erro := s.repository.GetReportesPagoRepository(filtro)

	if erro != nil {
		erro = fmt.Errorf("error en consultar reportes de pagos enviados para la conciliacion: %s", erro.Error())
		return
	}

	// si no hay resultados en la busqueda de reportes de pagos, no se debe continuar
	if reporte.ID == 0 {
		erro = fmt.Errorf("no existen reportes de pagos para ser conciliados")
		return
	}

	// comparar pagos exitosos con reportes
	valoresNoEncontrados = _conciliarPagosYReportes(successPayments, reporte)

	// si hay valores no encontrados o conciliados, se debe loguear y notificar email
	if len(valoresNoEncontrados) > 0 {

		// hacemos un log de ls valores no encontrados
		erro = fmt.Errorf("algunos pagos no se encontraron en la conciliacion del cliente %s: %v", reporte.Cliente, valoresNoEncontrados)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       erro.Error(),
			Funcionalidad: "ConciliacionPagosReportesService",
		}

		err := s.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro.Error())
			logs.Error(mensaje)
		}

		var email = []string{config.EMAIL_TELCO}
		fafafa := _armarMensajeVariable(valoresNoEncontrados)

		filtro := utildtos.RequestDatosMail{
			Email:            email,
			Asunto:           "Conciliacion Pagos con Reporte de Cobranzas",
			From:             "Wee.ar!",
			Nombre:           "Administrador",
			Mensaje:          "Las siguientes referencias de pagos para el cliente " + reporte.Cliente + " no se conciliaron: " + fafafa,
			CamposReemplazar: valoresNoEncontrados,
			AdjuntarEstado:   false,
			TipoEmail:        "template",
		}
		erro = s.utilService.EnviarMailService(filtro)
		logs.Info(erro)
	} // Fin de if se encuentran valores no conciliados

	/*  Si el cliente es DPEC, controlar los montos */
	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(filtro.CuentaId),
	}
	cuenta, erro := s.repository.GetCuenta(filtroCuenta)
	var montosIguales bool
	if commons.ContainStrings([]string{cuenta.Cliente.Cliente}, "dpec") {
		// conciliar el total de monto de pago con el total cobrado del reporte
		montosIguales, erro = _conciliarByMontos(successPayments, reporte.Totalcobrado)

		if !montosIguales {

			erro = fmt.Errorf("los montos de pagos exitosos no coinciden con el total cobrado reportado")

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       erro.Error(),
				Funcionalidad: "ConciliacionPagosReportesService",
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro.Error())
				logs.Error(mensaje)
			}
		}

	}
	/*  Fin del caso cliente DPEC */

	return
}

// compara las external references de los pagos exitosos de un determiando dia y cuenta, con los reportes de pagos enviados a ese cliente en ese dia
func _conciliarPagosYReportes(pagos []entities.Pago, reporte entities.Reporte) (valoresNoEncontrados []string) {

	var detalleExternalsIds []string
	externalReferences := _filtrarPorExternalReference(pagos)

	for _, detalle := range reporte.Reportedetalle {
		// guardar en un array para comparar de manera inversa los external con los detalles
		detalleExternalsIds = append(detalleExternalsIds, detalle.PagosId)

		// el pago_id del Reportedetalle contiene el valor del external reference de cada pago exitoso informado por el reporte
		if !commons.ContainStrings(externalReferences, detalle.PagosId) {
			valoresNoEncontrados = append(valoresNoEncontrados, detalle.PagosId)
		}

	} // end for detalle := range reporte.Reportedetalle

	for _, er := range externalReferences {

		if !commons.ContainStrings(detalleExternalsIds, er) {
			valoresNoEncontrados = append(valoresNoEncontrados, er)
		}

	} // end for er := range externalReferences

	return
}

func _filtrarPorExternalReference(pagos []entities.Pago) (externalReferences []string) {
	for _, pago := range pagos {
		externalReferences = append(externalReferences, pago.ExternalReference)
	}
	return
}

func _conciliarByMontos(pagos []entities.Pago, monto string) (result bool, erro error) {
	var amount entities.Monto
	// recorrer los pago
	for _, pago := range pagos {
		// recorrer los PagoIntentos
		for _, intento := range pago.PagoIntentos {
			amount += intento.Amount
		}
	} // fin de for de pagos
	// El monto del reporte esta en string, luego para comparar
	montoPagosString := util.Resolve().FormatNum(amount.Float64())

	if montoPagosString == monto {
		result = true
	}
	return
}

// recibe un array de string que se deben ocupar como campos a reemplazar en el mensaje del email
func _armarMensajeVariable(p_arrayString []string) (mensaje string) {
	mensaje = "<br>"
	for i := 0; i < len(p_arrayString); i++ {
		mensaje += fmt.Sprintf("<b>#%d</b><br>", i)
	}

	return
}

func (s *service) AsignarBancoIdRapipagoService(banco_id int64, rapipago_id int64) error {
	rapipago_cl, err := s.repository.ObtenerCierreLoteRapipago(rapipago_id)
	if err != nil {
		return errors.New("error: " + err.Error())
	}

	if rapipago_cl.BancoExternalId != 0 {
		return errors.New("error: cierre lote ya tiene id de mov. banco")
	}

	rapipago_cl.BancoExternalId = banco_id
	var array_rapigagoEntities []*entities.Rapipagocierrelote
	array_rapigagoEntities = append(array_rapigagoEntities, rapipago_cl)

	err = s.repository.UpdateCierreLoteRapipago(array_rapigagoEntities)
	if err != nil {
		return errors.New("error: " + err.Error())
	}

	return nil
}

func (s *service) CreateContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error) {
	//Valido el email
	isValid := commons.IsEmailValid(contactoReporte.ClienteEmail)
	if !isValid {
		return errors.New("debe enviar un correo valido")
	}
	//Verifico que el id sea mayor a 0
	if contactoReporte.ClienteID < 1 {
		return errors.New("debe enviar un cliente_id por valido")
	}

	entityContactos := contactoReporte.DtosToEntity()
	err = s.repository.CreateContactosReportesRepository(entityContactos)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ReadContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (contactosFormato administraciondtos.ResponseGetContactosReportes, err error) {

	if contactoReporte.ClienteID < 1 {
		err = errors.New("debe enviar un cliente_id válido")
		return
	}

	entityContatoReporte := contactoReporte.DtosToEntity()
	contactos, err := s.repository.ReadContactosReportesRepository(entityContatoReporte)
	if err != nil {
		return
	}
	//Recorro los contactos que traje de la base de datos, y le quito el created_at,deleted_at y updated_at
	for _, value := range contactos {
		contacto := administraciondtos.ResponseContactosReportes{
			ClienteEmail: value.Email,
			ClienteID:    uint(value.ClientesID),
		}
		contactosFormato.EmailsContacto = append(contactosFormato.EmailsContacto, contacto)
	}
	return contactosFormato, nil
}
func (s *service) DeleteContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error) {
	//Valido el email
	isValid := commons.IsEmailValid(contactoReporte.ClienteEmail)
	if !isValid {
		return errors.New("debe enviar un correo valido")
	}
	//Verifico que el id sea mayor a 0
	if contactoReporte.ClienteID < 1 {
		return errors.New("debe enviar un cliente_id por valido")
	}
	entityContactos := contactoReporte.DtosToEntity()
	err = s.repository.DeleteContactosReportesRepository(entityContactos)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateContactosReportesService(contactoReporte administraciondtos.RequestContactosReportes) (err error) {
	if contactoReporte.ClienteID < 1 {
		return errors.New("debe enviar un cliente_id valido")
	}
	if len(contactoReporte.ClienteEmail) == 0 || len(contactoReporte.ClienteEmailNuevo) == 0 {
		return errors.New("debe enviar un campo cliente_email y cliente_email_nuevo")
	}

	//Valido el email
	isValidEmail := commons.IsEmailValid(contactoReporte.ClienteEmail)
	isValidEmailNuevo := commons.IsEmailValid(contactoReporte.ClienteEmailNuevo)
	if !isValidEmail || !isValidEmailNuevo {
		return errors.New("los correos deben ser válidos")
	}
	entityContactos := contactoReporte.DtosToEntity()
	entityContactosNuevo := entities.Contactosreporte{
		Email:      contactoReporte.ClienteEmailNuevo,
		ClientesID: int64(contactoReporte.ClienteID),
	}
	//Controlo si existe un usuario con ese correo, si existe informo que no puede actualizar por uno existente
	cr, _ := s.repository.GetContactosReportesByIdEmailRepository(entityContactosNuevo)
	if cr.ID > 0 {
		return errors.New("el correo ya se encuentra registrado para el usuario")
	}
	//De lo contrario actualizo
	err = s.repository.UpdateContactosReportesRepository(entityContactos, entityContactosNuevo)

	if err != nil {
		return
	}
	return nil
}

func (s *service) CreateUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error) {
	//Valido el email
	isValid := commons.IsEmailValid(request.Email)
	if !isValid {
		erro = errors.New("debe enviar un correo valido")
		return
	}
	// //Verifico que logitud nombre sea mayor a 0
	// if len(request.Nombre) < 1 {
	// 	erro = errors.New("debe enviar un nombre")
	// 	return
	// }

	usuarioBloqueado := request.ToEntity()

	erro = s.repository.CreateUsuarioBloqueadoRepository(usuarioBloqueado)
	if erro != nil {
		logs.Error(erro)
		return
	}
	return
}

func (s *service) UpdateUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error) {
	//Valido el email
	isValid := commons.IsEmailValid(request.Email)
	if !isValid {
		erro = errors.New("debe enviar un correo valido")
		return
	}
	// //Verifico que logitud nombre sea mayor a 0
	// if len(request.Nombre) < 1 {
	// 	erro = errors.New("debe enviar un nombre")
	// 	return
	// }

	usuarioBloqueado := request.ToEntity()
	usuarioBloqueado.ID = request.Id

	erro = s.repository.UpdateUsuarioBloqueadoRepository(usuarioBloqueado)
	if erro != nil {
		logs.Error(erro)
		return
	}
	return
}

func (s *service) DeleteUsuarioBloqueadoService(request administraciondtos.RequestUserBloqueado) (erro error) {

	//Verifico que el id sea mayor a 0
	if request.Id < 1 {
		erro = errors.New("debe enviar un id de usuario bloqueado")
		return
	}

	usuarioBloqueado := request.ToEntity()
	usuarioBloqueado.ID = request.Id

	erro = s.repository.DeleteUsuarioBloqueadoRepository(usuarioBloqueado)
	if erro != nil {
		logs.Error(erro)
		return
	}
	return
}

func (s *service) GetUsuariosBloqueadoService(filtro filtros.UsuarioBloqueadoFiltro) (response administraciondtos.ResponseUsuariosBloqueados, erro error) {

	usuarios, totalUsuarios, erro := s.repository.GetUsuariosBloqueadoRepository(filtro)

	response.FromEntities(usuarios)

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, totalUsuarios)
	}

	return
}
func (s *service) CallFraudePersonas(cuil string) (object interface{}, erro error) {

	ruta := config.URL_FRAUDE + "fraude/verificacion-cliente"

	type busquedaCuil struct {
		Cuil string
	}

	newBusqueda := busquedaCuil{
		Cuil: cuil,
	}

	object, erro = s.utilService.RunEndpoint("POST", ruta, nil, newBusqueda, nil, false)

	return
}

// func (s *service) CreateSoporteService(soporte administraciondtos.RequestSoporte) ( error) {
// 	err:=soporte.IsValidCreate()
// 	if err != nil {
// 		return err
// 	}
// 	var archivonombre string
// 	if soporte.File != nil {
// 		//Abrir archivo , leer converir a array de buffer y guardar
// 		buffer, err := soporte.File.Open()
// 		if err != nil {
// 			return err
// 		}
// 		//Al finalizar cerrar el archivo
// 		defer buffer.Close()

// 		//Creo el contexto
// 		ctx:= context.Background()

// 		//Leo el archivo
// 		fileBytes, err := ioutil.ReadAll(buffer)
// 		if err != nil {
// 			msj := "error a leer datos del archivo:" + soporte.File.Filename
// 			logs.Error(msj)
// 			err = errors.New(msj)
// 			return err
// 		}
// 		//Obtengo la extension
// 		archivo_extension := strings.Split(soporte.File.Filename, ".")

// 		// Obtengo el tiempo en nanosegundos para poder crear el nombre del archivo
// 		currentTime := time.Now()
// 		nanoseconds := currentTime.UnixNano()

// 		url_minio:="wee/soporte"
// 		archivonombre = fmt.Sprintf("%s/%d-%s", url_minio,nanoseconds, archivo_extension[0])
// 		archivotipo := archivo_extension[len(archivo_extension)-1]

// 		// Se guarda en S3
// 		//fileBytes = contenido []byte ,archivonombre= carpeta y nombre del archivo donde se guarda ,archivotipo= extension
// 		err = s.store.PutObject(ctx, fileBytes, archivonombre, archivotipo)
// 		if err != nil {
// 			logs.Error("No se pudo guardar el archivo")
// 			return err
// 		}
// 	}
// 	soporteEntity:= entities.Soporte{
// 		Nombre:soporte.Nombre,
// 		Email:soporte.Email,
// 		Consulta: soporte.Consulta,
// 		Archivo: archivonombre,
// 		Estado: entities.EnumSoporte("espera"),
// 		Abierto: false,
// 		Visto: false,
// 	}
// 	err = s.repository.CreateSoporteRepository(soporteEntity)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (s *service) PutSoporteService(soporte  administraciondtos.RequestSoporte) (err error){

//		soporteEntity:= entities.Soporte{
//			Model:gorm.Model{
//				ID: uint(soporte.Id),
//			},
//			Visto: soporte.Visto,
//			Estado: entities.EnumSoporte(soporte.Estado),
//			Abierto: soporte.Abierto,
//		}
//		err = s.repository.UpdateSoporteRepository(soporteEntity)

//		if err != nil {
//			return
//		}
//		return
//	}
func (s *service) EstadoApiService() (err error) {
	filtroConfiguraciones := filtros.ConfiguracionFiltro{
		Nombre: "ESTADO_APLICACION",
	}
	_, _, err = s.repository.GetConfiguraciones(filtroConfiguraciones)
	if err != nil {
		return
	}
	return
}

func (s *service) GetHistorialOperacionesService(filtro filtros.RequestHistorial) (response administraciondtos.ResponseHistorial, erro error) {
	historial, total, erro := s.repository.GetHistorialOperacionesRepository(filtro)

	if erro != nil {
		return
	}

	response.FromEntities(historial)

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	return
}

func (s *service) UpsertEnvioService(request administraciondtos.RequestEnvios) (erro error) {

	erro = request.Validate()

	if erro != nil {
		return
	}

	envioModificado := request.ToEnvio()

	return s.repository.UpsertEnvioRepository(envioModificado)

}
