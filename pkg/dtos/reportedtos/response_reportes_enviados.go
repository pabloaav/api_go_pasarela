package reportedtos

import (
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type ResponseReportesEnviados struct {
	Reportes []ResponseReporteEnviado `json:"reportes"`
	Meta     dtos.Meta                `json:"meta"`
}

type ResponseReporteEnviado struct {
	Id             uint
	Cliente        string
	Tiporeporte    string
	Totalcobrado   string
	Totalrendido   string
	Fechacobranza  string
	Fecharendicion string
	Nro_reporte    string
	Detalles       []ResponseReporteDetalle
}

type ResponseReporteDetalle struct {
	Id         uint
	PagosId    string
	Monto      string
	Mediopago  string
	Estadopago string
}

func (rre *ResponseReporteEnviado) EntityToDto(entity entities.Reporte) {
	rre.Id = entity.ID
	rre.Cliente = entity.Cliente
	rre.Tiporeporte = entity.Tiporeporte
	rre.Totalcobrado = entity.Totalcobrado
	rre.Totalrendido = entity.Totalrendido
	rre.Fechacobranza = entity.Fechacobranza
	rre.Fecharendicion = entity.Fecharendicion
	rre.Nro_reporte = strconv.FormatUint(uint64(entity.Nro_reporte), 10)
	rre.Detalles = []ResponseReporteDetalle{}

	if len(entity.Reportedetalle) > 0 {
		respDetalleTemp := ResponseReporteDetalle{}
		for _, detalle := range entity.Reportedetalle {
			respDetalleTemp.ReporteDetalleToDto(detalle)
			rre.Detalles = append(rre.Detalles, respDetalleTemp)
		}
	}

}

func (rrd *ResponseReporteDetalle) ReporteDetalleToDto(entity entities.Reportedetalle) {
	rrd.Id = entity.ID
	rrd.PagosId = entity.PagosId
	rrd.Monto = entity.Monto
	rrd.Mediopago = entity.Mediopago
	rrd.Estadopago = entity.Estadopago
}
