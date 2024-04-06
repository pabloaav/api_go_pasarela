package apilink_tests

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
)

func _origenValidoCreateTransferencia() linktransferencia.OrigenTransferenciaLink {
	return linktransferencia.OrigenTransferenciaLink{
		Cbu: "0340218608218026437001",
	}
}

func _destinoValidoCreateTransferencia() linktransferencia.DestinoTransferenciaLink {
	return linktransferencia.DestinoTransferenciaLink{
		Cbu:            "0340218608218026437001",
		AliasCbu:       "",
		EsMismoTitular: false,
	}
}

func _CreateTransferenciaValida() linktransferencia.RequestTransferenciaCreateLink {
	return linktransferencia.RequestTransferenciaCreateLink{
		Origen:     _origenValidoCreateTransferencia(),
		Destino:    _destinoValidoCreateTransferencia(),
		Importe:    1,
		Moneda:     linkdtos.Pesos,
		Motivo:     linkdtos.AlquilerTransf,
		Referencia: "1235546658",
	}
}

func requestCreateTransferenciaOrigenCbuInvalido() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []string{"52r555", ".:___s", "1df88", "0342218608258026437008"}

	for _, c := range listaInvalidos {

		origenCbuInvalido := linktransferencia.OrigenTransferenciaLink{
			Cbu: c,
		}

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     origenCbuInvalido,
			Destino:    _destinoValidoCreateTransferencia(),
			Importe:    1,
			Moneda:     linkdtos.Pesos,
			Motivo:     linkdtos.AlquilerTransf,
			Referencia: "1235546658",
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}

func requestCreateTransferenciaDestinoCbuInvalido() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []string{"5porti55", ".:222222222222222222_s", "5554455665335555588555", "0342218608258026437008"}

	for _, c := range listaInvalidos {

		destinoCbuInvalido := linktransferencia.DestinoTransferenciaLink{
			Cbu:            c,
			AliasCbu:       "",
			EsMismoTitular: false,
		}

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     _origenValidoCreateTransferencia(),
			Destino:    destinoCbuInvalido,
			Importe:    1,
			Moneda:     linkdtos.Pesos,
			Motivo:     linkdtos.AlquilerTransf,
			Referencia: "1235546658",
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}

func requestCreateTransferenciaDestinoAliasInvalido() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []string{"23", ".,:222222222222222222_s", "  ", ""}

	for _, c := range listaInvalidos {

		destinoInvalido := linktransferencia.DestinoTransferenciaLink{
			Cbu:            "",
			AliasCbu:       c,
			EsMismoTitular: false,
		}

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     _origenValidoCreateTransferencia(),
			Destino:    destinoInvalido,
			Importe:    1,
			Moneda:     linkdtos.Pesos,
			Motivo:     linkdtos.AlquilerTransf,
			Referencia: "1235546658",
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}

func requestCreateTransferenciaMonedaInvalida() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []linkdtos.EnumMoneda{"Real", "R$", "  ", "", "123"}

	for _, c := range listaInvalidos {

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     _origenValidoCreateTransferencia(),
			Destino:    _destinoValidoCreateTransferencia(),
			Importe:    0,
			Moneda:     c,
			Motivo:     linkdtos.AlquilerTransf,
			Referencia: "1235546658",
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}

func requestCreateTransferenciaMotivoInvalida() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []linkdtos.EnumMotivoTransferencia{"ALQFFF", "12$", "  ", "", "Alquiler"}

	for _, c := range listaInvalidos {

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     _origenValidoCreateTransferencia(),
			Destino:    _destinoValidoCreateTransferencia(),
			Importe:    0,
			Moneda:     linkdtos.Pesos,
			Motivo:     c,
			Referencia: "1235546658",
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}

func requestCreateTransferenciaReferenciaInvalida() (requests []*linktransferencia.RequestTransferenciaCreateLink) {
	listaInvalidos := []string{"11112233.", "123rrr,", "  ", ":Alquiler"}

	for _, c := range listaInvalidos {

		createTransferenciaInvalida := linktransferencia.RequestTransferenciaCreateLink{
			Origen:     _origenValidoCreateTransferencia(),
			Destino:    _destinoValidoCreateTransferencia(),
			Importe:    0,
			Moneda:     linkdtos.Pesos,
			Motivo:     linkdtos.AlquilerTransf,
			Referencia: c,
		}
		requests = append(requests, &createTransferenciaInvalida)
	}

	return requests
}
