package mockrepository

import (
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/stretchr/testify/mock"
)

type MockRemoteRepositoryPrisma struct {
	mock.Mock
}

func (mock *MockRemoteRepositoryPrisma) GetHealthcheck() (response *prismatransacciones.HealthCheck, err error) {
	args := mock.Called()
	resultado := args.Get(0)
	return resultado.(*prismatransacciones.HealthCheck), args.Error(1)

}

func (mock *MockRemoteRepositoryPrisma) PostSolicitudTokenPago(card *prismatransacciones.Card) (response *prismatransacciones.PagoToken, err error) {
	args := mock.Called(*card)
	resultado := args.Get(0)
	return resultado.(*prismatransacciones.PagoToken), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) PostEjecutarPago(procesoPagoRequest *prismatransacciones.PaymentsSimpleRequest) (response *prismatransacciones.PaymentsSimpleResponse, err error) {
	args := mock.Called(*procesoPagoRequest)
	resultado := args.Get(0)

	return resultado.(*prismatransacciones.PaymentsSimpleResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) PostSolicitarTokenOffLine(tokenOffline *prismatransacciones.OfflineTokenRequest) (response *prismatransacciones.OfflineTokenResponse, err error) {
	args := mock.Called(*tokenOffline)
	resultado := args.Get(0)
	return resultado.(*prismatransacciones.OfflineTokenResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) PostEjecutarPagoOffLine(paymentOffline *prismatransacciones.PaymentsOfflineRequest) (response *prismatransacciones.PaymentsOfflineResponse, err error) {
	args := mock.Called(*paymentOffline)
	resultado := args.Get(0)
	return resultado.(*prismatransacciones.PaymentsOfflineResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) GetPrismaInformarPago(paymentId string) (response *prismainforme.UnPagoResponse, err error) {
	args := mock.Called(paymentId)
	resultado := args.Get(0)
	return resultado.(*prismainforme.UnPagoResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) ListarPagosPorFecha(request *prismainforme.ListaPagosRequest) (response *prismainforme.ListaPagosResponse, err error) {
	args := mock.Called(*request)
	resultado := args.Get(0)
	return resultado.(*prismainforme.ListaPagosResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) PostSolicitudAnulacionDevolucionPagoTotla(paymentId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, err error) {
	args := mock.Called(paymentId)
	resultado := args.Get(0)
	return resultado.(*prismaOperaciones.SolicitudAnulacionDevolucionResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, err error) {
	args := mock.Called(params)
	resultado := args.Get(0)
	return resultado.(*prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse), args.Error(1)
}

func (mock *MockRemoteRepositoryPrisma) DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, err error) {
	args := mock.Called(params)
	resultado := args.Get(0)
	return resultado.(*prismaOperaciones.DeletSolicitudPagoResponse), args.Error(1)
}
