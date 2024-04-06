package mockservice

import (
	"context"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	webhook "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/stretchr/testify/mock"
)

type MockAdministracionService struct {
	mock.Mock
}

func (mock *MockAdministracionService) GetPagoByID(pagoID int64) (*entities.Pago, error) {
	args := mock.Called(pagoID)
	resultado := args.Get(0)
	return resultado.(*entities.Pago), args.Error(1)
}
func (mock *MockAdministracionService) GetCuentasByCliente(cliente int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Cuenta, error) {
	args := mock.Called(cliente, number, size)
	resultadoMeta := args.Get(0)
	resultadoLinks := args.Get(1)
	resultadoCuentas := args.Get(2)
	return resultadoMeta.(*dtos.Meta), resultadoLinks.(*dtos.Links), resultadoCuentas.(*[]entities.Cuenta), args.Error(3)
}
func (mock *MockAdministracionService) PostCuenta(ctx context.Context, cuenta administraciondtos.CuentaRequest) (bool, error) {
	args := mock.Called(ctx, cuenta)
	resultado := args.Bool(0)
	return resultado, args.Error(1)
}
func (mock *MockAdministracionService) PostPagotipo(ctx context.Context, pagotipo *entities.Pagotipo) (bool, error) {
	args := mock.Called(ctx, pagotipo)
	resultado := args.Bool(0)
	return resultado, args.Error(1)
}
func (mock *MockAdministracionService) PostCuentaComision(ctx context.Context, comision *entities.Cuentacomision) error {
	args := mock.Called(ctx, comision)
	return args.Error(1)
}
func (mock *MockAdministracionService) GetSaldoCuentaService(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {
	args := mock.Called(cuentaId)
	resultado := args.Get(0)
	return resultado.(administraciondtos.SaldoCuentaResponse), args.Error(1)
}
func (mock *MockAdministracionService) GetSaldoClienteService(clientId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {
	args := mock.Called(clientId)
	resultado := args.Get(0)
	return resultado.(administraciondtos.SaldoClienteResponse), args.Error(1)
}
func (mock *MockAdministracionService) GetMovimientosAcumulados(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoAcumuladoResponsePaginado, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.MovimientoAcumuladoResponsePaginado), args.Error(1)
}

func (mock *MockAdministracionService) GetMovimientos(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoPorCuentaResponsePaginado, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.MovimientoPorCuentaResponsePaginado), args.Error(1)
}

func (mock *MockAdministracionService) CreateNotificacionService(notificacion entities.Notificacione) error {
	args := mock.Called(notificacion)
	return args.Error(1)
}

