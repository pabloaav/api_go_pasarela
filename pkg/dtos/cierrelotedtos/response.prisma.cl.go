package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePrismaCierresLotes struct {
	CierresLotes []ResponsePrismaCL `json:"cierreslotes"`
	Meta         dtos.Meta          `json:"meta"`
}

type ResponsePrismaCL struct {
	Id                         int64
	PagoestadoexternosId       int64                      `json:"pagoestadoexternos_id"`
	ChannelarancelesId         int64                      `json:"channelaranceles_id"`
	ImpuestosId                int64                      `json:"impuestos_id"`
	PrismamovimientodetallesId int64                      `json:"prismamovimientodetalles_id"`
	PrismamovimientodetalleId  int64                      `json:"prismamovimientodetalle_id"`
	PrismatrdospagosId         int64                      `json:"prismatrdospagos_id"`
	BancoExternalId            int64                      `json:"banco_external_id"`
	Tiporegistro               string                     `json:"tiporegistro"`
	PagosUuid                  string                     `json:"pagos_uuid"`
	ExternalmediopagoId        int64                      `json:"externalmediopago"`
	Nrotarjeta                 string                     `json:"nrotarjeta"`
	Tipooperacion              entities.EnumTipoOperacion `json:"tipooperacion"`
	Fechaoperacion             time.Time                  `json:"fechaoperacion"`
	Monto                      entities.Monto             `json:"monto"`
	Montofinal                 entities.Monto             `json:"montofinal"`
	Valorpresentado            entities.Monto             `json:"valorpresentado"`
	Diferenciaimporte          entities.Monto             `json:"difernciaimporte"`
	Coeficientecalculado       float64                    `json:"coeficientecalculado"`
	Costototalporcentaje       float64                    `json:"costototalporcentaje"`
	Importeivaarancel          float64                    `json:"importeivaarancel"`
	Descripcion                string                     `json:"descripcion"`
	Codigoautorizacion         string                     `json:"codigoautorizacion"`
	Nroticket                  int64                      `json:"nroticket"`
	SiteID                     int64                      `json:"site_id"`
	ExternalloteId             int64                      `json:"externallote_id"`
	Nrocuota                   int64                      `json:"nrocuota"`
	FechaCierre                time.Time                  `json:"fecha_cierre"`
	Nroestablecimiento         int64                      `json:"nroestablecimiento"`
	ExternalclienteID          string                     `json:"externalcliente_id"`
	Nombrearchivolote          string                     `json:"nombrearchivolote"`
	Match                      int                        `json:"match"`
	FechaPago                  time.Time                  `json:"fecha_pago"`
	Disputa                    bool                       `json:"disputa"`
	Cantdias                   int                        `json:"cantdias"`
	Enobservacion              bool                       `json:"enobservacion"`
	Channelarancel             ResponseChannelArancel     `json:"channelarancel"`
	Descripcionpresentacion    string                     `json:"descripcionpresentacion"`
	Istallmentsinfo            ResponseInstallmentInfo    `json:"istallmentsinfo"`
	DetallemovimientoId        int64                      `json:"detallemovimineto_id"`
	Reversion                  bool                       `json:"reversion"`
	DetallepagoId              int64                      `json:"detallepago_id"`
	Descripcioncontracargo     string                     `json:"descripcioncontracargo"`
	ExtbancoreversionId        int64                      `json:"extbancoreversion_id"`
	Conciliado                 bool                       `json:"conciliacion"`
	Estadomovimiento           bool                       `json:"estadomovimineto"`
	Descripcionbanco           string                     `json:"descripcionbanco"`
}

