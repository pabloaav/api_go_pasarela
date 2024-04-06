package cierrelote

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *repository) SaveCierreLote(detalleLote *entities.Prismacierrelote) (bool, error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		tx.Create(&detalleLote)
		return nil
	})

	return true, nil
}

func (r *repository) SaveCierreLoteBatch(detalleLote []entities.Prismacierrelote) (bool, error) {
	//println(detalleLote)
	tx := r.SQLClient.Begin()
	err := tx.Create(&detalleLote).Omit(clause.Associations, "prismamovimientodetalles_id").Error
	if err != nil {
		tx.Rollback()
		logs.Info(err)
		return false, errors.New("error: al inserte el valor en la base de datos")
	}
	err = tx.Commit().Error
	if err != nil {
		logs.Info(err)
		return false, errors.New("error: al confirmar una transacción")
	}
	return true, nil
}

func (r *repository) GetCierreLoteRepository(filtro filtrocl.FiltroCierreLote) (entityCL []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as cl")
	if filtro.ContraCargo {
		resp.Unscoped()
		if filtro.Reversion {
			resp.Where(" cl.reversion = 1")
		}
		if !filtro.Reversion {
			resp.Where(" cl.reversion = 0")
		}
		if filtro.ContraCargoMx {
			resp.Where(" cl.detallemovimiento_id <> 0")
		}
		if !filtro.ContraCargoMx {
			resp.Where(" cl.detallemovimiento_id = 0")
		}
		if filtro.ContraCargoPx {
			resp.Where(" cl.detallepago_id <> 0")
		}
		if !filtro.ContraCargoPx {
			resp.Where(" cl.detallepago_id = 0")
		}
	}
	if filtro.MovimientosMX {
		resp.Where("cl.prismamovimientodetalles_id is NULL")
	}
	if !filtro.MovimientosMX {
		resp.Where("cl.prismamovimientodetalles_id is NOT NULL ")
	}
	if filtro.PagosPx {
		resp.Where("cl.prismatrdospagos_id is NULL")
	}
	if !filtro.PagosPx {
		resp.Where("cl.prismatrdospagos_id is NOT NULL")
	}
	if filtro.MatchCl {
		resp.Where(" cl.match = 0")
	}
	if filtro.Banco {
		resp.Where(" cl.banco_external_id = 0")
	}
	if filtro.EstadoFechaPago {
		if filtro.FechaPago == "0000-00-00" {
			resp.Where("cl.fecha_pago <> ?", filtro.FechaPago)
		} else {
			resp.Where("cl.fecha_pago = ?", filtro.FechaPago)
		}

	}
	if filtro.PrismaPagoId > 0 {
		resp.Where("cl.prismatrdospagos_id = ?", filtro.PrismaPagoId)
	}
	if filtro.DetallePagoId > 0 {
		resp.Where("cl.detallepago_id = ?", filtro.DetallePagoId)
	}
	if filtro.Compras {
		resp.Where("cl.tipooperacion = 'C'")
	}
	if filtro.Devolucion {
		resp.Where("cl.tipooperacion = 'D'")
	}
	if filtro.Anulacion {
		resp.Where("cl.tipooperacion = 'A'")
	}

	resp.Preload("Channelarancel")
	resp.Find(&entityCL)

	if err := resp.Error; err != nil {
		logs.Info(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_CIERRE_LOTE)
		return
	}
	// if resp.RowsAffected <= 0 {
	// 	erro = errors.New("error: " + ERROR_LISTA_CIERRE_LOTE_VACIA)
	// 	logs.Info(erro.Error())
	// 	return
	// }
	return
}

