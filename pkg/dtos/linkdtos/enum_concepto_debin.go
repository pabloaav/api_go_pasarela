package linkdtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type EnumConceptoDebin string

const (
	Alquiler                           EnumConceptoDebin = "ALQ"
	AportesCapital                     EnumConceptoDebin = "APC"
	BieneRegistrablesHabituales        EnumConceptoDebin = "BRH"
	BieneRegistrablesNoHabituales      EnumConceptoDebin = "BRN"
	CompraPEI                          EnumConceptoDebin = "CCT"
	Cuotas                             EnumConceptoDebin = "CUO"
	DevolucionPEI                      EnumConceptoDebin = "DCT"
	EstadoExpropiaciones               EnumConceptoDebin = "ESE"
	Expensas                           EnumConceptoDebin = "EXP"
	Facturas                           EnumConceptoDebin = "FAC"
	Haberes                            EnumConceptoDebin = "HAB"
	Honorarios                         EnumConceptoDebin = "HON"
	OperacionesInmobiliarias           EnumConceptoDebin = "OIN"
	OperacionesInmobiliariasHabituales EnumConceptoDebin = "OIH"
	PagoPEI                            EnumConceptoDebin = "PCT"
	Prestamos                          EnumConceptoDebin = "PRE"
	ReitegrosObrasSociales             EnumConceptoDebin = "ROP"
	Seguros                            EnumConceptoDebin = "SEG"
	SiniestrosSeguros                  EnumConceptoDebin = "SIS"
	SuscripcionObligaciones            EnumConceptoDebin = "SON"
	Varios                             EnumConceptoDebin = "VAR"
)

func (e EnumConceptoDebin) IsValid() error {
	switch e {
	case Alquiler, AportesCapital, BieneRegistrablesHabituales, BieneRegistrablesNoHabituales, CompraPEI,
		Cuotas, DevolucionPEI, EstadoExpropiaciones, Expensas, Facturas, Haberes, Honorarios, OperacionesInmobiliarias,
		OperacionesInmobiliariasHabituales, PagoPEI, Prestamos, ReitegrosObrasSociales, Seguros, SiniestrosSeguros, SuscripcionObligaciones,
		Varios:
		return nil
	}
	return errors.New(tools.ERROR_ENUM_CONCEPTO)
}
