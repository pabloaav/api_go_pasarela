package cierrelotedtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type PrismaPxCuatroRegistro struct {
	Eclq02llEmpresa_04         string
	Eclq02llFpres_04           string
	Eclq02llTiporeg_04         string
	Eclq02llMoneda_04          string
	Eclq02llNumcom_04          string
	Eclq02llNumest_04          string
	Eclq02llNroliq_04          string
	Eclq02llFpag_04            string
	Eclq02llTipoliq_04         string
	Eclq02llCasacta            string
	Eclq02llTipcta             string
	Eclq02llCtabco             string
	Eclq02llCfExentoIva        string
	Eclq02llSigno_04_1         string
	Eclq02llLey25063           string
	Eclq02llSigno_04_2         string
	Eclq02llAliIngbru          string
	Eclq02llDtoCampania        string
	Eclq02llSigno_04_3         string
	Eclq02llIva1DtoCampania    string
	Eclq02llSigno_04_4         string
	Eclq02llRetIngbru2         string
	Eclq02llSigno_04_5         string
	Eclq02llAliIngbru2         string
	Filler1                    string
	Filler2                    string
	Filler3                    string
	Filler4                    string
	Filler5                    string
	Eclq02llTasaPex            string
	Eclq02llCargoXliq          string
	Eclq02llSigno_04_8         string
	Eclq02llIva1CargoXliq      string
	Eclq02llSigno_04_9         string
	Eclq02llDealer             string
	Eclq02llImpDbCr            string
	Eclq02llSigno_04_10        string
	Eclq02llCfNoReduceIva      string
	Eclq02llSigno_04_11        string
	Eclq02llPercepIbAgip       string
	Eclq02llSigno_04_12        string
	Eclq02llAlicPercepIbAgip   string
	Eclq02llRetenIbAgip        string
	Eclq02llSigno_04_13        string
	Eclq02llAlicRetenIbAgip    string
	Eclq02llSubtotRetivaRg3130 string
	Eclq02llSigno_04_14        string
	Eclq02llProvIngbru         string
	Eclq02llAdicPlancuo        string
	Eclq02llSigno_04_15        string
	Eclq02llIva1AdPlancuo      string
	Eclq02llSigno_04_16        string
	Eclq02llAdic_opinter       string
	Eclq02llSigno_04_17        string
	Eclq02llIva1Ad_opinter     string
	Eclq02llSigno_04_18        string
	Eclq02llAdicAltacom        string
	Eclq02llSigno_04_19        string
	Eclq02llIva1AdAltacom      string
	Eclq02llSigno_04_20        string
	Eclq02llAdicCupmanu        string
	Eclq02llSigno_04_21        string
	Eclq02llIva1AdCupmanu      string
	Eclq02llSigno_04_22        string
	Eclq02llAdicAltacomBco     string
	Eclq02llSigno_04_23        string
	Eclq02llIva1AdAltacomBco   string
	Eclq02llSigno_04_24        string
	Filler6                    string
	Filler7                    string
	Filler8                    string
	Filler9                    string
	Eclq02llAdicMovypag        string
	Eclq02llSigno_04_27        string
	Eclq02llIva1AdicMovypag    string
	Eclq02llSigno_04_28        string
	Eclq02llRetSellos          string
	Eclq02llSigno_04_29        string
	Eclq02llProvSellos         string
	Eclq02llRetIngbru3         string
	Eclq02llSigno_04_30        string
	Eclq02llAliIngbru3         string
	Eclq02llRetIngbru4         string
	Eclq02llSigno_04_31        string
	Eclq02llAliIngbru4         string
	Eclq02llRetIngbru5         string
	Eclq02llSigno_04_32        string
	Eclq02llAliIngbru5         string
	Eclq02llRetIngbru6         string
	Eclq02llSigno_04_33        string
	Eclq02llAliIngbru6         string
	Eclq02llFiller_04_10       string
	Eclq02llAster_04_11        string
}

