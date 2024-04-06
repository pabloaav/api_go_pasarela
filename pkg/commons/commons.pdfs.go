package commons

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type recaudacion struct {
	m pdf.Maroto
}

type ActorRecaudacion struct {
	RazonSocial    string `json:"razonSocial"`
	Domicilio      string `json:"domicilio"`
	Cuit           string `json:"cuit"`
	IngresosBrutos string `json:"ingresosBrutos"`
	Iva            string `json:"iva"`
}

type FechasRecaudacion struct {
	Cobro    string `json:"cobro"`
	Deposito string `json:"deposito"`
	Proceso  string `json:"proceso"`
}

type DataBodyRecaudacion struct {
	ChannelName       string `json:"channel_name"`
	ImporteCobrado    string `json:"importe_cobrado"`
	ImporteDepositado string `json:"importe_depositado"`
	CantidadBoletas   string `json:"cantidad_boletas"`
	Comisiones        string `json:"comisiones"`
	IvaComision       string `json:"iva_comision"`
	RetIva            string `json:"ret_iva"`
}
type DataFooterRecaudacion struct {
	RecuperoComisiones string `json:"recupero_comisiones"`
	IvaRecupero        string `json:"iva_recupero"`
	Totales            string `json:"totales"`
}

func (r *recaudacion) buildHeading(cliente ActorRecaudacion, fechas FechasRecaudacion, fileName string) {

	// Logo
	r.m.Row(20, func() {
		r.m.Col(3, func() {
			_ = r.m.FileImage(config.DIR_BASE+"/assets/images/wee_reduce.png", props.Rect{
				Center:  true,
				Percent: 80,
			})
		})
	})

	// Titulo
	r.m.Row(10, func() {

		r.m.Col(12, func() {
			r.m.Text("Recaudación WEE!", props.Text{
				Top:    3,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   12,
				Align:  consts.Center,
				Color:  getDarkGrayColor(),
			})
		})
	})

	r.m.Line(1)

	// Datos cliente
	r.m.Row(40, func() {
		r.m.ColSpace(1)

		r.m.Col(2, func() {

			r.m.Text("Cliente: ", props.Text{
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

		r.m.Col(6, func() {

			r.m.Text(cliente.RazonSocial, props.Text{
				Top:    7,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})

			r.m.Text("Domicilio: "+cliente.Domicilio, props.Text{
				Top:    14,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
			r.m.Text("C.U.I.T.: "+cliente.Cuit, props.Text{
				Left:   2,
				Top:    21,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
			r.m.Text("Ing. Brutos: "+cliente.IngresosBrutos, props.Text{
				Left:   2,
				Top:    28,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
			r.m.Text("I.V.A.: "+cliente.Iva, props.Text{
				Left:   2,
				Top:    35,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})

		})
		fecha_impresion_pdf := fmt.Sprintf(time.Now().Format("02-01-2006"))
		r.m.Col(3, func() {
			r.m.Text("Fecha: "+fecha_impresion_pdf, props.Text{
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

	})
	r.m.Line(1)
	// fila 2 fechas
	r.m.Row(10, func() {

		r.m.ColSpace(4)

		// col2 fecha de deposito
		r.m.Col(4, func() {
			r.m.Text("Fecha Depósito: "+fechas.Deposito, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

		// col3 fecha de proceso
		r.m.Col(4, func() {
			r.m.Text("Fecha Proceso: "+fechas.Proceso, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})
	})

	r.m.Row(10, func() {
		r.m.Col(12, func() {
			r.m.Text("Nombre de Archivo Transferido: "+fileName, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})
	})

}

func (r *recaudacion) buildContent(header []string, contenido [][]string) {
	grayColor := getGrayColor()
	r.m.Row(10, func() {
		r.m.Col(12, func() {
			r.m.Text("Detalle de Valores a Recaudar", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	//TODO: reemplazar con los medios de pagos
	medios := []string{"Debito", "Credito", "Debin"}

	// Para cada medio de pago se imprime barra de titulo, detalle de pagos y subtotales
	for _, medio := range medios {

		r.m.SetBackgroundColor(getDarkGrayColor())

		r.m.Row(7, func() {
			r.m.Col(3, func() {
				r.m.Text(medio, props.Text{
					Top:   1.5,
					Size:  9,
					Style: consts.Bold,
					Align: consts.Center,
					Color: color.NewWhite(),
				})
			})
			r.m.ColSpace(9)
		})

		r.m.SetBackgroundColor(color.NewWhite())

		r.m.TableList(header, contenido, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			ContentProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			Align:                consts.Center,
			AlternatedBackground: &grayColor,
			HeaderContentSpace:   1,
			Line:                 false,
		})
		// SUBTOTALES
		r.m.Row(5, func() {
			r.m.Col(7, func() {
				r.m.Text("Subtotales:", props.Text{
					Left:  2,
					Top:   1.5,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})

			r.m.Col(5, func() {
				r.m.Text("A cobrar:", props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A depositar: ", props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})

	}

}

func (r *recaudacion) buildFooter(dataFooter interface{}) {
	r.m.SetFirstPageNb(1)
	r.m.RegisterFooter(func() {

		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Total cobrado: "+"1000", props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 13})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Total a rendir: "+"1100", props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 16})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Cantidad de boletas: "+"15", props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 19})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Comisiones a cobrar: "+"20", props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 19})
			})
		})

	})
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   51,
		Green: 51,
		Blue:  51,
	}
}

func GetRecaudacionPdf(request interface{}, ruta, nombrearchivo string) error {
	recaudacionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	recaudacionPdf.SetPageMargins(20, 10, 20)

	var rec recaudacion
	rec.m = recaudacionPdf

	/************* Cabecera  *************/
	var cliente ActorRecaudacion
	var fechas FechasRecaudacion

	cliente.RazonSocial = "Direccion Prov de Energia de Corrientes"
	cliente.Domicilio = "Junin 1240"
	cliente.Cuit = ""
	cliente.IngresosBrutos = ""
	cliente.Iva = "Responsable Inscripto"

	fechas.Cobro = "24/11/22"
	fechas.Deposito = "29/11/22"
	fechas.Proceso = "25/11/22"

	fileName := "archivo de rendicion"

	// encabezdo
	rec.buildHeading(cliente, fechas, fileName)

	/************* Cuerpo *************/
	// TODO: cambiar por datos reales
	header, contenido := getContentsAndHeadersTable("datos reales")

	// contenido
	rec.buildContent(header, contenido)

	rec.m.Line(1)

	/************* Footer  *************/

	var dataFooter DataFooterRecaudacion
	dataFooter.RecuperoComisiones = "200"
	dataFooter.IvaRecupero = "21"
	dataFooter.Totales = "221"

	// pie de informe
	rec.buildFooter(dataFooter)

	// Se crea la carpeta temporal en caso de que no exista
	// carpetaTempImages := "./pdfs"
	// if _, err := os.Stat(carpetaTempImages); os.IsNotExist(err) {
	// 	err = os.MkdirAll(carpetaTempImages, 0755)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// crear archivo en la carpeta pdfs
	err := rec.m.OutputFileAndClose(ruta + "/" + nombrearchivo + ".pdf")
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
		return err
	}

	return nil
}

// 8 columnas
func getContentsAndHeadersTable(items interface{}) ([]string, [][]string) {
	header := []string{"Cuenta", "Referencia", "Fecha Cobro", "Importe Cobrado", "Importe a Depositar", "Cant Boletas", "Comision", "IVA"}

	contents := [][]string{}

	// 	for _, item := range items {
	// 	 	contents = append(contents, []string{item.Cantidad, item.Descripcion, item.Identificador, item.Monto})
	//  }
	for i := 0; i < 15; i++ {
		contents = append(contents, []string{"TASAS MUNICIPALES", "Jf5ftsAQG", "02/01/2023", "6670.2", "6549.14", "2", "100.05", "21.01"})
	}

	return header, contents
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
