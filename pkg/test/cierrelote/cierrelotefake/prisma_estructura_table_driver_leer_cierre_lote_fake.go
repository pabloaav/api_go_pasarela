package cierrelotefake

import (
	"io/fs"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

type TableDriverTestCierreLotePagoEstado struct {
	TituloPrueba    string
	WantError       string
	WantEstadosPago []entities.Pagoestado
}

type TableDriverTestCierreLoteArchivos struct {
	TituloPrueba         string
	RutaDirectorioOrigen string
	WantError            string
	WantEstadosPago      []fs.FileInfo
}

func EstructuraCierreLoteFakePagoEstadoError() (tableDriverTestPagoEstado TableDriverTestCierreLotePagoEstado) {
	tableDriverTestPagoEstado = TableDriverTestCierreLotePagoEstado{
		TituloPrueba: "verificar erro cuando la lista de estados se encuentra vacía",
		WantError:    cierrelote.ERROR_PAGO_ESTADO_VACIO,
	}
	return
}

func EstructuraCierreLoteFakePagoEstadoValido() (tableDriverTestPagoEstado TableDriverTestCierreLotePagoEstado) {
	tableDriverTestPagoEstado = TableDriverTestCierreLotePagoEstado{
		TituloPrueba: "verificar que recibe los estado finales de pagos",
		WantEstadosPago: []entities.Pagoestado{
			{
				Model:  gorm.Model{ID: 3},
				Estado: "Rejected",
				Final:  true,
			},
			{
				Model:  gorm.Model{ID: 5},
				Estado: "Reverted",
				Final:  true,
			},
			{
				Model:  gorm.Model{ID: 6},
				Estado: "Expired",
				Final:  true,
			},
			{
				Model:  gorm.Model{ID: 7},
				Estado: "Accredited",
				Final:  true,
			},
		},
	}
	return
}

func EstructuraCierreLoteFakeBuscarArchivosError() (tableDriverTestArchivos TableDriverTestCierreLoteArchivos) {
	tableDriverTestArchivos = TableDriverTestCierreLoteArchivos{
		TituloPrueba:         "verificar erro al intentar leer un directorio",
		RutaDirectorioOrigen: "C:/Users/Sergio/Downloads/archivosLotes/lotesinverificar", //config.RUTA_LOTES_SIN_VERIFICAR,
		WantError:            commons.ERROR_READ_ARCHIVO,                                 //cierrelote.ERROR_LEER_DIRECTORIO,
	}
	return
}

func EstructuraCierreLoteFakeDirectorioVacio() (tableDriverTestArchivos TableDriverTestCierreLoteArchivos) {
	tableDriverTestArchivos = TableDriverTestCierreLoteArchivos{
		TituloPrueba:         "verificar erro al leer un directorio vacio",
		RutaDirectorioOrigen: "C:/Users/Sergio/Downloads/archivosLotes/directoriovacio",
		WantError:            cierrelote.ERROR_GENERAL_ARCHIVO,
	}
	return
}
