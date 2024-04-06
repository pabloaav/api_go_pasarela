package linkdtos

import (
	"errors"
)

type EnumEstadoDebin string

const (
	Iniciado         EnumEstadoDebin = "INICIADO"
	EnCurso          EnumEstadoDebin = "EN_CURSO"
	RechazoCliente   EnumEstadoDebin = "RECHAZO_CLIENTE"
	SinSaldo         EnumEstadoDebin = "SIN_SALDO"
	ErrorDatos       EnumEstadoDebin = "ERROR_DATOS"
	ErrorDebito      EnumEstadoDebin = "ERROR_DEBITO"
	Vencido          EnumEstadoDebin = "VENCIDO"
	SinGarantia      EnumEstadoDebin = "SIN_GARANTIA"
	ErroAcreditacion EnumEstadoDebin = "ERROR_ACREDITACION"
	Acreditado       EnumEstadoDebin = "ACREDITADO"
	Desconocido      EnumEstadoDebin = "DESCONOCIDO"
)

func (e EnumEstadoDebin) IsValid() error {
	switch e {
	case Iniciado, EnCurso, RechazoCliente, SinSaldo, ErrorDatos, ErrorDebito, Vencido,
		SinGarantia, ErroAcreditacion, Acreditado, Desconocido:
		return nil
	}
	return errors.New("tipo EnumEstadoDebin con formato inválido")
}
