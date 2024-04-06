package checkout_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	requestNoValido = dtos.ResultadoRequest{}

	requestDebinValido = dtos.ResultadoRequest{
		Channel:           "debin",
		Uuid:              "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		Cbu:               "0133212221154658845",
		ConceptoAbreviado: "VAR",
		Moneda:            "ARS",
		HolderName:        "Fernando Castro",
		HolderEmail:       "fernando.castro@telco.com.ar",
		HolderCuit:        "23328803259",
		Importe:           105050,
		PaymentMethodID:   38,
		EsCuentaPropia:    true,
		Recurrente:        false,
	}

	requestCreditoValido = dtos.ResultadoRequest{
		Channel:        "credit",
		Uuid:           "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		HolderDocType:  "DNI",
		HolderDocNum:   "32880325",
		CardBrand:      "Visa",
		CardNumber:     "0133212221154658845",
		Importe:        105050,
		CardExpiration: "0830",
		CardMonth:      "08",
		CardYear:       "2030",
		CardCode:       "123",
		HolderName:     "Fernando Castro",
		HolderEmail:    "fernando.castro@telco.com.ar",
		Installments:   "1",
	}

	requestDebitoValido = dtos.ResultadoRequest{
		Channel:         "debit",
		Uuid:            "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		HolderDocType:   "DNI",
		HolderDocNum:    "32880325",
		CardBrand:       "Visa",
		CardNumber:      "0133212221154658845",
		Importe:         105050,
		CardExpiration:  "0830",
		CardMonth:       "08",
		CardYear:        "2030",
		CardCode:        "123",
		HolderName:      "Fernando Castro",
		HolderEmail:     "fernando.castro@telco.com.ar",
		Installments:    "1",
		PaymentMethodID: 31,
	}

	requestOfflineValido = dtos.ResultadoRequest{
		Channel:         "offline",
		Uuid:            "asdasdasd-asdasdasd-asdasdasd-asdasdasd",
		HolderName:      "Fernando Castro",
		HolderDocType:   "DNI",
		HolderDocNum:    "32880325",
		Importe:         105050,
		PaymentMethodID: 5,
	}

	mediopagoValido = entities.Mediopago{
		Model:          gorm.Model{ID: 1},
		ChannelsID:     1,
		AdquirientesID: 1,
		InstallmentsID: 1,
		ExternalID:     "1",
		LongitudPan:    16,
		LongitudCvv:    3,
	}

	mediopagoDebin = entities.Mediopago{
		Model:          gorm.Model{ID: 38},
		ChannelsID:     4,
		AdquirientesID: 1,
		InstallmentsID: 1,
		ExternalID:     "1",
		LongitudPan:    16,
		LongitudCvv:    3,
	}

	cuentaValida = entities.Cuenta{
		Model:      gorm.Model{ID: 1},
		ClientesID: 1,
		Cuenta:     "Principal",
		Cbu:        "0340218608218026437001",
		Apikey:     "123123123123123",
	}

	pagoIntentoValido = entities.Pagointento{
		Model:                gorm.Model{ID: 1},
		PagosID:              1,
		MediopagosID:         1,
		InstallmentdetailsID: 1,
		ExternalID:           "18987",
		PaidAt:               time.Now().Local(),
		ReportAt:             time.Now().Local(),
		IsAvailable:          false,
		Amount:               100050,
		StateComment:         "approved",
		HolderName:           "Castro Fernando",
		HolderEmail:          "castro.fernando@telco.com.ar",
		TicketNumber:         "3975",
		AuthorizationCode:    "083146",
		CardLastFourDigits:   "4905",
	}

	pagoIntentoDebin = entities.Pagointento{
		Model:                gorm.Model{ID: 1},
		PagosID:              1,
		MediopagosID:         38,
		InstallmentdetailsID: 1,
		ExternalID:           "18988",
		PaidAt:               time.Now().Local(),
		ReportAt:             time.Now().Local(),
		IsAvailable:          false,
		Amount:               100050,
		StateComment:         "ACEPTADO",
	}
)

type callsMocked struct {
	mock    string
	nombre  string
	filtro  []interface{}
	retorno interface{}
	mensaje error
}

type TableGetPagoResultadoTest struct {
	Nombre        string
	Error         error
	Request       *dtos.ResultadoRequest
	MetodosMocked []callsMocked
}

