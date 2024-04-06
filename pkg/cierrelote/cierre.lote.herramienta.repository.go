package cierrelote

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
	"gorm.io/gorm"
)

func (r *repository) GetAllCierreLotePrismaRepository(filtro filtrocl.CierreLoteFiltro) (listaCierreLote []entities.Prismacierrelote, totalFilas int64, erro error) {
	queryGorm := r.SQLClient.Model(entities.Prismacierrelote{}).Unscoped()
	var tipo_fecha = ""
	// EstadoexternoId
	if filtro.EstadoexternoId > 0 {
		queryGorm.Where("pagoestadoexternos_id = ?", filtro.EstadoexternoId)
	}
	// TipoOperacion
	if len(filtro.TipoOperacion) > 0 {
		queryGorm.Where("tipooperacion = ?", strings.ToUpper(filtro.TipoOperacion))
	}
	// Disputa
	if filtro.Disputa {
		queryGorm.Where("disputa = ?", filtro.Disputa)
	}
	// Observacion
	if filtro.Observacion {
		queryGorm.Where("enobservacion = ?", filtro.Observacion)
	}

	// codigo_autorizacion
	if len(filtro.CodigosAutorizacion) > 0 {
		queryGorm.Where("codigoautorizacion IN (?)", filtro.CodigosAutorizacion)
	}
	// tipo de fecha y rango
	if !filtro.FechaInicio.IsZero() {
		fecha_inicio := commons.GetDateFirstMoment(filtro.FechaInicio)
		// fecha de fin por default en now()
		fecha_fin := commons.GetDateLastMoment(time.Now())

		if !filtro.FechaFin.IsZero() {
			fecha_fin = commons.GetDateLastMoment(filtro.FechaFin)
		}

		if filtro.FechaOperacion {
			//fechaoperacion
			tipo_fecha = "fechaoperacion"
		}
		if filtro.FechaPago {
			//fecha_pago
			tipo_fecha = "fecha_pago"
		}
		if filtro.FechaCreacion {
			//created_at
			tipo_fecha = "created_at"
		}

		queryGorm.Where(tipo_fecha+" BETWEEN ? AND ?", fecha_inicio[0:10], fecha_fin[0:10])

	}
	// Paginacion: cuenta
	if filtro.Number > 0 && filtro.Size > 0 {

		// Ejecutar y contar las filas devueltas
		queryGorm.Count(&totalFilas)

		if queryGorm.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
			return
		}

		offset := (filtro.Number - 1) * filtro.Size
		queryGorm.Limit(int(filtro.Size))
		queryGorm.Offset(int(offset))
	}

	// se carga la informacion de la relacion con channelarancel referida en el DTO ResponsePrismaCL
	queryGorm.Preload("Channelarancel")

	// se carga informacion del movimiento relacionado con cierre lote
	if filtro.CargarMovimientoPrisma {
		queryGorm.Where("prismamovimientodetalles_id is not null")
		queryGorm.Preload("Prismamovimientodetalle")
		queryGorm.Preload("Prismamovimientodetalle.MovimientoCabecera")
		queryGorm.Preload("Prismamovimientodetalle.Contracargovisa")
		queryGorm.Preload("Prismamovimientodetalle.Contracargomaster")
		queryGorm.Preload("Prismamovimientodetalle.Tipooperacion")
		queryGorm.Preload("Prismamovimientodetalle.Rechazotransaccionprincipal")
		queryGorm.Preload("Prismamovimientodetalle.Rechazotransaccionsecundario")
		queryGorm.Preload("Prismamovimientodetalle.Motivoajuste")

	}
	// se carga informacion pago relacionado con cierre lote
	if filtro.CargarPagoPrisma {
		queryGorm.Where("prismatrdospagos_id is not null")
		queryGorm.Preload("Prismatrdospagos")
	}

	// Consulta a DB
	queryGorm.Find(&listaCierreLote)

	// capturar error query DB
	if queryGorm.Error != nil {

		erro = fmt.Errorf(ERROR_GET_CIERRE_LOTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetAllCierreLotePrismaRepository",
		}

		err := r.UtilRepository.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}
	}

	// Retorna []entities.Prismacierrelote
	return
}

