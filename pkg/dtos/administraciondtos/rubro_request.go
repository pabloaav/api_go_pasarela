package administraciondtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RubroRequest struct {
	Id    uint   `json:"id"`
	Rubro string `json:"rubro"`
}

func (r *RubroRequest) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		erro = fmt.Errorf(tools.ERROR_ID)
		return
	}

	if commons.StringIsEmpity(r.Rubro) {
		erro = fmt.Errorf(tools.ERROR_RUBRO)
		return
	}

	r.Rubro = strings.ToUpper(r.Rubro)

	return
}

func (c *RubroRequest) ToRubro(cargarId bool) (rubro entities.Rubro) {
	if cargarId {
		rubro.ID = c.Id
	}
	rubro.Rubro = c.Rubro

	return

}
