package entities

import "time"

type PaginationParams struct {
	Page           int64 `json:"page"`
	PageSize       int64 `json:"page_size"`
	WithTotalCount bool  `json:"with_total_count"`
}

type PeriodFilterPars struct {
	TsGTE *time.Time
	TsLTE *time.Time
}

type ChartVByTime struct {
	Ts time.Time `json:"ts"`
	V  int64     `json:"v"`
}
