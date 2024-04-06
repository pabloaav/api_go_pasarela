package checkout_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

var (
	comprador = linkdebin.CompradorCreateDebinLink{
		Cuit: requestDebinValido.HolderCuit,
		Cuenta: linkdebin.CuentaLink{
			Cbu:      requestDebinValido.Cbu,
			AliasCbu: requestDebinValido.Alias,
		},
	}

	vendedor = linkdebin.VendedorCreateLink{
		Cbu: cuentaValida.Cbu,
	}

	debin = linkdebin.DebinCreateLink{
		ComprobanteId:         pagoCreated.ExternalReference,
		EsCuentaPropia:        requestDebinValido.EsCuentaPropia,
		Concepto:              "VAR",
		TiempoExpiracion:      4320, //request.TiempoExpiracion,
		Descripcion:           "aaaa-bbb-cccc-ddddd",
		Importe:               requestDebinValido.Importe,
		Moneda:                "ARS",
		Recurrente:            requestDebinValido.Recurrente,
		DescripcionPrestacion: pagoCreated.Description,
	}

	debinRequest = linkdebin.RequestDebinCreateLink{
		Comprador: comprador,
		Vendedor:  vendedor,
		Debin:     debin,
	}

	debinResult = linkdebin.ResponseDebinCreateLink{
		Id:              "1",
		FechaOperacion:  time.Date(2021, 7, 8, 10, 00, 0, 0, time.Local),
		Estado:          "INICIADO",
		FechaExpiracion: time.Date(2021, 8, 9, 10, 00, 0, 0, time.Local),
	}
)

func TestDebinPayments(t *testing.T) {
	assertions := assert.New(t)
	serv := mockservice.MockApiLinkService{}
	util := mockservice.MockUtilService{}
	payment := checkout.NewDebinPayment(&serv, &util)

	t.Run("Cuando la api devuelve un error, debemos devolverlo", func(t *testing.T) {
		debinReqNoValido := debinRequest
		debinReqNoValido.Debin.Importe = 0
		reqNoValido := requestDebinValido
		reqNoValido.Importe = 0
		var installmentsDetails *dtos.InstallmentDetailsResponse
		serv.On("CreateDebinApiLinkService", reqNoValido.Uuid, debinReqNoValido).Return(&linkdebin.ResponseDebinCreateLink{}, fmt.Errorf("Error en proceso de pago"))

		_, err := payment.CreateResultado(&reqNoValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", installmentsDetails)

		assertions.EqualError(err, "Error en proceso de pago")
	})

	t.Run("Cuando el proceso es exitoso devuelvo un pagointento con el estado, el id de la operacion y la fecha paidAt distinta a 0", func(t *testing.T) {

		serv.On("CreateDebinApiLinkService", requestDebinValido.Uuid, debinRequest).Return(&debinResult, nil)

		res, _ := payment.CreateResultado(&requestDebinValido, &pagoCreated, &cuentaValida, "aaaa-bbb-cccc-ddddd", &dtos.InstallmentDetailsResponse{})

		assertions.NotZero(res.StateComment)
		assertions.NotZero(res.ExternalID)
		assertions.NotZero(res.PaidAt)
	})
}
