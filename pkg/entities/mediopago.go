package entities

import "gorm.io/gorm"

type Mediopago struct {
	gorm.Model
	ChannelsID              int64                `json:"channels_id"`
	AdquirientesID          int64                `json:"adquirientes_id"`
	InstallmentsID          int64                `json:"installments_id"`
	MediopagoinstallmentsID int64                `json:"mediopagoinstallments_id"`
	Mediopago               string               `json:"mediopago"`
	ExternalID              string               `json:"external_id"`
	LongitudPan             int32                `json:"longitud_pan"`
	LongitudCvv             int32                `json:"longitud_cvv"`
	Regexp                  string               `json:"regexp"`
	Installment             Installment          `json:"installmentdetail" gorm:"foreignKey:installments_id"`
	Channel                 Channel              `json:"channel" gorm:"foreignKey:ChannelsID"`
	Mediopagoinstallment    Mediopagoinstallment `json:"mediopagoinstallment" gorm:"foreignKey:MediopagoinstallmentsID"`
}
