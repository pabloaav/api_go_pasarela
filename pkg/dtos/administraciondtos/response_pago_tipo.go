package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagosTipo struct {
	PagosTipo []ResponsePagoTipo `json:"data"`
	Meta      dtos.Meta          `json:"meta"`
}

type ResponsePagoTipo struct {
	Id                       uint
	CuentasId                uint
	Cuenta                   string
	Pagotipo                 string
	BackUrlSuccess           string
	BackUrlPending           string
	BackUrlRejected          string
	BackUrlNotificacionPagos string
	IncludedChannels         []CanalesPago
	IncludedInstallments     []CuotasPago
}

type CanalesPago struct {
	ChannelsId uint
	Channel    string
	Nombre     string
}
type CuotasPago struct {
	Nro string
}

func (r *ResponsePagoTipo) FromPagoTipo(p entities.Pagotipo) {
	r.Id = p.ID
	r.CuentasId = p.CuentasID
	r.Cuenta = p.Cuenta.Cuenta
	r.Pagotipo = p.Pagotipo
	r.BackUrlSuccess = p.BackUrlSuccess
	r.BackUrlPending = p.BackUrlPending
	r.BackUrlRejected = p.BackUrlRejected
	r.BackUrlNotificacionPagos = p.BackUrlNotificacionPagos
	r.BackUrlNotificacionPagos = p.BackUrlNotificacionPagos
}
