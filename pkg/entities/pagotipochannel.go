package entities

import "gorm.io/gorm"

type Pagotipochannel struct {
	gorm.Model
	PagotiposId uint     `json:"pagotipos_id"`
	ChannelsId  uint     `json:"channels_id"`
	Channel     Channel  `json:"channel" gorm:"foreignKey:channels_id"`
	Pagotipo    Pagotipo `json:"pagotipo" gorm:"foreignKey:pagotipos_id"`
}
