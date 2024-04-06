package reportedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagosLiquidacion struct {
	Clientes     Clientes
	FechaCobro   time.Time // fecha cobro
	FechaProceso time.Time
	Pagos        []entities.Pagointento
	// PagLiq       PagoLiquidar
}

type Clientes struct {
	Id          uint     `json:"id"`      // Id cliente
	Cliente     string   `json:"cliente"` // nombre Cliente abreviado
	RazonSocial string   `json:"razon_social"`
	Email       []string `json:"email"`
	Cuit        string   `json:"cuit"`
}

type PagoLiquidar struct {
	Idpg          []uint `json:"idpg"` // Id cliente
	Idcliente     uint   `json:"idcliente"`
	Lote          int    `json:"lote"` // nombre Cliente abreviado
	Fechalote     string `json:"fecha_lote"`
	Cliente       string `json:"cliente"`
	NombreReporte string `json:"nombre_reporte"`
}

type ResultPagosLiquidacion struct {
	Cabeceras      Clientes
	FechaCobro     string
	FechaProceso   string
	NombreArchivo  string
	MedioPagoItems MedioPagoItems
	Totales        TotalesALiquidar
}

type MedioPagoItems struct {
	MedioPagoCredit
	MedioPagoDebit
	MedioPagoDebin
}

type MedioPagoCredit struct {
	ImporteCobrado    string
	ImporteADepositar string
	CantidadBoletas   string
	Comision          string
	Iva               string
}
type MedioPagoDebit struct {
	ImporteCobrado    string
	ImporteADepositar string
	CantidadBoletas   string
	Comision          string
	Iva               string
}

type MedioPagoDebin struct {
	ImporteCobrado    string
	ImporteADepositar string
	CantidadBoletas   string
	Comision          string
	Iva               string
}

type TotalesALiquidar struct {
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
