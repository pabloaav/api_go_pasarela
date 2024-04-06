package linkdebin

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestGetDebinesLink struct {
	Pagina      uint32                     `json:"pagina"`
	Tamanio     linkdtos.EnumPagiandoDebin `json:"tamanio"`
	Cbu         string                     `json:"cbu"` //cbu del vendedor o del comprador
	Estado      linkdtos.EnumEstadoDebin   `json:"estado,omitempty"`
	FechaDesde  time.Time                  `json:"fechadesde"`  //Formato UTC de ISO 8601
	FechaHasta  time.Time                  `json:"fechahasta"`  //Formato UTC de ISO 8601
	EsComprador bool                       `json:"escomprador"` //Indica si la invocación al metodo es como comprador
	Tipo        linkdtos.EnumTipoDebin     `json:"tipo,omitempty"`
}

func (r *RequestGetDebinesLink) IsValid() error {

	err := tools.EsCbuValido(r.Cbu, tools.ERROR_CBU)
	if err != nil {
		return err
	}
	err = r.Tamanio.IsValid()
	if err != nil {
		return errors.New(tools.ERROR_ENUM_PAGINADO_TAMANIO)
	}
	if len(r.Estado) > 0 {
		err = r.Estado.IsValid()
		if err != nil {
			return errors.New(tools.ERROR_ENUM_ESTADO_DEBIN)
		}
	}
	if len(r.Tipo) > 0 {
		err = r.Tipo.IsValid()
		if err != nil {
			return errors.New(tools.ERROR_ENUM_TIPO_DEBIN)
		}
	}
	return nil
}

func (r *RequestGetDebinesLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
