package administracion_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

type callsMocked struct {
	Nombre   string
	Request  interface{}
	Response interface{}
	Erro     error
}

// type TableRICuentasCliente struct {
// 	Nombre        string
// 	Request       administraciondtos.RICuentasClienteRequest
// 	Erro          error
// 	MetodosMocked []callsMocked
// }

// type TableRIDatosFondos struct {
// 	Nombre        string
// 	Request       administraciondtos.RiDatosFondosRequest
// 	Erro          error
// 	MetodosMocked []callsMocked
// }

// func _inicializarCuentasClientes() (table []TableRICuentasCliente) {

// 	riCuentaClienteInvalido := []administraciondtos.RiCuentaCliente{
// 		{CodigoPartida: "2500", Saldo: "1", Cantidad: "1", CBU: "1"},
// 		{CodigoPartida: "3400777777777", Saldo: "2", Cantidad: "2", CBU: "2"},
// 		{CodigoPartida: "3", Saldo: "3", Cantidad: "3", CBU: "3"},
// 	}

// 	riCuentaCliente := []administraciondtos.RiCuentaCliente{
// 		{CodigoPartida: "2500", Saldo: "1", Cantidad: "1", CBU: "1"},
// 		{CodigoPartida: "340777", Saldo: "2", Cantidad: "2", CBU: "2"},
// 		{CodigoPartida: "3", Saldo: "3", Cantidad: "3", CBU: "3"},
// 	}

// 	requestRutaInValida := administraciondtos.RICuentasClienteRequest{
// 		FechaInicio: time.Now(),
// 		FechaFin:    time.Now().Add(time.Hour * 24),
// 		Ruta:        ".png",
// 	}

// 	requestFechaInvalida := administraciondtos.RICuentasClienteRequest{
// 		FechaInicio: time.Now(),
// 	}

// 	guardarRequest := administraciondtos.RIGuardarArchivosRequest{
// 		Ruta: requestRutaInValida.Ruta, RI: riCuentaCliente,
// 	}

// 	table = []TableRICuentasCliente{
// 		{"Debe Retornar un error si algún parametro es inválido",
// 			requestFechaInvalida, fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA), []callsMocked{}},
// 		{"Debe Retornar un error si algun dato de la ri es invalido",
// 			requestRutaInValida, fmt.Errorf(administraciondtos.ERROR_RI_CODIGO_INVALIDO), []callsMocked{
// 				{"BuildRICuentasCliente", requestRutaInValida, riCuentaClienteInvalido, nil},
// 			},
// 		},
// 		{"Debe Retornar un error si no se puede guardar el archivo",
// 			requestRutaInValida, fmt.Errorf(commons.ERROR_FILE_NAME), []callsMocked{
// 				{"BuildRICuentasCliente", requestRutaInValida, riCuentaCliente, nil},
// 				{"CreateFile", guardarRequest.Ruta, &os.File{}, fmt.Errorf(commons.ERROR_FILE_NAME)},
// 			},
// 		},
// 	}

// 	return
// }

// func TestRICuentasClientes(t *testing.T) {
// 	table := _inicializarCuentasClientes()
// 	mockRepository := new(mockrepository.MockRepositoryAdministracion)
// 	mockApiLinkService := new(mockservice.MockApiLinkService)
// 	mockCommonsService := new(mockservice.MockCommonsService)

// 	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService)

// 	for _, v := range table {

// 		t.Run(v.Nombre, func(t *testing.T) {

// 			for _, m := range v.MetodosMocked {
// 				mockRepository.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
// 				if m.Nombre == "CreateFile" {
// 					mockCommonsService.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
// 				}
// 			}

// 			ri, err := service.RICuentasCliente(v.Request)

// 			if err != nil {
// 				assert.Equal(t, v.Erro.Error(), err.Error())
// 			}
// 			if ri != nil {
// 				assert.Equal(t, v.MetodosMocked[0].Response, ri)
// 			}

// 		})
// 	}
// }

// func _inicializarDatosFondos() (table []TableRIDatosFondos) {

