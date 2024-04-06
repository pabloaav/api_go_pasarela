package commons

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
)

type AlgoritmoVerificacion interface {
	/*
		ChequearTarjeta es la funcio que se debe llamar para validar el formato de una tarjeta
		puede retornanar uno de los siguientes valores
		True: si el formato de tarjeta es valido
		False: si el formato de tarjeta no es valido
	*/
	ChequearTarjeta(valorCheck string) bool
	/*
		permite validar un si un cbu es valido
	*/
	ValidarCBU(cbu string) error

	/*
		calcular cantida de dias entre dos fechas
		formato de defecha inicio y fin:(AAAA-MM-DD) 2006-01-02
	*/
	CalcularDiasEntreFechas(fechaInicio string, fechaFin string) (cantidadDias int, erro error)

	/*
		permite aplicar algoritmo par calcular el digito verificador de un codigo de barra para rapipago
	*/
	CalcularDigitoVerificador(codigo string) (digitoVerificador string, err error)

	/*
		permite validar el typo de una variable o una estructura
	*/
	VerificarType(val reflect.Type) (string, error)
}

type algoritmoVerificacion struct{}

func NewAlgoritmoVerificacion() AlgoritmoVerificacion {
	return &algoritmoVerificacion{}
}

func (algoritmoVerificacion) ChequearTarjeta(valorCheck string) bool {
	longitudValida := obtenerLongitud(valorCheck)
	if longitudValida != 0 {
		//validar solo trajetas naranja
		if strings.HasPrefix(valorCheck, "589562") {
			return true
		}
		valorVerificador := validarTarjeta(valorCheck)
		valorModulo := valorVerificador % 10
		if valorModulo == 0 {
			return true
		} else {
			return false
		}
	}
	return false
}

func (algoritmoVerificacion) ValidarCBU(cbu string) error {

	err := validarLargoCbu(cbu)
	if err != nil {
		return err
	}
	/* TODO validar que el cbu sea valido: verificar esta funcion: CBU TelCo no pasa esta validacion */
	err = validarCodigoBanco(cbu[0:8])
	if err != nil {
		return err
	}
	err = validarCuenta(cbu[8:22])
	if err != nil {
		return err
	}
	return nil

}

func (algoritmoVerificacion) CalcularDiasEntreFechas(fechaInicio string, fechaFin string) (cantidadDias int, erro error) {
	//se define un formato de fecha
	formatoFecha := "2006-01-02" // 15:04:05"

	// formateo de la fecha de inicio y fin
	fechaInicial, err := time.Parse(formatoFecha, fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_FORMATO_FECHA)
		return
	}
	fechaFinal, err := time.Parse(formatoFecha, fechaFin)
	if err != nil {
		erro = errors.New(ERROR_FORMATO_FECHA)
		return
	}
	// se calcula la diferencia de dias entre las fechas
	diferencia := fechaFinal.Sub(fechaInicial)
	// se obtiene la cantidad de dias
	cantidadDias = int(diferencia.Hours() / 24)
	return
}

func (algoritmoVerificacion) CalcularDigitoVerificador(codigo string) (digitoVerificador string, err error) {
	var SumaProducto int
	numeroSecuencia := config.RAPIPAGO_SERIE_NUMERICA
	// fmt.Println("numeroSecuencia: ", numeroSecuencia)
	// fmt.Println("numeroSecuencia: ", len(numeroSecuencia))
	// fmt.Println("codigo: ", codigo)
	// fmt.Println("codigo: ", len(codigo))
	if len(codigo) == 47 && len(numeroSecuencia) == 47 {
		for i := 0; i < len(codigo); i++ {
			digitoCodigoBarra, err := strconv.Atoi(codigo[i : i+1])
			if err != nil {
				return "", errors.New(ERROR_CALCULAR_DIGITO_VERIFICADOR)
			}
			digitoSecuenciaNumenro, err := strconv.Atoi(numeroSecuencia[i : i+1])
			if err != nil {
				return "", errors.New(ERROR_CALCULAR_DIGITO_VERIFICADOR)
			}
			SumaProducto = SumaProducto + (digitoCodigoBarra * digitoSecuenciaNumenro)
		}
		DivisionString := fmt.Sprint(float64(SumaProducto) / 2)
		resultadoArray := strings.Split(DivisionString, ".")
		parteEntera := resultadoArray[0]
		digitoVerificador = parteEntera[len(parteEntera)-1:]
		return digitoVerificador, nil
	}
	return "", errors.New(ERROR_CALCULAR_DIGITO_VERIFICADOR)
}

func (algoritmoVerificacion) VerificarType(val reflect.Type) (string, error) {
	switch val {
	case reflect.TypeOf(0):
		return "int", nil
	case reflect.TypeOf(int32(0)):
		return "int32", nil
	case reflect.TypeOf(int64(0)):
		return "int64", nil
	case reflect.TypeOf(float32(0.0)):
		return "float32", nil
	case reflect.TypeOf(0.0):
		return "float64", nil
	case reflect.TypeOf(""):
		return "string", nil
	case reflect.TypeOf(time.Time{}):
		return "time", nil
	default:
		return "", errors.New("no es de ningun tipo")
	}
}

////////funciones tarjeta////////////////
func obtenerLongitud(nro string) int {
	if StringIsEmpity(nro) && len(nro) >= 12 && len(nro) <= 16 {
		//return 0, fmt.Errorf("longitud no valida %v", nro)
		return 0
	} else {
		return int(len(nro))
	}
}

