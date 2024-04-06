package bancodtos

type ResponseMovimientosProceso struct {
	ArchivosProcesados int    `json:"archivos_procesados"`
	Error              string `json:"log_error"`
}
