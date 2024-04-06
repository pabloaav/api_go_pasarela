package util

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	// "image/color"
	"math"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type UtilService interface {
	CreateNotificacionService(notificacion entities.Notificacione) (erro error)
	CreateLogService(log entities.Log) (erro error)
	LogError(erro string, funcionalidad string)
	CrearPeticionesService(peticiones dtos.RequestWebServicePeticion) (erro error)

	//CONFIGURACIONES
	GetConfiguracionService(filtro filtros.ConfiguracionFiltro) (configuracion administraciondtos.ResponseConfiguracion, erro error)
	GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (configuraciones []administraciondtos.ResponseConfiguracion, erro error)
	CreateConfiguracionService(config administraciondtos.RequestConfiguracion) (id uint, erro error)
	FirstOrCreateConfiguracionService(nombre string, descripcion string, valor string) (key string, erro error)

	//Redondeo
	ToFixed(num float64, precision int) float64

	// Convertir Moneda
	FormatNum(num float64) string

	//Calcular Comisiones
	BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error)
	BuildComisionesTemporales(movimiento *entities.Movimientotemporale, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error)
	// impuestos
	GetImpuestoByIdService(id int64) (impuesto entities.Impuesto, erro error)
	CalcularValorCuponService(importe, coeficiente, impuesto float64) (valorCupon int64)
	// CalcularCostoFinancieroIvaService(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA float64) (importeCFMasIva int64)

	// manejo cadena
	RightStr(cadenaStr string, valueStr int) (rightStr string)
	LeftStr(cadenaStr string, valueStr int) (LeftStr string)
	BuildStr(cadenaStr string, valueStr int) (center string)

	GetMatenimietoSistemaService() (estado bool, fecha time.Time, erro error)

	EnviarMailService(params utildtos.RequestDatosMail) (erro error)

	CsvCreate(name string, data [][]string) error

	ValidarCBU(cbu string) (res bool, erro error)

	ValidarCalculoCF(RequestValidarCF utildtos.RequestValidarCF) (responseValidarCF utildtos.ResponseValidarCF)

	GetRecaudacionPdf(request reportedtos.ResultMovLiquidacion, ruta, nombrearchivo string) error

	// Retorna un string que es la concatenacion de la fecha mas cuatro ultimos digitos de la referencia de transferencia bancaria
	CrearReferenciaBancoStandard(referencia_banco string) (string, error)

	// Historial de operaciones del usuario
	CreateHistorialOperacionesService(bearer, observacion string, tipo entities.EnumHistorial)

	// Decodificar token del loginservice
	DecodeToken(bearer string) (userId int64, correo string)

	RunEndpoint(typeEndpoint string, urlEndpoint string, mapConHeaders map[string]string, bodyEndpoint interface{}, queryParameters map[string]string, registarPeticion bool) (objeto interface{}, erro error)
}

var util *utilService

func NewUtilService(r UtilRepository, httpcliente *http.Client) UtilService {
	util := utilService{
		repository: r,
		factory:    &crearMensajeFactory{},
		client:     httpcliente,
	}

	return &util
}

func NewUtilWithService(r UtilRepository, f CrearMensajeFactory, httpcliente *http.Client) UtilService {
	util := utilService{
		repository: r,
		factory:    f,
		client:     httpcliente,
	}

	return &util
}

// Resolve devuelve la instancia antes creada
func Resolve() *utilService {
	return util
}

type utilService struct {
	repository UtilRepository
	factory    CrearMensajeFactory
	client     *http.Client
}

func (r *utilService) CreateNotificacionService(notificacion entities.Notificacione) (erro error) {
	return r.repository.CreateNotificacion(notificacion)

}

func (r *utilService) CreateLogService(log entities.Log) (erro error) {
	return r.repository.CreateLog(log)
}

func (r *utilService) LogError(erro string, funcionalidad string) {

	log := entities.Log{
		Tipo:          entities.Error,
		Mensaje:       erro,
		Funcionalidad: funcionalidad,
	}

	err := r.CreateLogService(log)

	if err != nil {
		mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro)
		logs.Error(mensaje)
	}
}

func (r *utilService) CrearPeticionesService(peticiones dtos.RequestWebServicePeticion) (erro error) {
	peticionEntity := entities.Webservicespeticione{
		Operacion: peticiones.Operacion,
		Vendor:    peticiones.Vendor,
	}
	err := r.repository.CrearPeticionesRepository(peticionEntity)
	if err != nil {
		mensaje := fmt.Sprintf("error al registrar peticion: %s", err.Error())
		erro = errors.New(mensaje)
		logs.Error(mensaje)
		return
	}
	return nil

}

func (s *utilService) GetConfiguracionService(filtro filtros.ConfiguracionFiltro) (configuracion administraciondtos.ResponseConfiguracion, erro error) {

	response, erro := s.repository.GetConfiguracion(filtro)

	if erro != nil {
		return
	}

	configuracion.FromEntity(response)

	return
}

func (s *utilService) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (configuraciones []administraciondtos.ResponseConfiguracion, erro error) {
	response, erro := s.repository.GetConfiguracionesRepository(filtro)
	if erro != nil {
		return
	}
	for _, valueResponse := range response {
		configuracion := administraciondtos.ResponseConfiguracion{}
		configuracion.FromEntity(valueResponse)
		configuraciones = append(configuraciones, configuracion)
	}
	return
}

func (s *utilService) CreateConfiguracionService(config administraciondtos.RequestConfiguracion) (id uint, erro error) {

	erro = config.IsValid(false)

	if erro != nil {
		return
	}

	request := config.ToEntity(false)

	return s.repository.CreateConfiguracion(request)

}

func (s *utilService) FirstOrCreateConfiguracionService(nombre string, descripcion string, valor string) (key string, erro error) {

	if len(strings.TrimSpace(nombre)) < 1 || len(strings.TrimSpace(valor)) < 1 {
		erro = fmt.Errorf("el campo nombre o el campo valor es inválido")
		return
	}

	filtro := filtros.ConfiguracionFiltro{
		Nombre: nombre,
	}

	response, erro := s.GetConfiguracionService(filtro)

	if erro != nil || response.Id == 0 {

		configuracion := administraciondtos.RequestConfiguracion{
			Nombre:      nombre,
			Descripcion: descripcion,
			Valor:       valor,
		}
		_, erro = s.CreateConfiguracionService(configuracion)

		if erro != nil {
			return
		}

		response.Valor = valor
	}

	key = response.Valor

	return

}

