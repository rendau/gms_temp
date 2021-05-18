package rest

import (
	"encoding/json"
	"net/http"
)

func (a *St) hDicGet(w http.ResponseWriter, r *http.Request) {
	reqQuery := r.URL.Query()

	rhs := reqQuery.Get("hs")

	hs, data, err := a.ucs.DicGetJson(a.uGetRequestContext(r), rhs)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, struct {
		Hs   string          `json:"hs"`
		Data json.RawMessage `json:"data"`
	}{
		Hs:   hs,
		Data: data,
	})
}
