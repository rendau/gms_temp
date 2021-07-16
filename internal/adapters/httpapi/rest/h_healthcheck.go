package rest

import (
	"net/http"
)

func (a *St) hHealthCheck(w http.ResponseWriter, r *http.Request) {
	a.uRespondJSON(w, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
