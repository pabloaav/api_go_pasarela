package prismadtos

import (
	"errors"
	"regexp"

	tools "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type OfflineTokenRequest struct {
	Customer DataCustomer `json:"customer"`
}

func (tokenOffline *OfflineTokenRequest) ValidarSolicitudTokenOfflineRequest() error {
	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)
	if tokenOffline.Customer.Name == "" && tokenOffline.Customer.Identification.Type == "" && tokenOffline.Customer.Identification.Number == "" {
		return errors.New(ERROR_ESTRUCTURA_INCORRECTA)
	}

	if len(tokenOffline.Customer.Name) < 1 || tools.StringIsEmpity(tokenOffline.Customer.Name) {
		return errors.New(ERROR_HOLDER_NAME)
	}

	if len(tokenOffline.Customer.Identification.Number) < 7 || tools.StringIsEmpity(tokenOffline.Customer.Identification.Number) {
		return errors.New(ERROR_NRO_DOC)
	}
	if !digitCheckInt.MatchString(tokenOffline.Customer.Identification.Number) {
		return errors.New(ERROR_NRO_DOC)
	}

	err := tokenOffline.Customer.Identification.Type.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
