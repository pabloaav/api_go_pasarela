package administraciondtos

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type CuentaRequest struct {
	Id                      uint   `json:"id"`
	ClientesID              int64  `json:"clientes_id"`
	RubrosID                uint   `json:"rubros_id"`
	Cuenta                  string `json:"cuenta"`
	Cbu                     string `json:"cbu"`
	Cvu                     string `json:"cvu"`
	Apikey                  string `json:"apikey"`
	DiasRetiroAutomatico    int64  `json:"dias_retiro_automatico"`
	SplitCuentas            bool   `json:"split_cuentas"`
	TransferenciaAutomatica bool   `json:"transferencia_automatica"`
}

func (c *CuentaRequest) IsVAlid(isUpdate bool) (erro error) {

	DiasRetiro := reflect.TypeOf(c.DiasRetiroAutomatico)
	if isUpdate && c.Id < 1 {
		erro = fmt.Errorf(tools.ERROR_ID)
		return
	}

	if c.ClientesID < 1 {
		erro = fmt.Errorf(tools.ERROR_CLIENTE_ID)
		return
	}

	if c.RubrosID < 1 {
		erro = fmt.Errorf(tools.ERROR_RUBRO_ID)
		return
	}

	if commons.StringIsEmpity(c.Cuenta) {
		erro = fmt.Errorf(tools.ERROR_NOMBRE_CUENTA)
		return
	}
	serviceCheck := commons.NewAlgoritmoVerificacion()
	if !(commons.StringIsEmpity(c.Cbu)) {
		erro = serviceCheck.ValidarCBU(c.Cbu)
		if erro != nil {
			return
		}
		if !(commons.StringIsEmpity(c.Cvu)) {
			erro = fmt.Errorf(tools.ERROR_CUENTA)
		}
	} else {
		if commons.StringIsEmpity(c.Cvu) {
			erro = fmt.Errorf(tools.ERROR_CVU)
		}
	}
	if DiasRetiro != reflect.TypeOf(int64(0)) {
		erro = fmt.Errorf(tools.ERROR_DIAS_RETIRO)
	}

	c.Cuenta = strings.ToUpper(c.Cuenta)

	return
}

func (c *CuentaRequest) ToCuenta() (cuenta entities.Cuenta) {
	cuenta.ID = c.Id
	cuenta.ClientesID = c.ClientesID
	cuenta.RubrosID = c.RubrosID
	cuenta.Cuenta = c.Cuenta
	cuenta.Cbu = c.Cbu
	cuenta.Cvu = c.Cvu
	cuenta.Apikey = c.Apikey
	cuenta.DiasRetiroAutomatico = c.DiasRetiroAutomatico
	cuenta.SplitCuentas = c.SplitCuentas
	return

}
