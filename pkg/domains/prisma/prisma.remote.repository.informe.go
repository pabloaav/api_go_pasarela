package prisma

import (
	//"errors"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	config "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	commonds "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////información de un pago y de varios pagos///////////////////
func (r *remoteRepository) ListarPagosPorFecha(request *prismainforme.ListaPagosRequest) (response *prismainforme.ListaPagosResponse, err error) {
	var erro ErrorEstructura
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	//base.Path += config.URI_PAYMENTS
	params := url.Values{}
	params.Add("offset", fmt.Sprint(request.Offset))
	params.Add("pageSize", fmt.Sprint(request.PageSize))
	params.Add("siteId", "99999966")
	if !commonds.StringIsEmpity(request.SiteOperationId) {
		params.Add("siteOperationId", request.SiteOperationId)
	}
	if !commonds.StringIsEmpity(request.DateFrom) {
		params.Add("dateFrom", request.DateFrom)
	}
	if !commonds.StringIsEmpity(request.DateTo) {
		params.Add("dateTo", request.DateTo)
	}
	if !commonds.StringIsEmpity(request.MerchantId) {
		params.Add("merchantId", request.MerchantId)
	}
	println(params.Encode())
	base.Path += config.URI_PAYMENTS

	base.RawQuery = params.Encode()
	fmt.Println(base.String())
	req, _ := http.NewRequest("GET", base.String(), nil)
	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al obtener lista de pago: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return nil, &erro
	}
	//logs.Info(resp.Body)
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

func (r *remoteRepository) GetPrismaInformarPago(paymentId string) (response *prismainforme.UnPagoResponse, err error) {
	var erro ErrorEstructura
	payload := strings.NewReader("{}")
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path = base.Path + config.URI_PAYMENTS + "/" + paymentId
	fmt.Println(base.String())
	req, _ := http.NewRequest("GET", base.String(), payload)
	buildHeaderDefault(req, config.PRIVATE_APIKEY_PRISMA)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al obtener datos de pago: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return nil, &erro
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

////////////////////////////////////////////funciones/////////////////////////////////////////////////////
// func CrearLog()  {

// }
