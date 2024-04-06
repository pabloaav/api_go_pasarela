package auditoria

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type AuditoriaService interface {
	Create(l *entities.Auditoria) error
}

func New(r AuditoriaRepository) AuditoriaService {
	service := auditoriaService{
		repository: r,
	}

	return &service
}

type auditoriaService struct {
	repository AuditoriaRepository
}

func (s *auditoriaService) Create(l *entities.Auditoria) error {
	return s.repository.Create(l)
}
