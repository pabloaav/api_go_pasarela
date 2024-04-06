package utilfake

type TableDriverTestConsultarMoviento struct {
	TituloPrueba string
	WantTable    bool
	Cbu          string
}

const ERROR_CAMPO = "error de validación: el tipo de movimiento no es valido"
const ERROR_TIPO = "error de validación: la estructura del registro es incorrecto"
