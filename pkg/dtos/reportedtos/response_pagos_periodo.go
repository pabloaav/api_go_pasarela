package reportedtos

type ResponsePagosPeriodo struct {
	Cliente                 string  `json:"cliente"`
	Cuenta                  string  `json:"cuenta"`
	Pagotipo                string  `json:"tipopago"`
	ExternalReference       string  `json:"external_reference"`
	IdPago                  uint    `json:"idpago"`
	Estado                  string  `json:"estado"`
	ChannelId               int64   `json:"channelid"`
	ExternalId              string  `json:"externalid"`
	TransactionId           string  `json:"transactionid"`
	Barcode                 string  `json:"barcode"`
	IdExterno               string  `json:"idexterno"`
	MedioPago               string  `json:"mediopago"`
	Pagador                 string  `json:"pagador"`
	DniPagador              string  `json:"dni_pagador"`
	Cuotas                  uint    `json:"cuotas"`
	Nroestablecimiento      string  `json:"Número de establecimiento"`
	NroLiquidacion          string  `json:"nro_liquidacion"`
	FechaPago               string  `json:"fechapago"`
	FechaPresentacion       string  `json:"fechapresentacion"`
	FechaAcreditacion       string  `json:"fechaacreditacion"`
	FechaRendicion          string  `json:"fecharendicion"`
	Amount                  float64 `json:"importe_a_pagar"`
	AmountPagado            float64 `json:"importe_pagado"`
	ArancelPorcentaje       float64 `json:"arancel_porcentaje"`
	CftCoeficiente          uint    `json:"cft_coeficiente"`
	RetencionIva            string  `json:"retencion_iva"`
	ImporteMinimo           float64 `json:"importe_minimo"`
	ImporteMaximo           float64 `json:"importe_maximo"`
	ArancelPorcentajeMinimo float64 `json:"arancel_porcentaje_minimo"`
	ArancelPorcentajeMaximo float64 `json:"arancel_porcentaje_maximo"`
	Revertido               bool    `json:"revertido"`
	Enobservacion           bool    `json:"enobservacion"`
	Cantdias                int64   `json:"cantdias"`

	//esto corresponde a la cierrelote apilink
	CostoFijoTransaccion float64 `json:"costo_fijo_transaccion"`

	//PRISMA y RAPIPAGO ?
	ImporteArancel       float64 `json:"importe_arancel"`
	ImporteArancelIva    float64 `json:"importe_arancel_iva"`
	ImporteArancelIvaMov float64 `json:"importe_arancel_iva_mov"`
	ImporteCft           float64 `json:"importe_cft"`

	//SOLO SI ES MOVIMIENTO
	ImporteNetoCobrado      float64 `json:"importe_neto_cobrado"`
	ComisionPorcentaje      float64 `json:"comision_porcentaje"`
	ComisionPorcentajeIva   float64 `json:"comision_porcentaje_iva"`
	ImporteComisionSobreTap float64 `json:"importe_comision_sobre_tap"`
	ImporteIvaComisionTap   float64 `json:"importe_iva_comision_tap"`
	ImporteRendido          float64 `json:"importe_rendido"` // MONTO TRANSFERIDO AL CLIENTE
	// ImporteTelco   	    loat64 `json:"importe_telco"`       // MONTO TRANSFERIDO AL TELCO

	ReferenciaBancaria string `json:"referencia_bancaria"`
}

type ResultadoPagosPeriodo struct {
	Cliente                 string  `json:"Cliente"`
	Cuenta                  string  `json:"Cuenta"`
	Pagotipo                string  `json:"Tipo de pago"`
	ExternalReference       string  `json:"Referencia externa"`
	IdPago                  uint    `json:"Id de pago"`
	Estado                  string  `json:"Estado"`
	MedioPago               string  `json:"Medio de pago"`
	Pagador                 string  `json:"Pagador"`
	Dni                     string  `json:"Dni"`
	Cuotas                  uint    `json:"Cuotas"`
	FechaPago               string  `json:"Fecha de pago"`
	Nroestablecimiento      string  `json:"Número de establecimiento"`
	NroLiquidacion          string  `json:"Número de liquidación"`
	FechaPresentacion       string  `json:"Fecha de presentacion"`
	FechaAcreditacion       string  `json:"Fecha acreditacion"`
	FechaRendicion          string  `json:"Fecha rendicion"`
	Amount                  float64 `json:"Importe a pagar"`
	AmountPagado            float64 `json:"Importe pagado"`
	ArancelPorcentaje       float64 `json:"Arancel porcentaje"`
	CftCoeficiente          uint    `json:"Cft coeficiente"`
	RetencionIva            string  `json:"Retencion Iva"`
	ImporteMinimo           float64 `json:"Importe minimo"`
	ImporteMaximo           float64 `json:"importe maximo"`
	ArancelPorcentajeMinimo float64 `json:"Porcentaje arancel minimo"`
	ArancelPorcentajeMaximo float64 `json:"Porcentaje arancel maximo"`

	//esto corresponde a la cierrelote apilink
	CostoFijoTransaccion float64 `json:"Costo fijo de transaccion"`

	//PRISMA y RAPIPAGO ?
	ImporteArancel       float64 `json:"Importe arancel"`
	ImporteArancelIva    float64 `json:"Importe arancel iva"`
	ImporteArancelIvaMov float64 `json:"Importe arancel iva calculado"`
	ImporteCft           float64 `json:"Importe cft"`

	//SOLO SI ES MOVIMIENTO
	ComisionPorcentaje      float64 `json:"Comision porcentaje"`
	ComisionPorcentajeIva   float64 `json:"Comision porcentaje iva"`
	ImporteComisionSobreTap float64 `json:"Importe comision sobre tap"`
	ImporteIvaComisionTap   float64 `json:"Importe iva comision tap"`
	ImporteNetoCobrado      float64 `json:"Importe neto cobrado"`
	ImporteRendido          float64 `json:"Importe_rendido"` // MONTO TRANSFERIDO AL CLIENTE
	// ImporteTelco   	    loat64 `json:"importe_telco"`       // MONTO TRANSFERIDO AL TELCO
	Revertido          bool   `json:"Revertido"`
	Enobservacion      bool   `json:"En Observacion"`
	Cantdias           int64  `json:"Cantidad de dias"`
	ReferenciaBancaria string `json:"Referencia Bancaria"`
}
