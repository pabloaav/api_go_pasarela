package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"

type ResponseClMultipago struct {
	ClMultipago []ClMultipago `json:"data"`
	Meta        dtos.Meta     `json:"meta"`
}

type ClMultipago struct {
	IdClMultipago            uint64               `json:"id_clmultipago"`
	IdArchivo                string               `json:"nombre_archivo"`
	FechaProceso             string               `json:"fecha_proceso"`
	Detalles                 uint64               `json:"detalles"`
	ImporteTotal             uint64               `json:"importe_total"`
	ImporteTotalCalculado    float64              `json:"importe_total_calculado"`
	IdBanco                  uint64               `json:"id_banco"`
	FechaAcreditacion        string               `json:"fecha_acrditacion"`
	CantidadDiasAcreditacion uint64               `json:"cant_dias_acreditacion"`
	ImporteMinimo            uint64               `json:"importe_minimo_cobrado"`
	Coeficiente              float64              `json:"coeficiente"`
	EnObservacion            bool                 `json:"en_observacion"`
	DiferenciaBanco          float64              `json:"diferencia_banco"`
	FechaCreacion            string               `json:"fecha_creacion"`
	PagoActualizado          bool                 `json:"pago_actualizado"`
	ClMultipagoDetalle       []ClMultipagoDetalle `json:"detalles_cierre_lote"`
}

type ClMultipagoDetalle struct {
	FechaCobro       string  `json:"fecha_cobro"`
	ImporteCobrado   uint64  `json:"importe_cobrado"`
	ImporteCalculado float64 `json:"importe_calculado"`
	CodigoBarras     string  `json:"codigo_barra"`
	Conciliado       bool    `json:"conciliado"`
	Informado        bool    `json:"informado"`
}
