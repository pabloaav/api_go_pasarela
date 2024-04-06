package administraciondtos

import (
	"errors"
	"mime/multipart"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type RequestSoporte struct {
	Id		 uint					`json:"id"`
	Nombre   string                `json:"nombre"`
	Email    string                `json:"email"`
	Visto    bool                  `json:"visto"`
	Estado   EnumEstadoSoporte     `json:"estado"`
	Consulta string                `json:"consulta"`
	Abierto  bool                  `json:"abierto"`
	File     *multipart.FileHeader `json:"archivo"`
}

type EnumEstadoSoporte string

const (
	espera                EnumEstadoSoporte = "espera"
	resolviendo           EnumEstadoSoporte = "resolviendo"
	resuelta              EnumEstadoSoporte = "resuelta"
	rechazada             EnumEstadoSoporte = "rechazada"
	pendienteDeDesarrollo EnumEstadoSoporte = "pendiente de desarrollo"
)

func (e EnumEstadoSoporte) IsValid() bool {
	switch e {
	case espera, resolviendo, resuelta, rechazada, pendienteDeDesarrollo:
		return true
	}
	return false
}

// recuperar el string del EnumEstamultipart
func (e EnumEstadoSoporte) ToString() string {
	switch e {
	case espera:
		return "espera"
	case resolviendo:
		return "resolviendo"
	case resuelta:
		return "resuelta"
	case rechazada:
		return "rechazada"
	case pendienteDeDesarrollo:
		return "pendiente de desarrollo"
	default:
		return ""
	}
}
func (r *RequestSoporte) IsValidCreate() (err error){
	isEmailValid := commons.IsEmailValid(r.Email)
	if !isEmailValid {
		err= errors.New("debe enviar un correo valido")
		return 
	}
	if  len(r.Nombre)>44 || len(r.Nombre)<1{
		err= errors.New("debe enviar un nombre que no supere los 45 caracteres")
		return
	}
	return nil
}
func (r *RequestSoporte) IsValidPut() (err error){
	if !r.Visto && !r.Abierto && r.Estado==""{
		err= errors.New("debe enviar por lo menos un campo para actualizar")
		return
	}
	if  r.Id<1{
		err= errors.New("debe enviar un id valido")
		return
	}
	if r.Estado!= ""{
		isEnumValid := r.Estado.IsValid()
		if !isEnumValid {
			err= errors.New("debe enviar un estado valido")
			return

		}
	}
	return nil
}