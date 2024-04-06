package ribcradtos_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/stretchr/testify/assert"
)

type TableDatosFondoDto struct {
	Nombre  string
	Erro    error
	Request []ribcradtos.RiDatosFondosRequest
}

func _inicializarDatosFondosReques() (table []TableDatosFondoDto) {

	requestFechasInicioInvalidas := []ribcradtos.RiDatosFondosRequest{
		//Fecha inicio y fin cero
		{},
		//Fecha inicio cero
		{
			FechaFin: time.Now(),
		},
		//Fecha inicio superior a fecha fin
		{
			FechaInicio: time.Now().Add(time.Hour * 30),
			FechaFin:    time.Now(),
		},
	}

	requestFechasFinInvalidas := []ribcradtos.RiDatosFondosRequest{
		//Fecha fin es cero
		{
			FechaInicio: time.Now(),
		},
	}
	requestRutasInvalidas := []ribcradtos.RiDatosFondosRequest{
		//ruta vacia
		{
			FechaInicio: time.Now(),
			FechaFin:    time.Now().Add(time.Hour * 24),
		},
		{
			FechaInicio: time.Now(),
			FechaFin:    time.Now().Add(time.Hour * 24),
			Ruta:        "  ",
		},
	}

	table = []TableDatosFondoDto{
		{"Debe retornar un error si la fecha de inicio invalida", fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA), requestFechasInicioInvalidas},
		{"Debe retornar un error si la fecha de fin es invalida", fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA), requestFechasFinInvalidas},
		{"Debe retornar un error si la ruta es invalida", fmt.Errorf(administraciondtos.ERROR_RUTA_INVALIDA), requestRutasInvalidas},
	}
	return
}

func TestRIDatosFondosRequest(t *testing.T) {

	table := _inicializarDatosFondosReques()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for _, r := range v.Request {

				err := r.IsValid()

				if err != nil {
					assert.Equal(t, v.Erro.Error(), err.Error())
				}
				if err == nil {
					t.Fatalf("Todos los testes deben fallar")
				}

			}

		})
	}
}
