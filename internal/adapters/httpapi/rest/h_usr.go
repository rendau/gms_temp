package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

// swagger:parameters hUsrGet hUsrUpdate hUsrDelete
type docUsrPathParIdSt struct {
	// in: path
	Id int64 `json:"id"`
}

// swagger:parameters hUsrCreate hUsrUpdate
type docUsrBodyParObjSt struct {
	// in: body
	Body entities.UsrCUSt
}

// swagger:route GET /usrs usrs hUsrList
// Allowed to: `UsrTypeAdmin`
// Security:
//   token:
// Responses:
//   200: usrListRep
//   400: errRep
func (a *St) hUsrList(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hUsrList
	type docReqSt struct {
		// in:query
		entities.PaginationParams

		// in:query
		TypeId int `json:"type_id"`

		// in:query
		OnlyCount int `json:"only_count"`

		// in:query
		Search int `json:"search"`
	}

	// swagger:response usrListRep
	type docRepSt struct {
		// in:body
		Body struct {
			DocPaginatedListRepSt
			Results []*entities.UsrListSt `json:"results"`
		}
	}

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

// swagger:route POST /usrs usrs hUsrCreate
// Allowed to: `UsrTypeAdmin`
// Security:
//   token:
// Responses:
//   200: usrCreateRep
//   400: errRep
func (a *St) hUsrCreate(w http.ResponseWriter, r *http.Request) {
	// swagger:response usrCreateRep
	type docRepSt struct {
		// in:body
		Body struct {
			Id int64 `json:"id"`
		}
	}

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

// swagger:route GET /usrs/{id} usrs hUsrGet
// Security:
//   token:
// Responses:
//   200: usrRep
//   400: errRep
func (a *St) hUsrGet(w http.ResponseWriter, r *http.Request) {
	// swagger:response usrRep
	type docRepSt struct {
		// in:body
		Body *entities.UsrSt
	}

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	// qPars := r.URL.Query()

	result, err := a.ucs.UsrGet(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, result)
}

// swagger:route PUT /usrs/{id} usrs hUsrUpdate
// Allowed to: `UsrTypeAdmin`
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
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

// swagger:route DELETE /usrs/{id} usrs hUsrDelete
// Allowed to: `UsrTypeAdmin`
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hUsrDelete(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	err := a.ucs.UsrDelete(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
