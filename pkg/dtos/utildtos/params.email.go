package utildtos

type ParamsEmail struct {
	Email                 []string
	Nombre                string
	Mensaje               string
	Descripcion           DescripcionTemplate
	Totales               TotalesTemplate
	MensajeSegunMedioPago MensajeSegunMedioPagoStruct
}
