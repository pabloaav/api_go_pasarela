package cierrelotedtos

type PrismaLogArchivoResponse struct {
	NombreArchivo  string `json:"nombre_archivo"`
	ArchivoLeido   bool   `json:"archivo_leido"`
	ArchivoMovido  bool   `json:"archivo_movido"`
	LoteInsert     bool   `json:"lote_insert"`
	ErrorProducido string `json:"error_producido"`
}