func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {

	var descuentos entities.Monto

	/*
	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
	*/
	for _, cc := range *cuentacomisiones {
		if cc.Comision > -1 {
			// var totalMin float64
			// var totalMax float64
			// analizar
			var calculo_fijo bool
			var importeComisiones float64
			var importeComisionTelco float64
			var importeComisionProveedor float64
			var resultImporte float64
			var comisionGeneral float64
			totalMin := cc.Importeminimo + cc.ChannelArancel.Importeminimo
			totalMax := cc.Importemaximo + cc.ChannelArancel.Importemaximo

			switch cc.ChannelArancel.Tipocalculo {

			case "FIJO":
				calculo_fijo = true
				importeComisiones = (float64(movimiento.Monto) / 100 * cc.Comision)    //+ cc.ChannelArancel.Importe // 100
				importeComisionTelco = (float64(movimiento.Monto) / 100 * cc.Comision) // 98
				importeComisionProveedor = cc.ChannelArancel.Importe                   //2
				resultImporte = importeComisionTelco                                   // 100
			case "PORCENTAJE":

				// calcular comision general
				comisionGeneral = cc.Comision + cc.ChannelArancel.Importe
				importeComisiones = (float64(movimiento.Monto) / 100 * comisionGeneral)
				// calcular comision telco y proveedor
				importeComisionTelco = (float64(movimiento.Monto) / 100 * cc.Comision)
				importeComisionProveedor = (float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe)
				// controlar calculos
				// resultImporte = (((float64(movimiento.Monto)/100)*cc.Comision)+((float64(movimiento.Monto)/100)*cc.ChannelArancel.Importe) // ((importeComisionTelco + importeComisionProveedor), 4)
				resultImporte = (importeComisionTelco + importeComisionProveedor)

			}
			logs.Info(fmt.Sprintf("comision general %v - comision separada %v ", c.ToFixed(importeComisiones, 2), c.ToFixed(resultImporte, 2)))
			if c.ToFixed(importeComisiones, 2) == c.ToFixed(resultImporte, 2) {
				var valorPorcentajeTelco float64
				var valorPorcentajeProveedor float64

				if totalMax > 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionComisionConMaximo,
						Descripcion: fmt.Sprintf("comision con maximo. %s", erro.Error()),
					}
					c.CreateNotificacionService(notificacion)
				}

				min, max := VerificarCalculo(importeComisiones, totalMin, totalMax)

				RequestComision := RequestComision{
					ImporteComisionTelco:     importeComisionTelco,
					ImporteComisionProveedor: importeComisionProveedor,
					MinTelco:                 cc.Importeminimo,
					MaxTelco:                 cc.Importemaximo,
					MinProveedor:             cc.ChannelArancel.Importeminimo,
					MaxProveedor:             cc.ChannelArancel.Importemaximo,
					MinBool:                  min,
					MaxBool:                  max,
					Calculo_fijo:             calculo_fijo,
				}

				valorComisionTelcoVerif, valorComisionProveedorVerif := VerificarMinimoMaximo(RequestComision)

				// calculo de comision de telco
				valorPorcentajeTelco = cc.Comision
				if importeComisionTelco != valorComisionTelcoVerif {
					// if cc.ChannelArancel.Tipocalculo == "FIJO" {
					// 	valorComisionTelcoVerif = valorComisionTelcoVerif - cc.ChannelArancel.Importe
					// }
					importeComisionTelco = valorComisionTelcoVerif
					if movimiento.Monto > 0 {
						valorPorcentajeTelco = valorComisionTelcoVerif
					} else {
						valorPorcentajeTelco = -1.00 * valorComisionTelcoVerif
					}
				}

				// calculo de comision de proveedor
				valorPorcentajeProveedor = cc.ChannelArancel.Importe
				if importeComisionProveedor != valorComisionProveedorVerif {
					importeComisionProveedor = valorComisionProveedorVerif
					if movimiento.Monto > 0 {
						valorPorcentajeProveedor = valorComisionProveedorVerif
					} else {
						valorPorcentajeProveedor = -1.00 * valorComisionProveedorVerif
					}
				}

				// PAGOS CON DEBIN DONDE EL TIPO DE CALCULO ES FIJO Y SE CONSIDERA SOLO LA COMISION DE TELCO
				if cc.ChannelArancel.Tipocalculo == "FIJO" {
					importeComisionProveedor = 0
					valorPorcentajeProveedor = 0
				}

				/* Probar calculo de comisiones */
				// result := importeComisionTelco * 100
				// logs.Info(int64(util.ToFixed(result, 2)))
				// fmt.Printf("%T", int64(result))
				movimientoComision := entities.Movimientocomisiones{
					CuentacomisionsID:   cc.ID,
					Monto:               entities.Monto(int64(util.ToFixed(importeComisionTelco*100, 2))),
					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
					Montoproveedor:      entities.Monto(importeComisionProveedor * 100),
					Porcentajeproveedor: valorPorcentajeProveedor,
				}
				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

				//Calculo de Iva sobre comision
				if iva != nil && iva.Porcentaje > -1 {
					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
					movimientoImpuesto := entities.Movimientoimpuestos{
						ImpuestosID:    uint64(iva.ID),
						Monto:          entities.Monto(int64(impuestoIva * 100)),
						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
						Porcentaje:     iva.Porcentaje,
					}
					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
				}
			}
		}
	}
	if importeSolicitado < movimiento.Monto {
		movimiento.Monto = importeSolicitado
	}
	movimiento.Monto -= descuentos

	//Si es un movimiento normal le resto si es una devolucion le sumo.
	// if movimiento.Monto > 0 {
	// 	movimiento.Monto -= descuentos
	// } else {
	// 	movimiento.Monto += descuentos
	// }

	return
}

func (c *utilService) BuildComisionesTemporales(movimiento *entities.Movimientotemporale, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {

	var descuentos entities.Monto

	/*
	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
	*/
	for _, cc := range *cuentacomisiones {
		if cc.Comision > -1 {
			// var totalMin float64
			// var totalMax float64
			// analizar
			var calculo_fijo bool
			var importeComisiones float64
			var importeComisionTelco float64
			var importeComisionProveedor float64
			var resultImporte float64
			var comisionGeneral float64
			totalMin := cc.Importeminimo + cc.ChannelArancel.Importeminimo
			totalMax := cc.Importemaximo + cc.ChannelArancel.Importemaximo

			switch cc.ChannelArancel.Tipocalculo {

			case "FIJO":
				calculo_fijo = true
				importeComisiones = (float64(movimiento.Monto) / 100 * cc.Comision)    //+ cc.ChannelArancel.Importe // 100
				importeComisionTelco = (float64(movimiento.Monto) / 100 * cc.Comision) // 98
				importeComisionProveedor = cc.ChannelArancel.Importe                   //2
				resultImporte = importeComisionTelco                                   // 100
			case "PORCENTAJE":

				// calcular comision general
				comisionGeneral = cc.Comision + cc.ChannelArancel.Importe
				importeComisiones = (float64(movimiento.Monto) / 100 * comisionGeneral)
				// calcular comision telco y proveedor
				importeComisionTelco = (float64(movimiento.Monto) / 100 * cc.Comision)
				importeComisionProveedor = (float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe)
				// controlar calculos
				// resultImporte = (((float64(movimiento.Monto)/100)*cc.Comision)+((float64(movimiento.Monto)/100)*cc.ChannelArancel.Importe) // ((importeComisionTelco + importeComisionProveedor), 4)
				resultImporte = (importeComisionTelco + importeComisionProveedor)

			}
			logs.Info(fmt.Sprintf("comision general %v - comision separada %v ", c.ToFixed(importeComisiones, 2), c.ToFixed(resultImporte, 2)))
			if c.ToFixed(importeComisiones, 2) == c.ToFixed(resultImporte, 2) {
				var valorPorcentajeTelco float64
				var valorPorcentajeProveedor float64

				if totalMax > 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionComisionConMaximo,
						Descripcion: fmt.Sprintf("comision con maximo. %s", erro.Error()),
					}
					c.CreateNotificacionService(notificacion)
				}

				min, max := VerificarCalculo(importeComisiones, totalMin, totalMax)

				RequestComision := RequestComision{
					ImporteComisionTelco:     importeComisionTelco,
					ImporteComisionProveedor: importeComisionProveedor,
					MinTelco:                 cc.Importeminimo,
					MaxTelco:                 cc.Importemaximo,
					MinProveedor:             cc.ChannelArancel.Importeminimo,
					MaxProveedor:             cc.ChannelArancel.Importemaximo,
					MinBool:                  min,
					MaxBool:                  max,
					Calculo_fijo:             calculo_fijo,
				}

				valorComisionTelcoVerif, valorComisionProveedorVerif := VerificarMinimoMaximo(RequestComision)

				// calculo de comision de telco
				valorPorcentajeTelco = cc.Comision
				if importeComisionTelco != valorComisionTelcoVerif {
					// if cc.ChannelArancel.Tipocalculo == "FIJO" {
					// 	valorComisionTelcoVerif = valorComisionTelcoVerif - cc.ChannelArancel.Importe
					// }
					importeComisionTelco = valorComisionTelcoVerif
					if movimiento.Monto > 0 {
						valorPorcentajeTelco = valorComisionTelcoVerif
					} else {
						valorPorcentajeTelco = -1.00 * valorComisionTelcoVerif
					}
				}

				// calculo de comision de proveedor
				valorPorcentajeProveedor = cc.ChannelArancel.Importe
				if importeComisionProveedor != valorComisionProveedorVerif {
					importeComisionProveedor = valorComisionProveedorVerif
					if movimiento.Monto > 0 {
						valorPorcentajeProveedor = valorComisionProveedorVerif
					} else {
						valorPorcentajeProveedor = -1.00 * valorComisionProveedorVerif
					}
				}

				// PAGOS CON DEBIN DONDE EL TIPO DE CALCULO ES FIJO Y SE CONSIDERA SOLO LA COMISION DE TELCO
				if cc.ChannelArancel.Tipocalculo == "FIJO" {
					importeComisionProveedor = 0
					valorPorcentajeProveedor = 0
				}

				/* Probar calculo de comisiones */
				// result := importeComisionTelco * 100
				// logs.Info(int64(util.ToFixed(result, 2)))
				// fmt.Printf("%T", int64(result))
				movimientoComision := entities.Movimientocomisionetemporale{
					CuentacomisionsID:   cc.ID,
					Monto:               entities.Monto(int64(util.ToFixed(importeComisionTelco*100, 2))),
					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
					Montoproveedor:      entities.Monto(importeComisionProveedor * 100),
					Porcentajeproveedor: valorPorcentajeProveedor,
				}
				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

				//Calculo de Iva sobre comision
				if iva != nil && iva.Porcentaje > -1 {
					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
					movimientoImpuesto := entities.Movimientoimpuestotemporale{
						ImpuestosID:    uint64(iva.ID),
						Monto:          entities.Monto(int64(impuestoIva * 100)),
						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
						Porcentaje:     iva.Porcentaje,
					}
					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
				}
			}
		}
	}
	if importeSolicitado < movimiento.Monto {
		movimiento.Monto = importeSolicitado
	}
	movimiento.Monto -= descuentos

	//Si es un movimiento normal le resto si es una devolucion le sumo.
	// if movimiento.Monto > 0 {
	// 	movimiento.Monto -= descuentos
	// } else {
	// 	movimiento.Monto += descuentos
	// }

	return
}