func inicializar() []TableGetPagoResultadoTest {
	// variables para el segundo caso
	reqChannelError := requestCreditoValido
	reqChannelError.Channel = "efectivo"

	var channelVacio *entities.Channel
	// 3er caso
	reqCardBrandError := requestCreditoValido
	reqCardBrandError.CardBrand = "Visa Débito"

	var medioPagoVacio *entities.Mediopago
	//4to caso
	reqUuidError := requestCreditoValido
	reqUuidError.Uuid = noValidUuid

	var pagovacio *entities.Pago
	//5to caso
	reqPagoConTipoIdIncorrecto := requestCreditoValido
	reqPagoConTipoIdIncorrecto.Uuid = uuidEspecial
	pagoConTipoIdIncorrecto := pagoCreated
	pagoConTipoIdIncorrecto.PagostipoID = 999

	var pagotipoVacio *entities.Pagotipo
	//6to caso
	reqCuentaNoValida := requestCreditoValido
	reqCuentaNoValida.Uuid = "uuidparadevolcerpagotiposincuentaid"
	pagoConTipoIdSinCuenta := pagoCreated
	pagoConTipoIdSinCuenta.PagostipoID = 1000

	var cuentaVacia *entities.Cuenta
	//7mo caso
	reqPagoConCreditoIncompleto := requestCreditoValido
	reqPagoConCreditoIncompleto.CardExpiration = ""
	reqPagoConCreditoIncompleto.CardCode = ""

	var intentovacio *entities.Pagointento

	//8vo caso
	reqInstallmentsDetailsNoValido := requestCreditoValido
	reqInstallmentsDetailsNoValido.Installments = "0"

	pagoModificado := pagoCreated
	pagoModificado.PagoestadosID = 4

	// caso debin
	reqPagoConDebin := requestDebinValido
	pagoDebinModificado := pagoCreated
	pagoDebinModificado.PagoestadosID = 2

	table := []TableGetPagoResultadoTest{
		{
			"Ante un request no válido debe devolver un error de validación",
			fmt.Errorf("debe indicar el método por el cual va a pagar"),
			&requestNoValido,
			[]callsMocked{},
		},
		{
			"Cuando ingresa un channel que no existe debo devolver un error",
			fmt.Errorf("no se encontró metodo de pago con la descripción efectivo"),
			&reqChannelError,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{"efectivo"}, channelVacio, fmt.Errorf("no se encontró metodo de pago con la descripción %s", "efectivo")},
			},
		},
		{
			"Cuando la marca de tarjeta no coincide con el channel, no se encuentra mediopago en la base y devuelve un error",
			fmt.Errorf("error en medio de pago: no se encontraron medios de pago"),
			&reqCardBrandError,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqCardBrandError.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa Débito"}}, medioPagoVacio, fmt.Errorf("no se encontraron medios de pago")},
			},
		},
		{
			"Cuando el uuid no corresponde a ningún pago válido, debo devolver un error",
			fmt.Errorf("no existe pago con identificador 123456-891012345-7890123456789"),
			&reqUuidError,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqCardBrandError.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa"}}, &mediopagoValido, nil},
				{"REPO", "GetPagoByUuid", []interface{}{noValidUuid}, pagovacio, fmt.Errorf("no existe pago con identificador %s", noValidUuid)},
			},
		},
		{
			"Cuando el pago no tiene un pagotipo correcto devuelvo error",
			fmt.Errorf("no se encontró tipo de pago con el id: 999"),
			&reqPagoConTipoIdIncorrecto,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqCardBrandError.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa"}}, &mediopagoValido, nil},
				{"REPO", "GetPagoByUuid", []interface{}{uuidEspecial}, &pagoConTipoIdIncorrecto, nil},
				{"REPO", "GetPagotipoById", []interface{}{pagoConTipoIdIncorrecto.PagostipoID}, pagotipoVacio, fmt.Errorf("no se encontró tipo de pago con el id: %d", pagoConTipoIdIncorrecto.PagostipoID)},
			},
		},
		{
			"Cuando el tipo de pago no tiene una cuentaID correcta, debe devolver un error",
			fmt.Errorf("no se encontró cuenta con el id: %d", 999),
			&reqCuentaNoValida,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqCardBrandError.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa"}}, &mediopagoValido, nil},
				{"REPO", "GetPagoByUuid", []interface{}{reqCuentaNoValida.Uuid}, &pagoConTipoIdSinCuenta, nil},
				{"REPO", "GetPagotipoById", []interface{}{pagoConTipoIdSinCuenta.PagostipoID}, &entities.Pagotipo{Model: gorm.Model{ID: 1000}, CuentasID: 999}, nil},
				{"REPO", "GetCuentaById", []interface{}{int64(999)}, cuentaVacia, fmt.Errorf("no se encontró cuenta con el id: %d", 999)},
			},
		},
		{
			"Cuando el request está incompleto, el servicio de pagos devuelve un error",
			fmt.Errorf("error en el pagointento"),
			&reqPagoConCreditoIncompleto,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqPagoConCreditoIncompleto.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa"}}, &mediopagoValido, nil},
				{"REPO", "GetPagoByUuid", []interface{}{reqPagoConCreditoIncompleto.Uuid}, &pagoCreated, nil},
				{"REPO", "GetPagotipoById", []interface{}{pagoCreated.PagostipoID}, &entities.Pagotipo{Model: gorm.Model{ID: 1}, CuentasID: 1, Pagotipo: "sellos"}, nil},
				{"REPO", "GetCuentaById", []interface{}{int64(1)}, &cuentaValida, nil},
				{"CREDIT", "CreateResultado", []interface{}{&reqPagoConCreditoIncompleto, &pagoCreated, &cuentaValida}, intentovacio, fmt.Errorf("error en el pagointento")},
			},
		},
		{
			"Cuando se crea correctamente el pago debe almacenarse en la base de datos y actualizar el pago",
			nil,
			&reqInstallmentsDetailsNoValido,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqInstallmentsDetailsNoValido.Channel}, &entities.Channel{Model: gorm.Model{ID: 1}, Channel: "credit"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(1), "mediopago": "Visa"}}, &mediopagoValido, nil},
				{"REPO", "GetPagoByUuid", []interface{}{reqInstallmentsDetailsNoValido.Uuid}, &pagoCreated, nil},
				{"REPO", "GetPagotipoById", []interface{}{pagoCreated.PagostipoID}, &entities.Pagotipo{Model: gorm.Model{ID: 1}, CuentasID: 1, Pagotipo: "sellos"}, nil},
				{"REPO", "GetCuentaById", []interface{}{int64(1)}, &cuentaValida, nil},
				{"CREDIT", "CreateResultado", []interface{}{&reqInstallmentsDetailsNoValido, &pagoCreated, &cuentaValida}, &pagoIntentoValido, nil},
				{"REPO", "GetInstallmentDetailsID", []interface{}{int64(1), int64(0)}, int64(1), nil},
				{"REPO", "CreateResultado", []interface{}{&pagoIntentoValido}, true, nil},
				{"REPO", "UpdatePago", []interface{}{&pagoModificado}, true, nil},
			},
		},
		{
			"Cuando se paga con debin el estado final de exito es 2, procesando",
			nil,
			&reqPagoConDebin,
			[]callsMocked{
				{"REPO", "GetChannelByName", []interface{}{reqPagoConDebin.Channel}, &entities.Channel{Model: gorm.Model{ID: 4}, Channel: "debin"}, nil},
				{"REPO", "GetMediopago", []interface{}{map[string]interface{}{"channels_id": uint(4)}}, &mediopagoDebin, nil},
				{"REPO", "GetPagoByUuid", []interface{}{reqPagoConDebin.Uuid}, &pagoCreated, nil},
				{"REPO", "GetPagotipoById", []interface{}{pagoCreated.PagostipoID}, &entities.Pagotipo{Model: gorm.Model{ID: 1}, CuentasID: 1, Pagotipo: "sellos"}, nil},
				{"REPO", "GetCuentaById", []interface{}{int64(1)}, &cuentaValida, nil},
				{"DEBIN", "CreateResultado", []interface{}{&reqPagoConDebin, &pagoCreated, &cuentaValida}, &pagoIntentoDebin, nil},
				{"REPO", "GetInstallmentDetailsID", []interface{}{int64(1), int64(0)}, int64(1), nil},
				{"REPO", "CreateResultado", []interface{}{&pagoIntentoDebin}, true, nil},
				{"REPO", "UpdatePago", []interface{}{&pagoDebinModificado}, true, nil},
			},
		},
	}

	return table
}

