package prisma

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	config "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
)

func (r *remoteRepository) PostSolicitudAnulacionDevolucionPagoTotla(paymentId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, err error) {
	var erro ErrorEstructura
	paymentIdJson, _ := json.Marshal(paymentId)
	//logs.Info(paymentIdJson)
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_PAYMENTS + "/" + paymentId + config.URI_REFUNDS
	//logs.Info(base.String())
	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(paymentIdJson))

	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al solicitar anulación o devolucion de pago: " + err.Error())
	}
	defer resp.Body.Close()
	var opcion int = resp.StatusCode
	switch opcion {
	case 201:
		json.NewDecoder(resp.Body).Decode(&response)
		return response, nil
	case 402:
		json.NewDecoder(resp.Body).Decode(&response)
		return nil, response.StatusDetails.Error
	default:
		json.NewDecoder(resp.Body).Decode(&erro)
		err = errors.New(erro.BuscarMensajeError() + erro.CodeParamsError())
		fmt.Println(erro)
		return nil, err
	}

	// if resp.StatusCode != 201 {
	// 	json.NewDecoder(resp.Body).Decode(&erro)
	// 	err = errors.New(erro.BuscarMensajeError() + erro.CodeParamsError())
	// 	fmt.Println(erro)
	// 	return nil, err
	// }
	// json.NewDecoder(resp.Body).Decode(&response)
	// fmt.Printf("- Id : %v - Amount %v - SubPayments %v \n", response.ID, response.Amount, response.SubPayments)
	// return response, nil
}

//falta terminar esto eliminar
func (r *remoteRepository) DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, err error) {
	var erro ErrorEstructura
	//resultJson, _ := json.Marshal(erro)
	//logs.Info(resultJson)
	payload := strings.NewReader("{}")
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_PAYMENTS + "/" + params.ExternalId + config.URI_REFUNDS + "/" + params.RefundId
	//logs.Info(base.String())
	req, _ := http.NewRequest("DELETE", base.String(), payload)
	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al intentar eliminar anulación o devolucion de pago solicitado: " + err.Error())
	}
	defer resp.Body.Close()
	var opcion int = resp.StatusCode
	switch opcion {
	case 200:
		json.NewDecoder(resp.Body).Decode(&response)
		return response, nil
	default:
		json.NewDecoder(resp.Body).Decode(&erro)
		err = errors.New(erro.BuscarMensajeError() + erro.CodeParamsError())
		fmt.Println(erro)
		return nil, err
	}

	// if resp.StatusCode != 200 {
	// 	json.NewDecoder(resp.Body).Decode(&erro)
	// 	fmt.Println(erro)
	// 	return nil, &erro
	// }
	// json.NewDecoder(resp.Body).Decode(&response)
	// return response, nil
}

func (r *remoteRepository) PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, err error) {
	var erro ErrorEstructura
	var paramAmount prismaOperaciones.ParamAmountOperacion
	paramAmount.Amount = params.Monto
	resultJson, _ := json.Marshal(paramAmount)
	//logs.Info(resultJson)
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_PAYMENTS + "/" + params.ExternalId + config.URI_REFUNDS
	//logs.Info(base.String())
	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(resultJson))
	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al cancelar solicitud anulación o devolucion de pago: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		json.NewDecoder(resp.Body).Decode(&erro)
		err = errors.New(erro.BuscarMensajeError() + erro.CodeParamsError())
		fmt.Println(erro)
		return nil, err
	}
	json.NewDecoder(resp.Body).Decode(&response)
	fmt.Printf("- Id : %v - Amount %v - Status %v \n", response.ID, response.Amount, response.Status)
	return response, nil
}
