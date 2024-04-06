package administraciondtos

type ImpuestoResponse struct {
	Impuesto   string  `json:"impuesto"`
	Porcentaje float64 `json:"porcentaje"`
	Tipo       string  `json:"tipo"`
}
