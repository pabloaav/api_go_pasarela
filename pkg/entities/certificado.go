package entities

import (
	"time"

	"gorm.io/gorm"
)

type Certificado struct {
	gorm.Model
	ClienteRetencionsId uint             `json:"cliente_retencions_id"`
	Fecha_Presentacion  time.Time        `json:"fecha_presentacion,omitempty"`
	Fecha_Caducidad     time.Time        `json:"fecha_caducidad,omitempty"`
	Ruta_file           string           `json:"ruta_file"`
	ClienteRetencion    ClienteRetencion `json:"cliente_retencion" gorm:"foreignKey:cliente_retencions_id"`
}

func (c Certificado) IsExpired() (res bool) {
	if !c.Fecha_Caducidad.IsZero() {
		hoy := time.Now().Local()
		res = hoy.After(c.Fecha_Caducidad)
	}
	return
}
