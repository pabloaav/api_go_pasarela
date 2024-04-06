package cierrelotedtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
)

type MovimientoMxDetalleRegistro struct {
	Empresa                   string
	Fechapresentacion         string
	Tiporeg                   string
	Numcom                    string
	Numest                    string
	Codop                     string
	Tipoaplic                 string
	Lote                      string
	Codbco                    string
	Codcasa                   string
	Bcoest                    string
	Bcocasa                   string
	Numtar                    string
	ForigCompra               string
	Fechapag                  string
	Numcomp                   string
	Importe                   string
	Signo                     string
	Numaut                    string
	Numcuot                   string
	Plancuot                  string
	RecAcep                   string
	RechPrint                 string
	RechSecun                 string
	ImpPlan                   string
	Signo1                    string
	McaPex                    string
	Nroliq                    string
	CcoOrigen                 string
	CcoMotivo                 string
	IdCargoliq                string
	Moneda                    string
	PromoBonifUsu             string
	PromoBonifEst             string
	IdPromo                   string
	MovImporig                string
	SgImporig                 string
	IdCf                      string
	CfExentoIva               string
	Dealer                    string
	CuitEst                   string
	FechapagAjuLqe            string
	CodMotivoAjuLqe           string
	IdentifNroFactura         string
	PorcdtoArancel            string
	Arancel                   string
	SignoArancel              string
	TnaCf                     string
	ImporteCostoFin           string
	SigImporteCostoFinanciero string
	IdTx                      string
	Agencia                   string
	TipoPlan                  string
	BanderaEst                string
	Subcodigo                 string
	Filler                    string
	CcoMotivoMc               string
	NumtarXl                  string
	NumautXl                  string
	Filler1                   string
	IdFinRegistro             string
}

