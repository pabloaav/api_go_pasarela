package entities

import "gorm.io/gorm"

type ClienteRetencion struct {
	gorm.Model
	ClienteId    uint          `json:"cliente_id" gorm:"index:index_unique"`
	RetencionId  uint          `json:"retencion_id" gorm:"index:index_unique"`
	Cliente      Cliente       `json:"cliente" gorm:"foreignKey:ClienteId"`
	Retencion    Retencion     `json:"retencion" gorm:"foreignKey:RetencionId"`
	Certificados []Certificado `json:"certificados" gorm:"foreignKey:cliente_retencions_id"`
}
