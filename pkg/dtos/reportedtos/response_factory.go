package reportedtos

type ResponseFactory struct {
	Pago                    string
	NroEstablecimiento      string
	NroLiquidacion          string
	FechaPresentacion       string
	FechaAcreditacion       string
	ArancelPorcentaje       float64
	RetencionIva            string
	Importemaximo           float64 // solo offline
	Importeminimo           float64 // solo offline
	ArancelPorcentajeMinimo float64 // solo offline
	ArancelPorcentajeMaximo float64 // solo offline
	ImporteNetoCobrado      float64 // falta para prisma
	Revertido               bool    // `json:"revertido"`
	Enobservacion           bool    // `json:"enobservacion"`
	Cantdias                int64   // `json:"cantdias"`
	ImporteArancel          float64 // solo prisma  tabla prismamovimientosdetalles ->arancel
	ImporteArancelIva       float64 // solo prisma  tabla prismatrdospagos   ->retencion_iva
	ImporteArancalIvaMov    float64
	ImporteCft              float64 // solo prisma  tablaprismamoviemientosdetalles ->importe_costo_financiero
}
