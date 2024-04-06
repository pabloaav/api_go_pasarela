package ribcradtos

import (
	"mime/multipart"
	"github.com/gofiber/fiber/v2"
)


type CargarRiRequest struct {
	Rectificativa bool
	Opera         bool
	Periodo       string
}

type CargarRiCuentaClienteRequest struct {
	Ri []RiCuentaCliente
	CargarRiRequest
}

type CargarRiDatosFondosRequest struct {
	Ri []RiDatosFondos
	CargarRiRequest
}

type CargarRiInfEspecialRequest struct {
	InfEspecial *multipart.FileHeader
	Fiber       *fiber.Ctx
	CargarRiRequest
}

type CargarRiInfEstadisticaRequest struct {
	Ri []RiInfestadistica
	CargarRiRequest
}
