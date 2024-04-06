package cierrelote_test

// import (
// 	"context"
// 	"errors"
// 	"io/ioutil"
// 	"testing"

// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/cierrelote"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/cierrelote/cierrelotefake"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
// 	"github.com/stretchr/testify/assert"
// )

// func TestLeerArchivoCierrreLoteMinio(t *testing.T) {

// }
// func TestLeerCierreLoteFakePagosEstadoError(t *testing.T) {
// 	ctx := context.Background()
// 	tableDriverTestPagoEstadoError := cierrelotefake.EstructuraCierreLoteFakePagoEstadoError()
// 	//mockRemoteRepositoryPrisma := new(mockrepository.MockRemoteRepositoryPrisma)
// 	mockRepositoryPrisma := new(mockrepository.MockRepositoryPrisma)
// 	mockCommonds := new(mockservice.MockCommonsService)
// 	utilService := new(mockservice.MockUtilService)
// 	mockrepoStore := new(mockrepository.MockRepositoryAWSStore)
// 	service := cierrelote.NewService(mockRepositoryPrisma, mockCommonds, utilService, mockrepoStore)
// 	//service := prisma.NewService(mockRemoteRepositoryPrisma, mockRepositoryPrisma, mockCommonds, mockServiceAdministracion)
// 	t.Run(tableDriverTestPagoEstadoError.TituloPrueba, func(t *testing.T) {
// 		var estados []entities.Pagoestado
// 		//mockServiceAdministracion.On("GetPagosEstadosService", true, true).Return(estados, errors.New(cierrelote.ERROR_PAGO_ESTADO))
// 		errObtenerEstados := errors.New(tableDriverTestPagoEstadoError.WantError)
// 		want := errors.New(errObtenerEstados.Error())
// 		logError := entities.Log{
// 			Tipo:          entities.EnumLog("error"),
// 			Funcionalidad: "LeerCierreLoteTxt",
// 			Mensaje:       errObtenerEstados.Error() + "-" + want.Error(),
// 		}
// 		utilService.On("CreateLogService", logError).Return(want).Once()
// 		_, _, got := service.LeerCierreLoteTxt(ctx, "", estados)
// 		assert.Equal(t, got, want)
// 	})
// }

// func TestLeerCierreLoteFakeLeerDirectorioNoValido(t *testing.T) {
// 	ctx := context.Background()
// 	tableDriverTestPagoEstadoValido := cierrelotefake.EstructuraCierreLoteFakePagoEstadoValido()
// 	tableDriverTestArchivosError := cierrelotefake.EstructuraCierreLoteFakeBuscarArchivosError()
// 	//mockRemoteRepositoryPrisma := new(mockrepository.MockRemoteRepositoryPrisma)
// 	mockRepositoryPrisma := new(mockrepository.MockRepositoryPrisma)
// 	mockCommonds := new(mockservice.MockCommonsService)
// 	utilService := new(mockservice.MockUtilService)
// 	mockrepoStore := new(mockrepository.MockRepositoryAWSStore)
// 	service := cierrelote.NewService(mockRepositoryPrisma, mockCommonds, utilService, mockrepoStore)
// 	t.Run(tableDriverTestPagoEstadoValido.TituloPrueba, func(t *testing.T) {
// 		listaEstado := tableDriverTestPagoEstadoValido.WantEstadosPago
// 		want := errors.New(tableDriverTestArchivosError.WantError)
// 		//wantListaArchivos := []fs.FileInfo{}
// 		ruta := tableDriverTestArchivosError.RutaDirectorioOrigen
// 		wantListaArchivos, _ := ioutil.ReadDir(tableDriverTestArchivosError.RutaDirectorioOrigen)
// 		utilService.On("GetPagosEstadosService", true, true).Return(listaEstado, nil)
// 		mockCommonds.On("LeerDirectorio", ctx, ruta).Return(wantListaArchivos, errors.New(commons.ERROR_READ_ARCHIVO))
// 		_, _, got := service.LeerCierreLoteTxt(ctx, "", listaEstado)
// 		assert.Equal(t, got, want)
// 	})
// }

// func TestLeerCierreLoteFakeDirectorioVacio(t *testing.T) {
// 	ctx := context.Background()
// 	tableDriverTestPagoEstadoValido := cierrelotefake.EstructuraCierreLoteFakePagoEstadoValido()
// 	tableDriverTestDirectorioVacio := cierrelotefake.EstructuraCierreLoteFakeDirectorioVacio()
// 	mockRepositoryPrisma := new(mockrepository.MockRepositoryPrisma)
// 	mockCommonds := new(mockservice.MockCommonsService)
// 	utilService := new(mockservice.MockUtilService)
// 	mockrepoStore := new(mockrepository.MockRepositoryAWSStore)
// 	service := cierrelote.NewService(mockRepositoryPrisma, mockCommonds, utilService, mockrepoStore)
// 	t.Run(tableDriverTestDirectorioVacio.TituloPrueba, func(t *testing.T) {
// 		listaEstado := tableDriverTestPagoEstadoValido.WantEstadosPago
// 		want := errors.New(tableDriverTestDirectorioVacio.WantError)
// 		ruta := "C:/Users/Sergio/Downloads/archivosLotes/lotesinverificar"
// 		wantListaArchivos, _ := ioutil.ReadDir(tableDriverTestDirectorioVacio.RutaDirectorioOrigen)
// 		utilService.On("GetPagosEstadosService", true, true).Return(listaEstado, nil)
// 		mockCommonds.On("LeerDirectorio", ruta).Return(wantListaArchivos, nil)
// 		_, _, got := service.LeerCierreLoteTxt(ctx, "", listaEstado)
// 		assert.Equal(t, got, want)
// 	})
// }
