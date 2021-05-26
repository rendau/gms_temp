package rest

import (
	"net/http"
)

func (a *St) hSystemSmsBalanceAlarm(w http.ResponseWriter, r *http.Request) {
	reqObj := map[string]float64{}
	if !a.uParseRequestJSON(w, r, &reqObj) {
		return
	}

	a.ucs.SystemSmsBalanceAlarmCb(int64(reqObj["balance"]))
}

func (a *St) hSystemFilterUnusedFiles(w http.ResponseWriter, r *http.Request) {
	reqObj := make([]string, 0)
	if !a.uParseRequestJSON(w, r, &reqObj) {
		return
	}

	result := a.ucs.SystemFilterUnusedFiles(reqObj)

	a.uRespondJSON(w, result)
}

func (a *St) hSystemCronTick5m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick5m()
}

func (a *St) hSystemCronTick15m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick15m()
}

func (a *St) hSystemCronTick30m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick30m()
}
