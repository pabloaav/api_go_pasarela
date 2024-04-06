package cierrelotedtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type PrismaPxDosRegistro struct {
	Eclq02llEmpresa     string
	Eclq02llFpres       string
	Eclq02llTiporeg     string
	Eclq02llMoneda      string
	Eclq02llNumcom      string
	Eclq02llNumest      string
	Eclq02llNroliq      string
	Eclq02llFpag        string
	Eclq02llTipoliq     string
	Eclq02llImpbruto    string
	Eclq02llSigno_1     string
	Eclq02llImpret      string
	Eclq02llSigno_2     string
	Eclq02llImpneto     string
	Eclq02llSigno_3     string
	Eclq02llRetesp      string
	Eclq02llSigno_4     string
	Eclq02llRetivaEsp   string
	Eclq02llSigno_5     string
	Eclq02llPercepBa    string
	Eclq02llSigno_6     string
	Eclq02llRetivaD1    string
	Eclq02llSigno_7     string
	Filler1             string
	Filler2             string
	Eclq02llCargoPex    string
	Eclq02llSigno_9     string
	Eclq02llRetivaPex1  string
	Eclq02llSigno_10    string
	Filler3             string
	Filler4             string
	Eclq02llCostoCuoemi string
	Eclq02llSigno_12    string
	Eclq02llRetivaCuo1  string
	Eclq02llSigno_13    string
	Filler5             string
	Filler6             string
	Eclq02llImpServ     string
	Eclq02llSigno_15    string
	Eclq02llIva1Xlj     string
	Eclq02llSigno_16    string
	Filler7             string
	Filler8             string
	Eclq02llCargoEdcE   string
	Eclq02llSigno_18    string
	Eclq02llIva1EdcE    string
	Eclq02llSigno_19    string
	Filler9             string
	Filler10            string
	Eclq02llCargoEdcB   string
	Eclq02llSigno_21    string
	Eclq02llIva1EdcB    string
	Eclq02llSigno_22    string
	Filler11            string
	Filler12            string
	Eclq02llCargoCitE   string
	Eclq02llSigno_24    string
	Eclq02llIva1CitE    string
	Eclq02llllSigno_25  string
	Filler13            string
	Filler14            string
	Eclq02llCargoCitB   string
	Eclq02llSigno_27    string
	Eclq02llIva1CitB    string
	Eclq02llSigno_28    string
	Filler15            string
	Filler16            string
	Eclq02llRetIva      string
	Eclq02llSigno_30    string
	Eclq02llRetGcias    string
	Eclq02llSigno_31    string
	Eclq02llRetIngbru   string
	Eclq02llSigno_32    string
	Filler17            string
	Filler18            string
	Filler19            string
	Eclq02llAster       string
}

