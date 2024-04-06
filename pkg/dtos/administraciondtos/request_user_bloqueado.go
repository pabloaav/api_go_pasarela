package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestUserBloqueado struct {
	Id           uint      `json:"id"`
	Nombre       string    `json:"nombre"`
	Email        string    `json:"email"`
	DNI          string    `json:"dni"`
	Cuit         string    `json:"cuit"`
	FechaBloqueo time.Time `json:"fecha_bloqueo"`
	CantBloqueo  int       `json:"cant_bloqueo"`
	Permanente   bool      `json:"permanente"`
}

func (r *RequestUserBloqueado) ToEntity() entities.Usuariobloqueados {
	return entities.Usuariobloqueados{
		Nombre:       r.Nombre,
		Email:        r.Email,
		Dni:          r.DNI,
		Cuit:         r.Cuit,
		FechaBloqueo: r.FechaBloqueo,
		CantBloqueo:  r.CantBloqueo,
		Permanente:   r.Permanente,
	}
}
