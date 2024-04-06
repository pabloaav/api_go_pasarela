package administracionfake

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

var listaCierrLotePrisma []entities.Prismacierrelote = []entities.Prismacierrelote{
	{
		Model:                gorm.Model{ID: 34},
		PagosUuid:            "3b1c9d8cd6966w4",
		PagoestadoexternosId: 7,
		Tipooperacion:        "C",
		Monto:                68055,
		FechaCierre:          time.Now(),
		ExternalloteId:       789,
		Nombrearchivolote:    "lotetelco_060621.001.txt",
	},

	{
		Model:                gorm.Model{ID: 35},
		PagoestadoexternosId: 7,
		PagosUuid:            "ce86bbf2d468aw4",
		Tipooperacion:        "C",
		Monto:                76051,
		FechaCierre:          time.Now(),
		ExternalloteId:       789,
		Nombrearchivolote:    "lotetelco_060621.001.txt",
	},

	{
		Model:                gorm.Model{ID: 36},
		PagoestadoexternosId: 7,
		PagosUuid:            "d820cfb2db8d8w4",
		Tipooperacion:        "C",
		Monto:                76051,
		FechaCierre:          time.Now(),
		ExternalloteId:       789,
		Nombrearchivolote:    "lotetelco_060621.001.txt",
	},

	{
		Model:                gorm.Model{ID: 37},
		PagoestadoexternosId: 3,
		PagosUuid:            "77c86bb9d8cf5w4",
		Tipooperacion:        "A",
		Monto:                76051,
		FechaCierre:          time.Now(),
		ExternalloteId:       789,
		Nombrearchivolote:    "lotetelco_060621.001.txt",
	},
}
