package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	uuid "github.com/satori/go.uuid"
)

type AplinkService interface {
	CreateDebinApiLinkService(requerimientoId string, request linkdebin.RequestDebinCreateLink) (response linkdebin.ResponseDebinCreateLink, erro error)
	GetDebinesApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinesLink) (response linkdebin.ResponseGetDebinesLink, erro error)
	GetDebinesPendientesApiLinkService(requerimientoId string, cbu string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error)
	GetDebinApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, erro error)
	DeleteDebinApiLinkService(requerimientoId string, request linkdebin.RequestDeleteDebinLink) (response bool, erro error)

	//Cuentas
	CreateCuentaApiLinkService(request linkcuentas.LinkPostCuenta) (erro error)
	DeleteCuentaApiLinkService(request linkcuentas.LinkDeleteCuenta) (erro error)
	GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error)

	/*
		Este servicio fue creado para realizar transferencias de la cuenta de telco para sus clientes.
	*/
	CreateTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink) (response linktransferencia.ResponseTransferenciaCreateLink, erro error)
	GetTransferenciasApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink) (response linktransferencia.ResponseGetTransferenciasLink, erro error)
	GetTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink) (response linktransferencia.ResponseGetTransferenciaLink, erro error)
	/*
		Genera un UUid
	*/
	GenerarUUid() string

	GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error)

	/* Actualizar el campo estado de la tabla de apilinkcierrelote */
	/*
		Autor: Jose Luis Alarcon
		Fecha: 19/05/2022
		Descripción:  Este servicio fue creado para actualizar el campo estado de la tabla de apilinkcierrelote ,
		una vez que se coinciliaron  con los movimientos del banco
	*/
	PutApilinkCierrelote(listaDebinesId []string) (erro error)
}

// apilink variable que va a manejar la instancia del servicio
var apilink *aplinkService

type aplinkService struct {
	remoteRepository RemoteRepository
	repository       ApilinkRepository
}

func NewService(rm RemoteRepository, r ApilinkRepository) AplinkService {
	// al instanciar el servicio lo almaceno en la variable apilink
	apilink = &aplinkService{
		remoteRepository: rm,
		repository:       r,
	}
	return apilink
}

// Resolve devuelve la instancia antes creada
func Resolve() *aplinkService {
	return apilink
}

func (s *aplinkService) GenerarUUid() string {
	return uuid.NewV4().String()
}
