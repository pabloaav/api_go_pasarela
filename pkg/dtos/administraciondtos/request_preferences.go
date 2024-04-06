package administraciondtos

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type RequestPreferences struct {
	ClientId       string                `json:"clientId"`
	MainColor      string                `json:"mainColor"`
	SecondaryColor string                `json:"secondaryColor"`
	RutaLogo       string                `json:"ruta_logo"`
	File           *multipart.FileHeader `json:"archivo"`
}

func (request *RequestPreferences) Validar() (erro error) {
	mensaje := "los parametros enviados no son válidos"
	mensaje_archivo := "extensión del archivo no valido"
	clienteId, err := strconv.Atoi(request.ClientId)
	archivo := strings.Split(request.File.Filename, ".")[1]
	if err != nil {
		return fmt.Errorf(mensaje)
	}
	if clienteId < 1 {
		return fmt.Errorf(mensaje)
	}
	if archivo != "jpg" && archivo != "png" {
		return fmt.Errorf(mensaje_archivo)
	}

	if commons.StringIsEmpity(request.MainColor) {
		return fmt.Errorf(mensaje)
	}
	if commons.StringIsEmpity(request.SecondaryColor) {
		return fmt.Errorf(mensaje)
	}

	return
}
