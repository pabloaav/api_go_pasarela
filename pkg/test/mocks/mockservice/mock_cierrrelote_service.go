package mockservice

import (
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/stretchr/testify/mock"
)

type MockCierreLoteService struct {
	mock.Mock
}

func (mock *MockPrismaService) ArchivoLoteExterno() (totalArchivos int, err error) {
	args := mock.Called()
	resultado := args.Int(0)
	return resultado, args.Error(1)
}
func (mock *MockPrismaService) LeerCierreLote() (listaArchivo []prismaCierreLote.PrismaLogArchivoResponse, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]prismaCierreLote.PrismaLogArchivoResponse), args.Error(1)
}
