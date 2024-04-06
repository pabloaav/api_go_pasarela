package administraciondtos

import (
	"errors"
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestPagosConsulta struct {
	Uuid               string   `json:"uuid"`
	ExternalReference  string   `json:"external_reference"`
	FechaDesde         string   `json:"fecha_desde"`
	FechaHasta         string   `json:"fecha_hasta"`
	Uuids              []string `json:"uuids"`
	ExternalReferences []string `json:"external_references"`
}

type ParamsValidados struct {
	Uuuid              bool
	ExternalReference  bool
	RangoFecha         bool
	Uuids              bool
	ExternalReferences bool
}

func (r *RequestPagosConsulta) ToPago() entities.Pago {
	return entities.Pago{
		Uuid:              r.Uuid,
		ExternalReference: r.ExternalReference,
	}
}

func (r *RequestPagosConsulta) IsValid() error {
	if len(r.Uuid)+len(r.ExternalReference)+len(r.FechaDesde)+len(r.Uuids)+len(r.ExternalReferences) <= 0 {
		return fmt.Errorf("parámetros de búsqueda insuficientes")
	}
	return nil
}

func (r *RequestPagosConsulta) IsParamsValid() (parametros ParamsValidados, erro error) {
	var msg string
	var paramsRecibidos int64

	if len(r.Uuid) > 0 {
		paramsRecibidos++
		parametros.Uuuid = true

	}
	if len(r.ExternalReference) > 0 {
		paramsRecibidos++
		parametros.ExternalReference = true
	}
	if len(r.FechaDesde)+len(r.FechaHasta) > 0 {
		paramsRecibidos++
		parametros.RangoFecha = true
	}
	if len(r.Uuids) > 0 {
		paramsRecibidos++
		parametros.Uuids = true
	}
	if len(r.ExternalReferences) > 0 {
		paramsRecibidos++
		parametros.ExternalReferences = true
	}

	if paramsRecibidos < 1 {
		msg = "debe enviar parametro de consulta."
	}
	if paramsRecibidos > 1 {
		msg = "solo debe enviar un parametro de consulta."
	}

	if len(msg) > 0 {
		erro = errors.New(msg)
	}
	return
}
