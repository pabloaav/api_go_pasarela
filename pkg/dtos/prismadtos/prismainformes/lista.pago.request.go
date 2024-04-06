package prismadtos

type ListaPagosRequest struct {
	//offset nro de pagina
	Offset int64 `json:"offset,omitempty"`
	//cantodad elemento por pagina valor por defecto 50
	PageSize int64 `json:"pageSize,omitempty"`
	// es euivalente a site_transaction_id
	SiteOperationId string `json:"siteOperationId,omitempty"`
	// equivalente a site_id
	MerchantId string `json:"merchantId,omitempty"`
	DateFrom   string `json:"dateFrom,omitempty"`
	DateTo     string `json:"dateTos,omitempty"`
}
