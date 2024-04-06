package pagoofflinedtos

import "fmt"

var (
	ERROR_CAMPO_VACIO                    = "el campo %v es requerido"
	ERROR_CAMPO_LONGITUD_MAXIMA          = "el campo %v no debe superar la longitud de %v digito(s)"
	ERROR_VALIDACION_FECHA               = "el campo %v debe ser menor o igual que el campo %v"
	ERROR_VALIDACION_SEGUNDO_VENCIMIENTO = "el campo %v debe ser menor o igual que el campo %v"
	ERROR_VALIDACION_sIN_VENCIMIENTO     = "error el comprobante no debe poseer segundo vencioneto"
)

func ConstruirMensaje(msj string, valores ...interface{}) (mensaje string) {
	mensaje = fmt.Sprintf(msj, valores...)
	return mensaje
}
