package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

// type ResponsePagosEstado struct {
// 	Uuid              string    `json:"uuid"`
// 	ExternalReference string    `json:"external_reference"`
// 	Metadata          string    `json:"metadata"`
// 	Fecha             time.Time `json:"fecha"`
// 	Estado            string    `json:"estado"`
// }

type ResponseSolicitudPago struct {
	SolicitudpagoId   uint                  `json:"solicitudpago_id"`
	Uuid              string                `json:"uuid"`
	PaymentType       string                `json:"payment_type"`
	PagoEstado        string                `json:"pago_estado"`
	Estado            string                `json:"estado"`
	CreatedAt         string                `json:"created_at"`
	ExternalReference string                `json:"external_reference"`
	PayerName         string                `json:"payer_name"`
	PayerEmail        string                `json:"payer_email"`
	Description       string                `json:"description"`
	FirstDueDate      string                `json:"first_due_date"`
	FirstTotal        float64               `json:"first_total"`
	PagoIntento       []ResponsePagoIntento `json:"pago_Intento"`
}

type ResponsePagoIntento struct {
	PagoIntentoId uint    `json:"pagointento_id"`
	Channel       string  `json:"channel"`
	PaidAt        string  `json:"paid_at"`
	IsAvailable   bool    `json:"is_available"`
	Amount        float64 `json:"amount"`
	GrossFee      float64 `json:""`
	NetFee        float64 `json:"net_fee"`
	FeeIva        float64 `json:"fee_iva"`
	NetAmount     float64 `json:"net_amount"`
}

func (rs *ResponseSolicitudPago) SolicitudEntityToDtos(solicitudEntity entities.Pago) {
	rs.SolicitudpagoId = solicitudEntity.ID
	rs.Uuid = solicitudEntity.Uuid
	rs.PaymentType = solicitudEntity.PagosTipo.Pagotipo
	rs.PagoEstado = string(solicitudEntity.PagoEstados.Estado)
	rs.Estado = solicitudEntity.PagoEstados.Nombre
	rs.CreatedAt = solicitudEntity.Model.CreatedAt.String()
	rs.ExternalReference = solicitudEntity.ExternalReference
	rs.PayerName = solicitudEntity.PayerName
	rs.PayerEmail = solicitudEntity.PayerEmail
	rs.Description = solicitudEntity.Description
	rs.FirstDueDate = solicitudEntity.FirstDueDate.String()
	rs.FirstTotal = commons.ToFixedTool(solicitudEntity.FirstTotal.Float64(), 2)
	for _, value := range solicitudEntity.PagoIntentos {
		var pagoIntentoTemporal ResponsePagoIntento
		pagoIntentoTemporal.PagoIntentoId = value.ID
		pagoIntentoTemporal.Channel = value.Mediopagos.Channel.Channel
		pagoIntentoTemporal.PaidAt = value.PaidAt.String()
		pagoIntentoTemporal.Amount = commons.ToFixedTool(value.Amount.Float64(), 2)
		if len(value.Movimientos) == 0 {
			pagoIntentoTemporal.IsAvailable = false
			pagoIntentoTemporal.GrossFee = 0
			pagoIntentoTemporal.NetFee = 0
			pagoIntentoTemporal.FeeIva = 0
			pagoIntentoTemporal.NetAmount = 0
		}
		if len(value.Movimientos) > 0 {
			pagoIntentoTemporal.IsAvailable = true
			if len(value.Movimientos[0].Movimientocomisions) > 0 && len(value.Movimientos[0].Movimientoimpuestos) > 0 {
				pagoIntentoTemporal.NetFee = commons.ToFixedTool(value.Movimientos[0].Movimientocomisions[0].Monto.Float64()+value.Movimientos[0].Movimientocomisions[0].Montoproveedor.Float64(), 2)
				pagoIntentoTemporal.FeeIva = commons.ToFixedTool(value.Movimientos[0].Movimientoimpuestos[0].Monto.Float64()+value.Movimientos[0].Movimientoimpuestos[0].Montoproveedor.Float64(), 2)
				pagoIntentoTemporal.GrossFee = commons.ToFixedTool(pagoIntentoTemporal.NetFee+pagoIntentoTemporal.FeeIva, 2)
			}
			pagoIntentoTemporal.NetAmount = value.Movimientos[0].Monto.Float64()
		}
		rs.PagoIntento = append(rs.PagoIntento, pagoIntentoTemporal)
	}
}
