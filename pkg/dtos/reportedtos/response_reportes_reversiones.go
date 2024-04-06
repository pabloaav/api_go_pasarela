package reportedtos

// Seccion Para Reportes visuales y excel agrupados por fecha

type ResponseReversionesClientes struct {
	CantidadRegistros   int
	Total               float64
	DetallesReversiones []DetallesReversiones
}

type DetallesReversiones struct {
	Fecha               string
	Nombre              string
	CantidadOperaciones uint
	TotalMonto          float64
	Reversiones         []Reversiones
}
