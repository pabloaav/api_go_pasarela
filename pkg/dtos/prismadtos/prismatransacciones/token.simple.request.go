package prismadtos

import (
	"errors"
	"strconv"
	"time"

	tools "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type Card struct {
	CardNumber               string                   `json:"card_number"`
	CardExpirationMonth      string                   `json:"card_expiration_month"`
	CardExpirationYear       string                   `json:"card_expiration_year"`
	SecurityCode             string                   `json:"security_code"`
	CardHolderName           string                   `json:"card_holder_name"`
	CardHolderIdentification CardHolderIdentification `json:"card_holder_identification"`
}

type CardHolderIdentification struct {
	TypeDni   EnumTipoDocumento `json:"type"`
	NumberDni string            `json:"number"`
}

func (card *Card) Validar() error {
	if card.CardNumber == "" && card.CardExpirationMonth == "" && card.CardExpirationYear == "" && card.SecurityCode == "" && card.CardHolderName == "" && card.CardHolderIdentification.TypeDni == "" &&
		card.CardHolderIdentification.NumberDni == "" {
		return errors.New(ERROR_ESTRUCTURA_INCORRECTA)
	}
	if len(card.CardNumber) < 15 {
		return errors.New(ERROR_NUMBER_CARD)
	}
	check := tools.NewAlgoritmoVerificacion()
	if !check.ChequearTarjeta(card.CardNumber) {
		return errors.New(ERROR_NUMBER_CARD)
	}
	// if !ckeck.ChequearTarjeta(card.CardNumber) {
	// 	return errors.New(ERROR_NUMBER_CARD)
	// }
	f := time.Now()
	fechaFormat := f.Format("06-01")
	anioActual, _ := strconv.Atoi(fechaFormat[0:2])
	max := len(fechaFormat)
	mesActual, _ := strconv.Atoi(fechaFormat[3:max])

	monthCurrent, _ := strconv.Atoi(card.CardExpirationMonth)
	yearCurrent, _ := strconv.Atoi(card.CardExpirationYear)

	if yearCurrent < anioActual {
		return errors.New(ERROR_DATE_CARD)
	}
	if yearCurrent == anioActual && monthCurrent < mesActual {
		return errors.New(ERROR_DATE_CARD)
	}
	//strings.Trim(card.CardHolderName, " ") == ""
	if len(card.CardHolderName) < 1 || tools.StringIsEmpity(card.CardHolderName) {
		return errors.New(ERROR_HOLDER_NAME)
	}
	err := card.CardHolderIdentification.TypeDni.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}
	if len(card.SecurityCode) < 3 || len(card.SecurityCode) > 4 {
		return errors.New(ERROR_SECURITYCODE)
	}
	return nil
}
