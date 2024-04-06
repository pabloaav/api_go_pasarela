package reportedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseMovLiquidacion struct {
	Clientes       Cliente
	Cuenta         string
	FechaRendicion string // fecha cobro
	FechaProceso   time.Time
	Movimientos    []entities.Movimiento
	NroLiquidacion int
}

type Cliente struct {
	Id          uint     `json:"id"`      // Id cliente
	Cliente     string   `json:"cliente"` // nombre Cliente abreviado
	RazonSocial string   `json:"razon_social"`
	Domicilio   string   `json:"domicilio"`
	Respinsc    string   `json:"responsable_inscripto"`
	Iibbing     string   `json:"ingreso_bruto"`
	Email       []string `json:"email"`
	Cuit        string   `json:"cuit"`
	EnviarEmail bool
	EnviarPdf   bool
}

type MovLiquidar struct {
	Idpg          []uint `json:"idpg"` // Id cliente
	Idcliente     uint   `json:"idcliente"`
	Lote          int    `json:"lote"` // nombre Cliente abreviado
	Fechalote     string `json:"fecha_lote"`
	Cliente       string `json:"cliente"`
	NombreReporte string `json:"nombre_reporte"`
}

type ResultMovLiquidacion struct {
	Cabeceras      Cliente
	NroLiquidacion string
	Cuenta         string
	FechaRendicion string
	FechaProceso   string
	NombreArchivo  string
	MedioPagoItems MedioMovItems
	Totales        TotalesMovLiquidar
}

type MedioMovItems struct {
	MedioMovCredit
	MedioMovDebit
	MedioMovDebin
	MedioMovOffline
}

type MedioMovCredit struct {
	Detalle              []DetalleMov
	CantidaTotaldBoletas string
	TotalCobrado         string
	TotalaRendir         string
}
type MedioMovDebit struct {
	Detalle              []DetalleMov
	CantidaTotaldBoletas string
	TotalCobrado         string
	TotalaRendir         string
}

type MedioMovDebin struct {
	Detalle              []DetalleMov
	CantidaTotaldBoletas string
	TotalCobrado         string
	TotalaRendir         string
}

type MedioMovOffline struct {
	Detalle              []DetalleMov
	CantidaTotaldBoletas string
	TotalCobrado         string
	TotalaRendir         string
}

type DetalleMov struct {
	Cuenta            string
	Referencia        string
	FechaCobro        string
	ImporteCobrado    string
	ImporteADepositar string
	CantidadBoletas   string
	Comision          string
	Iva               string
}
type TotalesMovLiquidar struct {
	ImporteCobrado       string
	ImporteADepositar    string
	CantidadTotalBoletas string
	ComisionTotal        string
	IvaTotal             string
}

// func ToEntity(request PagLotes) (response []entities.Pagolotes) {

// 	for _, lot := range request.Idpg {
// 		response = append(response, entities.Pagolotes{
// 			PagosID:    uint64(lot),
// 			ClientesID: uint64(request.Idcliente),
// 			Lote:       int64(request.Lote),
// 			FechaEnvio: request.Fechalote,
// 		})
// 	}
// 	return
// }
