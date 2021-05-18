package rest

import (
	"net/http"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (a *St) hConfigUpdate(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.ConfigSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ConfigSet(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
