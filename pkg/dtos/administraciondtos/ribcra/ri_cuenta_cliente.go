package ribcradtos

import (
	"fmt"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)


type RiCuentaCliente struct {
	CodigoPartida string
	Saldo         string
	Cantidad      string
	CBU           string
	Orden         int
}

func (r *RiCuentaCliente) IsValid() error {
	if len(r.CodigoPartida) > 7 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CODIGO_INVALIDO)
	}
	if len(r.Saldo) > 11 {
		return fmt.Errorf(administraciondtos.ERROR_RI_SALDO)
	}
	if len(r.Cantidad) > 11 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CANTIDAD)
	}
	if len(r.CBU) > 22 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CBU)
	}
	return nil
}

func (r *RiCuentaCliente) ToString() string {

	return fmt.Sprintf("%s;%s;%s;%s", r.CodigoPartida, r.Saldo, r.Cantidad, r.CBU)

}
