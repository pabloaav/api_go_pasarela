package entities

import "gorm.io/gorm"

type Pagoitems struct {
	gorm.Model
	PagosID     int64  `json:"pagos_id"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
	Amount      Monto  `json:"amount"`
	Identifier  string `json:"identifier"`
	Pago        Pago   `json:"pago" gorm:"foreignKey:pagos_id"`
}
