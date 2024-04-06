package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseUsuariosBloqueados struct {
	Usuarios []UserBloqueado `json:"usuarios"`
	Meta     dtos.Meta       `json:"meta"`
}

type UserBloqueado struct {
	Id           uint      `json:"id"`
	Nombre       string    `json:"nombre"`
	Email        string    `json:"email"`
	DNI          string    `json:"dni"`
	Cuit         string    `json:"cuit"`
	FechaBloqueo time.Time `json:"fecha_bloqueo"`
	CantBloqueo  int       `json:"cant_bloqueo"`
	Permanente   bool      `json:"permanente"`
}

func (r *UserBloqueado) FromEntity(u entities.Usuariobloqueados) {
	r.Id = u.ID
	r.Nombre = u.Nombre
	r.Email = u.Email
	r.DNI = u.Dni
	r.FechaBloqueo = u.FechaBloqueo
	r.CantBloqueo = u.CantBloqueo
	r.Permanente = u.Permanente
}

func (ru *ResponseUsuariosBloqueados) FromEntities(users []entities.Usuariobloqueados) {
	for _, u := range users {
		ru.Usuarios = append(ru.Usuarios, UserBloqueado{
			Id:           u.ID,
			Nombre:       u.Nombre,
			Email:        u.Email,
			DNI:          u.Dni,
			Cuit:         u.Cuit,
			FechaBloqueo: u.FechaBloqueo,
			CantBloqueo:  u.CantBloqueo,
			Permanente:   u.Permanente,
		})
	}
	return
}
