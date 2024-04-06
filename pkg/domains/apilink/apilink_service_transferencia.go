package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
)

func (s *aplinkService) CreateTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestTransferenciaCreateLink) (response linktransferencia.ResponseTransferenciaCreateLink, erro error) {

	// 1 - Valido los datos de entrada
	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	response, erro = s.remoteRepository.CreateTransferenciaApiLink(requerimientoId, request, token.AccessToken)

	return
}

func (s *aplinkService) GetTransferenciasApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink) (response linktransferencia.ResponseGetTransferenciasLink, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	response, erro = s.remoteRepository.GetTransferenciasApiLink(requerimientoId, request, token.AccessToken)

	return

}

func (s *aplinkService) GetTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink) (response linktransferencia.ResponseGetTransferenciaLink, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	response, erro = s.remoteRepository.GetTransferenciaApiLink(requerimientoId, request, token.AccessToken)

	return
}