func (mxDetalles *MovimientoMxDetalleRegistro) ValidarMxDetalle(estructuraReg *EstructuraRegistros) (erro error) {
	err := errors.New("longitud del campo es incorrecto para el registro de detalles")
	if len(mxDetalles.Empresa) != estructuraReg.MxDetalleDescripcionRegistro()[0].Cantidad {
		return err
	}
	if len(mxDetalles.Fechapresentacion) != estructuraReg.MxDetalleDescripcionRegistro()[1].Cantidad {
		return err
	}
	if len(mxDetalles.Tiporeg) != estructuraReg.MxDetalleDescripcionRegistro()[2].Cantidad {
		return err
	}
	if len(mxDetalles.Numcom) != estructuraReg.MxDetalleDescripcionRegistro()[3].Cantidad {
		return err
	}
	if len(mxDetalles.Numest) != estructuraReg.MxDetalleDescripcionRegistro()[4].Cantidad {
		return err
	}
	if len(mxDetalles.Codop) != estructuraReg.MxDetalleDescripcionRegistro()[5].Cantidad {
		return err
	}
	if len(mxDetalles.Tipoaplic) != estructuraReg.MxDetalleDescripcionRegistro()[6].Cantidad {
		return err
	}
	if len(mxDetalles.Lote) != estructuraReg.MxDetalleDescripcionRegistro()[7].Cantidad {
		return err
	}
	if len(mxDetalles.Codbco) != estructuraReg.MxDetalleDescripcionRegistro()[8].Cantidad {
		return err
	}
	if len(mxDetalles.Codcasa) != estructuraReg.MxDetalleDescripcionRegistro()[9].Cantidad {
		return err
	}
	if len(mxDetalles.Bcoest) != estructuraReg.MxDetalleDescripcionRegistro()[10].Cantidad {
		return err
	}
	if len(mxDetalles.Bcocasa) != estructuraReg.MxDetalleDescripcionRegistro()[11].Cantidad {
		return err
	}
	// if len(mxDetalles.Numtar) != estructuraReg.MxDetalleDescripcionRegistro()[12].Cantidad {
	// 	return err
	// }
	if len(mxDetalles.ForigCompra) != estructuraReg.MxDetalleDescripcionRegistro()[13].Cantidad {
		return err
	}
	if len(mxDetalles.Fechapag) != estructuraReg.MxDetalleDescripcionRegistro()[14].Cantidad {
		return err
	}
	if len(mxDetalles.Numcomp) != estructuraReg.MxDetalleDescripcionRegistro()[15].Cantidad {
		return err
	}
	if len(mxDetalles.Importe) != estructuraReg.MxDetalleDescripcionRegistro()[16].Cantidad {
		return err
	}
	if len(mxDetalles.Signo) != estructuraReg.MxDetalleDescripcionRegistro()[17].Cantidad {
		return err
	}
	if len(mxDetalles.Numaut) != estructuraReg.MxDetalleDescripcionRegistro()[18].Cantidad {
		return err
	}
	if len(mxDetalles.Numcuot) != estructuraReg.MxDetalleDescripcionRegistro()[19].Cantidad {
		return err
	}
	if len(mxDetalles.Plancuot) != estructuraReg.MxDetalleDescripcionRegistro()[20].Cantidad {
		return err
	}
	if len(mxDetalles.RecAcep) != estructuraReg.MxDetalleDescripcionRegistro()[21].Cantidad {
		return err
	}
	if len(mxDetalles.RechPrint) != estructuraReg.MxDetalleDescripcionRegistro()[22].Cantidad {
		return err
	}
	if len(mxDetalles.RechSecun) != estructuraReg.MxDetalleDescripcionRegistro()[23].Cantidad {
		return err
	}
	if len(mxDetalles.ImpPlan) != estructuraReg.MxDetalleDescripcionRegistro()[24].Cantidad {
		return err
	}
	if len(mxDetalles.Signo1) != estructuraReg.MxDetalleDescripcionRegistro()[25].Cantidad {
		return err
	}
	if len(mxDetalles.McaPex) != estructuraReg.MxDetalleDescripcionRegistro()[26].Cantidad {
		return err
	}
	if len(mxDetalles.Nroliq) != estructuraReg.MxDetalleDescripcionRegistro()[27].Cantidad {
		return err
	}
	if len(mxDetalles.CcoOrigen) != estructuraReg.MxDetalleDescripcionRegistro()[28].Cantidad {
		return err
	}
	if len(mxDetalles.CcoMotivo) != estructuraReg.MxDetalleDescripcionRegistro()[29].Cantidad {
		return err
	}
	if len(mxDetalles.IdCargoliq) != estructuraReg.MxDetalleDescripcionRegistro()[30].Cantidad {
		return err
	}
	if len(mxDetalles.Moneda) != estructuraReg.MxDetalleDescripcionRegistro()[31].Cantidad {
		return err
	}
	if len(mxDetalles.PromoBonifUsu) != estructuraReg.MxDetalleDescripcionRegistro()[32].Cantidad {
		return err
	}
	if len(mxDetalles.PromoBonifEst) != estructuraReg.MxDetalleDescripcionRegistro()[33].Cantidad {
		return err
	}
	if len(mxDetalles.IdPromo) != estructuraReg.MxDetalleDescripcionRegistro()[34].Cantidad {
		return err
	}
	if len(mxDetalles.MovImporig) != estructuraReg.MxDetalleDescripcionRegistro()[35].Cantidad {
		return err
	}
	if len(mxDetalles.SgImporig) != estructuraReg.MxDetalleDescripcionRegistro()[36].Cantidad {
		return err
	}
	if len(mxDetalles.IdCf) != estructuraReg.MxDetalleDescripcionRegistro()[37].Cantidad {
		return err
	}
	if len(mxDetalles.CfExentoIva) != estructuraReg.MxDetalleDescripcionRegistro()[38].Cantidad {
		return err
	}
	if len(mxDetalles.Dealer) != estructuraReg.MxDetalleDescripcionRegistro()[39].Cantidad {
		return err
	}
	if len(mxDetalles.CuitEst) != estructuraReg.MxDetalleDescripcionRegistro()[40].Cantidad {
		return err
	}
	if len(mxDetalles.FechapagAjuLqe) != estructuraReg.MxDetalleDescripcionRegistro()[41].Cantidad {
		return err
	}
	if len(mxDetalles.CodMotivoAjuLqe) != estructuraReg.MxDetalleDescripcionRegistro()[42].Cantidad {
		return err
	}
	if len(mxDetalles.IdentifNroFactura) != estructuraReg.MxDetalleDescripcionRegistro()[43].Cantidad {
		return err
	}
	if len(mxDetalles.PorcdtoArancel) != estructuraReg.MxDetalleDescripcionRegistro()[44].Cantidad {
		return err
	}
	if len(mxDetalles.Arancel) != estructuraReg.MxDetalleDescripcionRegistro()[45].Cantidad {
		return err
	}
	if len(mxDetalles.SignoArancel) != estructuraReg.MxDetalleDescripcionRegistro()[46].Cantidad {
		return err
	}
	if len(mxDetalles.TnaCf) != estructuraReg.MxDetalleDescripcionRegistro()[47].Cantidad {
		return err
	}
	if len(mxDetalles.ImporteCostoFin) != estructuraReg.MxDetalleDescripcionRegistro()[48].Cantidad {
		return err
	}
	if len(mxDetalles.SigImporteCostoFinanciero) != estructuraReg.MxDetalleDescripcionRegistro()[49].Cantidad {
		return err
	}
	if len(mxDetalles.IdTx) != estructuraReg.MxDetalleDescripcionRegistro()[50].Cantidad {
		return err
	}
	if len(mxDetalles.Agencia) != estructuraReg.MxDetalleDescripcionRegistro()[51].Cantidad {
		return err
	}
	if len(mxDetalles.TipoPlan) != estructuraReg.MxDetalleDescripcionRegistro()[52].Cantidad {
		return err
	}
	if len(mxDetalles.BanderaEst) != estructuraReg.MxDetalleDescripcionRegistro()[53].Cantidad {
		return err
	}
	if len(mxDetalles.Subcodigo) != estructuraReg.MxDetalleDescripcionRegistro()[54].Cantidad {
		return err
	}
	if len(mxDetalles.Filler) != estructuraReg.MxDetalleDescripcionRegistro()[55].Cantidad {
		return err
	}
	if len(mxDetalles.CcoMotivoMc) != estructuraReg.MxDetalleDescripcionRegistro()[56].Cantidad {
		return err
	}
	// if len(mxDetalles.NumtarXl) != estructuraReg.MxDetalleDescripcionRegistro()[57].Cantidad {
	// 	return err
	// }
	if len(mxDetalles.NumautXl) != estructuraReg.MxDetalleDescripcionRegistro()[58].Cantidad {
		return err
	}
	if len(mxDetalles.Filler1) != estructuraReg.MxDetalleDescripcionRegistro()[59].Cantidad {
		return err
	}
	if len(mxDetalles.IdFinRegistro) != estructuraReg.MxDetalleDescripcionRegistro()[60].Cantidad {
		return err
	}
	return
}

