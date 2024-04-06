package entities

import "gorm.io/gorm"

type Prismamotivosajuste struct {
	gorm.Model
	ExternalId    string
	Motivoajustes string
}
