package checkout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type Repository interface {
	BeginTx()
	CommitTx()
	RollbackTx()
	CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error)
	UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error)
	UpdatePagoIntentoFallidos(ctx context.Context, pago *entities.Pago) error
	GetPagoByUuid(uuid string) (*entities.Pago, error)
	GetClienteByApikey(apikey string) (*entities.Cliente, error)
	GetCuentaByApikey(apikey string) (*entities.Cuenta, error)
	GetPagotipoById(id int64) (*entities.Pagotipo, error)
	GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error)
	GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error)
	GetChannelByName(nombre string) (*entities.Channel, error)
	GetCuentaById(id int64) (*entities.Cuenta, error)
	CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error)
	GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error)
	GetMediosDePagos() (*[]entities.Mediopago, error)
	GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error)
	GetInstallmentDetailsID(installmentID, numeroCuota int64) int64
	GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error)
	GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error)
	CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error
	// GetCuentaTelco() (*entities.Configuracione, error)
	GetPagoEstado(id int64) (*entities.Pagoestado, error)
	GetPreferencesByIdClienteRepository(id uint) (preferencia entities.Preference, erro error)
	GetChannelById(id uint) (channel entities.Channel, erro error)
	CheckUsuarioTRepository(usuario entities.Usuariobloqueados) (usuarioDB entities.Usuariobloqueados, err error)
	UpdateUsuarioBloqueoRepository(usuario entities.Usuariobloqueados) error
	DeleteUsuarioListaBloqueoRepository(usuario entities.Usuariobloqueados) error
	AgregarUsuarioListaBloqueoRepository(entities.Usuariobloqueados) error

	SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error)
	GetHasheado(hash string) (control bool, erro error)
}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService) Repository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
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

func (r *repository) auditarCheckout(ctx context.Context, resultado interface{}) error {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	audit.Operacion = strings.ToLower(audit.Query[:6])

	audit.Origen = "pasarela.checkout"

	res, _ := json.Marshal(resultado)
	audit.Resultado = string(res)

	err := r.auditoriaService.Create(&audit)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

func (r *repository) CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error) {
	res := r.SQLClient.WithContext(ctx).Create(&pago)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se pudo generar registro pago: " + res.Error.Error())
	}
	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return nil, err
	}
	return pago, nil
}

func (r *repository) UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error) {
	res := r.SQLClient.WithContext(ctx).Model(&pago).Updates(&pago)
	if res.Error != nil {
		return false, fmt.Errorf("al actualizar el estado del pago: %s", res.Error.Error())
	}
	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *repository) UpdatePagoIntentoFallidos(ctx context.Context, pago *entities.Pago) error {
	dataToUpdate := map[string]interface{}{"intento_fallido": pago.IntentoFallido}
	res := r.SQLClient.WithContext(ctx).Model(&pago).Select("intento_fallido").Updates(dataToUpdate)
	if res.Error != nil {
		return fmt.Errorf("al actualizar la columna del pago intento fallido: %s", res.Error.Error())
	}
	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetPagoByUuid(uuid string) (*entities.Pago, error) {
	var pago entities.Pago
	/*
		FIXME: se agrego el prelodad para que cargue los pagos intentos de un pago
	*/
	//res := r.SQLClient.Model(entities.Pago{}).Preload("Pagoitems").Where("pagoestados_id = 1 AND uuid = ?", uuid).Find(&pago)
	res := r.SQLClient.Model(entities.Pago{}).Preload("Pagoitems").Preload("PagoIntentos").Where("uuid = ?", uuid).Find(&pago)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no existe pago con identificador %s", uuid)
	}
	return &pago, nil
}

func (r *repository) GetClienteByApikey(apikey string) (*entities.Cliente, error) {

	var cliente entities.Cliente

	res := r.SQLClient.Preload("Cuentas.Pagotipos").Joins("JOIN cuentas on cuentas.clientes_id = clientes.id and cuentas.apikey = ?", apikey).Find(&cliente)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cliente con apikey: %s", apikey)
	}

	return &cliente, nil
}

func (r *repository) GetCuentaByApikey(apikey string) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta
	res := r.SQLClient.Preload("Pagotipos").Where("apikey = ?", apikey).Find(&cuenta)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con apikey: %s", apikey)
	}
	return &cuenta, nil
}
func (r *repository) GetPagotipoById(id int64) (*entities.Pagotipo, error) {
	var tipo entities.Pagotipo
	// res := r.SQLClient.Find(&tipo, id)
	res := r.SQLClient.Table("pagotipos").Where("id=?", id)
	res.Preload("Cuenta")
	res.Find(&tipo)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró tipo de pago con el id: %d", id)
	}
	return &tipo, nil
}

