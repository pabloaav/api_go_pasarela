package prisma

import "fmt"

type ErrorEstructura struct {
	ErrorType        string            `json:"error_type,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	Message          string            `json:"message,omitempty"`
	Code             string            `json:"code,omitempty"`
}

type ValidationError struct {
	Code   string `json:"code,omitempty"`
	Param  string `json:"param,omitempty"`
	Status string `json:"status,omitempty"`
}

func (e *ErrorEstructura) Error() string {
	return fmt.Sprintf(" {ErrorType:%s, ValidationErrors:{Code: %s,Param:  %s,Status:%s} Message:%s,Code:%s,}", e.ErrorType, e.ValidationErrors[0].Code, e.ValidationErrors[0].Param, e.ValidationErrors[0].Status, e.Message, e.Code)
}

func (e *ErrorEstructura) CodeParamsError() string {
	if e.ValidationErrors == nil {
		return fmt.Sprintf(" Code: %s, Param:  %s", e.Code, e.Message)
	}
	return fmt.Sprintf(" Code: %s, Param:  %s", e.ValidationErrors[0].Code, e.ValidationErrors[0].Param)
}

var StatusPrismaMessagge = map[string]string{
	"malformed_request_error": "la estructura de datos enviados no es correcta.", // "Error en el armado del json",
	"authentication_error":    "Error de autenticación.",                         // "ApiKey Inválido",
	"invalid_request_error":   "Los datos no son validos.",                       // "Error por datos inválidos",
	"invalid_status_error":    "Pago anulado o devuelto.",
	"not_found_error":         "Los datos no encontrados.", // "Error con datos no encontrados",
	"api_error":               "Error en el servicio.",     // "Error inesperado en la API REST",
}

func (e *ErrorEstructura) BuscarMensajeError() string {
	for k, v := range StatusPrismaMessagge {
		if k == e.ErrorType {
			return v
		}
	}
	return "Servicio de pagos con tarjeta no se encuentra disponible"
}

const (
	ERRR_TIPO_PAGO                      = "tipo de pago no valido"
	ERRR_INFO_PAGO                      = "el número de pago no existe o es incorrecto"
	ERROR_NUMBER_PAGE_SIZE              = "page size debe ser un valor mayor a 0"
	ERROR_FECHA                         = "formato de fecha no valido, formato valido es (AAAA-MM-DD)"
	ERROR_EXTERNAL_ID                   = "external id no es valida"
	ERROR_PETICION_ANULACION_DEVOLUCION = "error al intentar anular el pago"
	ERROR_CONSULTAR_MSG_ERROR_PRISMA    = "no se pudo obtener mensaje de error relacionado con el codigo"
	ERROR_MSG_VACIO                     = "no existe mensaje de error prisma"
)

// const ERROR_NUMBER_CARD = "el númer de tarjeta no es valido"
// const ERROR_DATE_CARD = "la fecha de vencimiento de la tarjeta es invalido"
// const ERROR_HOLDER_NAME = "nombre no valido"
// const ERROR_SITE_TRANSACTION_ID = "id de ttansacción es incorrecto "
// const ERROR_TOKEN_PAGO = "token de pago no valido"
// const ERROR_BIN = "el número de bien es incorrecto"
// const ERROR_AMOUNT = "el monto ingresado no es valido"
// const ERROR_CURRENCY = "la moneda seleccionada no es valida"
// const ERROR_INSTALLMENTS = "el valor ingresado es incorrecto"
// const ERROR_PAYMENT_TYPE = "el tipo de pago es incorrecto"
// const ERROR_ESTABLISHMENT_NAME = "el número de carcter es superior a 25 o no posee caracter no permitidos"
// const ERROR_NRO_DOC = "el numero de documento ingresado no es valido"
// const ERROR_ESTRUCTURA_INCORRECTA = "los datos recibidos son incorrectos"

// const ERROR_OPEN_ARCHIVO = "error al intentar abrir archivo"
// const ERROR_CREATE_ARCHIVO = "error al intentar crear archivo"
// const ERROR_READ_ARCHIVO = "error al intentar Leer archivos"
// const ERROR_REMOVER_ARCHIVO = "error al intentar eliminar archivo"
// const ERROR_AL_CREAR_NOTIFICACION = "error al intentar guardar notificaciones"
// const ERROR_MOVER_ARCHIVO = "error al intentar mover archivo"
// const ERROR_REMOVER_ARCHIVO = "error al intentar eliminar archivo"

// const ERROR_GENERAL_ARCHIVO = "error no existen archivo de cierre de lote en el directorio"
// const ERROR_LEER_CARACTER = "error al leer strings"
// const ERR_CONVERTIR_ENTERO = "error al convertir a númenro entero"
// const ERR_CONVERTIR_DECIMAL = "error al convertir a númenro decimal"

// const ERROR_PAGO_ESTADO = "error al solicitar lista de estados de pago"
// const ERROR_LEER_DIRECTORIO = "Error al intentar leer los archivos del directorio "
