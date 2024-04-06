package entities

import "gorm.io/gorm"

type Pagotipointallment struct {
	gorm.Model
	PagotiposId uint   `json:"pagotipos_id"`
	Cuota       string `json:"cuota"`
}
