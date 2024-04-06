package administraciondtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ClienteRequest struct {
	Id               uint   `json:"id"`
	Cuit             string `json:"cuit"`
	Razonsocial      string `json:"razon_social"`
	Cliente          string `json:"cliente"`
	Nombrefantasia   string `json:"nombre_fantasia"`
	Email            string `json:"email"`
	Emailcontacto    string `json:"email_contacto"`
	IvaID            int64  `json:"iva_id"`
	IibbID           int64  `json:"iibb_id"`
	Personeria       string `json:"personeria"`
	RetiroAutomatico bool   `json:"retiro_automatico"`
	ReporteBatch     bool   `json:"reporte_batch"`
	NombreReporte    string `json:"nombre_reporte"`
	Domicilio        string `json:"domicilio"`
	SujetoRetencion  bool   `json:"sujeto_retencion"`
	Formulario_8125  bool   `json:"formulario_8125"`
}

func (c *ClienteRequest) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && c.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}

	if commons.StringIsEmpity(c.Cliente) {
		return fmt.Errorf(tools.ERROR_NOMBRE_CLIENTE)
	}

	if commons.StringIsEmpity(c.Razonsocial) {
		return fmt.Errorf(tools.ERROR_RAZON_SOCIAL)
	}

	if !commons.IsEmailValid(c.Email) {
		return fmt.Errorf(tools.ERROR_EMAIL_INVALIDO)
	}

	erro = commons.EsCuilValido(c.Cuit)

	if erro != nil {
		return
	}

	if !(strings.ToUpper(c.Personeria) == "F" || strings.ToUpper(c.Personeria) == "J") {
		erro = fmt.Errorf(tools.ERROR_PERSONERIA)
		return
	}

	c.Cliente = strings.ToUpper(c.Cliente)
	c.Razonsocial = strings.ToUpper(c.Razonsocial)
	c.Nombrefantasia = strings.ToUpper(c.Nombrefantasia)
	c.Personeria = strings.ToUpper(c.Personeria)

	return
}

func (c *ClienteRequest) ToCliente(cargarId bool) (cliente entities.Cliente) {
	if cargarId {
		cliente.ID = c.Id
	}
	cliente.IvaID = c.IvaID
	cliente.IibbID = c.IibbID
	cliente.Cliente = c.Cliente
	cliente.Cuit = c.Cuit
	cliente.Razonsocial = c.Razonsocial
	cliente.Nombrefantasia = c.Nombrefantasia
	cliente.Email = c.Email
	cliente.Emailcontacto = c.Emailcontacto
	cliente.Personeria = c.Personeria
	cliente.RetiroAutomatico = c.RetiroAutomatico
	cliente.ReporteBatch = c.ReporteBatch
	cliente.NombreReporte = c.NombreReporte
	// cliente.SplitCuentas = c.SplitCuentas
	cliente.SujetoRetencion = c.SujetoRetencion
	cliente.Formulario_8125 = c.Formulario_8125

	return

}
