package checkout_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	apiKey               = "apiKeydeprueba"
	apiKeyErronea        = "apiKey de prueba erronea"
	apiKeyClienteErroneo = "apiKey prueba de cliente sin cuenta o pagotipo"

	ValidPaymentRequest = dtos.PagoRequest{
		PayerName:         "Fernando Castro",
		Description:       "Pago de obligaciones con DGR",
		FirstTotal:        100050,
		FirstDueDate:      "01-07-2021",
		ExternalReference: "111",
		SecondDueDate:     "01-08-2021",
		SecondTotal:       105050,
		PayerEmail:        "fernando.castro@telco.com.ar",
		PaymentType:       "sellos",
		Items: []entities.Pagoitems{
			{
				Quantity:    1,
				Description: "Item 1 impuestos",
				Amount:      50000,
			},
			{
				Quantity:    1,
				Description: "Item 2 intereses",
				Amount:      50050,
			},
		},
	}

	fechaVencimiento, _        = time.Parse("02-01-2006", ValidPaymentRequest.FirstDueDate)
	fechaSegundoVencimiento, _ = time.Parse("02-01-2006", ValidPaymentRequest.SecondDueDate)

	pago = entities.Pago{
		PagostipoID:       1,
		PagoestadosID:     1,
		Description:       ValidPaymentRequest.Description,
		FirstDueDate:      fechaVencimiento,
		FirstTotal:        entities.Monto(ValidPaymentRequest.FirstTotal),
		SecondDueDate:     fechaSegundoVencimiento,
		SecondTotal:       entities.Monto(ValidPaymentRequest.SecondTotal),
		PayerName:         ValidPaymentRequest.PayerName,
		PayerEmail:        ValidPaymentRequest.PayerEmail,
		ExternalReference: ValidPaymentRequest.ExternalReference,
		Metadata:          ValidPaymentRequest.Metadata,
		Uuid:              "123456-891012345-7890123456789-123456789",
		PdfUrl:            "http://localhost:3300/checkout/bill/123456-891012345-7890123456789-123456789",
		Pagoitems:         ValidPaymentRequest.Items,
	}

	pagoCreated = entities.Pago{
		Model: gorm.Model{
			ID: 1,
		},
		PagostipoID:       1,
		PagoestadosID:     1,
		Description:       ValidPaymentRequest.Description,
		FirstDueDate:      fechaVencimiento,
		FirstTotal:        entities.Monto(ValidPaymentRequest.FirstTotal),
		SecondDueDate:     fechaSegundoVencimiento,
		SecondTotal:       entities.Monto(ValidPaymentRequest.SecondTotal),
		PayerName:         ValidPaymentRequest.PayerName,
		PayerEmail:        ValidPaymentRequest.PayerEmail,
		ExternalReference: ValidPaymentRequest.ExternalReference,
		Metadata:          ValidPaymentRequest.Metadata,
		Uuid:              "123456-891012345-7890123456789-123456789",
		PdfUrl:            "http://localhost:3300/checkout/bill/123456-891012345-7890123456789-123456789",
		Pagoitems:         ValidPaymentRequest.Items,
	}
)

