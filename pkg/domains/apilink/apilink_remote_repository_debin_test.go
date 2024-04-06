package apilink_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkdebin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func InicializarDebines(scope linkdtos.EnumScopeLink) (repositorio apilink.RemoteRepository, uuidValido string, token linkdtos.TokenLink) {

	HTTPClient := http.DefaultClient
	clienteSQL := database.NewMySQLClient()
	utilRepository := util.NewUtilRepository(clienteSQL)
	utilService := util.NewUtilService(utilRepository, HTTPClient)
	repositorio = apilink.NewRemote(HTTPClient, utilService)

	uuidValido = uuid.NewV4().String()
	scopesValidos := []linkdtos.EnumScopeLink{scope}
	token, _ = repositorio.GetTokenApiLink(uuidValido, scopesValidos)

	return
}

func DebinValido() linkdebin.RequestDebinCreateLink {
	return linkdebin.RequestDebinCreateLink{
		Comprador: linkdebin.CompradorCreateDebinLink{
			Cuit: "30709585254", Cuenta: linkdebin.CuentaLink{
				Cbu: "3890001130005274052211", AliasCbu: "",
			},
		},
		Vendedor: linkdebin.VendedorCreateLink{Cbu: "0110599520000003855199", AliasCbu: ""},
		Debin: linkdebin.DebinCreateLink{
			ComprobanteId:         "123456",
			EsCuentaPropia:        false,
			Concepto:              "VAR",
			TiempoExpiracion:      1200,
			Importe:               6852,
			Moneda:                "ARS",
			Recurrente:            false,
			Descripcion:           "prueba",
			DescripcionPrestacion: "", //No se debe agregar valores descripcion prestación cuando no sea debin recurrencia
		},
	}
}

func TestCrearDebinesRecurrenteFalse(t *testing.T) {
	debinValido := DebinValido()
	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debinValido.String())
	resp, _ := repositorio.CreateDebinApiLink(uuidValido, debinValido, token.AccessToken)
	t.Log(resp.String())

	assert.Equal(t, resp.Estado, linkdtos.EnumEstadoDebin("INICIADO"))

}
func TestCrearDebinesDescripcionPrestacionValida(t *testing.T) {
	debinValido := DebinValido()
	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debinValido.String())
	resp, _ := repositorio.CreateDebinApiLink(uuidValido, debinValido, token.AccessToken)
	t.Log(resp.String())

	assert.Equal(t, resp.Estado, linkdtos.EnumEstadoDebin("INICIADO"))

}

func TestCrearDebinesDescripcionPrestacionInvalida(t *testing.T) {

	debin := DebinValido()
	for i := 0; i < 100; i++ {
		debin.Debin.DescripcionPrestacion += "c"
	}

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "descripcionPrestacion debin con formato invalido")

}

func TestCrearDebinesDescripcionPrestacionInexistente(t *testing.T) {

	debin := DebinValido()
	for i := 0; i < 100; i++ {
		debin.Debin.DescripcionPrestacion = ""
	}

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	resp, _ := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)
	t.Log(resp.String())

	assert.Equal(t, resp.Estado, linkdtos.EnumEstadoDebin("ACEPTADO"))

}

func TestCrearDebinesAliasCbuCompradorInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Comprador.Cuenta.Cbu = ""
	debin.Comprador.Cuenta.AliasCbu = "ALIAS_INVALIDO_COMPRADOR"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "aliasCbu cuenta comprador con formato invalido")

}

func TestCrearDebinesCbuCompradorInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Comprador.Cuenta.Cbu = "01403235014200002837"
	debin.Comprador.Cuenta.AliasCbu = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "cbu cuenta comprador con formato invalido")

}

func TestCrearDebinesCuitInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Comprador.Cuit = "306677543011"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "cuit comprador con formato invalido")

}

func TestCrearDebinesCbuVendedorInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Vendedor.Cbu = "011059952000000385519"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "cbu vendedor con formato invalido")

}

func TestCrearDebinesAliasCbuVendedorInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Vendedor.Cbu = ""
	debin.Vendedor.AliasCbu = "SAN_TERESA"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "aliasCbu vendedor con formato invalido")

}

