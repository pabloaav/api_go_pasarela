package util_test

import (
	"math"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

type TableDriverTestCalculoCupon struct {
	TituloPrueba string
	WantError    string
	Want         int64
	Importe      float64
	Coeficiente  float64
	Impuesto     float64
}

func EstructuraCierreLoteFakeValorCuponValido() (tableDriverTest TableDriverTestCalculoCupon) {
	tableDriverTest = TableDriverTestCalculoCupon{
		TituloPrueba: "obtener un valor de cupon de pago valido",
		WantError:    "",
		Want:         1537411, //1519194,
		Importe:      13274.75,
		Coeficiente:  1.1307,
		Impuesto:     0.21,
	}
	return
}

func TestCalcularCFValdo(t *testing.T) {
	tableDriveTestCalcularCupon := EstructuraCierreLoteFakeValorCuponValido()
	utilService := new(mockservice.MockUtilService)
	t.Run(tableDriveTestCalcularCupon.TituloPrueba, func(t *testing.T) {

		importeCupon := tableDriveTestCalcularCupon.Importe * tableDriveTestCalcularCupon.Coeficiente
		costoFinancieroNeto := importeCupon - tableDriveTestCalcularCupon.Importe
		valorACobrar := costoFinancieroNeto + (costoFinancieroNeto * tableDriveTestCalcularCupon.Impuesto)
		valorCFmasIVA := ToFixed(valorACobrar, 2)
		//inporteMasCfIva := importe + valorCFmasIVA
		//finalValor := s.ToFixed(inporteMasCfIva, 2)
		valorCupon := entities.Monto(ToFixed(tableDriveTestCalcularCupon.Importe+valorCFmasIVA, 4) * 100).Int64()
		want := tableDriveTestCalcularCupon.Want
		utilService.On("CalcularValorCuponService", tableDriveTestCalcularCupon.Importe, tableDriveTestCalcularCupon.Coeficiente, tableDriveTestCalcularCupon.Impuesto).Return(valorCupon)
		got := utilService.CalcularValorCuponService(tableDriveTestCalcularCupon.Importe, tableDriveTestCalcularCupon.Coeficiente, tableDriveTestCalcularCupon.Impuesto)
		assert.Equal(t, got, want)
	})

}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
