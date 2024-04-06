package linktransferencia

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestGetTransferenciasLink struct {
	Cbu        string                             `json:"cbu"`
	Tamanio    linkdtos.EnumPagiandoTransferencia `json:"tamanio"`
	Pagina     int32                              `json:"pagina"`
	FechaDesde time.Time                          `json:"fechaDesde"`
	FechaHasta time.Time                          `json:"fechaHasta"`
}

func (r *RequestGetTransferenciasLink) IsValid() error {
	err := tools.EsCbuValido(r.Cbu, tools.ERROR_CBU)
	if err != nil {
		return err
	}
	err = r.Tamanio.IsValid()
	if err != nil {
		return errors.New(tools.ERROR_ENUM_PAGINADO_TAMANIO)
	}

	return nil
}


func (r *RequestGetTransferenciasLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
