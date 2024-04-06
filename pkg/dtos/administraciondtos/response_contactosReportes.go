package administraciondtos

type ResponseContactosReportes struct {
	ClienteID    uint   `json:"cliente_id"`
	ClienteEmail string `json:"cliente_email"`
}
type ResponseGetContactosReportes struct {
	EmailsContacto []ResponseContactosReportes
}