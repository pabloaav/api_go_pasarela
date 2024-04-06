package usuario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/userdtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/usuario"
)

type UsuarioRemoteRepository interface {
	CreateUsuario(request userdtos.RequestUserAutorizacion) (id uint64, erro error)
	UpdateUsuario(request userdtos.RequestUserAutorizacion) (erro error)
	GetUsuario(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUser, erro error)
	GetUsuarios(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUsers, erro error)
}

type usuarioRemoteRepository struct {
	HTTPClient *http.Client
}

func NewRemote(http *http.Client) UsuarioRemoteRepository {
	return &usuarioRemoteRepository{
		HTTPClient: http,
	}
}

func (r *usuarioRemoteRepository) CreateUsuario(request userdtos.RequestUserAutorizacion) (id uint64, erro error) {

	base, erro := _buildUrlUser("adm/user-create")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(request.Request)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, request.Token)

	response := userdtos.ResponseUser{}

	erro = executeRequest(r, req, ERROR_CREAR_USUARIO, &response)

	if erro != nil {
		return
	}

	id = uint64(response.Id)

	return
}
func (r *usuarioRemoteRepository) UpdateUsuario(request userdtos.RequestUserAutorizacion) (erro error) {

	base, erro := _buildUrlUser("adm/user-update")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(request.RequestUpdate)

	req, _ := http.NewRequest("PUT", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, request.Token)

	return executeRequest(r, req, ERROR_MODIFICAR_USUARIO, nil)

}

func (r *usuarioRemoteRepository) GetUsuario(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUser, erro error) {

	base, erro := _buildUrlUser("adm/user")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(filtro.User)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, filtro.Token)

	erro = executeRequest(r, req, ERROR_CARGAR_USUARIO, &response)

	return
}

func (r *usuarioRemoteRepository) GetUsuarios(filtro filtros.UserFiltroAutenticacion) (response userdtos.ResponseUsers, erro error) {
	base, erro := _buildUrlUser("adm/users")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(filtro.User)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, filtro.Token)

	erro = executeRequest(r, req, ERROR_CARGAR_USUARIO, &response)

	return
}

func _buildUrlUser(ruta string) (*url.URL, error) {

	base, err := url.Parse(config.AUTH)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}

	base.Path += ruta

	return base, nil
}

func buildHeaderAutorizacion(request *http.Request, token string) {
	request.Header.Add("authorization", token)
}

func buildHeaderDefault(request *http.Request) {
	request.Header.Add("content-type", "application/json")
	request.Header.Add("accept", "application/json")
}

func executeRequest(r *usuarioRemoteRepository, req *http.Request, erro string, objeto interface{}) error {

	buildHeaderDefault(req)

	resp, err := r.HTTPClient.Do(req)

	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf(erro)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		apiError := ErrorUser{}
		err := json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			logs.Error(fmt.Sprintf("%s, %s", erro, resp.Status))
			return fmt.Errorf("%s, %s", erro, resp.Status)
		}

		logs.Error(apiError.Error())
		return &apiError
	}

	return json.NewDecoder(resp.Body).Decode(&objeto)

}
