package administracion

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/commonsdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

func (s *service) GetInformacionSupervision(request ribcradtos.GetInformacionSupervisionRequest) (ri ribcradtos.RiInformacionSupervisionReponse, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	filtroCbu := filtros.ConfiguracionFiltro {
		Nombre: "CBU_CUENTA_TELCO",
	}

	config, erro := s.utilService.GetConfiguracionService(filtroCbu)

	if erro != nil {
		return
	}

	requestCuentasClientes := ribcradtos.RICuentasClienteRequest{
		FechaInicio: request.FechaInicio,
		FechaFin:    request.FechaFin,
		CbuCuentaTelco: config.Valor,
	}

	ri.RiCuentaCliente, erro = s.repository.BuildRICuentasCliente(requestCuentasClientes)

	if erro != nil {
		return
	}

	requestDatosFondos := ribcradtos.RiDatosFondosRequest{
		FechaInicio: request.FechaInicio,
		FechaFin:    request.FechaFin,
	}

	ri.RiDatosFondos, erro = s.repository.BuildRIDatosFondo(requestDatosFondos)

	if erro != nil {
		return
	}

	return
}

func (s *service) BuildInformacionSupervision(request ribcradtos.BuildInformacionSupervisionRequest) (ruta string, erro error) {

	// Carga RiCuentaCliente
	var requerimientoCuentaCliente ribcradtos.RIRequerimiento
	var rutaCuentaCliente string
	var nombreArchivoCuentas commonsdtos.FileName

	requestCuentaCliente := ribcradtos.CargarRiCuentaClienteRequest{
		Ri: request.RiCuentaCliente,
		CargarRiRequest: ribcradtos.CargarRiRequest{
			Rectificativa: request.RectificativaCuentaCliente,
			Opera:         request.OperaCuentaCliente,
			Periodo:       request.Periodo,
		},
	}

	requerimientoCuentaCliente, rutaCuentaCliente, nombreArchivoCuentas, erro = s.CargarRiCuentaCliente(requestCuentaCliente)

	if erro != nil {
		return
	}
	//Carga RiDatosFondos
	var requerimientoDatosFondos ribcradtos.RIRequerimiento
	var rutaDatosFondos string
	var nombreArchivoDatos commonsdtos.FileName

	requestDatosFondos := ribcradtos.CargarRiDatosFondosRequest{
		Ri: request.RiDatosFondos,
		CargarRiRequest: ribcradtos.CargarRiRequest{
			Rectificativa: request.RectificativaDatosFondo,
			Opera:         request.OperaDatosFondos,
			Periodo:       request.Periodo,
		},
	}

	requerimientoDatosFondos, rutaDatosFondos, nombreArchivoDatos, erro = s.CargarRiDatosFondos(requestDatosFondos)

	if erro != nil {
		s._eliminarArchivo(nil, nil, rutaCuentaCliente)
		return
	}
	//Carga RiInfEspecial
	var requerimientoInfEspecial ribcradtos.RIRequerimiento
	var rutaInfEspecial string
	var nombreArchivoInfEspecial commonsdtos.FileName

	requestInfEspecial := ribcradtos.CargarRiInfEspecialRequest{
		InfEspecial: request.InfEspecial,
		Fiber:       request.Fiber,
		CargarRiRequest: ribcradtos.CargarRiRequest{
			Rectificativa: request.RectificativaInfEspecial,
			Opera:         request.OperaInfEspecial,
			Periodo:       request.Periodo,
		},
	}

	requerimientoInfEspecial, rutaInfEspecial, nombreArchivoInfEspecial, erro = s.CargarRiInfEspecial(requestInfEspecial)

	if erro != nil {
		s._eliminarArchivo(nil, nil, rutaCuentaCliente)
		s._eliminarArchivo(nil, nil, rutaDatosFondos)
		return
	}

	presentacion := ribcradtos.RIPresentacion{
		Informacion: ribcradtos.RIInformacion{
			Tipo: "RI",
			Especificacion: ribcradtos.RIEspecificacion{
				Regimen: []ribcradtos.RIRegimen{
					{
						Codigo: "100",
						Requerimiento: []ribcradtos.RIRequerimiento{
							requerimientoCuentaCliente,
							requerimientoDatosFondos,
							requerimientoInfEspecial,
						},
					},
				},
			},
		},
	}

	nombreArchivoDetalle := commonsdtos.FileName{
		RutaBase:  config.RUTA_RI_BCRA,
		Nombre:    "detalle",
		Extension: "xml",
		UsaFecha:  false,
	}

	rutaDetalle := s.commonsService.CreateFileName(nombreArchivoDetalle)

	erro = s._buildXmlFile(rutaDetalle, presentacion)

	if erro != nil {
		//En caso de error elimino los archivos de cuentas de cliente y datos fondos
		s._eliminarArchivo(nil, nil, rutaCuentaCliente)
		s._eliminarArchivo(nil, nil, rutaDatosFondos)
		s._eliminarArchivo(nil, nil, rutaInfEspecial)
		return
	}

	nombreArchivoZip := commonsdtos.FileName{
		RutaBase:  config.RUTA_RI_BCRA,
		Nombre:    "InfSupervision",
		Extension: "zip",
		UsaFecha:  true,
	}

	ruta = s.commonsService.CreateFileName(nombreArchivoZip)

	requestZip := commonsdtos.ZipFilesRequest{
		NombreArchivo: ruta,
		Rutas:         []commonsdtos.InfoFile{},
	}

	if len(rutaCuentaCliente) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaCuentaCliente, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoCuentas.Nombre, nombreArchivoCuentas.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	if len(rutaDatosFondos) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaDatosFondos, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoDatos.Nombre, nombreArchivoDatos.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	if len(rutaInfEspecial) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaInfEspecial, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoInfEspecial.Nombre, nombreArchivoInfEspecial.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	if len(rutaDetalle) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaDetalle, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoDetalle.Nombre, nombreArchivoDetalle.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	erro = s.commonsService.ZipFiles(requestZip)

	if len(rutaCuentaCliente) > 0 {
		s._eliminarArchivo(nil, nil, rutaCuentaCliente)
	}
	if len(rutaDatosFondos) > 0 {
		s._eliminarArchivo(nil, nil, rutaDatosFondos)
	}
	if len(rutaInfEspecial) > 0 {
		s._eliminarArchivo(nil, nil, rutaInfEspecial)
	}
	if len(rutaDetalle) > 0 {
		s._eliminarArchivo(nil, nil, rutaDetalle)
	}

	return
}

