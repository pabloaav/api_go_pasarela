package administracionfake

import (
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

type TableDriverTestConsultarPagos struct {
	TituloPrueba string
	WantTable    string
	Request      filtros.PagoFiltro
}

const ERROR_DATOS_REQUEST = "error de validación: los datos enviados son incorrectos"
const ERROR_TIPO = "error de validación: la estructura del registro es incorrecto"
