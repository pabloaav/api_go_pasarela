package linkdebin

import (
	"encoding/json"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

//Se usa para los gets de los debines
type ResponseGetDebinLink struct {
	Debin DebinDetalleLink `json:"debin"`
}

type DebinDetalleLink struct {
	ComprobanteId    string                     `json:"comprobanteId"`
	EsCuentaPropia   bool                       `json:"esCuentaPropia"`
	Concepto         linkdtos.EnumConceptoDebin `json:"concepto"`
	TiempoExpiracion int64                      `json:"tiempoExpiracion"`
	Importe          int64                      `json:"importe"`
	Moneda           string                     `json:"moneda"`
	UsuarioId        string                     `json:"usuarioId"`
	Estado           linkdtos.EnumEstadoDebin   `json:"estado"`
	FechaExpiracion  string                     `json:"fechaExpiracion"`
	FechaCreacion    string                     `json:"fechaCreacion"`
	Devuelto         bool                       `json:"devuelto"`
	Recurrente       bool                       `json:"recurrente"`
	// DescripcionPrestacion string               `json:"descripcionPrestacion"` //En la documentación está pero no tiene respuesta en la api
	Prestacion    string               `json:"prestacion"`
	ContracargoId string               `json:"contracargoid"`
	Comprador     CompradorDetalleLink `json:"comprador"`
	Vendedor      VendedorDetalleLink  `json:"vendedor"`
}

type CompradorDetalleLink struct {
	Cbu      string              `json:"cbu"`
	AliasCbu string              `json:"aliasCbu"`
	CuitCuil string              `json:"cuitcuil"`
	Tipo     string              `json:"tipo"`
	Moneda   linkdtos.EnumMoneda `json:"moneda"`
	Titular  string              `json:"titular"`
	Banco    string              `json:"banco"` //Identificación de entidad bancária. Código de 3 digitos BCRA
	Sucursal string              `json:"sucursal"`
	//Sucursal de la entidad bancaria. Identificador de la sucursal donde está radicada la cuenta.
}

func (c *CompradorDetalleLink) IsValid() error {

	err := tools.EsCuitValido(c.CuitCuil)

	if err != nil {
		return err
	}

	if len(c.Cbu) == 0 {
		err := tools.EsAliasCbuValido(c.AliasCbu)
		if err != nil {
			return err
		}
	} else {
		err := tools.EsCbuValido(c.Cbu, tools.ERROR_CBU_COMPRADOR)
		if err != nil {
			return err
		}
		//Limpio el alias porque solo debe tener uno el cbu o el alias
		c.AliasCbu = ""
	}

	return nil
}

type VendedorDetalleLink struct {
	Cbu      string              `json:"cbu"`
	AliasCbu string              `json:"aliasCbu"`
	CuitCuil string              `json:"cuitcuil"`
	Tipo     string              `json:"tipo"` //Tipo de cuenta bancaria
	Moneda   linkdtos.EnumMoneda `json:"moneda"`
	Titular  string              `json:"titular"`
	Banco    string              `json:"banco"` // Identificador de entidad bancaria 3 digitos del BCRA
	Sucursal string              `json:"sucursal"`
	//Sucursal de la entidad bancaria.
	//Identificador de la sucursal donde esta radicada la cuenta
}

func (c *VendedorDetalleLink) IsValid() error {

	if len(c.Cbu) == 0 {
		err := tools.EsAliasCbuValido(c.AliasCbu)
		if err != nil {
			return err
		}
	} else {
		err := tools.EsCbuValido(c.Cbu, tools.ERROR_CBU_COMPRADOR)
		if err != nil {
			return err
		}
		//Limpio el alias porque solo debe tener uno el cbu o el alias
		c.AliasCbu = ""
	}

	return nil
}



func (r *ResponseGetDebinLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}