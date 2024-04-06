package prismafake

import (
	"os"

	prismatransacciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

type TableDriverTest struct {
	TituloPrueba   string
	WantTable      string
	DataPrueba     []string
	TokenStructura prismatransacciones.StructToken
}

type TableDriverTestPayment struct {
	TituloPrueba     string
	WantTable        string
	DataPruebaString []string
	DataPruebaInt    []int64
	PaymentStructura prismatransacciones.StructPayments
}

type TableDriverDirArchivos struct {
	Path  string
	Files []os.FileInfo
}