func (pxCuatro *PrismaPxCuatroRegistro) ValidarPxCuatro(estructuraReg *EstructuraRegistros) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro tipo (04)")
	if len(pxCuatro.Eclq02llEmpresa_04) != estructuraReg.PxCuatroDescripcionRegistro()[0].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llFpres_04) != estructuraReg.PxCuatroDescripcionRegistro()[1].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llTiporeg_04) != estructuraReg.PxCuatroDescripcionRegistro()[2].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llMoneda_04) != estructuraReg.PxCuatroDescripcionRegistro()[3].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llNumcom_04) != estructuraReg.PxCuatroDescripcionRegistro()[4].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llNumest_04) != estructuraReg.PxCuatroDescripcionRegistro()[5].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llNroliq_04) != estructuraReg.PxCuatroDescripcionRegistro()[6].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llFpag_04) != estructuraReg.PxCuatroDescripcionRegistro()[7].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llTipoliq_04) != estructuraReg.PxCuatroDescripcionRegistro()[8].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llCasacta) != estructuraReg.PxCuatroDescripcionRegistro()[9].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llTipcta) != estructuraReg.PxCuatroDescripcionRegistro()[10].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llCtabco) != estructuraReg.PxCuatroDescripcionRegistro()[11].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llCfExentoIva) != estructuraReg.PxCuatroDescripcionRegistro()[12].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_1) != estructuraReg.PxCuatroDescripcionRegistro()[13].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llLey25063) != estructuraReg.PxCuatroDescripcionRegistro()[14].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_2) != estructuraReg.PxCuatroDescripcionRegistro()[15].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAliIngbru) != estructuraReg.PxCuatroDescripcionRegistro()[16].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llDtoCampania) != estructuraReg.PxCuatroDescripcionRegistro()[17].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_3) != estructuraReg.PxCuatroDescripcionRegistro()[18].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1DtoCampania) != estructuraReg.PxCuatroDescripcionRegistro()[19].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_4) != estructuraReg.PxCuatroDescripcionRegistro()[20].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llRetIngbru2) != estructuraReg.PxCuatroDescripcionRegistro()[21].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_5) != estructuraReg.PxCuatroDescripcionRegistro()[22].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAliIngbru2) != estructuraReg.PxCuatroDescripcionRegistro()[23].Cantidad {
		return err
	}
	if len(pxCuatro.Filler1) != estructuraReg.PxCuatroDescripcionRegistro()[24].Cantidad {
		return err
	}
	if len(pxCuatro.Filler2) != estructuraReg.PxCuatroDescripcionRegistro()[25].Cantidad {
		return err
	}
	if len(pxCuatro.Filler3) != estructuraReg.PxCuatroDescripcionRegistro()[26].Cantidad {
		return err
	}
	if len(pxCuatro.Filler4) != estructuraReg.PxCuatroDescripcionRegistro()[27].Cantidad {
		return err
	}
	if len(pxCuatro.Filler5) != estructuraReg.PxCuatroDescripcionRegistro()[28].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llTasaPex) != estructuraReg.PxCuatroDescripcionRegistro()[29].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llCargoXliq) != estructuraReg.PxCuatroDescripcionRegistro()[30].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_8) != estructuraReg.PxCuatroDescripcionRegistro()[31].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1CargoXliq) != estructuraReg.PxCuatroDescripcionRegistro()[32].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_9) != estructuraReg.PxCuatroDescripcionRegistro()[33].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llDealer) != estructuraReg.PxCuatroDescripcionRegistro()[34].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llImpDbCr) != estructuraReg.PxCuatroDescripcionRegistro()[35].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_10) != estructuraReg.PxCuatroDescripcionRegistro()[36].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llCfNoReduceIva) != estructuraReg.PxCuatroDescripcionRegistro()[37].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_11) != estructuraReg.PxCuatroDescripcionRegistro()[38].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llPercepIbAgip) != estructuraReg.PxCuatroDescripcionRegistro()[39].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_12) != estructuraReg.PxCuatroDescripcionRegistro()[40].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAlicPercepIbAgip) != estructuraReg.PxCuatroDescripcionRegistro()[41].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llRetenIbAgip) != estructuraReg.PxCuatroDescripcionRegistro()[42].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_13) != estructuraReg.PxCuatroDescripcionRegistro()[43].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAlicRetenIbAgip) != estructuraReg.PxCuatroDescripcionRegistro()[44].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSubtotRetivaRg3130) != estructuraReg.PxCuatroDescripcionRegistro()[45].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_14) != estructuraReg.PxCuatroDescripcionRegistro()[46].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llProvIngbru) != estructuraReg.PxCuatroDescripcionRegistro()[47].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdicPlancuo) != estructuraReg.PxCuatroDescripcionRegistro()[48].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_15) != estructuraReg.PxCuatroDescripcionRegistro()[49].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1AdPlancuo) != estructuraReg.PxCuatroDescripcionRegistro()[50].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_16) != estructuraReg.PxCuatroDescripcionRegistro()[51].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdic_opinter) != estructuraReg.PxCuatroDescripcionRegistro()[52].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_17) != estructuraReg.PxCuatroDescripcionRegistro()[53].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1Ad_opinter) != estructuraReg.PxCuatroDescripcionRegistro()[54].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_18) != estructuraReg.PxCuatroDescripcionRegistro()[55].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdicAltacom) != estructuraReg.PxCuatroDescripcionRegistro()[56].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_19) != estructuraReg.PxCuatroDescripcionRegistro()[57].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1AdAltacom) != estructuraReg.PxCuatroDescripcionRegistro()[58].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_20) != estructuraReg.PxCuatroDescripcionRegistro()[59].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdicCupmanu) != estructuraReg.PxCuatroDescripcionRegistro()[60].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_21) != estructuraReg.PxCuatroDescripcionRegistro()[61].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1AdCupmanu) != estructuraReg.PxCuatroDescripcionRegistro()[62].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_22) != estructuraReg.PxCuatroDescripcionRegistro()[63].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdicAltacomBco) != estructuraReg.PxCuatroDescripcionRegistro()[64].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_23) != estructuraReg.PxCuatroDescripcionRegistro()[65].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1AdAltacomBco) != estructuraReg.PxCuatroDescripcionRegistro()[66].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_24) != estructuraReg.PxCuatroDescripcionRegistro()[67].Cantidad {
		return err
	}
	if len(pxCuatro.Filler6) != estructuraReg.PxCuatroDescripcionRegistro()[68].Cantidad {
		return err
	}
	if len(pxCuatro.Filler7) != estructuraReg.PxCuatroDescripcionRegistro()[69].Cantidad {
		return err
	}
	if len(pxCuatro.Filler8) != estructuraReg.PxCuatroDescripcionRegistro()[70].Cantidad {
		return err
	}
	if len(pxCuatro.Filler9) != estructuraReg.PxCuatroDescripcionRegistro()[71].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAdicMovypag) != estructuraReg.PxCuatroDescripcionRegistro()[72].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_27) != estructuraReg.PxCuatroDescripcionRegistro()[73].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llIva1AdicMovypag) != estructuraReg.PxCuatroDescripcionRegistro()[74].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_28) != estructuraReg.PxCuatroDescripcionRegistro()[75].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llRetSellos) != estructuraReg.PxCuatroDescripcionRegistro()[76].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_29) != estructuraReg.PxCuatroDescripcionRegistro()[77].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llProvSellos) != estructuraReg.PxCuatroDescripcionRegistro()[78].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llRetIngbru3) != estructuraReg.PxCuatroDescripcionRegistro()[79].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llSigno_04_30) != estructuraReg.PxCuatroDescripcionRegistro()[80].Cantidad {
		return err
	}
	if len(pxCuatro.Eclq02llAliIngbru3) != estructuraReg.PxCuatroDescripcionRegistro()[81].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llRetIngbru4) != estructuraReg.PxCuatroDescripcionRegistro()[82].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llSigno_04_31) != estructuraReg.PxCuatroDescripcionRegistro()[83].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llAliIngbru4) != estructuraReg.PxCuatroDescripcionRegistro()[84].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llRetIngbru5) != estructuraReg.PxCuatroDescripcionRegistro()[85].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llSigno_04_32) != estructuraReg.PxCuatroDescripcionRegistro()[86].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llAliIngbru5) != estructuraReg.PxCuatroDescripcionRegistro()[87].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llRetIngbru6) != estructuraReg.PxCuatroDescripcionRegistro()[88].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llSigno_04_33) != estructuraReg.PxCuatroDescripcionRegistro()[89].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llAliIngbru6) != estructuraReg.PxCuatroDescripcionRegistro()[90].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llFiller_04_10) != estructuraReg.PxCuatroDescripcionRegistro()[91].Cantidad {
		return err
	}

	if len(pxCuatro.Eclq02llAster_04_11) != estructuraReg.PxCuatroDescripcionRegistro()[92].Cantidad {
		return err
	}
	return nil
}

