package commons

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type reversion struct {
	m pdf.Maroto
}

type ReversionData struct {
	Pago    PagoData
	Intento IntentoData
	Items   []ItemsReversionData
}

type ItemsReversionData struct {
	Cantidad      string
	Descripcion   string
	Identificador string
	Monto         string
}

type PagoData struct {
	ReferenciaExterna string
	MedioPago         string
	Monto             string
	IdPago            string
	Estado            string
}

type IntentoData struct {
	IdIntento     string
	IdTransaccion string
	FechaPago     string
	Importe       string
}

func (r *reversion) buildTitle() {

	fecha_impresion_pdf := fmt.Sprintf(time.Now().Format("02-01-2006"))

	r.m.RegisterHeader(func() {

		r.m.SetBorder(false)

		r.m.Line(1.0,
			props.Line{
				Color: color.Color{
					Red:   0,
					Green: 0,
					Blue:  255,
				},
			})

		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Reversiones de Pagos", props.Text{
					Top:   1,
					Align: consts.Center,
					Size:  12,
				})
			})
		})

		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Fecha: "+fecha_impresion_pdf, props.Text{
					Size:   10,
					Style:  consts.BoldItalic,
					Top:    1,
					Family: consts.Helvetica,
					Align:  consts.Right,
				})
			})
		})

		r.m.Line(1.0,
			props.Line{
				Color: color.Color{
					Red:   0,
					Green: 0,
					Blue:  255,
				},
			})
	}) // Fin de RegisterHEader
}

func (r *reversion) buildHeadingsReversiones(data []ReversionData) {

	for _, dato := range data {

		// fila 1 Datos on Top
		r.m.Row(7, func() {

			r.m.SetBorder(true)

			// col1 Ref Externa
			r.m.Col(4, func() {

				r.m.Text("Referencia Externa: "+dato.Pago.ReferenciaExterna, props.Text{
					// Top:    3,
					Left:   2,
					Style:  consts.Bold,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})
			})

			// col2 Medio de pago
			r.m.Col(4, func() {
				r.m.Text("Medio de Pago: "+dato.Pago.MedioPago, props.Text{
					// Top:    3,
					Left:   2,
					Style:  consts.Bold,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})
			})

			// col3 Monto
			r.m.Col(4, func() {
				r.m.Text("Monto: "+dato.Pago.Monto, props.Text{
					// Top:    3,
					Left:   2,
					Style:  consts.Bold,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})
			})
		})

		// fila 2 datos Pago e Intento
		r.m.Row(35, func() {

			// col 1 Pago
			r.m.Col(4, func() {

				r.m.Text("Pago", props.Text{

					Style:  consts.Bold,
					Family: consts.Courier,
					Size:   10,
					Align:  consts.Center,
					Color:  getDarkGrayColor(),
				})

				r.m.Text("Id: "+dato.Pago.IdPago, props.Text{
					Top:    7,
					Left:   2,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})

				r.m.Text("Estado: "+dato.Pago.Estado, props.Text{
					Top:  14,
					Left: 2,

					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})

			})

			// col 2 Intento
			r.m.Col(8, func() {

				r.m.Text("Intento", props.Text{
					Align:  consts.Center,
					Style:  consts.Bold,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})

				r.m.Text("Id Intento: "+dato.Intento.IdIntento, props.Text{
					Top:    7,
					Left:   2,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})

				r.m.Text("Transaccion: "+dato.Intento.IdTransaccion, props.Text{
					Top:    14,
					Left:   2,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})
				r.m.Text("Fecha: "+dato.Intento.FechaPago, props.Text{
					Left:   2,
					Top:    21,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})
				r.m.Text("Importe: "+dato.Intento.Importe, props.Text{
					Left:   2,
					Top:    28,
					Family: consts.Courier,
					Size:   10,
					Color:  getDarkGrayColor(),
				})

			})
		})

		r.m.SetBorder(false)

		// Mostrar los items de cada pago revertido
		buildBodyItems(r.m, dato.Items)

		// Agregar una nueva pagina para separar las reversiones
		r.m.AddPage()

	} // fin de for range

}

func buildBodyItems(m pdf.Maroto, items []ItemsReversionData) {
	header, contents := getMediumContent(items)
	m.Line(1)

	m.SetBackgroundColor(getSoftGrayColor())

	m.Row(7, func() {
		m.Col(12, func() {
			m.Text("Detalle de Pagos", props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
				Left:  4,
				Color: color.NewWhite(),
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.Line(1)

	m.TableList(header, contents, props.TableList{
		ContentProp: props.TableListContent{
			Family:    consts.Courier,
			Style:     consts.Italic,
			GridSizes: []uint{2, 3, 5, 2},
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{2, 3, 5, 2},
			Family:    consts.Courier,
			Style:     consts.BoldItalic,
			Color:     color.Color{100, 0, 0},
		},
		Line: true,
		LineProp: props.Line{
			Color: color.Color{
				Red:   128,
				Green: 221,
				Blue:  205,
			},
			Style: consts.Dashed,
		},
	})
}

func (r *reversion) buildFooter() {
	r.m.SetFirstPageNb(1)
	r.m.RegisterFooter(func() {
		r.m.Row(5, func() {
			r.m.Col(12, func() {
				r.m.Text(strconv.Itoa(r.m.GetCurrentPage()), props.Text{
					Align: consts.Right,
					Size:  8,
					Top:   10,
				})
			})
		})
	})
}

func getMediumContent(items []ItemsReversionData) ([]string, [][]string) {
	header := []string{"Cantidad", "Descripcion", "Identificador", "Monto"}

	contents := [][]string{}

	for _, item := range items {
		contents = append(contents, []string{item.Cantidad, item.Descripcion, item.Identificador, item.Monto})
	}

	return header, contents
}

func GetReversionesPdf(reversiones []ReversionData, nombreCliente, fecha string) error {
	reversionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	reversionPdf.SetPageMargins(10, 10, 10)

	var rev reversion
	rev.m = reversionPdf

	rev.buildTitle()

	// registrar el footer
	rev.buildFooter()

	// Las cabeceras de cada pago revertidos y los items de cada pago
	rev.buildHeadingsReversiones(reversiones)

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf(config.DIR_BASE + config.DOC_CL + "/reportes")
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			return err
		}
	}

	err := rev.m.OutputFileAndClose(tempFolder + "/" + nombreCliente + "-" + fecha + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func getSoftGrayColor() color.Color {
	return color.Color{
		Red:   68,
		Green: 68,
		Blue:  68,
	}
}
