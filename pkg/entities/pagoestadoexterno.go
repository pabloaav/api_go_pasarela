package entities

import "gorm.io/gorm"

type Pagoestadoexterno struct {
	gorm.Model
	Estado        string     `json:"estado"`
	Vendor        string     `json:"vendor"`
	PagoestadosId uint64     `json:"pagoestados_id"`
	PagoEstados   Pagoestado `json:"pago" gorm:"foreignKey:pagoestados_id"`
	// Apilinkcierrelote []Apilinkcierrelote `json:"apilinkcierrelote"`
}
