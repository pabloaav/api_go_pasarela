package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

type sendPagos struct {
	utilService util.UtilService
}

func SendPagos(util util.UtilService) Email {
	return &sendPagos{
		utilService: util,
	}
}

func (cl *sendPagos) SendReportes(ruta string, nombreArchivo string, request reportedtos.ResponseClientesReportes) (erro error) {

	if len(request.Pagos) > 0 {
		RutaFile := fmt.Sprintf("%s/%s.csv", ruta, nombreArchivo)
		/* estos datos son los que se van a escribir en el archivo */
		var slice_array = [][]string{
			{"TRANSACCIONES COBRADAS"},
			{"FECHA", request.Fecha, "", "", "", "", "", "", ""},
			{"CUENTA", "REFERENCIA", "FECHA COBRO", "MEDIO DE PAGO", "TIPO", "ESTADO", "MONTO"}, // columnas
		}
		for _, pago := range request.Pagos {
			slice_array = append(slice_array, []string{pago.Cuenta, pago.Id, pago.FechaPago, pago.MedioPago, pago.Tipo, pago.Estado, pago.Monto})
		}
		/* Crear archivo csv  en la carpeta documentos/reportes*/
		slice_array = append(slice_array, []string{"", "", "", "", "", "CANT OPERACIONES", request.CantOperaciones})
		slice_array = append(slice_array, []string{"", "", "", "", "", "TOTAL COBRADO $", request.TotalCobrado})
		erro = cl.utilService.CsvCreate(RutaFile, slice_array)
		if erro != nil {
			return
		}

	}
	return
}
