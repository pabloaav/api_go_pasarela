package administracionfake

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"gorm.io/gorm"
)

var estructuraValidasPagoIntento []entities.Pagointento = []entities.Pagointento{

	{
		Model:                gorm.Model{ID: 6},
		PagosID:              7,
		MediopagosID:         1,
		InstallmentdetailsID: 1,
		ExternalID:           "10676472",
		PaidAt:               time.Now(),
		ReportAt:             time.Now(),
		IsAvailable:          false,
		Amount:               66051,
		StateComment:         "approved",
		Barcode:              "",
		BarcodeUrl:           "",
		AvailableAt:          time.Now(),

		HolderName:         "Ignacio Fernandez",
		HolderEmail:        "Ignacio Fernandez",
		TicketNumber:       "",
		AuthorizationCode:  "152523",
		CardLastFourDigits: "4905",
	},
	{
		Model:                gorm.Model{ID: 8},
		PagosID:              8,
		MediopagosID:         1,
		InstallmentdetailsID: 1,
		ExternalID:           "10676500",
		PaidAt:               time.Now(),
		ReportAt:             time.Now(),
		IsAvailable:          false,
		Amount:               7651,
		StateComment:         "approved",
		Barcode:              "",
		BarcodeUrl:           "",
		AvailableAt:          time.Now(),

		HolderName:         "Castro",
		HolderEmail:        "Castro",
		TicketNumber:       "",
		AuthorizationCode:  "153031",
		CardLastFourDigits: "4905",
	},

	{
		Model:                gorm.Model{ID: 5},
		PagosID:              11,
		MediopagosID:         1,
		InstallmentdetailsID: 1,
		ExternalID:           "10676452",
		PaidAt:               time.Now(),
		ReportAt:             time.Now(),
		IsAvailable:          false,
		Amount:               76051,
		StateComment:         "approved",
		Barcode:              "",
		BarcodeUrl:           "",
		AvailableAt:          time.Now(),

		HolderName:         "Ignacio Fernandez",
		HolderEmail:        "Ignacio Fernandez",
		TicketNumber:       "",
		AuthorizationCode:  "152414",
		CardLastFourDigits: "4905",
	},

	{
		Model:                gorm.Model{ID: 3},
		PagosID:              12,
		MediopagosID:         1,
		InstallmentdetailsID: 1,
		ExternalID:           "10676296",
		PaidAt:               time.Now(),
		ReportAt:             time.Now(),
		IsAvailable:          false,
		Amount:               76051,
		StateComment:         "approved",
		Barcode:              "",
		BarcodeUrl:           "",
		// AvailableAt:        time.Now(),
		RevertedAt:         time.Now(),
		HolderName:         "De la Cruz N",
		HolderEmail:        "De la Cruz N",
		TicketNumber:       "",
		AuthorizationCode:  "123",
		CardLastFourDigits: "4905",
	},
}
