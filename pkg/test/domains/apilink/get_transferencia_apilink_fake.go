package apilink_tests

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
)

func requestGetTransferenciaCbuInvalido() (requests []*linktransferencia.RequestGetTransferenciaLink) {
	listaInvalidos := []string{"=??¡", "", " ", "2563478523685..", "111111111111"}

	for _, c := range listaInvalidos {

		requestInvalido := linktransferencia.RequestGetTransferenciaLink{
			NumeroReferenciaBancaria: "1123333",
			Cbu:                      c,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetTransferenciaNumeroReferenciaInvalido() (requests []*linktransferencia.RequestGetTransferenciaLink) {
	listaInvalidos := []string{"=??¡", "", " ", "2563478523685..", "111111111111", "33333333333333333333333333333"}

	for _, c := range listaInvalidos {

		requestInvalido := linktransferencia.RequestGetTransferenciaLink{
			NumeroReferenciaBancaria: c,
			Cbu:                      "0340218608218026437001",
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}
