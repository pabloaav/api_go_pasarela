package tools

import (
	"errors"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

func EsCbuValido(cbu string, erro string) error {
	serviceCheck := commons.NewAlgoritmoVerificacion()
	err := serviceCheck.ValidarCBU(cbu)

	if err != nil {
		return errors.New(erro)
	}
	return nil
}

func EsCuitValido(cuit string) error {

	err := commons.EsCuilValido(cuit)

	if err != nil {
		return err
	}

	return nil
}

func EsAliasCbuValido(alias string) error {
	if len(alias) < 6 || len(alias) > 20 {
		return errors.New(ERROR_ALIASCBULEN)
	}
	return nil
}

func EsStringVacio(valor string) bool {
	return len(strings.TrimSpace(valor)) == 0

}
