package cierrelotedtos

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type PrismaRegistroTrailer struct {
	TipoRegistro      string  //TIPOREGISTRO 	   [1-1] 	//"Tipo de Registro, Char, default ""T""."
	CantidadRegistros int64   //CANTIDADREGISTROS  [2-11] 	//"Cantidad Registros ""Detalle"", numérico de 10 dígitos, completando con ""0"" a la izquierda."
	IdMedioPago       int64   //IDMEDIOPAGO 	   [12-14] 	//"Medio de Pago, numérico de 3 dígitos, completando con ""0"" a la izquierda.Por ejemplo: 001 identifica a Visa""
	IdLote            int64   //IDLOTE 	           [15-17]     //Número de Lote, numérico de 3 dígitos (000...999).
	CantidadCompras   int64   //CANTCOMPRAS 	   [18-21] 	//Contador de Compras, numérico de 4 dígitos (0000...9999), cantidad de compras netas.
	MontoCompras      float64 //MONTOCOMPRAS 	   [22-33] 	//Monto de Compras, numérico de 12 dígitos, formato $$$$$$$$$$CC, monto total de compras netas.
	CantidadDevueltas int64   //CANTDEVUELTAS 	   [34-37] 	//Cantidad de Devoluciones, numérico de 4 dígitos (0000...9999), cantidad de devoluciones netas.
	MontoDevueltas    float64 //MONTODEVUELTAS 	   [38-49] 	//Contador de Anulaciones, numérico de 12, formato $$$$$$$$$$CC, cantidad de anulaciones.
	CantidadAnuladas  int64   //CANTANULADAS 	   [50-53] 	//Cantidad de Anulaciones, numérico de 4 dígitos (0000...9999), monto de anulaciones.
	MontoAnulacion    float64 //MONTOANULADAS 	   [54-65] 	//Monto de Anulaciones, numérico de 12 dígitos, formato, monto de anulaciones.
	Filler            string  //FILLER 	           [66-100] //"Filler, 35 caracteres completados con ""0""."
}

func (trailer *PrismaRegistroTrailer) Validar() error {

	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)
	digitCheckFloat := regexp.MustCompile(`[0-9][.]{1}[0-9]{2}$`)

	const ERROR_CAMPO = "la estructura del registro es incorrecto"
	if len(trailer.TipoRegistro) != 1 || commons.StringIsEmpity(trailer.TipoRegistro) {
		return errors.New(ERROR_CAMPO)
	}
	if len(trailer.Filler) != 35 || commons.StringIsEmpity(trailer.Filler) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.CantidadAnuladas, 10)) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.CantidadCompras, 10)) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.CantidadRegistros, 10)) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.CantidadDevueltas, 10)) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.IdMedioPago, 10)) {
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(trailer.IdLote, 10)) {
		return errors.New(ERROR_CAMPO)
	}

	valorStringAnulacion := strconv.FormatFloat(trailer.MontoAnulacion, 'f', 2, 64)
	if !digitCheckFloat.Match([]byte(valorStringAnulacion)) {
		return errors.New(ERROR_CAMPO)
	}
	valorStringDevolucion := strconv.FormatFloat(trailer.MontoDevueltas, 'f', 2, 64)
	if !digitCheckFloat.Match([]byte(valorStringDevolucion)) {
		return errors.New(ERROR_CAMPO)
	}
	valorStringCompra := strconv.FormatFloat(trailer.MontoCompras, 'f', 2, 64)
	if !digitCheckFloat.Match([]byte(valorStringCompra)) {
		return errors.New(ERROR_CAMPO)
	}

	return nil
}
