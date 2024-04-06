package linkcuentas

type GetCuentasResponse struct {
	Cbu         string                    `json:"cbu"`
	Tipo        string                    `json:"tipo"`
	AliasCbu    string                    `json:"aliasCbu"`
	Moneda      string                    `json:"moneda"`
	NroCuenta   string                    `json:"nroCuenta"`
	NombreBanco string                    `json:"nombreBanco"`
	Email       string                    `json:"email"`
	Adhesion    GetCuentaAdhesionREsponse `json:"adhesion"`
}

type GetCuentaAdhesionREsponse struct {
	FechaAdhesion string `json:"fechaAdhesion"`
	Adherida      string `json:"adherida"`
}

