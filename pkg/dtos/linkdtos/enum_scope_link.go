package linkdtos

import "errors"

type EnumScopeLink string

const (
	DebinRecurrencia                  EnumScopeLink = "DEBIN_RECURRENCIA"
	Debin                             EnumScopeLink = "DEBIN"
	Credin                            EnumScopeLink = "CREDIN"
	Pagos                             EnumScopeLink = "PAGOS"
	ConsultaTitularidad               EnumScopeLink = "CONSULTA_TITULARIDAD"
	ConsultaDestinatario              EnumScopeLink = "CONSULTA_DESTINATARIO"
	BiometriaFacial                   EnumScopeLink = "BIOMETRIA_FACIAL"
	TransferenciasBancariasInmediatas EnumScopeLink = "TRANSFERENCIAS_INMEDIATAS"
	AdhesionCuenta                    EnumScopeLink = "ADHE_VEND"
)

func (e EnumScopeLink) IsValid() error {
	switch e {
	case DebinRecurrencia, Debin, Credin, Pagos, ConsultaTitularidad,
		BiometriaFacial, TransferenciasBancariasInmediatas:
		return nil
	}
	return errors.New("tipo EnumScopeLink con formato inválido")
}
