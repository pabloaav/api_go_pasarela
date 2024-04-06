package apilink_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linktransferencia"
	"github.com/stretchr/testify/assert"
)

//FIXME ES MUY IMPORTANTE QUE SE GUARDEN ESTE ARCHIVO EN OTRO LADO O LO ELIMINE PORQUE SOLO SE USAN
//PARA LA HOMOLOGACIÓN Y HACEN PETICIONES A LA APILINK. NO SE QUE EN PRODUCCIÓN SE CORRA ESTE TEST
//Y SE HAGAN TRANFERENCIAS INEXISTENTES
type TablePostTransferencias struct {
	Nombre  string
	Erro    error
	Request RequestPostTransferencias
}

type RequestPostTransferencias struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linktransferencia.RequestTransferenciaCreateLink
	Token           string
}

func _inicializarCrearTransferenciaApiLink() (table []TablePostTransferencias) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linktransferencia.RequestTransferenciaCreateLink{
		Origen: linktransferencia.OrigenTransferenciaLink{
			Cbu: "0110599520000003855199",
		},
		Destino: linktransferencia.DestinoTransferenciaLink{
			Cbu:            "0290053700000000268714",
			AliasCbu:       "",
			EsMismoTitular: false,
		},
		Importe:    250,
		Moneda:     "ARS",
		Motivo:     "FAC",
		Referencia: "HomoApi",
	}

	requestValido := RequestPostTransferencias{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           valorValido,
		Token:           token.AccessToken,
	}

	// requestScopeInvalido := RequestPostTransferencias{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           valorValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	requestAliasMismoTitular := requestValido
	requestAliasMismoTitular.Valor.Destino.Cbu = ""
	requestAliasMismoTitular.Valor.Destino.AliasCbu = "Marzo2021"

	requestOtroTitular := requestValido
	requestOtroTitular.Valor.Destino.Cbu = "3890001130005274052211"
	requestOtroTitular.Valor.Destino.EsMismoTitular = false

	requestAliasOtroTitular := requestValido
	requestAliasOtroTitular.Valor.Destino.Cbu = ""
	requestAliasOtroTitular.Valor.Destino.AliasCbu = "PruebasDEBIN"
	requestAliasOtroTitular.Valor.Destino.EsMismoTitular = false

	requestRequerimientoInvalido := requestValido
	requestRequerimientoInvalido.RequerimientoId = "125552233-22233115-853333"

	requestTokenExpirado := requestValido
	requestTokenExpirado.Token = "eyJraWQiOiJSZWRMaW5rIiwiYWxnIjoiSFM1MTIifQ.eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiVFJBTlNGRVJFTkNJQVNfSU5NRURJQVRBUyIsImF1ZCI6ImQuYXBpLnJlZGxpbmsuY29tLmFyL3JlZGxpbmsvc2IvIiwiZXhwIjoxNjI5NTAwODA4LCJpYXQiOjE2Mjk0NjQ4MDh9.ojOHCo-5mbQxHrAI3FBupxw2-K4eFt1saKvGAu9z-3FLog0O4aNXEruoT8eNuPE3MOOWCJTWR6nnlgMSUHW8pg"

	requestCbuOrigenInvalido := requestValido
	requestCbuOrigenInvalido.Valor.Origen.Cbu = "01105995200000038551999"

	requestCbuDestinoInvalido := requestValido
	requestCbuDestinoInvalido.Valor.Destino.Cbu = "02900230100000500002933"

	requestAliasDestinoInvalido := requestValido
	requestAliasDestinoInvalido.Valor.Destino.Cbu = ""
	requestAliasDestinoInvalido.Valor.Destino.AliasCbu = "PruebasDEBIN.INVALIDO."

	requestMonedaInvalida := requestValido
	requestMonedaInvalida.Valor.Moneda = "REAL"

	requestMotivoInvalido := requestValido
	requestMotivoInvalido.Valor.Motivo = "MOTIVOINVALIDO"

	requestReferenciaInvalida := requestValido
	requestReferenciaInvalida.Valor.Referencia = "REFERENCIAINVALIDA_"

	requestRequerimientoVacio := requestValido
	requestRequerimientoVacio.RequerimientoId = ""

	requestTokenVacio := requestValido
	requestTokenVacio.Token = ""

	requestCbuOrigenVacio := requestValido
	requestCbuOrigenVacio.Valor.Origen.Cbu = ""

	requestCbuAliasDestinoVacio := requestValido
	requestCbuAliasDestinoVacio.Valor.Destino.Cbu = ""

	requestMonedaVacia := requestValido
	requestMonedaVacia.Valor.Moneda = ""

	requestMotivoVacio := requestValido
	requestMotivoVacio.Valor.Motivo = ""

	requestReferenciaVacia := requestValido
	requestReferenciaVacia.Valor.Referencia = ""

	requestCbuAliasDestino := requestValido
	requestCbuAliasDestino.Valor.Destino.AliasCbu = "DIA.HERMOSO"

	requestCbuOrigenDestinoIguales := requestValido
	requestCbuOrigenDestinoIguales.Valor.Destino.Cbu = "0110599520000003855199"

	requestMismoCbu := requestValido
	requestMismoCbu.Valor.Destino.Cbu = "0110599520000003855199"

	requestAliansInhabilitado := requestValido
	requestAliansInhabilitado.Valor.Destino.EsMismoTitular = false
	requestAliansInhabilitado.Valor.Destino.Cbu = ""
	requestAliansInhabilitado.Valor.Destino.AliasCbu = "Marbella.ES"

	requestAlianCbuInhabilitado := requestValido
	requestAlianCbuInhabilitado.Valor.Destino.EsMismoTitular = false
	requestAlianCbuInhabilitado.Valor.Destino.Cbu = ""
	requestAlianCbuInhabilitado.Valor.Destino.AliasCbu = "ApiLink2023"

	requestCbuInhabilitado := requestValido
	requestCbuInhabilitado.Valor.Destino.Cbu = "0110616530061602678801"

	requestMotivoDiferente := requestValido
	requestMotivoDiferente.Valor.Motivo = "PCT"

	requestCuentaDiferenteMoneda := requestValido
	requestCuentaDiferenteMoneda.Valor.Destino.Cbu = "0110599531000027805170"

	requestCbuCvu := requestValido
	requestCbuCvu.Valor.Destino.Cbu = "0000262402008276351202"

	requestCbuAliasCvu := requestValido
	requestCbuAliasCvu.Valor.Destino.Cbu = ""
	requestCbuAliasCvu.Valor.Destino.AliasCbu = "AliasNuevoXC123"

	requestCbuInexistente := requestValido
	requestCbuInexistente.Valor.Destino.Cbu = "0110616530061602678801"

	requestImporteMaximo := requestValido
	requestImporteMaximo.Valor.Importe = 999999999999999999

	requestSaldoInsuficiente := requestValido
	requestSaldoInsuficiente.Valor.Importe = 92640385200

	requestCbuNoVinculadoApi := requestValido
	requestCbuNoVinculadoApi.Valor.Origen.Cbu = "0940099372008402190012"

	requestMonedaDolar := requestValido
	requestMonedaDolar.Valor.Moneda = linkdtos.Dolar

	table = []TablePostTransferencias{
		// {"2 Transferencia Inmediata mismo titular", nil, requestValido},
		// {"3 Transferencia Inmediata Alias mismo titular", nil, requestAliasMismoTitular},
		//  {"4 Transferencia Inmediata otro titular", nil, requestOtroTitular},
		// {"5 Transferencia Inmediata Alias otro titular", nil, requestAliasOtroTitular},
		// {"6 Requerimiento invalido", fmt.Errorf("X-IdRequerimiento con formato invalid"), requestRequerimientoInvalido},
		// {"7 Token expirado", fmt.Errorf("Unauthorized"), requestTokenExpirado},
		// {"8 Escope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
		// {"9 Cbu de origen con formato invalido", fmt.Errorf("cbu origen con formato invalid"), requestCbuOrigenInvalido},
		// {"10 Cbu de destino con formato invalido", fmt.Errorf("cbu destino con formato invalido"), requestCbuDestinoInvalido},
		// {"11 Alias de destino con formato invalido", fmt.Errorf("alias destino con formato invalido"), requestAliasDestinoInvalido},
		// {"14 Moneda invalida", fmt.Errorf("moneda con formato invalido"), requestMonedaInvalida},
		// {"15 Motivo con formato invalido", fmt.Errorf("motivo con formato invalido"), requestMotivoInvalido},
		// {"16 Referencia con formato invalido", fmt.Errorf("referencia con formato invalido. Verificar el parametro"), requestReferenciaInvalida},
		// //Las prueba con id cliente se debe modificar el valor en variables del sitema
		// // {"17 X IBM Client id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// {"18 X IBM Client id vacio", fmt.Errorf("Unauthorized"), requestValido},
		// {"19 Requerimiento Vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), requestRequerimientoVacio},
		// {"20 Token vacio", fmt.Errorf("Unauthorized"), requestTokenVacio},
		// {"21 Cbu origen vacio", fmt.Errorf("Se debe indicar cbu origen"), requestCbuOrigenVacio},
		// {"22 Cbu y alias destino vacios", fmt.Errorf("Uno (y solo uno) de los siguientes campos debe estar completo"), requestCbuAliasDestinoVacio},
		// {"25 Moneda vacia", fmt.Errorf("Se debe indicar moneda"), requestMonedaVacia},
		// {"26 Motivo vacio", fmt.Errorf("Se debe indicar motivo"), requestMotivoVacio},
		// {"27 Referencia vacia", fmt.Errorf("Se debe indicar referencia"), requestReferenciaVacia},
		// {"28 Importe mayor maximo", fmt.Errorf("Ud. ha excedido el importe máximo del sistema. Por favor, ingrese uno menor"), requestImporteMaximo},
		// {"29 Saldo insuficiente cuenta origen", fmt.Errorf("Se ha alcanzado el límite de operación diario para la cuenta de origen"), requestSaldoInsuficiente},
		// {"30 Cbu no vinculado a la Api", fmt.Errorf("No posee Cuenta con autorizacion para realizar esta operacio"), requestCbuNoVinculadoApi},
		// {"31 mismo cbu origen y destino", fmt.Errorf("La operación no está permitida para las cuentas informadas"), requestCbuOrigenDestinoIguales},
		// {"32 Cbu y alias destino con valor", fmt.Errorf("Uno (y solo uno) de los siguientes campos debe estar completo"), requestCbuAliasDestino},
		// {"FE-034 Alias inhabilitada", fmt.Errorf("El alias de destino con el que se quiere operar es invalido o inexistente"), requestAliansInhabilitado},
		// {"36 Cbu iniexistente", fmt.Errorf("La cuenta destino se encuentra inhabilitada para transferencias"), requestCbuInhabilitado},
		// {"FE-036 Alias con CBU inhabilitada", fmt.Errorf("El alias de destino con el que se quiere operar es invalido o inexistente"), requestAlianCbuInhabilitado},
		// {"FE-037 Cbu valido iniexistente", fmt.Errorf("La búsqueda del CBU solicitado no arrojó resultados"), requestCbuInhabilitado},
		// {"39 motivo diferente no permitido", fmt.Errorf("El motivo de transferencia es inválido"), requestMotivoDiferente},
		// {"39 cuentas de diferentes monedas", fmt.Errorf("Las cuentas origen y destino no tienen el mismo tipo de moneda"), requestCuentaDiferenteMoneda},
		// {"41 transferencia de cbu a cvu", fmt.Errorf("La cuenta ingresada corresponde a una CVU"), requestCbuCvu},
		// {"42 transferencia de cbu a alias cvu", fmt.Errorf("ALIAS_CUENTA_VIRTUAL"), requestCbuAliasCvu},
		// {"43 Realizar una transferencia inmediata Moneda USD", fmt.Errorf("La operación es invalida"), requestMonedaDolar},

	}

	return
}

