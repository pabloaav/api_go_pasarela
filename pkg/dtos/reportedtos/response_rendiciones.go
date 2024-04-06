package reportedtos

import "time"

type ResponseRendiciones struct {
	AccountId        string                       `json:"account_id"`
	ReportDate       time.Time                    `json:"report_date"`
	TotalCredits     uint64                       `json:"total_credits"`
	CreditAmount     float64                      `json:"credit_amount"`
	TotalDebits      uint64                       `json:"total_debits"`
	DebitAmount      float64                      `json:"debit_amount"`
	SettlementAmount float64                      `json:"settlement_amount"`
	Data             []ResponseDetalleRendiciones `json:"data"`
}

type ResponseDetalleRendiciones struct {
	RequestId         int64   `json:"request_id"`
	ExternalReference string  `json:"external_reference"`
	Credit            float64 `json:"credit"`
	Debit             float64 `json:"debit"`
}