func (s *service) GetInformacionEstadistica(request ribcradtos.GetInformacionEstadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	requestInfEstadistica := ribcradtos.RiInfestadisticaRequest{
		FechaInicio: request.FechaInicio,
		FechaFin:    request.FechaFin,
	}

	ri, erro = s.repository.BuilRIInfestaditica(requestInfEstadistica)

	if erro != nil {
		return
	}

	return
}

func (s *service) BuildInformacionEstadistica(request ribcradtos.BuildInformacionEstadisticaRequest) (ruta string, erro error) {

	// Carga RiCuentaCliente
	var requerimientoInfestadistica ribcradtos.RIRequerimiento
	var rutaInfestadistica string
	var nombreArchivoInfestadistica commonsdtos.FileName

	requestInfEstadistica := ribcradtos.CargarRiInfEstadisticaRequest{
		Ri: request.RiInfestadistica,
		CargarRiRequest: ribcradtos.CargarRiRequest{
			Rectificativa: request.Rectificativa,
			Opera:         request.Opera,
			Periodo:       request.Periodo,
		},
	}

	requerimientoInfestadistica, rutaInfestadistica, nombreArchivoInfestadistica, erro = s.CargarRiInfEstadisitica(requestInfEstadistica)

	if erro != nil {
		return
	}

	presentacion := ribcradtos.RIPresentacion{
		Informacion: ribcradtos.RIInformacion{
			Tipo: "RI",
			Especificacion: ribcradtos.RIEspecificacion{
				Regimen: []ribcradtos.RIRegimen{
					{
						Codigo: "100",
						Requerimiento: []ribcradtos.RIRequerimiento{
							requerimientoInfestadistica,
						},
					},
				},
			},
		},
	}

	nombreArchivoDetalle := commonsdtos.FileName{
		RutaBase:  config.RUTA_RI_BCRA,
		Nombre:    "detalle",
		Extension: "xml",
		UsaFecha:  false,
	}

	rutaDetalle := s.commonsService.CreateFileName(nombreArchivoDetalle)

	erro = s._buildXmlFile(rutaDetalle, presentacion)

	if erro != nil {
		//En caso de error elimino los archivos de cuentas de cliente y datos fondos
		s._eliminarArchivo(nil, nil, rutaInfestadistica)
		return
	}

	nombreArchivoZip := commonsdtos.FileName{
		RutaBase:  config.RUTA_RI_BCRA,
		Nombre:    "InfEstadistica",
		Extension: "zip",
		UsaFecha:  true,
	}

	ruta = s.commonsService.CreateFileName(nombreArchivoZip)

	requestZip := commonsdtos.ZipFilesRequest{
		NombreArchivo: ruta,
		Rutas:         []commonsdtos.InfoFile{},
	}

	if len(rutaInfestadistica) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaInfestadistica, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoInfestadistica.Nombre, nombreArchivoInfestadistica.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	if len(rutaDetalle) > 0 {
		infoFile := commonsdtos.InfoFile{
			RutaCompleta: rutaDetalle, NombreArchivo: fmt.Sprintf("%s.%s", nombreArchivoDetalle.Nombre, nombreArchivoDetalle.Extension),
		}
		requestZip.Rutas = append(requestZip.Rutas, infoFile)
	}

	erro = s.commonsService.ZipFiles(requestZip)

	if len(rutaInfestadistica) > 0 {
		s._eliminarArchivo(nil, nil, rutaInfestadistica)
	}

	if len(rutaDetalle) > 0 {
		s._eliminarArchivo(nil, nil, rutaDetalle)
	}

	return
}