func (mock *MockAdministracionService) CreateLogService(log entities.Log) error {
	args := mock.Called(log)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetPagosEstadosService(buscarPorFinal, final bool) (estados []entities.Pagoestado, erro error) {
	args := mock.Called(buscarPorFinal, final)
	resultado := args.Get(0)
	return resultado.([]entities.Pagoestado), args.Error(1)
}
func (mock *MockAdministracionService) BuildTransferenciaCliente(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenicaCliente, cuentaId uint64) (response linktransferencia.ResponseTransferenciaCreateLink, erro error) {
	args := mock.Called(ctx, requerimientoId, request, cuentaId)
	resultado := args.Get(0)
	return resultado.(linktransferencia.ResponseTransferenciaCreateLink), args.Error(1)
}
func (mock *MockAdministracionService) BuildCierreLoteApiLinkService() (response administraciondtos.RegistroClPagosApilink, erro error) {
	args := mock.Called()
	resultado := args.Get(0)
	return resultado.(administraciondtos.RegistroClPagosApilink), args.Error(1)
}
func (mock *MockAdministracionService) BuildMovimientoApiLink(listaCierre []*entities.Apilinkcierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {
	args := mock.Called(listaCierre)
	resultadoListaPagosModificar := args.Get(0)
	resultadoListaPagosEstadosLog := args.Get(1)
	resultadoListaMovimientos := args.Get(2)
	resultadoListaPagoIntentos := args.Get(3)
	movimientoCierreLote.ListaPagos = resultadoListaPagosModificar.([]entities.Pago)
	movimientoCierreLote.ListaPagosEstadoLogs = resultadoListaPagosEstadosLog.([]entities.Pagoestadologs)
	movimientoCierreLote.ListaMovimientos = resultadoListaMovimientos.([]entities.Movimiento)
	movimientoCierreLote.ListaPagoIntentos = resultadoListaPagoIntentos.([]entities.Pagointento)
	return movimientoCierreLote, args.Error(1)
}
func (mock *MockAdministracionService) CreateMovimientosService(ctx context.Context, movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse) (erro error) {
	args := mock.Called(ctx, movimientoCierreLote)
	return args.Error(1)
}
func (mock *MockAdministracionService) GetPlanCuotas(idMedioPago uint) (response []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	args := mock.Called(idMedioPago)
	resultado := args.Get(0)
	return resultado.([]administraciondtos.PlanCuotasResponseDetalle), args.Error(1)
}
func (mock *MockAdministracionService) BuildPrismaMovimiento() (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {
	args := mock.Called()
	resultListaPrismaCierreLote := args.Get(0)
	resultListaPagos := args.Get(1)
	resultListaPagoEstadoLogs := args.Get(2)
	resultListamovimientos := args.Get(3)
	resultListaPagoIntentos := args.Get(4)
	movimientoCierreLote.ListaCLPrisma = resultListaPrismaCierreLote.([]entities.Prismacierrelote)
	movimientoCierreLote.ListaPagos = resultListaPagos.([]entities.Pago)
	movimientoCierreLote.ListaPagosEstadoLogs = resultListaPagoEstadoLogs.([]entities.Pagoestadologs)
	movimientoCierreLote.ListaMovimientos = resultListamovimientos.([]entities.Movimiento)
	movimientoCierreLote.ListaPagoIntentos = resultListaPagoIntentos.([]entities.Pagointento)
	return movimientoCierreLote, args.Error(1)
}
func (mock *MockAdministracionService) GetInteresesPlanes() (planes []administraciondtos.PlanCuotasResponse, erro error) {
	args := mock.Called()
	planes = args.Get(0).([]administraciondtos.PlanCuotasResponse)
	erro = args.Error(1)
	return
}

func (mock *MockAdministracionService) GetInformacionSupervision(request ribcradtos.GetInformacionSupervisionRequest) (ri ribcradtos.RiInformacionSupervisionReponse, erro error) {
	args := mock.Called(request)
	resultado := args.Get(0)
	return resultado.(ribcradtos.RiInformacionSupervisionReponse), args.Error(1)
}

func (mock *MockAdministracionService) GetInformacionEstadistica(request ribcradtos.GetInformacionEstadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {
	args := mock.Called(request)
	resultado := args.Get(0)
	return resultado.([]ribcradtos.RiInfestadistica), args.Error(1)
}

func (mock *MockAdministracionService) RIInfestadistica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {
	args := mock.Called(request)
	resultado := args.Get(0)
	return resultado.([]ribcradtos.RiInfestadistica), args.Error(1)
}

func (mock *MockAdministracionService) RIGuardarArchivos(request ribcradtos.RIGuardarArchivosRequest) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockAdministracionService) BuildInformacionSupervision(request ribcradtos.BuildInformacionSupervisionRequest) (ruta string, erro error) {
	args := mock.Called(request)
	return args.String(0), args.Error(1)
}

func (mock *MockAdministracionService) BuildInformacionEstadistica(request ribcradtos.BuildInformacionEstadisticaRequest) (ruta string, erro error) {
	args := mock.Called(request)
	return args.String(0), args.Error(1)
}

func (mock *MockAdministracionService) ModificarEstadoPagosExpirados() (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockAdministracionService) UpdateConfiguracionService(ctx context.Context, config administraciondtos.RequestConfiguracion) (erro error) {
	args := mock.Called(ctx, config)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (response administraciondtos.ResponseConfiguraciones, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseConfiguraciones), args.Error(1)
}

func (mock *MockAdministracionService) GetClienteService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacion, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseFacturacion), args.Error(1)
}

func (mock *MockAdministracionService) GetClientesService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacionPaginado, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseFacturacionPaginado), args.Error(1)
}

func (mock *MockAdministracionService) GetCuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseCuenta, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseCuenta), args.Error(1)
}

func (mock *MockAdministracionService) GetTransferencias(filtro filtros.TransferenciaFiltro) (response administraciondtos.TransferenciaResponsePaginado, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.TransferenciaResponsePaginado), args.Error(1)
}

