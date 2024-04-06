package pagoofflinedtos

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
)

type OffLineRapipagoDtos struct {
	errors            map[string]error
	codigoEmpresa     string    //longitud: 4
	numeroCliente     string    //longitud: 10
	numeroComprobante string    //longitud: 12
	importe           int64     //longitud: 8
	fechaPrimerVto    time.Time //longitud: 8
	importeRecargo    int64     //longitud: 6
	fechaSegundoVto   time.Time //longitud: 8

}

func New(reques_response OffLineRequestResponse) OffLineRapipagoDtos {
	offline_ov := OffLineRapipagoDtos{}
	offline_ov.errors = make(map[string]error)
	offline_ov.setCodigoEmpresa(reques_response.CodigoEmpresa)
	offline_ov.setNumeroCliente(reques_response.NumeroCliente)
	offline_ov.setNumeroComprobante(reques_response.NumeroComprobante)
	offline_ov.setImporte(reques_response.Importe)
	offline_ov.setFechaPrimerVto(reques_response.FechaPrimerVto)
	offline_ov.setImporteRecargo(reques_response.FechaSegundoVto, reques_response.ImporteRecargo)
	offline_ov.setFechaSegundoVto(reques_response.FechaSegundoVto)
	// offline_ov.SetDigitoVerificador(reques_response.DigitoVerificador)
	return offline_ov
}

func (or *OffLineRapipagoDtos) GetErrors() map[string]string {
	mensajesError := make(map[string]string)
	for key, value := range or.errors {
		mensajesError[key] = value.Error()
	}
	return mensajesError

}

func (or *OffLineRapipagoDtos) GetCodigoEmpresa() string {
	return or.codigoEmpresa
}
func (or *OffLineRapipagoDtos) setCodigoEmpresa(codigoEmpresa string) {
	err := validarIsEmptyAndLong(codigoEmpresa, "codigoEmpresa", 4)
	if err != nil {
		or.errors["codigoEmpresa"] = err
	}
	or.codigoEmpresa = codigoEmpresa

}

func (or *OffLineRapipagoDtos) GetNumeroCliente() string {
	return or.numeroCliente
}

func (or *OffLineRapipagoDtos) setNumeroCliente(numeroCliente string) {
	err := validarIsEmptyAndLong(numeroCliente, "numeroCliente", 10)
	if err != nil {
		or.errors["numeroCliente"] = err
	}
	or.numeroCliente = numeroCliente
}

func (or *OffLineRapipagoDtos) GetNumeroComprobante() string {
	return or.numeroComprobante
}
func (or *OffLineRapipagoDtos) setNumeroComprobante(numeroComprobante string) {
	err := validarIsEmptyAndLong(numeroComprobante, "numeroComprobante", 12)
	if err != nil {
		or.errors["numeroComprobante"] = err
	}
	or.numeroComprobante = numeroComprobante
}

func (or *OffLineRapipagoDtos) GetImporte() int64 {
	return or.importe
}
func (or *OffLineRapipagoDtos) setImporte(importe int64) {
	err := validarImporte(importe, "importe", 8)
	if err != nil {
		or.errors["importe"] = err
	}
	or.importe = importe
}

func (or *OffLineRapipagoDtos) GetFechaPrimerVto() time.Time {
	return or.fechaPrimerVto
}
func (or *OffLineRapipagoDtos) setFechaPrimerVto(fechaPrimerVto time.Time) {
	fechaFin := fechaPrimerVto
	fechaInicio := time.Now()
	err := validarFecha(fechaInicio, fechaFin, "fechaActual", "fechaPrimerVto")
	if err != nil {
		or.errors["fechaPrimerVto"] = err
	}
	or.fechaPrimerVto = fechaPrimerVto
}

func (or *OffLineRapipagoDtos) GetImporteRecargo() int64 {
	return or.importeRecargo
}
func (or *OffLineRapipagoDtos) setImporteRecargo(fechaSegundoVto time.Time, importeRecargo int64) {
	fechaInicio := or.fechaPrimerVto
	fechaFin := fechaSegundoVto
	err := validarSegundoimporte(fechaInicio, fechaFin, importeRecargo, "importeRecargo", 8)
	if err != nil {
		or.errors["importeRecargo"] = err
	}
	or.importeRecargo = importeRecargo
}

