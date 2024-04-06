package mockrepository

import (
	"context"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepositoryAdministracion struct {
	mock.Mock
}

func (mock *MockRepositoryAdministracion) BeginTx()    {}
func (mock *MockRepositoryAdministracion) CommitTx()   {}
func (mock *MockRepositoryAdministracion) RollbackTx() {}

func (mock *MockRepositoryAdministracion) PagoById(pagoID int64) (*entities.Pago, error) {
	args := mock.Called(pagoID)
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error) {
	args := mock.Called(cliente, limit, offset)
	result := args.Get(0)
	return result.(*[]entities.Cuenta), int64(args.Int(1)), args.Error(2)
}
func (mock *MockRepositoryAdministracion) CuentaByID(cuenta int64) (*entities.Cuenta, error) {
	args := mock.Called(cuenta)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error) {
	args := mock.Called(ctx, cuenta)
	return args.Bool(0), args.Error(1)
}

func (mock *MockRepositoryAdministracion) SavePagotipo(tipo *entities.Pagotipo) (bool, error) {
	args := mock.Called(tipo)
	return args.Bool(0), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error) {
	args := mock.Called(clienteId)
	result := args.Get(0)
	return result.([]entities.Cuenta), args.Error(1)

}

func (mock *MockRepositoryAdministracion) GetPagosByUUID(uuid []string) (pagos []*entities.Pago, erro error) {
	args := mock.Called(uuid)
	result := args.Get(0)
	return result.([]*entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), int64(args.Int(1)), args.Error(2)
}
func (mock *MockRepositoryAdministracion) GetPagosRepository(filtro filtros.PagoFiltro) (pagos []administraciondtos.ResponsePago, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]administraciondtos.ResponsePago), int64(args.Int(1)), args.Error(2)
}
func (mock *MockRepositoryAdministracion) GetItemsPagos(filtro filtros.PagoItemFiltro) ([]administraciondtos.PagoItems, error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]administraciondtos.PagoItems), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagos []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagoestado), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagoestadoexterno), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados entities.Pagoestado, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pagoestado), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {
	args := mock.Called(cuentaId)
	result := args.Get(0)
	return result.(administraciondtos.SaldoCuentaResponse), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {
	args := mock.Called(clienteId)
	result := args.Get(0)
	return result.(administraciondtos.SaldoClienteResponse), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMovimientos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Movimiento), int64(args.Int(1)), args.Error(2)
}
func (mock *MockRepositoryAdministracion) GetMovimientosTransferencias(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]entities.Movimiento), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosTransferencia(ctx context.Context, movimiento []*entities.Movimiento) error {
	args := mock.Called(ctx, movimiento)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateTransferencias(ctx context.Context, transferencias []*entities.Transferencia) (erro error) {
	args := mock.Called(ctx, transferencias)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateTransferenciasComisiones(ctx context.Context, transferencias []*entities.Transferenciacomisiones) (erro error) {
	args := mock.Called(ctx, transferencias)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) BajaMovimiento(ctx context.Context, movimientos []*entities.Movimiento, motivoBaja string) error {
	args := mock.Called(ctx, movimientos, motivoBaja)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error) {
	args := mock.Called(cierreLotes)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error) {
	args := mock.Called(pagos, pagoEstadoId)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	args := mock.Called(idMedioPago)
	result := args.Get(0)
	return result.([]administraciondtos.PlanCuotasResponseDetalle), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaCierreLotes(reversion bool) (prismaCierreLotes []entities.Prismacierrelote, erro error) {
	args := mock.Called(reversion)
	result := args.Get(0)
	return result.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaPagoIntentos(siteTransaccionId string) (pagos entities.Pago, erro error) {
	args := mock.Called(siteTransaccionId)
	result := args.Get(0)
	return result.(entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosCierreLote(ctx context.Context, movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse) (erro error) {

	args := mock.Called(ctx, movimientoCierreLote)

	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCLApilinkPagosRepository(ctx context.Context, pg administraciondtos.RegistroClPagosApilink) (erro error) {
	args := mock.Called(ctx, pg)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) ActualizarPagosClRapipagoRepository(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error) {

	args := mock.Called(pagosclrapiapgo)

	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) SaveCuentacomision(comision *entities.Cuentacomision) error {
	args := mock.Called(comision)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error) {
	args := mock.Called(id)
	return args.Get(0).([]entities.Installment), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) {
	args := mock.Called(fechaDesde)
	return args.Get(0).([]entities.Mediopagoinstallment), args.Error(1)
}
func (mock *MockRepositoryAdministracion) BuildRICuentasCliente(request ribcradtos.RICuentasClienteRequest) (ri []ribcradtos.RiCuentaCliente, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiCuentaCliente), args.Error(1)
}

func (mock *MockRepositoryAdministracion) BuildRIDatosFondo(request ribcradtos.RiDatosFondosRequest) (ri []ribcradtos.RiDatosFondos, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiDatosFondos), args.Error(1)
}
func (mock *MockRepositoryAdministracion) BuilRIInfestaditica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiInfestadistica), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Configuracione), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Cliente), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cliente), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetTransferencias(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferencia, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Transferencia), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetTransferenciasComisiones(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferenciacomisiones, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Transferenciacomisiones), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error) {
	args := mock.Called(ctx, cliente)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error) {
	args := mock.Called(ctx, cliente)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) DeleteCliente(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error) {
	args := mock.Called(ctx, rubro)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error) {
	args := mock.Called(ctx, rubro)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Rubro), args.Error(1)
}
func (mock *MockRepositoryAdministracion) GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Rubro), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error) {
	args := mock.Called(ctx, cuenta)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) DeleteCuenta(id uint64) (erro error) {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channels []int64, cuotas []string) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}
