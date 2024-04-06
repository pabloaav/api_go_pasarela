package filtros

import "errors"

type PagoCaducadoFiltro struct {
	PagoEstadosIds []uint64 `json:"id"`
	CuentaId       uint64   `json:"cuenta_id"`
	CuentaApikey   string   `json:"cuenta_apikey"`
	PagosTipoId    uint64   `json:"pagostipo_id"`
	PagoTipo       string   `json:"pagostipo"`
	MedioPagoId    uint64   `json:"mediopago_id"`
	Referencia     string   `json:"referencia"`
}

func (c *PagoCaducadoFiltro) Validar() (erro error) {
	validacion := false
	if c.CuentaId != 0 && c.PagosTipoId != 0 {
		validacion = true
	}

	if c.CuentaApikey != "" && c.PagoTipo != "" {
		validacion = true
	}

	if !validacion {
		erro = errors.New("Faltan enviar datos de los pagos asociados")
	}
	return
}
