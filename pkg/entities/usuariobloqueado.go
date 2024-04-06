package entities

import (
	"time"

	"gorm.io/gorm"
)

type Usuariobloqueados struct {
	gorm.Model
	Nombre       string    `json:"nombre"`
	Email        string    `json:"email"`
	Dni          string    `json:"dni"`
	Cuit         string    `json:"cuit"`
	FechaBloqueo time.Time `json:"fecha_bloqueo"`
	CantBloqueo  int       `json:"cant_bloqueo"`
	Permanente   bool      `json:"permanente"`
}
