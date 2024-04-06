package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseChannelsArancel struct {
	ChannelArancel []ResponseChannelsAranceles `json:"data"`
	Meta           dtos.Meta                   `json:"meta"`
}

type ResponseChannelsAranceles struct {
	Id            uint
	RubrosId      int64
	Rubro         string
	Importe       float64
	ChannelsId    uint
	Channel       string
	Fechadesde    string
	Tipocalculo   string
	Importeminimo float64
	Importemaximo float64
	Mediopagoid   int64
	Pagocuota     bool
}

func (r *ResponseChannelsAranceles) FromChannelArancel(c entities.Channelarancele) {

	r.Id = c.ID
	r.RubrosId = c.RubrosId
	r.Rubro = c.Rubro.Rubro
	r.Importe = c.Importe
	r.ChannelsId = c.Channel.ID
	r.Channel = c.Channel.Channel
	r.Fechadesde = c.Fechadesde
	r.Tipocalculo = c.Tipocalculo
	r.Importeminimo = c.Importeminimo
	r.Importemaximo = c.Importemaximo
	r.Mediopagoid = c.Mediopagoid
	r.Pagocuota = c.Pagocuota
}

func (r *ResponseChannelsAranceles) FromChArancel(c entities.Channelarancele) {
	r.Id = c.ID
	r.RubrosId = c.RubrosId
	r.Rubro = c.Rubro.Rubro
	r.Importe = c.Importe
	r.ChannelsId = c.Channel.ID
	r.Channel = c.Channel.Channel
	r.Fechadesde = c.Fechadesde
	r.Tipocalculo = c.Tipocalculo
	r.Importeminimo = c.Importeminimo
	r.Importemaximo = c.Importemaximo
	r.Mediopagoid = c.Mediopagoid
	r.Pagocuota = c.Pagocuota

}