// 	riDatosFondosInvalido := []administraciondtos.RiDatosFondos{
// 		{Numero: "150999996666666", Denominacion: "Prueba con Tílde", Agente: administraciondtos.AgenteAdministracion, DenominacionAgente: "prueba minuscula", CuitAgente: "27953043638"},
// 		{Numero: "151", Denominacion: "Prueba con Tílde", Agente: administraciondtos.AgenteAdministracion, DenominacionAgente: "prueba minuscula", CuitAgente: "27953043638"},
// 		{Numero: "152", Denominacion: "Prueba con Tílde", Agente: administraciondtos.AgenteAdministracion, DenominacionAgente: "prueba minuscula", CuitAgente: "27953043638"},
// 	}

// 	riDatosFondosNormalizar := []administraciondtos.RiDatosFondos{
// 		{Numero: "150", Denominacion: "Prueba con Tílde", Agente: administraciondtos.AgenteAdministracion, DenominacionAgente: "Prueba con Tílde", CuitAgente: "27953043638"},
// 	}

// 	requestRutaInValida := administraciondtos.RiDatosFondosRequest{
// 		FechaInicio: time.Now(),
// 		FechaFin:    time.Now().Add(time.Hour * 24),
// 		Ruta:        ".png",
// 	}

// 	requestFechaInvalida := administraciondtos.RiDatosFondosRequest{
// 		FechaInicio: time.Now(),
// 	}

// 	guardarRequest := administraciondtos.RIGuardarArchivosRequest{
// 		Ruta: requestRutaInValida.Ruta, RI: riDatosFondosNormalizar,
// 	}

// 	table = []TableRIDatosFondos{
// 		{"Debe Retornar un error si algún parametro es inválido",
// 			requestFechaInvalida, fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA), []callsMocked{}},
// 		{"Debe Retornar un error si algun dato de la ri es invalido",
// 			requestRutaInValida, fmt.Errorf(administraciondtos.ERROR_RI_NUMERO_FONDO), []callsMocked{
// 				{"BuildRIDatosFondo", requestRutaInValida, riDatosFondosInvalido, nil},
// 			},
// 		},
// 		{"Debe Retornar un error si no se puede normalizar los strings",
// 			requestRutaInValida, fmt.Errorf(commons.ERROR_NORMALIZAR), []callsMocked{
// 				{"BuildRIDatosFondo", requestRutaInValida, riDatosFondosNormalizar, nil},
// 				{"NormalizeStrings", "Prueba con Tílde", "Prueba con Tílde", fmt.Errorf(commons.ERROR_NORMALIZAR)},
// 			},
// 		},
// 		{"Debe Retornar un error si no se puede guardar el archivo",
// 			requestRutaInValida, fmt.Errorf(commons.ERROR_FILE_NAME), []callsMocked{
// 				{"BuildRIDatosFondo", requestRutaInValida, riDatosFondosNormalizar, nil},
// 				{"NormalizeStrings", "Prueba con Tílde", "PRUEBA CON TILDE", nil},

// 				{"CreateFile", guardarRequest.Ruta, &os.File{}, fmt.Errorf(commons.ERROR_FILE_NAME)},
// 			},
// 		},
// 	}

// 	return
// }

// func TestRIDatosFondos(t *testing.T) {
// 	table := _inicializarDatosFondos()
// 	mockRepository := new(mockrepository.MockRepositoryAdministracion)
// 	mockApiLinkService := new(mockservice.MockApiLinkService)
// 	mockCommonsService := new(mockservice.MockCommonsService)

// 	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService)

// 	for _, v := range table {

// 		t.Run(v.Nombre, func(t *testing.T) {

// 			for _, m := range v.MetodosMocked {
// 				mockRepository.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()

// 				if m.Nombre == "CreateFile" || m.Nombre == "NormalizeStrings" {
// 					mockCommonsService.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Maybe()
// 				}
// 			}

// 			ri, err := service.RIDatosFondos(v.Request)