func TestCrearDebinesRequerimientoInvalido(t *testing.T) {

	debin := DebinValido()

	uuidInvalido := "1234567891011"

	repositorio, _, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidInvalido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "X-IdRequerimiento con formato invalido")

}

func TestCrearDebinesAliasYCbuCompradorValidos(t *testing.T) {

	debin := DebinValido()

	debin.Comprador.Cuenta.AliasCbu = "alias.isa1"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Uno (y solo uno) de los siguientes campos debe estar completo")

}

func TestCrearDebinesAliasYCbuVendedorValidos(t *testing.T) {

	debin := DebinValido()

	debin.Vendedor.AliasCbu = "SANTERESA"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Uno (y solo uno) de los siguientes campos debe estar completo")

}

func TestCrearDebinesComprobanteInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Debin.ComprobanteId = "   12BB"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "comprobanteId debin con formato invalido. Verificar el parametro")

}

func TestCrearDebinesConceptoInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Debin.Concepto = "INVALIDO"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "concepto debin con formato invalido")

}

func TestCrearDebinesTiempoExpiracionInvalido(t *testing.T) {

	debin := DebinValido()
	debin.Debin.TiempoExpiracion = 6000

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "tiempoExpiracion debin con formato invalido. Verificar el parametro")

}

func TestCrearDebinesDescripcionInvalida(t *testing.T) {

	debin := DebinValido()
	for i := 0; i < 101; i++ {
		debin.Debin.Descripcion += "X"
	}

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "descripcion debin con formato invalido")

}

func TestCrearDebinesMonedaInvalida(t *testing.T) {

	debin := DebinValido()
	debin.Debin.Moneda = "REAL"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "moneda debin con formato invalido")

}

func TestCrearDebinesAliasCbuYCbuCompradorVacios(t *testing.T) {

	debin := DebinValido()
	debin.Comprador.Cuenta.AliasCbu = ""
	debin.Comprador.Cuenta.Cbu = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Uno (y solo uno) de los siguientes campos debe estar completo")

}

func TestCrearDebinesCuitCompradorVacio(t *testing.T) {

	debin := DebinValido()
	debin.Comprador.Cuit = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Se debe indicar cuit comprador")

}

func TestCrearDebinesAliasCbuYCbuVendedorVacios(t *testing.T) {

	debin := DebinValido()
	debin.Vendedor.AliasCbu = ""
	debin.Vendedor.Cbu = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Uno (y solo uno) de los siguientes campos debe estar completo: <request.body.vendedor.cbu>, <request.body.vendedor.aliasCbu>")

}

func TestCrearDebinesRequerimientoVacio(t *testing.T) {

	debin := DebinValido()

	repositorio, _, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink("", debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Se debe indicar X-IdRequerimiento")

}

func TestCrearDebinesConceptoVacio(t *testing.T) {

	debin := DebinValido()
	debin.Debin.Concepto = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Se debe indicar concepto debin")

}

func TestCrearDebinesMonedaVacia(t *testing.T) {

	debin := DebinValido()
	debin.Debin.Moneda = ""

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "Se debe indicar moneda debin")

}

// TODO Para lograr eso hay que ir en variables de sistema y modificar IDCLIENTLINK
func TestCrearDebinesXIBMClienteIdInvalido(t *testing.T) {

	debin := DebinValido()

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, erro := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, erro.Error(), "Unauthorized")

}

func TestCrearDebinesTokenInvalido(t *testing.T) {

	debin := DebinValido()

	repositorio, uuidValido, _ := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, "")

	assert.Contains(t, err.Error(), "401")

}

func TestCrearDebinesScopoInvalidoInvalido(t *testing.T) {

	debin := DebinValido()

	repositorio, uuidValido, token := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "403")

}

func TestCrearDebinesCbuNoRegistrado(t *testing.T) {

	debin := DebinValido()
	debin.Vendedor.Cbu = "0110046430004601454801"

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	t.Log(debin.String())
	_, err := repositorio.CreateDebinApiLink(uuidValido, debin, token.AccessToken)

	assert.Contains(t, err.Error(), "No posee Cuenta con autorizacion para realizar esta operacion")

}

//*******************************************************------------------------*************************************************//
//******************************************************BUSCAR DEBINES RECIBIDOS*************************************************//
//******************************************************------------------------*************************************************//

