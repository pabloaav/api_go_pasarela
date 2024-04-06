package mockrepository

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/stretchr/testify/mock"
)

type MockRemoteRepositoryApiLink struct {
	mock.Mock
}

func (mock *MockRemoteRepositoryApiLink) GetTokenApiLink(identificador string, scope []linkdtos.EnumScopeLink) (linkdtos.TokenLink, error) {
	args := mock.Called(identificador, scope)
	result := args.Get(0)
	return result.(linkdtos.TokenLink), args.Error(1)

}

func (mock *MockRemoteRepositoryApiLink) CreateDebinApiLink(requerimientoId string, request linkdebin.RequestDebinCreateLink, token string) (response linkdebin.ResponseDebinCreateLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linkdebin.ResponseDebinCreateLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) GetDebinesApiLink(requerimientoId string, request linkdebin.RequestGetDebinesLink, token string) (response linkdebin.ResponseGetDebinesLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinesLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) GetDebinApiLink(requerimientoId string, request linkdebin.RequestGetDebinLink, token string) (response linkdebin.ResponseGetDebinLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) GetDebinesPendientesApiLink(requerimientoId string, cbu string, token string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error) {
	args := mock.Called(requerimientoId, cbu, token)
	result := args.Get(0)
	return result.(linkdebin.ResponseGetDebinesPendientesLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) DeleteDebinApiLink(requerimientoId string, request linkdebin.RequestDeleteDebinLink, token string) (response bool, erro error) {
	args := mock.Called(requerimientoId, request, token)
	return args.Bool(0), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) CreateTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink, token string) (response linktransferencia.ResponseTransferenciaCreateLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linktransferencia.ResponseTransferenciaCreateLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) GetTransferenciasApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink, token string) (response linktransferencia.ResponseGetTransferenciasLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linktransferencia.ResponseGetTransferenciasLink), args.Error(1)
}
func (mock *MockRemoteRepositoryApiLink) GetTransferenciaApiLink(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink, token string) (response linktransferencia.ResponseGetTransferenciaLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linktransferencia.ResponseGetTransferenciaLink), args.Error(1)
}

func (mock *MockRemoteRepositoryApiLink) GetConsultaDestinatario(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink, token string) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {
	args := mock.Called(requerimientoId, request, token)
	result := args.Get(0)
	return result.(linkconsultadestinatario.ResponseConsultaDestinatarioLink), args.Error(1)
}

func (mock *MockRemoteRepositoryApiLink) CreateCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockRemoteRepositoryApiLink) DeleteCuentaApiLink(request linkcuentas.LinkCuentasRequest) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockRemoteRepositoryApiLink) GetCuentasApiLink(request linkcuentas.LinkGetCuentasRequest) (response []linkcuentas.GetCuentasResponse, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]linkcuentas.GetCuentasResponse), args.Error(1)
}