// 			if err != nil {
// 				assert.Equal(t, v.Erro.Error(), err.Error())
// 			}
// 			if ri != nil {
// 				assert.Equal(t, v.MetodosMocked[0].Response, ri)
// 			}

// 		})
// 	}
// }

type TableRIGuardarArchivos struct {
	Nombre        string
	Request       ribcradtos.RIGuardarArchivosRequest
	Erro          error
	MetodosMocked []callsMocked
}

func _inicializarGuardarArchivos() (table []TableRIGuardarArchivos) {

	rutasInvalidas := []string{"", " ", "      "}
	rutaValida := "C/Usuarios/prueba.txt"
	rutaExtensionInvalida := "C/Usuarios/prueba.png"
	riInfEstadistica := []ribcradtos.RiInfestadistica{
		{CodigoPartida: "100001", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 10, MontoTotal: "25"},
		{CodigoPartida: "100002", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 11, MontoTotal: "250"},
		{CodigoPartida: "100003", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 12, MontoTotal: "2500"},
	}

	table = []TableRIGuardarArchivos{
		{
			"Debe Retornar un error si los datos son nulos",
			ribcradtos.RIGuardarArchivosRequest{
				Ruta: rutaValida,
				RI:   nil,
			}, fmt.Errorf(administraciondtos.ERROR_RI_DATOS), []callsMocked{},
		},
		{
			"Debe Retornar un error si no puede crear el archivo",
			ribcradtos.RIGuardarArchivosRequest{
				Ruta: rutaExtensionInvalida,
				RI:   riInfEstadistica,
			}, fmt.Errorf(commons.ERROR_FILE_NAME), []callsMocked{
				{"CreateFile", rutaExtensionInvalida, &os.File{}, fmt.Errorf(commons.ERROR_FILE_NAME)},
			},
		},
	}

	for i := range rutasInvalidas {
		item := TableRIGuardarArchivos{
			"Debe Retornar un error si la ruta es inválida",
			ribcradtos.RIGuardarArchivosRequest{
				Ruta: rutasInvalidas[i],
				RI:   nil,
			}, fmt.Errorf(administraciondtos.ERROR_RUTA_INVALIDA), []callsMocked{},
		}
		table = append(table, item)
	}

	return
}

func TestGuardarArchivos(t *testing.T) {

	table := _inicializarGuardarArchivos()
	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)
	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for _, m := range v.MetodosMocked {
				mockCommonsService.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
			}

			err := service.RIGuardarArchivos(v.Request)

			if err != nil {
				assert.Equal(t, v.Erro.Error(), err.Error())
			}

		})
	}
}

type TableRIInfestadistica struct {
	Nombre        string
	Request       ribcradtos.RiInfestadisticaRequest
	Erro          error
	MetodosMocked []callsMocked
}

func _inicializarInfestadistica() (table []TableRIInfestadistica) {

	riInfestadisticaInvalido := []ribcradtos.RiInfestadistica{
		{CodigoPartida: "100001", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 10, MontoTotal: "25"},
		{CodigoPartida: "100002", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 11, MontoTotal: "250"},
		{CodigoPartida: "1000000000003", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 12, MontoTotal: "2500"},
	}

	riInfestadistica := []ribcradtos.RiInfestadistica{
		{CodigoPartida: "100001", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 10, MontoTotal: "25"},
		{CodigoPartida: "100002", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 11, MontoTotal: "250"},
		{CodigoPartida: "100003", MedioPago: "1", EsquemaPago: "1", CantOperaciones: 12, MontoTotal: "2500"},
	}
	requestRutaInValida := ribcradtos.RiInfestadisticaRequest{
		FechaInicio: time.Now(),
		FechaFin:    time.Now().Add(time.Hour * 24),
		Ruta:        ".png",
	}

	requestFechaInvalida := ribcradtos.RiInfestadisticaRequest{
		FechaInicio: time.Now(),
	}

	guardarRequest := ribcradtos.RIGuardarArchivosRequest{
		Ruta: requestRutaInValida.Ruta, RI: riInfestadistica,
	}

	table = []TableRIInfestadistica{
		{"Debe Retornar un error si algún parametro es inválido",
			requestFechaInvalida, fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA), []callsMocked{}},
		{"Debe Retornar un error si algun dato de la ri es invalido",
			requestRutaInValida, fmt.Errorf(administraciondtos.ERROR_RI_CODIGO_INVALIDO), []callsMocked{
				{"BuilRIInfestaditica", requestRutaInValida, riInfestadisticaInvalido, nil},
			},
		},
		{"Debe Retornar un error si no se puede guardar el archivo",
			requestRutaInValida, fmt.Errorf(commons.ERROR_FILE_NAME), []callsMocked{
				{"BuilRIInfestaditica", requestRutaInValida, riInfestadistica, nil},
				{"CreateFile", guardarRequest.Ruta, &os.File{}, fmt.Errorf(commons.ERROR_FILE_NAME)},
			},
		},
	}

	return
}

