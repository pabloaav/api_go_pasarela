package cierrelote

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtrocl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/cierrelote"
)

type Repository interface {
	// SaveCierreLote guarda los cierres de lostes em la base de datos
	SaveCierreLote(detalleLote *entities.Prismacierrelote) (bool, error)
	SaveCierreLoteBatch(detalleLote []entities.Prismacierrelote) (bool, error)
	GetCierreLoteRepository(filtro filtrocl.FiltroCierreLote) (entityCL []entities.Prismacierrelote, erro error)
	GetPagosIntentosByMedioPagoIdRepository(arraysMediosPagoIds []int64) (listaPagoIntentos []entities.Pagointento, erro error)
	GetCierreLoteGroupByRepository(nroCuota int64) ([]cierrelotedtos.PrismaClResultGroup, error)
	ActualizarCierreLoteMatch(reversion bool, clMatch []entities.Prismacierrelote) error
	GetCierreLoteByGroup(clMatch cierrelotedtos.PrismaClResultGroup) ([]entities.Prismacierrelote, error)
	GetCierreLotePrismaByExternalIdAndMacht() (listCierreLote []entities.Prismacierrelote, erro error)
	SaveTransactionPagoPx(pagoPx []entities.Prismapxcuatroregistro) (erro error)
	SaveTransactionMovimientoMx(movimientosMx []entities.Prismamxtotalesmovimiento) (erro error)
	SaveTransactionPagoRP(pagosRP []entities.Rapipagocierrelote) (erro error)
	GetMovimientosMxRepository() (movimientosMx []entities.Prismamxtotalesmovimiento, erro error)
	GetCodigosRechazoRepository() (codigosRechazo []entities.Prismacodigorechazo, erro error)
	GetVisaContracargoRepository() (visaContracargo []entities.Prismavisacontracargo, erro error)
	GetMotivosAjustesRepository() (motivosAjustes []entities.Prismamotivosajuste, erro error)
	GetOperacionesRepository() (operaciones []entities.Prismaoperacion, erro error)
	GetMasterContracargoRepository() (masterContracargo []entities.Prismamastercontracargo, erro error)
	SaveMovimientoMxRepository(movimientosMx []entities.Prismamovimientototale, movimientosMxEntity []entities.Prismamxtotalesmovimiento) (erro error)
	GetPagosPxRepository() (pagosPx []entities.Prismapxcuatroregistro, erro error)
	SavePagosPxRepository(pagosPx []entities.Prismatrcuatropago, entityPagoPxStr []entities.Prismapxcuatroregistro) (erro error)
	GetArancelByRubroIdChannelIdRepository(rubroId, channelId uint) (entityArancel entities.Channelarancele, erro error)
	GetPrismaMovimientosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityPrismaMovimientos []entities.Prismamovimientototale, erro error)
	GetContraCargoPrismaMovimientosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityPrismaMovimientos []entities.Prismamovimientototale, erro error)
	UpdateCierreloteAndMoviminetosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error)

	GetMovimientosConciliadosRepository(filtro filtrocl.FiltroPrismaMovimiento) (entityMovimientosConciliados []entities.Prismamovimientototale, erro error)
	GetMovimientosDetalleConciliadosRepository(filtro filtrocl.FiltroPrismaMovimientoDetalle) (entityMovimientosDetalleConciliados []entities.Prismamovimientodetalle, erro error)
	GetPrismaPagosRepository(filtro filtrocl.FiltroPrismaTrPagos) (entityPrismaPago []entities.Prismatrcuatropago, erro error)
	UpdateCierreloteAndPagosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error)
	GetCierreLoteMatch(filtro filtrocl.FiltroTablasConciliadas) (entityPrismaTr4Pagos []entities.Prismatrcuatropago, erro error)

	/*
		repositorios para herramienta cierre lote
	*/

	GetAllMovimientosPrismaRepository(filtro filtrocl.FiltroMovimientosPrisma) (entityMovimientoPrisma []entities.Prismamovimientototale, totalFilas int64, erro error)
	GetOneCierreLotePrismaRepository(oneFiltro filtrocl.OneCierreLoteFiltro) (oneCierreLote entities.Prismacierrelote, erro error)
	GetAllCierreLotePrismaRepository(filtro filtrocl.CierreLoteFiltro) (listaCierreLote []entities.Prismacierrelote, totalFilas int64, erro error)
	EditOneCierreLotePrismaRepository(entityPrismaCL entities.Prismacierrelote, filtroEdit cierrelotedtos.FiltroEditarCLHerramienta) (erro error)
	DeleteOneCierreLotePrismaRepository(id uint64) (erro error)

	/*
		repositorios para  cierre lote ApiLink
	*/
	GetAllCierreLoteApiLinkRepository(filtrosfechas cierrelotedtos.ApilinkRequest) (listaApilinkCierreLote []entities.Apilinkcierrelote, totalFilas int64, erro error)
}

type repository struct {
	SQLClient      *database.MySQLClient
	UtilRepository util.UtilRepository
}

func NewRepository(sqlClient *database.MySQLClient, ur util.UtilRepository) Repository {
	return &repository{
		SQLClient:      sqlClient,
		UtilRepository: ur,
	}
}
