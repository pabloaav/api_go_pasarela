package entities

import (
	"time"

	"gorm.io/gorm"
)

type Installment struct {
	gorm.Model

	Descripcion             string              `json:"descripcion"`
	Issuer                  string              `json:"issuer"`
	VigenciaDesde           time.Time           `json:"vigencia_desde"`
	VigenciaHasta           *time.Time          `json:"vigencia_hasta"`
	MediopagoinstallmentsID int64               `json:"mediopagoinstallments_id"`
	Installmentdetail       []Installmentdetail `json:"installmentdetail" gorm:"foreignKey:InstallmentsID"`
}
