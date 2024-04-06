package administracion

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

func (r *repository) GetPrismaCierreLotes(reversion bool) (prismaCierreLotes []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as pcl")

	if !reversion {
		resp.Where("pcl.match = 1")
	}
	if reversion {
		resp.Unscoped().Where("pcl.conciliado = 1 and pcl.detallepago_id > 0 and pcl.detallemovimiento_id > 0 and pcl.extbancoreversion_id <> 0 and pcl.estadomovimiento = 0")
	}
	// err := r.SQLClient.Table("prismacierrelotes as pcl").Where("pcl.match = 1").Find(&prismaCierreLotes).Error

	resp.Find(&prismaCierreLotes)
	if resp.Error != nil {
		logs.Error("error al querer obtener ierres de lotes")
		erro = errors.New(ERROR_OBTENER_CIERRE_LOTE)
		return
	}
	return
}

func (r *repository) ObtenerCierreLoteEnDisputaRepository(estadoDisputa int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as pcl").Unscoped().Where("pcl.match = 1 and pcl.disputa = ?", estadoDisputa)
	if filtro.FechaCreacion {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaOperacion {
		resp.Where("cast(fechaoperacion as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaPago {
		resp.Where("cast(fecha_pago as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaCierre {
		resp.Where("cast(fecha_cierre as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}

	err := resp.Find(&enttyClEnDsiputa).Error
	if err != nil {
		erro = errors.New(ERROR_OBTENER_CIERRE_LOTE)
		return
	}

	return
}
func (r *repository) ObtenerCierreLoteContraCargoRepository(estadoReversion int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as pcl").Unscoped().Where("pcl.match = 1 and pcl.reversion = ?", estadoReversion)
	if filtro.FechaCreacion {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaOperacion {
		resp.Where("cast(fechaoperacion as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaPago {
		resp.Where("cast(fecha_pago as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	if filtro.FechaCierre {
		resp.Where("cast(fecha_cierre as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
	}
	err := resp.Find(&enttyClEnDsiputa).Error
	if err != nil {
		erro = errors.New(ERROR_OBTENER_CIERRE_LOTE)
		return
	}
	return
}

/*
func (r *repository) GetPrismaCierreLotes() (prismaCierreLotes []entities.Prismacierrelote, erro error) {
	err := r.SQLClient.Find(&prismaCierreLotes).Error
	if err != nil {
		erro = errors.New(ERROR_OBTENER_CIERRE_LOTE)
		return
	}
	return
}
*/
