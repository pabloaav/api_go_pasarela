package filtros

type RendicionesClienteFiltro struct {
	FechaInicio string
	FechaFin    string
	ClienteId   int `json:"cliente_id"`
}
