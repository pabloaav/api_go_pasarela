package filtros

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type CierreLoteFiltro struct {
	EstadoexternoId        uint64    // query:estadoexternoId db:pagoestadoexternos_id //
	TipoOperacion          string    // query:tipoOperacion db:tipooperacion
	FechaOperacion         bool      // query:fechaOperacion db:fechaoperacion
	FechaPago              bool      // query:fechaPago db:fecha_pago
	FechaCreacion          bool      // query:fechaCreacion db:created_at
	Disputa                bool      // query:disputa db:disputa
	Observacion            bool      // query:observacion db:enobservacion
	FechaInicio            time.Time // query:fechaInicio db:no es campo
	FechaFin               time.Time // query:fechaFin db:no es campo
	CargarMovimientoPrisma bool      // query : no es campo, es para cargar moviminetos relacionado con en cl
	CargarPagoPrisma       bool      // query : no es campo, es para cargar pagos relacionado con en cl
	Number                 uint32    // query:number db:no es campo
	Size                   uint32    // query:size db:no es campo
	CodigosAutorizacion    []string  // query: CodigosAutorizacion: es un string que contiene valores separados por coma y se debe transformar a array
}

type OneCierreLoteFiltro struct {
	Id uint64 // query:id
}

func (clf *CierreLoteFiltro) Validar() (erro error) {
	tiposoperacion := []string{"D", "C", "A"}

	// Tipo de operacion
	if !commons.ContainStrings(tiposoperacion, clf.TipoOperacion) {
		erro = fmt.Errorf(ERROR_PARAM_TIPO_OPERACION)
		return erro
	}

	if !clf.FechaInicio.IsZero() {
		if !validarFechas(clf.FechaOperacion, clf.FechaCreacion, clf.FechaPago) {
			erro = fmt.Errorf(ERROR_PARAM_FECHA)
			return erro
		}
	}

	if !clf.FechaFin.IsZero() {
		if clf.FechaFin.Before(clf.FechaInicio) {
			erro = fmt.Errorf(ERROR_PARAM_FECHA_FIN)
			return erro
		}
	}

	return
}

// dados n parametros, uno y solo uno puede ser verdadero
func validarFechas(params ...bool) (result bool) {

	var count = 0
	for _, value := range params {
		if value {
			count++
		}
	}
	if count == 1 {
		return !result //existe una y solo una opcion de los parametros recibidos que es true
	}
	return
}
