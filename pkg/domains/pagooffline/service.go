package pagooffline

import (
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/pagoofflinedtos"
)

type Service interface {
	GenerarCodigoBarra(pagooffline_ov pagoofflinedtos.OffLineRapipagoDtos) (string, error)
}

var pagooffline *service

type service struct {
	algoritmosService commons.AlgoritmoVerificacion
}

func NewService(af commons.AlgoritmoVerificacion) Service {
	pagooffline = &service{
		algoritmosService: af,
	}
	return pagooffline
}

// Resolve devuelve la instancia antes creada
func Resolve() *service {
	return pagooffline
}

func (s *service) GenerarCodigoBarra(pagooffline_ov pagoofflinedtos.OffLineRapipagoDtos) (string, error) {
	var codigoBarraSinDigitoVerificador string
	// construir codigo de barra
	codigoBarraSinDigitoVerificador += AgregarDigitos(pagooffline_ov.GetCodigoEmpresa(), 4)
	codigoBarraSinDigitoVerificador += AgregarDigitos(pagooffline_ov.GetNumeroCliente(), 10)
	codigoBarraSinDigitoVerificador += AgregarDigitos(pagooffline_ov.GetNumeroComprobante(), 12)
	//convierto el importe a string
	importeString := convertirInt64ToString(pagooffline_ov.GetImporte())
	importeString = AgregarDigitos(importeString, 8)
	codigoBarraSinDigitoVerificador += importeString
	/*
		construir fecha en calendario juliano:
		- convierto fechaPrimerVto a string
		- obtengo el año de la fechaPrimerVto
		- construir fecha del primer dia del año y lo paso a string
		- obtengo cantidad de dias entre el primer dia de año y la fechaPrimerVto
		- convierto cantidad de dias a string y formateo a 3 digitos si es necesario
		- armo la fecha en formato juliano

	*/
	fechapriVtoMasUnDia := convertirTimeToString(pagooffline_ov.GetFechaPrimerVto().Add(24*time.Hour), "2006-01-02")
	fechaPrimerVtoString := convertirTimeToString(pagooffline_ov.GetFechaPrimerVto(), "2006-01-02")
	anioActual := pagooffline_ov.GetFechaPrimerVto().Year()
	fechaPrimerDiaAnio := time.Date(anioActual, 1, 1, 0, 0, 0, 0, time.UTC)
	fechaPrimerDiaAnioString := convertirTimeToString(fechaPrimerDiaAnio, "2006-01-02")
	cantidadDias, err := s.algoritmosService.CalcularDiasEntreFechas(fechaPrimerDiaAnioString, fechapriVtoMasUnDia)
	if err != nil {
		return "", err
	}
	cantidadDiasFormateado := AgregarDigitos(convertirIntToString(cantidadDias), 3)
	fechaJuliano := convertirIntToString(anioActual)[2:] + cantidadDiasFormateado
	codigoBarraSinDigitoVerificador += fechaJuliano
	// convierto importe recargo a string y lo formateo a 6 digitos
	importeRecargoString := convertirInt64ToString(pagooffline_ov.GetImporteRecargo())
	importeRecargoFormateado := AgregarDigitos(importeRecargoString, 6)
	codigoBarraSinDigitoVerificador += importeRecargoFormateado
	// obtengo la cantidad de dias entre fechaPrimerVto y fechaSegundoVto
	fechaSegundoVtoString := convertirTimeToString(pagooffline_ov.GetFechaSegundoVto(), "2006-01-02")
	cantidadDiasSegundoVto, err := s.algoritmosService.CalcularDiasEntreFechas(fechaPrimerVtoString, fechaSegundoVtoString)
	if err != nil {
		return "", err
	}
	cantidadDiasSegundoVtoFormateado := AgregarDigitos(convertirIntToString(cantidadDiasSegundoVto), 2)
	codigoBarraSinDigitoVerificador += cantidadDiasSegundoVtoFormateado
	/*
		calcular digito verificador:
		una vez que tengo el codigo de barra sin el digito verificador, se calcula el digito verificador.
	*/
	digitoVerificador, err := s.algoritmosService.CalcularDigitoVerificador(codigoBarraSinDigitoVerificador)
	if err != nil {
		return "", err
	}
	codigoBarra := codigoBarraSinDigitoVerificador + digitoVerificador
	return codigoBarra, nil
}

//// funciones auxiliares///////
func AgregarDigitos(numero string, longitud int) string {
	for len(numero) < longitud {
		numero = "0" + numero
	}
	return numero
}

func convertirInt64ToString(valor int64) string {
	numeroString := strconv.FormatInt(valor, 10)
	return numeroString
}
func convertirIntToString(valor int) string {
	numeroString := strconv.Itoa(valor)
	return numeroString
}

func convertirTimeToString(fecha time.Time, formatoFecha string) string {
	fechaString := fecha.Format(formatoFecha)
	return fechaString
}

/*
	valorTipo := reflect.ValueOf(pagooffline_ov)
	for i := 0; i < valorTipo.Type().NumField(); i++ {
		field := valorTipo.Type().Field(i).Type
		tipo, err := s.algoritmosService.VerificarType(field)
		if err != nil {
			return "", err
		}
		switch tipo {
		case "string":
			codigoBarra += AgregarDigitos(pagooffline_ov.GetCodigoEmpresa(), 4)
		case
		}
}
*/
