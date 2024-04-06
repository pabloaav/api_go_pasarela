package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseApilinkCierresLotes struct {
	CierresLotes []ResponseApilinkCL `json:"cierreslotes"`
	Meta         dtos.Meta           `json:"meta"`
}

type ResponseApilinkCL struct {
	Id                   int64
	Uuid                 string                     `json:"uuid"`
	DebinId              string                     `json:"debin_id"`
	Concepto             linkdtos.EnumConceptoDebin `json:"concepto"`
	Moneda               linkdtos.EnumMoneda        `json:"moneda"`
	Importe              entities.Monto             `json:"importe"`
	Estado               linkdtos.EnumEstadoDebin   `json:"estado"`
	Tipo                 linkdtos.EnumTipoDebin     `json:"tipo"`
	FechaExpiracion      time.Time                  `json:"fecha_expiracion"`
	Devuelto             bool                       `json:"devuelto"`
	ContracargoId        string                     `json:"contracargo_id"`
	CompradorCuit        string                     `json:"comprador_cuit"`
	VendedorCuit         string                     `json:"vendedor_cuit"`
	ReferenciaBanco      string                     `json:"referencia_banco"`
	Match                int                        `json:"match"`
	BancoExternalId      int                        `json:"banco_external_id"`
	PagoestadoexternosId uint64                     `json:"pagoestadoexternos_id"`
	Fechaacreditacion    time.Time                  `json:"fechaacreditacion"`
	CreatedAt            time.Time                  `json:"created_at"`
	PagoEstado           string                     `json:"pago_estado"`
}

func (apilinkcl *ResponseApilinkCL) EntityToDto(entityApilinkCL entities.Apilinkcierrelote) {
	apilinkcl.Id = 0
	if entityApilinkCL.ID > 0 {
		apilinkcl.Id = int64(entityApilinkCL.ID)
	}
	apilinkcl.Uuid = entityApilinkCL.Uuid
	apilinkcl.DebinId = entityApilinkCL.DebinId
	apilinkcl.Concepto = entityApilinkCL.Concepto
	apilinkcl.Moneda = entityApilinkCL.Moneda
	apilinkcl.Importe = entityApilinkCL.Importe
	apilinkcl.Estado = entityApilinkCL.Estado
	apilinkcl.Tipo = entityApilinkCL.Tipo
	apilinkcl.FechaExpiracion = entityApilinkCL.FechaExpiracion
	apilinkcl.Devuelto = entityApilinkCL.Devuelto
	apilinkcl.ContracargoId = entityApilinkCL.ContracargoId
	apilinkcl.CompradorCuit = entityApilinkCL.CompradorCuit
	apilinkcl.VendedorCuit = entityApilinkCL.VendedorCuit
	apilinkcl.ReferenciaBanco = entityApilinkCL.ReferenciaBanco
	apilinkcl.Match = entityApilinkCL.Match
	apilinkcl.BancoExternalId = 0
	if entityApilinkCL.BancoExternalId > 0 {
		apilinkcl.BancoExternalId = entityApilinkCL.BancoExternalId
	}
	apilinkcl.PagoestadoexternosId = entityApilinkCL.PagoestadoexternosId
	if entityApilinkCL.PagoestadoexternosId > 0 {
		apilinkcl.PagoestadoexternosId = entityApilinkCL.PagoestadoexternosId
	}
	apilinkcl.Fechaacreditacion = entityApilinkCL.Fechaacreditacion
	apilinkcl.CreatedAt = entityApilinkCL.CreatedAt
	// chequear si no viene vacio la relacion pago estado externo debtro de la entidad apilink
	if (entities.Pagoestadoexterno{} != entityApilinkCL.Pagoestadoexterno) {
		apilinkcl.PagoEstado = entityApilinkCL.Pagoestadoexterno.Estado
	}

}
