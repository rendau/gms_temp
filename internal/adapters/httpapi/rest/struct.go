package rest

type ErrRepSt struct {
	ErrorCode string `json:"error_code"`
}

type PaginatedListRepSt struct {
	Page       int64       `json:"page"`
	PageSize   int64       `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	Results    interface{} `json:"results"`
}
