package administraciondtos

import (
	"encoding/base64"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

type ResponseFacturacionPaginado struct {
	Clientes []ResponseFacturacion `json:"data"`
	Meta     dtos.Meta             `json:"meta"`
}

type ResponseFacturacion struct {
	Id             uint               `json:"id"`      // Id cliente
	Cliente        string             `json:"cliente"` // nombre Cliente abreviado
	RazonSocial    string             `json:"razon_social"`
	NombreFantasia string             `json:"nombre_fantasia"`
	Email          string             `json:"email"`
	Emailcontacto  string             `json:"email_contacto"`
	Emails         []string           `json:"emails"`
	Cuit           string             `json:"cuit"`
	Personeria     string             `json:"personeria"`
	Impuestos      []ResponseImpuesto `json:"impuestos"`
	// Iva              ResponseFacturacionIva      `json:"iva"`
	// Iibb             ResponseFacturacionIibb     `json:"iibb"`
	RetiroAutomatico bool   `json:"retiro_automatico"`
	ReporteBatch     bool   `json:"reporte_batch"`
	NombreReporte    string `json:"nombre_reporte"`
	OrdenDiaria      bool   `json:"orden_diaria"`
	// SplitCuentas     bool                        `json:"split_cuentas"`
	Cuenta          []ResponseFacturacionCuenta `json:"cuenta"`
	SujetoRetencion bool                        `json:"sujeto_retencion"`
	Formulario_8125 bool                        `json:"formulario_8125"`
	Envio 					entities.Envio 							`json:"envio"`
}

func (r *ResponseFacturacion) FromEntity(c entities.Cliente) {
	r.Id = c.ID
	r.Cliente = c.Cliente
	r.RazonSocial = c.Razonsocial
	r.NombreFantasia = c.Nombrefantasia
	r.Email = c.Email
	r.Emailcontacto = c.Emailcontacto
	r.Cuit = c.Cuit
	r.Personeria = c.Personeria
	r.RetiroAutomatico = c.RetiroAutomatico
	r.ReporteBatch = c.ReporteBatch
	r.NombreReporte = c.NombreReporte
	r.OrdenDiaria = c.OrdenDiaria
	r.SujetoRetencion = c.SujetoRetencion
	r.Formulario_8125 = c.Formulario_8125
	if c.Iva != nil {
		// var iva ResponseFacturacionIva
		// iva.FromEntity(*c.Iva)
		// r.Iva = iva
		var impuestoiva ResponseImpuesto
		iva := entities.Impuesto{
			Model:      gorm.Model{ID: c.Iva.ID},
			Impuesto:   c.Iva.Impuesto,
			Porcentaje: c.Iva.Porcentaje,
			Tipo:       c.Iva.Tipo,
			Fechadesde: c.Iva.Fechadesde,
		}
		impuestoiva.FromImpuesto(iva)
		r.Impuestos = append(r.Impuestos, impuestoiva)
	}
	if c.Iibb != nil {
		// var iibb ResponseFacturacionIibb
		// iibb.FromEntity(*c.Iibb)
		// r.Iibb = iibb
		var impuestoIibb ResponseImpuesto
		iibb := entities.Impuesto{
			Model:      gorm.Model{ID: c.Iibb.ID},
			Impuesto:   c.Iibb.Impuesto,
			Porcentaje: c.Iibb.Porcentaje,
			Tipo:       c.Iibb.Tipo,
			Fechadesde: c.Iibb.Fechadesde,
		}
		impuestoIibb.FromImpuesto(iibb)
		r.Impuestos = append(r.Impuestos, impuestoIibb)
	}
	if c.Cuentas != nil {
		for _, c := range *c.Cuentas {
			var cuenta ResponseFacturacionCuenta
			cuenta.FromEntity(c)
			r.Cuenta = append(r.Cuenta, cuenta)
		}
		for i := 0; i < len(r.Cuenta); i++ {

			for j := 0; j < len(r.Cuenta[i].Comisiones); j++ {

				r.Cuenta[i].Comisiones[j].ControlVigenteActual(r.Cuenta[i].Comisiones)

			}
		}
	}
	if c.Contactosreportes != nil {
		for _, valueEmail := range *c.Contactosreportes {
			r.Emails = append(r.Emails, valueEmail.Email)
		}
	}
	r.Envio = c.Envio
}

type ResponseFacturacionIva struct {
	Id         uint    `json:"id"`
	Impuesto   string  `json:"impuesto"`
	Tipo       string  `json:"tipo"`
	Porcentaje float64 `json:"porcentaje"`
}

func (r *ResponseFacturacionIva) FromEntity(c entities.Impuesto) {
	r.Id = c.ID
	r.Impuesto = c.Impuesto
	r.Tipo = c.Tipo
	r.Porcentaje = c.Porcentaje
}

type ResponseFacturacionIibb struct {
	Id         uint    `json:"id"`
	Impuesto   string  `json:"impuesto"`
	Tipo       string  `json:"tipo"`
	Porcentaje float64 `json:"porcentaje"`
}

func (r *ResponseFacturacionIibb) FromEntity(c entities.Impuesto) {
	r.Id = c.ID
	r.Impuesto = c.Impuesto
	r.Tipo = c.Tipo
	r.Porcentaje = c.Porcentaje
}

type ResponseFacturacionCuenta struct {
	Id                   uint                            `json:"id"`
	Cuenta               string                          `json:"cuenta"`
	Cbu                  string                          `json:"cbu"`
	Cvu                  string                          `json:"cvu"`
	Apikey               string                          `json:"apikey"`
	DiasRetiroAutomatico int64                           `json:"dias_retiro_automatico"`
	SplitCuentas         bool                            `json:"split_cuentas"`
	Rubro                ResponseFacturacionRubro        `json:"rubro"`
	Comisiones           []ResponseFacturacionComisiones `json:"comisiones"`
	TiposPago            []ResponseFacturacionTiposPago  `json:"tipos_pago"`
}

func (r *ResponseFacturacionCuenta) FromEntity(c entities.Cuenta) {
	hashApikey := base64.StdEncoding.EncodeToString([]byte(c.Apikey))
	r.Id = c.ID
	r.Cuenta = c.Cuenta
	r.Cbu = c.Cbu
	r.Cvu = c.Cvu
	r.Apikey = hashApikey
	r.DiasRetiroAutomatico = c.DiasRetiroAutomatico
	r.SplitCuentas = c.SplitCuentas
	if c.Rubro != nil {
		var rubro ResponseFacturacionRubro
		rubro.FromEntity(*c.Rubro)
		r.Rubro = rubro
	}
	if c.Cuentacomisions != nil {
		for _, c := range *c.Cuentacomisions {
			var comisiones ResponseFacturacionComisiones
			comisiones.FromEntity(c)
			r.Comisiones = append(r.Comisiones, comisiones)
		}
	}
	if c.Pagotipos != nil {
		for _, pt := range *c.Pagotipos {
			var pagotipos ResponseFacturacionTiposPago
			pagotipos.FromEntity(pt)
			r.TiposPago = append(r.TiposPago, pagotipos)
		}
	}
}

type ResponseFacturacionRubro struct {
	Id    uint   `json:"id"`
	Rubro string `json:"rubro"`
}

func (r *ResponseFacturacionRubro) FromEntity(c entities.Rubro) {
	r.Id = c.ID
	r.Rubro = c.Rubro
}

type ResponseFacturacionComisiones struct {
	Nombre        string  `json:"nombre"`
	Comision      float64 `json:"comision"`
	VigenciaDesde string  `json:"vigencia_desde"`
	VigenteActual bool    `json:"vigente_actual"`
	Canal         string  `json:"canal"`
	MedioPagoId   int     `json:"medio_pago_id"`
	PagoCuota     bool    `json:"pago_cuota"`
}

func (r *ResponseFacturacionComisiones) FromEntity(c entities.Cuentacomision) {
	r.Nombre = c.Cuentacomision
	if c.ChannelArancel.Tipocalculo == "FIJO" {
		r.Comision = c.Comision
	} else {
		r.Comision = c.Comision + c.ChannelArancel.Importe
	}
	r.VigenciaDesde = c.VigenciaDesde.Format("2006-01-02")
	r.Canal = c.Channel.Nombre
	r.MedioPagoId = int(c.Mediopagoid)
	r.PagoCuota = c.Pagocuota
}

type ResponseFacturacionTiposPago struct {
	PagoTipo                 string
	BackUrlSuccess           string
	BackUrlPending           string
	BackUrlRejected          string
	BackUrlNotificacionPagos string
}

func (r *ResponseFacturacionTiposPago) FromEntity(c entities.Pagotipo) {
	r.PagoTipo = c.Pagotipo
	r.BackUrlSuccess = c.BackUrlSuccess
	r.BackUrlPending = c.BackUrlPending
	r.BackUrlRejected = c.BackUrlRejected
	r.BackUrlNotificacionPagos = c.BackUrlNotificacionPagos
}

func (r *ResponseFacturacionComisiones) ControlVigenteActual(comparaciones []ResponseFacturacionComisiones) {

	lastComisionOfType := *LastComisionOfType(r, comparaciones)

	if *r == lastComisionOfType {
		r.VigenteActual = true
	}

	return

}

func LastComisionOfType(r *ResponseFacturacionComisiones, comparaciones []ResponseFacturacionComisiones) (last *ResponseFacturacionComisiones) {
	lastComision := r
	for y := 0; y < len(comparaciones); y++ {

		c := comparaciones[y]

		if MismaComision(lastComision, c) {
			if ComparacionComision(r, c) {
				lastComision = &c
			}

		}

	}
	last = lastComision
	return
}

func MismaComision(c *ResponseFacturacionComisiones, r ResponseFacturacionComisiones) (res bool) {

	if r.MedioPagoId != c.MedioPagoId {
		return
	}
	if r.Canal != c.Canal {
		return
	}
	if r.PagoCuota != c.PagoCuota {
		return
	}

	res = true
	return
}

func ComparacionComision(r *ResponseFacturacionComisiones, c ResponseFacturacionComisiones) (res bool) {

	fechaActual, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	fechaVigencia, err := time.Parse("2006-01-02", r.VigenciaDesde)
	fechaVigenciaComparacion, err := time.Parse("2006-01-02", c.VigenciaDesde)

	if err != nil {
		return
	}

	if fechaActual.After(fechaVigenciaComparacion) && fechaVigencia.Before(fechaVigenciaComparacion) {
		res = true
	}

	return
}
