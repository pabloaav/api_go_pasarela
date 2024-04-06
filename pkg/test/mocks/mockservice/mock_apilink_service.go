package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/stretchr/testify/mock"
)

type MockApiLinkService struct {
	mock.Mock
}

func (mock *MockApiLinkService) CreateDebinApiLinkService(requerimientoId string, request linkdebin.RequestDebinCreateLink) (response linkdebin.ResponseDebinCreateLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linkdebin.ResponseDebinCreateLink), args.Error(1)
}

func (mock *MockApiLinkService) GetDebinesApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinesLink) (response linkdebin.ResponseGetDebinesLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinesLink), args.Error(1)
}

func (mock *MockApiLinkService) GetDebinesPendientesApiLinkService(requerimientoId string, cbu string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error) {
	args := mock.Called(requerimientoId, cbu)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinesPendientesLink), args.Error(1)
}

func (mock *MockApiLinkService) GetDebinApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinLink), args.Error(1)
}

func (mock *MockApiLinkService) DeleteDebinApiLinkService(requerimientoId string, request linkdebin.RequestDeleteDebinLink) (response bool, erro error) {
	args := mock.Called(requerimientoId, request)
	return args.Bool(0), args.Error(1)
}

func (mock *MockApiLinkService) CreateTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink) (response linktransferencia.ResponseTransferenciaCreateLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linktransferencia.ResponseTransferenciaCreateLink), args.Error(1)
}
func (mock *MockApiLinkService) GetTransferenciasApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink) (response linktransferencia.ResponseGetTransferenciasLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linktransferencia.ResponseGetTransferenciasLink), args.Error(1)
}

func (mock *MockApiLinkService) GetTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink) (response linktransferencia.ResponseGetTransferenciaLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linktransferencia.ResponseGetTransferenciaLink), args.Error(1)
}

func (mock *MockApiLinkService) GenerarUUid() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *MockApiLinkService) GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {
	args := mock.Called(requerimientoId, request)
	result := args.Get(0)
	return result.(linkconsultadestinatario.ResponseConsultaDestinatarioLink), args.Error(1)
}

func (mock *MockApiLinkService) CreateCuentaApiLinkService(request linkcuentas.LinkPostCuenta) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockApiLinkService) DeleteCuentaApiLinkService(request linkcuentas.LinkDeleteCuenta) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockApiLinkService) GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]linkcuentas.GetCuentasResponse), args.Error(1)
}

func (mock *MockApiLinkService) PutApilinkCierrelote(listaDebinesId []string) (erro error) {
	args := mock.Called()
	return args.Error(0)
}


