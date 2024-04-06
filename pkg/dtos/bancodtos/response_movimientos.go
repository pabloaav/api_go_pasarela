package bancodtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"

type ResponseMovimientos struct {
	Movimientos []Movimientos `json:"data"`
	Meta        dtos.Meta     `json:"meta"`
}

type Movimientos struct {
	Id             uint   `json:"id"`
	NombreArchivo  string `json:"nombre_archivo"`
	Subcuenta      string `json:"subcuenta"`
	Referencia     string `json:"referencia"`
	DebinId        string `json:"debin_id"`
	DbCr           string `json:"db_cr"`
	Importe        uint64 `json:"importe"`
	Fecha          string `json:"fecha_acreditacion"`
	Conciliado     bool   `json:"conciliado"`
	TipoMovimiento uint64 `json:"tipo_movimiento"`
	Observacion    string `json:"observacion"`
	FechaCreacion  string `json:"fecha_creacion"`
}
