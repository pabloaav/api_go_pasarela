package reportedtos

type ResponseNotificaciones struct {
	FechaComienzo       string                          `json:"fecha_comienzo"`
	FechaFin            string                          `json:"fecha_fin"`
	TotalNotificaciones int                             `json:"total_notificaciones"`
	LastPage            int                             `json:"last_page"`
	Data                []ResponseDetalleNotificaciones `json:"data"`
}

type ResponseDetalleNotificaciones struct {
	Descripcion string `json:"descripcion"`
	Fecha       string `json:"fecha"`
	Tipo        string `json:"tipo"`
}