func (rcl *ResponsePrismaCL) EntityToDtos(entityPrismaCL entities.Prismacierrelote) {
	rcl.Id = 0
	if entityPrismaCL.ID > 0 {
		rcl.Id = int64(entityPrismaCL.ID)
	}
	rcl.PagoestadoexternosId = entityPrismaCL.PagoestadoexternosId
	rcl.ChannelarancelesId = entityPrismaCL.ChannelarancelesId
	rcl.ImpuestosId = entityPrismaCL.ImpuestosId
	rcl.PrismamovimientodetallesId = 0
	if entityPrismaCL.PrismamovimientodetallesId > 0 {
		rcl.PrismamovimientodetallesId = entityPrismaCL.PrismamovimientodetallesId
	}

	rcl.PrismamovimientodetalleId = 0
	if entityPrismaCL.PrismamovimientodetalleId > 0 {
		rcl.PrismamovimientodetalleId = entityPrismaCL.PrismamovimientodetalleId
	}

	rcl.PrismatrdospagosId = 0
	if entityPrismaCL.PrismatrdospagosId > 0 {
		rcl.PrismatrdospagosId = entityPrismaCL.PrismatrdospagosId
	}
	rcl.BancoExternalId = entityPrismaCL.BancoExternalId
	rcl.Tiporegistro = entityPrismaCL.Tiporegistro
	rcl.PagosUuid = entityPrismaCL.PagosUuid
	rcl.ExternalmediopagoId = entityPrismaCL.ExternalmediopagoId
	rcl.Nrotarjeta = entityPrismaCL.Nrotarjeta
	rcl.Tipooperacion = entityPrismaCL.Tipooperacion
	rcl.Fechaoperacion = entityPrismaCL.Fechaoperacion
	rcl.Monto = entityPrismaCL.Monto
	rcl.Montofinal = entityPrismaCL.Montofinal
	rcl.Codigoautorizacion = entityPrismaCL.Codigoautorizacion
	rcl.Nroticket = entityPrismaCL.Nroticket
	rcl.SiteID = entityPrismaCL.SiteID
	rcl.ExternalloteId = entityPrismaCL.ExternalloteId
	rcl.Nrocuota = entityPrismaCL.Nrocuota
	rcl.FechaCierre = entityPrismaCL.FechaCierre
	rcl.Nroestablecimiento = entityPrismaCL.Nroestablecimiento
	rcl.ExternalclienteID = entityPrismaCL.ExternalclienteID
	rcl.Nombrearchivolote = entityPrismaCL.Nombrearchivolote
	rcl.Match = entityPrismaCL.Match
	rcl.FechaPago = entityPrismaCL.FechaPago
	rcl.Disputa = entityPrismaCL.Disputa
	rcl.Enobservacion = entityPrismaCL.Enobservacion
	rcl.Valorpresentado = entityPrismaCL.Valorpresentado
	rcl.Diferenciaimporte = entityPrismaCL.Diferenciaimporte
	rcl.Coeficientecalculado = entityPrismaCL.Coeficientecalculado
	rcl.Costototalporcentaje = entityPrismaCL.Costototalporcentaje
	rcl.Importeivaarancel = entityPrismaCL.Importeivaarancel
	rcl.Descripcion = entityPrismaCL.Descripcion
	rcl.Descripcionpresentacion = entityPrismaCL.Descripcionpresentacion
	if entityPrismaCL.Channelarancel != nil {
		rcl.Channelarancel.EntityToDtos(*entityPrismaCL.Channelarancel)
	}
	rcl.DetallemovimientoId = entityPrismaCL.DetallemovimientoId
	rcl.DetallepagoId = entityPrismaCL.DetallepagoId
	rcl.Descripcioncontracargo = entityPrismaCL.Descripcioncontracargo
	rcl.Reversion = entityPrismaCL.Reversion
	rcl.ExtbancoreversionId = entityPrismaCL.ExtbancoreversionId
	rcl.Conciliado = entityPrismaCL.Conciliado
	rcl.Estadomovimiento = entityPrismaCL.Estadomovimiento
	rcl.Descripcionbanco = entityPrismaCL.Descripcionbanco

}

/*
PrismamovimientodetalleId
PrismapagotrdosId
*/

