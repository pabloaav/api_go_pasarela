package administracion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *repository) CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error) {

	res := r.SQLClient.Omit(clause.Associations).Create(&cierreLotes)

	if res.RowsAffected < 1 {
		logs.Error(res.Error.Error())

		erro = fmt.Errorf(ERROR_CREAR_CIERRE_LOTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       res.Error.Error(),
			Funcionalidad: "CreateCierreLoteApiLink",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}

	return
}

func (r *repository) CreateTransferencias(ctx context.Context, transferencias []*entities.Transferencia) (erro error) {

	resp := r.SQLClient.WithContext(ctx).Create(&transferencias)

	if resp.Error != nil {

		logs.Error(resp.Error.Error())

		erro = fmt.Errorf(ERROR_CREAR_TRANSFERENCIAS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateTransferencias",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}

		return
	}

	erro = r.auditarAdministracion(resp.Statement.Context, transferencias)
	if erro != nil {
		return erro
	}
	return nil
}

func (r *repository) CreateTransferenciasComisiones(ctx context.Context, transferencias []*entities.Transferenciacomisiones) (erro error) {

	resp := r.SQLClient.WithContext(ctx).Create(&transferencias)

	if resp.Error != nil {

		logs.Error(resp.Error.Error())

		erro = fmt.Errorf(ERROR_CREAR_TRANSFERENCIAS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateTransferenciasTelco",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}

		return
	}

	erro = r.auditarAdministracion(resp.Statement.Context, transferencias)
	if erro != nil {
		return erro
	}
	return nil
}

func (r *repository) GetTransferencias(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferencia, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Transferencia{})

	if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
		resp.Where("cast(transferencias.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}

	if len(filtro.ReferenciaBancaria) > 0 {
		resp.Where("id referencia_bancaria = ?", filtro.ReferenciaBancaria)
	}

	if len(filtro.MovimientosIds) > 0 {
		resp.Where("movimientos_id IN ?", filtro.MovimientosIds)
	}

	if filtro.CuentaId > 0 {
		resp.Joins("JOIN movimientos on movimientos.cuentas_id = ? and movimientos.id = transferencias.movimientos_id", filtro.CuentaId)
	}

	resp.Order("fecha_operacion DESC")

	resp.Preload("Movimiento.Pagointentos.Pago.PagosTipo")
	resp.Preload("Movimiento.Pagointentos.Mediopagos.Channel")

	if filtro.CargarPaginado {
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
	}

	resp.Find(&transferencias)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(ERROR_CARGAR_TRANSFERENCIA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetTransferencias",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}
	return
}

func (r *repository) GetTransferenciasComisiones(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferenciacomisiones, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Transferenciacomisiones{})

	if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
		resp.Where("cast(transferencias.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}

	if len(filtro.ReferenciaBancaria) > 0 {
		resp.Where("id referencia_bancaria = ?", filtro.ReferenciaBancaria)
	}

	if len(filtro.MovimientosIds) > 0 {
		resp.Where("movimientos_id IN ?", filtro.MovimientosIds)
	}

	if filtro.CargarPaginado {
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
	}

	resp.Find(&transferencias)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(ERROR_CARGAR_TRANSFERENCIA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetTransferencias",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}
	return
}

func (r *repository) GetMovimientos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if len(filtro.Ids) > 0 {
		resp.Where("movimientos.id IN ?", filtro.Ids)
	}

	if filtro.CuentaId > 0 {
		resp.Where("movimientos.cuentas_id", filtro.CuentaId)
	}

	if filtro.CargarPagoIntentos {
		resp.Preload("Pagointentos")
	}

	if filtro.CargarMovimientoSubcuenta {
		resp.Preload("Movimientosubcuentas.Subcuenta")
	}

	if filtro.CargarPago {
		if len(filtro.Concepto) > 0 || len(filtro.ReferenciaPago) > 0 {
			resp.Joins("INNER JOIN pagointentos as pi on movimientos.pagointentos_id = pi.id INNER JOIN pagos as p on pi.pagos_id = p.id INNER JOIN pagotipos as pt on p.pagostipo_id = pt.id")
			if len(filtro.ReferenciaPago) > 0 {
				resp.Where("p.external_reference LIKE ?", "%"+filtro.ReferenciaPago+"%")
			}
			if len(filtro.Concepto) > 0 {
				resp.Where("pt.pagotipo LIKE ?", "%"+filtro.Concepto+"%")
			}
		}
		resp.Preload("Pagointentos.Pago.PagosTipo")

	}

	if filtro.CargarPagoEstados {
		resp.Preload("Pagointentos.Pago.PagoEstados")
	}

	if filtro.CargarMedioPago {
		if filtro.MedioPagoId > 0 {
			// NOTE: 31-10-2022 SE MODIFICO LA CONSULTA
			// resp.Joins("INNER JOIN pagointentos as pi on movimientos.pagointentos_id = pi.id INNER JOIN mediopagos as mp on mp.id = pi.mediopagos_id").Where("mp.id = ?", filtro.MedioPagoId)
			resp.Joins("INNER JOIN pagointentos as pi on movimientos.pagointentos_id = pi.id INNER JOIN mediopagos as mp on mp.id = pi.mediopagos_id").Where("mp.channels_id = ?", filtro.MedioPagoId)
		}
		resp.Preload("Pagointentos.Mediopagos.Channel")
	}

	if filtro.CargarComision {
		resp.Preload("Movimientocomisions.Cuentacomisions")
	}

	if filtro.CargarComision {
		resp.Preload("Movimientoimpuestos.Impuesto")
	}

	if filtro.CargarTransferencias {
		resp.Preload("Movimientotransferencia")
	}

	if filtro.CargarRetenciones {
		resp.Preload("Movimientoretencions.Retencion.Condicion.Gravamen")
	}

	if filtro.AcumularPorPagoIntentos {
		if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
			if len(filtro.PagoIntentosIds) > 0 {
				resp.Select("id, cuentas_id, sum(monto) as monto, pagointentos_id, created_at").Group("pagointentos_id").Having("monto > 0").Having("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin).Having("pagointentos_id IN ?", filtro.PagoIntentosIds).Order("movimientos.created_at")
			} else {
				logs.Info(filtro.FechaInicio)
				logs.Info(filtro.FechaFin)
				resp.Select("id, cuentas_id, sum(monto) as monto, pagointentos_id, created_at").Group("pagointentos_id").Having("monto > 0").Having("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin).Order("movimientos.created_at")
			}
		} else {
			if len(filtro.PagoIntentosIds) > 0 {
				resp.Select("id, cuentas_id, sum(monto) as monto, pagointentos_id, created_at").Group("pagointentos_id").Having("monto > 0").Having("pagointentos_id IN ?", filtro.PagoIntentosIds).Order("movimientos.created_at")
			} else {
				resp.Select("id, cuentas_id, sum(monto) as monto, pagointentos_id, created_at").Group("pagointentos_id").Having("monto > 0").Order("movimientos.created_at")
			}
		}
	} else {
		if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
			resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
		}
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
	// devolver los movimientos en orden descendente por fecha de creacion
	resp.Order("created_at desc").Find(&movimiento)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(ERROR_MOVIMIENTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimientos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}

	return
}

func (r *repository) GetMovimientosTransferencias(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if request.PagoIntento > 0 {
		resp.Where("pagointentos_id = ?", request.PagoIntento)
	}

	if len(request.PagoIntentos) > 0 {
		resp.Where("pagointentos_id IN ?", request.PagoIntentos)
		resp.Preload("Pagointentos")
	}

	if len(request.TipoMovimiento) > 0 {
		resp.Where("tipo = ?", request.TipoMovimiento)
	}

	resp.Where("reversion = ?", request.CargarReversion)

	if request.CargarComisionImpuesto {
		resp.Preload("Movimientocomisions")
		resp.Preload("Movimientoimpuestos")
	}
	resp.Find(&movimientos)

	return
}
func (r *repository) GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestadoexterno{})

	if len(filtro.Vendor) > 0 {

		resp.Where("vendor", filtro.Vendor)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}
	if filtro.CargarEstadosInt {
		resp.Preload("PagoEstados")
	}

	resp.Find(&estados)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(ERROR_PAGO_ESTADO_EXTERNO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosEstadosExternos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}

	return

}

func (r *repository) BajaMovimiento(ctx context.Context, movimientos []*entities.Movimiento, motivoBaja string) error {

	resp := r.SQLClient.WithContext(ctx).Model(&movimientos).Omit(clause.Associations).UpdateColumns(map[string]interface{}{"updated_at": time.Now(), "deleted_at": time.Now(), "motivo_baja": motivoBaja})

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro := fmt.Errorf(ERROR_BAJAR_MOVIMIENTOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "BajaMovimiento",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
		return erro
	}

	err := r.auditarAdministracion(resp.Statement.Context, movimientos)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateMovimientosTransferencia(ctx context.Context, movimiento []*entities.Movimiento) (erro error) {

	res := r.SQLClient.WithContext(ctx).Create(&movimiento)

	if res.Error != nil {

		erro = fmt.Errorf(ERROR_CREAR_MOVIMIENTOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       res.Error.Error(),
			Funcionalidad: "CreateMovimientosTransferencia",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), res.Error)
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(res.Statement.Context, movimiento)
	if erro != nil {
		return erro
	}

	return nil
}

func (r *repository) CreateMovimientosCierreLote(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.

	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// 1 - creo los movimientos

		if len(mcl.ListaMovimientos) > 0 {

			res := tx.WithContext(ctx).Create(&mcl.ListaMovimientos)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_MOVIMIENTO_CREAR)
			}
			r.auditarCierreLote(res.Statement.Context, res.RowsAffected)

		}

		// 3 - Creo los pagos estados logs
		if len(mcl.ListaPagosEstadoLogs) > 0 {
			if err := tx.WithContext(ctx).Omit(clause.Associations).Create(&mcl.ListaPagosEstadoLogs).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_CREAR_ESTADO_LOGS)
			}
		}

		// 4 - Modifico los pagos intentos para indicar que se recibio el pago
		if len(mcl.ListaMovimientos) > 0 && len(mcl.ListaPagoIntentos) > 0 {
			res := tx.WithContext(ctx).Model(&mcl.ListaPagoIntentos).Omit(clause.Associations).Updates(entities.Pagointento{AvailableAt: mcl.ListaPagoIntentos[0].AvailableAt, RevertedAt: mcl.ListaPagoIntentos[0].RevertedAt})
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_MOVIMIENTO_PAGOS_INTENTOS)
			}
			r.auditarCierreLote(res.Statement.Context, res.RowsAffected)
		}

		// 5 - Le doy de baja en los cierres de lote
		if len(mcl.ListaCLApiLink) > 0 {
			if err := tx.Delete(&mcl.ListaCLApiLink).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_MOVIMIENTO_CIERRE)
			}
		}

		if len(mcl.ListaReversiones) > 0 {
			if len(mcl.ListaCLPrisma) > 0 {
				err := tx.Model(&mcl.ListaCLPrisma).Omit(clause.Associations).Unscoped().Updates(entities.Prismacierrelote{Estadomovimiento: true})
				if err.Error != nil {
					logs.Info(err.Error)
					return errors.New(ERROR_MOVIMIENTO_CIERRE)
				}
			}
		}

		if len(mcl.ListaReversiones) == 0 {
			if len(mcl.ListaCLPrisma) > 0 {
				if err := tx.Delete(&mcl.ListaCLPrisma).Error; err != nil {
					logs.Info(err.Error())
					return errors.New(ERROR_MOVIMIENTO_CIERRE)
				}
			}
		}

		if len(mcl.ListaCLRapipago) > 0 {
			if err := tx.Delete(&mcl.ListaCLRapipago).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_MOVIMIENTO_CIERRE)
			}
		}

		// FIXME verificar con las modificaciones
		// if len(mcl.ListaCLRapipagoHeaders) > 0 {
		// 	for _, valueCLR := range mcl.ListaCLRapipagoHeaders {
		// 		resp := tx.Model(entities.Rapipagocierrelote{}).Where("id = ?", valueCLR.ID).UpdateColumns(map[string]interface{}{"pago_match": valueCLR.PagoMatch})
		// 		if resp.Error != nil {
		// 			logs.Info(resp.Error)
		// 			erro = errors.New(ERROR_MOVIMIENTO_CIERRE)
		// 			return erro
		// 		}
		// 	}
		// }
		if len(mcl.ListaCLRapipagoHeaders) > 0 {
			if err := tx.Delete(&mcl.ListaCLRapipagoHeaders).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_MOVIMIENTO_CIERRE)
			}
		}

		// 6 - Modifico el estado de los pagos de acuerdo con el estado
		if len(mcl.ListaPagos) > 0 {
			filtroPagoEstado := filtros.PagoEstadoFiltro{}
			pagosEstados, err := r.GetPagosEstados(filtroPagoEstado)
			if err != nil {
				return err
			}
			for _, pe := range pagosEstados {
				var listaPagosPorEstado []entities.Pago
				for _, p := range mcl.ListaPagos {
					if p.PagoestadosID == int64(pe.ID) {
						listaPagosPorEstado = append(listaPagosPorEstado, p)
					}
				}
				if len(listaPagosPorEstado) > 0 {
					res := tx.WithContext(ctx).Model(&listaPagosPorEstado).Omit(clause.Associations).Update("pagoestados_id", &listaPagosPorEstado[0].PagoestadosID)
					if res.Error != nil {
						logs.Info(err.Error())
						return errors.New(ERROR_MOVIMIENTO_PAGOS)
					}
					r.auditarCierreLote(res.Statement.Context, res.RowsAffected)
				}
				listaPagosPorEstado = nil
			}

		}

		if len(mcl.ListaReversiones) > 0 {
			if err := tx.WithContext(ctx).Omit(clause.Associations).Create(&mcl.ListaReversiones).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_REVERSIONES_CREATE)
			}
		}

		return nil
	})
}

