package mockservice

import (
	"io/fs"
	"mime/multipart"
	"os"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/commonsdtos"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type MockCommonsService struct {
	mock.Mock
}

func (mock *MockCommonsService) NewUUID() string {

	args := mock.Called()
	return args.String(0)

}

func (mock *MockCommonsService) IsValidUUID(u string) (bool, error) {
	args := mock.Called(u)
	return args.Bool(0), args.Error(1)
}

func (mock *MockCommonsService) LeerDirectorio(rutaFTP string) ([]fs.FileInfo, error) {
	args := mock.Called(rutaFTP)
	resultado := args.Get(0)
	return resultado.([]fs.FileInfo), args.Error(1)
}

func (mock *MockCommonsService) MoverArchivos(rutaOrigen, rutaDestino, nombreArchivo string) error {
	args := mock.Called(rutaOrigen, rutaDestino, nombreArchivo)
	return args.Error(1)
}

func (mock *MockCommonsService) BorrarArchivo(rutaFTP, nombreArchivo string) error {
	args := mock.Called(rutaFTP, nombreArchivo)
	return args.Error(1)
}

func (mock *MockCommonsService) BorrarDirectorio(ruta string) error {
	args := mock.Called(ruta)
	return args.Error(1)
}

func (mock *MockCommonsService) CreateFile(ruta string) (archivo *os.File, erro error) {
	args := mock.Called(ruta)
	resultado := args.Get(0)
	return resultado.(*os.File), args.Error(1)
}

func (mock *MockCommonsService) CreateFileName(file commonsdtos.FileName) string {
	args := mock.Called(file)
	return args.String(0)
}

func (mock *MockCommonsService) RemoveFile(ruta string) (erro error) {
	args := mock.Called(ruta)
	return args.Error(0)
}

func (mock *MockCommonsService) NormalizeStrings(str string) (string, error) {
	args := mock.Called(str)
	return args.String(0), args.Error(1)
}

func (mock *MockCommonsService) ZipFiles(request commonsdtos.ZipFilesRequest) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockCommonsService) SaveFiberPdf(file *multipart.FileHeader, ruta string, fiber *fiber.Ctx) (erro error) {
	args := mock.Called(file, ruta, fiber)
	return args.Error(0)
}

func (mock *MockCommonsService) CreateMessage(to []string, from, value string, Subject string) string {
	args := mock.Called(to, from, value, Subject)
	return args.String(0)
}

func (mock *MockCommonsService) FormatFecha() (fechaI time.Time, fechaF time.Time, erro error) {
	args := mock.Called()
	return time.Now(), time.Now(), args.Error(0)
}

func (mock *MockCommonsService) CrearArchivo(patch string) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockCommonsService) LeerArchivo(patch string) (archivo *os.File, erro error) {
	args := mock.Called()
	return nil, args.Error(0)
}

func (mock *MockCommonsService) EscribirArchivo(datos string, file *os.File) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockCommonsService) GuardarCambios(file *os.File) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockCommonsService) ConvertirFormatoFecha(fecha string) string {
	args := mock.Called()
	return args.String()
}

func (mock *MockCommonsService) ConvertirFecha(fecha string) string {
	args := mock.Called()
	return args.String()
}
func (mock *MockCommonsService) ConvertirFechaYYYYMMDD(fecha string) string {
	args := mock.Called()
	return args.String()
}

func (mock *MockCommonsService) RemoveAccents(valor string) (string, error) {
	args := mock.Called(valor)
	return args.String(0), args.Error(1)
}

func (mock *MockCommonsService) GetDateLastMoment(fecha time.Time) (fechaISO string) {
	args := mock.Called()
	return args.String()
}

func (mock *MockCommonsService) GetDateFirstMoment(fecha time.Time) (fechaISO string) {
	args := mock.Called()
	return args.String()
}

func (mock *MockCommonsService) ConvertirFechaToDDMMYYYY(fecha string) string {
	args := mock.Called()
	return args.String()
}