func (r *repository) CheckUsuarioTRepository(usuario entities.Usuariobloqueados) (usuarioDB entities.Usuariobloqueados, err error) {

	resp := r.SQLClient.Model(entities.Usuariobloqueados{})
	// if len(request.HolderName) > 0 {
	// 	resp.Where("nombre = ?", request.HolderName)
	// }
	if len(usuario.Email) > 0 {
		resp.Where("email = ?", usuario.Email)
	}
	resp.Find(&usuarioDB)
	if resp.Error != nil {
		logs.Error(resp.Error)
		return usuarioDB, fmt.Errorf("error al buscar usuario en la lista de bloqueo")
	}
	if resp.RowsAffected <= 0 {
		return usuarioDB, nil
	}

	return usuarioDB, nil

}
func (r *repository) AgregarUsuarioListaBloqueoRepository(entitie entities.Usuariobloqueados) error {
	resp := r.SQLClient.Create(&entitie)
	if resp.RowsAffected <= 0 {
		return fmt.Errorf("no se pudo agregar usuario a la lista de bloqueo")
	}
	return nil
}
func (r *repository) UpdateUsuarioBloqueoRepository(usuario entities.Usuariobloqueados) error {
	resp := r.SQLClient.Model(entities.Usuariobloqueados{}).Where("email = ?", usuario.Email)
	if !time.Time.IsZero(usuario.FechaBloqueo) {
		resp.Update("fecha_bloqueo", usuario.FechaBloqueo)
	}
	if usuario.CantBloqueo > 0 {
		resp.Update("cant_bloqueo", usuario.CantBloqueo)
	}
	if usuario.Permanente {
		resp.Update("permanente", usuario.Permanente)
	}
	if resp.RowsAffected <= 0 {
		return fmt.Errorf("no se pudo actualizar usuario en la lista de bloqueo")
	}
	return nil
}

func (r *repository) DeleteUsuarioListaBloqueoRepository(entitie entities.Usuariobloqueados) error {
	resp := r.SQLClient.Delete(&entitie)
	if resp.RowsAffected <= 0 {
		return fmt.Errorf("no se pudo eliminar usuario de la lista de bloqueo")
	}
	return nil
}

// NOTE - revisar esta funcion
func (r *repository) GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error) {
	var channels []entities.Pagotipochannel
	res := r.SQLClient.Preload("Channel").Where("pagotipos_id = ?", id).Find(&channels)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró tipo de pago con el id: %d", id)
	}
	return &channels, nil
}

// NOTE - esta funcion se modifico controla erro
func (r *repository) GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error) {
	var installmentdetails []entities.Pagotipointallment
	res := r.SQLClient.Where("pagotipos_id = ?", id).Find(&installmentdetails)
	if res.Error != nil {
		return nil, fmt.Errorf("no se encontró cuotas para el tipo de pago con el id: %d", id)
	}
	// if res.RowsAffected <= 0 {
	// 	return nil, fmt.Errorf("no se encontró cuotas para el tipo de pago con el id: %d", id)
	// }
	return &installmentdetails, nil
}

func (r *repository) GetChannelByName(nombre string) (*entities.Channel, error) {
	var channel entities.Channel

	res := r.SQLClient.Where("channel = ?", nombre).Find(&channel)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró metodo de pago con la descripción %s", nombre)
	}

	return &channel, nil
}

func (r *repository) GetCuentaById(id int64) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta

	res := r.SQLClient.Preload("Cliente").Find(&cuenta, id)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con el id: %d", id)
	}
	// Incluir datos del cliente

	return &cuenta, nil
}

func (r *repository) CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error) {

	res := r.SQLClient.WithContext(ctx).Create(&resultado)
	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al guardar resultado: %s", res.Error.Error())
	}

	err := r.auditarCheckout(res.Statement.Context, resultado.StateComment)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error) {
	var intento entities.Pagointento
	res := r.SQLClient.Model(entities.Pagointento{}).Where("external_id != '0' AND pagos_id = ?", pagoId).Last(&intento)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró intento con el id de pago: %d", pagoId)
	}

	return &intento, nil
}

func (r *repository) GetMediosDePagos() (*[]entities.Mediopago, error) {
	var medios []entities.Mediopago
	res := r.SQLClient.Model(entities.Mediopago{})
	res.Preload("Mediopagoinstallment")
	res.Preload("Channel")
	res.Where("mediopagos.regexp != ''").Order("longitud_pan DESC")
	res.Order("codigo_bcra")
	res.Find(&medios)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron medios de pago")
	}
	return &medios, nil
}