func (rcl *ResponsePrismaCL) DtosToEntity() (entityPrismaCL entities.Prismacierrelote) {
	entityPrismaCL.ID = 0
	if rcl.Id > 0 {
		entityPrismaCL.ID = uint(rcl.Id)
	}
	entityPrismaCL.PagoestadoexternosId = rcl.PagoestadoexternosId
	entityPrismaCL.ChannelarancelesId = rcl.ChannelarancelesId
	entityPrismaCL.ImpuestosId = rcl.ImpuestosId

	entityPrismaCL.PrismamovimientodetallesId = 0
	if rcl.PrismamovimientodetallesId > 0 {
		entityPrismaCL.PrismamovimientodetallesId = rcl.PrismamovimientodetallesId
	}
	entityPrismaCL.PrismamovimientodetalleId = 0
	if rcl.PrismamovimientodetalleId > 0 {
		entityPrismaCL.PrismamovimientodetalleId = rcl.PrismamovimientodetalleId
	}
	entityPrismaCL.PrismatrdospagosId = 0
	if rcl.PrismatrdospagosId > 0 {
		entityPrismaCL.PrismatrdospagosId = rcl.PrismatrdospagosId
	}

	entityPrismaCL.BancoExternalId = rcl.BancoExternalId
	entityPrismaCL.Tiporegistro = rcl.Tiporegistro
	entityPrismaCL.PagosUuid = rcl.PagosUuid
	entityPrismaCL.ExternalmediopagoId = rcl.ExternalmediopagoId
	entityPrismaCL.Nrotarjeta = rcl.Nrotarjeta
	entityPrismaCL.Tipooperacion = rcl.Tipooperacion
	entityPrismaCL.Fechaoperacion = rcl.Fechaoperacion
	entityPrismaCL.Monto = rcl.Monto
	entityPrismaCL.Montofinal = rcl.Montofinal
	entityPrismaCL.Codigoautorizacion = rcl.Codigoautorizacion
	entityPrismaCL.Nroticket = rcl.Nroticket
	entityPrismaCL.SiteID = rcl.SiteID
	entityPrismaCL.ExternalloteId = rcl.ExternalloteId
	entityPrismaCL.Nrocuota = rcl.Nrocuota
	entityPrismaCL.FechaCierre = rcl.FechaCierre
	entityPrismaCL.Nroestablecimiento = rcl.Nroestablecimiento
	entityPrismaCL.ExternalclienteID = rcl.ExternalclienteID
	entityPrismaCL.Nombrearchivolote = rcl.Nombrearchivolote
	entityPrismaCL.Match = rcl.Match
	entityPrismaCL.FechaPago = rcl.FechaPago
	entityPrismaCL.Disputa = rcl.Disputa
	entityPrismaCL.Cantdias = int64(rcl.Cantdias)
	entityPrismaCL.Enobservacion = rcl.Enobservacion

	entityPrismaCL.Valorpresentado = rcl.Valorpresentado
	entityPrismaCL.Diferenciaimporte = rcl.Diferenciaimporte
	entityPrismaCL.Coeficientecalculado = rcl.Coeficientecalculado
	entityPrismaCL.Costototalporcentaje = rcl.Costototalporcentaje
	entityPrismaCL.Importeivaarancel = rcl.Importeivaarancel
	entityPrismaCL.Descripcion = rcl.Descripcion
	entityPrismaCL.Descripcionpresentacion = rcl.Descripcionpresentacion
	entityPrismaCL.DetallemovimientoId = rcl.DetallemovimientoId
	entityPrismaCL.DetallepagoId = rcl.DetallepagoId
	entityPrismaCL.Descripcioncontracargo = rcl.Descripcioncontracargo
	entityPrismaCL.Reversion = rcl.Reversion

	entityPrismaCL.ExtbancoreversionId = rcl.ExtbancoreversionId
	entityPrismaCL.Conciliado = rcl.Conciliado
	entityPrismaCL.Estadomovimiento = rcl.Estadomovimiento
	entityPrismaCL.Descripcionbanco = rcl.Descripcionbanco

	return
}
