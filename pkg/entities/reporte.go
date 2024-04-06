package entities

import (
	"gorm.io/gorm"
)

type Reporte struct {
	gorm.Model
	Cliente        string
	Tiporeporte    string
	Totalcobrado   string
	Totalrendido   string
	Fechacobranza  string
	Fecharendicion string
	Nro_reporte    uint
	Reportedetalle []Reportedetalle `gorm:"foreignkey:ReportesId"`
	TotalRetencionGanancias string
	TotalRetencionIva       string
	TotalRetencionIibb      string
	TotalRetenido           string // totalizador de las retenciones por reporte
}
