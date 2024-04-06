package administraciondtos

type RequestPlanCuotas struct {
	InstalmentsId string `json:"instalmentsId"`
	VigenciaDesde string `json:"vigenciaDesde"`
	RutaFile      string `json:"ruta_file"`
}
