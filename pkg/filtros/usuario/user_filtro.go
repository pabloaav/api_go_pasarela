package filtros

type UserFiltroAutenticacion struct {
	Token string
	User  UserFiltro
}

type UserFiltro struct {
	Paginacion
	Id                uint64    `json:"Id,omitempty"`
	Ids               *[]uint64 `json:"Ids,omitempty"`
	ClienteId         uint64    `json:"ClienteId,omitempty"`
	SistemaId         string
	Email             string
	Nombre            string
	Activo            string
	CargarSistema     bool
	CargarUserSistema bool
}
