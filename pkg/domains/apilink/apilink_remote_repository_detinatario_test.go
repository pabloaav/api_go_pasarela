package apilink_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/linkconsultadestinatario"
	"github.com/stretchr/testify/assert"
)

type TableGetConsultaDestinarario struct {
	Nombre  string
	Erro    error
	Request RequestGestConsultaDestinatario
}

type RequestGestConsultaDestinatario struct {
	Repositorio     apilink.RemoteRepository
	RequerimientoId string
	Valor           linkconsultadestinatario.RequestConsultaDestinatarioLink
	Token           string
}

func _inicializarGetConsultaDestinatario() (table []TableGetConsultaDestinarario) {

	repositorio, uuidValido, token := InicializarDebines(linkdtos.ConsultaDestinatario)

	// repositorioScopeInvalido, uuidValidoScopeInvalido, tokenScopeInvalido := InicializarDebines(linkdtos.DebinRecurrencia)

	valorValido := linkconsultadestinatario.RequestConsultaDestinatarioLink{
		Cbu:   "0110599520000003855199",
		Alias: "",
	}

	requestValido := RequestGestConsultaDestinatario{
		Repositorio:     repositorio,
		RequerimientoId: uuidValido,
		Valor:           valorValido,
		Token:           token.AccessToken,
	}

	// requestScopeInvalido := RequestGestConsultaDestinatario{
	// 	Repositorio:     repositorioScopeInvalido,
	// 	RequerimientoId: uuidValidoScopeInvalido,
	// 	Valor:           valorValido,
	// 	Token:           tokenScopeInvalido.AccessToken,
	// }

	requestMultiplesTitulares := requestValido
	requestMultiplesTitulares.Valor.Cbu = "0110017420001700884103"

	requestAliasProprioTitular := requestValido
	requestAliasProprioTitular.Valor.Cbu = ""
	requestAliasProprioTitular.Valor.Alias = "SANTERESA"

	requestCbuInvalido := requestValido
	requestCbuInvalido.Valor.Cbu = "3890001130005274052561"

	requestCbuInactivo := requestValido
	requestCbuInactivo.Valor.Cbu = "0200318211000011138080"

	requestAliasInvalido := requestValido
	requestAliasInvalido.Valor.Cbu = ""
	requestAliasInvalido.Valor.Alias = "SANTE"

	requestRequerimientoInvalido := requestValido
	requestRequerimientoInvalido.RequerimientoId = "7vvbbb2-223355-5kgiolff-hhhgg"

	requestRequerimientoVacio := requestValido
	requestRequerimientoVacio.RequerimientoId = " "

	requestAliasYCbu := requestValido
	requestAliasYCbu.Valor.Alias = "SANTERESA"

	requestAliasYCbuVacios := requestValido
	requestAliasYCbuVacios.Valor.Cbu = ""

	requestCbuInexistente := requestValido
	requestCbuInexistente.Valor.Cbu = "0110616530061602678801"

	requestAliasInactivo := requestValido
	requestAliasInactivo.Valor.Cbu = ""
	requestAliasInactivo.Valor.Alias = "ApiLink2023"

	requestAliasEliminado := requestValido
	requestAliasEliminado.Valor.Cbu = ""
	requestAliasEliminado.Valor.Alias = "Marbella.ES"

	requestCvu := requestValido
	requestCvu.Valor.Cbu = "0000262402008276351202"

	requestAliasCvu := requestValido
	requestAliasCvu.Valor.Cbu = ""
	requestAliasCvu.Valor.Alias = "AliasNuevoXC123"

	table = []TableGetConsultaDestinarario{
		// {"FB-001 Consulta cbu propio titular", nil, requestValido},
		// {"FB-002 Consulta cbu multiples titulares", nil, requestMultiplesTitulares},
		// {"FB-003 Consulta alias propio titular", nil, requestAliasProprioTitular},
		// {"FI-004 Consulta cbu invalido", fmt.Errorf("cbu con formato invalido"), requestCbuInvalido},
		// {"FI-005 Consulta alias invalido", fmt.Errorf("alias con formato invalido"), requestAliasInvalido},
		// {"FI-006 X idRequerimiento invalido", fmt.Errorf("X-IdRequerimiento con formato invalido"), requestRequerimientoInvalido},
		//Para ejecutar estas pruebas hay que modificar las variables de entorno
		// // {"FI-007 X-IBM-Client_id invalido", fmt.Errorf("Unauthorized"), requestValido},
		// // {"FI-008 X-IBM-Client_id vacio", fmt.Errorf("Unauthorized"), requestValido},
		// {"PR-009 X idRequerimiento vacio", fmt.Errorf("Se debe indicar X-IdRequerimiento"), requestRequerimientoVacio},
		// {"FI-010 consultar por alias y cbu", fmt.Errorf("Uno (y solo uno) de los siguientes campos debe estar completo"), requestAliasYCbu},
		// {"FI-011 consultar por alias y cbu vacios", fmt.Errorf("Uno (y solo uno) de los siguientes campos debe estar completo"), requestAliasYCbuVacios},
		// {"FI-012 consultar por cbu inexistente", fmt.Errorf("La búsqueda del CBU solicitado no arrojó resultados"), requestCbuInexistente},
		// {"FE-013 alias inactivo", fmt.Errorf("El alias de destino con el que se quiere operar es invalido o inexistente."), requestAliasInactivo},
		// {"FE-014 alias ELIMINADO", fmt.Errorf("El alias de destino con el que se quiere operar es invalido o inexistente."), requestAliasEliminado},
		{"FE-015 consultar por alias inactiva", nil, requestAliasInactivo},
		// {"FE-016 consultar por cbu inactiva", nil, requestCbuInactivo},
		// {"FE-017 consultar por cvu", fmt.Errorf("La cuenta ingresada corresponde a una CVU"), requestCvu},
		// {"FE-017 consultar por AliasCVU", fmt.Errorf("El alias ingresado corresponde a una\u00a0CVU"), requestAliasCvu},
		// {"FE-017 Authorizacion de otro Scope", fmt.Errorf("Forbidden"), requestScopeInvalido},

	}

	return
}

func TestGetConsultaDestinatario(t *testing.T) {

	table := _inicializarGetConsultaDestinatario()

	for _, v := range table {

		t.Run(v.Nombre, func(t *testing.T) {

			t.Log(v.Request.Valor.String())
			resp, err := v.Request.Repositorio.GetConsultaDestinatario(v.Request.RequerimientoId, v.Request.Valor, v.Request.Token)

			if err != nil {
				assert.Contains(t, err.Error(), v.Erro.Error())
			}

			if err == nil {
				assert.NotEmpty(t, resp.EntidadBancaria.Nombre)
				t.Log(resp.String())
			}

		})
	}

}
