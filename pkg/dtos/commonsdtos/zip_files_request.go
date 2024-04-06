package commonsdtos

import (
	"fmt"
)

type ZipFilesRequest struct {
	NombreArchivo string
	Rutas         []InfoFile
}

func (z *ZipFilesRequest) IsValid() error {
	if len(z.NombreArchivo) < 1 {
		return fmt.Errorf(ERROR_FILE_NAME)
	}
	if len(z.Rutas) < 1 {
		return fmt.Errorf(ERROR_RUTAS)
	}
	return nil
}

type InfoFile struct {
	RutaCompleta string
	NombreArchivo string
}
