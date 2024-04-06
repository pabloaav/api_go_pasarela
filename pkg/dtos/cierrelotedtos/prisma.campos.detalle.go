package cierrelotedtos

type CampoDetalle struct {
	NombreVariable string
	Cantidad       int
	Desde          int
	Hasta          int
}

func (ct *CampoDetalle) DescripcionCampos() []CampoDetalle {
	listaCampos := make([]CampoDetalle, 0, 11)
	listaCampos = append(listaCampos,
		CampoDetalle{NombreVariable: "TIPOREGISTRO", Cantidad: 1, Desde: 0, Hasta: 1},
		CampoDetalle{NombreVariable: "IDTRANSACCIONSITE", Cantidad: 15, Desde: 1, Hasta: 16},
		CampoDetalle{NombreVariable: "IDMEDIOPAGO", Cantidad: 3, Desde: 16, Hasta: 19},
		CampoDetalle{NombreVariable: "NROTARJETACOMPLETO", Cantidad: 16, Desde: 19, Hasta: 39},
		CampoDetalle{NombreVariable: "TIPOOPERACION", Cantidad: 1, Desde: 39, Hasta: 40},
		CampoDetalle{NombreVariable: "FECHA", Cantidad: 8, Desde: 40, Hasta: 48},
		CampoDetalle{NombreVariable: "MONTO", Cantidad: 12, Desde: 48, Hasta: 60},
		CampoDetalle{NombreVariable: "CODAUT", Cantidad: 6, Desde: 60, Hasta: 66},
		CampoDetalle{NombreVariable: "NROTICKET", Cantidad: 6, Desde: 66, Hasta: 72},
		CampoDetalle{NombreVariable: "IDSITE", Cantidad: 15, Desde: 72, Hasta: 87},
		CampoDetalle{NombreVariable: "IDLOTE", Cantidad: 3, Desde: 87, Hasta: 90},
		CampoDetalle{NombreVariable: "CUOTAS", Cantidad: 3, Desde: 90, Hasta: 93},
		CampoDetalle{NombreVariable: "FECHACIERRE", Cantidad: 8, Desde: 93, Hasta: 101},
		CampoDetalle{NombreVariable: "NROESTABLECIMIENTO", Cantidad: 30, Desde: 101, Hasta: 131},
		CampoDetalle{NombreVariable: "IDCLIENTE", Cantidad: 40, Desde: 131, Hasta: 171},
		CampoDetalle{NombreVariable: "FILLER", Cantidad: 19, Desde: 171, Hasta: 190},
	)
	return listaCampos
}