type TableGetDebines struct {
	Nombre  string
	Erro    error
	Request RequestDebines
}

type RequestDebines struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linkdebin.RequestGetDebinesLink
	Token           string
}

func _inicializarBuscarDebinesRecibidos() (table []TableGetDebines) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenEscopeInvalido := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	RequestGetDebinValido := linkdebin.RequestGetDebinesLink{
		Pagina:      5,
		Tamanio:     linkdtos.Diez,
		Cbu:         "0110599520000003855199",
		Estado:      "INICIADO",
		FechaDesde:  time.Date(2021, time.Month(9), 01, 0, 0, 0, 0, time.Local),
		FechaHasta:  time.Date(2021, time.Month(9), 30, 0, 0, 0, 0, time.Local),
		EsComprador: true,
		Tipo:        linkdtos.DebinDefault,
	}

	RequestDebinValido := RequestDebines{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           RequestGetDebinValido,
		Token:           token.AccessToken,
	}

	// RequestDebinScopeInvalido := RequestDebines{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           RequestGetDebinValido,
	// 	Token:           tokenEscopeInvalido.AccessToken,
	// }

	RequestRequerimientoInvalido := RequestDebinValido
	RequestRequerimientoInvalido.RequerimientoId = "1233662-22556632cc-ccc855dc"

	RequestPaginaInvalida := RequestDebinValido
	RequestPaginaInvalida.Valor.Pagina = 0

	RequestTamanioInvalido := RequestDebinValido
	RequestTamanioInvalido.Valor.Tamanio = "25"

	RequestTipoInvalido := RequestDebinValido
	RequestTipoInvalido.Valor.Tipo = "DEBININVALIDO"

	RequestEstadoInvalido := RequestDebinValido
	RequestEstadoInvalido.Valor.Estado = "INICIADO_"

	RequestEstadoInexistente := RequestDebinValido
	RequestEstadoInexistente.Valor.Estado = "INICIADO_INVALIDO"

	RequestCbuInvalido := RequestDebinValido
	RequestCbuInvalido.Valor.Cbu = "01105995200000038551991"

	RequestRequerimientoVacio := RequestDebinValido
	RequestRequerimientoVacio.RequerimientoId = ""

	RequestCbuVacio := RequestDebinValido
	RequestCbuVacio.Valor.Cbu = ""

	RequestFechaDesdeVacia := RequestDebinValido
	RequestFechaDesdeVacia.Valor.FechaDesde = time.Time{}

	RequestFechaHastaVacia := RequestDebinValido
	RequestFechaHastaVacia.Valor.FechaHasta = time.Time{}

	RequestFechaInvalida := RequestDebinValido
	RequestFechaInvalida.Valor.FechaDesde = time.Date(2021, time.Month(9), 01, 0, 0, 0, 0, time.Local)

	RequestSinDebines := RequestDebinValido
	RequestSinDebines.Valor.FechaDesde = time.Date(2021, time.Month(8), 01, 0, 0, 0, 0, time.Local)
	RequestSinDebines.Valor.FechaHasta = time.Date(2021, time.Month(8), 15, 0, 0, 0, 0, time.Local)

	RequestTokenVacio := RequestDebinValido
	RequestTokenVacio.Token = ""

	table = []TableGetDebines{
		// {"2 Campo es comprador true", nil, RequestDebinValido},
		// {"3 IdRequerimiento con formato invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), RequestRequerimientoInvalido},
		// {"4 Pagina con formato invalido", fmt.Errorf("pagina con formato invalido"), RequestPaginaInvalida},
		// {"5 Tamanio con formato invalido", fmt.Errorf("tamanio con formato invalido"), RequestTamanioInvalido},
		// {"6 Tipo con formato invalido", fmt.Errorf("tipo con formato invalido"), RequestTipoInvalido},
		////Para que se pueda ejecuta la prueba de fecha hay que comentar el codigo que convierte las fechas GetDebinesApiLink repositorio
		// // {"7 Fecha desde Formato invalido", fmt.Errorf("fechadesde con formato invalido"), RequestDebinValido},
		// {"8 Fecha hasta Formato invalido", fmt.Errorf("fechahasta con formato invalido"), RequestDebinValido},
		// {"9 Estado con formato invalido", fmt.Errorf("estado con formato invalido"), RequestEstadoInvalido},
		// {"10 Estado con valor inexistente", fmt.Errorf("estado con formato invalido"), RequestEstadoInexistente},
		// {"12 Cbu con formato invalido", fmt.Errorf("cbu con formato invalido"), RequestCbuInvalido},
		// {"13 Requerimiento Vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), RequestRequerimientoVacio},
		// {"15 Cbu Vacio", fmt.Errorf("Se debe indicar cbu"), RequestCbuVacio},
		// {"17 Fecha desde vacia", fmt.Errorf("Se debe indicar fechadesde"), RequestFechaDesdeVacia},
		// // {"18 Fecha hasta vacia", fmt.Errorf("Se debe indicar fechahasta"), RequestFechaHastaVacia},
		// {"19 Fecha con valor invalido", fmt.Errorf("Las fechas desde y hasta tienen valores incorrectos"), RequestFechaInvalida},
		// {"20 Fecha con valor valido", nil, RequestSinDebines},
		// //Para que funcione este test hay que modificar las variables de entorno antes de enviar
		// // {"21 ClienteId con valor invalido", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"22 ClienteId con valor vacio", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"25 Token Vacio", fmt.Errorf("Unauthorized"), RequestTokenVacio},
		// {"26 Token Vacio", fmt.Errorf("Forbidden"), RequestDebinScopeInvalido},
	}

	return
}

