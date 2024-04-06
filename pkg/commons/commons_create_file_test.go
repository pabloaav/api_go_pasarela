package commons_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

type callsMocked struct {
	Nombre         string
	Response       interface{}
	ResponseCreate *os.File
	Erro           error
}

type TableCreateFile struct {
	Nombre        string
	Request       []string
	Erro          error
	MetodosMocked []callsMocked
}

func _inicializarCreateFile() (table []TableCreateFile) {
	_listaNombresInvalidos := []string{"", " ", "Prueba", "1111", ".txt", "prueba/nueva/local.go"}
	_listaNombresValidos := []string{"prueba.txt", "prueba.doc", "prueba.csv", "prueba.pdf", "prueba.docx", "prueba/nueva/local.xlsx"}

	table = []TableCreateFile{
		{
			"Debe devolver un error si la ruta es incorrecta",
			_listaNombresInvalidos,
			fmt.Errorf(commons.ERROR_FILE_NAME),
			[]callsMocked{},
		},
		{
			"Debe devolver un error si ya encuentra la ruta especificada",
			_listaNombresValidos,
			fmt.Errorf(commons.ERROR_FILE_EXIST),
			[]callsMocked{
				{"ExisteArchivo", false, nil, nil},
			},
		},
		{
			"Debe devolver un error si no puede crear el archivo",
			_listaNombresValidos,
			fmt.Errorf(commons.ERROR_FILE_CREATE),
			[]callsMocked{
				{"ExisteArchivo", true, nil, nil},
				{"CrearArchivo", nil, &os.File{}, fmt.Errorf("commons.ERROR_FIL_CREATE")},
			},
		},
		{
			"Debe devolver un archivo en caso de succeso",
			_listaNombresValidos,
			nil,
			[]callsMocked{
				{"ExisteArchivo", true, nil, nil},
				{"CrearArchivo", nil, &os.File{}, nil},
			},
		},
	}

	return
}
func TestCreateFile(t *testing.T) {

	table := _inicializarCreateFile()

	mockFileRepository := new(mockrepository.MockFileRepository)
	service := commons.NewCommons(mockFileRepository)

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for _, r := range v.Request {

				for _, c := range v.MetodosMocked {
					switch c.Nombre {
					case "ExisteArchivo":
						mockFileRepository.On(c.Nombre, r).Return(c.Response).Once()
					case "CrearArchivo":
						mockFileRepository.On(c.Nombre, r).Return(c.ResponseCreate, c.Erro).Once()
					}

				}

				file, err := service.CreateFile(r)

				if err != nil {
					assert.Equal(t, v.Erro.Error(), err.Error())
				}
				if err == nil {
					assert.NotNil(t, file)
				}

			}

		})
	}

}




