package administracionfake

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

func PagoEstadoValido() entities.Pagoestado {

	processing := entities.Pagoestado{
		Final:  false,
		Estado: entities.Processing,
	}
	processing.ID = 2

	return processing
}

func ListaPagosSinProcessing() []entities.Pago {
	var listaPagos = make([]entities.Pago, 1)

	mediPagoInvalido := entities.Mediopago{
		ChannelsID: 3,
	}
	var listaPagosIntentos = make([]entities.Pagointento, 1)

	listaPagosIntentos[0] = entities.Pagointento{
		Mediopagos: mediPagoInvalido,
	}

	listaPagos[0] = entities.Pago{
		PagoIntentos: listaPagosIntentos,
	}

	return listaPagos
}

func ListaPagosProcessing() []entities.Pago {
	var listaPagos = make([]entities.Pago, 1)

	mediPagoInvalido := entities.Mediopago{
		ChannelsID: 4,
	}
	var listaPagosIntentos = make([]entities.Pagointento, 1)

	listaPagosIntentos[0] = entities.Pagointento{
		Mediopagos: mediPagoInvalido,
	}

	pago := entities.Pago{
		PagoIntentos: listaPagosIntentos,
	}
	pago.ID = 1
	pago.CreatedAt = time.Now()

	listaPagos[0] = pago

	return listaPagos
}

func ResponseGetDebinesValido() (response *linkdebin.ResponseGetDebinesLink) {

	var listaDebines = make([]linkdebin.DebinesListLink, 1)
	listaDebines[0] = linkdebin.DebinesListLink{
		Id:              "1idDebin",
		Concepto:        "ALQ",
		Moneda:          "ARS",
		Importe:         100000,
		Estado:          "ACREDITADO",
		Tipo:            "DEBIN",
		FechaExpiracion: time.Now(),
		Devuelto:        true,
		ContraCargoId:   "1idContracargo",
		Comprador: linkdebin.CompradorDebinesListLink{
			Cuit:    "20785695147",
			Titular: "Nombre Apellido",
		},
		Vendedor: linkdebin.VendedorDebinesListLink{
			Cuit:    "20546951227",
			Titular: "Otro Nombre Apellido",
		},
	}

	response = &linkdebin.ResponseGetDebinesLink{
		Paginado: linkdtos.PaginadoResponseLink{
			Pagina:          1,
			CantidadPaginas: 0,
		},
		Debines: listaDebines,
	}

	return

}

func ListaPagosExternos() []entities.Pagoestadoexterno {
	var listaPagosExternos = make([]entities.Pagoestadoexterno, 12)

	listaPagosExternos[0] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 1}, Estado: "INICIADO", Vendor: "APILINK", PagoestadosId: 2}
	listaPagosExternos[1] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 2}, Estado: "ERROR DEBITO", Vendor: "APILINK", PagoestadosId: 3}
	listaPagosExternos[2] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 3}, Estado: "SIN SALDO", Vendor: "APILINK", PagoestadosId: 3}
	listaPagosExternos[3] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 4}, Estado: "RECHAZADO DE CLIENTE", Vendor: "APILINK", PagoestadosId: 3}
	listaPagosExternos[4] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 5}, Estado: "ELIMINADO", Vendor: "APILINK", PagoestadosId: 3}
	listaPagosExternos[5] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 6}, Estado: "VENCIDO", Vendor: "APILINK", PagoestadosId: 6}
	listaPagosExternos[6] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 7}, Estado: "ERROR DATOS", Vendor: "APILINK", PagoestadosId: 3}
	listaPagosExternos[7] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 8}, Estado: "EN CURSO", Vendor: "APILINK", PagoestadosId: 4}
	listaPagosExternos[8] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 9}, Estado: "ERROR ACREDITACION", Vendor: "APILINK", PagoestadosId: 5}
	listaPagosExternos[9] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 10}, Estado: "SIN GARANTIA", Vendor: "APILINK", PagoestadosId: 5}
	listaPagosExternos[10] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 11}, Estado: "ACREDITADO", Vendor: "APILINK", PagoestadosId: 7}
	listaPagosExternos[11] = entities.Pagoestadoexterno{Model: gorm.Model{ID: 15}, Estado: "EN CURSO", Vendor: "APILINK", PagoestadosId: 2}

	return listaPagosExternos

}