func (mxDetalles *MovimientoMxDetalleRegistro) MxDetalleToEntities() (mxDetallesEntity entities.Prismamxdetallemovimiento) {
	mxDetallesEntity.Empresa = mxDetalles.Empresa
	mxDetallesEntity.Fechapresentacion = mxDetalles.Fechapresentacion
	mxDetallesEntity.Tiporeg = mxDetalles.Tiporeg
	mxDetallesEntity.Numcom = mxDetalles.Numcom
	mxDetallesEntity.Numest = mxDetalles.Numest
	mxDetallesEntity.Codop = mxDetalles.Codop
	mxDetallesEntity.Tipoaplic = mxDetalles.Tipoaplic
	mxDetallesEntity.Lote = mxDetalles.Lote
	mxDetallesEntity.Codbco = mxDetalles.Codbco
	mxDetallesEntity.Codcasa = mxDetalles.Codcasa
	mxDetallesEntity.Bcoest = mxDetalles.Bcoest
	mxDetallesEntity.Bcocasa = mxDetalles.Bcocasa
	mxDetallesEntity.Numtar = mxDetalles.Numtar
	mxDetallesEntity.ForigCompra = mxDetalles.ForigCompra
	mxDetallesEntity.Fechapag = mxDetalles.Fechapag
	mxDetallesEntity.Numcomp = mxDetalles.Numcomp
	mxDetallesEntity.Importe = mxDetalles.Importe
	mxDetallesEntity.Signo = mxDetalles.Signo
	mxDetallesEntity.Numaut = mxDetalles.Numaut
	mxDetallesEntity.Numcuot = mxDetalles.Numcuot
	mxDetallesEntity.Plancuot = mxDetalles.Plancuot
	mxDetallesEntity.RecAcep = mxDetalles.RecAcep
	mxDetallesEntity.RechPrint = mxDetalles.RechPrint
	mxDetallesEntity.RechSecun = mxDetalles.RechSecun
	mxDetallesEntity.ImpPlan = mxDetalles.ImpPlan
	mxDetallesEntity.Signo1 = mxDetalles.Signo1
	mxDetallesEntity.McaPex = mxDetalles.McaPex
	mxDetallesEntity.Nroliq = mxDetalles.Nroliq
	mxDetallesEntity.CcoOrigen = mxDetalles.CcoOrigen
	mxDetallesEntity.CcoMotivo = mxDetalles.CcoMotivo
	mxDetallesEntity.IdCargoliq = mxDetalles.IdCargoliq
	mxDetallesEntity.Moneda = mxDetalles.Moneda
	mxDetallesEntity.PromoBonifUsu = mxDetalles.PromoBonifUsu
	mxDetallesEntity.PromoBonifEst = mxDetalles.PromoBonifEst
	mxDetallesEntity.IdPromo = mxDetalles.IdPromo
	mxDetallesEntity.MovImporig = mxDetalles.MovImporig
	mxDetallesEntity.SgImporig = mxDetalles.SgImporig
	mxDetallesEntity.IdCf = mxDetalles.IdCf
	mxDetallesEntity.CfExentoIva = mxDetalles.CfExentoIva
	mxDetallesEntity.Dealer = mxDetalles.Dealer
	mxDetallesEntity.CuitEst = mxDetalles.CuitEst
	mxDetallesEntity.FechapagAjuLqe = mxDetalles.FechapagAjuLqe
	mxDetallesEntity.CodMotivoAjuLqe = mxDetalles.CodMotivoAjuLqe
	mxDetallesEntity.IdentifNroFactura = mxDetalles.IdentifNroFactura
	mxDetallesEntity.PorcdtoArancel = mxDetalles.PorcdtoArancel
	mxDetallesEntity.Arancel = mxDetalles.Arancel
	mxDetallesEntity.SignoArancel = mxDetalles.SignoArancel
	mxDetallesEntity.TnaCf = mxDetalles.TnaCf
	mxDetallesEntity.ImporteCostoFin = mxDetalles.ImporteCostoFin
	mxDetallesEntity.SigImporteCostoFinanciero = mxDetalles.SigImporteCostoFinanciero
	mxDetallesEntity.IdTx = mxDetalles.IdTx
	mxDetallesEntity.Agencia = mxDetalles.Agencia
	mxDetallesEntity.TipoPlan = mxDetalles.TipoPlan
	mxDetallesEntity.BanderaEst = mxDetalles.BanderaEst
	mxDetallesEntity.Subcodigo = mxDetalles.Subcodigo
	mxDetallesEntity.Filler = mxDetalles.Filler
	mxDetallesEntity.CcoMotivoMc = mxDetalles.CcoMotivoMc
	mxDetallesEntity.NumtarXl = mxDetalles.NumtarXl
	mxDetallesEntity.NumautXl = mxDetalles.NumautXl
	mxDetallesEntity.Filler1 = mxDetalles.Filler1
	mxDetallesEntity.IdFinRegistro = mxDetalles.IdFinRegistro
	return
}
