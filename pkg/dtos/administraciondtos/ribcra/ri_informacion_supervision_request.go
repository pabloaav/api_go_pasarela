package ribcradtos

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/gofiber/fiber/v2"
)

type GetInformacionSupervisionRequest struct {
	FechaInicio time.Time
	FechaFin    time.Time
}

func (r *GetInformacionSupervisionRequest) IsValid() (erro error) {

	if r.FechaInicio.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}

	if !r.FechaFin.IsZero() && r.FechaInicio.After(r.FechaFin) {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}
	if r.FechaFin.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA)
		return
	}

	return
}

type BuildInformacionSupervisionRequest struct {
	RiCuentaCliente            []RiCuentaCliente
	RiDatosFondos              []RiDatosFondos
	RectificativaCuentaCliente bool
	RectificativaDatosFondo    bool
	RectificativaInfEspecial   bool
	OperaCuentaCliente         bool
	OperaDatosFondos           bool
	OperaInfEspecial           bool
	Periodo                    string
	InfEspecial                *multipart.FileHeader
	Fiber                      *fiber.Ctx
}
