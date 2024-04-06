package reportedtos

type RequestConciliacion struct {
	TipoFactory TipoFactory
}

type TipoFactory struct {
	TipoApilink []string
	TipoOffline []string
	TipoPrisma  []string
}

// type TipoApilink struct {
// 	ListaApilink []string
// }

// type TipoOffline struct {
// 	ListaOffline []string
// }

// type TipoPrisma struct {
// 	ListaPrisma []string
// }
