package apilink_test

import (
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkcuentas"
	"github.com/stretchr/testify/assert"
)

type TablePostCuentaLink struct {
	Nombre  string
	Erro    error
	Request RequestPostCuentaLink
}

type RequestPostCuentaLink struct {
	Repositorio apilink.RemoteRepository
	Valor       linkcuentas.LinkCuentasRequest
}

func _inicializarPostCuentaLink() (table []TablePostCuentaLink) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.AdhesionCuenta)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linkcuentas.LinkPostCuenta{
		Cuit: "30546676427",
		Cbu:  "0110599520000003855199",
	}

	datosValidos := linkcuentas.LinkCuentasRequest{
		Token:           token,
		RequerimientoId: uuidValido,
		Request:         valorValido,
	}

	// datosScopeInvalido := linkcuentas.LinkCuentasRequest{
	// 	Token:           tokenScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Request:         valorValido,
	// }

	// requestScopeInvalido := RequestPostCuentaLink{
	// 	Repositorio: repositorioScopeInvalido,
	// 	Valor:       datosScopeInvalido,
	// }

	requestValido := RequestPostCuentaLink{
		Repositorio: repositorio,
		Valor:       datosValidos,
	}

	RequestRequerimientoInvalido := requestValido
	RequestRequerimientoInvalido.Valor.RequerimientoId = "__*"

	ValorInvalidoCbu := valorValido
	ValorInvalidoCbu.Cbu = "01105995200000038551"
	RequestCbuInvalido := requestValido
	RequestCbuInvalido.Valor.Request = ValorInvalidoCbu

	ValorInvalidoCuit := valorValido
	ValorInvalidoCuit.Cuit = "3054667642"
	RequestCuitInvalido := requestValido
	RequestCuitInvalido.Valor.Request = ValorInvalidoCuit

	RequestRequerimientoVacio := requestValido
	RequestRequerimientoVacio.Valor.RequerimientoId = " "

	ValorVacioCbu := valorValido
	ValorVacioCbu.Cbu = " "
	RequestCbuVacio := requestValido
	RequestCbuVacio.Valor.Request = ValorVacioCbu

	ValorVacioCuit := valorValido
	ValorVacioCuit.Cuit = " "
	RequestCuitVacio := requestValido
	RequestCuitVacio.Valor.Request = ValorVacioCuit

	RequestTokenInvalido := requestValido
	RequestTokenInvalido.Valor.Token.AccessToken = "eyJraWQiOiJSZWRMaW5rIiwiYWxnIjoiSFM1MTIifQ.eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiQURIRV9WRU5EIiwiYXVkIjoiaC5hcGkucmVkbGluay5jb20uYXIvcmVkbGluay9ob21vbG9nYWNpb24vIiwiZXhwIjoxNjMyOTI2MTcwLCJpYXQiOjE2MzI5MjI1NzB9.QnEPeIF8d2HS"

	table = []TablePostCuentaLink{
		// {"2 Debe crear una cuenta correctamente", nil, requestValido},
		// {"3 Error requerimiento invalido", fmt.Errorf("request.headers.X-IdRequerimiento"), RequestRequerimientoInvalido},
		// {"4 Error cbu invalido", fmt.Errorf("request.body.cbu"), RequestCbuInvalido},
		// {"5 Error cuit invalido", fmt.Errorf("request.body.cuit"), RequestCuitInvalido},
		// {"6 Error requerimiento vacio", fmt.Errorf("request.headers.X-IdRequerimiento"), RequestRequerimientoInvalido},
		// {"7 Error cbu vacio", fmt.Errorf("request.body.cbu"), RequestCbuVacio},
		// {"8 Error cuit vacio", fmt.Errorf("request.body.cuit"), RequestCuitVacio},
		{"10 X-IBM-Cliente_id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// {"12 Error token incompleto", fmt.Errorf("Unauthorized"), RequestTokenInvalido},
		// {"13 Error scope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
	}

	return
}

func TestPostCuentaLink(t *testing.T) {

	table := _inicializarPostCuentaLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {
			request := v.Request.Valor.Request.(linkcuentas.LinkPostCuenta)
			t.Log(request.String())
			err := v.Request.Repositorio.CreateCuentaApiLink(v.Request.Valor)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

		})
	}

}

type TableGettCuentaLink struct {
	Nombre  string
	Erro    error
	Request RequestGetCuentaLink
}

type RequestGetCuentaLink struct {
	Repositorio apilink.RemoteRepository
	Valor       linkcuentas.LinkGetCuentasRequest
}

