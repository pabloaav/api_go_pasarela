package entities

import "gorm.io/gorm"

type Envio struct {
	gorm.Model
	Cobranzas   bool
	Rendiciones bool
	Reversiones bool
	Batch       bool
	BatchPagos  bool
	ClientesId  uint
}
