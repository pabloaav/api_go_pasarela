package bancodtos

type RequestTokenBanco struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SistemaId string `json:"sistema_id"`
}
