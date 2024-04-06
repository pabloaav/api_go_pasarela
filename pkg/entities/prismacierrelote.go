package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Prismacierrelote struct {
	gorm.Model
	PagoestadoexternosId       int64                    `json:"pagoestadoexternos_id"`
	ChannelarancelesId         int64                    `json:"channelaranceles_id"`
	ImpuestosId                int64                    `json:"impuestos_id"`
	PrismamovimientodetallesId int64                    `gorm:"default:(null),foreignkey:prismamovimientodetalles_id"`
	PrismamovimientodetalleId  int64                    `json:"prismamovimientodetalle_id"`
	PrismatrdospagosId         int64                    `gorm:"default:(null),foreignkey:prismatrdospagos_id"`
	PrismapagotrdoId           int64                    `json:"prismapagotrdo_id"`
	BancoExternalId            int64                    `json:"banco_external_id"`
	Tiporegistro               string                   `json:"tiporegistro"`
	PagosUuid                  string                   `json:"pagos_uuid"`
	ExternalmediopagoId        int64                    `json:"externalmediopago"`
	Nrotarjeta                 string                   `json:"nrotarjeta"`
	Tipooperacion              EnumTipoOperacion        `json:"tipooperacion"`
	Fechaoperacion             time.Time                `json:"fechaoperacion"`
	Monto                      Monto                    `json:"monto"`
	Montofinal                 Monto                    `json:"montofinal"`
	Valorpresentado            Monto                    `json:"valorpresentado"`
	Diferenciaimporte          Monto                    `json:"difernciaimporte"`
	Coeficientecalculado       float64                  `json:"coeficientecalculado"`
	Costototalporcentaje       float64                  `json:"costototalporcentaje"`
	Importeivaarancel          float64                  `json:"importeivaarancel"`
	Descripcion                string                   `json:"descripcion"`
	Codigoautorizacion         string                   `json:"codigoautorizacion"`
	Nroticket                  int64                    `json:"nroticket"`
	SiteID                     int64                    `json:"site_id"`
	ExternalloteId             int64                    `json:"externallote_id"`
	Nrocuota                   int64                    `json:"nrocuota"`
	FechaCierre                time.Time                `json:"fecha_cierre"`
	Nroestablecimiento         int64                    `json:"nroestablecimiento"`
	ExternalclienteID          string                   `json:"externalcliente_id"`
	Nombrearchivolote          string                   `json:"nombrearchivolote"`
	Match                      int                      `json:"match"`
	FechaPago                  time.Time                `json:"fecha_pago"`
	Disputa                    bool                     `json:"disputa"`
	Cantdias                   int64                    `json:"cantdias"`
	Enobservacion              bool                     `josn:"enobservacion"`
	Descripcionpresentacion    string                   `json:"descripcionpresentacion"`
	Reversion                  bool                     `json:"reversion"`
	DetallemovimientoId        int64                    `json:"detallemovimineto_id"`
	DetallepagoId              int64                    `json:"detallepago_id"`
	Descripcioncontracargo     string                   `json:"descripcioncontracargo"`
	ExtbancoreversionId        int64                    `json:"extbancoreversion_id"`
	Conciliado                 bool                     `json:"conciliacion"`
	Estadomovimiento           bool                     `json:"estadomovimineto"`
	Descripcionbanco           string                   `json:"descripcionbanco"`
	Prismamovimientodetalle    *Prismamovimientodetalle `json:"prismamovimientodetalles_id" gorm:"foreignKey:PrismamovimientodetallesId"`
	Prismatrdospagos           *Prismatrdospago         `json:"prismatrdospagos_id" gorm:"foreignKey:PrismatrdospagosId"`
	Channelarancel             *Channelarancele         `json:"channelarancel" gorm:"foreignKey:channelaranceles_id"`
	//MovimientosId        uint64    `json:"movimientos_id,omitempty"`

	//PagoEstado           Pagoestado `json:"pago_estado" gorm:"foreignKey:PagoestadosId"`
}

type EnumTipoOperacion string

const (
	C EnumTipoOperacion = "C"
	A EnumTipoOperacion = "A"
	D EnumTipoOperacion = "D"
)

func (e EnumTipoOperacion) IsValid() error {
	switch e {
	case C, A, D:
		return nil
	}
	return errors.New("el tipo de operación es inválido")
}
