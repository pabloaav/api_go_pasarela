package linkdtos

import (
	"fmt"
)

type ErrorApiLink struct {
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
}

func (e *ErrorApiLink) Error() string {
	return fmt.Sprintf("error - Código: %s, Descripción: %s", e.Codigo, e.Descripcion)
}
