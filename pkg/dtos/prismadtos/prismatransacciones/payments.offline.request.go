package prismadtos

import (
	"errors"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type PaymentsOfflineRequest struct {
	Customer          DataCustomer   `json:"customer"`
	SiteTransactionID string         `json:"site_transaction_id"`
	Token             string         `json:"token"`
	PaymentMethodID   int64          `json:"payment_method_id"`
	Amount            int64          `json:"amount"`
	Currency          EnumTipoMoneda `json:"currency"`
	PaymentType       string         `json:"payment_type"`
	Email             string         `json:"email"`
	InvoiceExpiration string         `json:"invoice_expiration"`
	CodP3             string         `json:"cod_p3"`
	CodP4             string         `json:"cod_p4"`
	Client            string         `json:"client"`
	Surcharge         int64          `json:"surcharge"`
	PaymentMode       string         `json:"payment_mode"`
}

func (paymentsOffline *PaymentsOfflineRequest) Validar() error {

	err := paymentsOffline.Customer.Identification.Type.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}

	if len(paymentsOffline.Customer.Identification.Number) < 7 || commons.StringIsEmpity(paymentsOffline.Customer.Identification.Number) {
		return errors.New(ERROR_NRO_DOC)
	}

	_, err = strconv.Atoi(paymentsOffline.Customer.Identification.Number)
	if err != nil {
		return errors.New(ERROR_NRO_DOC)
	}

	if commons.StringIsEmpity(paymentsOffline.Customer.Name) {
		return errors.New(ERROR_NOMBRE_PAGADOR)
	}
	/*
		FIXME se debe verificar si el uuid es de 8 o 40 maximo
	*/
	if len(paymentsOffline.SiteTransactionID) > 40 || len(paymentsOffline.SiteTransactionID) < 1 {
		return errors.New(ERROR_SITE_TRANSACTION_ID)
	}

	if len(paymentsOffline.Token) != 36 {
		return errors.New(ERROR_TOKEN_PAGO)
	}

	sizeAmount := strconv.FormatInt(paymentsOffline.Amount, 10)
	if len(sizeAmount) <= 0 && len(sizeAmount) >= 9 {
		return errors.New(ERROR_AMOUNT)
	}

	if paymentsOffline.Amount <= 0 || paymentsOffline.Amount >= 99999999 {
		return errors.New(ERROR_AMOUNT)
	}

	if !commons.IsEmailValid(paymentsOffline.Email) {
		return errors.New(ERROR_EMAIL)
	}

	err = paymentsOffline.Currency.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}

	///FIXME PARA VALIDAR LAS CANTIDADES DE DIAS DE CODP3 Y CODP4
	/// TENER EN CONSIDERACION LAS FECHAS DE PRIMER Y SEGINDO VENCIMIENTO

	if len(paymentsOffline.CodP3) != 2 {
		return errors.New(ERROR_CODP3)
	}
	_, err = strconv.Atoi(paymentsOffline.CodP3)
	if err != nil {
		return errors.New(ERROR_CODP3)
	}

	if len(paymentsOffline.CodP4) != 3 {
		return errors.New(ERROR_CODP4)
	}
	_, err = strconv.Atoi(paymentsOffline.CodP4)
	if err != nil {
		return errors.New(ERROR_CODP4)
	}

	if len(paymentsOffline.Client) != 8 {
		return errors.New(ERROR_CLIENTE_NRO)
	}

	_, err = strconv.Atoi(paymentsOffline.Client)
	if err != nil {
		return errors.New(ERROR_CLIENTE_NRO)
	}

	sizeSurcharge := strconv.FormatInt(paymentsOffline.Surcharge, 10)
	if len(sizeSurcharge) <= 0 || len(sizeSurcharge) >= 8 {
		return errors.New(ERROR_SURCHANGE)
	}

	if paymentsOffline.Surcharge <= 0 || paymentsOffline.Surcharge >= 99999999 {
		return errors.New(ERROR_SURCHANGE)
	}

	if paymentsOffline.PaymentMode != "offline" {
		return errors.New(ERROR_MODO_PAGO)
	}

	return nil
}

// func DiasEntreDosFecha() {

// 	fechaActual := time.Now()
// 	fmt.Printf("la fecha actual es: %v\n", fechaActual)
// 	// fecha1 := "15-05-2021"
// 	// fecha2 := "20-05-2021"
// 	tTime := time.Date(2021, time.Month(06), 30, 0, 0, 0, 0, time.Local)
// 	fmt.Printf("%v\n", tTime)
// 	fmt.Printf("tTime es: %s\n", tTime.Format("2006/1/2"))

// }
