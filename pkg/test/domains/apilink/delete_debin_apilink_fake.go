package apilink_tests

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"

func requestDeleteDebinCbuInvalido() (requests []*linkdebin.RequestDeleteDebinLink) {
	listaInvalidos := []string{"23..8'", "55RfGhu", "__r554", "0340218608r18026437001", "0990218608213326437551"}

	for _, inv := range listaInvalidos {

		requestInvalido := linkdebin.RequestDeleteDebinLink{
			Cbu: inv,
			Id:  "1233335",
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestDeleteDebinIdInvalido() (requests []*linkdebin.RequestDeleteDebinLink) {
	listaInvalidos := []string{"", "  "}

	for _, inv := range listaInvalidos {

		requestInvalido := linkdebin.RequestDeleteDebinLink{
			Cbu: "0340218608218026437001",
			Id:  inv,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}
