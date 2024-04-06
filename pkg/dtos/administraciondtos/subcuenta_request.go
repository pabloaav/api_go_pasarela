package administraciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type SubcuentaRequest struct {
	Id                  uint    `json:"id"`
	Eliminado           bool    `json:"eliminado"`
	CuentaID            int64   `json:"cuentas_id"`
	Nombre              string  `json:"nombre"`
	Email               string  `json:"email"`
	Modificado          uint    `json:"modificado"`
	Tipo                string  `json:"tipo"`
	Porcentaje          float64 `json:"porcentaje"`
	Cbu                 string  `json:"cbu"`
	AplicaPorcentaje    bool    `json:"aplica_porcentaje"`
	AplicaCostoServicio bool    `json:"aplica_costo_servicio"`
}
type ArraySubcuentaRequest struct {
	ArraySubcuentas []SubcuentaRequest `json:"subcuentas"`
}

func (d *ArraySubcuentaRequest) ArrayIsValidReturnUpdate() (isUpdate bool, erro error) {

	lenArray := len(d.ArraySubcuentas)
	if lenArray <= 0 {
		erro = fmt.Errorf("No se enviaron subcuentas.")
		return false, erro
	}

	lenCreate := 0
	lenUpdate := 0
	for _, v := range d.ArraySubcuentas {

		if v.Id > 0 {
			lenUpdate++
		} else {
			lenCreate++
		}
	}

	if lenUpdate > 0 && lenArray != lenUpdate {
		erro = fmt.Errorf(tools.ERROR_SUBCUENTA_ACTUALIZAR_CREAR)
		return
	}
	if lenCreate > 0 && lenArray != lenCreate {
		erro = fmt.Errorf(tools.ERROR_SUBCUENTA_ACTUALIZAR_CREAR)
		return
	}

	if lenUpdate > 0 {
		return true, nil
	}
	if lenCreate > 0 {
		return false, nil
	}

	return
}

// Funcion que devuelve el tipo de operacion a realizar (ACTUALIZAR, CREAR, CREARACTUALIZAR)
func (d *ArraySubcuentaRequest) ArrayRequestTypeSaved() (typeSaved string, erro error) {

	lenArray := len(d.ArraySubcuentas)
	lenCreate := 0
	lenUpdate := 0
	for _, v := range d.ArraySubcuentas {
		if v.Id > 0 {
			lenUpdate++
		} else {
			lenCreate++
		}
	}

	if lenUpdate > 0 && lenArray == lenUpdate {
		return "ACTUALIZAR", nil
	}
	if lenCreate > 0 && lenArray == lenCreate {
		return "CREAR", nil
	}
	if lenCreate > 0 && lenUpdate > 0 {
		return "CREARACTUALIZAR", nil
	}

	return
}

func (c *SubcuentaRequest) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && c.Id < 1 {
		erro = fmt.Errorf(tools.ERROR_ID)
		return
	}
	if commons.StringIsEmpity(c.Cbu) && !isUpdate {
		erro = fmt.Errorf(tools.ERROR_SUBCUENTA_CBU)
		return
	}
	if commons.StringIsEmpity(c.Nombre) && !isUpdate {
		erro = fmt.Errorf(tools.ERROR_SUBCUENTA_NOMBRE)
		return
	}

	if c.CuentaID <= 0 && !isUpdate {
		erro = fmt.Errorf(tools.ERROR_CUENTAS_ID)
		return
	}
	if c.Porcentaje <= 0 {
		erro = fmt.Errorf(tools.ERROR_SUBCUENTA_PORCENTAJE)
		return
	}

	if commons.StringIsEmpity(c.Tipo) {
		erro = fmt.Errorf(tools.ERROR_TIPO_CUENTA)
		return
	}

	erro = EnumTipoSubcuenta(c.Tipo).IsValid()

	if erro != nil {
		return erro
	}

	serviceCheck := commons.NewAlgoritmoVerificacion()
	if !(commons.StringIsEmpity(c.Cbu)) {
		erro = serviceCheck.ValidarCBU(c.Cbu)
		if erro != nil {
			return
		}
	}

	return
}

func (c *SubcuentaRequest) ToCuenta() (subcuenta entities.Subcuenta) {
	subcuenta.ID = c.Id
	subcuenta.CuentasID = uint(c.CuentaID)
	subcuenta.Nombre = c.Nombre
	subcuenta.Email = c.Email
	subcuenta.Tipo = c.Tipo
	subcuenta.Porcentaje = float64(c.Porcentaje / 100)
	subcuenta.Cbu = c.Cbu
	subcuenta.AplicaCostoServicio = c.AplicaCostoServicio
	subcuenta.AplicaPorcentaje = c.AplicaPorcentaje

	return

}
