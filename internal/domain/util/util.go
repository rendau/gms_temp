package util

import (
	"strconv"
	"time"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
)

func RequirePageSize(pars entities.PaginationParams, allowedPageSize int64) error {
	if allowedPageSize == 0 {
		allowedPageSize = cns.MaximalPageSize
	}

	if pars.Limit == 0 || pars.Limit > allowedPageSize {
		return errs.IncorrectPageSize
	}

	return nil
}

func NewInt(v int) *int {
	return &v
}

func NewInt64(v int64) *int64 {
	return &v
}

func NewFloat64(v float64) *float64 {
	return &v
}

func NewString(v string) *string {
	return &v
}

func NewBool(v bool) *bool {
	return &v
}

func NewTime(v time.Time) *time.Time {
	return &v
}

func NewSliceInt64(v ...int64) *[]int64 {
	res := make([]int64, 0, len(v))
	res = append(res, v...)
	return &res
}

func NewSliceString(v ...string) *[]string {
	res := make([]string, 0, len(v))
	res = append(res, v...)
	return &res
}

func Int64SliceToString(src []int64, delimiter, emptyV string) string {
	if len(src) == 0 {
		return emptyV
	}

	res := ""

	for _, v := range src {
		if res != "" {
			res += delimiter
		}
		res += strconv.FormatInt(v, 10)
	}

	return res
}
