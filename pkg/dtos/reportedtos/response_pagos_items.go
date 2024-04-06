package reportedtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponsePagosItems struct {
	Clientes ClientesResponse
	Fecha    string
	Pagos    []entities.Pago
	PagLotes PagLotes
}

type ClientesResponse struct {
	Id          uint   `json:"id"`      // Id cliente
	Cliente     string `json:"cliente"` // nombre Cliente abreviado
	RazonSocial string `json:"razon_social"`
	Email       string `json:"email"`
}

type MovLotes struct {
	Idmov         []uint `json:"idmov"` // Id cliente
	Idcliente     uint   `json:"idcliente"`
	Lote          int    `json:"lote"` // nombre Cliente abreviado
	Fechalote     string `json:"fecha_lote"`
	Cliente       string `json:"cliente"`
	NombreReporte string `json:"nombre_reporte"`
}

type PagLotes struct {
	Idpg          []uint `json:"idpg"` // Id cliente
	Idcliente     uint   `json:"idcliente"`
	Lote          int    `json:"lote"` // nombre Cliente abreviado
	Fechalote     string `json:"fecha_lote"`
	Cliente       string `json:"cliente"`
	NombreReporte string `json:"nombre_reporte"`
}

type ResultPagosItems struct {
	PagLotes        PagLotes
	CabeceraArchivo CabeceraArchivo
	CabeceraLote    CabeceraLote
	// DetalleTransaccion []DetalleTransaccion
	// DetalleDescripcion []DetalleDescripcion
	ResultItems []ResultItems
	ColaArchivo ColaArchivo
}

type ResultItems struct {
	DetalleTransaccion DetalleTransaccion
	DetalleDescripcion DetalleDescripcion
}

type CabeceraArchivo struct { // longitud campo
	RecordCode   string `json:"record_code"`   // 1  codigo de registro inicia con 1
	CreateDate   string `json:"create_date"`   // 8  fecha de creacion
	OrigenName   string `json:"origen_name"`   // 25  origen del archivo
	ClientNumber string `json:"client_number"` // 9 este campo no es obligatorio
	ClientName   string `json:"client_name"`   // 35 este campo no es obligatorio
	Filler       string `json:"filler"`        // 50 este campo no es obligatorio
}

type CabeceraLote struct {
	RecordCodeLote string `json:"record_code_lote"`  // 1 codigo de registro  inicia con 3
	CreateDateLote string `json:"create_date_lote"`  // 8  fecha de creacion
	BatchNumber    string `json:"batch_number_lote"` // 6   numero de lote enviado
	Description    string `json:"description"`       // 35 este campo no es obligatorio
	Filler         string `json:"filler"`            // 78 este campo no es obligatorio
}

type DetalleTransaccion struct {
	RecordCodeTransaccion string `json:"record_code_lote"` // 1  inicia con 5
	RecordSequence        string `json:"record_sequence"`  // 5 este campo no es obligatorio
	TransactionCode       string `json:"transaction_code"` // 2 este campo no es obligatorio
	WorkDate              string `json:"work_date"`        // 8 este campo no es obligatorio
	TransferDate          string `json:"transfer_date"`    // 8 este campo no es obligatorio
	AccountNumber         string `json:"account_number"`   // 21  nro de referencia - nro de usuario : campo descripcion items del checkout
	CurrencyCode          string `json:"currency_code"`    // 3 este campo no es obligatorio
	Amount                string `json:"amount"`           // 10  importe cobrado
	TerminalId            string `json:"terminal_id"`      // 6 este campo no es obligatorio
	PaymentDate           string `json:"payment_date"`     // 6  fecha que se efectuo la transaccion
	PaymentTime           string `json:"payment_time"`     // 6  hora que se efectuo la transaccion
	SeqNumber             string `json:"seq_number"`       // 8 este campo no es obligatorio
	Filler                string `json:"filler"`           // 48 este campo no es obligatorio
}

type DetalleDescripcion struct {
	RecordCodeLote string `json:"record_code_lote"` // 1   codigo de registro inicia con 6
	BarCode        string `json:"bar_code"`         // 80   codigo de barras de la transaccion
	TypeCode       string `json:"type_code"`        // 1 este campo no es obligatorio
	Filler         string `json:"filler"`           // 46 este campo no es obligatorio
}

type ColaArchivo struct {
	RecordCodeCola    string `json:"record_code_cola"`    // 1   codigo de registro
	CreateDateCola    string `json:"create_date_cola"`    // 8   fecha de creacion del archivo
	TotalBatches      string `json:"total_batches"`       // 6 este campo no es obligatorio
	FilePaymentCount  string `json:"file_payment_count"`  // 7   cantidad total de transacciones del archivo
	FilePaymentAmount string `json:"file_payment_amount"` // 12  importe total cobrado del archivo
	Filler            string `json:"filler"`              // 38 este campo no es obligatorio
	FileCount         string `json:"file_count"`          // 7 este campo no es obligatorio
	Filler2           string `json:"filler2"`             // 49 este campo no es obligatorio
}

