package entities

import "gorm.io/gorm"

type Prismaerroresexterno struct {
	gorm.Model
	ExternalId            uint64
	Descripcion           string
	DescripcionmsgUsuario string
}
