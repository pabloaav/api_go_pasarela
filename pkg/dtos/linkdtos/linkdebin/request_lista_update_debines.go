package linkdebin

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type RequestListaUpdateDebines struct {
	DebinId              []uint64
	Debines              []*entities.Apilinkcierrelote
	DebinesNoAcreditados []*entities.Apilinkcierrelote
}
