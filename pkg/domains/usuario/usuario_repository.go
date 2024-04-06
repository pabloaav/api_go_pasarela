package usuario

import (
	"errors"
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/usuario"
	"gorm.io/gorm"
)

type UsuarioRepository interface {
	CreateClienteUser(request entities.Clienteuser) (id uint64, erro error)
	UpdateClienteUser(request entities.Clienteuser) (erro error)
	GetClienteUser(filtro filtros.ClienteUserFiltro) (response entities.Clienteuser, erro error)
	GetClienteUsers(filtro filtros.ClienteUserFiltro) (response []entities.Clienteuser, erro error)
}

type usuarioRepository struct {
	SQLClient *database.MySQLClient
	Util      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, u util.UtilService) UsuarioRepository {
	return &usuarioRepository{
		SQLClient: sqlClient,
		Util:      u,
	}
}

func (r *usuarioRepository) CreateClienteUser(request entities.Clienteuser) (id uint64, erro error) {

	result := r.SQLClient.Omit("id").Create(&request)

	if result.Error != nil {

		erro = fmt.Errorf(ErrorGuardar("ClienteUser", false))

		r.Util.LogError(result.Error.Error(), "CreateClienteUser")

		return

	}

	id = uint64(request.ID)

	return
}
func (r *usuarioRepository) UpdateClienteUser(request entities.Clienteuser) (erro error) {

	entidad := entities.Clienteuser{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CLIENTEUSER)

		r.Util.LogError(result.Error.Error(), "UpdateClienteUser")

		return
	}

	return
}

func (r *usuarioRepository) GetClienteUser(filtro filtros.ClienteUserFiltro) (response entities.Clienteuser, erro error) {

	resp := r.SQLClient.Model(entities.Clienteuser{})

	resp.Where("user_id = ? AND clientes_id = ?", filtro.UserId, filtro.ClienteId)

	_filtrosComunesUser(filtro, resp)

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_USUARIO)

		r.Util.LogError(resp.Error.Error(), "GetClienteUser")

	}

	return
}
func (r *usuarioRepository) GetClienteUsers(filtro filtros.ClienteUserFiltro) (response []entities.Clienteuser, erro error) {

	resp := r.SQLClient.Model(entities.Clienteuser{})

	if filtro.ClienteId > 0 {
		resp.Where("clientes_id = ?", filtro.ClienteId)
	}

	if filtro.UserId > 0 {
		resp.Where("user_id  = ?", filtro.UserId)
	}

	_filtrosComunesUser(filtro, resp)

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_USUARIO)

		r.Util.LogError(resp.Error.Error(), "GetClienteUsers")

	}

	return
}

func _filtrosComunesUser(filtro filtros.ClienteUserFiltro, resp *gorm.DB) {

	if filtro.CargarCliente {
		resp.Preload("Cliente")
	}

}
