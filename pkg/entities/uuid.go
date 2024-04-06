package entities

import (
	"time"

	"gorm.io/gorm"
)

type Uuid struct {
	gorm.Model
	Uuid          string    `json:"uuid"`
	Permanente    bool      `json:"permanente"`
	Fecha_bloqueo time.Time `json:"fecha_bloqueo"`
}
