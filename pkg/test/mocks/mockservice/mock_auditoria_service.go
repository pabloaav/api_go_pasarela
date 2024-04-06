package mockservice

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type MockAuditoriaService struct{}

func (m *MockAuditoriaService) Create(l *entities.Auditoria) error {
	return nil
}
