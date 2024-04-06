package cierrelotedtos

type CampoTrailer struct {
	NombreVariable string
	Cantidad       int
	Desde          int
	Hasta          int
}

func (ct *CampoTrailer) DescripcionCampos() []CampoTrailer {
	listaCampos := make([]CampoTrailer, 0, 11)
	listaCampos = append(listaCampos,
		CampoTrailer{NombreVariable: "TIPOREGISTRO", Cantidad: 1, Desde: 0, Hasta: 1},
		CampoTrailer{NombreVariable: "CANTIDADREGISTROS", Cantidad: 10, Desde: 1, Hasta: 11},
		CampoTrailer{NombreVariable: "IDMEDIOPAGO", Cantidad: 3, Desde: 11, Hasta: 14},
		CampoTrailer{NombreVariable: "IDLOTE", Cantidad: 3, Desde: 14, Hasta: 17},
		CampoTrailer{NombreVariable: "CANTCOMPRAS", Cantidad: 4, Desde: 17, Hasta: 21},
		CampoTrailer{NombreVariable: "MONTOCOMPRAS", Cantidad: 12, Desde: 21, Hasta: 33},
		CampoTrailer{NombreVariable: "CANTDEVUELTAS", Cantidad: 4, Desde: 33, Hasta: 37},
		CampoTrailer{NombreVariable: "MONTODEVUELTAS", Cantidad: 12, Desde: 37, Hasta: 49},
		CampoTrailer{NombreVariable: "CANTANULADAS", Cantidad: 4, Desde: 49, Hasta: 53},
		CampoTrailer{NombreVariable: "MONTOANULADAS", Cantidad: 12, Desde: 53, Hasta: 65},
		CampoTrailer{NombreVariable: "FILLER", Cantidad: 35, Desde: 65, Hasta: 100},
	)
	return listaCampos
}
