package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"text/template"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
)

type emailTemplateCrearMensaje struct {
}

func NewEmailTemplateCrearMensaje() CrearMensajeMethod {
	return &emailTemplateCrearMensaje{}
}

func (e *emailTemplateCrearMensaje) MensajeResultado(subject string, to []string, params utildtos.RequestDatosMail) (mensaje string, erro error) {
	paramsEmail := utildtos.ParamsEmail{
		Email:                 params.Email,
		Nombre:                params.Nombre,
		Mensaje:               params.Mensaje,
		Descripcion:           params.Descripcion,
		MensajeSegunMedioPago: params.MensajeSegunMedioPago,
	}
	var ruta_url string
	if params.FiltroReciboPago {
		ruta_url = config.URL_TEMPLATE + "/" + "recibo_pago.html"
	} else {
		ruta_url = config.URL_TEMPLATE + "/" + "send_mail.html"
	}
	t, err := template.ParseFiles(ruta_url)
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New("error al obtener template" + err.Error())
		return
	}
	buf := new(bytes.Buffer)
	erro = t.Execute(buf, paramsEmail)
	if erro != nil {
		erro = errors.New(err.Error())
		return
	}

	body := buf.String()
	header := make(map[string]string)
	header["From"] = params.From
	for _, valueTo := range to {
		header["To"] = valueTo
	}
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	mensaje = message
	return
}
