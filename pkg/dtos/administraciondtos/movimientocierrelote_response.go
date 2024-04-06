package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

type MovimientoCierreLoteResponse struct {
	ListaPagos             []entities.Pago                       `json:"pagos,omitempty"`
	ListaPagosEstadoLogs   []entities.Pagoestadologs             `json:"pago_estado_log,omitempty"`
	ListaMovimientos       []entities.Movimiento                 `json:"moviminetos,omitempty"`
	ListaCLApiLink         []entities.Apilinkcierrelote          `json:"apilinkcierrelote,omitempty"`
	ListaCLPrisma          []entities.Prismacierrelote           `json:"prismacierrelote,omitempty"`
	ListaCLRapipago        []entities.Rapipagocierrelotedetalles `json:"rapipagocierrelote,omitempty"`
	ListaCLRapipagoHeaders []*entities.Rapipagocierrelote        `json:"rapipagocierreloteheaders,omitempty"`
	ListaPagoIntentos      []entities.Pagointento                `json:"pagointento,omitempty"`
	ListaReversiones       []entities.Reversione                 `json:"reversion,omitempty"`
}