func (r *repository) CreateCLApilinkPagosRepository(ctx context.Context, pg administraciondtos.RegistroClPagosApilink) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.

	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// 1 - creo los registros en apilinkcierrelote
		if len(pg.ListaCLApiLink) > 0 {
			res := tx.WithContext(ctx).Create(&pg.ListaCLApiLink)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_CREAR_CLAPILINK)
			}
		}

		// actualizo el estado de los pagos debin
		if len(pg.ListaPagos) > 0 {
			for _, valuePago := range pg.ListaPagos {
				resp := tx.Model(entities.Pago{}).Where("id = ?", valuePago.Id).UpdateColumns(map[string]interface{}{"pagoestados_id": valuePago.Estadopago})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New(ACTUALIZAR_ESTADOS_CL_APILINK)
					return erro
				}
			}
		}

		return nil
	})
}

func (r *repository) ActualizarPagosClRapipagoRepository(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.

	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// actualizar estados del clrapipago
		if len(pagosclrapiapgo.ListaCLRapipagoHeaders) > 0 {
			for _, valueCLR := range pagosclrapiapgo.ListaCLRapipagoHeaders {
				resp := tx.Model(entities.Rapipagocierrelote{}).Where("id = ?", valueCLR).UpdateColumns(map[string]interface{}{"pago_actualizado": 1})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New(ACTUALIZAR_ESTADOS_CL_RAPIPAGO)
					return erro
				}
			}
		}

		// actualizar estados del pago
		if len(pagosclrapiapgo.ListaPagos) > 0 {
			for _, valuePago := range pagosclrapiapgo.ListaPagos {
				resp := tx.Model(entities.Pago{}).Where("id = ?", valuePago).UpdateColumns(map[string]interface{}{"pagoestados_id": pagosclrapiapgo.EstadoAprobado})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New(ACTUALIZAR_ESTADOS_CL_RAPIPAGO)
					return erro
				}
			}
		}

		return nil
	})
}