func (mock *MockRepositoryAdministracion) UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pagotipo), args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagotipo), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeletePagoTipo(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateChannel(ctx context.Context, request entities.Channel) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Channel), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Channel), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeleteChannel(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cuentacomision), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Cuentacomision), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeleteCuentaComision(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error) {
	args := mock.Called(ctx, cuenta)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateTransferencias(listas bancodtos.ResponseConciliacion) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Impuesto), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (erro error) {
	args := mock.Called(ctx, impuesto)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetConsultarDebines(request linkdebin.RequestDebines) (cierreLotes []*entities.Apilinkcierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]*entities.Apilinkcierrelote), args.Error(1)
}
func (mock *MockRepositoryAdministracion) CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error) {
	args := mock.Called(installmentActual, installmentNew, listaPlanCuotas)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetInstallmentById(id uint) (planCuotas entities.Installment, erro error) {
	args := mock.Called(id)
	result := args.Get(0)
	return result.(entities.Installment), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdatePagosNotificados(listaPagosNotificar []uint) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdateCierreloteApilink(request linkdebin.RequestListaUpdateDebines) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateImpuestoRepository(ctx context.Context, request entities.Impuesto) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]*entities.Rapipagocierrelote), args.Error(1)
}
func (mock *MockRepositoryAdministracion) GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Webservicespeticione), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (pagostipochannel []entities.Pagotipochannel, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Pagotipochannel), args.Error(1)
}

func (mock *MockRepositoryAdministracion) DeletePagoTipoChannel(id uint64) (erro error) {
	args := mock.Called()
	return args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreatePagoTipoChannel(ctx context.Context, request entities.Pagotipochannel) (id uint64, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Channelarancele), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) DeleteChannelsArancel(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Channelarancele), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Mediopago), args.Error(1)
}
func (mock *MockRepositoryAdministracion) GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismapxcuatroregistro), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismamxtotalesmovimiento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Movimiento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
	args := mock.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerCierreLoteEnDisputaRepository(estadoDisputa int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	args := mock.Called(estadoDisputa, filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Prismacierrelote), args.Error(1)
}
func (mock *MockRepositoryAdministracion) ObtenerCierreLoteContraCargoRepository(estadoReversion int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	args := mock.Called(estadoReversion, filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Prismacierrelote), args.Error(1)
}

// Preferences
func (mock *MockRepositoryAdministracion) PostPreferencesRepository(preferenceEntity entities.Preference) (erro error) {
	args := mock.Called(preferenceEntity)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetPreferencesRepository(clienteEntity entities.Cliente) (entityPreference entities.Preference, erro error) {
	args := mock.Called(clienteEntity)
	result := args.Get(0)
	return result.(entities.Preference), args.Error(1)
}
func (mock *MockRepositoryAdministracion) DeletePreferencesRepository(clienteEntity entities.Cliente) (erro error) {
	args := mock.Called(clienteEntity)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateSolicitudRepository(solicitudEntity entities.Solicitud) (erro error) {
	args := mock.Called(solicitudEntity)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdatePagosDev(pagos []uint) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (movimientos []entities.Rapipagocierrelote, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]entities.Rapipagocierrelote), 1, args.Error(1)
}

func (mock *MockRepositoryAdministracion) ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error) {
	args := mock.Called(parametrosVslido, filtro)
	resultado := args.Get(0)
	return resultado.([]entities.Pago), args.Error(1)

}

