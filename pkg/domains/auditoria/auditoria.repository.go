package auditoria

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type AuditoriaRepository interface {
	Create(audit *entities.Auditoria) error
}

func NewAuditoriaRepository(conn *database.MySQLClient) AuditoriaRepository {
	return &auditoriaRepository{
		SqlClient: conn,
	}
}

type auditoriaRepository struct {
	SqlClient *database.MySQLClient
}

func (r *auditoriaRepository) Create(audit *entities.Auditoria) error {
	if len(audit.Resultado) > 254 {
		audit.Resultado = audit.Resultado[0:254]
	}
	res := r.SqlClient.Create(audit)
	if res.RowsAffected <= 0 {
		return res.Error
	}
	return nil
}
