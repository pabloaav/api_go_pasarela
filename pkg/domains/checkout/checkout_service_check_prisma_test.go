package checkout_test

import (
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

func TestCheckPrismaService(t *testing.T) {
	assertions := assert.New(t)
	serv := mockservice.MockPrismaService{}
	offlineService := mockservice.MockPagoOffLineService{}
	mockUtilService := new(mockservice.MockUtilService)
	//checkout.PrismaServiceVar = &serv
	mockStoreService := new(mockservice.MockStoreService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)

	checker := checkout.NewService(&mockrepository.MockRepository{}, &mockservice.MockCommonsService{}, &serv, &offlineService, mockUtilService, mockRepositoryWebHook, mockStoreService)

	t.Run("Cuando el llamado a la api devuelve un error devuelvo el error tal cual", func(t *testing.T) {
		serv.On("CheckService").Return(false, fmt.Errorf("Error en el request a la api externa")).Once()

		err := checker.CheckPrisma()

		assertions.EqualError(err, "Error en el request a la api externa")
	})

	t.Run("Si la api devuelve false en el check devemos devolver un error", func(t *testing.T) {
		serv.On("CheckService").Return(false, nil).Once()

		err := checker.CheckPrisma()

		assertions.EqualError(err, "el servicio de prisma no está disponible")
	})

	t.Run("Si la consulta a la api devuelve true devolvemos error en nil", func(t *testing.T) {
		serv.On("CheckService").Return(true, nil).Once()

		err := checker.CheckPrisma()

		assertions.NoError(err)
	})
}
