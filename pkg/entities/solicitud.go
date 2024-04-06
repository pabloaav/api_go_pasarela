package entities

import "gorm.io/gorm"

type Solicitud struct {
	gorm.Model
	Impuestoivaid        string
	Impuestoiibbid       string
	Cliente              string
	Cuit                 string
	Razonsocial          string
	Nombrefantasia       string
	Email                string
	Personeria           string
	Retiroautomatico     uint
	Reportebatch         uint
	Nombrereporte        string
	Cuenta               string
	Cbu                  string
	Cvu                  string
	Apikey               string
	Diasretiroautomatico uint
	Pagotipo             string
	Urlsuccess           string
	Urlpending           string
	Urlrejected          string
	Urlnotificacionpagos string
	Canalpago            string
	Cuotas               string
	Telefono             string
}
