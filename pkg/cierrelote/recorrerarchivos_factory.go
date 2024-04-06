package cierrelote

import (
	"fmt"
	"os"
	adm "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

const (
	CLRECIBIDOS           = "/clrecibidos"
	PAGOPXRECIBIDOS       = "/pagopxrecibidos"
	MOVIMINETOMXRECIBIDOS = "/movimientomxrecibidos"
	RPRECIBIDOS           = "/rapipagorecibidos"
)

type RecorrerArchivosFactory interface {
	GetRecorrerArchivos(m string) (MetodoProcesarArtchivos, error)
}

type recorrerArchivosFactory struct{}

func NewRecorrerArchivos() RecorrerArchivosFactory {
	return &recorrerArchivosFactory{}
}

func (r *recorrerArchivosFactory) GetRecorrerArchivos(m string) (MetodoProcesarArtchivos, error) {
	switch m {
	case CLRECIBIDOS:
		return NewCierreLoteProcesarArchivo(util.Resolve()), nil
	case PAGOPXRECIBIDOS:
		return NewPXProcesarArchivo(util.Resolve()), nil
	case MOVIMINETOMXRECIBIDOS:
		return NewMXProcesarArchivo(util.Resolve()), nil
	case RPRECIBIDOS:
		return NewRPProcesarArchivo(util.Resolve(), adm.Resolve()), nil
	default:
		return nil, fmt.Errorf("el directorio de archivos %v, no es un directorio permitido para ser procesado", m)

	}
}

type MetodoProcesarArtchivos interface { // , rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
	ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo prismaCierreLote.PrismaLogArchivoResponse)
}