func (s *service) CargarRiInfEstadisitica(request ribcradtos.CargarRiInfEstadisticaRequest) (requerimiento ribcradtos.RIRequerimiento, ruta string, fileName commonsdtos.FileName, erro error) {

	// Valido que el ri Cuenta Cliente tiene datos validos
	for i := range request.Ri {
		erro = request.Ri[i].IsValid()
		if erro != nil {
			return
		}
	}

	if request.Opera {

		// Creo el archivo InfEstadistica.txt
		fileName = commonsdtos.FileName{
			RutaBase:  config.RUTA_RI_BCRA,
			Nombre:    "INFESTADISTICA",
			Extension: "txt",
			UsaFecha:  false,
		}

		ruta = s.commonsService.CreateFileName(fileName)

		guardarRequest := ribcradtos.RIGuardarArchivosRequest{
			Ruta: ruta, RI: request.Ri,
		}

		// Guardo el archivo InfEstadistica.txt
		erro = s.RIGuardarArchivos(guardarRequest)

		if erro != nil {
			return
		}

	}

	tipoPresentacion := ribcradtos.Normal
	if request.Rectificativa {
		tipoPresentacion = ribcradtos.Rectificativa
	}

	detalle := ribcradtos.RIDetalle{
		Opera:   request.Opera,
		Tipo:    string(tipoPresentacion),
		Periodo: request.Periodo,
		Archivo: []ribcradtos.RIArchivo{},
	}

	if detalle.Opera {
		detalle.Archivo = append(detalle.Archivo, ribcradtos.RIArchivo{Ruta: fmt.Sprintf("%s.%s", fileName.Nombre, fileName.Extension)})
	}

	requerimiento = ribcradtos.RIRequerimiento{
		Codigo:  "3",
		Detalle: detalle,
	}

	return
}