func (r *repository) GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error) {
	var medio entities.Mediopago
	res := r.SQLClient.Model(entities.Mediopago{})
	res.Where(filtro)
	res.First(&medio)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron medios de pago")
	}
	return &medio, nil
}

func (r *repository) GetInstallmentDetailsID(installmentID, numeroCuota int64) int64 {
	var response int64
	res := r.SQLClient.Model(entities.Installmentdetail{})
	res.Select("id")
	res.Where("installments_id = ? AND cuota = ?", installmentID, numeroCuota)
	res.First(&response)
	if res.RowsAffected <= 0 {
		return 1
	}
	return response
}

func (r *repository) GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error) {
	var result *entities.Installmentdetail
	res := r.SQLClient.Model(entities.Installmentdetail{})
	res.Where("installments_id = ? AND cuota = ?", installmentID, numeroCuota)
	res.First(&result)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se encontraron detalles de cuotas ")
		return
	}
	installmentDetails = &dtos.InstallmentDetailsResponse{
		Id:             result.Model.ID,
		InstallmentsID: result.InstallmentsID,
		NroCuota:       result.Cuota,
		Coeficiente:    result.Coeficiente,
	}
	return
}

func (r *repository) GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error) {
	// and vigencia_hasta is null
	res := r.SQLClient.Model(entities.Installment{})
	res.Where("mediopagoinstallments_id = ? ", id).Order("created_at asc")
	res.Find(&installments)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se encontro plan de cuotas ")
		return
	}
	return
}

func (r *repository) CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error {
	res := r.SQLClient.WithContext(ctx).Create(&pel)
	if res.RowsAffected <= 0 {
		return fmt.Errorf("error al guardar estado log: %s", res.Error.Error())
	}

	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return err
	}

	return nil
}

// func (r *repository) GetCuentaTelco() (*entities.Configuracione, error) {
// 	var cuentaTelco entities.Configuracione
// 	res := r.SQLClient.Model(&cuentaTelco).Where("nombre = ?", "CBU_CUENTA_TELCO").First(cuentaTelco)
// 	if res.RowsAffected <= 0 {
// 		err := errors.New("error al obtener la cuenta cbu telco")
// 		return nil, err
// 	}
// 	return &cuentaTelco, nil
// }

// obtener estado actual de un pago
func (r *repository) GetPagoEstado(id int64) (*entities.Pagoestado, error) {
	var estado entities.Pagoestado

	res := r.SQLClient.Find(&estado, id)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con el id: %d", id)
	}

	return &estado, nil
}

func (r *repository) GetPreferencesByIdClienteRepository(id uint) (preference entities.Preference, erro error) {
	res := r.SQLClient.Table("preferences").Where("clientes_id = ?", id).Preload("Cliente")

	res.Last(&preference)
	if res.RowsAffected == 0 {
		return
	}
	if res.Error != nil {
		erro = errors.New("error en la consulta a preferencias de cliente")
		return
	}

	return
}

func (r *repository) GetChannelById(id uint) (channel entities.Channel, erro error) {
	res := r.SQLClient.Model(entities.Channel{}).Where("id = ?", id).Find(&channel)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se pudo encontrar channel relacionados")
		return
	}
	return
}

func (r *repository) SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error) {
	res := r.SQLClient.Where("uuid = ?", hasheado.Uuid).FirstOrCreate(&hasheado)

	uuid_pagointento := entities.UuidsPagointento{
		UuidsId:        hasheado.ID,
		PagointentosId: pagointento_id,
	}

	res2 := r.SQLClient.Create(&uuid_pagointento)
	if res2.RowsAffected <= 0 {
		return fmt.Errorf("error al guardar uuid_pi: %s", res.Error.Error())
	}

	return
}

func (r *repository) GetHasheado(hash string) (control bool, erro error) {
	var coindidencias []entities.Uuid
	res := r.SQLClient.Model(entities.Uuid{}).Where("uuid = ?", hash)

	res.Where("fecha_bloqueo > '0000-00-00 00:00:00'")

	res.Find(&coindidencias).Limit(1)

	if res.Error != nil {
		return false, fmt.Errorf("error al buscar hash: %s", res.Error.Error())
	}
	if res.RowsAffected > 0 {
		control = true
	}

	return
}

/*!SECTION
var intento entities.Pagointento
res := r.SQLClient.Model(entities.Pagointento{}).Where("external_id != '0' AND pagos_id = ?", pagoId).Last(&intento)
if res.RowsAffected <= 0 {
	return nil, fmt.Errorf("no se encontró intento con el id de pago: %d", pagoId)
}

return &intento, nil

*/