func (pxDos *PrismaPxDosRegistro) ValidarPxDos(estructuraReg *EstructuraRegistros) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro tipo (02)")

	if len(pxDos.Eclq02llEmpresa) != estructuraReg.PxDosDescripcionRegistro()[0].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llFpres) != estructuraReg.PxDosDescripcionRegistro()[1].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llTiporeg) != estructuraReg.PxDosDescripcionRegistro()[2].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llMoneda) != estructuraReg.PxDosDescripcionRegistro()[3].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llNumcom) != estructuraReg.PxDosDescripcionRegistro()[4].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llNumest) != estructuraReg.PxDosDescripcionRegistro()[5].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llNroliq) != estructuraReg.PxDosDescripcionRegistro()[6].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llFpag) != estructuraReg.PxDosDescripcionRegistro()[7].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llTipoliq) != estructuraReg.PxDosDescripcionRegistro()[8].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llImpbruto) != estructuraReg.PxDosDescripcionRegistro()[9].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_1) != estructuraReg.PxDosDescripcionRegistro()[10].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llImpret) != estructuraReg.PxDosDescripcionRegistro()[11].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_2) != estructuraReg.PxDosDescripcionRegistro()[12].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llImpneto) != estructuraReg.PxDosDescripcionRegistro()[13].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_3) != estructuraReg.PxDosDescripcionRegistro()[14].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetesp) != estructuraReg.PxDosDescripcionRegistro()[15].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_4) != estructuraReg.PxDosDescripcionRegistro()[16].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetivaEsp) != estructuraReg.PxDosDescripcionRegistro()[17].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_5) != estructuraReg.PxDosDescripcionRegistro()[18].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llPercepBa) != estructuraReg.PxDosDescripcionRegistro()[19].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_6) != estructuraReg.PxDosDescripcionRegistro()[20].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetivaD1) != estructuraReg.PxDosDescripcionRegistro()[21].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_7) != estructuraReg.PxDosDescripcionRegistro()[22].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler1) != estructuraReg.PxDosDescripcionRegistro()[23].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler2) != estructuraReg.PxDosDescripcionRegistro()[24].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCargoPex) != estructuraReg.PxDosDescripcionRegistro()[25].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_9) != estructuraReg.PxDosDescripcionRegistro()[26].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetivaPex1) != estructuraReg.PxDosDescripcionRegistro()[27].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_10) != estructuraReg.PxDosDescripcionRegistro()[28].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler3) != estructuraReg.PxDosDescripcionRegistro()[29].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler4) != estructuraReg.PxDosDescripcionRegistro()[30].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCostoCuoemi) != estructuraReg.PxDosDescripcionRegistro()[31].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_12) != estructuraReg.PxDosDescripcionRegistro()[32].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetivaCuo1) != estructuraReg.PxDosDescripcionRegistro()[33].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_13) != estructuraReg.PxDosDescripcionRegistro()[34].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler5) != estructuraReg.PxDosDescripcionRegistro()[35].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler6) != estructuraReg.PxDosDescripcionRegistro()[36].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llImpServ) != estructuraReg.PxDosDescripcionRegistro()[37].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_15) != estructuraReg.PxDosDescripcionRegistro()[38].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llIva1Xlj) != estructuraReg.PxDosDescripcionRegistro()[39].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_16) != estructuraReg.PxDosDescripcionRegistro()[40].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler7) != estructuraReg.PxDosDescripcionRegistro()[41].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler8) != estructuraReg.PxDosDescripcionRegistro()[42].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCargoEdcE) != estructuraReg.PxDosDescripcionRegistro()[43].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_18) != estructuraReg.PxDosDescripcionRegistro()[44].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llIva1EdcE) != estructuraReg.PxDosDescripcionRegistro()[45].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_19) != estructuraReg.PxDosDescripcionRegistro()[46].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler9) != estructuraReg.PxDosDescripcionRegistro()[47].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler10) != estructuraReg.PxDosDescripcionRegistro()[48].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCargoEdcB) != estructuraReg.PxDosDescripcionRegistro()[49].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_21) != estructuraReg.PxDosDescripcionRegistro()[50].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llIva1EdcB) != estructuraReg.PxDosDescripcionRegistro()[51].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_22) != estructuraReg.PxDosDescripcionRegistro()[52].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler11) != estructuraReg.PxDosDescripcionRegistro()[53].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler12) != estructuraReg.PxDosDescripcionRegistro()[54].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCargoCitE) != estructuraReg.PxDosDescripcionRegistro()[55].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_24) != estructuraReg.PxDosDescripcionRegistro()[56].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llIva1CitE) != estructuraReg.PxDosDescripcionRegistro()[57].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llllSigno_25) != estructuraReg.PxDosDescripcionRegistro()[58].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler13) != estructuraReg.PxDosDescripcionRegistro()[59].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler14) != estructuraReg.PxDosDescripcionRegistro()[60].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llCargoCitB) != estructuraReg.PxDosDescripcionRegistro()[61].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_27) != estructuraReg.PxDosDescripcionRegistro()[62].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llIva1CitB) != estructuraReg.PxDosDescripcionRegistro()[63].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_28) != estructuraReg.PxDosDescripcionRegistro()[64].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler15) != estructuraReg.PxDosDescripcionRegistro()[65].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler16) != estructuraReg.PxDosDescripcionRegistro()[66].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetIva) != estructuraReg.PxDosDescripcionRegistro()[67].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_30) != estructuraReg.PxDosDescripcionRegistro()[68].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetGcias) != estructuraReg.PxDosDescripcionRegistro()[69].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_31) != estructuraReg.PxDosDescripcionRegistro()[70].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llRetIngbru) != estructuraReg.PxDosDescripcionRegistro()[71].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llSigno_32) != estructuraReg.PxDosDescripcionRegistro()[72].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler17) != estructuraReg.PxDosDescripcionRegistro()[73].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler18) != estructuraReg.PxDosDescripcionRegistro()[74].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Filler19) != estructuraReg.PxDosDescripcionRegistro()[75].Cantidad {
		erro = err
		return
	}
	if len(pxDos.Eclq02llAster) != estructuraReg.PxDosDescripcionRegistro()[76].Cantidad {
		erro = err
		return
	}
	return nil
}

