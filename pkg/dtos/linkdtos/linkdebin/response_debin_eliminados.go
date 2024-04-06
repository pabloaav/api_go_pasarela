package linkdebin

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"

type ResponseDebinesEliminados struct {
	Id              uint64                   `json:"id"`
	DebinId         string                   `json:"debin_id"`
	Estado          linkdtos.EnumEstadoDebin `json:"estado"`
	Match           int                      `json:"match"`
	BancoExternalId int                      `json:"banco_external_id"`
}
