package mockservice

import (
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/stretchr/testify/mock"
)

type MockPrismaService struct {
	mock.Mock
	
}
func (mock *MockPrismaService) CheckService() (estado bool, err error) {
	args := mock.Called()
	result := args.Bool(0)
	return result, args.Error(1)
}
func (mock *MockPrismaService) SolicitarToken(request prismadtos.StructToken) (response interface{}, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result, args.Error(1)
}
func (mock *MockPrismaService) Payments(request prismadtos.StructPayments) (response interface{}, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result, args.Error(1)
}
func (mock *MockPrismaService) GetInformePago(pagoId string) (response *prismainforme.UnPagoResponse, erro error) {
	args := mock.Called(pagoId)
	result := args.Get(0)
	return result.(*prismainforme.UnPagoResponse), args.Error(1)
}
func (mock *MockPrismaService) ListarPagosPorFecha(requestPago prismainforme.ListaPagosRequest) (listaPago []prismainforme.Result, erro error) {
	args := mock.Called(requestPago)
	result := args.Get(0)
	return result.([]prismainforme.Result), args.Error(1)
}

func (mock *MockPrismaService) ListarPagosService(estadoPago int, channel string) (pagoIntento []entities.Pagointento, erro error) {
	args := mock.Called(estadoPago, channel)
	resultado := args.Get(0)
	return resultado.([]entities.Pagointento), args.Error(1)
}

func (mock *MockPrismaService) PostAnulacionDevolucionTotalPago(ExternalId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, erro error) {
	args := mock.Called(ExternalId)
	result := args.Get(0)
	return result.(*prismaOperaciones.SolicitudAnulacionDevolucionResponse), args.Error(1)
}
func (mock *MockPrismaService) PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, erro error) {
	args := mock.Called(params)
	result := args.Get(0)
	return result.(*prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse), args.Error(1)
}
func (mock *MockPrismaService) DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, erro error) {
	args := mock.Called(params)
	result := args.Get(0)
	return result.(*prismaOperaciones.DeletSolicitudPagoResponse), args.Error(1)
}

// func (mock *MockPrismaService) ArchivoLoteExterno() (totalArchivos int, err error) {
// 	args := mock.Called()
// 	resultado := args.Int(0)
// 	return resultado, args.Error(1)
// }
// func (mock *MockPrismaService) LeerCierreLote() (listaArchivo []prismaCierreLote.PrismaLogArchivoResponse, err error) {
// 	args := mock.Called()
// 	result := args.Get(0)
// 	return result.([]prismaCierreLote.PrismaLogArchivoResponse), args.Error(1)
// }
