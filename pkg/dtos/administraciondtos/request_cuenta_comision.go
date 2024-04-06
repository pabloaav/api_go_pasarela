package administraciondtos

import (
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type RequestCuentaComision struct {
	Id             uint
	CuentasId      uint
	CuentaComision string
	Comision       float64
	ChannelsId     uint
	//Iva            float64
	VigenciaDesde      time.Time
	ImporteMinimo      float64
	ImporteMaximo      float64
	Mediopagoid        uint
	Pagocuota          bool
	ChannelarancelesId uint
}

func (r *RequestCuentaComision) ToCuentaComision(cargarId bool) (response entities.Cuentacomision) {
	if cargarId {
		response.ID = r.Id
	}

	response.CuentasID = r.CuentasId
	response.Cuentacomision = r.CuentaComision
	response.Comision = r.Comision
	//response.Iva = r.Iva
	response.ChannelsId = r.ChannelsId
	response.VigenciaDesde = &r.VigenciaDesde
	response.Importeminimo = r.ImporteMinimo
	response.Importemaximo = r.ImporteMaximo
	response.Mediopagoid = r.Mediopagoid
	response.Pagocuota = r.Pagocuota
	response.ChannelarancelesId = r.ChannelarancelesId

	return
}

func (r *RequestCuentaComision) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}
	if r.CuentasId < 1 {
		return fmt.Errorf(tools.ERROR_CUENTA_ID)
	}

	if commons.StringIsEmpity(r.CuentaComision) {
		return fmt.Errorf(tools.ERROR_CUENTA_COMISION)
	}

	if r.Comision < 0 {
		return fmt.Errorf(tools.ERROR_COMISION)
	}
	if r.ChannelsId < 1 {
		return fmt.Errorf(tools.ERROR_CHANNEL_ID)
	}

	if r.VigenciaDesde.IsZero() {
		return fmt.Errorf(tools.ERROR_VIGENCIA_DESDE)
	}
	if r.ChannelarancelesId == 0 {
		return fmt.Errorf(tools.ERROR_CHANNEL_ARANCEL)
	}

	// if r.Iva < 0 {
	// 	return fmt.Errorf(tools.ERROR_COMISION)
	// }

	r.CuentaComision = strings.ToUpper(r.CuentaComision)

	return
}
