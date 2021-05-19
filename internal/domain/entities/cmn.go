package entities

import "time"

type PaginationParams struct {
	Offset int64
	Limit  int64
}

type PeriodFilterPars struct {
	TsGTE *time.Time
	TsLTE *time.Time
}

type ChartVByTime struct {
	Ts time.Time `json:"ts"`
	V  int64     `json:"v"`
}