func (r *repository) GetOneCierreLotePrismaRepository(filtro filtrocl.OneCierreLoteFiltro) (oneCierreLote entities.Prismacierrelote, erro error) {
	queryGorm := r.SQLClient.Model(entities.Prismacierrelote{}).Unscoped()

	if filtro.Id > 0 {
		queryGorm.Where("id = ?", filtro.Id)
	}

	// se carga la informacion de la relacion con channelarancel referida en el DTO ResponsePrismaCL
	queryGorm.Preload("Channelarancel")

	queryGorm.First(&oneCierreLote)

	// el error en la consulta de gorm viene en el objeto queryGorm de tipo *gorm.DB
	if queryGorm.Error != nil {

		if errors.Is(queryGorm.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CIERRE_LOTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetOneCierreLotePrismaRepository",
		}

		erro := r.UtilRepository.CreateLog(log)

		if erro != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetOneCierreLotePrismaRepository: %s", erro.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetAllMovimientosPrismaRepository(filtro filtrocl.FiltroMovimientosPrisma) (entityMovimientoPrisma []entities.Prismamovimientototale, totalFilas int64, erro error) {
	resp := r.SQLClient.Table("prismamovimientototales as pmt")
	if filtro.BuscarPorFechaPresentacion {
		resp.Where("pmt.fecha_presentacion BETWEEN ? AND ? ", filtro.FechaInicio, filtro.FechaFin)
	}

	if filtro.BuscarPorFechaPago {
		fechaInicio := filtro.FechaInicio.Format("2006-01-02") // fmt.Sprintf("%v-%v-%v",filtro.FechaInicio.Year(), filtro.FechaInicio.Month(), filtro.FechaInicio.Day())
		fechaFin := filtro.FechaFin.Format("2006-01-02")       // fmt.Sprintf("%v-%v-%v",filtro.FechaFin.Year(), filtro.FechaFin.Month(), filtro.FechaFin.Day())
		resp.Where("pmt.fecha_pago BETWEEN ? AND ? ", fechaInicio, fechaFin)
		// resp.Preload("DetalleMovimientos").Joins("inner join prismamovimientodetalles as pmd on pmd.fecha_origen_compra BETWEEN ? AND ? ", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.BuscarPorFechaCreacion {
		resp.Where("pmt.created_at BETWEEN ? AND ? ", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.Number > 0 && filtro.Size > 0 {
		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}
	}
	if !filtro.BuscarPorFechaPago {
		resp.Preload("DetalleMovimientos")
	}
	resp.Preload("DetalleMovimientos.Contracargovisa")
	resp.Preload("DetalleMovimientos.Contracargomaster")
	resp.Preload("DetalleMovimientos.Tipooperacion")
	resp.Preload("DetalleMovimientos.Rechazotransaccionprincipal")
	resp.Preload("DetalleMovimientos.Rechazotransaccionsecundario")
	resp.Preload("DetalleMovimientos.Motivoajuste")
	resp.Preload("DetalleMovimientos.MovimientoCabecera")

	if filtro.Number > 0 && filtro.Size > 0 {

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&entityMovimientoPrisma)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: al consultar prsentaciones de prisma")
		return
	}
	return
}

func (r *repository) EditOneCierreLotePrismaRepository(entityPrismaCL entities.Prismacierrelote, filtroEdit cierrelotedtos.FiltroEditarCLHerramienta) (erro error) {
	const YYYYMMDD = "2006-01-02"
	mapStrIntfColumnsUpdate := make(map[string]interface{})
	queryGorm := r.SQLClient.Model(entities.Prismacierrelote{}).Unscoped()
	if filtroEdit.Monto {
		mapStrIntfColumnsUpdate["monto"] = entityPrismaCL.Monto
	}

	if filtroEdit.Disputa {
		mapStrIntfColumnsUpdate["disputa"] = entityPrismaCL.Disputa
	}
	if filtroEdit.Enobservacion {
		mapStrIntfColumnsUpdate["enobservacion"] = entityPrismaCL.Enobservacion
	}
	if filtroEdit.Fecha {
		mapStrIntfColumnsUpdate["fecha_cierre"] = entityPrismaCL.FechaCierre.Format(YYYYMMDD)
	}

	// mapStrIntfColumnsUpdate := map[string]interface{}{
	// 	"fecha_cierre":  entityPrismaCL.FechaCierre.Format(YYYYMMDD),
	// 	"disputa":       entityPrismaCL.Disputa,
	// 	"enobservacion": entityPrismaCL.Enobservacion,
	// }

	//result := r.SQLClient.Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)
	resp := queryGorm.Where("id = ?", entityPrismaCL.ID).Updates(mapStrIntfColumnsUpdate)
	if resp.Error != nil {
		logs.Info(resp.Error)
		erro = errors.New("error: al actualizar tabla de cierre de lote")
		return erro
	}
	if resp.RowsAffected == 0 {
		message := "no existe el cierre de lote con el id " + strconv.FormatUint(uint64(entityPrismaCL.ID), 10)
		logs.Info(message)
		erro = errors.New(message)
		return erro
	}
	return
}

func (r *repository) DeleteOneCierreLotePrismaRepository(id uint64) (erro error) {

	entidad := entities.Prismacierrelote{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.Delete(&entidad)

	// si resulta algun error de la base de datos
	if result.Error != nil {

		erro = fmt.Errorf("error en delete de cierre de lote")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteOneCierreLotePrismaRepository",
		}

		err := r.UtilRepository.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	return
}

func (r *repository) GetAllCierreLoteApiLinkRepository(filtro cierrelotedtos.ApilinkRequest) (listaApilinkCierreLote []entities.Apilinkcierrelote, totalFilas int64, erro error) {

	queryGorm := r.SQLClient.Model(entities.Apilinkcierrelote{})

	queryGorm.Unscoped().Where("created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)

	queryGorm.Preload("Pagoestadoexterno")

	// Paginacion
	if filtro.Number > 0 && filtro.Size > 0 {

		// Ejecutar y contar las filas devueltas
		queryGorm.Count(&totalFilas)

		if queryGorm.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
			return
		}

		offset := (filtro.Number - 1) * filtro.Size
		queryGorm.Limit(int(filtro.Size))
		queryGorm.Offset(int(offset))
	}

	queryGorm.Order("created_at desc").Find(&listaApilinkCierreLote)

	// capturar error query DB
	if queryGorm.Error != nil {

		erro = fmt.Errorf(ERROR_GET_CIERRE_LOTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetAllCierreLoteApiLinkRepository",
		}

		err := r.UtilRepository.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
