package dtos

type Meta struct {
	Page Page `json:"page"`
}

type Page struct {
	CurrentPage int32 `json:"currentPage"`
	From        int32 `json:"from"`
	LastPage    int32 `json:"lastPage"`
	PerPage     int32 `json:"perPage"`
	To          int32 `json:"to"`
	Total       int32 `json:"total"`
}

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}