func (r *repository) auditarCierreLote(ctx context.Context, resultado interface{}) error {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	audit.Operacion = strings.ToLower(audit.Query[:6])

	audit.Origen = "pasarela.CierreLote"

	res, _ := json.Marshal(resultado)
	audit.Resultado = string(res)

	err := r.auditoriaService.Create(&audit)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

func (r *repository) UpdateTransferencias(listas bancodtos.ResponseConciliacion) error {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		for _, TC := range listas.TransferenciasConciliadas {
			res := r.SQLClient.Table("transferencias").Where("id IN ?", TC.ListaIdsTransferenciasConciliadas).Updates(map[string]interface{}{"match": TC.Match, "banco_external_id": TC.BancoExternalId})
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_UPDATE_TRANSFERENCIAS)
			}
		}

		return nil
	})
}

func (r *repository) GetConsultarDebines(request linkdebin.RequestDebines) (cierreLotes []*entities.Apilinkcierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Apilinkcierrelote{})
	if len(request.Debines) > 0 {
		resp.Where("debin_id IN ?", request.Debines)
	}

	if request.BancoExternalId {
		resp.Where("banco_external_id != ?", 0)
	} else {
		resp.Where("banco_external_id = ?", 0)
	}

	if request.Pagoinformado {
		resp.Where("pagoinformado = ?", false)
	}

	if request.CargarPagoEstado {
		resp.Preload("Pagoestadoexterno")
	}

	resp.Find(&cierreLotes)

	if resp.Error != nil {
		erro = resp.Error
		return nil, erro
	}

	return cierreLotes, nil
}

