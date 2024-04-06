package cierrelotedtos

type ResponseTablasRelacionadas struct {
	ListaCodigosRechazados []ResponseCodigoRechazos    `json:"lista_codigos_rechazados"`
	ListaVisaContracargo   []ResponseVisaContracargo   `json:"lista_visa_contracargo"`
	ListaMotivosAjustes    []ResponseMotivosAjustes    `json:"lista_motivos_ajustes"`
	ListaOperaciones       []ResponseOperaciones       `json:"lista_operaciones"`
	ListaMasterContracargo []ResponseMasterContracargo `json:"lista_master_contracargo"`
}