func TestRIInfestadistica(t *testing.T) {
	table := _inicializarInfestadistica()
	mockRepository := new(mockrepository.MockRepositoryAdministracion)
	mockApiLinkService := new(mockservice.MockApiLinkService)
	mockCommonsService := new(mockservice.MockCommonsService)
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	mockStore := new(mockservice.MockStoreService)

	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService, mockUtilService, mockRepositoryWebHook, mockStore)

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			for _, m := range v.MetodosMocked {
				mockRepository.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
				if m.Nombre == "CreateFile" {
					mockCommonsService.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
				}
			}

			ri, err := service.RIInfestadistica(v.Request)

			if err != nil {
				assert.Equal(t, v.Erro.Error(), err.Error())
			}
			if ri != nil {
				assert.Equal(t, v.MetodosMocked[0].Response, ri)
			}

		})
	}
}

type TableBuildInformacionSupervision struct {
	Nombre        string
	Request       ribcradtos.GetInformacionSupervisionRequest
	Erro          error
	MetodosMocked []callsMocked
}

// func _inicializarInformacionSupervision() (table []TableBuildInformacionSupervision) {

// 	requestTipoInvalido := ribcradtos.GetInformacionSupervisionRequest{
// 		FechaInicio: time.Now(),
// 		FechaFin:    time.Now().Add(time.Hour * 25),
// 	}

// 	requestValido := ribcradtos.GetInformacionSupervisionRequest{
// 		FechaInicio: time.Now(),
// 		FechaFin:    time.Now().Add(time.Hour * 25),
// 	}

// 	nombreArchivoCuentas := commonsdtos.FileName{
// 		RutaBase:  config.RUTA_RI_BCRA,
// 		Nombre:    "CUENTASCLIENTES",
// 		Extension: "txt",
// 		UsaFecha:  false,
// 	}

// 	nombreArchivoDatos := commonsdtos.FileName{
// 		RutaBase:  config.RUTA_RI_BCRA,
// 		Nombre:    "DATOSFONDOS",
// 		Extension: "txt",
// 		UsaFecha:  false,
// 	}

// 	nombreArchivoDetalle := commonsdtos.FileName{
// 		RutaBase:  config.RUTA_RI_BCRA,
// 		Nombre:    "detalle",
// 		Extension: "xml",
// 		UsaFecha:  false,
// 	}

// 	requestCuentasClientes := ribcradtos.RICuentasClienteRequest{
// 		FechaInicio: requestValido.FechaInicio,
// 		FechaFin:    requestValido.FechaFin,
// 		Ruta:        "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt",
// 	}

// 	requestDatosFondos := ribcradtos.RiDatosFondosRequest{
// 		FechaInicio: requestValido.FechaInicio,
// 		FechaFin:    requestValido.FechaFin,
// 		Ruta:        "C:/Users/Alexandre/Documents/RIBCRA/DATOSFONDOS.txt",
// 	}

