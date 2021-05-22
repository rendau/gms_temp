package rest

// swagger:response error_reply
type docErrRepSt struct {
	// in:body
	Body ErrRepSt
}

type ErrRepSt struct {
	ErrorCode string `json:"error_code"`
}

type PaginatedListRepSt struct {
	Page       int64       `json:"page"`
	PageSize   int64       `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	Results    interface{} `json:"results"`
}
