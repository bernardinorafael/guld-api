package dto

type SearchParams struct {
	Sort  string `json:"sort"`
	Query string `json:"q"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