func (s *service) CargarRiCuentaCliente(request ribcradtos.CargarRiCuentaClienteRequest) (requerimiento ribcradtos.RIRequerimiento, ruta string, fileName commonsdtos.FileName, erro error) {

	// Valido que el ri Cuenta Cliente tiene datos validos
	for i := range request.Ri {
		erro = request.Ri[i].IsValid()
		if erro != nil {
			return
		}
	}

	if request.Opera {

		// Creo el archivo CuentasClientes.txt
		fileName = commonsdtos.FileName{
			RutaBase:  config.RUTA_RI_BCRA,
			Nombre:    "CUENTASCLIENTES",
			Extension: "txt",
			UsaFecha:  false,
		}

		ruta = s.commonsService.CreateFileName(fileName)

		guardarRequest := ribcradtos.RIGuardarArchivosRequest{
			Ruta: ruta, RI: request.Ri,
		}

		// Guardo el archivo CuentasClientes.txt
		erro = s.RIGuardarArchivos(guardarRequest)

		if erro != nil {
			return
		}

	}

	tipoPresentacionCuentas := ribcradtos.Normal
	if request.Rectificativa {
		tipoPresentacionCuentas = ribcradtos.Rectificativa
	}

	detalleCuentas := ribcradtos.RIDetalle{
		Opera:   request.Opera,
		Tipo:    string(tipoPresentacionCuentas),
		Periodo: request.Periodo,
		Archivo: []ribcradtos.RIArchivo{},
	}

	if detalleCuentas.Opera {
		detalleCuentas.Archivo = append(detalleCuentas.Archivo, ribcradtos.RIArchivo{Ruta: fmt.Sprintf("%s.%s", fileName.Nombre, fileName.Extension)})
	}

	requerimiento = ribcradtos.RIRequerimiento{
		Codigo:  "1",
		Detalle: detalleCuentas,
	}

	return
}

func (s *service) CargarRiDatosFondos(request ribcradtos.CargarRiDatosFondosRequest) (requerimiento ribcradtos.RIRequerimiento, ruta string, fileName commonsdtos.FileName, erro error) {

	// valido el ri de Datos Fondos
	for i := range request.Ri {
		erro = request.Ri[i].IsValid()
		if erro != nil {
			return
		}
		request.Ri[i].Denominacion, erro = s.commonsService.NormalizeStrings(request.Ri[i].Denominacion)
		if erro != nil {
			return
		}
		request.Ri[i].DenominacionAgente, erro = s.commonsService.NormalizeStrings(request.Ri[i].DenominacionAgente)
		if erro != nil {
			return
		}

	}

	if request.Opera {
		// creo el archivo DatosFondos.txt
		fileName = commonsdtos.FileName{
			RutaBase:  config.RUTA_RI_BCRA,
			Nombre:    "DATOSFONDOS",
			Extension: "txt",
			UsaFecha:  false,
		}

		ruta = s.commonsService.CreateFileName(fileName)

		guardarRequest := ribcradtos.RIGuardarArchivosRequest{
			Ruta: ruta, RI: request.Ri,
		}
		// guardo el archivo de DatosFondos.txt
		erro = s.RIGuardarArchivos(guardarRequest)

		if erro != nil {
			return
		}

	}

	tipoPresentacionDatosFondos := ribcradtos.Normal
	if request.Rectificativa {
		tipoPresentacionDatosFondos = ribcradtos.Rectificativa
	}

	detalleDatosFondos := ribcradtos.RIDetalle{
		Opera:   request.Opera,
		Tipo:    string(tipoPresentacionDatosFondos),
		Periodo: request.Periodo,
		Archivo: []ribcradtos.RIArchivo{},
	}

	if detalleDatosFondos.Opera {
		detalleDatosFondos.Archivo = append(detalleDatosFondos.Archivo, ribcradtos.RIArchivo{Ruta: fmt.Sprintf("%s.%s", fileName.Nombre, fileName.Extension)})
	}

	requerimiento = ribcradtos.RIRequerimiento{
		Codigo:  "1",
		Detalle: detalleDatosFondos,
	}

	return
}

func (s *service) CargarRiInfEspecial(request ribcradtos.CargarRiInfEspecialRequest) (requerimiento ribcradtos.RIRequerimiento, ruta string, fileName commonsdtos.FileName, erro error) {

	if request.Opera {
		fileName = commonsdtos.FileName{
			RutaBase:  config.RUTA_RI_BCRA,
			Nombre:    "INFESPECIAL",
			Extension: "pdf",
			UsaFecha:  false,
		}

		//En cada trimestre se guarda un archivo pdf en este caso hay que guardarlo
		periodo := string(request.Periodo[5:7])
		switch periodo {
		case "01", "04", "07", "10":
			if request.InfEspecial != nil {
				// creo el archivo InfEspecial.txt
				ruta = s.commonsService.CreateFileName(fileName)

				erro = s.commonsService.SaveFiberPdf(request.InfEspecial, ruta, request.Fiber)
				if erro != nil {
					return
				}
			}

		}

	}

	tipoPresentacionInfEspecial := ribcradtos.Normal
	if request.Rectificativa {
		tipoPresentacionInfEspecial = ribcradtos.Rectificativa
	}

	detalleInfEspecial := ribcradtos.RIDetalle{
		Opera:   request.Opera,
		Tipo:    string(tipoPresentacionInfEspecial),
		Periodo: request.Periodo,
		Archivo: []ribcradtos.RIArchivo{},
	}
	if detalleInfEspecial.Opera {
		detalleInfEspecial.Archivo = append(detalleInfEspecial.Archivo, ribcradtos.RIArchivo{Ruta: fmt.Sprintf("%s.%s", fileName.Nombre, fileName.Extension)})
	}

	requerimiento = ribcradtos.RIRequerimiento{
		Codigo:  "2",
		Detalle: detalleInfEspecial,
	}

	return
}

