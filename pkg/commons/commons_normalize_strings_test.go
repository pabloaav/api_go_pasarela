package commons_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

type TableNormalizedStrings struct {
	Nombre   string
	Request  []string
	Response []string
}

func _inicializarNormalizeStrings() (table []TableNormalizedStrings) {

	table = []TableNormalizedStrings{
		{
			"Debe devolver un string vacio",
			[]string{"", "  ", "?", "*//---¨{*¿", "&%·#@=)(/)|"},
			[]string{"", "  ", " ", "          ", "           "},
		},
		{
			"Debe eliminar las tildes",
			[]string{"início", "gestión", "López", "Méndez"},
			[]string{"INICIO", "GESTION", "LOPEZ", "MENDEZ"},
		},

		{
			"Debe eliminar las tildes y caracteres especiales",
			[]string{"início.", "gestión, persona:", "López?", "Méndez;"},
			[]string{"INICIO ", "GESTION  PERSONA ", "LOPEZ ", "MENDEZ "},
		},
	}

	return
}
func TestNormalizeStrings(t *testing.T) {

	table := _inicializarNormalizeStrings()

	mockFileRepository := new(mockrepository.MockFileRepository)
	service := commons.NewCommons(mockFileRepository)

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for i, r := range v.Request {

				file, _ := service.NormalizeStrings(r)

				assert.Equal(t, v.Response[i], file)

			}

		})
	}

}
