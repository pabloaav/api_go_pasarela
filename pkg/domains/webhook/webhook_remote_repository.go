package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/webhook"
)

type RemoteRepository interface {
	NotificarPago(Pago dtos.ResultadoResponseWebHook) (erro error)
	NotificarPagos(listaPagos webhook.WebhookResponse) (erro error)
}

type remoteRepository struct {
	HTTPClient *http.Client
}

func NewRemote(http *http.Client) RemoteRepository {
	return &remoteRepository{
		HTTPClient: http,
	}
}

func (r *remoteRepository) NotificarPago(Pago dtos.ResultadoResponseWebHook) (erro error) {

	requestJson, _ := json.Marshal(Pago.ResultadoResponse)

	urlBase, err := buildUrlBanco(Pago.Url)
	if err != nil {
		erro = fmt.Errorf(err.Error())
		return
	}
	req, _ := http.NewRequest("POST", urlBase.String(), bytes.NewBuffer(requestJson))
	buildHeaderDefault(req)
	if len(Pago.ApikeyExterno) > 0 {
		req.Header.Add("apikey", Pago.ApikeyExterno)
	}
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		erro = err
		logs.Error("error al realizar peticion: " + err.Error())
		return erro
	}
	if resp.StatusCode != 200 {
		erro = fmt.Errorf("error al realizar peticion: " + resp.Status)
		return erro
	}
	defer resp.Body.Close()

	return
}

func (r *remoteRepository) NotificarPagos(listaPagos webhook.WebhookResponse) (erro error) {

	requestJson, _ := json.Marshal(listaPagos.ResultadoResponseWebHook)

	urlBase, err := buildUrlBanco(listaPagos.Url)
	if err != nil {
		erro = fmt.Errorf(err.Error())
		return
	}
	req, _ := http.NewRequest("POST", urlBase.String(), bytes.NewBuffer(requestJson))
	buildHeaderDefault(req)
	if len(listaPagos.ApikeyExterno) > 0 {
		req.Header.Add("apikey", listaPagos.ApikeyExterno)
	}
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		erro = err
		logs.Error("error al realizar peticion: " + err.Error())
		return erro
	}
	if resp.StatusCode != 200 {
		erro = fmt.Errorf("error al realizar peticion: " + resp.Status)
		return erro
	}
	defer resp.Body.Close()

	return
}

func buildUrlBanco(ruta string) (*url.URL, error) {
	base, err := url.Parse(ruta)
	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}
	return base, nil
}

func buildHeaderDefault(request *http.Request) {
	request.Header.Add("content-type", "application/json")
}