func TestBuscarDebinesRecibidos(t *testing.T) {

	table := _inicializarBuscarDebinesRecibidos()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetDebinesApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				if len(resp.Debines) > 0 {
					assert.Equal(t, resp.Debines[0].Tipo, linkdtos.DebinDefault)
				}
				t.Log(resp.String())
			}

		})
	}

}

func _inicializarBuscarDebinesGenerados() (table []TableGetDebines) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenEscopeInvalido := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	RequestGetDebinValido := linkdebin.RequestGetDebinesLink{
		Pagina:      1,
		Tamanio:     "",
		Cbu:         "0110599520000003855199",
		Estado:      "",
		FechaDesde:  time.Date(2021, time.Month(11), 28, 1, 31, 35, 0, time.Local),
		FechaHasta:  time.Date(2021, time.Month(11), 29, 12, 8, 9, 0, time.Local),
		EsComprador: false,
		Tipo:        linkdtos.DebinDefault,
	}

	RequestDebinValido := RequestDebines{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           RequestGetDebinValido,
		Token:           token.AccessToken,
	}

	// RequestDebinScopeInvalido := RequestDebines{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           RequestGetDebinValido,
	// 	Token:           tokenEscopeInvalido.AccessToken,
	// }

	RequestRequerimientoInvalido := RequestDebinValido
	RequestRequerimientoInvalido.RequerimientoId = "1233662-22556632cc-ccc855dc"

	RequestPaginaInvalida := RequestDebinValido
	RequestPaginaInvalida.Valor.Pagina = 0

	RequestTamanioInvalido := RequestDebinValido
	RequestTamanioInvalido.Valor.Tamanio = "25"

	RequestTipoInvalido := RequestDebinValido
	RequestTipoInvalido.Valor.Tipo = "DEBININVALIDO"

	RequestEstadoInvalido := RequestDebinValido
	RequestEstadoInvalido.Valor.Estado = "1"

	RequestEstadoInexistente := RequestDebinValido
	RequestEstadoInexistente.Valor.Estado = "INICIADO_INVALIDO"

	RequestCbuInvalido := RequestDebinValido
	RequestCbuInvalido.Valor.Cbu = "011059952000000385519"

	RequestRequerimientoVacio := RequestDebinValido
	RequestRequerimientoVacio.RequerimientoId = ""

	RequestCbuVacio := RequestDebinValido
	RequestCbuVacio.Valor.Cbu = ""

	RequestFechaDesdeVacia := RequestDebinValido
	RequestFechaDesdeVacia.Valor.FechaDesde = time.Time{}

	RequestFechaHastaVacia := RequestDebinValido
	RequestFechaHastaVacia.Valor.FechaHasta = time.Time{}

	RequestFechaInvalida := RequestDebinValido
	RequestFechaInvalida.Valor.FechaDesde = time.Date(2021, time.Month(9), 01, 0, 0, 0, 0, time.Local)

	RequestTokenVacio := RequestDebinValido
	RequestTokenVacio.Token = ""

	table = []TableGetDebines{
		{"2 Campo es comprador false", nil, RequestDebinValido},
		// {"3 IdRequerimiento con formato invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), RequestRequerimientoInvalido},
		// {"4 Pagina con formato invalido", fmt.Errorf("pagina con formato invalido"), RequestPaginaInvalida},
		// {"5 Tamanio con formato invalido", fmt.Errorf("tamanio con formato invalido"), RequestTamanioInvalido},
		// {"6 Tipo con formato invalido", fmt.Errorf("tipo con formato invalido"), RequestTipoInvalido},
		////Para que se pueda ejecuta la prueba de fecha hay que comentar el codigo que convierte las fechas GetDebinesApiLink repositorio
		// // {"7 Fecha desde Formato invalido", fmt.Errorf("fechadesde con formato invalido"), RequestDebinValido},
		// // {"8 Fecha hasta Formato invalido", fmt.Errorf("fechahasta con formato invalido"), RequestDebinValido},
		// {"9 Estado con formato invalido", fmt.Errorf("estado con formato invalido"), RequestEstadoInvalido},
		// {"10 Estado con valor inexistente", fmt.Errorf("estado con formato invalido"), RequestEstadoInexistente},
		// {"12 Cbu con formato invalido", fmt.Errorf("cbu con formato invalido"), RequestCbuInvalido},
		// {"13 Requerimiento Vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), RequestRequerimientoVacio},
		// {"15 Cbu Vacio", fmt.Errorf("Se debe indicar cbu"), RequestCbuVacio},
		// // {"17 Fecha desde vacia", fmt.Errorf("Se debe indicar fechadesde"), RequestFechaDesdeVacia},
		// // {"18 Fecha hasta vacia", fmt.Errorf("Se debe indicar fechahasta"), RequestFechaHastaVacia},
		// // {"19 Fecha con valor invalido", fmt.Errorf("Las fechas desde y hasta tienen valores incorrectos"), RequestFechaInvalida},
		// // {"20 Fecha con valor valido", nil, RequestDebinValido},
		// //Para que funcione este test hay que modificar las variables de entorno antes de enviar
		// // {"21 ClienteId con valor invalido", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// // {"22 ClienteId con valor vacio", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"25 Token Vacio", fmt.Errorf("Unauthorized"), RequestTokenVacio},
		// {"26 Scope Invalido", fmt.Errorf("Forbidden"), RequestDebinScopeInvalido},
	}

	return
}

