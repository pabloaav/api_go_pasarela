package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
)

type sendRevertidos struct {
	utilService util.UtilService
}

func SendRevertidos(util util.UtilService) Email {
	return &sendRevertidos{
		utilService: util,
	}
}

func (cl *sendRevertidos) SendReportes(ruta string, nombreArchivo string, request reportedtos.ResponseClientesReportes) (erro error) {

	if len(request.Reversiones) > 0 {
		RutaFile := fmt.Sprintf("%s/%s.csv", ruta, nombreArchivo)
		var slice_array = [][]string{
			{"CUENTA", "IDPAGO", "MEDIOPAGO", "MONTO"}, // columnas
		}
		for _, pago := range request.Reversiones {
			slice_array = append(slice_array, []string{pago.Cuenta, pago.Id, pago.MedioPago, pago.Monto})
		}
		/* Crear archivo csv  en la carpeta documentos/reportes*/
		erro = cl.utilService.CsvCreate(RutaFile, slice_array)
		if erro != nil {
			return
		}
	}
	return
}
