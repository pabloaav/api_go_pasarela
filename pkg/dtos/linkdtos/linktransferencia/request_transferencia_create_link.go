package linktransferencia

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestTransferenciaCreateLink struct {
	Origen     OrigenTransferenciaLink          `json:"origen"`
	Destino    DestinoTransferenciaLink         `json:"destino"`
	Importe    entities.Monto                   `json:"importe"`
	Moneda     linkdtos.EnumMoneda              `json:"moneda"`
	Motivo     linkdtos.EnumMotivoTransferencia `json:"motivo"`
	Referencia string                           `json:"referencia"`
}

func (r *RequestTransferenciaCreateLink) IsValid() error {

	err := r.Origen.IsValid()
	if err != nil {
		return err
	}

	err = r.Destino.IsValid()

	if err != nil {
		return err
	}
	err = r.Moneda.IsValid()
	if err != nil {
		return errors.New(tools.ERROR_ENUM_MONEDA)
	}
	err = r.Motivo.IsValid()
	if err != nil {
		return errors.New(tools.ERROR_ENUM_MOTIVO)
	}

	referenciaValida, err := regexp.MatchString("^[a-zA-Z0-9]*$", r.Referencia)

	if err != nil || !referenciaValida {
		return errors.New(tools.ERROR_REFERENCIA)
	}

	return nil
}

type OrigenTransferenciaLink struct {
	Cbu string `json:"cbu"`
}

func (o *OrigenTransferenciaLink) IsValid() error {

	err := tools.EsCbuValido(o.Cbu, tools.ERROR_CBU)
	if err != nil {
		return err
	}

	return nil
}

type DestinoTransferenciaLink struct {
	Cbu            string `json:"cbu"`
	AliasCbu       string `json:"alias"`
	EsMismoTitular bool   `json:"esMismoTitular"`
}

func (d *DestinoTransferenciaLink) IsValid() error {

	if len(d.Cbu) == 0 {
		err := tools.EsAliasCbuValido(d.AliasCbu)
		if err != nil {
			return err
		}
	} else {
		err := tools.EsCbuValido(d.Cbu, tools.ERROR_CBU)
		if err != nil {
			return err
		}
		//Limpio el alias porque solo debe tener uno el cbu o el alias
		d.AliasCbu = ""
	}

	return nil
}

func (r *RequestTransferenciaCreateLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
