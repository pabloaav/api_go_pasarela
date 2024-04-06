package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
)

type emailAdjuntoCrearMensaje struct {
}

func NewEmailAdjuntoCrearMensaje() CrearMensajeMethod {
	return &emailAdjuntoCrearMensaje{}
}

func (e *emailAdjuntoCrearMensaje) MensajeResultado(subject string, to []string, params utildtos.RequestDatosMail) (mensaje string, erro error) {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = params.From
	Header["To"] = strings.Join(params.Email, ";")
	Header["Subject"] = params.Asunto
	Header["Mime-Version"] = "1.0"
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Date"] = time.Now().String()
	writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + params.Attachment.ContentType + "\r\n"
	body += "\r\n" + params.Mensaje + "\r\n"
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

func writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

func writeFile(buffer *bytes.Buffer, fileName string) {
	logs.Info(fileName)
	file, err := ioutil.ReadFile(fileName)
	logs.Info(err)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}
