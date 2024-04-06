package prisma

import (
	"errors"

	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

/// CheckService permite chequear si el servicio de prisma esta funcionando
/// devuelve un valor booleano true s i esta en servicio o false si no esta en servicio
func (s *service) CheckService() (estado bool, err error) {
	_, err = s.remoteRepository.GetHealthcheck()
	if err != nil {
		estado = false
		return estado, nil
	}
	estado = true
	return estado, nil
}

func (s *service) SolicitarToken(request prismadtos.StructToken) (response interface{}, erro error) {
	//fmt.Printf("%v", request.TypePay)
	var opcion string = string(request.TypePay)
	switch opcion {
	case "simple":
		var objetoRequest prismadtos.Card = request.Card
		err := objetoRequest.Validar()
		if err != nil {
			return nil, errors.New(err.Error())
		}
		resultado, err := s.remoteRepository.PostSolicitudTokenPago(&objetoRequest)
		if err != nil {
			erro = err
			return nil, erro
		}
		response = *resultado
	//	return response, nil

	// case "offline":
	// 	var objetoRequest prismadtos.OfflineTokenRequest = request.DataOffline
	// 	err := objetoRequest.ValidarSolicitudTokenOfflineRequest()
	// 	if err != nil {
	// 		return nil, errors.New(err.Error())
	// 	}
	// 	resultado, err := s.remoteRepository.PostSolicitarTokenOffLine(&objetoRequest)
	// 	if err != nil {
	// 		erro = err
	// 		return nil, erro
	// 	}
	// 	response = *resultado
	default:
		return response, errors.New(ERRR_TIPO_PAGO)
	}
	return response, nil
}

func (s *service) Payments(request prismadtos.StructPayments) (response interface{}, erro error) {
	//fmt.Printf("%v", request.TypePay)
	var opcion string = string(request.TypePay)
	switch opcion {
	case "simple":
		var objetoRequest prismadtos.PaymentsSimpleRequest = request.PagoSimple
		err := objetoRequest.ValidarProcesoPagoRequest()
		if err != nil {
			return nil, errors.New(err.Error())
		}
		resultado, err := s.remoteRepository.PostEjecutarPago(&objetoRequest)
		if err != nil {
			erro = err
			return nil, erro
		}
		response = *resultado

	// case "offline":
	// 	var objetoRequest prismadtos.PaymentsOfflineRequest = request.PagoOffline
	// 	err := objetoRequest.Validar()
	// 	if err != nil {
	// 		return nil, errors.New(err.Error())
	// 	}
	// 	resultado, err := s.remoteRepository.PostEjecutarPagoOffLine(&objetoRequest)
	// 	if err != nil {
	// 		erro = err
	// 		return nil, erro
	// 	}
	// 	response = *resultado
	default:
		return response, errors.New(ERRR_TIPO_PAGO)
	}

	return response, nil
}
