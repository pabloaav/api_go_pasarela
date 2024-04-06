package reportes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros_reportes "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/reportes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportesRepository interface {
	GetPagosReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error)
	GetPagosTemporalReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error)
	GetCierreLotePrisma(lista []string) (result []entities.Prismacierrelote, erro error)
	GetCierreLoteOffline(lista []string) (result []entities.Rapipagocierrelotedetalles, erro error)
	GetCierreLoteApilink(lista []string) (result []entities.Apilinkcierrelote, erro error)
	GetCierreLoteApilinkByFechaCobro(filtro reportedtos.RequestPagosPeriodo) (result []entities.Apilinkcierrelote, erro error)

	// GetCierreLoteApilinkReportes(requets filtros_reportes.ReportesFiltroApilink) (result []entities.Apilinkcierrelote, erro error)
	GetMovimiento(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error)
	GetRendicionReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error)
	GetReversionesReportes(request reportedtos.RequestPagosPeriodo, filtroValidacion reportedtos.ValidacionesFiltro) (pagos []entities.Reversione, erro error)
	GetTransferenciasReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Transferencia, erro error)

	GetPeticionesReportes(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error)
	GetPeticionesReportesByOperacion(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error)

	GetLogs(request reportedtos.RequestLogs) (logs []entities.Log, total int64, erro error)
	GetNotificaciones(request reportedtos.RequestNotificaciones) (notificaciones []entities.Notificacione, total int64, erro error)
	SaveLotes(ctx context.Context, lotes []entities.Movimientolotes) (erro error)
	BajaMovimientoLotes(ctx context.Context, movimientos []entities.Movimientolotes, motivo_baja string) error

	/* REPORTES RENTA */
	GetReversionesReportesRenta(request reportedtos.RequestPagosPeriodo) (pagos []entities.Transferencia, erro error)

	//
	GetPagosBatch(request reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error)
	GetLotes(request reportedtos.RequestPagosPeriodo) (lote []entities.Pagolotes, erro error)  // obtener datos del lote
	GetLastLote(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) // Retorna ultimo lote
	GetCantidadLotes(request reportedtos.RequestPagosPeriodo) (lote int64, erro error)
	SavePagosLotes(ctx context.Context, lotes []entities.Pagolotes) (erro error)
	BajaPagosLotes(ctx context.Context, pagos []entities.Pagolotes, motivo_baja string) error

	// generar orden de liquidacion
	SaveLiquidacion(movliquidacion entities.Movimientoliquidaciones) (id uint64, erro error)
	/* REPORTES MOVIMIENTOS-COMISIONES */
	MovimientosComisionesRepository(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error)
	MovimientosComisionesTemporales(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error)

	/* REPORTES COBRANZAS-CLIENTES */
	CobranzasClientesRepository(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)
	CobranzasApilink(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)
	CobranzasRapipago(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)
	CobranzasMultipago(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)

	CobranzasPrisma(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)

	// Guardar datos reportes de clientes : NOTE esto solo se registra para control en el envio
	SaveGuardarDatosReporte(reporte entities.Reporte) (erro error)
	GuardarReportesInfo(reporte []entities.Reporte) (erro error)

	GetLastReporteEnviadosRepository(request entities.Reporte, control bool) (siguiente uint, erro error)
	// Reportes enviados a clientes
	GetReportesEnviadosRepository(request reportedtos.RequestReportesEnviados) (listaReportes []entities.Reporte, totalFilas int64, erro error)
	//Obtener los pagosIntentos de un cierre de lote rapidopago a partir de su barcode
	getCierreLoteRapipago(filtro reportedtos.RequestPagosPeriodo) (rapipagoCierrelotes []entities.Rapipagocierrelotedetalles, erro error)
	getPagosByBarcode(filtro reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error)
	getPagosByExternalPagoIntento(filtro reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error)

	GetCuentaByApiKeyRepository(apikey string) (cuenta *entities.Cuenta, erro error)
}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
	utilService      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService, t util.UtilService) ReportesRepository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
		utilService:      t,
	}
}
func (r *repository) GetCuentaByApiKeyRepository(apikey string) (cuenta *entities.Cuenta, erro error) {
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

func (r *repository) GetPagosReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !request.FechaFin.IsZero() {
		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	// if request.BuscarFechaPagointento {
	// 	if !request.FechaFin.IsZero() {
	// 		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
	// 			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	// 	}
	// }

	resp.Where("pagoestados_id != ?", pendiente)
	resp.Preload("PagoEstados")

	if len(request.PagoEstados) > 0 {
		resp.Where("pagoestados_id IN ?", request.PagoEstados)
	}

	if request.ClienteId > 0 {
		resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	} else if len(request.ApiKey) > 0 {
		resp.Preload("PagosTipo.Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	} else {
		resp.Preload("PagosTipo.Cuenta.Cliente").Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Order("cli.id DESC")
	}
	resp.Preload("PagoIntentos.Mediopagos.Channel.Channelaranceles").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").
		Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO").
		Order("pi.created_at DESC")
	resp.Preload("PagoIntentos.Installmentdetail")
	resp.Preload("PagoIntentos.Movimientos.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientos.Movimientoimpuestos")
	resp.Preload("PagoIntentos.Movimientos.Movimientotransferencia")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPagosTemporalReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !request.FechaFin.IsZero() {
		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin).
			Where("pint.state_comment = ? OR pint.state_comment = ?", "approved", "INICIADO").
			Where("pint.card_last_four_digits != ''").
			Order("pint.created_at DESC")
	}

	// if request.BuscarFechaPagointento {
	// 	if !request.FechaFin.IsZero() {
	// 		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
	// 			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	// 	}
	// }

	if request.FiltroBarcode {
		resp.Where("pint.barcode = ?", "")
	}

	resp.Where("pagoestados_id != ?", pendiente)
	resp.Preload("PagoEstados")

	if len(request.PagoEstados) > 0 {
		resp.Where("pagoestados_id IN ?", request.PagoEstados)
	}

	if request.ClienteId > 0 {
		resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	} else if len(request.ApiKey) > 0 {
		resp.Preload("PagosTipo.Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	} else {
		resp.Preload("PagosTipo.Cuenta.Cliente").Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Order("cli.id DESC")
	}
	resp.Preload("PagoIntentos.Mediopagos.Channel.Channelaranceles")
	/* .Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").
	Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO").
	Order("pi.created_at DESC") */
	resp.Preload("PagoIntentos.Installmentdetail")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientoimpuestos")
	// resp.Preload("PagoIntentos.Movimientos.Movimientotransferencia")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTemporalReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCierreLotePrisma(lista []string) (result []entities.Prismacierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Prismacierrelote{})

	resp.Unscoped()

	resp.Where("externalcliente_id IN ?", lista)

	resp.Preload("Prismamovimientodetalle.MovimientoCabecera")

	resp.Preload("Prismatrdospagos")
	resp.Order("fecha_pago")
	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_PRISMA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLotePrisma",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCierreLoteOffline(lista []string) (result []entities.Rapipagocierrelotedetalles, erro error) {

	// resp := r.SQLClient.Model(entities.Rapipagocierrelotedetalles{}).Where("codigo_barras IN ?", lista)

	resp := r.SQLClient.Unscoped().Where("codigo_barras IN ?", lista)

	resp.Preload("RapipagoCabecera", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})

	// resp.Preload("RapipagoCabecera")

	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_OFFLINE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteOffline",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

// func (r *repository) GetCierreLoteApilinkReportes(requets filtros_reportes.ReportesFiltroApilink) (result []entities.Apilinkcierrelote, erro error) {

// 	resp := r.SQLClient.Unscoped()

// 	if requets.Conciliado {
// 		resp.Where("banco_external_id != ?", 0)
// 	}

// 	if !requets.Informado {
// 		resp.Where("pagoinformado = ?", requets.Informado)
// 	}

// 	resp.Find(&result)

// 	if resp.Error != nil {

// 		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_APILINK)

// 		log := entities.Log{
// 			Tipo:          entities.Error,
// 			Mensaje:       resp.Error.Error(),
// 			Funcionalidad: "GetCierreLoteApilinkReportes",
// 		}

// 		err := r.utilService.CreateLogService(log)

// 		if err != nil {
// 			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
// 			logs.Error(mensaje)
// 		}
// 	}

// 	return
// }

func (r *repository) GetCierreLoteApilink(lista []string) (result []entities.Apilinkcierrelote, erro error) {

	resp := r.SQLClient.Unscoped().Where("debin_id IN ?", lista)

	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_APILINK)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteApilink",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
func (r *repository) GetCierreLoteApilinkByFechaCobro(filtro reportedtos.RequestPagosPeriodo) (result []entities.Apilinkcierrelote, erro error) {

	resp := r.SQLClient.Unscoped().Where("estado = ? AND cast(fecha_cobro as date) BETWEEN cast(? as date) AND cast(? as date)", "ACREDITADO", filtro.FechaInicio, filtro.FechaFin)

	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_APILINK)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteApilinkByFechaCobro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
func (r *repository) GetRendicionReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if request.CuentaId > 0 {
		resp.Preload("Cuenta", "cuentas.id = ?", request.CuentaId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id").Where("c.id = ?", request.CuentaId)
	}

	if request.ClienteId > 0 {
		resp.Preload("Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	}

	if request.PagoIntento > 0 {
		resp.Where("pagointentos_id = ?", request.PagoIntento)
	}

	if len(request.TipoMovimiento) > 0 {
		resp.Where("tipo = ?", request.TipoMovimiento)
	}

	if request.CargarMedioPago {
		resp.Preload("Pagointentos.Mediopagos")
	}

	resp.Preload("Pagointentos.Pago.Pagoitems")
	resp.Preload("Movimientotransferencia")
	resp.Preload("Movimientocomisions")
	resp.Preload("Movimientoimpuestos")
	// resp.Preload("Movimientolotes")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
func (r *repository) GetReversionesReportesRenta(request reportedtos.RequestPagosPeriodo) (transferencias []entities.Transferencia, erro error) {
	resp := r.SQLClient.Model(entities.Transferencia{})

	if !request.FechaInicio.IsZero() {
		resp.Where("cast(transferencias.fecha_operacion as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.CargarReversion {
		resp.Where("transferencias.reversion = ?", 1)
	}

	resp.Preload("Movimiento").Joins("INNER JOIN movimientos as M on M.id = transferencias.movimientos_id").
		Joins("INNER JOIN cuentas as C on C.id = M.cuentas_id").Where("C.apikey = ?", request.ApiKey)

	resp.Find(&transferencias)

	return
}

func (r *repository) GetReversionesReportes(request reportedtos.RequestPagosPeriodo, filtroValidacion reportedtos.ValidacionesFiltro) (pagos []entities.Reversione, erro error) {

	resp := r.SQLClient.Model(entities.Reversione{})

	if filtroValidacion.Fin && filtroValidacion.Inicio {
		resp.Where("cast(reversiones.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.ClienteId > 0 && request.CuentaId > 0 {
		resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta", "cuentas.clientes_id = ? AND cuentas.id = ?", request.ClienteId, request.CuentaId).
			Joins("INNER JOIN pagointentos as pt on pt.id = reversiones.pagointentos_id INNER JOIN pagos as p on p.id = pt.pagos_id INNER JOIN pagotipos as pa on pa.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pa.cuentas_id").
			Where("cu.clientes_id = ? AND cu.id = ?", request.ClienteId, request.CuentaId)
	} else if request.ClienteId > 0 {
		resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta", "cuentas.clientes_id = ?", request.ClienteId).
			Joins("INNER JOIN pagointentos as pt on pt.id = reversiones.pagointentos_id INNER JOIN pagos as p on p.id = pt.pagos_id INNER JOIN pagotipos as pa on pa.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pa.cuentas_id").
			Where("cu.clientes_id = ?", request.ClienteId)
	} else if request.CuentaId > 0 {
		resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta", "cuentas.id = ?", request.CuentaId).
			Joins("INNER JOIN pagointentos as pt on pt.id = reversiones.pagointentos_id INNER JOIN pagos as p on p.id = pt.pagos_id INNER JOIN pagotipos as pa on pa.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pa.cuentas_id").
			Where("cu.id = ?", request.CuentaId)
	}

	resp.Preload("PagoIntento.Mediopagos")
	resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta.Cliente")
	resp.Preload("PagoIntento.Pago.Pagoitems")
	resp.Preload("PagoIntento.Pago.PagoEstados")

	if request.OrdenadoFecha {
		resp.Order("id DESC")
	}

	resp.Find(&pagos)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetReversionesReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetTransferenciasReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Transferencia, erro error) {

	resp := r.SQLClient.Model(entities.Transferencia{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(transferencias.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.ClienteId > 0 && request.CuentaId > 0 {
		resp.Preload("Movimiento.Cuenta", "cuentas.id = ?", request.CuentaId).
			Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id").
			Joins("INNER JOIN cuentas as c on c.id = mv.cuentas_id").
			Where("c.id = ? AND c.clientes_id = ?", request.CuentaId, request.ClienteId)
	} else {
		if request.ClienteId > 0 {
			resp.Preload("Movimiento.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id INNER JOIN cuentas as c on c.id = mv.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
		}
		if request.CuentaId > 0 {
			resp.Preload("Movimiento.Cuenta", "cuentas.id = ?", request.CuentaId).Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id INNER JOIN cuentas as c on c.id = mv.cuentas_id").Where("c.id = ?", request.CuentaId)
		}
	}

	if len(request.ApiKey) > 0 {
		resp.Preload("Movimiento.Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id INNER JOIN cuentas as c on c.id = mv.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	}
	// resp.Preload("Movimiento.Pagointentos.Pago.Pagoitems")
	// resp.Preload("Movimiento.Movimientocomisions")
	// resp.Preload("Movimiento.Movimientoimpuestos")
	// resp.Preload("Movimientolotes")

	if request.OrdenadoFecha {
		resp.Order("fecha_operacion DESC")
	}

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPeticionesReportes(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Webservicespeticione{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	resp.Where("vendor = ?", request.Vendor)
	resp.Where("operacion != ?", "Autenticacion(genera token)")
	resp.Order("created_at desc")

	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&peticiones)

	return
}
func (r *repository) GetPeticionesReportesByOperacion(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Webservicespeticione{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	resp.Where("vendor = ?", request.Vendor)
	resp.Where("operacion = ?", request.Operacion)
	resp.Order("created_at desc")
	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&peticiones)

	return
}
func (r *repository) GetLogs(request reportedtos.RequestLogs) (logs []entities.Log, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Log{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	resp.Order("created_at DESC")
	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&logs)

	return
}
func (r *repository) GetNotificaciones(request reportedtos.RequestNotificaciones) (notificaciones []entities.Notificacione, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Notificacione{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	resp.Order("created_at DESC")
	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&notificaciones)

	return
}

func (r *repository) GetLastLote(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) {

	resp := r.SQLClient.Model(entities.Pagolotes{})

	if request.ClienteId > 0 {
		resp.Where("clientes_id = ?", request.ClienteId)
	}
	resp.Last(&lote)

	if resp.RowsAffected <= 0 {
		lote = entities.Pagolotes{}
	}
	return
}

func (r *repository) SaveLotes(ctx context.Context, lotes []entities.Movimientolotes) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		// 1 - creo los movimientos lotes
		if len(lotes) > 0 {
			res := tx.WithContext(ctx).Create(&lotes)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_GUARDAR_LOTES)
			}
		}
		return nil
	})
}

func (r *repository) BajaMovimientoLotes(ctx context.Context, movimientos []entities.Movimientolotes, motivo_baja string) error {

	resp := r.SQLClient.WithContext(ctx).Model(&movimientos).Omit(clause.Associations).UpdateColumns(map[string]interface{}{"updated_at": time.Now(), "deleted_at": time.Now(), "motivo_baja": motivo_baja})

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro := fmt.Errorf(ERROR_BAJAR_MOVIMIENTOS_LOTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "BajaMovimientoLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
		return erro
	}
	return nil
}

func (r *repository) GetMovimiento(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if request.PagoIntento > 0 {
		resp.Where("pagointentos_id = ?", request.PagoIntento)
	}

	if len(request.PagoIntentos) > 0 {
		resp.Where("pagointentos_id IN ?", request.PagoIntentos)
	}

	if len(request.ApiKey) > 0 {
		resp.Preload("Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	}

	if request.CargarMedioPago {
		resp.Preload("Pagointentos.Mediopagos.Channel")
	}

	if len(request.TipoMovimiento) > 0 {
		resp.Where("tipo = ?", request.TipoMovimiento)
	}

	if request.CargarReversionReporte {
		resp.Where("reversion = ?", true)
		resp.Where("monto < ?", 0)
	}

	if request.CargarReversion {
		resp.Where("reversion = ?", true)
		resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if request.CargarCuenta {
		resp.Preload("Pagointentos.Pago.Pagoitems")
	}

	if request.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if request.CargarComisionImpuesto {
		resp.Preload("Movimientocomisions")
		resp.Preload("Movimientoimpuestos")
	}

	if request.CargarMovimientosTransferencias {
		resp.Preload("Movimientotransferencia")
	}

	if request.CargarCliente {
		resp.Preload("Cuenta.Cliente")
	}

	if request.CargarRetenciones {
		resp.Preload("Movimientoretencions.Retencion.Condicion.Gravamen")
	}

	if request.OrdenadoFecha {
		resp.Order("pagointentos_id Desc")
	}

	resp.Find(&pagos)

	return
}

// pagos items batch
func (r *repository) GetPagosBatch(request reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	// se cambio la busqueda al campo paid_at de pago intento
	// if !request.FechaFin.IsZero() {
	// 	resp.Where("cast(pagos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	// }

	if !request.FechaFin.IsZero() {
		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if len(request.PagoEstados) > 0 {
		resp.Where("pagoestados_id IN ?", request.PagoEstados)
	}

	if request.ClienteId > 0 {
		// resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
		resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", request.ClienteId)
	}
	// resp.Preload("PagoIntentos")
	// buscar solo sobre el pago intento aprobado
	resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").
		Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")

	resp.Preload("Pagoitems")
	resp.Preload("Pagolotes")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCantidadLotes(request reportedtos.RequestPagosPeriodo) (lote int64, erro error) {

	var pagoslotes []entities.Pagolotes
	resp := r.SQLClient.Model(entities.Pagolotes{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(pagolotes.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if len(request.Pagos) > 0 {
		resp.Where("pagos_id IN ?", request.Pagos)
	}

	// if request.ClienteId > 0 {
	// 	// resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	// 	resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", request.ClienteId)
	// }

	resp.Find(&pagoslotes)

	lote = resp.RowsAffected

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

// func (r *repository) GetPagosLotes(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) {

// 	resp := r.SQLClient.Model(entities.Pagolotes{})

// 	if request.ClienteId > 0 {
// 		resp.Where("clientes_id = ?", request.ClienteId)
// 	}
// 	resp.Last(&lote)

// 	if resp.RowsAffected <= 0 {
// 		lote = entities.Pagolotes{}
// 	}
// 	return
// }

func (r *repository) GetLotes(request reportedtos.RequestPagosPeriodo) (lote []entities.Pagolotes, erro error) {

	resp := r.SQLClient.Model(entities.Pagolotes{})

	if len(request.Pagos) > 0 {
		resp.Where("pagos_id = ?", request.Pagos)
	}

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	resp.Find(&lote)
	return
}

func (r *repository) SavePagosLotes(ctx context.Context, lotes []entities.Pagolotes) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		// 1 - creo los movimientos lotes
		if len(lotes) > 0 {
			res := tx.WithContext(ctx).Create(&lotes)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_GUARDAR_LOTES)
			}
		}
		return nil
	})
}

func (r *repository) BajaPagosLotes(ctx context.Context, pagos []entities.Pagolotes, motivo_baja string) error {

	resp := r.SQLClient.WithContext(ctx).Model(&pagos).Omit(clause.Associations).UpdateColumns(map[string]interface{}{"updated_at": time.Now(), "deleted_at": time.Now(), "motivo_baja": motivo_baja})

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro := fmt.Errorf(ERROR_BAJAR_MOVIMIENTOS_LOTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "BajaPagosLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
		return erro
	}
	return nil
}

func (r *repository) SaveLiquidacion(movliquidacion entities.Movimientoliquidaciones) (id uint64, erro error) {

	result := r.SQLClient.Create(&movliquidacion)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_MOVIMIENTOS_LIQUIDACIONES)
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
	id = uint64(movliquidacion.ID)

	return
}
func (r *repository) MovimientosComisionesRepository(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	resp.Joins("INNER JOIN pagointentos PI ON movimientos.pagointentos_id = PI.id INNER JOIN pagos P ON PI.pagos_id = P.id INNER JOIN movimientoimpuestos MI ON movimientos.id = MI.movimientos_id INNER JOIN movimientocomisiones AS MC ON movimientos.id = MC.movimientos_id INNER JOIN cuentacomisions CC ON MC.cuentacomisions_id = CC.id INNER JOIN cuentas C ON C.id = movimientos.cuentas_id INNER JOIN clientes CTS ON C.clientes_id = CTS.id INNER JOIN transferenciacomisiones TR ON TR.movimientos_id = movimientos.id ")
	resp.Select("movimientos.id AS id,case movimientos.tipo when 'C' then 'CREDITO' else 'DEBITO' end AS tipo, PI.amount AS monto_pago, movimientos.monto AS monto_movimiento, (MC.monto + MC.montoproveedor) AS monto_comision, (MC.porcentaje + MC.porcentajeproveedor) AS porcentaje_comision , (MI.monto + MI.montoproveedor) AS monto_impuesto, MI.porcentaje AS porcentaje_impuesto, MC.montoproveedor AS monto_comisionproveedor,MC.porcentajeproveedor AS porcentaje_comisionproveedor , MI.montoproveedor AS monto_impuestoproveedor, CTS.cliente AS nombre_cliente, movimientos.created_at,TR.fecha_operacion,PI.paid_at as fecha_pago")

	// resp.Where("movimientos.created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	if filtro.UsarFechaPago {
		resp.Where("PI.paid_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	} else {

		resp.Where("TR.fecha_operacion BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	}

	if filtro.ClienteId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
		if filtro.CuentaId != 0 {
			resp.Where("C.id = ?", filtro.CuentaId)
		}
	}

	if filtro.SoloReversiones {
		resp.Where("TR.reversion = ?", true)
	}

	resp.Order("movimientos.id DESC")

	resp.Find(&total)
	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_MOVIMIENTOS_COMISIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "MovimientosComisionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) MovimientosComisionesTemporales(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error) {

	resp := r.SQLClient.Model(entities.Movimientotemporale{})

	resp.Joins("INNER JOIN pagointentos PI ON movimientotemporales.pagointentos_id = PI.id INNER JOIN pagos P ON PI.pagos_id = P.id INNER JOIN movimientoimpuestotemporales MI ON movimientotemporales.id = MI.movimientotemporales_id INNER JOIN movimientocomisionetemporales AS MC ON movimientotemporales.id = MC.movimientotemporales_id INNER JOIN cuentacomisions CC ON MC.cuentacomisions_id = CC.id INNER JOIN cuentas C ON C.id = movimientotemporales.cuentas_id INNER JOIN clientes CTS ON C.clientes_id = CTS.id")
	resp.Select("movimientotemporales.id AS id,case movimientotemporales.tipo when 'C' then 'CREDITO' else 'DEBITO' end AS tipo, PI.amount AS monto_pago, movimientotemporales.monto AS monto_movimiento, (MC.monto + MC.montoproveedor) AS monto_comision, (MC.porcentaje + MC.porcentajeproveedor) AS porcentaje_comision , (MI.monto + MI.montoproveedor) AS monto_impuesto, MI.porcentaje AS porcentaje_impuesto, MC.montoproveedor AS monto_comisionproveedor,MC.porcentajeproveedor AS porcentaje_comisionproveedor , MI.montoproveedor AS monto_impuestoproveedor, CTS.cliente AS nombre_cliente, movimientotemporales.created_at,PI.paid_at as fecha_pago")

	// resp.Where("movimientos.created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	if filtro.UsarFechaPago {
		resp.Where("PI.paid_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	} else {

		resp.Where("TR.fecha_operacion BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	}

	if filtro.ClienteId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
		if filtro.CuentaId != 0 {
			resp.Where("C.id = ?", filtro.CuentaId)
		}
	}

	resp.Order("movimientotemporales.id DESC")

	resp.Find(&total)
	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_MOVIMIENTOS_COMISIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "MovimientosComisionesTemporales",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CobranzasClientesRepository(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})
	//var fechaCobro sql.NullString
	resp.Joins("INNER JOIN pagotipos PT ON pagos.pagostipo_id = PT.id")
	resp.Joins("INNER JOIN cuentas C ON C.id = PT.cuentas_id")
	resp.Joins("INNER JOIN clientes CTS ON C.clientes_id = CTS.id")
	resp.Joins("INNER JOIN pagoestados PE ON pagos.pagoestados_id = PE.id")
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")
	resp.Joins("INNER JOIN mediopagos MP ON MP.id = PI.mediopagos_id")
	resp.Joins("INNER JOIN channels CH ON CH.id = MP.channels_id")

	querySelect := "pagos.id, pagos.payer_name, pagos.payer_email, pagos.first_total as total_pago, pagos.external_reference as referencia , PE.nombre as Pagoestado ,pagos.description AS descripcion , C.cuenta, CTS.cliente, cast(PI.paid_at as date)  as fecha_pago, MP.mediopago as medio_pago, CH.nombre as canal_pago"

	//Si existen ambos filtros, se deben obtener la fecha de cobro de rapipago O apilink, entonces COALESCE nos sirve para obtener la primera fecha que no sea nula, como nunca deberian existir ambas.
	if filtro.ObtenerBarcodes && filtro.ObtenerApiLinkByFechaCobro && filtro.ObtenerPrismaByFechaOperacion {
		// // relacion entre prismacierrelotes y pagointento
		resp.Joins("LEFT JOIN prismacierrelotes PCL ON PI.card_last_four_digits IS NOT NULL AND PCL.externalcliente_id = PI.transaction_id")

		// Cuando exista el codigo de barra se debe consultar la tabla de rapipagocierre de lotes
		resp.Joins("LEFT JOIN rapipagocierrelotedetalles RP ON PI.barcode IS NOT NULL AND RP.codigo_barras = PI.barcode")
		//Busca si coincide el debin_id  y el external_id
		resp.Unscoped().Joins("LEFT JOIN apilinkcierrelotes ACL ON PI.external_id IS NOT NULL AND ACL.debin_id = PI.external_id")

		// select de fechas
		querySelect = fmt.Sprintf("%s,  COALESCE(ACL.fecha_cobro, RP.fecha_cobro, PCL.fechaoperacion ) AS fecha_cobro", querySelect)
	} else {
		//Si solo de activa uno  de los dos filtros, simplemente toma la fecha de cobro de la tabla correspondiente
		// Cuando exista el codigo de barra se debe consultar la tabla de rapipagocierre de lotes
		if filtro.ObtenerBarcodes {
			resp.Joins("LEFT JOIN rapipagocierrelotedetalles RP ON PI.barcode IS NOT NULL AND RP.codigo_barras = PI.barcode")
			querySelect = fmt.Sprintf("%s, RP.fecha_cobro", querySelect)
		}
		if filtro.ObtenerApiLinkByFechaCobro {
			resp.Joins("LEFT JOIN apilinkcierrelotes ACL ON PI.external_id IS NOT NULL AND ACL.debin_id = PI.external_id")
			querySelect = fmt.Sprintf("%s, ACL.fecha_cobro", querySelect)
		}
	}
	resp.Select(querySelect)
	if filtro.FiltrarPorFechaCobro {
		resp.Where("fecha_cobro BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	} else {
		// resp.Where("pagos.created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
		resp.Where("fecha_cobro BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	}
	resp.Where("PE.nombre = ? OR PE.nombre = ?", "AUTORIZADO", "APROBADO")
	resp.Where("PI.state_comment = ? OR PI.state_comment = ?", "approved", "INICIADO")

	// resp.Where("PI.id in ( SELECT  MAX(id) FROM pagointentos as PI GROUP BY PI.pagos_id	)")

	if filtro.ClienteId != 0 && filtro.CuentaId != 0 {
		// Consulta con filtro estricto por ClienteId y CuentaId
		resp.Where("CTS.id = ?", filtro.ClienteId)
		resp.Where("C.id = ?", filtro.CuentaId)
	} else {
		if filtro.ClienteId != 0 {
			resp.Where("CTS.id = ?", filtro.ClienteId)
		}
		if filtro.CuentaId != 0 {
			resp.Joins("INNER JOIN cuentas Cu ON Cu.id = PT.cuentas_id")
			resp.Where("Cu.id = ?", filtro.CuentaId)
		}
	}

	resp.Order("pagos.id DESC")

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_COBRANZAS_CLIENTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CobranzasClientesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) CobranzasApilink(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(filtro, "INNER JOIN apilinkcierrelotes ACL ON ACL.debin_id = PI.external_id", "ACL.fecha_cobro")

	resp.Find(&response)

	return
}

func (r *repository) CobranzasRapipago(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(filtro, "INNER JOIN rapipagocierrelotedetalles RP ON RP.codigo_barras = PI.barcode", "RP.fecha_cobro")
	resp.Find(&response)

	return
}

func (r *repository) CobranzasMultipago(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(filtro, "INNER JOIN multipagoscierrelotedetalles MTD ON MTD.codigo_barras = PI.barcode", "MTD.fecha_cobro")
	resp.Find(&response)

	return
}

func (r *repository) CobranzasPrisma(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(filtro, "", "PI.paid_at")
	resp.Where("PI.card_last_four_digits != '' AND PI.paid_at IS NOT NULL")
	resp.Find(&response)

	return
}

func (r *repository) constructQuery(filtro filtros_reportes.CobranzasClienteFiltro, joinTable string, fechaCobroColumn string) *gorm.DB {
	resp := r.SQLClient.Model(entities.Pago{})

	resp.Joins("INNER JOIN pagotipos PT ON pagos.pagostipo_id = PT.id")
	resp.Joins("INNER JOIN cuentas C ON C.id = PT.cuentas_id")
	resp.Joins("INNER JOIN clientes CTS ON C.clientes_id = CTS.id")
	resp.Joins("INNER JOIN pagoestados PE ON pagos.pagoestados_id = PE.id")
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")
	resp.Joins("INNER JOIN mediopagos MP ON MP.id = PI.mediopagos_id")
	resp.Joins("INNER JOIN channels CH ON CH.id = MP.channels_id")

	resp.Joins("LEFT JOIN movimientotemporales MT ON MT.pagointentos_id = PI.id")
	resp.Joins("LEFT JOIN movimientocomisionetemporales MCT ON MCT.movimientotemporales_id = MT.id")
	resp.Joins("LEFT JOIN movimientoimpuestotemporales MIT ON MIT.movimientotemporales_id = MT.id")
	resp.Joins("LEFT JOIN movimiento_retenciontemporales MRT ON MRT.movimientotemporales_id = MT.id")

	if len(joinTable) > 0 {
		resp.Joins(joinTable)
	}

	querySelect := `
		pagos.id, 
		pagos.payer_name, 
		pagos.payer_email, 
		pagos.first_total as total_pago, 
		pagos.external_reference as referencia , 
		PE.nombre as Pagoestado ,
		pagos.description AS descripcion , 
		C.cuenta, 
		CTS.cliente, 
		cast(PI.paid_at as date)  as fecha_pago, 
		MP.mediopago as medio_pago, 
		CH.nombre as canal_pago, 
		(MCT.monto + MCT.montoproveedor) AS comision, 
		(MIT.monto + MIT.montoproveedor) AS iva, 
		SUM(MRT.importe_retenido) AS retencion, 
		` + fechaCobroColumn

	resp.Select(querySelect)
	resp.Where(fechaCobroColumn+" BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	resp.Group("PI.id")

	if filtro.ClienteId != 0 && filtro.CuentaId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
		resp.Where("C.id = ?", filtro.CuentaId)
	} else {
		if filtro.ClienteId != 0 {
			resp.Where("CTS.id = ?", filtro.ClienteId)
		}
		if filtro.CuentaId != 0 {
			resp.Joins("INNER JOIN cuentas Cu ON Cu.id = PT.cuentas_id")
			resp.Where("Cu.id = ?", filtro.CuentaId)
		}
	}

	return resp
}

func (r *repository) SaveGuardarDatosReporte(reporte entities.Reporte) (erro error) {

	result := r.SQLClient.Create(&reporte)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_REGISTRO_REPORTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "SaveGuardarDatosReporte",
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

func (r *repository) GuardarReportesInfo(reportes []entities.Reporte) (erro error) {

	r.SQLClient.Transaction(func(tx *gorm.DB) error {

		for _, reporte := range reportes {
			if err := tx.Updates(&reporte).Error; err != nil {
				// return any error will rollback
				return err
			}
		}

		// return nil will commit the whole transaction
		return nil
	})

	return
}

func (r *repository) GetReportesEnviadosRepository(request reportedtos.RequestReportesEnviados) (listaReportes []entities.Reporte, totalFilas int64, erro error) {

	queryGorm := r.SQLClient.Model(entities.Reporte{})

	if request.FechaInicio != "" && request.FechaFin != "" {
		queryGorm.Unscoped()
		queryGorm.Where("created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin)
	}

	if len(request.Fecharendicion) > 0 {
		queryGorm.Where("reportes.fecharendicion IN ?", request.Fecharendicion)
	}

	// se checkea el filtro TipoReporte
	if request.TipoReporte != "todos" {
		queryGorm.Where("tiporeporte = ?", request.TipoReporte)
	}

	// filtro por cliente
	if len(request.Cliente) != 0 {
		queryGorm.Where("cliente LIKE ?", "%"+request.Cliente+"%")
	}

	// Paginacion
	if request.Number > 0 && request.Size > 0 {

		// Ejecutar y contar las filas devueltas
		queryGorm.Count(&totalFilas)

		if queryGorm.Error != nil {
			erro = fmt.Errorf("no se pudo cargar el total de filas de la consulta")
			return
		}

		offset := (request.Number - 1) * request.Size
		queryGorm.Limit(int(request.Size))
		queryGorm.Offset(int(offset))
	}

	if request.SinNumero {
		queryGorm.Where("nro_reporte = 0")
	}

	// Filtro para enumerar Reportes original agrupando ; sino funcionamiento normal
	if request.Enum {
		queryGorm.Group("cliente")
		if request.TipoReporte == "pagos" {
			queryGorm.Group("fechacobranza")
		}
		if request.TipoReporte == "rendiciones" {
			queryGorm.Group("fecharendicion")
		}
		queryGorm.Order("created_at asc")

	} else {
		queryGorm.Order("created_at desc")
		queryGorm.Preload("Reportedetalle")
	}
	queryGorm.Find(&listaReportes)

	// capturar error query DB
	if queryGorm.Error != nil {

		erro = fmt.Errorf("repositorio: no se puedieron obtener los registros de reportes enviados")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetReportesEnviadosRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetLastReporteEnviadosRepository(request entities.Reporte, control bool) (siguiente uint, erro error) {

	var controlEntity []entities.Reporte
	queryGorm := r.SQLClient.Model(entities.Reporte{})

	queryGorm.Where("tiporeporte = ?", request.Tiporeporte)

	if len(request.Cliente) != 0 {
		queryGorm.Where("cliente = ?", request.Cliente)
	}

	if request.Fechacobranza != "" {
		queryGorm.Where("fechacobranza = ?", request.Fechacobranza)
	}

	if request.Fecharendicion != "" {
		queryGorm.Where("fecharendicion = ?", request.Fecharendicion)
	}

	// Busca el último si no encuentra coincidencia

	queryGorm.Order("created_at asc")

	queryGorm.Find(&controlEntity)

	if len(controlEntity) > 0 {
		siguiente = controlEntity[0].Nro_reporte

		return
	}

	if control {
		return
	}

	var lastReporte entities.Reporte
	queryLast := r.SQLClient.Model(entities.Reporte{})

	queryLast.Where("tiporeporte = ?", request.Tiporeporte)

	// if len(request.Cliente) != 0 {
	// 	queryGorm.Where("cliente LIKE ?", "%"+request.Cliente+"%")
	// }

	queryLast.Order("created_at desc")

	queryLast.First(&lastReporte)

	siguiente = lastReporte.Nro_reporte + 1

	return

}

func (r *repository) getCierreLoteRapipago(filtro reportedtos.RequestPagosPeriodo) (rapipagoCierrelotes []entities.Rapipagocierrelotedetalles, erro error) {
	resp := r.SQLClient.Model(entities.Rapipagocierrelotedetalles{})

	resp.Where("fecha_cobro BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	// resp.Where("cast(rapipagocierrelotedetalles.fecha_cobro as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	resp.Unscoped().Find(&rapipagoCierrelotes)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}

func (r *repository) getPagosByBarcode(filtro reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error) {
	resp := r.SQLClient.Model(entities.Pago{})
	resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
		Where("barcode IN ?", filtro.Barcodes).
		Where("pint.state_comment = ? OR pint.state_comment = ?", "approved", "INICIADO").
		Order("pint.created_at DESC")
	if len(filtro.ApiKey) > 0 {
		resp.Preload("PagosTipo.Cuenta", "cuentas.apikey = ?", filtro.ApiKey).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.apikey = ?", filtro.ApiKey)
	}
	resp.Preload("PagoIntentos.Mediopagos.Channel.Channelaranceles")
	resp.Preload("PagoIntentos.Installmentdetail")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientoimpuestos")
	resp.Find(&pagos)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}

func (r *repository) getPagosByExternalPagoIntento(filtro reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error) {
	resp := r.SQLClient.Model(entities.Pago{})
	resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
		Where("external_id IN ?", filtro.ExternalIDs).
		Where("pint.state_comment = ? OR pint.state_comment = ?", "approved", "INICIADO").
		Order("pint.created_at DESC")
	if len(filtro.ApiKey) > 0 {
		resp.Preload("PagosTipo.Cuenta", "cuentas.apikey = ?", filtro.ApiKey).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.apikey = ?", filtro.ApiKey)
	}
	resp.Preload("PagoIntentos.Mediopagos.Channel.Channelaranceles")
	resp.Preload("PagoIntentos.Installmentdetail")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientoimpuestos")
	resp.Find(&pagos)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}
