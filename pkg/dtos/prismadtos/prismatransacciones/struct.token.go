package prismadtos

import "errors"

type StructToken struct {
	Card        Card                `json:"card,omitempty"`
	DataOffline OfflineTokenRequest `json:"data_offline,omitempty"`
	TypePay     EnumTipoPagoPrisma  `json:"type_pay"`
}

func (ts *StructToken) ValidarSolicitudTokenOfflineRequest() error {

	err := ts.TypePay.IsValid()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
