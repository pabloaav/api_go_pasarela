package prisma

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	commonds "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
)

func (s *service) PostAnulacionDevolucionTotalPago(ExternalId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, erro error) {
	if commonds.StringIsEmpity(ExternalId) {
		return nil, errors.New(ERROR_EXTERNAL_ID)
	}
	resultado, err := s.remoteRepository.PostSolicitudAnulacionDevolucionPagoTotla(ExternalId)
	if err != nil {
		//e := err.(*ErrorEstructura)
		//msjError := BuscarMensajeError(e.ErrorType)
		//return nil, errors.New(msjError)
		erro = err
		logs.Info(err.Error())
		return nil, errors.New(ERROR_PETICION_ANULACION_DEVOLUCION)
	}
	response = resultado
	return response, nil

}

func (s *service) PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, erro error) {

	err := params.ValidarParametros("SAPP")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	resultado, err := s.remoteRepository.PostSolicitudAnulacionDevolucionPagoParcial(params)
	if err != nil {
		erro = err
		return nil, erro
	}
	response = resultado
	return response, nil
}

func (s *service) DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, erro error) {
	err := params.ValidarParametros("ER")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	resultado, err := s.remoteRepository.DelAnulacionDevolucionPagoTotalParcial(params)
	if err != nil {
		erro = err
		return nil, erro
	}
	response = resultado
	return response, nil
}
