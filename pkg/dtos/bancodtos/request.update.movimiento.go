package bancodtos

import "errors"

type RequestUpdateMovimiento struct {
	ListaMovimientos []uint `json:"lista_movimientos"`
	EstadoMatch      bool   `json:"estado_match"`
}

func (listas *RequestUpdateMovimiento) ValidarListas() error {

	if len(listas.ListaMovimientos) == 0 {
		return errors.New("la lista de movimientos no pueden estar vacia")
	}

	return nil
}


