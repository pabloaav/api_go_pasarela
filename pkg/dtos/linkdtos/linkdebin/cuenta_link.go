package linkdebin

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type CuentaLink struct {
	Cbu      string `json:"cbu"`
	AliasCbu string `json:"aliasCbu"`
}

func (c *CuentaLink) IsValid() error {

	if len(c.Cbu) == 0 {
		err := tools.EsAliasCbuValido(c.AliasCbu)
		if err != nil {
			return err
		}
	} else {
		err := tools.EsCbuValido(c.Cbu, tools.ERROR_CBU_COMPRADOR)
		if err != nil {
			return err
		}
		//Limpio el alias porque solo debe tener uno el cbu o el alias
		c.AliasCbu = ""
	}

	return nil
}
