package util

import (
	"fmt"
	"regexp"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

func ValidateRequestUpdated(request *[]administraciondtos.SubcuentaRequest, subcuentasByCuentasId *[]*entities.Subcuenta) (ok bool, err error) {
	var (
		controlId                        = 0
		porcentajeRequestsModificado     = 0.0
		porcentajeRequestsNoModificado   = 0.0
		porcentajeSubcuentasNoModificado = 0.0 // Porcentaje de DB
		tipoPrincipalRequest             = 0
	)

	if len(*request) != len(*subcuentasByCuentasId) {
		err = fmt.Errorf("La cantidad de subcuentas no coinciden con las existentes.")
		return false, err
	}

	AplicaCostoServicio := false
	for _, v := range *request {
		//Si el porcentaje actual es igual a cero
		if v.Porcentaje == 0.0 {
			err = fmt.Errorf("Los porcentajes deben ser distinto de 0.")
			return false, err
		}
		if len(v.Nombre) > 35 {
			err = fmt.Errorf("El nombre no puede contener mas de 35 dígitos.")
			return false, err
		}

		//Si el id actual es igual a el anterior
		if controlId == int(v.Id) {
			err = fmt.Errorf("Ids de subcuentas duplicados o id = 0.")
			return false, err
		}
		if !validateEmail(v.Email) {
			err = fmt.Errorf("Email inválido.")
			return false, err
		}
		//actualizo id
		controlId = int(v.Id)
		for _, p := range *request {
			if v.Cbu == p.Cbu && len(p.Cbu) > 0 && p.Id != v.Id {
				err = fmt.Errorf("El CBU no puede repetirse. Cbu: %s", v.Cbu)
				return false, err
			}
			if p.CuentaID != v.CuentaID && v.CuentaID != 0 && v.Id != p.Id {
				err = fmt.Errorf("Solo se permite crear subcuentas para una cuenta, no mas. (verificar cuentas_id)")
				return false, err
			}
		}

		for _, va := range *subcuentasByCuentasId {
			//Controlo que los que los cbus que me envia el front-end no se repita en otras subcuentas que no sea la ella misma
			if va.Cbu == v.Cbu && va.ID != v.Id {
				err = fmt.Errorf("Ya existe una cuenta con ese cbu. Cbu: %s", v.Cbu)
				return false, err
			}
			if va.ID == v.Id {
				if v.Modificado != 1 {
					porcentajeSubcuentasNoModificado += va.Porcentaje
				}
			}

		}
		//Cuento el tipo de subcuenta y los porcentajes, en base a los que se modificaron y los que no, para realizar un control posterior
		if v.Modificado == 1 {
			if v.Tipo == "principal" {
				tipoPrincipalRequest++
			}
			porcentajeRequestsModificado += v.Porcentaje
		} else {
			if v.Tipo == "principal" {
				tipoPrincipalRequest++
			}
			porcentajeRequestsNoModificado += v.Porcentaje

		}

		if v.AplicaCostoServicio == v.AplicaPorcentaje {
			err = fmt.Errorf("Se puede aplicar solo a una opcion, 'costo servicio' o 'porcentaje'.(Por lo menos una de las dos opciones.)")
			return false, err
		} else {
			if v.AplicaCostoServicio && !v.AplicaPorcentaje {
				AplicaCostoServicio = true
			}
		}

	}

	controlId = 0
	if tipoPrincipalRequest != 1 {
		if tipoPrincipalRequest == 0 {
			err = fmt.Errorf("No se mandó tipo de subcuenta 'principal'. ")
		} else {
			err = fmt.Errorf("Solo se permite una subcuenta de tipo 'principal'. ")
		}

		return false, err
	}

	if porcentajeSubcuentasNoModificado+porcentajeRequestsModificado != 100.0 && !AplicaCostoServicio {
		err = fmt.Errorf("La suma de porcentajes de las cuentas deben alcanzar un 100 porciento. ")
		return false, err
	}

	if AplicaCostoServicio && len(*request) > 2 {
		err = fmt.Errorf("Aplicando a 'Costo de servicio' no es posible tener mas de 2 subcuentas. ")
		return false, err
	}

	return
}
func ValidateRequestCreated(request *[]administraciondtos.SubcuentaRequest, subcuentasByCuentasId *[]*entities.Subcuenta) (ok bool, err error) {
	// Variables auxiliares
	var (
		porcentajeRequestsNoModificado = 0.0
		contAplicaPorcentaje           = 0
		contAplicaCostoServicio        = 0
		tipoPrincipalRequest           = 0
		cuentaId                       = 0
		AplicaCostoServicio            = false
		AplicaPorcentaje               = false
	)

	if len(*subcuentasByCuentasId) > 0 {
		err = fmt.Errorf("Enviar todas las subcuentas mas las cuentas que se quieren agregar.")
		return false, err
	}

	tipoPrincipalRequest = 0

	for i, v := range *request {
		if len(v.Nombre) > 35 {
			err = fmt.Errorf("Los nombres de cuentas no puede contener mas de 35 dígitos.")
			return false, err
		}
		if len(v.Nombre) <= 3 {
			err = fmt.Errorf("Los nombres de cuentas deben contener mas de 3 dígitos.")
			return false, err
		}
		if v.Porcentaje == 0.0 {
			err = fmt.Errorf("Los porcentajes deben ser distinto de 0.")
			return false, err
		}
		if i == 0 {
			cuentaId = int(v.CuentaID)
		}

		if int64(cuentaId) != v.CuentaID {
			err = fmt.Errorf("Campo 'Cuentas_id' con distintos valores o, no se envió alguno. ")
			return false, err
		}
		if v.CuentaID == 0 {
			err = fmt.Errorf("Campo 'Cuentas_id' invalido. ", v.CuentaID)
			return false, err
		}

		if !validateEmail(v.Email) {
			err = fmt.Errorf("Email inválido.")
			return false, err
		}
		//Valido Cbu de las subcuentas que me manda
		for j, p := range *request {
			if v.Cbu == p.Cbu && j != i {
				err = fmt.Errorf("El CBU no puede repetirse. Cbu: %s", v.Cbu)
				return false, err
			}

		}
		//Valido Cbu de las subcuentas que me manda en comparacion con los de la DB
		for _, va := range *subcuentasByCuentasId {
			//Controlo que los que los cbus que me envia el front-end no se repita en otras subcuentas que no sea la ella misma
			if va.Cbu == v.Cbu && va.ID != v.Id {
				err = fmt.Errorf("Ya existe una cuenta con ese cbu. Cbu: %s", v.Cbu)
				return false, err
			}

		}
		if v.AplicaCostoServicio == v.AplicaPorcentaje && v.AplicaCostoServicio {
			err = fmt.Errorf("Se puede aplicar solo a una opcion, 'costo servicio' o 'porcentaje'.(Por lo menos una de las dos opciones.)")
			return false, err
		} else if v.AplicaCostoServicio == v.AplicaPorcentaje && !v.AplicaCostoServicio {
			err = fmt.Errorf("Se debe aplicar por lo menos a una opcion, 'costo servicio' o 'porcentaje'.(Por lo menos una de las dos opciones.)")
			return false, err
		} else if v.AplicaCostoServicio {
			AplicaCostoServicio = true
		} else if v.AplicaPorcentaje {
			AplicaPorcentaje = true
		}

		//Cuento las las subcuentas de tipo principal para luego controlar que solo haya una
		if v.Tipo == "principal" {
			tipoPrincipalRequest++
		}
		if v.AplicaCostoServicio {
			contAplicaCostoServicio++
		}
		if v.AplicaPorcentaje {
			contAplicaPorcentaje++
		}
		//Almaceno el total de porcentajes para luego hacer un control
		porcentajeRequestsNoModificado += v.Porcentaje
	}

	if AplicaCostoServicio {
		if len(*request) != contAplicaCostoServicio {
			err = fmt.Errorf("Existe/n subcuentas que no tiene/n seteado 'Aplica_costo_servicio'.")
			return false, err
		}
	}
	if AplicaPorcentaje {
		if len(*request) != contAplicaPorcentaje {
			err = fmt.Errorf("Existe/n subcuentas que no tiene/n seteado 'Aplica_porcentaje'.")
			return false, err
		}
	}

	if tipoPrincipalRequest != 1 {
		err = fmt.Errorf("Solo puede existir una cuenta de tipo 'principal'. Ni más, ni menos. ")
		return false, err
	}
	if porcentajeRequestsNoModificado != 100.0 && !AplicaCostoServicio {
		err = fmt.Errorf("La suma de porcentajes de las cuentas deben alcanzar un 100 porciento. ")
		return false, err
	}
	if AplicaCostoServicio && len(*request) > 2 {
		err = fmt.Errorf("Aplicando a 'Costo de servicio' no es posible tener mas de 2 subcuentas. ")
		return false, err
	}
	return
}
func ValidateRequestCreatedUpdated(request *[]administraciondtos.SubcuentaRequest, subcuentasByCuentasId *[]*entities.Subcuenta) (ok bool, err error) {
	var (
		porcentajeRequestsModificado   = 0.0
		porcentajeRequestsNoModificado = 0.0
		tipoPrincipalRequest           = 0
		AplicaCostoServicio            = false
	)

	if len(*request) < len(*subcuentasByCuentasId) {
		err = fmt.Errorf("Enviar todas las subcuentas mas las cuentas que se quieren agregar.")
		return false, err
	}

	cuentasExistentesRecibidas := 0
	array := *request

	for i, v := range array {
		if len(v.Nombre) > 35 {
			err = fmt.Errorf("El nombre no puede contener mas de 35 dígitos.")
			return false, err
		}
		if len(v.Nombre) <= 3 && v.Modificado == 0 {
			err = fmt.Errorf("El nombre debe contener mas de 3 dígitos.")
			return false, err
		}
		if v.Porcentaje == 0.0 {
			err = fmt.Errorf("Los porcentajes deben ser distinto de 0.")
			return false, err
		}
		if !validateEmail(v.Email) {
			err = fmt.Errorf("Email inválido.")
			return false, err
		}

		if v.Id != 0 {
			cuentasExistentesRecibidas++
		}

		if i > 0 {
			if array[i-1].Id != 0 && v.Id != 0 && v.Id == array[i-1].Id {
				err = fmt.Errorf("Id de subcuentas repetidos. ")
				return false, err
			}
		}
		for j, p := range *request {
			if v.Cbu == p.Cbu && j != i {
				err = fmt.Errorf("El CBU no puede repetirse. Cbu: %s", v.Cbu)
				return false, err
			}
		}
		for _, va := range *subcuentasByCuentasId {
			//Controlo que los que los cbus que me envia el front-end no se repita en otras subcuentas que no sea la ella misma
			if va.Cbu == v.Cbu && va.ID != v.Id {
				err = fmt.Errorf("Ya existe una cuenta con ese cbu. Cbu: %s.", v.Cbu)
				return false, err
			}
			if int64(va.CuentasID) != v.CuentaID && v.CuentaID != 0 {
				err = fmt.Errorf("Solo se permite crear subcuentas para una cuenta, no mas. (verificar cuentas_id)")
				return false, err
			}
		}

		if v.Modificado == 1 {
			if v.Tipo == "principal" {
				tipoPrincipalRequest++
			}
			porcentajeRequestsModificado += v.Porcentaje
		} else {
			if v.Tipo == "principal" {
				tipoPrincipalRequest++
			}
			porcentajeRequestsNoModificado += v.Porcentaje
		}
	}

	if cuentasExistentesRecibidas != len(*subcuentasByCuentasId) {
		err = fmt.Errorf("Lantidad de subcuentas enviadas no coinciden con las existentes.")
		return false, err
	}

	if tipoPrincipalRequest != 1 {
		err = fmt.Errorf("Solo puede existir una cuenta de tipo 'principal'.")
		return false, err
	}

	if porcentajeRequestsModificado+porcentajeRequestsNoModificado != 100.0 && !AplicaCostoServicio {
		err = fmt.Errorf("La suma de porcentajes de las cuentas deben alcanzar un 100 porciento. ")
		return false, err
	}

	if AplicaCostoServicio && len(*request) > 2 {
		err = fmt.Errorf("Aplicando a 'Costo de servicio' no es posible tener mas de 2 subcuentas. ")
		return false, err
	}
	return
}
func ValidateRequestDeleted(request *[]administraciondtos.SubcuentaRequest) (ok bool, erro error) {
	porcentajeTotal := 0.0
	subcuentasAEliminar := 0
	contSubcuentaPrincipal := 0
	arrayRequest := *request
	for _, v := range arrayRequest {
		if !v.Eliminado {
			if v.Id == 0 {
				erro = fmt.Errorf("Id de subcuenta en cero o no se envió.")
				return
			}
			if len(v.Nombre) > 35 {
				erro = fmt.Errorf("El nombre no puede contener mas de 35 dígitos.")
				return false, erro
			}
			if v.Tipo == "principal" {
				contSubcuentaPrincipal++
			}
			porcentajeTotal += v.Porcentaje
			if v.Porcentaje == 0 {
				erro = fmt.Errorf("No se admiten porcentajes en cero.")
				return
			}
			if !validateEmail(v.Email) {
				erro = fmt.Errorf("Email inválido.")
				return false, erro
			}
		} else {
			if v.Id == 0 {
				erro = fmt.Errorf("Id de subcuenta a eliminar en cero o no se envió.")
				return
			}
			subcuentasAEliminar++
		}
	}

	if contSubcuentaPrincipal != 1 {
		erro = fmt.Errorf("Solo se permite una subcuenta de tipo principal, ni mas ni menos.")
		return
	}
	if subcuentasAEliminar != 1 {
		erro = fmt.Errorf("Solo se permite eliminar una subcuenta, ni mas ni menos.")
		return
	}
	if porcentajeTotal != 100.0 {
		erro = fmt.Errorf("La suma de porcentajes debe ser de 100 porciento.")
		return
	}
	return
}
func validateEmail(email string) (ok bool) {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	return regexp.MustCompile(regex).MatchString(email)
}
