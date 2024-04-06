package usuario

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/userdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/usuario"
)

type UsuarioService interface {
	CreateUsuarioService(request userdtos.RequestUserAutorizacion) (id uint64, erro error)
	UpdateUsuarioService(request userdtos.RequestUserAutorizacion) (erro error)
	GetUsuarioService(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUser, erro error)
	GetUsuariosService(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUsers, erro error)
}

type usuarioService struct {
	remoteRepository UsuarioRemoteRepository
	repository       UsuarioRepository
}

func NewService(rm UsuarioRemoteRepository, r UsuarioRepository) UsuarioService {
	usuario := &usuarioService{
		remoteRepository: rm,
		repository:       r,
	}
	return usuario
}

func (s *usuarioService) CreateUsuarioService(request userdtos.RequestUserAutorizacion) (id uint64, erro error) {

	//Trata de crear el usuario en autenticacion
	//Si el usuario ya existe para el sistema devuelve su id
	//Si el usuario no existe en el sistema lo crea y devuelve el id
	id, erro = s.remoteRepository.CreateUsuario(request)

	if erro != nil {
		return
	}

	//Hago una consulta para ver si el usuario no esta registrado en la pasarela
	filtroClienteUser := filtros.ClienteUserFiltro{
		UserId:    id,
		ClienteId: request.Request.ClienteId,
	}

	response, erro := s.repository.GetClienteUser(filtroClienteUser)

	if erro != nil {
		//Si no está registrado lo registro y devuelvo su id de usuario
		if erro.Error() == RESULTADO_NO_ENCONTRADO {

			clienteUser := entities.Clienteuser{
				UserId:     id,
				ClientesId: uint64(request.Request.ClienteId),
			}

			_, erro = s.repository.CreateClienteUser(clienteUser)

			if erro != nil {
				return
			}

			return id, nil

		}

	}

	//Si el usuario ya está registrado entoces devuelve su id de usuario
	id = response.UserId

	return
}
func (s *usuarioService) UpdateUsuarioService(request userdtos.RequestUserAutorizacion) (erro error) {

	erro = s.remoteRepository.UpdateUsuario(request)

	if erro != nil {
		return
	}

	if request.RequestUpdate.ClienteIdNuevo > 0 {

		filtroClienteUser := filtros.ClienteUserFiltro{
			UserId:    request.RequestUpdate.Id,
			ClienteId: request.RequestUpdate.ClienteIdAnterior,
		}

		clienteUser, erro := s.repository.GetClienteUser(filtroClienteUser)

		if erro != nil {
			return erro
		}

		if clienteUser.ClientesId != request.RequestUpdate.ClienteIdNuevo {
			clienteUser.ClientesId = request.RequestUpdate.ClienteIdNuevo

			return s.repository.UpdateClienteUser(clienteUser)
		}

	}

	return
}
func (s *usuarioService) GetUsuarioService(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUser, erro error) {

	//Busco en la base de pasarela para saber si existe el usuario en este cliente
	//Si existe busco en la api de autenticacion caso contrario devuelve un error

	filtroClienteUser := filtros.ClienteUserFiltro{
		UserId:        filtro.User.Id,
		ClienteId:     filtro.User.ClienteId,
		CargarCliente: true,
	}

	if erro != nil {
		return
	}

	UserCliente, erro := s.repository.GetClienteUser(filtroClienteUser)

	if erro != nil {
		return
	}

	if UserCliente.ID > 0 {
		response, erro = s.remoteRepository.GetUsuario(filtro)

		if UserCliente.UserId == uint64(response.Id) {
			cliente := userdtos.ResponseUserCliente{}
			cliente.Id = UserCliente.ClientesId
			cliente.Cuit = UserCliente.Cliente.Cuit
			cliente.RazonSocial = UserCliente.Cliente.Razonsocial
			response.Cliente = &cliente
		}

		return
	}

	erro = fmt.Errorf(ERROR_CONSULTA_VACIA)

	return

}
func (s *usuarioService) GetUsuariosService(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUsers, erro error) {

	filtroClienteUser := filtros.ClienteUserFiltro{
		CargarCliente: true,
		ClienteId:     filtro.User.ClienteId,
	}

	clienteUsers, erro := s.repository.GetClienteUsers(filtroClienteUser)

	if erro != nil {
		return
	}
	if len(clienteUsers) > 0 {
		var listaIdsUsuario []uint64
		for i := range clienteUsers {
			listaIdsUsuario = append(listaIdsUsuario, clienteUsers[i].UserId)
		}

		filtro.User.Ids = &listaIdsUsuario
		filtro.User.CargarUserSistema = true

		response, erro = s.remoteRepository.GetUsuarios(filtro)
		for _, c := range clienteUsers {
			for i := range response.Data {
				if c.UserId == uint64(response.Data[i].Id) {
					cliente := userdtos.ResponseUserCliente{}
					cliente.Id = c.ClientesId
					cliente.Cuit = c.Cliente.Cuit
					cliente.RazonSocial = c.Cliente.Razonsocial
					response.Data[i].Cliente = &cliente
				}
			}
		}

		return
	}

	erro = fmt.Errorf(ERROR_CONSULTA_VACIA)

	return

}
