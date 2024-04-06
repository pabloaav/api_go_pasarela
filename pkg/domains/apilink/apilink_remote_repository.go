package apilink

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RemoteRepository interface {
	/* 	Genera los tokens con los que autorizarás las operaciones de las APIs dentro ApiLink
	https://portal.api.redlink.com.ar/redlink/sb/node/1126
	*/
	GetTokenApiLink(identificador string, scope []linkdtos.EnumScopeLink) (linkdtos.TokenLink, error)

	/*
		DEBIN es un medio de pago cuyo nombre alude a “débito inmediato”.
		Este medio permite a las entidades financieras y a nuevos actores de la industria de medios
		de pago a debitar fondos de las cuentas bancarias de sus clientes.
		Se incorpora el concepto de recurrencia del DEBIN, que permite crear un débito inmediato hacía el
		comprador con la intención de que se repita la operación de manera regular.
		La frecuencia con la que se estarán generando los DEBINes dependerá del ciclo de facturación del vendedor.
		https://portal.api.redlink.com.ar/redlink/sb/node/2679
	*/
	CreateDebinApiLink(requerimientoId string, request linkdebin.RequestDebinCreateLink, token string) (response linkdebin.ResponseDebinCreateLink, erro error)
	GetDebinesApiLink(requerimientoId string, request linkdebin.RequestGetDebinesLink, token string) (response linkdebin.ResponseGetDebinesLink, erro error)
	GetDebinApiLink(requerimientoId string, request linkdebin.RequestGetDebinLink, token string) (response linkdebin.ResponseGetDebinLink, erro error)
	GetDebinesPendientesApiLink(requerimientoId string, cbu string, token string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error)
	DeleteDebinApiLink(requerimientoId string, request linkdebin.RequestDeleteDebinLink, token string) (response bool, erro error)

	CreateCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error)
	DeleteCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error)
	GetCuentasApiLink(request linkcuentas.LinkGetCuentasRequest) (response []linkcuentas.GetCuentasResponse, erro error)

	CreateTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink, token string) (response linktransferencia.ResponseTransferenciaCreateLink, erro error)
	GetTransferenciasApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink, token string) (response linktransferencia.ResponseGetTransferenciasLink, erro error)
	GetTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink, token string) (response linktransferencia.ResponseGetTransferenciaLink, erro error)

	GetConsultaDestinatario(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink, token string) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error)
}

type remoteRepository struct {
	HTTPClient  *http.Client
	UtilService util.UtilService
}

func NewRemote(http *http.Client, u util.UtilService) RemoteRepository {
	return &remoteRepository{
		HTTPClient:  http,
		UtilService: u,
	}
}

func (r *remoteRepository) GetTokenApiLink(identificador string, scope []linkdtos.EnumScopeLink) (token linkdtos.TokenLink, erro error) {

	base, err := url.Parse(config.APILINKAUTENTICACIONHOST)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return token, err
	}

	base.Path += "clientcredential"

	var requestParams struct {
		Scope string `json:"scope"`
	}

	scopes := ""

	for _, s := range scope {
		scopes += fmt.Sprintf("%v,", s)
	}

	requestParams.Scope = scopes[:len(scopes)-1]

	json_data, _ := json.Marshal(requestParams)
	// logs.Info(json_data)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderToken(req, identificador)

	err = executeRequest(r, req, ERROR_TOKEN, &token)

	logs.Info(err)

	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "Autenticacion(genera token)",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}

	if err != nil {
		return token, err
	}
	return token, nil
}

func buildHeaderToken(request *http.Request, identificador string) {
	request.Header.Add("x-ibm-client-secret", config.SECRETLINK)
	request.Header.Add("x-idrequerimiento", identificador)
	request.Header.Add("x-ipcliente", config.APP_HOST)
	request.Header.Add("x-timestamp", time.Now().String())
}

func buildHeaderAutorizacion(request *http.Request, requerimientoId string, token string) {
	request.Header.Add("x-idrequerimiento", requerimientoId)
	request.Header.Add("authorization", "Bearer "+token)
}

func buildHeaderDefault(request *http.Request) {
	request.Header.Add("x-ibm-client-id", config.IDCLIENTLINK)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("accept", "application/json")
}

func executeRequest(r *remoteRepository, req *http.Request, erro string, objeto interface{}) error {

	buildHeaderDefault(req)

	resp, err := r.HTTPClient.Do(req)

	//TODO: hay que comentar este codigo porque solo sirve para las pruebas de homologacion
	// if erro != ERROR_TOKEN {
	// 	logs.Info(req)
	// 	logs.Info(resp)
	// }

	if err != nil {
		logs.Error(err.Error())
		return errors.New(erro)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		logs.Error(fmt.Sprint(resp.Status))
		return nil
	}

	if resp.StatusCode != 200 && resp.StatusCode != 202 {

		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "executeRequestCuentas",
		}

		apiError := linkdtos.ErrorApiLink{}

		if resp.StatusCode == 500 {

			apiError.Codigo = "500"

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprint(resp.Status)

		}

		err := json.NewDecoder(resp.Body).Decode(&apiError)

		if err != nil {

			apiError.Codigo = strconv.Itoa(resp.StatusCode)

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprintf("%s, %s", erro, resp.Status)

		}

		if resp.StatusCode == 401 {
			apiError.Codigo = "401"
			apiError.Descripcion = "Unauthorized"
		}

		r.UtilService.CreateLogService(log)

		return &apiError
	}

	err = json.NewDecoder(resp.Body).Decode(&objeto)

	if err != nil {
		return err
	}

	return nil

}
