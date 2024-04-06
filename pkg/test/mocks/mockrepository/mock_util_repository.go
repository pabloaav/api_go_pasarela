package mockrepository

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryUtil struct {
	mock.Mock
}

func (mock *MockRepositoryUtil) CreateNotificacion(notificacion entities.Notificacione) (erro error) {
	args := mock.Called(notificacion)
	return args.Error(0)
}
func (mock *MockRepositoryUtil) CreateLog(log entities.Log) (erro error) {
	args := mock.Called(log)
	return args.Error(0)
}

func (mock *MockRepositoryUtil) GetConfiguracion(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Configuracione), args.Error(1)
}

func (mock *MockRepositoryUtil) CreateConfiguracion(config entities.Configuracione) (id uint, erro error) {
	args := mock.Called(config)
	return uint(args.Int(0)), args.Error(1)
}

func (r *MockRepositoryUtil) GetImpuestoByIdRepository(id int64) (impuesto entities.Impuesto, erro error) {
	args := r.Called(id)
	return args.Get(0).(entities.Impuesto), args.Error(1)
}

func (r *MockRepositoryUtil) CrearPeticionesRepository(peticionWeb entities.Webservicespeticione) (erro error) {
	args := r.Called()
	return args.Error(1)
}

func (r *MockRepositoryUtil) GetConfiguracionesRepository(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, erro error) {
	args := r.Called()
	return args.Get(0).([]entities.Configuracione), args.Error(1)
}

func (r *MockRepositoryUtil) CreateHistorialOperacionesRepository(historial entities.HistorialOperaciones) error {
	args := r.Called()
	return args.Error(1)
}
