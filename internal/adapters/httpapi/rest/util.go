package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rendau/gms_temp/internal/domain/errs"
)

func (a *St) uParseRequestJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dst); err != nil {
		a.uHandleError(errs.BadJson, r, w)

		return false
	}

	return true
}

func (a *St) uRespondJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		a.lg.Infow("Fail to send response", "error", err)
	}
}

func (a *St) uHandleError(err error, r *http.Request, w http.ResponseWriter) bool {
	if err != nil {
		switch cErr := err.(type) {
		case errs.Err:
			a.uRespondJSON(w, http.StatusBadRequest, ErrRepSt{
				ErrorCode: cErr.Error(),
			})
		default:
			a.uLogErrorRequest(r, err, "Error in http handler")

			w.WriteHeader(http.StatusInternalServerError)
		}

		return true
	}

	return false
}

func (a *St) uLogErrorRequest(r *http.Request, err interface{}, msg string) {
	a.lg.Errorw(
		msg,
		err,
		"method", r.Method,
		"path", r.URL,
	)
}

func (a *St) uGetRequestToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" { // try from query parameter
		token = r.URL.Query().Get("auth_token")
	}

	return token
}

func (a *St) uExtractPaginationPars(pars url.Values) (offset int64, limit int64, page int64) {
	var err error

	qPar := pars.Get("page_size")
	if qPar != "" {
		limit, err = strconv.ParseInt(qPar, 10, 64)
		if err != nil {
			limit = 0
		}
	}

	qPar = pars.Get("page")
	if qPar != "" {
		page, err = strconv.ParseInt(qPar, 10, 64)
		if err != nil {
			page = 0
		}
	}
	if page == 0 {
		page = 1
	}

	offset = (page - 1) * limit

	return offset, limit, page
}

func (a *St) uQpParseBool(values url.Values, key string) *bool {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseBool(qp[0]); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseInt64(values url.Values, key string) *int64 {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseInt(qp[0], 10, 64); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseFloat64(values url.Values, key string) *float64 {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseFloat(qp[0], 64); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseInt(values url.Values, key string) *int {
	if qp, ok := values[key]; ok {
		if result, err := strconv.Atoi(qp[0]); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseString(values url.Values, key string) *string {
	if qp, ok := values[key]; ok {
		return &(qp[0])
	}
	return nil
}

func (a *St) uQpParseTime(values url.Values, key string) *time.Time {
	if qp, ok := values[key]; ok {
		if result, err := time.Parse(time.RFC3339, qp[0]); err == nil {
			return &result
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func (a *St) uQpParseInt64Slice(values url.Values, key string) *[]int64 {
	if _, ok := values[key]; ok {
		items := strings.Split(values.Get(key), ",")

		result := make([]int64, 0, len(items))

		for _, vStr := range items {
			if v, err := strconv.ParseInt(vStr, 10, 64); err == nil {
				result = append(result, v)
			}
		}

		return &result
	}

	return nil
}
