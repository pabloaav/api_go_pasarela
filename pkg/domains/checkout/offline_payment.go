package checkout

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type offlinePayment struct {
	service prisma.Service
}

func NewOfflinePayment(s prisma.Service) PaymentMethod {
	return &offlinePayment{service: s}
}

func (c *offlinePayment) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transactionID string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error) {
	//serv := prisma.Resolve()

	// Validaciones básicas
	err := request.Validar()
	if err != nil {
		return nil, err
	}

	// armo request del token
	customer := prismadtos.DataCustomer{
		Identification: prismadtos.IdentificationCustomer{
			Type:   prismadtos.EnumTipoDocumento(request.HolderDocType),
			Number: request.HolderDocNum,
		},
		Name: request.HolderName,
	}

	tokenReq := prismadtos.StructToken{
		DataOffline: prismadtos.OfflineTokenRequest{
			Customer: customer,
		},
		TypePay: "offline",
	}

	// llamo servicio de obtener el token
	tokenInterface, err := c.service.SolicitarToken(tokenReq)
	if err != nil {
		return nil, err
	}

	// casteo la respuesta a estructura de token
	token := tokenInterface.(prismadtos.OfflineTokenResponse)
	logs.Info(token)

	//firstDueDateString := pago.FirstDueDate.Format("2006-01-02")
	//invoiceExpiration := firstDueDateString[8:10] + firstDueDateString[5:7] + firstDueDateString[2:4]
	//daysBetweenDueDates := int((pago.SecondDueDate.Sub(pago.FirstDueDate).Hours() / 24))
	payReq := prismadtos.StructPayments{
		PagoOffline: prismadtos.PaymentsOfflineRequest{
			Customer:          customer,
			SiteTransactionID: transactionID,
			Token:             token.ID,
			PaymentMethodID:   request.PaymentMethodID,
			//PaymentMethodID: 26,
			Amount:      request.Importe,
			Currency:    "ARS",
			PaymentType: "single",
			Email:       request.HolderEmail,
			//InvoiceExpiration: invoiceExpiration,
			InvoiceExpiration: "191123",
			//CodP3:             fmt.Sprint(daysBetweenDueDates),
			CodP3: "10",
			//CodP4: fmt.Sprintf("%0*d", 3, daysBetweenDueDates+100),
			CodP4: "134",
			//Client: fmt.Sprint(cuenta.ID),
			Client: "12345678",
			//Surcharge:   int64(pago.SecondTotal) * 100,
			Surcharge:   1234567,
			PaymentMode: "offline",
		},
		TypePay: "offline",
	}

	// llamo al servicio del payment
	var payment prismadtos.PaymentsOfflineResponse
	paymentInterface, err := c.service.Payments(payReq)
	// casteo la respuesta a estructura de payment
	if err != nil {
		payment.Status = err.Error()
	} else {
		payment = paymentInterface.(prismadtos.PaymentsOfflineResponse)
	}

	logs.Info(payment)
	// armo la respuesta del metodo
	paidAt, _ := time.Parse("2006-01-02T15:04Z", payment.Date)
	if err != nil {
		paidAt = time.Time{}
	}
	response := entities.Pagointento{
		PagosID:      int64(pago.ID),
		MediopagosID: payReq.PagoOffline.PaymentMethodID,
		ExternalID:   fmt.Sprint(payment.ID),
		PaidAt:       paidAt,
		ReportAt:     time.Now().Local(),
		IsAvailable:  false,
		Amount:       entities.Monto(request.Importe),
		//Valorcupon:           entities.Monto(request.Importe),
		StateComment:         payment.Status,
		Barcode:              payment.Barcode,
		InstallmentdetailsID: 1,
		HolderType:           request.HolderDocType,
		HolderNumber:         request.HolderDocNum,
		HolderName:           request.HolderName,
		HolderEmail:          request.HolderEmail,
		TicketNumber:         payment.StatusDetails.Ticket,
		AuthorizationCode:    payment.StatusDetails.CardAuthorizationCode,
		TransactionID:        transactionID,
	}

	return &response, nil
}