// NOTE
// mismo comportamiento que pagos positivos
func (r *repository) GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if filtro.CuentaId > 0 {
		resp.Where("movimientos.cuentas_id = ?", filtro.CuentaId)
	}

	// if filtro.CargarMovimientosDetalles {
	// 	resp.Where("movimientos.tipo = ? AND monto < ?", "C", 0)
	// }

	if filtro.AcumularPorPagoIntentos {
		resp.Select("id, cuentas_id, sum(monto) as monto, pagointentos_id, created_at").Group("pagointentos_id").Having("monto < 0").Order("movimientos.created_at")
	}

	if filtro.CargarPagoIntentos {
		resp.Preload("Pagointentos.Pago.PagosTipo")
	}

	if filtro.CargarPagoEstados {
		resp.Preload("Pagointentos.Pago.PagoEstados")
	}

	if filtro.CargarMedioPago {
		resp.Preload("Pagointentos.Mediopagos")
	}

	if filtro.CargarComision {
		resp.Preload("Movimientocomisions.Cuentacomisions")
	}

	if filtro.CargarComision {
		resp.Preload("Movimientoimpuestos.Impuesto")
	}

	resp.Find(&movimiento)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(ERROR_MOVIMIENTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimientosNegativos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
	}

	return
}

func (r *repository) CreateMovimientosTemporalesCierreLote(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.

	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// 1 - creo los movimientos

		if len(mcl.ListaMovimientos) > 0 {

			res := tx.WithContext(ctx).Create(&mcl.ListaMovimientos)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_MOVIMIENTO_CREAR)
			}
			r.auditarCierreLote(res.Statement.Context, res.RowsAffected)

		}

		// FIXME verificar con las modificaciones
		if len(mcl.ListaPagosCalculado) > 0 {
			for _, valueCLR := range mcl.ListaPagosCalculado {
				resp := tx.Model(entities.Pagointento{}).Where("pagos_id = ?", valueCLR).UpdateColumns(map[string]interface{}{"calculado": 1})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New(ERROR_MOVIMIENTO_CIERRE)
					return erro
				}
			}
		}

		return nil
	})
}

