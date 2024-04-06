package administraciondtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestConfiguracion struct {
	Id          uint
	Nombre      string
	Descripcion string
	Valor       string
}

func (r *RequestConfiguracion) ToEntity(cargarId bool) (response entities.Configuracione) {
	if cargarId {
		response.ID = r.Id
	}
	response.Nombre = strings.ToUpper(r.Nombre)
	response.Nombre = strings.TrimSpace(r.Nombre)
	response.Descripcion = r.Descripcion
	response.Valor = r.Valor

	return
}

func (r *RequestConfiguracion) IsValid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}

	if commons.StringIsEmpity(r.Nombre) {
		return fmt.Errorf(tools.ERROR_NOMBRE_CONFIGURACION)
	}

	if commons.StringIsEmpity(r.Valor) {
		return fmt.Errorf(tools.ERROR_VALOR_CONFIGURACION)
	}

	return
}
