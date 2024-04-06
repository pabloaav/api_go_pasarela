package prisma

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	config "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

type RepositoryRemotePrisma interface {
	////////////////////////////////////////transacciones simples y offline//////////////////////////////////////////////////////////
	/*
		GetHealthcheck chequea si el servico de prisma esta en funcionamiento
	*/
	GetHealthcheck() (response *prismatransacciones.HealthCheck, err error)

	/*
		PostSolicitudTokenPago Solicita toquen de pago
	*/
	PostSolicitudTokenPago(card *prismatransacciones.Card) (response *prismatransacciones.PagoToken, err error)

	/*
		PostEjecutarPago Ejecucion de Pago
	*/
	PostEjecutarPago(procesoPagoRequest *prismatransacciones.PaymentsSimpleRequest) (response *prismatransacciones.PaymentsSimpleResponse, err error)

	// /*
	// 	PostSolicitarTokenOffLine permite solicitar un token de pago, para luego ejecutar un pago offline
	// */
	// PostSolicitarTokenOffLine(tokenOffline *prismatransacciones.OfflineTokenRequest) (response *prismatransacciones.OfflineTokenResponse, err error)

	// /*
	// 	PostEjecutarPagoOffLine
	// */
	// PostEjecutarPagoOffLine(paymentOffline *prismatransacciones.PaymentsOfflineRequest) (response *prismatransacciones.PaymentsOfflineResponse, err error)

	/////////////////////////////////////información de un pago y de varios pagos///////////////////////////////////////////////////////
	/*
		GetPrismaInformarPago realiza una peticion sobre un pago en prisma
	*/
	GetPrismaInformarPago(paymentId string) (response *prismainforme.UnPagoResponse, err error)

	/*
		ListarPagosPorFecha realiza una peticion a prisma, para obtener una lista de pagos realizados en un perido de fecha
	*/
	ListarPagosPorFecha(request *prismainforme.ListaPagosRequest) (response *prismainforme.ListaPagosResponse, err error)

	/////////////////////////////////////OPERACIONES SOBRE TRANSACCIONES///////////////////////////////////////////////////////
	/*
		PostSolicitudAnulacionDevolucionPagoTotla se realiza la solicitud de anulación o devolución total de pago a prisma
	*/
	PostSolicitudAnulacionDevolucionPagoTotla(paymentId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, err error)

	/*
		PostSolicitudAnulacionDevolucionPagoParcial se realiza la solicitud de anulación o devolución parcial de pago a prisma
		se le debe enviar el payment_id y el monto
		DelAnulacionDevolucionPagoTotalParcial perimte eliminar una solicitud de anulación o devolución total de pago a prisma
		se le debe enviar el payment_id y refund_id
	*/
	PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, err error)

	/*
		PostSolicitudAnulacionDevolucionPagoParcial se realiza la solicitud de anulación o devolución parcial de pago a prisma
		se le debe enviar el payment_id y el monto
		DelAnulacionDevolucionPagoTotalParcial perimte eliminar una solicitud de anulación o devolución total de pago a prisma
		se le debe enviar el payment_id y refund_id
	*/
	DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, err error)
}

type remoteRepository struct {
	HTTPClient *http.Client
	repository Repository
}

//NewRepository constructor
func NewRepoasitory(http *http.Client, r Repository) RepositoryRemotePrisma {
	return &remoteRepository{
		HTTPClient: http,
		repository: r,
	}

}

///////////////////////////////////////////////////////////////////////////////////////////
////////////// inicio chequear servico
func (r *remoteRepository) GetHealthcheck() (response *prismatransacciones.HealthCheck, err error) {
	var erro ErrorEstructura
	payload := strings.NewReader("{}")
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_HEALTHCHECK
	logs.Info(base.String())
	req, _ := http.NewRequest("GET", base.String(), payload)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al verificar servico de prisma: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return response, &erro
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

//////////////fin chequear servicio

// buildHeaderDefault agrega el Header
func buildHeaderDefault(request *http.Request, key string) {
	request.Header.Add("apikey", key)
	request.Header.Add("content-type", "application/json; charset=utf-8")
	request.Header.Add("Cache-Control", "no-cache")
}
