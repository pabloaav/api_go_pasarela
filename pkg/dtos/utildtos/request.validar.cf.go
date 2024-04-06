package utildtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type RequestValidarCF struct {
	Cupon        entities.Monto
	ArancelMonto float64
	Tna          float64
	Cuotas       float64
	Dias         float64
}