func (r *repository) GetPagosIntentosByMedioPagoIdRepository(arraysMediosPagoIds []int64) (listaPagoIntentos []entities.Pagointento, erro error) {

	resp := r.SQLClient.Table("pagointentos as pi")
	resp.Joins("JOIN mediopagos as mp ON mp.external_id IN (?) and pi.mediopagos_id = mp.id", arraysMediosPagoIds).
		Joins("JOIN pagos as p ON p.id = pi.pagos_id").
		Joins("JOIN pagotipos as pt ON pt.id = p.pagostipo_id").
		Joins("JOIN channels as cha on cha.id = mp.channels_id").
		Joins("JOIN cuentas as c ON c.id = pt.cuentas_id").
		Joins("join installmentdetails as instal on instal.id = pi.installmentdetails_id").
		Joins("JOIN rubros  as r on c.rubros_id = r.id").
		Joins("JOIN clientes as cli on cli.id = c.clientes_id").
		Joins("JOIN impuestos as impiva on impiva.id = cli.iva_id").
		Joins("JOIN impuestos as impiibb on impiibb.id = cli.iibb_id").
		Preload("Mediopagos").
		Preload("Mediopagos.Channel").
		Preload("Mediopagos.Channel.Channelaranceles").
		Preload("Pago.PagosTipo.Cuenta").
		Preload("Pago.PagosTipo.Cuenta.Rubro").
		Preload("Pago.PagosTipo.Cuenta.Cliente").
		Preload("Pago.PagosTipo.Cuenta.Cliente.Iva").
		Preload("Pago.PagosTipo.Cuenta.Cliente.Iibb").
		Preload("Installmentdetail").
		Where("state_comment = 'approved'").
		Find(&listaPagoIntentos)
	if resp.Error != nil {
		logs.Info(resp.Error)
		return nil, errors.New("error: al consultar los pagos intentos ")
	}
	return listaPagoIntentos, nil
}

func (r *repository) GetCierreLoteGroupByRepository(nroCuota int64) (ClMatch []cierrelotedtos.PrismaClResultGroup, erro error) {

	// resp := r.SQLClient.Table("prismacierrelotes as pcl").
	// 	Select("count(id) as cantidadregistro, pcl.externalmediopago_id, pcl.externallote_id, pcl.nroestablecimiento, pcl.nombrearchivolote, sum(pcl.montofinal) as monto, pcl.fechaoperacion, pcl.fecha_cierre, pcl.match, pcl.nrocuota").
	// 	Group("pcl.nombrearchivolote, pcl.nroestablecimiento, pcl.externallote_id").
	//  Having("pcl.match = 0 ")
	resp := r.SQLClient.Table("prismacierrelotes as pcl").
		Select("count(id) as cantidadregistro, pcl.externalmediopago_id, pcl.externallote_id, pcl.nroestablecimiento, pcl.nombrearchivolote, sum(pcl.montofinal) as monto, pcl.fechaoperacion, pcl.fecha_cierre, pcl.match, pcl.nrocuota").
		Where("pcl.match = 0 ").
		Group("pcl.nombrearchivolote, pcl.nroestablecimiento, pcl.externallote_id")
	if nroCuota == 1 {
		resp.Where("pcl.nrocuota = ?", nroCuota)
	}
	if nroCuota > 1 {
		resp.Where("pcl.nrocuota > 1")
	}
	resp.Find(&ClMatch)
	if resp.Error != nil {
		logs.Error(resp.Error)
		erro = errors.New("error: al consultar los cierres de lotes ")
		return
	}
	if resp.RowsAffected == 0 {
		erro = errors.New("no hay cierres de lotes para el nro de establecimiento")
		logs.Info(erro)
		return
	}
	return
}

func (r *repository) ActualizarCierreLoteMatch(reversion bool, clMatch []entities.Prismacierrelote) (erro error) {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, cl := range clMatch {
			if reversion {
				if err := tx.Model(&entities.Prismacierrelote{}).Where("id = ?", cl.ID).Unscoped().Updates(entities.Prismacierrelote{ExtbancoreversionId: cl.ExtbancoreversionId, Conciliado: cl.Conciliado, Descripcionbanco: cl.Descripcionbanco}).Error; err != nil {
					logs.Info(err)
					return errors.New("error: al actualizar el cierre de lote")
				}
			}
			if !reversion {
				if err := tx.Model(&entities.Prismacierrelote{}).Where("id = ?", cl.ID).Updates(entities.Prismacierrelote{Match: cl.Match, BancoExternalId: cl.BancoExternalId, Diferenciaimporte: cl.Diferenciaimporte, Descripcion: cl.Descripcion}).Error; err != nil {
					logs.Info(err)
					return errors.New("error: al actualizar el cierre de lote")
				}
			}
		}
		return nil
	})
}

