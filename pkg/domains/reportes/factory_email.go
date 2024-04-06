package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

const (
	PAGOS       = "pagos"
	RENDICIONES = "rendiciones"
	REVERTIDOS  = "revertidos"
)

type ReportesSendFactory interface {
	SendEnviarEmail(m string) (Email, error)
}

type enviaremailFactory struct{}

func NewEnviarReportes() ReportesSendFactory {
	return &enviaremailFactory{}
}

func (r *enviaremailFactory) SendEnviarEmail(m string) (Email, error) {
	switch m {
	case PAGOS:
		return SendPagos(util.Resolve()), nil
	case RENDICIONES:
		return SendRendiciones(util.Resolve()), nil
	case REVERTIDOS:
		return SendRevertidos(util.Resolve()), nil
	default:
		return nil, fmt.Errorf("el tipo de reportes a enviar  %v, no es valido", m)

	}
}

type Email interface {
	SendReportes(ruta string, nombreArchivo string, request reportedtos.ResponseClientesReportes) (erro error)
}
