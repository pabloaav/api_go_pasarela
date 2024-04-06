package apilink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
)

func (r *remoteRepository) CreateTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink, token string) (response linktransferencia.ResponseTransferenciaCreateLink, erro error) {

	base, erro := buildUrlTransferenica("transferencias")

	if erro != nil {
		return
	}

	json_data, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	buildHeaderAutorizacion(req, requerimientoId, token)

	erro = executeRequest(r, req, ERROR_CREATE_TRANSFERENCIA, &response)

	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "CreateTransferenciaApiLink",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}
	return
}

func (r *remoteRepository) GetTransferenciasApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink, token string) (response linktransferencia.ResponseGetTransferenciasLink, erro error) {

	base, erro := buildUrlTransferenica("transferencias")

	if erro != nil {
		return
	}

	params := url.Values{}
	params.Add("cbu", request.Cbu)
	params.Add("tamanio", string(request.Tamanio))
	params.Add("pagina", fmt.Sprint(request.Pagina))

	base.RawQuery = params.Encode() + fmt.Sprintf("&fechaDesde=%s", request.FechaDesde.Format("2006-01-02T15:04:05.000Z")) +
		fmt.Sprintf("&fechaHasta=%s", request.FechaHasta.Format("2006-01-02T15:04:05.000Z"))

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	erro = executeRequest(r, req, ERROR_GET_TRANSFERENCIAS, &response)

	return
}

func (r *remoteRepository) GetTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink, token string) (response linktransferencia.ResponseGetTransferenciaLink, erro error) {

	base, erro := buildUrlTransferenica("transferencias")

	if erro != nil {
		return
	}

	base.Path += fmt.Sprintf("/%s", request.NumeroReferenciaBancaria)

	params := url.Values{}
	params.Add("cbu", request.Cbu)

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	erro = executeRequest(r, req, ERROR_GET_TRANSFERENCIA, &response)

	return
}

func buildUrlTransferenica(ruta string) (*url.URL, error) {

	base, err := url.Parse(config.APILINKTRANSFERENCIA)

	if err != nil {
		logs.Error(ERROR_URL + err.Error())
		return nil, err
	}

	base.Path += ruta

	return base, nil
}
