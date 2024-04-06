package commons

import (
	"errors"
	"fmt"
	"math"
	"time"

	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func StringIsEmpity(e string) bool {
	return len(strings.TrimSpace(e)) == 0
}

// Recorrer un array de string y checkear si existe al menos uno vacio
func SomeStringIsEmpty(stringsRequired map[string]string) (bool, string) {

	for key, element := range stringsRequired {
		if len(strings.TrimSpace(element)) == 0 {
			return true, key
		}
	}
	return false, ""
}

func IsEmailValid(e string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	email := strings.TrimSpace(e)
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return pattern.MatchString(email)
}

func EsCuilValido(cuil string) error {
	if len(cuil) != 11 {
		return errors.New(ERROR_CUIL)
	}
	// var rv bool
	// var verificador int
	// resultado := 0
	// codes := "5432765432"
	// ultimoDigito := cuil[10:11]
	// verificador, err := strconv.Atoi(ultimoDigito)
	// if err != nil {
	// 	return errors.New(ERROR_CUIL)
	// }
	// for x := 0; x < 10; x++ {
	// 	digitoValidador, _ := strconv.Atoi(codes[x : x+1])
	// 	digito, _ := strconv.Atoi(cuil[x : x+1])
	// 	digitoValidacion := digitoValidador * digito
	// 	resultado += digitoValidacion
	// }
	// //resultado = resultado / 11
	// resto := resultado % 11
	// r2 := 11 - resto
	// rv = (r2 == verificador)
	// if !rv {
	// 	return errors.New(ERROR_CUIL)
	// }
	return nil
}

//// validara tarjeta de credito

// func obtenerLongitud(nro string) (int, error) {
// 	if StringIsEmpity(nro) && len(nro) >= 12 && len(nro) <= 16 {
// 		return 0, fmt.Errorf("longitud no valida %v", nro)
// 	} else {
// 		return int(len(nro)), nil
// 	}
// }
///////////////////////////////////////////////////////////////////////

// func obtenerLongitud(nro string) (int, error) {
// 	if StringIsEmpity(nro) && len(nro) >= 12 && len(nro) <= 16 {
// 		return 0, fmt.Errorf("longitud no valida %v", nro)
// 	} else {
// 		return int(len(nro)), nil
// 	}
// }

// func SumarDigitos(digito string) (int, error) {
// 	if len(digito) == 1 {
// 		suma, _ := strconv.Atoi(digito)
// 		return suma, nil
// 	} else if len(digito) == 2 {
// 		valor1, _ := strconv.Atoi(digito[0:1])
// 		valor2, _ := strconv.Atoi(digito[1:2])
// 		suma := valor1 + valor2
// 		return suma, nil
// 	} else {
// 		return 0, fmt.Errorf("longitud no valida %v", digito)
// 	}
// }

// func DuplicarValor(digito string) string {
// 	x, _ := strconv.Atoi(digito)
// 	z := x * 2
// 	y := strconv.Itoa(z)
// 	return y
// }

// func validarTarjeta(nroCard string) int {
// 	var suma int
// 	var sumaDoble int
// 	par := len(nroCard) % 2
// 	if par == 0 {
// 		for i := len(nroCard); i > 0; i-- {
// 			if i != len(nroCard) {
// 				if (i % 2) == 1 {
// 					digitoDoble := DuplicarValor(nroCard[i-1 : i])
// 					sumadigito, _ := SumarDigitos(digitoDoble)
// 					sumaDoble = sumaDoble + sumadigito
// 					//   fmt.Printf("%v - %v - %v - %v -e %v\n", i ,i%2, nroCard[i-1:i], digitoDoble, sumadigito  )
// 				} else {
// 					x, _ := strconv.Atoi(nroCard[i-1 : i])
// 					suma = suma + x
// 				}
// 			} else {
// 				x, _ := strconv.Atoi(nroCard[i-1 : i])
// 				suma = suma + x
// 			}
// 		}
// 		return suma + sumaDoble
// 		//   fmt.Printf("%v - %v\n", suma, sumaDoble)
// 		//   fmt.Printf("%v \n", suma + sumaDoble)
// 	} else {
// 		for i := len(nroCard); i >= 1; i-- {
// 			if i != len(nroCard) {
// 				if (i % 2) == 0 {
// 					digitoDoble := DuplicarValor(nroCard[i-1 : i])
// 					sumadigito, _ := SumarDigitos(digitoDoble)
// 					sumaDoble = sumaDoble + sumadigito
// 					//fmt.Printf("%v - %v - %v\n", i, i%2, nroCard[i-1:i])
// 				} else {
// 					x, _ := strconv.Atoi(nroCard[i-1 : i])
// 					suma = suma + x
// 				}
// 			} else {
// 				x, _ := strconv.Atoi(nroCard[i-1 : i])
// 				suma = suma + x
// 			}
// 		}
// 		// fmt.Printf("%v - %v\n", suma, sumaDoble)
// 		// fmt.Printf("%v \n", suma+sumaDoble)
// 		return suma + sumaDoble
// 	}
// }

// ChequearTarjeta es la funcio que se debe llamar para validar el formato de una tarjeta
// puede retornanar uno de los siguientes valores
// True: si el formato de tarjeta es valido
// False: si el formato de tarjeta no es valido
// func ChequearTarjeta(valorCheck string) bool {
// 	longitudValida, _ := obtenerLongitud(valorCheck)
// 	if longitudValida != 0 {
// 		valorVerificador := validarTarjeta(valorCheck)
// 		valorModulo := valorVerificador % 10
// 		if valorModulo == 0 {
// 			return true
// 		} else {
// 			return false
// 		}
// 	}
// 	return false
// }

////////////////////////////////////////////////////////////////////////

func Difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				if len(s1) > 0 {
					diff = append(diff, s1)
				}
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func DifferenceInteger(slice1 []uint64, slice2 []uint64) []uint64 {
	var diff []uint64

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				if s1 > 0 {
					diff = append(diff, s1)
				}
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

/*
funciones que se pueden utilizar para quitar caracteres especiales en una cadena
*/
func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

/*
remplazar un caracter por otro
*/
func ReplaceCharacters(str, valorBuscar, valorReemplazar string) string {
	resultadoString := strings.Replace(str, valorBuscar, valorReemplazar, -1)
	return resultadoString
}

func Concat(str1, str2 string) string {
	total := len(str1)
	resultado := str1[0:total-3] + str2
	return resultado
}

func ConcatReferencia(str1 *time.Time, str2 string) string {
	date := str1.Format("2006-01-02")
	resultado := date + str2
	return resultado
}

/*
buscar elementos duplicados en slice de tipo int64
*/
func RemoveDuplicateValues(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func DifferenceInt(slice1 []int64, slice2 []int64) (add []int64, delete []int64) {

	for _, s1 := range slice2 {
		v := existeEnArreglo(slice1, s1)
		if !v {
			delete = append(delete, s1)
		}
	}

	for _, s2 := range slice1 {
		v := existeEnArreglo(slice2, s2)
		if !v {
			add = append(add, s2)
		}
	}
	return add, delete
}

func DifferenceString(slice1 []string, slice2 []string) (add []string, delete []string) {

	for _, s1 := range slice2 {
		v := existeEnArregloString(slice1, s1)
		if !v {
			delete = append(delete, s1)
		}
	}

	for _, s2 := range slice1 {
		v := existeEnArregloString(slice2, s2)
		if !v {
			add = append(add, s2)
		}
	}
	return add, delete
}

func existeEnArreglo(arreglo []int64, busqueda int64) bool {
	for _, numero := range arreglo {
		if numero == busqueda {
			return true
		}
	}
	return false
}

func existeEnArregloString(arreglo []string, busqueda string) bool {
	for _, string := range arreglo {
		if string == busqueda {
			return true
		}
	}
	return false
}

func ConvertFechaString(str1 time.Time) string {
	date := str1.Format("20060102")
	resultado := date
	return resultado
}

// Permite agregar espacios en
func EspaciosBlanco(input string, padLen int, align string) string {
	inputLen := len(input)

	if inputLen >= padLen {
		return input
	}
	var output string
	switch align {
	case "RIGHT":
		output = fmt.Sprintf("% "+strconv.Itoa(-padLen)+"s", input)
	case "LEFT":
		output = fmt.Sprintf("% "+strconv.Itoa(padLen)+"s", input)
	}
	return output
}

func AgregarCeros(cantidad int, valor int) string {
	strngVal := "%0" + strconv.Itoa(cantidad) + "d"
	return fmt.Sprintf(strngVal, valor)
}

func AgregarCerosString(input string, padLen int, align string) string {
	inputLen := len(input)

	if inputLen >= padLen {
		return input
	}
	var output string
	switch align {
	case "RIGHT":
		output = fmt.Sprintf("%0"+strconv.Itoa(-padLen)+"s", input)
	case "LEFT":
		output = fmt.Sprintf("%0"+strconv.Itoa(padLen)+"s", input)
	}
	return output
}

func JoinString(words []string) string {
	var sb strings.Builder
	for _, w := range words {
		fmt.Fprintf(&sb, "%s", w)
	}
	return sb.String()
}

// Buscar un valor de tipo string en un slice de strings
// retorna true en caso de encontrarlo
func ContainStrings(arreglo []string, item string) bool {
	// el string vaciose considera incluido
	if len(item) == 0 {
		return true
	}

	// por ser case sensitive
	itemUpperCase := strings.ToUpper(item)
	for _, value := range arreglo {
		if strings.ToUpper(value) == itemUpperCase {
			return true
		}
	}

	return false
}

// devuelve una fehca string en formato ISO 8601 con la HH:mm:ss finales del dia.
// Uso: comparar limites de fechas
func GetDateLastMoment(fecha time.Time) (fechaISO string) {
	year, month, day := fecha.Date()
	t := time.Date(year, month, day, 23, 59, 59, 999, fecha.Location())
	return t.Format(time.RFC3339)
}

// devuelve una fehca string en formato ISO 8601 con la HH:mm:ss iniciales del dia.
// Uso: comparar limites de fechas
func GetDateFirstMoment(fecha time.Time) (fechaISO string) {
	year, month, day := fecha.Date()
	t := time.Date(year, month, day, 00, 00, 00, 000, fecha.Location())
	return t.Format(time.RFC3339)
}

// retorna time con el ultimo momento de la fecha
func GetDateLastMomentTime(fecha time.Time) (lastMomentDate time.Time) {
	year, month, day := fecha.Date()
	return time.Date(year, month, day, 23, 59, 59, 999, fecha.Location())
}

func ToFixedTool(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
