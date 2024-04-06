package mockrepository

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryPrisma struct {
	mock.Mock
}

func (mock *MockRepositoryPrisma) SaveCierreLote(detalleLote *entities.Prismacierrelote) (bool, error) {
	args := mock.Called(detalleLote)
	resultado := args.Get(0)
	fmt.Println(resultado)
	return true, args.Error(1)
}

func (mock *MockRepositoryPrisma) SaveCierreLoteBatch(detalleLote []entities.Prismacierrelote) (bool, error) {
	args := mock.Called(detalleLote)
	resultado := args.Get(0)
	fmt.Println(resultado)
	return true, args.Error(1)
}

func (mock *MockRepositoryPrisma) GetPagosPagosIntentosxChannel(estadoPago int, channel string) (pagos []entities.Pagointento, erro error) {
	args := mock.Called(estadoPago, channel)
	resultado := args.Get(0)
	return resultado.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryPrisma) GetMensajeErrorPrismaByExternalId(external_id uint64) (msgErrorPrisma entities.Prismaerroresexterno, erro error) {
	args := mock.Called(external_id)
	resultado := args.Get(0)
	return resultado.(entities.Prismaerroresexterno), args.Error(1)
}
