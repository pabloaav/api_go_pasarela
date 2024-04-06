package entities

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"gorm.io/gorm"
)

type Apilinkcierrelote struct {
	gorm.Model
	Uuid                 string                     `json:"uuid"`
	DebinId              string                     `json:"debin_id"`
	Concepto             linkdtos.EnumConceptoDebin `json:"concepto"`
	Moneda               linkdtos.EnumMoneda        `json:"moneda"`
	Importe              Monto                      `json:"importe"`
	Estado               linkdtos.EnumEstadoDebin   `json:"estado"`
	Tipo                 linkdtos.EnumTipoDebin     `json:"tipo"`
	FechaExpiracion      time.Time                  `json:"fecha_expiracion"`
	Devuelto             bool                       `json:"devuelto"`
	ContracargoId        string                     `json:"contracargo_id"`
	CompradorCuit        string                     `json:"comprador_cuit"`
	VendedorCuit         string                     `json:"vendedor_cuit"`
	ReferenciaBanco      string                     `json:"referencia_banco"`
	Match                int                        `json:"match"`
	Pagoinformado        bool                       `json:"pagoinformado"`
	BancoExternalId      int                        `json:"banco_external_id"`
	PagoestadoexternosId uint64                     `json:"pagoestadoexternos_id"`
	Fechaacreditacion    time.Time                  `json:"fechaacreditacion"`
	Pagoestadoexterno    Pagoestadoexterno          `json:"pagoestadoexterno" gorm:"foreignKey:pagoestadoexternos_id"`
	FechaCobro           time.Time                  `json:"fecha_cobro"`
}
