package reportedtos

type EstructuraRegistrosBatch struct {
	NombreVariable string
	Cantidad       int
}

func (px *EstructuraRegistrosBatch) CabeceraArchivo() []EstructuraRegistrosBatch {
	cabeceraArchivo := make([]EstructuraRegistrosBatch, 0, 6)
	cabeceraArchivo = append(cabeceraArchivo,
		EstructuraRegistrosBatch{NombreVariable: "RECORDCODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "CREATEDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "ORIGENNAME", Cantidad: 25},
		EstructuraRegistrosBatch{NombreVariable: "CLIENTNUMBER", Cantidad: 9},
		EstructuraRegistrosBatch{NombreVariable: "CLIENTNAME", Cantidad: 35},
		EstructuraRegistrosBatch{NombreVariable: "FILLER", Cantidad: 54})
	return cabeceraArchivo
}

func (px *EstructuraRegistrosBatch) CabeceraLote() []EstructuraRegistrosBatch {
	cabeceraArchivo := make([]EstructuraRegistrosBatch, 0, 5)
	cabeceraArchivo = append(cabeceraArchivo,
		EstructuraRegistrosBatch{NombreVariable: "RECORDCODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "CREATEDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "BATCHNUMBER", Cantidad: 6},
		EstructuraRegistrosBatch{NombreVariable: "DESCRIPTION", Cantidad: 35},
		EstructuraRegistrosBatch{NombreVariable: "FILLER", Cantidad: 82})
	return cabeceraArchivo
}

func (px *EstructuraRegistrosBatch) DetalleTransaccion() []EstructuraRegistrosBatch {
	cabeceraArchivo := make([]EstructuraRegistrosBatch, 0, 13)
	cabeceraArchivo = append(cabeceraArchivo,
		EstructuraRegistrosBatch{NombreVariable: "RECORDCODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "RECORDSEQUENCE", Cantidad: 5},
		EstructuraRegistrosBatch{NombreVariable: "TRANSACTIONCODE", Cantidad: 2},
		EstructuraRegistrosBatch{NombreVariable: "WORKDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "TRANSFERDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "ACCOUNTNUMBER", Cantidad: 21},
		EstructuraRegistrosBatch{NombreVariable: "CURRENCYCODE", Cantidad: 3},
		EstructuraRegistrosBatch{NombreVariable: "AMOUNT", Cantidad: 14},
		EstructuraRegistrosBatch{NombreVariable: "TERMINALID", Cantidad: 6},
		EstructuraRegistrosBatch{NombreVariable: "PAYMENTDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "PAYMENTTIME", Cantidad: 4},
		EstructuraRegistrosBatch{NombreVariable: "SEQNUMBER", Cantidad: 4},
		EstructuraRegistrosBatch{NombreVariable: "FILLER", Cantidad: 48})
	return cabeceraArchivo
}

func (px *EstructuraRegistrosBatch) DetalleDescripcion() []EstructuraRegistrosBatch {
	cabeceraArchivo := make([]EstructuraRegistrosBatch, 0, 4)
	cabeceraArchivo = append(cabeceraArchivo,
		EstructuraRegistrosBatch{NombreVariable: "RECORDCODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "BARCODE", Cantidad: 80},
		EstructuraRegistrosBatch{NombreVariable: "TYPECODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "FILLER", Cantidad: 50})
	return cabeceraArchivo
}

func (px *EstructuraRegistrosBatch) ColaArchivo() []EstructuraRegistrosBatch {
	cabeceraArchivo := make([]EstructuraRegistrosBatch, 0, 8)
	cabeceraArchivo = append(cabeceraArchivo,
		EstructuraRegistrosBatch{NombreVariable: "RECORDCODE", Cantidad: 1},
		EstructuraRegistrosBatch{NombreVariable: "CREATEDATE", Cantidad: 8},
		EstructuraRegistrosBatch{NombreVariable: "TOTALBATCHES", Cantidad: 6},
		EstructuraRegistrosBatch{NombreVariable: "FILECOUNT", Cantidad: 7},
		EstructuraRegistrosBatch{NombreVariable: "FILEAMOUNT", Cantidad: 12},
		EstructuraRegistrosBatch{NombreVariable: "FILLER", Cantidad: 38},
		EstructuraRegistrosBatch{NombreVariable: "FILECOUNT", Cantidad: 7},
		EstructuraRegistrosBatch{NombreVariable: "FILLER2", Cantidad: 53})
	return cabeceraArchivo
}
