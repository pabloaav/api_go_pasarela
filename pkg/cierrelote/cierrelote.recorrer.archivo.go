package cierrelote

import (
	"os"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type cierreLoteProcesarArchivos struct {
	utilService util.UtilService
}

func NewCierreLoteProcesarArchivo(util util.UtilService) MetodoProcesarArtchivos {
	return &cierreLoteProcesarArchivos{
		utilService: util,
	}
}

//, rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *cierreLoteProcesarArchivos) ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo prismaCierreLote.PrismaLogArchivoResponse) {
	var ErrorProducido string
	var estado = true
	var estadoInsert = true
	rutaArchivo := strings.Split(archivo.Name(), "/")
	registroDetalle, erro := RecorrerArchivo(archivo)
	if erro != nil {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_RECORRER_ARCHIVOS + erro.Error()
		logs.Error(ErrorProducido)
		logs.Error("no se realizo la insercion de registros de cierre de lote")
	} else {
		/* TODO: recorrer registroDetalle y obtener todos los idmediopago luego: consultar todos los medios de pagos junto con los channels y channelaranceles */
		var arraysMediosPagoIds []int64
		for _, value := range registroDetalle {
			arraysMediosPagoIds = append(arraysMediosPagoIds, value.IdMedioPago)
		}
		listaPagoIntentos, erro := clRepository.GetPagosIntentosByMedioPagoIdRepository(arraysMediosPagoIds)
		if erro != nil {
			estadoInsert = false
			ErrorProducido = ERROR_CONSULTAR_ARANCELES + erro.Error()
			logs.Error(ErrorProducido)
		}
		if len(listaPagoIntentos) == 0 {
			estadoInsert = false
			ErrorProducido = ERROR_MEDIO_PAGO_NO_EXISTE
			logs.Error(ErrorProducido)
		} else {
			listaCierreLote, erro := CrearListaCierreLote(listaPagoIntentos, estadosPagoExterno, archivo.Name(), registroDetalle)
			if erro != nil {
				estadoInsert = false
				ErrorProducido = ERROR_LISTA_CIERRE_LOTE_NO_CREADO + erro.Error()
				logs.Error(ErrorProducido)
			}
			estadoInt, err := clRepository.SaveCierreLoteBatch(listaCierreLote)
			if err != nil {
				estadoInsert = false
				ErrorProducido = ERROR_REGISTRO_EN_DB + err.Error()
				logs.Error(ErrorProducido)
			}
			estadoInsert = estadoInt
		}
	}
	listaLogArchivo = prismaCierreLote.PrismaLogArchivoResponse{
		NombreArchivo:  rutaArchivo[len(rutaArchivo)-1], //archivo.Name(),
		ArchivoLeido:   estado,
		ArchivoMovido:  false,
		LoteInsert:     estadoInsert,
		ErrorProducido: ErrorProducido,
	}
	return
}