/* update pagosm notificados*/
func (r *repository) UpdateCierreloteApilink(request linkdebin.RequestListaUpdateDebines) (erro error) {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if len(request.DebinId) > 0 {
			res := r.SQLClient.Table("apilinkcierrelotes").Where("id IN ?", request.DebinId).Updates(map[string]interface{}{"pagoinformado": 1})
			if res.Error != nil {
				logs.Info(res.Error)
				erro := fmt.Errorf("no se puedo actualizar los pagos apilink cierrelote")
				return erro
			}
		}

		if len(request.Debines) > 0 {
			for _, deb := range request.Debines {
				res := r.SQLClient.Table("apilinkcierrelotes").Where("id = ?", deb.ID).Updates(map[string]interface{}{"banco_external_id": deb.BancoExternalId, "match": deb.Match, "fechaacreditacion": deb.Fechaacreditacion})
				if res.Error != nil {
					logs.Info(res.Error)
					erro := fmt.Errorf("no se puedo actualizar los pagos apilink cierrelote")
					return erro
				}
			}
		}

		// 5 - Le doy de baja en los cierres de lote

		if len(request.DebinesNoAcreditados) > 0 {
			if err := tx.Delete(&request.DebinesNoAcreditados).Error; err != nil {
				logs.Info(err.Error())
				return errors.New(ERROR_CIERRELOTE_APILINK)
			}
		}

		return nil
	})

}

// func (r *repository) GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error) {

// 	resp := r.SQLClient.Model(entities.Movimiento{})

// 	if filtro.CuentaId > 0 {
// 		resp.Where("movimientos.cuentas_id = ? AND movimientos.monto < ? ", filtro.CuentaId, 0)
// 	}

// 	if filtro.CargarMovimientosNegativos {
// 		resp.Where("movimientos.tipo = ? ", "C")
// 	}

// 	if filtro.CargarPagoIntentos {
// 		resp.Preload("Pagointentos.Pago.PagosTipo")
// 	}

// 	if filtro.CargarPagoEstados {
// 		resp.Preload("Pagointentos.Pago.PagoEstados")
// 	}

// 	if filtro.CargarMedioPago {
// 		resp.Preload("Pagointentos.Mediopagos")
// 	}

// 	if filtro.CargarComision {
// 		resp.Preload("Movimientocomisions.Cuentacomisions")
// 	}

// 	if filtro.CargarComision {
// 		resp.Preload("Movimientoimpuestos.Impuesto")
// 	}

// 	resp.Find(&movimiento)

// 	if resp.Error != nil {

// 		logs.Error(resp.Error)

// 		erro = fmt.Errorf(ERROR_MOVIMIENTO)

// 		log := entities.Log{
// 			Tipo:          entities.Error,
// 			Mensaje:       resp.Error.Error(),
// 		}
// 			Funcionalidad: "GetMovimientosNegativos",

// 		err := r.utilService.CreateLogService(log)

// 		if err != nil {
// 			logs.Error(err)
// 		}
// 	}

// 	return
// }