func TestBuscarDebinesGenerados(t *testing.T) {

	table := _inicializarBuscarDebinesGenerados()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetDebinesApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				if len(resp.Debines) > 0 {
					assert.Equal(t, linkdtos.Iniciado, resp.Debines[0].Estado)
				}
				t.Log(resp.String())
			}

		})
	}

}

type TableGetDebin struct {
	Nombre  string
	Erro    error
	Request RequestDebin
}

type RequestDebin struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linkdebin.RequestGetDebinLink
	Token           string
}

func _inicializarConsultarDebinPorId() (table []TableGetDebin) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	RequestGetDebinValido := linkdebin.RequestGetDebinLink{
		Cbu: "0110599520000003855199",
		Id:  "Z86VRPQ2GDOVZYY2GLY0M1",
	}

	RequestDebinValido := RequestDebin{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           RequestGetDebinValido,
		Token:           token.AccessToken,
	}

	// RequestDebinScopeInvalido := RequestDebin{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           RequestGetDebinValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	RequestCbuInvalido := RequestDebinValido
	RequestCbuInvalido.Valor.Cbu = "99110599520000003855199"

	RequestIdInvalido := RequestDebinValido
	RequestIdInvalido.Valor.Id = "x_555llff"

	RequestRequerimientoInvalido := RequestDebinValido
	RequestRequerimientoInvalido.RequerimientoId = "ccvvv_fff---_ft666"

	RequestCbuVacio := RequestDebinValido
	RequestCbuVacio.Valor.Cbu = ""

	RequestRequerimientoVacio := RequestDebinValido
	RequestRequerimientoVacio.RequerimientoId = ""

	RequestDebinEliminado := RequestDebinValido
	RequestDebinEliminado.Valor.Id = "746YGOW9MJ58M3P9EXD8J5"

	RequestDebinVencido := RequestDebinValido
	RequestDebinVencido.Valor.Id = "V8D0Q619LY8O8QV27JZ5RG"

	RequestCbuOtroCliente := RequestDebinValido
	RequestCbuOtroCliente.Valor.Cbu = "3890001130005274052211"

	RequestTokenInvalido := RequestDebinValido
	RequestTokenInvalido.Token = ""

	table = []TableGetDebin{
		// {"2 Campos Correctos", nil, RequestDebinValido},
		// {"4 Cbu con formato invalido", fmt.Errorf("cbu con formato invalid"), RequestCbuInvalido},
		// {"5 Id con formato invalido", fmt.Errorf("id con formato invalido"), RequestIdInvalido},
		// {"6 Requerimiento con formato invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), RequestRequerimientoInvalido},
		// {"7 Cbu vacio", fmt.Errorf("Se debe indicar cbu"), RequestCbuVacio},
		// {"8 Requerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), RequestRequerimientoVacio},
		// //Para hacer esta prueba es necesario modificar variables del sistema
		// {"9 ClienteId formato invalido", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// // {"10 ClienteId vacio", fmt.Errorf("Unauthorized"), RequestDebinValido},
		{"12 cbu de otro cliente", nil, RequestCbuOtroCliente},
		// {"13 debin eliminado", nil, RequestDebinEliminado},
		// {"15 debin vencido", nil, RequestDebinVencido},
		// {"16 Token invalido", fmt.Errorf("Unauthorized"), RequestTokenInvalido},
		// {"17 Scope invalido", fmt.Errorf("Forbidden"), RequestDebinScopeInvalido},
	}

	return
}

