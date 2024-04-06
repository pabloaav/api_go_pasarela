package cierreloterapipagodtos

type RapipagoCamposDescripcion struct {
	NombreVariable string
	Cantidad       int
	Desde          int
	Hasta          int
}

func (rcd *RapipagoCamposDescripcion) RapipagoDescripcionHeader() (rapipagoHeader []RapipagoCamposDescripcion) {
	rapipagoHeader = append(rapipagoHeader,
		RapipagoCamposDescripcion{NombreVariable: "Id_header", Cantidad: 8, Desde: 0, Hasta: 9},
		RapipagoCamposDescripcion{NombreVariable: "Nombre de empresa", Cantidad: 20, Desde: 9, Hasta: 30},
		RapipagoCamposDescripcion{NombreVariable: "Fh_proceso", Cantidad: 8, Desde: 30, Hasta: 39},
		RapipagoCamposDescripcion{NombreVariable: "Id_archivo", Cantidad: 20, Desde: 39, Hasta: 59},
		RapipagoCamposDescripcion{NombreVariable: "Filler", Cantidad: 17, Desde: 59, Hasta: 76},
	)
	return rapipagoHeader
}

func (rcd *RapipagoCamposDescripcion) RapipagoDescripcionDetalle() (rapipagoDetalle []RapipagoCamposDescripcion) {
	rapipagoDetalle = append(rapipagoDetalle,
		RapipagoCamposDescripcion{NombreVariable: "Fecha de Cobro", Cantidad: 8, Desde: 0, Hasta: 9},
		RapipagoCamposDescripcion{NombreVariable: "Importe Cobrado", Cantidad: 15, Desde: 9, Hasta: 29},
		RapipagoCamposDescripcion{NombreVariable: "Código de barras", Cantidad: 0, Desde: 29, Hasta: 0},
	)
	return rapipagoDetalle
}

func (rcd *RapipagoCamposDescripcion) RapipagoDescripcionTrailer() (rapipagoTrailer []RapipagoCamposDescripcion) {
	rapipagoTrailer = append(rapipagoTrailer,
		RapipagoCamposDescripcion{NombreVariable: "Id_trailer", Cantidad: 8, Desde: 0, Hasta: 9},
		RapipagoCamposDescripcion{NombreVariable: "Cant_reg", Cantidad: 8, Desde: 9, Hasta: 18},
		RapipagoCamposDescripcion{NombreVariable: "Importe_tot", Cantidad: 18, Desde: 18, Hasta: 37},
		RapipagoCamposDescripcion{NombreVariable: "Filler", Cantidad: 39, Desde: 37, Hasta: 76},
	)
	return rapipagoTrailer
}
