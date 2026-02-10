package api

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginationOutput struct {
	Items      []interface{} `json:"items"`
	TotalCount int           `json:"totalCount"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}