func TestConsultarDebinPorId(t *testing.T) {

	table := _inicializarConsultarDebinPorId()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetDebinApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.Equal(t, resp.Debin.Estado, linkdtos.Iniciado)
				t.Log(resp.String())
			}

		})
	}

}

type TableEliminarDebines struct {
	Nombre  string
	Erro    error
	Request RequestEliminarDebines
}

type RequestEliminarDebines struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linkdebin.RequestDeleteDebinLink
	Token           string
}

func _inicializarEliminarDebines() (table []TableEliminarDebines) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	requestValido := linkdebin.RequestDeleteDebinLink{
		Cbu: "0110599520000003855199",
		Id:  "746YGOW9MJ58M3P9EXD8J5",
	}

	RequestDebinValido := RequestEliminarDebines{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           requestValido,
		Token:           token.AccessToken,
	}

	requestCbuInvalido := RequestDebinValido
	requestCbuInvalido.Valor.Cbu = "11059952000000385519"

	requestIdInvalido := RequestDebinValido
	requestIdInvalido.Valor.Id = "7586_0"

	requestRequerimientoInvalido := RequestDebinValido
	requestRequerimientoInvalido.RequerimientoId = "b34624f3-e2ff20-4ggb83-8a01-b2b47d87593f6"

	requestCbuVacio := RequestDebinValido
	requestCbuVacio.Valor.Cbu = ""

	requestRequerimientoVacio := RequestDebinValido
	requestRequerimientoVacio.RequerimientoId = ""

	requestTokenVacio := RequestDebinValido
	requestTokenVacio.Token = ""

	requestCbuComprador := RequestDebinValido
	requestCbuComprador.Valor.Cbu = "0140323501420000283729"
	requestCbuComprador.Valor.Id = "YZ6OLMDN3GDX0L52E7RQ5X"

	requestDebinEliminado := RequestDebinValido
	requestDebinEliminado.Valor.Id = "746YGOW9MJ58M3P9EXD8J5"

	requestDebinVencido := RequestDebinValido
	requestDebinVencido.Valor.Id = "V8D0Q619LY8O8QV27JZ5RG"

	// RequestDebinScopeInvalido := RequestEliminarDebines{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           requestValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	table = []TableEliminarDebines{
		// {"2 Eliminar debin creado con recurrencia false", nil, RequestDebinValido},
		// {"4 Cbu con formato invalido", fmt.Errorf("cbu con formato invalido"), requestCbuInvalido},
		// {"5 Id con formato invalido", fmt.Errorf("id con formato invalido"), requestIdInvalido},
		// {"6 RequerimientoId con formato invalido", fmt.Errorf("IdRequerimiento con formato invalido"), requestRequerimientoInvalido},
		// {"7 Cbu vacio", fmt.Errorf("Se debe indicar cbu"), requestCbuVacio},
		// {"8 Requerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), requestRequerimientoVacio},
		// //Para estas pruebas es necesario modificar los parametros del sistema
		// // {"9 Cliente Id invalido", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// // {"10 Cliente Id vacio", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"11 X-IBM Aprobación Pendiente", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"12 Eliminar Debin con cbu comprador", fmt.Errorf("El ID del DEBIN no pertenece a la cuenta informada."), requestCbuComprador},
		// {"13 Debin Eliminado", fmt.Errorf("El DEBIN se encuentra en un estado en el que no se puede modificar."), requestDebinEliminado},
		// {"14 Debin Vencido", fmt.Errorf("El DEBIN se encuentra en un estado en el que no se puede modificar."), requestDebinEliminado},
		// {"15 Token invalido", fmt.Errorf("Unauthorized"), requestTokenVacio},
		// {"16 Scope invalido", fmt.Errorf("Forbidden"), RequestDebinScopeInvalido},
	}

	return
}