// 	table = []TableBuildInformacionSupervision{
// 		{"Debe Retornar un error si algún parametro es inválido",
// 			requestTipoInvalido, fmt.Errorf(administraciondtos.ERROR_RI_TIPO_PRESENTACION), []callsMocked{}},
// 		{"Debe Retornar un error si no se puede crear ri cuentas cliente",
// 			requestValido, fmt.Errorf("no se pudo cargar ri cuentas cliente"), []callsMocked{
// 				{"CreateFileName", nombreArchivoCuentas, "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt", nil},
// 				{"BuildRICuentasCliente", requestCuentasClientes, []ribcradtos.RiCuentaCliente{}, fmt.Errorf("no se pudo cargar ri cuentas cliente")},
// 			},
// 		},
// 		{"Debe Retornar un error si no se puede crear ri datos fondo",
// 			requestValido, fmt.Errorf("no se pudo cargar ri datos fondos"), []callsMocked{
// 				{"CreateFileName", nombreArchivoCuentas, "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt", nil},
// 				{"BuildRICuentasCliente", requestCuentasClientes, []ribcradtos.RiCuentaCliente{}, nil},
// 				{"CreateFileName", nombreArchivoDatos, "C:/Users/Alexandre/Documents/RIBCRA/DATOSFONDOS.txt", nil},
// 				{"BuildRIDatosFondo", requestDatosFondos, []ribcradtos.RiDatosFondos{}, fmt.Errorf("no se pudo cargar ri datos fondos")},
// 				{"RemoveFile", "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt", nil, fmt.Errorf("error al eliminar")},
// 			},
// 		},
// 		{"Debe Retornar un error si no se puede crear el archivo de detalle",
// 			requestValido, fmt.Errorf(commons.ERROR_FILE_CREATE), []callsMocked{
// 				{"CreateFileName", nombreArchivoCuentas, "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt", nil},
// 				{"BuildRICuentasCliente", requestCuentasClientes, []ribcradtos.RiCuentaCliente{}, nil},
// 				{"CreateFileName", nombreArchivoDatos, "C:/Users/Alexandre/Documents/RIBCRA/DATOSFONDOS.txt", nil},
// 				{"BuildRIDatosFondo", requestDatosFondos, []ribcradtos.RiDatosFondos{}, nil},
// 				{"CreateFileName", nombreArchivoDetalle, "C:/Users/Alexandre/Documents/RIBCRA/detalle.xml", nil},
// 				{"CreateFile", "C:/Users/Alexandre/Documents/RIBCRA/detalle.xml", &os.File{}, fmt.Errorf(commons.ERROR_FILE_CREATE)},
// 				{"RemoveFile", "C:/Users/Alexandre/Documents/RIBCRA/CUENTASCLIENTES.txt", nil, fmt.Errorf("error al eliminar")},
// 				{"RemoveFile", "C:/Users/Alexandre/Documents/RIBCRA/DATOSFONDOS.txt", nil, fmt.Errorf("error al eliminar")},
// 			},
// 		},
// 	}

// 	return
// }

// func TestInformacionSupervision(t *testing.T) {
// 	table := _inicializarInformacionSupervision()
// 	mockRepository := new(mockrepository.MockRepositoryAdministracion)
// 	mockApiLinkService := new(mockservice.MockApiLinkService)
// 	mockCommonsService := new(mockservice.MockCommonsService)

// 	service := administracion.NewService(mockRepository, mockApiLinkService, mockCommonsService)

// 	for _, v := range table {

// 		t.Run(v.Nombre, func(t *testing.T) {

// 			for _, m := range v.MetodosMocked {
// 				if m.Nombre == "CreateFileName" {
// 					mockCommonsService.On(m.Nombre, m.Request).Return(m.Response).Maybe()
// 				}
// 				if m.Nombre == "RemoveFile" || m.Nombre == "CreateFile" {
// 					mockCommonsService.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Maybe()
// 				}
// 				mockRepository.On(m.Nombre, m.Request).Return(m.Response, m.Erro).Once()
// 			}

// 			err := service.BuildInformacionSupervision(v.Request)

// 			if err != nil {
// 				assert.Equal(t, v.Erro.Error(), err.Error())
// 			}

// 		})
// 	}
// }
