package mockrepository

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (mock *MockFileRepository) ExisteArchivo(nombre string) bool {
	args := mock.Called(nombre)
	return args.Bool(0)
}
func (mock *MockFileRepository) CrearArchivo(ruta string) (response *os.File, erro error) {
	args := mock.Called(ruta)
	result := args.Get(0)
	return result.(*os.File), args.Error(1)
}

func (mock *MockFileRepository) EliminarArchivo(ruta string) (erro error){
	args := mock.Called(ruta)
	return args.Error(0)
}
func (mock *MockFileRepository) AbrirArchivo(ruta string) (file *os.File, erro error){
	args := mock.Called(ruta)
	result := args.Get(0)
	return result.(*os.File), args.Error(1)
}
