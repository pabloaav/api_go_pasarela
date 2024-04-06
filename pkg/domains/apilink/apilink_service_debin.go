package apilink

import (
	"errors"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
)

/*
Crea un nuevo debin
requerimentoId es un identificador único para la operación.
*/

func (s *aplinkService) CreateDebinApiLinkService(requerimientoId string, request linkdebin.RequestDebinCreateLink) (response linkdebin.ResponseDebinCreateLink, erro error) {

	// NOTE descomentar esto en produccion
	// erro = request.IsValid()

	// if erro != nil {
	// 	return
	// }

	// //En esta version no se está implementando debin con recurrencia
	// //Por eso garantizo que no se haga peticiones con recurrencia

	// usaDebinRecurrencia, _ := strconv.ParseBool(config.USADEBINRECURRENCIA)

	// if !usaDebinRecurrencia {
	// 	request.Debin.Recurrente = false
	// }

	// scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	// token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)
	// logs.Info(token)
	// logs.Info(erro)

	// if erro != nil {
	// 	return
	// }

	// return s.remoteRepository.CreateDebinApiLink(requerimientoId, request, token.AccessToken)

	// NOTE descomentar esto en desarrollo
	erro = request.IsValid()

	if erro != nil {
		return
	}

	//En esta version no se está implementando debin con recurrencia
	//Por eso garantizo que no se haga peticiones con recurrencia

	usaDebinRecurrencia, _ := strconv.ParseBool(config.USADEBINRECURRENCIA)

	if !usaDebinRecurrencia {
		request.Debin.Recurrente = false
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}
	logs.Info(scopes)
	// token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)
	// logs.Info(token)
	// logs.Info(erro)

	// if erro != nil {
	// 	// return
	// 	logs.Info("Ocurrio un error al solicitar token de autenticacion")
	// }

	token := linkdtos.TokenLink{
		AccessToken: "eyJraWQiOiJSZWRMaW5rIiwiYWxnIjoiSFM1MTIifQ.eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiREVCSU5fUkVDVVJSRU5DSUEiLCJhdWQiOiJkLmFwaS5yZWRsaW5rLmNvbS5hci9yZWRsaW5rL3NiLyIsImV4cCI6MTY2NTExNDc2MSwiaWF0IjoxNjY1MDc4NzYxfQ.f3IthEvV7udPMXnd2GefWqsZIrVC72SeaPLDd8EDZH8sEZQBKIMs6ByN4jPTLKs9JdUDndle3mn7ZvE8HBSz7A",
		Scope:       "DEBIN_RECURRENCIA",
		Audience:    "d.api.redlink.com.ar/redlink/sb/",
		Expires_in:  "36000",
	}

	return s.remoteRepository.CreateDebinApiLink(requerimientoId, request, token.AccessToken)

}

/*
Consulta una lista de debines
requerimentoId es un identificador único para la operación.
*/
func (s *aplinkService) GetDebinesApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinesLink) (response linkdebin.ResponseGetDebinesLink, erro error) {
	// return linkdebin.ResponseGetDebinesLink{
	// 	Debines: []linkdebin.DebinesListLink{
	// 		{
	// 			Id:              "0V1JXON1R8163J7NZ64EL7",
	// 			Importe:         280000,
	// 			Estado:          "ACREDITADO",
	// 			Concepto:        "VAR",
	// 			Moneda:          "ARS",
	// 			Tipo:            "",
	// 			FechaExpiracion: time.Now(),
	// 			Devuelto:        false,
	// 			ContraCargoId:   "",
	// 			Comprador: linkdebin.CompradorDebinesListLink{
	// 				Cuit: "20326562468",
	// 			},
	// 			Vendedor: linkdebin.VendedorDebinesListLink{
	// 				Cuit: "30716550849",
	// 			},
	// 		},

	// 		{
	// 			Id:              "XJ8G7V95D01L4DM9EMPYR0",
	// 			Importe:         1000,
	// 			Estado:          "ACREDITADO",
	// 			Concepto:        "VAR",
	// 			Moneda:          "ARS",
	// 			Tipo:            "",
	// 			FechaExpiracion: time.Now(),
	// 			Devuelto:        false,
	// 			ContraCargoId:   "",
	// 			Comprador: linkdebin.CompradorDebinesListLink{
	// 				Cuit: "20326562468",
	// 			},
	// 			Vendedor: linkdebin.VendedorDebinesListLink{
	// 				Cuit: "30716550849",
	// 			},
	// 		},

	// 		// {
	// 		// 	Id:              "Z6OLMDN3VDQQ31W2E7RQ5X",
	// 		// 	Importe:         1000,
	// 		// 	Estado:          "ERROR_ACREDITACION",
	// 		// 	Concepto:        "VAR",
	// 		// 	Moneda:          "ARS",
	// 		// 	Tipo:            "",
	// 		// 	FechaExpiracion: time.Now(),
	// 		// 	Devuelto:        false,
	// 		// 	ContraCargoId:   "",
	// 		// 	Comprador: linkdebin.CompradorDebinesListLink{
	// 		// 		Cuit: "20326562468",
	// 		// 	},
	// 		// 	Vendedor: linkdebin.VendedorDebinesListLink{
	// 		// 		Cuit: "30716550849",
	// 		// 	},
	// 		// },

	// 		// {
	// 		// 	Id:              "D4RO172VP1GX7RD2KJ3QE6",
	// 		// 	Importe:         1000,
	// 		// 	Estado:          "ACREDITADO",
	// 		// 	Concepto:        "VAR",
	// 		// 	Moneda:          "ARS",
	// 		// 	Tipo:            "",
	// 		// 	FechaExpiracion: time.Now(),
	// 		// 	Devuelto:        false,
	// 		// 	ContraCargoId:   "",
	// 		// 	Comprador: linkdebin.CompradorDebinesListLink{
	// 		// 		Cuit: "20326562468",
	// 		// 	},
	// 		// 	Vendedor: linkdebin.VendedorDebinesListLink{
	// 		// 		Cuit: "30716550849",
	// 		// 	},
	// 		// },
	// 	},
	// }, nil

	/* DESCOMENTAR PARA PRODUCCION*/

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	/*CONSULTAR DEBINES A API: SE NECESITA EL ID REQUERIMIENTO , TOKEN Y REQUEST*/
	return s.remoteRepository.GetDebinesApiLink(requerimientoId, request, token.AccessToken)

}

func (s *aplinkService) GetDebinApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	return s.remoteRepository.GetDebinApiLink(requerimientoId, request, token.AccessToken)

}

func (s *aplinkService) DeleteDebinApiLinkService(requerimientoId string, request linkdebin.RequestDeleteDebinLink) (response bool, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}
	requestDebin := linkdebin.RequestGetDebinLink{}
	requestDebin.Cbu = request.Cbu
	requestDebin.Id = request.Id

	debin, erro := s.remoteRepository.GetDebinApiLink(requerimientoId, requestDebin, token.AccessToken)

	if erro != nil {
		return
	}

	if debin.Debin.Estado != linkdtos.Acreditado {
		return s.remoteRepository.DeleteDebinApiLink(requerimientoId, request, token.AccessToken)
	}

	return false, errors.New(ERROR_DELETE_DEBINES_ACREDITADOS)

}

func (s *aplinkService) GetDebinesPendientesApiLinkService(requerimientoId string, cbu string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error) {

	erro = tools.EsCbuValido(cbu, tools.ERROR_CBU)

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	return s.remoteRepository.GetDebinesPendientesApiLink(requerimientoId, cbu, token.AccessToken)

}