type RequestComision struct {
	ImporteComisionTelco     float64
	ImporteComisionProveedor float64
	MinTelco                 float64
	MaxTelco                 float64
	MinProveedor             float64
	MaxProveedor             float64
	MinBool                  bool
	MaxBool                  bool
	Calculo_fijo             bool
}

func (s *utilService) ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func (s *utilService) FormatNum(num float64) string {
	p := message.NewPrinter(language.Spanish)
	valor := p.Sprintf("%.2f", num)
	return valor
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (s *utilService) GetImpuestoByIdService(id int64) (impuesto entities.Impuesto, erro error) {

	impuesto, err := s.repository.GetImpuestoByIdRepository(id)
	if err != nil {
		erro = errors.New(ERROR_CONFIGURACIONES)
		return
	}
	return
}

func (s *utilService) CalcularValorCuponService(importe, coeficiente, impuesto float64) (valorCupon int64) {
	logs.Info("realizando calculo valor cupon")
	logs.Info(impuesto)
	// importeCupon := importe * coeficiente
	// costoFinancieroNeto := importeCupon - importe
	// valorACobrar := costoFinancieroNeto + (costoFinancieroNeto * impuesto)
	// valorCFmasIVA := s.ToFixed(valorACobrar, 2)
	// //inporteMasCfIva := importe + valorCFmasIVA
	// //finalValor := s.ToFixed(inporteMasCfIva, 2)
	// valorCupon = entities.Monto(s.ToFixed(importe+valorCFmasIVA, 2) * 100).Int64()
	valorCupon = entities.Monto(s.ToFixed(importe*coeficiente, 2) * 100).Int64()
	return
}

// func (s *utilService) CalcularCostoFinancieroIvaService(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA float64) (importeCFMasIva int64) {
// 	importeSinArancel := valorCupon - (valorCupon * porcentajeArancel)
// 	importeSinCostoFinanciero := importeSinArancel / coeficiente
// 	importeTemporalCostoFinanciero := importeSinArancel - importeSinCostoFinanciero
// 	importeTemporalIva := importeTemporalCostoFinanciero * porcentajeIVA
// 	importeTemporalCFMasIva := (importeTemporalCostoFinanciero + importeTemporalIva)
// 	importeredondeo := s.ToFixed(importeTemporalCFMasIva, 2)
// 	importeCFMasIva = entities.Monto(importeredondeo * 100).Int64()
// 	return
// }

func (s *utilService) RightStr(cadenaStr string, valueStr int) (rightStr string) {
	totalStr := len(cadenaStr) - valueStr
	rightStr = cadenaStr[totalStr:]
	return
}

func (s *utilService) LeftStr(cadenaStr string, valueStr int) (LeftStr string) {
	LeftStr = cadenaStr[0:valueStr]
	return
}

func (s *utilService) BuildStr(cadenaStr string, valueStr int) (center string) {
	totalStr := len(cadenaStr) - (valueStr * 2)
	for i := 0; i < totalStr; i++ {
		center += "0"
	}
	return
}

func (s *utilService) GetMatenimietoSistemaService() (estado bool, fecha time.Time, erro error) {
	filtro := filtros.ConfiguracionFiltro{
		Nombre: "ESTADO_APLICACION",
	}
	//estadoConfiguracion, err := s.GetConfiguracionService(filtro)
	response, err := s.repository.GetConfiguracion(filtro)
	logs.Info(fmt.Sprintf("fecha creacion: %v", response.CreatedAt))
	if err != nil {
		estado = true
		erro = fmt.Errorf("el servicio no está disponible")
		return
	}
	if response.Valor != "sin valor" {
		fecha, err = time.Parse("2006-01-02T15:04:00Z", response.Valor)
		if err != nil {
			estado = true

			logs.Error("error al convertir fecha de configuración")
			erro = fmt.Errorf("el servicio no está disponible")
			return
		}
		if !fecha.IsZero() {
			estado = true

			return
		}
	}
	estado = false
	return
}

func (s *utilService) EnviarMailService(params utildtos.RequestDatosMail) (erro error) {
	err := params.IsValid()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	typoEmail, err := params.TipoEmail.IsValid()
	if err != nil {
		return err
	}
	var mensajeCompleto string
	if len(params.CamposReemplazar) != 0 {
		mensajeCompleto = construirMensaje(params.Mensaje, params.CamposReemplazar)
	} else {
		mensajeCompleto = params.Mensaje
	}
	// paramsEmail := utildtos.ParamsEmail{
	// 	Email:   params.Email,
	// 	Nombre:  params.Nombre,
	// 	Mensaje: mensajeCompleto,
	// }
	// t, err := template.ParseFiles("../api/views/templateemail/send_mail.html")
	// //t, err := template.ParseFiles("./api/views/templateemail/send_mail.html")
	// if err != nil {
	// 	logs.Error(err.Error())
	// 	erro = errors.New("error al obtener template" + err.Error())
	// 	return
	// }
	// buf := new(bytes.Buffer)
	// erro = t.Execute(buf, paramsEmail)
	// if erro != nil {
	// 	return erro
	// }

	// Datos del remitente para el email.

	smtpUsername := config.SMTP_USERNAME
	smtpPassword := config.SMTP_PASSWORD

	from := config.EMAIL_FROM
	// password := config.EMAIL_PASS

	// Crear mensaje
	to := params.Email

	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort
	params.Mensaje = mensajeCompleto
	params.From = from
	factoryCrearMensaje, err := s.factory.GetCrearMensajeMethod(typoEmail)
	if err != nil {
		return err
	}
	message, err := factoryCrearMensaje.MensajeResultado(params.Asunto, to, params)
	if err != nil {

		logs.Info("Entro en error control envio de email. construccion del mensaje")
		logs.Info(err.Error())
		return err
	}

	//message := createMessage(params.Asunto, to, from, buf.String())

	//message := createMessage(params.Asunto, to, params.From+" "+from, buf.String())

	// Authentication.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Sending email.

	err = smtp.SendMail(address, auth, params.From, to, []byte(message))

	//err = smtp.SendMail(address, auth, params.From+" "+from, to, []byte(message))

	if err != nil {
		logs.Error("error enviar email - " + err.Error())
		return errors.New(ERROR_ENVIAR_EMAIL)
	}
	logs.Error("email enviado con exito.")

	return
}

// //Función para crear el mensage que se enviará al cliente
// func createMessage(subject string, to []string, from, value string) string {

// 	body := value
// 	header := make(map[string]string)
// 	header["From"] = from
// 	for _, valueTo := range to {
// 		header["To"] = valueTo
// 	}
// 	//header["To"] = to[0]
// 	header["Subject"] = subject
// 	header["MIME-Version"] = "1.0"
// 	header["Content-Type"] = "text/html; charset=\"utf-8\""
// 	header["Content-Transfer-Encoding"] = "base64"

// 	message := ""
// 	for k, v := range header {
// 		message += fmt.Sprintf("%s: %s\r\n", k, v)
// 	}
// 	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
// 	return message

// }

func construirMensaje(mensaje string, campos []string) (mensajeCompleto string) {
	for key, value := range campos {
		mensaje = strings.Replace(mensaje, "#"+strconv.Itoa(key), value, 1)
	}
	mensajeCompleto = mensaje
	return
}

func VerificarMinimoMaximo(request RequestComision) (importeComisionTelco, importeComisionProveedor float64) {

	if request.MinBool && !request.MaxBool {
		if request.MinTelco != 0 && request.MinProveedor != 0 {
			importeComisionTelco = request.MinTelco
			importeComisionProveedor = request.MinProveedor
			return
		}

		if request.MinTelco != 0 && request.MinProveedor == 0 {
			if request.Calculo_fijo {
				importeComisionTelco = request.MinTelco
				importeComisionProveedor = 0
			} else {
				importeComisionTelco = request.MinTelco - util.ToFixed(request.ImporteComisionProveedor, 2)
				// request.ImporteComisionProveedor
				importeComisionProveedor = util.ToFixed(request.ImporteComisionProveedor, 2)
			}
			return
		}

		// if request.MinTelco == 0 && request.MinProveedor != 0 {
		// 	importeComisionTelco = request.ImporteComisionTelco
		// 	importeComisionProveedor = request.MinProveedor
		// 	return
		// }
	}

	if request.MaxBool && !request.MinBool {
		// if request.MaxTelco != 0 && request.MaxProveedor != 0 {
		// 	importeComisionTelco = request.MaxTelco
		// 	importeComisionProveedor = request.MaxProveedor
		// 	return
		// }

		if request.MaxTelco != 0 && request.MaxProveedor == 0 {
			importeComisionTelco = request.MaxTelco
			importeComisionProveedor = request.ImporteComisionProveedor
			return
		}
	}

	// if !request.MinBool && !request.MaxBool {
	// 	if request.MinTelco != 0 && request.MinProveedor != 0 {
	// 		importeComisionTelco = request.MinTelco
	// 		importeComisionProveedor = request.MinProveedor
	// 		return
	// 	}
	// 	if request.MinTelco != 0 && request.MinProveedor == 0 {
	// 		importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
	// 		importeComisionProveedor = request.ImporteComisionProveedor
	// 		return
	// 	}
	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
	// 		importeComisionTelco = request.MaxTelco
	// 		importeComisionProveedor = request.MaxProveedor
	// 		return
	// 	}

	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
	// 		importeComisionTelco = request.MaxTelco
	// 		importeComisionProveedor = request.ImporteComisionProveedor
	// 		return
	// 	}
	// 	importeComisionTelco = request.ImporteComisionTelco
	// 	importeComisionProveedor = request.ImporteComisionProveedor
	// 	return
	// }
	importeComisionTelco = request.ImporteComisionTelco
	importeComisionProveedor = request.ImporteComisionProveedor
	return
}

func VerificarCalculo(comision, minimo, maximo float64) (resultMin, resultMax bool) {

	if minimo != 0 && maximo == 0 {
		if math.Abs(comision) < minimo {
			resultMin = true
		}
		return
	}
	// if minimo == 0 && maximo != 0 {
	// 	if math.Abs(comision) > maximo {
	// 		resultMax = true
	// 	}
	// 	return
	// }
	// if minimo != 0 && maximo != 0 {
	// 	if math.Abs(comision) < minimo {
	// 		resultMin = true
	// 	}
	// 	if math.Abs(comision) > maximo {
	// 		resultMax = true
	// 	}

	// 	return
	// }
	return
}

// convertir datos  a excel // utilizado para enviar reportes a clientes
func (s *utilService) CsvCreate(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()

	for _, d := range data {
		err := w.Write(d)
		if err != nil {
			return err
		}
	}

	return nil
}

// validar un CBU de acuerdo a sus diferentes posibilidades
func (s *utilService) ValidarCBU(cbu string) (res bool, erro error) {
	erro = validarLargoCbu(cbu)
	if erro != nil {
		return
	}
	erro = validarCodigoBanco(cbu[0:8])
	if erro != nil {
		return
	}
	erro = validarCuenta(cbu[8:22])
	if erro != nil {
		return
	}
	res = true
	return
}

func (s *utilService) ValidarCalculoCF(RequestValidarCF utildtos.RequestValidarCF) (responseValidarCF utildtos.ResponseValidarCF) {

	// se calcula el valor neto del cupon "Net"
	cuponNeto := RequestValidarCF.Cupon.Float64() * (1 - RequestValidarCF.ArancelMonto)
	// se calcula el valor de la cuota
	montoCuota := cuponNeto / RequestValidarCF.Cuotas

	air := RequestValidarCF.Tna / 100
	air_12 := air / 12

	//calculos intermedios
	intermedio := 1 + (air_12)
	exponente := RequestValidarCF.Cuotas - 1
	intermedio2 := (1 + (air)*((30-RequestValidarCF.Dias)/360))
	intermedio3 := 1 - 1/intermedio
	intermedio4 := 1 - 1/math.Pow(intermedio, exponente)
	// calcular terminso formula CF
	primerTermino := 1 / intermedio2

	segundoTermino := 1 / intermedio * (intermedio4 / intermedio3 / intermedio2)
	if math.IsNaN(segundoTermino) {
		segundoTermino = 0
	}

	// formula calculo CF
	responseValidarCF.CostoFinanciero = s.ToFixed(cuponNeto-montoCuota*(primerTermino+segundoTermino), 2)
	// calculo valor presente
	valorPresenteTemporal := cuponNeto - responseValidarCF.CostoFinanciero
	responseValidarCF.ValorPresente = s.ToFixed(valorPresenteTemporal, 2)
	// calculo ceficiente formula1
	coeficienteUno := s.ToFixed(cuponNeto/responseValidarCF.ValorPresente, 4)
	coeficienteDos := s.ToFixed(1/(1-responseValidarCF.CostoFinanciero/RequestValidarCF.Cupon.Float64()/(1-RequestValidarCF.ArancelMonto)), 4)
	//calculo costo total
	responseValidarCF.CostoTotalPorcentaje = s.ToFixed((responseValidarCF.CostoFinanciero/RequestValidarCF.Cupon.Float64())*100, 2)
	responseValidarCF.ValorCoeficiente = coeficienteUno
	logs.Info(coeficienteDos)

	return
}

// Retorna un string que es la concatenacion de la fecha mas cuatro ultimos digitos de la referencia de transferencia bancaria
func (s *utilService) CrearReferenciaBancoStandard(referencia_banco string) (string, error) {
	fecha := referencia_banco[0:10]
	cuatro_ultimos_digitos := referencia_banco[len(referencia_banco)-4:]
	referencia_banco = fecha + cuatro_ultimos_digitos
	if len(referencia_banco) != 14 {
		return "", fmt.Errorf("la longitud de la referencia banco %s no es correcta. Longitud: %d", referencia_banco, len(referencia_banco))
	}
	return referencia_banco, nil
}

// validar la longitud del CBU
func validarLargoCbu(cbu string) error {
	if commons.StringIsEmpity(cbu) {
		return fmt.Errorf("cbu está en blanco")
	}
	if len(cbu) != 22 {
		return fmt.Errorf("longitud de cbu no es válido: %d", len(cbu))
	}
	return nil
}

// validar la parte del Codigo de Banco del CBU que son los 8 primeros digitos
func validarCodigoBanco(codigo string) error {
	if len(codigo) != 8 {
		return fmt.Errorf("el código de banco es incorrecto")
	}
	banco := codigo[0:3] // numero de entidad

	digitoVerificador := codigo[3:4]

	sucursal := codigo[4:7] // numero de sucursal
	// fmt.Println("sucursal: " + sucursal)
	digitoVerificador2 := codigo[7:8]
	// fmt.Println("digito verificador 2: " + digitoVerificador2)

	var suma int
	var x int

	x, _ = strconv.Atoi(banco[0:1])
	suma = x * 7
	x, _ = strconv.Atoi(banco[1:2])
	suma = suma + x
	x, _ = strconv.Atoi(banco[2:3])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(digitoVerificador)
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(sucursal[0:1])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(sucursal[1:2])
	suma = suma + x
	x, _ = strconv.Atoi(sucursal[2:3])
	suma = suma + (x * 3)

	var diferencia = (10 - (suma % 10)) % 10
	digito, _ := strconv.Atoi(digitoVerificador2)
	if diferencia != digito {
		return fmt.Errorf("codigo de banco inválido")
	}
	return nil
}

// validar la parte del Cuenta del CBU que son los 14 ultimos digitos
func validarCuenta(cuenta string) error {
	if len(cuenta) != 14 {
		return fmt.Errorf("logitud de cuenta inválido: %d", len(cuenta))
	}
	digitoVerificador, _ := strconv.Atoi(cuenta[13:14])

	var suma int
	var x int

	x, _ = strconv.Atoi(cuenta[0:1])
	suma = x * 3
	x, _ = strconv.Atoi(cuenta[1:2])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[2:3])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[3:4])
	suma = suma + x
	x, _ = strconv.Atoi(cuenta[4:5])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(cuenta[5:6])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[6:7])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[7:8])
	suma = suma + (x * 1)
	x, _ = strconv.Atoi(cuenta[8:9])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(cuenta[9:10])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[10:11])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[11:12])
	suma = suma + (x * 1)
	x, _ = strconv.Atoi(cuenta[12:13])
	suma = suma + (x * 3)

	var diferencia = (10 - (suma % 10)) % 10

	if diferencia != digitoVerificador {
		return fmt.Errorf("error en cuenta bancaria")
	}

	return nil
}

