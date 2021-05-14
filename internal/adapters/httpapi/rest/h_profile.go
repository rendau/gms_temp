package rest

import (
	"net/http"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (a *St) hProfileSendPhoneValidatingCode(w http.ResponseWriter, r *http.Request) {
	reqObj := &(struct {
		Phone string `json:"phone"`
		ErrNE bool   `json:"err_ne"`
	}{})
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ProfileSendPhoneValidatingCode(a.uGetRequestContext(r), reqObj.Phone, reqObj.ErrNE)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

func (a *St) hProfileAuth(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.UsrAuthReqSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	usrId, token, err := a.ucs.ProfileAuth(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, struct {
		Id    int64  `json:"id"`
		Token string `json:"token"`
	}{usrId, token})
}

func (a *St) hProfileLogout(w http.ResponseWriter, r *http.Request) {
	err := a.ucs.ProfileLogout(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

func (a *St) hProfileGet(w http.ResponseWriter, r *http.Request) {
	profile, err := a.ucs.ProfileGet(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, profile)
}

func (a *St) hProfileUpdate(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ProfileUpdate(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

func (a *St) hProfileGetId(w http.ResponseWriter, r *http.Request) {
	usrId, err := a.ucs.ProfileGetId(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, map[string]int64{
		"id": usrId,
	})
}

func (a *St) hProfileDelete(w http.ResponseWriter, r *http.Request) {
	err := a.ucs.ProfileDelete(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
