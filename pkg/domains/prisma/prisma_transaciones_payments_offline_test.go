package prisma_test

import (
	"encoding/json"
	"testing"

	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/prisma/prismafake"
	"github.com/stretchr/testify/assert"
)

/*
	test validando el ingreso de tipo de documento erroneos
*/
func TestPaymentsPagoOffLineFakeTipoDocumento(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeTipoDocumento()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Customer.Identification.Type = prismadtos.EnumTipoDocumento(v)
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de número documento erroneos
*/
func TestPaymentsPagoOffLineFakeDocumento(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeDocumento()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Customer.Identification.Number = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)

		})
	}
}

/*
	test validando el ingreso de name erroneos
*/
func TestPaymentsPagoOffLineFakeCustomerName(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeCustomerName()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Customer.Name = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de Site transaction Id erroneos
*/
func TestPaymentsPagoOffLineFakeSiteTransactionId(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeSiteTransactionId()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.SiteTransactionID = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de token erroneos
*/
func TestPaymentsPagoOffLineFakeToken(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeToken()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Token = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de Amount erroneos
*/
func TestPaymentsPagoOffLineFakeAmount(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeAmount()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaInt {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Amount = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de email erroneos
*/
func TestPaymentsPagoOffLineFakeEmail(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeEmail()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Email = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de tipo moneda "currency" erroneos
*/
func TestPaymentsPagoOffLineFakeCurrency(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeCurrency()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Currency = prismadtos.EnumTipoMoneda(v)
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de cantidad de dia entre en 1º y 2º vencimiento, erroneos
*/
func TestPaymentsPagoOffLineFakeCodP3(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeCodP3()
	//service := IniciarMock()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.CodP3 = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)

		})
	}
}

/*
	test validando el ingreso de cantidad de Días después del 1º vencimiento y hasta que el cliente pueda abonar erroneos
*/
func TestPaymentsPagoOffLineFakeCodP4(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeCodP4()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.CodP4 = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test validando el ingreso de número de Cliente erroneo
*/
func TestPaymentsPagoOffLineFakeCliente(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeCliente()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Client = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)

		})
	}
}

/*
	test validando el ingreso de Surcharge erroneo
*/
func TestPaymentsPagoOffLineFakeSurcharge(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeSurcharge()
	for _, v := range TableDriverTest.DataPruebaInt {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.Surcharge = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)

		})
	}
}

/*
	test validando el ingreso de PaymentMode erroneo
*/
func TestPaymentsPagoOffLineFakePaymentMode(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakePaymentMode()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoOffline.PaymentMode = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)

		})
	}
}

/*
	Test validando paga offline valido
*/
func TestPaymentsPagoOffLineValido(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoOffLineFakeValido()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		var want prismadtos.PaymentsOfflineResponse
		json.Unmarshal([]byte(TableDriverTest.WantTable), &want)
		mockRemoteRepositoryPrisma.On("PostEjecutarPagoOffLine", TableDriverTest.PaymentStructura.PagoOffline).Return(&want, nil)
		got, err := service.Payments(TableDriverTest.PaymentStructura)
		assert.Equal(t, err, nil)
		assert.Equal(t, got, want)
	})
}
