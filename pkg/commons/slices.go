package commons

import "reflect"

// Index devuelve un entero positivo con el indice del objeto t dentro del slice vs
func Index(vs []interface{}, t interface{}) int {
	for i, v := range vs {
		if reflect.DeepEqual(v, t) {
			return i
		}
	}
	return -1
}

// Include busca dentro del slice [vs] un objeto [t] y devuelve true si lo encuentra
// ej: siendo lista []interface{} y p2 interface{}
// 	if Include(lista, p2) {
//	   fmt.Println("p2 se encuentra en la lista")
//	}
func Include(vs []interface{}, t interface{}) bool {
	return Index(vs, t) >= 0
}

// Any ejecuta una funcion a cada objeto dentro de un slice y devuelve true si alguno fue true
// ej: siendo lista []interface{} y Persona una struct{}
// fmt.Println(Any(lista, func(v interface{}) bool {
//		p := v.(Persona)
//		return strings.HasPrefix(p.Nombre, "f")
//	}))
func Any(vs []interface{}, f func(interface{}) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All ejecuta una funcion a cada objeto dentro de un slice y devuelve true si todos fueron true
// ej: siendo lista []interface{} y Persona una struct{}
// fmt.Println(All(lista, func(v interface{}) bool {
//		p := v.(Persona)
//		return strings.HasPrefix(p.Nombre, "f")
//	}))
func All(vs []interface{}, f func(interface{}) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter ayuda a filtrar un slice mediante cualquier condicion que agreguemos
// ej: busco en lista las personas con id == 1
// fmt.Println(Filter(lista, func(v interface{}) bool {
// 		p := v.(Persona)
// 		return p.ID == 1
// }))
func Filter(vs []interface{}, f func(interface{}) bool) []interface{} {
	vsf := make([]interface{}, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map transforma un slice corriendo una funcion cualquiera a cada objeto adentro
// ej: a cada persona en la lista le modifico el nombre
// fmt.Println(Map(lista, func(v interface{}) interface{} {
// 		p := v.(Persona)
// 		p.Nombre = ""
// 		return p
// }))
func Map(vs []interface{}, f func(interface{}) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