func TestNewPago(t *testing.T) {

	assertions := assert.New(t)

	commnsMock := mockservice.MockCommonsService{}
	commnsMock.On("NewUUID").Return("123456-891012345-7890123456789-123456788").Once()
	commnsMock.On("NewUUID").Return("123456-891012345-7890123456789-123456789")

	repositoryMock := mockrepository.MockRepository{}

	// Caso en que el repositorio devuelve algo con exito
	repositoryMock.On("GetClienteByApikey", apiKey).Return(&entities.Cliente{
		Model: gorm.Model{
			ID: 1,
		},
		IvaID:          1,
		IibbID:         2,
		Cliente:        "DGR",
		Razonsocial:    "Direccion General de Rentas",
		Nombrefantasia: "Rentas",
		Cuentas: &[]entities.Cuenta{
			{
				Model: gorm.Model{
					ID: 1,
				},
				ClientesID: 1,
				Cuenta:     "Cuenta recaudadora 04",
				Cbu:        "00321132554865",
				Cvu:        "00321132554865",
				Apikey:     "apiKeydeprueba",
				Pagotipos: &[]entities.Pagotipo{
					{
						Model: gorm.Model{
							ID: 1,
						},
						CuentasID:       1,
						Pagotipo:        "sellos",
						BackUrlSuccess:  "www.asdasd.com",
						BackUrlPending:  "www.asdasd.com",
						BackUrlRejected: "www.asdasd.com",
						// IncludedChannels:     "credit, debit, offline, debin",
						// IncludedInstallments: "1,3,6,12",
					},
				},
			},
		},
	}, nil)
	var cliente *entities.Cliente
	// caso en el que el repositorio no encuentra cliente con apiKey indicada
	repositoryMock.On("GetClienteByApikey", apiKeyErronea).Return(cliente, fmt.Errorf("no se encontró cliente con apikey: %s", apiKeyErronea))
	// caso en el que devuelve un cliente sin cuentas o pagotipo
	repositoryMock.On("GetClienteByApikey", apiKeyClienteErroneo).Return(&entities.Cliente{
		Model: gorm.Model{
			ID: 1,
		},
		IvaID:          1,
		IibbID:         2,
		Cliente:        "DGR",
		Razonsocial:    "Direccion General de Rentas",
		Nombrefantasia: "Rentas",
	}, nil)

	// CreatePago se llama al generar el registro en pago
	pagoUuidExistente := pago
	pagoUuidExistente.Uuid = "123456-891012345-7890123456789-123456788"
	pagoUuidExistente.PdfUrl = "http://localhost:3300/checkout/bill/123456-891012345-7890123456789-123456788"
	var pagoNil *entities.Pago
	pagoFirstTotalNoValido := pago
	pagoFirstTotalNoValido.FirstTotal = entities.Monto(999999999999)
	repositoryMock.On("CreatePago", &pagoUuidExistente).Return(pagoNil, fmt.Errorf("no se pudo generar registro pago: %s", "uuid_UNIQUE"))
	repositoryMock.On("CreatePago", &pagoFirstTotalNoValido).Return(pagoNil, fmt.Errorf("NewPago: %s", "Out of range value"))
	repositoryMock.On("CreatePago", &pago).Return(&pagoCreated, nil)

	prismaMock := mockservice.MockPrismaService{}
	offlineService := mockservice.MockPagoOffLineService{}
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStoreService := new(mockservice.MockStoreService)
	service := checkout.NewService(&repositoryMock, &commnsMock, &prismaMock, &offlineService, mockUtilService, mockRepositoryWebHook, mockStoreService)

	ctx := context.Background()

	t.Run("Ante un request no válido debe devolver un error de validación", func(t *testing.T) {

		request := ValidPaymentRequest
		request.PayerName = ""

		_, err := service.NewPago(ctx, &request, apiKey)

		assertions.Error(err)
	})

	t.Run("Si la apiKey no está relacionada a ningun cliente debe devolver un error", func(t *testing.T) {

		_, err := service.NewPago(ctx, &ValidPaymentRequest, "apiKey de prueba erronea")

		assertions.EqualErrorf(err, "no se encontró cliente con apikey: apiKey de prueba erronea", "no se encontró cliente con apikey: %s")
	})

	t.Run("Si el cliente no tiene configurado alguna cuenta con el tipo de pago indicado, debo devolver un error", func(t *testing.T) {

		_, err := service.NewPago(ctx, &ValidPaymentRequest, apiKeyClienteErroneo)

		assertions.EqualError(err, "en la configuración de cuentas, no hay tipo de pago correcto para sellos")
	})

	t.Run("Si la fecha de vencimiento cumple con el formato pero falla al ser parseada por algun valor incorrecto debo indicar el error", func(t *testing.T) {

		request := ValidPaymentRequest
		request.FirstDueDate = "30-30-1999"

		_, err := service.NewPago(ctx, &request, apiKey)

		assertions.Error(err)
		assertions.Contains(err.Error(), "error en fecha de vencimiento")
	})

	t.Run("Si la segunda fecha de vencimiento cumple con el formato pero falla al ser parseada por algun valor incorrecto debo indicar el error", func(t *testing.T) {

		request := ValidPaymentRequest
		request.SecondDueDate = "30-30-1999"

		_, err := service.NewPago(ctx, &request, apiKey)

		assertions.Error(err)
		assertions.Contains(err.Error(), "error en fecha de segundo vencimiento:")
	})

	t.Run("Al crear un pago con UUID ya existente, se debe ejecutar al menos dos veces el llamado CrearPago al repositorio y NewUUID a commons", func(t *testing.T) {

		_, err := service.NewPago(ctx, &ValidPaymentRequest, apiKey)

		commnsMock.AssertNumberOfCalls(t, "NewUUID", 2)
		repositoryMock.AssertNumberOfCalls(t, "CreatePago", 2)
		assertions.NoError(err)
	})

	t.Run("Al haber un error en insertar el pago en la base de datos, se debe devolver como respuesta", func(t *testing.T) {
		request := ValidPaymentRequest
		request.FirstTotal = entities.Monto(999999999999).Int64()

		_, err := service.NewPago(ctx, &request, apiKey)

		assertions.Error(err)
		assertions.Contains(err.Error(), "NewPago:")
	})
}

func TestConcurrenciaNewPago(t *testing.T) {
	t.Parallel()

	sqlClient := database.NewMySQLClient()
	auditoriaRepository := auditoria.NewAuditoriaRepository(sqlClient)
	auditoriaService := auditoria.New(auditoriaRepository)
	repository := checkout.NewRepository(sqlClient, auditoriaService)
	prismaMock := mockservice.MockPrismaService{}
	offlineService := mockservice.MockPagoOffLineService{}
	file := os.File{}
	fileRepostiory := commons.NewFileRepository(&file)
	commons := commons.NewCommons(fileRepostiory)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStoreService := new(mockservice.MockStoreService)
	checkoutService := checkout.NewService(repository, commons, &prismaMock, &offlineService, mockUtilService, mockRepositoryWebHook, mockStoreService)

	apiKey = "123123123123123"

	ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

	t.Run("Prueba 1 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "123456"
		time.Sleep(time.Millisecond * 1000)

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 1: %s", err.Error())
		}
	})

	t.Run("Prueba 2 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "223456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 2: %s", err.Error())
		}
	})

	t.Run("Prueba 3 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "323456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 3: %s", err.Error())
		}
	})

	t.Run("Prueba 4 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "423456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 4: %s", err.Error())
		}
	})

	t.Run("Prueba 5 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "523456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 5: %s", err.Error())
		}
	})

	t.Run("Prueba 6 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "623456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 6: %s", err.Error())
		}
	})

	t.Run("Prueba 7 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "723456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 7: %s", err.Error())
		}
	})

	t.Run("Prueba 8 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "823456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 8: %s", err.Error())
		}
	})

	t.Run("Prueba 9 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "923456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 9: %s", err.Error())
		}
	})

	t.Run("Prueba 10 de crear un pago", func(t *testing.T) {
		t.Parallel()
		request := ValidPaymentRequest
		request.ExternalReference = "1023456"

		_, err := checkoutService.NewPago(ctx, &request, apiKey)
		if err != nil {
			t.Errorf("Error al crear pago 10: %s", err.Error())
		}
	})

}
