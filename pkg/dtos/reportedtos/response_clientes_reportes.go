package reportedtos

import (
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ResponseClientesReportes struct {
	Clientes        string
	RazonSocial     string
	Cuit            string
	Email           []string
	Fecha           string
	Pagos           []PagosReportes
	Rendiciones     []*ResponseReportesRendiciones
	Reversiones     []Reversiones
	CantOperaciones string
	TotalCobrado    string
	RendicionTotal  string
	TotalIva        string
	TotalComision   string
	TotalRevertido  string
	TipoReporte     string
	TipoArchivoPdf  bool
}

// type FactoryEmail struct {
// 	Clientes    string
// 	Email       string
// 	Fecha       string
// 	Pagos       []PagosReportes
// 	Rendiciones []Rendiciones
// }

type PagosReportes struct {
	Cuenta    string
	Id        string
	FechaPago string
	MedioPago string
	Estado    string
	Tipo      string
	// Cuotas    string
	Monto string
}

type Rendiciones struct {
	Cuenta                  string // Nombre de la cuenta del cliente
	Id                      string // external_reference enviada por el cliente
	FechaCobro              string // fecha que el pagador realizo el pago
	FechaDeposito           string // fecha que se le envio el dinero al cliente(transferencia)
	ImporteCobrado          string // importe solicitud de pago
	ImporteDepositado       string // importe depositado al cliente
	CantidadBoletasCobradas string // pago items
	// ComisionPorcentaje      string // comision de telco cobrada al cliente
	// ComisionIva             string // iva Cobrado al cliente
	Comision string // comision de telco cobrada al cliente
	Iva      string // iva Cobrado al cliente
}

type Reversiones struct {
	Cuenta        string
	Id            string
	MedioPago     string
	Monto         string
	PagoRevertido PagoRevertido
}

type PagoRevertido struct {
	IdPago            string
	PagoEstado        string
	ReferenciaExterna string
	Items             []ItemsRevertidos
	IntentoPago       IntentoPagoRevertido
}

type ItemsRevertidos struct {
	IdItems       string
	Cantidad      string
	Descripcion   string
	Monto         string
	Identificador string
}

type IntentoPagoRevertido struct {
	IdIntentoPago string
	IdTransaccion string
	FechaPago     string
	ImportePagado string
}

func (rev *Reversiones) EntityToReversiones(entityReversion entities.Reversione) {
	rev.Cuenta = entityReversion.PagoIntento.Pago.PagosTipo.Cuenta.Cuenta
	rev.Id = entityReversion.PagoIntento.Pago.ExternalReference
	rev.MedioPago = entityReversion.PagoIntento.Mediopagos.Mediopago
	rev.Monto = strconv.FormatInt(entityReversion.Amount, 10) //entityReversion.Amount
}

func (pr *PagoRevertido) EntityToPagoRevertido(entityPago entities.Pago) {
	pr.IdPago = strconv.Itoa(int(entityPago.ID))
	pr.PagoEstado = string(entityPago.PagoEstados.Estado)
	pr.ReferenciaExterna = entityPago.ExternalReference

}
func (ir *ItemsRevertidos) EntityToItemsRevertidos(entityItems entities.Pagoitems) {
	ir.IdItems = strconv.Itoa(int(entityItems.ID))
	ir.Cantidad = strconv.Itoa(entityItems.Quantity)
	ir.Descripcion = entityItems.Description
	ir.Monto = strconv.Itoa(int(entityItems.Amount))
	ir.Identificador = entityItems.Identifier
}
func (ipr *IntentoPagoRevertido) EntityToIntentoPagoRevertido(entityPagoIntento entities.Pagointento) {
	ipr.IdIntentoPago = strconv.Itoa(int(entityPagoIntento.ID))
	ipr.IdTransaccion = entityPagoIntento.TransactionID
	ipr.FechaPago = entityPagoIntento.PaidAt.String()
	ipr.ImportePagado = strconv.Itoa(int(entityPagoIntento.Amount))
}

func ToEntityRegistroReporte(request ResponseClientesReportes) (response entities.Reporte) {

	var fechacobros string
	var fecharendicion string
	var reportedetalle []entities.Reportedetalle
	if len(request.Pagos) > 0 {
		for _, pg := range request.Pagos {
			reportedetalle = append(reportedetalle, entities.Reportedetalle{
				PagosId:    pg.Id,
				Monto:      pg.Monto,
				Mediopago:  pg.MedioPago,
				Estadopago: pg.Estado,
			})
		}
		fechacobros = request.Fecha
	}

	if len(request.Reversiones) > 0 {
		var totalMontoRevertido int64
		for _, reversion := range request.Reversiones {
			montoRevertidoInt64, _ := strconv.ParseInt(reversion.Monto, 10, 64)
			reportedetalle = append(reportedetalle, entities.Reportedetalle{
				PagosId:    reversion.PagoRevertido.ReferenciaExterna,
				Monto:      reversion.Monto,
				Mediopago:  reversion.MedioPago,
				Estadopago: reversion.PagoRevertido.PagoEstado,
			})
			// acumular los montos revertidos
			totalMontoRevertido += montoRevertidoInt64
		} // fin del for

		// formatear para que tenga el mismo formato que usa la tabla reportes
		p := message.NewPrinter(language.Spanish)
		request.TotalCobrado = p.Sprintf("%.2f", float64(totalMontoRevertido)/100)
		fechacobros = request.Fecha
	}

	if len(request.Rendiciones) > 0 {
		fecharendicion = request.Fecha
	}

	response = entities.Reporte{
		Cliente:        request.Clientes,
		Tiporeporte:    request.TipoReporte,
		Totalcobrado:   request.TotalCobrado,
		Totalrendido:   request.RendicionTotal,
		Fechacobranza:  fechacobros,
		Fecharendicion: fecharendicion,
		Reportedetalle: reportedetalle,
	}
	return
}
