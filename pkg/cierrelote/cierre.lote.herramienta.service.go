package cierrelote

import (
	"errors"

	"math"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
)

func (s *service) GetAllMovimientosPrismaServices(filtro filtrocl.FiltroMovimientosPrisma) (listaMovimientos []prismaCierreLote.ResponseMovimientoTotales, meta dtos.Meta, erro error) {

	listaMovimientoPrismaResult, totalFilas, err := s.repository.GetAllMovimientosPrismaRepository(filtro)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	if len(listaMovimientoPrismaResult) > 0 {
		for _, valuePrismaMovimiento := range listaMovimientoPrismaResult {
			var movimientoPrismaTemporal cierrelotedtos.ResponseMovimientoTotales
			var movimientoPrismaDetalleTemporal cierrelotedtos.ResponseMoviminetoDetalles
			movimientoPrismaTemporal.EntityToDtos(valuePrismaMovimiento)

			for _, valuePrismaDetalle := range valuePrismaMovimiento.DetalleMovimientos {
				movimientoPrismaDetalleTemporal.EntityToDtos(valuePrismaDetalle)
				movimientoPrismaTemporal.DetalleMovimientos = append(movimientoPrismaTemporal.DetalleMovimientos, movimientoPrismaDetalleTemporal)
			}
			listaMovimientos = append(listaMovimientos, movimientoPrismaTemporal)
		}
	}
	if filtro.Number > 0 && filtro.Size > 0 {
		meta = setPaginacion(filtro.Number, filtro.Size, totalFilas)
	}
	return
}

func setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))

	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	return

}

func (s *service) GetOneCierreLotePrismaService(filtro filtrocl.OneCierreLoteFiltro) (cierreLote prismaCierreLote.ResponsePrismaCL, erro error) {

	// hacer la cosulta al repositorio, que debe retornar una entidad Prismacierrelote
	oneCierreLote, erro := s.repository.GetOneCierreLotePrismaRepository(filtro)

	if erro != nil {
		return
	}

	// transforma el entity cierre de lote en tipo ResponsePrismaCL
	cierreLote.EntityToDtos(oneCierreLote)

	return
}

func (s *service) GetAllCierreLotePrismaService(filtro filtrocl.CierreLoteFiltro) (response prismaCierreLote.ResponsePrismaCLTools, erro error) {

	// hacer la cosulta al repositorio, que debe retornar un array de entidad Prismacierrelote, el total de filas y un error
	listaCierreLote, total, erro := s.repository.GetAllCierreLotePrismaRepository(filtro)

	if erro != nil {
		return
	}
	// Con el total de filas de la query se setea la paginacion
	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = setPaginacion(filtro.Number, filtro.Size, total)
	}

	// Es necesario inicializar la struct de tipo ResponsePrismaCL para poder acceder a sus metodos
	var clTemporal prismaCierreLote.PrismaCLTools
	var ChannelarancelTemporal prismaCierreLote.ResponseChannelArancel
	var movimientoTemporal prismaCierreLote.ResponseMoviminetoDetPrisma
	var pagoTemporal prismaCierreLote.ResponsePagoPrisma

	// Un for range para recorrer cada entidad y transformar en una struct de tipo ResponsePrismaCL
	for _, cl := range listaCierreLote {
		clTemporal.EntityCLToDto(&cl) // aqui se puede acceder al metodo (receiver function)
		ChannelarancelTemporal.EntityToDtos(*cl.Channelarancel)
		clTemporal.Channelarancel = ChannelarancelTemporal
		if filtro.CargarMovimientoPrisma {
			movimientoTemporal.EntityMovToDto(cl.Prismamovimientodetalle)
			clTemporal.MovimientoPrisma = movimientoTemporal
		}
		if filtro.CargarPagoPrisma {
			pagoTemporal.EntityPagoToDto(cl.Prismatrdospagos)
			clTemporal.PagoPrisma = pagoTemporal
		}

		response.CierresLotes = append(response.CierresLotes, clTemporal)
	}
	return
}

func (s *service) EditCierreLotePrismaService(request cierrelotedtos.RequestPrismaCL) (erro error) {
	erro = request.IsValid()
	if erro != nil {
		return
	}
	filtro := filtrocl.OneCierreLoteFiltro{Id: uint64(request.Id)}
	result, erro := s.repository.GetOneCierreLotePrismaRepository(filtro)
	if erro != nil {
		return
	}
	filtroEdit := request.ValidarCamposEdit(result)

	entity := request.RequestPrismaCLToEntity()
	logs.Info(filtroEdit)
	erro = s.repository.EditOneCierreLotePrismaRepository(entity, filtroEdit)
	if erro != nil {
		return
	}
	return nil
}

func (s *service) DeleteCierreLotePrismaService(id uint64) (erro error) {

	return s.repository.DeleteOneCierreLotePrismaRepository(id)

}

func (s *service) GetAllCierreLoteApiLinkService(filtro filtrocl.ApilinkCierreloteFiltro) (response prismaCierreLote.ResponseApilinkCierresLotes, erro error) {

	// Las fechas se pasan al formato Time
	filtrosfechas := filtro.ToFiltroRequest()

	// se recibe respuesta del repositorio con datos de base de datos. o un error
	apilinkEntities, totalFilas, erro := s.repository.GetAllCierreLoteApiLinkRepository(filtrosfechas)

	if erro != nil {
		return
	}
	// paginacion
	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = setPaginacion(filtro.Number, filtro.Size, totalFilas)
	}

	// tranformar la consulta en un tipo response
	var clTemporal prismaCierreLote.ResponseApilinkCL

	for _, apilinkcl := range apilinkEntities {

		clTemporal.EntityToDto(apilinkcl)

		response.CierresLotes = append(response.CierresLotes, clTemporal)
	}

	return

}
