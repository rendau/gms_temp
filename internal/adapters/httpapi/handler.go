package httpapi

import (
	"net/http"
)

func (a *St) hRoot(w http.ResponseWriter, r *http.Request) {
	a.uRespondJSON(w, map[string]string{"hello": "world"})
}
