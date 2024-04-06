package prismadtos

import (
	"errors"

	tools "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type ParamsPagoParcialTotalService struct {
	ExternalId string `json:"external_id,omitempty"`
	RefundId   string `json:"refund_id,omitempty"`
	Monto      int64  `json:"monto,omitempty"`
}

func (pagoParcialTotal *ParamsPagoParcialTotalService) ValidarParametros(opcion string) error {
	switch opcion {
	// case "SAPT":

	// case "DAPT":
	//case  "DAPP":

	case "SAPP":
		if tools.StringIsEmpity(pagoParcialTotal.ExternalId) {
			return errors.New(ERROR_ID_PAYMENT)
		}
		if pagoParcialTotal.Monto <= 0 {
			return errors.New(ERROR_MONTO)
		}
	case "ER":
		if tools.StringIsEmpity(pagoParcialTotal.ExternalId) {
			return errors.New(ERROR_ID_PAYMENT)
		}
		if tools.StringIsEmpity(pagoParcialTotal.RefundId) {
			return errors.New(ERROR_ID_REFUND)
		}
	default:
		return errors.New(ERROR_CASE)
	}
	return nil
}
