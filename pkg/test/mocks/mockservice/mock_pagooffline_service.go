package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/pagoofflinedtos"
	"github.com/stretchr/testify/mock"
)

type MockPagoOffLineService struct {
	mock.Mock
}

func (mock *MockPagoOffLineService) GenerarCodigoBarra(pagooffline_ov pagoofflinedtos.OffLineRapipagoDtos) (string, error) {
	args := mock.Called(pagooffline_ov)
	result := args.String(0)
	return result, args.Error(1)
}
