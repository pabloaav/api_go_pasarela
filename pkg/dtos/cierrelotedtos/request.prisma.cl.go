package cierrelotedtos

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestPrismaCL struct {
	Id                         int64 `json:"id"`
	PagoestadoexternosId       int64 `json:"pagoestadoexternos_id"`
	ChannelarancelesId         int64 `json:"channelaranceles_id"`
	ImpuestosId                int64 `json:"impuestos_id"`
	PrismamovimientodetallesId int64 `json:"prismamovimientodetalles_id"`
	// PrismamovimientodetalleId  int64                      `json:"prismamovimientodetalle_id"`
	PrismatrdospagosId  int64                      `json:"prismatrdospagos_id"`
	BancoExternalId     int64                      `json:"banco_external_id"`
	Tiporegistro        string                     `json:"tiporegistro"`
	PagosUuid           string                     `json:"pagos_uuid"`
	ExternalmediopagoId int64                      `json:"externalmediopago"`
	Nrotarjeta          string                     `json:"nrotarjeta"`
	Tipooperacion       entities.EnumTipoOperacion `json:"tipooperacion"`
	Fechaoperacion      time.Time                  `json:"fechaoperacion"`
	Monto               entities.Monto             `json:"monto"`
	Montofinal          entities.Monto             `json:"montofinal"`
	Codigoautorizacion  string                     `json:"codigoautorizacion"`
	Nroticket           int64                      `json:"nroticket"`
	SiteID              int64                      `json:"site_id"`
	ExternalloteId      int64                      `json:"externallote_id"`
	Nrocuota            int64                      `json:"nrocuota"`
	FechaCierre         time.Time                  `json:"fecha_cierre"`
	Nroestablecimiento  int64                      `json:"nroestablecimiento"`
	ExternalclienteID   string                     `json:"externalcliente_id"`
	Nombrearchivolote   string                     `json:"nombrearchivolote"`
	Match               int                        `json:"match"`
	FechaPago           time.Time                  `json:"fecha_pago"`
	Disputa             bool                       `json:"disputa"`
	Cantdias            int                        `json:"cantdias"`
	Enobservacion       bool                       `json:"enobservacion"`
}

func (rcl *RequestPrismaCL) RequestPrismaCLToEntity() (entityPrismaCL entities.Prismacierrelote) {
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
	// entityPrismaCL.PrismamovimientodetalleId = 0
	// if rcl.PrismamovimientodetalleId > 0 {
	// 	entityPrismaCL.PrismamovimientodetalleId = rcl.PrismamovimientodetalleId
	// }
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
	return
}

func (rcl *RequestPrismaCL) IsValid() (erro error) {

	if rcl.Id < 1 {
		return fmt.Errorf("id incorrecto")
	}
	if rcl.FechaCierre.IsZero() {
		return fmt.Errorf("la Fecha de Cierre es incorrecta")
	}
	if reflect.TypeOf(rcl.Disputa).Name() != "bool" {
		return fmt.Errorf("valor disputa incorrecto")
	}
	if reflect.TypeOf(rcl.Enobservacion).Name() != "bool" {
		return fmt.Errorf("valor observacion incorrecto")
	}
	if rcl.Monto <= 0 {
		return fmt.Errorf("monto recibido incorrecto")
	}

	return
}

func (rcl *RequestPrismaCL) ValidarCamposEdit(entityCL entities.Prismacierrelote) (filtroEdit FiltroEditarCLHerramienta) {

	if rcl.Monto == entityCL.Monto {
		filtroEdit.Monto = false
	} else {
		filtroEdit.Monto = true
	}

	diferencia := rcl.Monto - entityCL.Monto
	if diferencia > 2 || diferencia < -2 {
		filtroEdit.Monto = false
	} else {
		filtroEdit.Monto = true
	}

	if rcl.FechaCierre.Equal(entityCL.FechaCierre) {
		filtroEdit.Fecha = false
	} else {
		filtroEdit.Fecha = true
	}

	if rcl.Disputa == entityCL.Disputa {
		filtroEdit.Disputa = false
	} else {
		filtroEdit.Disputa = true
	}

	if rcl.Enobservacion == entityCL.Enobservacion {
		filtroEdit.Enobservacion = false
	} else {
		filtroEdit.Enobservacion = true
	}
	return
}
