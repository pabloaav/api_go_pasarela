package commons

import (
	"os"
)

type FileRepository interface {
	ExisteArchivo(nombre string) bool
	CrearArchivo(ruta string) (response *os.File, erro error)
	EliminarArchivo(ruta string) (erro error)
	AbrirArchivo(ruta string) (file *os.File, erro error)
}

type filerepository struct {
	File *os.File
}

func NewFileRepository(file *os.File) FileRepository {
	return &filerepository{
		File: file,
	}
}

func (r *filerepository) ExisteArchivo(nombre string) bool {
	_, erro := os.Stat(nombre)
	return os.IsNotExist(erro)
}

func (r *filerepository) CrearArchivo(ruta string) (response *os.File, erro error) {
	response, erro = os.Create(ruta)
	return
}

func (r *filerepository) EliminarArchivo(ruta string) (erro error) {
	erro = os.Remove(ruta)
	return
}

func (r *filerepository) AbrirArchivo(ruta string) (file *os.File, erro error) {
	file, erro = os.Open(ruta)
	return
}
