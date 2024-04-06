package linkdebin

import (
	"encoding/json"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
)

type ResponseDebinCreateLink struct {
	Id              string                   `json:"id"`
	FechaOperacion  time.Time                `json:"fechaOperacion"`
	Estado          linkdtos.EnumEstadoDebin `json:"estado"`
	FechaExpiracion time.Time                `json:"fechaExpiracion"`
}

func (r *ResponseDebinCreateLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
