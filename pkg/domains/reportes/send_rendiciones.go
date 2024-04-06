package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

type sendRendiciones struct {
	utilService util.UtilService
}

func SendRendiciones(util util.UtilService) Email {
	return &sendRendiciones{
		utilService: util,
	}
}

func (cl *sendRendiciones) SendReportes(ruta string, nombreArchivo string, request reportedtos.ResponseClientesReportes) (erro error) {

	if len(request.Rendiciones) > 0 {
		RutaFile := fmt.Sprintf("%s/%s.csv", ruta, nombreArchivo)
		/* estos datos son los que se van a escribir en el archivo */
		var slice_array = [][]string{
			{"", "", "", "", "RECAUDACION WEE", "", "", "", ""},
			{"FECHA", request.Fecha, "", "", "", "", "", "", ""},
			{"CUENTA", "REFERENCIA", "CONCEPTO", "FECHA COBRO", "FECHA DEPOSITO", "IMPORTE COBRADO", "IMPORTE DEPOSITADO", "CANT BOLETAS COBRADAS", "COMISION", "IVA"}, // columnas
		}
		for _, pago := range request.Rendiciones {
			slice_array = append(slice_array, []string{pago.Cuenta, pago.Id, pago.Concepto, pago.FechaCobro, pago.FechaDeposito, pago.ImporteCobrado, pago.ImporteDepositado, pago.CantidadBoletasCobradas, pago.Comision, pago.Iva})
		}
		slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL COBRADO $", request.TotalCobrado})
		slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL IVA $", request.TotalIva})
		slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL COMISION $", request.TotalComision})
		slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL REVERTIDO $", request.TotalRevertido})
		slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL RENDIDO $", request.RendicionTotal})

		/* Crear archivo csv  en la carpeta documentos/reportes*/
		erro = cl.utilService.CsvCreate(RutaFile, slice_array)
		if erro != nil {
			return
		}

	}
	return
}
