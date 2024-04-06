package cierrelotedtos

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type PrismaRegistroDetalle struct {
	TipoRegistro       string  // 0 TIPOREGISTRO 	    [1-1] 	Tipo de Registro, Char default "D".
	IdTransacciones    string  // 1 IDTRANSACCIONSITE [2-16] 	Id de Transacción, Alfanumérico de 15 dígitos, completando con "0" a la izquierda.
	IdMedioPago        int64   // 2 IDMEDIOPAGO       [17-19] 	Medio de Pago, numérico, 3 dígitos completando con "0" a la izquierda. Por ejemplo: 001 identifica a Visa.
	NroTarjetaCompleto string  // int64   // 3 NROTARJETACOMPLETO[20-39] 	"Número de Tarjeta, numérico de 20 dígitos. Se informan los seis primeros digitos (BIN), últimos 4 dígitos del número de tarjeta y se completa con ""0"" los digitos restantes."
	TipoOperacion      string  // 4 TIPOOPERACION 	[40-40] 	Operación, Char valores posibles:“C”:Compra, “D”:Devolución, “A”:Anulación.
	Fecha              string  // 5 FECHA 	        [41-48] 	"Fecha de Operación, numérico de 8 dígitos, formato ""DDMMYYYY""."
	Monto              float64 // 6 MONTO 	        [49-60] 	Monto de Operación, numérico de 12 dígitos, 10 enteros + 2 decimales (sin punto decimal).
	CodAut             string  // 7 CODAUT* 	        [61-66] 	"Código de Autorización, numérico de 6 dígitos, completando con ""0"" a la izquierda."
	NroTicket          int64   // 8 NROTICKET 	    [67-72] 	Número de Cupón, numérico de 6 dígitos.
	IdSite             int64   // 9 IDSITE 	        [73-87] 	Id Site Decidir, numérico de 15 dígitos, el Id Site siempre es de 8 dígitos y se completa con 7 ceros a la izquierda.
	IdLote             int64   // 10 IDLOTE 	        [88-90] 	Número de lote, numérico de 3 dígitos.
	Cuotas             int64   // 11 CUOTAS 	        [91-93] 	Cantidad de coutas, numérico de 3 dígitos.
	FechaCierre        string  // 12 FECHACIERRE       [94-101] 	Fecha de cierre,numérico de 8 dígitos.
	NroEstablecimiento int64   // 13 NROESTABLECIMIENTO[102-131] 	Número de establecimiento, numérico de 30 dígitos.
	IdCliente          string  // 14 IDCLIENTE 	    [132-171] 	IDCLIENTE 40 caracteres completados con "0".
	Filler             string  // 15 FILLER 	        [172-190] 	Filler, 19 caracteres completados con "0".

}

const ERROR_CAMPO = "la estructura del registro es incorrecto"

func (detalle *PrismaRegistroDetalle) Validar() error {
	digitCheckInt := regexp.MustCompile(`^[0-9]+$`)
	digitCheckFloat := regexp.MustCompile(`[0-9][.]{1}[0-9]{2}$`)
	regularCheckFecha := regexp.MustCompile(`([0-2][0-9]|3[0-1])(-)(0[1-9]|1[0-2])(-)(\d{4})$`)

	if len(detalle.TipoRegistro) != 1 || commons.StringIsEmpity(detalle.TipoRegistro) {
		logs.Error("procesando cl - nombre campo TipoRegistro " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if len(detalle.IdTransacciones) != 15 || commons.StringIsEmpity(detalle.IdTransacciones) { //if len(detalle.IdTransacciones) != 40 || commons.StringIsEmpity(detalle.IdTransacciones) {
		logs.Error("procesando cl - nombre campo IdTransacciones " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if len(detalle.TipoOperacion) != 1 || commons.StringIsEmpity(detalle.TipoOperacion) {
		logs.Error("procesando cl - nombre campo TipoOperacion " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if len(detalle.Fecha) != 10 {
		logs.Error("procesando cl - nombre campo Fecha " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if commons.StringIsEmpity(detalle.Fecha) {
		logs.Error("procesando cl - nombre campo Fecha " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if !regularCheckFecha.MatchString(detalle.Fecha) {
		logs.Error("procesando cl - nombre campo Fecha " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if len(detalle.FechaCierre) != 10 || commons.StringIsEmpity(detalle.FechaCierre) || !regularCheckFecha.MatchString(detalle.FechaCierre) {
		logs.Error("procesando cl - nombre campo FechaCierre " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if len(detalle.IdCliente) != 40 || commons.StringIsEmpity(detalle.IdCliente) {
		logs.Error("procesando cl - nombre campo IdCliente " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if len(detalle.Filler) != 19 || commons.StringIsEmpity(detalle.Filler) {
		logs.Error("procesando cl - nombre campo Filler " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.IdMedioPago, 10)) {
		logs.Error("procesando cl - nombre campo IdMedioPago " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	// if !digitCheckInt.MatchString(strconv.FormatInt(detalle.NroTarjetaCompleto, 10)) {
	// 	return errors.New(ERROR_CAMPO)
	// }

	if len(detalle.CodAut) != 6 { //digitCheckInt.MatchString(strconv.FormatInt(detalle.CodAut, 10))
		logs.Error("procesando cl - nombre campo CodAut " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	// logs.Info(len(detalle.CodAut))
	// if !digitCheckInt.MatchString(detalle.CodAut) { //digitCheckInt.MatchString(strconv.FormatInt(detalle.CodAut, 10))
	// 	logs.Error("procesando cl - nombre campo TipoRegistro "+ERROR_CAMPO)
	// 	return errors.New(ERROR_CAMPO)
	// }

	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.NroTicket, 10)) {
		logs.Error("procesando cl - nombre campo NroTicket " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.IdSite, 10)) {
		logs.Error("procesando cl - nombre campo IdSite " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.IdLote, 10)) {
		logs.Error("procesando cl - nombre campo IdLote " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.Cuotas, 10)) {
		logs.Error("procesando cl - nombre campo Cuotas " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	if !digitCheckInt.MatchString(strconv.FormatInt(detalle.NroEstablecimiento, 10)) {
		logs.Error("procesando cl - nombre campo NroEstablecimiento " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}
	valorStringMonto := strconv.FormatFloat(detalle.Monto, 'f', 2, 64)
	if !digitCheckFloat.Match([]byte(valorStringMonto)) {
		logs.Error("procesando cl - nombre campo valorStringMonto " + ERROR_CAMPO)
		return errors.New(ERROR_CAMPO)
	}

	return nil
}
