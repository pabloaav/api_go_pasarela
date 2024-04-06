package administraciondtos


type ResponseSoporteCreate struct {
	Id 				uint						`json:"id"`
}
type ResponseSoporteRead struct {
	Id 				uint						`json:"id"`
	Nombre   		string                		`json:"nombre"`
	Email    		string                		`json:"email"`
	Visto    		bool                  		`json:"visto"`
	Estado   		EnumEstadoSoporte     		`json:"estado"`
	Consulta 		string                		`json:"consulta"`
	Abierto  		bool                  		`json:"abierto"`
	File     		string						`json:"archivo"`
	FechaCreacion 	string						`json:"fechaCreacion"`
	Respuestas 		[]ReponseSoporteRespuesta	`json:"respuestas"`

}
type ReponseSoporteRespuesta struct {
	Id uint
	Respuesta string
	Visto bool
	FechaCreacion string
}
