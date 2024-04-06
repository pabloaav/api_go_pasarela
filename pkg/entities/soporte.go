package entities

import "gorm.io/gorm"

type Soporte struct {
	gorm.Model
	Nombre  	string 			`json:"nombre"`
	Email   	string 			`json:"email"`
	Visto   	bool 			`json:"visto"`
	Estado  	EnumSoporte 	`json:"estado"`
	Consulta 	string 			`json:"consulta"`
	Abierto 	bool 			`json:"abierto"`
	Archivo		string 			`json:"archivo"`
}

type EnumSoporte string

const (
	espera EnumSoporte = "espera"
	resolviendo EnumSoporte = "resolviendo"
	resuelta EnumSoporte = "resuelta"
	rechazada EnumSoporte = "rechazada"
	pendienteDeDesarrollo EnumSoporte = "pendiente de desarrollo"
)
