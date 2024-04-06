package apilink_tests

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
)

func requestGetDebinCbuInvalido() (requests []*linkdebin.RequestGetDebinLink) {
	listaInvalidos := []string{"=??¡", "8685663322558861223366", "2563478523685216325888", "111111111111111111111"}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinLink{
			Cbu: c,
			Id:  "1123333",
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetDebinIdInvalido() (requests []*linkdebin.RequestGetDebinLink) {
	listaInvalidos := []string{"", "   "}

	for _, c := range listaInvalidos {

		requestInvalido := linkdebin.RequestGetDebinLink{
			Cbu: "0340218608218026437001",
			Id:  c,
		}
		requests = append(requests, &requestInvalido)
	}

	return requests
}

func requestGetDebinesValidos() (response linkdebin.ResponseGetDebinesLink) {
	return linkdebin.ResponseGetDebinesLink{
		Debines: []linkdebin.DebinesListLink{
			{
				Id:              "G1LMP68NK6RPMK5NR7OEV4",
				Importe:         9999,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "4XJ8G7V95JG5VER9EMPYR0",
				Importe:         7755,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "V8D0Q619LY8X6PL27JZ5RG",
				Importe:         35622,
				Estado:          "ERROR_DEBITO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "EZ4K6DVNOG8V4X195J8LQ7",
				Importe:         50022,
				Estado:          "SIN_SALDO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "MRD06ZO9WEYOM0395GP7XY",
				Importe:         88833,
				Estado:          "RECHAZO_CLIENTE",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "DLMORZP90K03OZ69EGJ468",
				Importe:         12586,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "MRD06ZO9WEYJ5K195GP7XY",
				Importe:         1111,
				Estado:          "ERROR_DATOS",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "G1LMP68NK6R8WLKNR7OEV4",
				Importe:         5599,
				Estado:          "ERROR_ACREDITACION",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "M67REZ8NP1YP5OK24KVGOP",
				Importe:         220055,
				Estado:          "SIN_GARANTIA",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "O7L8GYKNXRYLXPKNMPRZ50",
				Importe:         333333,
				Estado:          "VENCIDO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "WZ0KV8794KM3Y6QNPEYDX4",
				Importe:         885500,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "8PDX4OGNYJYDXZMN0L6EY5",
				Importe:         332211,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "GWY7ZEPN6R7EWVW9Q0M51O",
				Importe:         452511,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
			{
				Id:              "WORD6LEN8Q3W6JO9M1Y30V",
				Importe:         222222,
				Estado:          "ACREDITADO",
				Concepto:        "VAR",
				Moneda:          "ARS",
				Tipo:            "",
				FechaExpiracion: time.Now(),
				Devuelto:        false,
				ContraCargoId:   "",
				Comprador: linkdebin.CompradorDebinesListLink{
					Cuit: "30667754301",
				},
				Vendedor: linkdebin.VendedorDebinesListLink{
					Cuit: "30546676427",
				},
			},
		},
	}
}
