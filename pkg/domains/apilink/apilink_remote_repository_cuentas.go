package apilink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

func (r *remoteRepository) GetCuentasApiLink(request linkcuentas.LinkGetCuentasRequest) (response []linkcuentas.GetCuentasResponse, erro error) {

	base, err := buildUrlCuentas("cuentas")

	if err != nil {
		return response, err
	}

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, request.RequerimientoId, request.Token.AccessToken)

	err = executeRequestCuentas(r, req, ERROR_GET_CUENTAS, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "GetCuentasApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *remoteRepository) CreateCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error) {

	base, erro := buildUrlCuentas("cuentas")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(request.Request)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, request.RequerimientoId, request.Token.AccessToken)

	erro = executeRequestCuentas(r, req, ERROR_CREAR_CUENTA, nil)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "CreateCuentaApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}

	return
}

func (r *remoteRepository) DeleteCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error) {

	base, erro := buildUrlCuentas("cuentas")

	if erro != nil {
		return
	}

	valor := request.Request.(linkcuentas.LinkDeleteCuenta)
	params := url.Values{}
	params.Add("cbu", string(valor.Cbu))

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("DELETE", base.String(), nil)

	buildHeaderAutorizacion(req, request.RequerimientoId, request.Token.AccessToken)

	erro = executeRequestCuentas(r, req, ERROR_BAJAR_CUENTA, nil)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "DeleteCuentaApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}

	return
}

func executeRequestCuentas(r *remoteRepository, req *http.Request, erro string, objeto interface{}) error {

	buildHeaderDefault(req)

	req.Header.Add("x-timestamp", time.Now().Format("2006-01-02T15:04:05.000Z"))

	resp, err := r.HTTPClient.Do(req)

	//Todo hay que comentar este codigo porque solo sirve para las pruebas de homologacion
	if erro != ERROR_TOKEN {
		logs.Info(req)
		logs.Info(resp)
	}

	if err != nil {

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       err.Error(),
			Funcionalidad: "executeRequestCuentas",
		}

		r.UtilService.CreateLogService(log)

		return fmt.Errorf(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 204 {

		if resp.StatusCode == 200 {

			return json.NewDecoder(resp.Body).Decode(&objeto)

		}

		apiError := linkdtos.ErrorApiLink{}

		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "executeRequestCuentas",
		}

		if resp.StatusCode == 401 {
			apiError.Codigo = "401"
			apiError.Descripcion = "Unauthorized"

			log.Mensaje = apiError.Error()
			r.UtilService.CreateLogService(log)

			return &apiError
		}

		if resp.StatusCode == 403 {
			apiError.Codigo = "403"
			apiError.Descripcion = "Forbidden"

			log.Mensaje = apiError.Error()
			r.UtilService.CreateLogService(log)

			return &apiError
		}

		err := json.NewDecoder(resp.Body).Decode(&apiError)

		if err != nil {

			log.Mensaje = fmt.Sprintf("%s, %s", erro, resp.Status)

			r.UtilService.CreateLogService(log)

			return fmt.Errorf("%s, %s", erro, resp.Status)
		}

		log.Mensaje = apiError.Error()

		r.UtilService.CreateLogService(log)

		return &apiError
	}

	return nil

}

func buildUrlCuentas(ruta string) (*url.URL, error) {

	base, err := url.Parse(config.APILINKCUENTASHOST)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}

	base.Path += ruta

	return base, nil
}