/*
Omit("pagoestadoexternos_id", "channelaranceles_id", "impuestos_id", "tiporegistro", "pagos_uuid", "externalmediopago", "nrotarjeta", "tipooperacion", "fechaoperacion", "monto", "montofinal", "codigoautorizacion", "nroticket", "site_id", "externallote_id", "nrocuota", "fecha_cierre", "nroestablecimiento", "externalcliente_id", "nombrearchivolote")
*/

func (r *repository) GetCierreLoteByGroup(cl cierrelotedtos.PrismaClResultGroup) (clPrisma []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as cl").
		Where("cl.nroestablecimiento = ? and cl.externallote_id = ? and cl.fechaoperacion = ? and cl.fecha_cierre = ? and cl.nombrearchivolote = ? and cl.nrocuota = ?", cl.Nroestablecimiento, cl.ExternalloteId, cl.Fechaoperacion, cl.FechaCierre, cl.Nombrearchivolote, cl.Nrocuota).Find(&clPrisma)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: al consultar los cierres de lotes")
		return
	}
	return
}

func (r *repository) GetCierreLotePrismaByExternalIdAndMacht() (listCierreLote []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as cl").Where("cl.match = 1 and cl.banco_external_id <> 0").Find(&listCierreLote)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: al consultar los cierres de lotes")
	}
	return
}

func (r *repository) SaveTransactionPagoPx(pagoPx []entities.Prismapxcuatroregistro) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pagoPx).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al guardar los registros 02 y 04 del archivo pago px")
			return erro
		}
		return nil
	})
	return
}
func (r *repository) SaveTransactionMovimientoMx(movimientosMx []entities.Prismamxtotalesmovimiento) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&movimientosMx).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al guardar los registros totales y detalle del archivo movimiento mx")
			return erro
		}
		return nil
	})
	return
}

func (r *repository) SaveTransactionPagoRP(pagoRP []entities.Rapipagocierrelote) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pagoRP).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al guardar los registros rapipago cierre lote")
			return erro
		}
		return nil
	})
	return
	// tx := r.SQLClient.Begin()
	// err := tx.Create(&pagoRP).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	logs.Info(err)
	// 	return errors.New("error: al inserte el valor en la base de datos ")
	// }
	// err = tx.Commit().Error
	// if err != nil {
	// 	logs.Info(err)
	// 	return errors.New("error: al confirmar una transacción ")
	// }
	// return nil
}

func (r *repository) GetMovimientosMxRepository() (movimientosMx []entities.Prismamxtotalesmovimiento, erro error) {
	resp := r.SQLClient.Model(&entities.Prismamxtotalesmovimiento{}).Preload("MovimientosDetalle").Find(&movimientosMx)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_MOVIMIENTOS_MX)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_MOVIMIENTOS_MX)
		return
	}
	return
}

