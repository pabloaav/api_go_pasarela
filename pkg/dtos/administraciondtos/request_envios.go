package administraciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestEnvios struct {
	Id          uint
	Cobranzas   bool
	Rendiciones bool
	Reversiones bool
	Batch       bool
	BatchPagos  bool `json:"batch_pagos"`
	ClientesId  uint `json:"clientes_id"`
}

func (r *RequestEnvios) Validate() (erro error) {
	if r.Id == 0 {
		return fmt.Errorf("error: parametro id cero")
	}

	if r.ClientesId == 0 {
		return fmt.Errorf("error: parametro clientes_id cero")
	}
	return
}

func (r *RequestEnvios) ToEnvio() (envio entities.Envio) {
	envio.ID = r.Id
	envio.Cobranzas = r.Cobranzas
	envio.Rendiciones = r.Rendiciones
	envio.Reversiones = r.Reversiones
	envio.Batch = r.Batch
	envio.BatchPagos = r.BatchPagos
	envio.ClientesId = r.ClientesId
	return
}
