package cierrelotedtos

import "errors"

type EnumDiaSemana string

const (
	Lunes EnumDiaSemana = "Monday"
	Martes  EnumDiaSemana = "Tuesday"
	Miercoles  EnumDiaSemana = "Wednesday"
	Jueves  EnumDiaSemana = "Thursday"
	Viernes  EnumDiaSemana = "Friday"
	Sabado  EnumDiaSemana = "Saturday"
	Domingo  EnumDiaSemana = "Sunday"
)




func (e EnumDiaSemana) IsValid() error {
	switch e {
	case Lunes, Martes, Miercoles, Jueves, Viernes:
		return nil
	}
	return errors.New("sabado y domingo no se concilia con banco")
}