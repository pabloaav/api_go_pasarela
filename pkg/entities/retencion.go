package entities

import "gorm.io/gorm"

type Retencion struct {
	gorm.Model
	CondicionsId      uint               `json:"condicions_id"`
	Condicion         Condicion          `gorm:"foreignKey:CondicionsId"`
	ChannelsId        uint               `json:"channels_id"`
	Channel           Channel            `gorm:"foreignKey:ChannelsId"`
	Alicuota          float64            `json:"alicuota"`
	AlicuotaOpcional  float64            `json:"alicuota_opcional"`
	Rg2854            bool               `json:"rg2854"`
	Minorista         bool               `json:"minorista"`
	MontoMinimo       float64            `json:"monto_minimo"`
	Descripcion       string             `json:"descripcion"`
	CodigoRegimen     string             `json:"codigo"`
	Clientes          []Cliente          `gorm:"many2many:cliente_retencions;"`
	ClienteRetencions []ClienteRetencion `json:"cliente_retencions" gorm:"foreignKey:retencion_id"`
}
