package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

func (a *St) hUsrList(w http.ResponseWriter, r *http.Request) {
	qPars := r.URL.Query()

	pars := &entities.UsrListParsSt{
		TypeId:    a.uQpParseInt(qPars, "type_id"),
		OnlyCount: a.uQpParseBoolV(qPars, "only_count"),
		Search:    a.uQpParseString(qPars, "search"),
	}

	if idPar := a.uQpParseInt64(qPars, "id"); idPar != nil {
		pars.Ids = util.NewSliceInt64(*idPar)
	}

	offset, limit, page := a.uExtractPaginationPars(qPars)
	pars.Offset = offset
	pars.Limit = limit

	paginated := pars.Limit > 0

	result, tCount, err := a.ucs.UsrList(a.uGetRequestContext(r), pars)
	if a.uHandleError(err, r, w) {
		return
	}

	if paginated {
		a.uRespondJSON(w, 0, &PaginatedListRepSt{
			Page:       page,
			PageSize:   limit,
			TotalCount: tCount,
			Results:    result,
		})
	} else {
		a.uRespondJSON(w, 0, result)
	}
}

func (a *St) hUsrCreate(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	newId, err := a.ucs.UsrCreate(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, map[string]int64{"id": newId})
}

func (a *St) hUsrGet(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	// qPars := r.URL.Query()

	result, err := a.ucs.UsrGet(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, result)
}

func (a *St) hUsrUpdate(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.UsrUpdate(a.uGetRequestContext(r), id, reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

func (a *St) hUsrDelete(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	err := a.ucs.UsrDelete(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
