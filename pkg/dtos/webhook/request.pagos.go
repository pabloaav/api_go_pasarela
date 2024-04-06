package webhook

type RequestWebhook struct {
	DiasPago         int64  `json:"dias_pago"`
	PagosNotificado  bool   `json:"pagos_notificado"`
	EstadoFinalPagos bool   `json:"estado_final_pagos"`
	CuentaId         uint64 `json:"cuenta_id"`
}
