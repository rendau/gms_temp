package rest

import (
	"net/http"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

// swagger:route PUT /config config hConfigUpdate
// Конфикурационные данные системы.
// Allowed to: `UsrTypeAdmin`
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hConfigUpdate(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hConfigUpdate
	type docReqSt struct {
		// in: body
		Body entities.ConfigSt
	}

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
