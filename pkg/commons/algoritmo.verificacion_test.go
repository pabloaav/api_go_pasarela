package commons_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/commons/commonsfake"
	"github.com/stretchr/testify/assert"
)

var (
	service = commons.NewAlgoritmoVerificacion()
)

func TestValidarTarjeta(t *testing.T) {
	TableDriverTest := commonsfake.EstructuraValidarTarjeta()
	for _, test := range TableDriverTest {
		t.Run(test.TituloPrueba, func(t *testing.T) {
			want := test.WantTable
			logs.Info(test.TituloPrueba)
			got := service.ChequearTarjeta(test.Tarjeta)
			assert.Equal(t, got, want)
		})

	}

}

func TestDiferenciaUint(t *testing.T) {
	TableDriverTest := commonsfake.EstructuraValidarTarjeta()
	for _, test := range TableDriverTest {
		t.Run(test.TituloPrueba, func(t *testing.T) {
			want := test.WantTable
			logs.Info(test.TituloPrueba)
			got := commons.DifferenceInteger([]uint64{574,575,576}, []uint64{576})
			assert.Equal(t, got, want)
		})

	}

}
