package mockrepository

import (
	"context"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx()    {}
func (m *MockRepository) CommitTx()   {}
func (m *MockRepository) RollbackTx() {}

func (m *MockRepository) CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error) {
	args := m.Called(ctx, pago)
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}
func (m *MockRepository) UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error) {
	args := m.Called(pago)
	result := args.Bool(0)
	return result, args.Error(1)
}
func (m *MockRepository) GetPagoByUuid(uuid string) (*entities.Pago, error) {
	args := m.Called(uuid)
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}
func (m *MockRepository) GetClienteByApikey(apikey string) (*entities.Cliente, error) {
	args := m.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cliente), args.Error(1)
}
func (m *MockRepository) GetCuentaByApikey(apikey string) (*entities.Cuenta, error) {
	args := m.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}
func (m *MockRepository) GetPagotipoById(id int64) (*entities.Pagotipo, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Pagotipo), args.Error(1)
}
func (m *MockRepository) GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*[]entities.Pagotipochannel), args.Error(1)
}
func (m *MockRepository) GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*[]entities.Pagotipointallment), args.Error(1)
}
func (m *MockRepository) GetChannelByName(nombre string) (*entities.Channel, error) {
	args := m.Called(nombre)
	result := args.Get(0)
	return result.(*entities.Channel), args.Error(1)
}
func (m *MockRepository) GetCuentaById(id int64) (*entities.Cuenta, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}
func (m *MockRepository) CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error) {
	args := m.Called(ctx, resultado)
	result := args.Bool(0)
	return result, args.Error(1)
}
func (m *MockRepository) GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error) {
	args := m.Called(pagoId)
	result := args.Get(0)
	return result.(*entities.Pagointento), args.Error(1)
}
func (m *MockRepository) GetMediosDePagos() (*[]entities.Mediopago, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*[]entities.Mediopago), args.Error(1)
}
func (m *MockRepository) GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.(*entities.Mediopago), args.Error(1)
}
func (m *MockRepository) GetInstallmentDetailsID(installmentID, numeroCuota int64) int64 {
	args := m.Called(installmentID, numeroCuota)
	result := args.Get(0)
	return result.(int64)
}
func (m *MockRepository) GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error) {
	args := m.Called(installmentID, numeroCuota)
	result := args.Get(0)
	return result.(*dtos.InstallmentDetailsResponse), args.Error(1)

}
func (m *MockRepository) CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error {
	args := m.Called(ctx, pel)
	return args.Error(0)
}

func (m *MockRepository) GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.([]entities.Installment), args.Error(1)
}

func (m *MockRepository) GetPagoEstado(id int64) (*entities.Pagoestado, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Pagoestado), args.Error(1)
}
func (m *MockRepository) GetPreferencesByIdClienteRepository(id uint) (preferencia entities.Preference, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(entities.Preference), args.Error(1)

}

func (m *MockRepository) GetChannelById(id uint) (channel entities.Channel, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(entities.Channel), args.Error(1)

}

func (m *MockRepository) CheckUsuarioTRepository(usuario entities.Usuariobloqueados) (usuarioDB entities.Usuariobloqueados, err error) {
	args := m.Called(usuario)
	result := args.Get(0)
	return result.(entities.Usuariobloqueados), args.Error(1)
}
func (m *MockRepository) AgregarUsuarioListaBloqueoRepository(entities.Usuariobloqueados) error {
	return nil
}
func (m *MockRepository) DeleteUsuarioListaBloqueoRepository(entitie entities.Usuariobloqueados) error {
	return nil
}
func (m *MockRepository) UpdatePagoIntentoFallidos(ctx context.Context, pago *entities.Pago) error {
	return nil

}
func (m *MockRepository) UpdateUsuarioBloqueoRepository(usuario entities.Usuariobloqueados) error {
	return nil
}
func (m *MockRepository) SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error) {
	args := m.Called(hasheado, pagointento_id)
	return args.Error(1)
}
func (m *MockRepository) GetHasheado(hash string) (control bool, erro error) {
	args := m.Called(hash)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}