type recaudacion struct {
	m pdf.Maroto
}

type ActorRecaudacion struct {
	RazonSocial    string `json:"razonSocial"`
	Domicilio      string `json:"domicilio"`
	Cuit           string `json:"cuit"`
	IngresosBrutos string `json:"ingresosBrutos"`
	Iva            string `json:"iva"`
	Cuenta         string `json:"cuenta"`
	NroOrden       string `json:"nro_orden"`
}

type FechasRecaudacion struct {
	Cobro    string `json:"cobro"`
	Deposito string `json:"deposito"`
	Proceso  string `json:"proceso"`
}

type DataBodyRecaudacion struct {
	ChannelName       string `json:"channel_name"`
	ImporteCobrado    string `json:"importe_cobrado"`
	ImporteDepositado string `json:"importe_depositado"`
	CantidadBoletas   string `json:"cantidad_boletas"`
	Comisiones        string `json:"comisiones"`
	IvaComision       string `json:"iva_comision"`
	RetIva            string `json:"ret_iva"`
}
type DataFooterRecaudacion struct {
	ImporteCobrado    string `json:"Importe cobrado"`
	CantidadBoletas   string `json:"Cantidad de boletas"`
	ComisionTotal     string `json:"Comisión total"`
	IvaTotal          string `json:"IVA total"`
	ImporteADepositar string `json:"Importe a depositar"`
}

