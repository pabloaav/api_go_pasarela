package administracion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	BeginTx()
	RollbackTx()
	CommitTx()

	//CUENTAS
	CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error)
	GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error)
	CuentaByID(cuenta int64) (*entities.Cuenta, error)
	SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error)
	UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error)
	DeleteCuenta(id uint64) (erro error)
	SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error)
	GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error)

	GetSubcuenta(filtro filtros.CuentaFiltro) (cuenta entities.Subcuenta, erro error)
	SubcuentaByCuentaPage(cuentaId int64, limit, offset int) (*[]entities.Subcuenta, int64, error)
	GetSubcuentasByCuentaId(cuentaId uint) (subcuentas []*entities.Subcuenta, erro error)
	DeleteSubcuenta(ctx context.Context, id uint64) (erro error)
	SaveSubcuentaInTransaction(ctx context.Context, tx *gorm.DB, subcuenta *entities.Subcuenta) (bool, error)
	GuardarSubcuentas(ctx context.Context, request []administraciondtos.SubcuentaRequest) (bool, error)

	//CLIENTES
	CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error)
	UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error)
	DeleteCliente(ctx context.Context, id uint64) (erro error)
	GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error)
	GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error)
	GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error)

	//PAGOS
	GetPagosByUUID(uuid []string) (pagos []*entities.Pago, erro error)
	GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error)
	GetPagosRepository(filtro filtros.PagoFiltro) (pagos []administraciondtos.ResponsePago, totalFilas int64, erro error)
	GetItemsPagos(filtro filtros.PagoItemFiltro) ([]administraciondtos.PagoItems, error)
	GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error)
	GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagos []entities.Pagointento, erro error)
	GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error)
	GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error)
	GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados entities.Pagoestado, erro error)
	PagoById(pagoID int64) (*entities.Pago, error)
	SavePagotipo(tipo *entities.Pagotipo) (bool, error)
	ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error)

	//ABM RUBROS
	CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error)
	UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error)
	GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error)
	GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error)

	//ABM PAGOS TIPOS
	CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channels []int64, cuotas []string) (id uint64, erro error)
	UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error)
	GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error)
	GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error)
	DeletePagoTipo(ctx context.Context, id uint64) (erro error)

	//ABM CHANNELS
	CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error)
	UpdateChannel(ctx context.Context, request entities.Channel) (erro error)
	GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error)
	GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error)
	DeleteChannel(ctx context.Context, id uint64) (erro error)

	//ABM CUENTAS COMISIONES
	CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error)
	UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error)
	GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error)
	GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error)
	DeleteCuentaComision(ctx context.Context, id uint64) (erro error)

	// ABM IMPUESTOS
	GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error)
	CreateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (id uint64, erro error)
	UpdateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (erro error)

	// ABM CHANNELS ARANCELES
	GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error)
	GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error)
	CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error)
	UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error)
	DeleteChannelsArancel(ctx context.Context, id uint64) (erro error)
	/*
		Devuelve el saldo actual de una cuenta específica.
		Se debe informar el id de la cuenta.
	*/
	GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error)

	/*
		Devuelve el saldo actual de un cliente específico.
		Se debe informar una lista de cuentas del cliente.
	*/
	GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error)

	//MOVIMIENTOS
	GetMovimientos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, totalFilas int64, erro error)
	BajaMovimiento(ctx context.Context, movimientos []*entities.Movimiento, motivoBaja string) error
	SaveCuentacomision(comision *entities.Cuentacomision) error
	GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error)
	GetMovimientosTransferencias(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error)

	//TRANSFERENCIAS
	CreateMovimientosTransferencia(ctx context.Context, movimiento []*entities.Movimiento) error
	CreateTransferencias(ctx context.Context, transferencias []*entities.Transferencia) (erro error)
	GetTransferencias(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferencia, totalFilas int64, erro error)
	CreateTransferenciasComisiones(ctx context.Context, transferencias []*entities.Transferenciacomisiones) (erro error)
	GetTransferenciasComisiones(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferenciacomisiones, totalFilas int64, erro error)

	/* actualizar el estado de una transferencia con los datos conciliados del banco*/
	UpdateTransferencias(listas bancodtos.ResponseConciliacion) error

	/*
		Modifica el estado de una lista de pagos y además crea un pago estado log
	*/
	UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error)

	//ABM PLAN DE CUOTAS
	/*
		REVIEW: este codigo se debe revisar teniendo en cuenta los cambios que se hicieron en la BD para actualizar installmentdetails
		"CreatePlanCuotasByInstallmenIdRepository"
	*/
	/* Obtener el plan de cuotas */
	GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error)
	// obtiene todos los planes de cuotas
	GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error)
	// obtengo todos los planes por id
	GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error)
	// obtengo un plan de cuotas por id
	GetInstallmentById(id uint) (planCuotas entities.Installment, erro error)
	CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error)

	//CIERRE LOTE
	CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error)
	CreateMovimientosCierreLote(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error)
	GetPrismaCierreLotes(reversion bool) (prismaCierreLotes []entities.Prismacierrelote, erro error)
	CreateMovimientosTemporalesCierreLote(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error)

	// & Actualizar estados pagos y clrapipago
	ActualizarPagosClRapipagoRepository(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error)

	// & Crear y actualizar pagos
	CreateCLApilinkPagosRepository(ctx context.Context, pg administraciondtos.RegistroClPagosApilink) (erro error)
	// consultar debines eliminados para ser procesados en el cierre de lote
	GetConsultarDebines(request linkdebin.RequestDebines) (cierreLotes []*entities.Apilinkcierrelote, erro error)
	//&end apilinkcierrelote

	//RI BCRA
	BuildRICuentasCliente(request ribcradtos.RICuentasClienteRequest) (ri []ribcradtos.RiCuentaCliente, erro error)
	BuildRIDatosFondo(request ribcradtos.RiDatosFondosRequest) (ri []ribcradtos.RiDatosFondos, erro error)
	BuilRIInfestaditica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)

	//CONFIGURACIONES
	GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error)
	UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error)

	// UPDATE PAGOS NOTIFICACIDOS
	UpdatePagosNotificados(listaPagosNotificar []uint) (erro error)

	// consultar movimientos de la tabla rapipagos para luego ser procesados en el cierre de lote
	GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error)
	UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error)

	//PagoTipoChannel
	GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error)
	DeletePagoTipoChannel(id uint64) (erro error)
	CreatePagoTipoChannel(ctx context.Context, pagotipochannel entities.Pagotipochannel) (id uint64, erro error)

	//PETICIONES WEBSERVICES
	GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error)

	// MEDIO-PAGO
	GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error)

	// ARCHIVOS SIBIDOS DE "CIERRE LOTE, PRISMA PX Y PRISMA MX"
	GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error)
	GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error)
	GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error)

	// Obtener registro de cierre de lote rapipago
	ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error)
	ObtenerCierreLoteEnDisputaRepository(estadoDisputa int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error)
	ObtenerCierreLoteContraCargoRepository(estadoReversion int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error)

	ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error)

	ObtenerCierreLoteRapipago(id int64) (response *entities.Rapipagocierrelote, erro error)

	// Preferencias
	PostPreferencesRepository(preferenceEntity entities.Preference) (erro error)
	GetPreferencesRepository(clienteEntity entities.Cliente) (entityPreference entities.Preference, erro error)
	DeletePreferencesRepository(clienteEntity entities.Cliente) (erro error)

	// UPDATE PAGOS MOVIMIENTOS DEV // solo se utiliza para generar movimientos en ambiente sandbox y dev
	UpdatePagosDev(pagos []uint) (erro error)
	// Solicitud de Cuenta
	CreateSolicitudRepository(solicitudEntity entities.Solicitud) (erro error)

	// ? consultar repository CLlotes para herramienta wee
	// consultar movimeintos para herramienta wee
	GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (clrapiapgo []entities.Rapipagocierrelote, totalFilas int64, erro error)

	// consultar cierres de lote para herramienta wee Multipago
	GetConsultarClMultipagoRepository(filtro filtros.RequestClMultipago) (clmultipago []entities.Multipagoscierrelote, totalFilas int64, erro error)

	// consultar pagosintentos calculo de comisiones temporales
	GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error)

	// CL apilink
	UpdateCierreloteApilink(request linkdebin.RequestListaUpdateDebines) (erro error)

	GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error)
	GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error)

	// USUARIOS BLOQUEADOS
	CreateUsuarioBloqueadoRepository(userbloqueado entities.Usuariobloqueados) (erro error)
	GetUsuariosBloqueadoRepository(filtro filtros.UsuarioBloqueadoFiltro) (bloqueados []entities.Usuariobloqueados, totalBloqueados int64, erro error)
	UpdateUsuarioBloqueadoRepository(userbloqueado entities.Usuariobloqueados) (erro error)
	DeleteUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error

	//CONTACTOSREPORTES
	CreateContactosReportesRepository(contactoReportes entities.Contactosreporte) error
	DeleteContactosReportesRepository(contactoReportes entities.Contactosreporte) error
	ReadContactosReportesRepository(contactoReportes entities.Contactosreporte) ([]entities.Contactosreporte, error)
	UpdateContactosReportesRepository(contactoReportes, contactoReportesNuevo entities.Contactosreporte) error
	GetContactosReportesByIdEmailRepository(contactoReportes entities.Contactosreporte) (entities.Contactosreporte, error)

	GetHistorialOperacionesRepository(filtro filtros.RequestHistorial) (historial []entities.HistorialOperaciones, totalFilas int64, erro error)

	UpsertEnvioRepository(envio entities.Envio) error
	// //Soporte
	// CreateSoporteRepository(soporte entities.Soporte) (erro error)
	// UpdateSoporteRepository(soporte entities.Soporte) (erro error)

}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
	utilService      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService, t util.UtilService) Repository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
		utilService:      t,
	}
}

