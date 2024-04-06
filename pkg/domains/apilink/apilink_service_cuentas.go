package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
)

func (s *aplinkService) GetCuentasApiLinkService() (response []linkcuentas.GetCuentasResponse, erro error) {
	

	
	requerimientoId := s.GenerarUUid()

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}


	repRequest := linkcuentas.LinkGetCuentasRequest{
		Token: token,
		RequerimientoId: requerimientoId,
	}


	return s.remoteRepository.GetCuentasApiLink(repRequest)

}


func (s *aplinkService) CreateCuentaApiLinkService(request linkcuentas.LinkPostCuenta) (erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	requerimientoId := s.GenerarUUid()

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}


	repRequest := linkcuentas.LinkCuentasRequest{
		Token: token,
		RequerimientoId: requerimientoId,
		Request: request,
	}

	return s.remoteRepository.CreateCuentaApiLink(repRequest)
}


func (s *aplinkService) DeleteCuentaApiLinkService(request linkcuentas.LinkDeleteCuenta) (erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	requerimientoId := s.GenerarUUid()

	scopes := []linkdtos.EnumScopeLink{linkdtos.DebinRecurrencia}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}


	repRequest := linkcuentas.LinkCuentasRequest{
		Token: token,
		RequerimientoId: requerimientoId,
		Request: request,
	}

	return s.remoteRepository.DeleteCuentaApiLink(repRequest)
}