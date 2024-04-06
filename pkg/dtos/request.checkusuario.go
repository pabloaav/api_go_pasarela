package dtos

type RequestCheckusuario struct {
	HolderName            string `json:"holder_name"`
	HolderEmail           string `json:"holder_email"`
	HolderDocNum          string `json:"holder_docNum"`
}