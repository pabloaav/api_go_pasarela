package dtos

import (
	"fmt"
	"strconv"
	"strings"
)

type ResultadoRequest struct {
	Channel               string `json:"channel"`
	HolderName            string `json:"holder_name"`
	HolderEmail           string `json:"holder_email"`
	HolderDocType         string `json:"holder_docType"`
	HolderDocNum          string `json:"holder_docNum"`
	HolderCuit            string `json:"holder_cuit"`
	CardBrand             string `json:"card_brand"`
	CardNumber            string `json:"card_number"`
	CardExpiration        string `json:"card_expiration"`
	CardYear              string `json:"card_year"`
	CardMonth             string `json:"card_month"`
	CardCode              string `json:"card_code"`
	Cbu                   string `json:"cbu"`
	Alias                 string `json:"alias"`
	Installments          string `json:"installments"`
	Uuid                  string `json:"uuid"`
	ComprobanteID         string `json:"id"`
	EsCuentaPropia        bool   `json:"es_cuenta_propia"`
	ConceptoAbreviado     string `json:"concepto_abreviado"`
	TiempoExpiracion      int64  `json:"tiempo_expiracion"`
	Importe               int64  `json:"importe"`
	Valorcupon            int64  `json:"valorcupon"`
	Moneda                string `json:"moneda"`
	Recurrente            bool   `json:"recurrente"`
	DescripcionPrestacion string `json:"descripcion_prestacion"`
	PaymentMethodID       int64  `json:"mediopago_id"`
	Ip 				      string `json:"ip"`
}

func (rr *ResultadoRequest) Validar() error {
	if len(rr.Channel) <= 0 {
		return fmt.Errorf("debe indicar el método por el cual va a pagar")
	}
	if len(rr.Uuid) <= 0 {
		return fmt.Errorf("debe indicar el código identificador del pago")
	}
	switch rr.Channel {
	case "debin":
		if len(rr.Cbu) <= 0 && len(rr.Alias) <= 0 {
			return fmt.Errorf("debe indicar el número de CBU o alias")
		}
		installmentID, _ := strconv.ParseInt(rr.Installments, 10, 64)
		if installmentID != 1 {
			return fmt.Errorf("error al pagar con debin, vuelva a intentarlo")
		}
	case "credit":
		if len(rr.CardNumber) <= 0 {
			return fmt.Errorf("debe indicar el número de la tarjeta")
		}
	case "debit":
		if len(rr.CardNumber) <= 0 {
			return fmt.Errorf("debe indicar el número de la tarjeta")
		}
	case "offline":
		if len(rr.HolderDocNum) <= 0 && len(rr.HolderCuit) <= 0 {
			return fmt.Errorf("debe indicar un número de identificación")
		}
	}
	return nil
}

func (rr *ResultadoRequest) ToFormatStr() {
	rr.HolderName = strings.ToUpper(rr.HolderName)
}
