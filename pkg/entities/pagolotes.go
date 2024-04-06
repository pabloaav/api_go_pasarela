package entities

import "gorm.io/gorm"

type Pagolotes struct {
	gorm.Model
	PagosID    uint64   `json:"pagos_id"`
	ClientesID uint64   `json:"clientes_id"`
	Lote       int64    `json:"lote"`
	FechaEnvio string   `json:"fecha_envio"`
	MotivoBaja string   `json:"motivi_baja"`
	Pago       *Pago    `json:"pago" gorm:"foreignKey:pagos_id"`
	Cliente    *Cliente `json:"cliente" gorm:"foreignKey:clientes_id"`
}
