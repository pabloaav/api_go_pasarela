package prismadtos

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type PaymentsSimpleRequest struct {
	Customerid        Customerid      `json:"customerid"`
	SiteTransactionID string          `json:"site_transaction_id"`
	SiteID            string          `json:"site_id,omitempty"`
	Token             string          `json:"token"`
	PaymentMethodID   int64           `json:"payment_method_id"`
	Bin               string          `json:"bin"`
	Amount            int64           `json:"amount"`
	Currency          EnumTipoMoneda  `json:"currency"`
	Installments      int64           `json:"installments"`
	Description       string          `json:"description"`
	PaymentType       EnumPaymentType `json:"payment_type"`
	EstablishmentName string          `json:"establishment_name,omitempty"`
	Customeremail     Customeremail   `json:"customer"`
	//Customeremail Customeremail `json:"customeremail"`
	SubPayments []interface{} `json:"sub_payments"`
}

type Customeremail struct {
	Email string `json:"email"`
}

type Customerid struct {
	ID string `json:"id"`
}

func (peymentsSimple *PaymentsSimpleRequest) ValidarProcesoPagoRequest() error {
	regularEstablishmentName := regexp.MustCompile(`^[0-9a-zA-Z!¡¿?@#$%&()=+-.,_]{1,25}$`)

	if len(peymentsSimple.SiteTransactionID) > 40 || len(peymentsSimple.SiteTransactionID) < 1 {
		return errors.New(ERROR_SITE_TRANSACTION_ID)
	}
	if commons.StringIsEmpity(peymentsSimple.SiteTransactionID) {
		return errors.New(ERROR_SITE_TRANSACTION_ID)
	}

	if len(peymentsSimple.Token) > 36 || commons.StringIsEmpity(peymentsSimple.Token) {
		return errors.New(ERROR_TOKEN_PAGO)
	}
	_, err1 := strconv.Atoi(peymentsSimple.Bin)
	if len(peymentsSimple.Bin) != 6 || err1 != nil {
		return errors.New(ERROR_BIN)
	}

	if peymentsSimple.Amount < 1 || peymentsSimple.Amount > 99999999999999 {
		return errors.New(ERROR_AMOUNT)
	}
	// valorAmount := fmt.Sprintln(reflect.TypeOf(peymentsSimple.Amount))
	// if strings.TrimSpace(valorAmount) != "int64" {
	// 	return errors.New(ERROR_AMOUNT)
	// }

	err := peymentsSimple.Currency.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}
	if peymentsSimple.Installments < 1 || peymentsSimple.Installments > 99 {
		return errors.New(ERROR_INSTALLMENTS)
	}
	if peymentsSimple.PaymentType != "single" {
		return errors.New(ERROR_PAYMENT_TYPE)
	}
	// if len(peymentsSimple.EstablishmentName) > 0 && len(peymentsSimple.EstablishmentName) < 26 {
	// 	runes := []rune(peymentsSimple.EstablishmentName)
	// 	//var result []string
	// 	for i := 0; i < len(runes); i++ {
	// 		if runes[i] < 32 || runes[i] > 126 {
	// 			return errors.New(ERROR_ESTABLISHMENT_NAME)
	// 		}
	// 		//result = append(result, int(runes[i]))
	// 	}
	// } else if len(peymentsSimple.EstablishmentName) == 0 {
	// 	return nil
	// } else {
	// 	return errors.New(ERROR_ESTABLISHMENT_NAME)
	// }
	// valorbool := regularEstablishmentName.MatchString(peymentsSimple.EstablishmentName)
	// if commons.StringIsEmpity(peymentsSimple.EstablishmentName) || !valorbool {
	// 	return errors.New(ERROR_NOMBRE_ESTABLECIMIENTO)
	// }
	//cadena := commons.SpaceStringsBuilder(peymentsSimple.EstablishmentName)
	cadena := commons.ReplaceCharacters(peymentsSimple.EstablishmentName, " ", "_")
	valorbool := regularEstablishmentName.MatchString(cadena)
	if commons.StringIsEmpity(peymentsSimple.EstablishmentName) || !valorbool {
		return errors.New(ERROR_NOMBRE_ESTABLECIMIENTO)
	}

	if !commons.IsEmailValid(peymentsSimple.Customeremail.Email) {
		return errors.New(ERROR_EMAIL)
	}

	return nil
}
