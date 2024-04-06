package administraciondtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestPagoTipo struct {
	Id                       uint
	CuentasId                uint
	Pagotipo                 string
	BackUrlSuccess           string
	BackUrlPending           string
	BackUrlRejected          string
	BackUrlNotificacionPagos string
	IncludedChannels         []int64
	IncludedInstallments     []string
}

func (r *RequestPagoTipo) ToPagoTipo(cargarId bool) (response entities.Pagotipo) {
	if cargarId {
		response.ID = r.Id
	}
	response.CuentasID = r.CuentasId
	response.Pagotipo = r.Pagotipo
	response.BackUrlSuccess = r.BackUrlSuccess
	response.BackUrlPending = r.BackUrlPending
	response.BackUrlRejected = r.BackUrlRejected
	response.BackUrlNotificacionPagos = r.BackUrlNotificacionPagos
	// response.IncludedChannels = r.IncludedChannels
	// response.IncludedInstallments = r.IncludedInstallments
	return
}

func (r *RequestPagoTipo) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}

	// if isUpdate && len(r.IncludedInstallments) < 1 {
	// 	return fmt.Errorf(tools.ERROR_INCLUDED_INSTALLMENTS)
	// }

	if r.CuentasId == 0 {
		return fmt.Errorf(tools.ERROR_CUENTA_ID)
	}

	if commons.StringIsEmpity(r.Pagotipo) {
		return fmt.Errorf(tools.ERROR_PAGO_TIPO)
	}

	if commons.StringIsEmpity(r.BackUrlSuccess) {
		return fmt.Errorf(tools.ERROR_BACK_URL)
	}

	if commons.StringIsEmpity(r.BackUrlPending) {
		return fmt.Errorf(tools.ERROR_BACK_URL)
	}

	if len(r.IncludedChannels) < 1 {
		return fmt.Errorf(tools.ERROR_INCLUDED_CHANNELS)
	}

	// if len(r.IncludedInstallments) < 1 {
	// 	return fmt.Errorf(tools.ERROR_INCLUDED_INSTALLMENTS)
	// }

	// if commons.StringIsEmpity(r.IncludedChannels) {
	// 	return fmt.Errorf(tools.ERROR_INCLUDED_CHANNELS)
	// }

	// if commons.StringIsEmpity(r.IncludedInstallments) {
	// 	return fmt.Errorf(tools.ERROR_INCLUDED_INSTALLMENTS)
	// }

	r.Pagotipo = strings.ToUpper(r.Pagotipo)

	return
}

type RequestPagoTipoChannels struct {
	Add    []int64
	Delete []int64
}

type RequestPagoTipoCuotas struct {
	Add    []string
	Delete []string
}
