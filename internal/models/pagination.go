package models

type Pagination struct {
	PageSize   string `json:"PageSize,omitempty"`
	PageNumber string `json:"PageNumber,omitempty"`
}