func (mock *MockRepositoryAdministracion) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosTemporalesCierreLote(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Reporte), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerCierreLoteRapipago(id int64) (response *entities.Rapipagocierrelote, erro error) {
	args := mock.Called(id)
	result := args.Get(0)
	return result.(*entities.Rapipagocierrelote), args.Error(1)
}

// CRUD contactosReportes
func (mock *MockRepositoryAdministracion) CreateContactosReportesRepository(contactoReportes entities.Contactosreporte) (err error) {
	args := mock.Called(contactoReportes)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) UpdateContactosReportesRepository(contactoReportes, contactoReportesNuevo entities.Contactosreporte) (err error) {
	args := mock.Called(contactoReportes, contactoReportesNuevo)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) ReadContactosReportesRepository(contactoReportes entities.Contactosreporte) ([]entities.Contactosreporte, error) {
	args := mock.Called(contactoReportes)
	result := args.Get(0)
	return result.([]entities.Contactosreporte), args.Error(1)
}
func (mock *MockRepositoryAdministracion) DeleteContactosReportesRepository(contactoReportes entities.Contactosreporte) (err error) {
	args := mock.Called(contactoReportes)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetContactosReportesByIdEmailRepository(contactoReportes entities.Contactosreporte) (entities.Contactosreporte, error) {
	args := mock.Called(contactoReportes)
	result := args.Get(0)
	return result.(entities.Contactosreporte), args.Error(1)
}
func (mock *MockRepositoryAdministracion) GetSubcuenta(filtro filtros.CuentaFiltro) (cuenta entities.Subcuenta, erro error) {
	return
}
func (mock *MockRepositoryAdministracion) GuardarSubcuentas(ctx context.Context, request []administraciondtos.SubcuentaRequest) (bool, error) {
	return true, nil
}
func (mock *MockRepositoryAdministracion) SubcuentaByCuentaPage(cuentaId int64, limit, offset int) (*[]entities.Subcuenta, int64, error) {
	return nil, 0, nil
}
func (mock *MockRepositoryAdministracion) GetSubcuentasByCuentaId(cuentaId uint) (subcuentas []*entities.Subcuenta, erro error) {
	return nil, nil
}
func (mock *MockRepositoryAdministracion) DeleteSubcuenta(ctx context.Context, id uint64) (erro error) {
	return
}
func (mock *MockRepositoryAdministracion) SaveSubcuentaInTransaction(ctx context.Context, tx *gorm.DB, subcuenta *entities.Subcuenta) (bool, error) {

	return true, nil
}

func (mock *MockRepositoryAdministracion) CreateUsuarioBloqueadoRepository(usuario entities.Usuariobloqueados) (erro error) {
	args := mock.Called(usuario)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdateUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error {
	args := mock.Called(usuarioBloqueado)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) DeleteUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error {
	args := mock.Called(usuarioBloqueado)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetUsuariosBloqueadoRepository(filtro filtros.UsuarioBloqueadoFiltro) (bloqueados []entities.Usuariobloqueados, totalBloqueados int64, erro error) {
	args := mock.Called(filtro)
	return args.Get(0).([]entities.Usuariobloqueados), args.Get(1).(int64), args.Error(0)
}

// CRUD soporte
func (mock *MockRepositoryAdministracion) CreateSoporteRepository(soporte entities.Soporte) (erro error) {
	args := mock.Called(soporte)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) UpdateSoporteRepository(soporte entities.Soporte) (erro error) {
	args := mock.Called(soporte)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetHistorialOperacionesRepository(filtro filtros.RequestHistorial) (historial []entities.HistorialOperaciones, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	return args.Get(0).([]entities.HistorialOperaciones), args.Get(1).(int64), args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpsertEnvioRepository(envio entities.Envio) error {
	args := mock.Called(envio)
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) GetConsultarClMultipagoRepository(filtro filtros.RequestClMultipago) (clmultipago []entities.Multipagoscierrelote, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	return args.Get(0).([]entities.Multipagoscierrelote), args.Get(1).(int64), args.Error(0)
}
