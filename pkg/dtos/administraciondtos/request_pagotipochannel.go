package administraciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestPagoTipoChannel struct {
	Id         uint
	ChannelId  uint
	PagoTipoId uint
}

func (r *RequestPagoTipoChannel) ToPagoTipoChannel(cargarId bool) (response entities.Pagotipochannel) {
	if cargarId {
		response.ID = r.Id
	}
	response.ChannelsId = r.ChannelId
	response.PagotiposId = r.PagoTipoId

	return
}

func (r *RequestPagoTipoChannel) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}

	if r.ChannelId == 0 {
		return fmt.Errorf(tools.ERROR_CHANNEL_ID)
	}
	if r.PagoTipoId == 0 {
		return fmt.Errorf(tools.ERROR_PAGO_TIPO)
	}

	return
}
