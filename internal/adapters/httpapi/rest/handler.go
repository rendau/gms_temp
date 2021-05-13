package rest

import (
	"net/http"
)

func (a *St) hRoot(w http.ResponseWriter, r *http.Request) {
	a.uRespondJSON(w, 0, map[string]string{"hello": "world"})
}