func _inicializarGetCuentaLink() (table []TableGettCuentaLink) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.AdhesionCuenta)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	datosValidos := linkcuentas.LinkGetCuentasRequest{
		Token:           token,
		RequerimientoId: uuidValido,
	}

	// datosScopeInvalido := linkcuentas.LinkGetCuentasRequest{
	// 	Token:           tokenScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// }

	// requestScopeInvalido := RequestGetCuentaLink{
	// 	Repositorio: repositorioScopeInvalido,
	// 	Valor:       datosScopeInvalido,
	// }

	requestValido := RequestGetCuentaLink{
		Repositorio: repositorio,
		Valor:       datosValidos,
	}

	requestRequerimientoInvalido := requestValido
	requestRequerimientoInvalido.Valor.RequerimientoId = "__;;;;"

	requestRequerimientoVacio := requestValido
	requestRequerimientoVacio.Valor.RequerimientoId = ""

	requestTokenVacio := requestValido
	requestTokenVacio.Valor.Token.AccessToken = " "

	table = []TableGettCuentaLink{

		// {"2 Consulta valida", nil, requestValido},
		// {"3 Requerimiento-Id invalido", fmt.Errorf("<request.headers.X-IdRequerimiento>"), requestRequerimientoInvalido},
		// {"3 Requerimiento-Id vacio", fmt.Errorf("<request.headers.X-IdRequerimiento>"), requestRequerimientoInvalido},
		// {"8 Token Formato Invalido vacio", fmt.Errorf("Unauthorized"), requestTokenVacio},
		// {"9 Token con scope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
		{"10 X-IBM-Cliente_id invalido", fmt.Errorf("Unauthorized"), requestValido},
	}

	return
}

func TestGetCuentaLink(t *testing.T) {

	table := _inicializarGetCuentaLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			resp, err := v.Request.Repositorio.GetCuentasApiLink(v.Request.Valor)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				t.Log(resp)
				assert.NotNil(t, resp)
			}

		})
	}

}

type TableDeleteCuentaLink struct {
	Nombre  string
	Erro    error
	Request RequestDeleteCuentaLink
}

type RequestDeleteCuentaLink struct {
	Repositorio apilink.RemoteRepository
	Valor       linkcuentas.LinkCuentasRequest
}

func _inicializarDeleteCuentaLink() (table []TableDeleteCuentaLink) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.AdhesionCuenta)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linkcuentas.LinkDeleteCuenta{
		Cbu: "0110599520000003855199",
	}

	datosValidos := linkcuentas.LinkCuentasRequest{
		Token:           token,
		RequerimientoId: uuidValido,
		Request:         valorValido,
	}

	// datosScopeInvalido := linkcuentas.LinkCuentasRequest{
	// 	Token:           tokenScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Request:         valorValido,
	// }

	// requestScopeInvalido := RequestDeleteCuentaLink{
	// 	Repositorio: repositorioScopeInvalido,
	// 	Valor:       datosScopeInvalido,
	// }

	requestValido := RequestDeleteCuentaLink{
		Repositorio: repositorio,
		Valor:       datosValidos,
	}

	RequestRequerimientoInvalido := requestValido
	RequestRequerimientoInvalido.Valor.RequerimientoId = ":;;"

	ValorInvalidoCbu := valorValido
	ValorInvalidoCbu.Cbu = "01105995200000038555551"
	RequestCbuInvalido := requestValido
	RequestCbuInvalido.Valor.Request = ValorInvalidoCbu

	RequestRequerimientoVacio := requestValido
	RequestRequerimientoVacio.Valor.RequerimientoId = " "

	ValorVacioCbu := valorValido
	ValorVacioCbu.Cbu = " "
	RequestCbuVacio := requestValido
	RequestCbuVacio.Valor.Request = ValorVacioCbu

	RequestTokenInvalido := requestValido
	RequestTokenInvalido.Valor.Token.AccessToken = "eyJpc3MiOiJBUElMaW5rIiwic3ViIjoiQURIRV9WRU5EIiwiYXVkIjoiaC5hcGkucmVkbGluay5jb20uYXIvcmVkbGluay9ob21vbG9nYWNpb24vIiwiZXhwIjoxNjMyOTI2MTcwLCJpYXQiOjE2MzI5MjI1NzB9.QnEPeIF8d2HS"

	table = []TableDeleteCuentaLink{
		{"2 Debe eliminar una cuenta correctamente", nil, requestValido},
		// {"3 Error requerimiento invalido", fmt.Errorf("request.headers.X-IdRequerimiento"), RequestRequerimientoInvalido},
		// {"4 Error cbu invalido", fmt.Errorf("request.body.cbu"), RequestCbuInvalido},
		// {"6 Error requerimiento vacio", fmt.Errorf("request.headers.X-IdRequerimiento"), RequestRequerimientoInvalido},
		// {"8 Error cbu vacio", fmt.Errorf("request.body.cbu"), RequestCbuVacio},
		{"10 X-IBM-Cliente_id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// {"13 Error token incompleto", fmt.Errorf("Unauthorized"), RequestTokenInvalido},
		// {"14 Error scope invalido", fmt.Errorf("Forbidden"), requestScopeInvalido},
	}

	return
}

func TestDeleteCuentaLink(t *testing.T) {

	table := _inicializarDeleteCuentaLink()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {
			request := v.Request.Valor.Request.(linkcuentas.LinkDeleteCuenta)
			t.Log(request.String())
			err := v.Request.Repositorio.DeleteCuentaApiLink(v.Request.Valor)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

		})
	}

}
