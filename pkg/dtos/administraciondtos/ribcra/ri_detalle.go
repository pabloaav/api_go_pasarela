package ribcradtos


import "encoding/xml"

//Se usa para generar el archivo de detalle para RI BCRA

type RIPresentacion struct {
	XMLName     xml.Name `xml:"PRESENTACION"`
	Informacion RIInformacion
}

type RIInformacion struct {
	XMLName        xml.Name `xml:"INFORMACION"`
	Tipo           string   `xml:"tipo,attr"`
	Especificacion RIEspecificacion
}

type RIEspecificacion struct {
	XMLName xml.Name `xml:"ESPECIFICACION"`
	Regimen []RIRegimen
}

type RIRegimen struct {
	XMLName       xml.Name `xml:"REGIMEN"`
	Codigo        string   `xml:"codigo,attr"`
	Requerimiento []RIRequerimiento
}

type RIRequerimiento struct {
	XMLName xml.Name `xml:"REQUERIMIENTO"`
	Codigo  string   `xml:"codigo,attr"`
	Detalle RIDetalle
}
type RIDetalle struct {
	XMLName xml.Name    `xml:"DETALLE"`
	Opera   bool        `xml:"opera,attr"`
	Tipo    string      `xml:"tipo,attr"`
	Periodo string      `xml:"periodo,attr"`
	Archivo []RIArchivo `xml:"RIArchivo"`
}

type RIArchivo struct {
	XMLName xml.Name `xml:"ARCHIVO"`
	Ruta    string   `xml:"ruta,attr"`
}
