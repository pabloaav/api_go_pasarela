package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseOperacionesContracargo struct {
	Cuenta ResponseCuentaContraRecargo `json:"cuenta"`
}

type ResponseCuentaContraRecargo struct {
	Id         uint                 `json:"id"`
	ClientesID int64                `json:"clientes_id"`
	RubrosID   uint                 `json:"rubros_id"`
	Cuenta     string               `json:"cuenta"`
	Cbu        string               `json:"cbu"`
	Cvu        string               `json:"cvu"`
	PagoTipo   []ResponsePagotipoCC `joson:"pago_tipo"`
}

type ResponsePagotipoCC struct {
	Id       uint             `json:"id"`
	Pagotipo string           `json:"pagotipo"`
	Pagos    []ResponsePagoCC `json:"pagos"`
	// BackUrlSuccess           string
	// BackUrlPending           string
	// BackUrlRejected          string
	// BackUrlNotificacionPagos string
	// IncludedChannels         []CanalesPago
	// IncludedInstallments     []CuotasPago
}

type ResponsePagoCC struct {
	Id                  uint                  `json:"id"`
	Fecha               time.Time             `json:"fecha"`
	PagostipoID         int64                 `json:"pagostipo_id"`
	ExternalReference   string                `json:"external_reference"`
	PayerName           string                `json:"payer_name"`
	Estado              string                `json:"estado"`
	NombreEstado        string                `json:"nombre_estado"`
	Amount              entities.Monto        `json:"amount"`
	FechaPago           time.Time             `json:"fecha_pago"`
	Channel             string                `json:"channel"`
	NombreChannel       string                `json:"nombre_channel"`
	UltimoPagoIntentoId uint64                `json:"ultimo_pago_intento_id"`
	TransferenciaId     uint64                `json:"transferencia_id"`
	ReferenciaBancaria  string                `json:"referencia_bancaria"`
	PagoIntento         ResponsePagoIntentoCC `json:"pago_intento"`
}

type ResponsePagoIntentoCC struct {
	Id                   uint           `json:"id"`
	MediopagosId         uint           `json:"mediopagos_id"`
	InstallmentdetailsId uint           `json:"installmentdetails_id"`
	ExternalId           string         `json:"external_id"`
	PaidAt               time.Time      `json:"paid_at"`
	ReortAt              time.Time      `json:"report_at"`
	IsAvailable          bool           `json:"is_available"`
	Amount               entities.Monto `json:"amount"`
	Valorcupon           entities.Monto `json:"valorcupon"`
	StateComent          string         `json:"state_comment"`
	Barcode              string         `json:"barcode"`
	BarcodeUrl           string         `json:"barcode_url"`
	AvailableAt          time.Time      `json:"available_at"`
	RevertedAt           time.Time      `json:"reverted_at"`
	HolderName           string         `json:"holder_name"`
	HolderEmail          string         `json:"holder_email"`
	HolderType           string         `json:"holder_type"`
	HolderNumber         string         `json:"holder_number"`
	HolderCbu            string         `json:"holder_cbu"`
	TicketNumber         string         `json:"ticket_number"`
	AuthorizationCode    string         `json:"authorization_code"`
	CardLastFourDigits   string         `json:"card_last_four_digits"`
	TransactionId        string         `json:"transaction_id"`
	SiteId               string         `json:"site_id"`
	CierreLote           ResponseCLCC   `json:"cierre_lote"`
}

type ResponseCLCC struct {
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

	Reversion                  bool                     `json:"reversion"`
	DetallemovimientoId        int64                    `json:"detallemovimineto_id"`
	DetallepagoId              int64                    `json:"detallepago_id"`
	Descripcioncontracargo     string                   `json:"descripcioncontracargo"`
	ExtbancoreversionId        int64                    `json:"extbancoreversion_id"`
	Conciliado                 bool                     `json:"conciliacion"`
	Estadomovimiento           bool                     `json:"estadomovimineto"`
	Descripcionbanco           string                   `json:"descripcionbanco"`
}
