package banco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type RemoteRepository interface {
	GetTokenBancoRemotRepository(token bancodtos.RequestTokenBanco) (responseToken bancodtos.ResponseTokenBanco, err error)
	GetMovimientosBancoRemotRepository(filtrosMovimientos filtros.MovimientosBancoFiltro, token string) (response []bancodtos.ResponseMovimientosBanco, erro error)
	ActualizarRegistrosMatchRepository(listaMovimientoMatch bancodtos.RequestUpdateMovimiento, token string) (response bool, erro error)
	GetConsultarMovimientosRemoteRepository(filtro filtros.RequestMovimientos, token string) (response bancodtos.ResponseMovimientos, erro error)
	GetProcesarMovimientosRemoteRepository(token string) (response bancodtos.ResponseMovimientosProceso, erro error)
}

type remoteRepository struct {
	HTTPClient *http.Client
}

func NewRemote(http *http.Client) RemoteRepository {
	return &remoteRepository{
		HTTPClient: http,
	}
}

func (r *remoteRepository) GetTokenBancoRemotRepository(token bancodtos.RequestTokenBanco) (responseToken bancodtos.ResponseTokenBanco, err error) {

	base, err := url.Parse(config.AUTH)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return responseToken, err
	}

	base.Path += "/users/login"

	json_data, _ := json.Marshal(token)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderDefault(req)

	resp, err := r.HTTPClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("")
		return
	}

	json.NewDecoder(resp.Body).Decode(&responseToken)
	return responseToken, nil
}

func (r *remoteRepository) GetMovimientosBancoRemotRepository(filtrosMovimientos filtros.MovimientosBancoFiltro, token string) (response []bancodtos.ResponseMovimientosBanco, erro error) {

	base, err := buildUrlBanco("/movimientos")
	if err != nil {
		return
	}
	json_data, _ := json.Marshal(filtrosMovimientos)

	req, _ := http.NewRequest("GET", base.String(), bytes.NewBuffer(json_data))

	buildHeaderDefault(req)
	buildHeaderAutorizacion(req, token)

	resp, erro := r.HTTPClient.Do(req)

	if erro != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		erro = fmt.Errorf("error al obtener los movimientos api banco")
		return
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

func (r *remoteRepository) ActualizarRegistrosMatchRepository(listaMovimientoMatch bancodtos.RequestUpdateMovimiento, token string) (response bool, erro error) {
	var respuesta bancodtos.ActualizacionResponse
	requestJson, _ := json.Marshal(listaMovimientoMatch)

	urlBase, err := buildUrlBanco("/movimientos")
	if err != nil {
		erro = fmt.Errorf(err.Error())
		return
	}
	req, _ := http.NewRequest("POST", urlBase.String(), bytes.NewBuffer(requestJson))
	buildHeaderDefault(req)
	buildHeaderAutorizacion(req, token)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al realizar peticion: " + err.Error())
	}
	if resp.StatusCode != 200 {
		erro = fmt.Errorf("error al actualizar los movimientos api banco")
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&respuesta)
	response = respuesta.Status
	return
}

func (r *remoteRepository) GetConsultarMovimientosRemoteRepository(filtro filtros.RequestMovimientos, token string) (response bancodtos.ResponseMovimientos, erro error) {
	params := url.Values{}
	params.Add("Number", fmt.Sprint(filtro.Paginacion.Number))
	params.Add("Size", fmt.Sprint(filtro.Paginacion.Size))
	params.Add("FechaInicio", fmt.Sprint(filtro.FechaInicio))
	params.Add("FechaFin", fmt.Sprint(filtro.FechaFin))
	if len(filtro.ListaIdsTipoMovimientos) > 0 {
		params.Add("TipoMovimientosIds", fmt.Sprint(filtro.ListaIdsTipoMovimientos))
	}
	if len(filtro.DebitoCredito) > 0 {
		params.Add("DebitoCredito", fmt.Sprint(filtro.DebitoCredito))
	}

	urlBase, err := buildUrlBanco("/movimientos-banco-cuenta")
	if err != nil {
		logs.Error(err.Error())
	}
	urlBase.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", urlBase.String(), nil)
	buildHeaderDefault(req)
	buildHeaderAutorizacion(req, token)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		erro = fmt.Errorf(err.Error())
		return
	}
	defer resp.Body.Close()
	if !strings.Contains(resp.Status, "200") {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return
}

func (r *remoteRepository) GetProcesarMovimientosRemoteRepository(token string) (response bancodtos.ResponseMovimientosProceso, erro error) {

	urlBase, err := buildUrlBanco("/movimientos-banco")
	if err != nil {
		logs.Error(err.Error())
	}
	req, _ := http.NewRequest("GET", urlBase.String(), nil)
	buildHeaderDefault(req)
	buildHeaderAutorizacion(req, token)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		erro = fmt.Errorf(err.Error())
		return
	}
	defer resp.Body.Close()
	if !strings.Contains(resp.Status, "200") {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return
}

func buildUrlBanco(ruta string) (*url.URL, error) {
	base, err := url.Parse(config.URL_BANCO)
	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}
	base.Path += ruta

	return base, nil
}

func buildHeaderAutorizacion(request *http.Request, token string) {
	request.Header.Add("authorization", "Bearer "+token)
}

func buildHeaderDefault(request *http.Request) {
	request.Header.Add("content-type", "application/json")
}
