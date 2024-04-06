package entities

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"gorm.io/gorm"
)

type Reversione struct {
	gorm.Model
	PagointentosID      uint        `json:"pagointentos_id"`
	ExternalreversionID string      `json:"externalreversion_id"`
	Amount              int64       `json:"amount"`
	Status              string      `json:"status"`
	PagoIntento         Pagointento `json:"pagointento" gorm:"foreignKey:pagointentos_id"`
}

func (r *Reversione) ValidarReversion() error {

	if r.PagointentosID == 0 {
		return errors.New(ERROR_PAGOINTENTOID)
	}
	if tools.EsStringVacio(r.ExternalreversionID) {
		return errors.New(ERROR_EXTERNAL_ID)
	}

	if r.Amount == 0 {
		return errors.New(ERROR_AMOUNT)
	}

	return nil

}

const (
	ERROR_PAGOINTENTOID = "pago intento no debe ser nulo "
	ERROR_EXTERNAL_ID   = "el campo no debe estar vacio"
	ERROR_AMOUNT        = "el monto no debe ser 0"
	ERROR_STATUS        = "el estado no debe ser vacio"
)

func (r *Reversione) AddReversion(pagoIntentoId uint, amount int64, externalReversionId, status string) error {
	r.PagointentosID, r.ExternalreversionID, r.Amount, r.Status = pagoIntentoId, externalReversionId, amount, status
	err := r.ValidarReversion()
	if err != nil {
		return err
	}
	return nil
}
