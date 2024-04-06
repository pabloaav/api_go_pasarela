package utildtos

import (
	"errors"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type RequestDatosMail struct {
	Asunto           string
	Email            []string `json:"email"` //Obligatorio
	From             string   //Obligatorio
	Nombre           string   //Opcional
	Mensaje          string   //Obligatorio
	CamposReemplazar []string //opcional
	FiltroReciboPago bool

	Descripcion           DescripcionTemplate // opcional
	Totales               TotalesTemplate     // opcional
	AdjuntarEstado        bool
	Attachment            Attachment
	TipoEmail             EnumTipoEmail
	MensajeSegunMedioPago MensajeSegunMedioPagoStruct
	FiltroReporte         int
}

// MensajeSegunMedioPagoStruct sirve para determinar el titulo del email segun el medio de pago, y un breve comentario al respecto
type MensajeSegunMedioPagoStruct struct {
	Title   string
	Content string
}
type DescripcionTemplate struct {
	Fecha         string
	Cliente       string
	Cuit          string
	EmailContacto string
	Detalles      []DetallesPago
	TotalPagado   string
}

type TotalesTemplate struct {
	Titulo       string
	TipoReporte  string
	Cantidad     string
	Elemento     string
	TotalCobrado string
	TotalRendido string
}

type DetallesPago struct {
	Descripcion string
	Cantidad    string
	Monto       string
}

func (s *RequestDatosMail) IsValid() error {
	var message = errors.New(PARAMS_INVALID)
	for _, value := range s.Email {
		if !commons.IsEmailValid(value) {
			return message
		}
	}
	if len(s.Mensaje) == 0 {
		return message
	}
	if len(s.CamposReemplazar) != 0 {
		if !strings.Contains(s.Mensaje, "#") {
			return message
		}
	}

	if len(s.From) == 0 {
		return message
	}

	if s.AdjuntarEstado {
		err := s.Attachment.IsValid()
		if err != nil {
			return err
		}
	}

	// if !strings.Contains(s.Mensaje, "$") {
	// 	return message
	// }
	// if len(s.CamposReemplazar) == 0 {
	// 	return message
	// }
	return nil
}

type Attachment struct {
	Name        string
	ContentType string
	WithFile    bool
}

func (at *Attachment) IsValid() error {
	var message = errors.New(PARAMS_INVALID)
	if len(at.Name) == 0 {
		return message
	}
	if len(at.ContentType) == 0 {
		return message
	}

	return nil
}