func TestEliminarDebines(t *testing.T) {

	table := _inicializarEliminarDebines()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.DeleteDebinApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.Equal(t, resp, true)
				t.Log(resp)
			}

		})
	}

}

type TableGetDebinesPendiente struct {
	Nombre  string
	Erro    error
	Request RequestDebinesPendientes
}

type RequestDebinesPendientes struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           string
	Token           string
}

func _inicializarBuscarDebinesPendientes() (table []TableGetDebinesPendiente) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.DebinRecurrencia)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.TransferenciasBancariasInmediatas)

	cbuValido := "0110599520000003855199"

	RequestDebinValido := RequestDebinesPendientes{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           cbuValido,
		Token:           token.AccessToken,
	}

	// RequestDebinScopeInvalido := RequestDebinesPendientes{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           cbuValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	RequestCbuInvalido := RequestDebinValido
	RequestCbuInvalido.Valor = "01105995200000038455199"

	RequestCbuVacio := RequestDebinValido
	RequestCbuVacio.Valor = ""

	RequestRequerimientoInvalido := RequestDebinValido
	RequestRequerimientoInvalido.RequerimientoId = "2233f545_00oofjjjf_oiirjjmnf"

	RequestRequerimientoVacio := RequestDebinValido
	RequestRequerimientoVacio.RequerimientoId = ""

	RequestTokenVacio := RequestDebinValido
	RequestTokenVacio.Token = ""

	// RequestDebinScopeInvalido := RequestDebinesPendientes{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           cbuValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	table = []TableGetDebinesPendiente{
		{"2 Datos validos", nil, RequestDebinValido},
		// {"4 Cbu con formato invalido", fmt.Errorf("cbu con formato invalido"), RequestCbuInvalido},
		// {"5 Cbu vacio", fmt.Errorf("Se debe indicar cbu"), RequestCbuVacio},
		// {"6 Requerimiento con formato invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), RequestRequerimientoInvalido},
		// {"7 Requerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), RequestRequerimientoVacio},
		// //Para hacer esta prueba es necesario modificar variables de entorno
		// {"8 ClienteId con formato invalido", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// // {"9 ClienteId vacio", fmt.Errorf("Unauthorized"), RequestDebinValido},
		// {"12 token vacio", fmt.Errorf("Unauthorized"), RequestTokenVacio},
		// {"13 escopo invalido", fmt.Errorf("Forbidden"), RequestDebinScopeInvalido},
	}

	return
}

func TestBuscarDebinesPendientes(t *testing.T) {

	table := _inicializarBuscarDebinesPendientes()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor)
			resp, err := v.Request.Repositorio.GetDebinesPendientesApiLink(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				if len(resp.Debines) > 0 {
					assert.Equal(t, resp.Debines[0].Estado, linkdtos.Iniciado)
					t.Log(resp)
				}
			}

		})
	}

}
