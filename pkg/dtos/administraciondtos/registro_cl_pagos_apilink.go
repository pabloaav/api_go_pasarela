package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RegistroClPagosApilink struct {
	ListaPagos     []cierrelotedtos.ResponsePagosApilink
	ListaCLApiLink []*entities.Apilinkcierrelote `json:"apilinkcierrelote,omitempty"`
	// ListaCLApiLinkNoAcreditados []*entities.Apilinkcierrelote `json:"apilinkcierrelote_no_acreditados,omitempty"`
}
