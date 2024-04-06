package prisma

import (
	"errors"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	commonds "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

func (s *service) GetInformePago(pagoId string) (response *prismainforme.UnPagoResponse, erro error) {
	if commonds.StringIsEmpity(pagoId) {
		//return nil, erro
		return nil, errors.New(ERRR_TIPO_PAGO)
	}
	resultado, err := s.remoteRepository.GetPrismaInformarPago(pagoId)
	if err != nil {
		//return nil, erro
		return nil, errors.New(ERRR_INFO_PAGO)
	}
	if resultado == nil {
		//return nil, erro
		return nil, errors.New(ERRR_TIPO_PAGO)
	}

	response = resultado
	return response, nil
}

func (s *service) ListarPagosPorFecha(requestPago prismainforme.ListaPagosRequest) (listaPago []prismainforme.Result, erro error) {
	estadoPeticion := true
	var listaPagos []prismainforme.Result
	err := validarListaPagosRequest(requestPago)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for estadoPeticion {
		resultado, err := s.remoteRepository.ListarPagosPorFecha(&requestPago)
		if err != nil {
			//return nil, erro
			return nil, errors.New(ERRR_INFO_PAGO)
		}
		logs.Info(resultado.HasMore)
		if !resultado.HasMore {
			estadoPeticion = false
		}
		requestPago.Offset += 1

		listaPagos = append(listaPagos, resultado.Results...)
		time.Sleep(3 * time.Millisecond)
	}
	//response = resultado
	return listaPagos, nil
}

func (s *service) ListarPagosService(estadoPago int, channel string) (pagoIntentos []entities.Pagointento, erro error) {
	pagoIntento, err := s.repository.GetPagosPagosIntentosxChannel(estadoPago, channel)
	if err != nil {
		return nil, err
	}
	return pagoIntento, nil
}

/////////////////validar///////////////////////
func validarListaPagosRequest(pago prismainforme.ListaPagosRequest) error {

	if pago.PageSize == 0 {
		return errors.New(ERROR_NUMBER_PAGE_SIZE)
	}
	err := validaFecha(pago.DateFrom)
	if err != nil {
		return errors.New(err.Error())
	}
	erro := validaFecha(pago.DateTo)
	if erro != nil {
		return errors.New(erro.Error())
	}
	return nil

}

func validaFecha(fecha string) error {
	res := strings.Split(fecha, "-")
	if len(res) != 3 {
		return errors.New(ERROR_FECHA)
	}
	if len(res[0]) != 4 {
		return errors.New(ERROR_FECHA)
	}
	if len(res[1]) != 2 {
		return errors.New(ERROR_FECHA)
	}
	if len(res[2]) != 2 {
		return errors.New(ERROR_FECHA)
	}
	return nil
}
