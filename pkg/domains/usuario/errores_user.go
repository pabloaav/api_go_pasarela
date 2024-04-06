package usuario

import "fmt"

type ErrorUser struct {
	Codigo  string `json:"codigo"`
	Message string `json:"message"`
}

func (e *ErrorUser) Error() string {
	return fmt.Sprintf("error - Código: %s, Descripción: %s", e.Codigo, e.Message)
}

const ERROR_URL = "error al crear base url"
const ERROR_CREAR_USUARIO = "no se pudo crear el usuario, intente nuevamente más tarte"
const ERROR_CONSULTA_VACIA = "no se encontraron resultados para esta consulta"
const RESULTADO_NO_ENCONTRADO = "no se encontraron resultados para la busqueda"
const ERROR_CARGAR_USUARIO = "no se pudo cargar el usuario, intente nuevamente más tarde"
const ERROR_MODIFICAR_USUARIO = "no se pudo modificar el usuario, intente nuevamente más tarde"
const ERROR_MODIFICAR_CLIENTEUSER = "no se pudo modificar el usuario del cliente, intente nuevamente más tarde"
const ERROR_ID = "el id es invalido"
const INTENTE_NUEVAMENTE = "Intente nuevamente mas tarte"

func ErrorGuardar(entidad string, generoFeminino bool) string {
	if generoFeminino {
		return fmt.Sprintf("no se pudo guardar la %s. %s", entidad, INTENTE_NUEVAMENTE)
	} else {
		return fmt.Sprintf("no se pudo guardar el %s. %s", entidad, INTENTE_NUEVAMENTE)
	}
}