func (r *repository) GetCodigosRechazoRepository() (codigosRechazo []entities.Prismacodigorechazo, erro error) {
	resp := r.SQLClient.Model(&entities.Prismacodigorechazo{}).Find(&codigosRechazo)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_TABLA_CODIGOS_RECHAZOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_TABLA_CODIGOS_RECHAZOS)
		return
	}
	return
}
func (r *repository) GetVisaContracargoRepository() (visaContracargo []entities.Prismavisacontracargo, erro error) {
	resp := r.SQLClient.Model(&entities.Prismavisacontracargo{}).Find(&visaContracargo)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_TABLA_VISA_CONTRACARGO)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_TABLA_VISA_CONTRACARGO)
		return
	}
	return
}
func (r *repository) GetMotivosAjustesRepository() (motivosAjustes []entities.Prismamotivosajuste, erro error) {
	resp := r.SQLClient.Model(&entities.Prismamotivosajuste{}).Find(&motivosAjustes)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_TABLA_MOTIVOS_AJUSTE)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_TABLA_MOTIVOS_AJUSTE)
		return
	}
	return
}
func (r *repository) GetOperacionesRepository() (operaciones []entities.Prismaoperacion, erro error) {
	resp := r.SQLClient.Model(&entities.Prismaoperacion{}).Find(&operaciones)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_TABLA_OPERACIONES)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_TABLA_OPERACIONES)
		return
	}
	return
}
func (r *repository) GetMasterContracargoRepository() (masterContracargo []entities.Prismamastercontracargo, erro error) {
	resp := r.SQLClient.Model(&entities.Prismamastercontracargo{}).Find(&masterContracargo)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_TABLA_MASTER_CONTRACARGO)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_TABLA_MASTER_CONTRACARGO)
		return
	}
	return
}
func (r *repository) SaveMovimientoMxRepository(movimientosMx []entities.Prismamovimientototale, movimientosMxEntity []entities.Prismamxtotalesmovimiento) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&movimientosMx).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al guardar los registros de las tablas movimineto totales y detalle")
			return erro
		}
		for _, valueDetalle := range movimientosMxEntity {
			if err := tx.Delete(&valueDetalle.MovimientosDetalle).Error; err != nil {
				logs.Info(err)
				erro = errors.New("error: al borrar los registros procesados de la tabla prismamxmoviminetodetalle")
				return erro
			}
		}
		if err := tx.Delete(&movimientosMxEntity).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al borrar los registros procesados de la tabla prismamxmoviminetototales")
			return erro
		}
		return nil
	})
	return
}

func (r *repository) GetPagosPxRepository() (pagosPx []entities.Prismapxcuatroregistro, erro error) {

	resp := r.SQLClient.Model(&entities.Prismapxcuatroregistro{}).Preload("PxDosRegistros").Find(&pagosPx)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: " + ERROR_CONSULTAR_PAGOS_PX)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New(RESULT_CONSULTAR_PAGOS_PX)
		return
	}
	return
}

func (r *repository) SavePagosPxRepository(pagosPx []entities.Prismatrcuatropago, entityPagoPxStr []entities.Prismapxcuatroregistro) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pagosPx).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al guardar los registros de las tablas prisma pagos ")
			return erro
		}
		for _, valueDetalle := range entityPagoPxStr {
			if err := tx.Delete(&valueDetalle.PxDosRegistros).Error; err != nil {
				logs.Info(err)
				erro = errors.New("error: al borrar los registros procesados de la tabla prismapxdospago")
				return erro
			}
		}
		if err := tx.Delete(&entityPagoPxStr).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al borrar los registros procesados de la tabla prismapxcuatropago")
			return erro
		}

		return nil
	})
	return
}

func (r *repository) GetArancelByRubroIdChannelIdRepository(rubroId, channelId uint) (entityArancel entities.Channelarancele, erro error) {
	resp := r.SQLClient.Table("channelaranceles as cha").Where("cha.channels_id = ? and rubros_id = ?", channelId, rubroId).Find(&entityArancel)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: al consultar el arancel relacionado con el pago")
		return
	}
	return

}

