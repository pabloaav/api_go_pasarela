package prisma

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	// SaveCierreLote guarda los cierres de lostes em la base de datos
	SaveCierreLote(detalleLote *entities.Prismacierrelote) (bool, error)
	SaveCierreLoteBatch(detalleLote []entities.Prismacierrelote) (bool, error)
	GetPagosPagosIntentosxChannel(estadoPago int, channel string) (pagoIntento []entities.Pagointento, erro error)
	GetMensajeErrorPrismaByExternalId(external_id uint64) (msgErrorPrisma entities.Prismaerroresexterno, erro error)
}

type repository struct {
	SQLClient *database.MySQLClient
}

func NewRepository(sqlClient *database.MySQLClient) Repository {
	return &repository{
		SQLClient: sqlClient,
	}
}

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
	err := tx.Create(&detalleLote).Error
	if err != nil {
		tx.Rollback()
		logs.Info(err)
		return false, errors.New("error: al inserte el valor en la base de datos ")
	}
	err = tx.Commit().Error
	if err != nil {
		logs.Info(err)
		return false, errors.New("error: al confirmar una transacción ")
	}
	return true, nil
}

func (r *repository) GetPagosPagosIntentosxChannel(estadoPago int, channel string) (pagoIntento []entities.Pagointento, erro error) {
	resp := r.SQLClient.Table("pagointentos as pi").Select("*").
		Joins("INNER JOIN pagos as p ON p.id = pi.pagos_id and p.pagoestados_id = ?", estadoPago).Preload("Pago").
		Joins("INNER JOIN mediopagos as mp ON mp.id = pi.mediopagos_id ").Preload("Mediopagos").
		Joins("INNER JOIN channels as ch ON ch.id = mp.channels_id  and ch.channel = ?", channel).Preload("Mediopagos.Channel").
		Where("pi.ticket_number <> '' and pi.authorization_code <> ''").Find(&pagoIntento)
	if resp != nil {
		logs.Error(resp)
	}
	return
}

func (r *repository) GetMensajeErrorPrismaByExternalId(external_id uint64) (msgErrorPrisma entities.Prismaerroresexterno, erro error) {
	resp := r.SQLClient.Model(entities.Prismaerroresexterno{}).Where("external_id = ?", external_id).Find(&msgErrorPrisma)
	if resp.Error != nil {
		erro = errors.New(ERROR_CONSULTAR_MSG_ERROR_PRISMA)
		logs.Error(resp.Error)
		return
	}
	if resp.Error != nil {
		erro = errors.New(ERROR_MSG_VACIO)
		logs.Error(erro.Error())
		return
	}
	return
}
