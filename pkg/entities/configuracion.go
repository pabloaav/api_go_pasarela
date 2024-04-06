package entities

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"gorm.io/gorm"
)

type Configuracione struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Valor       string `json:"valor"`
}

func (c *Configuracione) IsValid() error {

	if tools.EsStringVacio(c.Nombre) {
		return fmt.Errorf("el campo nombre es obligatorio")
	}
	if tools.EsStringVacio(c.Valor) {
		return fmt.Errorf("el campo valor es obligatorio")
	}

	c.Nombre = strings.ToUpper(c.Nombre)
	c.Nombre = strings.TrimSpace(c.Nombre)

	return nil
}