func (r *repository) GetPrismaMovimientosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityPrismaMovimientos []entities.Prismamovimientototale, erro error) {
	resp := r.SQLClient.Table("prismamovimientototales as pmt")
	if !filtro.Match {
		resp.Where("pmt.match = 0 ")
	}
	if filtro.Match {
		resp.Where("pmt.match = 1 ")
	}
	if len(filtro.CodigosOperacion) > 0 {
		resp.Where("pmt.codop in ?", filtro.CodigosOperacion)
	}
	if filtro.TipoAplicacion != "" {
		resp.Where("pmt.tipo_aplicacion = ?", filtro.TipoAplicacion)
	}
	if !filtro.FechaPresentacion.IsZero() {
		resp.Where("pmt.fecha_presentacion = ?", filtro.FechaPresentacion)
	}
	if len(filtro.EstablecimientoNro) != 0 {
		resp.Where("pmt.establecimiento_nro like ?", "%"+filtro.EstablecimientoNro)
	}
	if filtro.CargarDetalle {
		if filtro.ContraCargo {
			resp.Preload("DetalleMovimientos", "prismamovimientodetalles.match = ?", 0)
		}
		if !filtro.ContraCargo {
			resp.Preload("DetalleMovimientos")
		}
		if filtro.Contracargovisa {
			resp.Preload("DetalleMovimientos.Contracargovisa")
		}
		if filtro.Contracargomaster {

			resp.Preload("DetalleMovimientos.Contracargomaster")
		}
		if filtro.Tipooperacion {
			resp.Preload("DetalleMovimientos.Tipooperacion")
		}
		if filtro.Rechazotransaccionprincipal {
			resp.Preload("DetalleMovimientos.Rechazotransaccionprincipal")
		}
		if filtro.Rechazotransaccionsecundario {
			resp.Preload("DetalleMovimientos.Rechazotransaccionsecundario")
		}
		if filtro.Motivoajuste {
			resp.Preload("DetalleMovimientos.Motivoajuste")
		}
	}
	resp.Find(&entityPrismaMovimientos)
	if err := resp.Error; err != nil {
		logs.Info(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_MOVIMIENTOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_MOVIMIENTOS_VACIA)
		logs.Info(erro.Error())
		return
	}
	return
}

func (r *repository) GetContraCargoPrismaMovimientosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityPrismaMovimientos []entities.Prismamovimientototale, erro error) {
	resp := r.SQLClient.Table("prismamovimientototales as pmt")
	resp.Where("pmt.match = 0")
	if len(filtro.CodigosOperacion) > 0 {
		resp.Where("pmt.codop in ?", filtro.CodigosOperacion)
	}
	if filtro.TipoAplicacion != "" {
		resp.Where("pmt.tipo_aplicacion = ?", filtro.TipoAplicacion)
	}

	if filtro.CargarDetalle {
		resp.Preload("DetalleMovimientos", "prismamovimientodetalles.match = ?", 0) //.Joins("inner join prismamovimientodetalles as pmd on pmt.id = pmd.prismamovimientototales_id and pmd.match = ?", filtro.Match)
		if filtro.Contracargovisa {
			resp.Preload("DetalleMovimientos.Contracargovisa")
		}
		if filtro.Contracargomaster {

			resp.Preload("DetalleMovimientos.Contracargomaster")
		}
		if filtro.Tipooperacion {
			resp.Preload("DetalleMovimientos.Tipooperacion")
		}
		if filtro.Rechazotransaccionprincipal {
			resp.Preload("DetalleMovimientos.Rechazotransaccionprincipal")
		}
		if filtro.Rechazotransaccionsecundario {
			resp.Preload("DetalleMovimientos.Rechazotransaccionsecundario")
		}
		if filtro.Motivoajuste {
			resp.Preload("DetalleMovimientos.Motivoajuste")
		}
	}

	resp.Find(&entityPrismaMovimientos)
	if err := resp.Error; err != nil {
		logs.Info(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_MOVIMIENTOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_MOVIMIENTOS_VACIA)
		logs.Info(erro.Error())
		return
	}
	return

}

