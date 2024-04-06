package userdtos

type RequestUserAutorizacion struct {
	Token         string
	Request       *RequestUser
	RequestUpdate *RequestUserUpdate
}

type RequestUser struct {
	User      string
	Nombre    string
	Password  string
	Activo    bool
	SistemaId string
	ClienteId uint64
}

type RequestUserUpdate struct {
	Id                uint64
	Nombre            string
	Password          string
	Activo            bool
	ClienteIdNuevo    uint64
	ClienteIdAnterior uint64
	SistemaId         string
}