func (r *repository) BeginTx() {
	r.SQLClient.TX = r.SQLClient.DB
	r.SQLClient.DB = r.SQLClient.Begin()
}
func (r *repository) CommitTx() {
	r.SQLClient.Commit()
	r.SQLClient.DB = r.SQLClient.TX
}
func (r *repository) RollbackTx() {
	r.SQLClient.Rollback()
	r.SQLClient.DB = r.SQLClient.TX
}

/* ****************************************************************************************** */

func (r *repository) GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error) {

	resp := r.SQLClient.Model(entities.Cuenta{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.DistintoId > 0 {
		resp.Where("id != ?", filtro.DistintoId)
	}

	if len(filtro.Cbu) > 0 {
		resp.Where("cbu = ?", filtro.Cbu)
	}

	if len(filtro.Cvu) > 0 {
		resp.Where("cvu = ?", filtro.Cvu)
	}

	if len(filtro.ApiKey) > 0 {
		resp.Where("apikey = ?", filtro.ApiKey)
	}

	resp.Preload("Cliente")

	resp.First(&cuenta)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			/*REVIEW en el caso de que no devuelva un elemento */
			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCuenta: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetSubcuenta(filtro filtros.CuentaFiltro) (subcuenta entities.Subcuenta, erro error) {

	resp := r.SQLClient.Model(entities.Subcuenta{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}
	fmt.Println(resp.Error)

	if filtro.DistintoId > 0 {
		resp.Where("id != ?", filtro.DistintoId)
	}

	if len(filtro.Cbu) > 0 {
		resp.Where("cbu = ?", filtro.Cbu)
	}

	resp.Preload("Cuenta")

	resp.Find(&subcuenta)

	fmt.Println(resp.Error)

	if resp.Error != nil {
		if strings.Contains(resp.Error.Error(), "record not found") {
			erro = fmt.Errorf("Registro no encontrado con ese ID. Id: %s", strconv.Itoa(int(filtro.Id)))
			return
		}
	}

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			/*REVIEW en el caso de que no devuelva un elemento */
			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetSubcuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetSubcuenta: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) PagoById(pagoID int64) (*entities.Pago, error) {
	var pago entities.Pago

	res := r.SQLClient.Model(entities.Pago{}).Preload("PagosResultados.Pagoresultadodetalle").Find(&pago, pagoID)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron pagos con id %d", pagoID)
	}

	return &pago, nil
}

func (r *repository) CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error) {
	var cuentas []entities.Cuenta
	var count int64

	res := r.SQLClient.Model(entities.Cuenta{}).Preload("Pagotipos")
	res.Where("clientes_id = ?", cliente)
	res.Count(&count)
	res.Limit(limit)
	res.Offset(offset)
	res.Find(&cuentas)

	if res.RowsAffected <= 0 {
		return nil, 0, fmt.Errorf("no se encontraron cuentas para el cliente con id %d", cliente)
	}

	return &cuentas, count, nil
}

func (r *repository) SubcuentaByCuentaPage(cuentaId int64, limit, offset int) (*[]entities.Subcuenta, int64, error) {
	var subcuentas []entities.Subcuenta
	var count int64

	res := r.SQLClient.Model(entities.Subcuenta{}).
		Preload("Cuenta").
		Where("cuentas_id = ?", cuentaId)

	res.Count(&count)
	res.Limit(limit)
	res.Offset(offset)
	res.Find(&subcuentas)

	if res.RowsAffected <= 0 {
		return nil, 0, fmt.Errorf("no se encontraron subcuentas para la cuenta con id %d", cuentaId)
	}

	return &subcuentas, count, nil
}

func (r *repository) CuentaByID(cuentaID int64) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta

	res := r.SQLClient.Model(entities.Cuenta{}).Preload("Pagotipos").Preload("Cuentacomisiones").Find(&cuenta, cuentaID)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con id %d", cuentaID)
	}

	return &cuenta, nil
}

func (r *repository) GetSubcuentasByCuentaId(cuentaId uint) (subcuentas []*entities.Subcuenta, erro error) {
	resp := r.SQLClient.Model(entities.Subcuenta{})

	if cuentaId > 0 {
		resp.Where("cuentas_id = ?", cuentaId)
	}

	resp.Find(&subcuentas)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_SUBCUENTAS)
		return
	}
	return
}

func (r *repository) SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error) {
	res := r.SQLClient.WithContext(ctx)
	if cuenta.ID == 0 {
		res = res.Create(&cuenta)
	} else {
		res = res.Model(&cuenta).Updates(cuenta)
	}

	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al guardar cuenta: %s", res.Error.Error())
	}

	err := r.auditarAdministracion(res.Statement.Context, true)
	if err != nil {
		return false, fmt.Errorf("no es posible la auditoría: %v", err)
	}

	return true, nil
}

func (r *repository) GuardarSubcuentas(ctx context.Context, request []administraciondtos.SubcuentaRequest) (bool, error) {
	tx := r.SQLClient.Begin()

	for _, v := range request {
		if v.Modificado == 1 || v.Id == 0 {
			subcuenta := v.ToCuenta()
			ok, err := r.SaveSubcuentaInTransaction(ctx, tx, &subcuenta)
			if err != nil {
				tx.Rollback()
				return false, err
			}

			if !ok {
				tx.Rollback()
				return false, nil
			}
		}
	}

	tx.Commit()

	return true, nil
}
func (r *repository) SaveSubcuentaInTransaction(ctx context.Context, tx *gorm.DB, subcuenta *entities.Subcuenta) (bool, error) {
	res := tx.WithContext(ctx)

	if subcuenta.ID == 0 {
		res = res.Create(&subcuenta)
		if res.Error != nil {
			return false, fmt.Errorf("error al crear subcuenta: %v", res.Error)
		}
	} else {
		entidad := entities.Subcuenta{
			Model: gorm.Model{ID: subcuenta.ID},
		}

		if len(subcuenta.Cbu) > 0 {
			res = res.Model(&entidad).Select("tipo", "cuentas_id", "cbu", "nombre", "email", "porcentaje", "cuentas_id", "aplica_porcentaje", "aplica_costo_servicio").Updates(&subcuenta)
		} else {
			res = res.Model(&entidad).Select("tipo", "cuentas_id", "nombre", "email", "porcentaje", "cuentas_id", "aplica_porcentaje", "aplica_costo_servicio").Updates(&subcuenta)
		}

		if res.Error != nil {
			return false, fmt.Errorf("error al actualizar subcuenta: %v", res.Error)
		}
		if res.Error != nil {
			return false, fmt.Errorf("error al actualizar subcuenta: %v", res.Error)
		}
	}

	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al guardar subcuenta: %s", res.Error.Error())
	}

	err := r.auditarAdministracion(res.Statement.Context, true)
	if err != nil {
		return false, fmt.Errorf("no es posible la auditoría: %v", err)
	}

	return true, nil
}

func (r *repository) SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: cuenta.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Select("apikey").Updates(cuenta)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_APIKEY)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "SetApiKey",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
	resp := r.SQLClient.Model(entities.Cuenta{}).Where("apikey = ?", apikey)
	resp.Preload("Pagotipos")
	resp.Find(&cuenta)
	if resp.Error != nil {
		logs.Error("error al consultar cuenta: " + resp.Error.Error())
		erro = errors.New(ERROR_CONSULTAR_CUENTA)
		return
	}
	if resp.RowsAffected <= 0 {
		logs.Error("no existe cuenta")
		erro = errors.New(ERROR_CONSULTAR_CUENTA)
		return
	}
	return
}

func (r *repository) UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: cuenta.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,apikey,created_at,deleted_at").Select("*").Updates(cuenta)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) DeleteCuenta(id uint64) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) SavePagotipo(tipo *entities.Pagotipo) (bool, error) {
	if tipo.ID == 0 {
		res := r.SQLClient.Create(tipo)
		if res.RowsAffected <= 0 {
			return false, fmt.Errorf("error al crear tipo de pago: %s", res.Error.Error())
		}
		return true, nil
	}
	res := r.SQLClient.Model(entities.Pagotipo{}).Where("id = ?", tipo.ID).Updates(tipo)
	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al actualizar tipo de pago: %s", res.Error.Error())
	}

	return true, nil
}

func (r *repository) ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error) {
	/*
		los fitros que recibe son:
		- 1 uuid
		- arrays de uuid
		- rango de fecha
		- external reference
	*/
	resp := r.SQLClient.Model(entities.Pago{})

	if filtro.PagoEstadosId != 0 {
		resp.Where("pagoestados_id <> ?", filtro.PagoEstadosId)
	}
	if len(filtro.PagosTipoIds) > 0 {
		resp.Where("pagostipo_id in (?)", filtro.PagosTipoIds)
	}

	if parametrosVslido.Uuuid || parametrosVslido.Uuids {
		resp.Where("uuid in ? ", filtro.Uuids)
	}

	if parametrosVslido.ExternalReference {
		resp.Where("external_reference = ?", filtro.ExternalReference)
	}

	if parametrosVslido.RangoFecha {
		resp.Where("created_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}
	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}
	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}
	if filtro.CargaPagoIntentos {
		stateComments := []string{"approved", "INICIADO"}
		resp.Preload("PagoIntentos", "state_comment in ? ", stateComments)
		resp.Preload("PagoIntentos.Mediopagos")
		resp.Preload("PagoIntentos.Mediopagos.Channel")
		resp.Preload("PagoIntentos.Movimientos", "tipo = ?", "C")
		resp.Preload("PagoIntentos.Movimientos.Movimientocomisions")
		resp.Preload("PagoIntentos.Movimientos.Movimientoimpuestos")

	}
	resp.Find(&entityPagos)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_PAGO)
		return
	}
	return
}

