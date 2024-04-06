package prisma_test

import (
	"encoding/json"
	"testing"

	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/domains/prisma/prismafake"
	"github.com/stretchr/testify/assert"
)

func TestCheckService(t *testing.T) {

}

/*
	test solicitad token pago con tajeta - validación de datos
*/
func TestSolicitarTokenPagoTarjetaFakeNroTarjeta(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeNroTarjeta()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.Card.CardNumber = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable
			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)
			assert.Equal(t, got.Error(), want)
		})
	}
}
func TestSolicitarTokenPagoTarjetaFakeExpirationMonth(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeExpirationMonth()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.Card.CardExpirationMonth = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)
		})
	}
}

func TestSolicitarTokenPagoTarjetaFakeExpirationYear(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeExpirationYear()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.Card.CardExpirationYear = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)

		})
	}
}

func TestSolicitarTokenPagoTarjetaFakeHolderName(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeHolderName()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.Card.CardHolderName = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)

		})
	}
}

func TestSolicitarTokenPagoTarjetaFakeTypePay(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeTypePay()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.TypePay = prismatransacciones.EnumTipoPagoPrisma(v)
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)

		})
	}
}
func TestSolicitarTokenPagoTarjetaFakeCardEmpty(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFakeCardEmpty()
	want := tableDriverTest.WantTable

	_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

	assert.Equal(t, got.Error(), want)
}

/*
	test solocitar token pago offline - validacion de datos
*/
func TestSolicitarTokenPagoofflineFakeDni(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoOffLineFakeDni()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.DataOffline.Customer.Identification.Number = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)
		})
	}
}

func TestSolicitarTokenPagoofflineFakeNonbre(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoOffLineFakeNonbre()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.DataOffline.Customer.Name = v
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)
		})
	}
}

func TestSolicitarTokenPagoTarjetaFakePagoOfflineEmpty(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoOffLineFakePagoOffline()
	want := tableDriverTest.WantTable

	_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

	assert.Equal(t, got.Error(), want)
}

func TestSolicitarTokenPagoofflineFakeTipoDni(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoOffLineFakeTipoDni()
	for _, v := range tableDriverTest.DataPrueba {
		tableDriverTest.TokenStructura.DataOffline.Customer.Identification.Type = prismatransacciones.EnumTipoDocumento(v)
		t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
			want := tableDriverTest.WantTable

			_, got := service.SolicitarToken(tableDriverTest.TokenStructura)

			assert.Equal(t, got.Error(), want)
		})
	}
}

/*
	test solicitud de token pago con tarjeta - validacion repositorio remoto
*/
func TestSolicitarTokenPagoTarjetaValido(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoTarjetaFaValido()
	t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
		var want prismatransacciones.PagoToken
		json.Unmarshal([]byte(tableDriverTest.WantTable), &want)

		mockRemoteRepositoryPrisma.On("PostSolicitudTokenPago", tableDriverTest.TokenStructura.Card).Return(&want, nil)
		got, err := service.SolicitarToken(tableDriverTest.TokenStructura)

		assert.Equal(t, err, nil)
		assert.Equal(t, got, want)
	})
}

func TestSolicitarTokenPagoOffLineValido(t *testing.T) {
	tableDriverTest := prismafake.EstructuraTokenPagoOffLineValido()
	t.Run(tableDriverTest.TituloPrueba, func(t *testing.T) {
		var want prismatransacciones.OfflineTokenResponse
		json.Unmarshal([]byte(tableDriverTest.WantTable), &want)
		mockRemoteRepositoryPrisma.On("PostSolicitarTokenOffLine", tableDriverTest.TokenStructura.DataOffline).Return(&want, nil)
		got, err := service.SolicitarToken(tableDriverTest.TokenStructura)
		mockRemoteRepositoryPrisma.AssertExpectations(t)
		assert.Equal(t, err, nil)
		assert.Equal(t, got, want)
	})
}