func (pxDos *PrismaPxDosRegistro) PxDosToEntities() (pxDosEntity entities.Prismapxdosregistro) {

	pxDosEntity.Eclq02llEmpresa = pxDos.Eclq02llEmpresa
	pxDosEntity.Eclq02llFpres = pxDos.Eclq02llFpres
	pxDosEntity.Eclq02llTiporeg = pxDos.Eclq02llTiporeg
	pxDosEntity.Eclq02llMoneda = pxDos.Eclq02llMoneda
	pxDosEntity.Eclq02llNumcom = pxDos.Eclq02llNumcom
	pxDosEntity.Eclq02llNumest = pxDos.Eclq02llNumest
	pxDosEntity.Eclq02llNroliq = pxDos.Eclq02llNroliq
	pxDosEntity.Eclq02llFpag = pxDos.Eclq02llFpag
	pxDosEntity.Eclq02llTipoliq = pxDos.Eclq02llTipoliq
	pxDosEntity.Eclq02llImpbruto = pxDos.Eclq02llImpbruto
	pxDosEntity.Eclq02llSigno_1 = pxDos.Eclq02llSigno_1
	pxDosEntity.Eclq02llImppret = pxDos.Eclq02llImpret
	pxDosEntity.Eclq02llSigno_2 = pxDos.Eclq02llSigno_2
	pxDosEntity.Eclq02llImpneto = pxDos.Eclq02llImpneto
	pxDosEntity.Eclq02llSigno_3 = pxDos.Eclq02llSigno_3
	pxDosEntity.Eclq02llRetesp = pxDos.Eclq02llRetesp
	pxDosEntity.Eclq02llSigno_4 = pxDos.Eclq02llSigno_4
	pxDosEntity.Eclq02llRetivaEsp = pxDos.Eclq02llRetivaEsp
	pxDosEntity.Eclq02llSigno_5 = pxDos.Eclq02llSigno_5
	pxDosEntity.Eclq02llPercepBa = pxDos.Eclq02llPercepBa
	pxDosEntity.Eclq02llSigno_6 = pxDos.Eclq02llSigno_6
	pxDosEntity.Eclq02llRetivaD1 = pxDos.Eclq02llRetivaD1
	pxDosEntity.Eclq02llSigno_7 = pxDos.Eclq02llSigno_7
	pxDosEntity.Filler1 = pxDos.Filler1
	pxDosEntity.Filler2 = pxDos.Filler2
	pxDosEntity.Eclq02llCargoPex = pxDos.Eclq02llCargoPex
	pxDosEntity.Eclq02llSigno_9 = pxDos.Eclq02llSigno_9
	pxDosEntity.Eclq02llRetivaPex1 = pxDos.Eclq02llRetivaPex1
	pxDosEntity.Eclq02llSigno_10 = pxDos.Eclq02llSigno_10
	pxDosEntity.Filler3 = pxDos.Filler3
	pxDosEntity.Filler4 = pxDos.Filler4
	pxDosEntity.Eclq02llCostoCuoemi = pxDos.Eclq02llCostoCuoemi
	pxDosEntity.Eclq02llSigno_12 = pxDos.Eclq02llSigno_12
	pxDosEntity.Eclq02llRetivaCuo1 = pxDos.Eclq02llRetivaCuo1
	pxDosEntity.Eclq02llSigno_13 = pxDos.Eclq02llSigno_13
	pxDosEntity.Filler5 = pxDos.Filler5
	pxDosEntity.Filler6 = pxDos.Filler6
	pxDosEntity.Eclq02llImpServ = pxDos.Eclq02llImpServ
	pxDosEntity.Eclq02llSigno_15 = pxDos.Eclq02llSigno_15
	pxDosEntity.Eclq02llIva1Xlj = pxDos.Eclq02llIva1Xlj
	pxDosEntity.Eclq02llSigno_16 = pxDos.Eclq02llSigno_16
	pxDosEntity.Filler7 = pxDos.Filler7
	pxDosEntity.Filler8 = pxDos.Filler8
	pxDosEntity.Eclq02llCargoEdcE = pxDos.Eclq02llCargoEdcE
	pxDosEntity.Eclq02llSigno_18 = pxDos.Eclq02llSigno_18
	pxDosEntity.Eclq02llIva1EdcE = pxDos.Eclq02llIva1EdcE
	pxDosEntity.Eclq02llSigno_19 = pxDos.Eclq02llSigno_19
	pxDosEntity.Filler9 = pxDos.Filler9
	pxDosEntity.Filler10 = pxDos.Filler10
	pxDosEntity.Eclq02llCargoEdcB = pxDos.Eclq02llCargoEdcB
	pxDosEntity.Eclq02llSigno_21 = pxDos.Eclq02llSigno_21
	pxDosEntity.Eclq02llIva1EdcB = pxDos.Eclq02llIva1EdcB
	pxDosEntity.Eclq02llSigno_22 = pxDos.Eclq02llSigno_22
	pxDosEntity.Filler11 = pxDos.Filler11
	pxDosEntity.Filler12 = pxDos.Filler12
	pxDosEntity.Eclq02llCargoCitE = pxDos.Eclq02llCargoCitE
	pxDosEntity.Eclq02llSigno_24 = pxDos.Eclq02llSigno_24
	pxDosEntity.Eclq02llIva1CitE = pxDos.Eclq02llIva1CitE
	pxDosEntity.Eclq02llSigno_25 = pxDos.Eclq02llllSigno_25
	pxDosEntity.Filler13 = pxDos.Filler13
	pxDosEntity.Filler14 = pxDos.Filler14
	pxDosEntity.Eclq02llCargoCitB = pxDos.Eclq02llCargoCitB
	pxDosEntity.Eclq02llSigno_27 = pxDos.Eclq02llSigno_27
	pxDosEntity.Eclq02llIva1CitB = pxDos.Eclq02llIva1CitB
	pxDosEntity.Eclq02llSigno_28 = pxDos.Eclq02llSigno_28
	pxDosEntity.Filler15 = pxDos.Filler15
	pxDosEntity.Filler16 = pxDos.Filler16
	pxDosEntity.Eclq02llRetIva = pxDos.Eclq02llRetIva
	pxDosEntity.Eclq02llSigno_30 = pxDos.Eclq02llSigno_30
	pxDosEntity.Eclq02llRetGcias = pxDos.Eclq02llRetGcias
	pxDosEntity.Eclq02llSigno_31 = pxDos.Eclq02llSigno_31
	pxDosEntity.Eclq02llRetIngbru = pxDos.Eclq02llRetIngbru
	pxDosEntity.Eclq02llSigno_32 = pxDos.Eclq02llSigno_32
	pxDosEntity.Filler17 = pxDos.Filler17
	pxDosEntity.Filler18 = pxDos.Filler18
	pxDosEntity.Filler19 = pxDos.Filler19
	pxDosEntity.Eclq02llAster = pxDos.Eclq02llAster
	return
}