func (r *repository) SaveCuentacomision(comision *entities.Cuentacomision) error {
	if comision.ID == 0 {
		res := r.SQLClient.Create(comision)
		if res.RowsAffected <= 0 {
			return fmt.Errorf("error al crear comision: %s", res.Error.Error())
		}
		return nil
	}
	res := r.SQLClient.Model(entities.Cuentacomision{}).Where("id = ?", comision.ID).Updates(comision)
	if res.RowsAffected <= 0 {
		return fmt.Errorf("error al actualizar comision: %s", res.Error.Error())
	}
	return nil
}

func (r *repository) GetPagosByUUID(uuids []string) (pagos []*entities.Pago, erro error) {

	if len(uuids) > 0 {
		resp := r.SQLClient.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC")
		}).Preload("PagoTipos").Where("uuid IN ?", uuids).Find(&pagos)
		if resp.Error != nil {
			erro = resp.Error
		}
	}

	return
}

func (r *repository) GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestado{})

	if filtro.BuscarPorFinal {
		resp.Where("final = ?", filtro.Final)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}
	if filtro.EstadoId != 0 {
		resp.Where("id = ?", filtro.EstadoId)
	}
	resp.Find(&estados)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosEstados",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estado entities.Pagoestado, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestado{})

	if filtro.BuscarPorFinal {
		if filtro.Final {

			resp.Where("final = ?", filtro.Final)
		}
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}

	if filtro.EstadoId > 0 {
		resp.Where("id = ?", filtro.EstadoId)
	}

	resp.First(&estado)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagoEstado",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {

	resp := r.SQLClient.Table("movimientos as m").Select("m.cuentas_id, sum(m.monto) as total").
		Where("m.deleted_at IS NULL").Group("m.cuentas_id").Having("m.cuentas_id", cuentaId).Scan(&saldo)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_SALDO_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetSaldoCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {

	resp := r.SQLClient.Table("movimientos as m").
		Select("m.cuentas_id, sum(m.monto) as total").
		Joins("inner join cuentas as c on c.id = m.cuentas_id").
		Joins("left join clientes as cl on cl.id = c.clientes_id").
		Where("m.deleted_at IS NULL").Where("cl.id", clienteId).Group("m.cuentas_id").Scan(&saldo)

	if resp.Error != nil {
		erro = resp.Error
	}

	return
}

func (r *repository) GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error) {

	resp := r.SQLClient.Where("clientes_id", clienteId).Find(&cuentas)

	if resp.Error != nil {
		erro = resp.Error
	}

	return
}

func (r *repository) UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error) {

	var estadosLogs []entities.Pagoestadologs

	for i := range pagos {
		pagoEstado := entities.Pagoestadologs{
			PagosID:       int64(pagos[i].ID),
			PagoestadosID: pagos[i].PagoestadosID,
		}
		estadosLogs = append(estadosLogs, pagoEstado)
	}

	erro = r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// Creo los logs de estados
		if err := tx.Create(&estadosLogs).Error; err != nil {

			erro := fmt.Errorf(ERROR_CREAR_ESTADO_LOGS)

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       err.Error(),
				Funcionalidad: "UpdateEstadoPagos",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				logs.Error(err.Error())
			}

			return erro
		}
		// Modifico los estados de los pagos
		if err := tx.Model(&pagos).Omit(clause.Associations).UpdateColumns(entities.Pago{PagoestadosID: int64(pagoEstadoId), Model: gorm.Model{UpdatedAt: time.Now()}}).Error; err != nil {

			erro := fmt.Errorf(ERROR_UPDATE_PAGO)

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       err.Error(),
				Funcionalidad: "UpdateEstadoPagos",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				logs.Error(err.Error())
			}

			return erro
		}

		return nil
	})

	return
}