func TestGetPagoResultado(t *testing.T) {
	assertions := assert.New(t)

	repositoryMock := mockrepository.MockRepository{}

	commnsMock := mockservice.MockCommonsService{}

	mockPrismaService := mockservice.MockPrismaService{}
	mockApiLinkService := mockservice.MockApiLinkService{}

	paymentMock := mockservice.MockPaymentFactory{
		MockPrisma:  &mockPrismaService,
		MockApiLink: &mockApiLinkService,
	}

	creditPaymentMock := mockservice.MockCreditPaymentMethod{}
	debinPaymentMock := mockservice.MockDebinPaymentMethod{}

	ctx := context.Background()

	paymentMock.On("GetPaymentMethod", 1).Return(&creditPaymentMock, nil)
	paymentMock.On("GetPaymentMethod", 4).Return(&debinPaymentMock, nil)
	//mockUtilService := new(mockservice.MockUtilService)
	service := checkout.NewServiceWithPayment(&repositoryMock, &commnsMock, &paymentMock) //, mockUtilService)

	table := inicializar()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {
			for _, c := range v.MetodosMocked {
				switch c.mock {
				case "REPO":
					repositoryMock.On(c.nombre, c.filtro...).Return(c.retorno, c.mensaje).Once()
				case "CREDIT":
					creditPaymentMock.On(c.nombre, c.filtro...).Return(c.retorno, c.mensaje).Once()
				case "DEBIN":
					debinPaymentMock.On(c.nombre, c.filtro...).Return(c.retorno, c.mensaje).Once()
				}
			}

			res, err := service.GetPagoResultado(ctx, v.Request)

			if err != nil {
				assertions.EqualError(err, v.Error.Error())
			}

			if err == nil {
				assertions.True(res.Exito)
			}
		})
	}
}
