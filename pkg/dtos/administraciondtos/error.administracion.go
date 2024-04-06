package administraciondtos

const ERROR_STRING = "es necesario ingresar "
const ERROR_FECHA_INICIO_INVALIDA = "fecha inicio inválida"
const ERROR_FECHA_FIN_INVALIDA = "fecha fin inválida"
const ERROR_RUTA_INVALIDA = "la ruta informada es inválida"
const ERROR_RI_CODIGO_INVALIDO = "el código de partida es inválido"
const ERROR_RI_SALDO = "el saldo el inválido"
const ERROR_RI_CANTIDAD = "la cantidad es inválida"
const ERROR_RI_CBU = "el CBU es inválido"
const ERROR_RI_NUMERO_FONDO = "el numero del fondo de dinero es inválido"
const ERROR_RI_DENOMINACION_FONDO = "la denominación del fondo es inválida"
const ERROR_RI_AGENTE = "el numero del agente es inválido"
const ERROR_RI_DENOMINACION_AGENTE = "la denominación del agente es inválida"
const ERROR_RI_CUIT_AGENTE = "el cuit del agente es inválido"
const ERROR_RI_MEDIO_PAGO = "el medio de pago es inválido"
const ERROR_RI_ESQUEMA_PAGO = "el esquema de pago es inválido"
const ERROR_RI_MONTO = "el monto es inválido"
const ERROR_RI_TIPO_PRESENTACION = "el tipo de presentación es inválido"
const ERROR_RI_DATOS = "los datos a guardar no pueden ser nulos"

type ErrorSubcuenta struct {
	Error string
	Cbu   string
}