func (r *repository) UpdateCierreloteAndMoviminetosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, valueCL := range entityCierreLote {
			resp := tx.Model(entities.Prismacierrelote{}).Unscoped().Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{
				"prismamovimientodetalles_id": valueCL.PrismamovimientodetallesId,
				"fecha_pago":                  valueCL.FechaPago,
				"cantdias":                    valueCL.Cantdias,
				"enobservacion":               valueCL.Enobservacion,
				"valorpresentado":             valueCL.Valorpresentado,
				"diferenciaimporte":           valueCL.Diferenciaimporte,
				"coeficientecalculado":        valueCL.Coeficientecalculado,
				"costototalporcentaje":        valueCL.Costototalporcentaje,
				"importeivaarancel":           valueCL.Importeivaarancel,
				"descripcionpresentacion":     valueCL.Descripcionpresentacion,
				"detallemovimiento_id":        valueCL.DetallemovimientoId,
				"detallepago_id":              valueCL.DetallepagoId,
				"reversion":                   valueCL.Reversion,
				"descripcioncontracargo":      valueCL.Descripcioncontracargo,
			})
			if resp.Error != nil {
				logs.Info(resp.Error)
				erro = errors.New("error: al actualizar tabla de cierre de lote")
				return erro
			}
		}

		if err := tx.Model(&entities.Prismamovimientototale{}).Where("id in (?)", listaIdsCabecera).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al actualizar tabla Prisma Movimientos cabecera")
			return erro
		}
		if err := tx.Model(&entities.Prismamovimientodetalle{}).Where("id in (?)", listaIdsDetalle).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
			logs.Info(err)
			erro = errors.New("error: al actualizar tabla Prisma Movimientos detalle")
			return erro
		}
		return nil
	})

	return
}

func (r *repository) GetMovimientosConciliadosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityMovimientosCabeceraConciliados []entities.Prismamovimientototale, erro error) { //listaMovimientoDetalle []entities.Prismamovimientodetalle
	resp := r.SQLClient.Table("prismamovimientototales as pmt")
	var valorMatch int64
	if filtro.Match {
		valorMatch = 1
	}
	if filtro.ContraCargo {
		resp.Where("pmt.codop in (1507, 6000, 1517)")
	}
	resp.Where("pmt.match = ?", valorMatch)
	resp.Find(&entityMovimientosCabeceraConciliados)
	if err := resp.Error; err != nil {
		logs.Info(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_MOVIMIENTOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_MOVIMIENTOS_VACIA)
		logs.Info(erro.Error())
		return
	}
	return
}
func (r *repository) GetMovimientosDetalleConciliadosRepository(filtro filtrocl.FiltroPrismaMovimientoDetalle) (entityMovimientosDetalleConciliados []entities.Prismamovimientodetalle, erro error) {
	resp := r.SQLClient.Table("prismamovimientodetalles as pmd")

	if len(filtro.FechaPago) > 0 {
		var valorMatchCl int64
		if filtro.MatchCl {
			valorMatchCl = 1
		}
		resp.Joins("join prismacierrelotes as cl on cl.prismamovimientodetalles_id = pmd.id and cl.fecha_pago > ? and cl.match = ?", filtro.FechaPago, valorMatchCl)
		resp.Preload("CierreLote")
	}

	if filtro.Contracargovisa {
		resp.Preload("Contracargovisa")
	}
	if filtro.Contracargomaster {
		resp.Preload("Contracargomaster")
	}

	if filtro.Tipooperacion {
		resp.Preload("Tipooperacion")
	}
	if filtro.Rechazotransaccionprincipal {
		resp.Preload("Rechazotransaccionprincipal")
	}
	if filtro.Rechazotransaccionsecundario {
		resp.Preload("Rechazotransaccionsecundario")
	}
	if filtro.Motivoajuste {
		resp.Preload("Motivoajuste")
	}
	if filtro.Match {
		var valorMatch int64
		if filtro.Match {
			valorMatch = 1
		}
		resp.Where("pmd.match = ?", valorMatch)
	}
	if len(filtro.ListIdsCabecera) > 0 {
		resp.Where("pmd.prismamovimientototales_id in ?", filtro.ListIdsCabecera)
	}
	resp.Find(&entityMovimientosDetalleConciliados)
	if err := resp.Error; err != nil {
		logs.Error(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_MOVIMIENTOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_MOVIMIENTOS_VACIA)
		logs.Error(erro.Error())
		return
	}
	return
}

