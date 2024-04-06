package ribcradtos

import (
	"fmt"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
)

type RiInfestadistica struct {
	CodigoPartida   string
	MedioPago       string
	EsquemaPago     string
	CantOperaciones int64
	MontoTotal      string
}

func (r *RiInfestadistica) IsValid() error {

	mediPago, err := strconv.Atoi(r.MedioPago)

	if err != nil {
		return err
	}

	esquemaPago, err := strconv.Atoi(r.EsquemaPago)

	if err != nil {
		return err
	}

	if len(r.CodigoPartida) > 7 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CODIGO_INVALIDO)
	}

	if mediPago > 5 || len(r.MedioPago) > 3 {
		return fmt.Errorf(administraciondtos.ERROR_RI_MEDIO_PAGO)
	}

	if esquemaPago > 9 || len(r.EsquemaPago) > 3 {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	if len(fmt.Sprint(r.CantOperaciones)) > 11 {
		return fmt.Errorf(administraciondtos.ERROR_RI_CANTIDAD)
	}

	if len(r.MontoTotal) > 11 {
		return fmt.Errorf(administraciondtos.ERROR_RI_MONTO)
	}

	if mediPago == 1 && esquemaPago > 6 {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	if (mediPago == 2 && esquemaPago < 2) || (mediPago == 2 && esquemaPago == 4) || (mediPago == 2 && esquemaPago > 6) {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	if (mediPago == 3 && esquemaPago < 2) || (mediPago == 3 && esquemaPago == 4) || (mediPago == 3 && esquemaPago > 6) {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	if mediPago == 4 && esquemaPago < 7 {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	if mediPago == 5 && esquemaPago != 9 {
		return fmt.Errorf(administraciondtos.ERROR_RI_ESQUEMA_PAGO)
	}

	return nil

}

func (r *RiInfestadistica) ToString() string {

	return fmt.Sprintf("%s;%s;%s;%d;%s", r.CodigoPartida, r.MedioPago, r.EsquemaPago, r.CantOperaciones, r.MontoTotal)

}
