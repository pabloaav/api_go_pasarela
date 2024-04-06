package mockservice

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type MockUtilService struct {
	mock.Mock
}

func (mock *MockUtilService) CreateNotificacionService(notificacion entities.Notificacione) (erro error) {
	args := mock.Called(notificacion)
	return args.Error(0)
}

func (mock *MockUtilService) CreateLogService(log entities.Log) (erro error) {
	args := mock.Called(log)
	return args.Error(0)
}

func (mock *MockUtilService) GetConfiguracionService(filtro filtros.ConfiguracionFiltro) (configuracion administraciondtos.ResponseConfiguracion, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(administraciondtos.ResponseConfiguracion), args.Error(1)
}

func (mock *MockUtilService) CreateConfiguracionService(config administraciondtos.RequestConfiguracion) (id uint, erro error) {
	args := mock.Called(config)
	return uint(args.Int(0)), args.Error(1)
}

func (mock *MockUtilService) FirstOrCreateConfiguracionService(nombre string, descripcion string, valor string) (key string, erro error) {
	args := mock.Called(nombre, descripcion, valor)
	return args.String(0), args.Error(1)
}

func (mock *MockUtilService) LogError(erro string, funcionalidad string) {

}

func (mock *MockUtilService) ToFixed(num float64, precision int) float64 {
	args := mock.Called(num, precision)
	return float64(args.Int(0))
}

func (mock *MockUtilService) FormatNum(num float64) string {
	p := message.NewPrinter(language.Spanish)
	valor := p.Sprintf("%.2f", num)
	logs.Info(valor)
	return valor
}

func (mock *MockUtilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {
	args := mock.Called(movimiento, cuentacomisiones, iva, importeSolicitado)
	return args.Error(0)
}

func (mock *MockUtilService) CrearPeticionesService(peticiones dtos.RequestWebServicePeticion) (erro error) {
	args := mock.Called(peticiones)
	return args.Error(0)

}

func (mock *MockUtilService) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (configuraciones []administraciondtos.ResponseConfiguracion, erro error) {
	args := mock.Called(filtro)
	return args.Get(0).([]administraciondtos.ResponseConfiguracion), args.Error(1)
}

func (mock *MockUtilService) GetImpuestoByIdService(id int64) (impuesto entities.Impuesto, erro error) {
	args := mock.Called(id)
	return args.Get(0).(entities.Impuesto), args.Error(1)
}

func (mock *MockUtilService) CalcularValorCuponService(importe, coeficiente, impuesto float64) (valorCupon int64) {
	args := mock.Called(importe, coeficiente, impuesto)
	return args.Get(0).(int64)
}

func (mock *MockUtilService) CalcularCostoFinancieroIvaService(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA float64) (importeCFMasIva int64) {
	args := mock.Called(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA)
	return args.Get(0).(int64)
}

func (mock *MockUtilService) RightStr(cadenaStr string, valueStr int) (rightStr string) {
	args := mock.Called(cadenaStr, valueStr)
	return args.String(0)
}

func (mock *MockUtilService) LeftStr(cadenaStr string, valueStr int) (LeftStr string) {
	args := mock.Called(cadenaStr, valueStr)
	return args.String(0)
}

func (mock *MockUtilService) BuildStr(cadenaStr string, valueStr int) (center string) {
	args := mock.Called(cadenaStr, valueStr)
	return args.String(0)
}

func (mock *MockUtilService) GetMatenimietoSistemaService() (estado bool, fecha time.Time, erro error) {
	args := mock.Called()
	return args.Bool(0), args.Get(1).(time.Time), args.Error(2)
}

func (mock *MockUtilService) EnviarMailService(params utildtos.RequestDatosMail) (erro error) {
	args := mock.Called(params)
	args.Get(0)
	return args.Error(1)
}

func (mock *MockUtilService) CsvCreate(name string, data [][]string) error {
	args := mock.Called(name)
	args.Get(0)
	return args.Error(1)
}

func (mock *MockUtilService) ValidarCBU(cbu string) (res bool, erro error) {
	args := mock.Called(cbu)
	args.Get(0)
	return args.Bool(0), args.Error(1)
}

func (mock *MockUtilService) ValidarCalculoCF(RequestValidarCF utildtos.RequestValidarCF) (responseValidarCF utildtos.ResponseValidarCF) {
	args := mock.Called(RequestValidarCF)
	args.Get(0)
	return args.Get(0).(utildtos.ResponseValidarCF)
}

func (mock *MockUtilService) BuildComisionesTemporales(movimiento *entities.Movimientotemporale, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {
	args := mock.Called(movimiento, cuentacomisiones, iva, importeSolicitado)
	return args.Error(0)
}

func (mock *MockUtilService) GetRecaudacionPdf(request reportedtos.ResultMovLiquidacion, ruta, nombrearchivo string) error {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockUtilService) CrearReferenciaBancoStandard(referencia_banco string) (string, error) {
	args := mock.Called(referencia_banco)
	return args.String(0), args.Error(1)
}

func (mock *MockUtilService) CreateHistorialOperacionesService(bearer, observacion string, tipo entities.EnumHistorial) {
	mock.Called(bearer, observacion, tipo)
	return
}

func (mock *MockUtilService) DecodeToken(bearer string) (int64, string) {
	args := mock.Called(bearer)
	return int64(args.Int(0)), args.String(1)
}

func (mock *MockUtilService) RunEndpoint(typeEndpoint string, urlEndpoint string, mapConHeaders map[string]string, bodyEndpoint interface{}, queryParameters map[string]string, registarPeticion bool) (objeto interface{}, erro error) {
	args := mock.Called(typeEndpoint, urlEndpoint, mapConHeaders, bodyEndpoint, queryParameters, registarPeticion)
	return (args.Get(0)), args.Error(1)
}
