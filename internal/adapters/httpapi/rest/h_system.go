package rest

import (
	"net/http"
)

func (a *St) hSystemCronTick5m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick5m()
}

func (a *St) hSystemCronTick15m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick15m()
}

func (a *St) hSystemCronTick30m(w http.ResponseWriter, r *http.Request) {
	a.ucs.SystemCronTick30m()
}
