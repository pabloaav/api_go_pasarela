package utilfake

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"

type TableDriverTestEmailSend struct {
	TituloPrueba string
	WantTable    string
	Request      utildtos.RequestDatosMail
}

// const ERROR_CAMPO = "error de validación: el tipo de movimiento no es valido"
// const ERROR_TIPO = "error de validación: la estructura del registro es incorrecto"