func (r *repository) GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !filtro.FechaPagoFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	}

	if len(filtro.Ids) > 0 {
		resp.Where("id IN ?", filtro.Ids)
	}

	if filtro.PagoEstadosId > 0 {
		resp.Where("pagoestados_id", filtro.PagoEstadosId)
	}

	if len(filtro.TiempoExpiracion) > 0 {
		resp.Where("timestampdiff(day, created_at, now() ) >= ?", filtro.TiempoExpiracion)
	}

	if filtro.CargaPagoIntentos {
		resp.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC")
		})
	}
	if len(filtro.Uuids) > 0 {
		resp.Where("uuid IN ?", filtro.Uuids)
	}

	if len(filtro.ExternalReference) > 0 {
		resp.Where("external_reference = ?", filtro.ExternalReference)
	}

	if filtro.CargaMedioPagos {
		if filtro.CargarChannel {
			resp.Preload("PagoIntentos.Mediopagos.Channel")
		} else {
			resp.Preload("PagoIntentos.Mediopagos")
		}
	}

	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}

	if filtro.CargarCuenta {
		if filtro.CuentaId > 0 {
			resp.Preload("PagosTipo.Cuenta", "id = ?", filtro.CuentaId)
		} else {
			resp.Preload("PagosTipo.Cuenta")
		}
	}

	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	resp.First(&pago)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPago",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPago: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !filtro.FechaPagoFin.IsZero() {
		if filtro.FiltroFechaPaid {
			resp.Preload("PagoIntentos").Joins("JOIN pasarela.pagointentos as p1 ON (pagos.id = p1.pagos_id) LEFT OUTER JOIN pasarela.pagointentos p2 ON (pagos.id = p2.pagos_id AND (p1.created_at < p2.created_at OR (p1.created_at = p2.created_at AND p1.id < p2.id)))").Where("p2.id IS NULL").Where("cast(p1.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin).Order("p1.paid_at ASC")
			if filtro.Ordenar {
				if filtro.Descendente {
					resp.Order("pagos.PagoIntentos.paid_at DESC")
				}
				if !filtro.Descendente {
					resp.Order("pagos.PagoIntentos.paid_at ASC")
				}
			}

		} else {
			resp.Where("pagos.created_at BETWEEN cast(? as datetime) AND cast(? as datetime)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
			if filtro.Ordenar {
				if filtro.Descendente {
					resp.Order("pagos.created_at DESC")
				}
				if !filtro.Descendente {
					resp.Order("pagos.created_at ASC")
				}
			}
		}
	}

	if filtro.BuscarNotificado {
		if !filtro.Notificado {
			resp.Where("notificado = ?", filtro.Notificado)
		}
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	if len(filtro.Ids) > 0 {

		resp.Where("id IN ?", filtro.Ids)
	}

	if len(filtro.Uuids) > 0 {
		resp.Where("uuid IN ?", filtro.Uuids)
	}

	if len(filtro.PagoEstadosIds) > 0 {
		resp.Where("pagoestados_id IN ?", filtro.PagoEstadosIds)
	}

	if filtro.PagoEstadosId > 0 {
		resp.Where("pagoestados_id", filtro.PagoEstadosId)
	}

	if filtro.PagosTipoId > 0 {
		resp.Where("pagostipo_id = ?", filtro.PagosTipoId)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("payer_name LIKE ?", "%"+filtro.Nombre+"%")
	}

	if len(filtro.ExternalReference) > 0 {
		resp.Where("external_reference LIKE ?", "%"+filtro.ExternalReference+"%")

	}
	if len(filtro.ExternalReferences) > 0 {
		resp.Where("external_reference in (?)", filtro.ExternalReferences)
	}

	if len(filtro.PagosTipoIds) > 0 {
		resp.Where("pagostipo_id IN ?", filtro.PagosTipoIds)
	}

	if !filtro.VisualizarPendientes && len(filtro.PagoEstadosIds) == 0 {
		filtro := filtros.PagoEstadoFiltro{
			Nombre: "pending",
		}

		estadoPendiente, err := r.GetPagoEstado(filtro)

		if err != nil {
			erro = err
			return
		}

		resp.Where("pagoestados_id != ?", estadoPendiente.ID)
	}

	if len(filtro.TiempoExpiracion) > 0 {
		resp.Where("timestampdiff(day, created_at, now() ) >= ?", filtro.TiempoExpiracion)
	}

	if filtro.CargarCuenta {
		if filtro.CuentaId > 0 {
			resp.Preload("PagosTipo.Cuenta", "cuentas.id = ?", filtro.CuentaId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.id = ?", filtro.CuentaId)
		} else {
			resp.Preload("PagosTipo.Cuenta")
		}
	}

	if filtro.CargaPagoIntentos {
		if filtro.CargaMedioPagos {
			if filtro.MedioPagoId > 0 {
				if filtro.CargarChannel {
					//NOTE: SE MODIFICO LA CONSULTA m.id por m.channels_id
					resp.Preload("PagoIntentos.Mediopagos.Channel").Joins("LEFT JOIN pagointentos as pi on pagos.id = pi.pagos_id  INNER JOIN mediopagos as m on m.id = pi.mediopagos_id INNER JOIN channels as ch on ch.id = m.channels_id").Where("m.channels_id = ?", filtro.MedioPagoId).Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")
				} else {
					resp.Preload("PagoIntentos.Mediopagos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id inner join INNER JOIN mediopagos as m on m.id = pi.mediopagos_id").Where("m.channels_id = ?", filtro.MedioPagoId).Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")
				}
			} else {
				if filtro.CargarChannel {
					resp.Preload("PagoIntentos", "state_comment = ? OR state_comment = ?", "approved", "INICIADO")
					resp.Preload("PagoIntentos.Mediopagos.Channel")
				} else {
					resp.Preload("PagoIntentos.Mediopagos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").Where("pi.external_id != ''").Order("pi.paid_at ASC")
				}
			}
		} else {
			resp.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
				return db.Order("id DESC")
			})
		}

	}

	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}

	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}

	if filtro.CargarPagosItems {
		resp.Preload("Pagoitems")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPagosRepository(filtro filtros.PagoFiltro) (pagos []administraciondtos.ResponsePago, totalFilas int64, erro error) {

	resp := r.SQLClient.Table("pagos")

	// pago intentos
	subquery := r.SQLClient.
		Table("pagointentos PAI").
		Select("MAX(id) as max_id, pagos_id, mediopagos_id, amount, paid_at, holder_number ").
		Where("paid_at != '0000-00-00 00:00:00'").
		Group("pagos_id")

	resp.Joins("INNER JOIN (?) PI ON PI.pagos_id = pagos.id", subquery)

	// mediopagos
	resp.Joins("INNER JOIN mediopagos M ON M.id = PI.mediopagos_id")

	// channels
	resp.Joins("INNER JOIN channels CH ON CH.id = M.channels_id")

	//pago tipos
	resp.Joins("INNER JOIN pagotipos PT ON PT.id = pagos.pagostipo_id")

	//pago estado
	resp.Joins("INNER JOIN pagoestados PE ON PE.id = pagos.pagoestados_id")

	//cuentas
	resp.Joins("INNER JOIN cuentas C ON C.id = PT.cuentas_id")

	// movimientos
	subquery1 := r.SQLClient.
		Table("movimientos").
		Select("MAX(id) as max_id_mov, pagointentos_id").
		Group("pagointentos_id")

	resp.Joins("LEFT JOIN (?) MOV ON MOV.pagointentos_id = PI.max_id", subquery1)

	// transferencias
	resp.Joins("LEFT JOIN transferencias TR ON TR.movimientos_id = MOV.max_id_mov")

	if filtro.CuentaId > 0 {
		resp.Where("C.id = ?", filtro.CuentaId)
	}

	if len(filtro.ExternalReference) > 0 {
		resp.Where("pagos.external_reference = ?", filtro.ExternalReference)
	}

	if len(filtro.PagoEstadosIds) > 0 {
		resp.Where("pagos.pagoestados_id IN (?)", filtro.PagoEstadosIds)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("pagos.payer_name LIKE ?", "%"+filtro.Nombre+"%")
	}

	if len(filtro.HolderNumber) > 0 {
		resp.Where("PI.holder_number = ?", filtro.HolderNumber)
	}

	if filtro.MedioPagoId > 0 {
		resp.Where("CH.id = ?", filtro.MedioPagoId)
	}

	if !filtro.FechaPagoFin.IsZero() && !filtro.FechaPagoInicio.IsZero() {
		resp.Where("pagos.created_at BETWEEN cast(? as datetime) AND cast(? as datetime)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	}

	if !filtro.VisualizarPendientes && len(filtro.PagoEstadosIds) == 0 {
		filtro := filtros.PagoEstadoFiltro{
			Nombre: "pending",
		}

		estadoPendiente, err := r.GetPagoEstado(filtro)

		if err != nil {
			erro = err
			return
		}

		resp.Where("pagos.pagoestados_id != ?", estadoPendiente.ID)
	}

	resp.Select(`
			pagos.id as identificador, 
			pagos.created_at as fecha, 
			C.cuenta, 
			PT.pagotipo, 
			pagos.external_reference, 
			pagos.payer_name, 
			PI.amount,
			PI.paid_at as fecha_pago,
			CH.channel, 
			CH.nombre AS nombre_channel, 
			PE.estado, 
			PE.nombre as nombre_estado, 
			PI.max_id AS ultimo_pago_intento_id, 
			TR.id as transferencia_id, 
			TR.referencia_bancaria, 
			TR.fecha_operacion as fecha_transferencia
		`)

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	if filtro.Descendente {
		resp.Order("pagos.created_at DESC")
	} else {
		resp.Order("pagos.created_at ASC")
	}

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetItemsPagos(filtro filtros.PagoItemFiltro) ([]administraciondtos.PagoItems, error) {

	var pagoItems []administraciondtos.PagoItems
	resp := r.SQLClient.Table("pagoitems")

	if filtro.PagoId > 0 {
		resp.Where("pagos_id = ?", filtro.PagoId)
	}

	resp.Select("description as descripcion, identifier as identificador, quantity as cantidad, amount as monto")

	resp.Find(&pagoItems)

	if resp.Error != nil {

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetItemsPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return pagoItems, nil
}

func (r *repository) GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	var details []entities.Installmentdetail
	response := r.SQLClient.Model(entities.Installmentdetail{}).Joins("Installment").Joins("LEFT JOIN mediopagos ON (Installment.id = mediopagos.installments_id) AND mediopagos.id = ?", idMedioPago).Find(&details)
	if response.Error != nil {
		erro = response.Error
	}
	for _, v := range details {
		planCuotas = append(planCuotas, administraciondtos.PlanCuotasResponseDetalle{
			InstallmentsID: v.InstallmentsID,
			Cuota:          uint(v.Cuota),
			Tna:            v.Tna,
			Tem:            v.Tem,
			Coeficiente:    v.Coeficiente,
		})
	}
	return
}

func (r *repository) GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) {
	res := r.SQLClient.Table("mediopagoinstallments as mpi")
	res.Preload("Installments")
	res.Preload("Installments.Installmentdetail")
	res.Find(&medioPagoInstallments)
	if res.Error != nil {
		logs.Info(res.Error)
		erro = errors.New(ERROR_CREAR_INSTALLMENT_DETAILS)
		return
	}
	return
}

func (r *repository) GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error) {
	res := r.SQLClient.Model(entities.Installment{}).Where("mediopagoinstallments_id = ?", id).Find(&installment)
	if res.Error != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	return
}

func (r *repository) GetInstallmentById(id uint) (installment entities.Installment, erro error) {
	res := r.SQLClient.Model(entities.Installment{}).Where("mediopagoinstallments_id = ?", id).Order("created_at desc").First(&installment)
	if res.Error != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	return
}

func (r *repository) CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error) {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(entities.Installment{}).Omit(clause.Associations).Where("id = ?", installmentActual.ID).Update("vigencia_hasta", installmentActual.VigenciaHasta)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(ERROR_ACTUALIZAR_INSTALLMENT)
		}
		installmentNew.Installmentdetail = listaPlanCuotas

		res = tx.Create(&installmentNew)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(ERROR_CREAR_INSTALLMENT_DETAILS)
		}
		return nil
	})
}

