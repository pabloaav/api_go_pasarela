package prisma

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	config "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

////////////////////////////////////////transacciones simples y offline/////////////////////
///////////////////////////////////PROCESO DE PAGO CON TARJETA//////////////////////////////////////
/////////////Inicio solicitud token de pago
func (r *remoteRepository) PostSolicitudTokenPago(card *prismadtos.Card) (response *prismadtos.PagoToken, err error) {
	var erro ErrorEstructura
	cardJson, _ := json.Marshal(card)
	//logs.Info(string(cardJson))
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_TOKENS
	//logs.Info(base.String())
	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(cardJson))
	buildHeaderDefault(req, config.PUBLIC_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al solicitar token: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		json.NewDecoder(resp.Body).Decode(&erro)
		logs.Error("Error en la peticion" + err.Error())
		return response, &erro
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

/////////////Fin solicitud token de pago
/////////////Inio ejecucion de pago
func (r *remoteRepository) PostEjecutarPago(procesoPagoRequest *prismadtos.PaymentsSimpleRequest) (response *prismadtos.PaymentsSimpleResponse, err error) {
	var erro ErrorEstructura
	requestJson, _ := json.Marshal(procesoPagoRequest)
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_PAYMENTS
	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(requestJson))
	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al ejecutar el pago: " + err.Error())
	}
	defer resp.Body.Close()
	var opcion int = resp.StatusCode
	switch opcion {
	case 201:
		json.NewDecoder(resp.Body).Decode(&response)
		return response, nil
	case 402:
		json.NewDecoder(resp.Body).Decode(&response)
		respuesta, _ := json.Marshal(response)
		logs.Info(string(respuesta))
		external_id := response.StatusDetails.Error.Reason.ID
		resultadoMsgError, _ := r.repository.GetMensajeErrorPrismaByExternalId(uint64(external_id))
		msjErrorStr := fmt.Sprintf(" error datos tarjeta- 402 - Code: %v, Param:  %v , Descripcion: %v", resultadoMsgError.ID, resultadoMsgError.Descripcion, resultadoMsgError.DescripcionmsgUsuario)
		return nil, errors.New(msjErrorStr)
	default:
		json.NewDecoder(resp.Body).Decode(&erro)
		err = errors.New(erro.BuscarMensajeError() + " - " + erro.CodeParamsError())
		return nil, err
	}

}

/*
400 	malformed_request_error 	Error en el armado del json
401 	authentication_error 	ApiKey Inválido
402 	invalid_request_error 	Error por datos inválidos
404 	not_found_error 	Error con datos no encontrados
409 	api_error 	Error inesperado en la API REST
*/

/////////////fin ejecucion de pago
/////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////PROCESO DE PAGO OFFLINE//////////////////////////////////////
/////////////inicio solicitud de token offline
// func (r *remoteRepository) PostSolicitarTokenOffLine(tokenOffline *prismadtos.OfflineTokenRequest) (response *prismadtos.OfflineTokenResponse, err error) {

// 	var erro ErrorEstructura
// 	dataJson, _ := json.Marshal(tokenOffline)
// 	//logs.Info(dataJson)
// 	base, err := url.Parse(config.URL_PRISMA)
// 	if err != nil {
// 		logs.Error("Error al crear base url" + err.Error())
// 	}
// 	base.Path += config.URI_TOKENS
// 	//logs.Info(base.String())
// 	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(dataJson))
// 	buildHeaderDefault(req, config.PUBLIC_OFFLINE_APIKEY_PRISMA)
// 	resp, err := r.HTTPClient.Do(req)
// 	if err != nil {
// 		logs.Error("error al solicitar token: " + err.Error())
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 201 {
// 		json.NewDecoder(resp.Body).Decode(&erro)
// 		//modifique esto
// 		logs.Error("Error en la peticion" + err.Error())
// 		return nil, &erro

// 		// return response, &erro
// 	}
// 	json.NewDecoder(resp.Body).Decode(&response)
// 	return response, nil
// }

/////////////Fin solicitud token de pago offline
/////////////Inio ejecucion de pago offline
// func (r *remoteRepository) PostEjecutarPagoOffLine(paymentOffline *prismadtos.PaymentsOfflineRequest) (response *prismadtos.PaymentsOfflineResponse, err error) {
// 	var erro ErrorEstructura
// 	dataJson, _ := json.Marshal(paymentOffline)
// 	//logs.Info(string(dataJson))
// 	base, err := url.Parse(config.URL_PRISMA)
// 	if err != nil {
// 		logs.Error("Error al crear base url" + err.Error())
// 	}
// 	base.Path += config.URI_PAYMENTS
// 	//logs.Info(base.String())
// 	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(dataJson))
// 	buildHeaderDefault(req, config.PRIVATE_OFFLINE_APIKEY_PRISMA)
// 	resp, err := r.HTTPClient.Do(req)
// 	if err != nil {
// 		logs.Error("error al ejecutar el pago off-line: " + err.Error())
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 201 {
// 		json.NewDecoder(resp.Body).Decode(&erro)
// 		fmt.Println(erro)
// 		err = errors.New(erro.BuscarMensajeError() + erro.CodeParamsError())
// 		return nil, err
// 	}
// 	json.NewDecoder(resp.Body).Decode(&response)
// 	return response, nil
// }

/////////////fin ejecucion de pago offline
