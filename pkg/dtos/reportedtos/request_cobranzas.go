package reportedtos

import (
	"errors"
	"regexp"
)

type RequestCobranzas struct {
	Date string `json:"date"`
}

func (r *RequestCobranzas) Validar() error {

	/* expresion regular para velidar fecha -> formato: año/mes/dia (20210330)*/
	regularCheckFecha := regexp.MustCompile(`([0-2][0-9]|3[0-1])(-)(0[1-9]|1[0-2])(-)(\d{4})$`)
	// regexp.MustCompile(`(\d{4})(-)(0[1-9]|1[0-2])(-)([0-2][0-9]|3[0-1])$`)
	// regularCheckHora := regexp.MustCompile(`([0-2][0-9])(:)([0-5][0-9])(:)([0-5][0-9])$`)

	if len(r.Date) <= 0 {
		return errors.New("se debe indicar una fecha")
	}
	/* FECHA */
	if len(r.Date) != 10 || !regularCheckFecha.MatchString(r.Date) {
		return errors.New("error en el formato de fecha enviado")
	}
	return nil
}
