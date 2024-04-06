package userdtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"

type ResponseUsers struct {
	Data []ResponseUser `json:"data"`
	Meta dtos.Meta      `json:"meta"`
}

type ResponseUser struct {
	Id          uint64
	User        string
	Nombre      string
	PersonaId   uint64
	Activo      bool
	UserSistema *[]ResponseUserSistema `json:"UserSistema,omitempty"`
	Sistemas    *[]ResponseSistema     `json:"sistema,omitempty"`
	Cliente     *ResponseUserCliente   `json:"cliente,omitempty"`
}

type ResponseUserCliente struct {
	Id          uint64
	Cuit        string
	RazonSocial string
}

type ResponseUserSistema struct {
	ID         uint
	UsersID    uint64
	SistemasID uint64
	Activo     bool
}
