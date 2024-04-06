package administraciondtos

type ArchivoResponse struct {
	NombreArchivo  string `json:"nombre_archivo"`
	ErrorProducido string `json:"error_producido"`
}
