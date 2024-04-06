package util

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
)

type emailReporteCrearMensaje struct {
}

func NewEmailReporteCrearMensaje() CrearMensajeMethod {
	return &emailReporteCrearMensaje{}
}

func (e *emailReporteCrearMensaje) MensajeResultado(subject string, to []string, params utildtos.RequestDatosMail) (mensaje string, erro error) {
	paramsEmail := utildtos.ParamsEmail{
		Email:                 params.Email,
		Nombre:                params.Nombre,
		Mensaje:               params.Mensaje,
		Descripcion:           params.Descripcion,
		Totales:               params.Totales,
		MensajeSegunMedioPago: params.MensajeSegunMedioPago,
	}
	var ruta_url string

	ruta_url = config.URL_TEMPLATE + "/" + "reporte_mail.html"

	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = params.From
	for _, valueTo := range to {
		Header["To"] = valueTo
	}
	Header["Subject"] = params.Asunto
	Header["Mime-Version"] = "1.0"
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	writeHeader(buffer, Header)

	text := "\r\n--" + boundary + "\r\n"
	text += "Content-Type:" + "text/html;" + "\r\n"
	buffer.WriteString(text)
	t, err := template.ParseFiles(ruta_url)
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New("error al obtener template" + err.Error())
		return
	}
	erro = t.Execute(buffer, paramsEmail)
	if erro != nil {
		erro = errors.New(err.Error())
		return
	}

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + params.Attachment.ContentType + "\r\n"
	buffer.WriteString(body)
	if params.Attachment.WithFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + params.Attachment.ContentType + ";name=\"" + params.Attachment.Name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				erro = errors.New("error al adjuntar archivo")
				return
			}
		}()
		//writeFile(buffer, ".."+config.DOC_CL+config.DIR_REPORTE+"/"+params.Attachment.Name) // descomentar en local
		//writeFile(buffer, ".."+config.DIR_REPORTE+"/"+params.Attachment.Name) // descomentar en local prueba
		writeFile(buffer, config.DIR_BASE+config.DIR_REPORTE+"/"+params.Attachment.Name) // descomentar en produccion
	}
	buffer.WriteString("\r\n--" + boundary + "--")

	mensaje = buffer.String()

	return
}
