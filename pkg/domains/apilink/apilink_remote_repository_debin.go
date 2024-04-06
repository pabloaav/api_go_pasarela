package apilink

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
)

func (r *remoteRepository) CreateDebinApiLink(requerimientoId string, request linkdebin.RequestDebinCreateLink, token string) (response linkdebin.ResponseDebinCreateLink, erro error) {

	//NOTE descomentar esto en produccion
	// base, erro := buildUrlDebines("debines")

	// if erro != nil {
	// 	return
	// }

	// json_data, _ := json.Marshal(request)

	// req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	// buildHeaderAutorizacion(req, requerimientoId, token)

	// erro = executeRequest(r, req, ERROR_DEBIN, &response)
	// /*
	// 	 se registra la peticion realizada a la api de apilink
	// 		-	armo el request para registrar la peticion realizada
	// 		-	registro la peticion realizada
	// */
	// peticionApiLink := dtos.RequestWebServicePeticion{
	// 	Operacion: "CreateDebinApiLink",
	// 	Vendor:    "ApiLink",
	// }
	// err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	// if err1 != nil {
	// 	logs.Error(ERROR_CREAR_PETICION + err1.Error())
	// }
	// return

	// NOTE descomentar esto en desarrollo
	base, erro := buildUrlDebines("debines")
	logs.Info(base)

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(request)
	logs.Info(json_data)

	// req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))
	// logs.Info(req)

	// buildHeaderAutorizacion(req, requerimientoId, token)

	// erro = executeRequest(r, req, ERROR_DEBIN, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/

	response = linkdebin.ResponseDebinCreateLink{
		Id:              "53",
		FechaOperacion:  time.Now(),
		Estado:          "ACEPTADO",
		FechaExpiracion: time.Now(),
	}
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "CreateDebinApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}
	return
}

func (r *remoteRepository) GetDebinesApiLink(requerimientoId string, request linkdebin.RequestGetDebinesLink, token string) (response linkdebin.ResponseGetDebinesLink, erro error) {

	base, err := buildUrlDebines("debines")

	if err != nil {
		return response, err
	}

	params := url.Values{}
	params.Add("pagina", fmt.Sprint(request.Pagina))
	params.Add("tamanio", string(request.Tamanio))
	params.Add("cbu", request.Cbu)
	params.Add("estado", string(request.Estado))
	params.Add("escomprador", fmt.Sprint(request.EsComprador))
	params.Add("tipo", string(request.Tipo))
	// En este momento estamos guardando la fecha sin los milisegundos y por eso hago esta transformación
	// En caso de que se guarde en la base de datos lo microsegundo hay que modificar esta peticion.
	//Todo Hay que validar como vamos a trabajar con las fechas para que eso funcione correctamente.
	base.RawQuery = params.Encode() + fmt.Sprintf("&fechadesde=%s", request.FechaDesde.Format("2006-01-02T15:04:05.000Z")) +
		fmt.Sprintf("&fechahasta=%s", request.FechaHasta.Format("2006-01-02T15:04:05.000Z"))

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	err = executeRequest(r, req, ERROR_GET_DEBINES, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "GetDebinesApiLink",
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

func (r *remoteRepository) GetDebinApiLink(requerimientoId string, request linkdebin.RequestGetDebinLink, token string) (response linkdebin.ResponseGetDebinLink, erro error) {

	base, err := buildUrlDebines("debines")

	if err != nil {
		return response, err
	}

	base.Path += "/" + request.Id
	params := url.Values{}
	params.Add("cbu", request.Cbu)

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	err = executeRequest(r, req, ERROR_GET_DEBINES, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "GetDebinApiLink",
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

func (r *remoteRepository) GetDebinesPendientesApiLink(requerimientoId string, cbu string, token string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error) {

	base, err := buildUrlDebines("debines/pendientes")

	if err != nil {
		return response, err
	}

	params := url.Values{}
	params.Add("cbu", cbu)

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	err = executeRequest(r, req, ERROR_GET_DEBINES_PENDIENTES, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "GetDebinesPendientesApiLink",
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

func (r *remoteRepository) DeleteDebinApiLink(requerimientoId string, request linkdebin.RequestDeleteDebinLink, token string) (response bool, erro error) {

	base, err := buildUrlDebines("debines")

	if err != nil {
		return response, err
	}

	base.Path += "/" + request.Id
	params := url.Values{}
	params.Add("cbu", request.Cbu)

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("DELETE", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	buildHeaderDefault(req)

	resp, err := r.HTTPClient.Do(req)
	//Todo hay que comentar este codigo porque solo sirve para las pruebas de homologacion
	// logs.Info(req)
	// logs.Info(resp)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "DeleteDebinApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}

	if err != nil {
		logs.Error(err.Error())
		return response, errors.New(ERROR_DELETE_DEBINES)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		apiError := linkdtos.ErrorApiLink{}
		err := json.NewDecoder(resp.Body).Decode(&apiError)

		if resp.StatusCode == 401 {
			apiError.Codigo = "401"
			apiError.Descripcion = "Unauthorized"
			logs.Error(fmt.Sprintf("%s, %s", apiError.Error(), resp.Status))
			return false, &apiError
		}

		if resp.StatusCode == 403 {
			apiError.Codigo = "403"
			apiError.Descripcion = "Forbidden"
			logs.Error(fmt.Sprintf("%s, %s", apiError.Error(), resp.Status))
			return false, &apiError
		}

		if err != nil {
			logs.Error(fmt.Sprintf("%s, %s", err.Error(), resp.Status))
			return false, errors.New(ERROR_DELETE_DEBINES)
		}

		logs.Error(fmt.Sprintf("%s, %s", apiError.Error(), resp.Status))
		return false, &apiError
	}

	return true, nil
}

func buildUrlDebines(ruta string) (*url.URL, error) {

	base, err := url.Parse(config.APILINKDEBINHOST)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}

	base.Path += ruta

	return base, nil
}
