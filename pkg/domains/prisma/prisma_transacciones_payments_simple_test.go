package prisma_test

import (
	"encoding/json"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/prisma/prismafake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

/*
	inicializa todos los mocks
*/
var (
	mockRemoteRepositoryPrisma = new(mockrepository.MockRemoteRepositoryPrisma)
	mockRepositoryPrisma       = new(mockrepository.MockRepositoryPrisma)
	mockCommonds               = new(mockservice.MockCommonsService)
	//mockServiceAdministracion  = new(mockservice.MockAdministracionService)
	service = prisma.NewService(mockRemoteRepositoryPrisma, mockRepositoryPrisma, mockCommonds)
)

/*
	Test validando datos incorrectos  SiteTransactionId
*/
func TestPaymentsPagoTarjetaSiteTransaccionId(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeSiteTransactionId()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.SiteTransactionID = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos token
*/
func TestPaymentsPagoTarjetaToken(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeToken()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Token = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})

	}
}

/*
	Test validando datos incorrectos bin
*/
func TestPaymentsPagoTarjetaBin(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeBin()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Bin = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos amount
*/
func TestPaymentsPagoTarjetaAmount(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeAmount()
	for _, v := range TableDriverTest.DataPruebaInt {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Amount = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos currency
*/
func TestPaymentsPagoTarjetaCurrency(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeCurrency()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Currency = prismadtos.EnumTipoMoneda(v)
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos installments
*/
func TestPaymentsPagoTarjetaInstallments(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeInstallments()
	for _, v := range TableDriverTest.DataPruebaInt {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Installments = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos type
*/
func TestPaymentsPagoTarjetaPaymentType(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakePaymentType()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.PaymentType = prismadtos.EnumPaymentType(v)
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos establishments name
*/
func TestPaymentsPagoTarjetaEstablishmentName(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeEstablishmentName()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.EstablishmentName = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando datos incorrectos email
*/
func TestPaymentsPagoTarjetaEmail(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeEmail()
	for _, v := range TableDriverTest.DataPruebaString {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			TableDriverTest.PaymentStructura.PagoSimple.Customeremail.Email = v
			_, got := service.Payments(TableDriverTest.PaymentStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	Test validando respuesta de un pago con tarjeta valida
*/
func TestPaymentsPagoTarjetaValido(t *testing.T) {
	TableDriverTest := prismafake.EstructuraPaymentsPagoTarjetaFakeValido()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		var want prismadtos.PaymentsSimpleResponse
		json.Unmarshal([]byte(TableDriverTest.WantTable), &want)
		mockRemoteRepositoryPrisma.On("PostEjecutarPago", TableDriverTest.PaymentStructura.PagoSimple).Return(&want, nil)
		got, err := service.Payments(TableDriverTest.PaymentStructura)
		assert.Equal(t, err, nil)
		assert.Equal(t, got, want)
	})
}