func validarTarjeta(nroCard string) int {
	var suma int
	var sumaDoble int
	par := len(nroCard) % 2
	if par == 0 {
		for i := len(nroCard); i > 0; i-- {
			if i != len(nroCard) {
				if (i % 2) == 1 {
					digitoDoble := DuplicarValor(nroCard[i-1 : i])
					sumadigito := SumarDigitos(digitoDoble)
					sumaDoble = sumaDoble + sumadigito
				} else {
					x, _ := strconv.Atoi(nroCard[i-1 : i])
					suma = suma + x
				}
			} else {
				x, _ := strconv.Atoi(nroCard[i-1 : i])
				suma = suma + x
			}
		}
		return suma + sumaDoble
	} else {
		for i := len(nroCard); i >= 1; i-- {
			if i != len(nroCard) {
				if (i % 2) == 0 {
					digitoDoble := DuplicarValor(nroCard[i-1 : i])
					sumadigito := SumarDigitos(digitoDoble)
					sumaDoble = sumaDoble + sumadigito
				} else {
					x, _ := strconv.Atoi(nroCard[i-1 : i])
					suma = suma + x
				}
			} else {
				x, _ := strconv.Atoi(nroCard[i-1 : i])
				suma = suma + x
			}
		}
		return suma + sumaDoble
	}
}

func DuplicarValor(digito string) string {
	x, _ := strconv.Atoi(digito)
	z := x * 2
	y := strconv.Itoa(z)
	return y
}

func SumarDigitos(digito string) int {
	if len(digito) == 1 {
		suma, _ := strconv.Atoi(digito)
		return suma //, nil
	} else if len(digito) == 2 {
		valor1, _ := strconv.Atoi(digito[0:1])
		valor2, _ := strconv.Atoi(digito[1:2])
		suma := valor1 + valor2
		return suma //, nil
	} else {
		return 0 //, fmt.Errorf("longitud no valida %v", digito)
	}
}

//////funciones cbu
func validarLargoCbu(cbu string) error {
	if StringIsEmpity(cbu) {
		return fmt.Errorf("cbu está en blanco")
	}
	if len(cbu) != 22 {
		return fmt.Errorf("longitud de cbu no es válido: %d", len(cbu))
	}
	return nil
}

func validarCodigoBanco(codigo string) error {
	if len(codigo) != 8 {
		return fmt.Errorf("el código de banco es incorrecto")
	}
	/* TODO: descometar esta funcion para validar el contendio del cbu */
	// banco := codigo[0:3]
	// digitoVerificador := codigo[3:4]
	// sucursal := codigo[4:7]
	// digitoVerificador2 := codigo[7:8]
	// var suma int
	// var x int
	// x, _ = strconv.Atoi(banco[0:1])
	// suma = x * 7
	// x, _ = strconv.Atoi(banco[1:2])
	// suma = suma + x
	// x, _ = strconv.Atoi(banco[2:3])
	// suma = suma + (x * 3)
	// x, _ = strconv.Atoi(digitoVerificador)
	// suma = suma + (x * 9)
	// x, _ = strconv.Atoi(sucursal[0:1])
	// suma = suma + (x * 7)
	// x, _ = strconv.Atoi(sucursal[1:2])
	// suma = suma + x
	// x, _ = strconv.Atoi(sucursal[2:3])
	// suma = suma + (x * 3)
	// diferencia := 10 - (suma % 10)
	// digito, _ := strconv.Atoi(digitoVerificador2)
	// if diferencia != digito {
	// 	return fmt.Errorf("codigo de banco inválido")
	// }
	return nil
}

func validarCuenta(cuenta string) error {
	if len(cuenta) != 14 {
		return fmt.Errorf("logitud de cuenta inválido: %d", len(cuenta))
	}
	/* TODO : descomentar esta funcion para validar el contendio del cbu */
	// digitoVerificador, _ := strconv.Atoi(cuenta[13:14])
	// var suma int
	// var x int
	// x, _ = strconv.Atoi(cuenta[0:1])
	// suma = x * 3
	// x, _ = strconv.Atoi(cuenta[1:2])
	// suma = suma + (x * 9)
	// x, _ = strconv.Atoi(cuenta[2:3])
	// suma = suma + (x * 7)
	// x, _ = strconv.Atoi(cuenta[3:4])
	// suma = suma + x
	// x, _ = strconv.Atoi(cuenta[4:5])
	// suma = suma + (x * 3)
	// x, _ = strconv.Atoi(cuenta[5:6])
	// suma = suma + (x * 9)
	// x, _ = strconv.Atoi(cuenta[6:7])
	// suma = suma + (x * 7)
	// x, _ = strconv.Atoi(cuenta[7:8])
	// suma = suma + (x * 1)
	// x, _ = strconv.Atoi(cuenta[8:9])
	// suma = suma + (x * 3)
	// x, _ = strconv.Atoi(cuenta[9:10])
	// suma = suma + (x * 9)
	// x, _ = strconv.Atoi(cuenta[10:11])
	// suma = suma + (x * 7)
	// x, _ = strconv.Atoi(cuenta[11:12])
	// suma = suma + (x * 1)
	// x, _ = strconv.Atoi(cuenta[12:13])
	// suma = suma + (x * 3)
	// diferencia := 10 - (suma % 10)
	// if diferencia != digitoVerificador {
	// 	return fmt.Errorf("error en cuenta bancaria")
	// }
	return nil
}
