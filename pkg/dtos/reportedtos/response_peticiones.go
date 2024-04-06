package reportedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/enumsdtos"

type ResponsePeticiones struct {
	FechaComienzo   string                    `json:"fecha_comienzo"`
	FechaFin        string                    `json:"fecha_fin"`
	TotalPeticiones int                       `json:"total_peticiones"`
	LastPage        int                       `json:"last_page"`
	Data            []ResponseDetallePeticion `json:"data"`
}

type ResponseDetallePeticion struct {
	Operacion string 				`json:"operacion"`
	Fecha     string 				`json:"fecha"`
	Vendor    enumsdtos.EnumVendor  `json:"vendor"`
}
