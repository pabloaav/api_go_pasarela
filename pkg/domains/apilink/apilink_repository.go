package apilink

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
)

type ApilinkRepository interface {
	PutApilinkCierrelote(listaDebinesId []string) (erro error)
}

type apilinkRepository struct {
	SQLClient *database.MySQLClient
	Util      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, u util.UtilService) ApilinkRepository {
	return &apilinkRepository{
		SQLClient: sqlClient,
		Util:      u,
	}
}

func (r *apilinkRepository) PutApilinkCierrelote(listaDebinesId []string) (erro error) {
	tx := r.SQLClient.Begin()
	err := tx.Table("apilinkcierrelotes").Where("debin_id IN (?)", listaDebinesId).Updates(map[string]interface{}{"match": 1}).Error
	if err != nil {
		tx.Rollback()
		logs.Info(err)
		return errors.New("error: al actualizar el cierre de lote apilink")
	}
	err = tx.Commit().Error
	if err != nil {
		logs.Info(err)
		return errors.New("error: al confirmar una transacción")
	}
	return nil
}