func (or *OffLineRapipagoDtos) GetFechaSegundoVto() time.Time {
	return or.fechaSegundoVto
}
func (or *OffLineRapipagoDtos) setFechaSegundoVto(fechaSegundoVto time.Time) {
	fechaInicio := or.fechaPrimerVto
	fechaFin := fechaSegundoVto
	erro := validarFecha(fechaInicio, fechaFin, "fechaPrimerVencimiento", "fechaSegundoVto")
	if erro != nil {
		or.errors["fechaPrimerVto"] = erro
	}
	or.fechaSegundoVto = fechaSegundoVto
}

func validarIsEmptyAndLong(valor string, nombreCampo string, longitud int) error {
	if commons.StringIsEmpity(valor) {
		mensaje := ConstruirMensaje(ERROR_CAMPO_VACIO, nombreCampo)
		return fmt.Errorf("%v", mensaje)
	}
	if len(valor) > longitud {
		return fmt.Errorf("%v", ConstruirMensaje(ERROR_CAMPO_LONGITUD_MAXIMA, nombreCampo, longitud))
	}
	return nil
}

func validarImporte(importe int64, nombreCampo string, longitud int) error {
	if importe <= 0 {
		mensaje := ConstruirMensaje(ERROR_CAMPO_VACIO, nombreCampo)
		return fmt.Errorf("%v", mensaje)
	}
	importeStr := strconv.FormatInt(importe, 10)
	if len(importeStr) > longitud {
		return fmt.Errorf("%v", ConstruirMensaje(ERROR_CAMPO_LONGITUD_MAXIMA, nombreCampo, longitud))
	}
	return nil
}

func validarSegundoimporte(primerFecha, segundaFecha time.Time, importe int64, nombreCampo string, longitud int) error {
	validar := commons.NewAlgoritmoVerificacion()
	if importe < 0 {
		mensaje := ConstruirMensaje(ERROR_CAMPO_VACIO, nombreCampo)
		return fmt.Errorf("%v", mensaje)
	}
	importeStr := strconv.FormatInt(importe, 10)
	if len(importeStr) > longitud {
		return fmt.Errorf("%v", ConstruirMensaje(ERROR_CAMPO_LONGITUD_MAXIMA, nombreCampo, longitud))
	}
	if importe == 0 {
		fechaInicio := primerFecha.Format("2006-01-02")
		fechaFin := segundaFecha.Format("2006-01-02")
		cantidadDias, erro := validar.CalcularDiasEntreFechas(fechaInicio, fechaFin)
		if erro != nil {
			return erro
		}
		if cantidadDias != 0 {
			return fmt.Errorf("%s", ERROR_VALIDACION_sIN_VENCIMIENTO)
		}
	}
	return nil
}

func validarFecha(primerFecha, segundaFecha time.Time, nombreParam1, nombreParam2 string) error {
	validar := commons.NewAlgoritmoVerificacion()
	fechaInicio := primerFecha.Format("2006-01-02")
	fechaFin := segundaFecha.Format("2006-01-02")
	cantidadDias, erro := validar.CalcularDiasEntreFechas(fechaInicio, fechaFin)
	if erro != nil {
		return erro
	}
	if cantidadDias < 0 {
		msg := fmt.Sprintf("%v", ConstruirMensaje(ERROR_VALIDACION_FECHA, nombreParam1, nombreParam2))
		return errors.New(msg)
	}
	return nil
}

// func validarSegundoVencimiento(primerFecha, segundaFecha time.Time, nombreParam1, nombreParam2 string, importe int64, nombreCampo string, longitud int) error {
// 	validar := commons.NewAlgoritmoVerificacion()
// 	fechaInicio := primerFecha.Format("2006-01-02")
// 	fechaFin := segundaFecha.Format("2006-01-02")
// 	cantidadDias, erro := validar.CalcularDiasEntreFechas(fechaInicio, fechaFin)
// 	if erro != nil {
// 		return erro
// 	}
// 	if cantidadDias == 0 && importe == 0 {
// 		return nil
// 	}
// 	mensaje := ConstruirMensaje(ERROR_CAMPO_VACIO, nombreCampo)
// 	return fmt.Errorf("%v", mensaje)

// }
