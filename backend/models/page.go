package models

type PageResponse struct {
	Data  []Article `json:"data"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int64     `json:"total"`
}