func (r *repository) GetPrismaPagosRepository(filtro filtrocl.FiltroPrismaTrPagos) (entityPrismaPago []entities.Prismatrcuatropago, erro error) {
	resp := r.SQLClient.Table("prismatrcuatropagos as ppc")
	// if !filtro.Devolucion {
	// 	if filtro.Match {
	// 		resp.Where("ppc.match = 1")
	// 	} else {
	// 		resp.Where("ppc.match = 0")
	// 	}
	// }
	if len(filtro.FechaPagos) > 0 {
		resp.Where("ppc.fecha_pago in (?)", filtro.FechaPagos)
	}
	if filtro.CargarDetalle {
		resp.Preload("Pagostrdos")
	}
	//resp.Where("ppc.fecha_presentacion = '2022-06-30' and ppc.fecha_pago = '2022-07-01' and ppc.establecimiento_nro = '0092163468'")
	resp.Find(&entityPrismaPago)
	if err := resp.Error; err != nil {
		logs.Error(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_PAGOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_PAGOS_VACIA)
		logs.Error(erro.Error())
		return
	}

	return
}

func (r *repository) UpdateCierreloteAndPagosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, valueCL := range entityCierreLote {
			resp := tx.Model(entities.Prismacierrelote{})
			if valueCL.Reversion {
				resp.Unscoped().Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"detallepago_id": valueCL.DetallepagoId})
			}
			if !valueCL.Reversion {
				resp.Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"prismatrdospagos_id": valueCL.PrismatrdospagosId})
			}
			if resp.Error != nil {
				logs.Info(resp.Error)
				erro = errors.New("error: al actualizar tabla de cierre de lote")
				return erro
			}
		}
		if err := tx.Model(&entities.Prismatrcuatropago{}).Where("id in (?)", listaIdsCabecera).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
			logs.Error(err)
			erro = errors.New("error: al actualizar tabla Prisma Pagos cabecera")
			return erro
		}
		if err := tx.Model(&entities.Prismatrdospago{}).Where("id in (?)", listaIdsDetalle).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
			logs.Error(err)
			erro = errors.New("error: al actualizar tabla Prisma Pagos detalle")
			return erro
		}
		return nil
	})
	return
}

func (r *repository) GetCierreLoteMatch(filtro filtrocl.FiltroTablasConciliadas) (entityPrismaTr4Pagos []entities.Prismatrcuatropago, erro error) {

	resp := r.SQLClient.Table("prismatrcuatropagos as ptr4pago")
	resp.Joins("join prismatrdospagos as ptr2pago on ptr4pago.id = ptr2pago.prismatrcuatropagos_id")
	resp.Preload("Pagostrdos")
	if !filtro.Reversion {
		resp.Joins("join prismacierrelotes as cl on cl.prismatrdospagos_id = ptr2pago.id  ")
	}
	if filtro.Reversion {
		resp.Joins("join prismacierrelotes as cl on cl.reversion = 1 and  cl.detallepago_id = ptr2pago.id  ")
	}

	if len(filtro.FechaPresentacion) > 0 {
		resp.Where("ptr4pago.fecha_presentacion = ?", filtro.FechaPresentacion)
	}
	if len(filtro.FechaPago) > 0 {
		resp.Where("ptr4pago.fecha_pago = ?", filtro.FechaPago)
	}
	if len(filtro.NroEstablecimiento) > 0 {
		resp.Where("ptr4pago.establecimiento_nro = ?", filtro.NroEstablecimiento)
	}
	if filtro.Match {
		resp.Where("ptr4pago.match = 1")
	}
	if !filtro.Match {
		resp.Where("ptr4pago.match = 0")
	}

	resp.Group("ptr4pago.id")
	resp.Find(&entityPrismaTr4Pagos)
	if err := resp.Error; err != nil {
		logs.Error(err)
		erro = errors.New("error: " + ERROR_CONSULTAR_PRISMA_PAGOS)
		return
	}
	if resp.RowsAffected <= 0 {
		erro = errors.New("error: " + ERROR_LISTA_PRISMA_PAGOS_VACIA)
		logs.Error(erro.Error())
		return
	}
	return
}
