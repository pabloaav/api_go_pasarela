package banco

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/banco"
)

type BancoService interface {
	GetGenerarTokenBancoService() (token bancodtos.ResponseTokenBanco, erro error)
	/*
		Autor: Jose Alarcon
		Fecha: 27/04/2022
		Descripción: Consultar movimientos del banco RECIBE:
		- subcuenta
		- tipo (debin, prisma)
		- lista de movimientos
	*/
	BuildCierreLoteApiLinkBancoService(filtro filtros.MovimientosBancoFiltro) (listaMovimientoBanco []bancodtos.ResponseMovimientosBanco, erro error) //listaCierreApiLink []*entities.Apilinkcierrelote
	ActualizarRegistrosMatchBancoService(listaIds []uint, estadoMatch bool) (estadoResponse bool, erro error)

	/*
		Autor: Jose Alarcon
		Fecha: 26/05/2022
		Descripción: coinciliacion de pagos con servicio de banco
	*/
	ConciliacionPasarelaBanco(request bancodtos.RequestConciliacion) (response bancodtos.ResponseConciliacion, movimientosIds []uint, erro error)

	// servicios para herramienta wee
	GetConsultarMovimientosService(filtro filtros.RequestMovimientos) (response bancodtos.ResponseMovimientos, erro error)
	// Procesar movimientos (background_service) de banco
	GetProcesarMovimientosGeneralService() (response bancodtos.ResponseMovimientosProceso, erro error)
}

type bancoService struct {
	remoteRepository RemoteRepository
	utilService      util.UtilService
	// apilinkService   apilink.AplinkService
	administracion administracion.Service
	factory        ConciliacionFactory
}

// func NewService(rm RemoteRepository, util util.UtilService, adm administracion.Service) BancoService {
// 	return &bancoService{
// 		remoteRepository: rm,
// 		utilService:      util,
// 		// apilinkService:   link,
// 		administracion: adm,
// 		factory:        &procesarConciliacionFactory{},
// 	}
// }

func NewService(rm RemoteRepository, util util.UtilService, adm administracion.Service) BancoService {
	banco := &bancoService{
		remoteRepository: rm,
		utilService:      util,
		// apilinkService:   link,
		administracion: adm,
		factory:        &procesarConciliacionFactory{},
	}
	return banco

}

// func NewServiceConFactory(rm RemoteRepository, util util.UtilService, adm administracion.Service, factory procesarConciliacionFactory) BancoService {
// 	banco := bancoService{
// 		remoteRepository: rm,
// 		utilService:      util,
// 		// apilinkService:   link,
// 		administracion: adm,
// 		factory:        &factory,
// 	}
// 	return &banco
// }

func (s *bancoService) GetGenerarTokenBancoService() (token bancodtos.ResponseTokenBanco, erro error) {
	/* 1 solicitar token de autentificacion */
	/* 	SE DEBE CREAR UN USUARIO PARA BANCO */
	requestTokenBanco := bancodtos.RequestTokenBanco{
		Username:  config.BANCO_USER,
		Password:  config.BANCO_PASS,
		SistemaId: config.BANCO_SISTEMA,
	}

	token, erro = s.remoteRepository.GetTokenBancoRemotRepository(requestTokenBanco)
	if erro != nil {
		return
	}
	return token, nil
}

func (s *bancoService) BuildCierreLoteApiLinkBancoService(filtro filtros.MovimientosBancoFiltro) (listaMovimientoBanco []bancodtos.ResponseMovimientosBanco, erro error) {

	// Paso 1: obtener token de autentificacion
	token, erro := s.GetGenerarTokenBancoService()
	if erro != nil {
		return
	}

	// Paso 2: consultar movimientos
	response, erro := s.remoteRepository.GetMovimientosBancoRemotRepository(filtro, token.Token)

	if erro != nil || len(response) < 1 {
		return
	}
	listaMovimientoBanco = response
	return
}

func (s *bancoService) ActualizarRegistrosMatchBancoService(listaIds []uint, estadoMatch bool) (estadoResponse bool, erro error) {

	// La lista de ids no deberia venir vacia
	if len(listaIds) < 1 {
		erro = errors.New(ERROR_ACTUALIZAR_MATHC)
		return
	}

	/* 1 obtener token de autentificacion */
	token, erro := s.GetGenerarTokenBancoService()
	if erro != nil {
		return
	}

	/*2 actualizar registros de match */
	listaMovimientosMatch := bancodtos.RequestUpdateMovimiento{

		ListaMovimientos: listaIds,
		EstadoMatch:      estadoMatch,
	}
	estadoResponse, err := s.remoteRepository.ActualizarRegistrosMatchRepository(listaMovimientosMatch, token.Token)
	if err != nil {
		erro = errors.New(ERROR_ACTUALIZAR_MATHC_BANCO)
		return
	}
	return
}

func (s *bancoService) ConciliacionPasarelaBanco(request bancodtos.RequestConciliacion) (response bancodtos.ResponseConciliacion, movimientosIds []uint, erro error) {

	/* Factory: recibe las listas para conciliar de apilink ,tranferencias y rapipago */
	metodoConciliar, err := s.factory.GetProcesarConciliacion(request.TipoConciliacion)
	if err != nil {
		erro = err
		logs.Error(err)
		return
	}

	// devuelve un tipo MovimientosBancoFiltro. Incluye una lista de las referencia_banco
	filtro := metodoConciliar.FiltroRequestConsultaBanco(request)

	// consultar al servicio de proyecto banco los movimientos que coinciden con las transferencias teniendo en cuenta el campo referencia_banco
	listaMovimientosBanco, err := s.BuildCierreLoteApiLinkBancoService(filtro)

	erro = err
	if erro != nil {
		erro = errors.New(ERROR_OBTENER_MOVIMIENTOS_BANCO + erro.Error())
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BuildCierreLoteApiLinkBancoService",
			Mensaje:       erro.Error(),
		}
		er := s.utilService.CreateLogService(log)
		if er != nil {
			logs.Error("error: al crear logs: " + er.Error())
		}
		return
	}
	response, movimientosIds = metodoConciliar.ConciliacionBanco(request, listaMovimientosBanco)

	return response, movimientosIds, nil
}

func (s *bancoService) GetConsultarMovimientosService(filtro filtros.RequestMovimientos) (response bancodtos.ResponseMovimientos, erro error) {
	token, erro := s.GetGenerarTokenBancoService()
	if erro != nil {
		return
	}
	response, err := s.remoteRepository.GetConsultarMovimientosRemoteRepository(filtro, token.Token)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_MOVIMIENTOS_BANCO + " " + err.Error())
		return
	}
	return
}

func (s *bancoService) GetProcesarMovimientosGeneralService() (response bancodtos.ResponseMovimientosProceso, erro error) {
	token, erro := s.GetGenerarTokenBancoService()
	if erro != nil {
		return
	}
	response, err := s.remoteRepository.GetProcesarMovimientosRemoteRepository(token.Token)
	if err != nil {
		erro = errors.New(ERROR_OBTENER_MOVIMIENTOS_BANCO + " " + err.Error())
		return
	}
	return
}
