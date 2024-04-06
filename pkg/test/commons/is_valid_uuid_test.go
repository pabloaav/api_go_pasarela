package commons_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUuid(t *testing.T) {
	mockFileRepository := new(mockrepository.MockFileRepository)
	RequerimientosInvalidos := []string{"", "  ", "4%rrrreeee", "0b949e69-12d3-4a39-b87e-8c632d34c2zz", "00-33-RRRRRRRRRRR440000555555"}
	RequerimientosValidos := []string{
		"4b241667-9a94-41e1-95fe-63681a971ecd", "ecff59f5-c227-459b-ad62-56dc786fba65",
		"3d52d242-9f52-49bd-a230-6eb7dccf0aea",
		"8bee9e67-8c06-4f10-8ecf-05b3877db6ad",
	}

	t.Run("Debe retornar un error si el requerimiento es inválido", func(t *testing.T) {

		for _, r := range RequerimientosInvalidos {

			service := commons.NewCommons(mockFileRepository)

			res, err := service.IsValidUUID(r)

			assert.Equal(t, commons.ERROR_UUID, err.Error())
			assert.Equal(t, false, res)
		}

	})

	t.Run("Debe retornar nil si el uuid es valido", func(t *testing.T) {

		for _, r := range RequerimientosValidos {

			service := commons.NewCommons(mockFileRepository)

			res, err := service.IsValidUUID(r)

			assert.Equal(t, nil, err)
			assert.Equal(t, true, res)
		}

	})
}
