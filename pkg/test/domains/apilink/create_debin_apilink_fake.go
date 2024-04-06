package apilink_tests

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
)

func _cuentaValidaCreateDebin() linkdebin.CuentaLink {
	return linkdebin.CuentaLink{
		Cbu:      "0340218608218026437001",
		AliasCbu: ""}
}

func _compradorValidoCreateDebin() linkdebin.CompradorCreateDebinLink {
	return linkdebin.CompradorCreateDebinLink{
		Cuit:   "20953043638",
		Cuenta: _cuentaValidaCreateDebin(),
	}
}

func _vendedorValidoCreateDebin() linkdebin.VendedorCreateLink {
	return linkdebin.VendedorCreateLink{
		Cbu:      "0340218608218026437001",
		AliasCbu: "",
	}
}

func _debinValidoCreateDebin() linkdebin.DebinCreateLink {
	return linkdebin.DebinCreateLink{
		ComprobanteId:         "4738999090413568",
		EsCuentaPropia:        true,
		Concepto:              "HON",
		TiempoExpiracion:      5,
		Importe:               1,
		Moneda:                "ARS",
		Recurrente:            false,
		DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
	}
}

func requestCreateCompradorCuentaCbuInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaCbusInvalidos := []string{"123", "ssss", "1df88", "111111111111111111111"}

	for _, c := range listaCbusInvalidos {
		cuentaCbuInvalido := linkdebin.CuentaLink{
			Cbu:      c,
			AliasCbu: ""}
		compradorCuentaCbuInvalido := linkdebin.CompradorCreateDebinLink{
			Cuit:   "20953043638",
			Cuenta: cuentaCbuInvalido,
		}
		debinCompradorCuentaCbuInvalida := linkdebin.RequestDebinCreateLink{
			Comprador: compradorCuentaCbuInvalido,
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     _debinValidoCreateDebin()}
		requests = append(requests, &debinCompradorCuentaCbuInvalida)
	}

	return requests
}

func requestCreateCompradorCuentaAliasCbuInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaAliasCbusInvalidos := []string{"", "ssss", "1df88", "222222222222222222211111111111111", "...."}

	for _, c := range listaAliasCbusInvalidos {
		cuentaAliasCbuInvalido := linkdebin.CuentaLink{
			Cbu:      "",
			AliasCbu: c}
		compradorCuentaAliasCbuInvalido := linkdebin.CompradorCreateDebinLink{
			Cuit:   "20953043638",
			Cuenta: cuentaAliasCbuInvalido,
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: compradorCuentaAliasCbuInvalido,
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     _debinValidoCreateDebin()}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateCompradorCuitInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []string{"", "ssss", "1df88", "222222222222222222211111111111111", "....", "20953043632"}

	for _, c := range listaInvalidos {
		cuentaAliasCbuInvalido := linkdebin.CuentaLink{
			Cbu:      "0340218608218026437001",
			AliasCbu: ""}
		compradorCuentaAliasCbuInvalido := linkdebin.CompradorCreateDebinLink{
			Cuit:   c,
			Cuenta: cuentaAliasCbuInvalido,
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: compradorCuentaAliasCbuInvalido,
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     _debinValidoCreateDebin()}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateVendedorCbuInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaCbusInvalidos := []string{"_(=?", "ADFFGG", "38", "852635f11111111111111"}

	for _, c := range listaCbusInvalidos {
		vendedorCbuInvalido := linkdebin.VendedorCreateLink{
			Cbu:      c,
			AliasCbu: "",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  vendedorCbuInvalido,
			Debin:     _debinValidoCreateDebin()}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateVendedorAliasCbuInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaAliasCbusInvalidos := []string{"", "GTHH", "25FFG", "2886622633222222222211111111111111", "...."}

	for _, c := range listaAliasCbusInvalidos {
		vendedorAliasCbuInvalido := linkdebin.VendedorCreateLink{
			Cbu:      "",
			AliasCbu: c,
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  vendedorAliasCbuInvalido,
			Debin:     _debinValidoCreateDebin()}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateComprobanteIdInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []string{"", "  "}

	for _, c := range listaInvalidos {
		debinComprobanteIdInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         c,
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      5,
			Importe:               1,
			Moneda:                "ARS",
			Recurrente:            false,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinComprobanteIdInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateConceptoInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []linkdtos.EnumConceptoDebin{"Prueba", "  ", "", "11233"}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              c,
			TiempoExpiracion:      5,
			Importe:               1,
			Moneda:                "ARS",
			Recurrente:            false,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateMonedaInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []linkdtos.EnumMoneda{"R$", "Dolar", "", "11233"}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      5,
			Importe:               1,
			Moneda:                c,
			Recurrente:            false,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateImporteInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []int64{-1, -850}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      5,
			Importe:               c,
			Moneda:                "ARS",
			Recurrente:            false,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateTiempoExpiracionInvalido() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []int64{-1, 0, 4325}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      c,
			Importe:               1,
			Moneda:                "ARS",
			Recurrente:            false,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateDescripcionInvalida() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []string{"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"0, 4325bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhfffffffffffffffffffffffffbbbbbbbb"}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      5,
			Importe:               1,
			Moneda:                "ARS",
			Recurrente:            false,
			Descripcion:           c,
			DescripcionPrestacion: "Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignole",
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestCreateDescripcionPrestacionInvalida() (requests []*linkdebin.RequestDebinCreateLink) {
	listaInvalidos := []string{"Topu omoli revfe izalapnuz afufid hem metmik nav uci vuk riuvle ojsaona kodel adoju pawheuh pignoleddd",
		"0, 4325bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhfffffffffffffffffffffffffbbbbbbbb"}

	for _, c := range listaInvalidos {
		debinInvalido := linkdebin.DebinCreateLink{
			ComprobanteId:         "ff222fffff",
			EsCuentaPropia:        true,
			Concepto:              "HON",
			TiempoExpiracion:      5,
			Importe:               1,
			Moneda:                "ARS",
			Recurrente:            false,
			Descripcion:           "",
			DescripcionPrestacion: c,
		}
		requestInvalido := linkdebin.RequestDebinCreateLink{
			Comprador: _compradorValidoCreateDebin(),
			Vendedor:  _vendedorValidoCreateDebin(),
			Debin:     debinInvalido}
		requests = append(requests, &requestInvalido)
	}

	return requests
}
