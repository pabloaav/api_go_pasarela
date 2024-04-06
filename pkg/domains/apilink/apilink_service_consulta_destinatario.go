package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
)

func (s *aplinkService) GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}
	//NOTE Decomentar en PRODUCCION
	// scopes := []linkdtos.EnumScopeLink{linkdtos.ConsultaDestinatario}

	// token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	// if erro != nil {
	// 	return
	// }

	//NOTE - Descomentar en DESARROLLO
	token := linkdtos.TokenLink{
		AccessToken: "eyJraWQiOiJSZWRMaW5rIiwiYWxnIjoiSFM1MTIifQ.eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiQ09OU1VMVEFfREVTVElOQVRBUklPIiwiYXVkIjoiZC5hcGkucmVkbGluay5jb20uYXIvcmVkbGluay9zYi8iLCJleHAiOjE2OTExMDAwNzMsImlhdCI6MTY5MTA2NDA3M30.6zQIvPL9amI21XDiT62CpTH84-GmQ7R2n_76ybyrJgzmaZILv0CKZ-NDzM2LOI5Sy7mneI2Dn-rah92ZCuATVA",
		Scope:       "CONSULTA_DESTINATARIO",
		Audience:    "d.api.redlink.com.ar/redlink/sb/",
		Expires_in:  "36000",
	}

	return s.remoteRepository.GetConsultaDestinatario(requerimientoId, request, token.AccessToken)

}
