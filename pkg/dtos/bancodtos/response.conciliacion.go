package bancodtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseConciliacion struct {
	Transferencias            []MovimientosTransferenciasResponse
	ListaApilink              []*entities.Apilinkcierrelote
	ListaApilinkNoAcreditados []*entities.Apilinkcierrelote
	TransferenciasConciliadas []TransferenciasConciliadasConBanco
	ListaRapipago             []*entities.Rapipagocierrelote
}

type MovimientosTransferenciasResponse struct {
	Id              uint `json:"id"`
	Match           int  `json:"match"`
	BancoExternalId int  `json:"banco_external_id"`
}

type TransferenciasConciliadasConBanco struct {
	ListaIdsTransferenciasConciliadas []uint `json:"lista_ids_transferencias_conciliadas"`
	Match                             int    `json:"match"`
	BancoExternalId                   int    `json:"banco_external_id"`
}
