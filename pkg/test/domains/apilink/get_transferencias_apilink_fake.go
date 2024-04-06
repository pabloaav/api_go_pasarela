package apilink_tests

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
)

func requestGetTransferenciasCbuInvalido() (requests []*linktransferencia.RequestGetTransferenciasLink) {
	listaInvalidos := []string{"dddddddddddddddddddddd", "", " ", "888888588858885888588", "111111111111"}

	for _, c := range listaInvalidos {

		requestInvalido := linktransferencia.RequestGetTransferenciasLink{

			Cbu:        c,
			Tamanio:    linkdtos.CienTransf,
			Pagina:     1,
			FechaDesde: time.Now(),
			FechaHasta: time.Now(),
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetTransferenciasTamanioInvalido() (requests []*linktransferencia.RequestGetTransferenciasLink) {
	listaInvalidos := []linkdtos.EnumPagiandoTransferencia{" ", "Cinco", "Cuarenta"}

	for _, c := range listaInvalidos {

		requestInvalido := linktransferencia.RequestGetTransferenciasLink{

			Cbu:        "0340218608218026437001",
			Tamanio:    c,
			Pagina:     1,
			FechaDesde: time.Now(),
			FechaHasta: time.Now(),
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}
