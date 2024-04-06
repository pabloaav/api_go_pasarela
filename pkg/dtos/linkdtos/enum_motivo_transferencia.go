package linkdtos

import "errors"

type EnumMotivoTransferencia string

const (
	AlquilerTransf   EnumMotivoTransferencia = "ALQ"
	CuotasTransf     EnumMotivoTransferencia = "CUO"
	ExpensasTransf   EnumMotivoTransferencia = "EXP"
	FacturasTransf   EnumMotivoTransferencia = "FAC"
	PrestamosTransf  EnumMotivoTransferencia = "PRE"
	SegurosTransf    EnumMotivoTransferencia = "SEG"
	HonorariosTransf EnumMotivoTransferencia = "HON"
	VariosTransf     EnumMotivoTransferencia = "VAR"
)

func (e EnumMotivoTransferencia) IsValid() error {
	switch e {
	case AlquilerTransf, CuotasTransf, ExpensasTransf, FacturasTransf, PrestamosTransf,
		SegurosTransf, HonorariosTransf, VariosTransf:
		return nil
	}
	return errors.New("tipo EnumMotivoTransferencia con formato inválido")
}
