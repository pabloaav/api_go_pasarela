package cierrelote

import (
	"context"
	"io/fs"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
)

type Service interface {
	///////////////////////////////////servicios cierre lote/////////////////////////////////////////////////////////////
	/*
		ArchivoLoteExterno permite obtener los archivos de cierre de lote recibidos por ftp y mover
		el/los archivo/s a un directorio temporal de destino dentro del proyecto
	*/
	LeerArchivoLoteExterno(ctx context.Context, nombreDirectorio string) (archivos []fs.FileInfo, rutaArchivos string, totalArchivos int, err error)

	/*
		LeerCierreLote es un servicio que permite leer los cierre de loste eniviados por prisma y los carga a la base de datos
	*/
	LeerCierreLoteTxt([]fs.FileInfo, string, []entities.Pagoestadoexterno) (listaArchivo []prismaCierreLote.PrismaLogArchivoResponse, err error)

	/*
		peromite subir los cierre de lotes procesados y los que no fueron procesados al minio
	*/
	MoverArchivos(ctx context.Context, rutaArchivos string, listaArchivo []prismaCierreLote.PrismaLogArchivoResponse) (countArchivo int, erro error)

	/*
		BorrarArchivos: se encarga de eliminar los archivos temporales creado para procesar su contenido y de borrar los archivos recibidos en S3 de aws.
		recibe los siguientes parametros:
		- contexto.
		- "nombreDirectorio" nombre del directorio remoto del s3 donde se ubican los archivos
		- "rutaArchivos" ruta local donde se ubican los archivos temporales
		- "listaArchivo" y un array de objetos que contienen nombre del archivo y el estado del archivo (que indica si fue leido, movido e insertado en la base de datos)
	*/
	BorrarArchivos(ctx context.Context, nombreDirectorio string, rutaArchivos string, listaArchivo []prismaCierreLote.PrismaLogArchivoResponse) (erro error)

	/*
		permite obtener los movimientos del banco y los cierres de lotes de prisma para poder conciliar los movimientos del banco con los cierres de lotes de prisma
	*/
	ConciliacionBancoPrisma(fechaPagoProcesar string, reversion bool, responseListprismaTrPagos []prismaCierreLote.ResponseTrPagosCabecera) (listaMovimientosBanco []bancodtos.ResponseMovimientosBanco, erro error)

	/*
		este servicio es un proceso que se ejecuta una vez al dia para actualizar los los movimientos del banco
		que durante el proceso de conciliacion no fueron actualizados.
	*/
	ActualizarMovimientosBanco() (estadoResponse bool, erro error)

	/*este servicio obtiene la informacion de las tablas prismamxtotalesmovimientos prismamxdetallemoviminetos*/
	ObtenerMxMoviminetosServices() (movimientosMx []prismaCierreLote.ResponseMovimientoMx, entityMovimientoMxStr []entities.Prismamxtotalesmovimiento, erro error)
	ObtenerTablasRelacionadasServices() (tablasRelacionadas prismaCierreLote.ResponseTablasRelacionadas, erro error)
	ProcesarMovimientoMxServices(movimientosMx []prismaCierreLote.ResponseMovimientoMx, tablasRelacionadas prismaCierreLote.ResponseTablasRelacionadas) (resultMovimientosMx []prismaCierreLote.ResponseMovimientoMx)
	SaveMovimientoMxServices(movimientosMx []prismaCierreLote.ResponseMovimientoMx, movimientosMxEntity []entities.Prismamxtotalesmovimiento) (erro error)

	/*este servicio obtiene la informacion de las tablas prismapxdosregistros prismapxcuatroregistros*/
	ObtenerPxPagosServices() (pagosPx []prismaCierreLote.ResponsePagoPx, entityPagoPxStr []entities.Prismapxcuatroregistro, erro error)
	SavePagoPxServices(pagoPx []prismaCierreLote.ResponsePagoPx, entityPagoPxStr []entities.Prismapxcuatroregistro) (erro error)

	/*
		proceso para preparar la tabla prisma cierre de lote con la informacion necesaria par conciliar con los moviminetos del bamco
	*/
	ObtenerCierreloteServices(filtro filtrocl.FiltroCierreLote, codigoautorizacion []string) (listaCierreLote []prismaCierreLote.ResponsePrismaCL, erro error)
	ObtenerPrismaMovimientosServices(filtro filtrocl.FiltroPrismaMovimiento) (listaPrismaMovimientos []prismaCierreLote.ResponseMovimientoTotales, codigoautorizacion []string, erro error)
	ConciliarCierreLotePrismaMovimientoServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaPrismaMovimientos []prismaCierreLote.ResponseMovimientoTotales) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error)
	ActualizarCierreloteMoviminetosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error)
	ObtenerPrismaMovimientoConciliadosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, filtroCabecera filtrocl.FiltroPrismaMovimiento) (listaMovimientoDetalle []prismaCierreLote.ResponseCierrreLotePrismaMovimiento, erro error)
	ObtenerPrismaPagosServices(iltro filtrocl.FiltroPrismaTrPagos) (listaPrismaPago []prismaCierreLote.ResponseTrPagosCabecera, erro error)
	ConciliarCierreLotePrismaPagoServices(listaCierreLoteMovimientos []prismaCierreLote.ResponseCierrreLotePrismaMovimiento, listaPrismaPago []prismaCierreLote.ResponseTrPagosCabecera) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detallePagosIdArray []int64, cabeceraPagosIdArray []int64, erro error)
	ActualizarCierrelotePagosServices(listaCierreLote []prismaCierreLote.ResponsePrismaCL, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error)
	ObtenerRepoPagosPrisma(filtro filtrocl.FiltroTablasConciliadas) (responseListprismaTrPagos []prismaCierreLote.ResponseTrPagosCabecera, erro error)

	/*
		servicios para herramienta Cierre de Lote
	*/

	GetAllMovimientosPrismaServices(filtro filtrocl.FiltroMovimientosPrisma) (listaMovimiento []prismaCierreLote.ResponseMovimientoTotales, meta dtos.Meta, erro error)
	GetOneCierreLotePrismaService(oneFiltro filtrocl.OneCierreLoteFiltro) (cierreLote prismaCierreLote.ResponsePrismaCL, erro error)
	GetAllCierreLotePrismaService(filtro filtrocl.CierreLoteFiltro) (response prismaCierreLote.ResponsePrismaCLTools, erro error)
	EditCierreLotePrismaService(request cierrelotedtos.RequestPrismaCL) error
	DeleteCierreLotePrismaService(id uint64) error
	/*
		servicios para  Cierre de Lote APILINK
	*/
	GetAllCierreLoteApiLinkService(filtro filtrocl.ApilinkCierreloteFiltro) (response prismaCierreLote.ResponseApilinkCierresLotes, erro error)
}

