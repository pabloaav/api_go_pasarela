package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseCuentasComision struct {
	CuentasComision []ResponseCuentaComision `json:"data"`
	Meta            dtos.Meta                `json:"meta"`
}

type ResponseCuentaComision struct {
	Id             uint
	CuentasId      uint
	Cuenta         string
	CuentaComision string
	Comision       float64
	ChannelsId     uint
	Channel        string
	//Iva            float64
	VigenciaDesde  *time.Time
	ImporteMinimo  float64
	ImporteMaximo  float64
	Mediopagoid    uint
	Channelarancel ResponseChannelsAranceles
}

func (r *ResponseCuentaComision) FromCuentaComision(c entities.Cuentacomision) {
	var ChannelarancelDTO ResponseChannelsAranceles
	r.Id = c.ID
	r.CuentasId = c.CuentasID
	r.CuentaComision = c.Cuentacomision
	r.Comision = c.Comision
	//r.Iva = c.Iva
	r.VigenciaDesde = c.VigenciaDesde
	r.Cuenta = c.Cuenta.Cuenta
	r.ChannelsId = c.Channel.ID
	r.Channel = c.Channel.Channel
	r.ImporteMinimo = c.Importeminimo
	r.ImporteMaximo = c.Importemaximo
	r.Mediopagoid = c.Mediopagoid
	ChannelarancelDTO.FromChArancel(c.ChannelArancel)
	r.Channelarancel = ChannelarancelDTO
}
