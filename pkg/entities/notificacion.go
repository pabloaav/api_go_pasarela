package entities

import (
	"errors"

	"gorm.io/gorm"
)

type Notificacione struct {
	gorm.Model
	UserId      uint64               `json:"user_id"`
	Tipo        EnumTipoNotificacion `json:"tipo"`
	Descripcion string               `json:"descripcion"`
}

type EnumTipoNotificacion string

const (
	NotificacionTransferenciaAutomatica EnumTipoNotificacion = "TransferenciaAutomatica"
	NotificacionTransferencia           EnumTipoNotificacion = "Transferencia"
	NotificacionCierreLote              EnumTipoNotificacion = "CierreLote"
	NotificacionPagoExpirado            EnumTipoNotificacion = "PagoExpirado"
	NotificacionConfiguraciones         EnumTipoNotificacion = "Configuraciones"
	NotificacionSolicitudCuenta         EnumTipoNotificacion = "SolicitudCuenta"
	NotivicacionEnvioEmail              EnumTipoNotificacion = "EnvioEmail"
	NotificacionProcesoMx               EnumTipoNotificacion = "ProcesoMovimientosMx"
	NotificacionProcesoPx               EnumTipoNotificacion = "ProcesoPagosPx"
	NotificacionConciliacionCLMx        EnumTipoNotificacion = "ConciliacionClMx"
	NotificacionConciliacionCLPx        EnumTipoNotificacion = "ConciliacionClPx"
	NotificacionConciliacionBancoCL     EnumTipoNotificacion = "ConciliacionBancoCl"
	NotificacionWebhook                 EnumTipoNotificacion = "Webhook"
	NotificacionSendEmailCsv            EnumTipoNotificacion = "NotificacionSendEmailCsv"
	NotificacionArchivoBatchCliente     EnumTipoNotificacion = "NotificacionArchivoBatchCliente"
	NotificacionComisionConMaximo       EnumTipoNotificacion = "NotificacionComisionConMaximo"
	NotificacionPagoOfflineExpirado     EnumTipoNotificacion = "NotificacionPagoOfflineExpirado"
)

func (e EnumTipoNotificacion) IsValid() error {
	switch e {
	case NotificacionTransferencia, NotificacionCierreLote, NotificacionConciliacionBancoCL, NotificacionConciliacionCLMx, NotificacionConciliacionCLPx:
		return nil
	}
	return errors.New("tipo EnumTipoNotificacion con formato inválido")
}
