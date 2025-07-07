package dto

type Meta struct {
	Page      int `json:"page,omitempty"`
	Limit     int `json:"limit,omitempty"`
	Total     int `json:"total,omitempty"`
	TotalPage int `json:"total_page,omitempty"`
}

type Response[T any] struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    T      `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}
