package prismadtos

type StructPayments struct {
	PagoSimple  PaymentsSimpleRequest  `json:"pago_simple,omitempty"`
	PagoOffline PaymentsOfflineRequest `json:"pago_offline,omitempty"`
	TypePay     EnumTipoPagoPrisma     `json:"type_pay"`
}
