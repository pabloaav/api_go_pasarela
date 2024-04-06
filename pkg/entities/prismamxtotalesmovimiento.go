package entities

import "gorm.io/gorm"

type Prismamxtotalesmovimiento struct {
	gorm.Model
	Empresa            string
	Fechapres          string
	Tiporeg            string
	Numcom             string
	Numest             string
	Codop              string
	Tipoaplic          string
	Filler             string
	Fechapago          string
	Libre              string
	ImporteTotal       string
	SignoImporteTotal  string
	Filler1            string
	McaPex             string
	Filler2            string
	Aster              string
	Nombrearchivo      string
	MovimientosDetalle []Prismamxdetallemovimiento `gorm:"foreignkey:PrismamxtotalesmoviminetosId"`
}
