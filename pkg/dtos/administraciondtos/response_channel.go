package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseChannels struct {
	Channels []ResponseChannel `json:"data"`
	Meta     dtos.Meta         `json:"meta"`
}

type ResponseChannel struct {
	Id         uint
	Channel    string
	Nombre     string
	CodigoBcra int32
	MedioPago  []ResponseMediopagoChannel
}

type ResponseMediopagoChannel struct {
	Id          uint
	Mediopago   string
	LongitudPan int32
	LongitudCvv int32
	Regexp      string
}

func (r *ResponseChannel) FromChannel(c entities.Channel) {
	r.Id = c.ID
	r.Channel = c.Channel
	r.Nombre = c.Nombre
	r.CodigoBcra = c.CodigoBcra
	if c.Mediopagos != nil {
		r.MedioPago = r.fromMedioPago(c.Mediopagos)
	}
}

func (r *ResponseChannel) fromMedioPago(mediosPago []entities.Mediopago) (channelMedioPago []ResponseMediopagoChannel) {
	for _, valueMedioPago := range mediosPago {
		medioPagoResponse := ResponseMediopagoChannel{
			Id:          valueMedioPago.ID,
			Mediopago:   valueMedioPago.Mediopago,
			LongitudPan: valueMedioPago.LongitudPan,
			LongitudCvv: valueMedioPago.LongitudCvv,
			Regexp:      valueMedioPago.Regexp,
		}
		channelMedioPago = append(channelMedioPago, medioPagoResponse)
	}
	return
}
