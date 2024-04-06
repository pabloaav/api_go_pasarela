package dtos

type HttpResponse struct {
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Status      int         `json:"status"`
	Detail      string      `json:"detail"`
	Parametters interface{} `json:"parametters"`
}