type Subtotales struct {
	ImporteCobradoDebito    string
	ImporteADepositarDebito string
	CantidadBoletasDebito   string

	ImporteCobradoCredito    string
	ImporteADepositarCredito string
	CantidadBoletasCredito   string

	ImporteCobradoDebin    string
	ImporteADepositarDebin string
	CantidadBoletasDebin   string
}

func (r *recaudacion) buildHeading(cliente ActorRecaudacion, fechas FechasRecaudacion, fileName string) {

	// Logo
	r.m.Row(20, func() {
		r.m.Col(3, func() {
			_ = r.m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "wee_por_telco.png"), props.Rect{
				Center:  true,
				Percent: 75,
			})
		})
	})

	// Titulo
	r.m.Row(10, func() {

		r.m.Col(12, func() {
			r.m.Text("Recaudación WEE! Liq Nro: "+cliente.NroOrden, props.Text{
				Top:    3,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   12,
				Align:  consts.Center,
				Color:  getDarkGrayColor(),
			})
		})
	})

	r.m.Line(1)

	// Datos cliente
	r.m.Row(40, func() {
		r.m.ColSpace(1)

		r.m.Col(2, func() {

			r.m.Text("Cliente: ", props.Text{
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

		r.m.Col(6, func() {

			r.m.Text(cliente.RazonSocial, props.Text{
				Top:    7,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})

			r.m.Text("Domicilio: "+cliente.Domicilio, props.Text{
				Top:    14,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
			r.m.Text("C.U.I.T.: "+cliente.Cuit, props.Text{
				Left:   2,
				Top:    21,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
			r.m.Text("Ing. Brutos: "+cliente.IngresosBrutos, props.Text{
				Left:   2,
				Top:    28,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})

			r.m.Text("IVA: "+"RESPONSABLE INSCRIPTO", props.Text{
				Left:   2,
				Top:    35,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})

		})
		//fecha_impresion_pdf := fmt.Sprintf(time.Now().Format("02-01-2006"))
		r.m.Col(3, func() {
			r.m.Text("Fecha: "+fechas.Proceso, props.Text{
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

	})
	r.m.Line(1)
	// fila 2 fechas
	r.m.Row(10, func() {

		// r.m.ColSpace(4)
		// col3 fecha de proceso
		r.m.Col(4, func() {
			r.m.Text("Fecha Proceso: "+fechas.Proceso, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

		// col2 fecha de deposito
		r.m.Col(4, func() {
			r.m.Text("Fecha Depósito: "+fechas.Deposito, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

		r.m.Col(4, func() {
			r.m.Text("Cuenta: "+cliente.Cuenta, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})

	})

	r.m.Row(10, func() {
		r.m.Col(12, func() {
			r.m.Text("Nombre de Archivo Transferido: "+fileName, props.Text{
				Top:    3,
				Left:   2,
				Style:  consts.Bold,
				Family: consts.Helvetica,
				Size:   8,
				Color:  getDarkGrayColor(),
			})
		})
	})

}

func (r *recaudacion) buildContent(header []string, contenido [][]string, contenidocredit [][]string, contenidodebin [][]string, subtotales Subtotales) {
	grayColor := getGrayColor()
	r.m.Row(10, func() {
		r.m.Col(12, func() {
			r.m.Text("Detalle de valores a Recaudar", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	if len(contenido) == 0 && len(contenidocredit) == 0 && len(contenidodebin) == 0 {
		r.m.Row(5, func() {
			r.m.Col(7, func() {
				r.m.Text("No se registraron movimientos", props.Text{
					Left:  2,
					Top:   1.5,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
	}

	if len(contenido) > 0 {
		r.m.SetBackgroundColor(getDarkGrayColor())

		r.m.Row(7, func() {
			r.m.Col(3, func() {
				r.m.Text("Debito", props.Text{
					Top:   1.5,
					Size:  9,
					Style: consts.Bold,
					Align: consts.Center,
					Color: color.NewWhite(),
				})
			})
			r.m.ColSpace(9)
		})

		r.m.SetBackgroundColor(color.NewWhite())

		r.m.TableList(header, contenido, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			ContentProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			Align:                consts.Center,
			AlternatedBackground: &grayColor,
			HeaderContentSpace:   1,
			Line:                 false,
		})
		// SUBTOTALES
		r.m.Row(5, func() {
			r.m.Col(7, func() {
				r.m.Text("Subtotales:", props.Text{
					Left:  2,
					Top:   1.5,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("Cantidad de boletas:"+subtotales.CantidadBoletasDebito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A cobrar:"+subtotales.ImporteCobradoDebito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A depositar: "+subtotales.ImporteADepositarDebito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
	}

	if len(contenidocredit) > 0 {
		r.m.SetBackgroundColor(getDarkGrayColor())

		r.m.Row(7, func() {
			r.m.Col(3, func() {
				r.m.Text("Credito", props.Text{
					Top:   1.5,
					Size:  9,
					Style: consts.Bold,
					Align: consts.Center,
					Color: color.NewWhite(),
				})
			})
			r.m.ColSpace(9)
		})

		r.m.SetBackgroundColor(color.NewWhite())

		r.m.TableList(header, contenidocredit, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			ContentProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			Align:                consts.Center,
			AlternatedBackground: &grayColor,
			HeaderContentSpace:   1,
			Line:                 false,
		})
		// SUBTOTALES
		r.m.Row(5, func() {
			r.m.Col(7, func() {
				r.m.Text("Subtotales:", props.Text{
					Left:  2,
					Top:   1.5,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})

		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("Cantidad de boletas:"+subtotales.CantidadBoletasCredito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A cobrar:"+subtotales.ImporteCobradoCredito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A depositar: "+subtotales.ImporteADepositarCredito, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
	}

	if len(contenidodebin) > 0 {
		r.m.SetBackgroundColor(getDarkGrayColor())

		r.m.Row(7, func() {
			r.m.Col(3, func() {
				r.m.Text("Debin", props.Text{
					Top:   1.5,
					Size:  9,
					Style: consts.Bold,
					Align: consts.Center,
					Color: color.NewWhite(),
				})
			})
			r.m.ColSpace(9)
		})

		r.m.SetBackgroundColor(color.NewWhite())

		r.m.TableList(header, contenidodebin, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			ContentProp: props.TableListContent{
				Size:      8,
				GridSizes: []uint{2, 2, 2, 2, 1, 1, 1, 1},
			},
			Align:                consts.Center,
			AlternatedBackground: &grayColor,
			HeaderContentSpace:   1,
			Line:                 false,
		})
		// SUBTOTALES
		r.m.Row(5, func() {
			r.m.Col(7, func() {
				r.m.Text("Subtotales:", props.Text{
					Left:  2,
					Top:   1.5,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})

		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("Cantidad de boletas:"+subtotales.CantidadBoletasDebin, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A cobrar:"+subtotales.ImporteCobradoDebin, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
		r.m.Row(5, func() {
			r.m.ColSpace(7)
			r.m.Col(5, func() {
				r.m.Text("A depositar: "+subtotales.ImporteADepositarDebin, props.Text{
					Left:  2,
					Size:  8,
					Align: consts.Left,
					Color: getDarkGrayColor(),
				})
			})
		})
	}
}

func (r *recaudacion) buildFooter(dataFooter DataFooterRecaudacion) {
	r.m.SetFirstPageNb(1)
	r.m.RegisterFooter(func() {

		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Importe total cobrado: "+dataFooter.ImporteCobrado, props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 13})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Importe total a depositar: "+dataFooter.ImporteADepositar, props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 16})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Cantidad de boletas: "+dataFooter.CantidadBoletas, props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 19})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("Comisiones a cobrar: "+dataFooter.ComisionTotal, props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 19})
			})
		})
		r.m.Row(7, func() {
			r.m.Col(12, func() {
				r.m.Text("IVA a cobrar: "+dataFooter.IvaTotal, props.Text{Align: consts.Left, Left: 2, Size: 8, Top: 19})
			})
		})

	})
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   51,
		Green: 51,
		Blue:  51,
	}
}

func (s *utilService) GetRecaudacionPdf(request reportedtos.ResultMovLiquidacion, ruta, nombrearchivo string) error {
	recaudacionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	recaudacionPdf.SetPageMargins(20, 10, 20)

	var rec recaudacion
	rec.m = recaudacionPdf

	/************* Cabecera  *************/
	var cliente ActorRecaudacion
	var fechas FechasRecaudacion

	cliente.RazonSocial = request.Cabeceras.RazonSocial
	cliente.Domicilio = "San Martín 907, Corrientes"
	cliente.Cuit = request.Cabeceras.Cuit
	cliente.IngresosBrutos = ""
	cliente.Iva = "RESPONSABLE INSCRIPTO"
	cliente.NroOrden = request.NroLiquidacion
	cliente.Cuenta = request.Cuenta

	// fechas.Cobro = "24/11/22"
	fechas.Deposito = request.FechaRendicion
	fechas.Proceso = request.FechaProceso

	fileName := nombrearchivo

	// encabezdo
	rec.buildHeading(cliente, fechas, fileName)

	/************* Cuerpo *************/
	// TODO: cambiar por datos reales
	header, debito, credito, debin := getContentsAndHeadersTable(request.MedioPagoItems.MedioMovDebit, request.MedioPagoItems.MedioMovCredit, request.MedioPagoItems.MedioMovDebin)

	var subtotales Subtotales
	subtotales.ImporteCobradoDebito = request.MedioPagoItems.MedioMovDebit.TotalCobrado
	subtotales.ImporteADepositarDebito = request.MedioPagoItems.MedioMovDebit.TotalaRendir
	subtotales.CantidadBoletasDebito = request.MedioPagoItems.MedioMovDebit.CantidaTotaldBoletas

	subtotales.ImporteCobradoCredito = request.MedioPagoItems.MedioMovCredit.TotalCobrado
	subtotales.ImporteADepositarCredito = request.MedioPagoItems.MedioMovCredit.TotalaRendir
	subtotales.CantidadBoletasCredito = request.MedioPagoItems.MedioMovCredit.CantidaTotaldBoletas

	subtotales.ImporteCobradoDebin = request.MedioPagoItems.MedioMovDebin.TotalCobrado
	subtotales.ImporteADepositarDebin = request.MedioPagoItems.MedioMovDebin.TotalaRendir
	subtotales.CantidadBoletasDebin = request.MedioPagoItems.MedioMovDebin.CantidaTotaldBoletas

	// contenido
	rec.buildContent(header, debito, credito, debin, subtotales)

	rec.m.Line(1)

	/************* Footer  *************/

	var dataFooter DataFooterRecaudacion
	dataFooter.ImporteCobrado = request.Totales.ImporteCobrado
	dataFooter.CantidadBoletas = request.Totales.CantidadTotalBoletas
	dataFooter.ComisionTotal = request.Totales.ComisionTotal
	dataFooter.IvaTotal = request.Totales.IvaTotal
	dataFooter.ImporteADepositar = request.Totales.ImporteADepositar
	// pie de informe
	rec.buildFooter(dataFooter)

	// Se crea la carpeta temporal en caso de que no exista
	// carpetaTempImages := "./pdfs"
	// if _, err := os.Stat(carpetaTempImages); os.IsNotExist(err) {
	// 	err = os.MkdirAll(carpetaTempImages, 0755)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// crear archivo en la carpeta pdfs
	err := rec.m.OutputFileAndClose(ruta + "/" + nombrearchivo + ".pdf")
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
		return err
	}

	return nil
}

// 8 columnas
func getContentsAndHeadersTable(debito reportedtos.MedioMovDebit, credito reportedtos.MedioMovCredit, debin reportedtos.MedioMovDebin) ([]string, [][]string, [][]string, [][]string) {
	header := []string{"Cuenta", "Referencia", "Fecha Cobro", "Importe Cobrado", "Importe a Depositar", "Cant Boletas", "Comision", "IVA"}

	contentsdebito := [][]string{}
	contentscredito := [][]string{}
	contentsdebin := [][]string{}

	// debito
	if len(debito.Detalle) > 0 {
		for _, item := range debito.Detalle {
			contentsdebito = append(contentsdebito, []string{item.Cuenta, item.Referencia, item.FechaCobro, item.ImporteCobrado, item.ImporteADepositar, item.CantidadBoletas, item.Comision, item.Iva})
		}
	}

	//credito
	if len(credito.Detalle) > 0 {
		for _, cr := range credito.Detalle {
			contentscredito = append(contentscredito, []string{cr.Cuenta, cr.Referencia, cr.FechaCobro, cr.ImporteCobrado, cr.ImporteADepositar, cr.CantidadBoletas, cr.Comision, cr.Iva})
		}
	}
	//debin
	if len(debin.Detalle) > 0 {
		for _, de := range debin.Detalle {
			contentsdebin = append(contentsdebin, []string{de.Cuenta, de.Referencia, de.FechaCobro, de.ImporteCobrado, de.ImporteADepositar, de.CantidadBoletas, de.Comision, de.Iva})
		}
	}
	// contenido de la tabla
	// for i := 0; i < 15; i++ {
	// 	contents = append(contents, []string{"TASAS MUNICIPALES", "Jf5ftsAQG", "02/01/2023", "6670.2", "6549.14", "2", "100.05", "21.01"})
	// }

	return header, contentsdebito, contentscredito, contentsdebin
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func (s *utilService) CreateHistorialOperacionesService(bearer, observacion string, tipo entities.EnumHistorial) {
	userId, correo := s.DecodeToken(bearer)

	historial := entities.HistorialOperaciones{
		UsersId:       int(userId),
		Correo:        correo,
		TipoOperacion: tipo,
		Observaciones: observacion,
	}

	err := s.repository.CreateHistorialOperacionesRepository(historial)

	if err != nil {
		mensaje := fmt.Sprintf("Crear Log: %s", err.Error())
		logs.Error(mensaje)
	}

	return
}

func (s utilService) DecodeToken(bearer string) (userId int64, email string) {
	// Separar el token de la palabra "Bearer".
	parts := strings.Split(bearer, "Bearer ")

	if len(parts) != 2 {
		fmt.Println("El formato del token no es válido.")
		return
	}

	// El token real se encuentra en parts[1].
	token := parts[1]

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.TOKEN_KEY), nil
	})

	if err != nil {
		err = errors.New("El token ingresado no es válido")
		return
	}

	if t.Valid {
		claims := t.Claims.(jwt.MapClaims)
		userId, _ := strconv.ParseInt(claims["id"].(string), 10, 64)
		userName, _ := claims["user"].(string)

		return userId, userName
	}

	return
}

func (r *utilService) RunEndpoint(typeEndpoint string, urlEndpoint string, mapConHeaders map[string]string, bodyEndpoint interface{}, queryParameters map[string]string, registarPeticion bool) (objeto interface{}, erro error) {
	// Construye la URL base
	base, err := url.Parse(urlEndpoint)
	if err != nil {
		logs.Error("ERROR_URL" + err.Error())
		return objeto, err
	}

	// Agrega parámetros de consulta a la URL
	q := base.Query()
	for key, value := range queryParameters {
		q.Add(key, value)
	}
	base.RawQuery = q.Encode()

	json_data, err := json.Marshal(bodyEndpoint)
	if err != nil {
		fmt.Println("Error al convertir a JSON: ", err.Error())
		return objeto, err
	}

	// Crea la solicitud HTTP
	req, _ := http.NewRequest(typeEndpoint, base.String(), bytes.NewBuffer(json_data))

	buildHeaderUniversal(req, mapConHeaders)
	// Ejecuta la solicitud y maneja la respuesta
	err = executeRequest(r, req, ERROR_RUN_ENDPOINTS, &objeto)

	// Registra la petición realizada
	if registarPeticion {
		peticionApiLink := dtos.RequestWebServicePeticion{
			Operacion: urlEndpoint,
			Vendor:    "Vendor",
		}
		err1 := r.CrearPeticionesService(peticionApiLink)
		if err1 != nil {
			logs.Error("ERROR_CREAR_PETICION" + err1.Error())
		}
	}

	if err != nil {
		return objeto, err
	}
	return objeto, nil
}

func buildHeaderUniversal(request *http.Request, nameValueHeaders map[string]string) {
	_, hasAccept := nameValueHeaders["accept"]
	_, hasContentType := nameValueHeaders["content-type"]

	if !hasAccept {
		request.Header.Add("accept", "application/json")
	}

	if !hasContentType {
		request.Header.Add("content-type", "application/json")
	}

	// Agrega los demás headers del map
	for name, value := range nameValueHeaders {
		request.Header.Add(name, value)
	}

}

func executeRequest(r *utilService, req *http.Request, erro string, objeto interface{}) error {

	resp, err := r.client.Do(req)

	if err != nil {
		logs.Error(err.Error())
		return errors.New(erro)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		logs.Error(fmt.Sprint(resp.Status))
		return nil
	}

	if resp.StatusCode != 200 && resp.StatusCode != 202 {

		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "executeRequestCuentas",
		}

		apiError := linkdtos.ErrorApiLink{}

		if resp.StatusCode == 500 {

			apiError.Codigo = "500"

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprint(resp.Status)

		}

		err := json.NewDecoder(resp.Body).Decode(&apiError)

		if err != nil {

			apiError.Codigo = strconv.Itoa(resp.StatusCode)

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprintf("%s, %s", erro, resp.Status)

		}

		if resp.StatusCode == 401 {
			apiError.Codigo = "401"
			apiError.Descripcion = "Unauthorized"
		}

		r.CreateLogService(log)

		return &apiError
	}

	err = json.NewDecoder(resp.Body).Decode(&objeto)

	if err != nil {
		return err
	}

	return nil

}

// mes de junio
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto) (erro error) {

// 	var descuentos entities.Monto
// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {

// 			comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)

// 			movimientoComision := entities.Movimientocomisiones{
// 				CuentacomisionsID: cc.ID,
// 				Monto:             entities.Monto(int64(comision * 100)),
// 				Porcentaje:        cc.Comision,
// 			}
// 			descuentos += movimientoComision.Monto
// 			movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 			//Calculo de Iva sobre comision
// 			if iva != nil && iva.Porcentaje > -1 {
// 				impuestoIva := c.ToFixed((comision * iva.Porcentaje), 2)
// 				movimientoImpuesto := entities.Movimientoimpuestos{
// 					ImpuestosID: uint64(iva.ID),
// 					Monto:       entities.Monto(int64(impuestoIva * 100)),
// 					Porcentaje:  iva.Porcentaje,
// 				}
// 				descuentos += movimientoImpuesto.Monto
// 				movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 			}
// 		}
// 	}
// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	if movimiento.Monto > 0 {
// 		movimiento.Monto -= descuentos
// 	} else {
// 		movimiento.Monto += descuentos
// 	}

// 	return
// }

/*
	// if comision < cc.Importeminimo && cc.Importeminimo != 0 {
		// 	importeComision = cc.Importeminimo
		// } else if comision > cc.Importemaximo && cc.Importemaximo != 0 {
		// 	importeComision = cc.Importemaximo
		// } else {
		// 	importeComision = comision
		// }
*/

//06-10-2022
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto) (erro error) {

// 	var descuentos entities.Monto
// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {
// 			// calcular comision general

// 			// calcular comision telco y proveedor

// 			// controlar calculos

// 			// verificar si el importe de la comision no supera el minimo,
// 			// supera el maximo o se encuetra entre el minimo y el maximo
// 			// de telco y proveedor

// 			var valorPorcentaje float64
// 			comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)
// 			importeComision := VerificarMinimoMaximo(comision, cc.Importeminimo, cc.Importemaximo)

// 			valorPorcentaje = cc.Comision
// 			if comision != importeComision {
// 				if movimiento.Monto > 0 {
// 					valorPorcentaje = importeComision
// 				} else {
// 					valorPorcentaje = -1.00 * importeComision
// 				}
// 			}

// 			movimientoComision := entities.Movimientocomisiones{
// 				CuentacomisionsID: cc.ID,
// 				Monto:             entities.Monto(int64(importeComision * 100)),
// 				Porcentaje:        valorPorcentaje, //cc.Comision,
// 			}
// 			descuentos += movimientoComision.Monto
// 			movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 			//Calculo de Iva sobre comision
// 			if iva != nil && iva.Porcentaje > -1 {
// 				impuestoIva := c.ToFixed((importeComision * iva.Porcentaje), 2)
// 				movimientoImpuesto := entities.Movimientoimpuestos{
// 					ImpuestosID: uint64(iva.ID),
// 					Monto:       entities.Monto(int64(impuestoIva * 100)),
// 					Porcentaje:  iva.Porcentaje,
// 				}
// 				descuentos += movimientoImpuesto.Monto
// 				movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 			}
// 		}
// 	}
// 	movimiento.Monto -= descuentos

// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	// if movimiento.Monto > 0 {
// 	// 	movimiento.Monto -= descuentos
// 	// } else {
// 	// 	movimiento.Monto += descuentos
// 	// }

// 	return
// }

// 14-10-2022
// // valores porcentual del proveedor
// 			// calcular comision general
// 			comisionGeneral := cc.Comision + cc.ChannelArancel.Importe
// 			importeComisiones := c.ToFixed((float64(movimiento.Monto) / 100 * comisionGeneral), 4)
// 			// calcular comision telco y proveedor
// 			importeComisionTelco := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)
// 			importeComisionProveedor := c.ToFixed((float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe), 4)
// 			// controlar calculos
// 			resultImporte := c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)
// 			if importeComisiones == resultImporte {
// 				//var valorPorcentaje float64
// 				var valorPorcentajeTelco float64
// 				var valorPorcentajeProveedor float64
// 				// comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)
// 				// importeComision := VerificarMinimoMaximo(comisionGeneral, cc.Importeminimo, cc.Importemaximo)

// 				// valorPorcentaje = cc.Comision
// 				// if importeComisiones != importeComision {
// 				// 	if movimiento.Monto > 0 {
// 				// 		valorPorcentaje = importeComision
// 				// 	} else {
// 				// 		valorPorcentaje = -1.00 * importeComision
// 				// 	}
// 				// }
// 				//verifico valor minimo y maximo de telco
// 				importeComisionTelcoVerif := VerificarMinimoMaximo(cc.Comision, cc.Importeminimo, cc.Importemaximo)
// 				valorPorcentajeTelco = cc.Comision
// 				if importeComisionTelco != importeComisionTelcoVerif {
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeTelco = importeComisionTelcoVerif
// 					} else {
// 						valorPorcentajeTelco = -1.00 * importeComisionTelcoVerif
// 					}
// 				}
// 				//verifico valor minimo y maximo del proveedor
// 				// cc.ChannelArancel.Importe es igual a la comision del proveedor
// 				importeComisionProveedorVerif := VerificarMinimoMaximo(cc.ChannelArancel.Importe, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo)
// 				valorPorcentajeProveedor = cc.ChannelArancel.Importe
// 				if importeComisionProveedor != importeComisionProveedorVerif {
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeProveedor = importeComisionProveedorVerif
// 					} else {
// 						valorPorcentajeProveedor = -1.00 * importeComisionProveedorVerif
// 					}
// 				}

// 				movimientoComision := entities.Movimientocomisiones{
// 					CuentacomisionsID:   cc.ID,
// 					Monto:               entities.Monto(int64(importeComisionTelco * 100)),
// 					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
// 					Montoproveedor:      entities.Monto(int64(importeComisionProveedor * 100)),
// 					Porcentajeproveedor: valorPorcentajeProveedor,
// 				}
// 				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
// 				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 				//Calculo de Iva sobre comision
// 				if iva != nil && iva.Porcentaje > -1 {
// 					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
// 					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
// 					movimientoImpuesto := entities.Movimientoimpuestos{
// 						ImpuestosID:    uint64(iva.ID),
// 						Monto:          entities.Monto(int64(impuestoIva * 100)),
// 						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
// 						Porcentaje:     iva.Porcentaje,
// 					}
// 					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
// 					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 				}
// 			}

// 		}

// 27/12/22
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {

// 	var descuentos entities.Monto

// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {
// 			// analizar
// 			var importeComisiones float64
// 			var importeComisionTelco float64
// 			var importeComisionProveedor float64
// 			var resultImporte float64
// 			switch cc.ChannelArancel.Tipocalculo {
// 			case "FIJO":

// 				importeComisiones = c.ToFixed(((float64(movimiento.Monto) / 100 * cc.Comision) + cc.ChannelArancel.Importe), 4) //+ cc.ChannelArancel.Importe // 100
// 				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)                            // 98
// 				importeComisionProveedor = cc.ChannelArancel.Importe                                                            //2
// 				resultImporte = c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)                                 // 100
// 			case "PORCENTAJE":
// 				// valores porcentual del proveedor
// 				// calcular comision general

// 				comisionGeneral := cc.Comision + cc.ChannelArancel.Importe
// 				importeComisiones = c.ToFixed((float64(movimiento.Monto) / 100 * comisionGeneral), 4)
// 				// calcular comision telco y proveedor
// 				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)
// 				importeComisionProveedor = c.ToFixed((float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe), 4)
// 				// controlar calculos
// 				resultImporte = c.ToFixed((float64(movimiento.Monto)/100*cc.Comision)+(float64(movimiento.Monto)/100*cc.ChannelArancel.Importe), 4) // c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)

// 			}
// 			if importeComisiones == resultImporte {

// 				var valorPorcentajeTelco float64
// 				var valorPorcentajeProveedor float64

// 				valorComisionTelcoVerif := VerificarMinimoMaximo(importeComisionTelco, cc.Importeminimo, cc.Importemaximo, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo)
// 				valorPorcentajeTelco = cc.Comision
// 				if importeComisionTelco != valorComisionTelcoVerif {
// 					if cc.ChannelArancel.Tipocalculo == "FIJO" {
// 						valorComisionTelcoVerif = valorComisionTelcoVerif - cc.ChannelArancel.Importe
// 					}
// 					importeComisionTelco = valorComisionTelcoVerif
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeTelco = valorComisionTelcoVerif
// 					} else {
// 						valorPorcentajeTelco = -1.00 * valorComisionTelcoVerif
// 					}
// 				}
// 				//verifico valor minimo y maximo del proveedor
// 				// cc.ChannelArancel.Importe es igual a la comision del proveedor
// 				// importeComisionProveedorVerif puede contener un importe o un coeficiente

// 				valorComisionProveedorVerif := VerificarMinimoMaximo(importeComisionProveedor, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo, 0, 0)
// 				valorPorcentajeProveedor = cc.ChannelArancel.Importe
// 				if importeComisionProveedor != valorComisionProveedorVerif {
// 					importeComisionProveedor = valorComisionProveedorVerif
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeProveedor = valorComisionProveedorVerif
// 					} else {
// 						valorPorcentajeProveedor = -1.00 * valorComisionProveedorVerif
// 					}
// 				}

// 				movimientoComision := entities.Movimientocomisiones{
// 					CuentacomisionsID:   cc.ID,
// 					Monto:               entities.Monto(int64(importeComisionTelco * 100)),
// 					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
// 					Montoproveedor:      entities.Monto(int64(importeComisionProveedor * 100)),
// 					Porcentajeproveedor: valorPorcentajeProveedor,
// 				}
// 				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
// 				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 				//Calculo de Iva sobre comision
// 				if iva != nil && iva.Porcentaje > -1 {
// 					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
// 					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
// 					movimientoImpuesto := entities.Movimientoimpuestos{
// 						ImpuestosID:    uint64(iva.ID),
// 						Monto:          entities.Monto(int64(impuestoIva * 100)),
// 						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
// 						Porcentaje:     iva.Porcentaje,
// 					}
// 					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
// 					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 				}
// 			}
// 		}
// 	}
// 	if importeSolicitado < movimiento.Monto {
// 		movimiento.Monto = importeSolicitado
// 	}
// 	movimiento.Monto -= descuentos

// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	// if movimiento.Monto > 0 {
// 	// 	movimiento.Monto -= descuentos
// 	// } else {
// 	// 	movimiento.Monto += descuentos
// 	// }

// 	return
// }

// func VerificarMinimoMaximo(comision, minimo, maximo float64) (importeComision float64) {

// 	if minimo != 0 && maximo == 0 {
// 		if math.Abs(comision) < minimo {
// 			importeComision = minimo
// 		} else {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	if minimo == 0 && maximo != 0 {
// 		if math.Abs(comision) > maximo {
// 			importeComision = maximo
// 		} else {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	if minimo != 0 && maximo != 0 {
// 		if math.Abs(comision) < minimo {
// 			importeComision = minimo
// 		}
// 		if math.Abs(comision) > maximo {
// 			importeComision = maximo
// 		}
// 		if math.Abs(comision) > minimo && math.Abs(comision) < maximo {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	importeComision = comision
// 	return
// }

// func VerificarMinimoMaximo(request RequestComision) (importeComisionTelco, importeComisionProveedor float64) {

// 	if request.MinBool && !request.MaxBool {
// 		if request.MinTelco != 0 && request.MinProveedor != 0 {
// 			importeComisionTelco = request.MinTelco
// 			importeComisionProveedor = request.MinProveedor
// 			return
// 		}

// 		if request.MinTelco != 0 && request.MinProveedor == 0 {
// 			importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
// 			importeComisionProveedor = request.ImporteComisionProveedor
// 			return
// 		}

// 		// if request.MinTelco == 0 && request.MinProveedor != 0 {
// 		// 	importeComisionTelco = request.ImporteComisionTelco
// 		// 	importeComisionProveedor = request.MinProveedor
// 		// 	return
// 		// }
// 	}

// 	// if request.MaxBool && !request.MinBool {
// 	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.MaxProveedor
// 	// 		return
// 	// 	}

// 	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// }

// 	// if !request.MinBool && !request.MaxBool {
// 	// 	if request.MinTelco != 0 && request.MinProveedor != 0 {
// 	// 		importeComisionTelco = request.MinTelco
// 	// 		importeComisionProveedor = request.MinProveedor
// 	// 		return
// 	// 	}
// 	// 	if request.MinTelco != 0 && request.MinProveedor == 0 {
// 	// 		importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.MaxProveedor
// 	// 		return
// 	// 	}

// 	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// 	importeComisionTelco = request.ImporteComisionTelco
// 	// 	importeComisionProveedor = request.ImporteComisionProveedor
// 	// 	return
// 	// }
// 	importeComisionTelco = request.ImporteComisionTelco
// 	importeComisionProveedor = request.ImporteComisionProveedor
// 	return
// }
