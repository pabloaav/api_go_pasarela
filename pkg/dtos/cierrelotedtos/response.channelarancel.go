package cierrelotedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type ResponseChannelArancel struct {
	ChannelsId    uint
	RubrosId      int64
	Importe       float64
	Fechadesde    string
	Tipocalculo   string
	Importeminimo float64
	Importemaximo float64
	Mediopagoid   int64
	Pagocuota     bool
}

func (rch *ResponseChannelArancel) EntityToDtos(entityChannelArancel entities.Channelarancele) {
	rch.ChannelsId = entityChannelArancel.ChannelsId
	rch.RubrosId = entityChannelArancel.RubrosId
	rch.Importe = entityChannelArancel.Importe
	rch.Fechadesde = entityChannelArancel.Fechadesde
	rch.Tipocalculo = entityChannelArancel.Tipocalculo
	rch.Importeminimo = entityChannelArancel.Importeminimo
	rch.Importemaximo = entityChannelArancel.Importemaximo
	rch.Mediopagoid = entityChannelArancel.Mediopagoid
	rch.Pagocuota = entityChannelArancel.Pagocuota
}
