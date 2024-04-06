package reportedtos

type ResponseLogs struct {
	FechaComienzo string               `json:"fecha_comienzo"`
	FechaFin      string               `json:"fecha_fin"`
	TotalLogs     int                  `json:"total_logs"`
	LastPage      int                  `json:"last_page"`
	Data          []ResponseDetalleLog `json:"data"`
}

type ResponseDetalleLog struct {
	Mensaje       string `json:"mensaje"`
	Fecha         string `json:"fecha"`
	Funcionalidad string `json:"funcionalidad"`
	Tipo          string `json:"tipo"`
}