func (s *service) RIInfestadistica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {
	erro = request.IsValid()

	if erro != nil {
		return
	}

	ri, erro = s.repository.BuilRIInfestaditica(request)

	if erro != nil {
		return
	}
	if len(ri) > 0 {

		for i := range ri {
			erro = ri[i].IsValid()
			if erro != nil {
				return
			}
		}

		guardarRequest := ribcradtos.RIGuardarArchivosRequest{
			Ruta: request.Ruta, RI: ri,
		}
		erro = s.RIGuardarArchivos(guardarRequest)

	}

	return
}

func (s *service) _eliminarArchivo(windows *transform.Writer, archivo *os.File, ruta string) {
	if windows != nil {
		windows.Close()
	}
	if archivo != nil {
		archivo.Close()
	}
	erro := s.commonsService.RemoveFile(ruta)
	if erro != nil {
		logs.Error(erro.Error())
	}
}

func (s *service) _buildXmlFile(ruta string, data interface{}) (erro error) {

	archivo, erro := s.commonsService.CreateFile(ruta)

	if erro != nil {
		return
	}

	defer archivo.Close()

	archivo.WriteString("<?xml version=\"1.0\"?>\n")
	encoder := xml.NewEncoder(archivo)
	encoder.Indent("", "\t")
	erro = encoder.Encode(&data)

	if erro != nil {
		s._eliminarArchivo(nil, archivo, ruta)
		return
	}

	return
}

func (s *service) RIGuardarArchivos(request ribcradtos.RIGuardarArchivosRequest) (erro error) {
	erro = request.IsValid()
	if erro != nil {
		return
	}
	archivo, erro := s.commonsService.CreateFile(request.Ruta)

	if erro != nil {
		return
	}

	defer archivo.Close()

	// Se debe transformar en el formato ANSI 1252
	buffwriter := bufio.NewWriter(archivo)
	enc, _ := htmlindex.Get("windows-1252")

	encoder := enc.NewEncoder()
	windows1252writer := transform.NewWriter(buffwriter, encoder)
	defer windows1252writer.Close()


	switch v := request.RI.(type) { 
	case []ribcradtos.RiCuentaCliente:
		for i := range v {
			_, erro = fmt.Fprintln(windows1252writer, v[i].ToString())
			if erro != nil {
				s._eliminarArchivo(windows1252writer, archivo, request.Ruta)
				return
			}

		}
	case []ribcradtos.RiDatosFondos:
		for i := range v {
			_, erro = fmt.Fprintln(windows1252writer, v[i].ToString())
			if erro != nil {
				s._eliminarArchivo(windows1252writer, archivo, request.Ruta)
				return
			}

		}
	case []ribcradtos.RiInfestadistica:
		for i := range v {
			_, erro = fmt.Fprintln(windows1252writer, v[i].ToString())
			if erro != nil {
				s._eliminarArchivo(windows1252writer, archivo, request.Ruta)
				return
			}

		}
	default:
		erro = fmt.Errorf("el archivo seleccionado no se pudo guardar porque no tiene un implementación definida")
		s._eliminarArchivo(windows1252writer, archivo, request.Ruta)
		return

	}

	// Aqui recién escribo en el archivo ya con el formato windows-1252
	erro = buffwriter.Flush()
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "RIGuardarArchivos",
			Mensaje:       erro.Error(),
		}
		s.utilService.CreateLogService(log)
		erro = fmt.Errorf(ERROR_ESCRIBIR_ARCHIVO)
		s._eliminarArchivo(windows1252writer, archivo, request.Ruta)
	}

	return
}
