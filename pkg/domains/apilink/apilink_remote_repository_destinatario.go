package apilink

import (
	//NOTE Descomentar en PRODUCCION

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
)

func (r *remoteRepository) GetConsultaDestinatario(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink, token string) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {

	//NOTE Descomentar en PRODUCCION
	// base, erro := url.Parse(config.APILINKCONSULTADESTINATARIOHOST)

	// if erro != nil {
	// 	return
	// }
	// base.Path += "destinatarios"
	// params := url.Values{}
	// params.Add("cbu", request.Cbu)
	// params.Add("alias", request.Alias)
	// base.RawQuery = params.Encode()

	// req, _ := http.NewRequest("GET", base.String(), nil)

	// buildHeaderAutorizacion(req, requerimientoId, token)

	// erro = executeRequest(r, req, ERROR_GET_CONSULTA_DESTINATARIOS, &response)
	// /*
	// 	 se registra la peticion realizada a la api de apilink
	// 		-	armo el request para registrar la peticion realizada
	// 		-	registro la peticion realizada
	// */
	// peticionApiLink := dtos.RequestWebServicePeticion{
	// 	Operacion: "GetConsultaDestinatario",
	// 	Vendor:    "ApiLink",
	// }
	// err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	// if err1 != nil {
	// 	logs.Error(ERROR_CREAR_PETICION + err1.Error())
	// }

	//NOTE Descomentar en DESARROLLO
	titularesLink := []linkconsultadestinatario.TitularLink{
		{
			IdTributario: "30000000015",
			Denominacion: "Un comercio para APILINK",
		},
	}
	entidadBancaria := linkconsultadestinatario.EntidadBancariaLink{
		Nombre:          "NUEVO BANCO DEL CHACO S. A.",
		NombreAbreviado: "NBDC",
	}
	response = linkconsultadestinatario.ResponseConsultaDestinatarioLink{
		Titulares:       titularesLink,
		EntidadBancaria: entidadBancaria,
		Cbu:             "3110030211000063689055",
	}
	return

}
