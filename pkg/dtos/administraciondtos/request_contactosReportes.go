package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type RequestContactosReportes struct {
	ClienteID    uint   `json:"cliente_id"`
	ClienteEmail string `json:"cliente_email"`
	ClienteEmailNuevo string `json:"cliente_email_nuevo"`
}

func (rcr *RequestContactosReportes) DtosToEntity() (entityContatos entities.Contactosreporte) {
	entityContatos.Email = rcr.ClienteEmail
	entityContatos.ClientesID = int64(rcr.ClienteID)
	return 
}