func (mock *MockAdministracionService) CreateClienteService(ctx context.Context, request administraciondtos.ClienteRequest) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockAdministracionService) UpdateClienteService(ctx context.Context, cliente administraciondtos.ClienteRequest) (erro error) {
	args := mock.Called(ctx, cliente)
	return args.Error(0)
}
func (mock *MockAdministracionService) DeleteClienteService(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockAdministracionService) CreateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}
func (mock *MockAdministracionService) UpdateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}
func (mock *MockAdministracionService) GetRubroService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubro, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseRubro), args.Error(1)
}
func (mock *MockAdministracionService) GetRubrosService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubros, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseRubros), args.Error(1)
}

func (mock *MockAdministracionService) UpdateCuentaService(ctx context.Context, request administraciondtos.CuentaRequest) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}
func (mock *MockAdministracionService) DeleteCuentaService(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockAdministracionService) SendSolicitudCuenta(request administraciondtos.SolicitudCuentaRequest) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockAdministracionService) CreatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockAdministracionService) UpdatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetPagoTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagoTipo, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponsePagoTipo), args.Error(1)
}

func (mock *MockAdministracionService) GetPagosTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagosTipo, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponsePagosTipo), args.Error(1)
}

func (mock *MockAdministracionService) DeletePagoTipoService(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockAdministracionService) CreateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockAdministracionService) UpdateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetChannelService(filtro filtros.ChannelFiltro) (channel administraciondtos.ResponseChannel, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseChannel), args.Error(1)
}
func (mock *MockAdministracionService) GetChannelsService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannels, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseChannels), args.Error(1)
}
func (mock *MockAdministracionService) DeleteChannelService(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockAdministracionService) CreateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}
func (mock *MockAdministracionService) UpdateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}
func (mock *MockAdministracionService) GetCuentaComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentaComision, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseCuentaComision), args.Error(1)
}
func (mock *MockAdministracionService) GetCuentasComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentasComision, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponseCuentasComision), args.Error(1)
}
func (mock *MockAdministracionService) DeleteCuentaComisionService(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockAdministracionService) UpdateConfiguracionSendEmailService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {
	args := mock.Called(requerimientoId, request)
	resultado := args.Get(0)
	return resultado.(linkconsultadestinatario.ResponseConsultaDestinatarioLink), args.Error(1)
}

func (mock *MockAdministracionService) GetPagosConsulta(req administraciondtos.RequestPagosConsulta) (*[]administraciondtos.ResponsePagosConsulta, error) {
	args := mock.Called(req)
	return args.Get(0).(*[]administraciondtos.ResponsePagosConsulta), args.Error(1)
}

func (mock *MockAdministracionService) RetiroAutomaticoClientes(ctx context.Context) (erro error) {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *MockAdministracionService) CreateCuentaApilinkService(request linkcuentas.LinkPostCuenta) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockAdministracionService) DeleteCuentaApilinkService(request linkcuentas.LinkDeleteCuenta) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockAdministracionService) GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error) {
	args := mock.Called()
	resultado := args.Get(0)
	return resultado.([]linkcuentas.GetCuentasResponse), args.Error(1)
}

func (mock *MockAdministracionService) GetPagosService(filtro filtros.PagoFiltro) (response administraciondtos.ResponsePagos, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.(administraciondtos.ResponsePagos), args.Error(1)
}

func (mock *MockAdministracionService) SetApiKeyService(ctx context.Context, request *administraciondtos.CuentaRequest) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockAdministracionService) ActualizarPagosClRapipagoService(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error) {
	args := mock.Called(pagosclrapiapgo)
	return args.Error(0)
}

func (mock *MockAdministracionService) BuildNotificacionPagosCLApilink(request []linkdebin.ResponseDebinesEliminados) (response []webhook.WebhookResponse, debinID []uint64, erro error) {
	args := mock.Called()
	resultado := args.Get(0)
	debin := []uint64{1, 2, 3}
	return resultado.([]webhook.WebhookResponse), debin, args.Error(2)
}

func (mock *MockAdministracionService) BuildNotificacionPagosCLRapipago(filtro filtros.PagoEstadoFiltro) (response []webhook.WebhookResponse, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]webhook.WebhookResponse), args.Error(1)
}