func TestPostTransferencia(t *testing.T) {

	table := _inicializarCrearTransferenciaApiLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.CreateTransferenciaApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.NotEmpty(t, resp.NumeroReferenciaBancaria)
				t.Log(resp.String())
			}

		})
	}

}

type TableGetTransferencia struct {
	Nombre  string
	Erro    error
	Request RequestGetTransferencia
}

type RequestGetTransferencia struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linktransferencia.RequestGetTransferenciaLink
	Token           string
}

func _inicializarGetTransferenciaApiLink() (table []TableGetTransferencia) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linktransferencia.RequestGetTransferenciaLink{
		Cbu:                      "0110599520000003855199",
		NumeroReferenciaBancaria: "2021090815102317200000762",
	}

	requestValido := RequestGetTransferencia{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           valorValido,
		Token:           token.AccessToken,
	}

	// requestScopeInvalido := RequestGetTransferencia{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           valorValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	requestTokenInvalido := requestValido
	requestTokenInvalido.Token = "Bearer eyJraWQiOiJSZWRMaW5rIiwiYWxnIjoiSFM1MTIifQ.eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiVFJBTlNGRVJFTkNJQVNfSU5NRURJQVRBUyIsImF1ZCI6ImQuYXBpLnJlZGxpbmsuY29tLmFyL3JlZGxpbmsvc2IvIiwiZXhwIjoxNjI5NTAwODA4LCJpYXQiOjE2Mjk0NjQ4MDh9.ojOHCo-5mbQxHrAI3FBupxw2-K4eFt1saKvGAu9z-3FLog0O4aNXEruoT8eNuPE3MOOWCJTWR6nnlgMSUHW8pg"

	requestCbuDiferenteReferencia := requestValido
	requestCbuDiferenteReferencia.Valor.NumeroReferenciaBancaria = "202109081520537590000076"

	requestRequerimientoInvalido := requestValido
	requestRequerimientoInvalido.RequerimientoId = "88uuutjjd-999kkkfm-5900kkkd"

	requestCbuInvalido := requestValido
	requestCbuInvalido.Valor.Cbu = "0340219999999999437001"

	requestReferenciaInvalida := requestValido
	requestReferenciaInvalida.Valor.NumeroReferenciaBancaria = "034021_444"

	requestRequerimientoVacio := requestValido
	requestRequerimientoVacio.RequerimientoId = ""

	requestTokenVacio := requestValido
	requestTokenVacio.Token = ""

	requestCbuVacio := requestValido
	requestCbuVacio.Valor.Cbu = ""

	requestReferenciaVacia := requestValido
	requestReferenciaVacia.Valor.NumeroReferenciaBancaria = ""

	requestCbuOtroCliente := requestValido
	requestCbuOtroCliente.Valor.Cbu = "0110046430004601454801"

	table = []TableGetTransferencia{
		// {"2 Buscar transferencia por numero referencia", nil, requestValido},
		// {"3 Cbu no coincide con numero referencia", fmt.Errorf("No Content"), requestCbuDiferenteReferencia},
		// {"5 Requerimiento con formato invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), requestRequerimientoInvalido},
		// //Las prueba con id cliente se debe modificar el valor en variables del sitema
		// // {"6 X IBM Client id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// {"7 Autorizacion Formato Invalido", fmt.Errorf("401"), requestTokenInvalido},
		// {"8 Cbu invalido", fmt.Errorf("cbu con formato invalido"), requestCbuInvalido},
		// {"9 Escope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
		// {"10 Referencia Bancaria formato invalido", fmt.Errorf("numeroReferenciaBancaria con formato invalido"), requestReferenciaInvalida},
		// {"11 Requerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), requestRequerimientoVacio},
		// // {"7 X IBM Client id vacio", fmt.Errorf("Unauthorized"), requestValido},
		// {"13 Token vacio", fmt.Errorf("Unauthorized"), requestTokenVacio},
		// {"14 Cbu vacio", fmt.Errorf("Se debe indicar cbu"), requestCbuVacio},
		// {"15 Referencia vacia", fmt.Errorf("Se debe indicar numeroReferenciaBancaria"), requestReferenciaVacia},
		{"16 cbu que no corresponde al cliente id que invoca", fmt.Errorf("No posee Cuenta con autorizacion para realizar esta operacion"), requestCbuOtroCliente},
	}

	return
}