func (r *repository) GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error) {

	resp := r.SQLClient.Model(entities.Pagointento{})

	if len(filtro.ExternalIds) > 0 {

		resp.Where("external_id IN ?", filtro.ExternalIds)
	}

	if filtro.PagoIntentoAprobado {
		resp.Where("paid_at <> ?", "0000-00-00 00:00:00")
	}

	if filtro.ExternalId {
		resp.Where("external_id <>  ? OR external_id <>  ?", "", "0")
		// resp.Where("external_id <> 0")
	}

	if len(filtro.TransaccionesId) > 0 {
		resp.Where("transaction_id IN (?)", filtro.TransaccionesId)
	}
	if len(filtro.TicketNumber) > 0 {
		resp.Where("ticket_number IN (?)", filtro.TicketNumber)
	}

	if len(filtro.CodigoAutorizacion) > 0 {
		resp.Where("authorization_code IN (?)", filtro.CodigoAutorizacion)
	}

	if len(filtro.Barcode) > 0 {
		resp.Where("barcode IN (?)", filtro.Barcode)
	}

	if len(filtro.PagosId) > 0 {
		resp.Where("pagos_id IN (?)", filtro.PagosId)
	}

	if filtro.ChannelIdFiltro != 0 {
		resp.Joins("INNER JOIN mediopagos as mp ON mp.id = pagointentos.mediopagos_id and mp.channels_id = ?", filtro.ChannelIdFiltro)
	}

	if filtro.Channel {

		resp.Preload("Mediopagos")
		resp.Preload("Mediopagos.Channel")
	}

	if filtro.PagoEstadoIdFiltro != 0 {
		resp.Joins("INNER JOIN pagos as p ON p.id = pagointentos.pagos_id and p.pagoestados_id = ?", filtro.PagoEstadoIdFiltro)
	}

	if filtro.PagoTipoid != 0 {
		resp.Joins("INNER JOIN pagos as p ON p.id = pagointentos.pagos_id and p.pagostipo_id = ?", filtro.PagoTipoid)
	}

	if filtro.CargarPago {

		resp.Preload("Pago")

	}

	if filtro.CargarPagoTipo {
		resp.Preload("Pago.PagosTipo")
		if filtro.CargarCuenta {
			resp.Preload("Pago.PagosTipo.Cuenta")
			if filtro.CargarCliente {
				resp.Preload("Pago.PagosTipo.Cuenta.Cliente")
				if filtro.CargarImpuestos {
					resp.Preload("Pago.PagosTipo.Cuenta.Cliente.Iva")
					resp.Preload("Pago.PagosTipo.Cuenta.Cliente.Iibb")
				}
			}
			if filtro.CargarCuentaComision {
				resp.Preload("Pago.PagosTipo.Cuenta.Cuentacomisions")
			}
		}
	}

	if filtro.CargarPagoEstado {
		resp.Preload("Pago.PagoEstados")
	}

	if filtro.CargarMovimientos {
		resp.Preload("Movimientos")
	}

	if filtro.CargarInstallmentdetail {
		resp.Preload("Installmentdetail")
	}

	resp.Find(&pagosIntentos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_INTENTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosIntentos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Configuracione{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("nombre like ?", fmt.Sprintf("%%%s%%", filtro.Nombre))
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&configuraciones)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguraciones",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error) {

	entidad := entities.Configuracione{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at,nombre").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateConfiguracion",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

// ABM CLIENTES
func (r *repository) GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Cliente{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.ClientesIds) > 0 {
		resp.Where("id in ?", filtro.ClientesIds)
	}

	if filtro.CargarImpuestos {
		resp.Preload("Iva")
		resp.Preload("Iibb")
	}

	if filtro.CargarCuentas {
		resp.Preload("Cuentas")
	}

	if filtro.CargarRubros {
		resp.Preload("Cuentas.Rubro")
	}

	if filtro.RetiroAutomatico {
		resp.Where("retiro_automatico = ?", filtro.RetiroAutomatico)
	}

	if filtro.CargarContactos {
		resp.Preload("Contactosreportes")
	}

	// if filtro.CargarCuentaComision {
	// 	resp.Preload("Cuentas.Cuentacomisions")
	// }

	// if filtro.CargarTiposPago {
	// 	resp.Preload("Cuentas.Pagotipos")
	// }

	if filtro.CargarEnvio {
		resp.Preload("Envio")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&clientes)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetClientes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetClientes: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error) {

	resp := r.SQLClient.Model(entities.Cliente{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.DistintoId > 0 {
		resp.Where("id = ?", filtro.DistintoId)
	}

	if filtro.UserId > 0 {
		resp.Joins("JOIN clienteusers ON clientes.id = clienteusers.clientes_id AND clienteusers.user_id = ?", filtro.UserId)

	}

	if len(filtro.Cuit) > 0 {
		resp.Where("cuit = ?", filtro.Cuit)
	}

	if filtro.RetiroAutomatico {
		resp.Where("retiro_automatico = ?", filtro.RetiroAutomatico)
	}

	if filtro.CargarImpuestos {
		resp.Preload("Iva")
		resp.Preload("Iibb")
	}

	if filtro.CargarCuentas {
		resp.Preload("Cuentas")
	}

	if filtro.CargarRubros {
		resp.Preload("Cuentas.Rubro")
	}
	if filtro.CargarCuentaComision {
		resp.Preload("Cuentas.Cuentacomisions.Channel")
		resp.Preload("Cuentas.Cuentacomisions.ChannelArancel")
	}

	// if filtro.CargarCuentaComision {
	// 	resp.Preload("Cuentas.Cuentacomisions")
	// }

	if filtro.CargarTiposPago {
		resp.Preload("Cuentas.Pagotipos")
	}

	resp.First(&cliente)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			/*REVIEW en el caso de que no devuelva un elemento */
			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCliente: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&cliente)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CLIENTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(cliente.ID)

	erro = r.auditarAdministracion(result.Statement.Context, id)
	if erro != nil {
		return id, erro
	}

	return
}

func (r *repository) UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error) {

	entidad := entities.Cliente{
		Model: gorm.Model{ID: cliente.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(cliente)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_MODIFICAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, cliente)

	return
}
func (r *repository) DeleteCliente(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Cliente{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM RUBROS
func (r *repository) GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Rubro{})

	if len(filtro.Rubro) > 0 {
		resp.Where("rubro like ?", fmt.Sprintf("%%%s%%", filtro.Rubro))
	}

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&rubros)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_RUBROS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRubros",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetRubros: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error) {

	resp := r.SQLClient.Model(entities.Rubro{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Rubro) > 0 {
		resp.Where("rubro", filtro.Rubro)
	}

	resp.First(&rubro)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_RUBROS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetRubro: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error) {
	if rubro.ID > 0 {
		rubro.ID = 0
	}
	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&rubro)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_RUBRO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(rubro.ID)

	erro = r.auditarAdministracion(result.Statement.Context, id)

	return
}

func (r *repository) UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error) {

	entidad := entities.Rubro{
		Model: gorm.Model{ID: rubro.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(rubro)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_RUBRO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, rubro)

	return
}

// ABM PAGO TIPOS
func (r *repository) GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipo{})

	if len(filtro.PagoTipo) > 0 {
		resp.Where("pagotipo = ?", filtro.PagoTipo)
	}

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if filtro.IdCuenta > 0 {
		resp.Where("cuentas_id = ?", filtro.IdCuenta)
	}

	if filtro.CargarTipoPagoChannels {
		resp.Preload("Pagotipochannel.Channel")
		resp.Preload("Pagotipoinstallment")
	}

	//  filtro cargar los pagos y para filtrar por estado , pagos de los ultimos 3 dias para notificar al usuario
	if filtro.CargarPagos {
		resp.Preload("Pagos", "pagoestados_id IN ? AND notificado = ? AND cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.PagoEstadosIds, filtro.CargarPagosNotificado, filtro.FechaPagoInicio, filtro.FechaPagoFin)
		resp.Preload("Pagos.PagoEstados")
		resp.Preload("Pagos.PagoIntentos.Mediopagos.Channel")
	}

	// if filtro.CargarPagosIntentos {
	// 	if len(filtro.ExternalId) > 0 {
	// 		resp.Joins("INNER JOIN pagos as pg on pagotipos.id = pg.pagostipo_id INNER JOIN pagointentos as pi on pg.id = pi.pagos_id").
	// 			Where("pi.external_id IN ?", filtro.ExternalId)
	// 	}
	// 	resp.Preload("Pagos.PagoIntentos", func(db *gorm.DB) *gorm.DB {
	// 		return db.Where("pagointentos.external_id IN ?", filtro.ExternalId)
	// 	})

	// }

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPagosTipo: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipo{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if len(filtro.PagoTipo) > 0 {
		resp.Where("pagotipo", filtro.PagoTipo)
	}

	if filtro.CargarTipoPagoChannels {
		resp.Preload("Pagotipochannel.Channel")
		resp.Preload("Pagotipoinstallment")
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPagoTipo: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channel []int64, cuotas []string) (id uint64, erro error) {

	r.BeginTx()
	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_PAGO_TIPO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreatePagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		r.RollbackTx()
		return
	}

	id = uint64(request.ID)

	// agregar pagotipochannels
	for _, ch := range channel {
		entidadChannel := entities.Pagotipochannel{
			PagotiposId: uint(id),
			ChannelsId:  uint(ch),
		}

		resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

		if resultpagotipochannels.Error != nil {
			erro = fmt.Errorf(ERROR_CREAR_RUBRO)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       result.Error.Error(),
				Funcionalidad: "Createpagotipochannels",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
				logs.Error(mensaje)
			}
			r.RollbackTx()
			return
		}

	}

	// agregar cuotas
	for _, c := range cuotas {
		entidadCuota := entities.Pagotipointallment{
			PagotiposId: uint(id),
			Cuota:       c,
		}

		resultcuotas := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadCuota)

		if resultcuotas.Error != nil {
			erro = fmt.Errorf(ERROR_CREAR_RUBRO)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       result.Error.Error(),
				Funcionalidad: "Createpagotipoinstallment",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
				logs.Error(mensaje)
			}
			r.RollbackTx()
			return
		}

	}

	erro = r.auditarAdministracion(result.Statement.Context, request)
	r.CommitTx()

	return
}

