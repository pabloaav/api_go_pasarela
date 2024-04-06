package administraciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

// type SolicitudCuentaRequest struct {
// 	ImpuestoIvaId        string `json:"impuesto_iva_id"`
// 	ImpuestoIibbId       string `json:"impuesto_iibb_id"`
// 	Cliente              string `json:"cliente"`        // requerido
// 	Cuit                 string `json:"cuit"`           // requerido
// 	Razonsocial          string `json:"razonsocial"`    // requerido
// 	Nombrefantasia       string `json:"nombrefantasia"` // requerido
// 	Email                string `json:"email"`          // requerido
// 	Personeria           string `json:"personeria"`
// 	RetiroAutomatico     uint   `json:"retiro_automatico"` // requerido
// 	Reportebatch         uint   `json:"reporte_batch"`
// 	NombreReporte        string `json:"nombre_reporte"`
// 	Cuenta               string `json:"cuenta"` // requerido
// 	Cbu                  string `json:"cbu"`    // requerido
// 	Cvu                  string `json:"cvu"`
// 	Apikey               string `json:"apikey"`
// 	DiasRetiroAutomatico uint   `json:"dias_retiro_automatico"`
// 	Pagotipo             string `json:"pagotipo"` // requerido
// 	UrlSuccess           string `json:"url_success"`
// 	UrlPending           string `json:"url_pending"`
// 	UrlRejected          string `json:"url_rejected"`
// 	UrlNotificacionPagos string `json:"url_notificacion_pagos"`
// 	CanalPago            string `json:"canal_pago"` // requerido
// 	Cuotas               string `json:"cuotas"`
// }

type SolicitudCuentaRequest struct {
	Nombre      string `json:"nombre"`
	Apellido    string `json:"apellido"`
	Cuit        string `json:"cuit"`
	Razonsocial string `json:"razonsocial"`
	Email       string `json:"email"`
	Telefono    string `json:"telefono"`
}

func (solicitud SolicitudCuentaRequest) IsValid() error {

	stringsRequired := map[string]string{
		"Nombre":       solicitud.Nombre,
		"Apellido":     solicitud.Apellido,
		"Cuit":         solicitud.Cuit,
		"Razon Social": solicitud.Razonsocial,
		"Telefono":     solicitud.Telefono,
	}

	result, field := commons.SomeStringIsEmpty(stringsRequired)

	if result {
		return fmt.Errorf("el campo %s es requerido", field)
	}

	if !commons.IsEmailValid(solicitud.Email) {
		return fmt.Errorf("el campo email no es valido")
	}

	err := commons.EsCuilValido(solicitud.Cuit)
	if err != nil {
		return err
	}

	return nil
}

func (request SolicitudCuentaRequest) ToSolicitudEntity() (s entities.Solicitud) {
	// s.Impuestoivaid = request.ImpuestoIvaId
	// s.Impuestoiibbid = request.ImpuestoIibbId
	// s.Cliente = request.Cliente
	// s.Cuit = request.Cuit
	// s.Razonsocial = request.Razonsocial
	// s.Nombrefantasia = request.Nombrefantasia
	// s.Email = request.Email
	// s.Personeria = "F"
	// s.Retiroautomatico = request.RetiroAutomatico
	// s.Reportebatch = request.Reportebatch
	// s.Nombrereporte = request.NombreReporte
	// s.Cuenta = request.Cuenta
	// s.Cbu = request.Cbu
	// s.Cvu = request.Cvu
	// s.Apikey = request.Apikey
	// s.Diasretiroautomatico = request.DiasRetiroAutomatico
	// s.Pagotipo = request.Pagotipo
	// s.Urlsuccess = request.UrlSuccess
	// s.Urlpending = request.UrlPending
	// s.Urlrejected = request.UrlRejected
	// s.Urlnotificacionpagos = request.UrlNotificacionPagos
	// s.Canalpago = request.CanalPago
	// s.Cuotas = request.Cuotas

	s.Cliente = request.Nombre + " " + request.Apellido
	s.Cuit = request.Cuit
	s.Razonsocial = request.Razonsocial
	s.Email = request.Email
	s.Telefono = request.Telefono
	s.Personeria = "F"

	return
}
