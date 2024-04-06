package cierrelotedtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type MovimientoMxTotalesRegistro struct {
	Empresa           string
	Fechapres         string
	Tiporeg           string
	Numcom            string
	Numest            string
	Codop             string
	Tipoaplic         string
	Filler            string
	Fechapago         string
	Libre             string
	ImporteTotal      string
	SignoImporteTotal string
	Filler1           string
	McaPex            string
	Filler2           string
	Aster             string
}

func (mxTotales *MovimientoMxTotalesRegistro) ValidarMxTotales(estructuraReg *EstructuraRegistros) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro de totales")
	if len(mxTotales.Empresa) != estructuraReg.MxTotalesDescripcionRegistro()[0].Cantidad {
		return err
	}
	if len(mxTotales.Fechapres) != estructuraReg.MxTotalesDescripcionRegistro()[1].Cantidad {
		return err
	}
	if len(mxTotales.Tiporeg) != estructuraReg.MxTotalesDescripcionRegistro()[2].Cantidad {
		return err
	}
	if len(mxTotales.Numcom) != estructuraReg.MxTotalesDescripcionRegistro()[3].Cantidad {
		return err
	}
	if len(mxTotales.Numest) != estructuraReg.MxTotalesDescripcionRegistro()[4].Cantidad {
		return err
	}
	if len(mxTotales.Codop) != estructuraReg.MxTotalesDescripcionRegistro()[5].Cantidad {
		return err
	}
	if len(mxTotales.Tipoaplic) != estructuraReg.MxTotalesDescripcionRegistro()[6].Cantidad {
		return err
	}
	if len(mxTotales.Filler) != estructuraReg.MxTotalesDescripcionRegistro()[7].Cantidad {
		return err
	}
	if len(mxTotales.Fechapago) != estructuraReg.MxTotalesDescripcionRegistro()[8].Cantidad {
		return err
	}
	if len(mxTotales.Libre) != estructuraReg.MxTotalesDescripcionRegistro()[9].Cantidad {
		return err
	}
	if len(mxTotales.ImporteTotal) != estructuraReg.MxTotalesDescripcionRegistro()[10].Cantidad {
		return err
	}
	if len(mxTotales.SignoImporteTotal) != estructuraReg.MxTotalesDescripcionRegistro()[11].Cantidad {
		return err
	}
	if len(mxTotales.Filler1) != estructuraReg.MxTotalesDescripcionRegistro()[12].Cantidad {
		return err
	}
	if len(mxTotales.McaPex) != estructuraReg.MxTotalesDescripcionRegistro()[13].Cantidad {
		return err
	}
	if len(mxTotales.Filler2) != estructuraReg.MxTotalesDescripcionRegistro()[14].Cantidad {
		return err
	}
	if len(mxTotales.Aster) != estructuraReg.MxTotalesDescripcionRegistro()[15].Cantidad {
		return err
	}
	return
}

func (mxTotales *MovimientoMxTotalesRegistro) MxTotalesToEntities(nombreArchivo string) (mxTotalesEntity entities.Prismamxtotalesmovimiento) {
	mxTotalesEntity.Empresa = mxTotales.Empresa
	mxTotalesEntity.Fechapres = mxTotales.Fechapres
	mxTotalesEntity.Tiporeg = mxTotales.Tiporeg
	mxTotalesEntity.Numcom = mxTotales.Numcom
	mxTotalesEntity.Numest = mxTotales.Numest
	mxTotalesEntity.Codop = mxTotales.Codop
	mxTotalesEntity.Tipoaplic = mxTotales.Tipoaplic
	mxTotalesEntity.Filler = mxTotales.Filler
	mxTotalesEntity.Fechapago = mxTotales.Fechapago
	mxTotalesEntity.Libre = mxTotales.Libre
	mxTotalesEntity.ImporteTotal = mxTotales.ImporteTotal
	mxTotalesEntity.SignoImporteTotal = mxTotales.SignoImporteTotal
	mxTotalesEntity.Filler1 = mxTotales.Filler1
	mxTotalesEntity.McaPex = mxTotales.McaPex
	mxTotalesEntity.Filler2 = mxTotales.Filler2
	mxTotalesEntity.Aster = mxTotales.Aster
	mxTotalesEntity.Nombrearchivo = nombreArchivo
	return
}