var cierrelote *service

type service struct {
	repository                    Repository
	commonsService                commons.Commons
	utilService                   util.UtilService
	store                         Store
	administracionService         administracion.Service
	bancoService                  banco.BancoService
	factory                       RecorrerArchivosFactory
	factoryConciliacionMovimiento ConciliarClMovimientosFactory
}

func NewService(r Repository, c commons.Commons, util util.UtilService, s Store, adm administracion.Service, bancoService banco.BancoService) Service {
	cierrelote := &service{
		commonsService:                c,
		utilService:                   util,
		repository:                    r,
		store:                         s,
		administracionService:         adm,
		bancoService:                  bancoService,
		factory:                       &recorrerArchivosFactory{},
		factoryConciliacionMovimiento: &conciliarClMovimientosFactory{},
	}
	return cierrelote

}

func NewServiceConFactory(r Repository, c commons.Commons, util util.UtilService, s Store, bancoService banco.BancoService, factory recorrerArchivosFactory, factoryClMov conciliarClMovimientosFactory) Service {
	cierrelote := service{
		commonsService:                c,
		utilService:                   util,
		repository:                    r,
		store:                         s,
		bancoService:                  bancoService,
		factory:                       &factory,
		factoryConciliacionMovimiento: &factoryClMov,
	}
	return &cierrelote
}

func Resolve() Service {
	return cierrelote
}
