package apilink_tests

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
)

func requestGetDebinesCbuInvalido() (requests []*linkdebin.RequestGetDebinesLink) {
	listaInvalidos := []string{"=??¡", "8685663322558861223366", "2563478523685216325888", "111111111111111111111"}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinesLink{
			Pagina:      5,
			Tamanio:     linkdtos.Cinco,
			Cbu:         c,
			Estado:      linkdtos.Iniciado,
			FechaDesde:  time.Now(),
			FechaHasta:  time.Now(),
			EsComprador: true,
			Tipo:        linkdtos.DebinDefault,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetDebinesTamanioInvalido() (requests []*linkdebin.RequestGetDebinesLink) {
	listaInvalidos := []linkdtos.EnumPagiandoDebin{"", "  ", "325888", "Prueba12", "Prueba"}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinesLink{
			Pagina:      5,
			Tamanio:     c,
			Cbu:         "0340218608218026437001",
			Estado:      linkdtos.Iniciado,
			FechaDesde:  time.Now(),
			FechaHasta:  time.Now(),
			EsComprador: true,
			Tipo:        linkdtos.DebinDefault,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetDebinesEstadoInvalido() (requests []*linkdebin.RequestGetDebinesLink) {
	listaInvalidos := []linkdtos.EnumEstadoDebin{"", "  ", "Aceptado", ".ttt", "Iniciado1"}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinesLink{
			Pagina:      5,
			Tamanio:     linkdtos.Cien,
			Cbu:         "0340218608218026437001",
			Estado:      c,
			FechaDesde:  time.Now(),
			FechaHasta:  time.Now(),
			EsComprador: true,
			Tipo:        linkdtos.DebinDefault,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetDebinesTipoInvalido() (requests []*linkdebin.RequestGetDebinesLink) {
	listaInvalidos := []linkdtos.EnumTipoDebin{"", "  ", "Debinf", "12335", "Defecto"}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinesLink{
			Pagina:      5,
			Tamanio:     linkdtos.Cien,
			Cbu:         "0340218608218026437001",
			Estado:      linkdtos.Iniciado,
			FechaDesde:  time.Now(),
			FechaHasta:  time.Now(),
			EsComprador: true,
			Tipo:        c,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}
