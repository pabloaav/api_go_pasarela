package linkdebin

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

type RequestDebinCreateLink struct {
	Comprador CompradorCreateDebinLink `json:"comprador"`
	Vendedor  VendedorCreateLink       `json:"vendedor"`
	Debin     DebinCreateLink          `json:"debin"`
}

func (c *RequestDebinCreateLink) IsValid() error {

	err := c.Comprador.IsValid()
	if err != nil {
		return err
	}
	err = c.Vendedor.IsValid()

	if err != nil {
		return err
	}
	err = c.Debin.IsValid()
	if err != nil {
		return err
	}

	return nil
}

type CompradorCreateDebinLink struct {
	Cuit   string     `json:"cuit"`
	Cuenta CuentaLink `json:"cuenta"`
}

func (c *CompradorCreateDebinLink) IsValid() error {

	err := tools.EsCuitValido(c.Cuit)

	if err != nil {
		return err
	}
	err = c.Cuenta.IsValid()

	if err != nil {
		return err
	}

	return nil
}

type VendedorCreateLink struct {
	Cbu      string `json:"cbu"`
	AliasCbu string `json:"aliasCbu"`
}

func (c *VendedorCreateLink) IsValid() error {

	if len(c.Cbu) == 0 {
		err := tools.EsAliasCbuValido(c.AliasCbu)
		if err != nil {
			return err
		}
	} else {
		if len(c.Cbu) != 22 {
			err := errors.New("longitud de cbu no es válido")
			return err
		}
		// err := tools.EsCbuValido(c.Cbu, tools.ERROR_CBU_VENDEDOR)
		// if err != nil {
		// 	return err
		// }
		//Limpio el alias porque solo debe tener uno el cbu o el alias
		c.AliasCbu = ""
	}

	return nil
}

type DebinCreateLink struct {
	ComprobanteId         string                     `json:"comprobanteId"`
	EsCuentaPropia        bool                       `json:"esCuentaPropia"`
	Concepto              linkdtos.EnumConceptoDebin `json:"concepto"`
	TiempoExpiracion      int64                      `json:"tiempoExpiracion"`
	Descripcion           string                     `json:"descripcion"`
	Importe               int64                      `json:"importe"`
	Moneda                linkdtos.EnumMoneda        `json:"moneda"`
	Recurrente            bool                       `json:"recurrente"`
	DescripcionPrestacion string                     `json:"descripcionPrestacion"`
}

func (d *DebinCreateLink) IsValid() error {

	if tools.EsStringVacio(d.ComprobanteId) {
		return errors.New(tools.ERROR_IDENTIFICADORDEBIN)
	}

	if len(d.ComprobanteId) > 10 {
		return errors.New(tools.ERROR_IDENTIFICADORDEBIN_INVALIDO)
	}

	_, err := strconv.Atoi(d.ComprobanteId)

	if err != nil {
		return errors.New(tools.ERROR_IDENTIFICADORDEBIN_INVALIDO)
	}

	if len(d.Descripcion) > 100 {
		return errors.New(tools.ERROR_DESCRIPCIONDEBIN)
	}
	if d.TiempoExpiracion < 1 || d.TiempoExpiracion > 4320 {
		return errors.New(tools.ERROR_TIEMPOEXPIRACIONDEBIN)
	}
	if d.Importe < 0 {
		return errors.New(tools.ERROR_IMPORTE)
	}
	err = d.Moneda.IsValid()
	if err != nil {
		return err
	}
	err = d.Concepto.IsValid()
	if err != nil {
		return err
	}

	//Lo hago eso porque solo utilizaremos debin spot y en este caso este campo deber estar vacío.
	d.DescripcionPrestacion = ""

	return nil
}

func (r *RequestDebinCreateLink) String() string {
	jsonFormat, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonFormat)
}
