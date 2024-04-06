package cierreloterapipagodtos

import (
	"errors"
	"regexp"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

// type RapipagoResponse struct {
// 	Rapipago []Rapipago
// }

type Rapipago struct {
	RapipagoHeader   HeaderTrailer
	RapipagoDetalles []Detalles
}

type HeaderTrailer struct {
	Header  Header
	Trailer Trailler
}
type Header struct {
	IdHeader      string
	NombreEmpresa string
	FechaProceso  string
	IdArchivo     string
	FillerHeader  string
}

func (header *Header) ValidarHeader() error {

	const ERROR_CAMPO = "la estructura del header es incorrecto"
	/*validar cantidad de caracteres y que sea string */
	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)

	/* expresion regular para velidar fecha -> formato: año/mes/dia (20210330)*/
	regularCheckFecha := regexp.MustCompile(`(\d{4})()(0[1-9]|1[0-2])()([0-2][0-9]|3[0-1])$`)

	/*validar cantidad de campos y  su contenido */
	/* IDHEADER */
	if !digitCheckInt.MatchString(header.IdHeader) || len(header.IdHeader) != 8 {
		return errors.New(ERROR_CAMPO)
	}
	/* NOMBREEMPRESA */
	if len(header.NombreEmpresa) != 20 || commons.StringIsEmpity(header.NombreEmpresa) {
		return errors.New(ERROR_CAMPO)
	}

	/* FECHAPROCESO */
	if !digitCheckInt.MatchString(header.FechaProceso) || len(header.FechaProceso) != 8 || !regularCheckFecha.MatchString(header.FechaProceso) {
		return errors.New(ERROR_CAMPO)
	}

	/* IDARCHIVO */
	if len(header.IdArchivo) != 20 || commons.StringIsEmpity(header.IdArchivo) {
		return errors.New(ERROR_CAMPO)
	}

	/* FILLERHEADER */
	if len(header.FillerHeader) != 17 {
		return errors.New(ERROR_CAMPO)
	}

	return nil
}

type Trailler struct {
	IdTrailler    string
	CantDetalles  string
	ImporteTotal  string
	FillerTrailer string
}

func (trailler *Trailler) ValidarTrailer() error {

	const ERROR_CAMPO = "la estructura del trailer rapipago es incorrecto"
	/*validar cantidad de caracteres y que sea string */
	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)

	/* IDTRAILER */
	if !digitCheckInt.MatchString(trailler.IdTrailler) || len(trailler.IdTrailler) != 8 {
		return errors.New(ERROR_CAMPO)
	}

	/* CANTDETALLES */
	if !digitCheckInt.MatchString(trailler.CantDetalles) || len(trailler.CantDetalles) != 8 {
		return errors.New(ERROR_CAMPO)
	}

	/* CANTDETALLES */
	if !digitCheckInt.MatchString(trailler.ImporteTotal) || len(trailler.ImporteTotal) != 18 {
		return errors.New(ERROR_CAMPO)
	}

	/* CODIGOBARRAS */
	if len(trailler.FillerTrailer) != 39 {
		return errors.New(ERROR_CAMPO)
	}
	return nil
}

type Detalles struct {
	FechaCobro     string
	ImporteCobrado string
	CodigoBarras   string
	Clearing       string
}

func (detalles *Detalles) ValidarDetalle() error {

	const ERROR_CAMPO = "la estructura del detalle rapipago es incorrecto"
	/*validar cantidad de caracteres y que sea string */
	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)

	/* expresion regular para velidar fecha -> formato: año/mes/dia (20210330)*/
	regularCheckFecha := regexp.MustCompile(`(\d{4})()(0[1-9]|1[0-2])()([0-2][0-9]|3[0-1])$`)

	/* expresion regular para velidar fecha -> formato: dia/mes/año (03032021)*/
	regularCheckFechaClearing := regexp.MustCompile(`([0-2][0-9]|3[0-1])()(0[1-9]|1[0-2])()(\d{4})$`)

	/* FECHACOBRO */
	if !digitCheckInt.MatchString(detalles.FechaCobro) || len(detalles.FechaCobro) != 8 || !regularCheckFecha.MatchString(detalles.FechaCobro) {
		return errors.New(ERROR_CAMPO)
	}

	/* IMPORTECOBRADO */
	if !digitCheckInt.MatchString(detalles.ImporteCobrado) || len(detalles.ImporteCobrado) != 15 {
		return errors.New(ERROR_CAMPO)
	}

	/* CODIGOBARRAS */
	if len(detalles.CodigoBarras) != 48 || commons.StringIsEmpity(detalles.CodigoBarras) {
		return errors.New(ERROR_CAMPO)
	}

	/* CLEARING*/
	if !digitCheckInt.MatchString(detalles.Clearing) || len(detalles.Clearing) != 8 || !regularCheckFechaClearing.MatchString(detalles.Clearing) {
		return errors.New(ERROR_CAMPO)
	}

	return nil
}