func (pxCuatro *PrismaPxCuatroRegistro) PxCuatroToEntities(nombreArchivo string) (pxCuatroEntity entities.Prismapxcuatroregistro) {
	pxCuatroEntity.Eclq02llEmpresa_04 = pxCuatro.Eclq02llEmpresa_04
	pxCuatroEntity.Eclq02llFpres_04 = pxCuatro.Eclq02llFpres_04
	pxCuatroEntity.Eclq02llTiporeg_04 = pxCuatro.Eclq02llTiporeg_04
	pxCuatroEntity.Eclq02llMoneda_04 = pxCuatro.Eclq02llMoneda_04
	pxCuatroEntity.Eclq02llNumcom_04 = pxCuatro.Eclq02llNumcom_04
	pxCuatroEntity.Eclq02llNumest_04 = pxCuatro.Eclq02llNumest_04
	pxCuatroEntity.Eclq02llNroliq_04 = pxCuatro.Eclq02llNroliq_04
	pxCuatroEntity.Eclq02llFpag_04 = pxCuatro.Eclq02llFpag_04
	pxCuatroEntity.Eclq02llTipoliq_04 = pxCuatro.Eclq02llTipoliq_04
	pxCuatroEntity.Eclq02llCasacta = pxCuatro.Eclq02llCasacta
	pxCuatroEntity.Eclq02llTipcta = pxCuatro.Eclq02llTipcta
	pxCuatroEntity.Eclq02llCtabco = pxCuatro.Eclq02llCtabco
	pxCuatroEntity.Eclq02llCfExentoIva = pxCuatro.Eclq02llCfExentoIva
	pxCuatroEntity.Eclq02llSigno_04_1 = pxCuatro.Eclq02llSigno_04_1
	pxCuatroEntity.Eclq02llLey25063 = pxCuatro.Eclq02llLey25063
	pxCuatroEntity.Eclq02llSigno_04_2 = pxCuatro.Eclq02llSigno_04_2
	pxCuatroEntity.Eclq02llAliIngbru = pxCuatro.Eclq02llAliIngbru
	pxCuatroEntity.Eclq02llDtoCampania = pxCuatro.Eclq02llDtoCampania
	pxCuatroEntity.Eclq02llSigno_04_3 = pxCuatro.Eclq02llSigno_04_3
	pxCuatroEntity.Eclq02llIva1DtoCampania = pxCuatro.Eclq02llIva1DtoCampania
	pxCuatroEntity.Eclq02llSigno_04_4 = pxCuatro.Eclq02llSigno_04_4
	pxCuatroEntity.Eclq02llRetIngbru2 = pxCuatro.Eclq02llRetIngbru2
	pxCuatroEntity.Eclq02llSigno_04_5 = pxCuatro.Eclq02llSigno_04_5
	pxCuatroEntity.Eclq02llAliIngbru2 = pxCuatro.Eclq02llAliIngbru2
	pxCuatroEntity.Filler1 = pxCuatro.Filler1
	pxCuatroEntity.Filler2 = pxCuatro.Filler2
	pxCuatroEntity.Filler3 = pxCuatro.Filler3
	pxCuatroEntity.Filler4 = pxCuatro.Filler4
	pxCuatroEntity.Filler5 = pxCuatro.Filler5
	pxCuatroEntity.Eclq02llTasaPex = pxCuatro.Eclq02llTasaPex
	pxCuatroEntity.Eclq02llCargoXLiq = pxCuatro.Eclq02llCargoXliq
	pxCuatroEntity.Eclq02llSigno_04_8 = pxCuatro.Eclq02llSigno_04_8
	pxCuatroEntity.Eclq02llIva1CargoXLiq = pxCuatro.Eclq02llIva1CargoXliq
	pxCuatroEntity.Eclq02llSigno_04_9 = pxCuatro.Eclq02llSigno_04_9
	pxCuatroEntity.Eclq02llDealer = pxCuatro.Eclq02llDealer
	pxCuatroEntity.Eclq02llImpDbCr = pxCuatro.Eclq02llImpDbCr
	pxCuatroEntity.Eclq02llSigno_04_10 = pxCuatro.Eclq02llSigno_04_10
	pxCuatroEntity.Eclq02llCfNoReduceIva = pxCuatro.Eclq02llCfNoReduceIva
	pxCuatroEntity.Eclq02llSigno_04_11 = pxCuatro.Eclq02llSigno_04_11
	pxCuatroEntity.Eclq02llPercepIbAgip = pxCuatro.Eclq02llPercepIbAgip
	pxCuatroEntity.Eclq02llSigno_04_12 = pxCuatro.Eclq02llSigno_04_12
	pxCuatroEntity.Eclq02llAlicPercepIbAgip = pxCuatro.Eclq02llAlicPercepIbAgip
	pxCuatroEntity.Eclq02llRetenIbAgip = pxCuatro.Eclq02llRetenIbAgip
	pxCuatroEntity.Eclq02llSigno_04_13 = pxCuatro.Eclq02llSigno_04_13
	pxCuatroEntity.Eclq02llAlicRetenIbAgip = pxCuatro.Eclq02llAlicRetenIbAgip
	pxCuatroEntity.Eclq02llSubtotalRetivaRg3130 = pxCuatro.Eclq02llSubtotRetivaRg3130
	pxCuatroEntity.Eclq02llSigno_04_14 = pxCuatro.Eclq02llSigno_04_14
	pxCuatroEntity.Eclq02llProvIngbru = pxCuatro.Eclq02llProvIngbru
	pxCuatroEntity.Eclq02llAdicPlancuo = pxCuatro.Eclq02llAdicPlancuo
	pxCuatroEntity.Eclq02llSigno_04_15 = pxCuatro.Eclq02llSigno_04_15
	pxCuatroEntity.Eclq02llIva1AdPlancuo = pxCuatro.Eclq02llIva1AdPlancuo
	pxCuatroEntity.Eclq02llSigno_04_16 = pxCuatro.Eclq02llSigno_04_16
	pxCuatroEntity.Eclq02llAdicOpinter = pxCuatro.Eclq02llAdic_opinter
	pxCuatroEntity.Eclq02llSigno_04_17 = pxCuatro.Eclq02llSigno_04_17
	pxCuatroEntity.Eclq02llIva1AdOpinter = pxCuatro.Eclq02llIva1Ad_opinter
	pxCuatroEntity.Eclq02llSigno_04_18 = pxCuatro.Eclq02llSigno_04_18
	pxCuatroEntity.Eclq02llAdicAltacom = pxCuatro.Eclq02llAdicAltacom
	pxCuatroEntity.Eclq02llSigno_04_19 = pxCuatro.Eclq02llSigno_04_19
	pxCuatroEntity.Eclq02llIva1AdAltacom = pxCuatro.Eclq02llIva1AdAltacom
	pxCuatroEntity.Eclq02llSigno_04_20 = pxCuatro.Eclq02llSigno_04_20
	pxCuatroEntity.Eclq02llAdicCupmanu = pxCuatro.Eclq02llAdicCupmanu
	pxCuatroEntity.Eclq02llSigno_04_21 = pxCuatro.Eclq02llSigno_04_21
	pxCuatroEntity.Eclq02llIva1AdCupmanu = pxCuatro.Eclq02llIva1AdCupmanu
	pxCuatroEntity.Eclq02llSigno_04_22 = pxCuatro.Eclq02llSigno_04_22
	pxCuatroEntity.Eclq02llAdicAltacomBco = pxCuatro.Eclq02llAdicAltacomBco
	pxCuatroEntity.Eclq02llSgno_04_23 = pxCuatro.Eclq02llSigno_04_23
	pxCuatroEntity.Eclq02llIva1AdAltacomBco = pxCuatro.Eclq02llIva1AdAltacomBco
	pxCuatroEntity.Eclq02llSigno_04_24 = pxCuatro.Eclq02llSigno_04_24
	pxCuatroEntity.Filler6 = pxCuatro.Filler6
	pxCuatroEntity.Filler7 = pxCuatro.Filler7
	pxCuatroEntity.Filler8 = pxCuatro.Filler8
	pxCuatroEntity.Filler9 = pxCuatro.Filler9
	pxCuatroEntity.Eclq02llAdicMovpag = pxCuatro.Eclq02llAdicMovypag
	pxCuatroEntity.Eclq02llSigno_04_27 = pxCuatro.Eclq02llSigno_04_27
	pxCuatroEntity.Eclq02llIva1AdicMovpag = pxCuatro.Eclq02llIva1AdicMovypag
	pxCuatroEntity.Eclq02llSigno_04_28 = pxCuatro.Eclq02llSigno_04_28
	pxCuatroEntity.Eclq02llRetSellos = pxCuatro.Eclq02llRetSellos
	pxCuatroEntity.Eclq02llSigno_29 = pxCuatro.Eclq02llSigno_04_29
	pxCuatroEntity.Eclq02llProvSellos = pxCuatro.Eclq02llProvSellos
	pxCuatroEntity.Eclq02llRetIngbru3 = pxCuatro.Eclq02llRetIngbru3
	pxCuatroEntity.Eclq02llSigno_04_30 = pxCuatro.Eclq02llSigno_04_30
	pxCuatroEntity.Eclq02llAliIngbru3 = pxCuatro.Eclq02llAliIngbru3
	pxCuatroEntity.Eclq02llRetIngbru4 = pxCuatro.Eclq02llRetIngbru4
	pxCuatroEntity.Eclq02llSigno_04_31 = pxCuatro.Eclq02llSigno_04_31
	pxCuatroEntity.Eclq02llAliIngbru4 = pxCuatro.Eclq02llAliIngbru4
	pxCuatroEntity.Eclq02llRetIngbru5 = pxCuatro.Eclq02llRetIngbru5
	pxCuatroEntity.Eclq02llSigno_04_32 = pxCuatro.Eclq02llSigno_04_32
	pxCuatroEntity.Eclq02llAliIngbru5 = pxCuatro.Eclq02llAliIngbru5
	pxCuatroEntity.Eclq02llRetIngbru6 = pxCuatro.Eclq02llRetIngbru6
	pxCuatroEntity.Eclq02llSigno_04_33 = pxCuatro.Eclq02llSigno_04_33
	pxCuatroEntity.Eclq02llAliIngbru6 = pxCuatro.Eclq02llAliIngbru6
	pxCuatroEntity.Eclq02llFiller_04_10 = pxCuatro.Eclq02llFiller_04_10
	pxCuatroEntity.Eclq02llAster_04 = pxCuatro.Eclq02llAster_04_11
	pxCuatroEntity.Nombrearchivo = nombreArchivo
	return
}
