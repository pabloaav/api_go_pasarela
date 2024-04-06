package checkout_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/checkout"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	pagoExpirado = pagoCreated

	noValidUuid          = "123456-891012345-7890123456789"
	uuidEspecial         = "123456-891012345-7890"
	uuidPagoExpirado     = "123456-891012345"
	validUuidVencimiento = "123456-891012345-7890-54"
	validUuid            = "123456-891012345-7890123456789-123456789"

	pagoTipo = entities.Pagotipo{
		Model: gorm.Model{
			ID: 1,
		},
		CuentasID:       1,
		Pagotipo:        "sellos",
		BackUrlSuccess:  "http://www.dgrcorrientes.gov.ar",
		BackUrlPending:  "http://www.dgrcorrientes.gov.ar",
		BackUrlRejected: "http://www.dgrcorrientes.gov.ar",
		// IncludedChannels:     "credit, debit, offline, debin",
		// IncludedInstallments: "1,3,6,12",
	}
)

func TestGetPago(t *testing.T) {
	assertions := assert.New(t)

	repositoryMock := mockrepository.MockRepository{}

	var pagovacio *entities.Pago

	// cuando se le manda un uuid no valido el repo devuelve un error
	repositoryMock.On("GetPagoByUuid", noValidUuid).Return(pagovacio, fmt.Errorf("no existe pago con identificador %s", noValidUuid))

	// cuando se le manda un uuid de pago expirado el repo devuelve el pago con el campo createdAt de hace una hora
	pagoExpirado.CreatedAt = time.Now().Local().Add(-time.Hour)
	repositoryMock.On("GetPagoByUuid", uuidPagoExpirado).Return(&pagoExpirado, nil)

	// cuando se le manda un uuid valido el pago q devuelve el repo debe ser valido
	pagoCreated.CreatedAt = time.Now().Local()
	pagoCreated.FirstDueDate = time.Now().Local().Add(time.Hour * 24)
	repositoryMock.On("GetPagoByUuid", validUuid).Return(&pagoCreated, nil)

	// cuando mando un uuid especial el repo va a devolver un pago con un tipo de pago incorrecto
	pagoConTipoIdIncorrecto := pagoCreated
	pagoConTipoIdIncorrecto.CreatedAt = time.Now().Local()
	pagoConTipoIdIncorrecto.PagostipoID = 999
	repositoryMock.On("GetPagoByUuid", uuidEspecial).Return(&pagoConTipoIdIncorrecto, nil)
	var pagotipoVacio *entities.Pagotipo
	repositoryMock.On("GetPagotipoById", pagoConTipoIdIncorrecto.PagostipoID).Return(pagotipoVacio, fmt.Errorf("no se encontró tipo de pago con el id: %d", pagoConTipoIdIncorrecto.PagostipoID))

	// para testear q el importe se asigne según la fecha de vencimiento:
	pagoConPrimeraFechaVencida := pagoCreated
	fechaVencimiento = time.Now().Local().Add(-time.Hour * 48)
	pagoConPrimeraFechaVencida.FirstDueDate = fechaVencimiento
	pagoConPrimeraFechaVencida.Pagoitems = []entities.Pagoitems{}
	repositoryMock.On("GetPagoByUuid", validUuidVencimiento).Return(&pagoConPrimeraFechaVencida, nil)
	repositoryMock.On("GetPagotipoById", pagoConPrimeraFechaVencida.PagostipoID).Return(&pagoTipo, nil)

	commnsMock := mockservice.MockCommonsService{}
	prismaService := mockservice.MockPrismaService{}
	offlineService := mockservice.MockPagoOffLineService{}
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStoreService := new(mockservice.MockStoreService)
	service := checkout.NewService(&repositoryMock, &commnsMock, &prismaService, &offlineService, mockUtilService, mockRepositoryWebHook, mockStoreService)

	t.Run("Ante un request con uuid vacío se debe devolver un error", func(t *testing.T) {

		request := ""

		_, err := service.GetPago(request)

		assertions.EqualErrorf(err, "debe enviar código único del pago, envió: ", "debe enviar código único del pago, envió: %s", request)
	})

	t.Run("Ante un uuid sin pagos asociados se debe devolver un error", func(t *testing.T) {

		_, err := service.GetPago(noValidUuid)

		assertions.EqualErrorf(err, "error al obtener pago: no existe pago con identificador 123456-891012345-7890123456789", "error al obtener pago: no existe pago con identificador %s", noValidUuid)
	})

	t.Run("Ante un pago Creado hace mas de media hora se debe devolver un error", func(t *testing.T) {

		_, err := service.GetPago(uuidPagoExpirado)

		assertions.EqualError(err, "el pago expiró, vuelva a generarlo")
	})

	t.Run("Cuando el repositorio manda un pago con pagoTipoID incorrecto debe devolver un error", func(t *testing.T) {

		_, err := service.GetPago(uuidEspecial)

		assertions.EqualErrorf(err, "no se encontró tipo de pago con el id: 999", "no se encontró tipo de pago con el id: %d", pagoConTipoIdIncorrecto.PagostipoID)
	})

	t.Run("Cuando la fecha de hoy supera a la primera fecha de vencimiento del pago, se debe asignar el segundo importe a la respuesta del servicio", func(t *testing.T) {

		res, err := service.GetPago(validUuidVencimiento)

		assertions.NoError(err)
		assertions.Equal(res.Total, pagoConPrimeraFechaVencida.SecondTotal)
		// el pago con vencimiento ademas viene sin items, por lo q controlo q al parsearlo a string traiga solo corchetes
		assertions.Equal(res.Items, "[]")
	})

	t.Run("Cuando la fecha de hoy no supera a la primera fecha de vencimiento del pago, se debe asignar el primer importe a la respuesta del servicio", func(t *testing.T) {

		res, err := service.GetPago(validUuid)

		assertions.NoError(err)
		assertions.Equal(res.Total, pagoCreated.FirstTotal)
		// el pago con vencimiento ademas viene con 2 items, por lo q controlo q al parsearlo el string sea el correcto
		assertions.Equal(res.Items, "[{\"quantity\":1,\"description\":\"Item 1 impuestos\",\"amount\":500},{\"quantity\":1,\"description\":\"Item 2 intereses\",\"amount\":500.5}]")
	})
}