func (r *repository) UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error) {

	r.BeginTx()
	entidad := entities.Pagotipo{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdatePagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		r.RollbackTx()
		return
	}

	if len(channels.Add) > 0 {

		for _, ch := range channels.Add {
			entidadChannel := entities.Pagotipochannel{
				PagotiposId: entidad.ID,
				ChannelsId:  uint(ch),
			}

			resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

			if resultpagotipochannels.Error != nil {
				erro = fmt.Errorf(ERROR_CREAR_RUBRO)
				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "Createpagotipochannels",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(channels.Delete) > 0 {
		for _, cdelete := range channels.Delete {
			result := r.SQLClient.WithContext(ctx).Where("pagotipos_id = ? AND channels_id = ?", entidad.ID, cdelete).Delete(&entities.Pagotipochannel{})

			if result.Error != nil {

				erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "DeletePagoTipo",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(cuotas.Add) > 0 {

		for _, c := range cuotas.Add {
			entidadChannel := entities.Pagotipointallment{
				PagotiposId: entidad.ID,
				Cuota:       c,
			}

			resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

			if resultpagotipochannels.Error != nil {
				erro = fmt.Errorf(ERROR_CREAR_RUBRO)
				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "CreatePagotipointallment",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(cuotas.Delete) > 0 {
		for _, cudelete := range cuotas.Delete {
			result := r.SQLClient.WithContext(ctx).Where("pagotipos_id = ? AND cuota = ?", entidad.ID, cudelete).Delete(&entities.Pagotipointallment{})

			if result.Error != nil {

				erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "DeletePagoTipo",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	erro = r.auditarAdministracion(result.Statement.Context, request)
	r.CommitTx()
	return
}

func (r *repository) DeletePagoTipo(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Pagotipo{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeletePagotipointallment",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}
	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM CHANNELS
func (r *repository) GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Channel{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Channel) > 0 {
		resp.Where("channel = ?", filtro.Channel)
	} else if len(filtro.Channels) > 0 {
		resp.Where("channel IN ?", filtro.Channels)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannels",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetChannels: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error) {

	resp := r.SQLClient.Model(entities.Channel{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarMedioPago {
		resp.Preload("Mediopagos")
	}

	if len(filtro.Channel) > 0 {
		resp.Where("channel = ?", strings.ToUpper(filtro.Channel))
	} else if len(filtro.Channels) > 0 {
		resp.Where("channel IN ?", filtro.Channels)
	}

	resp.First(&channel)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CHANNEL)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) UpdateChannel(ctx context.Context, request entities.Channel) (erro error) {

	entidad := entities.Channel{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) DeleteChannel(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Channel{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM CUENTA COMISION
func (r *repository) GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Cuentacomision{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CuentaId > 0 {
		resp.Where("cuentas_id = ?", filtro.CuentaId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.Channelarancel {
		resp.Preload("ChannelArancel")
	}

	if len(filtro.CuentaComision) > 0 {
		resp.Where("cuentacomision = ?", filtro.CuentaComision)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuentasComisiones",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCuentasComisiones: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error) {

	resp := r.SQLClient.Model(entities.Cuentacomision{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.CuentaComision) > 0 {
		resp.Where("cuentacomision = ?", filtro.CuentaComision)
	}

	if filtro.CuentaId > 0 {
		resp.Where("cuentas_id = ?", filtro.CuentaId)
	}
	// if filtro.ChannelId > 0 {
	// 	resp.Where("channels_id = ?", filtro.ChannelId).Order("vigencia_desde desc")
	// }

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if !filtro.FechaPagoVigencia.IsZero() {
		resp.Where("vigencia_desde <= ?", filtro.FechaPagoVigencia).Order("vigencia_desde desc")
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}
	if filtro.Channelarancel {
		resp.Preload("ChannelArancel")
	}
	// if filtro.Mediopagoid > 0 {
	// 	resp.Where("mediopagoid = ?", filtro.Mediopagoid)
	// }

	if filtro.ExaminarPagoCuota {
		resp.Where("mediopagoid = ? AND pagocuota = ?", filtro.Mediopagoid, filtro.PagoCuota)
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CUENTA_COMISION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error) {

	entidad := entities.Cuentacomision{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Impuesto{})
	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}
	if len(filtro.Tipo) > 0 {
		resp.Where("tipo = ? and activo = ?", strings.ToUpper(filtro.Tipo), 1)
	}
	if filtro.OrdenarPorFecha {
		resp.Order("fechadesde asc")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_IMPUESTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetImpuestosRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetImpuestosRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) CreateImpuestoRepository(ctx context.Context, request entities.Impuesto) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_IMPUESTO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateImpuestoRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	return
}

func (r *repository) UpdateImpuestoRepository(ctx context.Context, request entities.Impuesto) (erro error) {

	entidad := entities.Impuesto{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateImpuestoRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}
	return
}
func (r *repository) DeleteSubcuenta(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Subcuenta{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_SUBCUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteSubcuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

func (r *repository) DeleteCuentaComision(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Cuentacomision{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

func (r *repository) auditarAdministracion(ctx context.Context, resultado interface{}) error {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	if audit.Query == "" {
		audit.Operacion = "delete"
	} else {
		audit.Operacion = strings.ToLower(audit.Query[:6])
	}

	audit.Origen = "pasarela.administracion"

	res, _ := json.Marshal(resultado)
	audit.Resultado = string(res)

	err := r.auditoriaService.Create(&audit)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

/* update pagosm notificados*/
func (r *repository) UpdatePagosNotificados(listaPagosNotificar []uint) (erro error) {

	result := r.SQLClient.Table("pagos").Where("id IN ?", listaPagosNotificar).Updates(map[string]interface{}{"notificado": 1})
	if result.Error != nil {
		erro := fmt.Errorf("no se puedo actualizar los pagos notificados")
		return erro
	}
	if result.RowsAffected <= 0 {
		logs.Info("caso de no actualizacion de pagos notificados 0")
		return nil
	}
	logs.Info("cantidad de pagos actualizados con exito " + fmt.Sprintf("%v", result.RowsAffected))

	return nil

}

func (r *repository) GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Rapipagocierrelote{})

	if filtro.CargarMovConciliados {
		resp.Where("banco_external_id != ?", 0)
	} else {
		resp.Where("banco_external_id = ?", 0)
	}

	if filtro.PagosNotificado {
		resp.Where("pago_actualizado != ?", 0)
	} else {
		resp.Where("pago_actualizado = ?", 0)
	}

	resp.Preload("RapipagoDetalle")

	resp.Find(&response)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_RAPIPAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConsultarMovimientosRapipago",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetConsultarMovimientosRapipago: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return response, erro
}

func (r *repository) UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error) {

	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, valueCL := range cierreLotes {
			resp := tx.Model(entities.Rapipagocierrelote{}).Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"banco_external_id": valueCL.BancoExternalId, "enobservacion": valueCL.Enobservacion, "difbancocl": valueCL.Difbancocl})
			if resp.Error != nil {
				logs.Info(resp.Error)
				erro = errors.New("error: al actualizar tabla de cierre de lote rapipago")
				return erro
			}
		}

		for _, valueCL := range cierreLotes {
			for _, detalle := range valueCL.RapipagoDetalle {
				resp := tx.Model(entities.Rapipagocierrelotedetalles{}).Where("id = ?", detalle.ID).UpdateColumns(map[string]interface{}{"match": detalle.Match, "enobservacion": detalle.Enobservacion})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New("error: al actualizar tabla de cierre de lote rapipago detalle")
					return erro
				}
			}
		}
		return nil
	})

	return
}

// func (r *repository) UpdateCierreloteAndMoviminetosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
// 	r.SQLClient.Transaction(func(tx *gorm.DB) error {
// 		for _, valueCL := range entityCierreLote {
// 			resp := tx.Model(entities.Prismacierrelote{}).Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"prismamovimientodetalles_id": valueCL.PrismamovimientodetallesId, "fecha_pago": valueCL.FechaPago})
// 			if resp.Error != nil {
// 				logs.Info(resp.Error)
// 				erro = errors.New("error: al actualizar tabla de cierre de lote")
// 				return erro
// 			}
// 		}

// 		if err := tx.Model(&entities.Prismamovimientototale{}).Where("id in (?)", listaIdsCabecera).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
// 			logs.Info(err)
// 			erro = errors.New("error: al actualizar tabla Prisma Movimientos cabecera")
// 			return erro
// 		}
// 		if err := tx.Model(&entities.Prismamovimientodetalle{}).Where("id in (?)", listaIdsDetalle).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
// 			logs.Info(err)
// 			erro = errors.New("error: al actualizar tabla Prisma Movimientos detalle")
// 			return erro
// 		}
// 		return nil
// 	})

// 	return
// }

func (r *repository) GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Webservicespeticione{})

	if filtro.Id != 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.OrdenarPorFechaInv {
		resp.Order("updated_at asc")
	}

	if filtro.Operacion != "" {
		resp.Where("operacion", filtro.Operacion)
	}
	if filtro.Vendor != "" {
		resp.Where("vendor = ?", filtro.Vendor)
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&peticiones)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPeticionesWebServices",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (pagostipochannel []entities.Pagotipochannel, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipochannel{})

	if filtro.PagoTipoId > 0 {
		resp.Where("pagotipos_id", filtro.PagoTipoId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id= ?", filtro.ChannelId)
	}

	resp.Find(&pagostipochannel)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTipoChannelRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) DeletePagoTipoChannel(id uint64) (erro error) {

	entidad := entities.Pagotipochannel{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) CreatePagoTipoChannel(ctx context.Context, request entities.Pagotipochannel) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_IMPUESTO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreatePagoTipoChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	return
}

// ABM CHANNELS ARANCELES
func (r *repository) GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Channelarancele{})

	if filtro.RubrosId > 0 {
		resp.Where("rubros_id = ?", filtro.RubrosId)
	}
	if filtro.CargarRubro {
		resp.Preload("Rubro")
	}
	if filtro.PagoCuota {
		resp.Where("pagocuota = ? ", filtro.PagoCuota)
	} else {
		resp.Where("pagocuota = ? ", filtro.PagoCuota)
	}
	if !filtro.CargarAllMedioPago {
		if filtro.MedioPagoId > 0 {
			resp.Where("mediopagoid = ?", filtro.MedioPagoId)
		}
		if filtro.MedioPagoId == 0 {
			resp.Where("mediopagoid = 0")
		}
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannelsAranceles",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetChannelsAranceles: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CUENTA_COMISION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error) {

	entidad := entities.Channelarancele{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) DeleteChannelsArancel(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Channelarancele{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

func (r *repository) GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error) {

	resp := r.SQLClient.Model(entities.Channelarancele{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.RubrosId > 0 {
		resp.Where("rubros_id = ?", filtro.RubrosId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if filtro.CargarRubro {
		resp.Preload("Rubro")
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.OrdernarChannel {
		resp.Order("fechadesde desc")
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CHANNEL_ARANCEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannelArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error) {
	resp := r.SQLClient.Model(entities.Mediopago{})
	if filtro.IdMedioPago > 0 {
		resp.Where("id = ?", filtro.IdMedioPago)
	}
	resp.Find(&mediopago)
	if resp.Error != nil {
		erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMedioPagoRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMedioPagoRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as cl")
	resp.Select("cl.nombrearchivolote, cl.created_at, cl.deleted_at")
	resp.Unscoped()
	resp.Group("cl.nombrearchivolote, cl.created_at")
	resp.Order("cl.created_at desc")
	resp.Find(&entityCl)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_CIERRE_LOTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error) {
	resp := r.SQLClient.Table("prismapxcuatroregistros as px")
	resp.Select("px.nombrearchivo, px.created_at, px.deleted_at")
	resp.Unscoped()
	resp.Group("px.nombrearchivo, px.created_at")
	resp.Order("px.created_at desc")
	resp.Find(&entityPx)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_PRISMA_PX)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPrismaPxSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error) {
	resp := r.SQLClient.Table("prismamxtotalesmovimientos as mx")
	resp.Select("mx.nombrearchivo, mx.created_at, mx.deleted_at")
	resp.Unscoped()
	resp.Group("mx.nombrearchivo, mx.created_at")
	resp.Order("mx.created_at desc")
	resp.Find(&entityMx)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_PRISMA_MX)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPrismaPxSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error) {
	//consultar movimiento por nombre de archivo- verificar si existe
	//si existe retornar el nombre del archivo
	var result entities.Rapipagocierrelote
	res := r.SQLClient.Table("rapipagocierrelotes").Where("nombre_archivo = ?", nombre).Find(&result)

	if res.Error != nil {
		erro = res.Error
		return false, erro
	}
	if res.RowsAffected > 0 {
		existeArchivo = true
		return existeArchivo, nil
	} else {
		existeArchivo = false
		return existeArchivo, nil
	}
}

func (r *repository) ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error) {
	resp := r.SQLClient.Table("pagointentos as pi")

	if len(filtro.TransactionId) > 0 {
		resp.Where("pi.transaction_id in (?)", filtro.TransactionId)

	}
	if filtro.CargarPagos {
		resp.Joins("inner join pagos as p on p.id = pi.pagos_id")
		resp.Preload("Pago")
	}

	if filtro.CargarTiposPago {
		resp.Joins("inner join pagotipos as ptip on ptip.id = p.pagostipo_id")
		resp.Preload("Pago.PagosTipo")
	}

	if filtro.CargarCuentas {
		resp.Joins("inner join cuentas as c on c.id = ptip.cuentas_id and c.id = ? and c.clientes_id = ? ", filtro.IdCuenta, filtro.IdCliente)
		resp.Preload("Pago.PagosTipo.Cuenta")
	}

	resp.Find(&pagosEnDisputa)
	if resp.Error != nil {
		erro = errors.New(ERROR_OBTENER_PAGOS_DISPUTA)
	}

	return
}

func (r *repository) PostPreferencesRepository(preferenceEntity entities.Preference) (erro error) {
	err := r.SQLClient.Create(&preferenceEntity).Error
	if err != nil {
		erro = errors.New("no se pudo guardar la preferencia en la base de datos")
		return
	}

	return
}
func (r *repository) GetPreferencesRepository(clienteEntity entities.Cliente) (entityPreference entities.Preference, erro error) {
	res := r.SQLClient.Table("preferences").Where("clientes_id = ?", clienteEntity.ID)
	res.Last(&entityPreference)
	if res.RowsAffected == 0 {
		return
	}
	if res.Error != nil {
		erro = errors.New("no se pudo obtener las preferencias del usuario")
		return
	}
	return
}
func (r *repository) DeletePreferencesRepository(clienteEntity entities.Cliente) (erro error) {
	res := r.SQLClient.Where("clientes_id = ?", clienteEntity.ID).Delete(&entities.Preference{})
	if res.Error != nil {
		erro = errors.New("ocurrio un error al eliminar el registro")
		return
	}
	if res.RowsAffected == 0 {
		erro = errors.New("no se encontro ningun registro a eliminar")
		return
	}
	return nil
}

/* update pagosm notificados*/
func (r *repository) UpdatePagosDev(pagos []uint) (erro error) {

	result := r.SQLClient.Table("pagos").Where("id IN ?", pagos).Updates(map[string]interface{}{"pagoestados_id": 7})
	if result.Error != nil {
		erro := fmt.Errorf("no se puedo actualizar los pagos")
		return erro
	}
	if result.RowsAffected <= 0 {
		logs.Info("caso de no actualizacion de pagos notificados 0")
		return nil
	}
	logs.Info("cantidad de pagos actualizados con exito " + fmt.Sprintf("%v", result.RowsAffected))

	return nil
}

func (r *repository) CreateSolicitudRepository(solicitud entities.Solicitud) (erro error) {
	// Guardar la solicitud en la base de datos en su tabla correspondiente
	result := r.SQLClient.Model(entities.Solicitud{}).Create(&solicitud)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_SOLICITUD_CUENTA)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateSolicitud",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	return
}

// ? consultar clrapipago para herramietna wee
func (r *repository) GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (clrapiapgo []entities.Rapipagocierrelote, totalFilas int64, erro error) {

	resp := r.SQLClient.Unscoped().Model(entities.Rapipagocierrelote{})

	if filtro.CodigoBarra != "" {
		resp.Joins("INNER JOIN rapipagocierrelotedetalles as rpcl_detalles on rpcl_detalles.rapipagocierrelotes_id = rapipagocierrelotes.id")
		resp.Where("rpcl_detalles.codigo_barras = ? ", filtro.CodigoBarra)
	} else {

		if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
			resp.Where("cast(rapipagocierrelotes.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
		}
	}
	resp.Preload("RapipagoDetalle", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})

	// if filtro.Id > 0 {
	// 	resp.Where("id = ?", filtro.Id)
	// }

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf("error al cargar total filas de la columna")
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Order("created_at desc").Find(&clrapiapgo)

	return
}

func (r *repository) GetConsultarClMultipagoRepository(filtro filtros.RequestClMultipago) (clmultipago []entities.Multipagoscierrelote, totalFilas int64, erro error) {
	resp := r.SQLClient.Unscoped().Model(entities.Multipagoscierrelote{})

	if filtro.CodigoBarra != "" {
		resp.Joins("INNER JOIN multipagoscierrelotedetalles as mpcl_detalles on mpcl_detalles.multipagoscierrelotes_id = multipagoscierrelotes.id")
		resp.Where("mpcl_detalles.codigo_barras = ?", filtro.CodigoBarra)
	} else {
		if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
			resp.Where("cast(multipagoscierrelotes.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
		}
	}

	resp.Preload("MultipagoDetalle", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf("error al cargar total filas de la columna")
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Order("created_at desc").Find(&clmultipago)

	return
}

func (r *repository) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {

	resp := r.SQLClient.Model(entities.Pagointento{})

	if !filtro.FechaPagoFin.IsZero() {
		resp.Where("cast(pagointentos.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin)

	}
	if filtro.CargarPagoCalculado {
		resp.Where("calculado = ?", 1)
	} else {
		resp.Where("calculado = ?", 0)
	}

	if filtro.ClienteId > 0 {
		resp.Preload("Pago.PagosTipo.Cuenta", "cuentas.clientes_id = ?", filtro.ClienteId).Joins("INNER JOIN pagos as p on p.id = pagointentos.pagos_id INNER JOIN pagotipos as pt on pt.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", filtro.ClienteId)
	}

	if len(filtro.PagoEstadosIds) > 0 {
		resp.Preload("Pago", "pagos.pagoestados_id IN ?", filtro.PagoEstadosIds).Joins("INNER JOIN pagos as pg on pg.id = pagointentos.pagos_id").Where("pg.pagoestados_id IN ?", filtro.PagoEstadosIds)
		// resp.Where("pagoestados_id IN ?", filtro.PagoEstadosIds)
	}
	if filtro.PagoIntentoAprobado {

		resp.Where("paid_at <> ?", "0000-00-00 00:00:00")
	}

	if filtro.CargarMovimientosTemporales {
		resp.Preload("Movimientotemporale.Movimientocomisions")
		resp.Preload("Movimientotemporale.Movimientoimpuestos")
	}

	if filtro.CargarPagoItems {
		resp.Preload("Pago.Pagoitems")
	}

	if filtro.Channel {
		resp.Preload("Mediopagos.Channel")
	}
	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	// esto hace que cargue los pagointentos del pago, pero solo aquellos que tienen algun valor setado en el campo paid_at, y no los que tienen todo en cero ese campo
	resp.Preload("PagoIntentos", "paid_at")
	resp.Joins("INNER JOIN pagoestados AS PE ON pagos.pagoestados_id = PE.id INNER JOIN pagointentos AS PI ON pagos.id = PI.pagos_id INNER JOIN pagotipos AS PT ON pagos.pagostipo_id = PT.id INNER JOIN cuentas AS C ON PT.cuentas_id = C.id INNER JOIN pagoitems AS PIT ON pagos.id = PIT.pagos_id")
	// es necesario el DISTINCT porque al hacer join con pagoitems se duplican los registros del resultado de la consulta
	resp.Select("DISTINCT pagos.*")
	// Los estados de pagos exitosos son 4 y 7
	resp.Where("pagos.pagoestados_id IN ?", []uint{4, 7}).Where("PI.paid_at LIKE ?", filtro.Fecha[0]+"%").Where("C.id = ?", filtro.CuentaId)

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetSuccessPaymentsRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error) {

	resp := r.SQLClient.Model(entities.Reporte{})

	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(filtro.CuentaId),
	}

	// Averiguar el cliente segun la cuenta
	cuenta, err := r.GetCuenta(filtroCuenta)

	if err != nil {
		mensaje := fmt.Sprintf("se produjo el siguiente error: %s.", err.Error())
		logs.Error(mensaje)
		return
	}

	// si no existe la cuenta, el objeto se encuentra vacio
	if cuenta.ID == 0 {
		erro = fmt.Errorf("no se encuentra la cuenta con el id requerido")
		logs.Error(erro.Error())
		return
	}

	// El nombre del cliente para pregunatr en la tabla reportes por ese cliente
	cliente := cuenta.Cliente.Cliente

	resp.Preload("Reportedetalle")

	resp.Where("tiporeporte = ? AND cliente = ? AND fechacobranza = ?", "pagos", cliente, filtro.Fecha)

	resp.Find(&reportes)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetReportesPagoRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) ObtenerCierreLoteRapipago(id int64) (response *entities.Rapipagocierrelote, erro error) {
	resp := r.SQLClient.Model(entities.Rapipagocierrelote{})

	resp.Where("id = ?", id)

	resp.First(&response)
	if resp.Error != nil {
		erro = resp.Error
		return
	}

	return response, erro
}

func (r *repository) checkIfModelsExistsByValue(tableName string, column string, value interface{}) (exists bool, erro error) {

	erro = r.SQLClient.Table(tableName).Select("count(*) > 0").Where(column+" = ?", value).Find(&exists).Error

	if erro != nil {
		erro = errors.New(ERROR_DB_ACCESS)
	}
	return
}

/*
	resp := r.SQLClient.Table("cuentas as c")
	if filtro.IdCuenta > 0 {
		resp.Where("c.id = ? ", filtro.IdCuenta)
	}
	if filtro.CargarTiposPago {
		resp.Joins("inner join pagotipos as ptip on c.id = ptip.cuentas_id")
		resp.Preload("Pagotipos")
	}
	if filtro.CargarPagos {
		resp.Joins("inner join pagos as p on p.pagostipo_id = ptip.id ")
		resp.Preload("Pagotipos.Pagos")
	}
	if filtro.CargarPagosIntentos {
		//resp.Joins("join pagointentos as pt on p.id = pt.pagos_id and pt.state_comment = 'approved' and pt.transaction_id in (?)", filtro.TransactionId)
		resp.Preload("Pagotipos.Pagos.PagoIntentos", "transaction_id in (?)", filtro.TransactionId)
	}
	resp.Find(&pagosEnDisputa)
	if resp.Error != nil {
		erro = errors.New(ERROR_OBTENER_PAGOS_DISPUTA)
	}
*/

func (r *repository) CreateContactosReportesRepository(contactoReportes entities.Contactosreporte) error {

	result := r.SQLClient.Model(entities.Contactosreporte{}).Where(contactoReportes).FirstOrCreate(&contactoReportes)

	if result.Error != nil {
		//Controlo si el error que recibo de la creacion contiene un error de clave foranea, que se produce si no existe un cliente_id valido
		if strings.Contains(result.Error.Error(), "foreign key constraint") {
			return errors.New("el id no corresponde a un cliente")
		}
		return errors.New("error al crear contactoReportes en la base de datos")
	}
	return nil
}
func (r *repository) ReadContactosReportesRepository(contactoReportes entities.Contactosreporte) ([]entities.Contactosreporte, error) {
	var contactosreportes []entities.Contactosreporte
	result := r.SQLClient.Find(&contactosreportes, "clientes_id = ?", contactoReportes.ClientesID)

	if result.Error != nil {
		return nil, errors.New("error al obtener los contactoReportes de la base de datos")
	}

	return contactosreportes, nil
}

func (r *repository) UpdateContactosReportesRepository(contactoReportes, contactoReportesNuevo entities.Contactosreporte) error {
	result := r.SQLClient.Model(entities.Contactosreporte{}).Where("clientes_id = ? AND email = ?", contactoReportes.ClientesID, contactoReportes.Email).Update("email", contactoReportesNuevo.Email)
	if result.Error != nil {
		return errors.New("error al actualizar el correo del cliente")
	}
	if result.RowsAffected == 0 {
		return errors.New("error al actualizar, verifique los campos enviados")
	}
	return nil
}
func (r *repository) GetContactosReportesByIdEmailRepository(contactoReportes entities.Contactosreporte) (entities.Contactosreporte, error) {
	result := r.SQLClient.Find(&contactoReportes, "clientes_id = ? AND email = ? ", contactoReportes.ClientesID, contactoReportes.Email)

	if result.Error != nil {
		return contactoReportes, errors.New("error al obtener los contactoReportes de la base de datos")
	}
	return contactoReportes, nil
}
func (r *repository) DeleteContactosReportesRepository(contactoReportes entities.Contactosreporte) error {

	result := r.SQLClient.Where("clientes_id = ? AND email = ?", contactoReportes.ClientesID, contactoReportes.Email).Delete(&entities.Contactosreporte{})

	if result.Error != nil {
		return errors.New("error al eliminar contactoReportes en la base de datos")
	}
	if result.RowsAffected == 0 {
		return errors.New("no existe un registro a eliminar, verifique los datos enviados")
	}
	return nil
}

func (r *repository) CreateUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error {

	result := r.SQLClient.Model(entities.Usuariobloqueados{}).Create(&usuarioBloqueado)

	if result.Error != nil {
		return errors.New("error al crear usuario bloqueado en la base de datos")
	}
	return nil
}

func (r *repository) UpdateUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error {

	result := r.SQLClient.Model(&usuarioBloqueado).Updates(&usuarioBloqueado)

	if usuarioBloqueado.CantBloqueo == 0 {
		result.Update("cant_bloqueo", 0)
	}

	if usuarioBloqueado.Permanente == false {
		result.Update("permanente", 0)
	}

	if result.Error != nil {
		return errors.New("error al actualizar usuario bloqueado en la base de datos")
	}
	return nil
}

func (r *repository) DeleteUsuarioBloqueadoRepository(usuarioBloqueado entities.Usuariobloqueados) error {

	result := r.SQLClient.Delete(&usuarioBloqueado)

	if result.Error != nil {
		return errors.New("error al eliminar usuario bloqueado en la base de datos")
	}
	return nil
}

func (r *repository) GetUsuariosBloqueadoRepository(filtro filtros.UsuarioBloqueadoFiltro) (bloqueados []entities.Usuariobloqueados, totalBloqueados int64, erro error) {

	result := r.SQLClient.Model(entities.Usuariobloqueados{})

	if filtro.Id > 0 {
		result.Where("id = ?", filtro.Id)
	}

	if len(filtro.Ids) > 0 {
		result.Where("id in (?)", filtro.Ids)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		result.Count(&totalBloqueados)

		if result.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		result.Limit(int(filtro.Size))
		result.Offset(int(offset))

	}

	result.Find(&bloqueados)

	if result.Error != nil {
		erro = errors.New("error al buscar usuarios bloqueados en la base de datos")
		return
	}

	return
}

func (r *repository) GetHistorialOperacionesRepository(filtro filtros.RequestHistorial) (historial []entities.HistorialOperaciones, totalFilas int64, erro error) {
	resp := r.SQLClient.Unscoped().Model(entities.HistorialOperaciones{})

	if filtro.Correo != "" {
		resp.Where("correo = ?", filtro.Correo)
	} else {
		if !filtro.FechaInicio.IsZero() && !filtro.FechaInicio.IsZero() {
			resp.Where("cast(historial_operaciones.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
		}
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf("error al cargar total filas de la columna")
		}
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Order("created_at desc").Find(&historial)

	return
}

func (r *repository) UpsertEnvioRepository(envio entities.Envio) error {
	var envioBuscado entities.Envio
	found := r.SQLClient.Where("clientes_id = ?", envio.ClientesId).First(&envioBuscado)
	if found.RowsAffected == 1 {
		envioBuscado.Cobranzas = envio.Cobranzas
		envioBuscado.Rendiciones = envio.Rendiciones
		envioBuscado.Reversiones = envio.Reversiones
		envioBuscado.Batch = envio.Batch
		envioBuscado.BatchPagos = envio.BatchPagos
		result := r.SQLClient.Omit("created_at", "id").Save(&envioBuscado)
		if result.Error != nil {
			return errors.New("error al actualizar configuracion de envio de archivo")
		}
	}
	if found.RowsAffected == 0 {
		result := r.SQLClient.Omit("id").Create(&envio)
		if result.Error != nil {
			return errors.New("error al actualizar configuracion de envio de archivo")
		}
	}
	return nil
}

// 	//Guardo todos los campos que debo actualizar, para hacerlo todo en una consulta
// 	if soporte.Visto {
// 		updates["visto"] = soporte.Visto
// 	}
// 	if soporte.Estado != "" {
// 		updates["estado"] = soporte.Estado
// 	}
// 	if soporte.Abierto {
// 		updates["abierto"] = soporte.Abierto
// 	}
// 	//Realizo la unica consulta para actualizar todos los campos
// 	if len(updates) > 0 {
// 		result.Updates(updates)
// 	}
// 	if result.Error != nil{
// 		erro = fmt.Errorf("ocurrio un error al actualizar en la base de datos")
// 		logs.Error(erro)
// 		return
// 	}
// 	if result.RowsAffected == 0 {
// 		erro = fmt.Errorf("no se encontro ningun registro a actualizar")
// 		logs.Error(erro)
// 		return
// 	}
// 	return
// }