// VALIDAR ESTRUCTURA DE REGISTROS CREADA
func (cabArchivo *CabeceraArchivo) ValidarCabeceraArchivo(estructuraReg *EstructuraRegistrosBatch) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro cabecera de archivo")
	if len(cabArchivo.RecordCode) != estructuraReg.CabeceraArchivo()[0].Cantidad {
		return err
	}
	if len(cabArchivo.CreateDate) != estructuraReg.CabeceraArchivo()[1].Cantidad {
		return err
	}
	if len(cabArchivo.OrigenName) != estructuraReg.CabeceraArchivo()[2].Cantidad {
		return err
	}
	if len(cabArchivo.ClientNumber) != estructuraReg.CabeceraArchivo()[3].Cantidad {
		return err
	}
	if len(cabArchivo.ClientName) != estructuraReg.CabeceraArchivo()[4].Cantidad {
		return err
	}
	if len(cabArchivo.Filler) != estructuraReg.CabeceraArchivo()[5].Cantidad {
		return err
	}
	return
}

func (cabLote *CabeceraLote) ValidarCabeceraLote(estructuraReg *EstructuraRegistrosBatch) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro cabecera de lote")
	if len(cabLote.RecordCodeLote) != estructuraReg.CabeceraLote()[0].Cantidad {
		return err
	}
	if len(cabLote.CreateDateLote) != estructuraReg.CabeceraLote()[1].Cantidad {
		return err
	}
	if len(cabLote.BatchNumber) != estructuraReg.CabeceraLote()[2].Cantidad {
		return err
	}
	if len(cabLote.Description) != estructuraReg.CabeceraLote()[3].Cantidad {
		return err
	}
	if len(cabLote.Filler) != estructuraReg.CabeceraLote()[4].Cantidad {
		return err
	}
	return
}

func (detTrans *DetalleTransaccion) ValidarDetalleTransaccion(estructuraReg *EstructuraRegistrosBatch) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro detalle transacción")
	if len(detTrans.RecordCodeTransaccion) != estructuraReg.DetalleTransaccion()[0].Cantidad {
		return err
	}
	if len(detTrans.RecordSequence) != estructuraReg.DetalleTransaccion()[1].Cantidad {
		return err
	}
	if len(detTrans.TransactionCode) != estructuraReg.DetalleTransaccion()[2].Cantidad {
		return err
	}
	if len(detTrans.WorkDate) != estructuraReg.DetalleTransaccion()[3].Cantidad {
		return err
	}

	if len(detTrans.TransferDate) != estructuraReg.DetalleTransaccion()[4].Cantidad {
		return err
	}
	if len(detTrans.AccountNumber) != estructuraReg.DetalleTransaccion()[5].Cantidad {
		return err
	}
	if len(detTrans.CurrencyCode) != estructuraReg.DetalleTransaccion()[6].Cantidad {
		return err
	}
	if len(detTrans.Amount) != estructuraReg.DetalleTransaccion()[7].Cantidad {
		return err
	}

	if len(detTrans.TerminalId) != estructuraReg.DetalleTransaccion()[8].Cantidad {
		return err
	}
	if len(detTrans.PaymentDate) != estructuraReg.DetalleTransaccion()[9].Cantidad {
		return err
	}
	if len(detTrans.PaymentTime) != estructuraReg.DetalleTransaccion()[10].Cantidad {
		return err
	}
	if len(detTrans.SeqNumber) != estructuraReg.DetalleTransaccion()[11].Cantidad {
		return err
	}

	if len(detTrans.Filler) != estructuraReg.DetalleTransaccion()[12].Cantidad {
		return err
	}
	return
}

func (cabDesc *DetalleDescripcion) ValidarDetalleDescripcion(estructuraReg *EstructuraRegistrosBatch) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro detalle descripción")
	if len(cabDesc.RecordCodeLote) != estructuraReg.DetalleDescripcion()[0].Cantidad {
		return err
	}
	if len(cabDesc.BarCode) != estructuraReg.DetalleDescripcion()[1].Cantidad {
		return err
	}
	if len(cabDesc.TypeCode) != estructuraReg.DetalleDescripcion()[2].Cantidad {
		return err
	}
	if len(cabDesc.Filler) != estructuraReg.DetalleDescripcion()[3].Cantidad {
		return err
	}
	return
}

func (colaArchivo *ColaArchivo) ValidarColaArchivo(estructuraReg *EstructuraRegistrosBatch) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro cabecera de lote")
	if len(colaArchivo.RecordCodeCola) != estructuraReg.ColaArchivo()[0].Cantidad {
		return err
	}
	if len(colaArchivo.CreateDateCola) != estructuraReg.ColaArchivo()[1].Cantidad {
		return err
	}
	if len(colaArchivo.TotalBatches) != estructuraReg.ColaArchivo()[2].Cantidad {
		return err
	}
	if len(colaArchivo.FilePaymentCount) != estructuraReg.ColaArchivo()[3].Cantidad {
		return err
	}
	if len(colaArchivo.FilePaymentAmount) != estructuraReg.ColaArchivo()[4].Cantidad {
		return err
	}
	if len(colaArchivo.Filler) != estructuraReg.ColaArchivo()[5].Cantidad {
		return err
	}
	if len(colaArchivo.FileCount) != estructuraReg.ColaArchivo()[6].Cantidad {
		return err
	}
	if len(colaArchivo.Filler2) != estructuraReg.ColaArchivo()[7].Cantidad {
		return err
	}
	return
}

func ToEntity(request PagLotes) (response []entities.Pagolotes) {

	for _, lot := range request.Idpg {
		response = append(response, entities.Pagolotes{
			PagosID:    uint64(lot),
			ClientesID: uint64(request.Idcliente),
			Lote:       int64(request.Lote),
			FechaEnvio: request.Fechalote,
		})
	}
	return
}