func TestGetTransferenciaReferencia(t *testing.T) {

	table := _inicializarGetTransferenciaApiLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetTransferenciaApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.NotEmpty(t, resp.NumeroReferenciaBancaria)
				t.Log(resp.String())
			}

		})
	}

}

type TableGetTransferencias struct {
	Nombre  string
	Erro    error
	Request RequestGetTransferencias
}

type RequestGetTransferencias struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linktransferencia.RequestGetTransferenciasLink
	Token           string
}

func _inicializarGetTransferenciasApiLink() (table []TableGetTransferencias) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linktransferencia.RequestGetTransferenciasLink{
		Cbu:        "0110599520000003855199",
		Tamanio:    linkdtos.DiezTransf,
		Pagina:     1,
		FechaDesde: time.Now().Add(time.Hour * -24),
		FechaHasta: time.Now(),
	}

	requestValido := RequestGetTransferencias{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           valorValido,
		Token:           token.AccessToken,
	}

	requestPeriodoSinTransferencia := requestValido
	requestPeriodoSinTransferencia.Valor.FechaDesde = time.Now()
	requestPeriodoSinTransferencia.Valor.FechaHasta = time.Now().Add(time.Millisecond * 100)

	requestRequerimientoInvalido := requestValido
	requestRequerimientoInvalido.RequerimientoId = "cccccc-vv998889-9iiitkkf"

	requestTokenInvalido := requestValido
	requestTokenInvalido.Token = "cccccc-vv998889-9iiitkkf"

	requestCbuInvalido := requestValido
	requestCbuInvalido.Valor.Cbu = "01105995200000038551991"

	requestTamanioInvalido := requestValido
	requestTamanioInvalido.Valor.Tamanio = "2"

	requestPaginaInvalido := requestValido
	requestPaginaInvalido.Valor.Pagina = 0

	requestRequerimientoVacio := requestValido
	requestRequerimientoVacio.RequerimientoId = ""

	requestTokenVacio := requestValido
	requestTokenVacio.Token = ""

	requestCbuVacio := requestValido
	requestCbuVacio.Valor.Cbu = ""

	requestTamanioVacio := requestValido
	requestTamanioVacio.Valor.Tamanio = ""

	// requestScopeInvalido := RequestGetTransferencias{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           valorValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	requestFechaDesdeMayorFechaHasta := requestValido
	requestFechaDesdeMayorFechaHasta.Valor.FechaDesde = time.Now().Add(time.Hour * 48)
	requestFechaDesdeMayorFechaHasta.Valor.FechaHasta = time.Now()

	requestFechaDesdeMayorFechaActual := requestValido
	requestFechaDesdeMayorFechaActual.Valor.FechaDesde = time.Now().Add(time.Hour * 48)
	requestFechaDesdeMayorFechaHasta.Valor.FechaHasta = time.Now().Add(time.Hour * 96)

	requestRangoFechaExcedeLimite := requestValido
	requestRangoFechaExcedeLimite.Valor.FechaDesde = time.Now().Add(time.Hour * -2160)
	requestFechaDesdeMayorFechaHasta.Valor.FechaHasta = time.Now()

	table = []TableGetTransferencias{
		{"2 Buscar transferencia por fecha", nil, requestValido},
		// {"3 Buscar por fecha que no tenga transferencia", fmt.Errorf("no se encontraron resultados para la consulta"), requestPeriodoSinTransferencia},
		// {"4 Buscar transferencia tamanio valido", nil, requestValido},
		// {"6 Fecha desde menor que fecha hasta", nil, requestValido},
		// {"7 Fecha hasta mayor que fecha desde", nil, requestValido},
		// {"8 Requerimiento invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), requestRequerimientoInvalido},
		// //Las prueba con id cliente se debe modificar el valor en variables del sitema
		// // {"9 X IBM Client id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// {"10 Token invalido invalido", fmt.Errorf("401"), requestTokenInvalido},
		// {"11 cbu invalido", fmt.Errorf("cbu con formato invalido"), requestCbuInvalido},
		// {"12 tamanio invalido", fmt.Errorf("tamanio con formato invalido"), requestTamanioInvalido},
		// {"13 pagina invalida", fmt.Errorf("pagina con formato invalido"), requestPaginaInvalido},
		// //La fecha tiene que ser modificada en el repositorio
		// // {"14 Fecha Desde invalida", fmt.Errorf("fechaDesde con formato invalido"), requestValido},
		// // {"15 Fecha Hasta invalida", fmt.Errorf("fechaHasta con formato invalido"), requestValido},
		// {"16 Requerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), requestRequerimientoVacio},
		// //Las prueba con id cliente se debe modificar el valor en variables del sitema
		// // {"17 X IBM Client id vacio", fmt.Errorf("Unauthorized"), requestValido},
		// {"18 Token vacio", fmt.Errorf("Unauthorized"), requestTokenVacio},
		// {"19 Cbu vacio", fmt.Errorf("Se debe indicar cbu"), requestCbuVacio},
		// {"20 Tamanio vacio", nil, requestTamanioVacio},
		// {"24 Scope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
		// {"25 Fecha desde mayor que fecha hasta", fmt.Errorf("Las fechas desde y hasta tienen valores incorrectos"), requestFechaDesdeMayorFechaHasta},
		// {"26 Fecha desde mayor que fecha actual", fmt.Errorf("Las fechas desde y hasta tienen valores incorrectos"), requestFechaDesdeMayorFechaActual},
		// {"27 Rango de fecha excede al limite", fmt.Errorf("No se pudo realizar la consulta ya que la cantidad de días entre la fecha desde y hasta es mayor a 30 días. Acotar el rango de fechas"), requestRangoFechaExcedeLimite},
	}

	return
}

func TestGetTransferenciasPorFecha(t *testing.T) {

	table := _inicializarGetTransferenciasApiLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetTransferenciasApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.Len(t, resp.Transferencias, 1)
				t.Log(resp.String())
			}
		})
	}

}
